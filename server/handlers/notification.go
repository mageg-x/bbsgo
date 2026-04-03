package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"strconv"
)

// GetNotifications 获取当前用户的通知列表处理器
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var notifications []models.Notification
	var total int64

	offset := (page - 1) * pageSize

	// 统计通知总数
	if err := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		log.Printf("get notifications: failed to count notifications, userID: %d, error: %v", userID, err)
	}

	// 查询通知列表
	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications).Error; err != nil {
		log.Printf("get notifications: failed to query notifications, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取通知列表失败")
		return
	}

	// 构建响应，包含勋章详细信息
	type NotificationWithBadge struct {
		models.Notification
		Badge *models.Badge `json:"badge,omitempty"`
	}

	notificationsWithBadge := make([]NotificationWithBadge, len(notifications))
	for i, n := range notifications {
		notificationsWithBadge[i] = NotificationWithBadge{Notification: n}
		// 如果是勋章通知，查询勋章详细信息
		if n.Type == "badge" && n.RelatedType == "badge" && n.RelatedID > 0 {
			var badge models.Badge
			if err := database.DB.First(&badge, n.RelatedID).Error; err == nil {
				notificationsWithBadge[i].Badge = &badge
			}
		}
	}

	utils.Success(w, map[string]interface{}{
		"list":      notificationsWithBadge,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUnreadNotificationCount 获取未读通知数量处理器
func GetUnreadNotificationCount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var count int64
	if err := database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
		log.Printf("get unread notification count: failed to count notifications, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取未读数量失败")
		return
	}

	utils.Success(w, map[string]int64{"count": count})
}

// MarkAllNotificationsRead 标记所有通知已读处理器
func MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	if err := database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error; err != nil {
		log.Printf("mark all notifications read: failed to mark notifications as read, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "标记已读失败")
		return
	}

	log.Printf("mark all notifications read: notifications marked as read, userID: %d", userID)
	utils.Success(w, nil)
}

// CreateNotification 创建通知（内部函数）
// userID: 接收通知的用户ID
// notifType: 通知类型
// content: 通知内容
// link: 点击跳转链接
func CreateNotification(userID uint, notifType, content, link string) {
	notification := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Content: content,
		Link:    link,
		IsRead:  false,
	}
	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("create notification: failed to create notification, userID: %d, type: %s, error: %v", userID, notifType, err)
	}
}
