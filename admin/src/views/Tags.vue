<template>
  <div class="tags-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
          </div>
          <el-button type="primary" @click="openAddModal">
            <Plus :size="16" />
            {{ t('tag.addOfficialTag') }}
          </el-button>
        </div>
      </template>

      <div class="info-tip">
        <Info :size="16" />
        <span>{{ t('tag.infoTip') }}</span>
      </div>

      <el-table :data="tags" stripe style="width: 100%" v-loading="loading" :row-class-name="getRowClass">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('tag.name')" min-width="180">
          <template #default="{ row }">
            <div class="tag-name">
              <span v-if="row.icon" class="tag-icon">{{ row.icon }}</span>
              <span :class="{ 'line-through text-gray-400': row.is_banned }">{{ row.name }}</span>
              <el-tag v-if="row.is_official" type="primary" size="small">{{ t('tag.isOfficial') }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="usage_count" :label="t('tag.usageCount')" width="120" sortable />
        <el-table-column :label="t('user.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_banned ? 'danger' : 'success'" size="small">
              {{ row.is_banned ? t('tag.disabled') : t('tag.normal') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('comment.publishTime')" width="120">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" min-width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editTag(row)">
              <Edit :size="14" />
              {{ t('common.edit') }}
            </el-button>
            <el-button link :type="row.is_official ? 'info' : 'success'" @click="toggleOfficial(row)">
              <Star :size="14" />
              {{ row.is_official ? t('tag.cancelOfficial') : t('tag.setOfficial') }}
            </el-button>
            <el-button link :type="row.is_banned ? 'success' : 'warning'" @click="toggleBan(row)">
              <Ban :size="14" />
              {{ row.is_banned ? t('user.unban') : t('tag.disable') }}
            </el-button>
            <el-button link type="info" @click="openMergeModal(row)">
              <GitMerge :size="14" />
              {{ t('tag.merge') }}
            </el-button>
            <el-button link type="danger" @click="deleteTag(row)">
              <Trash2 :size="14" />
              {{ t('tag.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editingTag ? t('tag.editTag') : t('tag.addOfficialTag')" width="420px" :close-on-click-modal="false">
      <el-form ref="formRef" :model="form" label-position="top">
        <el-form-item :label="t('tag.name')" prop="name" :rules="[{ required: true, message: t('tag.nameRequired'), trigger: 'blur' }]">
          <el-input v-model="form.name" :placeholder="t('tag.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('tag.icon') + '（Emoji）'" prop="icon">
          <el-input v-model="form.icon" :placeholder="t('tag.iconPlaceholder')">
            <template #append>
              <span class="emoji-preview">{{ form.icon || '👁' }}</span>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item v-if="!editingTag">
          <el-checkbox v-model="form.is_official">{{ t('tag.setOfficialRecommend') }}</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('tag.cancel') }}</el-button>
        <el-button type="primary" @click="saveTag" :loading="saving">{{ editingTag ? t('tag.update') : t('tag.create') }}</el-button>
      </template>
    </el-dialog>

    <!-- 合并弹窗 -->
    <el-dialog v-model="mergeDialogVisible" :title="t('tag.mergeTag')" width="420px">
      <div class="merge-tip">
        <AlertCircle :size="16" />
        <span>{{ t('tag.mergeTip').replace('%s', mergeSource?.name) }}</span>
      </div>
      <el-select v-model="mergeTargetId" :placeholder="t('tag.selectTargetTag')" style="width: 100%; margin-top: 12px">
        <el-option v-for="tag in availableTags" :key="tag.id" :label="`${tag.name} (${t('tag.usageTimes').replace('%d', tag.usage_count)})`" :value="tag.id" />
      </el-select>
      <template #footer>
        <el-button @click="mergeDialogVisible = false">{{ t('tag.cancel') }}</el-button>
        <el-button type="primary" @click="mergeTags" :disabled="!mergeTargetId" :loading="saving">{{ t('tag.confirmMerge') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { Tag, Plus, Edit, Star, Ban, GitMerge, Trash2, Info, AlertCircle } from 'lucide-vue-next'

const { t } = useI18n()
const tags = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const mergeDialogVisible = ref(false)
const saving = ref(false)
const editingTag = ref(null)
const mergeSource = ref(null)
const mergeTargetId = ref('')
const formRef = ref(null)

const form = ref({
  name: '',
  icon: '',
  is_official: true
})

const availableTags = computed(() => tags.value.filter(t => t.id !== mergeSource.value?.id))

function getRowClass({ row }) {
  return row.is_banned ? 'banned-row' : ''
}

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}

async function loadTags() {
  loading.value = true
  try {
    const res = await api.get('/admin/tags')
    tags.value = res || []
  } catch (e) {
    console.error('load tags failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  editingTag.value = null
  form.value = { name: '', icon: '', is_official: true }
  dialogVisible.value = true
}

function editTag(tag) {
  editingTag.value = tag
  form.value = { name: tag.name, icon: tag.icon || '', is_official: tag.is_official }
  dialogVisible.value = true
}

async function saveTag() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    saving.value = true
    try {
      if (editingTag.value) {
        await api.put(`/admin/tags/${editingTag.value.id}`, form.value)
        Object.assign(editingTag.value, form.value)
        ElMessage.success(t('common.success'))
      } else {
        await api.post('/admin/tags', form.value)
        ElMessage.success(t('common.success'))
        loadTags()
      }
      dialogVisible.value = false
    } catch (e) {
      console.error('save tag failed', e)
      ElMessage.error(t('common.failed'))
    } finally {
      saving.value = false
    }
  })
}

async function toggleOfficial(tag) {
  try {
    await api.put(`/admin/tags/${tag.id}`, { ...tag, is_official: !tag.is_official })
    tag.is_official = !tag.is_official
    ElMessage.success(tag.is_official ? t('tag.setOfficial') : t('tag.cancelOfficial'))
  } catch (e) {
    console.error('toggle official failed', e)
    ElMessage.error(t('common.failed'))
  }
}

async function toggleBan(tag) {
  const action = tag.is_banned ? t('tag.enable') : t('tag.disable')
  try {
    await ElMessageBox.confirm(`${t('common.confirm')} ${action} tag "${tag.name}"?`, action, {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.put(`/admin/tags/${tag.id}`, { ...tag, is_banned: !tag.is_banned })
    tag.is_banned = !tag.is_banned
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('toggle ban failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

async function deleteTag(tag) {
  try {
    await ElMessageBox.confirm(t('tag.cannotRecover'), t('tag.delete'), {
      confirmButtonText: t('common.delete'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.delete(`/admin/tags/${tag.id}`)
    tags.value = tags.value.filter(t => t.id !== tag.id)
    ElMessage.success(t('common.success'))
  } catch (e) {
    if (e !== 'cancel') {
      console.error('delete tag failed', e)
      ElMessage.error(t('common.failed'))
    }
  }
}

function openMergeModal(tag) {
  mergeSource.value = tag
  mergeTargetId.value = ''
  mergeDialogVisible.value = true
}

async function mergeTags() {
  if (!mergeTargetId.value) return

  saving.value = true
  try {
    await api.post('/admin/tags/merge', {
      source_id: mergeSource.value.id,
      target_id: parseInt(mergeTargetId.value)
    })
    ElMessage.success(t('common.success'))
    mergeDialogVisible.value = false
    loadTags()
  } catch (e) {
    console.error('merge tags failed', e)
    ElMessage.error(t('common.failed'))
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.tags-page {
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

.header-left h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.info-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(102, 126, 234, 0.08);
  color: #667eea;
  border-radius: 10px;
  margin-bottom: 16px;
  font-size: 13px;
}

.tag-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-icon {
  font-size: 18px;
}

.date-text {
  font-size: 13px;
  color: #6b7280;
}

.merge-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 14px;
}

.emoji-preview {
  font-size: 16px;
}

:deep(.el-table) {
  --el-table-border-color: #f3f4f6;
  --el-table-header-bg-color: #f9fafb;
}

:deep(.el-table th) {
  font-weight: 600;
  color: #374151;
}

:deep(.el-table .banned-row) {
  background-color: rgba(254, 226, 226, 0.3);
}

:deep(.el-button.is-link) {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
}

:deep(.el-button--purple) {
  color: #8b5cf6;
}
</style>
