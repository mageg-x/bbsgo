package models

import "time"

// Report 举报模型
// 存储用户发起的举报信息
type Report struct {
	ID         uint       `gorm:"primarykey" json:"id"`                  // 举报唯一标识
	ReporterID uint       `gorm:"not null;index" json:"reporter_id"`     // 举报者用户ID
	Reporter   User       `gorm:"foreignKey:ReporterID" json:"reporter"` // 举报者用户信息
	TargetType string     `gorm:"size:20;not null" json:"target_type"`   // 被举报对象类型：topic=话题, post=帖子, user=用户
	TargetID   uint       `gorm:"not null" json:"target_id"`             // 被举报对象ID
	Reason     string     `gorm:"type:text;not null" json:"reason"`      // 举报原因
	Status     int        `gorm:"default:0;index" json:"status"`         // 处理状态：0=待处理, 1=已处理, 2=已忽略
	HandledAt  *time.Time `json:"handled_at"`                            // 处理时间
	HandlerID  *uint      `json:"handler_id"`                            // 处理人用户ID
	Handler    *User      `gorm:"foreignKey:HandlerID" json:"handler"`   // 处理人用户信息
	CreatedAt  time.Time  `json:"created_at"`                            // 举报时间
}
