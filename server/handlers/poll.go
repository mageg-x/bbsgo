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
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreatePoll 创建投票
func CreatePoll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("create poll: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	var req struct {
		TopicID     uint   `json:"topic_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PollType    string `json:"poll_type"`
		MaxChoices  int    `json:"max_choices"`
		EndTime     string `json:"end_time"`
		Options     []struct {
			Text string `json:"text"`
		} `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create poll: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.TopicID == 0 {
		log.Printf("create poll: topic ID is empty, userID: %d", userID)
		utils.Error(w, 400, "话题ID不能为空")
		return
	}

	if len(req.Options) < 2 || len(req.Options) > 10 {
		log.Printf("create poll: invalid options count: %d, userID: %d", len(req.Options), userID)
		utils.Error(w, 400, "选项数量必须在2-10个之间")
		return
	}

	var topic models.Topic
	if err := database.DB.First(&topic, req.TopicID).Error; err != nil {
		log.Printf("create poll: topic not found, topicID: %d, error: %v", req.TopicID, err)
		utils.Error(w, 404, "话题不存在")
		return
	}

	if topic.UserID != userID {
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil || user.Role < 1 {
			log.Printf("create poll: permission denied, userID: %d, topicUserID: %d", userID, topic.UserID)
			utils.Error(w, 403, "无权限创建投票")
			return
		}
	}

	var existingPoll models.Poll
	if err := database.DB.Where("topic_id = ?", req.TopicID).First(&existingPoll).Error; err == nil {
		log.Printf("create poll: poll already exists for topic, topicID: %d", req.TopicID)
		utils.Error(w, 400, "该话题已存在投票")
		return
	}

	pollType := req.PollType
	if pollType == "" {
		pollType = "single"
	}
	if pollType != "single" && pollType != "multiple" {
		log.Printf("create poll: invalid poll type: %s, userID: %d", pollType, userID)
		utils.Error(w, 400, "投票类型无效")
		return
	}

	maxChoices := req.MaxChoices
	if maxChoices < 1 {
		maxChoices = 1
	}
	if pollType == "single" {
		maxChoices = 1
	}

	var endTime *time.Time
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			log.Printf("create poll: invalid end time format: %s, error: %v", req.EndTime, err)
			utils.Error(w, 400, "截止时间格式无效")
			return
		}
		if t.Before(time.Now()) {
			log.Printf("create poll: end time is in the past: %s", req.EndTime)
			utils.Error(w, 400, "截止时间不能早于当前时间")
			return
		}
		endTime = &t
	}

	poll := models.Poll{
		TopicID:     req.TopicID,
		Title:       req.Title,
		Description: req.Description,
		PollType:    pollType,
		MaxChoices:  maxChoices,
		EndTime:     endTime,
		Options:     []models.PollOption{},
	}

	if err := database.DB.Create(&poll).Error; err != nil {
		log.Printf("create poll: failed to create poll, error: %v", err)
		utils.Error(w, 500, "创建投票失败")
		return
	}

	for i, opt := range req.Options {
		if opt.Text == "" {
			continue
		}
		option := models.PollOption{
			PollID:    poll.ID,
			Text:      opt.Text,
			SortOrder: i,
		}
		if err := database.DB.Create(&option).Error; err != nil {
			log.Printf("create poll: failed to create option, error: %v", err)
		} else {
			poll.Options = append(poll.Options, option)
		}
	}

	utils.Success(w, poll)
}

// GetPoll 获取投票详情
func GetPoll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var poll models.Poll
	if err := database.DB.Preload("Options").First(&poll, id).Error; err != nil {
		log.Printf("get poll: poll not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "投票不存在")
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())
	if userID > 0 {
		var votes []models.PollVote
		database.DB.Where("poll_id = ? AND user_id = ?", poll.ID, userID).Find(&votes)
		if len(votes) > 0 {
			votedOptionIDs := make([]uint, len(votes))
			for i, v := range votes {
				votedOptionIDs[i] = v.OptionID
			}
			response := map[string]interface{}{
				"poll":             poll,
				"has_voted":        true,
				"voted_option_ids": votedOptionIDs,
			}
			utils.Success(w, response)
			return
		}
	}

	response := map[string]interface{}{
		"poll":      poll,
		"has_voted": false,
	}
	utils.Success(w, response)
}

// GetPollByTopic 根据话题ID获取投票
func GetPollByTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["topic_id"])

	var poll models.Poll
	if err := database.DB.Where("topic_id = ?", topicID).Preload("Options").First(&poll).Error; err != nil {
		log.Printf("get poll by topic: poll not found, topicID: %d, error: %v", topicID, err)
		utils.Error(w, 404, "该话题没有投票")
		return
	}

	// 尝试从 token 获取用户ID（可选认证）
	userID := tryGetUserIDFromRequest(r)

	if userID > 0 {
		var votes []models.PollVote
		database.DB.Where("poll_id = ? AND user_id = ?", poll.ID, userID).Find(&votes)
		if len(votes) > 0 {
			votedOptionIDs := make([]uint, len(votes))
			for i, v := range votes {
				votedOptionIDs[i] = v.OptionID
			}
			response := map[string]interface{}{
				"poll":             poll,
				"has_voted":        true,
				"voted_option_ids": votedOptionIDs,
			}
			utils.Success(w, response)
			return
		}
	}

	response := map[string]interface{}{
		"poll":      poll,
		"has_voted": false,
	}
	utils.Success(w, response)
}

// tryGetUserIDFromRequest 尝试从请求中获取用户ID（不强制认证）
func tryGetUserIDFromRequest(r *http.Request) uint {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return 0
	}

	return claims.UserID
}

// SubmitVote 提交投票
func SubmitVote(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("submit vote: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	var req struct {
		PollID    uint   `json:"poll_id"`
		OptionIDs []uint `json:"option_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("submit vote: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	if req.PollID == 0 || len(req.OptionIDs) == 0 {
		log.Printf("submit vote: incomplete params, pollID: %d, optionIDs count: %d", req.PollID, len(req.OptionIDs))
		utils.Error(w, 400, "参数不完整")
		return
	}

	var poll models.Poll
	if err := database.DB.Preload("Options").First(&poll, req.PollID).Error; err != nil {
		log.Printf("submit vote: poll not found, pollID: %d, error: %v", req.PollID, err)
		utils.Error(w, 404, "投票不存在")
		return
	}

	if poll.IsEnded {
		log.Printf("submit vote: poll already ended, pollID: %d", poll.ID)
		utils.Error(w, 400, "投票已结束")
		return
	}

	if poll.EndTime != nil && poll.EndTime.Before(time.Now()) {
		log.Printf("submit vote: poll deadline passed, pollID: %d, endTime: %v", poll.ID, poll.EndTime)
		utils.Error(w, 400, "投票已截止")
		return
	}

	var existingVote models.PollVote
	if err := database.DB.Where("poll_id = ? AND user_id = ?", poll.ID, userID).First(&existingVote).Error; err == nil {
		log.Printf("submit vote: user already voted, pollID: %d, userID: %d", poll.ID, userID)
		utils.Error(w, 400, "您已经投过票了")
		return
	}

	if poll.PollType == "single" && len(req.OptionIDs) > 1 {
		log.Printf("submit vote: single poll with multiple options, pollID: %d, userID: %d", poll.ID, userID)
		utils.Error(w, 400, "单选投票只能选择一个选项")
		return
	}

	if poll.PollType == "multiple" && len(req.OptionIDs) > poll.MaxChoices {
		log.Printf("submit vote: too many options selected, pollID: %d, userID: %d, selected: %d, max: %d", poll.ID, userID, len(req.OptionIDs), poll.MaxChoices)
		utils.Error(w, 400, "选择的选项数量超过限制")
		return
	}

	optionMap := make(map[uint]bool)
	for _, opt := range poll.Options {
		optionMap[opt.ID] = true
	}

	for _, optID := range req.OptionIDs {
		if !optionMap[optID] {
			log.Printf("submit vote: invalid option ID, pollID: %d, optionID: %d", poll.ID, optID)
			utils.Error(w, 400, "无效的选项ID")
			return
		}
	}

	tx := database.DB.Begin()

	for _, optID := range req.OptionIDs {
		vote := models.PollVote{
			PollID:   poll.ID,
			OptionID: optID,
			UserID:   userID,
		}
		if err := tx.Create(&vote).Error; err != nil {
			tx.Rollback()
			log.Printf("submit vote: failed to create vote, error: %v", err)
			utils.Error(w, 500, "投票失败")
			return
		}

		if err := tx.Model(&models.PollOption{}).Where("id = ?", optID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + 1")).Error; err != nil {
			tx.Rollback()
			log.Printf("submit vote: failed to update option count, error: %v", err)
			utils.Error(w, 500, "投票失败")
			return
		}
	}

	if err := tx.Model(&poll).UpdateColumn("total_votes", poll.TotalVotes+len(req.OptionIDs)).Error; err != nil {
		tx.Rollback()
		log.Printf("submit vote: failed to update poll count, error: %v", err)
		utils.Error(w, 500, "投票失败")
		return
	}

	tx.Commit()

	database.DB.Preload("Options").First(&poll, poll.ID)
	utils.Success(w, poll)
}

// UpdatePoll 更新投票（管理员）
func UpdatePoll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("update poll: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || user.Role < 1 {
		log.Printf("update poll: permission denied, userID: %d", userID)
		utils.Error(w, 403, "无权限")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var poll models.Poll
	if err := database.DB.Preload("Options").First(&poll, id).Error; err != nil {
		log.Printf("update poll: poll not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "投票不存在")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		EndTime     string `json:"end_time"`
		MaxChoices  int    `json:"max_choices"`
		Options     []struct {
			ID   uint   `json:"id"`
			Text string `json:"text"`
		} `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("update poll: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			log.Printf("update poll: invalid end time format: %s, error: %v", req.EndTime, err)
			utils.Error(w, 400, "截止时间格式无效")
			return
		}
		if poll.EndTime != nil && t.Before(*poll.EndTime) {
			log.Printf("update poll: cannot shorten end time, pollID: %d, currentEndTime: %v, newEndTime: %s", poll.ID, poll.EndTime, req.EndTime)
			utils.Error(w, 400, "只能延长截止时间，不能缩短")
			return
		}
		updates["end_time"] = &t
	}

	if poll.PollType == "multiple" && req.MaxChoices > 0 {
		updates["max_choices"] = req.MaxChoices
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&poll).Updates(updates).Error; err != nil {
			log.Printf("update poll: failed to update poll, error: %v", err)
			utils.Error(w, 500, "更新失败")
			return
		}
	}

	if len(req.Options) > 0 {
		if poll.TotalVotes > 0 {
			for _, optReq := range req.Options {
				if optReq.ID > 0 && optReq.Text != "" {
					database.DB.Model(&models.PollOption{}).Where("id = ?", optReq.ID).Update("text", optReq.Text)
				}
			}
		} else {
			database.DB.Unscoped().Where("poll_id = ?", poll.ID).Delete(&models.PollOption{})
			for i, optReq := range req.Options {
				if optReq.Text != "" {
					option := models.PollOption{
						PollID:    poll.ID,
						Text:      optReq.Text,
						SortOrder: i,
					}
					database.DB.Create(&option)
				}
			}
		}
	}

	database.DB.Preload("Options").First(&poll, poll.ID)
	utils.Success(w, poll)
}

// EndPoll 结束投票（管理员）
func EndPoll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("end poll: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || user.Role < 1 {
		log.Printf("end poll: permission denied, userID: %d", userID)
		utils.Error(w, 403, "无权限")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var poll models.Poll
	if err := database.DB.First(&poll, id).Error; err != nil {
		log.Printf("end poll: poll not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "投票不存在")
		return
	}

	if err := database.DB.Model(&poll).Update("is_ended", true).Error; err != nil {
		log.Printf("end poll: failed to end poll, error: %v", err)
		utils.Error(w, 500, "结束投票失败")
		return
	}

	utils.Success(w, map[string]string{"message": "投票已结束"})
}

// DeletePoll 删除投票（管理员）
func DeletePoll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Printf("delete poll: user not authenticated")
		utils.Error(w, 401, "未认证")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || user.Role < 1 {
		log.Printf("delete poll: permission denied, userID: %d", userID)
		utils.Error(w, 403, "无权限")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var poll models.Poll
	if err := database.DB.First(&poll, id).Error; err != nil {
		log.Printf("delete poll: poll not found, id: %d, error: %v", id, err)
		utils.Error(w, 404, "投票不存在")
		return
	}

	database.DB.Unscoped().Where("poll_id = ?", poll.ID).Delete(&models.PollVote{})
	database.DB.Unscoped().Where("poll_id = ?", poll.ID).Delete(&models.PollOption{})
	database.DB.Unscoped().Delete(&poll)

	utils.Success(w, map[string]string{"message": "投票已删除"})
}

// GetAdminPolls 获取投票列表（管理员）
func GetAdminPolls(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var polls []models.Poll
	var total int64

	database.DB.Model(&models.Poll{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Preload("Options").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&polls).Error; err != nil {
		log.Printf("get admin polls: failed to query polls, error: %v", err)
		utils.Error(w, 500, "获取投票列表失败")
		return
	}

	var pollsWithTopic []map[string]interface{}
	for _, poll := range polls {
		pollMap := map[string]interface{}{
			"poll": poll,
		}
		var topic models.Topic
		if err := database.DB.First(&topic, poll.TopicID).Error; err == nil {
			pollMap["topic"] = topic
		}
		pollsWithTopic = append(pollsWithTopic, pollMap)
	}

	utils.Success(w, map[string]interface{}{
		"list":      pollsWithTopic,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
