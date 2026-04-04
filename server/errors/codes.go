/*
错误码规范文档
==================

本文档定义了系统所有错误码的规范，必须严格遵守。

【错误码结构】
- 错误码为整数，范围 1000-9999
- 格式：ABCD
  - A: 大类 (1=认证, 2=用户, 3=内容, 4=业务, 5=系统)
  - BCD: 具体错误序号

【错误码分类】

1xxx - 认证注册类
    1001: 注册功能已关闭
    1002: 无效的请求参数
    1003: 请填写完整信息
    1004: 用户名已存在
    1005: 邮箱已被注册
    1006: 密码加密失败
    1007: 生成令牌失败
    1008: 用户名或密码错误
    1009: 邮箱格式错误
    1010: 验证码错误
    1011: 验证码已过期
    1012: 邮箱未注册

2xxx - 用户相关类
    2001: 用户不存在
    2002: 无权限操作
    2003: 无法关注自己
    2004: 已经关注过
    2005: 尚未关注该用户
    2006: 积分不足
    2007: 信誉分过低

3xxx - 内容操作类
    3001: 话题不存在
    3002: 评论不存在
    3003: 版块不存在
    3004: 标签不存在
    3005: 收藏不存在
    3006: 消息不存在
    3007: 通知不存在
    3008: 投票不存在
    3009: 草稿不存在

4xxx - 业务限制类
    4001: 发帖功能已关闭
    4002: 评论功能已关闭
    4003: 发帖间隔过快
    4004: 内容包含敏感词
    4005: 标题长度超出限制
    4006: 内容长度超出限制
    4007: 投票选项数量不足
    4008: 投票已结束
    4009: 已投过票
    4010: 超出投票选项限制
    4011: 文件大小超出限制
    4012: 文件类型不支持
    4013: 图片尺寸超出限制

5xxx - 系统错误类
    5001: 服务器内部错误
    5002: 数据库操作失败
    5003: 文件上传失败
    5004: 缓存操作失败
    5005: 第三方服务调用失败

【使用规则】

1. 所有 API 错误必须使用错误码，不得硬编码中文消息
2. 错误消息字段 message 统一传空字符串，由前端根据错误码查找翻译
3. 日志中应同时记录错误码和中文描述，方便后端排查
4. 新增错误码必须在本文档和 errors.go 常量中同时添加
5. 错误码一旦确定不应变更，废弃的错误码标记为 deprecated

【前端处理】

前端通过 errors.js 维护错误码到翻译 key 的映射：
{
  1001: 'register.closed',
  1004: 'register.usernameExists',
  ...
}

*/

package errors

// 错误码常量
const (
	// 1xxx - 认证注册类
	CodeRegisterDisabled    = 1001 // 注册功能已关闭
	CodeInvalidParams       = 1002 // 无效的请求参数
	CodeIncompleteInfo      = 1003 // 请填写完整信息
	CodeUsernameExists      = 1004 // 用户名已存在
	CodeEmailExists         = 1005 // 邮箱已被注册
	CodePasswordHashFailed  = 1006 // 密码加密失败
	CodeTokenGenerateFailed = 1007 // 生成令牌失败
	CodeUsernameOrPassword  = 1008 // 用户名或密码错误
	CodeInvalidEmail        = 1009 // 邮箱格式错误
	CodeVerifyCodeError     = 1010 // 验证码错误
	CodeVerifyCodeExpired   = 1011 // 验证码已过期
	CodeEmailNotRegistered  = 1012 // 邮箱未注册

	// 2xxx - 用户相关类
	CodeUserNotFound      = 2001 // 用户不存在
	CodeNoPermission     = 2002 // 无权限操作
	CodeCannotFollowSelf  = 2003 // 无法关注自己
	CodeAlreadyFollowed  = 2004 // 已经关注过
	CodeNotFollowed      = 2005 // 尚未关注该用户
	CodeCreditsInsufficient = 2006 // 积分不足
	CodeReputationLow     = 2007 // 信誉分过低
	CodeUnauthorized      = 2008 // 未认证，请先登录
	CodeUserBanned        = 2009 // 用户已被禁言

	// 3xxx - 内容操作类
	CodeTopicNotFound   = 3001 // 话题不存在
	CodeCommentNotFound = 3002 // 评论不存在
	CodeForumNotFound   = 3003 // 版块不存在
	CodeTagNotFound     = 3004 // 标签不存在
	CodeFavoriteNotFound = 3005 // 收藏不存在
	CodeMessageNotFound  = 3006 // 消息不存在
	CodeNotificationNotFound = 3007 // 通知不存在
	CodePollNotFound     = 3008 // 投票不存在
	CodeDraftNotFound    = 3009 // 草稿不存在

	// 4xxx - 业务限制类
	CodePostDisabled       = 4001 // 发帖功能已关闭
	CodeCommentDisabled    = 4002 // 评论功能已关闭
	CodePostIntervalFast   = 4003 // 发帖间隔过快
	CodeSensitiveContent   = 4004 // 内容包含敏感词
	CodeTitleTooLong       = 4005 // 标题长度超出限制
	CodeContentTooLong     = 4006 // 内容长度超出限制
	CodePollOptionsFew     = 4007 // 投票选项数量不足
	CodePollEnded          = 4008 // 投票已结束
	CodeAlreadyVoted       = 4009 // 已投过票
	CodeVoteExceedMax      = 4010 // 超出投票选项限制
	CodeFileSizeExceeded   = 4011 // 文件大小超出限制
	CodeFileTypeUnsupported = 4012 // 文件类型不支持
	CodeImageSizeExceeded  = 4013 // 图片尺寸超出限制
	CodeContentTooShort   = 4014 // 内容太短
	CodeOperationTooFast   = 4015 // 操作过快，请稍后再试
	CodeDailyLimitExceeded = 4016 // 今日操作次数已达上限
	CodeNoSubstantiveContent = 4017 // 内容无实质信息
	CodeSymbolsOrEmojiOnly  = 4018 // 内容仅包含符号或表情
	CodeRepeatingChars      = 4019 // 内容包含重复字符
	CodeTooManyLinks       = 4020 // 内容包含过多外部链接

	// 5xxx - 系统错误类
	CodeServerInternal = 5001 // 服务器内部错误
	CodeDatabaseError  = 5002 // 数据库操作失败
	CodeUploadFailed   = 5003 // 文件上传失败
	CodeCacheError     = 5004 // 缓存操作失败
	CodeThirdPartyError = 5005 // 第三方服务调用失败
)

// 获取错误码对应的中文消息（仅用于日志记录）
var CodeMessages = map[int]string{
	CodeRegisterDisabled:    "注册功能已关闭",
	CodeInvalidParams:       "无效的请求参数",
	CodeIncompleteInfo:      "请填写完整信息",
	CodeUsernameExists:       "用户名已存在",
	CodeEmailExists:          "邮箱已被注册",
	CodePasswordHashFailed:   "密码加密失败",
	CodeTokenGenerateFailed:  "生成令牌失败",
	CodeUsernameOrPassword:   "用户名或密码错误",
	CodeInvalidEmail:         "邮箱格式错误",
	CodeVerifyCodeError:      "验证码错误",
	CodeVerifyCodeExpired:    "验证码已过期",
	CodeEmailNotRegistered:   "邮箱未注册",
	CodeUserNotFound:         "用户不存在",
	CodeNoPermission:        "无权限操作",
	CodeUnauthorized:        "未认证，请先登录",
	CodeUserBanned:           "用户已被禁言",
	CodeCannotFollowSelf:     "无法关注自己",
	CodeAlreadyFollowed:      "已经关注过",
	CodeNotFollowed:          "尚未关注该用户",
	CodeCreditsInsufficient:  "积分不足",
	CodeReputationLow:        "信誉分过低",
	CodeTopicNotFound:        "话题不存在",
	CodeCommentNotFound:      "评论不存在",
	CodeForumNotFound:        "版块不存在",
	CodeTagNotFound:          "标签不存在",
	CodeFavoriteNotFound:     "收藏不存在",
	CodeMessageNotFound:      "消息不存在",
	CodeNotificationNotFound: "通知不存在",
	CodePollNotFound:         "投票不存在",
	CodeDraftNotFound:        "草稿不存在",
	CodePostDisabled:         "发帖功能已关闭",
	CodeCommentDisabled:      "评论功能已关闭",
	CodePostIntervalFast:     "发帖间隔过快",
	CodeSensitiveContent:     "内容包含敏感词",
	CodeTitleTooLong:         "标题长度超出限制",
	CodeContentTooLong:       "内容长度超出限制",
	CodePollOptionsFew:       "投票选项数量不足",
	CodePollEnded:            "投票已结束",
	CodeAlreadyVoted:         "已投过票",
	CodeVoteExceedMax:        "超出投票选项限制",
	CodeFileSizeExceeded:     "文件大小超出限制",
	CodeFileTypeUnsupported:  "文件类型不支持",
	CodeImageSizeExceeded:    "图片尺寸超出限制",
	CodeContentTooShort:     "内容太短",
	CodeOperationTooFast:   "操作过快，请稍后再试",
	CodeDailyLimitExceeded:  "今日操作次数已达上限",
	CodeNoSubstantiveContent: "内容无实质信息",
	CodeSymbolsOrEmojiOnly:  "内容仅包含符号或表情",
	CodeRepeatingChars:      "内容包含重复字符",
	CodeTooManyLinks:       "内容包含过多外部链接",
	CodeServerInternal:       "服务器内部错误",
	CodeDatabaseError:        "数据库操作失败",
	CodeUploadFailed:         "文件上传失败",
	CodeCacheError:           "缓存操作失败",
	CodeThirdPartyError:      "第三方服务调用失败",
}

// GetMessage 根据错误码获取中文消息（仅用于日志）
func GetMessage(code int) string {
	if msg, ok := CodeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
