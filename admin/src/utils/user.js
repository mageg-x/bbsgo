export function getUserAvatar(userData) {
  if (userData?.avatar) {
    return userData.avatar
  }
  const username = userData?.username || 'default'
  return `https://api.dicebear.com/9.x/adventurer/svg?seed=${encodeURIComponent(username)}`
}

export function getUserDisplayName(userData) {
  return userData?.nickname || userData?.username || '未知用户'
}
