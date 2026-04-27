package models

import (
	"time"

	"gorm.io/gorm"
)

// Topic 话题/帖子模型
// 论坛中的主帖子，每个话题可以有多个回复帖子
type Topic struct {
	ID           uint           `gorm:"primarykey" json:"id"`                 // 话题唯一标识
	Title        string         `gorm:"size:200;not null" json:"title"`       // 话题标题
	Content      string         `gorm:"type:text;not null" json:"content"`    // 话题正文内容
	UserID       uint           `gorm:"not null;index" json:"user_id"`        // 发布者用户ID
	User         User           `gorm:"foreignKey:UserID" json:"user"`        // 发布者用户信息
	ForumID      uint           `gorm:"not null;index" json:"forum_id"`       // 所属版块ID
	Forum        Forum          `gorm:"foreignKey:ForumID" json:"forum"`      // 所属版块信息
	IsPinned     bool           `gorm:"default:false;index" json:"is_pinned"` // 是否置顶（管理员）
	IsUserPinned bool           `gorm:"default:false;index" json:"is_user_pinned"` // 作者是否置顶
	IsLocked     bool           `gorm:"default:false" json:"is_locked"`       // 是否锁定（禁止回复）
	IsEssence    bool           `gorm:"default:false" json:"is_essence"`      // 是否加精
	IsHidden     bool           `gorm:"default:false;index" json:"is_hidden"` // 是否隐藏（被举报）
	HotScore     float64        `gorm:"default:0;index" json:"hot_score"`     // 热度分数
	LikeCount    int            `gorm:"default:0" json:"like_count"`          // 点赞数
	ViewCount    int            `gorm:"default:0" json:"view_count"`          // 浏览数
	ReplyCount   int            `gorm:"default:0" json:"reply_count"`         // 回复数
	FavoriteCount int           `gorm:"default:0" json:"favorite_count"`      // 收藏数
	LastReplyAt  *time.Time     `json:"last_reply_at"`                        // 最后回复时间
	AllowComment bool           `gorm:"default:true" json:"allow_comment"`    // 是否允许评论
	CreatedAt    time.Time      `json:"created_at"`                           // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"`                           // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                       // 软删除时间

	// 关联关系
	Comments  []Comment  `gorm:"foreignKey:TopicID" json:"-"`      // 话题下的所有评论
	Likes     []Like     `gorm:"foreignKey:TargetID" json:"-"`      // 话题的点赞记录
	Favorites []Favorite `gorm:"foreignKey:TopicID" json:"-"`       // 话题的收藏记录
	Tags      []Tag      `gorm:"many2many:topic_tags;" json:"tags"` // 话题关联的标签
}
