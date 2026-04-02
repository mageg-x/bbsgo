package handlers

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetUser 获取指定用户公开资料处理器
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("get user: invalid user id, id: %s, error: %v", vars["id"], err)
		utils.Error(w, 400, "无效的用户ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("get user: user not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	utils.Success(w, user)
}

// GetUserStats 获取指定用户的统计数据处理器
// 包括话题数、回复数、排名
func GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("get user stats: invalid user id, id: %s, error: %v", vars["id"], err)
		utils.Error(w, 400, "无效的用户ID")
		return
	}

	var topicCount int64
	var postCount int64
	var rank int64

	// 统计用户话题数
	if err := database.DB.Model(&models.Topic{}).Where("user_id = ?", id).Count(&topicCount).Error; err != nil {
		log.Printf("get user stats: failed to count topics, id: %d, error: %v", id, err)
	}

	// 统计用户回复数
	if err := database.DB.Model(&models.Post{}).Where("user_id = ?", id).Count(&postCount).Error; err != nil {
		log.Printf("get user stats: failed to count posts, id: %d, error: %v", id, err)
	}

	// 计算用户排名（根据ID顺序估算）
	if err := database.DB.Table("users").Where("id <= ?", id).Count(&rank).Error; err != nil {
		log.Printf("get user stats: failed to calculate rank, id: %d, error: %v", id, err)
	}

	utils.Success(w, map[string]interface{}{
		"topic_count": topicCount,
		"post_count":  postCount,
		"rank":        rank,
	})
}

// GetUserFollowers 获取指定用户的粉丝列表处理器
func GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("get user followers: invalid user id, id: %s, error: %v", vars["id"], err)
		utils.Error(w, 400, "无效的用户ID")
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

	var follows []models.Follow
	var total int64

	offset := (page - 1) * pageSize

	// 统计粉丝数量
	if err := database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", id).Count(&total).Error; err != nil {
		log.Printf("get user followers: failed to count followers, id: %d, error: %v", id, err)
	}

	// 查询粉丝列表
	if err := database.DB.Where("follow_user_id = ?", id).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&follows).Error; err != nil {
		log.Printf("get user followers: failed to query followers, id: %d, error: %v", id, err)
		utils.Error(w, 500, "获取粉丝列表失败")
		return
	}

	// 提取用户信息
	followers := make([]models.User, 0)
	for _, follow := range follows {
		followers = append(followers, follow.User)
	}

	utils.Success(w, map[string]interface{}{
		"list":      followers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUserTopics 获取指定用户发布的话题列表处理器
func GetUserTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("get user topics: invalid user id, id: %s, error: %v", vars["id"], err)
		utils.Error(w, 400, "无效的用户ID")
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
	if err := database.DB.Model(&models.Topic{}).Where("user_id = ?", id).Count(&total).Error; err != nil {
		log.Printf("get user topics: failed to count topics, id: %d, error: %v", id, err)
	}

	// 查询话题列表
	if err := database.DB.Where("user_id = ?", id).
		Preload("User").
		Preload("Forum").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&topics).Error; err != nil {
		log.Printf("get user topics: failed to query topics, id: %d, error: %v", id, err)
		utils.Error(w, 500, "获取话题列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
