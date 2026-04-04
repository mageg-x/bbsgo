<template>
  <div class="follow-management">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
          </div>
          <div class="header-right">
            <el-input v-model="searchKeyword" :placeholder="t('follow.searchPlaceholder')" @keyup.enter="handleSearch" clearable style="width: 300px">
              <template #prefix>
                <Search :size="16" />
              </template>
            </el-input>
          </div>
        </div>
      </template>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="t('follow.followList')" name="follows">
        <el-table :data="follows" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column :label="t('follow.user')" width="200">
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
          <el-table-column :label="t('follow.followTarget')" width="200">
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
          <el-table-column prop="created_at" :label="t('follow.followTime')" width="180">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.actions')" width="120">
            <template #default="{ row }">
              <el-button type="danger" size="small" @click="handleDeleteFollow(row)">
                {{ t('follow.delete') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane :label="t('follow.fanList')" name="followers">
        <el-table :data="followers" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column :label="t('follow.followers')" width="200">
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
          <el-table-column :label="t('follow.followTarget')" width="200">
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
          <el-table-column prop="created_at" :label="t('follow.followTime')" width="180">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.actions')" width="120">
            <template #default="{ row }">
              <el-button type="danger" size="small" @click="handleDeleteFollow(row)">
                {{ t('follow.delete') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Users } from 'lucide-vue-next'
import api from '@/api'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'

const { t } = useI18n()
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
    console.error('Load follows failed', e)
    ElMessage.error(t('follow.loadFailed'))
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
    console.error('Load followers failed', e)
    ElMessage.error(t('follow.loadFanFailed'))
  } finally {
    loading.value = false
  }
}

async function handleDeleteFollow(row) {
  try {
    await ElMessageBox.confirm(t('follow.confirmDelete'), t('follow.deleteTitle'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.delete(`/admin/follows/${row.id}`)
    ElMessage.success(t('follow.deleteSuccess'))

    if (activeTab.value === 'follows') {
      loadFollows()
    } else {
      loadFollowers()
    }
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Delete follow failed', e)
      ElMessage.error(t('follow.deleteFailed'))
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
  return new Date(time).toLocaleString()
}

onMounted(() => {
  loadFollows()
})
</script>

<style scoped>
.follow-management {
  max-width: 1400px;
}

.main-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
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
  color: #9ca3af;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-table) {
  --el-table-border-color: #f3f4f6;
  --el-table-header-bg-color: #f9fafb;
}

:deep(.el-table th) {
  font-weight: 600;
  color: #374151;
}
</style>
