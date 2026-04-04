// 错误码翻译
// 后端错误码与前端翻译 key 的映射
// 错误码规范: https://github.com/xxx/bbsgo/blob/main/server/errors/codes.go

export default {
  // 1xxx - 认证注册类
  1001: 'errors.registerDisabled',
  1002: 'errors.invalidParams',
  1003: 'errors.incompleteInfo',
  1004: 'errors.usernameExists',
  1005: 'errors.emailExists',
  1006: 'errors.passwordHashFailed',
  1007: 'errors.tokenGenerateFailed',
  1008: 'errors.usernameOrPassword',
  1009: 'errors.invalidEmail',
  1010: 'errors.verifyCodeError',
  1011: 'errors.verifyCodeExpired',
  1012: 'errors.emailNotRegistered',

  // 2xxx - 用户相关类
  2001: 'errors.userNotFound',
  2002: 'errors.noPermission',
  2003: 'errors.cannotFollowSelf',
  2004: 'errors.alreadyFollowed',
  2005: 'errors.notFollowed',
  2006: 'errors.creditsInsufficient',
  2007: 'errors.reputationLow',

  // 3xxx - 内容操作类
  3001: 'errors.topicNotFound',
  3002: 'errors.commentNotFound',
  3003: 'errors.forumNotFound',
  3004: 'errors.tagNotFound',
  3005: 'errors.favoriteNotFound',
  3006: 'errors.messageNotFound',
  3007: 'errors.notificationNotFound',
  3008: 'errors.pollNotFound',
  3009: 'errors.draftNotFound',

  // 4xxx - 业务限制类
  4001: 'errors.postDisabled',
  4002: 'errors.commentDisabled',
  4003: 'errors.postIntervalFast',
  4004: 'errors.sensitiveContent',
  4005: 'errors.titleTooLong',
  4006: 'errors.contentTooLong',
  4007: 'errors.pollOptionsFew',
  4008: 'errors.pollEnded',
  4009: 'errors.alreadyVoted',
  4010: 'errors.voteExceedMax',
  4011: 'errors.fileSizeExceeded',
  4012: 'errors.fileTypeUnsupported',
  4013: 'errors.imageSizeExceeded',

  // 5xxx - 系统错误类
  5001: 'errors.serverInternal',
  5002: 'errors.databaseError',
  5003: 'errors.uploadFailed',
  5004: 'errors.cacheError',
  5005: 'errors.thirdPartyError',
}
