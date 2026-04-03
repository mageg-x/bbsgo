<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
    <div v-if="!configStore.state.allow_register"
      class="bg-white p-8 rounded-2xl shadow-xl w-full max-w-md text-center">
      <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>
      </svg>
      <h2 class="text-2xl font-bold text-gray-900 mb-2">注册功能已关闭</h2>
      <p class="text-gray-500 mb-6">管理员暂时关闭了注册功能，请稍后再试。</p>
      <router-link to="/login"
        class="inline-block w-full bg-blue-500 text-white py-3 rounded-lg hover:bg-blue-600 transition-colors font-medium">
        返回登录
      </router-link>
    </div>
    <div v-else class="bg-white p-8 rounded-2xl shadow-xl w-full max-w-md">
      <div class="text-center mb-8">
        <h2 class="text-3xl font-bold text-gray-800">注册账号</h2>
        <p class="text-gray-500 mt-2">加入我们，开始您的社区之旅</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-5">
        <div>
          <label class="block text-gray-700 text-sm font-medium mb-2">用户名</label>
          <input type="text" v-model="form.username"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            placeholder="请输入用户名" required />
        </div>

        <div>
          <label class="block text-gray-700 text-sm font-medium mb-2">昵称</label>
          <input type="text" v-model="form.nickname"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            placeholder="请输入昵称" required />
        </div>

        <div>
          <label class="block text-gray-700 text-sm font-medium mb-2">邮箱</label>
          <input type="email" v-model="form.email"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            placeholder="请输入邮箱" required />
        </div>

        <div v-if="emailEnabled">
          <label class="block text-gray-700 text-sm font-medium mb-2">验证码</label>
          <div class="flex gap-3">
            <input type="text" v-model="form.code"
              class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
              placeholder="请输入验证码" maxlength="6" required />
            <button type="button" @click="sendCode" :disabled="countdown > 0 || !form.email"
              class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap">
              {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
            </button>
          </div>
        </div>

        <div>
          <label class="block text-gray-700 text-sm font-medium mb-2">密码</label>
          <input type="password" v-model="form.password"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            placeholder="请输入密码" required />
        </div>

        <div>
          <label class="block text-gray-700 text-sm font-medium mb-2">确认密码</label>
          <input type="password" v-model="form.confirm_password"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
            placeholder="请再次输入密码" required />
        </div>

        <button type="submit" :disabled="loading"
          class="w-full bg-blue-500 text-white py-3 rounded-lg hover:bg-blue-600 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>

      <p class="text-center mt-6 text-gray-600">
        已有账号？
        <router-link to="/login" class="text-blue-500 hover:underline font-medium">立即登录</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useConfigStore } from '@/stores/config'
import { ElMessage } from 'element-plus'
import api from '@/api'

const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()

const form = ref({
  username: '',
  nickname: '',
  email: '',
  code: '',
  password: '',
  confirm_password: ''
})

const loading = ref(false)
const countdown = ref(0)
const emailEnabled = ref(false)
let timer = null

async function checkEmailEnabled() {
  try {
    const res = await api.get('/config')
    emailEnabled.value = res.email_enabled === 'true'
  } catch (e) {
    console.error('获取配置失败', e)
  }
}

async function sendCode() {
  if (!form.value.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.value.email)) {
    ElMessage.warning('请输入有效的邮箱地址')
    return
  }

  try {
    await api.post('/send-code', {
      email: form.value.email,
      type: 'register'
    })
    ElMessage.success('验证码已发送到您的邮箱')
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (e) {
    // 统一提示验证码发送失败
    ElMessage.error('验证码发送失败，请重试')
  }
}

async function handleRegister() {
  if (form.value.password !== form.value.confirm_password) {
    ElMessage.warning('两次密码输入不一致')
    return
  }

  if (emailEnabled.value && !form.value.code) {
    ElMessage.warning('请输入验证码')
    return
  }

  loading.value = true
  try {
    await userStore.register(form.value)
    ElMessage.success('注册成功！')
    router.push('/')
  } catch (e) {
    // 错误已在 interceptor 中显示
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  checkEmailEnabled()
})
</script>
