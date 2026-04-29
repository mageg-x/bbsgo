package services

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"time"
)

// PrivateTopicScheduler 私密帖子定时任务调度器
// 用于检查过期的私密帖子并自动转为公开
type PrivateTopicScheduler struct {
	ticker   *time.Ticker
	stopChan chan struct{}
}

// NewPrivateTopicScheduler 创建新的私密帖子调度器
func NewPrivateTopicScheduler() *PrivateTopicScheduler {
	return &PrivateTopicScheduler{
		stopChan: make(chan struct{}),
	}
}

// Start 启动定时任务
// interval: 检查间隔，建议使用每分钟检查一次
func (s *PrivateTopicScheduler) Start(interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute // 默认每分钟检查一次
	}

	s.ticker = time.NewTicker(interval)

	go func() {
		log.Println("private topic scheduler started")
		for {
			select {
			case <-s.ticker.C:
				s.checkAndExpirePrivateTopics()
			case <-s.stopChan:
				log.Println("private topic scheduler stopped")
				return
			}
		}
	}()
}

// Stop 停止定时任务
func (s *PrivateTopicScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
}

// checkAndExpirePrivateTopics 检查并处理过期的私密帖子
func (s *PrivateTopicScheduler) checkAndExpirePrivateTopics() {
	now := time.Now()

	// 查询所有需要处理的私密帖子：
	// 1. 是私密状态 (is_private = true)
	// 2. 未结束 (private_ended = false)
	// 3. 有过期时间且已过期 (private_expire_at <= now)
	var topics []models.Topic
	if err := database.DB.
		Where("is_private = ? AND private_ended = ? AND private_expire_at IS NOT NULL AND private_expire_at <= ?", true, false, now).
		Find(&topics).Error; err != nil {
		log.Printf("private topic scheduler: failed to query expired private topics, error: %v", err)
		return
	}

	if len(topics) == 0 {
		return
	}

	log.Printf("private topic scheduler: found %d expired private topics to process", len(topics))

	// 批量处理过期的私密帖子
	for _, topic := range topics {
		// 更新帖子状态为公开
		updates := map[string]interface{}{
			"is_private":    false,
			"private_ended": true,
			"allow_comment": true, // 解禁后允许评论
		}

		if err := database.DB.Model(&topic).Updates(updates).Error; err != nil {
			log.Printf("private topic scheduler: failed to update topic %d, error: %v", topic.ID, err)
			continue
		}

		// 清除缓存
		cache.TopicCache.Invalidate(topic.ID)

		log.Printf("private topic scheduler: topic %d expired and made public", topic.ID)
	}

	// 清除首页缓存
	if len(topics) > 0 {
		cache.HomePageCache.InvalidateTopics()
	}
}

// StartPrivateTopicScheduler 启动全局私密帖子调度器
func StartPrivateTopicScheduler() {
	scheduler := NewPrivateTopicScheduler()
	scheduler.Start(time.Minute) // 每分钟检查一次
}
