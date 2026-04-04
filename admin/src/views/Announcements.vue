<template>
  <div class="announcements-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="total-count">{{ t('announcement.totalAnnouncements').replace('%d', announcements.length) }}</span>
          </div>
          <el-button type="primary" @click="openCreateModal">
            <Plus :size="16" />
            {{ t('announcement.publishAnnouncement') }}
          </el-button>
        </div>
      </template>

      <el-table :data="announcements" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('announcement.title2')" min-width="200">
          <template #default="{ row }">
            <div class="title-cell">
              <span v-if="row.is_pinned" class="pin-icon">
                <Pin :size="14" />
              </span>
              <span class="title-text">{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('announcement.isPinned')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_pinned ? 'danger' : 'info'" size="small">
              {{ row.is_pinned ? t('common.yes') : t('common.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('announcement.expiresAt')" width="150">
          <template #default="{ row }">
            <span class="date-text">{{ row.expires_at ? formatDate(row.expires_at) : t('announcement.permanent') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('announcement.publishTime')" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="140" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editAnnouncement(row)">
              <Edit :size="14" />
              {{ t('common.edit') }}
            </el-button>
            <el-button link type="danger" @click="deleteAnnouncement(row)">
              <Trash2 :size="14" />
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editingAnnouncement ? t('announcement.editAnnouncement') : t('announcement.publishAnnouncement')" width="520px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="form" label-position="top">
        <el-form-item :label="t('announcement.title2')" prop="title" :rules="[{ required: true, message: t('announcement.titlePlaceholder'), trigger: 'blur' }]">
          <el-input v-model="form.title" :placeholder="t('announcement.titlePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('announcement.content')" prop="content" :rules="[{ required: true, message: t('announcement.contentPlaceholder'), trigger: 'blur' }]">
          <el-input v-model="form.content" type="textarea" :rows="5" :placeholder="t('announcement.contentPlaceholder')" />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="form.is_pinned">{{ t('announcement.isPinned') }}</el-checkbox>
        </el-form-item>
        <el-form-item :label="t('announcement.expiresAtOptional')">
          <el-date-picker v-model="form.expires_at" type="datetime" :placeholder="t('announcement.expiresAtPlaceholder')" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('announcement.cancel') }}</el-button>
        <el-button type="primary" @click="saveAnnouncement" :loading="saving">{{ t('announcement.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Bell, Plus, Edit, Trash2, Pin } from 'lucide-vue-next'

const { t } = useI18n()
const announcements = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const saving = ref(false)
const editingAnnouncement = ref(null)
const formRef = ref(null)

const form = ref({
  title: '',
  content: '',
  is_pinned: false,
  expires_at: ''
})

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}

async function loadAnnouncements() {
  loading.value = true
  try {
    const res = await api.get('/announcements')
    announcements.value = res || []
  } catch (e) {
    console.error('Load announcements failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingAnnouncement.value = null
  form.value = { title: '', content: '', is_pinned: false, expires_at: '' }
  dialogVisible.value = true
}

function editAnnouncement(announcement) {
  editingAnnouncement.value = announcement
  form.value = { ...announcement }
  dialogVisible.value = true
}

async function saveAnnouncement() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    saving.value = true
    try {
      if (editingAnnouncement.value) {
        await api.put(`/admin/announcements/${editingAnnouncement.value.id}`, form.value)
        Object.assign(editingAnnouncement.value, form.value)
        ElMessage.success(t('common.success'))
      } else {
        const res = await api.post('/admin/announcements', form.value)
        announcements.value.unshift(res)
        ElMessage.success(t('common.success'))
      }
      dialogVisible.value = false
    } catch (e) {
      console.error('Save announcement failed', e)
      ElMessage.error(t('common.failed'))
    } finally {
      saving.value = false
    }
  })
}

async function deleteAnnouncement(announcement) {
  try {
    await ElMessageBox.confirm(`${t('common.confirm')} "${announcement.title}"?`, t('announcement.delete'), {
      confirmButtonText: t('common.delete'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.delete(`/admin/announcements/${announcement.id}`)
    announcements.value = announcements.value.filter(a => a.id !== announcement.id)
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Delete announcement failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

onMounted(() => {
  loadAnnouncements()
})
</script>

<style scoped>
.announcements-page {
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

.title-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pin-icon {
  color: #f87171;
}

.title-text {
  font-weight: 500;
  color: #1f2937;
}

.date-text {
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
