<template>
  <div class="max-w-5xl mx-auto px-3 sm:px-4 py-4 sm:py-6">
    <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6 mb-4">
      <h1 class="text-lg sm:text-xl font-bold text-gray-900">{{ t('favorites.myFavorites') }}</h1>
    </div>
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
    </div>
    <div v-else-if="favorites.length === 0" class="bg-white rounded-lg shadow-sm p-6 sm:p-8 text-center">
      <svg class="w-12 h-12 sm:w-16 sm:h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor"
        viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
      </svg>
      <p class="text-gray-500 text-sm">{{ t('favorites.noFavorites') }}</p>
      <router-link to="/" class="inline-block mt-4 text-blue-500 hover:underline text-sm">
        {{ t('favorites.goHomepage') }}
      </router-link>
    </div>
    <div v-else class="space-y-3 sm:space-y-4">
      <div v-for="fav in favorites" :key="fav.id"
        class="bg-white rounded-lg shadow-sm p-3 sm:p-4 hover:shadow-md transition-shadow">
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <router-link :to="`/topic/${fav.topic?.id}`" class="block">
              <TopicCard :topic="fav.topic" :external-has-poll="fav.topic_has_poll" />
            </router-link>
            <div class="flex items-center flex-wrap gap-2 text-xs sm:text-sm text-gray-500 mt-2">
              <router-link :to="`/user/${fav.topic?.user_id}`" class="hover:text-blue-500">
                {{ getUserDisplayName(fav.topic?.user) }}
              </router-link>
              <div v-if="getAuthorBadges(fav.topic).length > 0" class="flex items-center gap-0.5">
                <SvgBadge v-for="badge in getAuthorBadges(fav.topic)" :key="badge.id" :type="badge.icon" :size="14"
                  :title="badge.name" />
              </div>
              <span v-if="fav.topic?.forum" class="text-xs bg-blue-100 text-blue-600 px-1.5 py-0.5 rounded">
                {{ fav.topic.forum.name }}
              </span>
              <span>{{ formatTime(fav.created_at) }}</span>
            </div>
          </div>
          <button @click="removeFavorite(fav)" class="text-gray-400 hover:text-red-500 ml-3 flex-shrink-0">
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
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { getUserDisplayName } from '@/utils/user'
import { getDisplayBadges } from '@/utils/badge'
import { getErrorI18nKey } from '@/utils/error'
import SvgBadge from '@/components/SvgBadge.vue'
import TopicCard from '@/components/TopicCard.vue'

const { t } = useI18n()
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

  if (seconds < 60) return t('notifications.justNow')
  if (minutes < 60) return t('notifications.minutesAgo', { 0: minutes })
  if (hours < 24) return t('notifications.hoursAgo', { 0: hours })
  if (days < 30) return t('notifications.daysAgo', { 0: days })

  return date.toLocaleDateString('zh-CN')
}

// 获取帖子作者的展示勋章（最多2枚）
function getAuthorBadges(topic) {
  if (!topic || !topic.author_badges || topic.author_badges.length === 0) return []
  return getDisplayBadges(topic.author_badges, 'post-list')
}

async function loadFavorites() {
  loading.value = true
  try {
    const res = await api.get('/user/favorites')
    favorites.value = res || []
  } catch (e) {
    console.error(e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  } finally {
    loading.value = false
  }
}

async function removeFavorite(fav) {
  try {
    await ElMessageBox.confirm(t('favorites.confirmRemove'), t('favorites.tips'), {
      confirmButtonText: t('favorites.confirm'),
      cancelButtonText: t('favorites.cancel'),
      type: 'warning'
    })
  } catch {
    return
  }

  try {
    await api.delete(`/favorites?topic_id=${fav.topic_id}`)
    favorites.value = favorites.value.filter(f => f.id !== fav.id)
    ElMessage.success(t('favorites.removed'))
  } catch (e) {
    console.error(e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

onMounted(() => {
  loadFavorites()
})
</script>
