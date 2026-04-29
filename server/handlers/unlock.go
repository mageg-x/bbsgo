package handlers

import (
	"bbsgo/database"
	"bbsgo/models"
	"log"
)

// checkAndCreateUnlockRecord 检查并创建解锁记录
// 当用户完成指定互动（点赞或评论）后，检查是否满足解锁条件
// 如果满足条件且用户尚未解锁，则创建解锁记录
func checkAndCreateUnlockRecord(userID uint, topicID uint, actionType string) {
	// 查询话题信息
	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		log.Printf("check unlock: topic not found, topicID: %d, error: %v", topicID, err)
		return
	}

	// 检查是否启用了解锁功能
	if !topic.IsUnlockEnabled {
		return
	}

	// 检查是否是作者本人（作者不需要解锁）
	if userID == topic.UserID {
		return
	}

	// 检查用户是否已经解锁
	var existingUnlock models.UnlockRecord
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&existingUnlock).Error; err == nil {
		// 已经解锁，直接返回
		return
	}

	// 检查是否满足解锁条件
	var isUnlocked bool
	var unlockType string

	switch topic.UnlockType {
	case "like":
		// 只需要点赞
		if actionType == "like" && topic.LikeCount >= topic.UnlockLikeCount {
			isUnlocked = true
			unlockType = "like"
		}
	case "comment":
		// 只需要评论
		if actionType == "comment" && topic.ReplyCount >= topic.UnlockCommentCount {
			isUnlocked = true
			unlockType = "comment"
		}
	case "both":
		// 需要两者都满足
		likeSatisfied := topic.LikeCount >= topic.UnlockLikeCount
		commentSatisfied := topic.ReplyCount >= topic.UnlockCommentCount
		
		// 检查当前操作是否帮助满足了条件
		if actionType == "like" {
			// 如果是点赞操作，检查点赞后是否满足条件
			likeSatisfied = (topic.LikeCount) >= topic.UnlockLikeCount
		} else if actionType == "comment" {
			// 如果是评论操作，检查评论后是否满足条件
			commentSatisfied = (topic.ReplyCount) >= topic.UnlockCommentCount
		}
		
		if likeSatisfied && commentSatisfied {
			isUnlocked = true
			unlockType = "both"
		}
	}

	// 如果满足解锁条件，创建解锁记录
	if isUnlocked {
		unlockRecord := models.UnlockRecord{
			UserID:     userID,
			TopicID:    topicID,
			UnlockType: unlockType,
		}

		if err := database.DB.Create(&unlockRecord).Error; err != nil {
			log.Printf("check unlock: failed to create unlock record, userID: %d, topicID: %d, error: %v", userID, topicID, err)
			return
		}

		log.Printf("check unlock: unlock record created successfully, userID: %d, topicID: %d, unlockType: %s", userID, topicID, unlockType)
	}
}
