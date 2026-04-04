package handlers

import (
	"bbsgo/antispam"
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
	"time"

	"github.com/gorilla/mux"
)

var commentBadgeService = services.NewBadgeService()

// GetComments 获取话题的评论列表处理器
// 支持分页，返回话题下的一级评论
func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])

	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var comments []models.Comment
	var total int64

	offset := (page - 1) * pageSize

	// 统计评论数量
	if err := database.DB.Model(&models.Comment{}).Where("topic_id = ?", topicID).Count(&total).Error; err != nil {
		log.Printf("get comments: failed to count comments, topicID: %d, error: %v", topicID, err)
	}

	// 查询所有评论，按置顶和创建时间排序
	if err := database.DB.Where("topic_id = ?", topicID).
		Preload("User").
		Preload("ReplyTo").
		Preload("ReplyTo.User").
		Order("is_pinned DESC, created_at ASC").
		Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		log.Printf("get comments: failed to query comments, topicID: %d, error: %v", topicID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 收集所有评论者的用户ID
	userIDs := make(map[uint]bool)
	for _, comment := range comments {
		userIDs[comment.UserID] = true
	}

	// 批量查询用户的勋章
	userBadgesMap := make(map[uint][]models.UserBadge)
	if len(userIDs) > 0 {
		var userBadges []models.UserBadge
		ids := make([]uint, 0, len(userIDs))
		for id := range userIDs {
			ids = append(ids, id)
		}
		if err := database.DB.Where("user_id IN ? AND is_revoked = ?", ids, false).
			Preload("Badge").
			Find(&userBadges).Error; err != nil {
			log.Printf("get comments: failed to query user badges, error: %v", err)
		}
		for _, ub := range userBadges {
			userBadgesMap[ub.UserID] = append(userBadgesMap[ub.UserID], ub)
		}
	}

	// 构建响应结构，包含用户勋章和被回复用户信息
	type CommentWithUserBadges struct {
		models.Comment
		User struct {
			ID       uint               `json:"id"`
			Username string             `json:"username"`
			Nickname string             `json:"nickname"`
			Avatar   string             `json:"avatar"`
			Badges   []models.UserBadge `json:"badges"`
		} `json:"user"`
		ReplyUser *struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Nickname string `json:"nickname"`
		} `json:"reply_user,omitempty"`
	}

	response := make([]CommentWithUserBadges, len(comments))
	for i, comment := range comments {
		response[i] = CommentWithUserBadges{Comment: comment}
		response[i].User.ID = comment.User.ID
		response[i].User.Username = comment.User.Username
		response[i].User.Nickname = comment.User.Nickname
		response[i].User.Avatar = comment.User.Avatar
		response[i].User.Badges = userBadgesMap[comment.UserID]
		// 如果是回复评论，添加被回复用户信息
		if comment.ReplyTo != nil && comment.ReplyTo.User.ID != 0 {
			response[i].ReplyUser = &struct {
				ID       uint   `json:"id"`
				Username string `json:"username"`
				Nickname string `json:"nickname"`
			}{
				ID:       comment.ReplyTo.User.ID,
				Username: comment.ReplyTo.User.Username,
				Nickname: comment.ReplyTo.User.Nickname,
			}
		}
	}

	errors.Success(w, map[string]interface{}{
		"list":      response,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateComment 创建评论处理器
// 在指定话题下创建评论
func CreateComment(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create comment: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 检查是否允许评论
	if !utils.GetConfigBool("allow_comment", true) {
		log.Printf("create comment: comment disabled")
		errors.Error(w, errors.CodeCommentDisabled, "")
		return
	}

	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])

	// 解析请求体
	var req struct {
		Content   string `json:"content"`     // 评论内容
		ReplyToID *uint  `json:"reply_to_id"` // 回复给哪个评论的ID
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create comment: failed to decode request body, topicID: %d, error: %v", topicID, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 验证内容
	if req.Content == "" {
		log.Printf("create comment: content is empty, topicID: %d", topicID)
		errors.Error(w, errors.CodeIncompleteInfo, "")
		return
	}

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		log.Printf("create comment: topic not found, topicID: %d, error: %v", topicID, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 检查话题是否允许评论
	if topic.IsLocked || !topic.AllowComment {
		log.Printf("create comment: topic is locked or not allowing comments, topicID: %d", topicID)
		errors.Error(w, errors.CodeCommentDisabled, "")
		return
	}

	// 防刷检查
	antispamMiddleware := antispam.GetAntiSpamMiddleware()
	checkResult := antispamMiddleware.CheckCommentCreate(userID, req.Content)
	if !checkResult.Allowed {
		log.Printf("create comment: antispam check failed, userID: %d, reason: %s", userID, checkResult.Reason)
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

	// 创建评论
	comment := models.Comment{
		TopicID:   uint(topicID),
		UserID:    userID,
		Content:   req.Content,
		ReplyToID: req.ReplyToID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		log.Printf("create comment: failed to create comment, topicID: %d, userID: %d, error: %v", topicID, userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 记录用户操作（用于防刷统计）
	antispamMiddleware.RecordUserOperation(userID, "comment", comment.ID, req.Content)

	// 更新话题的评论数和最后回复时间
	now := time.Now()
	if err := database.DB.Model(&topic).Updates(map[string]interface{}{
		"reply_count":   topic.ReplyCount + 1,
		"last_reply_at": now,
	}).Error; err != nil {
		log.Printf("create comment: failed to update topic reply count, topicID: %d, error: %v", topicID, err)
	}

	// 给评论用户增加积分
	var user models.User
	if err := database.DB.First(&user, userID).Error; err == nil {
		creditAmount := utils.GetConfigInt("credit_post", 5)
		user.Credits += creditAmount
		if err := database.DB.Save(&user).Error; err != nil {
			log.Printf("create comment: failed to add credits, userID: %d, error: %v", userID, err)
		}
	}

	// 重新加载评论关联数据
	if err := database.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
		log.Printf("create comment: failed to reload comment, id: %d, error: %v", comment.ID, err)
	}

	log.Printf("create comment: comment created successfully, id: %d, topicID: %d, userID: %d", comment.ID, topicID, userID)
	errors.Success(w, comment)

	// 检查并授予勋章
	go commentBadgeService.CheckAndAwardBadges(userID)
}

// UpdateComment 更新评论处理器
// 仅评论作者可以更新
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update comment: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询评论
	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		log.Printf("update comment: comment not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeCommentNotFound, "")
		return
	}

	// 验证权限：仅作者可以更新
	if comment.UserID != userID {
		log.Printf("update comment: permission denied, commentID: %d, userID: %d", id, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update comment: failed to decode request body, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "topic_id")
	delete(updates, "reply_id")
	delete(updates, "created_at")
	delete(updates, "like_count")
	delete(updates, "is_pinned")

	// 执行更新
	if err := database.DB.Model(&comment).Updates(updates).Error; err != nil {
		log.Printf("update comment: failed to update comment, id: %d, userID: %d, error: %v", id, userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 重新加载数据
	if err := database.DB.Preload("User").First(&comment, id).Error; err != nil {
		log.Printf("update comment: failed to reload comment, id: %d, error: %v", id, err)
	}

	log.Printf("update comment: comment updated successfully, id: %d, userID: %d", id, userID)
	errors.Success(w, comment)
}

// DeleteComment 删除评论处理器
// 作者或管理员可以删除
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete comment: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询评论
	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		log.Printf("delete comment: comment not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeCommentNotFound, "")
		return
	}

	// 查询用户信息用于权限判断
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("delete comment: user not found, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	// 验证权限：作者或管理员(role>=1)可以删除
	if comment.UserID != userID && user.Role < 1 {
		log.Printf("delete comment: permission denied, commentID: %d, userID: %d", id, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 查询关联话题用于更新评论数
	var topic models.Topic
	if err := database.DB.First(&topic, comment.TopicID).Error; err != nil {
		log.Printf("delete comment: topic not found, topicID: %d, error: %v", comment.TopicID, err)
	}

	// 物理删除评论
	if err := database.DB.Unscoped().Delete(&comment).Error; err != nil {
		log.Printf("delete comment: failed to delete comment, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 更新话题评论数
	if topic.ID != 0 && topic.ReplyCount > 0 {
		if err := database.DB.Model(&topic).UpdateColumn("reply_count", topic.ReplyCount-1).Error; err != nil {
			log.Printf("delete comment: failed to update topic reply count, topicID: %d, error: %v", topic.ID, err)
		}
	}

	log.Printf("delete comment: comment deleted successfully, id: %d, userID: %d", id, userID)
	errors.Success(w, nil)
}

// PinComment 置顶/取消置顶评论处理器
// 仅帖子作者可以操作
func PinComment(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("pin comment: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["topic_id"])
	commentID, _ := strconv.Atoi(vars["comment_id"])

	// 查询评论
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		log.Printf("pin comment: comment not found, commentID: %d, error: %v", commentID, err)
		errors.Error(w, errors.CodeCommentNotFound, "")
		return
	}

	// 验证评论属于指定话题
	if comment.TopicID != uint(topicID) {
		log.Printf("pin comment: comment does not belong to topic, commentID: %d, topicID: %d", commentID, topicID)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 查询话题，验证是否为帖子作者
	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		log.Printf("pin comment: topic not found, topicID: %d, error: %v", topicID, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 仅帖子作者可以置顶评论
	if topic.UserID != userID {
		log.Printf("pin comment: permission denied, topicID: %d, userID: %d", topicID, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	// 解析请求体
	var req struct {
		Pinned bool `json:"pinned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("pin comment: failed to decode request body, commentID: %d, error: %v", commentID, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 更新置顶状态
	if err := database.DB.Model(&comment).UpdateColumn("is_pinned", req.Pinned).Error; err != nil {
		log.Printf("pin comment: failed to update pin status, commentID: %d, error: %v", commentID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("pin comment: comment pin updated, commentID: %d, pinned: %v", commentID, req.Pinned)
	errors.Success(w, map[string]interface{}{
		"id":        comment.ID,
		"is_pinned": req.Pinned,
	})
}

// BestComment 标记/取消最佳评论处理器
// 仅帖子作者可以操作
func BestComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("best comment: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["topic_id"])
	commentID, _ := strconv.Atoi(vars["comment_id"])

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		log.Printf("best comment: comment not found, commentID: %d, error: %v", commentID, err)
		errors.Error(w, errors.CodeCommentNotFound, "")
		return
	}

	if comment.TopicID != uint(topicID) {
		log.Printf("best comment: comment does not belong to topic, commentID: %d, topicID: %d", commentID, topicID)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		log.Printf("best comment: topic not found, topicID: %d, error: %v", topicID, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	if topic.UserID != userID {
		log.Printf("best comment: permission denied, topicID: %d, userID: %d", topicID, userID)
		errors.Error(w, errors.CodeNoPermission, "")
		return
	}

	var req struct {
		Best bool `json:"best"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("best comment: failed to decode request body, commentID: %d, error: %v", commentID, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	if req.Best {
		if comment.UserID == userID {
			log.Printf("best comment: cannot mark own comment as best, userID: %d", userID)
			errors.Error(w, errors.CodeNoPermission, "")
			return
		}

		var existingBest models.Comment
		if err := database.DB.Where("topic_id = ? AND is_best = ? AND id != ?", topicID, true, commentID).First(&existingBest).Error; err == nil {
			if err := database.DB.Model(&existingBest).UpdateColumn("is_best", false).Error; err != nil {
				log.Printf("best comment: failed to clear previous best comment, commentID: %d, error: %v", existingBest.ID, err)
				errors.Error(w, errors.CodeServerInternal, "")
				return
			}
			log.Printf("best comment: cleared previous best comment, commentID: %d", existingBest.ID)
		}
	}

	if err := database.DB.Model(&comment).UpdateColumn("is_best", req.Best).Error; err != nil {
		log.Printf("best comment: failed to update best status, commentID: %d, error: %v", commentID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	if req.Best {
		go commentBadgeService.CheckAndAwardBadges(comment.UserID)
		CreateNotificationWithRelated(
			comment.UserID,
			"best_comment",
			"notifications.comment_best",
			fmt.Sprintf("/topic/%d", topicID),
			comment.ID,
			"comment",
		)
		// 最佳评论积分奖励
		var bestUser models.User
		if err := database.DB.First(&bestUser, comment.UserID).Error; err == nil {
			creditAmount := utils.GetConfigInt("credit_best_comment", 1)
			bestUser.Credits += creditAmount
			if err := database.DB.Save(&bestUser).Error; err != nil {
				log.Printf("best comment: failed to add credits, userID: %d, error: %v", comment.UserID, err)
			} else {
				log.Printf("best comment: awarded %d credits to userID: %d", creditAmount, comment.UserID)
			}
		}
	}

	log.Printf("best comment: comment best updated, commentID: %d, is_best: %v", commentID, req.Best)
	errors.Success(w, map[string]interface{}{
		"id":      comment.ID,
		"is_best": req.Best,
	})
}
