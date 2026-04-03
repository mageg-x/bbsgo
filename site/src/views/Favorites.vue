<template>
  <div class="max-w-4xl mx-auto">
    <div class="bg-white rounded-lg shadow-sm p-6 mb-4">
      <h1 class="text-xl font-bold text-gray-900">我的收藏</h1>
    </div>
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
    </div>
    <div v-else-if="favorites.length === 0" class="bg-white rounded-lg shadow-sm p-8 text-center">
      <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
      </svg>
      <p class="text-gray-500">还没有收藏任何帖子</p>
      <router-link to="/" class="inline-block mt-4 text-blue-500 hover:underline">
        去首页逛逛吧
      </router-link>
    </div>
    <div v-else class="space-y-4">
      <div v-for="fav in favorites" :key="fav.id"
        class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <router-link :to="`/topic/${fav.topic?.id}`" class="block">
              <TopicCard :topic="fav.topic" :external-has-poll="fav.topic_has_poll" />
            </router-link>
            <div class="flex items-center space-x-3 text-sm text-gray-500">
              <router-link :to="`/user/${fav.topic?.user_id}`" class="hover:text-blue-500">
                {{ getUserDisplayName(fav.topic?.user) }}
              </router-link>
              <span v-if="fav.topic?.forum" class="text-xs bg-blue-100 text-blue-600 px-2 py-0.5 rounded">
                {{ fav.topic.forum.name }}
              </span>
              <span>{{ formatTime(fav.created_at) }}</span>
            </div>
          </div>
          <button @click="removeFavorite(fav)" class="text-gray-400 hover:text-red-500 ml-4">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
              </path>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { getUserDisplayName } from '@/utils/user'
import TopicCard from '@/components/TopicCard.vue'

const userStore = useUserStore()
const loading = ref(false)
const favorites = ref([])

function formatTime(timeStr) {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (seconds < 60) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`

  return date.toLocaleDateString('zh-CN')
}

async function loadFavorites() {
  loading.value = true
  try {
    const res = await api.get('/user/favorites')
    favorites.value = res || []
  } catch (e) {
    console.error(e)
    ElMessage.error('加载收藏列表失败')
  } finally {
    loading.value = false
  }
}

async function removeFavorite(fav) {
  try {
    await ElMessageBox.confirm('确认取消收藏?', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
  } catch {
    return
  }

  try {
    await api.delete(`/favorites?topic_id=${fav.topic_id}`)
    favorites.value = favorites.value.filter(f => f.id !== fav.id)
    ElMessage.success('已取消收藏')
  } catch (e) {
    console.error(e)
    ElMessage.error('取消收藏失败')
  }
}

onMounted(() => {
  loadFavorites()
})
</script>
