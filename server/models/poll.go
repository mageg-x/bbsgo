package models

import (
	"time"

	"gorm.io/gorm"
)

// Poll 投票模型
type Poll struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	TopicID     uint           `gorm:"not null;uniqueIndex" json:"topic_id"`
	Topic       Topic          `gorm:"foreignKey:TopicID" json:"-"`
	Title       string         `gorm:"size:200" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	PollType    string         `gorm:"size:20;default:'single'" json:"poll_type"`
	MaxChoices  int            `gorm:"default:1" json:"max_choices"`
	EndTime     *time.Time     `json:"end_time"`
	IsEnded     bool           `gorm:"default:false" json:"is_ended"`
	TotalVotes  int            `gorm:"default:0" json:"total_votes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Options     []PollOption   `gorm:"foreignKey:PollID" json:"options"`
}

// PollOption 投票选项模型
type PollOption struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	PollID    uint           `gorm:"not null;index" json:"poll_id"`
	Poll      Poll           `gorm:"foreignKey:PollID" json:"-"`
	Text      string         `gorm:"size:200;not null" json:"text"`
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	VoteCount int            `gorm:"default:0" json:"vote_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PollVote 投票记录模型
type PollVote struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	PollID     uint           `gorm:"not null;index:idx_poll_user,unique" json:"poll_id"`
	Poll       Poll           `gorm:"foreignKey:PollID" json:"-"`
	OptionID   uint           `gorm:"not null;index" json:"option_id"`
	Option     PollOption     `gorm:"foreignKey:OptionID" json:"-"`
	UserID     uint           `gorm:"not null;index:idx_poll_user,unique" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Poll) TableName() string {
	return "polls"
}

func (PollOption) TableName() string {
	return "poll_options"
}

func (PollVote) TableName() string {
	return "poll_votes"
}
