package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/models"
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
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
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
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
		"list":  follows,
		"total": total,
		"page":  page,
	})
}

func DeleteAdminFollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	followID, err := strconv.Atoi(vars["id"])
	if err != nil {
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	if err := database.DB.Unscoped().Delete(&models.Follow{}, followID).Error; err != nil {
		log.Printf("delete admin follow: failed to delete follow, id: %d, error: %v", followID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, nil)
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
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
		"list":  comments,
		"total": total,
		"page":  page,
	})
}

func UpdateCommentBest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	var req struct {
		IsBest bool `json:"is_best"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update comment best: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		log.Printf("update comment best: comment not found, id: %d", commentID)
		errors.Error(w, errors.CodeCommentNotFound, "")
		return
	}

	if err := database.DB.Model(&comment).UpdateColumn("is_best", req.IsBest).Error; err != nil {
		log.Printf("update comment best: failed to update, id: %d, error: %v", commentID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, nil)
}
