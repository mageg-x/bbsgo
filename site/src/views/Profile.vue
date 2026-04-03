<template>
  <div>
    <div class= "bg-white rounded-lg shadow-sm mb-6 overflow-hidden">
      <div class="h-48 bg-cover bg-center relative"
        :style="{ backgroundImage: `url(${user?.background || getUserBackground(user?.username)})` }">
        <div class="flex justify-end p-4">
          <input ref="backgroundInput" type="file" accept="image/*" class="hidden" @change="handleBackgroundUpload" />
          <button @click="$refs.backgroundInput.click()"
            class="bg-white/90 px-3 py-1 rounded text-sm text-gray-600 hover:bg-white">
            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z">
              </path>
            </svg>
            设置背景
          </button>
        </div>
      </div>
      <div class="px-6 pb-6">
        <div class="flex items-end -mt-16">
          <div class="relative">
            <img :src="getUserAvatar(user)" class="w-32 h-32 rounded-full border-4 border-white shadow-lg bg-gray-200">
            <input ref="avatarInput" type="file" accept="image/*" class="hidden" @change="handleAvatarUpload" />
            <button v-if="isCurrentUser" @click="$refs.avatarInput.click()"
              class="absolute bottom-0 right-0 bg-blue-500 text-white rounded-full w-8 h-8 flex items-center justify-center hover:bg-blue-600 shadow-lg">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z">
                </path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
            </button>
          </div>
          <div class="ml-4 mb-2 flex-1">
            <h1 class="text-xl font-bold text-gray-900">{{ getUserDisplayName(user) }}</h1>
            <p class="text-gray-500 text-sm">{{ user?.signature || '这家伙很懒，什么都没留下...' }}</p>
          </div>
          <div class="mb-2 flex gap-2">
            <button v-if="isCurrentUser" @click="openEditDialog"
              class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-blue-600">
              编辑资料
            </button>
            <template v-else-if="userStore.isLoggedIn">
              <FollowButton :user-id="parseInt(route.params.id)" />
              <button @click="sendMessage"
                class="bg-green-500 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-green-600">
                发私信
              </button>
            </template>
          </div>
        </div>
      </div>
    </div>
    <div class="flex gap-6">
      <aside class="w-72 flex-shrink-0">
        <div class="space-y-4">
          <div class="bg-white rounded-lg shadow-sm p-4">
            <h3 class="font-medium text-gray-900 mb-4 border-b pb-2">个人成就</h3>
            <div class="grid grid-cols-4 gap-4 text-center">
              <div>
                <div class="text-2xl font-bold text-gray-700">{{ userStats.topic_count || 0 }}</div>
                <div class="text-xs text-gray-400">帖子</div>
              </div>
              <div>
                <div class="text-2xl font-bold text-gray-700">{{ userStats.comment_count || 0 }}</div>
                <div class="text-xs text-gray-400">评论</div>
              </div>
              <router-link :to="`/user/${user?.id}/follows?type=follows`" class="block hover:bg-gray-50 rounded">
                <div class="text-2xl font-bold text-gray-700">{{ userStats.follow_count || 0 }}</div>
                <div class="text-xs text-gray-400">关注</div>
              </router-link>
              <router-link :to="`/user/${user?.id}/follows?type=followers`" class="block hover:bg-gray-50 rounded">
                <div class="text-2xl font-bold text-gray-700">{{ userStats.follower_count || 0 }}</div>
                <div class="text-xs text-gray-400">粉丝</div>
              </router-link>
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="font-medium text-gray-900">勋章墙</h3>
              <router-link v-if="userBadges.length > 0" :to="`/user/${user?.id}/badges`" class="text-blue-500 text-sm hover:underline">
                查看全部
              </router-link>
            </div>
            <div v-if="userBadges.length > 0" class="grid grid-cols-6 gap-3">
              <div v-for="ub in userBadges.slice(0, 6)" :key="ub.id" 
                class="flex flex-col items-center p-2 rounded-lg hover:bg-gray-50 cursor-pointer"
                :title="`${ub.badge?.name} - ${ub.badge?.description}`">
                <SvgBadge :type="ub.badge?.icon" :size="32" />
                <span class="text-xs text-gray-600 text-center truncate w-full mt-1">{{ ub.badge?.name }}</span>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-6 text-sm">
              暂无勋章
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="font-medium text-gray-900">个人资料</h3>
              <button v-if="isCurrentUser" @click="openEditDialog" class="text-blue-500 text-sm hover:underline">
                编辑资料
              </button>
            </div>
            <div class="space-y-3">
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">昵称</span>
                <span class="text-gray-900 text-sm">{{ getUserDisplayName(user) }}</span>
              </div>
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">签名</span>
                <span class="text-gray-900 text-sm">{{ user?.signature || '-' }}</span>
              </div>
              <div class="flex">
                <span class="w-20 text-gray-500 text-sm">主页</span>
                <span class="text-blue-500 text-sm">{{ user?.intro ? user.intro : 'https://mlog.club/user/' + (user?.id
                  || '') }}</span>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="font-medium text-gray-900">粉丝 {{ followers.length }}</h3>
              <button class="text-blue-500 text-sm hover:underline">更多</button>
            </div>
            <div v-if="followers.length > 0" class="space-y-3">
              <div v-for="follower in followers" :key="follower.id" class="flex items-center space-x-3">
                <img :src="getUserAvatar(follower.user)" class="w-10 h-10 rounded-full bg-gray-200">
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-medium text-gray-900 truncate">{{ getUserDisplayName(follower.user) }}
                  </div>
                  <div class="text-xs text-gray-400 truncate">{{ follower.user?.signature || '这家伙很懒，什么都没留下' }}</div>
                </div>
                <button class="bg-blue-500 text-white text-xs px-3 py-1 rounded hover:bg-blue-600">+ 关注</button>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-4 text-sm">
              暂无粉丝
            </div>
          </div>
        </div>
      </aside>
      <div class="flex-1 min-w-0">
        <div class=" bg-slate-100 rounded-lg shadow-sm">
          <div class="p-4">
            <div v-if="userTopics.length > 0" class="space-y-3">
              <div v-for="topic in userTopics" :key="topic.id"
                class="bg-white rounded-lg border border-gray-100 p-4 hover:border-blue-200 hover:shadow-md transition-all">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center space-x-2">
                    <span class="text-sm text-gray-500">{{ getUserDisplayName(topic.user) }}</span>
                    <div v-if="getAuthorBadges(topic).length > 0" class="flex items-center gap-0.5">
                      <SvgBadge v-for="badge in getAuthorBadges(topic)" :key="badge.id"
                        :type="badge.icon" :size="16" :title="badge.name" />
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <button v-if="isCurrentUser" @click="toggleTopicPin(topic)"
                      :class="['text-xs transition-colors', topic.is_user_pinned ? 'text-red-500 hover:text-red-600' : 'text-gray-400 hover:text-red-500']">
                      {{ topic.is_user_pinned ? '取消置顶' : '置顶' }}
                    </button>
                    <span class="text-xs text-gray-400">{{ formatTime(topic.created_at) }}</span>
                  </div>
                </div>
                <router-link :to="`/topic/${topic.id}`" class="block">
                  <TopicCard :topic="topic" />
                </router-link>
                <div class="flex items-center space-x-4 text-xs text-gray-500">
                  <span>👍 {{ topic.like_count || 0 }}</span>
                  <span>💬 {{ topic.comment_count || 0 }}</span>
                  <span>👁 {{ topic.view_count || 0 }}</span>
                  <span v-if="topic.forum" class="bg-gray-100 px-2 py-0.5 rounded">{{ topic.forum.name }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-12">
              暂无帖子
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showEditDialog" title="编辑资料" width="500px" :close-on-click-modal="false">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" maxlength="20" show-word-limit />
        </el-form-item>
        <el-form-item label="签名">
          <el-input v-model="editForm.signature" type="textarea" :rows="3" placeholder="请输入个性签名" maxlength="100"
            show-word-limit />
        </el-form-item>
        <el-form-item label="个人主页">
          <el-input v-model="editForm.intro" placeholder="请输入个人主页链接" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSaveProfile" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import api, { topicApi } from '@/api'
import { uploadImage } from '@/utils/upload'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'
import { getDisplayBadges } from '@/utils/badge'
import TopicCard from '@/components/TopicCard.vue'
import SvgBadge from '@/components/SvgBadge.vue'
import FollowButton from '@/components/FollowButton.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const user = ref(null)
const userStats = ref({
  topic_count: 0,
  comment_count: 0,
  follow_count: 0,
  follower_count: 0
})
const userTopics = ref([])
const followers = ref([])
const userBadges = ref([])
const showEditDialog = ref(false)
const saving = ref(false)
const editForm = ref({
  nickname: '',
  signature: '',
  intro: ''
})

const isCurrentUser = computed(() => {
  return userStore.user?.id === parseInt(route.params.id)
})

// 监听路由参数变化，重新加载数据
watch(() => route.params.id, async (newId) => {
  if (newId) {
    await Promise.all([
      loadUser(),
      loadUserTopics(),
      loadFollowers(),
      loadUserStats(),
      loadUserBadges()
    ])
  }
}, { immediate: false })

function sendMessage() {
  // 跳转到私信页面，并带上目标用户信息
  router.push({
    path: '/messages',
    query: { userId: user.value.id, username: user.value.username }
  })
}

function getUserBackground(username) {
  if (!username) {
    return 'https://picsum.photos/id/1015/1200/400'
  }
  let hash = 0
  for (let i = 0; i < username.length; i++) {
    hash = ((hash << 5) - hash) + username.charCodeAt(i)
    hash = hash & hash
  }
  const imageId = (Math.abs(hash) % 1000) + 1
  return `https://picsum.photos/id/${imageId}/1200/400`
}

// 获取帖子作者的展示勋章（最多2枚）
function getAuthorBadges(topic) {
  if (!topic || !topic.author_badges || topic.author_badges.length === 0) return []
  return getDisplayBadges(topic.author_badges, 'post-list')
}

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  if (diff < 2592000000) return Math.floor(diff / 86400000) + '天前'
  return date.toLocaleDateString()
}

async function handleBackgroundUpload(event) {
  const file = event.target.files[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    ElMessage.error('只能上传图片文件')
    return
  }

  try {
    ElMessage.info('正在上传背景图片...')

    const url = await uploadImage(file, {
      dir: 'backgrounds',
      onInstant: () => ElMessage.success('背景图片更新成功（秒传）')
    })

    ElMessage.success('背景图片更新成功')
    user.value.background = url
    await updateProfile({ background: url })
  } catch (error) {
    console.error('Background upload error:', error)
    ElMessage.error('背景图片上传失败')
  }

  event.target.value = ''
}

async function handleAvatarUpload(event) {
  const file = event.target.files[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    ElMessage.error('只能上传图片文件')
    return
  }

  try {
    ElMessage.info('正在上传头像...')

    const url = await uploadImage(file, {
      dir: 'avatars',
      onInstant: () => ElMessage.success('头像更新成功（秒传）')
    })

    ElMessage.success('头像更新成功')
    user.value.avatar = url
    await updateProfile({ avatar: url })
  } catch (error) {
    console.error('Avatar upload error:', error)
    ElMessage.error('头像上传失败')
  }

  event.target.value = ''
}

async function updateProfile(data) {
  try {
    await api.put('/user/profile', data)
    if (userStore.user?.id === user.value?.id) {
      userStore.user = { ...userStore.user, ...data }
      localStorage.setItem('user', JSON.stringify(userStore.user))
    }
  } catch (error) {
    console.error('Update profile error:', error)
    throw error
  }
}

function openEditDialog() {
  editForm.value = {
    nickname: user.value?.nickname || '',
    signature: user.value?.signature || '',
    intro: user.value?.intro || ''
  }
  showEditDialog.value = true
}

async function handleSaveProfile() {
  saving.value = true

  try {
    await updateProfile(editForm.value)
    user.value = { ...user.value, ...editForm.value }

    if (isCurrentUser.value) {
      userStore.user = { ...userStore.user, ...editForm.value }
      localStorage.setItem('user', JSON.stringify(userStore.user))
    }

    showEditDialog.value = false
    ElMessage.success('资料更新成功')
  } catch (error) {
    ElMessage.error('资料更新失败')
  } finally {
    saving.value = false
  }
}

async function toggleTopicPin(topic) {
  try {
    await ElMessageBox.confirm(
      topic.is_user_pinned ? '确定要取消置顶这篇帖子吗？' : '确定要置顶这篇帖子吗？',
      topic.is_user_pinned ? '取消置顶' : '置顶帖子',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await topicApi.pinTopic(topic.id, !topic.is_user_pinned)

    // 记录原始索引
    const originalIndex = userTopics.value.findIndex(t => t.id === topic.id)
    if (originalIndex === -1) return

    // 更新置顶状态
    topic.is_user_pinned = !topic.is_user_pinned

    // 从原位置移除
    userTopics.value.splice(originalIndex, 1)

    if (topic.is_user_pinned) {
      // 置顶：插入到最前面
      userTopics.value.unshift(topic)
    } else {
      // 取消置顶：插入到非置顶帖子之后
      let insertIndex = userTopics.value.length
      for (let i = 0; i < userTopics.value.length; i++) {
        if (!userTopics.value[i].is_user_pinned) {
          insertIndex = i
          break
        }
      }
      userTopics.value.splice(insertIndex, 0, topic)
    }

    ElMessage.success(topic.is_user_pinned ? '帖子已置顶' : '帖子已取消置顶')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('操作失败', e)
      ElMessage.error('操作失败')
    }
  }
}

async function loadUser() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}`)
    user.value = res
  } catch (e) {
    console.error('加载用户信息失败', e)
  }
}

async function loadUserTopics() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/topics`)
    userTopics.value = res?.list || []
  } catch (e) {
    console.error('加载用户帖子失败', e)
  }
}

async function loadFollowers() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/followers`)
    followers.value = res?.list || []
  } catch (e) {
    console.error('加载粉丝列表失败', e)
  }
}

async function loadUserStats() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/stats`)
    userStats.value = res || {
      topic_count: 0,
      comment_count: 0,
      rank: 0
    }
  } catch (e) {
    console.error('加载用户统计失败', e)
  }
}

async function loadUserBadges() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/badges`)
    userBadges.value = res || []
  } catch (e) {
    console.error('加载用户勋章失败', e)
  }
}

onMounted(async () => {
  await Promise.all([
    loadUser(),
    loadUserTopics(),
    loadFollowers(),
    loadUserStats(),
    loadUserBadges()
  ])
})
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
