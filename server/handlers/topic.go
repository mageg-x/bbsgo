package handlers

import (
	"bbsgo/antispam"
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/services"
	"bbsgo/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var badgeService = services.NewBadgeService()

// GetTopics 获取话题列表处理器
// 支持按版块、标签筛选，支持分页
func GetTopics(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	forumID, _ := strconv.Atoi(r.URL.Query().Get("forum_id"))
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 构建缓存 key（基于查询参数）
	cacheKey := cache.BuildKey(cache.TopicListPrefix, fmt.Sprintf("%d:%d:%s:%d:%d", forumID, tagID, sort, page, pageSize))

	// 尝试从缓存获取
	data, err := cache.GetData(cacheKey, func() (interface{}, error) {
		return fetchTopics(forumID, tagID, sort, page, pageSize)
	}, cache.TopicListExpire)

	if err != nil {
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, data)
}

// fetchTopics 获取话题列表核心逻辑
func fetchTopics(forumID, tagID int, sort string, page, pageSize int) (map[string]interface{}, error) {
	var topics []models.Topic
	var total int64

	// 构建查询
	query := database.DB.Model(&models.Topic{})
	if forumID > 0 {
		query = query.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("get topics: failed to count topics, error: %v", err)
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	dbQuery := database.DB.Preload("User").Preload("Forum").Preload("Tags")
	if forumID > 0 {
		dbQuery = dbQuery.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		dbQuery = dbQuery.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	// 排序方式
	switch sort {
	case "hot":
		dbQuery = dbQuery.Order("is_pinned DESC, (like_count + reply_count * 2) DESC, created_at DESC")
	case "reply":
		dbQuery = dbQuery.Order("is_pinned DESC, last_reply_at DESC NULLS LAST, created_at DESC")
	default:
		dbQuery = dbQuery.Order("is_pinned DESC, created_at DESC")
	}

	// 执行查询
	if err := dbQuery.Offset(offset).Limit(pageSize).Find(&topics).Error; err != nil {
		log.Printf("get topics: failed to query topics, error: %v", err)
		return nil, err
	}

	// 给每个 topic 添加 has_poll 字段
	var topicIDs []uint
	for _, t := range topics {
		topicIDs = append(topicIDs, t.ID)
	}

	// 批量查询投票存在的 topic IDs
	hasPollMap := make(map[uint]bool)
	if len(topicIDs) > 0 {
		var polls []models.Poll
		database.DB.Model(&models.Poll{}).Where("topic_id IN ?", topicIDs).Select("topic_id").Find(&polls)
		for _, p := range polls {
			hasPollMap[p.TopicID] = true
		}
	}

	// 收集所有作者用户ID
	userIDs := make(map[uint]bool)
	for _, t := range topics {
		userIDs[t.UserID] = true
	}

	// 批量查询用户的勋章
	userBadgesMap := make(map[uint][]models.UserBadge)
	if len(userIDs) > 0 {
		ids := make([]uint, 0, len(userIDs))
		for id := range userIDs {
			ids = append(ids, id)
		}
		var userBadges []models.UserBadge
		if err := database.DB.Where("user_id IN ? AND is_revoked = ?", ids, false).
			Preload("Badge").
			Find(&userBadges).Error; err != nil {
			log.Printf("get topics: failed to query user badges, error: %v", err)
		}
		for _, ub := range userBadges {
			userBadgesMap[ub.UserID] = append(userBadgesMap[ub.UserID], ub)
		}
	}

	// 构建返回数据，添加 has_poll 和 author_badges
	type TopicWithPoll struct {
		models.Topic
		HasPoll      bool               `json:"has_poll"`
		AuthorBadges []models.UserBadge `json:"author_badges"`
	}
	var response []TopicWithPoll
	for _, t := range topics {
		response = append(response, TopicWithPoll{
			Topic:        t,
			HasPoll:      hasPollMap[t.ID],
			AuthorBadges: userBadgesMap[t.UserID],
		})
	}

	return map[string]interface{}{
		"list":      response,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

// GetTopic 获取单个话题详情处理器
func GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 尝试从缓存获取
	cacheKey := cache.TopicCache.BuildKey(id)
	data, err := cache.GetData(cacheKey, func() (interface{}, error) {
		var topic models.Topic
		if err := database.DB.Preload("User").Preload("Forum").First(&topic, id).Error; err != nil {
			return nil, err
		}
		return topic, nil
	}, cache.TopicExpire)

	if err != nil {
		log.Printf("get topic: topic not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	topic := data.(models.Topic)

	// 增加浏览数（异步更新）
	go func(topicID uint) {
		database.DB.Model(&models.Topic{}).Where("id = ?", topicID).
			UpdateColumn("view_count", database.DB.Raw("view_count + 1"))
	}(topic.ID)

	// 处理内容解锁逻辑
	response := prepareTopicResponse(&topic, r)

	errors.Success(w, response)
}

// prepareTopicResponse 准备话题响应，处理内容解锁逻辑
func prepareTopicResponse(topic *models.Topic, r *http.Request) map[string]interface{} {
	// 构建基础响应
	response := map[string]interface{}{
		"id":                    topic.ID,
		"title":                 topic.Title,
		"content":               topic.Content,
		"user_id":               topic.UserID,
		"user":                  topic.User,
		"forum_id":              topic.ForumID,
		"forum":                 topic.Forum,
		"is_pinned":             topic.IsPinned,
		"is_user_pinned":        topic.IsUserPinned,
		"is_locked":             topic.IsLocked,
		"is_essence":            topic.IsEssence,
		"is_hidden":             topic.IsHidden,
		"hot_score":             topic.HotScore,
		"like_count":            topic.LikeCount,
		"view_count":            topic.ViewCount,
		"reply_count":           topic.ReplyCount,
		"last_reply_at":         topic.LastReplyAt,
		"allow_comment":         topic.AllowComment,
		"created_at":            topic.CreatedAt,
		"updated_at":            topic.UpdatedAt,
		"tags":                  topic.Tags,
		"is_unlock_enabled":     topic.IsUnlockEnabled,
		"unlock_type":           topic.UnlockType,
		"unlock_like_count":     topic.UnlockLikeCount,
		"unlock_comment_count":  topic.UnlockCommentCount,
		"preview_length":        topic.PreviewLength,
		"is_unlocked":           true, // 默认已解锁
		"unlock_progress":       map[string]interface{}{},
	}

	// 如果没有启用解锁功能，直接返回完整内容
	if !topic.IsUnlockEnabled {
		return response
	}

	// 获取当前用户ID
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	
	// 检查是否是作者本人（作者总是可以看到完整内容）
	if ok && userID == topic.UserID {
		return response
	}

	// 检查用户是否已经解锁
	var isUnlocked bool
	if ok {
		var unlockRecord models.UnlockRecord
		if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, topic.ID).First(&unlockRecord).Error; err == nil {
			isUnlocked = true
		}
	}

	// 计算解锁进度
	unlockProgress := calculateUnlockProgress(topic)
	response["unlock_progress"] = unlockProgress

	// 如果未解锁，只返回部分内容
	if !isUnlocked {
		response["is_unlocked"] = false
		// 截取内容预览
		content := topic.Content
		if len(content) > topic.PreviewLength {
			response["content"] = content[:topic.PreviewLength] + "..."
		} else {
			response["content"] = content
		}
	}

	return response
}

// calculateUnlockProgress 计算解锁进度
func calculateUnlockProgress(topic *models.Topic) map[string]interface{} {
	progress := map[string]interface{}{
		"like_progress":    0,
		"comment_progress": 0,
		"like_needed":      topic.UnlockLikeCount,
		"comment_needed":   topic.UnlockCommentCount,
		"like_current":     topic.LikeCount,
		"comment_current":  topic.ReplyCount,
		"unlock_type":      topic.UnlockType,
	}

	// 计算点赞进度
	if topic.UnlockType == "like" || topic.UnlockType == "both" {
		if topic.UnlockLikeCount > 0 {
			likeProgress := float64(topic.LikeCount) / float64(topic.UnlockLikeCount)
			if likeProgress > 1 {
				likeProgress = 1
			}
			progress["like_progress"] = likeProgress
		}
	}

	// 计算评论进度
	if topic.UnlockType == "comment" || topic.UnlockType == "both" {
		if topic.UnlockCommentCount > 0 {
			commentProgress := float64(topic.ReplyCount) / float64(topic.UnlockCommentCount)
			if commentProgress > 1 {
				commentProgress = 1
			}
			progress["comment_progress"] = commentProgress
		}
	}

	return progress
}

// CreateTopic 创建话题处理器
// 需要用户登录
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	// 添加 recover 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("CreateTopic PANIC: %v", err)
			errors.Error(w, errors.CodeServerInternal, "")
		}
	}()

	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create topic: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 检查是否允许发帖
	if !utils.GetConfigBool("allow_post", true) {
		log.Printf("create topic: create topic disabled")
		errors.Error(w, errors.CodePostDisabled, "")
		return
	}

	// 解析请求体
	var req struct {
		Title               string   `json:"title"`                // 话题标题
		Content             string   `json:"content"`              // 话题内容
		ForumID             uint     `json:"forum_id"`             // 版块ID
		TagNames            []string `json:"tag_names"`            // 标签名称列表
		IsUnlockEnabled     bool     `json:"is_unlock_enabled"`    // 是否启用解锁功能
		UnlockType          string   `json:"unlock_type"`          // 解锁类型：like=点赞解锁, comment=评论解锁, both=两者都需要
		UnlockLikeCount     int      `json:"unlock_like_count"`    // 解锁所需点赞数门槛
		UnlockCommentCount  int      `json:"unlock_comment_count"` // 解锁所需回复数门槛
		PreviewLength       int      `json:"preview_length"`       // 内容预览长度（字符数）
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create topic: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 验证必填字段
	if req.Title == "" || req.Content == "" || req.ForumID == 0 {
		log.Printf("create topic: incomplete topic info, title: %s, forumID: %d", req.Title, req.ForumID)
		errors.Error(w, errors.CodeIncompleteInfo, "")
		return
	}

	// XSS 安全检查
	if middleware.ContainsXSS(req.Title) || middleware.ContainsXSS(req.Content) {
		log.Printf("create topic: potential XSS attack detected, title: %s, userID: %d", req.Title, userID)
		errors.Error(w, errors.CodeSensitiveContent, "内容包含不安全字符")
		return
	}

	// 验证标签数量
	if len(req.TagNames) > 3 {
		log.Printf("create topic: too many tags, count: %d", len(req.TagNames))
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 防刷检查
	antispamMiddleware := antispam.GetAntiSpamMiddleware()
	checkResult := antispamMiddleware.CheckTopicCreate(userID, req.Content)
	if !checkResult.Allowed {
		log.Printf("create topic: antispam check failed, userID: %d, reason: %s", userID, checkResult.Reason)
		// 根据具体原因返回对应错误码
		reason := checkResult.Reason
		switch {
		case strings.Contains(reason, "禁言"):
			errors.Error(w, errors.CodeUserBanned, reason)
		case strings.Contains(reason, "操作过快"):
			errors.Error(w, errors.CodeOperationTooFast, reason)
		case strings.Contains(reason, "已达上限"):
			errors.Error(w, errors.CodeDailyLimitExceeded, reason)
		case strings.Contains(reason, "太短"):
			errors.Error(w, errors.CodeContentTooShort, reason)
		case strings.Contains(reason, "广告"):
			errors.Error(w, errors.CodeSensitiveContent, reason)
		case strings.Contains(reason, "敏感词"):
			errors.Error(w, errors.CodeSensitiveContent, reason)
		case strings.Contains(reason, "重复字符"):
			errors.Error(w, errors.CodeRepeatingChars, reason)
		case strings.Contains(reason, "无实质信息"):
			errors.Error(w, errors.CodeNoSubstantiveContent, reason)
		case strings.Contains(reason, "符号或表情"):
			errors.Error(w, errors.CodeSymbolsOrEmojiOnly, reason)
		case strings.Contains(reason, "过多外部链接"):
			errors.Error(w, errors.CodeTooManyLinks, reason)
		case strings.Contains(reason, "重复"):
			errors.Error(w, errors.CodeSensitiveContent, reason)
		default:
			errors.Error(w, errors.CodeSensitiveContent, reason)
		}
		return
	}

	// 验证解锁条件参数
	if req.IsUnlockEnabled {
		// 验证解锁类型
		if req.UnlockType != "like" && req.UnlockType != "comment" && req.UnlockType != "both" {
			log.Printf("create topic: invalid unlock type, unlockType: %s", req.UnlockType)
			errors.Error(w, errors.CodeInvalidParams, "无效的解锁类型")
			return
		}
		
		// 验证解锁门槛
		if req.UnlockType == "like" || req.UnlockType == "both" {
			if req.UnlockLikeCount <= 0 {
				log.Printf("create topic: invalid unlock like count, unlockLikeCount: %d", req.UnlockLikeCount)
				errors.Error(w, errors.CodeInvalidParams, "点赞数门槛必须大于0")
				return
			}
		}
		
		if req.UnlockType == "comment" || req.UnlockType == "both" {
			if req.UnlockCommentCount <= 0 {
				log.Printf("create topic: invalid unlock comment count, unlockCommentCount: %d", req.UnlockCommentCount)
				errors.Error(w, errors.CodeInvalidParams, "回复数门槛必须大于0")
				return
			}
		}
		
		// 验证预览长度
		if req.PreviewLength <= 0 {
			req.PreviewLength = 200 // 默认预览长度
		}
	}

	// 创建话题
	topic := models.Topic{
		Title:               req.Title,
		Content:             req.Content,
		UserID:              userID,
		ForumID:             req.ForumID,
		AllowComment:        true,                  // 默认允许评论
		Tags:                []models.Tag{},        // 初始化 Tags 为空切片，避免 nil 导致的空指针异常
		IsUnlockEnabled:     req.IsUnlockEnabled,   // 是否启用解锁功能
		UnlockType:          req.UnlockType,        // 解锁类型
		UnlockLikeCount:     req.UnlockLikeCount,   // 点赞数门槛
		UnlockCommentCount:  req.UnlockCommentCount,// 回复数门槛
		PreviewLength:       req.PreviewLength,     // 预览长度
	}

	if err := database.DB.Create(&topic).Error; err != nil {
		log.Printf("create topic: failed to create topic, userID: %d, forumID: %d, error: %v", userID, req.ForumID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 记录用户操作（用于防刷统计）
	antispamMiddleware.RecordUserOperation(userID, "topic", topic.ID, req.Content)

	// 处理标签关联
	if len(req.TagNames) > 0 {
		var tags []models.Tag
		for _, name := range req.TagNames {
			if name == "" {
				continue
			}
			tag, err := GetOrCreateTagByName(name)
			if err != nil {
				log.Printf("create topic: failed to get or create tag, name: %s, error: %v", name, err)
				continue
			}
			if tag == nil {
				log.Printf("create topic: tag is nil, name: %s", name)
				continue
			}
			if tag.IsBanned {
				log.Printf("create topic: tag is banned, name: %s", name)
				continue
			}
			// 只保留 ID，避免 GORM 尝试插入完整对象
			tags = append(tags, models.Tag{ID: tag.ID})
			IncrementTagUsage(tag.ID)
		}
		if len(tags) > 0 {
			if err := database.DB.Model(&topic).Association("Tags").Replace(tags); err != nil {
				log.Printf("create topic: failed to associate tags, topicID: %d, error: %v", topic.ID, err)
			}
		}
	}

	// 给发帖用户增加积分
	var user models.User
	if err := database.DB.First(&user, userID).Error; err == nil {
		log.Printf("create topic: [DEBUG] before GetConfigInt, database.DB=%v", database.DB != nil)
		creditAmount := utils.GetConfigInt("credit_topic", 20)
		log.Printf("create topic: [DEBUG] after GetConfigInt, creditAmount=%d", creditAmount)
		user.Credits += creditAmount
		if err := database.DB.Save(&user).Error; err != nil {
			log.Printf("create topic: failed to add credits, userID: %d, error: %v", userID, err)
		}
	}

	// 重新加载话题关联数据
	log.Printf("create topic: [DEBUG] before Preload, topic.ID=%d", topic.ID)
	if err := database.DB.Preload("User").Preload("Forum").Preload("Tags").First(&topic, topic.ID).Error; err != nil {
		log.Printf("create topic: failed to reload topic, id: %d, error: %v", topic.ID, err)
	}
	log.Printf("create topic: [DEBUG] after Preload, topic.ID=%d", topic.ID)

	// 确保 Tags 不为 nil
	if topic.Tags == nil {
		topic.Tags = []models.Tag{}
	}

	log.Printf("create topic: [DEBUG] before Success, w=%v, topic=%+v", w != nil, topic)
	errors.Success(w, topic)
	log.Printf("create topic: [DEBUG] after Success")

	// 清除首页缓存
	cache.HomePageCache.InvalidateTopics()

	// 检查并授予勋章
	go badgeService.CheckAndAwardBadges(userID)
}

// UpdateTopic 更新话题处理器
// 仅话题作者可以更新
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update topic: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("update topic: topic not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 验证权限：仅作者可以更新
	if topic.UserID != userID {
		log.Printf("update topic: permission denied, topicID: %d, userID: %d", id, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update topic: failed to decode request body, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")
	delete(updates, "is_pinned")
	delete(updates, "is_user_pinned")
	delete(updates, "is_locked")
	delete(updates, "is_essence")
	delete(updates, "like_count")
	delete(updates, "view_count")
	delete(updates, "reply_count")

	// 执行更新
	if err := database.DB.Model(&topic).Updates(updates).Error; err != nil {
		log.Printf("update topic: failed to update topic, id: %d, userID: %d, error: %v", id, userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 重新加载数据
	if err := database.DB.Preload("User").Preload("Forum").First(&topic, id).Error; err != nil {
		log.Printf("update topic: failed to reload topic, id: %d, error: %v", id, err)
	}

	log.Printf("update topic: topic updated successfully, id: %d, userID: %d", id, userID)
	errors.Success(w, topic)
}

// DeleteTopic 删除话题处理器
// 作者或管理员可以删除
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete topic: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("delete topic: topic not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 查询用户信息用于权限判断
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("delete topic: user not found, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	// 验证权限：作者或管理员(role>=1)可以删除
	if topic.UserID != userID && user.Role < 1 {
		log.Printf("delete topic: permission denied, topicID: %d, userID: %d, topicUserID: %d", id, userID, topic.UserID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 物理删除话题（不保留软删除）
	if err := database.DB.Unscoped().Delete(&topic).Error; err != nil {
		log.Printf("delete topic: failed to delete topic, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("delete topic: topic deleted successfully, id: %d, userID: %d", id, userID)
	errors.Success(w, nil)

	// 清除首页缓存
	cache.HomePageCache.InvalidateTopics()
}

// AdminPinTopic 管理员置顶/取消置顶话题处理器
// 影响首页帖子列表排序
func AdminPinTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 解析请求体
	var req struct {
		Pinned bool `json:"pinned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("admin pin topic: failed to decode request body, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("admin pin topic: topic not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 更新置顶状态
	if err := database.DB.Model(&topic).UpdateColumn("is_pinned", req.Pinned).Error; err != nil {
		log.Printf("admin pin topic: failed to update pin status, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("admin pin topic: topic pin updated, id: %d, pinned: %v", id, req.Pinned)
	errors.Success(w, map[string]interface{}{
		"id":        topic.ID,
		"is_pinned": req.Pinned,
	})

	// 清除首页缓存
	cache.HomePageCache.InvalidateTopics()
}

// UserPinTopic 作者置顶/取消置顶话题处理器
// 影响个人主页帖子排序
func UserPinTopic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("user pin topic: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("user pin topic: topic not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 验证权限：仅作者可以操作
	if topic.UserID != userID {
		log.Printf("user pin topic: permission denied, topicID: %d, userID: %d", id, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 解析请求体
	var req struct {
		Pinned bool `json:"pinned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("user pin topic: failed to decode request body, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 更新作者置顶状态
	if err := database.DB.Model(&topic).UpdateColumn("is_user_pinned", req.Pinned).Error; err != nil {
		log.Printf("user pin topic: failed to update user pin status, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("user pin topic: topic user pin updated, id: %d, is_user_pinned: %v", id, req.Pinned)
	errors.Success(w, map[string]interface{}{
		"id":             topic.ID,
		"is_user_pinned": req.Pinned,
	})
}
