package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/models"
	"log"
	"net/http"
	"strconv"
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

	var topics []models.Topic
	var total int64

	offset := (page - 1) * pageSize
	searchPattern := "%" + keyword + "%"

	// 统计匹配的话题数量
	if err := database.DB.Model(&models.Topic{}).Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).Count(&total).Error; err != nil {
		log.Printf("search: failed to count topics, keyword: %s, error: %v", keyword, err)
	}

	// 搜索话题
	if err := database.DB.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Preload("User").
		Preload("Forum").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&topics).Error; err != nil {
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
