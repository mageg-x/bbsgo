package seed

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
)

func seedDataZh() {
	seedForumsZh()
	seedConfigsZh()
	seedTagsZh()
	seedUsersZh()
	seedTopicsZh()
	seedTopicTagsZh()
}

func seedDataEn() {
	seedForumsEn()
	seedConfigsEn()
	seedTagsEn()
	seedUsersEn()
	seedTopicsEn()
	seedTopicTagsEn()
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
			{"credit_topic", "5"},
			{"credit_post", "1"},
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
			{"credit_topic", "5"},
			{"credit_post", "1"},
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
			{"Chill Vibes", "😌", "Share slow living, reject anxiety", 1},
			{"Love Yourself", "💖", "Self-care and ways to be kind to yourself", 2},
			{"Real Life", "🫠", "Authentic, imperfect life fragments", 3},
			{"Hack It", "⚡", "Find shortcuts, efficient slacking, anti-hustle", 4},
			{"Externalize", "😤", "Instead of internalizing, externalize", 5},
			{"Small Wins", "✨", "Tiny but certain happy moments", 6},
			{"Mouthpiece", "🗣️", "Said what I wanted to but couldn't", 7},
			{"Advice/Avoid", "❓", "Life help, consumer traps", 8},
			{"LOL", "😂", "Funny jokes, hot takes, memes", 9},
			{"Plot Twist", "🔥", "Facepalm moments, unexpected wins", 10},
			{"Emotional", "💔", "Touched, heartbroken, relatable moments", 11},
			{"Rate This", "🤔", "Ask for evaluation, show achievements", 12},
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
			{10001, "yuexia", "yuexia@example.com", "月下独酌"},
			{10002, "xingchen", "xingchen@example.com", "星辰大海"},
			{10003, "yulan", "yulan@example.com", "雨后玉兰"},
			{10004, "moli", "moli@example.com", "茉莉清茶"},
			{10005, "yunqi", "yunqi@example.com", "云起之时"},
			{10006, "fengqing", "fengqing@example.com", "风轻云淡"},
			{10007, "huakai", "huakai@example.com", "花开半夏"},
			{10008, "xueying", "xueying@example.com", "雪映梅花"},
			{10009, "mengxing", "mengxing@example.com", "梦醒时分"},
			{10010, "chenguang", "chenguang@example.com", "晨曦微露"},
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
			{10001, "alex", "alex@example.com", "Alex Chen"},
			{10002, "emma", "emma@example.com", "Emma Wang"},
			{10003, "mike", "mike@example.com", "Mike Li"},
			{10004, "sarah", "sarah@example.com", "Sarah Zhang"},
			{10005, "david", "david@example.com", "David Liu"},
			{10006, "jessica", "jessica@example.com", "Jessica Wu"},
			{10007, "tom", "tom@example.com", "Tom Chen"},
			{10008, "amy", "amy@example.com", "Amy Huang"},
			{10009, "developer", "dev@example.com", "Code Ninja"},
			{10010, "dreamer", "dream@example.com", "Dream Chaser"},
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
			{"今天终于躺平了，感觉人生都升华了", "终于下定决心不卷了！早上睡到自然醒，慢悠悠地吃了个早餐，晒晒太阳看看书。\n\n![躺平日常](https://picsum.photos/seed/1/800/450)\n\n以前总觉得要努力要奋斗，现在才发现，允许自己不那么完美，也是一种自我关怀。\n\n推荐一个超治愈的视频：https://www.youtube.com/watch?v=dQw4w9WgXcQ\n\n今天也是爱自己的一天！💖", 10001, 1, 234, 45, 3241},
			{"分享我的极简护肤流程，对自己好一点", "以前总买一堆护肤品，现在简化了，皮肤反而更好了！\n\n![护肤日常](https://picsum.photos/seed/2/800/450)\n\n晨间流程：\n1. 温水洁面\n2. 爽肤水\n3. 精华\n4. 防晒\n\n晚间流程：\n1. 卸妆+洁面\n2. 爽肤水\n3. 精华\n4. 面霜\n\n最重要的是：多喝水，早睡！", 10002, 1, 189, 32, 2891},
			{"记录一下今天的小确幸", "今天下班路上看到了超美的夕阳！\n\n![夕阳](https://picsum.photos/seed/3/800/450)\n\n生活中的小确幸：\n- 买到了喜欢的奶茶\n- 地铁刚好赶上\n- 今天的饭菜很好吃\n- 同事夸我今天好看\n- 回家路上看到了好看的云\n\n这些小小的瞬间，组成了我们的生活呀！", 10001, 1, 178, 29, 2456},
			{"今天摸鱼了一天，太爽了", "今天工作效率特别低，索性就摸鱼了！\n\n![摸鱼](https://picsum.photos/seed/4/800/450)\n\n摸鱼清单：\n- 刷了半小时小红书\n- 看了几个短视频\n- 跟同事聊了八卦\n- 喝了一杯咖啡\n- 摸了摸鱼（真的在看鱼缸里的鱼）\n\n偶尔摸鱼，心情真的会好很多！", 10003, 1, 201, 34, 2876},
			{"今天把锅甩出去了，心情舒畅", "以前总怕得罪人，什么锅都自己背。今天终于学会了：这锅我不背！\n\n![甩锅](https://picsum.photos/seed/5/800/450)\n\n心得：\n1. 明确自己的职责范围\n2. 不是自己的问题，礼貌但坚定地拒绝\n3. 提供解决方案，但不替别人承担责任\n\n与其内耗自己，不如外耗别人！", 10004, 2, 56, 23, 1234},
			{"今天发生了一件超好笑的事", "今天在地铁上，有个小朋友指着我的头发说：阿姨，你的头发好像棉花糖！\n\n![搞笑](https://picsum.photos/seed/6/800/450)\n\n我：\n\n然后旁边的人都笑了，小朋友还特别认真地问我能不能吃一口。\n\n小朋友的脑洞真的太大了，笑哭！😂", 10005, 2, 89, 34, 2134},
			{"今天看了一个超感人的视频，破防了", "今天刷到一个视频，是关于爷爷奶奶的爱情故事，看得我眼泪直流。\n\n![感动](https://picsum.photos/seed/7/800/450)\n\n视频里，爷爷得了老年痴呆，忘了很多事，但还记得奶奶爱吃的糖。\n\n这种长久的陪伴，真的太好哭了。\n\n珍惜眼前人啊！", 10006, 3, 123, 45, 3456},
			{"之前说再也不买了，结果又真香了", "上个月还在说：今年再也不买衣服了！结果这周就...\n\n![真香](https://picsum.photos/seed/8/800/450)\n\n新衣服真的太好看了！忍不住就下单了。\n\n打脸来得太快就像龙卷风。\n\n没关系，钱没有消失，它只是换了一种方式陪在我身边！", 10007, 3, 98, 38, 2891},
			{"求大家帮我看看，这套穿搭什么水平？", "今天穿了新买的衣服，想让大家帮我看看怎么样！\n\n![穿搭](https://picsum.photos/seed/9/800/450)\n\n上衣：新买的卫衣\n裤子：牛仔裤\n鞋子：小白鞋\n\n大家觉得可以打几分？有什么改进建议吗？", 10008, 4, 156, 52, 3210},
			{"分享一下我的省钱小妙招", "作为一个精致的穷鬼，我太有发言权了！\n\n![省钱](https://picsum.photos/seed/10/800/450)\n\n省钱小妙招：\n1. 不买就是省\n2. 买之前先问自己：真的需要吗？\n3. 延迟满足，放购物车冷静几天\n4. 找替代品\n5. 利用优惠券和活动\n\n钱要花在刀刃上！", 10009, 4, 189, 67, 4521},
			{"我的互联网嘴替找到了！", "今天刷到一条评论，说得太对了！完全就是我想说但说不出来的话！\n\n![嘴替](https://picsum.photos/seed/11/800/450)\n\n这条评论简直是在我的脑子里装了监控。\n\n终于有人把我想说的话说出来了！", 10010, 5, 234, 78, 5678},
			{"记录真实的一天，不完美但很真实", "今天过得有点乱糟糟的，但想记录一下真实的生活。\n\n![真实日常](https://picsum.photos/seed/12/800/450)\n\n- 早上起晚了，差点迟到\n- 中午的外卖有点难吃\n- 工作上犯了个小错误\n- 但下班路上买了好吃的\n- 晚上追了喜欢的剧\n\n生活就是这样，不完美但很真实呀！", 10001, 5, 178, 56, 3456},
			{"分享一下今天拍的照片", "今天天气不错，拍了几张照片！\n\n![照片](https://picsum.photos/seed/13/800/450)\n\n第一张：路边的小花\n第二张：天上的云\n第三张：街角的咖啡店\n第四张：夕阳\n\n拍照技术一般，但记录生活很开心！", 10002, 6, 145, 43, 2341},
			{"求推荐！最近有什么好看的剧吗？", "最近剧荒了，大家有什么推荐吗？\n\n![追剧](https://picsum.photos/seed/14/800/450)\n\n我喜欢的类型：\n- 悬疑推理\n- 治愈系\n- 搞笑轻松\n- 古装剧\n\n大家有什么好看的剧推荐吗？电影也行！", 10003, 6, 123, 38, 2156},
			{"社区新手指南，欢迎新朋友！", "欢迎加入我们的社区！给新朋友们介绍一下怎么玩～\n\n![新手指南](https://picsum.photos/seed/15/800/450)\n\n## 社区功能\n1. 发帖分享你的生活\n2. 给喜欢的帖子点赞评论\n3. 关注你喜欢的用户\n4. 给帖子打标签，方便分类\n\n有任何问题都可以在评论区问我哦！", 10000, 7, 89, 23, 1876},
			{"社区公告：请大家友善发言哦", "为了维护良好的社区氛围，请大家注意：\n\n![公告](https://picsum.photos/seed/16/800/450)\n\n## 社区规范\n1. 友善发言，尊重他人\n2. 不发布广告和垃圾信息\n3. 不传播不实信息\n4. 保护他人隐私\n\n希望大家一起营造一个温暖友好的社区！", 10000, 7, 156, 45, 2567},
			{"今天摸鱼的时候写了首小诗", "今天摸鱼的时候灵感来了，写了首小诗～\n\n![写诗](https://picsum.photos/seed/17/800/450)\n\n《摸鱼之歌》\n\n工作堆如山，\n我心已飘然。\n鼠标轻轻点，\n摸鱼乐无边。\n\n哈哈，写得不好，大家见笑了！", 10004, 1, 167, 38, 2654},
			{"分享我的摸鱼高效摆烂指南", "作为一个资深摸鱼选手，我来分享一下如何高效摆烂！\n\n![摆烂](https://picsum.photos/seed/18/800/450)\n\n高效摆烂指南：\n1. 先做完必须做的事\n2. 剩下的时间，想干嘛干嘛\n3. 不要有负罪感\n4. 摆烂也是为了更好地工作\n\n适当摆烂，有益身心健康！", 10005, 1, 145, 32, 2345},
		}
		for i, t := range topics {
			database.DB.Exec(`INSERT INTO topics (id, title, content, user_id, forum_id, like_count, reply_count, view_count, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now', '-'||abs(random())%72||' hours'), datetime('now', '-'||abs(random())%72||' hours'))`,
				10000+i, t.Title, t.Content, t.UserID, t.ForumID, t.LikeCount, t.ReplyCount, t.ViewCount)
		}
		log.Println("[seed] topics created (zh)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('topics', 10019)")
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
			{"Finally decided to stop hustling, life feels elevated", "Finally made the decision to stop the grind! Woke up naturally, had a slow breakfast, sunbathed and read.\n\n![Chill Day](https://picsum.photos/seed/18/800/450)\n\nUsed to think I had to strive and struggle, now I realize that allowing yourself to be imperfect is also self-care.\n\nRecommend a super healing video: https://www.youtube.com/watch?v=dQw4w9WgXcQ\n\nAnother day of loving myself! 💖", 10001, 1, 234, 45, 3241},
			{"Sharing my minimal skincare routine, be kind to yourself", "Used to buy so many products, now simplified and my skin got better!\n\n![Skincare](https://picsum.photos/seed/19/800/450)\n\nMorning routine:\n1. Warm water cleanse\n2. Toner\n3. Serum\n4. Sunscreen\n\nEvening routine:\n1. Cleanse + double cleanse\n2. Toner\n3. Serum\n4. Moisturizer\n\nMost importantly: drink water, sleep early!", 10002, 1, 189, 32, 2891},
			{"Documenting my small win today", "Saw the most beautiful sunset on my way home from work!\n\n![Sunset](https://picsum.photos/seed/20/800/450)\n\nSmall wins in life:\n- Got my favorite drink\n- Just caught the subway\n- Food was delicious today\n- Colleague said I looked nice\n- Saw pretty clouds on the way home\n\nThese little moments make up our lives!", 10001, 1, 178, 29, 2456},
			{"Slacked off all day, it was awesome", "Was super unproductive today, so I just slacked off!\n\n![Slacking](https://picsum.photos/seed/21/800/450)\n\nSlack list:\n- Scrolled Instagram for half an hour\n- Watched some short videos\n- Gossiped with coworkers\n- Had a coffee\n- Looked at fish in the fish tank\n\nSlacking occasionally really improves your mood!", 10003, 1, 201, 34, 2876},
			{"Successfully passed the buck today, feeling great", "Used to be afraid of offending people, took all the blame. Today I finally learned: Not my problem!\n\n![Passing the Buck](https://picsum.photos/seed/22/800/450)\n\nTakeaways:\n1. Clarify your responsibilities\n2. Politely but firmly reject what's not your problem\n3. Offer solutions but don't take responsibility for others\n\nInstead of internalizing, externalize!", 10004, 2, 56, 23, 1234},
			{"Something super funny happened today", "On the subway today, a kid pointed at my hair and said: Auntie, your hair looks like cotton candy!\n\n![Funny](https://picsum.photos/seed/23/800/450)\n\nMe:\n\nThen everyone laughed, and the kid seriously asked if he could have a bite.\n\nKids have such big imaginations, crying laughing! 😂", 10005, 2, 89, 34, 2134},
			{"Watched a super emotional video today, got me in my feelings", "Saw a video today about grandparents' love story, tears were flowing.\n\n![Emotional](https://picsum.photos/seed/24/800/450)\n\nIn the video, grandpa had dementia and forgot many things, but still remembered grandma's favorite candy.\n\nThis kind of long-term companionship really hits different.\n\nCherish the people in front of you!", 10006, 3, 123, 45, 3456},
			{"Said I wouldn't buy anything, then it happened again", "Last month was still saying: No more clothes this year! Then this week...\n\n![Plot Twist](https://picsum.photos/seed/25/800/450)\n\nThe new clothes are just too pretty! Couldn't help but order.\n\nThe打脸 came so fast like a tornado.\n\nIt's okay, money didn't disappear, it just stayed with me in a different form!", 10007, 3, 98, 38, 2891},
			{"Help me rate this outfit, what do you think?", "Wore my new clothes today, want everyone's opinion!\n\n![Outfit](https://picsum.photos/seed/26/800/450)\n\nTop: new sweatshirt\nPants: jeans\nShoes: white sneakers\n\nWhat score would you give? Any improvement suggestions?", 10008, 4, 156, 52, 3210},
			{"Sharing my money-saving tips", "As a sophisticated broke person, I have so much to say!\n\n![Saving Money](https://picsum.photos/seed/27/800/450)\n\nMoney saving tips:\n1. Not buying is saving\n2. Ask yourself before buying: Do I really need this?\n3. Delayed gratification, leave in cart for a few days\n4. Find alternatives\n5. Use coupons and sales\n\nMoney should be spent on what matters!", 10009, 4, 189, 67, 4521},
			{"Found my internet mouthpiece!", "Saw a comment today that was so right! Exactly what I wanted to say but couldn't!\n\n![Mouthpiece](https://picsum.photos/seed/28/800/450)\n\nThis comment was basically monitoring my brain.\n\nFinally someone said what I wanted to say!", 10010, 5, 234, 78, 5678},
			{"Documenting a real day, imperfect but real", "Today was a bit messy, but wanted to document real life.\n\n![Real Life](https://picsum.photos/seed/29/800/450)\n\n- Woke up late, almost late\n- Lunch takeout was kinda bad\n- Made a small mistake at work\n- But got something nice on the way home\n- Binged my favorite show in the evening\n\nLife is like this, imperfect but real!", 10001, 5, 178, 56, 3456},
			{"Sharing photos I took today", "Nice weather today, took some photos!\n\n![Photos](https://picsum.photos/seed/30/800/450)\n\nFirst: flowers by the road\nSecond: clouds in the sky\nThird: coffee shop on the corner\nFourth: sunset\n\nPhotography skills average, but documenting life is fun!", 10002, 6, 145, 43, 2341},
			{"Recommendations! Any good shows lately?", "Running out of shows, any recommendations?\n\n![Binging](https://picsum.photos/seed/31/800/450)\n\nGenres I like:\n- Mystery thriller\n- Healing\n- Funny and light\n- Historical drama\n\nAny good show recommendations? Movies too!", 10003, 6, 123, 38, 2156},
			{"Community Newbie Guide, welcome new friends!", "Welcome to our community! Introducing how to play for new friends～\n\n![Newbie Guide](https://picsum.photos/seed/32/800/450)\n\n## Community Features\n1. Post to share your life\n2. Like and comment on posts you like\n3. Follow users you like\n4. Tag posts for easy categorization\n\nAny questions can be asked in the comments!", 10000, 7, 89, 23, 1876},
			{"Community Announcement: Please be kind", "To maintain a good community atmosphere, please note:\n\n![Announcement](https://picsum.photos/seed/33/800/450)\n\n## Community Guidelines\n1. Be kind, respect others\n2. No ads or spam\n3. Don't spread misinformation\n4. Protect others' privacy\n\nHope everyone helps create a warm and friendly community!", 10000, 7, 156, 45, 2567},
			{"Wrote a little poem while slacking today", "Inspiration hit while slacking today, wrote a little poem～\n\n![Poem](https://picsum.photos/seed/34/800/450)\n\n\"The Slacking Song\"\n\nWork piles up like a mountain,\nMy mind has already wandered.\nMouse clicks softly,\nSlacking is boundless.\n\nHaha, not that good, don't laugh too hard!", 10004, 1, 167, 38, 2654},
			{"Sharing my efficient slacking guide", "As a veteran slacker, let me share how to slack efficiently!\n\n![Slacking Guide](https://picsum.photos/seed/35/800/450)\n\nEfficient slacking guide:\n1. Finish what must be done first\n2. Remaining time, do whatever you want\n3. Don't feel guilty\n4. Slacking is for better work\n\nProper slacking is good for physical and mental health!", 10005, 1, 145, 32, 2345},
		}
		for i, t := range topics {
			database.DB.Exec(`INSERT INTO topics (id, title, content, user_id, forum_id, like_count, reply_count, view_count, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now', '-'||abs(random())%72||' hours'), datetime('now', '-'||abs(random())%72||' hours'))`,
				10000+i, t.Title, t.Content, t.UserID, t.ForumID, t.LikeCount, t.ReplyCount, t.ViewCount)
		}
		log.Println("[seed] topics created (en)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('topics', 10019)")
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('comments', 9999)")
	})
}

// ========== 话题-标签关联数据 ==========

func seedTopicTagsZh() {
	var count int64
	database.DB.Model(&struct{}{}).Table("topic_tags").Count(&count)
	checkAndInsert("topic_tags", &count, func() {
		var tags []models.Tag
		database.DB.Find(&tags)

		var topics []models.Topic
		database.DB.Find(&topics)

		tagMap := make(map[string]uint)
		for _, tag := range tags {
			tagMap[tag.Name] = tag.ID
		}

		topicTagAssociations := []struct {
			TopicIndex int
			TagNames   []string
		}{
			{0, []string{"今日份松弛", "外耗模式"}},
			{1, []string{"爱你老己", "今日小确幸"}},
			{2, []string{"今日份松弛", "今日小确幸"}},
			{3, []string{"邪修一下", "活人感日常"}},
			{4, []string{"外耗模式", "我的互联网嘴替"}},
			{5, []string{"笑死我了", "活人感日常"}},
			{6, []string{"破防了", "今日小确幸"}},
			{7, []string{"真香现场", "活人感日常"}},
			{8, []string{"什么水平？", "求建议/避雷"}},
			{9, []string{"求建议/避雷", "邪修一下"}},
			{10, []string{"我的互联网嘴替", "今日份松弛"}},
			{11, []string{"活人感日常", "今日小确幸"}},
			{12, []string{"今日份松弛", "活人感日常"}},
			{13, []string{"求建议/避雷", "什么水平？"}},
			{14, []string{"今日小确幸", "今日份松弛"}},
			{15, []string{"今日份松弛", "今日小确幸"}},
			{16, []string{"邪修一下", "今日份松弛"}},
			{17, []string{"邪修一下", "今日份松弛"}},
			{18, []string{"今日份松弛", "邪修一下"}},
			{19, []string{"邪修一下", "今日份松弛"}},
		}

		for _, assoc := range topicTagAssociations {
			if assoc.TopicIndex < len(topics) {
				var topicTags []models.Tag
				for _, tagName := range assoc.TagNames {
					if tagID, ok := tagMap[tagName]; ok {
						topicTags = append(topicTags, models.Tag{ID: tagID})
					}
				}
				if len(topicTags) > 0 {
					database.DB.Model(&topics[assoc.TopicIndex]).Association("Tags").Replace(topicTags)
				}
			}
		}
		log.Println("[seed] topic tags created (zh)")
	})
}

func seedTopicTagsEn() {
	var count int64
	database.DB.Model(&struct{}{}).Table("topic_tags").Count(&count)
	checkAndInsert("topic_tags", &count, func() {
		var tags []models.Tag
		database.DB.Find(&tags)

		var topics []models.Topic
		database.DB.Find(&topics)

		tagMap := make(map[string]uint)
		for _, tag := range tags {
			tagMap[tag.Name] = tag.ID
		}

		topicTagAssociations := []struct {
			TopicIndex int
			TagNames   []string
		}{
			{0, []string{"Chill Vibes", "Externalize"}},
			{1, []string{"Love Yourself", "Small Wins"}},
			{2, []string{"Chill Vibes", "Small Wins"}},
			{3, []string{"Hack It", "Real Life"}},
			{4, []string{"Externalize", "Mouthpiece"}},
			{5, []string{"LOL", "Real Life"}},
			{6, []string{"Emotional", "Small Wins"}},
			{7, []string{"Plot Twist", "Real Life"}},
			{8, []string{"Rate This", "Advice/Avoid"}},
			{9, []string{"Advice/Avoid", "Hack It"}},
			{10, []string{"Mouthpiece", "Chill Vibes"}},
			{11, []string{"Real Life", "Small Wins"}},
			{12, []string{"Chill Vibes", "Real Life"}},
			{13, []string{"Advice/Avoid", "Rate This"}},
			{14, []string{"Small Wins", "Chill Vibes"}},
			{15, []string{"Chill Vibes", "Small Wins"}},
			{16, []string{"Hack It", "Chill Vibes"}},
			{17, []string{"Hack It", "Chill Vibes"}},
			{18, []string{"Chill Vibes", "Hack It"}},
			{19, []string{"Hack It", "Chill Vibes"}},
		}

		for _, assoc := range topicTagAssociations {
			if assoc.TopicIndex < len(topics) {
				var topicTags []models.Tag
				for _, tagName := range assoc.TagNames {
					if tagID, ok := tagMap[tagName]; ok {
						topicTags = append(topicTags, models.Tag{ID: tagID})
					}
				}
				if len(topicTags) > 0 {
					database.DB.Model(&topics[assoc.TopicIndex]).Association("Tags").Replace(topicTags)
				}
			}
		}
		log.Println("[seed] topic tags created (en)")
	})
}
