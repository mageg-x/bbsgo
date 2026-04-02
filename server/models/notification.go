package models

import "time"

// Notification 通知模型
// 存储系统向用户发送的通知消息
type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`               // 通知唯一标识
	UserID    uint      `gorm:"not null;index" json:"user_id"`      // 接收通知的用户ID
	User      User      `gorm:"foreignKey:UserID" json:"-"`         // 接收通知的用户信息
	Type      string    `gorm:"size:50;not null" json:"type"`       // 通知类型：message=私信, system=系统, etc.
	Content   string    `gorm:"type:text;not null" json:"content"`  // 通知内容
	Link      string    `gorm:"size:255" json:"link"`               // 点击通知跳转的链接
	IsRead    bool      `gorm:"default:false;index" json:"is_read"` // 是否已读
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
}
