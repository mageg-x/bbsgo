<template>
  <div class="polls-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="total-count">{{ t('poll.totalPolls').replace('%d', total) }}</span>
          </div>
          <div class="header-right">
            <el-select v-model="statusFilter" :placeholder="t('poll.statusFilter')" clearable @change="loadPolls" style="width: 120px">
              <el-option :label="t('poll.ongoing')" value="active" />
              <el-option :label="t('poll.ended')" value="ended" />
            </el-select>
            <el-input
              v-model="searchKeyword"
              :placeholder="t('poll.searchPlaceholder')"
              clearable
              @clear="loadPolls"
              @keyup.enter="loadPolls"
              style="width: 220px"
            >
              <template #prefix>
                <Search :size="16" />
              </template>
            </el-input>
            <el-button type="primary" @click="loadPolls">
              <Search :size="16" />
              {{ t('poll.search') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="polls" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('poll.pollTitle')" min-width="200">
          <template #default="{ row }">
            <div class="title-cell">
              <el-link :href="`/topic/${row.topic?.id}`" target="_blank" type="primary">
                {{ row.poll?.title || t('poll.noTitle') }}
              </el-link>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('poll.pollType')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.poll?.poll_type === 'single' ? 'primary' : 'success'" size="small">
              {{ row.poll?.poll_type === 'single' ? t('poll.single') : t('poll.multiple') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('poll.optionCount')" width="80">
          <template #default="{ row }">
            <span>{{ row.poll?.options?.length || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('poll.voterCount')" width="100">
          <template #default="{ row }">
            <span class="vote-count">{{ row.poll?.total_votes || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('poll.status')" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.poll?.is_ended" type="info" size="small">{{ t('poll.ended2') }}</el-tag>
            <el-tag v-else-if="isEnded(row.poll?.end_time)" type="info" size="small">{{ t('poll.ended2') }}</el-tag>
            <el-tag v-else type="success" size="small">{{ t('poll.ongoing') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('poll.endTime')" width="160">
          <template #default="{ row }">
            <span class="date-text">{{ row.poll?.end_time ? formatDateTime(row.poll.end_time) : t('poll.permanent') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('poll.createdAt')" width="160">
          <template #default="{ row }">
            <span class="date-text">{{ formatDateTime(row.poll?.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewPoll(row)">
              <ExternalLink :size="14" />
              {{ t('poll.view') }}
            </el-button>
            <el-button link type="warning" @click="endPoll(row)" v-if="!row.poll?.is_ended && !isEnded(row.poll?.end_time)">
              <Square :size="14" />
              {{ t('poll.end') }}
            </el-button>
            <el-button link type="danger" @click="deletePoll(row)">
              <Trash2 :size="14" />
              {{ t('poll.delete') }}
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
          @current-change="loadPolls"
        />
      </div>
    </el-card>

    <el-dialog v-model="pollDialogVisible" :title="t('poll.pollDetails')" width="600px">
      <div v-if="currentPoll" class="poll-detail">
        <div class="poll-info">
          <h4>{{ currentPoll.title || t('poll.noTitle') }}</h4>
          <p class="poll-meta">
            <el-tag :type="currentPoll.poll_type === 'single' ? 'primary' : 'success'" size="small">
              {{ currentPoll.poll_type === 'single' ? t('poll.single') : t('poll.multiple') + '(' + t('poll.maxChoices').replace('%d', currentPoll.max_choices) + ')' }}
            </el-tag>
            <span class="vote-total">{{ t('poll.totalVoters').replace('%d', currentPoll.total_votes || 0) }}</span>
          </p>
        </div>

        <div class="options-list">
          <div v-for="option in currentPoll.options" :key="option.id" class="option-item">
            <div class="option-header">
              <span class="option-text">{{ option.text }}</span>
              <span class="option-votes">{{ option.vote_count }} {{ t('poll.votes') }} ({{ getPercentage(option.vote_count) }}%)</span>
            </div>
            <el-progress :percentage="getPercentage(option.vote_count)" :show-text="false" />
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Vote, Search, ExternalLink, Trash2, Square } from 'lucide-vue-next'

const { t } = useI18n()
const polls = ref([])
const searchKeyword = ref('')
const statusFilter = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)
const pollDialogVisible = ref(false)
const currentPoll = ref(null)

function formatDateTime(date) {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

function isEnded(endTime) {
  if (!endTime) return false
  return new Date(endTime) < new Date()
}

function getPercentage(voteCount) {
  if (!currentPoll.value || currentPoll.value.total_votes === 0) return 0
  return Math.round((voteCount / currentPoll.value.total_votes) * 100)
}

async function loadPolls() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize,
      keyword: searchKeyword.value
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    const res = await api.get('/admin/polls', { params })
    polls.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('load polls failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

function viewPoll(row) {
  currentPoll.value = row.poll
  pollDialogVisible.value = true
}

async function endPoll(row) {
  try {
    await ElMessageBox.confirm(t('poll.confirmEnd').replace('%s', row.poll?.title || t('poll.noTitle')), t('poll.end'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.post(`/admin/polls/${row.poll.id}/end`)
    row.poll.is_ended = true
    ElMessage.success(t('common.success'))
    loadPolls()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('end poll failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

async function deletePoll(row) {
  try {
    await ElMessageBox.confirm(t('poll.confirmDelete').replace('%s', row.poll?.title || t('poll.noTitle')), t('poll.delete'), {
      confirmButtonText: t('common.delete'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.delete(`/admin/polls/${row.poll.id}`)
    polls.value = polls.value.filter(p => p.poll.id !== row.poll.id)
    total.value--
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('delete poll failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

onMounted(() => {
  loadPolls()
})
</script>

<style scoped>
.polls-page {
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

.title-cell {
  max-width: 300px;
}

.vote-count {
  font-weight: 500;
  color: #3b82f6;
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

.poll-detail {
  padding: 16px 0;
}

.poll-info h4 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 12px 0;
}

.poll-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0 0 20px 0;
}

.vote-total {
  font-size: 14px;
  color: #6b7280;
}

.options-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.option-item {
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.option-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.option-text {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.option-votes {
  font-size: 13px;
  color: #6b7280;
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
