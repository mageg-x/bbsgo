package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetTopics 获取话题列表处理器
// 支持按版块、标签筛选，支持分页
func GetTopics(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	forumID, _ := strconv.Atoi(r.URL.Query().Get("forum_id"))
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var topics []models.Topic
	var total int64

	// 构建查询
	query := database.DB.Model(&models.Topic{})
	if forumID > 0 {
		query = query.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("get topics: failed to count topics, error: %v", err)
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	dbQuery := database.DB.Preload("User").Preload("Forum").Preload("Tags")
	if forumID > 0 {
		dbQuery = dbQuery.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		dbQuery = dbQuery.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	// 排序方式
	switch sort {
	case "hot":
		dbQuery = dbQuery.Order("is_pinned DESC, (like_count + reply_count * 2) DESC, created_at DESC")
	case "reply":
		dbQuery = dbQuery.Order("is_pinned DESC, last_reply_at DESC NULLS LAST, created_at DESC")
	default:
		dbQuery = dbQuery.Order("is_pinned DESC, created_at DESC")
	}

	// 执行查询
	if err := dbQuery.Offset(offset).Limit(pageSize).Find(&topics).Error; err != nil {
		log.Printf("get topics: failed to query topics, error: %v", err)
		utils.Error(w, 500, "获取话题失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetTopic 获取单个话题详情处理器
func GetTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var topic models.Topic
	if err := database.DB.Preload("User").Preload("Forum").First(&topic, id).Error; err != nil {
		log.Printf("get topic: topic not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	// 增加浏览数
	if err := database.DB.Model(&topic).UpdateColumn("view_count", topic.ViewCount+1).Error; err != nil {
		log.Printf("get topic: failed to increment view count, id: %d, error: %v", id, err)
	}

	utils.Success(w, topic)
}

// CreateTopic 创建话题处理器
// 需要用户登录
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	// 添加 recover 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("CreateTopic PANIC: %v", err)
			utils.Error(w, 500, "服务器内部错误")
		}
	}()

	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create topic: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	// 检查是否允许发帖
	if !utils.GetConfigBool("allow_post", true) {
		log.Printf("create topic: create topic disabled")
		utils.Error(w, 400, "发帖功能已关闭")
		return
	}

	// 解析请求体
	var req struct {
		Title    string   `json:"title"`     // 话题标题
		Content  string   `json:"content"`   // 话题内容
		ForumID  uint     `json:"forum_id"`  // 版块ID
		TagNames []string `json:"tag_names"` // 标签名称列表
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create topic: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.Title == "" || req.Content == "" || req.ForumID == 0 {
		log.Printf("create topic: incomplete topic info, title: %s, forumID: %d", req.Title, req.ForumID)
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	// 验证标签数量
	if len(req.TagNames) > 3 {
		log.Printf("create topic: too many tags, count: %d", len(req.TagNames))
		utils.Error(w, 400, "最多只能添加3个标签")
		return
	}

	// 创建话题
	topic := models.Topic{
		Title:        req.Title,
		Content:      req.Content,
		UserID:       userID,
		ForumID:      req.ForumID,
		AllowComment: true,           // 默认允许评论
		Tags:         []models.Tag{}, // 初始化 Tags 为空切片，避免 nil 导致的空指针异常
	}

	if err := database.DB.Create(&topic).Error; err != nil {
		log.Printf("create topic: failed to create topic, userID: %d, forumID: %d, error: %v", userID, req.ForumID, err)
		utils.Error(w, 500, "发布失败")
		return
	}

	// 处理标签关联
	if len(req.TagNames) > 0 {
		var tags []models.Tag
		for _, name := range req.TagNames {
			if name == "" {
				continue
			}
			tag, err := GetOrCreateTagByName(name)
			if err != nil {
				log.Printf("create topic: failed to get or create tag, name: %s, error: %v", name, err)
				continue
			}
			if tag == nil {
				log.Printf("create topic: tag is nil, name: %s", name)
				continue
			}
			if tag.IsBanned {
				log.Printf("create topic: tag is banned, name: %s", name)
				continue
			}
			// 只保留 ID，避免 GORM 尝试插入完整对象
			tags = append(tags, models.Tag{ID: tag.ID})
			IncrementTagUsage(tag.ID)
		}
		if len(tags) > 0 {
			if err := database.DB.Model(&topic).Association("Tags").Replace(tags); err != nil {
				log.Printf("create topic: failed to associate tags, topicID: %d, error: %v", topic.ID, err)
			}
		}
	}

	// 给发帖用户增加积分
	var user models.User
	if err := database.DB.First(&user, userID).Error; err == nil {
		log.Printf("create topic: [DEBUG] before GetConfigInt, database.DB=%v", database.DB != nil)
		creditAmount := utils.GetConfigInt("credit_topic", 20)
		log.Printf("create topic: [DEBUG] after GetConfigInt, creditAmount=%d", creditAmount)
		user.Credits += creditAmount
		if err := database.DB.Save(&user).Error; err != nil {
			log.Printf("create topic: failed to add credits, userID: %d, error: %v", userID, err)
		}
	}

	// 重新加载话题关联数据
	log.Printf("create topic: [DEBUG] before Preload, topic.ID=%d", topic.ID)
	if err := database.DB.Preload("User").Preload("Forum").Preload("Tags").First(&topic, topic.ID).Error; err != nil {
		log.Printf("create topic: failed to reload topic, id: %d, error: %v", topic.ID, err)
	}
	log.Printf("create topic: [DEBUG] after Preload, topic.ID=%d", topic.ID)

	// 确保 Tags 不为 nil
	if topic.Tags == nil {
		topic.Tags = []models.Tag{}
	}

	log.Printf("create topic: [DEBUG] before Success, w=%v, topic=%+v", w != nil, topic)
	utils.Success(w, topic)
	log.Printf("create topic: [DEBUG] after Success")
}

// UpdateTopic 更新话题处理器
// 仅话题作者可以更新
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update topic: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("update topic: topic not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	// 验证权限：仅作者可以更新
	if topic.UserID != userID {
		log.Printf("update topic: permission denied, topicID: %d, userID: %d", id, userID)
		utils.Error(w, 403, "无权限编辑")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update topic: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")
	delete(updates, "is_pinned")
	delete(updates, "is_locked")
	delete(updates, "is_essence")
	delete(updates, "like_count")
	delete(updates, "view_count")
	delete(updates, "reply_count")

	// 执行更新
	if err := database.DB.Model(&topic).Updates(updates).Error; err != nil {
		log.Printf("update topic: failed to update topic, id: %d, userID: %d, error: %v", id, userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	// 重新加载数据
	if err := database.DB.Preload("User").Preload("Forum").First(&topic, id).Error; err != nil {
		log.Printf("update topic: failed to reload topic, id: %d, error: %v", id, err)
	}

	log.Printf("update topic: topic updated successfully, id: %d, userID: %d", id, userID)
	utils.Success(w, topic)
}

// DeleteTopic 删除话题处理器
// 作者或管理员可以删除
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete topic: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("delete topic: topic not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	// 查询用户信息用于权限判断
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("delete topic: user not found, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 验证权限：作者或管理员(role>=1)可以删除
	if topic.UserID != userID && user.Role < 1 {
		log.Printf("delete topic: permission denied, topicID: %d, userID: %d, topicUserID: %d", id, userID, topic.UserID)
		utils.Error(w, 403, "无权限删除")
		return
	}

	// 物理删除话题（不保留软删除）
	if err := database.DB.Unscoped().Delete(&topic).Error; err != nil {
		log.Printf("delete topic: failed to delete topic, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	log.Printf("delete topic: topic deleted successfully, id: %d, userID: %d", id, userID)
	utils.Success(w, nil)
}
