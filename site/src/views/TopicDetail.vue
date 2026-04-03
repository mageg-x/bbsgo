<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 bg-white rounded-lg shadow-sm">
    <div v-if="topic" class="mb-6">
      <div class="flex items-start justify-between mb-4">
        <h1 class="text-2xl font-bold text-gray-900">{{ topic.title }}</h1>
        <button v-if="canDeleteTopic" @click="handleDeleteTopic"
          class="flex items-center space-x-1 px-3 py-1.5 text-sm text-red-600 hover:text-white hover:bg-red-500 border border-red-300 rounded-lg transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
            </path>
          </svg>
          <span>删除</span>
        </button>
      </div>
      <div v-if="topic.tags && topic.tags.length > 0" class="flex items-center flex-wrap gap-2 mb-4">
        <router-link v-for="tag in topic.tags" :key="tag.id" :to="`/?tag=${tag.id}`"
          class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded-full hover:bg-blue-200">
          #{{ tag.name }}
        </router-link>
      </div>
      <div class="flex items-center space-x-4 mb-6 pb-6 border-b">
        <router-link :to="`/user/${topic.user_id}`">
          <img :src="getUserAvatar(topic.user)" class="w-12 h-12 rounded-full">
        </router-link>
        <div>
          <router-link :to="`/user/${topic.user_id}`" class="font-medium text-gray-900 hover:text-blue-500">{{
            getUserDisplayName(topic.user) }}</router-link>
          <div class="text-sm text-gray-500">{{ formatTime(topic.created_at) }} · {{ topic.view_count }} 浏览</div>
        </div>
      </div>
      <div class="prose max-w-none mb-6 topic-content" v-html="topic.content"></div>
      
      <div v-if="poll && configStore.state.allow_poll" class="mb-6 p-6 bg-gradient-to-r from-blue-50 to-purple-50 rounded-xl border border-blue-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-bold text-gray-900">{{ poll.title || topic.title }}</h3>
          <span v-if="isPollEnded" class="px-3 py-1 text-xs font-medium bg-gray-200 text-gray-600 rounded-full">已结束</span>
          <span v-else-if="poll.end_time" class="text-sm text-gray-500">
            剩余 {{ getRemainingTime(poll.end_time) }}
          </span>
        </div>
        
        <div v-if="!hasVoted && !isPollEnded" class="space-y-3">
          <div v-for="option in poll.options" :key="option.id" 
            @click="poll.poll_type === 'single' && selectOption(option.id)"
            :class="['p-4 rounded-lg border-2 cursor-pointer transition-all', 
              selectedOptions.includes(option.id) ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-blue-300 bg-white']">
            <div class="flex items-center">
              <input v-if="poll.poll_type === 'single'" type="radio" :checked="selectedOptions.includes(option.id)" 
                class="w-4 h-4 text-blue-600" @click.stop>
              <input v-else type="checkbox" :checked="selectedOptions.includes(option.id)" 
                @change="toggleOption(option.id)" @click.stop class="w-4 h-4 text-blue-600 rounded">
              <span class="ml-3 text-gray-700">{{ option.text }}</span>
            </div>
          </div>
          
          <div class="flex items-center justify-between pt-4">
            <span class="text-sm text-gray-500">
              {{ poll.poll_type === 'single' ? '单选' : `多选，最多选${poll.max_choices}项` }}
            </span>
            <button @click="submitVote" :disabled="selectedOptions.length === 0 || submittingVote"
              class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
              {{ submittingVote ? '提交中...' : '提交投票' }}
            </button>
          </div>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="option in poll.options" :key="option.id" class="relative">
            <div :class="['p-4 rounded-lg border-2 transition-all', 
              votedOptionIds.includes(option.id) ? 'border-blue-500 bg-blue-50' : 'border-gray-200 bg-white']">
              <div class="flex items-center justify-between mb-2">
                <span class="text-gray-700 font-medium">{{ option.text }}</span>
                <span class="text-sm font-medium text-gray-600">
                  {{ getPercentage(option.vote_count) }}%
                  <span v-if="votedOptionIds.includes(option.id)" class="ml-2 text-blue-500">✓ 已选</span>
                </span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                <div class="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full transition-all duration-500"
                  :style="{ width: getPercentage(option.vote_count) + '%' }"></div>
              </div>
              <div class="text-xs text-gray-500 mt-1">{{ option.vote_count }} 票</div>
            </div>
          </div>
          
          <div class="text-center text-sm text-gray-500 pt-2">
            共 {{ poll.total_votes }} 票
            <span v-if="hasVoted" class="ml-2 text-blue-500">· 您已投票</span>
          </div>
        </div>
      </div>
      
      <div class="flex items-center space-x-4 pt-4 border-t">
        <button @click="toggleLike"
          :class="['flex items-center space-x-2 transition-colors', liked ? 'text-red-500' : 'text-gray-500 hover:text-red-500']">
          <svg class="w-5 h-5" :fill="liked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
            </path>
          </svg>
          <span>{{ topic.like_count }}</span>
        </button>
        <button @click="toggleFavorite"
          :class="['flex items-center space-x-2 transition-colors', favorited ? 'text-yellow-500' : 'text-gray-500 hover:text-yellow-500']">
          <svg class="w-5 h-5" :fill="favorited ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
          </svg>
          <span>{{ favorited ? '已收藏' : '收藏' }}</span>
        </button>
        <button @click="shareTopic"
          class="flex items-center space-x-2 text-gray-500 hover:text-green-500 transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z">
            </path>
          </svg>
          <span>分享</span>
        </button>
      </div>
    </div>
    <div class="mt-8">
      <h3 class="text-lg font-medium text-gray-900 mb-4">{{ posts.length }} 条评论</h3>
      <div v-if="userStore.isLoggedIn && configStore.state.allow_comment && topic?.allow_comment" class="mb-6">
        <textarea v-model="newPost" rows="3"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
          placeholder="写下你的评论..."></textarea>
        <div class="flex justify-end mt-2">
          <button @click="submitPost"
            class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600">发表评论</button>
        </div>
      </div>
      <div v-else-if="!configStore.state.allow_comment"
        class="mb-6 p-4 bg-gray-100 rounded-lg text-center text-gray-500">
        评论功能已关闭
      </div>
      <div v-else-if="topic && !topic.allow_comment" class="mb-6 p-4 bg-gray-100 rounded-lg text-center text-gray-500">
        本话题已关闭评论
      </div>
      <div class="space-y-4">
        <div v-for="post in posts" :key="post.id" class="flex space-x-4 p-4 bg-gray-50 rounded-lg">
          <img :src="getUserAvatar(post.user)" class="w-10 h-10 rounded-full">
          <div class="flex-1">
            <div class="flex items-center space-x-2 mb-1">
              <span class="font-medium text-gray-900">{{ getUserDisplayName(post.user) }}</span>
              <span class="text-sm text-gray-500">{{ formatTime(post.created_at) }}</span>
            </div>
            <p class="text-gray-700">{{ post.content }}</p>
            <div class="flex items-center space-x-4 mt-2 text-sm">
              <button @click="togglePostLike(post)"
                :class="['transition-colors', getPostLiked(post.id) ? 'text-red-500' : 'text-gray-500 hover:text-red-500']">
                {{ getPostLiked(post.id) ? '❤️' : '🤍' }} {{ post.like_count }}
              </button>
              <button class="text-gray-500 hover:text-blue-500">回复</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片查看器 -->
    <div v-if="showLightbox" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-90"
      @click="closeLightbox">
      <button @click="closeLightbox"
        class="absolute top-4 right-4 text-white text-4xl hover:text-gray-300 transition-colors z-10">
        ×
      </button>
      <img :src="lightboxImage" class="max-w-full max-h-full object-contain" @click.stop>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useConfigStore } from '@/stores/config'
import api, { pollApi, topicApi } from '@/api'
import { ElMessage } from 'element-plus'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'

const route = useRoute()
const userStore = useUserStore()
const configStore = useConfigStore()
const topic = ref(null)
const posts = ref([])
const newPost = ref('')
const liked = ref(false)
const favorited = ref(false)
const postLikes = ref({})
const showLightbox = ref(false)
const lightboxImage = ref('')

const poll = ref(null)
const selectedOptions = ref([])
const submittingVote = ref(false)
const votedOptionIds = ref([])
const hasVotedFromServer = ref(false)

const hasVoted = computed(() => hasVotedFromServer.value || votedOptionIds.value.length > 0)
const isPollEnded = computed(() => {
  if (!poll.value?.end_time) return false
  return new Date(poll.value.end_time) < new Date()
})

const canDeleteTopic = computed(() => {
  if (!userStore.isLoggedIn || !topic.value) return false
  const isAuthor = topic.value.user_id === userStore.user?.id
  const isAdmin = userStore.user?.role === 2
  return isAuthor || isAdmin
})

function getPercentage(voteCount) {
  if (!poll.value || poll.value.total_votes === 0) return 0
  return Math.round((voteCount / poll.value.total_votes) * 100)
}

function getRemainingTime(endTime) {
  const end = new Date(endTime)
  const now = new Date()
  const diff = end - now
  
  if (diff <= 0) return '已结束'
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  
  if (days > 0) return `${days}天${hours}小时`
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  if (hours > 0) return `${hours}小时${minutes}分钟`
  return `${minutes}分钟`
}

function selectOption(optionId) {
  selectedOptions.value = [optionId]
}

function toggleOption(optionId) {
  const index = selectedOptions.value.indexOf(optionId)
  if (index > -1) {
    selectedOptions.value.splice(index, 1)
  } else {
    if (selectedOptions.value.length < poll.value.max_choices) {
      selectedOptions.value.push(optionId)
    } else {
      ElMessage.warning(`最多只能选择${poll.value.max_choices}项`)
    }
  }
}

async function loadPoll() {
  try {
    const res = await pollApi.getPollByTopic(route.params.id)
    console.log('loadPoll response:', res)
    if (res && res.poll) {
      poll.value = res.poll
      hasVotedFromServer.value = res.has_voted || false
      votedOptionIds.value = res.voted_option_ids || []
      console.log('hasVotedFromServer:', hasVotedFromServer.value, 'votedOptionIds:', votedOptionIds.value)
    }
  } catch (e) {
    console.log('No poll for this topic')
  }
}

async function handleDeleteTopic() {
  if (!confirm('确定要删除这篇帖子吗？删除后无法恢复。')) {
    return
  }
  
  try {
    await topicApi.deleteTopic(topic.value.id)
    ElMessage.success('帖子已删除')
    window.location.href = '/'
  } catch (e) {
    console.error('删除失败', e)
    ElMessage.error(e.response?.data?.message || '删除失败')
  }
}

async function submitVote() {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  if (selectedOptions.value.length === 0) {
    ElMessage.warning('请选择投票选项')
    return
  }

  submittingVote.value = true
  try {
    const res = await pollApi.submitVote({
      poll_id: poll.value.id,
      option_ids: selectedOptions.value
    })
    console.log('submitVote response:', res)

    // 立即更新状态
    hasVotedFromServer.value = true
    votedOptionIds.value = [...selectedOptions.value]

    // 再刷新获取最新数据
    await loadPoll()
    ElMessage.success('投票成功')
  } catch (e) {
    console.error(e)
    ElMessage.error(e.response?.data?.message || '投票失败')
  } finally {
    submittingVote.value = false
  }
}

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

async function loadTopic() {
  try {
    const id = route.params.id
    topic.value = await api.get(`/topics/${id}`)
    const postsRes = await api.get(`/topics/${id}/posts`)
    posts.value = postsRes.list || postsRes || []

    if (userStore.isLoggedIn) {
      await checkLikeStatus()
      await checkFavoriteStatus()
    }

    await loadPoll()

    // 设置媒体查看器
    setupMediaViewers()
  } catch (e) {
    console.error(e)
  }
}

async function checkLikeStatus() {
  try {
    const res = await api.post('/likes/check', {
      target_type: 'topic',
      target_id: topic.value.id
    })
    liked.value = res.liked
  } catch (e) {
    console.error(e)
  }
}

async function checkFavoriteStatus() {
  try {
    const res = await api.post('/favorites/check', {
      topic_id: topic.value.id
    })
    favorited.value = res.favorited
  } catch (e) {
    console.error(e)
  }
}

async function toggleLike() {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (liked.value) {
      await api.delete(`/likes?target_type=topic&target_id=${topic.value.id}`)
      topic.value.like_count--
    } else {
      await api.post('/likes', {
        target_type: 'topic',
        target_id: topic.value.id
      })
      topic.value.like_count++
    }
    liked.value = !liked.value
  } catch (e) {
    console.error(e)
    ElMessage.error('操作失败')
  }
}

async function toggleFavorite() {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (favorited.value) {
      await api.delete(`/favorites?topic_id=${topic.value.id}`)
    } else {
      await api.post('/favorites', {
        topic_id: topic.value.id
      })
    }
    favorited.value = !favorited.value
  } catch (e) {
    console.error(e)
    ElMessage.error('操作失败')
  }
}

async function shareTopic() {
  const url = window.location.href
  try {
    if (navigator.share) {
      await navigator.share({
        title: topic.value.title,
        url: url
      })
    } else {
      await navigator.clipboard.writeText(url)
      ElMessage.success('链接已复制到剪贴板')
    }
  } catch (e) {
    console.error(e)
  }
}

function getPostLiked(postId) {
  return postLikes.value[postId] || false
}

async function togglePostLike(post) {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (getPostLiked(post.id)) {
      await api.delete(`/likes?target_type=post&target_id=${post.id}`)
      post.like_count--
    } else {
      await api.post('/likes', {
        target_type: 'post',
        target_id: post.id
      })
      post.like_count++
    }
    postLikes.value[post.id] = !getPostLiked(post.id)
  } catch (e) {
    console.error(e)
    ElMessage.error('操作失败')
  }
}

async function submitPost() {
  if (!newPost.value.trim()) return
  try {
    await api.post(`/topics/${route.params.id}/posts`, { content: newPost.value })
    newPost.value = ''
    loadTopic()
  } catch (e) {
    console.error(e)
  }
}

function openLightbox(src) {
  lightboxImage.value = src
  showLightbox.value = true
  document.body.style.overflow = 'hidden'
}

function closeLightbox() {
  showLightbox.value = false
  document.body.style.overflow = ''
}

function setupMediaViewers() {
  setTimeout(() => {
    const content = document.querySelector('.topic-content')
    if (!content) return

    // 为图片添加点击放大功能
    const images = content.querySelectorAll('img')
    images.forEach(img => {
      img.style.cursor = 'pointer'
      img.style.transition = 'transform 0.2s, box-shadow 0.2s'
      
      img.addEventListener('mouseenter', () => {
        img.style.transform = 'scale(1.02)'
        img.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)'
      })
      
      img.addEventListener('mouseleave', () => {
        img.style.transform = ''
        img.style.boxShadow = ''
      })
      
      img.addEventListener('click', () => {
        openLightbox(img.src)
      })
    })

    // 为视频添加全屏播放功能
    const videos = content.querySelectorAll('video')
    videos.forEach(video => {
      video.style.cursor = 'pointer'
      video.style.transition = 'transform 0.2s, box-shadow 0.2s'
      
      video.addEventListener('mouseenter', () => {
        video.style.transform = 'scale(1.02)'
        video.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)'
      })
      
      video.addEventListener('mouseleave', () => {
        video.style.transform = ''
        video.style.boxShadow = ''
      })
      
      video.addEventListener('click', () => {
        if (video.requestFullscreen) {
          video.requestFullscreen()
        } else if (video.webkitRequestFullscreen) {
          video.webkitRequestFullscreen()
        } else if (video.msRequestFullscreen) {
          video.msRequestFullscreen()
        }
      })
    })
  }, 100)
}

onMounted(() => {
  loadTopic()
})
</script>

<style scoped>
.topic-content img,
.topic-content video {
  border-radius: 8px;
  margin: 1rem 0;
  max-width: 100%;
  height: auto;
}

.topic-content img:hover,
.topic-content video:hover {
  position: relative;
}

.topic-content img::after,
.topic-content video::after {
  content: '🔍';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 24px;
  opacity: 0;
  transition: opacity 0.2s;
}

.topic-content img:hover::after,
.topic-content video:hover::after {
  opacity: 1;
}
</style>
