package routes

import (
	"bbsgo/handlers"
	"bbsgo/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes 配置并返回路由实例
// 设置所有 API 路由规则，包括公开接口、认证接口和管理员接口
func SetupRoutes() *mux.Router {
	// 创建路由实例
	r := mux.NewRouter()

	// 应用 CORS 中间件
	r.Use(middleware.CORS)

	// 全局 OPTIONS 处理器，处理 CORS 预检请求
	r.HandleFunc("/api/v1/{rest:.*}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	// ========== API v1 版本 ==========
	api := r.PathPrefix("/api/v1").Subrouter()

	// ========== 公开接口（无需认证）==========

	// 用户注册和登录
	api.HandleFunc("/register", handlers.RegisterWithCode).Methods("POST")      // 邮箱注册
	api.HandleFunc("/send-code", handlers.SendVerificationCode).Methods("POST") // 发送验证码
	api.HandleFunc("/login", handlers.Login).Methods("POST")                    // 登录

	// 公开内容获取
	api.HandleFunc("/forums", handlers.GetForums).Methods("GET")                      // 获取版块列表
	api.HandleFunc("/config", handlers.GetSiteConfig).Methods("GET")                  // 获取网站配置
	api.HandleFunc("/topics", handlers.GetTopics).Methods("GET")                      // 获取话题列表
	api.HandleFunc("/topics/{id}", handlers.GetTopic).Methods("GET")                  // 获取话题详情
	api.HandleFunc("/topics/{id}/comments", handlers.GetComments).Methods("GET")      // 获取话题评论
	api.HandleFunc("/tags", handlers.GetTags).Methods("GET")                          // 获取标签列表
	api.HandleFunc("/tags/search", handlers.SearchTags).Methods("GET")                // 搜索标签
	api.HandleFunc("/tags/{id}", handlers.GetTag).Methods("GET")                      // 获取标签详情
	api.HandleFunc("/announcements", handlers.GetAnnouncements).Methods("GET")        // 获取公告列表
	api.HandleFunc("/users/credit", handlers.GetCreditUsers).Methods("GET")           // 获取积分排行
	api.HandleFunc("/users/search", handlers.SearchUsers).Methods("GET")              // 搜索用户
	api.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")                    // 获取用户公开信息
	api.HandleFunc("/users/{id}/stats", handlers.GetUserStats).Methods("GET")         // 获取用户统计
	api.HandleFunc("/users/{id}/followers", handlers.GetUserFollowers).Methods("GET") // 获取用户粉丝
	api.HandleFunc("/users/{id}/topics", handlers.GetUserTopics).Methods("GET")       // 获取用户话题
	api.HandleFunc("/users/{id}/badges", handlers.GetUserBadgesByID).Methods("GET")  // 获取指定用户勋章
	api.HandleFunc("/search", handlers.Search).Methods("GET")                         // 搜索

	// 投票公开接口
	api.HandleFunc("/polls/{id}", handlers.GetPoll).Methods("GET")                    // 获取投票详情
	api.HandleFunc("/topics/{topic_id}/poll", handlers.GetPollByTopic).Methods("GET") // 根据话题获取投票

	// ========== 认证接口（需要登录）==========
	auth := api.PathPrefix("").Subrouter()
	auth.Use(middleware.Auth) // 应用认证中间件

	// 用户个人资料
	auth.HandleFunc("/user/profile", handlers.GetProfile).Methods("GET")            // 获取个人资料
	auth.HandleFunc("/user/profile", handlers.UpdateProfile).Methods("PUT")         // 更新个人资料
	auth.HandleFunc("/user/topics", handlers.GetCurrentUserTopics).Methods("GET")   // 获取我的话题
	auth.HandleFunc("/user/signin", handlers.SignIn).Methods("POST")                // 签到
	auth.HandleFunc("/user/signin/status", handlers.GetSignInStatus).Methods("GET") // 签到状态
	auth.HandleFunc("/user/favorites", handlers.GetFavorites).Methods("GET")        // 获取收藏列表
	auth.HandleFunc("/user/follows", handlers.GetFollows).Methods("GET")            // 获取关注列表
	auth.HandleFunc("/user/followers", handlers.GetFollowers).Methods("GET")        // 获取粉丝列表
	auth.HandleFunc("/user/follow-topics", handlers.GetFollowTopics).Methods("GET") // 获取关注动态
	auth.HandleFunc("/user/badges", handlers.GetUserBadges).Methods("GET")          // 获取用户勋章
	auth.HandleFunc("/user/reports", handlers.GetUserReports).Methods("GET")        // 获取我的举报

	// 话题操作
	auth.HandleFunc("/topics", handlers.CreateTopic).Methods("POST")          // 创建话题
	auth.HandleFunc("/topics/{id}", handlers.UpdateTopic).Methods("PUT")      // 更新话题
	auth.HandleFunc("/topics/{id}", handlers.DeleteTopic).Methods("DELETE")   // 删除话题
	auth.HandleFunc("/topics/{id}/pin", handlers.UserPinTopic).Methods("PUT") // 作者置顶/取消置顶

	// 评论操作
	auth.HandleFunc("/topics/{id}/comments", handlers.CreateComment).Methods("POST")                      // 创建评论
	auth.HandleFunc("/comments/{id}", handlers.UpdateComment).Methods("PUT")                              // 更新评论
	auth.HandleFunc("/comments/{id}", handlers.DeleteComment).Methods("DELETE")                           // 删除评论
	auth.HandleFunc("/topics/{topic_id}/comments/{comment_id}/pin", handlers.PinComment).Methods("PUT")   // 置顶/取消置顶评论
	auth.HandleFunc("/topics/{topic_id}/comments/{comment_id}/best", handlers.BestComment).Methods("PUT") // 标记/取消最佳评论

	// 点赞操作
	auth.HandleFunc("/likes", handlers.CreateLike).Methods("POST")      // 创建点赞
	auth.HandleFunc("/likes", handlers.DeleteLike).Methods("DELETE")    // 删除点赞
	auth.HandleFunc("/likes/check", handlers.CheckLike).Methods("POST") // 检查点赞状态

	// 收藏操作
	auth.HandleFunc("/favorites", handlers.CreateFavorite).Methods("POST")      // 创建收藏
	auth.HandleFunc("/favorites", handlers.DeleteFavorite).Methods("DELETE")    // 删除收藏
	auth.HandleFunc("/favorites/check", handlers.CheckFavorite).Methods("POST") // 检查收藏状态

	// 关注操作
	auth.HandleFunc("/follows", handlers.CreateFollow).Methods("POST")     // 创建关注
	auth.HandleFunc("/follows", handlers.DeleteFollow).Methods("DELETE")   // 删除关注
	auth.HandleFunc("/follows/check", handlers.CheckFollow).Methods("GET") // 检查关注状态

	// 私信操作
	auth.HandleFunc("/messages", handlers.GetMessages).Methods("GET")                           // 获取私信列表
	auth.HandleFunc("/messages", handlers.SendMessage).Methods("POST")                          // 发送私信
	auth.HandleFunc("/messages/unread-count", handlers.GetUnreadMessageCount).Methods("GET")    // 未读数量
	auth.HandleFunc("/messages/read", handlers.MarkMessagesRead).Methods("PUT")                 // 标记已读
	auth.HandleFunc("/messages/with/{user_id}", handlers.GetMessageConversation).Methods("GET") // 会话详情

	// 通知操作
	auth.HandleFunc("/notifications", handlers.GetNotifications).Methods("GET")                        // 获取通知列表
	auth.HandleFunc("/notifications/unread-count", handlers.GetUnreadNotificationCount).Methods("GET") // 未读数量
	auth.HandleFunc("/notifications/read-all", handlers.MarkAllNotificationsRead).Methods("PUT")       // 全部已读

	// 草稿箱操作
	auth.HandleFunc("/drafts", handlers.GetDrafts).Methods("GET")           // 获取草稿列表
	auth.HandleFunc("/drafts", handlers.CreateDraft).Methods("POST")        // 创建草稿
	auth.HandleFunc("/drafts/{id}", handlers.GetDraft).Methods("GET")       // 获取草稿详情
	auth.HandleFunc("/drafts/{id}", handlers.UpdateDraft).Methods("PUT")    // 更新草稿
	auth.HandleFunc("/drafts/{id}", handlers.DeleteDraft).Methods("DELETE") // 删除草稿

	// 举报操作
	auth.HandleFunc("/reports", handlers.CreateReport).Methods("POST") // 创建举报

	// 勋章操作
	auth.HandleFunc("/badges", handlers.GetBadges).Methods("GET")                     // 获取勋章列表
	auth.HandleFunc("/badges/progress", handlers.GetUserBadgeProgress).Methods("GET") // 获取用户勋章进度

	// 文件上传
	auth.HandleFunc("/upload", handlers.UploadFile).Methods("POST")           // 上传文件
	auth.HandleFunc("/upload/check", handlers.CheckFileExists).Methods("GET") // 检查文件是否存在（秒传）

	// 投票操作
	auth.HandleFunc("/polls", handlers.CreatePoll).Methods("POST")      // 创建投票
	auth.HandleFunc("/polls/vote", handlers.SubmitVote).Methods("POST") // 提交投票

	// ========== 管理后台接口（需要管理员权限）==========
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.Auth)      // 认证中间件
	admin.Use(middleware.AdminAuth) // 管理员权限中间件

	// 用户管理
	admin.HandleFunc("/users", handlers.GetAdminUsers).Methods("GET")            // 获取用户列表
	admin.HandleFunc("/users/{id}/role", handlers.UpdateUserRole).Methods("PUT") // 更新用户角色
	admin.HandleFunc("/users/{id}/ban", handlers.BanUser).Methods("PUT")         // 封禁用户
	admin.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")       // 删除用户

	// 版块管理
	admin.HandleFunc("/forums", handlers.CreateForum).Methods("POST")        // 创建版块
	admin.HandleFunc("/forums/{id}", handlers.UpdateForum).Methods("PUT")    // 更新版块
	admin.HandleFunc("/forums/{id}", handlers.DeleteForum).Methods("DELETE") // 删除版块

	// 标签管理
	admin.HandleFunc("/tags", handlers.GetAdminTags).Methods("GET")      // 获取标签列表
	admin.HandleFunc("/tags", handlers.CreateTag).Methods("POST")        // 创建标签
	admin.HandleFunc("/tags/merge", handlers.MergeTags).Methods("POST")  // 合并标签
	admin.HandleFunc("/tags/{id}", handlers.UpdateTag).Methods("PUT")    // 更新标签
	admin.HandleFunc("/tags/{id}", handlers.DeleteTag).Methods("DELETE") // 删除标签

	// 话题管理
	admin.HandleFunc("/topics", handlers.GetAdminTopics).Methods("GET")           // 获取话题列表
	admin.HandleFunc("/topics/{id}", handlers.DeleteAdminTopic).Methods("DELETE") // 删除话题
	admin.HandleFunc("/topics/{id}/pin", handlers.AdminPinTopic).Methods("PUT")   // 管理员置顶/取消置顶

	// 评论管理
	admin.HandleFunc("/comments", handlers.GetAdminComments).Methods("GET")           // 获取评论列表
	admin.HandleFunc("/comments/{id}", handlers.DeleteAdminComment).Methods("DELETE") // 删除评论

	// 举报管理
	admin.HandleFunc("/reports", handlers.GetAdminReports).Methods("GET")          // 获取举报列表
	admin.HandleFunc("/reports/{id}/handle", handlers.HandleReport).Methods("PUT") // 处理举报

	// 公告管理
	admin.HandleFunc("/announcements", handlers.CreateAnnouncement).Methods("POST")        // 创建公告
	admin.HandleFunc("/announcements/{id}", handlers.UpdateAnnouncement).Methods("PUT")    // 更新公告
	admin.HandleFunc("/announcements/{id}", handlers.DeleteAnnouncement).Methods("DELETE") // 删除公告

	// 网站配置
	admin.HandleFunc("/config", handlers.UpdateSiteConfig).Methods("PUT") // 更新配置

	// 防刷系统配置
	admin.HandleFunc("/antispam/config", handlers.GetAntiSpamConfig).Methods("GET")        // 获取防刷配置
	admin.HandleFunc("/antispam/config", handlers.UpdateAntiSpamConfig).Methods("POST")   // 更新防刷配置
	admin.HandleFunc("/antispam/stats", handlers.GetAntiSpamStats).Methods("GET")         // 获取防刷统计
	admin.HandleFunc("/antispam/keywords", handlers.GetSpamKeywords).Methods("GET")        // 获取敏感词列表
	admin.HandleFunc("/antispam/keywords", handlers.AddSpamKeyword).Methods("POST")       // 添加敏感词
	admin.HandleFunc("/antispam/keywords", handlers.DeleteSpamKeyword).Methods("DELETE")   // 删除敏感词
	admin.HandleFunc("/users/{id}/reputation", handlers.AdjustUserReputation).Methods("POST") // 调整用户信誉分
	admin.HandleFunc("/users/{id}/unban", handlers.UnbanUser).Methods("POST")             // 解禁用户
	admin.HandleFunc("/users/{id}/ban-status", handlers.GetUserBanStatus).Methods("GET")  // 获取用户禁言状态
	admin.HandleFunc("/users/{id}/reputation-logs", handlers.GetUserReputationLogs).Methods("GET") // 获取信誉分日志
	admin.HandleFunc("/users/{id}/ban", handlers.AdminBanUser).Methods("POST")            // 管理员禁言用户

	// 版块分类管理
	admin.HandleFunc("/forum-categories", handlers.GetAllForumCategories).Methods("GET")       // 获取分类列表
	admin.HandleFunc("/forum-categories", handlers.CreateForumCategory).Methods("POST")        // 创建分类
	admin.HandleFunc("/forum-categories/{id}", handlers.UpdateForumCategory).Methods("PUT")    // 更新分类
	admin.HandleFunc("/forum-categories/{id}", handlers.DeleteForumCategory).Methods("DELETE") // 删除分类

	// 管理员密码
	admin.HandleFunc("/change-password", handlers.ChangeAdminPassword).Methods("POST") // 修改密码

	// 投票管理
	admin.HandleFunc("/polls", handlers.GetAdminPolls).Methods("GET")      // 获取投票列表
	admin.HandleFunc("/polls/{id}", handlers.UpdatePoll).Methods("PUT")    // 更新投票
	admin.HandleFunc("/polls/{id}/end", handlers.EndPoll).Methods("POST")  // 结束投票
	admin.HandleFunc("/polls/{id}", handlers.DeletePoll).Methods("DELETE") // 删除投票

	// 勋章管理
	admin.HandleFunc("/badges", handlers.GetAdminBadges).Methods("GET")           // 获取勋章列表
	admin.HandleFunc("/badges", handlers.CreateBadge).Methods("POST")             // 创建勋章
	admin.HandleFunc("/badges/init", handlers.InitBadges).Methods("POST")         // 初始化勋章
	admin.HandleFunc("/badges/{id}", handlers.UpdateBadge).Methods("PUT")         // 更新勋章
	admin.HandleFunc("/badges/{id}", handlers.DeleteBadge).Methods("DELETE")      // 删除勋章
	admin.HandleFunc("/badges/{id}/users", handlers.GetBadgeUsers).Methods("GET") // 获取勋章用户列表
	admin.HandleFunc("/badges/award", handlers.AwardBadge).Methods("POST")        // 授予勋章
	admin.HandleFunc("/badges/{id}/revoke", handlers.RevokeBadge).Methods("PUT")  // 撤销勋章

	admin.HandleFunc("/follows", handlers.GetAdminFollows).Methods("GET")           // 获取关注列表
	admin.HandleFunc("/followers", handlers.GetAdminFollowers).Methods("GET")       // 获取粉丝列表
	admin.HandleFunc("/follows/{id}", handlers.DeleteAdminFollow).Methods("DELETE") // 删除关注

	admin.HandleFunc("/best-comments", handlers.GetAdminBestComments).Methods("GET")   // 获取最佳评论列表
	admin.HandleFunc("/comments/{id}/best", handlers.UpdateCommentBest).Methods("PUT") // 更新最佳评论状态

	return r
}
