package models

import "time"

// Announcement 公告模型
// 存储网站公告信息
type Announcement struct {
	ID        uint       `gorm:"primarykey" json:"id"`              // 公告唯一标识
	Title     string     `gorm:"size:200;not null" json:"title"`    // 公告标题
	Content   string     `gorm:"type:text;not null" json:"content"` // 公告正文内容
	IsPinned  bool       `gorm:"default:false" json:"is_pinned"`    // 是否置顶显示
	ExpiresAt *time.Time `json:"expires_at"`                        // 过期时间，nil表示永不过期
	CreatedAt time.Time  `json:"created_at"`                        // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                        // 更新时间
}
