package models

import (
	"time"
)

// UserBadge 用户勋章关联模型
// 记录用户获得的勋章
type UserBadge struct {
	ID        uint       `gorm:"primarykey" json:"id"`           // 记录唯一标识
	UserID    uint       `gorm:"not null;index" json:"user_id"`  // 用户ID
	BadgeID   uint       `gorm:"not null;index" json:"badge_id"` // 勋章ID
	Badge     Badge      `gorm:"foreignKey:BadgeID" json:"badge"` // 勋章信息
	AwardedAt time.Time  `json:"awarded_at"`                     // 获得时间

	// 撤销相关
	IsRevoked     bool       `gorm:"default:false" json:"is_revoked"`       // 是否已撤销
	RevokedAt     *time.Time `json:"revoked_at"`                           // 撤销时间
	RevokedReason string     `gorm:"size:255" json:"revoked_reason"`        // 撤销原因
	RevokedBy     uint       `json:"revoked_by"`                            // 撤销操作人ID
}

// TableName 指定表名为 user_badges（默认情况下 GORM 会将 UserBadge 转换为 user_badges）
func (UserBadge) TableName() string {
	return "user_badges"
}
