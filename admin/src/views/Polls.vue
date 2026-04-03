<template>
  <div class="polls-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>
              <Vote :size="18" />
              投票管理
            </h3>
            <span class="total-count">共 {{ total }} 个投票</span>
          </div>
          <div class="header-right">
            <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="loadPolls" style="width: 120px">
              <el-option label="进行中" value="active" />
              <el-option label="已结束" value="ended" />
            </el-select>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索投票标题"
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
              搜索
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="polls" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="投票标题" min-width="200">
          <template #default="{ row }">
            <div class="title-cell">
              <el-link :href="`/topic/${row.topic?.id}`" target="_blank" type="primary">
                {{ row.poll?.title || '(无标题)' }}
              </el-link>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="row.poll?.poll_type === 'single' ? 'primary' : 'success'" size="small">
              {{ row.poll?.poll_type === 'single' ? '单选' : '多选' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="选项数" width="80">
          <template #default="{ row }">
            <span>{{ row.poll?.options?.length || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="投票人数" width="100">
          <template #default="{ row }">
            <span class="vote-count">{{ row.poll?.total_votes || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.poll?.is_ended" type="info" size="small">已结束</el-tag>
            <el-tag v-else-if="isEnded(row.poll?.end_time)" type="info" size="small">已结束</el-tag>
            <el-tag v-else type="success" size="small">进行中</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="截止时间" width="160">
          <template #default="{ row }">
            <span class="date-text">{{ row.poll?.end_time ? formatDateTime(row.poll.end_time) : '永久有效' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            <span class="date-text">{{ formatDateTime(row.poll?.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewPoll(row)">
              <ExternalLink :size="14" />
              查看
            </el-button>
            <el-button link type="warning" @click="endPoll(row)" v-if="!row.poll?.is_ended && !isEnded(row.poll?.end_time)">
              <Square :size="14" />
              结束
            </el-button>
            <el-button link type="danger" @click="deletePoll(row)">
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
          @current-change="loadPolls"
        />
      </div>
    </el-card>

    <el-dialog v-model="pollDialogVisible" title="投票详情" width="600px">
      <div v-if="currentPoll" class="poll-detail">
        <div class="poll-info">
          <h4>{{ currentPoll.title || '(无标题)' }}</h4>
          <p class="poll-meta">
            <el-tag :type="currentPoll.poll_type === 'single' ? 'primary' : 'success'" size="small">
              {{ currentPoll.poll_type === 'single' ? '单选' : `多选(最多${currentPoll.max_choices}项)` }}
            </el-tag>
            <span class="vote-total">共 {{ currentPoll.total_votes || 0 }} 人参与</span>
          </p>
        </div>
        
        <div class="options-list">
          <div v-for="option in currentPoll.options" :key="option.id" class="option-item">
            <div class="option-header">
              <span class="option-text">{{ option.text }}</span>
              <span class="option-votes">{{ option.vote_count }} 票 ({{ getPercentage(option.vote_count) }}%)</span>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Vote, Search, ExternalLink, Trash2, Square } from 'lucide-vue-next'

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
  return new Date(date).toLocaleString('zh-CN')
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
    console.error('加载投票失败', e)
    ElMessage.error('加载投票失败')
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
    await ElMessageBox.confirm(`确定要结束投票 "${row.poll?.title || '(无标题)'}" 吗？`, '结束投票', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await api.post(`/admin/polls/${row.poll.id}/end`)
    row.poll.is_ended = true
    ElMessage.success('投票已结束')
    loadPolls()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('操作失败', e)
      ElMessage.error('操作失败')
    }
  }
}

async function deletePoll(row) {
  try {
    await ElMessageBox.confirm(`确定要删除投票 "${row.poll?.title || '(无标题)'}" 吗？`, '删除投票', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await api.delete(`/admin/polls/${row.poll.id}`)
    polls.value = polls.value.filter(p => p.poll.id !== row.poll.id)
    total.value--
    ElMessage.success('投票已删除')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('删除失败', e)
      ElMessage.error('删除失败')
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
