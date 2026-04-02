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

// GetBadges 获取所有勋章列表处理器
func GetBadges(w http.ResponseWriter, r *http.Request) {
	var badges []models.Badge
	if err := database.DB.Find(&badges).Error; err != nil {
		log.Printf("get badges: failed to query badges, error: %v", err)
		utils.Error(w, 500, "获取勋章列表失败")
		return
	}

	utils.Success(w, badges)
}

// GetUserBadges 获取当前用户获得的勋章列表处理器
func GetUserBadges(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var userBadges []models.UserBadge
	if err := database.DB.Where("user_id = ?", userID).Preload("Badge").Find(&userBadges).Error; err != nil {
		log.Printf("get user badges: failed to query user badges, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取用户勋章列表失败")
		return
	}

	utils.Success(w, userBadges)
}

// CreateBadge 创建勋章处理器（管理员）
func CreateBadge(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var badge models.Badge
	if err := json.NewDecoder(r.Body).Decode(&badge); err != nil {
		log.Printf("create badge: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证勋章名称
	if badge.Name == "" {
		log.Printf("create badge: badge name is empty")
		utils.Error(w, 400, "请填写勋章名称")
		return
	}

	// 创建勋章
	if err := database.DB.Create(&badge).Error; err != nil {
		log.Printf("create badge: failed to create badge, name: %s, error: %v", badge.Name, err)
		utils.Error(w, 500, "创建失败")
		return
	}

	log.Printf("create badge: badge created successfully, id: %d, name: %s", badge.ID, badge.Name)
	utils.Success(w, badge)
}

// AwardBadge 授予用户勋章处理器（管理员）
func AwardBadge(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var req struct {
		UserID  uint `json:"user_id"`  // 用户ID
		BadgeID uint `json:"badge_id"` // 勋章ID
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("award badge: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		log.Printf("award badge: user not found, userID: %d, error: %v", req.UserID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 检查勋章是否存在
	var badge models.Badge
	if err := database.DB.First(&badge, req.BadgeID).Error; err != nil {
		log.Printf("award badge: badge not found, badgeID: %d, error: %v", req.BadgeID, err)
		utils.Error(w, 404, "勋章不存在")
		return
	}

	// 检查用户是否已获得该勋章
	var existing models.UserBadge
	if err := database.DB.Where("user_id = ? AND badge_id = ?", req.UserID, req.BadgeID).First(&existing).Error; err == nil {
		log.Printf("award badge: user already has this badge, userID: %d, badgeID: %d", req.UserID, req.BadgeID)
		utils.Error(w, 400, "用户已获得该勋章")
		return
	}

	// 创建用户勋章记录
	userBadge := models.UserBadge{
		UserID:  req.UserID,
		BadgeID: req.BadgeID,
	}

	if err := database.DB.Create(&userBadge).Error; err != nil {
		log.Printf("award badge: failed to create user badge, userID: %d, badgeID: %d, error: %v", req.UserID, req.BadgeID, err)
		utils.Error(w, 500, "授予失败")
		return
	}

	log.Printf("award badge: badge awarded successfully, userID: %d, badgeID: %d", req.UserID, req.BadgeID)
	utils.Success(w, userBadge)
}

// DeleteBadge 删除勋章处理器（管理员）
func DeleteBadge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 检查勋章是否存在
	var badge models.Badge
	if err := database.DB.First(&badge, id).Error; err != nil {
		log.Printf("delete badge: badge not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "勋章不存在")
		return
	}

	// 物理删除勋章
	if err := database.DB.Unscoped().Delete(&badge).Error; err != nil {
		log.Printf("delete badge: failed to delete badge, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	log.Printf("delete badge: badge deleted successfully, id: %d", id)
	utils.Success(w, nil)
}
