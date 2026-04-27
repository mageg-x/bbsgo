<template>
  <div :class="['comment-item', level > 0 ? 'ml-4 sm:ml-8 border-l-2 border-gray-200 pl-3 sm:pl-4' : '']">
    <div :class="[
      'flex space-x-3 sm:space-x-4 p-3 sm:p-4 rounded-lg transition-all',
      comment.is_best ? 'bg-gradient-to-r from-yellow-50 to-orange-50 border-2 border-yellow-300 shadow-md' : 'bg-gray-50'
    ]" :id="'post-' + comment.id">
      
      <div class="flex-1">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-1 gap-2">
          <div class="flex items-center flex-wrap gap-1 sm:gap-2">
            <img :src="getUserAvatar(comment.user)" class="w-8 h-8 sm:w-10 sm:h-10 rounded-full flex-shrink-0">
            <span v-if="comment.is_best" class="px-2 py-0.5 text-xs font-bold bg-yellow-500 text-white rounded-full animate-pulse">{{ t('topic.best') }}</span>
            <span v-if="comment.is_pinned && level === 0" class="text-xs text-red-500 font-medium">{{ t('topic.pin') }}</span>
            <SvgBadge v-if="comment.is_best" type="gold-comment" :size="16" :title="t('topic.bestComment')" />
            <span class="font-medium text-gray-900 text-sm sm:text-base">{{ getUserDisplayName(comment.user) }}</span>
            <span v-if="comment.reply_user" class="text-gray-400 text-xs sm:text-sm">{{ t('topic.reply') }} @{{ comment.reply_user.nickname || comment.reply_user.username }}</span>
            <span v-if="comment.user_id === topicAuthorId" class="px-1.5 py-0.5 text-xs bg-red-500 text-white rounded">{{ t('topic.author') }}</span>
            <div v-if="getCommentAuthorTopBadge(comment)" class="flex items-center gap-0.5">
              <SvgBadge :type="getCommentAuthorTopBadge(comment).icon" :size="14" :title="getCommentAuthorTopBadge(comment).name" />
            </div>
            <span class="text-xs sm:text-sm text-gray-500">{{ formatTime(comment.created_at) }}</span>
          </div>
          <div class="flex flex-wrap gap-1 sm:gap-2">
            <button v-if="canBestComment(comment)" @click="$emit('toggle-best', comment)"
              :class="['text-xs transition-colors', comment.is_best ? 'text-yellow-500 hover:text-yellow-600' : 'text-gray-400 hover:text-yellow-500']">
              {{ comment.is_best ? t('topic.cancelBest') : t('topic.markBest') }}
            </button>
            <button v-if="canPinComment(comment) && level === 0" @click="$emit('toggle-pin', comment)"
              :class="['text-xs transition-colors', comment.is_pinned ? 'text-red-500 hover:text-red-600' : 'text-gray-400 hover:text-red-500']">
              {{ comment.is_pinned ? t('topic.cancelPin') : t('topic.setPin') }}
            </button>
            <button v-if="canDeletePost(comment)" @click="$emit('delete', comment)"
              class="text-xs text-gray-400 hover:text-red-500 transition-colors">{{ t('common.delete') }}</button>
            <button v-if="canReportPost(comment)"
              @click="$emit('report', comment)"
              class="text-xs text-gray-400 hover:text-red-500 transition-colors">{{ t('topic.report') }}</button>
          </div>
        </div>

        <p class="text-gray-700 text-sm sm:text-base mt-2">{{ comment.content }}</p>
        <div class="flex items-center gap-3 sm:space-x-4 mt-2 text-sm">
          <button @click="$emit('toggle-like', comment)"
            :class="['transition-colors', getPostLiked(comment.id) ? 'text-red-500' : 'text-gray-500 hover:text-red-500']">
            {{ getPostLiked(comment.id) ? '❤️' : '🤍' }} {{ comment.like_count }}
          </button>
          <button @click="$emit('reply', comment)" class="text-gray-500 hover:text-blue-500">{{ t('topic.reply') }}</button>
        </div>
      </div>
    </div>

    <div v-if="replies.length > 0" class="mt-2 space-y-2">
      <div v-for="reply in visibleReplies" :key="reply.id">
        <CommentItem 
          :comment="reply" 
          :level="level + 1"
          :topic-author-id="topicAuthorId"
          @toggle-best="$emit('toggle-best', $event)"
          @toggle-pin="$emit('toggle-pin', $event)"
          @delete="$emit('delete', $event)"
          @report="$emit('report', $event)"
          @toggle-like="$emit('toggle-like', $event)"
          @reply="$emit('reply', $event)"
          :get-user-avatar="getUserAvatar"
          :get-user-display-name="getUserDisplayName"
          :get-comment-author-top-badge="getCommentAuthorTopBadge"
          :can-best-comment="canBestComment"
          :can-pin-comment="canPinComment"
          :can-delete-post="canDeletePost"
          :can-report-post="canReportPost"
          :get-post-liked="getPostLiked"
          :format-time="formatTime"
        />
      </div>
      
      <div v-if="hiddenCount > 0" class="text-center py-2">
        <button 
          v-if="!isFullyExpanded"
          @click="expandMore" 
          class="text-sm text-blue-500 hover:text-blue-600 transition-colors"
        >
          {{ t('topic.expandMore', { count: Math.min(5, hiddenCount) }) }}
        </button>
        <button 
          v-else-if="replies.length > 5"
          @click="collapse" 
          class="text-sm text-gray-500 hover:text-gray-600 transition-colors"
        >
          {{ t('topic.collapse') }}
        </button>
        <span v-else-if="hiddenCount > 0" class="text-xs text-gray-400">
          {{ t('topic.moreReplies', { count: hiddenCount }) }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import SvgBadge from './SvgBadge.vue'

const { t } = useI18n()

const props = defineProps({
  comment: {
    type: Object,
    required: true
  },
  level: {
    type: Number,
    default: 0
  },
  topicAuthorId: {
    type: [Number, String],
    default: null
  },
  getUserAvatar: {
    type: Function,
    required: true
  },
  getUserDisplayName: {
    type: Function,
    required: true
  },
  getCommentAuthorTopBadge: {
    type: Function,
    required: true
  },
  canBestComment: {
    type: Function,
    required: true
  },
  canPinComment: {
    type: Function,
    required: true
  },
  canDeletePost: {
    type: Function,
    required: true
  },
  canReportPost: {
    type: Function,
    required: true
  },
  getPostLiked: {
    type: Function,
    required: true
  },
  formatTime: {
    type: Function,
    required: true
  }
})

defineEmits(['toggle-best', 'toggle-pin', 'delete', 'report', 'toggle-like', 'reply'])

const expandedCount = ref(5)

const replies = computed(() => {
  return props.comment.replies || []
})

const visibleReplies = computed(() => {
  return replies.value.slice(0, expandedCount.value)
})

const hiddenCount = computed(() => {
  return replies.value.length - expandedCount.value
})

const isFullyExpanded = computed(() => {
  return expandedCount.value >= replies.value.length
})

function expandMore() {
  expandedCount.value += 5
}

function collapse() {
  expandedCount.value = 5
}
</script>

<style scoped>
.comment-item {
  position: relative;
}
</style>
