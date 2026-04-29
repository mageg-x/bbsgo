package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
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

	// 获取当前查看者信息（用于匿名判断）
	viewerID, viewerRole, _ := middleware.GetOptionalUserInfo(r)

	// 处理匿名数据
	type ProcessedTopic struct {
		ID               uint                 `json:"id"`
		Title            string               `json:"title"`
		Content          string               `json:"content"`
		UserID           uint                 `json:"user_id"`
		User             interface{}          `json:"user"`
		ForumID          uint                 `json:"forum_id"`
		Forum            models.Forum         `json:"forum"`
		IsPinned         bool                 `json:"is_pinned"`
		IsUserPinned     bool                 `json:"is_user_pinned"`
		IsLocked         bool                 `json:"is_locked"`
		IsEssence        bool                 `json:"is_essence"`
		IsHidden         bool                 `json:"is_hidden"`
		HotScore         float64              `json:"hot_score"`
		LikeCount        int                  `json:"like_count"`
		ViewCount        int                  `json:"view_count"`
		ReplyCount       int                  `json:"reply_count"`
		LastReplyAt      *time.Time           `json:"last_reply_at"`
		AllowComment     bool                 `json:"allow_comment"`
		CreatedAt        time.Time            `json:"created_at"`
		UpdatedAt        time.Time            `json:"updated_at"`
		IsAnonymous      bool                 `json:"is_anonymous"`
		AnonymousType    models.AnonymousType `json:"anonymous_type"`
		AnonymousUntil   *time.Time           `json:"anonymous_until"`
		IsAnonymousEnded bool                 `json:"is_anonymous_ended"`
		Tags             []models.Tag         `json:"tags"`
	}

	var processedList []ProcessedTopic
	for _, t := range topics {
		shouldShowAnonymous := utils.ShouldShowAnonymous(viewerID, t.UserID, viewerRole, &t)

		item := ProcessedTopic{
			ID:               t.ID,
			Title:            t.Title,
			Content:          t.Content,
			UserID:           t.UserID,
			ForumID:          t.ForumID,
			Forum:            t.Forum,
			IsPinned:         t.IsPinned,
			IsUserPinned:     t.IsUserPinned,
			IsLocked:         t.IsLocked,
			IsEssence:        t.IsEssence,
			IsHidden:         t.IsHidden,
			HotScore:         t.HotScore,
			LikeCount:        t.LikeCount,
			ViewCount:        t.ViewCount,
			ReplyCount:       t.ReplyCount,
			LastReplyAt:      t.LastReplyAt,
			AllowComment:     t.AllowComment,
			CreatedAt:        t.CreatedAt,
			UpdatedAt:        t.UpdatedAt,
			IsAnonymous:      t.IsAnonymous,
			AnonymousType:    t.AnonymousType,
			AnonymousUntil:   t.AnonymousUntil,
			IsAnonymousEnded: t.IsAnonymousEnded,
			Tags:             t.Tags,
		}

		if shouldShowAnonymous {
			item.User = utils.GetAnonymousUserInfo()
			item.UserID = 0
		} else {
			item.User = t.User
		}

		processedList = append(processedList, item)
	}

	log.Printf("search: search completed, keyword: %s, results: %d", keyword, total)
	errors.Success(w, map[string]interface{}{
		"list":      processedList,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"keyword":   keyword,
	})
}
