package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
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
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 构建响应，包含关联详细信息
	type NotificationWithRelated struct {
		models.Notification
		Badge       *models.Badge  `json:"badge,omitempty"`
		Message     *models.Message `json:"message,omitempty"`
		FromUser    *models.User   `json:"from_user,omitempty"`
		RelatedUser *models.User   `json:"related_user,omitempty"`
	}

	notificationsWithRelated := make([]NotificationWithRelated, len(notifications))
	for i, n := range notifications {
		notificationsWithRelated[i] = NotificationWithRelated{Notification: n}

		// 根据类型预加载关联数据
		switch n.Type {
		case "badge":
			// 勋章通知
			if n.RelatedType == "badge" && n.RelatedID > 0 {
				var badge models.Badge
				if err := database.DB.First(&badge, n.RelatedID).Error; err == nil {
					notificationsWithRelated[i].Badge = &badge
				}
			}
		case "message":
			// 私信通知，预加载发送者信息
			if n.RelatedID > 0 {
				var message models.Message
				if err := database.DB.First(&message, n.RelatedID).Error; err == nil {
					notificationsWithRelated[i].Message = &message
					// 加载发送者信息
					var fromUser models.User
					if err := database.DB.Select("id, username, nickname, avatar").First(&fromUser, message.FromUserID).Error; err == nil {
						notificationsWithRelated[i].FromUser = &fromUser
					}
				}
			}
		case "follow":
			// 关注通知，预加载关注者信息
			if n.RelatedType == "user" && n.RelatedID > 0 {
				var relatedUser models.User
				if err := database.DB.Select("id, username, nickname, avatar").First(&relatedUser, n.RelatedID).Error; err == nil {
					notificationsWithRelated[i].RelatedUser = &relatedUser
				}
			}
		case "best_comment":
			// 最佳评论通知，预加载评论信息
			if n.RelatedType == "comment" && n.RelatedID > 0 {
				var comment models.Comment
				if err := database.DB.First(&comment, n.RelatedID).Error; err == nil {
					// 加载评论者信息
					var commentUser models.User
					if err := database.DB.Select("id, username, nickname, avatar").First(&commentUser, comment.UserID).Error; err == nil {
						notificationsWithRelated[i].FromUser = &commentUser
					}
				}
			}
		}
	}

	errors.Success(w, map[string]interface{}{
		"list":      notificationsWithRelated,
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
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]int64{"count": count})
}

// MarkAllNotificationsRead 标记所有通知已读处理器
func MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	if err := database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error; err != nil {
		log.Printf("mark all notifications read: failed to mark notifications as read, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("mark all notifications read: notifications marked as read, userID: %d", userID)
	errors.Success(w, nil)
}

// CreateNotification 创建通知（内部函数）
// userID: 接收通知的用户ID
// notifType: 通知类型
// content: 通知内容（i18n key）
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

// CreateNotificationWithRelated 创建带关联数据的通知
// userID: 接收通知的用户ID
// notifType: 通知类型
// content: 通知内容（i18n key）
// link: 点击跳转链接
// relatedID: 关联对象ID
// relatedType: 关联对象类型
func CreateNotificationWithRelated(userID uint, notifType, content, link string, relatedID uint, relatedType string) {
	notification := models.Notification{
		UserID:      userID,
		Type:        notifType,
		Content:     content,
		Link:        link,
		RelatedID:   relatedID,
		RelatedType: relatedType,
		IsRead:      false,
	}
	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("create notification: failed to create notification, userID: %d, type: %s, error: %v", userID, notifType, err)
	}
}
