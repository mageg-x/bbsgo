<template>
  <div class="max-w-5xl mx-auto px-3 sm:px-4 py-4 sm:py-6">
    <div class="bg-white rounded-lg shadow-sm p-4 sm:p-6 mb-4">
      <div class="flex items-center justify-between">
        <h1 class="text-lg sm:text-xl font-bold text-gray-900">{{ t('drafts.myDrafts') }}</h1>
        <span class="text-sm text-gray-500">{{ t('drafts.draftCount', { count: drafts.length }) }}</span>
      </div>
    </div>

    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
    </div>

    <div v-else-if="drafts.length === 0" class="bg-white rounded-lg shadow-sm p-6 sm:p-8 text-center">
      <svg class="w-12 h-12 sm:w-16 sm:h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor"
        viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
      </svg>
      <p class="text-gray-500 text-sm">{{ t('drafts.noDrafts') }}</p>
      <router-link to="/new-topic" class="inline-block mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors text-sm">
        {{ t('drafts.createDraft') }}
      </router-link>
    </div>

    <div v-else class="space-y-3 sm:space-y-4">
      <div v-for="draft in drafts" :key="draft.id"
        class="bg-white rounded-lg shadow-sm p-3 sm:p-4 hover:shadow-md transition-shadow">
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0 cursor-pointer" @click="openDraft(draft)">
            <h3 class="text-base sm:text-lg font-medium text-gray-900 hover:text-blue-500 transition-colors truncate">
              {{ draft.title || t('drafts.untitledDraft') }}
            </h3>
            <div class="mt-2 text-sm text-gray-600 line-clamp-2" v-html="renderContentPreview(draft.content)">
            </div>
            <div class="flex items-center flex-wrap gap-2 mt-3 text-xs text-gray-500">
              <span>{{ t('drafts.updatedAt', { time: formatTime(draft.updated_at) }) }}</span>
              <span v-if="draft.forum" class="text-xs bg-blue-100 text-blue-600 px-1.5 py-0.5 rounded">
                {{ draft.forum.name }}
              </span>
              <span v-if="draft.tags && draft.tags.length > 0" class="flex items-center gap-1">
                <span v-for="(tag, index) in draft.tags" :key="index"
                  class="text-xs bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded">
                  #{{ tag }}
                </span>
              </span>
            </div>
          </div>
          <div class="flex flex-col gap-2 ml-3 flex-shrink-0">
            <button @click="openVersionHistory(draft)" class="text-blue-500 hover:text-blue-700 p-1"
              :title="t('drafts.versionHistory')">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </button>
            <button @click="deleteDraft(draft)" class="text-red-500 hover:text-red-700 p-1"
              :title="t('common.delete')">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
                </path>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 草稿版本历史对话框 -->
    <DraftVersionHistory
      v-model="showVersionHistory"
      :draft-id="selectedDraftId"
      @restore="handleRestoreVersion"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { draftApi } from '@/api'
import { stripMarkdown } from '@/utils/markdown'
import { getErrorI18nKey } from '@/utils/error'
import DraftVersionHistory from '@/components/DraftVersionHistory.vue'

const { t } = useI18n()
const router = useRouter()
const loading = ref(false)
const drafts = ref([])
const showVersionHistory = ref(false)
const selectedDraftId = ref(null)

function formatTime(timeStr) {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (seconds < 60) return t('notifications.justNow')
  if (minutes < 60) return t('notifications.minutesAgo', { 0: minutes })
  if (hours < 24) return t('notifications.hoursAgo', { 0: hours })
  if (days < 30) return t('notifications.daysAgo', { 0: days })

  return date.toLocaleDateString('zh-CN')
}

function renderContentPreview(content) {
  if (!content) return t('drafts.noContent')
  const preview = stripMarkdown(content)
  return preview.length > 100 ? preview.substring(0, 100) + '...' : preview
}

async function loadDrafts() {
  loading.value = true
  try {
    const res = await draftApi.getDrafts()
    drafts.value = res || []
  } catch (e) {
    console.error('Failed to load drafts:', e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  } finally {
    loading.value = false
  }
}

function openDraft(draft) {
  router.push(`/new-topic?draft_id=${draft.id}`)
}

function openVersionHistory(draft) {
  selectedDraftId.value = draft.id
  showVersionHistory.value = true
}

function handleRestoreVersion(version) {
  showVersionHistory.value = false
  ElMessage.success(t('newTopic.versionRestored'))
  loadDrafts()
}

async function deleteDraft(draft) {
  try {
    await ElMessageBox.confirm(
      t('drafts.deleteConfirm'),
      t('newTopic.confirmTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
  } catch {
    return
  }

  try {
    await draftApi.deleteDraft(draft.id)
    drafts.value = drafts.value.filter(d => d.id !== draft.id)
    ElMessage.success(t('drafts.deleteSuccess'))
  } catch (e) {
    console.error('Failed to delete draft:', e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  }
}

onMounted(() => {
  loadDrafts()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
