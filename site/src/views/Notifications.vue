<template>
  <div class="max-w-4xl mx-auto px-4 py-6">
    <div class="bg-white rounded-lg shadow-sm">
      <div class="p-4 border-b flex justify-between items-center">
        <h2 class="text-xl font-bold">通知</h2>
        <button @click="markAllRead" class="text-blue-500 text-sm hover:text-blue-600">全部标记已读</button>
      </div>

      <div class="divide-y">
        <div v-if="notifications.length === 0" class="p-8 text-center text-gray-500">
          暂无通知
        </div>
        <div v-for="notif in notifications" :key="notif.id"
          :class="['p-4 hover:bg-gray-50 cursor-pointer', !notif.is_read ? 'bg-blue-50' : '']"
          @click="handleClick(notif)">
          <div class="flex items-start">
            <div :class="['w-2 h-2 rounded-full mt-2 mr-3 flex-shrink-0',
              notif.type === 'like' ? 'bg-red-500' :
                notif.type === 'comment' ? 'bg-blue-500' :
                  notif.type === 'follow' ? 'bg-green-500' :
                    notif.type === 'badge' ? 'bg-yellow-500' : 'bg-gray-500']"></div>
            <div v-if="notif.type === 'badge' && notif.badge" class="mr-2">
              <SvgBadge :type="notif.badge.icon" :size="28" :title="notif.badge.name" />
            </div>
            <div class="flex-1">
              <p class="text-gray-900">{{ notif.content }}</p>
              <span class="text-xs text-gray-400">{{ formatTime(notif.created_at) }}</span>
            </div>
            <span v-if="!notif.is_read" class="w-2 h-2 bg-blue-500 rounded-full"></span>
          </div>
        </div>
      </div>

      <div v-if="total > pageSize" class="p-4 border-t flex justify-center">
        <button v-if="page * pageSize < total" @click="loadMore" class="text-blue-500 hover:text-blue-600">
          加载更多
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import SvgBadge from '@/components/SvgBadge.vue'

const router = useRouter()
const notifications = ref([])
const page = ref(1)
const pageSize = 20
const total = ref(0)

function formatTime(date) {
  const d = new Date(date)
  const now = new Date()
  const diff = now - d

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前`
  return d.toLocaleDateString('zh-CN')
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
  } catch (e) {
    console.error(e)
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
