<template>
  <div class="reports-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span v-if="pendingCount > 0" class="pending-badge">
              <Bell :size="12" />
              {{ t('report.pendingCount').replace('%d', pendingCount) }}
            </span>
          </div>
          <div class="header-right">
            <el-select v-model="filterStatus" :placeholder="t('report.selectStatus')" clearable style="width: 140px">
              <el-option :label="t('report.allStatus')" value="" />
              <el-option :label="t('report.pending')" value="0" />
              <el-option :label="t('report.approved')" value="1" />
              <el-option :label="t('report.rejected')" value="2" />
            </el-select>
            <el-button type="primary" @click="loadReports">
              <RefreshCw :size="16" />
              {{ t('report.refresh') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="reports" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('report.reporter')" width="120">
          <template #default="{ row }">
            <div class="reporter-cell">
              <el-avatar :size="24" :src="row.reporter?.avatar">
                <User :size="12" />
              </el-avatar>
              <span class="text-sm">{{ row.reporter?.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.reportedUser')" width="120">
          <template #default="{ row }">
            <div v-if="row.target_user" class="flex items-center">
              <el-avatar :size="24" :src="row.target_user.avatar">
                <User :size="12" />
              </el-avatar>
              <span class="text-sm ml-1">{{ row.target_user.username }}</span>
            </div>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.type')" width="80">
          <template #default="{ row }">
            <el-tag :type="getTypeType(row.target_type)" size="small">
              {{ getTargetType(row.target_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.targetContent')" min-width="200">
          <template #default="{ row }">
            <template v-if="row.target_type === 'topic'">
              <div>
                <el-tooltip :content="row.topic_title || t('report.topic')" placement="top">
                  <span class="text-sm line-clamp-1">{{ row.topic_title || t('report.topic') }}</span>
                </el-tooltip>
                <br>
                <a :href="`/topic/${row.target_id}`" target="_blank" class="text-blue-500 hover:text-blue-700 text-xs">
                  {{ t('report.viewTopic') }}
                </a>
              </div>
            </template>
            <template v-else-if="row.target_type === 'comment'">
              <div>
                <el-tooltip :content="row.comment_content || t('report.contentDeleted')" placement="top">
                  <span class="text-sm line-clamp-2 text-gray-700">{{ row.comment_content ? row.comment_content.substring(0, 50) +
                    (row.comment_content.length > 50 ? '...' : '') : t('report.contentDeleted') }}</span>
                </el-tooltip>
                <br>
                <a v-if="row.topic_id" :href="`/topic/${row.topic_id}#post-${row.target_id}`" target="_blank"
                  class="text-blue-500 hover:text-blue-700 text-xs">
                  {{ t('report.viewComment') }}
                </a>
                <span v-else class="text-gray-400 text-xs">{{ t('report.comment') }}ID: {{ row.target_id }}</span>
              </div>
            </template>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.reason')" width="100">
          <template #default="{ row }">
            <span class="text-sm">{{ getReasonText(row.reason) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.description')" min-width="120">
          <template #default="{ row }">
            <span v-if="row.detail" class="text-sm text-gray-600">{{ row.detail }}</span>
            <span v-else class="text-gray-400">{{ t('report.noDetail') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.reportTime')" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('report.actions')" width="140" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 0">
              <el-button link type="success" @click="handleReport(row, true)">
                <Check :size="14" />
                {{ t('report.approve') }}
              </el-button>
              <el-button link type="danger" @click="handleReport(row, false)">
                <X :size="14" />
                {{ t('report.reject') }}
              </el-button>
            </template>
            <span v-else class="handled-text">{{ t('report.handled') }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { AlertTriangle, Bell, User, Check, X, RefreshCw } from 'lucide-vue-next'

const { t } = useI18n()
const reports = ref([])
const filterStatus = ref('')
const loading = ref(false)

const pendingCount = computed(() => (Array.isArray(reports.value) ? reports.value.filter(r => r.status === 0).length : 0))

function getTargetType(type) {
  const types = { topic: t('report.topic'), comment: t('report.comment'), message: t('report.message') }
  return types[type] || type
}

function getTypeType(type) {
  const types = { topic: 'primary', comment: 'primary', message: 'warning' }
  return types[type] || 'info'
}

function getStatusName(status) {
  const statuses = { 0: t('report.pending'), 1: t('report.approved'), 2: t('report.rejected') }
  return statuses[status] || 'Unknown'
}

function getStatusType(status) {
  const types = { 0: 'warning', 1: 'success', 2: 'info' }
  return types[status] || 'info'
}

function getReasonText(reason) {
  const reasons = {
    spam: t('report.reasons.spam'),
    illegal: t('report.reasons.illegal'),
    attack: t('report.reasons.attack'),
    rumor: t('report.reasons.rumor'),
    other: t('report.reasons.other')
  }
  return reasons[reason] || reason
}

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}

async function loadReports() {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value !== '') {
      params.status = filterStatus.value
    }
    const res = await api.get('/admin/reports', { params })
    reports.value = Array.isArray(res) ? res : (res?.list || [])
  } catch (e) {
    console.error('load reports failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

async function handleReport(report, approved) {
  const action = approved ? t('report.approve') : t('report.reject')
  try {
    await ElMessageBox.confirm(
      approved ? t('report.confirmApprove') : t('report.confirmReject'),
      `${action} ${t('report.title')}`,
      { confirmButtonText: action, cancelButtonText: t('common.cancel'), type: approved ? 'warning' : 'info' }
    )

    await api.put(`/admin/reports/${report.id}/handle`, { status: approved ? 1 : 2 })
    report.status = approved ? 1 : 2
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('handle report failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

onMounted(() => {
  loadReports()
})
</script>

<style scoped>
.reports-page {
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

.pending-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: rgba(248, 113, 113, 0.1);
  color: #f87171;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.reporter-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
}

.handled-text {
  font-size: 12px;
  color: #9ca3af;
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
  padding: 4px 8px;
}
</style>
