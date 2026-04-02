package models

import "time"

// Message 私信模型
// 存储用户之间的私信消息
type Message struct {
	ID         uint      `gorm:"primarykey" json:"id"`                   // 私信唯一标识
	FromUserID uint      `gorm:"not null;index" json:"from_user_id"`     // 发送者用户ID
	FromUser   User      `gorm:"foreignKey:FromUserID" json:"from_user"` // 发送者用户信息
	ToUserID   uint      `gorm:"not null;index" json:"to_user_id"`       // 接收者用户ID
	ToUser     User      `gorm:"foreignKey:ToUserID" json:"to_user"`     // 接收者用户信息
	Content    string    `gorm:"type:text;not null" json:"content"`      // 私信正文内容
	IsRead     bool      `gorm:"default:false;index" json:"is_read"`     // 是否已读
	CreatedAt  time.Time `json:"created_at"`                             // 发送时间
}
