<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6">{{ t('login.title') }}</h2>
      <el-form @submit.prevent="handleLogin" :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item :label="t('login.username')" prop="username">
          <el-input v-model="form.username" :placeholder="t('login.usernamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('login.password')" prop="password">
          <el-input v-model="form.password" type="password" :placeholder="t('login.passwordPlaceholder')" show-password />
        </el-form-item>
        <el-form-item>
          <button type="submit" class="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition-colors">
            {{ t('login.loginBtn') }}
          </button>
        </el-form-item>
      </el-form>
      <p class="text-center mt-4 text-gray-600">
        {{ t('login.noAccount') }}<router-link to="/register" class="text-blue-500 hover:underline">{{ t('login.goRegister') }}</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getErrorI18nKey } from '@/utils/error'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)

const form = ref({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '', trigger: 'blur' }],
  password: [{ required: true, message: '', trigger: 'blur' }]
}

async function handleLogin() {
  try {
    await formRef.value.validate()
    await userStore.login(form.value)
    router.push('/')
  } catch (e) {
    if (e.code) {
      const errorKey = getErrorI18nKey(e.code)
      ElMessage.error(t(errorKey))
    }
  }
}
</script>
