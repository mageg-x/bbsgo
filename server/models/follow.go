package models

import "time"

// Follow 关注关系模型
// 记录用户之间的关注关系
type Follow struct {
	ID           uint      `gorm:"primarykey" json:"id"`                       // 关注关系唯一标识
	UserID       uint      `gorm:"not null;index" json:"user_id"`              // 关注者用户ID
	User         User      `gorm:"foreignKey:UserID" json:"user"`             // 关注者用户信息
	FollowUserID uint      `gorm:"not null;index" json:"follow_user_id"`       // 被关注者用户ID
	FollowUser   User      `gorm:"foreignKey:FollowUserID" json:"follow_user"` // 被关注者用户信息
	CreatedAt    time.Time `json:"created_at"`                                 // 关注时间
}
