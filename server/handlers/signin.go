package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"time"
)

// SignIn 用户签到处理器
// 每日签到可获得积分，连续签到可获得额外积分
func SignIn(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("sign in: user not found, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 检查今日是否已签到
	today := time.Now().Format("2006-01-02")
	if user.LastSignAt != nil && user.LastSignAt.Format("2006-01-02") == today {
		log.Printf("sign in: already signed in today, userID: %d", userID)
		utils.Error(w, 400, "今日已签到")
		return
	}

	// 计算签到积分
	credits := 10

	var lastSign time.Time
	if user.LastSignAt != nil {
		lastSign = *user.LastSignAt
	}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	// 如果昨天也签到了，获得连续签到奖励
	if lastSign.Format("2006-01-02") == yesterday {
		credits = 15
	}

	// 更新用户签到信息
	now := time.Now()
	user.LastSignAt = &now
	user.Credits += credits

	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("sign in: failed to update user, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "签到失败")
		return
	}

	log.Printf("sign in: signed in successfully, userID: %d, credits: %d", userID, credits)
	utils.Success(w, map[string]interface{}{
		"credits":       credits,
		"total_credits": user.Credits,
		"message":       "签到成功",
	})
}

// GetSignInStatus 获取签到状态处理器
func GetSignInStatus(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get sign in status: user not found, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	today := time.Now().Format("2006-01-02")
	signedToday := user.LastSignAt != nil && user.LastSignAt.Format("2006-01-02") == today

	utils.Success(w, map[string]interface{}{
		"signed_today": signedToday,
		"last_sign_at": user.LastSignAt,
		"credits":      user.Credits,
	})
}
