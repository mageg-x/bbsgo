package utils

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"strconv"
	"time"
)

// GetConfigString 获取字符串类型的配置值
func GetConfigString(key, defaultValue string) string {
	log.Printf("GetConfigString: [DEBUG] start, key=%s, defaultValue=%s", key, defaultValue)
	cacheKey := "config:" + key
	if cached, ok := cache.Get(cacheKey); ok {
		if val, ok := cached.(string); ok {
			log.Printf("GetConfigString: [DEBUG] cache hit, key=%s, value=%s", key, val)
			return val
		}
	}

	if database.DB == nil {
		log.Printf("GetConfigString: [DEBUG] database.DB is nil, returning default")
		return defaultValue
	}

	var config models.SiteConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err == nil {
		cache.Set(cacheKey, config.Value, 5*time.Minute)
		log.Printf("GetConfigString: [DEBUG] db hit, key=%s, value=%s", key, config.Value)
		return config.Value
	}

	log.Printf("GetConfigString: [DEBUG] not found, key=%s, returning default=%s", key, defaultValue)
	return defaultValue
}

// GetConfigBool 获取布尔类型的配置值
func GetConfigBool(key string, defaultValue bool) bool {
	strVal := GetConfigString(key, "")
	if strVal == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

// GetConfigInt 获取整数类型的配置值
func GetConfigInt(key string, defaultValue int) int {
	strVal := GetConfigString(key, "")
	if strVal == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		log.Printf("get config int: failed to parse %s, error: %v", key, err)
		return defaultValue
	}
	return intVal
}

// InvalidateConfigCache 使配置缓存失效
func InvalidateConfigCache() {
	// 这里可以清空所有 config 前缀的缓存
}
