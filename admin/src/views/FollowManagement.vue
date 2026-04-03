<template>
  <div class="follow-management">
    <div class="header">
      <h2>关注管理</h2>
      <div class="search-box">
        <el-input v-model="searchKeyword" placeholder="搜索用户" @keyup.enter="handleSearch" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="关注列表" name="follows">
        <el-table :data="follows" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column label="用户" width="200">
            <template #default="{ row }">
              <div class="user-info">
                <img :src="getUserAvatar(row.user)" class="avatar" />
                <div>
                  <div class="username">{{ getUserDisplayName(row.user) }}</div>
                  <div class="user-id">ID: {{ row.user_id }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="关注对象" width="200">
            <template #default="{ row }">
              <div class="user-info">
                <img :src="getUserAvatar(row.follow_user)" class="avatar" />
                <div>
                  <div class="username">{{ getUserDisplayName(row.follow_user) }}</div>
                  <div class="user-id">ID: {{ row.follow_user_id }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="关注时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button type="danger" size="small" @click="handleDeleteFollow(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="粉丝列表" name="followers">
        <el-table :data="followers" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column label="粉丝" width="200">
            <template #default="{ row }">
              <div class="user-info">
                <img :src="getUserAvatar(row.user)" class="avatar" />
                <div>
                  <div class="username">{{ getUserDisplayName(row.user) }}</div>
                  <div class="user-id">ID: {{ row.user_id }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="关注对象" width="200">
            <template #default="{ row }">
              <div class="user-info">
                <img :src="getUserAvatar(row.follow_user)" class="avatar" />
                <div>
                  <div class="username">{{ getUserDisplayName(row.follow_user) }}</div>
                  <div class="user-id">ID: {{ row.follow_user_id }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="关注时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button type="danger" size="small" @click="handleDeleteFollow(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import api from '@/api'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'

const activeTab = ref('follows')
const follows = ref([])
const followers = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

async function loadFollows() {
  loading.value = true
  try {
    const res = await api.get('/admin/follows', {
      params: {
        page: currentPage.value,
        keyword: searchKeyword.value
      }
    })
    follows.value = res.list || []
    total.value = res.total || 0
  } catch (e) {
    console.error('加载关注列表失败', e)
    ElMessage.error('加载关注列表失败')
  } finally {
    loading.value = false
  }
}

async function loadFollowers() {
  loading.value = true
  try {
    const res = await api.get('/admin/followers', {
      params: {
        page: currentPage.value,
        keyword: searchKeyword.value
      }
    })
    followers.value = res.list || []
    total.value = res.total || 0
  } catch (e) {
    console.error('加载粉丝列表失败', e)
    ElMessage.error('加载粉丝列表失败')
  } finally {
    loading.value = false
  }
}

async function handleDeleteFollow(row) {
  try {
    await ElMessageBox.confirm('确定要删除这条关注记录吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await api.delete(`/admin/follows/${row.id}`)
    ElMessage.success('删除成功')
    
    if (activeTab.value === 'follows') {
      loadFollows()
    } else {
      loadFollowers()
    }
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除失败', e)
      ElMessage.error('删除失败')
    }
  }
}

function handleTabChange(tab) {
  currentPage.value = 1
  if (tab === 'follows') {
    loadFollows()
  } else {
    loadFollowers()
  }
}

function handleSearch() {
  currentPage.value = 1
  if (activeTab.value === 'follows') {
    loadFollows()
  } else {
    loadFollowers()
  }
}

function handlePageChange(page) {
  currentPage.value = page
  if (activeTab.value === 'follows') {
    loadFollows()
  } else {
    loadFollowers()
  }
}

function formatTime(time) {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadFollows()
})
</script>

<style scoped>
.follow-management {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.search-box {
  width: 300px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.username {
  font-weight: 500;
}

.user-id {
  font-size: 12px;
  color: #999;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
