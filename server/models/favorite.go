package models

import "time"

// Favorite 收藏模型
// 记录用户对话题的收藏关系
type Favorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`            // 收藏唯一标识
	UserID    uint      `gorm:"not null;index" json:"user_id"`   // 收藏者用户ID
	User      User      `gorm:"foreignKey:UserID" json:"-"`      // 收藏者用户信息
	TopicID   uint      `gorm:"not null;index" json:"topic_id"`  // 被收藏的话题ID
	Topic     Topic     `gorm:"foreignKey:TopicID" json:"topic"` // 被收藏的话题信息
	CreatedAt time.Time `json:"created_at"`                      // 收藏时间
}
