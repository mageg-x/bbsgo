package services

import (
	"bbsgo/database"
	"bbsgo/models"
	"encoding/json"
	"log"
	"time"
)

// BadgeService 勋章服务
type BadgeService struct{}

// NewBadgeService 创建勋章服务实例
func NewBadgeService() *BadgeService {
	return &BadgeService{}
}

// CheckAndAwardBadges 检查并授予用户勋章
// 在用户行为后调用，检查是否满足勋章条件
func (s *BadgeService) CheckAndAwardBadges(userID uint) {
	var badges []models.Badge
	if err := database.DB.Where("1=1").Order("sort_order ASC, id ASC").Find(&badges).Error; err != nil {
		log.Printf("check badges: failed to query badges, userID: %d, error: %v", userID, err)
		return
	}

	for _, badge := range badges {
		s.checkAndAwardBadge(userID, badge)
	}
}

// checkAndAwardBadge 检查单个勋章条件并授予
func (s *BadgeService) checkAndAwardBadge(userID uint, badge models.Badge) {
	// 检查用户是否已获得该勋章（且未被撤销）
	var existingBadge models.UserBadge
	err := database.DB.Where("user_id = ? AND badge_id = ? AND is_revoked = ?", userID, badge.ID, false).
		First(&existingBadge).Error
	if err == nil {
		return
	}

	// 检查是否满足条件
	if s.checkCondition(userID, badge) {
		s.awardBadge(userID, badge.ID)
	}
}

// checkCondition 检查用户是否满足勋章条件
func (s *BadgeService) checkCondition(userID uint, badge models.Badge) bool {
	switch badge.ConditionType {
	case "register":
		return s.checkRegister(userID)
	case "topic_count":
		return s.checkTopicCount(userID, badge.ConditionValue)
	case "comment_count":
		return s.checkCommentCount(userID, badge.ConditionValue)
	case "like_count":
		return s.checkLikeCount(userID, badge.ConditionValue)
	case "best_comment":
		return s.checkBestComment(userID, badge.ConditionValue)
	case "follower_count":
		return s.checkFollowerCount(userID, badge.ConditionValue)
	case "combination":
		return s.checkCombination(userID, badge.ConditionData)
	default:
		return false
	}
}

// checkRegister 检查注册条件（注册成功即获得）
func (s *BadgeService) checkRegister(userID uint) bool {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return false
	}
	return user.ID > 0
}

// checkTopicCount 检查发帖数
func (s *BadgeService) checkTopicCount(userID uint, count int) bool {
	var topicCount int64
	database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&topicCount)
	return topicCount >= int64(count)
}

// checkCommentCount 检查评论数
func (s *BadgeService) checkCommentCount(userID uint, count int) bool {
	var commentCount int64
	database.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentCount)
	return commentCount >= int64(count)
}

// checkLikeCount 检查获赞数
func (s *BadgeService) checkLikeCount(userID uint, count int) bool {
	var likeCount int64
	database.DB.Model(&models.Like{}).
		Joins("JOIN topics ON topics.id = likes.target_id AND likes.target_type = 'topic'").
		Where("topics.user_id = ?", userID).
		Count(&likeCount)

	var commentLikeCount int64
	database.DB.Model(&models.Like{}).
		Joins("JOIN comments ON comments.id = likes.target_id AND likes.target_type = 'comment'").
		Where("comments.user_id = ?", userID).
		Count(&commentLikeCount)

	return (likeCount + commentLikeCount) >= int64(count)
}

// checkBestComment 检查最佳评论数
func (s *BadgeService) checkBestComment(userID uint, count int) bool {
	var bestCommentCount int64
	database.DB.Model(&models.Comment{}).
		Where("user_id = ? AND is_best = ?", userID, true).
		Count(&bestCommentCount)
	return bestCommentCount >= int64(count)
}

// checkFollowerCount 检查粉丝数
func (s *BadgeService) checkFollowerCount(userID uint, count int) bool {
	var followerCount int64
	database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", userID).Count(&followerCount)
	return followerCount >= int64(count)
}

// checkCombination 检查组合条件
func (s *BadgeService) checkCombination(userID uint, conditionData string) bool {
	var conditions map[string]int
	if err := json.Unmarshal([]byte(conditionData), &conditions); err != nil {
		log.Printf("check combination: failed to parse condition data, error: %v", err)
		return false
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return false
	}

	// 检查注册天数
	if days, ok := conditions["register_days"]; ok {
		registerDays := int(time.Since(user.CreatedAt).Hours() / 24)
		if registerDays < days {
			return false
		}
	}

	// 检查发帖数
	if count, ok := conditions["topic_count"]; ok {
		if !s.checkTopicCount(userID, count) {
			return false
		}
	}

	// 检查获赞数
	if count, ok := conditions["like_count"]; ok {
		if !s.checkLikeCount(userID, count) {
			return false
		}
	}

	// 检查最佳评论数
	if count, ok := conditions["best_comment"]; ok {
		if !s.checkBestComment(userID, count) {
			return false
		}
	}

	return true
}

// awardBadge 授予用户勋章
func (s *BadgeService) awardBadge(userID uint, badgeID uint) error {
	userBadge := models.UserBadge{
		UserID:    userID,
		BadgeID:   badgeID,
		AwardedAt: time.Now(),
	}

	if err := database.DB.Create(&userBadge).Error; err != nil {
		log.Printf("award badge: failed to create user badge, userID: %d, badgeID: %d, error: %v", userID, badgeID, err)
		return err
	}

	log.Printf("award badge: badge awarded successfully, userID: %d, badgeID: %d", userID, badgeID)

	// 发送通知
	s.SendBadgeNotification(userID, badgeID)

	return nil
}

// SendBadgeNotification 发送勋章获得通知
func (s *BadgeService) SendBadgeNotification(userID uint, badgeID uint) {
	var badge models.Badge
	if err := database.DB.First(&badge, badgeID).Error; err != nil {
		log.Printf("send badge notification: badge not found, badgeID: %d", badgeID)
		return
	}

	notification := models.Notification{
		UserID:      userID,
		Type:        "badge",
		Content:     "notifications.badge_earned",
		RelatedID:   badgeID,
		RelatedType: "badge",
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("send badge notification: failed to create notification, userID: %d, error: %v", userID, err)
	}
}

// GetUserBadgeProgress 获取用户勋章进度
func (s *BadgeService) GetUserBadgeProgress(userID uint) ([]map[string]interface{}, error) {
	var badges []models.Badge
	if err := database.DB.Order("sort_order ASC, id ASC").Find(&badges).Error; err != nil {
		return nil, err
	}

	var userBadges []models.UserBadge
	database.DB.Where("user_id = ? AND is_revoked = ?", userID, false).Find(&userBadges)

	userBadgeMap := make(map[uint]models.UserBadge)
	for _, ub := range userBadges {
		userBadgeMap[ub.BadgeID] = ub
	}

	var result []map[string]interface{}
	for _, badge := range badges {
		progress := s.calculateProgress(userID, badge)

		item := map[string]interface{}{
			"badge":      badge,
			"awarded":    false,
			"awarded_at": nil,
			"progress":   progress,
		}

		if ub, ok := userBadgeMap[badge.ID]; ok {
			item["awarded"] = true
			item["awarded_at"] = ub.AwardedAt
		}

		result = append(result, item)
	}

	return result, nil
}

// calculateProgress 计算勋章进度
func (s *BadgeService) calculateProgress(userID uint, badge models.Badge) map[string]interface{} {
	progress := make(map[string]interface{})

	switch badge.ConditionType {
	case "register":
		progress["current"] = 1
		progress["target"] = 1
	case "topic_count":
		var count int64
		database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&count)
		progress["current"] = count
		progress["target"] = badge.ConditionValue
	case "comment_count":
		var count int64
		database.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&count)
		progress["current"] = count
		progress["target"] = badge.ConditionValue
	case "like_count":
		var likeCount int64
		database.DB.Model(&models.Like{}).
			Joins("JOIN topics ON topics.id = likes.target_id AND likes.target_type = 'topic'").
			Where("topics.user_id = ?", userID).
			Count(&likeCount)

		var commentLikeCount int64
		database.DB.Model(&models.Like{}).
			Joins("JOIN comments ON comments.id = likes.target_id AND likes.target_type = 'comment'").
			Where("comments.user_id = ?", userID).
			Count(&commentLikeCount)

		progress["current"] = likeCount + commentLikeCount
		progress["target"] = badge.ConditionValue
	case "best_comment":
		var count int64
		database.DB.Model(&models.Comment{}).
			Where("user_id = ? AND is_best = ?", userID, true).
			Count(&count)
		progress["current"] = count
		progress["target"] = badge.ConditionValue
	case "follower_count":
		var count int64
		database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", userID).Count(&count)
		progress["current"] = count
		progress["target"] = badge.ConditionValue
	case "combination":
		var conditions map[string]int
		if err := json.Unmarshal([]byte(badge.ConditionData), &conditions); err == nil {
			details := make(map[string]interface{})

			if days, ok := conditions["register_days"]; ok {
				var user models.User
				database.DB.First(&user, userID)
				registerDays := int(time.Since(user.CreatedAt).Hours() / 24)
				details["register_days"] = map[string]int{
					"current": registerDays,
					"target":  days,
				}
			}

			if count, ok := conditions["topic_count"]; ok {
				var topicCount int64
				database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&topicCount)
				details["topic_count"] = map[string]int64{
					"current": topicCount,
					"target":  int64(count),
				}
			}

			if count, ok := conditions["like_count"]; ok {
				var likeCount int64
				database.DB.Model(&models.Like{}).
					Joins("JOIN topics ON topics.id = likes.target_id AND likes.target_type = 'topic'").
					Where("topics.user_id = ?", userID).
					Count(&likeCount)

				var commentLikeCount int64
				database.DB.Model(&models.Like{}).
					Joins("JOIN comments ON comments.id = likes.target_id AND likes.target_type = 'comment'").
					Where("comments.user_id = ?", userID).
					Count(&commentLikeCount)

				details["like_count"] = map[string]int64{
					"current": likeCount + commentLikeCount,
					"target":  int64(count),
				}
			}

			if count, ok := conditions["best_comment"]; ok {
				var bestCommentCount int64
				database.DB.Model(&models.Comment{}).
					Where("user_id = ? AND is_best = ?", userID, true).
					Count(&bestCommentCount)
				details["best_comment"] = map[string]int64{
					"current": bestCommentCount,
					"target":  int64(count),
				}
			}

			progress["details"] = details
		}
	}

	return progress
}
