<template>
  <button
    v-if="!isCurrentUser"
    @click="toggleFollow"
    :class="[
      'follow-button',
      isFollowing ? 'following' : 'not-following',
      { loading: loading }
    ]"
    :disabled="loading"
  >
    <span v-if="loading" class="loading-spinner"></span>
    <span v-else>{{ isFollowing ? '已关注' : '关注' }}</span>
  </button>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import api from '@/api'

const props = defineProps({
  userId: {
    type: Number,
    required: true
  }
})

const userStore = useUserStore()
const isFollowing = ref(false)
const loading = ref(false)

const isCurrentUser = computed(() => {
  return userStore.user?.id === props.userId
})

async function checkFollowStatus() {
  if (isCurrentUser.value) return
  
  try {
    const res = await api.get('/follows/check', {
      params: { user_id: props.userId }
    })
    isFollowing.value = res.is_following
  } catch (e) {
    console.error('检查关注状态失败', e)
  }
}

async function toggleFollow() {
  if (!userStore.user) {
    return
  }

  loading.value = true
  try {
    if (isFollowing.value) {
      await api.delete('/follows', {
        data: { follow_user_id: props.userId }
      })
      isFollowing.value = false
    } else {
      await api.post('/follows', {
        follow_user_id: props.userId
      })
      isFollowing.value = true
    }
  } catch (e) {
    console.error('关注操作失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  checkFollowStatus()
})
</script>

<style scoped>
.follow-button {
  padding: 6px 16px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  outline: none;
  min-width: 80px;
}

.follow-button.not-following {
  background-color: #3b82f6;
  color: white;
}

.follow-button.not-following:hover {
  background-color: #2563eb;
}

.follow-button.following {
  background-color: #f3f4f6;
  color: #6b7280;
  border: 1px solid #e5e7eb;
}

.follow-button.following:hover {
  background-color: #fee2e2;
  color: #dc2626;
  border-color: #fecaca;
}

.follow-button.loading {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: currentColor;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
