package models

import "time"

// Tag 标签模型
// 用于对话题进行分类和标记
type Tag struct {
	ID          uint      `gorm:"primarykey" json:"id"`                     // 标签唯一标识
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"` // 标签名称，唯一索引
	Icon        string    `gorm:"size:10" json:"icon"`                      // 标签图标（emoji）
	Description string    `gorm:"type:text" json:"description"`             // 标签描述
	SortOrder   int       `gorm:"default:0" json:"sort_order"`              // 排序顺序
	UsageCount  int       `gorm:"default:0" json:"usage_count"`             // 使用次数（被多少话题使用）
	IsOfficial  bool      `gorm:"default:false" json:"is_official"`         // 是否官方标签
	IsBanned    bool      `gorm:"default:false" json:"is_banned"`           // 是否被禁用
	CreatedAt   time.Time `json:"created_at"`                               // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                               // 更新时间
}

// TableName 指定表名为 tags（默认情况下 GORM 会将 Tag 转换为 tags）
func (Tag) TableName() string {
	return "tags"
}
