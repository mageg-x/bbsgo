package antispam

import (
	"bbsgo/database"
	"bbsgo/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
)

// RateLimitService 频率限制服务
// 负责检查用户的发帖/评论频率限制
type RateLimitService struct {
	config *ConfigService
}

// NewRateLimitService 创建频率限制服务实例
func NewRateLimitService() *RateLimitService {
	return &RateLimitService{
		config: GetConfigService(),
	}
}

// RateLimitResult 频率限制检查结果
type RateLimitResult struct {
	Allowed    bool
	Reason     string
	RetryAfter int
}

// CheckTopicRateLimit 检查话题发布的频率限制
func (s *RateLimitService) CheckTopicRateLimit(userID uint, createdAt time.Time) *RateLimitResult {
	minInterval := s.config.GetInt(ConfigTopicMinInterval, 60)
	maxPerDay := s.config.GetInt(ConfigMaxTopicsPerDay, 10)
	newUserMaxPerDay := s.config.GetInt(ConfigNewUserMaxTopicsPerDay, 3)
	newUserHours := s.config.GetInt(ConfigNewUserHours, 24)

	log.Printf("[ratelimit] checking topic, userID: %d, minInterval: %ds, maxPerDay: %d, newUserMax: %d, newUserHours: %d",
		userID, minInterval, maxPerDay, newUserMaxPerDay, newUserHours)

	return s.checkRateLimit(userID, "topic", minInterval, maxPerDay, newUserMaxPerDay, newUserHours, createdAt)
}

// CheckCommentRateLimit 检查评论发布的频率限制
func (s *RateLimitService) CheckCommentRateLimit(userID uint, createdAt time.Time) *RateLimitResult {
	minInterval := s.config.GetInt(ConfigCommentMinInterval, 30)
	maxPerDay := s.config.GetInt(ConfigMaxCommentsPerDay, 50)
	newUserMaxPerDay := s.config.GetInt(ConfigNewUserMaxCommentsPerDay, 10)
	newUserHours := s.config.GetInt(ConfigNewUserHours, 24)

	log.Printf("[ratelimit] checking comment, userID: %d, minInterval: %ds, maxPerDay: %d, newUserMax: %d, newUserHours: %d",
		userID, minInterval, maxPerDay, newUserMaxPerDay, newUserHours)

	return s.checkRateLimit(userID, "comment", minInterval, maxPerDay, newUserMaxPerDay, newUserHours, createdAt)
}

// checkRateLimit 通用频率限制检查
func (s *RateLimitService) checkRateLimit(userID uint, operation string, minInterval, maxPerDay, newUserMaxPerDay, newUserHours int, _ time.Time) *RateLimitResult {
	// 1. 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("[ratelimit] user not found, userID: %d, error: %v", userID, err)
		return &RateLimitResult{Allowed: false, Reason: "用户不存在"}
	}

	// 2. 根据信誉分调整限制（信誉分<40时，日上限降为3）
	if user.Reputation < 40 {
		maxPerDay = min(maxPerDay, 3)
		log.Printf("[ratelimit] low reputation, userID: %d, rep: %d, adjusted maxPerDay: %d",
			userID, user.Reputation, maxPerDay)
	}

	// 3. 检查是否是新用户
	isNewUser := time.Since(user.CreatedAt) < time.Duration(newUserHours)*time.Hour
	if isNewUser {
		maxPerDay = min(maxPerDay, newUserMaxPerDay)
		log.Printf("[ratelimit] new user, userID: %d, createdAt: %v, using new user limit: %d",
			userID, user.CreatedAt, maxPerDay)
	}

	// 4. 检查最小间隔（如果 minInterval > 0 才检查）
	if minInterval > 0 {
		var lastOp models.UserOperation
		result := database.DB.Where("user_id = ? AND operation = ?", userID, operation).
			Order("created_at DESC").
			First(&lastOp)

		if result.Error == nil {
			elapsed := time.Since(lastOp.CreatedAt)
			if elapsed < time.Duration(minInterval)*time.Second {
				retryAfter := minInterval - int(elapsed.Seconds())
				log.Printf("[ratelimit] interval too short, userID: %d, op: %s, elapsed: %.0fs, minInterval: %ds, retryAfter: %ds",
					userID, operation, elapsed.Seconds(), minInterval, retryAfter)
				return &RateLimitResult{
					Allowed:    false,
					Reason:     "操作过快，请稍后再试",
					RetryAfter: retryAfter,
				}
			}
		}
	}

	// 5. 检查当日次数限制
	today := time.Now().Format("2006-01-02")
	var count int64
	if err := database.DB.Model(&models.UserOperation{}).
		Where("user_id = ? AND operation = ? AND DATE(created_at) = ?", userID, operation, today).
		Count(&count).Error; err != nil {
		log.Printf("[ratelimit] query daily count failed, userID: %d, op: %s, error: %v",
			userID, operation, err)
	}

	if int(count) >= maxPerDay {
		log.Printf("[ratelimit] daily limit reached, userID: %d, op: %s, count: %d, maxPerDay: %d",
			userID, operation, count, maxPerDay)
		return &RateLimitResult{
			Allowed: false,
			Reason:  "今日操作次数已达上限",
		}
	}

	log.Printf("[ratelimit] passed, userID: %d, op: %s, todayCount: %d, remaining: %d",
		userID, operation, count, maxPerDay-int(count))
	return &RateLimitResult{Allowed: true}
}

// RecordOperation 记录用户操作
func (s *RateLimitService) RecordOperation(userID uint, operation string, targetID uint, targetType string, content string) error {
	contentHash := ""
	if content != "" {
		hash := sha256.Sum256([]byte(content))
		contentHash = hex.EncodeToString(hash[:])
	}

	op := models.UserOperation{
		UserID:      userID,
		Operation:   operation,
		TargetID:    targetID,
		TargetType:  targetType,
		ContentHash: contentHash,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&op).Error; err != nil {
		log.Printf("[ratelimit] record failed, userID: %d, op: %s, targetID: %d, error: %v",
			userID, operation, targetID, err)
		return err
	}

	log.Printf("[ratelimit] recorded, userID: %d, op: %s, targetID: %d, targetType: %s",
		userID, operation, targetID, targetType)
	return nil
}

// CheckDuplicate 检查重复内容（1小时内相同内容哈希不能发布）
func (s *RateLimitService) CheckDuplicate(userID uint, content string, operation string, hours int) error {
	if content == "" {
		return nil
	}

	hash := sha256.Sum256([]byte(content))
	contentHash := hex.EncodeToString(hash[:])

	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	var recentOps []models.UserOperation
	if err := database.DB.Where("user_id = ? AND operation = ? AND created_at > ?", userID, operation, since).
		Find(&recentOps).Error; err != nil {
		log.Printf("[ratelimit] query duplicate failed, userID: %d, op: %s, error: %v", userID, operation, err)
		return nil
	}

	for _, op := range recentOps {
		if op.ContentHash == contentHash {
			log.Printf("[ratelimit] duplicate detected, userID: %d, op: %s, hours: %d, existingID: %d",
				userID, operation, hours, op.ID)
			return errors.New("请勿重复发布相同内容")
		}
	}

	return nil
}

// GetDailyCount 获取用户当日指定操作的次数
func (s *RateLimitService) GetDailyCount(userID uint, operation string) (int, error) {
	today := time.Now().Format("2006-01-02")
	var count int64
	err := database.DB.Model(&models.UserOperation{}).
		Where("user_id = ? AND operation = ? AND DATE(created_at) = ?", userID, operation, today).
		Count(&count).Error
	if err != nil {
		log.Printf("[ratelimit] get daily count failed, userID: %d, op: %s, error: %v", userID, operation, err)
		return 0, err
	}
	return int(count), nil
}

// CleanupOldRecords 清理7天前的操作记录
func (s *RateLimitService) CleanupOldRecords() error {
	threshold := time.Now().AddDate(0, 0, -7)
	result := database.DB.Where("created_at < ?", threshold).Delete(&models.UserOperation{})
	if result.Error != nil {
		log.Printf("[ratelimit] cleanup failed, threshold: %v, error: %v", threshold, result.Error)
		return result.Error
	}
	log.Printf("[ratelimit] cleanup done, deleted: %d, threshold: %v", result.RowsAffected, threshold)
	return nil
}

// CheckBanStatus 检查用户是否被禁言
func (s *RateLimitService) CheckBanStatus(userID uint) (bool, string) {
	var ban models.UserBan
	result := database.DB.Where("user_id = ? AND is_active = ?", userID, true).
		Where("end_time IS NULL OR end_time > ?", time.Now()).
		Order("created_at DESC").
		First(&ban)

	if result.Error != nil {
		log.Printf("[ratelimit] ban check done - not banned, userID: %d", userID)
		return false, ""
	}

	if ban.EndTime == nil {
		log.Printf("[ratelimit] permanently banned, userID: %d, reason: %s, start: %v",
			userID, ban.Reason, ban.CreatedAt)
		return true, "您已被永久禁言，原因：" + ban.Reason
	}

	remaining := time.Until(*ban.EndTime)
	if remaining > 0 {
		log.Printf("[ratelimit] temporarily banned, userID: %d, reason: %s, remaining: %v, end: %v",
			userID, ban.Reason, remaining, *ban.EndTime)
		return true, "您已被禁言，剩余时间：" + formatDuration(remaining) + "，原因：" + ban.Reason
	}

	return false, ""
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if days > 0 {
		return fmt.Sprintf("%d天%d小时", days, hours)
	}
	return fmt.Sprintf("%d小时%d分钟", hours, int(d.Minutes())%60)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
