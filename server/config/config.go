package config

import (
	"bbsgo/database"
	"bbsgo/models"
	"log"
	"strconv"
)

// GetConfig 获取配置项的值
// key: 配置键名
// 返回: 配置值字符串，如果不存在则返回空字符串
func GetConfig(key string) string {
	var config models.SiteConfig
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		log.Printf("get config failed, key: %s, error: %v", key, err)
		return ""
	}
	return config.Value
}

// GetConfigInt 获取配置项的值作为整数
// key: 配置键名
// defaultValue: 默认值，当配置不存在或解析失败时使用
// 返回: 配置值整数
func GetConfigInt(key string, defaultValue int) int {
	val := GetConfig(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("parse config int failed, key: %s, value: %s, error: %v", key, val, err)
		return defaultValue
	}
	return i
}

// GetConfigBool 获取配置项的值作为布尔值
// key: 配置键名
// defaultValue: 默认值
// 返回: 配置值布尔值
func GetConfigBool(key string, defaultValue bool) bool {
	val := GetConfig(key)
	if val == "" {
		return defaultValue
	}
	return val == "true" || val == "1"
}

// SetConfig 设置配置项的值
// key: 配置键名
// value: 配置值
// 返回: 错误信息
func SetConfig(key, value string) error {
	var config models.SiteConfig
	// 尝试查找现有配置
	if err := database.DB.Where("key = ?", key).First(&config).Error; err != nil {
		// 配置不存在，创建新记录
		config = models.SiteConfig{
			Key:   key,
			Value: value,
		}
		if err := database.DB.Create(&config).Error; err != nil {
			log.Printf("create config failed, key: %s, value: %s, error: %v", key, value, err)
			return err
		}
		return nil
	}
	// 配置存在，更新值
	if err := database.DB.Model(&config).Update("value", value).Error; err != nil {
		log.Printf("update config failed, key: %s, value: %s, error: %v", key, value, err)
		return err
	}
	return nil
}
