package handlers

import (
	"bbsgo/database"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetDrafts 获取当前用户的草稿列表处理器
func GetDrafts(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var drafts []models.Draft
	if err := database.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&drafts).Error; err != nil {
		log.Printf("get drafts: failed to query drafts, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "获取草稿列表失败")
		return
	}

	utils.Success(w, drafts)
}

// GetDraft 获取单个草稿详情处理器
func GetDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&draft).Error; err != nil {
		log.Printf("get draft: draft not found, id: %d, userID: %d, error: %v", id, userID, err)
		utils.Error(w, 404, "草稿不存在")
		return
	}

	utils.Success(w, draft)
}

// CreateDraft 创建草稿处理器
func CreateDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 解析请求体
	var draft models.Draft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		log.Printf("create draft: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	draft.UserID = userID

	// 创建草稿
	if err := database.DB.Create(&draft).Error; err != nil {
		log.Printf("create draft: failed to create draft, userID: %d, error: %v", userID, err)
		utils.Error(w, 500, "保存失败")
		return
	}

	log.Printf("create draft: draft created successfully, id: %d, userID: %d", draft.ID, userID)
	utils.Success(w, draft)
}

// UpdateDraft 更新草稿处理器
func UpdateDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 检查草稿是否存在且属于当前用户
	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&draft).Error; err != nil {
		log.Printf("update draft: draft not found, id: %d, userID: %d, error: %v", id, userID, err)
		utils.Error(w, 404, "草稿不存在")
		return
	}

	// 解析请求体
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update draft: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 过滤不允许更新的字段
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")

	// 执行更新
	if err := database.DB.Model(&draft).Updates(updates).Error; err != nil {
		log.Printf("update draft: failed to update draft, id: %d, userID: %d, error: %v", id, userID, err)
		utils.Error(w, 500, "更新失败")
		return
	}

	utils.Success(w, draft)
}

// DeleteDraft 删除草稿处理器
func DeleteDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 检查草稿是否存在且属于当前用户
	var draft models.Draft
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&draft).Error; err != nil {
		log.Printf("delete draft: draft not found, id: %d, userID: %d, error: %v", id, userID, err)
		utils.Error(w, 404, "草稿不存在")
		return
	}

	// 物理删除草稿
	if err := database.DB.Unscoped().Delete(&draft).Error; err != nil {
		log.Printf("delete draft: failed to delete draft, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除失败")
		return
	}

	log.Printf("delete draft: draft deleted successfully, id: %d, userID: %d", id, userID)
	utils.Success(w, nil)
}
