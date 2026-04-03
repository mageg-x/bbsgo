<template>
  <div class="max-w-6xl mx-auto px-4 py-6">
    <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">{{ title }}</h1>
          <p class="text-gray-500 mt-1">共 {{ total }} 人</p>
        </div>
        <router-link :to="`/user/${userId}`" class="text-blue-500 hover:underline">
          返回个人主页
        </router-link>
      </div>

      <div v-if="loading" class="text-center py-8">
        <el-icon class="is-loading" :size="32">
          <Loading />
        </el-icon>
      </div>

      <div v-else-if="list.length === 0" class="text-center py-12 text-gray-500">
        暂无{{ type === 'follows' ? '关注' : '粉丝' }}
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="item in list" :key="item.id" 
          class="flex items-center gap-3 p-4 rounded-lg border hover:shadow-md transition-shadow">
          <router-link :to="`/user/${getUserID(item)}`">
            <img :src="getUserAvatar(getUserInfo(item))" 
              class="w-12 h-12 rounded-full object-cover">
          </router-link>
          <div class="flex-1 min-w-0">
            <router-link :to="`/user/${getUserID(item)}`" 
              class="font-medium text-gray-900 hover:text-blue-500 truncate block">
              {{ getUserDisplayName(getUserInfo(item)) }}
            </router-link>
            <p class="text-sm text-gray-500 truncate">
              {{ getUserInfo(item)?.signature || '这个人很懒，什么都没写' }}
            </p>
          </div>
          <FollowButton v-if="type === 'follows'" :user-id="getUserID(item)" />
        </div>
      </div>

      <div v-if="total > pageSize" class="flex justify-center mt-6">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import api from '@/api'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'
import FollowButton from '@/components/FollowButton.vue'

const route = useRoute()
const userId = computed(() => route.params.id)
const type = computed(() => route.query.type || 'follows')

const title = computed(() => {
  return type.value === 'follows' ? '我的关注' : '我的粉丝'
})

const list = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

function getUserInfo(item) {
  return type.value === 'follows' ? item.follow_user : item.user
}

function getUserID(item) {
  const user = getUserInfo(item)
  return user?.id || 0
}

async function loadList() {
  loading.value = true
  try {
    const url = type.value === 'follows' 
      ? '/user/follows' 
      : `/users/${userId.value}/followers`
    
    const res = await api.get(url, {
      params: { page: page.value }
    })
    
    list.value = res.list || res || []
    total.value = res.total || list.value.length
  } catch (e) {
    console.error('加载列表失败', e)
    ElMessage.error('加载列表失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(newPage) {
  page.value = newPage
  loadList()
}

onMounted(() => {
  loadList()
})
</script>
