package models

import "time"

// Like 点赞模型
// 记录用户对话题或帖子的点赞行为
type Like struct {
	ID         uint      `gorm:"primarykey" json:"id"`                      // 点赞唯一标识
	UserID     uint      `gorm:"not null;index" json:"user_id"`             // 点赞用户ID
	User       User      `gorm:"foreignKey:UserID" json:"-"`                // 点赞用户信息
	TargetType string    `gorm:"size:20;not null;index" json:"target_type"` // 被点赞对象类型：topic=话题, post=帖子
	TargetID   uint      `gorm:"not null;index" json:"target_id"`           // 被点赞对象ID
	CreatedAt  time.Time `json:"created_at"`                                // 点赞时间
}
