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

	"github.com/gorilla/mux"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get profile: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get profile: user not found, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	utils.Success(w, user)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update profile: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update profile: failed to decode request body, userID: %d, error: %v", userID, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if password, ok := updates["password"].(string); ok && password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Printf("update profile: failed to hash password, userID: %d, error: %v", userID, err)
			utils.Error(w, 500, "密码加密失败")
			return
		}
		updates["password_hash"] = hashedPassword
		delete(updates, "password")
	}

	delete(updates, "id")
	delete(updates, "username")
	delete(updates, "email")
	delete(updates, "role")
	delete(updates, "credits")
	delete(updates, "level")
	delete(updates, "created_at")

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		log.Printf("update profile: failed to update profile, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("update profile: user not found after update, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	log.Printf("update profile: profile updated successfully, userID: %d", userID)
	utils.Success(w, user)
}

func GetCurrentUserTopics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get current user topics: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var topics []models.Topic
	var total int64

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		log.Printf("get current user topics: failed to count topics, userID: %d, error: %v", userID, err)
	}

	if err := database.DB.Where("user_id = ?", userID).Preload("User").Preload("Forum").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&topics).Error; err != nil {
		log.Printf("get current user topics: failed to query topics, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取话题失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetCreditUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := database.DB.Order("credits DESC").Limit(10).Find(&users).Error; err != nil {
		log.Printf("get credit users: failed to query users, error: %v", err)
		utils.Error(w, 500, "获取排行榜失败")
		return
	}

	utils.Success(w, users)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		utils.Success(w, []models.User{})
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	searchPattern := "%" + keyword + "%"

	var users []models.User
	var total int64

	if err := database.DB.Model(&models.User{}).Where("username LIKE ? OR nickname LIKE ?", searchPattern, searchPattern).Count(&total).Error; err != nil {
		log.Printf("search users: failed to count users, keyword: %s, error: %v", keyword, err)
	}

	if err := database.DB.Where("username LIKE ? OR nickname LIKE ?", searchPattern, searchPattern).
		Select("id, username, nickname, avatar, signature, created_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		log.Printf("search users: failed to search users, keyword: %s, error: %v", keyword, err)
		utils.Error(w, 500, "搜索失败")
		return
	}

	log.Printf("search users: search completed, keyword: %s, results: %d", keyword, total)
	utils.Success(w, map[string]interface{}{
		"list":  users,
		"total": total,
		"page":  page,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get user: user not found, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	utils.Success(w, user)
}

func GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var topicCount, commentCount int64
	database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&topicCount)
	database.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentCount)

	var followCount, followerCount int64
	database.DB.Model(&models.Follow{}).Where("user_id = ?", userID).Count(&followCount)
	database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", userID).Count(&followerCount)

	utils.Success(w, map[string]interface{}{
		"topic_count":    topicCount,
		"comment_count":  commentCount,
		"follow_count":   followCount,
		"follower_count": followerCount,
	})
}

func GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var followers []models.Follow
	var total int64

	database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", userID).Count(&total)

	if err := database.DB.Where("follow_user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&followers).Error; err != nil {
		log.Printf("get user followers: failed to query followers, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取粉丝列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":  followers,
		"total": total,
		"page":  page,
	})
}

func GetUserTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var topics []models.Topic
	var total int64

	database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&total)

	if err := database.DB.Where("user_id = ?", userID).
		Preload("User").Preload("Forum").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&topics).Error; err != nil {
		log.Printf("get user topics: failed to query topics, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取话题列表失败")
		return
	}

	// 查询用户的勋章
	var userBadges []models.UserBadge
	if err := database.DB.Where("user_id = ? AND is_revoked = ?", userID, false).
		Preload("Badge").
		Find(&userBadges).Error; err != nil {
		log.Printf("get user topics: failed to query user badges, userID: %d, error: %v", userID, err)
	}

	// 为每个话题添加 author_badges
	type TopicWithBadges struct {
		models.Topic
		AuthorBadges []models.UserBadge `json:"author_badges"`
	}
	response := make([]TopicWithBadges, len(topics))
	for i, t := range topics {
		response[i] = TopicWithBadges{Topic: t}
		if t.UserID == uint(userID) {
			response[i].AuthorBadges = userBadges
		}
	}

	utils.Success(w, map[string]interface{}{
		"list":      response,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetFollowTopics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get follow topics: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var followIDs []uint
	database.DB.Model(&models.Follow{}).Where("user_id = ?", userID).Pluck("follow_user_id", &followIDs)

	if len(followIDs) == 0 {
		utils.Success(w, map[string]interface{}{
			"list":      []models.Topic{},
			"total":     0,
			"page":      page,
			"page_size": pageSize,
		})
		return
	}

	var topics []models.Topic
	var total int64

	database.DB.Model(&models.Topic{}).Where("user_id IN ?", followIDs).Count(&total)

	if err := database.DB.Where("user_id IN ?", followIDs).
		Preload("User").Preload("Forum").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&topics).Error; err != nil {
		log.Printf("get follow topics: failed to query topics, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取关注动态失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
