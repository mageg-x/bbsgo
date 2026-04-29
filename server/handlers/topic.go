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

	// 获取当前用户ID（用于私密帖子过滤）
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 构建缓存 key（基于查询参数和用户ID，因为私密帖子对不同用户可见性不同）
	cacheKey := cache.BuildKey(cache.TopicListPrefix, fmt.Sprintf("%d:%d:%s:%d:%d:%d", forumID, tagID, sort, page, pageSize, userID))

	// 尝试从缓存获取
	data, err := cache.GetData(cacheKey, func() (interface{}, error) {
		return fetchTopics(forumID, tagID, sort, page, pageSize, userID)
	}, cache.TopicListExpire)

	if err != nil {
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, data)
}

// fetchTopics 获取话题列表核心逻辑
func fetchTopics(forumID, tagID int, sort string, page, pageSize int, currentUserID uint) (map[string]interface{}, error) {
	var topics []models.Topic
	var total int64

	// 构建查询条件：过滤活跃的私密帖子（仅作者和管理员可见）
	// 条件：非私密 或 已结束私密 或 是自己的帖子
	now := time.Now()
	privateCondition := "(is_private = ? OR private_ended = ? OR (is_private = ? AND (private_expire_at IS NULL OR private_expire_at < ?)))"
	if currentUserID > 0 {
		// 已登录用户可以看到自己的私密帖子
		privateCondition = fmt.Sprintf("(%s OR user_id = %d)", privateCondition, currentUserID)
	}

	// 构建查询
	query := database.DB.Model(&models.Topic{}).
		Where(privateCondition, false, true, true, now)
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
	dbQuery := database.DB.Preload("User").Preload("Forum").Preload("Tags").
		Where(privateCondition, false, true, true, now)
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

	// 检查私密权限
	if topic.IsPrivateActive() {
		userID, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok || (topic.UserID != userID) {
			// 检查是否是管理员
			if ok {
				var user models.User
				if err := database.DB.First(&user, userID).Error; err == nil && user.Role >= 1 {
					// 管理员可以查看
					goto allowAccess
				}
			}
			log.Printf("get topic: private topic access denied, id: %d", id)
			errors.Error(w, errors.CodeTopicNotFound, "话题不存在或无权访问")
			return
		}
	}
allowAccess:

	// 增加浏览数（异步更新）
	go func(topicID uint) {
		database.DB.Model(&models.Topic{}).Where("id = ?", topicID).
			UpdateColumn("view_count", database.DB.Raw("view_count + 1"))
	}(topic.ID)

	errors.Success(w, topic)
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
		Title           string   `json:"title"`            // 话题标题
		Content         string   `json:"content"`          // 话题内容
		ForumID         uint     `json:"forum_id"`         // 版块ID
		TagNames        []string `json:"tag_names"`        // 标签名称列表
		IsPrivate       bool     `json:"is_private"`       // 是否私密
		PrivateDuration int64    `json:"private_duration"` // 私密时长
		PrivateUnit     string   `json:"private_unit"`     // 时长单位：second, minute, hour, day
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

	// 创建话题
	topic := models.Topic{
		Title:           req.Title,
		Content:         req.Content,
		UserID:          userID,
		ForumID:         req.ForumID,
		AllowComment:    true,           // 默认允许评论
		Tags:            []models.Tag{}, // 初始化 Tags 为空切片，避免 nil 导致的空指针异常
		IsPrivate:       req.IsPrivate,
		PrivateDuration: req.PrivateDuration,
		PrivateUnit:     req.PrivateUnit,
	}
	
	// 如果设置了私密且有时长，计算过期时间
	if req.IsPrivate && req.PrivateDuration > 0 {
		topic.PrivateExpireAt = topic.CalculatePrivateExpireAt(req.PrivateDuration, req.PrivateUnit)
	}
	
	// 私密帖子默认禁止评论
	if req.IsPrivate {
		topic.AllowComment = false
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

// UnlockPrivateTopic 手动解禁私密话题处理器
// 作者可以提前结束私密状态
func UnlockPrivateTopic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("unlock private topic: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("unlock private topic: topic not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 验证权限：仅作者可以操作
	if topic.UserID != userID {
		log.Printf("unlock private topic: permission denied, topicID: %d, userID: %d", id, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 检查是否是私密状态
	if !topic.IsPrivate {
		errors.Error(w, errors.CodeInvalidParams, "该话题不是私密状态")
		return
	}

	if topic.PrivateEnded {
		errors.Error(w, errors.CodeInvalidParams, "该私密话题已过期或已解禁")
		return
	}

	// 更新私密状态
	updates := map[string]interface{}{
		"is_private":    false,
		"private_ended": true,
		"allow_comment": true, // 解禁后允许评论
	}

	if err := database.DB.Model(&topic).Updates(updates).Error; err != nil {
		log.Printf("unlock private topic: failed to update topic, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("unlock private topic: topic unlocked, id: %d", id)
	errors.Success(w, map[string]interface{}{
		"id":          topic.ID,
		"is_private":  false,
		"private_ended": true,
	})

	// 清除缓存
	cache.TopicCache.Invalidate(topic.ID)
	cache.HomePageCache.InvalidateTopics()
}

// ExtendPrivateTopic 延长私密时长处理器
// 作者可以延长私密话题的时长
func ExtendPrivateTopic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("extend private topic: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("extend private topic: topic not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 验证权限：仅作者可以操作
	if topic.UserID != userID {
		log.Printf("extend private topic: permission denied, topicID: %d, userID: %d", id, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 检查是否是私密状态
	if !topic.IsPrivate {
		errors.Error(w, errors.CodeInvalidParams, "该话题不是私密状态")
		return
	}

	if topic.PrivateEnded {
		errors.Error(w, errors.CodeInvalidParams, "该私密话题已过期或已解禁")
		return
	}

	// 解析请求体
	var req struct {
		Duration int64  `json:"duration"` // 延长的时长
		Unit     string `json:"unit"`     // 时长单位：second, minute, hour, day
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("extend private topic: failed to decode request body, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	if req.Duration <= 0 {
		errors.Error(w, errors.CodeInvalidParams, "延时时长必须大于0")
		return
	}

	// 计算新的过期时间
	var newExpireAt time.Time
	if topic.PrivateExpireAt != nil && time.Now().Before(*topic.PrivateExpireAt) {
		// 在原有过期时间基础上延长
		newExpireAt = *topic.PrivateExpireAt
	} else {
		// 如果已过期或没有设置过期时间，从现在开始计算
		newExpireAt = time.Now()
	}

	// 添加延时时长
	switch req.Unit {
	case "minute":
		newExpireAt = newExpireAt.Add(time.Minute * time.Duration(req.Duration))
	case "hour":
		newExpireAt = newExpireAt.Add(time.Hour * time.Duration(req.Duration))
	case "day":
		newExpireAt = newExpireAt.Add(time.Hour * 24 * time.Duration(req.Duration))
	default: // second
		newExpireAt = newExpireAt.Add(time.Second * time.Duration(req.Duration))
	}

	// 更新私密状态
	updates := map[string]interface{}{
		"private_expire_at": newExpireAt,
		"private_ended":     false,
	}

	if err := database.DB.Model(&topic).Updates(updates).Error; err != nil {
		log.Printf("extend private topic: failed to update topic, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("extend private topic: topic extended, id: %d, new expire at: %v", id, newExpireAt)
	errors.Success(w, map[string]interface{}{
		"id":                topic.ID,
		"is_private":        true,
		"private_expire_at": newExpireAt,
		"private_ended":     false,
	})

	// 清除缓存
	cache.TopicCache.Invalidate(topic.ID)
}
