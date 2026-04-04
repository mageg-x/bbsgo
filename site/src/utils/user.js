export function getUserAvatar(userData) {
  if (userData?.avatar) {
    return userData.avatar
  }
  const username = userData?.username || 'default'
  return `https://api.dicebear.com/9.x/adventurer/svg?seed=${encodeURIComponent(username)}`
}

export function getUserDisplayName(userData, t) {
  const name = userData?.nickname || userData?.username
  if (name) return name
  return t ? t('user.unknown') : '未知用户'
}
