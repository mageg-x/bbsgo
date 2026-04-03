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

// GetProfile 获取当前用户个人资料处理器
func GetProfile(w http.ResponseWriter, r *http.Request) {
	// 获取当前用户ID
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get profile: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get profile: user not found, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	utils.Success(w, user)
}

// UpdateProfile 更新当前用户个人资料处理器
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// 获取当前用户ID
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update profile: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update profile: failed to decode request body, userID: %d, error: %v", userID, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 处理密码修改
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

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "username")
	delete(updates, "email")
	delete(updates, "role")
	delete(updates, "credits")
	delete(updates, "level")
	delete(updates, "created_at")

	// 执行更新
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		log.Printf("update profile: failed to update profile, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	// 重新加载用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("update profile: user not found after update, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	log.Printf("update profile: profile updated successfully, userID: %d", userID)
	utils.Success(w, user)
}

// GetCurrentUserTopics 获取当前用户发布的话题列表处理器
func GetCurrentUserTopics(w http.ResponseWriter, r *http.Request) {
	// 获取当前用户ID
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get current user topics: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 解析分页参数
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

	// 统计话题数量
	if err := database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		log.Printf("get current user topics: failed to count topics, userID: %d, error: %v", userID, err)
	}

	// 查询话题
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

// GetCreditUsers 获取积分排行榜处理器
// 返回积分最高的前10名用户
func GetCreditUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := database.DB.Order("credits DESC").Limit(10).Find(&users).Error; err != nil {
		log.Printf("get credit users: failed to query users, error: %v", err)
		utils.Error(w, 500, "获取排行榜失败")
		return
	}

	utils.Success(w, users)
}

// SearchUsers 搜索用户处理器
// 根据关键词搜索用户名或昵称
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

	// 统计匹配的用户数量
	if err := database.DB.Model(&models.User{}).Where("username LIKE ? OR nickname LIKE ?", searchPattern, searchPattern).Count(&total).Error; err != nil {
		log.Printf("search users: failed to count users, keyword: %s, error: %v", keyword, err)
	}

	// 搜索用户
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
