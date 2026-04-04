package seed

import (
	"bbsgo/database"
	"bbsgo/utils"
	"log"
)

func seedDataZh() {
	seedForumsZh()
	seedConfigsZh()
	seedTagsZh()
	seedUsersZh()
	seedTopicsZh()
}

func seedDataEn() {
	seedForumsEn()
	seedConfigsEn()
	seedTagsEn()
	seedUsersEn()
	seedTopicsEn()
}

func checkAndInsert(table string, count *int64, insertFunc func()) {
	if *count == 0 {
		insertFunc()
	} else {
		log.Printf("[seed] %s already has data, skipping", table)
	}
}

// ========== 版块数据 ==========

func seedForumsZh() {
	var count int64
	database.DB.Model(&struct{}{}).Table("forums").Count(&count)
	checkAndInsert("forums", &count, func() {
		forums := []struct {
			Name        string
			Description string
			SortOrder   int
		}{
			{"技术交流", "编程语言、框架、架构等纯技术讨论", 1},
			{"提问求助", "发帖求助、解答问题", 2},
			{"业界资讯", "科技新闻、技术动态、行业趋势", 3},
			{"资源分享", "工具、教程、电子书、代码片段", 4},
			{"求职招聘", "内推、招聘信息、面经", 5},
			{"灌水闲聊", "生活、娱乐、非技术话题", 6},
			{"站务管理", "公告、反馈、版务", 7},
		}
		for _, f := range forums {
			database.DB.Exec("INSERT INTO forums (name, description, sort_order, allow_post, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
				f.Name, f.Description, f.SortOrder, true)
		}
		log.Println("[seed] forums created (zh)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('forums', 9999)")
	})
}

func seedForumsEn() {
	var count int64
	database.DB.Model(&struct{}{}).Table("forums").Count(&count)
	checkAndInsert("forums", &count, func() {
		forums := []struct {
			Name        string
			Description string
			SortOrder   int
		}{
			{"Tech Talk", "Programming languages, frameworks, architecture discussions", 1},
			{"Q&A", "Post questions and get answers", 2},
			{"Tech News", "Tech news, trends, industry updates", 3},
			{"Resources", "Tools, tutorials, ebooks, code snippets", 4},
			{"Jobs", "Job referrals, recruitments, interviews", 5},
			{"Chat", "Life, entertainment, off-topic", 6},
			{"Meta", "Announcements, feedback, site management", 7},
		}
		for _, f := range forums {
			database.DB.Exec("INSERT INTO forums (name, description, sort_order, allow_post, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
				f.Name, f.Description, f.SortOrder, true)
		}
		log.Println("[seed] forums created (en)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('forums', 9999)")
	})
}

// ========== 配置数据 ==========

func seedConfigsZh() {
	var count int64
	database.DB.Model(&struct{}{}).Table("site_configs").Count(&count)
	checkAndInsert("site_configs", &count, func() {
		configs := []struct {
			Key   string
			Value string
		}{
			{"site_name", "彩虹BBS"},
			{"site_logo", ""},
			{"site_icon", ""},
			{"site_description", "一个现代化的社区论坛系统"},
			{"allow_register", "true"},
			{"allow_post", "true"},
			{"allow_comment", "true"},
			{"email_enabled", "false"},
			{"email_host", ""},
			{"email_port", "465"},
			{"email_user", ""},
			{"email_password", ""},
			{"email_from", ""},
			{"email_from_name", "彩虹BBS"},
			{"qiniu_access_key", ""},
			{"qiniu_secret_key", ""},
			{"qiniu_bucket", ""},
			{"qiniu_domain", ""},
			{"jwt_secret", "change-this-secret-in-production"},
			{"jwt_expire_days", "7"},
			{"credit_signin", "2"},
			{"credit_signin_consecutive", "3"},
			{"cache_num_counters", "10000"},
			{"cache_max_cost", "10000000"},
		}
		for _, c := range configs {
			database.DB.Exec("INSERT INTO site_configs (key, value, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))", c.Key, c.Value)
		}
		log.Println("[seed] site configs created (zh)")
	})
}

func seedConfigsEn() {
	var count int64
	database.DB.Model(&struct{}{}).Table("site_configs").Count(&count)
	checkAndInsert("site_configs", &count, func() {
		configs := []struct {
			Key   string
			Value string
		}{
			{"site_name", "RainbowBBS"},
			{"site_logo", ""},
			{"site_icon", ""},
			{"site_description", "A modern community forum system"},
			{"allow_register", "true"},
			{"allow_post", "true"},
			{"allow_comment", "true"},
			{"email_enabled", "false"},
			{"email_host", ""},
			{"email_port", "465"},
			{"email_user", ""},
			{"email_password", ""},
			{"email_from", ""},
			{"email_from_name", "RainbowBBS"},
			{"qiniu_access_key", ""},
			{"qiniu_secret_key", ""},
			{"qiniu_bucket", ""},
			{"qiniu_domain", ""},
			{"jwt_secret", "change-this-secret-in-production"},
			{"jwt_expire_days", "7"},
			{"credit_signin", "2"},
			{"credit_signin_consecutive", "3"},
			{"cache_num_counters", "10000"},
			{"cache_max_cost", "10000000"},
		}
		for _, c := range configs {
			database.DB.Exec("INSERT INTO site_configs (key, value, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))", c.Key, c.Value)
		}
		log.Println("[seed] site configs created (en)")
	})
}

// ========== 标签数据 ==========

func seedTagsZh() {
	var count int64
	database.DB.Model(&struct{}{}).Table("tags").Count(&count)
	checkAndInsert("tags", &count, func() {
		tags := []struct {
			Name        string
			Icon        string
			Description string
			SortOrder   int
		}{
			{"今日份松弛", "😌", "分享慢生活、拒绝焦虑的瞬间", 1},
			{"爱你老己", "💖", "对自己好的方式、自我关怀", 2},
			{"活人感日常", "🫠", "真实、不完美的生活碎片", 3},
			{"邪修一下", "⚡", "找捷径、高效摆烂、反内卷", 4},
			{"外耗模式", "😤", "与其内耗自己，不如外耗别人", 5},
			{"今日小确幸", "✨", "微小而确定的幸福瞬间", 6},
			{"我的互联网嘴替", "🗣️", "说出了我想说但说不出的话", 7},
			{"求建议/避雷", "❓", "生活求助、消费避坑", 8},
			{"笑死我了", "😂", "搞笑段子、神评论、趣图", 9},
			{"真香现场", "🔥", "打脸时刻、意外真香", 10},
			{"破防了", "💔", "感动、扎心、被戳中的瞬间", 11},
			{"什么水平？", "🤔", "求评价、求鉴定、秀成果", 12},
		}
		for _, t := range tags {
			database.DB.Exec("INSERT INTO tags (name, icon, description, sort_order, usage_count, is_official, is_banned, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))",
				t.Name, t.Icon, t.Description, t.SortOrder, 0, true, false)
		}
		log.Println("[seed] tags created (zh)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('tags', 9999)")
	})
}

func seedTagsEn() {
	var count int64
	database.DB.Model(&struct{}{}).Table("tags").Count(&count)
	checkAndInsert("tags", &count, func() {
		tags := []struct {
			Name        string
			Icon        string
			Description string
			SortOrder   int
		}{
			{"Chill Mode", "😌", "Share slow living, reject anxiety", 1},
			{"Self Love", "💖", "Ways to be good to yourself", 2},
			{"Real Life", "🫠", "Real, imperfect life moments", 3},
			{"Life Hack", "⚡", "Shortcuts, efficient slacking", 4},
			{"Vent Out", "😤", "Externalize stress instead of internalizing", 5},
			{"Little Joy", "✨", "Small but certain moments of happiness", 6},
			{"My Voice", "🗣️", "Says what I wanted to say", 7},
			{"Advice Needed", "❓", "Life advice, shopping tips", 8},
			{"So Funny", "😂", "Funny posts, comments, images", 9},
			{"Mind Blown", "🔥", "Unexpected pleasant surprises", 10},
			{"Touched", "💔", "Moved, stabbed, hit by feelings", 11},
			{"How Is It?", "🤔", "Rate my work, show results", 12},
		}
		for _, t := range tags {
			database.DB.Exec("INSERT INTO tags (name, icon, description, sort_order, usage_count, is_official, is_banned, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))",
				t.Name, t.Icon, t.Description, t.SortOrder, 0, true, false)
		}
		log.Println("[seed] tags created (en)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('tags', 9999)")
	})
}

// ========== 用户数据 ==========

func seedUsersZh() {
	var count int64
	database.DB.Model(&struct{}{}).Table("users").Count(&count)
	checkAndInsert("users", &count, func() {
		adminPassword, _ := utils.HashPassword("12345678")
		database.DB.Exec(`INSERT INTO users (id, username, email, nickname, password_hash, role, credits, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
			10000, "admin", "admin@example.com", "管理员", adminPassword, 2, 10000)
		log.Println("[seed] admin user created (zh): username=admin, password=12345678, id=10000")

		users := []struct {
			ID       int
			Username string
			Email    string
			Nickname string
		}{
			{10001, "testuser1", "test1@example.com", "测试用户1"},
			{10002, "testuser2", "test2@example.com", "测试用户2"},
			{10003, "testuser3", "test3@example.com", "测试用户3"},
			{10004, "testuser4", "test4@example.com", "测试用户4"},
			{10005, "testuser5", "test5@example.com", "测试用户5"},
			{10006, "testuser6", "test6@example.com", "测试用户6"},
			{10007, "testuser7", "test7@example.com", "测试用户7"},
			{10008, "testuser8", "test8@example.com", "测试用户8"},
			{10009, "testuser9", "test9@example.com", "测试用户9"},
			{10010, "testuser10", "test10@example.com", "测试用户10"},
		}
		for _, u := range users {
			hashedPassword, _ := utils.HashPassword("123456")
			database.DB.Exec(`INSERT INTO users (id, username, email, nickname, password_hash, credits, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
				u.ID, u.Username, u.Email, u.Nickname, hashedPassword, 1000)
		}
		log.Println("[seed] test users created (zh)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('users', 10010)")
	})
}

func seedUsersEn() {
	var count int64
	database.DB.Model(&struct{}{}).Table("users").Count(&count)
	checkAndInsert("users", &count, func() {
		adminPassword, _ := utils.HashPassword("12345678")
		database.DB.Exec(`INSERT INTO users (id, username, email, nickname, password_hash, role, credits, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
			10000, "admin", "admin@example.com", "Administrator", adminPassword, 2, 10000)
		log.Println("[seed] admin user created (en): username=admin, password=12345678, id=10000")

		users := []struct {
			ID       int
			Username string
			Email    string
			Nickname string
		}{
			{10001, "testuser1", "test1@example.com", "Test User 1"},
			{10002, "testuser2", "test2@example.com", "Test User 2"},
			{10003, "testuser3", "test3@example.com", "Test User 3"},
			{10004, "testuser4", "test4@example.com", "Test User 4"},
			{10005, "testuser5", "test5@example.com", "Test User 5"},
			{10006, "testuser6", "test6@example.com", "Test User 6"},
			{10007, "testuser7", "test7@example.com", "Test User 7"},
			{10008, "testuser8", "test8@example.com", "Test User 8"},
			{10009, "testuser9", "test9@example.com", "Test User 9"},
			{10010, "testuser10", "test10@example.com", "Test User 10"},
		}
		for _, u := range users {
			hashedPassword, _ := utils.HashPassword("123456")
			database.DB.Exec(`INSERT INTO users (id, username, email, nickname, password_hash, credits, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
				u.ID, u.Username, u.Email, u.Nickname, hashedPassword, 1000)
		}
		log.Println("[seed] test users created (en)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('users', 10010)")
	})
}

// ========== 话题数据 ==========

func seedTopicsZh() {
	var count int64
	database.DB.Model(&struct{}{}).Table("topics").Count(&count)
	checkAndInsert("topics", &count, func() {
		topics := []struct {
			Title      string
			Content    string
			UserID     uint
			ForumID    uint
			LikeCount  int
			ReplyCount int
			ViewCount  int
		}{
			{"bbs-go v3.5.0 发布，升级 go1.18", "文档地址：https://docs.bbs-go.com/\n官网交流：https://mlog.club\n问题反馈：https://mlog.club/topic/node/3\n\n本次更新内容：\n1. 升级 Go 1.18 版本\n2. 优化数据库查询性能\n3. 修复已知 bug\n4. 新增配置管理功能", 10001, 2, 12, 8, 352},
			{"Vue3 + TypeScript 项目实践分享", "最近用 Vue3 + TypeScript 做了一个项目，分享一些实践经验：\n\n1. 组合式 API 真的很香，逻辑复用更方便了\n2. TypeScript 的类型推导需要好好配置\n3. Pinia 比 Vuex 更简洁好用\n\n有问题的朋友欢迎留言讨论~", 10002, 2, 45, 18, 892},
			{"今天天气不错，适合摸鱼", "周末到了，阳光明媚，正是摸鱼好时节。大家最近都在看什么书？有什么好剧推荐吗？\n\n我最近在看《三体》，真的很精彩！强烈推荐给还没看过的朋友。", 10003, 7, 32, 21, 687},
			{"求助：MySQL 慢查询优化", "公司有个 MySQL 表数据量大概 500 万，查询越来越慢了。\n\n表结构大概是：\n- id (主键)\n- user_id (索引)\n- created_at (索引)\n- content (text)\n\n查询语句：SELECT * FROM table WHERE user_id = ? ORDER BY created_at DESC LIMIT 20\n\n请问有什么优化建议吗？", 10004, 3, 8, 12, 234},
			{"分享一些 Linux 常用命令", "整理了一些常用的 Linux 命令，希望对大家有帮助：\n\n查看端口占用情况：\nnetstat -tunlp | grep 端口号\n\n通过 ssh 将远程端口映射到本地端口：\nssh -L 13306:127.0.0.1:3306 用户名@远程地址 -N\n\n这样远程服务器就不需要开放需要的端口到公网了，更安全。", 10005, 5, 18, 5, 221},
			{"C++ 程序返回 value 3221226356 求教！", "return value 3221226356 求教求教！\n\n#include <iostream>\nusing namespace std;\nint main() {\n    int n;\n    double *p=new double[n];\n    cin>>n;\n    for(int i=0;i<n;i++) { cin>>p[i]; }\n    for(int i=0;i<n;i++) { cout<<p[i]<<\" \"; }\n    return 0;\n}\n\n程序运行时出现这个错误，请问是什么原因？", 10006, 3, 3, 7, 126},
			{"分享一张今天拍的美照", "今天去公园玩了，随手拍了一张照片，分享给大家~\n\n[图片]\n\n摄影器材：Sony A7M3\n参数：f/2.8, 1/500s, ISO100", 10007, 8, 156, 43, 2341},
			{"网站有个 BUG 反馈", "在使用网站时发现一个问题：\n\n当我在移动端浏览帖子时，点击回复按钮后键盘会遮挡输入框，需要手动收起键盘才能看到输入内容。\n\n浏览器：Safari\n系统：iOS 16\n设备：iPhone 13 Pro\n\n希望能修复一下，谢谢！", 10008, 8, 5, 3, 89},
			{"推荐一个很好用的开源项目", "最近发现一个很棒的开源项目：\n\n项目名称：VSCode\nGitHub 地址：https://github.com/microsoft/vscode\n\n功能强大，插件生态丰富，支持几乎所有编程语言。强烈推荐给各位开发者！\n\n大家还有什么好用的工具欢迎分享~", 10009, 5, 67, 29, 1523},
			{"2024 年前端技术趋势预测", "随着 AI 的快速发展，前端领域也在不断变化。以下是我对 2024 年前端技术趋势的一些预测：\n\n1. AI 辅助开发将成为标配\n2. Server Components 会更加流行\n3. TypeScript 使用率继续上升\n4. Rust 在前端工具链中的应用会更广泛\n5. Web Components 可能会迎来第二春\n\n大家怎么看？欢迎讨论！", 10010, 4, 89, 32, 1523},
		}
		for i, t := range topics {
			database.DB.Exec(`INSERT INTO topics (id, title, content, user_id, forum_id, like_count, reply_count, view_count, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now', '-'||abs(random())%72||' hours'), datetime('now', '-'||abs(random())%72||' hours'))`,
				10000+i, t.Title, t.Content, t.UserID, t.ForumID, t.LikeCount, t.ReplyCount, t.ViewCount)
		}
		log.Println("[seed] topics created (zh)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('topics', 10010)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('comments', 9999)")
	})
}

func seedTopicsEn() {
	var count int64
	database.DB.Model(&struct{}{}).Table("topics").Count(&count)
	checkAndInsert("topics", &count, func() {
		topics := []struct {
			Title      string
			Content    string
			UserID     uint
			ForumID    uint
			LikeCount  int
			ReplyCount int
			ViewCount  int
		}{
			{"bbs-go v3.5.0 Released, upgraded to go1.18", "Docs: https://docs.bbs-go.com/\nForum: https://mlog.club\nFeedback: https://mlog.club/topic/node/3\n\nChangelog:\n1. Upgraded to Go 1.18\n2. Optimized database query performance\n3. Fixed known bugs\n4. Added configuration management", 10001, 2, 12, 8, 352},
			{"Vue3 + TypeScript Project Practice", "Recently built a project with Vue3 + TypeScript, sharing some experiences:\n\n1. Composition API is great, logic reuse is more convenient\n2. TypeScript type inference needs proper configuration\n3. Pinia is simpler and better than Vuex\n\nFeel free to leave comments if you have questions!", 10002, 2, 45, 18, 892},
			{"Nice weather today, perfect for slacking off", "It's the weekend, sunny and bright, perfect time to slack off. What books are you reading lately? Any good shows to recommend?\n\nI'm reading \"The Three-Body Problem\" recently, it's really exciting! Highly recommend to those who haven't read it.", 10003, 7, 32, 21, 687},
			{"Help: MySQL Slow Query Optimization", "Our MySQL table has about 5 million rows and queries are getting slower.\n\nTable structure:\n- id (primary key)\n- user_id (index)\n- created_at (index)\n- content (text)\n\nQuery: SELECT * FROM table WHERE user_id = ? ORDER BY created_at DESC LIMIT 20\n\nAny optimization suggestions?", 10004, 3, 8, 12, 234},
			{"Sharing Some Common Linux Commands", "Compiled some common Linux commands, hope it helps:\n\nCheck port usage:\nnetstat -tunlp | grep port\n\nMap remote port to local via ssh:\nssh -L 13306:127.0.0.1:3306 user@remote -N\n\nThis way the remote server doesn't need to expose ports to the internet, more secure.", 10005, 5, 18, 5, 221},
			{"C++ Program Returns value 3221226356 Help!", "return value 3221226356 please help!\n\n#include <iostream>\nusing namespace std;\nint main() {\n    int n;\n    double *p=new double[n];\n    cin>>n;\n    for(int i=0;i<n;i++) { cin>>p[i]; }\n    for(int i=0;i<n;i++) { cout<<p[i]<<\" \"; }\n    return 0;\n}\n\nWhat causes this error?", 10006, 3, 3, 7, 126},
			{"Sharing a beautiful photo I took today", "Went to the park today, took a photo, sharing with everyone~\n\n[Image]\n\nCamera: Sony A7M3\nSettings: f/2.8, 1/500s, ISO100", 10007, 8, 156, 43, 2341},
			{"Website BUG Feedback", "Found an issue when using the website:\n\nWhen browsing posts on mobile, clicking the reply button the keyboard blocks the input box, need to manually close the keyboard to see the content.\n\nBrowser: Safari\nOS: iOS 16\nDevice: iPhone 13 Pro\n\nHope it can be fixed, thanks!", 10008, 8, 5, 3, 89},
			{"Recommending a Great Open Source Project", "Found a great open source project recently:\n\nProject Name: VSCode\nGitHub: https://github.com/microsoft/vscode\n\nPowerful, rich plugin ecosystem, supports almost all programming languages. Highly recommended for developers!\n\nWhat other useful tools do you recommend?", 10009, 5, 67, 29, 1523},
			{"2024 Frontend Technology Trends", "With the rapid development of AI, the frontend landscape is constantly changing. Here are my predictions for 2024 frontend trends:\n\n1. AI-assisted development will become standard\n2. Server Components will become more popular\n3. TypeScript usage will continue to rise\n4. Rust will be more widely used in frontend toolchains\n5. Web Components might see a resurgence\n\nWhat do you think? Feel free to discuss!", 10010, 4, 89, 32, 1523},
		}
		for i, t := range topics {
			database.DB.Exec(`INSERT INTO topics (id, title, content, user_id, forum_id, like_count, reply_count, view_count, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now', '-'||abs(random())%72||' hours'), datetime('now', '-'||abs(random())%72||' hours'))`,
				10000+i, t.Title, t.Content, t.UserID, t.ForumID, t.LikeCount, t.ReplyCount, t.ViewCount)
		}
		log.Println("[seed] topics created (en)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('topics', 10010)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('comments', 9999)")
	})
}
