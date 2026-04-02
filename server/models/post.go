package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 帖子/回复模型
// 论坛中的回复帖子，隶属于某个话题，可以嵌套回复
type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`              // 回复唯一标识
	TopicID   uint           `gorm:"not null;index" json:"topic_id"`    // 所属话题ID
	Topic     Topic          `gorm:"foreignKey:TopicID" json:"topic"`   // 所属话题信息
	UserID    uint           `gorm:"not null;index" json:"user_id"`     // 发布者用户ID
	User      User           `gorm:"foreignKey:UserID" json:"user"`     // 发布者用户信息
	ParentID  *uint          `gorm:"index" json:"parent_id"`            // 父回复ID，用于嵌套评论
	Parent    *Post          `gorm:"foreignKey:ParentID" json:"-"`      // 父回复信息
	Content   string         `gorm:"type:text;not null" json:"content"` // 回复正文内容
	LikeCount int            `gorm:"default:0" json:"like_count"`       // 点赞数
	CreatedAt time.Time      `json:"created_at"`                        // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                        // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                    // 软删除时间

	// 关联关系
	Children []Post `gorm:"foreignKey:ParentID" json:"children"` // 子回复列表（嵌套评论）
	Likes    []Like `gorm:"foreignKey:TargetID" json:"-"`        // 回复的点赞记录
}
