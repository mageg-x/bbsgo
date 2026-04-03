<template>
  <div class="layout-container">
    <!-- 侧边栏 -->
    <aside class="sidebar no-scrollbar">
      <!-- Logo 区域 -->
      <div class="logo-area">
        <img src="/src/assets/bbs.png" class="logo-icon" alt="logo" />
        <span class="logo-text">BBS Go</span>
      </div>

      <!-- 菜单分组 -->
      <nav class="menu-nav">
        <div class="menu-group">
          <span class="menu-label">主菜单</span>
          <router-link to="/console/dashboard" class="menu-item" :class="{ active: route.name === 'Dashboard' }">
            <span class="menu-icon blue">
              <LayoutDashboard :size="18" />
            </span>
            <span class="menu-text">仪表盘</span>
          </router-link>
          <router-link to="users" class="menu-item" :class="{ active: route.name === 'Users' }">
            <span class="menu-icon purple">
              <Users :size="18" />
            </span>
            <span class="menu-text">用户管理</span>
          </router-link>
          <router-link to="forums" class="menu-item" :class="{ active: route.name === 'Forums' }">
            <span class="menu-icon pink">
              <FolderOpen :size="18" />
            </span>
            <span class="menu-text">版块管理</span>
          </router-link>
          <router-link to="topics" class="menu-item" :class="{ active: route.name === 'Topics' }">
            <span class="menu-icon green">
              <FileText :size="18" />
            </span>
            <span class="menu-text">帖子管理</span>
          </router-link>
          <router-link to="comments" class="menu-item" :class="{ active: route.name === 'Comments' }">
            <span class="menu-icon cyan">
              <MessageSquare :size="18" />
            </span>
            <span class="menu-text">评论管理</span>
          </router-link>
          <router-link to="tags" class="menu-item" :class="{ active: route.name === 'Tags' }">
            <span class="menu-icon orange">
              <Tag :size="18" />
            </span>
            <span class="menu-text">话题管理</span>
          </router-link>
          <router-link to="polls" class="menu-item" :class="{ active: route.name === 'Polls' }">
            <span class="menu-icon indigo">
              <Vote :size="18" />
            </span>
            <span class="menu-text">投票管理</span>
          </router-link>
          <router-link to="badges" class="menu-item" :class="{ active: route.name === 'Badges' }">
            <span class="menu-icon yellow">
              <Award :size="18" />
            </span>
            <span class="menu-text">勋章管理</span>
          </router-link>
        </div>

        <div class="menu-group">
          <span class="menu-label">系统</span>
          <router-link to="reports" class="menu-item" :class="{ active: route.name === 'Reports' }">
            <span class="menu-icon red">
              <AlertTriangle :size="18" />
            </span>
            <span class="menu-text">举报管理</span>
            <span v-if="reportCount > 0" class="badge badge-red">{{ reportCount }}</span>
          </router-link>
          <router-link to="announcements" class="menu-item" :class="{ active: route.name === 'Announcements' }">
            <span class="menu-icon yellow">
              <Bell :size="18" />
            </span>
            <span class="menu-text">公告管理</span>
          </router-link>
          <router-link to="config" class="menu-item" :class="{ active: route.name === 'Config' }">
            <span class="menu-icon blue">
              <Settings :size="18" />
            </span>
            <span class="menu-text">网站配置</span>
          </router-link>
          <router-link to="settings" class="menu-item" :class="{ active: route.name === 'Settings' }">
            <span class="menu-icon gray">
              <Sliders :size="18" />
            </span>
            <span class="menu-text">系统设置</span>
          </router-link>
        </div>

        <div class="menu-group">
          <span class="menu-label">账户</span>
          <router-link to="change-password" class="menu-item" :class="{ active: route.name === 'ChangePassword' }">
            <span class="menu-icon teal">
              <Key :size="18" />
            </span>
            <span class="menu-text">修改密码</span>
          </router-link>
          <a @click="handleLogout" class="menu-item">
            <span class="menu-icon red">
              <LogOut :size="18" />
            </span>
            <span class="menu-text">退出登录</span>
          </a>
        </div>
      </nav>
    </aside>

    <!-- 主内容区 -->
    <div class="main-wrapper">
      <!-- 顶部栏 -->
      <header class="topbar">
        <div class="topbar-left">
          <h2 class="page-title">{{ pageTitle }}</h2>
        </div>
        <div class="topbar-right">
          <span class="datetime">
            <Calendar :size="14" />
            {{ currentDate }}
          </span>
          <div class="user-info-header">
            <div class="user-avatar-header">
              <User :size="16" />
            </div>
            <div class="user-details-header">
              <span class="user-name-header">{{ adminStore.user?.username }}</span>
              <span class="user-role-header">{{ getRoleName(adminStore.user?.role) }}</span>
            </div>
          </div>
        </div>
      </header>

      <!-- 内容区 -->
      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import api from '@/api'
import {
  LayoutDashboard, User, Users, FolderOpen, FileText, MessageSquare,
  Tag, AlertTriangle, Bell, Settings, Sliders, Key, LogOut, Calendar, Vote, Award
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminStore()

const reportCount = ref(0)

function getRoleName(role) {
  const roles = { 0: '普通用户', 1: '版主', 2: '管理员' }
  return roles[role] || '未知'
}

const currentDate = computed(() => {
  return new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})

const pageTitle = computed(() => {
  const titles = {
    'Dashboard': '仪表盘',
    'Users': '用户管理',
    'Forums': '版块管理',
    'Topics': '帖子管理',
    'Posts': '评论管理',
    'Tags': '话题管理',
    'Polls': '投票管理',
    'Reports': '举报管理',
    'Announcements': '公告管理',
    'Config': '网站配置',
    'Settings': '系统设置',
    'ChangePassword': '修改密码'
  }
  return titles[route.name] || '管理后台'
})

async function loadReportCount() {
  try {
    const res = await api.get('/admin/reports', { params: { status: 0 } })
    reportCount.value = Array.isArray(res) ? res.length : 0
  } catch (e) {
    console.error('加载举报数量失败', e)
  }
}

function handleLogout() {
  adminStore.logout()
  router.push('/console/login')
}

onMounted(() => {
  loadReportCount()
})
</script>

<style scoped>
.layout-container {
  display: flex;
  min-height: 100vh;
  background: #f0f2f5;
}

/* 侧边栏 */
.sidebar {
  width: 260px;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 菜单导航 */
.menu-nav {
  flex: 1;
  padding: 8px 12px;
  overflow-y: auto;
}

.menu-group {
  margin-bottom: 24px;
}

.menu-label {
  display: block;
  font-size: 11px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.35);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0 12px;
  margin-bottom: 8px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 4px;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.menu-item.active {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.menu-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.menu-icon.blue { background: rgba(102, 126, 234, 0.2); color: #6691ff; }
.menu-icon.purple { background: rgba(171, 122, 224, 0.2); color: #c084fc; }
.menu-icon.pink { background: rgba(244, 114, 182, 0.2); color: #f472b6; }
.menu-icon.green { background: rgba(52, 211, 153, 0.2); color: #34d399; }
.menu-icon.cyan { background: rgba(34, 211, 238, 0.2); color: #22d3ee; }
.menu-icon.orange { background: rgba(251, 146, 60, 0.2); color: #fb923c; }
.menu-icon.indigo { background: rgba(129, 140, 248, 0.2); color: #818cf8; }
.menu-icon.red { background: rgba(248, 113, 113, 0.2); color: #f87171; }
.menu-icon.yellow { background: rgba(251, 191, 36, 0.2); color: #fbbf24; }
.menu-icon.teal { background: rgba(45, 212, 191, 0.2); color: #2dd4bf; }
.menu-icon.gray { background: rgba(156, 163, 175, 0.2); color: #9ca3af; }

.menu-text {
  font-size: 14px;
  font-weight: 500;
  flex: 1;
}

.badge {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  min-width: 20px;
  text-align: center;
}

.badge-red {
  background: rgba(248, 113, 113, 0.2);
  color: #f87171;
}

/* 主内容区 */
.main-wrapper {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  overflow-x: hidden;
}

.topbar {
  height: 64px;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 10;
}

.topbar-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.datetime {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #6b7280;
}

.user-info-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px 6px 6px;
  background: #f3f4f6;
  border-radius: 24px;
}

.user-avatar-header {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.user-details-header {
  display: flex;
  flex-direction: column;
}

.user-name-header {
  font-size: 13px;
  font-weight: 600;
  color: #1f2937;
}

.user-role-header {
  font-size: 11px;
  color: #9ca3af;
}

.main-content {
  flex: 1;
  padding: 24px;
  overflow-x: auto;
}
</style>
