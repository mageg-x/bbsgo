<template>
  <el-config-provider :locale="elementLocale">
    <div class="min-h-screen bg-gray-100">
      <nav class="bg-white shadow-sm sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-3 sm:px-4 lg:px-8">
        <div class="flex justify-between h-14 sm:h-16">
          <div class="flex items-center">
            <router-link to="/" class="flex items-center">
              <img :src="siteConfig.site_logo || defaultLogo" class="w-7 h-7 sm:w-8 sm:h-8 object-contain">
              <span class="ml-2 text-lg sm:text-xl font-bold text-gray-800">{{ siteConfig.site_name || '彩虹BBS' }}</span>
            </router-link>
          </div>
          <div v-if="!isAuthPage" class="flex items-center space-x-2 sm:space-x-4">
            <div class="relative">
              <input type="text" :placeholder="t('common.search')"
                class="w-32 sm:w-64 px-3 sm:px-4 py-1.5 sm:py-2 rounded-full border border-gray-200 text-xs sm:text-sm focus:outline-none focus:border-blue-400"
                @keypress.enter="handleSearch" v-model="searchKeyword">
              <button @click="handleSearch" class="absolute right-3 top-1.5 sm:top-2.5 text-gray-400 hover:text-blue-500">
                <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
              </button>
            </div>
            <!-- 语言切换 - 所有人都能看到 -->
            <el-dropdown @command="switchLanguage" trigger="click">
              <span class="language-btn">
                <Globe :size="16" />
                <span>{{ locale === 'zh' ? '中文' : 'EN' }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="'zh'" :class="{ 'is-active': locale === 'zh' }">
                    中文
                  </el-dropdown-item>
                  <el-dropdown-item :command="'en'" :class="{ 'is-active': locale === 'en' }">
                    English
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <template v-if="userStore.isLoggedIn">
              <router-link v-if="configStore.state.allow_post" to="/new-topic"
                class="bg-blue-500 whitespace-nowrap text-white px-2.5 sm:px-3 py-1.5 rounded-lg text-xs sm:text-sm font-medium hover:bg-blue-600">
                {{ t('common.publish') }}
              </router-link>
              <button @click="$router.push('/notifications')" class="relative text-gray-600 hover:text-blue-500">
                <svg class="w-5 h-5 sm:w-6 sm:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9">
                  </path>
                </svg>
                <span v-if="unreadCount > 0"
                  class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4.5 h-4.5 sm:w-5 sm:h-5 flex items-center justify-center">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
              </button>
              <div class="relative" ref="userMenuRef">
                <button @click="toggleUserMenu" class="flex items-center space-x-1 sm:space-x-2 hover:bg-gray-50 rounded-lg px-2 sm:px-3 py-1.5 sm:py-2 transition-colors">
                  <img :src="getUserAvatar()" class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-gray-200">
                  <span class="text-xs sm:text-sm font-medium text-gray-700 hidden sm:block">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                  <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </button>
                <div v-if="showUserMenu" class="absolute right-0 mt-2 w-44 sm:w-48 bg-white rounded-lg shadow-lg py-2 z-50">
                  <router-link
                    :to="`/user/${userStore.user?.id}`"
                    @click="closeUserMenu"
                    class="flex items-center px-3 sm:px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 mr-2 sm:mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                    </svg>
                    <span class="text-sm">{{ t('nav.personalCenter') }}</span>
                  </router-link>
                  <router-link
                    to="/messages"
                    @click="closeUserMenu"
                    class="flex items-center px-3 sm:px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 mr-2 sm:mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path>
                    </svg>
                    <span class="text-sm">{{ t('nav.myMessages') }}</span>
                  </router-link>
                  <router-link
                    to="/notifications"
                    @click="closeUserMenu"
                    class="flex items-center px-3 sm:px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 mr-2 sm:mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path>
                    </svg>
                    <span class="text-sm">{{ t('nav.systemNotifications') }}</span>
                  </router-link>
                  <router-link
                    to="/favorites"
                    @click="closeUserMenu"
                    class="flex items-center px-3 sm:px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 mr-2 sm:mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
                    </svg>
                    <span class="text-sm">{{ t('nav.myFavorites') }}</span>
                  </router-link>
                  <div class="border-t border-gray-100 my-2"></div>
                  <button
                    @click="handleLogout"
                    class="flex items-center w-full px-3 sm:px-4 py-2 text-red-600 hover:bg-red-50 transition-colors">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 mr-2 sm:mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                    </svg>
                    <span class="text-sm">{{ t('nav.logout') }}</span>
                  </button>
                </div>
              </div>

            </template>
            <template v-else>
              <router-link to="/login" class="text-gray-600 hover:text-blue-500 text-sm">{{ t('common.login') }}</router-link>
              <router-link v-if="configStore.state.allow_register" to="/register"
                class="bg-blue-500 text-white px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium hover:bg-blue-600">{{ t('common.register') }}</router-link>
            </template>
          </div>
        </div>
      </div>
      <div v-if="!isAuthPage" class="border-t border-gray-100 bg-gray-50">
        <div class="max-w-7xl mx-auto px-3 sm:px-4 lg:px-8">
          <div class="flex items-center space-x-1 py-1.5 sm:py-2 overflow-x-auto scrollbar-hide">
            <router-link
              key="all"
              to="/"
              :class="[
                'px-3 sm:px-4 py-1.5 rounded-full text-xs sm:text-sm font-medium whitespace-nowrap transition-colors',
                !route.query.forum
                  ? 'bg-blue-500 text-white'
                  : 'text-gray-600 hover:bg-gray-200'
              ]">
              {{ t('nav.allForums') }}
            </router-link>
            <router-link
              v-for="forum in forums"
              :key="forum.id"
              :to="`/?forum=${forum.id}`"
              :class="[
                'px-3 sm:px-4 py-1.5 rounded-full text-xs sm:text-sm font-medium whitespace-nowrap transition-colors',
                currentForumId === forum.id
                  ? 'bg-blue-500 text-white'
                  : 'text-gray-600 hover:bg-gray-200'
              ]">
              {{ forum.name }}
            </router-link>
          </div>
        </div>
      </div>
    </nav>
    <main :class="isAuthPage ? '' : 'max-w-7xl mx-auto px-3 sm:px-4 md:px-6 lg:px-8 py-4 sm:py-6'">
      <router-view />
    </main>
  </div>
  </el-config-provider>
</template>

<script setup>
import { ref, onMounted, computed, onBeforeUnmount, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElConfigProvider } from 'element-plus'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import en from 'element-plus/dist/locale/en.mjs'
import { useUserStore } from '@/stores/user'
import { useConfigStore } from '@/stores/config'
import api from '@/api'
import defaultLogo from '@/assets/bbs.png'
import { Globe } from 'lucide-vue-next'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const configStore = useConfigStore()
const searchKeyword = ref('')
const forums = ref([])
const siteConfig = ref({
  site_name: '',
  site_logo: '',
  site_icon: ''
})
const showUserMenu = ref(false)
const userMenuRef = ref(null)
const unreadCount = ref(0)

// Element Plus 语言映射
const elementLocale = computed(() => {
  return locale.value === 'en' ? en : zhCn
})

const currentForumId = computed(() => {
  const forumId = route.query.forum
  return forumId ? parseInt(forumId) : null
})

const isAuthPage = computed(() => {
  return route.name === 'Login' || route.name === 'Register'
})

function toggleUserMenu() {
  showUserMenu.value = !showUserMenu.value
}

function closeUserMenu() {
  showUserMenu.value = false
}

function switchLanguage(lang) {
  locale.value = lang
  localStorage.setItem('site_locale', lang)
}

function handleLogout() {
  userStore.logout()
  showUserMenu.value = false
  router.push('/')
}

function handleClickOutside(event) {
  if (userMenuRef.value && !userMenuRef.value.contains(event.target)) {
    showUserMenu.value = false
  }
}

function handleSearch() {
  if (searchKeyword.value.trim()) {
    router.push({ name: 'Search', query: { keyword: searchKeyword.value } })
  }
}

function updateFavicon(iconUrl) {
  const icon = iconUrl || defaultLogo
  
  let link = document.querySelector("link[rel*='icon']")
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.href = icon
}

function updatePageTitle(title) {
  if (title) {
    document.title = title
  }
}

function getUserAvatar() {
  if (userStore.user?.avatar) {
    return userStore.user.avatar
  }
  const username = userStore.user?.username || 'default'
  return `https://api.dicebear.com/9.x/adventurer/svg?seed=${encodeURIComponent(username)}`
}

async function loadSiteConfig() {
  try {
    const res = await api.get('/config')
    if (res) {
      siteConfig.value = {
        site_name: res.site_name || '',
        site_logo: res.site_logo || '',
        site_icon: res.site_icon || ''
      }
      
      updateFavicon(res.site_icon)
      updatePageTitle(res.site_name)
    }
  } catch (e) {
    console.error(e)
  }
}

async function loadForums() {
  try {
    const res = await api.get('/forums')
    // 在板块列表最前面加一个"全部"选项
    forums.value = res || []
  } catch (e) {
    console.error(e)
  }
}

async function loadUnreadCount() {
  if (!userStore.isLoggedIn) {
    unreadCount.value = 0
    return
  }
  try {
    const [notifRes, msgRes] = await Promise.all([
      api.get('/notifications/unread-count'),
      api.get('/messages/unread-count')
    ])
    unreadCount.value = (notifRes?.count || 0) + (msgRes?.count || 0)
  } catch (e) {
    console.error(t('common.failed'), e)
  }
}

// 监听登录状态变化，刷新未读数
watch(() => userStore.isLoggedIn, (isLoggedIn) => {
  if (isLoggedIn) {
    loadUnreadCount()
  } else {
    unreadCount.value = 0
  }
})

onMounted(async () => {
  await configStore.loadConfig()
  console.log('🔧 App onMount - 配置状态:', {
    allow_post: configStore.state.allow_post,
    loading: configStore.state.loading
  })
  loadSiteConfig()
  loadForums()
  loadUnreadCount()
  document.addEventListener('click', handleClickOutside)

  // 监听标记全部已读事件
  window.addEventListener('notifications-read-all', loadUnreadCount)

  // 每30秒刷新一次未读数量
  setInterval(loadUnreadCount, 30000)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('notifications-read-all', loadUnreadCount)
})
</script>

<style scoped>
.language-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f3f4f6;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: #6b7280;
  transition: all 0.2s;
}

.language-btn:hover {
  background: #e5e7eb;
  color: #374151;
}

:deep(.el-dropdown-menu__item.is-active) {
  color: #409eff;
  font-weight: 600;
}
</style>
