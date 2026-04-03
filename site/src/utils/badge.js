// 勋章展示优先级定义（数值越小优先级越高）
export const BADGE_PRIORITY = {
  'legend': 1,           // 社区传奇
  'opinion-leader': 2,   // 意见领袖
  'gold-comment': 3,     // 金牌评论
  'popular': 4,          // 广受欢迎
  'writer': 5,           // 笔耕不辍
  'community-star': 5,   // 社区活宝（与笔耕不辍同级）
  'first-post': 6,       // 首次发声
  'first-comment': 7,     // 热心回复
  'newcomer': 8,         // 初来乍到
}

// 基础勋章列表（用于判断是否只拥有基础勋章）
const BASIC_BADGES = ['first-post', 'first-comment', 'newcomer']

/**
 * 获取用户已获得的勋章中优先级最高的N个
 * @param {Array} userBadges - 用户勋章数组，每个元素包含 badge 对象
 * @param {number} limit - 最大展示数量
 * @returns {Array} 排序后的前N个勋章
 */
export function getTopBadges(userBadges, limit = 2) {
  if (!userBadges || userBadges.length === 0) {
    return []
  }

  // 提取所有勋章并按优先级排序
  const sortedBadges = [...userBadges]
    .filter(ub => ub.badge && ub.badge.icon)
    .map(ub => ({
      ...ub.badge,
      awarded_at: ub.awarded_at,
      user_badge_id: ub.id
    }))
    .sort((a, b) => {
      const priorityA = BADGE_PRIORITY[a.icon] || 999
      const priorityB = BADGE_PRIORITY[b.icon] || 999
      return priorityA - priorityB
    })

  if (sortedBadges.length === 0) {
    return []
  }

  // 检查是否只拥有基础勋章
  const hasNonBasicBadge = sortedBadges.some(badge => !BASIC_BADGES.includes(badge.icon))

  // 如果只拥有基础勋章，只展示1个
  if (!hasNonBasicBadge) {
    return sortedBadges.slice(0, 1)
  }

  return sortedBadges.slice(0, limit)
}

/**
 * 根据位置获取勋章展示数量
 * @param {string} location - 展示位置: 'post-list', 'author-card', 'comment'
 * @returns {number} 最大展示数量
 */
export function getBadgeLimit(location) {
  switch (location) {
    case 'post-list':
      return 2
    case 'author-card':
      return 3
    case 'comment':
      return 1
    default:
      return 2
  }
}

/**
 * 获取用户展示用的勋章
 * @param {Array} userBadges - 用户勋章数组
 * @param {string} location - 展示位置
 * @returns {Array} 要展示的勋章数组
 */
export function getDisplayBadges(userBadges, location = 'post-list') {
  const limit = getBadgeLimit(location)
  return getTopBadges(userBadges, limit)
}
