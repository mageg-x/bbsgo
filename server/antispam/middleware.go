package antispam

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// AntiSpamMiddleware 防刷中间件
// 整合频率限制、内容质量、信誉分、举报处理等多个服务
type AntiSpamMiddleware struct {
	rateLimit  *RateLimitService
	quality    *ContentQualityService
	reputation *ReputationService
	report     *ReportService
	config     *ConfigService
}

var antiSpamMiddleware *AntiSpamMiddleware

// NewAntiSpamMiddleware 创建新的防刷中间件实例
func NewAntiSpamMiddleware() *AntiSpamMiddleware {
	return &AntiSpamMiddleware{
		rateLimit:  NewRateLimitService(),
		quality:    NewContentQualityService(),
		reputation: NewReputationService(),
		report:     NewReportService(),
		config:     GetConfigService(),
	}
}

// GetAntiSpamMiddleware 获取全局防刷中间件单例
func GetAntiSpamMiddleware() *AntiSpamMiddleware {
	if antiSpamMiddleware == nil {
		log.Printf("[antispam] creating middleware singleton")
		antiSpamMiddleware = NewAntiSpamMiddleware()
	}
	return antiSpamMiddleware
}

// AntiSpamCheckResult 防刷检查结果
type AntiSpamCheckResult struct {
	Allowed      bool   `json:"allowed"`
	Reason       string `json:"reason"`
	RetryAfter   int    `json:"retry_after,omitempty"`
	NeedsCaptcha bool   `json:"needs_captcha,omitempty"`
}

// CheckTopicCreate 检查创建话题
func (m *AntiSpamMiddleware) CheckTopicCreate(userID uint, content string) *AntiSpamCheckResult {
	log.Printf("[antispam] checking topic create, userID: %d, content length: %d", userID, len(content))

	// 1. 禁言状态检查
	isBanned, banReason := m.rateLimit.CheckBanStatus(userID)
	if isBanned {
		log.Printf("[antispam] topic rejected - user banned, userID: %d, reason: %s", userID, banReason)
		return &AntiSpamCheckResult{
			Allowed: false,
			Reason:  banReason,
		}
	}

	// 2. 频率限制检查
	rateResult := m.rateLimit.CheckTopicRateLimit(userID, time.Now())
	if !rateResult.Allowed {
		log.Printf("[antispam] topic rejected - rate limit, userID: %d, reason: %s, retryAfter: %d",
			userID, rateResult.Reason, rateResult.RetryAfter)
		return &AntiSpamCheckResult{
			Allowed:    false,
			Reason:     rateResult.Reason,
			RetryAfter: rateResult.RetryAfter,
		}
	}

	// 3. 内容质量检查
	if err := m.quality.ValidateContent(content); err != nil {
		log.Printf("[antispam] topic rejected - quality failed, userID: %d, reason: %s, will deduct rep",
			userID, err.Error())
		m.reputation.ChangeReputation(userID, -1, "发布低质量内容", 0)
		return &AntiSpamCheckResult{
			Allowed: false,
			Reason:  err.Error(),
		}
	}

	// 4. 重复内容检查
	if err := m.rateLimit.CheckDuplicate(userID, content, "topic", 1); err != nil {
		log.Printf("[antispam] topic rejected - duplicate, userID: %d, reason: %s", userID, err.Error())
		return &AntiSpamCheckResult{
			Allowed: false,
			Reason:  err.Error(),
		}
	}

	// 5. 验证码检查
	needsCaptcha := m.reputation.NeedsCaptcha(userID)
	if needsCaptcha {
		log.Printf("[antispam] topic needs captcha, userID: %d", userID)
	}

	log.Printf("[antispam] topic check passed, userID: %d, needsCaptcha: %v", userID, needsCaptcha)
	return &AntiSpamCheckResult{
		Allowed:      true,
		NeedsCaptcha: needsCaptcha,
	}
}

// CheckCommentCreate 检查创建评论
func (m *AntiSpamMiddleware) CheckCommentCreate(userID uint, content string) *AntiSpamCheckResult {
	log.Printf("[antispam] checking comment create, userID: %d, content length: %d", userID, len(content))

	// 1. 禁言状态检查
	isBanned, banReason := m.rateLimit.CheckBanStatus(userID)
	if isBanned {
		log.Printf("[antispam] comment rejected - user banned, userID: %d, reason: %s", userID, banReason)
		return &AntiSpamCheckResult{
			Allowed: false,
			Reason:  banReason,
		}
	}

	// 2. 频率限制检查
	rateResult := m.rateLimit.CheckCommentRateLimit(userID, time.Now())
	if !rateResult.Allowed {
		log.Printf("[antispam] comment rejected - rate limit, userID: %d, reason: %s, retryAfter: %d",
			userID, rateResult.Reason, rateResult.RetryAfter)
		return &AntiSpamCheckResult{
			Allowed:    false,
			Reason:     rateResult.Reason,
			RetryAfter: rateResult.RetryAfter,
		}
	}

	// 3. 内容质量检查
	if err := m.quality.ValidateContent(content); err != nil {
		log.Printf("[antispam] comment rejected - quality failed, userID: %d, reason: %s, will deduct rep",
			userID, err.Error())
		m.reputation.ChangeReputation(userID, -1, "发布低质量内容", 0)
		return &AntiSpamCheckResult{
			Allowed: false,
			Reason:  err.Error(),
		}
	}

	// 4. 重复内容检查
	if err := m.rateLimit.CheckDuplicate(userID, content, "comment", 1); err != nil {
		log.Printf("[antispam] comment rejected - duplicate, userID: %d, reason: %s", userID, err.Error())
		return &AntiSpamCheckResult{
			Allowed: false,
			Reason:  err.Error(),
		}
	}

	// 5. 验证码检查
	needsCaptcha := m.reputation.NeedsCaptcha(userID)
	if needsCaptcha {
		log.Printf("[antispam] comment needs captcha, userID: %d", userID)
	}

	log.Printf("[antispam] comment check passed, userID: %d, needsCaptcha: %v", userID, needsCaptcha)
	return &AntiSpamCheckResult{
		Allowed:      true,
		NeedsCaptcha: needsCaptcha,
	}
}

// RecordTopicCreation 记录话题发布
func (m *AntiSpamMiddleware) RecordTopicCreation(userID uint, topicID uint, content string) {
	log.Printf("[antispam] recording topic, userID: %d, topicID: %d", userID, topicID)
	m.rateLimit.RecordOperation(userID, "topic", topicID, "topic", content)

	qualityResult := m.quality.CheckContent(content)
	m.quality.RecordQuality(topicID, "topic", qualityResult)

	log.Printf("[antispam] topic quality recorded, topicID: %d, score: %.2f, lowQuality: %v, reason: %s",
		topicID, qualityResult.Score, qualityResult.IsLowQuality, qualityResult.Reason)
}

// RecordCommentCreation 记录评论发布
func (m *AntiSpamMiddleware) RecordCommentCreation(userID uint, commentID uint, content string) {
	log.Printf("[antispam] recording comment, userID: %d, commentID: %d", userID, commentID)
	m.rateLimit.RecordOperation(userID, "comment", commentID, "comment", content)

	qualityResult := m.quality.CheckContent(content)
	m.quality.RecordQuality(commentID, "comment", qualityResult)

	log.Printf("[antispam] comment quality recorded, commentID: %d, score: %.2f, lowQuality: %v, reason: %s",
		commentID, qualityResult.Score, qualityResult.IsLowQuality, qualityResult.Reason)
}

// RecordUserOperation 记录用户操作
func (m *AntiSpamMiddleware) RecordUserOperation(userID uint, operation string, targetID uint, content string) {
	log.Printf("[antispam] recording operation, userID: %d, operation: %s, targetID: %d", userID, operation, targetID)
	m.rateLimit.RecordOperation(userID, operation, targetID, operation, content)

	qualityResult := m.quality.CheckContent(content)
	m.quality.RecordQuality(targetID, operation, qualityResult)
}

// HandleReport 处理用户举报
func (m *AntiSpamMiddleware) HandleReport(reporterID uint, targetType string, targetID uint, reason string) error {
	log.Printf("[antispam] handling report, reporterID: %d, targetType: %s, targetID: %d, reason: %s",
		reporterID, targetType, targetID, reason)
	return m.report.CreateReport(reporterID, targetType, targetID, reason)
}

// InitializeAntiSpamSystem 初始化防刷系统
func InitializeAntiSpamSystem() {
	log.Printf("[antispam] initializing system")
	config := GetConfigService()
	config.InitializeDefaults()
	log.Printf("[antispam] initialization completed")
}

// RegisterAntiSpamRoutes 注册防刷相关API路由
func RegisterAntiSpamRoutes(r *mux.Router) {
	log.Printf("[antispam] registering routes")
	middleware := NewAntiSpamMiddleware()

	r.HandleFunc("/api/v1/antispam/check/topic", func(w http.ResponseWriter, req *http.Request) {
		var reqBody struct {
			UserID  uint   `json:"user_id"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			log.Printf("[antispam] parse topic check request failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("[antispam] topic check request, userID: %d", reqBody.UserID)
		result := middleware.CheckTopicCreate(reqBody.UserID, reqBody.Content)
		json.NewEncoder(w).Encode(result)
	}).Methods("POST")

	r.HandleFunc("/api/v1/antispam/check/comment", func(w http.ResponseWriter, req *http.Request) {
		var reqBody struct {
			UserID  uint   `json:"user_id"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			log.Printf("[antispam] parse comment check request failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("[antispam] comment check request, userID: %d", reqBody.UserID)
		result := middleware.CheckCommentCreate(reqBody.UserID, reqBody.Content)
		json.NewEncoder(w).Encode(result)
	}).Methods("POST")

	r.HandleFunc("/api/v1/antispam/report", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ReporterID uint   `json:"reporter_id"`
			TargetType string `json:"target_type"`
			TargetID   uint   `json:"target_id"`
			Reason     string `json:"reason"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[antispam] parse report request failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := middleware.HandleReport(req.ReporterID, req.TargetType, req.TargetID, req.Reason); err != nil {
			log.Printf("[antispam] handle report failed, reporterID: %d, targetType: %s, targetID: %d, error: %v",
				req.ReporterID, req.TargetType, req.TargetID, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("[antispam] report handled, reporterID: %d, targetType: %s, targetID: %d",
			req.ReporterID, req.TargetType, req.TargetID)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "举报成功",
		})
	}).Methods("POST")

	log.Printf("[antispam] routes registered")
}

// StartScheduledTasks 启动定时任务
func StartScheduledTasks() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("[antispam-scheduler] running hourly tasks")

			rateLimit := NewRateLimitService()
			if err := rateLimit.CleanupOldRecords(); err != nil {
				log.Printf("[antispam-scheduler] cleanup failed: %v", err)
			} else {
				log.Printf("[antispam-scheduler] cleanup done")
			}

			reputation := NewReputationService()
			if err := reputation.RecoverFromBan(0); err != nil {
				log.Printf("[antispam-scheduler] ban recover check failed: %v", err)
			} else {
				log.Printf("[antispam-scheduler] ban recover check done")
			}

			log.Printf("[antispam-scheduler] hourly tasks done")
		}
	}()

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("[antispam-scheduler] running daily tasks")

			reputation := NewReputationService()
			if err := reputation.AwardDailyRecovery(); err != nil {
				log.Printf("[antispam-scheduler] daily recovery failed: %v", err)
			} else {
				log.Printf("[antispam-scheduler] daily recovery done")
			}

			hotScore := NewHotScoreService()
			if err := hotScore.RecalculateAllScores(); err != nil {
				log.Printf("[antispam-scheduler] recalculate hot scores failed: %v", err)
			} else {
				log.Printf("[antispam-scheduler] recalculate hot scores done")
			}

			log.Printf("[antispam-scheduler] daily tasks done")
		}
	}()

	log.Printf("[antispam] scheduled tasks started")
}
