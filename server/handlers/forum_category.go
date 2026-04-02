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

// GetForumCategories 获取启用的版块分类列表处理器
func GetForumCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.ForumCategory
	if err := database.DB.Where("is_active = ?", true).Order("sort_order").Find(&categories).Error; err != nil {
		log.Printf("get forum categories: failed to query categories, error: %v", err)
		utils.Error(w, 500, "获取分类列表失败")
		return
	}

	utils.Success(w, categories)
}

// GetAllForumCategories 获取所有版块分类列表处理器（管理员）
func GetAllForumCategories(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("get all forum categories: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	var categories []models.ForumCategory
	if err := database.DB.Order("sort_order").Find(&categories).Error; err != nil {
		log.Printf("get all forum categories: failed to query categories, error: %v", err)
		utils.Error(w, 500, "获取分类列表失败")
		return
	}

	utils.Success(w, categories)
}

// CreateForumCategory 创建版块分类处理器（管理员）
func CreateForumCategory(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("create forum category: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var req models.ForumCategory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create forum category: failed to decode request body, error: %v", err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 创建分类
	if err := database.DB.Create(&req).Error; err != nil {
		log.Printf("create forum category: failed to create category, name: %s, error: %v", req.Name, err)
		utils.Error(w, 500, "创建分类失败")
		return
	}

	log.Printf("create forum category: category created successfully, id: %d, name: %s", req.ID, req.Name)
	utils.Success(w, req)
}

// UpdateForumCategory 更新版块分类处理器（管理员）
func UpdateForumCategory(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("update forum category: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 解析请求体
	var req models.ForumCategory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update forum category: failed to decode request body, id: %d, error: %v", id, err)
		utils.Error(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 查询要更新的分类
	var category models.ForumCategory
	if result := database.DB.First(&category, id); result.Error != nil {
		log.Printf("update forum category: category not found, id: %d, error: %v", id, result.Error)
		utils.Error(w, http.StatusNotFound, "分类不存在")
		return
	}

	// 更新字段
	category.Name = req.Name
	category.Icon = req.Icon
	category.Description = req.Description
	category.SortOrder = req.SortOrder
	category.IsActive = req.IsActive

	// 保存更新
	if err := database.DB.Save(&category).Error; err != nil {
		log.Printf("update forum category: failed to update category, id: %d, error: %v", id, err)
		utils.Error(w, 500, "更新分类失败")
		return
	}

	log.Printf("update forum category: category updated successfully, id: %d, name: %s", id, req.Name)
	utils.Success(w, category)
}

// DeleteForumCategory 删除版块分类处理器（管理员）
func DeleteForumCategory(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	_, ok := middleware.GetAdminIDFromContext(r.Context())
	if !ok {
		log.Printf("delete forum category: unauthorized access")
		utils.Error(w, http.StatusUnauthorized, "未授权")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// 物理删除分类
	if err := database.DB.Unscoped().Delete(&models.ForumCategory{}, id).Error; err != nil {
		log.Printf("delete forum category: failed to delete category, id: %d, error: %v", id, err)
		utils.Error(w, 500, "删除分类失败")
		return
	}

	log.Printf("delete forum category: category deleted successfully, id: %d", id)
	utils.Success(w, nil)
}
