package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Search 搜索话题处理器
// 根据关键词搜索话题标题和内容
func Search(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		log.Printf("search: keyword is empty")
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取当前用户ID（用于私密帖子过滤）
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var topics []models.Topic
	var total int64

	offset := (page - 1) * pageSize
	searchPattern := "%" + keyword + "%"

	// 构建查询条件：过滤活跃的私密帖子
	now := time.Now()
	privateCondition := "(is_private = ? OR private_ended = ? OR (is_private = ? AND (private_expire_at IS NULL OR private_expire_at < ?)))"
	if userID > 0 {
		// 已登录用户可以搜索到自己的私密帖子
		privateCondition = fmt.Sprintf("(%s OR user_id = %d)", privateCondition, userID)
	}

	// 统计匹配的话题数量（排除活跃的私密帖子）
	countQuery := database.DB.Model(&models.Topic{}).
		Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Where(privateCondition, false, true, true, now)
	if err := countQuery.Count(&total).Error; err != nil {
		log.Printf("search: failed to count topics, keyword: %s, error: %v", keyword, err)
	}

	// 搜索话题（排除活跃的私密帖子）
	searchQuery := database.DB.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Where(privateCondition, false, true, true, now).
		Preload("User").
		Preload("Forum").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize)
	if err := searchQuery.Find(&topics).Error; err != nil {
		log.Printf("search: failed to search topics, keyword: %s, error: %v", keyword, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("search: search completed, keyword: %s, results: %d", keyword, total)
	errors.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"keyword":   keyword,
	})
}
