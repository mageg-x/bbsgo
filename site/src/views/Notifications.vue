<template>
  <div class="max-w-5xl mx-auto px-3 sm:px-4 py-4 sm:py-6">
    <div class="bg-white rounded-lg shadow-sm">
      <div class="p-3 sm:p-4 border-b flex justify-between items-center">
        <h2 class="text-lg sm:text-xl font-bold">{{ t('notifications.title') }}</h2>
        <button @click="markAllRead" class="text-blue-500 text-xs sm:text-sm hover:text-blue-600">{{ t('notifications.markAllRead') }}</button>
      </div>

      <div class="divide-y">
        <div v-if="notifications.length === 0" class="p-6 sm:p-8 text-center text-gray-500 text-sm">
          {{ t('notifications.noNotifications') }}
        </div>
        <div v-for="notif in notifications" :key="notif.id"
          :class="['p-3 sm:p-4 hover:bg-gray-50 cursor-pointer', !notif.is_read ? 'bg-blue-50' : '']"
          @click="handleClick(notif)">
          <div class="flex items-start">
            <div :class="['w-2 h-2 rounded-full mt-2 mr-2 sm:mr-3 flex-shrink-0',
              notif.type === 'like' ? 'bg-red-500' :
                notif.type === 'comment' ? 'bg-blue-500' :
                  notif.type === 'follow' ? 'bg-green-500' :
                    notif.type === 'badge' ? 'bg-yellow-500' :
                      notif.type === 'message' ? 'bg-purple-500' :
                        notif.type === 'best_comment' ? 'bg-orange-500' : 'bg-gray-500']"></div>
            <div v-if="notif.type === 'badge' && notif.badge" class="mr-2">
              <SvgBadge :type="notif.badge.icon" :size="24" :title="notif.badge.name" />
            </div>
            <div class="flex-1 min-w-0">
              <!-- 勋章通知 -->
              <p v-if="notif.type === 'badge'" class="text-gray-900 text-sm">
                {{ t('notifications.badgeEarned', { badgeName: notif.badge?.name || '' }) }}
              </p>
              <!-- 私信通知 -->
              <p v-else-if="notif.type === 'message'" class="text-gray-900 text-sm">
                {{ t('notifications.newMessage') }}
                <span v-if="notif.from_user" class="text-blue-500">@{{ notif.from_user.username }}</span>
              </p>
              <!-- 关注通知 -->
              <p v-else-if="notif.type === 'follow'" class="text-gray-900 text-sm">
                {{ t('notifications.userFollowedYou', { username: notif.related_user?.username || '' }) }}
              </p>
              <!-- 最佳评论通知 -->
              <p v-else-if="notif.type === 'best_comment'" class="text-gray-900 text-sm">
                {{ t('notifications.commentBest') }}
              </p>
              <!-- 其他通知直接显示 content -->
              <p v-else class="text-gray-900 text-sm">{{ notif.content }}</p>
              <span class="text-xs text-gray-400">{{ formatTime(notif.created_at) }}</span>
            </div>
            <span v-if="!notif.is_read" class="w-2 h-2 bg-blue-500 rounded-full flex-shrink-0 ml-2"></span>
          </div>
        </div>
      </div>

      <div v-if="total > pageSize" class="p-3 sm:p-4 border-t flex justify-center">
        <button v-if="page * pageSize < total" @click="loadMore" class="text-blue-500 hover:text-blue-600 text-sm">
          {{ t('notifications.loadMore') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { getErrorI18nKey } from '@/utils/error'
import SvgBadge from '@/components/SvgBadge.vue'

const { t } = useI18n()
const router = useRouter()
const notifications = ref([])
const page = ref(1)
const pageSize = 20
const total = ref(0)

function formatTime(date) {
  const d = new Date(date)
  const now = new Date()
  const diff = now - d

  if (diff < 60000) return t('notifications.justNow')
  if (diff < 3600000) return t('notifications.minutesAgo', { minutes: Math.floor(diff / 60000) })
  if (diff < 86400000) return t('notifications.hoursAgo', { hours: Math.floor(diff / 3600000) })
  if (diff < 604800000) return t('notifications.daysAgo', { days: Math.floor(diff / 86400000) })
  return d.toLocaleDateString()
}

async function loadNotifications() {
  try {
    const res = await api.get('/notifications', {
      params: { page: page.value, page_size: pageSize }
    })
    notifications.value = [...notifications.value, ...res.list]
    total.value = res.total
  } catch (e) {
    console.error(e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function loadMore() {
  page.value++
  await loadNotifications()
}

async function markAllRead() {
  try {
    await api.put('/notifications/read-all')
    notifications.value.forEach(n => n.is_read = true)
    ElMessage.success(t('notifications.markAllReadSuccess'))
    // 通知 App.vue 更新未读数
    window.dispatchEvent(new CustomEvent('notifications-read-all'))
  } catch (e) {
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

function handleClick(notif) {
  if (notif.link) {
    router.push(notif.link)
  }
}

onMounted(() => {
  loadNotifications()
})
</script>
