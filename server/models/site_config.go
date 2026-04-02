package models

import "time"

// SiteConfig 网站配置模型
// 以键值对形式存储网站的各种配置项
type SiteConfig struct {
	ID        uint      `gorm:"primarykey" json:"id"`                    // 配置项唯一标识
	Key       string    `gorm:"size:50;uniqueIndex;not null" json:"key"` // 配置键，唯一索引
	Value     string    `gorm:"type:text" json:"value"`                  // 配置值
	CreatedAt time.Time `json:"created_at"`                              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                              // 更新时间
}
