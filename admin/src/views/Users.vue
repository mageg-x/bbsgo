<template>
  <div class="users-page">
    <!-- 角色选择对话框 -->
    <el-dialog v-model="roleDialogVisible" title="修改用户角色" width="400px" :close-on-click-modal="false">
      <div class="role-dialog-content">
        <p class="role-dialog-user">当前用户：<strong>{{ selectedUser?.username }}</strong></p>
        <div class="role-options">
          <el-radio-group v-model="selectedRole">
            <el-radio :value="0" class="role-option">普通用户</el-radio>
            <el-radio :value="1" class="role-option">版主</el-radio>
            <el-radio :value="2" class="role-option">管理员</el-radio>
          </el-radio-group>
        </div>
      </div>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmRoleChange">确定</el-button>
      </template>
    </el-dialog>

    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>
              <Users :size="18" />
              用户列表
            </h3>
            <span class="total-count">共 {{ total }} 位用户</span>
          </div>
          <div class="header-right">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索用户名或邮箱"
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
              搜索
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="users" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户信息" min-width="180">
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
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" size="small">
              {{ getRoleName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="credits" label="积分" width="100" sortable />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_banned" type="danger" size="small">已封禁</el-tag>
            <el-tag v-else type="success" size="small">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editRole(row)">
              <Shield :size="14" />
              角色
            </el-button>
            <el-button link :type="row.is_banned ? 'success' : 'danger'" @click="toggleBan(row)">
              <Ban v-if="!row.is_banned" :size="14" />
              <CheckCircle v-else :size="14" />
              {{ row.is_banned ? '解封' : '封禁' }}
            </el-button>
            <el-button link type="danger" @click="deleteUser(row)">
              <Trash2 :size="14" />
              删除
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
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Users, User, Search, Shield, Ban, CheckCircle, Trash2 } from 'lucide-vue-next'

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
  const roles = { 0: '普通用户', 1: '版主', 2: '管理员' }
  return roles[role] || '未知'
}

function getRoleType(role) {
  const types = { 0: 'info', 1: 'warning', 2: 'danger' }
  return types[role] || 'info'
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('zh-CN')
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
    console.error('加载用户失败', e)
    ElMessage.error('加载用户失败')
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
    ElMessage.success('角色已更新')
    roleDialogVisible.value = false
  } catch (e) {
    console.error('更新角色失败', e)
    ElMessage.error('更新角色失败')
  }
}

async function toggleBan(user) {
  const action = user.is_banned ? '解封' : '封禁'
  try {
    await ElMessageBox.confirm(
      `确定要${action}用户 "${user.username}" 吗？`,
      `${action}用户`,
      { confirmButtonText: action, cancelButtonText: '取消', type: 'warning' }
    )

    await api.put(`/admin/users/${user.id}/ban`, { is_banned: !user.is_banned })
    user.is_banned = !user.is_banned
    ElMessage.success(`用户已${action}`)
  } catch (e) {
    if (e !== 'cancel') {
      console.error(`${action}用户失败`, e)
      ElMessage.error(`${action}失败`)
    }
  }
}

async function deleteUser(user) {
  try {
    await ElMessageBox.confirm(
      `确定要彻底删除用户 "${user.username}" 吗？此操作不可恢复，相关帖子、收藏、点赞等数据也将被删除！`,
      '删除用户',
      { confirmButtonText: '删除', cancelButtonText: '取消', type: 'error' }
    )

    await api.delete(`/admin/users/${user.id}`)
    users.value = users.value.filter(u => u.id !== user.id)
    total.value--
    ElMessage.success('用户已删除')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除用户失败', e)
      ElMessage.error(e.response?.data?.message || '删除用户失败')
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
  display: block;
  margin-bottom: 16px;
  height: 32px;
  line-height: 32px;
}

.role-option:last-child {
  margin-bottom: 0;
}
</style>
