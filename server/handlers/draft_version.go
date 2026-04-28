package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	DefaultDraftAutoSaveInterval = 60 // 默认自动保存间隔（秒）
	MaxDraftVersions             = 20 // 每个草稿最多保留的版本数
)

// GetDraftVersions 获取草稿版本列表处理器
func GetDraftVersions(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	draftID, _ := strconv.Atoi(vars["id"])

	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", draftID, userID).First(&draft).Error; err != nil {
		log.Printf("get draft versions: draft not found, id: %d, userID: %d, error: %v", draftID, userID, err)
		errors.Error(w, errors.CodeDraftNotFound, "")
		return
	}

	var versions []models.DraftVersion
	if err := database.DB.Where("draft_id = ? AND user_id = ?", draftID, userID).
		Order("version DESC").Find(&versions).Error; err != nil {
		log.Printf("get draft versions: failed to query versions, draftID: %d, userID: %d, error: %v", draftID, userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 转换为带预览的版本列表
	versionsWithPreview := make([]models.DraftVersionWithPreview, len(versions))
	for i, v := range versions {
		versionsWithPreview[i] = models.DraftVersionWithPreview{
			DraftVersion: v,
			Preview:      generatePreview(v.Content),
		}
	}

	errors.Success(w, versionsWithPreview)
}

// GetDraftVersion 获取单个草稿版本详情处理器
func GetDraftVersion(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	draftID, _ := strconv.Atoi(vars["id"])
	versionID, _ := strconv.Atoi(vars["version_id"])

	var version models.DraftVersion
	if err := database.DB.Where("id = ? AND draft_id = ? AND user_id = ?", versionID, draftID, userID).
		First(&version).Error; err != nil {
		log.Printf("get draft version: version not found, versionID: %d, draftID: %d, userID: %d, error: %v", versionID, draftID, userID, err)
		errors.Error(w, errors.CodeDraftNotFound, "")
		return
	}

	errors.Success(w, version)
}

// CreateDraftVersion 创建草稿版本处理器
func CreateDraftVersion(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	draftID, _ := strconv.Atoi(vars["id"])

	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", draftID, userID).First(&draft).Error; err != nil {
		log.Printf("create draft version: draft not found, id: %d, userID: %d, error: %v", draftID, userID, err)
		errors.Error(w, errors.CodeDraftNotFound, "")
		return
	}

	var req struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		ForumID uint     `json:"forum_id"`
		Tags    []int    `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create draft version: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 获取当前最大版本号
	var maxVersion int
	database.DB.Model(&models.DraftVersion{}).
		Where("draft_id = ?", draftID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion)

	// 创建新版本
	version := models.DraftVersion{
		DraftID: draft.ID,
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		ForumID: req.ForumID,
		Tags:    req.Tags,
		Version: maxVersion + 1,
	}

	if err := database.DB.Create(&version).Error; err != nil {
		log.Printf("create draft version: failed to create version, draftID: %d, userID: %d, error: %v", draftID, userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 检查是否超过最大版本数，超过则删除最旧的版本
	var versionCount int64
	database.DB.Model(&models.DraftVersion{}).Where("draft_id = ?", draftID).Count(&versionCount)

	if versionCount > MaxDraftVersions {
		var oldestVersions []models.DraftVersion
		database.DB.Where("draft_id = ?", draftID).
			Order("version ASC").
			Limit(int(versionCount - MaxDraftVersions)).
			Find(&oldestVersions)

		for _, v := range oldestVersions {
			database.DB.Unscoped().Delete(&v)
		}
	}

	log.Printf("create draft version: version created successfully, draftID: %d, version: %d", draftID, version.Version)
	errors.Success(w, version)
}

// RestoreDraftVersion 回退到指定版本处理器
func RestoreDraftVersion(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	draftID, _ := strconv.Atoi(vars["id"])
	versionID, _ := strconv.Atoi(vars["version_id"])

	var version models.DraftVersion
	if err := database.DB.Where("id = ? AND draft_id = ? AND user_id = ?", versionID, draftID, userID).
		First(&version).Error; err != nil {
		log.Printf("restore draft version: version not found, versionID: %d, draftID: %d, userID: %d, error: %v", versionID, draftID, userID, err)
		errors.Error(w, errors.CodeDraftNotFound, "")
		return
	}

	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", draftID, userID).First(&draft).Error; err != nil {
		log.Printf("restore draft version: draft not found, id: %d, userID: %d, error: %v", draftID, userID, err)
		errors.Error(w, errors.CodeDraftNotFound, "")
		return
	}

	// 在回退前，先保存当前状态为一个新版本
	var maxVersion int
	database.DB.Model(&models.DraftVersion{}).
		Where("draft_id = ?", draftID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion)

	currentVersion := models.DraftVersion{
		DraftID: draft.ID,
		UserID:  userID,
		Title:   draft.Title,
		Content: draft.Content,
		ForumID: draft.ForumID,
		Tags:    draft.Tags,
		Version: maxVersion + 1,
	}

	if err := database.DB.Create(&currentVersion).Error; err != nil {
		log.Printf("restore draft version: failed to save current version, draftID: %d, error: %v", draftID, err)
	}

	// 更新草稿为指定版本的内容
	updates := map[string]interface{}{
		"title":   version.Title,
		"content": version.Content,
		"forum_id": version.ForumID,
		"tags":    version.Tags,
	}

	if err := database.DB.Model(&draft).Updates(updates).Error; err != nil {
		log.Printf("restore draft version: failed to update draft, id: %d, error: %v", draftID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("restore draft version: draft restored to version %d, draftID: %d", version.Version, draftID)
	errors.Success(w, draft)
}

// DeleteDraftVersion 删除草稿版本处理器
func DeleteDraftVersion(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	draftID, _ := strconv.Atoi(vars["id"])
	versionID, _ := strconv.Atoi(vars["version_id"])

	var version models.DraftVersion
	if err := database.DB.Where("id = ? AND draft_id = ? AND user_id = ?", versionID, draftID, userID).
		First(&version).Error; err != nil {
		log.Printf("delete draft version: version not found, versionID: %d, draftID: %d, userID: %d, error: %v", versionID, draftID, userID, err)
		errors.Error(w, errors.CodeDraftNotFound, "")
		return
	}

	if err := database.DB.Unscoped().Delete(&version).Error; err != nil {
		log.Printf("delete draft version: failed to delete version, id: %d, error: %v", versionID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("delete draft version: version deleted successfully, id: %d, draftID: %d", versionID, draftID)
	errors.Success(w, nil)
}

// CleanupDraftVersions 清理草稿版本（发布后调用）
func CleanupDraftVersions(draftID uint, userID uint) error {
	// 保留最近3个版本，删除其余的
	var versions []models.DraftVersion
	if err := database.DB.Where("draft_id = ? AND user_id = ?", draftID, userID).
		Order("version DESC").Find(&versions).Error; err != nil {
		return err
	}

	if len(versions) <= 3 {
		return nil
	}

	// 删除除了最近3个版本之外的所有版本
	for i := 3; i < len(versions); i++ {
		if err := database.DB.Unscoped().Delete(&versions[i]).Error; err != nil {
			log.Printf("cleanup draft versions: failed to delete version %d, error: %v", versions[i].ID, err)
		}
	}

	log.Printf("cleanup draft versions: cleaned up %d versions for draft %d", len(versions)-3, draftID)
	return nil
}

// GetDraftAutoSaveInterval 获取自动保存间隔配置
func GetDraftAutoSaveInterval() int {
	var config models.SiteConfig
	if err := database.DB.Where("key = ?", "draft_auto_save_interval").First(&config).Error; err != nil {
		return DefaultDraftAutoSaveInterval
	}

	interval, err := strconv.Atoi(config.Value)
	if err != nil || interval <= 0 {
		return DefaultDraftAutoSaveInterval
	}

	return interval
}

// generatePreview 生成内容预览
func generatePreview(content string) string {
	// 移除 Markdown 格式标记
	preview := content
	preview = strings.ReplaceAll(preview, "#", "")
	preview = strings.ReplaceAll(preview, "*", "")
	preview = strings.ReplaceAll(preview, "`", "")
	preview = strings.ReplaceAll(preview, "~", "")
	preview = strings.ReplaceAll(preview, "\n", " ")
	preview = strings.TrimSpace(preview)

	// 截取前100个字符
	runes := []rune(preview)
	if len(runes) > 100 {
		preview = string(runes[:100]) + "..."
	}

	return preview
}

// GetDraftConfig 获取草稿相关配置处理器
func GetDraftConfig(w http.ResponseWriter, r *http.Request) {
	interval := GetDraftAutoSaveInterval()
	config := map[string]interface{}{
		"auto_save_interval": interval,
		"max_versions":       MaxDraftVersions,
	}
	errors.Success(w, config)
}

// PreviewDraftVersion 预览草稿版本处理器
func PreviewDraftVersion(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	draftID, _ := strconv.Atoi(vars["id"])
	versionID, _ := strconv.Atoi(vars["version_id"])

	var version models.DraftVersion
	if err := database.DB.Where("id = ? AND draft_id = ? AND user_id = ?", versionID, draftID, userID).
		First(&version).Error; err != nil {
		log.Printf("preview draft version: version not found, versionID: %d, draftID: %d, userID: %d, error: %v", versionID, draftID, userID, err)
		errors.Error(w, errors.CodeDraftNotFound, "")
		return
	}

	// 返回版本内容供预览
	preview := map[string]interface{}{
		"id":         version.ID,
		"draft_id":   version.DraftID,
		"title":      version.Title,
		"content":    version.Content,
		"forum_id":   version.ForumID,
		"tags":       version.Tags,
		"version":    version.Version,
		"created_at": version.CreatedAt.Format(time.RFC3339),
		"preview":    generatePreview(version.Content),
	}

	errors.Success(w, preview)
}
