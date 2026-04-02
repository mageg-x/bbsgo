package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
)

// LikeRequest 点赞请求结构
type LikeRequest struct {
	TargetType string `json:"target_type"` // 目标类型：topic=话题, post=帖子
	TargetID   uint   `json:"target_id"`   // 目标ID
}

// CreateLike 创建点赞处理器
// 用户对话题或帖子进行点赞
func CreateLike(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create like: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 解析请求体
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create like: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证目标类型
	if req.TargetType != "topic" && req.TargetType != "post" {
		log.Printf("create like: invalid target type, targetType: %s", req.TargetType)
		utils.Error(w, 400, "无效的目标类型")
		return
	}

	// 检查是否已经点赞
	var existingLike models.Like
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).First(&existingLike).Error; err == nil {
		log.Printf("create like: already liked, userID: %d, targetType: %s, targetID: %d", userID, req.TargetType, req.TargetID)
		utils.Error(w, 400, "已经点赞过了")
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
		utils.Error(w, 500, "点赞失败")
		return
	}

	// 更新目标对象的点赞数
	if req.TargetType == "topic" {
		var topic models.Topic
		if err := database.DB.First(&topic, req.TargetID).Error; err != nil {
			log.Printf("create like: topic not found, topicID: %d, error: %v", req.TargetID, err)
		} else {
			if err := database.DB.Model(&topic).UpdateColumn("like_count", topic.LikeCount+1).Error; err != nil {
				log.Printf("create like: failed to increment topic like count, topicID: %d, error: %v", req.TargetID, err)
			}
		}
	} else if req.TargetType == "post" {
		var post models.Post
		if err := database.DB.First(&post, req.TargetID).Error; err != nil {
			log.Printf("create like: post not found, postID: %d, error: %v", req.TargetID, err)
		} else {
			if err := database.DB.Model(&post).UpdateColumn("like_count", post.LikeCount+1).Error; err != nil {
				log.Printf("create like: failed to increment post like count, postID: %d, error: %v", req.TargetID, err)
			}
		}
	}

	log.Printf("create like: like created successfully, userID: %d, targetType: %s, targetID: %d", userID, req.TargetType, req.TargetID)
	utils.Success(w, like)
}

// DeleteLike 取消点赞处理器
// 用户取消对话题或帖子的点赞
func DeleteLike(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete like: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 解析请求体
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("delete like: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 查询点赞记录
	var like models.Like
	if err := database.DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).First(&like).Error; err != nil {
		log.Printf("delete like: like not found, userID: %d, targetType: %s, targetID: %d, error: %v", userID, req.TargetType, req.TargetID, err)
		utils.Error(w, 404, "点赞记录不存在")
		return
	}

	// 删除点赞记录
	if err := database.DB.Unscoped().Delete(&like).Error; err != nil {
		log.Printf("delete like: failed to delete like, id: %d, error: %v", like.ID, err)
		utils.Error(w, 500, "取消点赞失败")
		return
	}

	// 更新目标对象的点赞数
	if req.TargetType == "topic" {
		var topic models.Topic
		if err := database.DB.First(&topic, req.TargetID).Error; err != nil {
			log.Printf("delete like: topic not found, topicID: %d, error: %v", req.TargetID, err)
		} else if topic.LikeCount > 0 {
			if err := database.DB.Model(&topic).UpdateColumn("like_count", topic.LikeCount-1).Error; err != nil {
				log.Printf("delete like: failed to decrement topic like count, topicID: %d, error: %v", req.TargetID, err)
			}
		}
	} else if req.TargetType == "post" {
		var post models.Post
		if err := database.DB.First(&post, req.TargetID).Error; err != nil {
			log.Printf("delete like: post not found, postID: %d, error: %v", req.TargetID, err)
		} else if post.LikeCount > 0 {
			if err := database.DB.Model(&post).UpdateColumn("like_count", post.LikeCount-1).Error; err != nil {
				log.Printf("delete like: failed to decrement post like count, postID: %d, error: %v", req.TargetID, err)
			}
		}
	}

	log.Printf("delete like: like deleted successfully, userID: %d, targetType: %s, targetID: %d", userID, req.TargetType, req.TargetID)
	utils.Success(w, nil)
}
