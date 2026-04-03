<template>
  <div class="max-w-6xl mx-auto px-4 py-6">
    <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">勋章陈列馆</h1>
          <p class="text-gray-500 mt-1">{{ user?.username }} 的勋章成就</p>
        </div>
        <router-link :to="`/user/${userId}`" class="text-blue-500 hover:underline">
          返回个人主页
        </router-link>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="item in badgeProgress" :key="item.badge.id" 
          :class="['rounded-lg border-2 p-6 transition-all', 
            item.awarded ? 'border-blue-200 bg-blue-50' : 'border-gray-200 bg-gray-50']">
          <div class="flex items-start gap-4">
            <div :class="[item.awarded ? '' : 'grayscale opacity-50']">
              <SvgBadge :type="item.badge.icon" :size="64" />
            </div>
            <div class="flex-1">
              <h3 class="font-bold text-gray-900 mb-1">{{ item.badge.name }}</h3>
              <p class="text-sm text-gray-600 mb-2">{{ item.badge.description }}</p>
              <div class="flex items-center gap-2">
                <el-tag :type="getTypeColor(item.badge.type)" size="small">
                  {{ getTypeName(item.badge.type) }}
                </el-tag>
                <el-tag v-if="item.awarded" type="success" size="small">
                  已获得
                </el-tag>
              </div>
            </div>
          </div>

          <div v-if="item.awarded" class="mt-4 pt-4 border-t border-blue-200">
            <div class="text-xs text-gray-500">
              获得时间：{{ formatDateTime(item.awarded_at) }}
            </div>
          </div>

          <div v-else class="mt-4 pt-4 border-t border-gray-200">
            <div v-if="item.progress" class="space-y-2">
              <div v-if="item.progress.current !== undefined" class="flex items-center justify-between text-sm">
                <span class="text-gray-600">进度</span>
                <span class="font-medium">
                  {{ item.progress.current }} / {{ item.progress.target }}
                </span>
              </div>
              <el-progress 
                v-if="item.progress.current !== undefined"
                :percentage="Math.min(100, (item.progress.current / item.progress.target) * 100)" 
                :show-text="false"
              />
              <div v-if="item.progress.details" class="text-xs text-gray-500 space-y-1">
                <div v-for="(detail, key) in item.progress.details" :key="key">
                  {{ getProgressLabel(key) }}: {{ detail.current }} / {{ detail.target }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="loading" class="text-center py-12">
        <el-icon class="is-loading" :size="40"><Loading /></el-icon>
        <p class="text-gray-500 mt-4">加载中...</p>
      </div>

      <div v-if="!loading && badgeProgress.length === 0" class="text-center py-12">
        <p class="text-gray-500">暂无勋章数据</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import api from '@/api'
import SvgBadge from '@/components/SvgBadge.vue'

const route = useRoute()
const userId = route.params.id
const user = ref(null)
const badgeProgress = ref([])
const loading = ref(false)

function getTypeName(type) {
  const types = { 
    basic: '基础入门', 
    advanced: '进阶成就', 
    top: '顶级荣耀' 
  }
  return types[type] || '未知'
}

function getTypeColor(type) {
  const colors = { 
    basic: 'info', 
    advanced: 'warning', 
    top: 'danger' 
  }
  return colors[type] || 'info'
}

function getProgressLabel(key) {
  const labels = {
    register_days: '注册天数',
    topic_count: '发帖数',
    like_count: '获赞数',
    best_comment: '最佳评论'
  }
  return labels[key] || key
}

function formatDateTime(date) {
  return new Date(date).toLocaleString('zh-CN')
}

async function loadUser() {
  try {
    user.value = await api.get(`/users/${userId}`)
  } catch (e) {
    console.error('加载用户信息失败', e)
  }
}

async function loadBadgeProgress() {
  loading.value = true
  try {
    badgeProgress.value = await api.get('/badges/progress')
  } catch (e) {
    console.error('加载勋章进度失败', e)
    ElMessage.error('加载勋章进度失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    loadUser(),
    loadBadgeProgress()
  ])
})
</script>

<style scoped>
.grayscale {
  filter: grayscale(100%);
}
</style>
