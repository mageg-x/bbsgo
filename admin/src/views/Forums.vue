<template>
  <div class="forums-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="total-count">{{ t('forum.totalForums').replace('%d', forums.length) }}</span>
          </div>
          <el-button type="primary" @click="openCreateModal">
            <Plus :size="16" />
            {{ t('forum.create') }}
          </el-button>
        </div>
      </template>

      <el-table :data="forums" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('forum.name')" min-width="150">
          <template #default="{ row }">
            <div class="forum-name">
              <div class="forum-icon">
                <FolderOpen :size="16" />
              </div>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('forum.description')" min-width="200">
          <template #default="{ row }">
            <span class="desc-text">{{ row.description || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" :label="t('forum.sortOrder')" width="100" sortable />
        <el-table-column :label="t('forum.allowPost')" width="120">
          <template #default="{ row }">
            <el-tag :type="row.allow_post ? 'success' : 'danger'" size="small">
              {{ row.allow_post ? t('common.yes') : t('common.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('forum.actions')" width="140" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editForum(row)">
              <Edit :size="14" />
              {{ t('common.edit') }}
            </el-button>
            <el-button link type="danger" @click="deleteForum(row)">
              <Trash2 :size="14" />
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editingForum ? t('forum.edit') : t('forum.create')" width="480px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="form" label-position="top">
        <el-form-item :label="t('forum.name')" prop="name" :rules="[{ required: true, message: t('forum.nameRequired'), trigger: 'blur' }]">
          <el-input v-model="form.name" :placeholder="t('forum.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('forum.description')" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" :placeholder="t('forum.descPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('forum.sortOrder')" prop="sort_order">
          <el-input-number v-model="form.sort_order" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="form.allow_post">{{ t('forum.allowPost') }}</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="saveForum" :loading="saving">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { FolderOpen, Plus, Edit, Trash2 } from 'lucide-vue-next'

const { t } = useI18n()
const forums = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const saving = ref(false)
const editingForum = ref(null)
const formRef = ref(null)

const form = ref({
  name: '',
  description: '',
  sort_order: 0,
  allow_post: true
})

async function loadForums() {
  loading.value = true
  try {
    const res = await api.get('/forums')
    forums.value = res || []
  } catch (e) {
    console.error('load forums failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingForum.value = null
  form.value = { name: '', description: '', sort_order: 0, allow_post: true }
  dialogVisible.value = true
}

function editForum(forum) {
  editingForum.value = forum
  form.value = { ...forum }
  dialogVisible.value = true
}

async function saveForum() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    saving.value = true
    try {
      if (editingForum.value) {
        await api.put(`/admin/forums/${editingForum.value.id}`, form.value)
        Object.assign(editingForum.value, form.value)
        ElMessage.success(t('common.success'))
      } else {
        const res = await api.post('/admin/forums', form.value)
        forums.value.push(res)
        ElMessage.success(t('common.success'))
      }
      dialogVisible.value = false
    } catch (e) {
      console.error('save forum failed', e)
      ElMessage.error(t('common.failed'))
    } finally {
      saving.value = false
    }
  })
}

async function deleteForum(forum) {
  try {
    await ElMessageBox.confirm(`${t('forum.confirmDelete')} "${forum.name}"?`, t('forum.delete'), {
      confirmButtonText: t('common.delete'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.delete(`/admin/forums/${forum.id}`)
    forums.value = forums.value.filter(f => f.id !== forum.id)
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('delete forum failed', e)
    }
  }
}

onMounted(() => {
  loadForums()
})
</script>

<style scoped>
.forums-page {
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

.forum-name {
  display: flex;
  align-items: center;
  gap: 10px;
}

.forum-icon {
  width: 32px;
  height: 32px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.desc-text {
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
