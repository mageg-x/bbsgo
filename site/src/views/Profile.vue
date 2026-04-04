<template>
  <div class="max-w-5xl mx-auto px-3 sm:px-4">
    <div class="bg-white rounded-lg shadow-sm mb-4 sm:mb-6 overflow-hidden">
      <div class="h-32 sm:h-48 bg-cover bg-center relative"
        :style="{ backgroundImage: `url(${user?.background || getUserBackground(user?.username)})` }">
        <div class="flex justify-end p-3 sm:p-4">
          <input ref="backgroundInput" type="file" accept="image/*" class="hidden" @change="handleBackgroundUpload" />
          <button v-if="isCurrentUser" @click="$refs.backgroundInput.click()"
            class="bg-white/90 px-2.5 sm:px-3 py-1 rounded text-xs sm:text-sm text-gray-600 hover:bg-white">
            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z">
              </path>
            </svg>
            {{ t('topic.setBackground') }}
          </button>
        </div>
      </div>
      <div class="px-4 sm:px-6 pb-4 sm:pb-6">
        <div class="flex flex-col sm:flex-row sm:items-end -mt-10 sm:-mt-16 gap-3">
          <div class="relative flex-shrink-0">
            <img :src="getUserAvatar(user)" class="w-20 h-20 sm:w-32 sm:h-32 rounded-full border-4 border-white shadow-lg bg-gray-200">
            <input ref="avatarInput" type="file" accept="image/*" class="hidden" @change="handleAvatarUpload" />
            <button v-if="isCurrentUser" @click="$refs.avatarInput.click()"
              class="absolute bottom-0 right-0 bg-blue-500 text-white rounded-full w-7 h-7 sm:w-8 sm:h-8 flex items-center justify-center hover:bg-blue-600 shadow-lg">
              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812-1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z">
                </path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
            </button>
          </div>
          <div class="flex-1 min-w-0">
            <h1 class="text-lg sm:text-xl font-bold text-gray-900">{{ getUserDisplayName(user, t) }}</h1>
            <p class="text-gray-500 text-xs sm:text-sm mt-1">{{ user?.signature || t('topic.noSignature') }}</p>
          </div>
          <div class="flex gap-2">
            <button v-if="isCurrentUser" @click="openEditDialog"
              class="bg-blue-500 text-white px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium hover:bg-blue-600">
              {{ t('topic.editProfile') }}
            </button>
            <template v-else-if="userStore.isLoggedIn">
              <FollowButton :user-id="parseInt(route.params.id)" />
              <button @click="sendMessage"
                class="bg-green-500 text-white px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium hover:bg-green-600">
                {{ t('topic.sendMessage') }}
              </button>
            </template>
          </div>
        </div>
      </div>
    </div>
    <div class="flex flex-col lg:flex-row gap-3 lg:gap-4">
      <aside class="order-2 lg:order-1 lg:w-64 lg:flex-shrink-0">
        <div class="space-y-3 lg:space-y-4">
          <div class="bg-white rounded-lg shadow-sm p-3 sm:p-4">
            <h3 class="font-medium text-gray-900 mb-3 lg:mb-4 border-b pb-2">{{ t('topic.personalAchievements') }}</h3>
            <div class="grid grid-cols-2 lg:grid-cols-4 gap-2 lg:gap-3 text-center">
              <div class="p-1">
                <div class="text-xl sm:text-2xl font-bold text-gray-700">{{ userStats.topic_count || 0 }}</div>
                <div class="text-xs sm:text-[10px] text-gray-400 leading-tight mt-0.5">{{ t('topic.posts') }}</div>
              </div>
              <div class="p-1">
                <div class="text-xl sm:text-2xl font-bold text-gray-700">{{ userStats.comment_count || 0 }}</div>
                <div class="text-xs sm:text-[10px] text-gray-400 leading-tight mt-0.5">{{ t('topic.comments') }}</div>
              </div>
              <router-link :to="`/user/${user?.id}/follows?type=follows`" class="block hover:bg-gray-50 rounded p-1">
                <div class="text-xl sm:text-2xl font-bold text-gray-700">{{ userStats.follow_count || 0 }}</div>
                <div class="text-xs sm:text-[10px] text-gray-400 leading-tight mt-0.5">{{ t('topic.follows') }}</div>
              </router-link>
              <router-link :to="`/user/${user?.id}/follows?type=followers`" class="block hover:bg-gray-50 rounded p-1">
                <div class="text-xl sm:text-2xl font-bold text-gray-700">{{ userStats.follower_count || 0 }}</div>
                <div class="text-xs sm:text-[10px] text-gray-400 leading-tight mt-0.5">{{ t('topic.fans') }}</div>
              </router-link>
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-3 sm:p-4 hidden lg:block">
            <div class="flex justify-between items-center mb-3 lg:mb-4">
              <h3 class="font-medium text-gray-900">{{ t('topic.badgeWall') }}</h3>
              <router-link v-if="userBadges.length > 0" :to="`/user/${user?.id}/badges`" class="text-blue-500 text-xs sm:text-sm hover:underline">
                {{ t('topic.viewAll') }}
              </router-link>
            </div>
            <div v-if="userBadges.length > 0" class="grid grid-cols-3 lg:grid-cols-4 gap-2 lg:gap-3">
              <div v-for="ub in userBadges.slice(0, 8)" :key="ub.id" 
                class="flex flex-col items-center p-1.5 rounded-lg hover:bg-gray-50 cursor-pointer"
                :title="`${ub.badge?.name} - ${ub.badge?.description}`">
                <SvgBadge :type="ub.badge?.icon" :size="28" />
                <span class="text-[10px] sm:text-xs text-gray-600 text-center truncate w-full mt-0.5">{{ ub.badge?.name }}</span>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-4 lg:py-6 text-xs sm:text-sm">
              {{ t('topic.noBadges') }}
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-3 sm:p-4 hidden lg:block">
            <div class="flex justify-between items-center mb-3 lg:mb-4">
              <h3 class="font-medium text-gray-900">{{ t('profile.title') }}</h3>
              <button v-if="isCurrentUser" @click="openEditDialog" class="text-blue-500 text-xs sm:text-sm hover:underline">
                {{ t('topic.editProfile') }}
              </button>
            </div>
            <div class="space-y-2 lg:space-y-3">
              <div class="flex">
                <span class="w-16 lg:w-20 text-gray-500 text-xs sm:text-sm">{{ t('topic.nickname') }}</span>
                <span class="text-gray-900 text-xs sm:text-sm">{{ getUserDisplayName(user, t) }}</span>
              </div>
              <div class="flex">
                <span class="w-16 lg:w-20 text-gray-500 text-xs sm:text-sm">{{ t('topic.signature') }}</span>
                <span class="text-gray-900 text-xs sm:text-sm">{{ user?.signature || '-' }}</span>
              </div>
              <div class="flex">
                <span class="w-16 lg:w-20 text-gray-500 text-xs sm:text-sm">{{ t('topic.homepage') }}</span>
                <span class="text-blue-500 text-xs sm:text-sm truncate">{{ user?.intro ? user.intro : 'https://mlog.club/user/' + (user?.id || '') }}</span>
              </div>
            </div>
          </div>
          <div class="bg-white rounded-lg shadow-sm p-3 sm:p-4">
            <div class="flex justify-between items-center mb-3 lg:mb-4">
              <h3 class="font-medium text-gray-900">{{ t('topic.fans') }} {{ followers.length }}</h3>
              <button class="text-blue-500 text-xs sm:text-sm hover:underline">{{ t('topic.loadMore') }}</button>
            </div>
            <div v-if="followers.length > 0" class="space-y-3">
              <div v-for="follower in followers" :key="follower.id" class="flex items-center space-x-2">
                <img :src="getUserAvatar(follower.user)" class="w-8 h-8 sm:w-10 sm:h-10 rounded-full bg-gray-200 flex-shrink-0">
                <div class="flex-1 min-w-0">
                  <div class="text-xs sm:text-sm font-medium text-gray-900 truncate">{{ getUserDisplayName(follower.user, t) }}
                  </div>
                  <div class="text-xs text-gray-400 truncate">{{ follower.user?.signature || t('topic.noSignature') }}</div>
                </div>
                <button class="bg-blue-500 text-white text-[10px] sm:text-xs px-2 sm:px-3 py-0.5 sm:py-1 rounded hover:bg-blue-600 flex-shrink-0 whitespace-nowrap">+ {{ t('profile.follow') }}</button>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-4 text-xs sm:text-sm">
              {{ t('topic.noFollowers') }}
            </div>
          </div>
        </div>
      </aside>
      <div class="order-1 lg:order-2 flex-1 min-w-0">
        <div class="bg-slate-100 rounded-lg shadow-sm">
          <div class="p-3 sm:p-4">
            <div v-if="userTopics.length > 0" class="space-y-3">
              <div v-for="topic in userTopics" :key="topic.id"
                class="bg-white rounded-lg border border-gray-100 p-3 sm:p-4 hover:border-blue-200 hover:shadow-md transition-all">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center space-x-2">
                    <span class="text-xs sm:text-sm text-gray-500">{{ getUserDisplayName(topic.user, t) }}</span>
                    <div v-if="getAuthorBadges(topic).length > 0" class="flex items-center gap-0.5">
                      <SvgBadge v-for="badge in getAuthorBadges(topic)" :key="badge.id"
                        :type="badge.icon" :size="14" :title="badge.name" />
                    </div>
                  </div>
                  <div class="flex items-center space-x-2">
                    <button v-if="isCurrentUser" @click="toggleTopicPin(topic)"
                      :class="['text-xs transition-colors', topic.is_user_pinned ? 'text-red-500 hover:text-red-600' : 'text-gray-400 hover:text-red-500']">
                      {{ topic.is_user_pinned ? t('topic.cancelPin') : t('topic.setPin') }}
                    </button>
                    <span class="text-xs text-gray-400">{{ formatTime(topic.created_at) }}</span>
                  </div>
                </div>
                <router-link :to="`/topic/${topic.id}`" class="block">
                  <TopicCard :topic="topic" />
                </router-link>
                <div class="flex items-center space-x-4 sm:space-x-6 text-xs sm:text-sm text-gray-500 mt-2">
                  <button @click="toggleLike(topic)"
                    :class="['flex items-center space-x-1 transition-colors', topic.liked ? 'text-red-500' : 'hover:text-red-500']">
                    <svg class="w-4 h-4" :fill="topic.liked ? 'currentColor' : 'none'" stroke="currentColor"
                      viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                      </path>
                    </svg>
                    <span>{{ topic.like_count || 0 }}</span>
                  </button>
                  <router-link :to="`/topic/${topic.id}`" class="flex items-center space-x-1 hover:text-blue-500">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                      </path>
                    </svg>
                    <span>{{ topic.reply_count || 0 }}</span>
                  </router-link>
                  <button class="flex items-center space-x-1 hover:text-green-500">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z">
                      </path>
                    </svg>
                    <span>{{ topic.view_count || 0 }}</span>
                  </button>
                  <span v-if="topic.forum" class="bg-gray-100 px-2 py-0.5 rounded">{{ topic.forum.name }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-gray-400 py-8 lg:py-12 text-sm">
              {{ t('topic.noTopics') }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showEditDialog" :title="t('topic.editProfile')" width="500px" :close-on-click-modal="false">
      <el-form :model="editForm" label-width="80px">
        <el-form-item :label="t('topic.nickname')">
          <el-input v-model="editForm.nickname" :placeholder="t('topic.nicknamePlaceholder')" maxlength="20" show-word-limit />
        </el-form-item>
        <el-form-item :label="t('topic.signature')">
          <el-input v-model="editForm.signature" type="textarea" :rows="3" :placeholder="t('topic.signaturePlaceholder')" maxlength="100"
            show-word-limit />
        </el-form-item>
        <el-form-item :label="t('topic.homepage')">
          <el-input v-model="editForm.intro" :placeholder="t('topic.homepagePlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSaveProfile" :loading="saving">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getErrorI18nKey } from '@/utils/error'

const { t } = useI18n()
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

function getAuthorBadges(topic) {
  if (!topic || !topic.author_badges || topic.author_badges.length === 0) return []
  return getDisplayBadges(topic.author_badges, 'post-list')
}

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return t('topic.justNow')
  if (diff < 3600000) return t('topic.minutesAgo', { minutes: Math.floor(diff / 60000) })
  if (diff < 86400000) return t('topic.hoursAgo', { hours: Math.floor(diff / 3600000) })
  if (diff < 2592000000) return t('topic.daysAgo', { days: Math.floor(diff / 86400000) })
  return date.toLocaleDateString()
}

async function handleBackgroundUpload(event) {
  const file = event.target.files[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('topic.onlyImageAllowed'))
    return
  }

  try {
    ElMessage.info(t('topic.uploadingBackground'))

    const url = await uploadImage(file, {
      dir: 'backgrounds',
      onInstant: () => ElMessage.success(t('topic.backgroundUpdated') + t('topic.instantUploadSuccess'))
    })

    ElMessage.success(t('topic.backgroundUpdated'))
    user.value.background = url
    await updateProfile({ background: url })
  } catch (error) {
    console.error('Background upload error:', error)
    ElMessage.error(t(getErrorI18nKey(error?.code)))
  }

  event.target.value = ''
}

async function handleAvatarUpload(event) {
  const file = event.target.files[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('topic.onlyImageAllowed'))
    return
  }

  try {
    ElMessage.info(t('topic.uploadingAvatar'))

    const url = await uploadImage(file, {
      dir: 'avatars',
      onInstant: () => ElMessage.success(t('topic.avatarUpdated') + t('topic.instantUploadSuccess'))
    })

    ElMessage.success(t('topic.avatarUpdated'))
    user.value.avatar = url
    await updateProfile({ avatar: url })
  } catch (error) {
    console.error('Avatar upload error:', error)
    ElMessage.error(t(getErrorI18nKey(error?.code)))
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
    ElMessage.success(t('topic.profileUpdated'))
  } catch (error) {
    ElMessage.error(t(getErrorI18nKey(error?.code)))
  } finally {
    saving.value = false
  }
}

async function toggleTopicPin(topic) {
  try {
    await ElMessageBox.confirm(
      topic.is_user_pinned ? t('topic.confirmUnpin') : t('topic.confirmPin'),
      topic.is_user_pinned ? t('topic.unpinTopic') : t('topic.pinTopic'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    await topicApi.pinTopic(topic.id, !topic.is_user_pinned)

    const originalIndex = userTopics.value.findIndex(t => t.id === topic.id)
    if (originalIndex === -1) return

    topic.is_user_pinned = !topic.is_user_pinned

    userTopics.value.splice(originalIndex, 1)

    if (topic.is_user_pinned) {
      userTopics.value.unshift(topic)
    } else {
      let insertIndex = userTopics.value.length
      for (let i = 0; i < userTopics.value.length; i++) {
        if (!userTopics.value[i].is_user_pinned) {
          insertIndex = i
          break
        }
      }
      userTopics.value.splice(insertIndex, 0, topic)
    }

    ElMessage.success(topic.is_user_pinned ? t('topic.topicPinned') : t('topic.topicUnpinned'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Operation failed', e)
      ElMessage.error(t(getErrorI18nKey(e?.code)))
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
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function loadUserTopics() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/topics`)
    userTopics.value = res?.list || []
  } catch (e) {
    console.error('加载用户帖子失败', e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function loadFollowers() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/followers`)
    followers.value = res?.list || []
  } catch (e) {
    console.error('加载粉丝列表失败', e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
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
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function loadUserBadges() {
  try {
    const userId = route.params.id
    const res = await api.get(`/users/${userId}/badges`)
    userBadges.value = res || []
  } catch (e) {
    console.error('加载用户勋章失败', e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function toggleLike(topic) {
  if (!userStore.isLoggedIn) {
    ElMessage.warning(t('common.pleaseLoginFirst'))
    return
  }

  try {
    if (topic.liked) {
      await api.delete(`/likes?target_type=topic&target_id=${topic.id}`)
      topic.like_count--
    } else {
      await api.post('/likes', {
        target_type: 'topic',
        target_id: topic.id
      })
      topic.like_count++
    }
    topic.liked = !topic.liked
  } catch (e) {
    console.error(e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function checkTopicLikes() {
  if (!userStore.isLoggedIn || userTopics.value.length === 0) return

  try {
    const topicIds = userTopics.value.map(t => t.id)
    const res = await api.post('/likes/check', {
      target_type: 'topic',
      target_ids: topicIds
    })
    if (res.liked_map) {
      for (const topic of userTopics.value) {
        topic.liked = res.liked_map[topic.id] || false
      }
    }
  } catch (e) {
    console.error(e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
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

watch(() => userTopics.value.length, (newLen) => {
  if (newLen > 0) {
    checkTopicLikes()
  }
})

watch(() => userStore.isLoggedIn, (isLoggedIn) => {
  if (isLoggedIn) {
    checkTopicLikes()
  }
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
