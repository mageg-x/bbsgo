package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/models"
	"log"
	"net/http"
	"time"
)

// GetAnnouncements 获取公告列表处理器
// 支持缓存，缓存时间5分钟
// 只返回未过期的公告
func GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	// 尝试从缓存获取
	if cached, ok := cache.Get("announcements:list"); ok {
		errors.Success(w, cached)
		return
	}

	now := time.Now()

	// 查询未过期的公告
	var announcements []models.Announcement
	if err := database.DB.Where("(expires_at IS NULL OR expires_at > ?)", now).
		Order("is_pinned DESC, created_at DESC").
		Find(&announcements).Error; err != nil {
		log.Printf("get announcements: failed to query announcements, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 设置缓存
	cache.Set("announcements:list", announcements, 5*time.Minute)

	errors.Success(w, announcements)
}
