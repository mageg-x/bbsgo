package antispam

import (
	"bbsgo/database"
	"bbsgo/models"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"
)

// ConfigService 配置服务
// 负责管理系统防刷相关配置
type ConfigService struct {
	cache map[string]string
	mu    sync.RWMutex
}

var configService *ConfigService
var configOnce sync.Once

// GetConfigService 获取配置服务单例
func GetConfigService() *ConfigService {
	configOnce.Do(func() {
		log.Printf("[config] initializing config service singleton")
		configService = &ConfigService{
			cache: make(map[string]string),
		}
		configService.LoadFromDB()
	})
	return configService
}

// LoadFromDB 从数据库加载配置到缓存
func (s *ConfigService) LoadFromDB() {
	log.Printf("[config] loading config from database")

	var configs []models.AntiSpamConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		log.Printf("[config] load failed: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, config := range configs {
		s.cache[config.Key] = config.Value
	}

	log.Printf("[config] loaded %d configs from database", len(configs))
}

// Get 获取字符串配置值
func (s *ConfigService) Get(key string, defaultValue string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.cache[key]; ok {
		return val
	}
	log.Printf("[config] key not found: %s, using default: %s", key, defaultValue)
	return defaultValue
}

// GetInt 获取整数配置值
func (s *ConfigService) GetInt(key string, defaultValue int) int {
	val := s.Get(key, "")
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("[config] parse int failed, key: %s, value: %s, error: %v, using default: %d",
			key, val, err, defaultValue)
		return defaultValue
	}
	return intVal
}

// GetFloat 获取浮点数配置值
func (s *ConfigService) GetFloat(key string, defaultValue float64) float64 {
	val := s.Get(key, "")
	if val == "" {
		return defaultValue
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Printf("[config] parse float failed, key: %s, value: %s, error: %v, using default: %.2f",
			key, val, err, defaultValue)
		return defaultValue
	}
	return floatVal
}

// GetBool 获取布尔配置值
func (s *ConfigService) GetBool(key string, defaultValue bool) bool {
	val := s.Get(key, "")
	if val == "" {
		return defaultValue
	}

	boolVal := val == "true"
	log.Printf("[config] get bool, key: %s, value: %s, result: %v", key, val, boolVal)
	return boolVal
}

// GetStringSlice 获取字符串数组配置值
func (s *ConfigService) GetStringSlice(key string, defaultValue []string) []string {
	val := s.Get(key, "")
	if val == "" {
		return defaultValue
	}

	var result []string
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		log.Printf("[config] parse slice failed, key: %s, value: %s, error: %v, using default",
			key, val, err)
		return defaultValue
	}
	return result
}

// Set 设置配置值
func (s *ConfigService) Set(key string, value string) error {
	log.Printf("[config] setting, key: %s, value: %s", key, value)

	s.mu.Lock()
	defer s.mu.Unlock()

	config := models.AntiSpamConfig{
		Key:   key,
		Value: value,
	}

	result := database.DB.Where("key = ?", key).Assign(config).FirstOrCreate(&config)
	if result.Error != nil {
		log.Printf("[config] save failed, key: %s, value: %s, error: %v", key, value, result.Error)
		return result.Error
	}

	s.cache[key] = value
	log.Printf("[config] saved, key: %s, value: %s", key, value)
	return nil
}

// GetAll 获取所有配置项
func (s *ConfigService) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]string)
	for k, v := range s.cache {
		result[k] = v
	}
	return result
}

// 配置项常量定义
const (
	ConfigTopicMinInterval           = "topic_min_interval"
	ConfigCommentMinInterval         = "comment_min_interval"
	ConfigMaxTopicsPerDay            = "max_topics_per_day"
	ConfigMaxCommentsPerDay          = "max_comments_per_day"
	ConfigNewUserMaxTopicsPerDay     = "new_user_max_topics_per_day"
	ConfigNewUserMaxCommentsPerDay   = "new_user_max_comments_per_day"
	ConfigNewUserHours               = "new_user_hours"
	ConfigMinContentLength           = "min_content_length"
	ConfigSimilarityThreshold        = "similarity_threshold"
	ConfigRepeatCharThreshold        = "repeat_char_threshold"
	ConfigReportThreshold            = "report_threshold"
	ConfigReportBanThreshold         = "report_ban_threshold"
	ConfigReportBanDays             = "report_ban_days"
	ConfigMaxReportsPerDay           = "max_reports_per_day"
	ConfigLowQualityHotMultiplier    = "low_quality_hot_multiplier"
	ConfigLowReputationHotMultiplier = "low_reputation_hot_multiplier"
	ConfigLowReputationThreshold     = "low_reputation_threshold"
	ConfigBanLowReputation           = "ban_low_reputation"
	ConfigBanReputationThreshold     = "ban_reputation_threshold"
	ConfigSpamKeywords             = "spam_keywords"
)

// GetDefaultConfigs 获取默认配置项
func (s *ConfigService) GetDefaultConfigs() map[string]string {
	return map[string]string{
		ConfigTopicMinInterval:           "60",
		ConfigCommentMinInterval:         "30",
		ConfigMaxTopicsPerDay:            "10",
		ConfigMaxCommentsPerDay:          "50",
		ConfigNewUserMaxTopicsPerDay:     "3",
		ConfigNewUserMaxCommentsPerDay:   "10",
		ConfigNewUserHours:               "24",
		ConfigMinContentLength:           "10",
		ConfigSimilarityThreshold:        "0.8",
		ConfigRepeatCharThreshold:        "5",
		ConfigReportThreshold:            "3",
		ConfigReportBanThreshold:         "5",
		ConfigReportBanDays:             "3",
		ConfigMaxReportsPerDay:           "10",
		ConfigLowQualityHotMultiplier:    "0.3",
		ConfigLowReputationHotMultiplier: "0.5",
		ConfigLowReputationThreshold:     "60",
		ConfigBanLowReputation:           "true",
		ConfigBanReputationThreshold:     "20",
	}
}

// InitializeDefaults 初始化默认配置
func (s *ConfigService) InitializeDefaults() {
	log.Printf("[config] initializing defaults")

	defaults := s.GetDefaultConfigs()
	createdCount := 0

	for key, value := range defaults {
		var config models.AntiSpamConfig
		result := database.DB.Where("key = ?", key).First(&config)

		if result.Error != nil {
			config = models.AntiSpamConfig{
				Key:   key,
				Value: value,
			}
			if err := database.DB.Create(&config).Error; err != nil {
				log.Printf("[config] create default failed, key: %s, value: %s, error: %v",
					key, value, err)
			} else {
				createdCount++
				log.Printf("[config] created default, key: %s, value: %s", key, value)
			}
		}
	}

	s.LoadFromDB()
	log.Printf("[config] initialized, created %d new configs", createdCount)
}

func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if configService != nil {
				log.Printf("[config] reloading config periodically")
				configService.LoadFromDB()
			}
		}
	}()
}
