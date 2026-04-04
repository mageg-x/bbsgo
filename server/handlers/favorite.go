package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
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
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 解析请求体
	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create favorite: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 检查是否已经收藏
	var existingFavorite models.Favorite
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, req.TopicID).First(&existingFavorite).Error; err == nil {
		log.Printf("create favorite: already favorited, userID: %d, topicID: %d", userID, req.TopicID)
		errors.Error(w, errors.CodeAlreadyFollowed, "")
		return
	}

	// 创建收藏记录
	favorite := models.Favorite{
		UserID:  userID,
		TopicID: req.TopicID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		log.Printf("create favorite: failed to create favorite, userID: %d, topicID: %d, error: %v", userID, req.TopicID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("create favorite: favorite created successfully, userID: %d, topicID: %d", userID, req.TopicID)
	errors.Success(w, favorite)
}

// DeleteFavorite 取消收藏处理器
// 用户取消对话题的收藏
func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete favorite: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 从 URL 查询参数获取
	topicIDStr := r.URL.Query().Get("topic_id")
	topicID, _ := strconv.Atoi(topicIDStr)

	if topicID == 0 {
		log.Printf("delete favorite: invalid parameters, topicID: %d", topicID)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 查询收藏记录
	var favorite models.Favorite
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&favorite).Error; err != nil {
		log.Printf("delete favorite: favorite not found, userID: %d, topicID: %d, error: %v", userID, topicID, err)
		errors.Error(w, errors.CodeFavoriteNotFound, "")
		return
	}

	// 删除收藏记录
	if err := database.DB.Unscoped().Delete(&favorite).Error; err != nil {
		log.Printf("delete favorite: failed to delete favorite, id: %d, error: %v", favorite.ID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("delete favorite: favorite deleted successfully, userID: %d, topicID: %d", userID, topicID)
	errors.Success(w, nil)
}

// CheckFavorite 检查收藏状态处理器
func CheckFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		errors.Success(w, map[string]interface{}{"favorited": false})
		return
	}

	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.Success(w, map[string]interface{}{"favorited": false})
		return
	}

	var count int64
	if err := database.DB.Model(&models.Favorite{}).
		Where("user_id = ? AND topic_id = ?", userID, req.TopicID).
		Count(&count).Error; err != nil {
		errors.Success(w, map[string]interface{}{"favorited": false})
		return
	}

	errors.Success(w, map[string]interface{}{"favorited": count > 0})
}

// GetFavorites 获取当前用户的收藏列表处理器
func GetFavorites(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get favorites: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 查询收藏列表
	var favorites []models.Favorite
	if err := database.DB.Where("user_id = ?", userID).Preload("Topic").Preload("Topic.User").Preload("Topic.Forum").
		Order("created_at DESC").Find(&favorites).Error; err != nil {
		log.Printf("get favorites: failed to query favorites, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
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

	// 收集所有作者用户ID
	userIDs := make(map[uint]bool)
	for _, fav := range favorites {
		if fav.Topic.UserID > 0 {
			userIDs[fav.Topic.UserID] = true
		}
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
			log.Printf("get favorites: failed to query user badges, error: %v", err)
		}
		for _, ub := range userBadges {
			userBadgesMap[ub.UserID] = append(userBadgesMap[ub.UserID], ub)
		}
	}

	// 构建返回数据，添加 has_poll 和 author_badges
	type FavoriteWithPoll struct {
		models.Favorite
		TopicHasPoll  bool                    `json:"topic_has_poll"`
		AuthorBadges  []models.UserBadge      `json:"author_badges"`
	}
	var response []FavoriteWithPoll
	for _, fav := range favorites {
		response = append(response, FavoriteWithPoll{
			Favorite:     fav,
			TopicHasPoll: hasPollMap[fav.TopicID],
			AuthorBadges: userBadgesMap[fav.Topic.UserID],
		})
	}

	errors.Success(w, response)
}
