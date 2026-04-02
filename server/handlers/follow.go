package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// GetFollows 获取当前用户关注的人列表处理器
func GetFollows(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var follows []models.Follow
	if err := database.DB.Where("user_id = ?", userID).Preload("FollowUser").Find(&follows).Error; err != nil {
		log.Printf("get follows: failed to query follows, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取关注列表失败")
		return
	}

	utils.Success(w, follows)
}

// GetFollowers 获取当前用户的粉丝列表处理器
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var followers []models.Follow
	if err := database.DB.Where("follow_user_id = ?", userID).Preload("User").Find(&followers).Error; err != nil {
		log.Printf("get followers: failed to query followers, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取粉丝列表失败")
		return
	}

	utils.Success(w, followers)
}

// CreateFollow 创建关注处理器
// 用户关注另一个用户
func CreateFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 解析请求体
	var req struct {
		FollowUserID uint `json:"follow_user_id"` // 要关注的用户ID
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create follow: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 不能关注自己
	if req.FollowUserID == userID {
		log.Printf("create follow: cannot follow yourself, userID: %d", userID)
		utils.Error(w, 400, "不能关注自己")
		return
	}

	// 检查要关注的用户是否存在
	var followUser models.User
	if err := database.DB.First(&followUser, req.FollowUserID).Error; err != nil {
		log.Printf("create follow: user not found, followUserID: %d, error: %v", req.FollowUserID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 检查是否已经关注
	var existing models.Follow
	if err := database.DB.Where("user_id = ? AND follow_user_id = ?", userID, req.FollowUserID).First(&existing).Error; err == nil {
		log.Printf("create follow: already following, userID: %d, followUserID: %d", userID, req.FollowUserID)
		utils.Error(w, 400, "已关注该用户")
		return
	}

	// 创建关注记录
	follow := models.Follow{
		UserID:       userID,
		FollowUserID: req.FollowUserID,
	}

	if err := database.DB.Create(&follow).Error; err != nil {
		log.Printf("create follow: failed to create follow, userID: %d, followUserID: %d, error: %v", userID, req.FollowUserID, err)
		utils.Error(w, 500, "关注失败")
		return
	}

	log.Printf("create follow: follow created successfully, userID: %d, followUserID: %d", userID, req.FollowUserID)
	utils.Success(w, follow)
}

// DeleteFollow 取消关注处理器
// 用户取消关注另一个用户
func DeleteFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 解析请求体
	var req struct {
		FollowUserID uint `json:"follow_user_id"` // 要取消关注的用户ID
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("delete follow: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 删除关注记录
	if err := database.DB.Unscoped().Where("user_id = ? AND follow_user_id = ?", userID, req.FollowUserID).Delete(&models.Follow{}).Error; err != nil {
		log.Printf("delete follow: failed to delete follow, userID: %d, followUserID: %d, error: %v", userID, req.FollowUserID, err)
		utils.Error(w, 500, "取消关注失败")
		return
	}

	log.Printf("delete follow: follow deleted successfully, userID: %d, followUserID: %d", userID, req.FollowUserID)
	utils.Success(w, nil)
}

// CheckFollow 检查是否关注了指定用户处理器
func CheckFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	targetUserID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	var follow models.Follow
	isFollowing := database.DB.Where("user_id = ? AND follow_user_id = ?", userID, targetUserID).First(&follow).Error == nil

	utils.Success(w, map[string]bool{"is_following": isFollowing})
}
