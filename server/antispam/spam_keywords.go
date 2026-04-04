package antispam

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
)

// SpamKeywordService 敏感词服务
// 从数据库加载和管理敏感词库
type SpamKeywordService struct {
	keywords []string
	mu       sync.RWMutex
}

var spamKeywordService *SpamKeywordService
var spamKeywordOnce sync.Once

// GetSpamKeywordService 获取敏感词服务单例
func GetSpamKeywordService() *SpamKeywordService {
	spamKeywordOnce.Do(func() {
		log.Printf("[spam-keyword] initializing spam keyword service")
		spamKeywordService = &SpamKeywordService{
			keywords: getDefaultSpamKeywords(),
		}
		spamKeywordService.LoadFromDB()
	})
	return spamKeywordService
}

// getDefaultSpamKeywords 获取默认敏感词列表
func getDefaultSpamKeywords() []string {
	return []string{
		"加VX", "加微信", "加V", "加Q", "加QQ",
		"赚钱", "日赚", "月入", "兼职赚钱",
		"代刷", "代写", "代做",
		"免费领取", "点击领取", "扫码领取",
		"低价出售", "高价收购",
	}
}

// LoadFromDB 从数据库加载敏感词
func (s *SpamKeywordService) LoadFromDB() {
	configService := GetConfigService()
	keywordsStr := configService.Get("spam_keywords", "")

	if keywordsStr == "" {
		log.Printf("[spam-keyword] no keywords in db, using defaults (%d words)", len(s.keywords))
		return
	}

	var keywords []string
	if err := json.Unmarshal([]byte(keywordsStr), &keywords); err != nil {
		log.Printf("[spam-keyword] parse keywords failed: %v, using defaults", err)
		return
	}

	s.mu.Lock()
	s.keywords = keywords
	s.mu.Unlock()

	log.Printf("[spam-keyword] loaded %d keywords from db", len(keywords))
}

// GetKeywords 获取所有敏感词
func (s *SpamKeywordService) GetKeywords() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, len(s.keywords))
	copy(result, s.keywords)
	return result
}

// Check 检查内容是否包含敏感词
// 返回敏感词列表和是否包含敏感词
func (s *SpamKeywordService) Check(content string) (bool, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lowerContent := strings.ToLower(content)
	var matched []string

	for _, keyword := range s.keywords {
		if strings.Contains(lowerContent, strings.ToLower(keyword)) {
			matched = append(matched, keyword)
		}
	}

	return len(matched) > 0, matched
}

// AddKeyword 添加敏感词
func (s *SpamKeywordService) AddKeyword(keyword string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已存在
	for _, k := range s.keywords {
		if strings.ToLower(k) == strings.ToLower(keyword) {
			return
		}
	}

	s.keywords = append(s.keywords, keyword)
	s.saveToDB()
	log.Printf("[spam-keyword] added keyword: %s, total: %d", keyword, len(s.keywords))
}

// RemoveKeyword 删除敏感词
func (s *SpamKeywordService) RemoveKeyword(keyword string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, k := range s.keywords {
		if strings.ToLower(k) == strings.ToLower(keyword) {
			s.keywords = append(s.keywords[:i], s.keywords[i+1:]...)
			s.saveToDB()
			log.Printf("[spam-keyword] removed keyword: %s, total: %d", keyword, len(s.keywords))
			return
		}
	}
}

// saveToDB 保存敏感词到数据库
func (s *SpamKeywordService) saveToDB() {
	keywordsStr, err := json.Marshal(s.keywords)
	if err != nil {
		log.Printf("[spam-keyword] marshal keywords failed: %v", err)
		return
	}

	configService := GetConfigService()
	if err := configService.Set("spam_keywords", string(keywordsStr)); err != nil {
		log.Printf("[spam-keyword] save to db failed: %v", err)
	}
}

// SetKeywords 设置敏感词列表（批量替换）
func (s *SpamKeywordService) SetKeywords(keywords []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.keywords = keywords
	s.saveToDB()
	log.Printf("[spam-keyword] set %d keywords", len(keywords))
}
