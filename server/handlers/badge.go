package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// GetBadges 获取所有勋章列表处理器
func GetBadges(w http.ResponseWriter, r *http.Request) {
	var badges []models.Badge
	if err := database.DB.Order("sort_order ASC, id ASC").Find(&badges).Error; err != nil {
		log.Printf("get badges: failed to query badges, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, badges)
}

// GetUserBadges 获取当前用户获得的勋章列表处理器
func GetUserBadges(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var userBadges []models.UserBadge
	if err := database.DB.Where("user_id = ? AND is_revoked = ?", userID, false).
		Preload("Badge").Order("awarded_at DESC").Find(&userBadges).Error; err != nil {
		log.Printf("get user badges: failed to query user badges, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, userBadges)
}

// GetUserBadgeProgress 获取用户勋章进度处理器
func GetUserBadgeProgress(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	progress, err := badgeService.GetUserBadgeProgress(userID)
	if err != nil {
		log.Printf("get user badge progress: failed, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, progress)
}

// GetUserBadgesByID 获取指定用户的勋章列表处理器（公开接口）
func GetUserBadgesByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var userBadges []models.UserBadge
	if err := database.DB.Where("user_id = ? AND is_revoked = ?", userID, false).
		Preload("Badge").Order("awarded_at DESC").Find(&userBadges).Error; err != nil {
		log.Printf("get user badges by id: failed to query user badges, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, userBadges)
}

// CreateBadge 创建勋章处理器（管理员）
func CreateBadge(w http.ResponseWriter, r *http.Request) {
	var badge models.Badge
	if err := json.NewDecoder(r.Body).Decode(&badge); err != nil {
		log.Printf("create badge: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	if badge.Name == "" {
		log.Printf("create badge: badge name is empty")
		errors.Error(w, errors.CodeIncompleteInfo, "")
		return
	}

	badge.CreatedAt = time.Now()
	if err := database.DB.Create(&badge).Error; err != nil {
		log.Printf("create badge: failed to create badge, name: %s, error: %v", badge.Name, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("create badge: badge created successfully, id: %d, name: %s", badge.ID, badge.Name)
	errors.Success(w, badge)
}

// UpdateBadge 更新勋章处理器（管理员）
func UpdateBadge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var badge models.Badge
	if err := database.DB.First(&badge, id).Error; err != nil {
		log.Printf("update badge: badge not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	var req models.Badge
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update badge: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	updates := map[string]interface{}{
		"name":            req.Name,
		"description":     req.Description,
		"icon":            req.Icon,
		"condition":       req.Condition,
		"type":            req.Type,
		"condition_type":  req.ConditionType,
		"condition_value": req.ConditionValue,
		"condition_data":  req.ConditionData,
		"sort_order":      req.SortOrder,
	}

	if err := database.DB.Model(&badge).Updates(updates).Error; err != nil {
		log.Printf("update badge: failed to update badge, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	database.DB.First(&badge, id)
	log.Printf("update badge: badge updated successfully, id: %d", id)
	errors.Success(w, badge)
}

// DeleteBadge 删除勋章处理器（管理员）
func DeleteBadge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var badge models.Badge
	if err := database.DB.First(&badge, id).Error; err != nil {
		log.Printf("delete badge: badge not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	if err := database.DB.Unscoped().Delete(&badge).Error; err != nil {
		log.Printf("delete badge: failed to delete badge, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("delete badge: badge deleted successfully, id: %d", id)
	errors.Success(w, nil)
}

// AwardBadge 授予用户勋章处理器（管理员）
func AwardBadge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  uint `json:"user_id"`
		BadgeID uint `json:"badge_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("award badge: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		log.Printf("award badge: user not found, userID: %d, error: %v", req.UserID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	var badge models.Badge
	if err := database.DB.First(&badge, req.BadgeID).Error; err != nil {
		log.Printf("award badge: badge not found, badgeID: %d, error: %v", req.BadgeID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	var existing models.UserBadge
	if err := database.DB.Where("user_id = ? AND badge_id = ? AND is_revoked = ?", req.UserID, req.BadgeID, false).
		First(&existing).Error; err == nil {
		log.Printf("award badge: user already has this badge, userID: %d, badgeID: %d", req.UserID, req.BadgeID)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	userBadge := models.UserBadge{
		UserID:    req.UserID,
		BadgeID:   req.BadgeID,
		AwardedAt: time.Now(),
	}

	if err := database.DB.Create(&userBadge).Error; err != nil {
		log.Printf("award badge: failed to create user badge, userID: %d, badgeID: %d, error: %v", req.UserID, req.BadgeID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 发送勋章获得通知
	badgeSvc := services.NewBadgeService()
	badgeSvc.SendBadgeNotification(req.UserID, req.BadgeID)

	log.Printf("award badge: badge awarded successfully, userID: %d, badgeID: %d", req.UserID, req.BadgeID)
	errors.Success(w, userBadge)
}

// RevokeBadge 撤销用户勋章处理器（管理员）
func RevokeBadge(w http.ResponseWriter, r *http.Request) {
	adminID, _ := middleware.GetUserIDFromContext(r.Context())

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("revoke badge: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	var userBadge models.UserBadge
	if err := database.DB.First(&userBadge, id).Error; err != nil {
		log.Printf("revoke badge: user badge not found, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	if userBadge.IsRevoked {
		log.Printf("revoke badge: badge already revoked, id: %d", id)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"is_revoked":     true,
		"revoked_at":     &now,
		"revoked_reason": req.Reason,
		"revoked_by":     adminID,
	}

	if err := database.DB.Model(&userBadge).Updates(updates).Error; err != nil {
		log.Printf("revoke badge: failed to revoke badge, id: %d, error: %v", id, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("revoke badge: badge revoked successfully, id: %d", id)
	errors.Success(w, nil)
}

// GetAdminBadges 获取勋章列表（管理员）
func GetAdminBadges(w http.ResponseWriter, r *http.Request) {
	var badges []models.Badge
	if err := database.DB.Order("sort_order ASC, id ASC").Find(&badges).Error; err != nil {
		log.Printf("get admin badges: failed to query badges, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	var result []map[string]interface{}
	for _, badge := range badges {
		var count int64
		database.DB.Model(&models.UserBadge{}).
			Where("badge_id = ? AND is_revoked = ?", badge.ID, false).
			Count(&count)

		result = append(result, map[string]interface{}{
			"id":              badge.ID,
			"name":            badge.Name,
			"description":     badge.Description,
			"icon":            badge.Icon,
			"condition":       badge.Condition,
			"type":            badge.Type,
			"condition_type":  badge.ConditionType,
			"condition_value": badge.ConditionValue,
			"condition_data":  badge.ConditionData,
			"sort_order":      badge.SortOrder,
			"created_at":      badge.CreatedAt,
			"award_count":     count,
		})
	}

	errors.Success(w, result)
}

// GetBadgeUsers 获取获得指定勋章的用户列表（管理员）
func GetBadgeUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	badgeID, _ := strconv.Atoi(vars["id"])

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.UserBadge{}).
		Where("badge_id = ? AND is_revoked = ?", badgeID, false).
		Count(&total)

	var userBadges []models.UserBadge
	if err := database.DB.Where("badge_id = ? AND is_revoked = ?", badgeID, false).
		Preload("User").Preload("Badge").
		Order("awarded_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&userBadges).Error; err != nil {
		log.Printf("get badge users: failed to query user badges, badgeID: %d, error: %v", badgeID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
		"list":  userBadges,
		"total": total,
		"page":  page,
	})
}

// InitBadges 初始化系统勋章（默认中文，仅供管理接口调用）
func InitBadges(w http.ResponseWriter, r *http.Request) {
	SeedBadgesWithLang("zh")
	errors.Success(w, map[string]string{"message": "勋章初始化成功"})
}

// SeedBadgesWithLang 使用指定语言初始化系统勋章
func SeedBadgesWithLang(lang string) {
	var badges []models.Badge
	if lang == "en" {
		badges = getBadgeDataEn()
	} else {
		badges = getBadgeDataZh()
	}

	for _, badge := range badges {
		var existing models.Badge
		if err := database.DB.Where("name = ?", badge.Name).First(&existing).Error; err != nil {
			database.DB.Create(&badge)
			log.Printf("seed badges: created badge: %s", badge.Name)
		}
	}

	log.Printf("seed badges: initialization completed (%s)", lang)
}

func getBadgeDataZh() []models.Badge {
	return []models.Badge{
		{
			Name:           "初来乍到",
			Description:    "注册成功即获得，欢迎加入社区！",
			Icon:           "newcomer",
			Condition:      "注册成功",
			Type:           "basic",
			ConditionType:  "register",
			ConditionValue: 1,
			SortOrder:      1,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "首次发声",
			Description:    "发布第1个帖子，开始你的社区之旅！",
			Icon:           "first-post",
			Condition:      "发布1个帖子",
			Type:           "basic",
			ConditionType:  "topic_count",
			ConditionValue: 1,
			SortOrder:      2,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "热心回复",
			Description:    "发布第1条评论，积极参与讨论！",
			Icon:           "first-comment",
			Condition:      "发布1条评论",
			Type:           "basic",
			ConditionType:  "comment_count",
			ConditionValue: 1,
			SortOrder:      3,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "笔耕不辍",
			Description:    "发布50个帖子，持续输出优质内容！",
			Icon:           "writer",
			Condition:      "发布50个帖子",
			Type:           "advanced",
			ConditionType:  "topic_count",
			ConditionValue: 50,
			SortOrder:      4,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "社区活宝",
			Description:    "发布1000条评论，成为社区活跃分子！",
			Icon:           "community-star",
			Condition:      "发布1000条评论",
			Type:           "advanced",
			ConditionType:  "comment_count",
			ConditionValue: 1000,
			SortOrder:      5,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "广受欢迎",
			Description:    "累计获得1000个点赞，内容质量被认可！",
			Icon:           "popular",
			Condition:      "获得1000个点赞",
			Type:           "advanced",
			ConditionType:  "like_count",
			ConditionValue: 1000,
			SortOrder:      6,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "金牌评论",
			Description:    "获得5次最佳评论，被楼主高度认可！",
			Icon:           "gold-comment",
			Condition:      "获得5次最佳评论",
			Type:           "advanced",
			ConditionType:  "best_comment",
			ConditionValue: 5,
			SortOrder:      7,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "意见领袖",
			Description:    "被500个用户关注，拥有社交影响力！",
			Icon:           "opinion-leader",
			Condition:      "被500个用户关注",
			Type:           "top",
			ConditionType:  "follower_count",
			ConditionValue: 500,
			SortOrder:      8,
			CreatedAt:      time.Now(),
		},
		{
			Name:          "社区传奇",
			Description:   "注册满2年 + 发帖≥200 + 获赞≥500 + 最佳评论≥10，终身成就！",
			Icon:          "legend",
			Condition:     "注册满2年 + 发帖≥200 + 获赞≥500 + 最佳评论≥10",
			Type:          "top",
			ConditionType: "combination",
			ConditionData: `{"register_days": 730, "topic_count": 200, "like_count": 500, "best_comment": 10}`,
			SortOrder:     9,
			CreatedAt:     time.Now(),
		},
	}
}

func getBadgeDataEn() []models.Badge {
	return []models.Badge{
		{
			Name:           "Newcomer",
			Description:    "Welcome to the community! Earned upon registration.",
			Icon:           "newcomer",
			Condition:      "Register successfully",
			Type:           "basic",
			ConditionType:  "register",
			ConditionValue: 1,
			SortOrder:      1,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "First Post",
			Description:    "Post your first topic and start your journey!",
			Icon:           "first-post",
			Condition:      "Post 1 topic",
			Type:           "basic",
			ConditionType:  "topic_count",
			ConditionValue: 1,
			SortOrder:      2,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "Helpful Reply",
			Description:    "Post your first comment and join the discussion!",
			Icon:           "first-comment",
			Condition:      "Post 1 comment",
			Type:           "basic",
			ConditionType:  "comment_count",
			ConditionValue: 1,
			SortOrder:      3,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "Prolific Writer",
			Description:    "Post 50 topics and keep producing quality content!",
			Icon:           "writer",
			Condition:      "Post 50 topics",
			Type:           "advanced",
			ConditionType:  "topic_count",
			ConditionValue: 50,
			SortOrder:      4,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "Community Star",
			Description:    "Post 1000 comments and become an active member!",
			Icon:           "community-star",
			Condition:      "Post 1000 comments",
			Type:           "advanced",
			ConditionType:  "comment_count",
			ConditionValue: 1000,
			SortOrder:      5,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "Popular",
			Description:    "Receive 1000 likes, your content is recognized!",
			Icon:           "popular",
			Condition:      "Receive 1000 likes",
			Type:           "advanced",
			ConditionType:  "like_count",
			ConditionValue: 1000,
			SortOrder:      6,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "Golden Comment",
			Description:    "Receive 5 best comments, highly recognized by thread owners!",
			Icon:           "gold-comment",
			Condition:      "Receive 5 best comments",
			Type:           "advanced",
			ConditionType:  "best_comment",
			ConditionValue: 5,
			SortOrder:      7,
			CreatedAt:      time.Now(),
		},
		{
			Name:           "Opinion Leader",
			Description:    "Followed by 500 users, have social influence!",
			Icon:           "opinion-leader",
			Condition:      "Followed by 500 users",
			Type:           "top",
			ConditionType:  "follower_count",
			ConditionValue: 500,
			SortOrder:      8,
			CreatedAt:      time.Now(),
		},
		{
			Name:          "Community Legend",
			Description:   "2 years registration + 200 posts + 500 likes + 10 best comments - Lifetime Achievement!",
			Icon:          "legend",
			Condition:     "2 years + 200 posts + 500 likes + 10 best comments",
			Type:          "top",
			ConditionType: "combination",
			ConditionData: `{"register_days": 730, "topic_count": 200, "like_count": 500, "best_comment": 10}`,
			SortOrder:     9,
			CreatedAt:     time.Now(),
		},
	}
}
