<template>
  <div class="comments-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="total-count">{{ t('comment.totalComments').replace('%d', total) }}</span>
          </div>
          <div class="header-right">
            <el-input
              v-model="searchKeyword"
              :placeholder="t('comment.searchPlaceholder')"
              clearable
              @clear="loadComments"
              @keyup.enter="loadComments"
              style="width: 220px"
            >
              <template #prefix>
                <Search :size="16" />
              </template>
            </el-input>
            <el-button type="primary" @click="loadComments">
              <Search :size="16" />
              {{ t('common.search') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="comments" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('comment.content')" min-width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.content" placement="top" :disabled="row.content.length < 60">
              <span class="content-text">{{ row.content }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column :label="t('comment.author')" width="120">
          <template #default="{ row }">
            <div class="author-cell">
              <el-avatar :size="24" :src="row.user?.avatar">
                <User :size="12" />
              </el-avatar>
              <span>{{ row.user?.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('comment.topic')" width="120">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewTopic(row.topic_id)">
              <ExternalLink :size="12" />
              {{ t('common.view') }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column :label="t('comment.like')" width="80">
          <template #default="{ row }">
            <span class="like-count">
              <Heart :size="12" />
              {{ row.like_count }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="t('comment.publishTime')" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="danger" @click="deleteComment(row)">
              <Trash2 :size="14" />
              {{ t('common.delete') }}
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
          @current-change="loadComments"
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
import { MessageCircle, Search, User, ExternalLink, Heart, Trash2 } from 'lucide-vue-next'

const { t } = useI18n()
const comments = ref([])
const searchKeyword = ref('')
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}

async function loadComments() {
  loading.value = true
  try {
    const res = await api.get('/admin/comments', {
      params: {
        page: page.value,
        page_size: pageSize,
        keyword: searchKeyword.value
      }
    })
    comments.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    console.error('load comments failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

function viewTopic(topicId) {
  window.open(`/topic/${topicId}`, '_blank')
}

async function deleteComment(comment) {
  try {
    await ElMessageBox.confirm(t('comment.confirmDelete'), t('comment.delete'), {
      confirmButtonText: t('common.delete'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.delete(`/admin/comments/${comment.id}`)
    comments.value = comments.value.filter(c => c.id !== comment.id)
    total.value--
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('delete comment failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

onMounted(() => {
  loadComments()
})
</script>

<style scoped>
.comments-page {
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

.content-text {
  font-size: 13px;
  color: #374151;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.author-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.like-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #f472b6;
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
