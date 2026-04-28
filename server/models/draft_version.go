package models

import "time"

// DraftVersion 草稿版本模型
// 存储草稿的历史版本记录
type DraftVersion struct {
	ID        uint      `gorm:"primarykey" json:"id"`          // 版本唯一标识
	DraftID   uint      `gorm:"not null;index" json:"draft_id"` // 所属草稿ID
	UserID    uint      `gorm:"not null;index" json:"user_id"`  // 所属用户ID
	Title     string    `gorm:"size:200" json:"title"`         // 版本标题
	Content   string    `gorm:"type:text" json:"content"`      // 版本正文内容
	ForumID   uint      `json:"forum_id"`                      // 目标版块ID
	Tags      IntSlice  `gorm:"type:json" json:"tags"`         // 标签ID列表
	Version   int       `gorm:"not null;default:1" json:"version"` // 版本号
	CreatedAt time.Time `json:"created_at"`                    // 创建时间
}

// DraftVersionWithPreview 带预览信息的草稿版本
type DraftVersionWithPreview struct {
	DraftVersion
	Preview string `json:"preview"` // 内容预览（前100字符）
}
