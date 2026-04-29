package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get profile: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get profile: user not found, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	errors.Success(w, user)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update profile: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("update profile: failed to decode request body, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	if password, ok := updates["password"].(string); ok && password != "" {
		// 修改密码需要验证旧密码
		oldPassword, ok := updates["old_password"].(string)
		if !ok || oldPassword == "" {
			log.Printf("update profile: old password is required, userID: %d", userID)
			errors.Error(w, errors.CodeInvalidParams, "请输入旧密码")
			return
		}

		// 查询当前用户验证旧密码
		var currentUser models.User
		if err := database.DB.Select("password_hash", "token_version").First(&currentUser, userID).Error; err != nil {
			log.Printf("update profile: user not found, userID: %d", userID)
			errors.Error(w, errors.CodeUserNotFound, "")
			return
		}
		if !utils.CheckPassword(oldPassword, currentUser.PasswordHash) {
			log.Printf("update profile: old password mismatch, userID: %d", userID)
			errors.Error(w, errors.CodeUsernameOrPassword, "旧密码错误")
			return
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Printf("update profile: failed to hash password, userID: %d, error: %v", userID, err)
			errors.Error(w, errors.CodePasswordHashFailed, "")
			return
		}
		updates["password_hash"] = hashedPassword
		// 递增 TokenVersion，使旧token失效
		updates["token_version"] = currentUser.TokenVersion + 1
		log.Printf("update profile: password changed, userID: %d, new token_version: %d", userID, currentUser.TokenVersion+1)
		delete(updates, "password")
		delete(updates, "old_password")
	}

	delete(updates, "id")
	delete(updates, "username")
	delete(updates, "email")
	delete(updates, "role")
	delete(updates, "credits")
	delete(updates, "level")
	delete(updates, "created_at")

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		log.Printf("update profile: failed to update profile, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("update profile: user not found after update, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	log.Printf("update profile: profile updated successfully, userID: %d", userID)
	errors.Success(w, user)
}

func GetCurrentUserTopics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get current user topics: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var topics []models.Topic
	var total int64

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		log.Printf("get current user topics: failed to count topics, userID: %d, error: %v", userID, err)
	}

	if err := database.DB.Where("user_id = ?", userID).Preload("User").Preload("Forum").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&topics).Error; err != nil {
		log.Printf("get current user topics: failed to query topics, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
		"list":      topics,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetCreditUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := database.DB.Order("credits DESC").Limit(10).Find(&users).Error; err != nil {
		log.Printf("get credit users: failed to query users, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, users)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		errors.Success(w, []models.User{})
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	searchPattern := "%" + keyword + "%"

	var users []models.User
	var total int64

	if err := database.DB.Model(&models.User{}).Where("username LIKE ? OR nickname LIKE ?", searchPattern, searchPattern).Count(&total).Error; err != nil {
		log.Printf("search users: failed to count users, keyword: %s, error: %v", keyword, err)
	}

	if err := database.DB.Where("username LIKE ? OR nickname LIKE ?", searchPattern, searchPattern).
		Select("id, username, nickname, avatar, signature, created_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		log.Printf("search users: failed to search users, keyword: %s, error: %v", keyword, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	log.Printf("search users: search completed, keyword: %s, results: %d", keyword, total)
	errors.Success(w, map[string]interface{}{
		"list":  users,
		"total": total,
		"page":  page,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get user: user not found, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeUserNotFound, "")
		return
	}

	errors.Success(w, user)
}

func GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var topicCount, commentCount int64
	database.DB.Model(&models.Topic{}).Where("user_id = ?", userID).Count(&topicCount)
	database.DB.Model(&models.Comment{}).Where("user_id = ?", userID).Count(&commentCount)

	var followCount, followerCount int64
	database.DB.Model(&models.Follow{}).Where("user_id = ?", userID).Count(&followCount)
	database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", userID).Count(&followerCount)

	errors.Success(w, map[string]interface{}{
		"topic_count":    topicCount,
		"comment_count":  commentCount,
		"follow_count":   followCount,
		"follower_count": followerCount,
	})
}

// GetUserFollows 获取指定用户的关注列表
func GetUserFollows(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var follows []models.Follow
	var total int64

	database.DB.Model(&models.Follow{}).Where("user_id = ?", userID).Count(&total)

	if err := database.DB.Where("user_id = ?", userID).
		Preload("FollowUser").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&follows).Error; err != nil {
		log.Printf("get user follows: failed to query follows, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
		"list":  follows,
		"total": total,
		"page":  page,
	})
}

func GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var followers []models.Follow
	var total int64

	database.DB.Model(&models.Follow{}).Where("follow_user_id = ?", userID).Count(&total)

	if err := database.DB.Where("follow_user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&followers).Error; err != nil {
		log.Printf("get user followers: failed to query followers, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	errors.Success(w, map[string]interface{}{
		"list":  followers,
		"total": total,
		"page":  page,
	})
}

func GetUserTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID, _ := strconv.Atoi(vars["id"])

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var topics []models.Topic
	var total int64

	database.DB.Model(&models.Topic{}).Where("user_id = ?", targetUserID).Count(&total)

	if err := database.DB.Where("user_id = ?", targetUserID).
		Preload("User").Preload("Forum").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&topics).Error; err != nil {
		log.Printf("get user topics: failed to query topics, userID: %d, error: %v", targetUserID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 查询用户的勋章
	var userBadges []models.UserBadge
	if err := database.DB.Where("user_id = ? AND is_revoked = ?", targetUserID, false).
		Preload("Badge").
		Find(&userBadges).Error; err != nil {
		log.Printf("get user topics: failed to query user badges, userID: %d, error: %v", targetUserID, err)
	}

	// 获取当前查看者信息（用于匿名判断）
	viewerID, viewerRole, _ := middleware.GetOptionalUserInfo(r)

	// 为每个话题添加 author_badges 并处理匿名
	type ProcessedTopic struct {
		ID               uint                 `json:"id"`
		Title            string               `json:"title"`
		Content          string               `json:"content"`
		UserID           uint                 `json:"user_id"`
		User             interface{}          `json:"user"`
		ForumID          uint                 `json:"forum_id"`
		Forum            models.Forum         `json:"forum"`
		IsPinned         bool                 `json:"is_pinned"`
		IsUserPinned     bool                 `json:"is_user_pinned"`
		IsLocked         bool                 `json:"is_locked"`
		IsEssence        bool                 `json:"is_essence"`
		IsHidden         bool                 `json:"is_hidden"`
		HotScore         float64              `json:"hot_score"`
		LikeCount        int                  `json:"like_count"`
		ViewCount        int                  `json:"view_count"`
		ReplyCount       int                  `json:"reply_count"`
		LastReplyAt      *time.Time           `json:"last_reply_at"`
		AllowComment     bool                 `json:"allow_comment"`
		CreatedAt        time.Time            `json:"created_at"`
		UpdatedAt        time.Time            `json:"updated_at"`
		IsAnonymous      bool                 `json:"is_anonymous"`
		AnonymousType    models.AnonymousType `json:"anonymous_type"`
		AnonymousUntil   *time.Time           `json:"anonymous_until"`
		IsAnonymousEnded bool                 `json:"is_anonymous_ended"`
		AuthorBadges     interface{}          `json:"author_badges"`
		Tags             []models.Tag         `json:"tags"`
	}

	response := make([]ProcessedTopic, len(topics))
	for i, t := range topics {
		shouldShowAnonymous := utils.ShouldShowAnonymous(viewerID, t.UserID, viewerRole, &t)

		item := ProcessedTopic{
			ID:               t.ID,
			Title:            t.Title,
			Content:          t.Content,
			UserID:           t.UserID,
			ForumID:          t.ForumID,
			Forum:            t.Forum,
			IsPinned:         t.IsPinned,
			IsUserPinned:     t.IsUserPinned,
			IsLocked:         t.IsLocked,
			IsEssence:        t.IsEssence,
			IsHidden:         t.IsHidden,
			HotScore:         t.HotScore,
			LikeCount:        t.LikeCount,
			ViewCount:        t.ViewCount,
			ReplyCount:       t.ReplyCount,
			LastReplyAt:      t.LastReplyAt,
			AllowComment:     t.AllowComment,
			CreatedAt:        t.CreatedAt,
			UpdatedAt:        t.UpdatedAt,
			IsAnonymous:      t.IsAnonymous,
			AnonymousType:    t.AnonymousType,
			AnonymousUntil:   t.AnonymousUntil,
			IsAnonymousEnded: t.IsAnonymousEnded,
			Tags:             t.Tags,
		}

		if shouldShowAnonymous {
			item.User = utils.GetAnonymousUserInfo()
			item.AuthorBadges = []models.UserBadge{}
			item.UserID = 0
		} else {
			item.User = t.User
			if t.UserID == uint(targetUserID) {
				item.AuthorBadges = userBadges
			} else {
				item.AuthorBadges = []models.UserBadge{}
			}
		}

		response[i] = item
	}

	errors.Success(w, map[string]interface{}{
		"list":      response,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetFollowTopics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("get follow topics: user not authenticated")
		errors.ErrorWithStatus(w, 401, errors.CodeUnauthorized, "")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var followIDs []uint
	database.DB.Model(&models.Follow{}).Where("user_id = ?", userID).Pluck("follow_user_id", &followIDs)

	if len(followIDs) == 0 {
		errors.Success(w, map[string]interface{}{
			"list":      []models.Topic{},
			"total":     0,
			"page":      page,
			"page_size": pageSize,
		})
		return
	}

	var topics []models.Topic
	var total int64

	database.DB.Model(&models.Topic{}).Where("user_id IN ?", followIDs).Count(&total)

	if err := database.DB.Where("user_id IN ?", followIDs).
		Preload("User").Preload("Forum").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&topics).Error; err != nil {
		log.Printf("get follow topics: failed to query topics, userID: %d, error: %v", userID, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 获取当前查看者信息（用于匿名判断）
	// 获取当前用户角色
	var currentUser models.User
	viewerRole := 0
	if err := database.DB.First(&currentUser, userID).Error; err == nil {
		viewerRole = currentUser.Role
	}

	// 处理匿名数据
	type ProcessedTopic struct {
		ID               uint                 `json:"id"`
		Title            string               `json:"title"`
		Content          string               `json:"content"`
		UserID           uint                 `json:"user_id"`
		User             interface{}          `json:"user"`
		ForumID          uint                 `json:"forum_id"`
		Forum            models.Forum         `json:"forum"`
		IsPinned         bool                 `json:"is_pinned"`
		IsUserPinned     bool                 `json:"is_user_pinned"`
		IsLocked         bool                 `json:"is_locked"`
		IsEssence        bool                 `json:"is_essence"`
		IsHidden         bool                 `json:"is_hidden"`
		HotScore         float64              `json:"hot_score"`
		LikeCount        int                  `json:"like_count"`
		ViewCount        int                  `json:"view_count"`
		ReplyCount       int                  `json:"reply_count"`
		LastReplyAt      *time.Time           `json:"last_reply_at"`
		AllowComment     bool                 `json:"allow_comment"`
		CreatedAt        time.Time            `json:"created_at"`
		UpdatedAt        time.Time            `json:"updated_at"`
		IsAnonymous      bool                 `json:"is_anonymous"`
		AnonymousType    models.AnonymousType `json:"anonymous_type"`
		AnonymousUntil   *time.Time           `json:"anonymous_until"`
		IsAnonymousEnded bool                 `json:"is_anonymous_ended"`
		Tags             []models.Tag         `json:"tags"`
	}

	var processedList []ProcessedTopic
	for _, t := range topics {
		shouldShowAnonymous := utils.ShouldShowAnonymous(userID, t.UserID, viewerRole, &t)

		item := ProcessedTopic{
			ID:               t.ID,
			Title:            t.Title,
			Content:          t.Content,
			UserID:           t.UserID,
			ForumID:          t.ForumID,
			Forum:            t.Forum,
			IsPinned:         t.IsPinned,
			IsUserPinned:     t.IsUserPinned,
			IsLocked:         t.IsLocked,
			IsEssence:        t.IsEssence,
			IsHidden:         t.IsHidden,
			HotScore:         t.HotScore,
			LikeCount:        t.LikeCount,
			ViewCount:        t.ViewCount,
			ReplyCount:       t.ReplyCount,
			LastReplyAt:      t.LastReplyAt,
			AllowComment:     t.AllowComment,
			CreatedAt:        t.CreatedAt,
			UpdatedAt:        t.UpdatedAt,
			IsAnonymous:      t.IsAnonymous,
			AnonymousType:    t.AnonymousType,
			AnonymousUntil:   t.AnonymousUntil,
			IsAnonymousEnded: t.IsAnonymousEnded,
			Tags:             t.Tags,
		}

		if shouldShowAnonymous {
			item.User = utils.GetAnonymousUserInfo()
			item.UserID = 0
		} else {
			item.User = t.User
		}

		processedList = append(processedList, item)
	}

	errors.Success(w, map[string]interface{}{
		"list":      processedList,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
