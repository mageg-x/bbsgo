package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"encoding/json"
	"log"
	"net/http"
)

// 清除配置相关的缓存
func clearConfigCache() {
	keys := []string{"config:allow_register", "config:allow_post", "config:allow_comment",
		"config:credit_topic", "config:credit_post", "config:credit_signin"}
	for _, key := range keys {
		cache.Delete(key)
	}
}

// GetSiteConfig 获取网站配置处理器
// 返回所有配置项的键值对
func GetSiteConfig(w http.ResponseWriter, r *http.Request) {
	var configs []models.SiteConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		log.Printf("get site config: failed to query configs, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 转换为 map 返回
	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	errors.Success(w, configMap)
}

// UpdateSiteConfig 更新网站配置处理器
// 批量更新多个配置项，需要管理员权限
func UpdateSiteConfig(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("update site config: unauthorized access")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	// 解析请求体
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update site config: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 遍历更新每个配置项
	for key, value := range req {
		var config models.SiteConfig
		result := database.DB.Where("key = ?", key).First(&config)
		if result.Error == nil {
			// 配置存在，更新值
			config.Value = value
			if err := database.DB.Save(&config).Error; err != nil {
				log.Printf("update site config: failed to update config, key: %s, value: %s, error: %v", key, value, err)
			}
		} else {
			// 配置不存在，创建新记录
			newConfig := models.SiteConfig{
				Key:   key,
				Value: value,
			}
			if err := database.DB.Create(&newConfig).Error; err != nil {
				log.Printf("update site config: failed to create config, key: %s, value: %s, error: %v", key, value, err)
			}
		}
	}

	log.Printf("update site config: updated %d configs", len(req))
	clearConfigCache()
	errors.Success(w, nil)
}
