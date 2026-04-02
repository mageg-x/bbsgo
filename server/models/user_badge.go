package models

import "time"

// UserBadge 用户勋章关联模型
// 记录用户获得的勋章
type UserBadge struct {
	ID        uint      `gorm:"primarykey" json:"id"`           // 记录唯一标识
	UserID    uint      `gorm:"not null;index" json:"user_id"`  // 用户ID
	BadgeID   uint      `gorm:"not null;index" json:"badge_id"` // 勋章ID
	AwardedAt time.Time `json:"awarded_at"`                     // 获得时间
}

// TableName 指定表名为 user_badges（默认情况下 GORM 会将 UserBadge 转换为 user_badges）
func (UserBadge) TableName() string {
	return "user_badges"
}
