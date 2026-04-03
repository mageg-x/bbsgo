package handlers

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAdminFollows(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	keyword := r.URL.Query().Get("keyword")

	var follows []models.Follow
	var total int64

	query := database.DB.Model(&models.Follow{})

	if keyword != "" {
		subQuery := database.DB.Model(&models.User{}).Where("nickname LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Select("id")
		query = query.Where("user_id IN (?) OR follow_user_id IN (?)", subQuery, subQuery)
	}

	query.Count(&total)

	if err := query.Preload("User").Preload("FollowUser").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&follows).Error; err != nil {
		log.Printf("get admin follows: failed to query follows, error: %v", err)
		utils.Error(w, 500, "获取关注列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":  follows,
		"total": total,
		"page":  page,
	})
}

func GetAdminFollowers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	keyword := r.URL.Query().Get("keyword")

	var follows []models.Follow
	var total int64

	query := database.DB.Model(&models.Follow{})

	if keyword != "" {
		subQuery := database.DB.Model(&models.User{}).Where("nickname LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Select("id")
		query = query.Where("user_id IN (?) OR follow_user_id IN (?)", subQuery, subQuery)
	}

	query.Count(&total)

	if err := query.Preload("User").Preload("FollowUser").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&follows).Error; err != nil {
		log.Printf("get admin followers: failed to query followers, error: %v", err)
		utils.Error(w, 500, "获取粉丝列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":  follows,
		"total": total,
		"page":  page,
	})
}

func DeleteAdminFollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, 400, "无效的关注ID")
		return
	}

	if err := database.DB.Unscoped().Delete(&models.Follow{}, followID).Error; err != nil {
		log.Printf("delete admin follow: failed to delete follow, id: %d, error: %v", followID, err)
		utils.Error(w, 500, "删除关注失败")
		return
	}

	utils.Success(w, nil)
}

func GetAdminBestComments(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	keyword := r.URL.Query().Get("keyword")

	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{}).Where("is_best = ?", true)

	if keyword != "" {
		subQuery := database.DB.Model(&models.Topic{}).Where("title LIKE ?", "%"+keyword+"%").Select("id")
		query = query.Where("topic_id IN (?) OR content LIKE ?", subQuery, "%"+keyword+"%")
	}

	query.Count(&total)

	if err := query.Preload("User").Preload("Topic").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error; err != nil {
		log.Printf("get admin best comments: failed to query best comments, error: %v", err)
		utils.Error(w, 500, "获取最佳评论列表失败")
		return
	}

	utils.Success(w, map[string]interface{}{
		"list":  comments,
		"total": total,
		"page":  page,
	})
}

func UpdateCommentBest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.Error(w, 400, "无效的评论ID")
		return
	}

	var req struct {
		IsBest bool `json:"is_best"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update comment best: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		log.Printf("update comment best: comment not found, id: %d", commentID)
		utils.Error(w, 404, "评论不存在")
		return
	}

	if err := database.DB.Model(&comment).UpdateColumn("is_best", req.IsBest).Error; err != nil {
		log.Printf("update comment best: failed to update, id: %d, error: %v", commentID, err)
		utils.Error(w, 500, "操作失败")
		return
	}

	utils.Success(w, nil)
}
