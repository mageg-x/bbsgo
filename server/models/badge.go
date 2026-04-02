package models

import "time"

// Badge 勋章模型
// 存储系统中可授予用户的勋章信息
type Badge struct {
	ID          uint      `gorm:"primarykey" json:"id"`         // 勋章唯一标识
	Name        string    `gorm:"size:50;not null" json:"name"` // 勋章名称
	Description string    `gorm:"type:text" json:"description"` // 勋章描述
	Icon        string    `gorm:"size:255" json:"icon"`         // 勋章图标URL
	Condition   string    `gorm:"type:text" json:"condition"`   // 获得条件
	CreatedAt   time.Time `json:"created_at"`                   // 创建时间
}
