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

// FavoriteRequest 收藏请求结构
type FavoriteRequest struct {
	TopicID uint `json:"topic_id"` // 话题ID
}

// CreateFavorite 创建收藏处理器
// 用户收藏指定话题
func CreateFavorite(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create favorite: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 解析请求体
	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create favorite: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 检查是否已经收藏
	var existingFavorite models.Favorite
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, req.TopicID).First(&existingFavorite).Error; err == nil {
		log.Printf("create favorite: already favorited, userID: %d, topicID: %d", userID, req.TopicID)
		utils.Error(w, 400, "已经收藏过了")
		return
	}

	// 创建收藏记录
	favorite := models.Favorite{
		UserID:  userID,
		TopicID: req.TopicID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		log.Printf("create favorite: failed to create favorite, userID: %d, topicID: %d, error: %v", userID, req.TopicID, err)
		utils.Error(w, 500, "收藏失败")
		return
	}

	log.Printf("create favorite: favorite created successfully, userID: %d, topicID: %d", userID, req.TopicID)
	utils.Success(w, favorite)
}

// DeleteFavorite 取消收藏处理器
// 用户取消对话题的收藏
func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete favorite: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 解析请求体
	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("delete favorite: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 查询收藏记录
	var favorite models.Favorite
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, req.TopicID).First(&favorite).Error; err != nil {
		log.Printf("delete favorite: favorite not found, userID: %d, topicID: %d, error: %v", userID, req.TopicID, err)
		utils.Error(w, 404, "收藏记录不存在")
		return
	}

	// 删除收藏记录
	if err := database.DB.Unscoped().Delete(&favorite).Error; err != nil {
		log.Printf("delete favorite: failed to delete favorite, id: %d, error: %v", favorite.ID, err)
		utils.Error(w, 500, "取消收藏失败")
		return
	}

	log.Printf("delete favorite: favorite deleted successfully, userID: %d, topicID: %d", userID, req.TopicID)
	utils.Success(w, nil)
}

// GetFavorites 获取当前用户的收藏列表处理器
func GetFavorites(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get favorites: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 查询收藏列表
	var favorites []models.Favorite
	if err := database.DB.Where("user_id = ?", userID).Preload("Topic").Preload("Topic.User").Preload("Topic.Forum").
		Order("created_at DESC").Find(&favorites).Error; err != nil {
		log.Printf("get favorites: failed to query favorites, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取收藏列表失败")
		return
	}

	utils.Success(w, favorites)
}
