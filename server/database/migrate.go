package database

import (
	"bbsgo/models"
	"log"
)

// AutoMigrate 自动执行数据库迁移
// 根据模型定义自动创建或更新数据表结构
func AutoMigrate() {
	// 定义需要迁移的所有模型
	err := DB.AutoMigrate(
		&models.User{},             // 用户表
		&models.Forum{},            // 版块表
		&models.Topic{},            // 话题表
		&models.Post{},             // 帖子表
		&models.Like{},             // 点赞表
		&models.Favorite{},         // 收藏表
		&models.Follow{},           // 关注表
		&models.Message{},          // 私信表
		&models.Notification{},     // 通知表
		&models.Tag{},              // 标签表
		&models.Report{},           // 举报表
		&models.Badge{},            // 勋章表
		&models.UserBadge{},        // 用户勋章关联表
		&models.SiteConfig{},       // 网站配置表
		&models.Draft{},            // 草稿箱表
		&models.Announcement{},     // 公告表
		&models.VerificationCode{}, // 验证码表
		&models.Poll{},             // 投票表
		&models.PollOption{},       // 投票选项表
		&models.PollVote{},         // 投票记录表
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("database migrated successfully")
}
