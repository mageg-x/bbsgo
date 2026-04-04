package antispam

import (
	"bbsgo/database"
	"bbsgo/models"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// ContentQualityService 内容质量检测服务
// 负责检测用户发布内容的质量
type ContentQualityService struct {
	config *ConfigService
}

// NewContentQualityService 创建内容质量检测服务实例
func NewContentQualityService() *ContentQualityService {
	return &ContentQualityService{
		config: GetConfigService(),
	}
}

// QualityCheckResult 内容质量检查结果
type QualityCheckResult struct {
	IsLowQuality bool
	Reason       string
	Score        float64
}

// CheckContent 检测内容质量
func (s *ContentQualityService) CheckContent(content string) *QualityCheckResult {
	result := &QualityCheckResult{
		IsLowQuality: false,
		Score:        100.0,
	}

	log.Printf("[quality] checking content, length: %d", len(content))

	if s.isEmptyContent(content) {
		result.IsLowQuality = true
		result.Reason = "内容无实质信息"
		result.Score = 0
		log.Printf("[quality] result: empty content, score: %.2f", result.Score)
		return result
	}

	if s.isOnlySymbols(content) {
		result.IsLowQuality = true
		result.Reason = "内容仅包含符号或表情"
		result.Score = 0
		log.Printf("[quality] result: symbols only, score: %.2f", result.Score)
		return result
	}

	if s.isTooShort(content) {
		result.IsLowQuality = true
		minLen := s.config.GetInt(ConfigMinContentLength, 10)
		result.Reason = fmt.Sprintf("内容太短，最少需要%d个字符", minLen)
		result.Score = 20
		log.Printf("[quality] result: too short, minLength: %d, score: %.2f", minLen, result.Score)
		return result
	}

	if s.hasRepeatingChars(content) {
		result.IsLowQuality = true
		result.Reason = "内容包含大量无意义重复字符"
		result.Score = 10
		log.Printf("[quality] result: repeating chars, score: %.2f", result.Score)
		return result
	}

	if s.hasSpamKeywords(content) {
		result.IsLowQuality = true
		result.Reason = "内容包含广告特征"
		result.Score = 0
		log.Printf("[quality] result: spam keywords, score: %.2f", result.Score)
		return result
	}

	linkCount := s.countExternalLinks(content)
	if linkCount >= 3 {
		result.IsLowQuality = true
		result.Reason = "内容包含过多外部链接"
		result.Score = 40
		log.Printf("[quality] result: too many external links (%d), score: %.2f", linkCount, result.Score)
		return result
	}

	log.Printf("[quality] result: passed, score: %.2f", result.Score)
	return result
}

func (s *ContentQualityService) isEmptyContent(content string) bool {
	return strings.TrimSpace(content) == ""
}

func (s *ContentQualityService) isOnlySymbols(content string) bool {
	hasLetterOrDigit := false
	for _, r := range content {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			hasLetterOrDigit = true
			break
		}
	}
	return !hasLetterOrDigit
}

func (s *ContentQualityService) isTooShort(content string) bool {
	minLength := s.config.GetInt(ConfigMinContentLength, 10)

	codeBlockRegex := regexp.MustCompile("```[\\s\\S]*?```")
	contentWithoutCode := codeBlockRegex.ReplaceAllString(content, "")

	linkRegex := regexp.MustCompile(`https?://[^\s]+`)
	contentWithoutLinks := linkRegex.ReplaceAllString(contentWithoutCode, "")

	plainText := strings.TrimSpace(contentWithoutLinks)

	// 统计有效字符数量（汉字、英文字母、标点都算）
	validCount := 0
	for _, r := range plainText {
		if unicode.IsLetter(r) || unicode.Is(unicode.Han, r) || unicode.IsPunct(r) {
			validCount++
		}
	}

	isTooShort := validCount < minLength
	log.Printf("[quality-length] valid chars: %d, min: %d, too_short: %v",
		validCount, minLength, isTooShort)

	return isTooShort
}

func (s *ContentQualityService) hasRepeatingChars(content string) bool {
	// 检测连续重复字符（5个或以上相同字符）
	if len(content) >= 5 {
		count := 1
		prev := rune(content[0])
		for i := 1; i < len(content); i++ {
			if rune(content[i]) == prev {
				count++
				if count >= 5 {
					log.Printf("[quality-repeating] detected repeating char, count: %d", count)
					return true
				}
			} else {
				count = 1
				prev = rune(content[i])
			}
		}
	}
	return false
}

func (s *ContentQualityService) hasSpamKeywords(content string) bool {
	spamService := GetSpamKeywordService()
	hasSpam, matchedKeywords := spamService.Check(content)
	if hasSpam {
		log.Printf("[quality-keyword] detected spam keywords: %v", matchedKeywords)
	}
	return hasSpam
}

func (s *ContentQualityService) countExternalLinks(content string) int {
	linkRegex := regexp.MustCompile(`https?://[^\s]+`)
	links := linkRegex.FindAllString(content, -1)

	whitelistDomains := []string{
		"localhost",
		"127.0.0.1",
		"bbsgo.com",
	}

	externalCount := 0
	for _, link := range links {
		isWhitelisted := false
		for _, domain := range whitelistDomains {
			if strings.Contains(link, domain) {
				isWhitelisted = true
				break
			}
		}
		if !isWhitelisted {
			externalCount++
		}
	}

	log.Printf("[quality-links] total: %d, external: %d", len(links), externalCount)
	return externalCount
}

// RecordQuality 记录内容质量数据
func (s *ContentQualityService) RecordQuality(targetID uint, targetType string, result *QualityCheckResult) error {
	quality := models.ContentQuality{
		TargetID:     targetID,
		TargetType:   targetType,
		QualityScore: result.Score,
		IsLowQuality: result.IsLowQuality,
		Reasons:      result.Reason,
		CreatedAt:    time.Now(),
	}

	if err := database.DB.Create(&quality).Error; err != nil {
		log.Printf("[quality] record failed, targetID: %d, targetType: %s, error: %v",
			targetID, targetType, err)
		return err
	}

	log.Printf("[quality] record success, targetID: %d, targetType: %s, score: %.2f, low: %v",
		targetID, targetType, result.Score, result.IsLowQuality)
	return nil
}

// IsLowQuality 检查指定内容是否为低质量
func (s *ContentQualityService) IsLowQuality(targetID uint, targetType string) bool {
	var quality models.ContentQuality
	result := database.DB.Where("target_id = ? AND target_type = ?", targetID, targetType).
		Order("created_at DESC").
		First(&quality)

	if result.Error != nil {
		log.Printf("[quality] query failed, targetID: %d, targetType: %s, error: %v",
			targetID, targetType, result.Error)
		return false
	}

	log.Printf("[quality] query success, targetID: %d, targetType: %s, isLowQuality: %v",
		targetID, targetType, quality.IsLowQuality)
	return quality.IsLowQuality
}

// ValidateContent 验证内容是否合格
func (s *ContentQualityService) ValidateContent(content string) error {
	result := s.CheckContent(content)
	if result.IsLowQuality && result.Score < 30 {
		log.Printf("[quality] validate failed, reason: %s, score: %.2f", result.Reason, result.Score)
		return errors.New(result.Reason)
	}
	log.Printf("[quality] validate passed, score: %.2f", result.Score)
	return nil
}
