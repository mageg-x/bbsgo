package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/services"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var followBadgeService = services.NewBadgeService()

// GetFollows 获取当前用户关注的人列表处理器
func GetFollows(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var follows []models.Follow
	if err := database.DB.Where("user_id = ?", userID).Preload("FollowUser").Find(&follows).Error; err != nil {
		log.Printf("get follows: failed to query follows, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, follows)
}

// GetFollowers 获取当前用户的粉丝列表处理器
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var followers []models.Follow
	if err := database.DB.Where("follow_user_id = ?", userID).Preload("User").Find(&followers).Error; err != nil {
		log.Printf("get followers: failed to query followers, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, followers)
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
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 不能关注自己
	if req.FollowUserID == userID {
		log.Printf("create follow: cannot follow yourself, userID: %d", userID)
		errors.Error(w, errors.CodeCannotFollowSelf, "")
		return
	}

	// 检查要关注的用户是否存在
	var followUser models.User
	if err := database.DB.First(&followUser, req.FollowUserID).Error; err != nil {
		log.Printf("create follow: user not found, followUserID: %d, error: %v", req.FollowUserID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	// 检查是否已经关注
	var existing models.Follow
	if err := database.DB.Where("user_id = ? AND follow_user_id = ?", userID, req.FollowUserID).First(&existing).Error; err == nil {
		log.Printf("create follow: already following, userID: %d, followUserID: %d", userID, req.FollowUserID)
		errors.Error(w, errors.CodeAlreadyFollowed, "")
		return
	}

	// 创建关注记录
	follow := models.Follow{
		UserID:       userID,
		FollowUserID: req.FollowUserID,
	}

	if err := database.DB.Create(&follow).Error; err != nil {
		log.Printf("create follow: failed to create follow, userID: %d, followUserID: %d, error: %v", userID, req.FollowUserID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 触发被关注者的勋章检查（因为增加了粉丝）
	go followBadgeService.CheckAndAwardBadges(req.FollowUserID)

	// 发送关注通知
	var follower models.User
	if err := database.DB.First(&follower, userID).Error; err == nil {
		CreateNotification(
			req.FollowUserID,
			"follow",
			"用户 "+follower.Username+" 关注了你",
			"/user/"+strconv.FormatUint(uint64(userID), 10),
		)
	}

	// 获得关注积分奖励
	var followedUser models.User
	if err := database.DB.First(&followedUser, req.FollowUserID).Error; err == nil {
		creditAmount := utils.GetConfigInt("credit_follow", 1)
		followedUser.Credits += creditAmount
		if err := database.DB.Save(&followedUser).Error; err != nil {
			log.Printf("create follow: failed to add credits, userID: %d, error: %v", req.FollowUserID, err)
		} else {
			log.Printf("create follow: awarded %d credits to userID: %d", creditAmount, req.FollowUserID)
		}
	}

	log.Printf("create follow: follow created successfully, userID: %d, followUserID: %d", userID, req.FollowUserID)
	errors.Success(w, follow)
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
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 删除关注记录
	if err := database.DB.Unscoped().Where("user_id = ? AND follow_user_id = ?", userID, req.FollowUserID).Delete(&models.Follow{}).Error; err != nil {
		log.Printf("delete follow: failed to delete follow, userID: %d, followUserID: %d, error: %v", userID, req.FollowUserID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("delete follow: follow deleted successfully, userID: %d, followUserID: %d", userID, req.FollowUserID)
	errors.Success(w, nil)
}

// CheckFollow 检查是否关注了指定用户处理器
func CheckFollow(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	targetUserID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	var follow models.Follow
	isFollowing := database.DB.Where("user_id = ? AND follow_user_id = ?", userID, targetUserID).First(&follow).Error == nil

	errors.Success(w, map[string]bool{"is_following": isFollowing})
}
