<template>
  <div class="users-page">
    <!-- 角色选择对话框 -->
    <el-dialog v-model="roleDialogVisible" :title="t('user.editRole')" width="400px" :close-on-click-modal="false">
      <div class="role-dialog-content">
        <p class="role-dialog-user">{{ t('user.currentUser') }}：<strong>{{ selectedUser?.username }}</strong></p>
        <div class="role-options">
          <el-radio-group v-model="selectedRole">
            <el-radio :value="0" class="role-option">{{ t('user.roleOptions.0') }}</el-radio>
            <el-radio :value="1" class="role-option">{{ t('user.roleOptions.1') }}</el-radio>
            <el-radio :value="2" class="role-option">{{ t('user.roleOptions.2') }}</el-radio>
          </el-radio-group>
        </div>
      </div>
      <template #footer>
        <el-button @click="roleDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="confirmRoleChange">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="total-count">{{ t('user.totalUsers').replace('%d', total) }}</span>
          </div>
          <div class="header-right">
            <el-input
              v-model="searchKeyword"
              :placeholder="t('user.searchPlaceholder')"
              clearable
              @clear="loadUsers"
              @keyup.enter="loadUsers"
              style="width: 220px"
            >
              <template #prefix>
                <Search :size="16" />
              </template>
            </el-input>
            <el-button type="primary" @click="loadUsers">
              <Search :size="16" />
              {{ t('common.search') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="users" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('user.userInfo')" min-width="180">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="36" :src="row.avatar">
                <User :size="18" />
              </el-avatar>
              <div class="user-info">
                <span class="username">{{ row.username }}</span>
                <span class="email">{{ row.email }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('user.role')" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" size="small">
              {{ getRoleName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="credits" :label="t('user.credits')" width="100" sortable />
        <el-table-column :label="t('user.status')" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_banned" type="danger" size="small">{{ t('user.banned') }}</el-tag>
            <el-tag v-else type="success" size="small">{{ t('user.normal') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('user.createdAt')" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('user.actions')" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editRole(row)">
              <Shield :size="14" />
              {{ t('user.role') }}
            </el-button>
            <el-button link :type="row.is_banned ? 'success' : 'danger'" @click="toggleBan(row)">
              <Ban v-if="!row.is_banned" :size="14" />
              <CheckCircle v-else :size="14" />
              {{ row.is_banned ? t('user.unban') : t('user.ban') }}
            </el-button>
            <el-button link type="danger" @click="deleteUser(row)">
              <Trash2 :size="14" />
              {{ t('user.remove') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadUsers"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Users, User, Search, Shield, Ban, CheckCircle, Trash2 } from 'lucide-vue-next'

const { t } = useI18n()
const users = ref([])
const searchKeyword = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)
const roleDialogVisible = ref(false)
const selectedUser = ref(null)
const selectedRole = ref(0)

function getRoleName(role) {
  return t(`user.roleOptions.${role}`) || t('common.noData')
}

function getRoleType(role) {
  const types = { 0: 'info', 1: 'warning', 2: 'danger' }
  return types[role] || 'info'
}

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}

async function loadUsers() {
  loading.value = true
  try {
    const res = await api.get('/admin/users', {
      params: {
        page: page.value,
        page_size: pageSize,
        keyword: searchKeyword.value
      }
    })
    users.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('load users failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

function editRole(user) {
  selectedUser.value = user
  selectedRole.value = user.role
  roleDialogVisible.value = true
}

async function confirmRoleChange() {
  if (!selectedUser.value) return
  try {
    await api.put(`/admin/users/${selectedUser.value.id}/role`, { role: selectedRole.value })
    selectedUser.value.role = selectedRole.value
    ElMessage.success(t('common.success'))
    roleDialogVisible.value = false
  } catch (e) {
    console.error('update role failed', e)
    ElMessage.error(t('common.failed'))
  }
}

async function toggleBan(user) {
  const action = user.is_banned ? t('user.unban') : t('user.ban')
  try {
    await ElMessageBox.confirm(
      `${t('common.confirm')} ${action} "${user.username}"?`,
      action,
      { confirmButtonText: action, cancelButtonText: t('common.cancel'), type: 'warning' }
    )

    await api.put(`/admin/users/${user.id}/ban`, { is_banned: !user.is_banned })
    user.is_banned = !user.is_banned
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error(`${action} user failed`, e)
      ElMessage.error(t('common.failed'))
    }
  }
}

async function deleteUser(user) {
  try {
    await ElMessageBox.confirm(
      t('user.deleteConfirm'),
      t('user.remove'),
      { confirmButtonText: t('user.remove'), cancelButtonText: t('common.cancel'), type: 'error' }
    )

    await api.delete(`/admin/users/${user.id}`)
    users.value = users.value.filter(u => u.id !== user.id)
    total.value--
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('delete user failed', e)
      ElMessage.error(e.response?.data?.message || t('common.failed'))
    }
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.users-page {
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

:deep(.el-avatar) {
  flex-shrink: 0;
}

.email {
  font-size: 12px;
  color: #9ca3af;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
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

:deep(.el-button.is-link) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.role-dialog-content {
  padding: 10px 0;
}

.role-dialog-user {
  margin-bottom: 24px;
  font-size: 14px;
  color: #606266;
}

.role-options {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px 20px;
}

.role-option {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  height: 32px;
}

:deep(.role-option .el-radio) {
  display: flex;
  align-items: center;
}

:deep(.role-option .el-radio__label) {
  margin-left: 8px;
}

.role-option:last-child {
  margin-bottom: 0;
}
</style>
