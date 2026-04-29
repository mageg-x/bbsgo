package utils

import (
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"time"
)

// AnonymousUserInfo 匿名用户信息
type AnonymousUserInfo struct {
	ID           uint               `json:"id"`
	Username     string             `json:"username"`
	Nickname     string             `json:"nickname"`
	Avatar       string             `json:"avatar"`
	IsAnonymous  bool               `json:"is_anonymous"`
	Badges       []models.UserBadge `json:"badges"`
}

// GetAnonymousUserInfo 获取匿名用户信息
func GetAnonymousUserInfo() AnonymousUserInfo {
	return AnonymousUserInfo{
		ID:          0,
		Username:    "anonymous",
		Nickname:    "匿名用户",
		Avatar:      "",
		IsAnonymous: true,
		Badges:      []models.UserBadge{},
	}
}

// ShouldShowAnonymous 判断当前用户是否应该看到匿名状态
// viewerID: 当前查看者用户ID，如果为0表示未登录
// topicAuthorID: 帖子作者用户ID
// viewerRole: 当前查看者角色（0普通用户, 1版主, 2管理员）
// topic: 话题对象
func ShouldShowAnonymous(viewerID uint, topicAuthorID uint, viewerRole int, topic *models.Topic) bool {
	// 管理员可以看到真实身份
	if viewerRole >= 2 {
		return false
	}

	// 帖子作者可以看到自己的真实身份
	if viewerID == topicAuthorID && viewerID != 0 {
		return false
	}

	// 检查话题当前是否处于匿名状态
	return topic.IsCurrentlyAnonymous()
}

// ShouldShowAnonymousForComment 判断评论是否应该显示匿名
// viewerID: 当前查看者用户ID
// commentAuthorID: 评论作者用户ID
// viewerRole: 当前查看者角色
// topic: 话题对象
func ShouldShowAnonymousForComment(viewerID uint, commentAuthorID uint, viewerRole int, topic *models.Topic) bool {
	// 管理员可以看到真实身份
	if viewerRole >= 2 {
		return false
	}

	// 评论作者可以看到自己的真实身份
	if viewerID == commentAuthorID && viewerID != 0 {
		return false
	}

	// 话题作者可以看到所有评论的真实身份（方便管理）
	if viewerID == topic.UserID && viewerID != 0 {
		return false
	}

	// 检查话题当前是否处于匿名状态
	return topic.IsCurrentlyAnonymous()
}

// ApplyAnonymityToUser 对用户应用匿名处理
func ApplyAnonymityToUser(user *models.User) AnonymousUserInfo {
	if user == nil {
		return GetAnonymousUserInfo()
	}
	return AnonymousUserInfo{
		ID:          0,
		Username:    "anonymous",
		Nickname:    "匿名用户",
		Avatar:      "",
		IsAnonymous: true,
		Badges:      []models.UserBadge{},
	}
}

// ApplyAnonymityToReplyUser 对回复用户应用匿名处理
func ApplyAnonymityToReplyUser(user *models.User) *struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
} {
	return &struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
	}{
		ID:       0,
		Username: "anonymous",
		Nickname: "匿名用户",
	}
}

// ProcessExpiredAnonymousTopics 处理过期的定时匿名话题
// 将已达到解匿时间的定时匿名话题标记为已结束匿名结束
func ProcessExpiredAnonymousTopics() {
	now := time.Now()

	// 查询所有需要解匿的话题：
	// 1. is_anonymous = true
	// 2. anonymous_type = 'timed'
	// 3. is_anonymous_ended = false
	// 4. anonymous_until <= now
	result := database.DB.Model(&models.Topic{}).
		Where("is_anonymous = ?", true).
		Where("anonymous_type = ?", models.AnonymousTypeTimed).
		Where("is_anonymous_ended = ?", false).
		Where("anonymous_until <= ?", now).
		Update("is_anonymous_ended", true)

	if result.Error != nil {
		log.Printf("[anonymity] failed to process expired anonymous topics: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("[anonymity] processed %d expired anonymous topics", result.RowsAffected)
	}
}

// StartAnonymityCleanupTask 启动匿名清理定时任务
// 每分钟检查一次是否有需要解匿的话题
func StartAnonymityCleanupTask() {
	go func() {
		// 立即执行一次
		ProcessExpiredAnonymousTopics()

		// 然后每分钟检查一次
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			ProcessExpiredAnonymousTopics()
		}
	}()
	log.Printf("[anonymity] anonymity cleanup task started")
}

// init 初始化定时清理任务
func init() {
	// 使用延迟启动，确保数据库已初始化
	// 实际启动将在应用初始化后由 main 函数手动调用
}
