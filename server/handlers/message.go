package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetMessages 获取当前用户的私信列表处理器
// 包括发送和接收的私信
func GetMessages(w http.ResponseWriter, r *http.Request) {
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

	var messages []models.Message
	var total int64

	offset := (page - 1) * pageSize

	// 统计私信总数
	if err := database.DB.Model(&models.Message{}).Where("to_user_id = ? OR from_user_id = ?", userID, userID).Count(&total).Error; err != nil {
		log.Printf("get messages: failed to count messages, userID: %d, error: %v", userID, err)
	}

	// 查询私信列表
	if err := database.DB.Where("to_user_id = ? OR from_user_id = ?", userID, userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Preload("FromUser").
		Preload("ToUser").
		Find(&messages).Error; err != nil {
		log.Printf("get messages: failed to query messages, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
		"list":      messages,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetMessageConversation 获取与指定用户的私信会话处理器
func GetMessageConversation(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	otherUserID, _ := strconv.Atoi(vars["user_id"])

	// 查询与该用户的所有私信
	var messages []models.Message
	if err := database.DB.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		userID, otherUserID, otherUserID, userID).
		Order("created_at ASC").
		Preload("FromUser").
		Preload("ToUser").
		Find(&messages).Error; err != nil {
		log.Printf("get message conversation: failed to query messages, userID: %d, otherUserID: %d, error: %v", userID, otherUserID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 将与该用户的私信标记为已读
	if err := database.DB.Model(&models.Message{}).
		Where("from_user_id = ? AND to_user_id = ? AND is_read = ?", otherUserID, userID, false).
		Update("is_read", true).Error; err != nil {
		log.Printf("get message conversation: failed to mark messages as read, userID: %d, otherUserID: %d, error: %v", userID, otherUserID, err)
	}

	errors.Success(w, messages)
}

// SendMessage 发送私信处理器
func SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 解析请求体
	var req struct {
		ToUserID uint   `json:"to_user_id"` // 接收者用户ID
		Content  string `json:"content"`    // 私信内容
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("send message: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 验证内容
	if req.Content == "" {
		log.Printf("send message: content is empty, fromUserID: %d, toUserID: %d", userID, req.ToUserID)
		errors.Error(w, errors.CodeIncompleteInfo, "")
		return
	}

	// 检查接收者是否存在
	var toUser models.User
	if err := database.DB.First(&toUser, req.ToUserID).Error; err != nil {
		log.Printf("send message: recipient not found, toUserID: %d, error: %v", req.ToUserID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	// 创建私信
	message := models.Message{
		FromUserID: userID,
		ToUserID:   req.ToUserID,
		Content:    req.Content,
		IsRead:     false,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		log.Printf("send message: failed to create message, fromUserID: %d, toUserID: %d, error: %v", userID, req.ToUserID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 创建通知
	notification := models.Notification{
		UserID:  req.ToUserID,
		Type:    "message",
		Content: "您收到了一条新私信",
		Link:    "/messages",
		IsRead:  false,
	}
	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("send message: failed to create notification, toUserID: %d, error: %v", req.ToUserID, err)
	}

	log.Printf("send message: message sent successfully, fromUserID: %d, toUserID: %d", userID, req.ToUserID)
	errors.Success(w, message)
}

// MarkMessagesRead 标记私信已读处理器
// 标记与指定用户的所有私信为已读
func MarkMessagesRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 解析请求体
	var req struct {
		FromUserID uint `json:"from_user_id"` // 发送者用户ID
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("mark messages read: failed to decode request body, error: %v", err)
	}

	// 标记为已读
	if err := database.DB.Model(&models.Message{}).
		Where("to_user_id = ? AND from_user_id = ? AND is_read = ?", userID, req.FromUserID, false).
		Update("is_read", true).Error; err != nil {
		log.Printf("mark messages read: failed to mark messages as read, userID: %d, fromUserID: %d, error: %v", userID, req.FromUserID, err)
	}

	errors.Success(w, nil)
}

// GetUnreadMessageCount 获取未读私信数量处理器
func GetUnreadMessageCount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var count int64
	if err := database.DB.Model(&models.Message{}).Where("to_user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
		log.Printf("get unread message count: failed to count messages, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]int64{"count": count})
}
