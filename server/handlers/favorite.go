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

	// 从 URL 查询参数获取
	topicIDStr := r.URL.Query().Get("topic_id")
	topicID, _ := strconv.Atoi(topicIDStr)

	if topicID == 0 {
		log.Printf("delete favorite: invalid parameters, topicID: %d", topicID)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 查询收藏记录
	var favorite models.Favorite
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&favorite).Error; err != nil {
		log.Printf("delete favorite: favorite not found, userID: %d, topicID: %d, error: %v", userID, topicID, err)
		utils.Error(w, 404, "收藏记录不存在")
		return
	}

	// 删除收藏记录
	if err := database.DB.Unscoped().Delete(&favorite).Error; err != nil {
		log.Printf("delete favorite: failed to delete favorite, id: %d, error: %v", favorite.ID, err)
		utils.Error(w, 500, "取消收藏失败")
		return
	}

	log.Printf("delete favorite: favorite deleted successfully, userID: %d, topicID: %d", userID, topicID)
	utils.Success(w, nil)
}

// CheckFavorite 检查收藏状态处理器
func CheckFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.Success(w, map[string]interface{}{"favorited": false})
		return
	}

	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Success(w, map[string]interface{}{"favorited": false})
		return
	}

	var count int64
	if err := database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND topic_id = ?", userID, req.TopicID).
		Count(&count).Error; err != nil {
		utils.Success(w, map[string]interface{}{"favorited": false})
		return
	}

	utils.Success(w, map[string]interface{}{"favorited": count > 0})
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

	// 给每个 topic 添加 has_poll 字段
	var topicIDs []uint
	for _, fav := range favorites {
		if fav.TopicID > 0 {
			topicIDs = append(topicIDs, fav.TopicID)
		}
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

	// 构建返回数据，添加 has_poll
	type FavoriteWithPoll struct {
		models.Favorite
		TopicHasPoll bool `json:"topic_has_poll"`
	}
	var response []FavoriteWithPoll
	for _, fav := range favorites {
		response = append(response, FavoriteWithPoll{
			Favorite:     fav,
			TopicHasPoll: hasPollMap[fav.TopicID],
		})
	}

	utils.Success(w, response)
}
