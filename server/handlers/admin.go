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

// ========== 用户管理 ==========

// GetAdminUsers 获取用户列表处理器（管理员）
// 支持分页，返回所有用户信息
func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var users []models.User
	var total int64

	offset := (page - 1) * pageSize

	// 统计用户总数
	if err := database.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		log.Printf("get admin users: failed to count users, error: %v", err)
	}

	// 查询用户列表
	if err := database.DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		log.Printf("get admin users: failed to query users, error: %v", err)
		utils.Error(w, 500, "获取用户列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateUserRole 更新用户角色处理器（管理员）
// 角色：0=普通用户, 1=版主, 2=管理员
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 解析请求体
	var req struct {
		Role int `json:"role"` // 新角色
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update user role: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证角色值
	if req.Role < 0 || req.Role > 2 {
		log.Printf("update user role: invalid role value, id: %d, role: %d", id, req.Role)
		utils.Error(w, 400, "无效的角色")
		return
	}

	// 查询要更新的用户
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("update user role: user not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 获取当前管理员ID
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 不能修改自己的角色
	if uint(id) == userID {
		log.Printf("update user role: cannot modify own role, id: %d, userID: %d", id, userID)
		utils.Error(w, 400, "不能修改自己的角色")
		return
	}

	// 更新角色
	if err := database.DB.Model(&user).Update("role", req.Role).Error; err != nil {
		log.Printf("update user role: failed to update role, id: %d, role: %d, error: %v", id, req.Role, err)
		utils.Error(w, 500, "更新角色失败")
		return
	}

	log.Printf("update user role: role updated successfully, id: %d, newRole: %d", id, req.Role)
	utils.Success(w, user)
}

// BanUser 封禁/解封用户处理器（管理员）
func BanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要封禁的用户
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("ban user: user not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 获取当前管理员ID
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 不能封禁自己
	if uint(id) == userID {
		log.Printf("ban user: cannot ban yourself, id: %d, userID: %d", id, userID)
		utils.Error(w, 400, "不能封禁自己")
		return
	}

	// 解析请求体
	var req struct {
		Banned bool `json:"banned"` // 是否封禁
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ban user: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 更新封禁状态
	if err := database.DB.Model(&user).Update("is_banned", req.Banned).Error; err != nil {
		log.Printf("ban user: failed to update banned status, id: %d, banned: %v, error: %v", id, req.Banned, err)
		utils.Error(w, 500, "操作失败")
		return
	}

	log.Printf("ban user: user banned status updated, id: %d, banned: %v", id, req.Banned)
	utils.Success(w, user)
}

// DeleteUser 删除用户处理器（管理员）
// 执行物理删除，同时删除用户的所有相关数据
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 获取当前管理员ID
	adminID, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("delete user: unauthorized access, id: %d", id)
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 不能删除自己
	if uint(id) == adminID {
		log.Printf("delete user: cannot delete yourself, id: %d, adminID: %d", id, adminID)
		utils.Error(w, 400, "不能删除自己")
		return
	}

	// 查询要删除的用户
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Printf("delete user: user not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "用户不存在")
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("delete user: failed to begin transaction, error: %v", tx.Error)
		utils.Error(w, 500, "操作失败")
		return
	}

	// 删除用户的帖子
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Post{}).Error; err != nil {
		tx.Rollback()
		log.Printf("delete user: failed to delete user posts, userId: %d, error: %v", id, err)
		utils.Error(w, 500, "删除用户帖子失败")
		return
	}

	// 删除用户的收藏
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Favorite{}).Error; err != nil {
		tx.Rollback()
		log.Printf("delete user: failed to delete user favorites, userId: %d, error: %v", id, err)
		utils.Error(w, 500, "删除用户收藏失败")
		return
	}

	// 删除用户的点赞
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Like{}).Error; err != nil {
		tx.Rollback()
		log.Printf("delete user: failed to delete user likes, userId: %d, error: %v", id, err)
		utils.Error(w, 500, "删除用户点赞失败")
		return
	}

	// 删除用户的通知
	if err := tx.Unscoped().Where("user_id = ?", id).Delete(&models.Notification{}).Error; err != nil {
		tx.Rollback()
		log.Printf("delete user: failed to delete user notifications, userId: %d, error: %v", id, err)
		utils.Error(w, 500, "删除用户通知失败")
		return
	}

	// 删除用户的消息
	if err := tx.Unscoped().Where("from_user_id = ? OR to_user_id = ?", id, id).Delete(&models.Message{}).Error; err != nil {
		tx.Rollback()
		log.Printf("delete user: failed to delete user messages, userId: %d, error: %v", id, err)
		utils.Error(w, 500, "删除用户消息失败")
		return
	}

	// 删除用户本身
	if err := tx.Unscoped().Delete(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("delete user: failed to delete user, userId: %d, error: %v", id, err)
		utils.Error(w, 500, "删除用户失败")
		return
	}

	// 提交事务
	tx.Commit()
	log.Printf("delete user: user deleted successfully, userId: %d, username: %s", id, user.Username)
	utils.Success(w, nil)
}

// ========== 版块管理 ==========

// CreateForum 创建版块处理器（管理员）
func CreateForum(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var forum models.Forum
	if err := json.NewDecoder(r.Body).Decode(&forum); err != nil {
		log.Printf("create forum: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证版块名称
	if forum.Name == "" {
		log.Printf("create forum: forum name is empty")
		utils.Error(w, 400, "请填写版块名称")
		return
	}

	// 创建版块
	if err := database.DB.Create(&forum).Error; err != nil {
		log.Printf("create forum: failed to create forum, name: %s, error: %v", forum.Name, err)
		utils.Error(w, 500, "创建失败")
		return
	}

	log.Printf("create forum: forum created successfully, id: %d, name: %s", forum.ID, forum.Name)
	utils.Success(w, forum)
}

// UpdateForum 更新版块处理器（管理员）
func UpdateForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要更新的版块
	var forum models.Forum
	if err := database.DB.First(&forum, id).Error; err != nil {
		log.Printf("update forum: forum not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "版块不存在")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update forum: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "created_at")

	// 执行更新
	if err := database.DB.Model(&forum).Updates(updates).Error; err != nil {
		log.Printf("update forum: failed to update forum, id: %d, error: %v", id, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	log.Printf("update forum: forum updated successfully, id: %d", id)
	utils.Success(w, forum)
}

// DeleteForum 删除版块处理器（管理员）
func DeleteForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要删除的版块
	var forum models.Forum
	if err := database.DB.First(&forum, id).Error; err != nil {
		log.Printf("delete forum: forum not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "版块不存在")
		return
	}

	// 物理删除版块
	if err := database.DB.Unscoped().Delete(&forum).Error; err != nil {
		log.Printf("delete forum: failed to delete forum, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	log.Printf("delete forum: forum deleted successfully, id: %d", id)
	utils.Success(w, nil)
}

// ========== 话题管理 ==========

// GetAdminTopics 获取话题列表处理器（管理员）
func GetAdminTopics(w http.ResponseWriter, r *http.Request) {
	// 解析分页参数
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

	offset := (page - 1) * pageSize

	// 统计话题总数
	if err := database.DB.Model(&models.Topic{}).Count(&total).Error; err != nil {
		log.Printf("get admin topics: failed to count topics, error: %v", err)
	}

	// 查询话题列表
	if err := database.DB.Offset(offset).Limit(pageSize).Preload("User").Preload("Forum").Find(&topics).Error; err != nil {
		log.Printf("get admin topics: failed to query topics, error: %v", err)
		utils.Error(w, 500, "获取话题列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteAdminTopic 删除话题处理器（管理员）
func DeleteAdminTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要删除的话题
	var topic models.Topic
	if err := database.DB.First(&topic, id).Error; err != nil {
		log.Printf("delete admin topic: topic not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	// 物理删除话题
	if err := database.DB.Unscoped().Delete(&topic).Error; err != nil {
		log.Printf("delete admin topic: failed to delete topic, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	log.Printf("delete admin topic: topic deleted successfully, id: %d", id)
	utils.Success(w, nil)
}

// ========== 帖子管理 ==========

// GetAdminPosts 获取帖子列表处理器（管理员）
func GetAdminPosts(w http.ResponseWriter, r *http.Request) {
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

	// 统计帖子总数
	if err := database.DB.Model(&models.Post{}).Count(&total).Error; err != nil {
		log.Printf("get admin posts: failed to count posts, error: %v", err)
	}

	// 查询帖子列表
	if err := database.DB.Offset(offset).Limit(pageSize).Preload("User").Find(&posts).Error; err != nil {
		log.Printf("get admin posts: failed to query posts, error: %v", err)
		utils.Error(w, 500, "获取帖子列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteAdminPost 删除帖子处理器（管理员）
func DeleteAdminPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要删除的帖子
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		log.Printf("delete admin post: post not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "帖子不存在")
		return
	}

	// 物理删除帖子
	if err := database.DB.Unscoped().Delete(&post).Error; err != nil {
		log.Printf("delete admin post: failed to delete post, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	log.Printf("delete admin post: post deleted successfully, id: %d", id)
	utils.Success(w, nil)
}

// ========== 举报管理 ==========

// GetAdminReports 获取举报列表处理器（管理员）
func GetAdminReports(w http.ResponseWriter, r *http.Request) {
	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var reports []models.Report
	var total int64

	offset := (page - 1) * pageSize

	// 统计举报总数
	if err := database.DB.Model(&models.Report{}).Count(&total).Error; err != nil {
		log.Printf("get admin reports: failed to count reports, error: %v", err)
	}

	// 查询举报列表
	if err := database.DB.Offset(offset).Limit(pageSize).Preload("Reporter").Find(&reports).Error; err != nil {
		log.Printf("get admin reports: failed to query reports, error: %v", err)
		utils.Error(w, 500, "获取举报列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":      reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// HandleReport 处理举报处理器（管理员）
func HandleReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要处理的举报
	var report models.Report
	if err := database.DB.First(&report, id).Error; err != nil {
		log.Printf("handle report: report not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "举报不存在")
		return
	}

	// 解析请求体
	var req struct {
		Status int `json:"status"` // 处理状态：1=已处理, 2=已忽略
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("handle report: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 获取当前管理员ID
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	now := time.Now()

	// 更新举报状态
	updates := map[string]interface{}{
		"status":     req.Status,
		"handled_at": now,
		"handler_id": userID,
	}
	if err := database.DB.Model(&report).Updates(updates).Error; err != nil {
		log.Printf("handle report: failed to update report, id: %d, error: %v", id, err)
		utils.Error(w, 500, "处理举报失败")
		return
	}

	log.Printf("handle report: report handled successfully, id: %d, status: %d", id, req.Status)
	utils.Success(w, report)
}

// ========== 公告管理 ==========

// CreateAnnouncement 创建公告处理器（管理员）
func CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var announcement models.Announcement
	if err := json.NewDecoder(r.Body).Decode(&announcement); err != nil {
		log.Printf("create announcement: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证必填字段
	if announcement.Title == "" || announcement.Content == "" {
		log.Printf("create announcement: incomplete announcement info, title: %s", announcement.Title)
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	// 创建公告
	if err := database.DB.Create(&announcement).Error; err != nil {
		log.Printf("create announcement: failed to create announcement, title: %s, error: %v", announcement.Title, err)
		utils.Error(w, 500, "创建失败")
		return
	}

	log.Printf("create announcement: announcement created successfully, id: %d, title: %s", announcement.ID, announcement.Title)
	utils.Success(w, announcement)
}

// UpdateAnnouncement 更新公告处理器（管理员）
func UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要更新的公告
	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		log.Printf("update announcement: announcement not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "公告不存在")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update announcement: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 执行更新
	if err := database.DB.Model(&announcement).Updates(updates).Error; err != nil {
		log.Printf("update announcement: failed to update announcement, id: %d, error: %v", id, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	log.Printf("update announcement: announcement updated successfully, id: %d", id)
	utils.Success(w, announcement)
}

// DeleteAnnouncement 删除公告处理器（管理员）
func DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 查询要删除的公告
	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		log.Printf("delete announcement: announcement not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "公告不存在")
		return
	}

	// 物理删除公告
	if err := database.DB.Unscoped().Delete(&announcement).Error; err != nil {
		log.Printf("delete announcement: failed to delete announcement, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	log.Printf("delete announcement: announcement deleted successfully, id: %d", id)
	utils.Success(w, nil)
}

// ========== 标签管理 ==========

// GetAdminTags 获取标签列表处理器（管理员）
func GetAdminTags(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("get admin tags: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var tags []models.Tag
	if err := database.DB.Order("usage_count DESC, sort_order ASC").Find(&tags).Error; err != nil {
		log.Printf("get admin tags: failed to query tags, error: %v", err)
		utils.Error(w, 500, "获取标签列表失败")
		return
	}

	utils.Success(w, tags)
}

// CreateTag 创建标签处理器（管理员）
func CreateTag(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("create tag: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var req struct {
		Name       string `json:"name"`        // 标签名称
		Icon       string `json:"icon"`        // 标签图标
		IsOfficial bool   `json:"is_official"` // 是否官方标签
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create tag: failed to decode request body, error: %v", err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 验证标签名称
	if req.Name == "" {
		log.Printf("create tag: tag name is empty")
		utils.Error(w, http.StatusBadRequest, "标签名称不能为空")
		return
	}

	// 创建标签
	tag := models.Tag{
		Name:       req.Name,
		Icon:       req.Icon,
		IsOfficial: req.IsOfficial,
	}
	if err := database.DB.Create(&tag).Error; err != nil {
		log.Printf("create tag: failed to create tag, name: %s, error: %v", req.Name, err)
		utils.Error(w, 500, "创建标签失败")
		return
	}

	log.Printf("create tag: tag created successfully, id: %d, name: %s", tag.ID, req.Name)
	utils.Success(w, tag)
}

// UpdateTag 更新标签处理器（管理员）
func UpdateTag(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("update tag: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 解析请求体
	var req struct {
		Name       string `json:"name"`        // 标签名称
		Icon       string `json:"icon"`        // 标签图标
		IsOfficial bool   `json:"is_official"` // 是否官方标签
		IsBanned   bool   `json:"is_banned"`   // 是否禁用
		SortOrder  int    `json:"sort_order"`  // 排序顺序
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update tag: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 查询要更新的标签
	var tag models.Tag
	if result := database.DB.First(&tag, id); result.Error != nil {
		log.Printf("update tag: tag not found, id: %d, error: %v", id, result.Error)
		utils.Error(w, http.StatusNotFound, "标签不存在")
		return
	}

	// 更新字段
	tag.Name = req.Name
	tag.Icon = req.Icon
	tag.IsOfficial = req.IsOfficial
	tag.IsBanned = req.IsBanned
	tag.SortOrder = req.SortOrder

	// 保存更新
	if err := database.DB.Save(&tag).Error; err != nil {
		log.Printf("update tag: failed to update tag, id: %d, error: %v", id, err)
		utils.Error(w, 500, "更新标签失败")
		return
	}

	log.Printf("update tag: tag updated successfully, id: %d, name: %s", id, req.Name)
	utils.Success(w, tag)
}

// DeleteTag 删除标签处理器（管理员）
func DeleteTag(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("delete tag: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 物理删除标签
	if err := database.DB.Unscoped().Delete(&models.Tag{}, id).Error; err != nil {
		log.Printf("delete tag: failed to delete tag, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除标签失败")
		return
	}

	log.Printf("delete tag: tag deleted successfully, id: %d", id)
	utils.Success(w, nil)
}

// MergeTags 合并标签处理器（管理员）
// 将源标签合并到目标标签，更新关联关系并删除源标签
func MergeTags(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("merge tags: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var req struct {
		SourceID uint `json:"source_id"` // 源标签ID
		TargetID uint `json:"target_id"` // 目标标签ID
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("merge tags: failed to decode request body, error: %v", err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 验证参数
	if req.SourceID == 0 || req.TargetID == 0 || req.SourceID == req.TargetID {
		log.Printf("merge tags: invalid tag ids, sourceID: %d, targetID: %d", req.SourceID, req.TargetID)
		utils.Error(w, http.StatusBadRequest, "无效的标签ID")
		return
	}

	// 查询源标签
	var sourceTag, targetTag models.Tag
	if err := database.DB.First(&sourceTag, req.SourceID).Error; err != nil {
		log.Printf("merge tags: source tag not found, sourceID: %d, error: %v", req.SourceID, err)
		utils.Error(w, http.StatusNotFound, "源标签不存在")
		return
	}
	if err := database.DB.First(&targetTag, req.TargetID).Error; err != nil {
		log.Printf("merge tags: target tag not found, targetID: %d, error: %v", req.TargetID, err)
		utils.Error(w, http.StatusNotFound, "目标标签不存在")
		return
	}

	// 查找使用源标签的话题
	var topics []models.Topic
	database.DB.Model(&sourceTag).Association("Topics").Find(&topics)

	// 将话题关联到目标标签
	if len(topics) > 0 {
		if err := database.DB.Model(&targetTag).Association("Topics").Append(topics).Error; err != nil {
			log.Printf("merge tags: failed to append topics to target tag, error: %v", err)
		}
		// 清除源标签的关联
		if err := database.DB.Model(&sourceTag).Association("Topics").Clear().Error; err != nil {
			log.Printf("merge tags: failed to clear source tag associations, error: %v", err)
		}
	}

	// 更新目标标签的使用次数
	targetTag.UsageCount += sourceTag.UsageCount
	if err := database.DB.Save(&targetTag).Error; err != nil {
		log.Printf("merge tags: failed to update target tag usage count, error: %v", err)
	}

	// 删除源标签
	if err := database.DB.Unscoped().Delete(&sourceTag).Error; err != nil {
		log.Printf("merge tags: failed to delete source tag, error: %v", err)
	}

	log.Printf("merge tags: tags merged successfully, sourceID: %d, targetID: %d, mergedCount: %d", req.SourceID, req.TargetID, len(topics))
	utils.Success(w, map[string]interface{}{
		"merged_count": len(topics),
		"target_tag":   targetTag,
	})
}

// ========== 管理员密码 ==========

// ChangeAdminPassword 修改管理员密码处理器
func ChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
	// 获取当前管理员ID
	adminID, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("change admin password: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var req struct {
		OldPassword string `json:"old_password"` // 原密码
		NewPassword string `json:"new_password"` // 新密码
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("change admin password: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.OldPassword == "" || req.NewPassword == "" {
		log.Printf("change admin password: incomplete password info")
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	// 验证新密码长度
	if len(req.NewPassword) < 6 {
		log.Printf("change admin password: new password too short, length: %d", len(req.NewPassword))
		utils.Error(w, 400, "新密码长度至少为6位")
		return
	}

	// 查询管理员信息
	var admin models.User
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		log.Printf("change admin password: admin not found, id: %d, error: %v", adminID, err)
		utils.Error(w, 404, "管理员不存在")
		return
	}

	// 验证原密码
	if !utils.CheckPassword(req.OldPassword, admin.PasswordHash) {
		log.Printf("change admin password: old password mismatch, adminID: %d", adminID)
		utils.Error(w, 400, "原密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("change admin password: failed to hash new password, error: %v", err)
		utils.Error(w, 500, "密码加密失败")
		return
	}

	// 更新密码
	if err := database.DB.Model(&admin).Update("password_hash", hashedPassword).Error; err != nil {
		log.Printf("change admin password: failed to update password, adminID: %d, error: %v", adminID, err)
		utils.Error(w, 500, "更新密码失败")
		return
	}

	log.Printf("change admin password: password changed successfully, adminID: %d", adminID)
	utils.Success(w, nil)
}
