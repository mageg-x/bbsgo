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
	"time"

	"github.com/gorilla/mux"
)

// GetPosts 获取话题的回复列表处理器
// 支持分页，返回话题下的一级回复
func GetPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])

	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var posts []models.Post
	var total int64

	offset := (page - 1) * pageSize

	// 统计一级回复数量
	if err := database.DB.Model(&models.Post{}).Where("topic_id = ? AND parent_id IS NULL", topicID).Count(&total).Error; err != nil {
		log.Printf("get posts: failed to count posts, topicID: %d, error: %v", topicID, err)
	}

	// 查询一级回复
	if err := database.DB.Where("topic_id = ? AND parent_id IS NULL", topicID).Preload("User").
		Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		log.Printf("get posts: failed to query posts, topicID: %d, error: %v", topicID, err)
		utils.Error(w, 500, "获取回复失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreatePost 创建回复处理器
// 在指定话题下创建回复
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create post: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])

	// 解析请求体
	var req struct {
		Content  string `json:"content"`   // 回复内容
		ParentID *uint  `json:"parent_id"` // 父回复ID（用于嵌套回复）
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create post: failed to decode request body, topicID: %d, error: %v", topicID, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证内容
	if req.Content == "" {
		log.Printf("create post: content is empty, topicID: %d", topicID)
		utils.Error(w, 400, "请填写内容")
		return
	}

	// 查询话题
	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		log.Printf("create post: topic not found, topicID: %d, error: %v", topicID, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	// 检查话题是否允许评论
	if topic.IsLocked || !topic.AllowComment {
		log.Printf("create post: topic is locked or not allowing comments, topicID: %d", topicID)
		utils.Error(w, 400, "该话题已关闭评论")
		return
	}

	// 创建回复
	post := models.Post{
		TopicID:  uint(topicID),
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		log.Printf("create post: failed to create post, topicID: %d, userID: %d, error: %v", topicID, userID, err)
		utils.Error(w, 500, "发布失败")
		return
	}

	// 更新话题的回复数和最后回复时间
	now := time.Now()
	if err := database.DB.Model(&topic).Updates(map[string]interface{}{
		"reply_count":   topic.ReplyCount + 1,
		"last_reply_at": now,
	}).Error; err != nil {
		log.Printf("create post: failed to update topic reply count, topicID: %d, error: %v", topicID, err)
	}

	// 重新加载回复关联数据
	if err := database.DB.Preload("User").First(&post, post.ID).Error; err != nil {
		log.Printf("create post: failed to reload post, id: %d, error: %v", post.ID, err)
	}

	log.Printf("create post: post created successfully, id: %d, topicID: %d, userID: %d", post.ID, topicID, userID)
	utils.Success(w, post)
}

// UpdatePost 更新回复处理器
// 仅回复作者可以更新
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update post: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询回复
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		log.Printf("update post: post not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "评论不存在")
		return
	}

	// 验证权限：仅作者可以更新
	if post.UserID != userID {
		log.Printf("update post: permission denied, postID: %d, userID: %d", id, userID)
		utils.Error(w, 403, "无权限编辑")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update post: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "topic_id")
	delete(updates, "parent_id")
	delete(updates, "created_at")
	delete(updates, "like_count")

	// 执行更新
	if err := database.DB.Model(&post).Updates(updates).Error; err != nil {
		log.Printf("update post: failed to update post, id: %d, userID: %d, error: %v", id, userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	// 重新加载数据
	if err := database.DB.Preload("User").First(&post, id).Error; err != nil {
		log.Printf("update post: failed to reload post, id: %d, error: %v", id, err)
	}

	log.Printf("update post: post updated successfully, id: %d, userID: %d", id, userID)
	utils.Success(w, post)
}

// DeletePost 删除回复处理器
// 作者或管理员可以删除
func DeletePost(w http.ResponseWriter, r *http.Request) {
	// 验证用户登录
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete post: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询回复
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		log.Printf("delete post: post not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "评论不存在")
		return
	}

	// 查询用户信息用于权限判断
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("delete post: user not found, userID: %d, error: %v", userID, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 验证权限：作者或管理员(role>=1)可以删除
	if post.UserID != userID && user.Role < 1 {
		log.Printf("delete post: permission denied, postID: %d, userID: %d", id, userID)
		utils.Error(w, 403, "无权限删除")
		return
	}

	// 查询关联话题用于更新回复数
	var topic models.Topic
	if err := database.DB.First(&topic, post.TopicID).Error; err != nil {
		log.Printf("delete post: topic not found, topicID: %d, error: %v", post.TopicID, err)
	}

	// 物理删除回复
	if err := database.DB.Unscoped().Delete(&post).Error; err != nil {
		log.Printf("delete post: failed to delete post, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	// 更新话题回复数
	if topic.ID != 0 && topic.ReplyCount > 0 {
		if err := database.DB.Model(&topic).UpdateColumn("reply_count", topic.ReplyCount-1).Error; err != nil {
			log.Printf("delete post: failed to update topic reply count, topicID: %d, error: %v", topic.ID, err)
		}
	}

	log.Printf("delete post: post deleted successfully, id: %d, userID: %d", id, userID)
	utils.Success(w, nil)
}
