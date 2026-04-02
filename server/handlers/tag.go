package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// GetTags 获取热门标签列表处理器
// 支持缓存，缓存时间10分钟
// 返回未被禁用的标签，按官方标签优先、使用次数排序
func GetTags(w http.ResponseWriter, r *http.Request) {
	// 尝试从缓存获取
	if cached, ok := cache.Get("tags:hot"); ok {
		utils.Success(w, cached)
		return
	}

	// 查询未禁用的标签
	var tags []models.Tag
	if err := database.DB.Where("is_banned = ?", false).
		Order("is_official DESC, usage_count DESC, sort_order ASC").
		Limit(20).
		Find(&tags).Error; err != nil {
		log.Printf("get tags: failed to query tags, error: %v", err)
		utils.Error(w, 500, "获取标签失败")
		return
	}

	// 设置缓存
	cache.Set("tags:hot", tags, 10*time.Minute)

	utils.Success(w, tags)
}

// GetTag 获取单个标签详情处理器
func GetTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var tag models.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		log.Printf("get tag: tag not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "标签不存在")
		return
	}

	utils.Success(w, tag)
}

// SearchTags 搜索标签处理器
// 根据关键词模糊搜索标签
func SearchTags(w http.ResponseWriter, r *http.Request) {
	keyword := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(keyword) < 1 {
		utils.Success(w, []models.Tag{})
		return
	}

	var tags []models.Tag
	if err := database.DB.Where("is_banned = ? AND name LIKE ?", false, "%"+keyword+"%").
		Order("usage_count DESC").
		Limit(10).
		Find(&tags).Error; err != nil {
		log.Printf("search tags: failed to search tags, keyword: %s, error: %v", keyword, err)
		utils.Error(w, 500, "搜索标签失败")
		return
	}

	utils.Success(w, tags)
}

// GetOrCreateTagByName 根据名称获取或创建标签
// name: 标签名称
// 返回: 标签指针和错误
func GetOrCreateTagByName(name string) (*models.Tag, error) {
	name = strings.TrimSpace(name)
	// 验证标签名称长度
	if len(name) < 2 || len(name) > 20 {
		return nil, nil
	}

	var tag models.Tag
	if err := database.DB.Where("name = ?", name).First(&tag).Error; err != nil {
		// 标签不存在，创建新标签
		tag = models.Tag{
			Name:       name,
			UsageCount: 0,
			IsOfficial: false,
			IsBanned:   false,
		}
		if err := database.DB.Create(&tag).Error; err != nil {
			log.Printf("get or create tag: failed to create tag, name: %s, error: %v", name, err)
			return nil, err
		}
	}
	return &tag, nil
}

// IncrementTagUsage 增加标签使用次数
// tagID: 标签ID
func IncrementTagUsage(tagID uint) {
	if err := database.DB.Model(&models.Tag{}).Where("id = ?", tagID).
		UpdateColumn("usage_count", database.DB.Raw("usage_count + 1")).Error; err != nil {
		log.Printf("increment tag usage: failed to increment usage count, tagID: %d, error: %v", tagID, err)
	}
}
