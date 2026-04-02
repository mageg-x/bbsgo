package models

import "time"

// ForumCategory 版块分类模型
// 用于对版块进行分组管理
type ForumCategory struct {
	ID          uint      `gorm:"primarykey" json:"id"`          // 分类唯一标识
	Name        string    `gorm:"size:50;not null" json:"name"`  // 分类名称
	Icon        string    `gorm:"size:10" json:"icon"`           // 分类图标（emoji）
	Description string    `gorm:"size:200" json:"description"`   // 分类描述
	SortOrder   int       `gorm:"default:0" json:"sort_order"`   // 排序顺序
	IsActive    bool      `gorm:"default:true" json:"is_active"` // 是否启用
	CreatedAt   time.Time `json:"created_at"`                    // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                    // 更新时间
}
