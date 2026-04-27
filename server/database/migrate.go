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
		&models.Comment{},          // 评论表
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
		&models.AntiSpamConfig{},   // 防刷配置表
		&models.UserOperation{},    // 用户操作记录表
		&models.ReputationLog{},    // 信誉分日志表
		&models.ContentQuality{},   // 内容质量表
		&models.UserBan{},          // 用户禁言表
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 创建复合索引优化查询性能
	createCompositeIndexes()

	// 同步现有收藏数到 topic 表的 favorite_count 字段
	syncFavoriteCounts()

	log.Println("database migrated successfully")
}

// syncFavoriteCounts 同步现有收藏数到 topic 表
func syncFavoriteCounts() {
	log.Println("syncing favorite counts to topics...")

	// 更新所有话题的 favorite_count
	result := DB.Exec(`
		UPDATE topics 
		SET favorite_count = (
			SELECT COUNT(*) 
			FROM favorites 
			WHERE favorites.topic_id = topics.id
		)
	`)
	if result.Error != nil {
		log.Printf("failed to sync favorite counts: %v", result.Error)
		return
	}
	if result.RowsAffected > 0 {
		log.Printf("synced favorite counts for %d topics", result.RowsAffected)
	}
}

// createCompositeIndexes 创建复合索引优化查询性能
func createCompositeIndexes() {
	// PollVote 表索引迁移：删除旧索引，创建新索引支持多选投票
	DB.Exec("DROP INDEX IF EXISTS idx_poll_user")
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_poll_user_option ON poll_votes(poll_id, user_id, option_id)")
	// Topic 表索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_topic_forum_created ON topics(forum_id, created_at DESC)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_topic_pinned_hot ON topics(is_pinned DESC, (like_count + reply_count * 2) DESC)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_topic_pinned_reply ON topics(is_pinned DESC, last_reply_at DESC)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_topic_user_created ON topics(user_id, created_at DESC)")

	// Comment 表索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_comment_topic_created ON comments(topic_id, created_at DESC)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_comment_user_created ON comments(user_id, created_at DESC)")

	// TopicTags 表索引（用于标签筛选话题）
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_topic_tags_tag ON topic_tags(tag_id)")

	// Like 表索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_like_user_topic ON likes(user_id, topic_id)")

	// Favorite 表索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_favorite_user_topic ON favorites(user_id, topic_id)")

	// Follow 表索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_follow_user ON follows(user_id)")

	// UserOperation 表索引（用于防刷检查）
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_user_op_user_type ON user_operations(user_id, operation_type)")

	// Notification 表索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_notification_user_read ON notifications(user_id, is_read)")

	log.Println("composite indexes created successfully")
}
