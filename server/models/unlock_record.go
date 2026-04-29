package models

import "time"

// UnlockRecord 解锁记录模型
// 记录用户对话题的解锁状态，实现永久解锁
type UnlockRecord struct {
	ID        uint      `gorm:"primarykey" json:"id"`                      // 记录唯一标识
	UserID    uint      `gorm:"not null;index" json:"user_id"`             // 用户ID
	TopicID   uint      `gorm:"not null;index" json:"topic_id"`            // 话题ID
	Topic     Topic     `gorm:"foreignKey:TopicID" json:"-"`                // 话题信息
	User      User      `gorm:"foreignKey:UserID" json:"-"`                 // 用户信息
	UnlockType string    `gorm:"size:20;not null" json:"unlock_type"`     // 解锁类型：like=点赞解锁, comment=评论解锁, both=两者都满足
	CreatedAt time.Time `json:"created_at"`                                // 创建时间
}
