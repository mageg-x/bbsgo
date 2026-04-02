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

// CreateReport 创建举报处理器
// 用户可以举报话题、帖子或其他用户
func CreateReport(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 解析请求体
	var report models.Report
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Printf("create report: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证举报原因
	if report.Reason == "" {
		log.Printf("create report: reason is empty, reporterID: %d", userID)
		utils.Error(w, 400, "请填写举报原因")
		return
	}

	// 设置举报者ID和状态
	report.ReporterID = userID
	report.Status = 0 // 待处理

	// 创建举报记录
	if err := database.DB.Create(&report).Error; err != nil {
		log.Printf("create report: failed to create report, reporterID: %d, error: %v", userID, err)
		utils.Error(w, 500, "举报失败")
		return
	}

	log.Printf("create report: report created successfully, id: %d, reporterID: %d", report.ID, userID)
	utils.Success(w, report)
}

// GetUserReports 获取当前用户发起的举报列表处理器
func GetUserReports(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var reports []models.Report
	if err := database.DB.Where("reporter_id = ?", userID).Order("created_at DESC").Find(&reports).Error; err != nil {
		log.Printf("get user reports: failed to query reports, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取举报列表失败")
		return
	}

	utils.Success(w, reports)
}
