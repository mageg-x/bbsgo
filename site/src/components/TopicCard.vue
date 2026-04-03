<template>
  <div class="topic-card">
    <!-- 标题行：置顶/精华/锁定标签 + 标题 -->
    <h3 class="text-lg font-semibold mb-2 hover:text-blue-500 line-clamp-2 flex items-center flex-wrap gap-1">
      <!-- 管理员置顶 -->
      <span v-if="topic.is_pinned" class="text-xs bg-red-500 text-white px-1.5 py-0.5 rounded font-medium">置顶</span>
      <!-- 精华 -->
      <span v-if="topic.is_essence" class="text-xs bg-yellow-500 text-white px-1.5 py-0.5 rounded font-medium">精华</span>
      <!-- 锁定 -->
      <span v-if="topic.is_locked" class="text-xs bg-gray-500 text-white px-1.5 py-0.5 rounded font-medium">锁定</span>
      <!-- 作者置顶（个人主页用） -->
      <span v-if="topic.is_user_pinned" class="text-xs bg-red-500 text-white px-1.5 py-0.5 rounded font-medium">置顶</span>
      <span class="text-gray-900">{{ topic.title }}</span>
    </h3>

    <!-- 图片预览 -->
    <div v-if="firstImage" class="mb-3">
      <img :src="firstImage" class="w-full max-h-64 object-cover rounded-lg" loading="lazy">
    </div>

    <!-- 视频/投票标识 -->
    <div v-if="hasVideoFlag || hasPoll" class="mb-3 flex items-center gap-2">
      <span v-if="hasVideoFlag" class="inline-flex items-center gap-1 px-2 py-1 bg-purple-100 text-purple-600 rounded text-xs">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        视频
      </span>
      <span v-if="hasPoll" class="inline-flex items-center gap-1 px-2 py-1 bg-green-100 text-green-600 rounded text-xs">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
        </svg>
        投票
      </span>
    </div>

    <!-- 文字预览 -->
    <p v-if="!firstImage && !hasVideoFlag && !hasPoll" class="text-gray-600 text-sm mb-3 line-clamp-3">
      {{ plainContent }}
    </p>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { stripMarkdown, extractFirstImage, hasVideo } from '@/utils/markdown'

const props = defineProps({
  topic: {
    type: Object,
    required: true
  },
  // 支持外部传入 has_poll（用于嵌套结构如收藏页面的 favorite.topic）
  externalHasPoll: {
    type: Boolean,
    default: null
  }
})

const firstImage = computed(() => extractFirstImage(props.topic.content))
const hasVideoFlag = computed(() => hasVideo(props.topic.content))
const hasPoll = computed(() => {
  if (props.externalHasPoll !== null) {
    return props.externalHasPoll
  }
  return props.topic.has_poll || false
})
const plainContent = computed(() => {
  const content = props.topic.content || ''
  return stripMarkdown(content).substring(0, 200)
})
</script>
