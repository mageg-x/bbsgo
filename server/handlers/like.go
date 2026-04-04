package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var likeBadgeService = services.NewBadgeService()

// LikeRequest 点赞请求结构
type LikeRequest struct {
	TargetType string `json:"target_type"` // 目标类型：topic=话题, comment=评论
	TargetID   uint   `json:"target_id"`   // 目标ID
}

// CreateLike 创建点赞处理器
// 用户对话题或帖子进行点赞
func CreateLike(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create like: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 解析请求体
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create like: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 验证目标类型
	if req.TargetType != "topic" && req.TargetType != "comment" {
		log.Printf("create like: invalid target type, targetType: %s", req.TargetType)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 检查是否已经点赞
	var existingLike models.Like
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).First(&existingLike).Error; err == nil {
		log.Printf("create like: already liked, userID: %d, targetType: %s, targetID: %d", userID, req.TargetType, req.TargetID)
		errors.Error(w, errors.CodeAlreadyFollowed, "")
		return
	}

	// 创建点赞记录
	like := models.Like{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
	}

	if err := database.DB.Create(&like).Error; err != nil {
		log.Printf("create like: failed to create like, userID: %d, targetType: %s, targetID: %d, error: %v", userID, req.TargetType, req.TargetID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 更新目标对象的点赞数，并触发被点赞者的勋章检查
	var contentOwnerID uint
	if req.TargetType == "topic" {
		var topic models.Topic
		if err := database.DB.First(&topic, req.TargetID).Error; err != nil {
			log.Printf("create like: topic not found, topicID: %d, error: %v", req.TargetID, err)
		} else {
			contentOwnerID = topic.UserID
			if err := database.DB.Model(&topic).UpdateColumn("like_count", topic.LikeCount+1).Error; err != nil {
				log.Printf("create like: failed to increment topic like count, topicID: %d, error: %v", req.TargetID, err)
			}
		}
	} else if req.TargetType == "comment" {
		var comment models.Comment
		if err := database.DB.First(&comment, req.TargetID).Error; err != nil {
			log.Printf("create like: comment not found, commentID: %d, error: %v", req.TargetID, err)
		} else {
			contentOwnerID = comment.UserID
			if err := database.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount+1).Error; err != nil {
				log.Printf("create like: failed to increment comment like count, commentID: %d, error: %v", req.TargetID, err)
			}
		}
	}

	// 如果成功获取内容作者ID，触发其勋章检查（因为收到了点赞）
	if contentOwnerID > 0 {
		go likeBadgeService.CheckAndAwardBadges(contentOwnerID)
	}

	log.Printf("create like: like created successfully, userID: %d, targetType: %s, targetID: %d", userID, req.TargetType, req.TargetID)
	errors.Success(w, like)
}

// DeleteLike 取消点赞处理器
// 用户取消对话题或帖子的点赞
func DeleteLike(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete like: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 从 URL 查询参数获取
	targetType := r.URL.Query().Get("target_type")
	targetIDStr := r.URL.Query().Get("target_id")
	targetID, _ := strconv.Atoi(targetIDStr)

	if targetType == "" || targetID == 0 {
		log.Printf("delete like: invalid parameters, targetType: %s, targetID: %d", targetType, targetID)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 查询点赞记录
	var like models.Like
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).First(&like).Error; err != nil {
		log.Printf("delete like: like not found, userID: %d, targetType: %s, targetID: %d, error: %v", userID, targetType, targetID, err)
		errors.Error(w, errors.CodeTopicNotFound, "")
		return
	}

	// 删除点赞记录
	if err := database.DB.Unscoped().Delete(&like).Error; err != nil {
		log.Printf("delete like: failed to delete like, id: %d, error: %v", like.ID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 更新目标对象的点赞数
	if targetType == "topic" {
		var topic models.Topic
		if err := database.DB.First(&topic, targetID).Error; err != nil {
			log.Printf("delete like: topic not found, topicID: %d, error: %v", targetID, err)
		} else if topic.LikeCount > 0 {
			if err := database.DB.Model(&topic).UpdateColumn("like_count", topic.LikeCount-1).Error; err != nil {
				log.Printf("delete like: failed to decrement topic like count, topicID: %d, error: %v", targetID, err)
			}
		}
	} else if targetType == "comment" {
		var comment models.Comment
		if err := database.DB.First(&comment, targetID).Error; err != nil {
			log.Printf("delete like: comment not found, commentID: %d, error: %v", targetID, err)
		} else if comment.LikeCount > 0 {
			if err := database.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount-1).Error; err != nil {
				log.Printf("delete like: failed to decrement comment like count, commentID: %d, error: %v", targetID, err)
			}
		}
	}

	log.Printf("delete like: like deleted successfully, userID: %d, targetType: %s, targetID: %d", userID, targetType, targetID)
	errors.Success(w, nil)
}

// CheckLike 检查点赞状态处理器
// 支持单条和批量检查
func CheckLike(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.Success(w, map[string]interface{}{"liked": false})
		return
	}

	var req struct {
		TargetType string `json:"target_type"`
		TargetIDs  []uint `json:"target_ids"` // 批量检查时使用
		TargetID   uint   `json:"target_id"`  // 单条检查时使用
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.Success(w, map[string]interface{}{"liked": false})
		return
	}

	// 批量检查模式
	if len(req.TargetIDs) > 0 {
		var likes []models.Like
		if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id IN ?", userID, req.TargetType, req.TargetIDs).Find(&likes).Error; err != nil {
			errors.Success(w, map[string]interface{}{"liked_map": map[uint]bool{}})
			return
		}
		likedMap := make(map[uint]bool)
		for _, like := range likes {
			likedMap[like.TargetID] = true
		}
		errors.Success(w, map[string]interface{}{"liked_map": likedMap})
		return
	}

	// 单条检查模式
	var count int64
	if err := database.DB.Model(&models.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).
		Count(&count).Error; err != nil {
		errors.Success(w, map[string]interface{}{"liked": false})
		return
	}

	errors.Success(w, map[string]interface{}{"liked": count > 0})
}
