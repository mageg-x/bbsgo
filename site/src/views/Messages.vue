<template>
  <div class="max-w-5xl mx-auto px-3 sm:px-4 py-4 sm:py-6">
    <div class="bg-white rounded-xl shadow-sm overflow-hidden">
      <!-- 标题栏 -->
      <div class="px-4 sm:px-6 py-3 sm:py-4 border-b border-gray-100 flex items-center justify-between">
        <h1 class="text-lg sm:text-xl font-bold text-gray-900">{{ t('messages.privateMessages') }}</h1>
        <button @click="showSearchDialog = true" class="text-blue-500 hover:text-blue-600 text-xs sm:text-sm font-medium">
          + {{ t('messages.startPrivate') }}
        </button>
      </div>

      <div class="flex flex-col lg:flex-row h-[500px] sm:h-[600px]">
        <!-- 左侧会话列表 -->
        <div :class="['border-r border-gray-100 flex flex-col', selectedUser ? 'hidden lg:flex' : 'flex', 'lg:w-80']">
          <div v-if="conversations.length === 0" class="flex-1 flex items-center justify-center text-gray-400 p-4">
            <div class="text-center">
              <svg class="w-12 h-12 sm:w-16 sm:h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
              </svg>
              <p class="text-sm">{{ t('messages.noConversations') }}</p>
            </div>
          </div>
          <div v-else class="flex-1 overflow-y-auto">
            <div
              v-for="conv in conversations"
              :key="conv.user_id"
              @click="selectConversation(conv)"
              :class="[
                'p-3 sm:p-4 cursor-pointer hover:bg-gray-50 border-b border-gray-50 transition-colors',
                selectedUser?.id === conv.user_id ? 'bg-blue-50 border-l-4 border-l-blue-500' : ''
              ]"
            >
              <div class="flex items-center">
                <div class="relative">
                  <img :src="getUserAvatar(conv.user)" class="w-10 h-10 sm:w-12 sm:h-12 rounded-full bg-gray-200">
                  <span v-if="conv.unread_count > 0" class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4.5 h-4.5 sm:w-5 sm:h-5 flex items-center justify-center">
                    {{ conv.unread_count > 9 ? '9+' : conv.unread_count }}
                  </span>
                </div>
                <div class="ml-3 flex-1 min-w-0">
                  <div class="flex justify-between items-center">
                    <span class="font-medium text-gray-900 truncate text-sm">{{ conv.user?.nickname || conv.user?.username }}</span>
                    <span class="text-xs text-gray-400">{{ formatTime(conv.last_message?.created_at) }}</span>
                  </div>
                  <p class="text-xs sm:text-sm text-gray-500 truncate mt-1">{{ conv.last_message?.content || t('messages.noMessages') }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧消息区域 -->
        <div :class="['flex-1 flex flex-col', selectedUser ? 'flex' : 'hidden lg:flex']">
          <template v-if="selectedUser">
            <!-- 聊天对象信息 -->
            <div class="p-3 sm:p-4 border-b border-gray-100 flex items-center bg-gray-50">
              <button @click="selectedUser = null" class="lg:hidden mr-2 text-gray-500">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                </svg>
              </button>
              <img :src="getUserAvatar(selectedUser)" class="w-9 h-9 sm:w-10 sm:h-10 rounded-full bg-gray-200">
              <div class="ml-2 sm:ml-3 flex-1 min-w-0">
                <div class="font-medium text-gray-900 text-sm">{{ getUserDisplayName(selectedUser) }}</div>
                <div class="text-xs text-gray-500 hidden sm:block">{{ t('messages.clickAvatarHint') }}</div>
              </div>
              <router-link :to="`/user/${selectedUser.id}`" class="text-blue-500 hover:text-blue-600 text-xs sm:text-sm ml-auto">
                {{ t('messages.viewProfile') }}
              </router-link>
            </div>

            <!-- 消息列表 -->
            <div ref="messageListRef" class="flex-1 overflow-y-auto p-3 sm:p-4 space-y-3 sm:space-y-4 bg-gray-50">
              <div v-if="messages.length === 0" class="flex items-center justify-center h-full text-gray-400 text-sm">
                {{ t('messages.startChat', { 0: getUserDisplayName(selectedUser) }) }}
              </div>
              <div v-for="msg in messages" :key="msg.id">
                <!-- 自己发送的消息 -->
                <div v-if="msg.from_user_id === currentUserId" class="flex justify-end">
                  <div class="max-w-[80%] sm:max-w-[70%]">
                    <div class="bg-blue-500 text-white rounded-2xl rounded-tr-sm px-3 sm:px-4 py-2 shadow-sm">
                      <p class="whitespace-pre-wrap break-words text-sm">{{ msg.content }}</p>
                    </div>
                    <div class="text-xs text-gray-400 mt-1 text-right">{{ formatTime(msg.created_at) }}</div>
                  </div>
                </div>
                <!-- 收到的消息 -->
                <div v-else class="flex justify-start">
                  <img :src="getUserAvatar(selectedUser)" class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-gray-200 flex-shrink-0">
                  <div class="ml-2 max-w-[80%] sm:max-w-[70%]">
                    <div class="bg-white text-gray-800 rounded-2xl rounded-tl-sm px-3 sm:px-4 py-2 shadow-sm">
                      <p class="whitespace-pre-wrap break-words text-sm">{{ msg.content }}</p>
                    </div>
                    <div class="text-xs text-gray-400 mt-1">{{ formatTime(msg.created_at) }}</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 输入框 -->
            <div class="p-3 sm:p-4 border-t border-gray-100 bg-white">
              <form @submit.prevent="sendMessage" class="flex items-center gap-2 sm:gap-3">
                <input
                  type="text"
                  v-model="newMessage"
                  :placeholder="t('messages.typeMessage')"
                  class="flex-1 px-3 sm:px-4 py-2.5 sm:py-3 border border-gray-200 rounded-full focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-50 transition-all text-sm"
                >
                <button
                  type="submit"
                  :disabled="!newMessage.trim()"
                  class="bg-blue-500 text-white px-4 sm:px-6 py-2.5 sm:py-3 rounded-full hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm"
                >
                  {{ t('messages.send') }}
                </button>
              </form>
            </div>
          </template>

          <div v-else class="flex-1 flex items-center justify-center bg-gray-50">
            <div class="text-center text-gray-400">
              <svg class="w-16 h-16 sm:w-20 sm:h-20 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
              </svg>
              <p class="text-base sm:text-lg">{{ t('messages.selectConversation') }}</p>
              <p class="text-xs sm:text-sm mt-2">{{ t('messages.or') }} <button @click="showSearchDialog = true" class="text-blue-500 hover:underline">{{ t('messages.startNewConversation') }}</button></p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 搜索用户对话框 -->
    <el-dialog v-model="showSearchDialog" :title="t('messages.startConversation')" width="450px" :close-on-click-modal="false">
      <div class="mb-4">
        <el-input
          v-model="searchKeyword"
          :placeholder="t('messages.searchUser')"
          prefix-icon="Search"
          @input="handleSearchUser"
          clearable
        />
      </div>
      <div v-if="searchResults.length > 0" class="space-y-2 max-h-64 overflow-y-auto">
        <div
          v-for="user in searchResults"
          :key="user.id"
          @click="startConversation(user)"
          class="flex items-center p-3 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors"
        >
          <img :src="getUserAvatar(user)" class="w-10 h-10 rounded-full bg-gray-200">
          <div class="ml-3">
            <div class="font-medium text-gray-900">{{ getUserDisplayName(user) }}</div>
            <div class="text-xs text-gray-500">@{{ user.username }}</div>
          </div>
        </div>
      </div>
      <div v-else-if="searchKeyword && !searching" class="text-center text-gray-400 py-8">
        {{ t('messages.userNotFound') }}
      </div>
      <div v-else class="text-center text-gray-400 py-8">
        {{ t('messages.searchUsername') }}
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onActivated } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'

const { t } = useI18n()

const router = useRouter()
const route = useRoute()
const conversations = ref([])
const messages = ref([])
const selectedUser = ref(null)
const newMessage = ref('')
const currentUserId = ref(0)
const messageListRef = ref(null)
const showSearchDialog = ref(false)
const searchKeyword = ref('')
const searchResults = ref([])
const searching = ref(false)
let searchTimer = null

function formatTime(date) {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now - d
  const oneDay = 24 * 60 * 60 * 1000

  if (diff < 60000) return t('notifications.justNow')
  if (diff < 3600000) return t('notifications.minutesAgo', { 0: Math.floor(diff / 60000) })
  if (diff < oneDay) return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  if (diff < 7 * oneDay) return t('notifications.daysAgo', { 0: Math.floor(diff / oneDay) })
  return d.toLocaleDateString('zh-CN')
}

async function loadConversations() {
  try {
    const res = await api.get('/messages')
    const userMap = new Map()

    if (!res?.list) return

    // 按用户分组，获取每个用户最后一条消息和未读数
    res.list.forEach(msg => {
      // 确定对方用户
      const isSent = msg.from_user_id === currentUserId.value
      const otherUser = isSent ? msg.to_user : msg.from_user
      const otherId = isSent ? msg.to_user_id : msg.from_user_id

      if (!otherId || !otherUser) return

      if (!userMap.has(otherId)) {
        userMap.set(otherId, {
          user_id: otherId,
          user: otherUser,
          last_message: msg,
          unread_count: 0
        })
      }

      // 更新最后一条消息
      const conv = userMap.get(otherId)
      if (new Date(msg.created_at) > new Date(conv.last_message.created_at)) {
        conv.last_message = msg
      }

      // 统计未读数
      if (!isSent && !msg.is_read) {
        userMap.get(otherId).unread_count++
      }
    })

    // 转换为数组并按最后消息时间排序
    conversations.value = Array.from(userMap.values())
      .sort((a, b) => new Date(b.last_message.created_at) - new Date(a.last_message.created_at))

    // 更新用户头像
    conversations.value.forEach(conv => {
      if (!conv.user.avatar) {
        conv.user.avatar = getUserAvatar(conv.user)
      }
    })
  } catch (e) {
    console.error(t('messages.loadFailed'), e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function selectConversation(conv) {
  selectedUser.value = conv.user
  await loadMessages(conv.user_id)
  conv.unread_count = 0
}

async function loadMessages(userId) {
  try {
    messages.value = await api.get(`/messages/with/${userId}`)
    await nextTick()
    scrollToBottom()
  } catch (e) {
    console.error(t('messages.loadFailed'), e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

async function sendMessage() {
  if (!newMessage.value.trim() || !selectedUser.value) return

  try {
    const msg = await api.post('/messages', {
      to_user_id: selectedUser.value.id,
      content: newMessage.value
    })
    messages.value.push(msg)
    newMessage.value = ''
    await nextTick()
    scrollToBottom()

    // 更新会话列表中的最后一条消息
    const conv = conversations.value.find(c => c.user_id === selectedUser.value.id)
    if (conv) {
      conv.last_message = msg
      // 移到顶部
      conversations.value = [conv, ...conversations.value.filter(c => c.user_id !== selectedUser.value.id)]
    }
  } catch (e) {
    console.error(t('messages.sendFailed'), e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

function scrollToBottom() {
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}

function handleSearchUser() {
  if (searchTimer) clearTimeout(searchTimer)

  if (!searchKeyword.value.trim()) {
    searchResults.value = []
    return
  }

  searchTimer = setTimeout(async () => {
    searching.value = true
    try {
      const res = await api.get('/users/search', { params: { q: searchKeyword.value } })
      searchResults.value = (res?.list || []).filter(u => u.id !== currentUserId.value)
      searchResults.value.forEach(user => {
        if (!user.avatar) {
          user.avatar = getUserAvatar(user)
        }
      })
    } catch (e) {
      console.error(t('messages.searchFailed'), e)
      ElMessage.error(t(getErrorI18nKey(e?.code)))
      searchResults.value = []
    } finally {
      searching.value = false
    }
  }, 300)
}

async function startConversation(user) {
  // 检查是否已有会话
  let conv = conversations.value.find(c => c.user_id === user.id)

  if (!conv) {
    // 创建新会话
    conv = {
      user_id: user.id,
      user: user,
      last_message: { content: '', created_at: new Date().toISOString() },
      unread_count: 0
    }
  }

  showSearchDialog.value = false
  searchKeyword.value = ''
  searchResults.value = []

  await selectConversation(conv)
}

// 暴露刷新方法供外部调用
defineExpose({ loadConversations })

onMounted(async () => {
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  currentUserId.value = user.id
  await loadConversations()

  // 处理从个人主页跳转过来的情况
  const userId = route.query.userId
  if (userId) {
    const targetUserId = parseInt(userId)
    // 检查是否已在会话列表中
    let conv = conversations.value.find(c => c.user_id === targetUserId)
    if (!conv) {
      // 获取用户信息并创建会话
      try {
        const userInfo = await api.get(`/users/${targetUserId}`)
        conv = {
          user_id: targetUserId,
          user: userInfo,
          last_message: { content: '', created_at: new Date().toISOString() },
          unread_count: 0
        }
      } catch (e) {
        console.error(t('messages.cannotSend'), e)
        ElMessage.error(t(getErrorI18nKey(e?.code)))
        return
      }
    }
    await selectConversation(conv)
  }
})

// 每次激活组件时刷新
onActivated(async () => {
  if (currentUserId.value) {
    await loadConversations()
    if (selectedUser.value) {
      await loadMessages(selectedUser.value.id)
    }
  }
})
</script>

<style scoped>
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background-color: rgba(0, 0, 0, 0.2);
}
</style>
