<template>
  <div class="topics-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="total-count">{{ t('topic.totalTopics').replace('%d', total) }}</span>
          </div>
          <div class="header-right">
            <el-input
              v-model="searchKeyword"
              :placeholder="t('topic.searchPlaceholder')"
              clearable
              @clear="loadTopics"
              @keyup.enter="loadTopics"
              style="width: 220px"
            >
              <template #prefix>
                <Search :size="16" />
              </template>
            </el-input>
            <el-button type="primary" @click="loadTopics">
              <Search :size="16" />
              {{ t('common.search') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="topics" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('topic.title')" min-width="200">
          <template #default="{ row }">
            <div class="title-cell">
              <el-link :href="`/topic/${row.id}`" target="_blank" type="primary">
                {{ row.title }}
              </el-link>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('topic.author')" width="120">
          <template #default="{ row }">
            <div class="author-cell">
              <el-avatar :size="24" :src="row.user?.avatar">
                <User :size="12" />
              </el-avatar>
              <span>{{ row.user?.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('topic.forum')" width="120">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.forum?.name || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('topic.stats')" width="150">
          <template #default="{ row }">
            <div class="stats-cell">
              <span class="stat-item">
                <Eye :size="12" />
                {{ row.view_count }}
              </span>
              <span class="stat-item">
                <Heart :size="12" />
                {{ row.like_count }}
              </span>
              <span class="stat-item">
                <MessageCircle :size="12" />
                {{ row.reply_count }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('topic.status')" width="150">
          <template #default="{ row }">
            <div class="status-cell">
              <el-tag v-if="row.is_pinned" type="danger" size="small">{{ t('topic.pinned') }}</el-tag>
              <el-tag v-if="row.is_essence" type="warning" size="small">{{ t('topic.essence') }}</el-tag>
              <el-tag v-if="row.is_locked" type="info" size="small">{{ t('topic.locked') }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('topic.publishTime')" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="180" fixed="right">
          <template #default="{ row }">
            <div class=" inline-flex">
              <el-button link type="primary" @click="togglePin(row)">
                <Pin :size="14" />
                {{ row.is_pinned ? t('topic.unpin') : t('topic.pin') }}
              </el-button>
              <el-button link type="primary" @click="viewTopic(row)">
                <ExternalLink :size="14" />
                {{ t('topic.view') }}
              </el-button>
              <el-button link type="danger" @click="deleteTopic(row)">
                <Trash2 :size="14" />
                {{ t('common.delete') }}
              </el-button>
            </div>

          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadTopics"
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
import { FileText, Search, User, Eye, Heart, MessageCircle, ExternalLink, Trash2, Pin } from 'lucide-vue-next'

const { t } = useI18n()
const topics = ref([])
const searchKeyword = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}

async function loadTopics() {
  loading.value = true
  try {
    const res = await api.get('/admin/topics', {
      params: {
        page: page.value,
        page_size: pageSize,
        keyword: searchKeyword.value
      }
    })
    topics.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('load topics failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

function viewTopic(topic) {
  window.open(`/topic/${topic.id}`, '_blank')
}

async function togglePin(topic) {
  try {
    await ElMessageBox.confirm(
      topic.is_pinned ? t('topic.confirmUnpin') : t('topic.confirmPin'),
      topic.is_pinned ? t('topic.unpin') : t('topic.pin'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )

    await api.put(`/admin/topics/${topic.id}/pin`, {
      pinned: !topic.is_pinned
    })
    topic.is_pinned = !topic.is_pinned
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('toggle pin failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

async function deleteTopic(topic) {
  try {
    await ElMessageBox.confirm(`${t('topic.confirmDelete')} "${topic.title}"?`, t('common.delete'), {
      confirmButtonText: t('common.delete'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.delete(`/admin/topics/${topic.id}`)
    topics.value = topics.value.filter(t => t.id !== topic.id)
    total.value--
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('delete topic failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

onMounted(() => {
  loadTopics()
})
</script>

<style scoped>
.topics-page {
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

.author-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.stats-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
}

.status-cell {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
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
</style>
