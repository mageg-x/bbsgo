<template>
  <div class="badges-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>
              <Award :size="18" />
              勋章管理
            </h3>
            <span class="total-count">共 {{ badges.length }} 枚勋章</span>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="initBadges">
              <RefreshCw :size="16" />
              初始化勋章
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="badges" stripe style="width: 100%; font-size: 0.75rem;" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="勋章" min-width="200">
          <template #default="{ row }">
            <div class="badge-cell">
              <SvgBadge :type="row.icon" :size="32" />
              <div class="badge-info">
                <span class="badge-name">{{ row.name || '未命名' }}</span>
                <span class="badge-type">{{ getTypeName(row.type) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="描述" min-width="250">
          <template #default="{ row }">
            {{ row.description || '暂无描述' }}
          </template>
        </el-table-column>
        <el-table-column label="获得条件" min-width="180">
          <template #default="{ row }">
            {{ row.condition || '暂无条件' }}
          </template>
        </el-table-column>
        <el-table-column label="获得人数" width="120">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.award_count || 0 }} 人</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewBadgeUsers(row)">
              <Users :size="14" />
              查看用户
            </el-button>
            <el-button link type="danger" @click="deleteBadge(row)">
              <Trash2 :size="14" />
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="usersDialogVisible" :title="`${currentBadge?.name || '勋章'} - 获得用户`" width="800px">
      <el-table :data="badgeUsers" stripe style="width: 100%" v-loading="usersLoading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户信息" min-width="180">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="36" :src="row.user?.avatar">
                <User :size="18" />
              </el-avatar>
              <div class="user-info">
                <span class="username">{{ row.user?.username || '未知用户' }}</span>
                <span class="email">{{ row.user?.email || '-' }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="获得时间" width="180">
          <template #default="{ row }">
            <span class="date-text">{{ row.awarded_at ? formatDateTime(row.awarded_at) : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_revoked" type="danger" size="small">已撤销</el-tag>
            <el-tag v-else type="success" size="small">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right" v-if="!currentBadge?.is_revoked">
          <template #default="{ row }">
            <el-button v-if="!row.is_revoked" link type="danger" @click="revokeBadge(row)">
              <XCircle :size="14" />
              撤销
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="usersTotal > 10" class="pagination-wrapper">
        <el-pagination
          :current-page="usersPage"
          :page-size="10"
          :total="usersTotal"
          layout="total, prev, pager, next"
          @current-change="loadBadgeUsers"
        />
      </div>
    </el-dialog>

    <el-dialog v-model="revokeDialogVisible" title="撤销勋章" width="400px">
      <el-form :model="revokeForm" label-width="80px">
        <el-form-item label="撤销原因">
          <el-input 
            v-model="revokeForm.reason" 
            type="textarea" 
            :rows="3"
            placeholder="请输入撤销原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="revokeDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmRevoke">确定撤销</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Award, Users, User, RefreshCw, Trash2, XCircle } from 'lucide-vue-next'
import SvgBadge from '@/components/SvgBadge.vue'

const badges = ref([])
const loading = ref(false)
const usersDialogVisible = ref(false)
const usersLoading = ref(false)
const badgeUsers = ref([])
const currentBadge = ref(null)
const usersPage = ref(1)
const usersTotal = ref(0)
const revokeDialogVisible = ref(false)
const revokeForm = ref({
  id: null,
  reason: ''
})

function getTypeName(type) {
  const types = { 
    basic: '基础入门', 
    advanced: '进阶成就', 
    top: '顶级荣耀' 
  }
  return types[type] || '未知'
}

function formatDateTime(date) {
  return new Date(date).toLocaleString('zh-CN')
}

async function loadBadges() {
  loading.value = true
  try {
    const res = await api.get('/admin/badges')
    badges.value = res || []
  } catch (e) {
    console.error('加载勋章列表失败', e)
    ElMessage.error('加载勋章列表失败')
  } finally {
    loading.value = false
  }
}

async function initBadges() {
  try {
    await ElMessageBox.confirm(
      '确定要初始化系统勋章吗？这将创建9枚默认勋章。',
      '初始化勋章',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'info' }
    )
    
    await api.post('/admin/badges/init')
    ElMessage.success('勋章初始化成功')
    await loadBadges()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('初始化勋章失败', e)
      ElMessage.error('初始化勋章失败')
    }
  }
}

async function viewBadgeUsers(badge) {
  currentBadge.value = badge
  usersPage.value = 1
  await loadBadgeUsers()
  usersDialogVisible.value = true
}

async function loadBadgeUsers(page = 1) {
  usersPage.value = page
  usersLoading.value = true
  try {
    const res = await api.get(`/admin/badges/${currentBadge.value.id}/users`, {
      params: { page: usersPage.value }
    })
    badgeUsers.value = res?.list || []
    usersTotal.value = res?.total || 0
  } catch (e) {
    console.error('加载用户列表失败', e)
    ElMessage.error('加载用户列表失败')
  } finally {
    usersLoading.value = false
  }
}

async function deleteBadge(badge) {
  try {
    await ElMessageBox.confirm(
      `确定要删除勋章 "${badge.name}" 吗？此操作不可恢复！`,
      '删除勋章',
      { confirmButtonText: '删除', cancelButtonText: '取消', type: 'error' }
    )

    await api.delete(`/admin/badges/${badge.id}`)
    badges.value = badges.value.filter(b => b.id !== badge.id)
    ElMessage.success('勋章已删除')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除勋章失败', e)
      ElMessage.error('删除勋章失败')
    }
  }
}

function revokeBadge(userBadge) {
  revokeForm.value = {
    id: userBadge.id,
    reason: ''
  }
  revokeDialogVisible.value = true
}

async function confirmRevoke() {
  if (!revokeForm.value.reason) {
    ElMessage.warning('请输入撤销原因')
    return
  }

  try {
    await api.put(`/admin/badges/${revokeForm.value.id}/revoke`, revokeForm.value)
    ElMessage.success('勋章已撤销')
    revokeDialogVisible.value = false
    await loadBadgeUsers()
  } catch (e) {
    console.error('撤销勋章失败', e)
    ElMessage.error('撤销勋章失败')
  }
}

onMounted(() => {
  loadBadges()
})
</script>

<style scoped>
.badges-page {
  padding: 20px;
}

.main-card {
  border-radius: 8px;
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

.total-count {
  font-size: 13px;
  color: #6b7280;
  background: #f3f4f6;
  padding: 4px 10px;
  border-radius: 12px;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.badge-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.badge-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  font-size: 0.75rem;
}

.badge-name {
  font-weight: 500;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.badge-type {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.username {
  font-weight: 500;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.email {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

:deep(.el-table) {
  --el-table-border-color: #f3f4f6;
  --el-table-header-bg-color: #f9fafb;
}

:deep(.el-table th) {
  font-weight: 600;
  color: #374151;
}

:deep(.el-button.is-link) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>
