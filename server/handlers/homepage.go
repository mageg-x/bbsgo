package handlers

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/middleware"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TopicWithPoll 首页话题结构，与 GetTopics 返回格式一致
type TopicWithPoll struct {
	models.Topic
	HasPoll      bool               `json:"has_poll"`
	AuthorBadges []models.UserBadge `json:"author_badges"`
}

// GetHomePage 获取首页聚合数据
// 一次性返回首页所需的所有数据，减少前端多次请求
func GetHomePage(w http.ResponseWriter, r *http.Request) {
	// 获取分区缓存
	forums, _ := homePageForums()
	tags, _ := homePageTags()
	announcements, _ := homePageAnnouncements()
	topics, total, _ := fetchTopicsForHome(1, 20)

	// 获取当前查看者信息（用于匿名判断）
	viewerID, viewerRole, _ := middleware.GetOptionalUserInfo(r)

	// 处理匿名数据
	processedTopics := processTopicsForAnonymity(topics, viewerID, viewerRole)

	homeData := map[string]interface{}{
		"forums":        forums,
		"tags":          tags,
		"announcements": announcements,
		"topics": map[string]interface{}{
			"list":      processedTopics,
			"total":     total,
			"page":      1,
			"page_size": 20,
		},
		"updated_at": time.Now(),
	}

	errors.Success(w, homeData)
}

// homePageForums 获取首页板块列表（带缓存）
func homePageForums() ([]models.Forum, error) {
	data, ok := cache.HomePageCache.Forums().Get()
	if ok && data != nil {
		return data.([]models.Forum), nil
	}

	var forums []models.Forum
	if err := database.DB.Where("name != ?", "全部").
		Order("sort_order ASC, id ASC").
		Find(&forums).Error; err != nil {
		log.Printf("home page forums: query error, %v", err)
		return forums, err
	}

	cache.HomePageCache.Forums().Set(forums)
	return forums, nil
}

// homePageTags 获取首页标签列表（带缓存）
func homePageTags() ([]models.Tag, error) {
	data, ok := cache.HomePageCache.Tags().Get()
	if ok && data != nil {
		return data.([]models.Tag), nil
	}

	var tags []models.Tag
	if err := database.DB.Where("is_banned = ?", false).
		Order("is_official DESC, usage_count DESC, sort_order ASC").
		Limit(20).
		Find(&tags).Error; err != nil {
		log.Printf("home page tags: query error, %v", err)
		return tags, err
	}

	cache.HomePageCache.Tags().Set(tags)
	return tags, nil
}

// homePageAnnouncements 获取首页公告列表（带缓存）
func homePageAnnouncements() ([]models.Announcement, error) {
	data, ok := cache.HomePageCache.Announcements().Get()
	if ok && data != nil {
		return data.([]models.Announcement), nil
	}

	var announcements []models.Announcement
	now := time.Now()
	if err := database.DB.Where("expires_at IS NULL OR expires_at > ?", now).
		Order("is_pinned DESC, created_at DESC").
		Find(&announcements).Error; err != nil {
		log.Printf("home page announcements: query error, %v", err)
		return announcements, err
	}

	cache.HomePageCache.Announcements().Set(announcements)
	return announcements, nil
}

// fetchTopicsForHome 获取首页话题列表（带缓存）
func fetchTopicsForHome(page, pageSize int) ([]TopicWithPoll, int64, error) {
	// 尝试从缓存获取（只缓存第1页）
	if page == 1 && pageSize == 20 {
		data, ok := cache.HomePageCache.Topics().TopicsPageGet(page, pageSize)
		if ok && data != nil {
			result := data.(map[string]interface{})
			topics := result["list"].([]TopicWithPoll)
			total := result["total"].(int64)
			return topics, total, nil
		}
	}

	topics, total, err := doFetchTopicsForHome(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 只缓存第1页
	if page == 1 && pageSize == 20 {
		cache.HomePageCache.Topics().TopicsPageSet(page, pageSize, map[string]interface{}{
			"list":  topics,
			"total": total,
		})
	}

	return topics, total, nil
}

// doFetchTopicsForHome 获取首页话题列表核心逻辑
func doFetchTopicsForHome(page, pageSize int) ([]TopicWithPoll, int64, error) {
	var topics []models.Topic
	var total int64

	// 统计总数
	if err := database.DB.Model(&models.Topic{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询列表
	if err := database.DB.Preload("User").Preload("Forum").Preload("Tags").
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&topics).Error; err != nil {
		return nil, 0, err
	}

	// 收集 topic IDs 和 user IDs
	topicIDs := make([]uint, len(topics))
	userIDs := make(map[uint]bool)
	for i, t := range topics {
		topicIDs[i] = t.ID
		userIDs[t.UserID] = true
	}

	// 批量查询投票
	hasPollMap := make(map[uint]bool)
	if len(topicIDs) > 0 {
		var polls []models.Poll
		database.DB.Model(&models.Poll{}).Where("topic_id IN ?", topicIDs).Select("topic_id").Find(&polls)
		for _, p := range polls {
			hasPollMap[p.TopicID] = true
		}
	}

	// 批量查询用户勋章
	userBadgesMap := make(map[uint][]models.UserBadge)
	if len(userIDs) > 0 {
		ids := make([]uint, 0, len(userIDs))
		for id := range userIDs {
			ids = append(ids, id)
		}
		var userBadges []models.UserBadge
		if err := database.DB.Where("user_id IN ? AND is_revoked = ?", ids, false).
			Preload("Badge").
			Find(&userBadges).Error; err != nil {
			log.Printf("fetch home page topics: failed to query user badges, error: %v", err)
		}
		for _, ub := range userBadges {
			userBadgesMap[ub.UserID] = append(userBadgesMap[ub.UserID], ub)
		}
	}

	// 构建返回数据，与 GetTopics 格式一致
	var response []TopicWithPoll
	for _, t := range topics {
		response = append(response, TopicWithPoll{
			Topic:        t,
			HasPoll:      hasPollMap[t.ID],
			AuthorBadges: userBadgesMap[t.UserID],
		})
	}

	return response, total, nil
}

// GetHomePageWithQuery 支持按forum/tag筛选的首页数据
func GetHomePageWithQuery(w http.ResponseWriter, r *http.Request) {
	forumID, _ := strconv.Atoi(r.URL.Query().Get("forum_id"))
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取基础数据
	forums, _ := homePageForums()
	tags, _ := homePageTags()
	announcements, _ := homePageAnnouncements()

	// 获取话题列表（按筛选条件，不走缓存）
	topics, total, err := fetchHomePageTopicsByQuery(forumID, tagID, sort, page, pageSize)
	if err != nil {
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 获取当前查看者信息（用于匿名判断）
	viewerID, viewerRole, _ := middleware.GetOptionalUserInfo(r)

	// 处理匿名数据
	processedTopics := processTopicsForAnonymity(topics, viewerID, viewerRole)

	homeData := map[string]interface{}{
		"forums":        forums,
		"tags":          tags,
		"announcements": announcements,
		"topics": map[string]interface{}{
			"list":      processedTopics,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"updated_at": time.Now(),
	}

	errors.Success(w, homeData)
}

// processTopicsForAnonymity 处理话题列表的匿名数据
func processTopicsForAnonymity(topics []TopicWithPoll, viewerID uint, viewerRole int) []interface{} {
	type ProcessedTopic struct {
		ID             uint                 `json:"id"`
		Title          string               `json:"title"`
		Content        string               `json:"content"`
		UserID         uint                 `json:"user_id"`
		User           interface{}          `json:"user"`
		ForumID        uint                 `json:"forum_id"`
		Forum          models.Forum         `json:"forum"`
		IsPinned       bool                 `json:"is_pinned"`
		IsUserPinned   bool                 `json:"is_user_pinned"`
		IsLocked       bool                 `json:"is_locked"`
		IsEssence      bool                 `json:"is_essence"`
		IsHidden       bool                 `json:"is_hidden"`
		HotScore       float64              `json:"hot_score"`
		LikeCount      int                  `json:"like_count"`
		ViewCount      int                  `json:"view_count"`
		ReplyCount     int                  `json:"reply_count"`
		LastReplyAt    *time.Time           `json:"last_reply_at"`
		AllowComment   bool                 `json:"allow_comment"`
		CreatedAt      time.Time            `json:"created_at"`
		UpdatedAt      time.Time            `json:"updated_at"`
		IsAnonymous    bool                 `json:"is_anonymous"`
		AnonymousType  models.AnonymousType `json:"anonymous_type"`
		AnonymousUntil *time.Time           `json:"anonymous_until"`
		IsAnonymousEnded bool               `json:"is_anonymous_ended"`
		HasPoll        bool                 `json:"has_poll"`
		AuthorBadges   interface{}          `json:"author_badges"`
		Tags           []models.Tag         `json:"tags"`
	}

	var processedList []interface{}
	for _, t := range topics {
		topic := t.Topic
		shouldShowAnonymous := utils.ShouldShowAnonymous(viewerID, topic.UserID, viewerRole, &topic)

		item := ProcessedTopic{
			ID:               topic.ID,
			Title:            topic.Title,
			Content:          topic.Content,
			UserID:           topic.UserID,
			ForumID:          topic.ForumID,
			Forum:            topic.Forum,
			IsPinned:         topic.IsPinned,
			IsUserPinned:     topic.IsUserPinned,
			IsLocked:         topic.IsLocked,
			IsEssence:        topic.IsEssence,
			IsHidden:         topic.IsHidden,
			HotScore:         topic.HotScore,
			LikeCount:        topic.LikeCount,
			ViewCount:        topic.ViewCount,
			ReplyCount:       topic.ReplyCount,
			LastReplyAt:      topic.LastReplyAt,
			AllowComment:     topic.AllowComment,
			CreatedAt:        topic.CreatedAt,
			UpdatedAt:        topic.UpdatedAt,
			IsAnonymous:      topic.IsAnonymous,
			AnonymousType:    topic.AnonymousType,
			AnonymousUntil:   topic.AnonymousUntil,
			IsAnonymousEnded: topic.IsAnonymousEnded,
			HasPoll:          t.HasPoll,
			Tags:             topic.Tags,
		}

		if shouldShowAnonymous {
			item.User = utils.GetAnonymousUserInfo()
			item.AuthorBadges = []models.UserBadge{}
			item.UserID = 0
		} else {
			item.User = topic.User
			item.AuthorBadges = t.AuthorBadges
		}

		processedList = append(processedList, item)
	}

	return processedList
}

// fetchHomePageTopicsByQuery 获取筛选条件的话题列表
func fetchHomePageTopicsByQuery(forumID, tagID int, sort string, page, pageSize int) ([]TopicWithPoll, int64, error) {
	var topics []models.Topic
	var total int64

	// 构建查询
	query := database.DB.Model(&models.Topic{})
	if forumID > 0 {
		query = query.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		query = query.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	dbQuery := database.DB.Preload("User").Preload("Forum").Preload("Tags")
	if forumID > 0 {
		dbQuery = dbQuery.Where("forum_id = ?", forumID)
	}
	if tagID > 0 {
		dbQuery = dbQuery.Joins("JOIN topic_tags ON topic_tags.topic_id = topics.id").
			Where("topic_tags.tag_id = ?", tagID)
	}

	// 排序方式
	switch sort {
	case "hot":
		dbQuery = dbQuery.Order("is_pinned DESC, (like_count + reply_count * 2) DESC, created_at DESC")
	case "reply":
		dbQuery = dbQuery.Order("is_pinned DESC, last_reply_at DESC NULLS LAST, created_at DESC")
	default:
		dbQuery = dbQuery.Order("is_pinned DESC, created_at DESC")
	}

	// 执行查询
	if err := dbQuery.Offset(offset).Limit(pageSize).Find(&topics).Error; err != nil {
		return nil, 0, err
	}

	// 收集 topic IDs 和 user IDs
	topicIDs := make([]uint, len(topics))
	userIDs := make(map[uint]bool)
	for i, t := range topics {
		topicIDs[i] = t.ID
		userIDs[t.UserID] = true
	}

	// 批量查询投票
	hasPollMap := make(map[uint]bool)
	if len(topicIDs) > 0 {
		var polls []models.Poll
		database.DB.Model(&models.Poll{}).Where("topic_id IN ?", topicIDs).Select("topic_id").Find(&polls)
		for _, p := range polls {
			hasPollMap[p.TopicID] = true
		}
	}

	// 批量查询用户勋章
	userBadgesMap := make(map[uint][]models.UserBadge)
	if len(userIDs) > 0 {
		ids := make([]uint, 0, len(userIDs))
		for id := range userIDs {
			ids = append(ids, id)
		}
		var userBadges []models.UserBadge
		if err := database.DB.Where("user_id IN ? AND is_revoked = ?", ids, false).
			Preload("Badge").
			Find(&userBadges).Error; err != nil {
			log.Printf("fetch home page topics: failed to query user badges, error: %v", err)
		}
		for _, ub := range userBadges {
			userBadgesMap[ub.UserID] = append(userBadgesMap[ub.UserID], ub)
		}
	}

	// 组装返回数据，与 GetTopics 格式一致
	var response []TopicWithPoll
	for _, t := range topics {
		response = append(response, TopicWithPoll{
			Topic:        t,
			HasPoll:      hasPollMap[t.ID],
			AuthorBadges: userBadgesMap[t.UserID],
		})
	}

	return response, total, nil
}
