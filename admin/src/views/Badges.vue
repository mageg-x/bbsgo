<template>
  <div class="badges-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="total-count">{{ t('badge.totalBadges').replace('%d', badges.length) }}</span>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="initBadges">
              <RefreshCw :size="16" />
              {{ t('badge.initBadges') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="badges" stripe style="width: 100%; font-size: 0.75rem;" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('badge.badge')" min-width="200">
          <template #default="{ row }">
            <div class="badge-cell">
              <SvgBadge :type="row.icon" :size="32" />
              <div class="badge-info">
                <span class="badge-name">{{ row.name || '-' }}</span>
                <span class="badge-type">{{ getTypeName(row.type) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('badge.description')" min-width="250">
          <template #default="{ row }">
            {{ row.description || t('badge.noDescription') }}
          </template>
        </el-table-column>
        <el-table-column :label="t('badge.condition')" min-width="180">
          <template #default="{ row }">
            {{ row.condition || t('badge.noCondition') }}
          </template>
        </el-table-column>
        <el-table-column :label="t('badge.awardCount')" width="120">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.award_count || 0 }} {{ t('badge.people') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewBadgeUsers(row)">
              <Users :size="14" />
              {{ t('badge.viewUsers') }}
            </el-button>
            <el-button link type="danger" @click="deleteBadge(row)">
              <Trash2 :size="14" />
              {{ t('badge.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="usersDialogVisible" :title="t('badge.badgeUsers').replace('%s', currentBadge?.name || t('badge.badge'))" width="800px">
      <el-table :data="badgeUsers" stripe style="width: 100%" v-loading="usersLoading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('badge.userInfo')" min-width="180">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="36" :src="row.user?.avatar">
                <User :size="18" />
              </el-avatar>
              <div class="user-info">
                <span class="username">{{ row.user?.username || t('badge.unknownUser') }}</span>
                <span class="email">{{ row.user?.email || '-' }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('badge.awardedAt')" width="180">
          <template #default="{ row }">
            <span class="date-text">{{ row.awarded_at ? formatDateTime(row.awarded_at) : '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('badge.status')" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_revoked" type="danger" size="small">{{ t('badge.revoked2') }}</el-tag>
            <el-tag v-else type="success" size="small">{{ t('badge.normal') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="120" fixed="right" v-if="!currentBadge?.is_revoked">
          <template #default="{ row }">
            <el-button v-if="!row.is_revoked" link type="danger" @click="revokeBadge(row)">
              <XCircle :size="14" />
              {{ t('badge.revoke') }}
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

    <el-dialog v-model="revokeDialogVisible" :title="t('badge.revokeDialogTitle')" width="400px">
      <el-form :model="revokeForm" label-width="80px">
        <el-form-item :label="t('badge.revokeReason2')">
          <el-input
            v-model="revokeForm.reason"
            type="textarea"
            :rows="3"
            :placeholder="t('badge.revokeReasonPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="revokeDialogVisible = false">{{ t('badge.cancel') }}</el-button>
        <el-button type="danger" @click="confirmRevoke">{{ t('badge.confirmRevoke2') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Award, Users, User, RefreshCw, Trash2, XCircle } from 'lucide-vue-next'
import SvgBadge from '@/components/SvgBadge.vue'

const { t } = useI18n()
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
    basic: t('badge.type.basic'),
    advanced: t('badge.type.advanced'),
    top: t('badge.type.top')
  }
  return types[type] || t('badge.type.unknown')
}

function formatDateTime(date) {
  return new Date(date).toLocaleString()
}

async function loadBadges() {
  loading.value = true
  try {
    const res = await api.get('/admin/badges')
    badges.value = res || []
  } catch (e) {
    console.error('Load badges failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

async function initBadges() {
  try {
    await ElMessageBox.confirm(
      t('badge.confirmInit'),
      t('badge.init'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'info' }
    )

    await api.post('/admin/badges/init')
    ElMessage.success(t('badge.initSuccess'))
    await loadBadges()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Init badges failed', e)
      ElMessage.error(t('common.failed'))
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
    console.error('Load badge users failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    usersLoading.value = false
  }
}

async function deleteBadge(badge) {
  try {
    await ElMessageBox.confirm(
      t('badge.confirmDeleteBadge').replace('%s', badge.name),
      t('badge.deleteBadgeTitle'),
      { confirmButtonText: t('common.delete'), cancelButtonText: t('common.cancel'), type: 'error' }
    )

    await api.delete(`/admin/badges/${badge.id}`)
    badges.value = badges.value.filter(b => b.id !== badge.id)
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Delete badge failed', e)
      ElMessage.error(t('common.failed'))
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
    ElMessage.warning(t('badge.revokeReasonPlaceholder'))
    return
  }

  try {
    await api.put(`/admin/badges/${revokeForm.value.id}/revoke`, revokeForm.value)
    ElMessage.success(t('common.success'))
    revokeDialogVisible.value = false
    await loadBadgeUsers()
  } catch (e) {
    console.error('Revoke badge failed', e)
    ElMessage.error(t('common.failed'))
  }
}

onMounted(() => {
  loadBadges()
})
</script>

<style scoped>
.badges-page {
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
