package models

import "time"

// Draft 草稿箱模型
// 存储用户未发布的话题草稿
type Draft struct {
	ID        uint      `gorm:"primarykey" json:"id"`          // 草稿唯一标识
	UserID    uint      `gorm:"not null;index" json:"user_id"` // 所属用户ID
	User      User      `gorm:"foreignKey:UserID" json:"-"`    // 所属用户信息
	Title     string    `gorm:"size:200" json:"title"`         // 草稿标题
	Content   string    `gorm:"type:text" json:"content"`      // 草稿正文内容
	ForumID   uint      `json:"forum_id"`                      // 目标版块ID
	Tags      IntSlice  `gorm:"type:json" json:"tags"`         // 标签ID列表
	CreatedAt time.Time `json:"created_at"`                    // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                    // 更新时间
}
