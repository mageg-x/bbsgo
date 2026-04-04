<template>
  <div class="max-w-6xl mx-auto px-4 py-6">
    <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">{{ t('followTopics.followingUpdates') }}</h1>
      <p class="text-gray-500">{{ t('followTopics.followingUpdatesDesc') }}</p>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <el-icon class="is-loading" :size="32">
        <Loading />
      </el-icon>
    </div>

    <div v-else-if="topics.length === 0" class="bg-white rounded-lg shadow-sm p-12 text-center">
      <p class="text-gray-500 mb-4">{{ t('followTopics.noFollowingOrNoContent') }}</p>
      <router-link to="/" class="text-blue-500 hover:underline">{{ t('followTopics.goDiscover') }}</router-link>
    </div>

    <div v-else class="space-y-4">
      <TopicCard v-for="topic in topics" :key="topic.id" :topic="topic" />
    </div>

    <div v-if="total > pageSize" class="flex justify-center mt-6">
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import api from '@/api'
import { getErrorI18nKey } from '@/utils/error'
import TopicCard from '@/components/TopicCard.vue'

const { t } = useI18n()
const topics = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

async function loadFollowTopics() {
  loading.value = true
  try {
    const res = await api.get('/user/follow-topics', {
      params: {
        page: page.value
      }
    })
    topics.value = res.list || []
    total.value = res.total || 0
  } catch (e) {
    console.error(t('followTopics.loadFailed'), e)
    ElMessage.error(t(getErrorI18nKey(e?.code)))
  } finally {
    loading.value = false
  }
}

function handlePageChange(newPage) {
  page.value = newPage
  loadFollowTopics()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  loadFollowTopics()
})
</script>
