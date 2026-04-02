package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
)

// GetSiteConfig 获取网站配置处理器
// 返回所有配置项的键值对
func GetSiteConfig(w http.ResponseWriter, r *http.Request) {
	var configs []models.SiteConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		log.Printf("get site config: failed to query configs, error: %v", err)
		utils.Error(w, 500, "获取配置失败")
		return
	}

	// 转换为 map 返回
	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	utils.Success(w, configMap)
}

// UpdateSiteConfig 更新网站配置处理器
// 批量更新多个配置项，需要管理员权限
func UpdateSiteConfig(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("update site config: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update site config: failed to decode request body, error: %v", err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
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
	utils.Success(w, nil)
}
