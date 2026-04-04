package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
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
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 验证举报原因
	if report.Reason == "" {
		log.Printf("create report: reason is empty, reporterID: %d", userID)
		errors.Error(w, errors.CodeIncompleteInfo, "")
		return
	}

	// 设置举报者ID和状态
	report.ReporterID = userID
	report.Status = 0 // 待处理

	// 创建举报记录
	if err := database.DB.Create(&report).Error; err != nil {
		log.Printf("create report: failed to create report, reporterID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("create report: report created successfully, id: %d, reporterID: %d", report.ID, userID)
	errors.Success(w, report)
}

// GetUserReports 获取当前用户发起的举报列表处理器
func GetUserReports(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var reports []models.Report
	if err := database.DB.Where("reporter_id = ?", userID).Order("created_at DESC").Find(&reports).Error; err != nil {
		log.Printf("get user reports: failed to query reports, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, reports)
}
