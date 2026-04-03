package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// IntSlice 自定义类型，用于存储整数数组
// 实现 sql.Scanner 和 driver.Valuer 接口以支持 JSON 序列化
type IntSlice []int

// Value 实现 driver.Valuer 接口，将 IntSlice 转换为 JSON 字节数组
func (s IntSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口，从数据库读取 JSON 数据并转换为 IntSlice
func (s *IntSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	if len(bytes) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// Forum 版块模型
// 论坛的主要分类单元，每个版块可以包含多个话题
type Forum struct {
	ID           uint      `gorm:"primarykey" json:"id"`           // 版块唯一标识
	Name         string    `gorm:"size:100;not null" json:"name"`  // 版块名称
	Description  string    `gorm:"type:text" json:"description"`   // 版块描述
	SortOrder    int       `gorm:"default:0" json:"sort_order"`    // 排序顺序，数字越小越靠前
	Icon         string    `gorm:"size:255" json:"icon"`           // 版块图标URL
	ModeratorIDs IntSlice  `gorm:"type:json" json:"moderator_ids"` // 版主用户ID列表
	AllowPost    bool      `gorm:"default:true" json:"allow_post"` // 是否允许发布话题
	CreatedAt    time.Time `json:"created_at"`                     // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`                     // 更新时间

	Topics []Topic `gorm:"foreignKey:ForumID" json:"-"` // 版块下的话题列表
}
