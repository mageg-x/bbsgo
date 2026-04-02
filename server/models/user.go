package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// 存储论坛用户的基本信息和状态
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`                         // 用户唯一标识
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"` // 用户名，唯一索引
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`   // 邮箱，唯一索引
	Nickname     string         `gorm:"size:50" json:"nickname"`                      // 昵称
	PasswordHash string         `gorm:"size:255;not null" json:"-"`                   // 密码哈希，json序列化时隐藏
	Role         int            `gorm:"default:0" json:"role"`                        // 用户角色：0普通用户 1版主 2管理员
	Avatar       string         `gorm:"size:255" json:"avatar"`                       // 头像URL
	Background   string         `gorm:"size:255" json:"background"`                   // 个人主页背景图
	Signature    string         `gorm:"size:255" json:"signature"`                    // 个性签名
	Intro        string         `gorm:"type:text" json:"intro"`                       // 个人简介
	Credits      int            `gorm:"default:0" json:"credits"`                     // 积分
	Level        int            `gorm:"default:1" json:"level"`                       // 等级
	LastSignAt   *time.Time     `json:"last_sign_at"`                                 // 最后签到时间
	CreatedAt    time.Time      `json:"created_at"`                                   // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"`                                   // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                               // 软删除时间

	// 关联关系
	Topics           []Topic        `gorm:"foreignKey:UserID" json:"-"`            // 用户发布的话题
	Posts            []Post         `gorm:"foreignKey:UserID" json:"-"`            // 用户发布的帖子/回复
	Likes            []Like         `gorm:"foreignKey:UserID" json:"-"`            // 用户点赞记录
	Favorites        []Favorite     `gorm:"foreignKey:UserID" json:"-"`            // 用户收藏记录
	Follows          []Follow       `gorm:"foreignKey:UserID" json:"-"`            // 用户关注的人（正向）
	Followers        []Follow       `gorm:"foreignKey:FollowUserID" json:"-"`      // 用户的粉丝（反向）
	SentMessages     []Message      `gorm:"foreignKey:FromUserID" json:"-"`        // 用户发送的私信
	ReceivedMessages []Message      `gorm:"foreignKey:ToUserID" json:"-"`          // 用户接收的私信
	Notifications    []Notification `gorm:"foreignKey:UserID" json:"-"`            // 用户通知
	Reports          []Report       `gorm:"foreignKey:ReporterID" json:"-"`        // 用户发起的举报
	Drafts           []Draft        `gorm:"foreignKey:UserID" json:"-"`            // 用户草稿
	UserBadges       []UserBadge    `gorm:"foreignKey:UserID" json:"-"`            // 用户获得的勋章
	FollowedTopics   []Topic        `gorm:"many2many:user_follow_topics" json:"-"` // 用户关注的话题
}
