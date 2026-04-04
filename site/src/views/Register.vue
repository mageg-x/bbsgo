<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 py-6">
    <div v-if="!configStore.state.allow_register"
      class="bg-white p-6 rounded-xl shadow-lg w-full max-w-md text-center">
      <svg class="w-12 h-12 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>
      </svg>
      <h2 class="text-xl font-bold text-gray-900 mb-2">{{ t('register.closed') }}</h2>
      <p class="text-gray-500 mb-4">{{ t('register.closedTip') }}</p>
      <router-link to="/login"
        class="inline-block w-full bg-blue-500 text-white py-2.5 rounded-lg hover:bg-blue-600 transition-colors font-medium">
        {{ t('register.backToLogin') }}
      </router-link>
    </div>
    <div v-else class="bg-white p-6 rounded-xl shadow-lg w-full max-w-md">
      <div class="text-center mb-5">
        <h2 class="text-2xl font-bold text-gray-800">{{ t('register.title') }}</h2>
        <p class="text-gray-500 mt-1 text-sm">{{ t('register.subtitle') }}</p>
      </div>

      <el-form @submit.prevent="handleRegister" :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item :label="t('register.username')" prop="username">
          <el-input v-model="form.username" :placeholder="t('register.usernamePlaceholder')" />
        </el-form-item>

        <el-form-item :label="t('register.nickname')" prop="nickname">
          <el-input v-model="form.nickname" :placeholder="t('register.nicknamePlaceholder')" />
        </el-form-item>

        <el-form-item :label="t('register.email')" prop="email">
          <el-input v-model="form.email" type="email" :placeholder="t('register.emailPlaceholder')" />
        </el-form-item>

        <el-form-item v-if="emailEnabled" :label="t('register.emailCode')" prop="code">
          <div class="flex gap-2">
            <el-input v-model="form.code" :placeholder="t('register.codePlaceholder')" maxlength="6" />
            <button type="button" @click="sendCode" :disabled="countdown > 0 || !form.email"
              class="px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap text-sm">
              {{ countdown > 0 ? `${countdown}s` : t('register.sendCode') }}
            </button>
          </div>
        </el-form-item>

        <el-form-item :label="t('register.password')" prop="password">
          <el-input v-model="form.password" type="password" :placeholder="t('register.passwordPlaceholder')" show-password />
        </el-form-item>

        <el-form-item :label="t('register.confirmPassword')" prop="confirm_password">
          <el-input v-model="form.confirm_password" type="password" :placeholder="t('register.confirmPasswordPlaceholder')" show-password />
        </el-form-item>

        <el-form-item>
          <button type="submit" :disabled="loading"
            class="w-full bg-blue-500 text-white py-2.5 rounded-lg hover:bg-blue-600 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed text-sm">
            {{ loading ? t('register.registering') : t('register.registerBtn') }}
          </button>
        </el-form-item>
      </el-form>

      <p class="text-center mt-4 text-gray-600 text-sm">
        {{ t('register.hasAccount') }}
        <router-link to="/login" class="text-blue-500 hover:underline font-medium">{{ t('register.goLogin') }}</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { useConfigStore } from '@/stores/config'
import { ElMessage } from 'element-plus'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const configStore = useConfigStore()
const formRef = ref(null)

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

// 自定义验证：确认密码
const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.value.password) {
    callback(new Error(''))
  } else {
    callback()
  }
}

const rules = {
  username: [{ required: true, message: '', trigger: 'blur' }],
  nickname: [{ required: true, message: '', trigger: 'blur' }],
  email: [
    { required: true, message: '', trigger: 'blur' },
    { type: 'email', message: '', trigger: 'blur' }
  ],
  code: [{ required: emailEnabled.value, message: '', trigger: 'blur' }],
  password: [{ required: true, message: '', trigger: 'blur' }],
  confirm_password: [
    { required: true, message: '', trigger: 'blur' },
    { validator: validateConfirmPassword, message: '', trigger: 'blur' }
  ]
}

async function checkEmailEnabled() {
  try {
    const res = await api.get('/config')
    emailEnabled.value = res.email_enabled === 'true'
  } catch (e) {
    console.error('get config failed', e)
  }
}

async function sendCode() {
  if (!form.value.email) {
    ElMessage.warning(t('register.enterEmailFirst'))
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.value.email)) {
    ElMessage.warning(t('register.invalidEmail'))
    return
  }

  try {
    await api.post('/send-code', {
      email: form.value.email,
      type: 'register'
    })
    ElMessage.success(t('register.codeSent'))
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (e) {
    ElMessage.error(t('register.codeSendFailed'))
  }
}

async function handleRegister() {
  if (form.value.password !== form.value.confirm_password) {
    ElMessage.warning(t('register.passwordMismatch'))
    return
  }

  if (emailEnabled.value && !form.value.code) {
    ElMessage.warning(t('register.enterCode'))
    return
  }

  loading.value = true
  try {
    await formRef.value.validate()
    await userStore.register(form.value)
    ElMessage.success(t('register.success'))
    router.push('/')
  } catch (e) {
    // 验证失败或注册失败
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  checkEmailEnabled()
})
</script>
