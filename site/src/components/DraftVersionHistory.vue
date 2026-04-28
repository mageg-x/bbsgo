<template>
  <el-dialog v-model="visible" :title="t('newTopic.versionHistory')" width="700px" class="version-history-dialog">
    <div v-if="loading" class="text-center py-8">
      <i class="el-icon-loading text-2xl text-blue-500"></i>
      <p class="mt-2 text-gray-500">{{ t('newTopic.loadingVersions') }}</p>
    </div>
    <div v-else-if="versions.length === 0" class="text-center py-8">
      <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <p class="text-gray-500">{{ t('newTopic.noVersions') }}</p>
    </div>
    <div v-else class="version-list">
      <div v-for="version in versions" :key="version.id" 
           :class="['version-item', { 'active': selectedVersion?.id === version.id }]"
           @click="selectVersion(version)">
        <div class="version-header">
          <div class="version-info">
            <span class="version-badge">v{{ version.version }}</span>
            <span class="version-time">{{ formatTime(version.created_at) }}</span>
          </div>
          <div class="version-actions" v-if="selectedVersion?.id === version.id">
            <el-button type="primary" size="small" @click.stop="previewVersion(version)">
              {{ t('newTopic.preview') }}
            </el-button>
            <el-button type="success" size="small" @click.stop="restoreVersion(version)">
              {{ t('newTopic.restore') }}
            </el-button>
            <el-button type="danger" size="small" @click.stop="deleteVersion(version)">
              {{ t('newTopic.delete') }}
            </el-button>
          </div>
        </div>
        <div class="version-preview" v-if="version.title">
          <span class="version-title">{{ version.title || t('newTopic.untitled') }}</span>
        </div>
        <div class="version-content-preview">
          {{ version.preview || t('newTopic.noContent') }}
        </div>
      </div>
    </div>

    <!-- 版本预览对话框 -->
    <el-dialog v-model="previewVisible" :title="t('newTopic.versionPreview')" width="800px">
      <div v-if="previewLoading" class="text-center py-8">
        <i class="el-icon-loading text-2xl text-blue-500"></i>
      </div>
      <div v-else class="preview-content">
        <div class="preview-header mb-4 pb-4 border-b">
          <h3 class="text-lg font-bold text-gray-800">{{ previewData?.title || t('newTopic.untitled') }}</h3>
          <p class="text-sm text-gray-500 mt-1">
            {{ t('newTopic.version') }} v{{ previewData?.version }} · 
            {{ formatTime(previewData?.created_at) }}
          </p>
        </div>
        <div class="preview-body markdown-body" v-html="renderMarkdown(previewData?.content || '')">
        </div>
      </div>
      <template #footer>
        <el-button @click="previewVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="success" @click="restoreFromPreview" :loading="restoring">
          {{ t('newTopic.restoreThisVersion') }}
        </el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { draftApi } from '@/api'
import { renderMarkdown as renderMd } from '@/utils/markdown'
import { getErrorI18nKey } from '@/utils/error'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  draftId: {
    type: [Number, String],
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'restore'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const versions = ref([])
const selectedVersion = ref(null)
const previewVisible = ref(false)
const previewLoading = ref(false)
const previewData = ref(null)
const restoring = ref(false)

const loadVersions = async () => {
  if (!props.draftId) return
  
  loading.value = true
  try {
    const res = await draftApi.getVersions(props.draftId)
    versions.value = res || []
  } catch (e) {
    console.error('Failed to load versions:', e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  } finally {
    loading.value = false
  }
}

const selectVersion = (version) => {
  selectedVersion.value = selectedVersion.value?.id === version.id ? null : version
}

const previewVersion = async (version) => {
  previewLoading.value = true
  previewVisible.value = true
  
  try {
    const res = await draftApi.previewVersion(props.draftId, version.id)
    previewData.value = res
  } catch (e) {
    console.error('Failed to preview version:', e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  } finally {
    previewLoading.value = false
  }
}

const restoreVersion = async (version) => {
  try {
    await ElMessageBox.confirm(
      t('newTopic.restoreConfirm'),
      t('newTopic.confirmTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    
    restoring.value = true
    const res = await draftApi.restoreVersion(props.draftId, version.id)
    ElMessage.success(t('newTopic.restoreSuccess'))
    emit('restore', res)
    visible.value = false
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Failed to restore version:', e)
      ElMessage.error(t(getErrorI18nKey(e?.code)))
    }
  } finally {
    restoring.value = false
  }
}

const restoreFromPreview = async () => {
  if (previewData.value) {
    await restoreVersion(previewData.value)
  }
}

const deleteVersion = async (version) => {
  try {
    await ElMessageBox.confirm(
      t('newTopic.deleteVersionConfirm'),
      t('newTopic.confirmTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    
    await draftApi.deleteVersion(props.draftId, version.id)
    ElMessage.success(t('newTopic.deleteSuccess'))
    
    const index = versions.value.findIndex(v => v.id === version.id)
    if (index > -1) {
      versions.value.splice(index, 1)
    }
    if (selectedVersion.value?.id === version.id) {
      selectedVersion.value = null
    }
  } catch (e) {
    if (e !== 'cancel') {
      console.error('Failed to delete version:', e)
      ElMessage.error(t(getErrorI18nKey(e?.code)))
    }
  }
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const renderMarkdown = (content) => {
  if (!content) return ''
  return renderMd(content)
}

watch(visible, (val) => {
  if (val && props.draftId) {
    loadVersions()
    selectedVersion.value = null
  }
})

watch(() => props.draftId, () => {
  if (visible.value) {
    loadVersions()
  }
})
</script>

<style scoped>
.version-history-dialog :deep(.el-dialog__body) {
  padding: 0;
  max-height: 600px;
  overflow-y: auto;
}

.version-list {
  padding: 16px;
}

.version-item {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.version-item:hover {
  border-color: #3b82f6;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.1);
}

.version-item.active {
  border-color: #3b82f6;
  background: #eff6ff;
}

.version-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.version-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.version-badge {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.version-time {
  color: #6b7280;
  font-size: 13px;
}

.version-actions {
  display: flex;
  gap: 8px;
}

.version-preview {
  margin-bottom: 4px;
}

.version-title {
  font-weight: 600;
  color: #1f2937;
  font-size: 14px;
}

.version-content-preview {
  color: #6b7280;
  font-size: 13px;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.preview-content {
  max-height: 500px;
  overflow-y: auto;
}

.preview-body {
  padding: 0;
}

.preview-body :deep(.markdown-body) {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 15px;
  line-height: 1.7;
  color: #1f2937;
}

.preview-body :deep(.markdown-body h1),
.preview-body :deep(.markdown-body h2),
.preview-body :deep(.markdown-body h3),
.preview-body :deep(.markdown-body h4),
.preview-body :deep(.markdown-body h5),
.preview-body :deep(.markdown-body h6) {
  color: #111827;
  font-weight: 700;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
}

.preview-body :deep(.markdown-body p) {
  margin: 0.75rem 0;
}

.preview-body :deep(.markdown-body img),
.preview-body :deep(.markdown-body video) {
  max-width: 100%;
  border-radius: 8px;
  margin: 1rem 0;
}
</style>
