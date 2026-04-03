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

	// 统计用户评论数
	if err := database.DB.Model(&models.Comment{}).Where("user_id = ?", id).Count(&postCount).Error; err != nil {
		log.Printf("get user stats: failed to count comments, id: %d, error: %v", id, err)
	}

	// 计算用户排名（根据ID顺序估算）
	if err := database.DB.Table("users").Where("id <= ?", id).Count(&rank).Error; err != nil {
		log.Printf("get user stats: failed to calculate rank, id: %d, error: %v", id, err)
	}

	utils.Success(w, map[string]interface{}{
		"topic_count":   topicCount,
		"comment_count": postCount,
		"rank":          rank,
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
		Order("is_user_pinned DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&topics).Error; err != nil {
		log.Printf("get user topics: failed to query topics, id: %d, error: %v", id, err)
		utils.Error(w, 500, "获取话题列表失败")
		return
	}

	// 给每个 topic 添加 has_poll 字段
	var topicIDs []uint
	for _, t := range topics {
		topicIDs = append(topicIDs, t.ID)
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
	for _, t := range topics {
		userIDs[t.UserID] = true
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
			log.Printf("get user topics: failed to query user badges, error: %v", err)
		}
		for _, ub := range userBadges {
			userBadgesMap[ub.UserID] = append(userBadgesMap[ub.UserID], ub)
		}
	}

	// 构建返回数据，添加 has_poll 和 author_badges
	type TopicWithPoll struct {
		models.Topic
		HasPoll      bool                    `json:"has_poll"`
		AuthorBadges []models.UserBadge      `json:"author_badges"`
	}
	var response []TopicWithPoll
	for _, t := range topics {
		response = append(response, TopicWithPoll{
			Topic:       t,
			HasPoll:     hasPollMap[t.ID],
			AuthorBadges: userBadgesMap[t.UserID],
		})
	}

	utils.Success(w, map[string]interface{}{
		"list":      response,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
