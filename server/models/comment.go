package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`              // 评论唯一标识
	TopicID   uint           `gorm:"not null;index" json:"topic_id"`   // 所属话题ID
	Topic     Topic          `gorm:"foreignKey:TopicID" json:"topic"`   // 所属话题信息
	UserID    uint           `gorm:"not null;index" json:"user_id"`     // 发布者用户ID
	User      User           `gorm:"foreignKey:UserID" json:"user"`     // 发布者用户信息
	ReplyToID *uint          `gorm:"index" json:"reply_to_id"`         // 回复给哪个评论的ID
	ReplyTo   *Comment       `gorm:"foreignKey:ReplyToID" json:"-"`    // 被回复的评论
	Content   string         `gorm:"type:text;not null" json:"content"` // 评论正文内容
	LikeCount int            `gorm:"default:0" json:"like_count"`       // 点赞数
	IsPinned  bool           `gorm:"default:false" json:"is_pinned"`    // 是否置顶
	IsBest    bool           `gorm:"default:false" json:"is_best"`      // 是否最佳评论
	CreatedAt time.Time      `json:"created_at"`                        // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                        // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                    // 软删除时间

	// 关联关系
	Likes []Like `gorm:"foreignKey:TargetID" json:"-"` // 评论的点赞记录
}
