package models

import "time"

// Badge 勋章模型
// 存储系统中可授予用户的勋章信息
type Badge struct {
	ID          uint   `gorm:"primarykey" json:"id"`         // 勋章唯一标识
	Name        string `gorm:"size:50;not null" json:"name"` // 勋章名称
	Description string `gorm:"type:text" json:"description"` // 勋章描述
	Icon        string `gorm:"size:255" json:"icon"`         // 勋章图标URL
	Condition   string `gorm:"type:text" json:"condition"`   // 获得条件描述

	// 勋章类型：basic(基础入门), advanced(进阶成就), top(顶级荣耀)
	Type string `gorm:"size:20;default:'basic'" json:"type"`

	// 条件类型：topic_count(发帖数), comment_count(评论数), like_count(获赞数),
	// best_comment(最佳评论), follower_count(粉丝数), register_days(注册天数),
	// combination(组合条件)
	ConditionType string `gorm:"size:50" json:"condition_type"`

	// 条件值（数值型条件）
	ConditionValue int `gorm:"default:0" json:"condition_value"`

	// 组合条件（JSON格式，用于"社区传奇"等复杂条件）
	// 例如：{"register_days": 730, "topic_count": 200, "like_count": 500, "best_comment": 10}
	ConditionData string `gorm:"type:text" json:"condition_data"`

	// 排序权重（数字越小越靠前）
	SortOrder int `gorm:"default:0" json:"sort_order"`

	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// TableName 指定表名为 badges
func (Badge) TableName() string {
	return "badges"
}
