<template>
  <div class="flex gap-6 max-w-6xl mx-auto">
    <aside class="w-56 flex-shrink-0 hidden lg:block">
      <div class="bg-white rounded-lg shadow-sm overflow-hidden">
        <div class="px-4 py-3 border-b bg-gray-50">
          <h3 class="font-semibold text-gray-700">热门话题</h3>
        </div>
        <router-link v-for="tag in tags" :key="tag.id" :to="tag.id ? `/?tag=${tag.id}` : '/'" :class="['px-4 py-3 flex items-center justify-between transition-colors',
          currentTagId === tag.id ? 'bg-blue-50 text-blue-600' : 'text-gray-600 hover:bg-gray-50']">
          <div class="flex items-center space-x-2">
            <span v-if="tag.icon" class="text-lg">{{ tag.icon }}</span>
            <span class="font-medium">{{ tag.name }}</span>
          </div>
          <span class="text-xs text-gray-400">{{ tag.usage_count }}</span>
        </router-link>
      </div>
    </aside>
    <div class="flex-1 min-w-0">
      <div v-if="announcements.length > 0" class="mb-4">
        <div v-for="announcement in announcements" :key="announcement.id"
          class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-2">
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <svg class="w-5 h-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z">
                </path>
              </svg>
            </div>
            <div class="ml-3 flex-1">
              <h4 class="font-medium text-blue-900">{{ announcement.title }}</h4>
              <p v-if="announcement.content" class="mt-1 text-sm text-blue-800">{{ announcement.content }}</p>
            </div>
          </div>
        </div>
      </div>
      <div class="space-y-4">
        <div v-for="topic in topics" :key="topic.id"
          class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex space-x-4">
            <router-link :to="`/user/${topic.user_id}`">
              <img :src="getUserAvatar(topic.user)" class="w-12 h-12 rounded-full">
            </router-link>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between mb-1">
                <div class="flex items-center space-x-2">
                  <router-link :to="`/user/${topic.user_id}`" class="font-medium text-gray-900 hover:text-blue-500">
                    {{ getUserDisplayName(topic.user) }}
                  </router-link>
                  <span v-if="topic.forum" class="text-xs bg-blue-100 text-blue-600 px-2 py-0.5 rounded">
                    {{ topic.forum.name }}
                  </span>
                </div>
                <span class="text-xs text-gray-400">{{ formatTime(topic.created_at) }}</span>
              </div>
              <router-link :to="`/topic/${topic.id}`" class="block">
                <h3 class="text-lg font-semibold mb-2 hover:text-blue-500 line-clamp-2 flex items-center flex-wrap gap-1">
                  <span v-if="topic.is_pinned" class="text-xs bg-red-500 text-white px-1.5 py-0.5 rounded font-medium">置顶</span>
                  <span v-if="topic.is_essence" class="text-xs bg-yellow-500 text-white px-1.5 py-0.5 rounded font-medium">精华</span>
                  <span v-if="topic.is_locked" class="text-xs bg-gray-500 text-white px-1.5 py-0.5 rounded font-medium">锁定</span>
                  <span class="text-gray-900">{{ topic.title }}</span>
                </h3>
                <!-- 图片预览 -->
                <div v-if="extractFirstImage(topic.content)" class="mb-3">
                  <img :src="extractFirstImage(topic.content)" class="w-full max-h-64 object-cover rounded-lg" loading="lazy">
                </div>
                <!-- 视频/投票标识（图片存在时也显示） -->
                <div v-if="hasVideo(topic.content) || topic.has_poll" class="mb-3 flex items-center gap-2">
                  <span v-if="hasVideo(topic.content)" class="inline-flex items-center gap-1 px-2 py-1 bg-purple-100 text-purple-600 rounded text-xs">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                    视频
                  </span>
                  <span v-if="topic.has_poll" class="inline-flex items-center gap-1 px-2 py-1 bg-green-100 text-green-600 rounded text-xs">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                    </svg>
                    投票
                  </span>
                </div>
                <!-- 无图片时显示文字预览 -->
                <p v-if="!extractFirstImage(topic.content) && !hasVideo(topic.content) && !topic.has_poll" class="text-gray-600 text-sm mb-3 line-clamp-3">{{ stripMarkdown(topic.content).substring(0, 200) }}</p>
              </router-link>
              <div class="flex items-center flex-wrap gap-2 mb-2" v-if="topic.tags && topic.tags.length > 0">
                <router-link v-for="tag in topic.tags" :key="tag.id" :to="`/?tag=${tag.id}`"
                  class="px-2 py-0.5 text-xs bg-blue-100 text-blue-600 rounded hover:bg-blue-200 hover:text-blue-700">
                  #{{ tag.name }}
                </router-link>
              </div>
              <div class="flex items-center space-x-6 text-sm text-gray-500">
                <button @click="toggleLike(topic)"
                  :class="['flex items-center space-x-1 transition-colors', topic.liked ? 'text-red-500' : 'hover:text-red-500']">
                  <svg class="w-4 h-4" :fill="topic.liked ? 'currentColor' : 'none'" stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
                    </path>
                  </svg>
                  <span>{{ topic.like_count }}</span>
                </button>
                <router-link :to="`/topic/${topic.id}`" class="flex items-center space-x-1 hover:text-blue-500">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z">
                    </path>
                  </svg>
                  <span>{{ topic.reply_count }}</span>
                </router-link>
                <button class="flex items-center space-x-1 hover:text-green-500">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z">
                    </path>
                  </svg>
                  <span>{{ topic.view_count }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <div ref="loadMoreTrigger" class="py-8 text-center">
          <div v-if="loading" class="flex items-center justify-center space-x-2">
            <svg class="animate-spin w-5 h-5 text-blue-500" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
              </path>
            </svg>
            <span class="text-gray-500">加载中...</span>
          </div>
          <div v-else-if="noMore" class="text-gray-400 text-sm">
            已经到底啦~
          </div>
        </div>
      </div>
    </div>
    <aside class="w-52 flex-shrink-0 hidden xl:block">
      <div class="bg-white rounded-lg shadow-sm p-3 mb-4" v-if="userStore.isLoggedIn">
        <h3 class="font-semibold text-gray-900 mb-2 text-sm">签到</h3>
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-800 mb-1">{{ signInStatus.credits || userStore.user?.credits || 0 }}
          </div>
          <div class="text-xs text-gray-500 mb-2">积分</div>
          <button @click="handleSignIn" :disabled="signInLoading || signInStatus.signed_today" :class="['w-full py-1.5 rounded-lg text-sm font-medium transition-colors',
            signInStatus.signed_today
              ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
              : 'bg-blue-500 text-white hover:bg-blue-600']">
            {{ signInLoading ? '签到中...' : (signInStatus.signed_today ? '今日已签到' : '立即签到') }}
          </button>
          <div v-if="signInStatus.last_sign_at" class="mt-2 text-xs text-gray-500">
            你已经连续签到 {{ getStreakDays() }} 天啦！
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 mb-4">
        <h3 class="font-semibold text-gray-900 mb-3">热门帖子</h3>
        <div class="space-y-3">
          <router-link v-for="t in hotTopics" :key="t.id" :to="`/topic/${t.id}`" class="block group">
            <div class="text-sm text-gray-700 group-hover:text-blue-500 line-clamp-2">{{ t.title }}</div>
            <div class="text-xs text-gray-400 mt-1">{{ t.view_count }} 浏览</div>
          </router-link>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <h3 class="font-semibold text-gray-900 mb-3">活跃用户</h3>
        <div class="space-y-3">
          <router-link v-for="(user, index) in creditUsers" :key="user.id" :to="`/user/${user.id}`"
            class="flex items-center justify-between hover:bg-gray-50 -mx-2 px-2 py-1 rounded transition-colors">
            <div class="flex items-center space-x-2">
              <img :src="getUserAvatar(user)" class="w-6 h-6 rounded-full">
              <span class="text-sm text-gray-700">{{ getUserDisplayName(user) }}</span>
            </div>
            <span class="text-xs font-medium text-gray-600">{{ user.credits }}</span>
          </router-link>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useIntersectionObserver } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { useUserStore } from '@/stores/user'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'
import { stripMarkdown, extractFirstImage, hasVideo } from '@/utils/markdown'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const tags = ref([])
const topics = ref([])
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)
const noMore = ref(false)
const loadMoreTrigger = ref(null)

const hotTopics = ref([])
const creditUsers = ref([])
const announcements = ref([])
const signInStatus = ref({
  signed_today: false,
  last_sign_at: null,
  credits: 0
})
const signInLoading = ref(false)

const currentTagId = computed(() => {
  const tagId = route.query.tag
  return tagId ? parseInt(tagId) : null
})

const currentForum = computed(() => {
  const forumId = route.query.forum
  return forumId ? parseInt(forumId) : null
})

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

async function loadTags() {
  try {
    const res = await api.get('/tags')
    tags.value = res || []
  } catch (e) {
    console.error(e)
  }
}

async function loadTopics(isLoadMore = false) {
  if (loading.value || noMore.value) return

  loading.value = true

  try {
    const params = {
      page: page.value,
      page_size: pageSize
    }
    if (currentForum.value) {
      params.forum_id = currentForum.value
    }
    if (currentTagId.value) {
      params.tag_id = currentTagId.value
    }
    const res = await api.get('/topics', { params })

    if (isLoadMore) {
      topics.value = [...topics.value, ...(res.list || [])]
    } else {
      topics.value = res.list || []
    }

    total.value = res.total || 0

    if (topics.value.length >= total.value) {
      noMore.value = true
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadHotTopics() {
  try {
    const res = await api.get('/topics', {
      params: {
        page: 1,
        page_size: 5,
        order_by: 'view_count'
      }
    })
    hotTopics.value = res.list || []
  } catch (e) {
    console.error(e)
  }
}

async function loadCreditUsers() {
  try {
    const res = await api.get('/users/credit')
    creditUsers.value = res || []
  } catch (e) {
    console.error(e)
  }
}

async function loadAnnouncements() {
  try {
    const res = await api.get('/announcements')
    announcements.value = res || []
  } catch (e) {
    console.error(e)
  }
}

function loadMore() {
  if (!loading.value && !noMore.value) {
    page.value++
    loadTopics(true)
  }
}

useIntersectionObserver(
  loadMoreTrigger,
  ([{ isIntersecting }]) => {
    if (isIntersecting) {
      loadMore()
    }
  },
  { threshold: 0.1 }
)

watch([() => route.query.forum, () => route.query.tag], () => {
  page.value = 1
  noMore.value = false
  topics.value = []
  loadTopics()
})

async function toggleLike(topic) {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
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
    ElMessage.error('操作失败')
  }
}

async function loadSignInStatus() {
  if (!userStore.isLoggedIn) return
  try {
    const res = await api.get('/user/signin/status')
    signInStatus.value.signed_today = res.signed_today || false
    signInStatus.value.last_sign_at = res.last_sign_at || null
    signInStatus.value.credits = res.credits || 0
  } catch (e) {
    console.error(e)
  }
}

async function handleSignIn() {
  if (!userStore.isLoggedIn) return
  signInLoading.value = true
  try {
    const res = await api.post('/user/signin')
    ElMessage.success(`签到成功，获得 ${res.credits} 积分，总积分 ${res.total_credits}`)
    signInStatus.value.signed_today = true
    signInStatus.value.last_sign_at = new Date().toISOString()
    signInStatus.value.credits = res.total_credits
    if (userStore.user) {
      userStore.user.credits = res.total_credits
    }
  } catch (e) {
    console.error(e)
    ElMessage.error(e.response?.data?.message || '签到失败')
  } finally {
    signInLoading.value = false
  }
}

function getStreakDays() {
  if (!signInStatus.value.last_sign_at) return 0

  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const lastSignDate = new Date(signInStatus.value.last_sign_at)
  lastSignDate.setHours(0, 0, 0, 0)

  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)

  if (lastSignDate.getTime() === today.getTime() || lastSignDate.getTime() === yesterday.getTime()) {
    return 1
  }
  return 1
}

async function checkTopicLikes() {
  if (!userStore.isLoggedIn || topics.value.length === 0) return

  try {
    const topicIds = topics.value.map(t => t.id)
    const res = await api.post('/likes/check', {
      target_type: 'topic',
      target_ids: topicIds
    })
    if (res.liked_map) {
      for (const topic of topics.value) {
        topic.liked = res.liked_map[topic.id] || false
      }
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadTags()
  loadTopics()
  loadHotTopics()
  loadCreditUsers()
  loadAnnouncements()
  if (userStore.isLoggedIn) {
    loadSignInStatus()
  }
})

// 监听 topics 加载完成后检查点赞状态（避免 deep watch 导致的无限循环）
watch(() => topics.value.length, (newLen) => {
  if (newLen > 0) {
    checkTopicLikes()
  }
})

watch(() => userStore.isLoggedIn, (isLoggedIn) => {
  if (isLoggedIn) {
    loadSignInStatus()
    checkTopicLikes()
  }
})
</script>
