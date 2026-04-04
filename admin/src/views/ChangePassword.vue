<template>
  <div class="change-password-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
        </div>
      </template>

      <div class="form-wrapper">
        <el-form ref="formRef" :model="form" :rules="rules" label-position="top" class="password-form">
          <el-form-item :label="t('password.oldPassword')" prop="old_password">
            <el-input v-model="form.old_password" type="password" :placeholder="t('password.oldPassword')" show-password size="large">
              <template #prefix>
                <Lock :size="18" />
              </template>
            </el-input>
          </el-form-item>

          <el-form-item :label="t('password.newPassword')" prop="new_password">
            <el-input v-model="form.new_password" type="password" :placeholder="t('password.newPasswordPlaceholder')" show-password size="large">
              <template #prefix>
                <Key :size="18" />
              </template>
            </el-input>
          </el-form-item>

          <el-form-item :label="t('password.confirmPassword')" prop="confirm_password">
            <el-input v-model="form.confirm_password" type="password" :placeholder="t('password.confirmPasswordPlaceholder')" show-password size="large">
              <template #prefix>
                <KeyRound :size="18" />
              </template>
            </el-input>
          </el-form-item>

          <div class="form-actions">
            <el-button size="large" @click="resetForm">{{ t('common.reset') }}</el-button>
            <el-button type="primary" size="large" @click="handleSubmit" :loading="loading">
              <Save :size="16" />
              {{ t('password.submit') }}
            </el-button>
          </div>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAdminStore } from '@/stores/admin'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { Key, Lock, KeyRound, Save } from 'lucide-vue-next'

const { t } = useI18n()
const router = useRouter()
const adminStore = useAdminStore()
const formRef = ref(null)
const loading = ref(false)

const form = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirm = (rule, value, callback) => {
  if (value !== form.value.new_password) {
    callback(new Error(t('password.passwordMismatch')))
  } else {
    callback()
  }
}

const rules = {
  old_password: [
    { required: true, message: () => t('password.oldPassword'), trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: () => t('password.newPassword'), trigger: 'blur' },
    { min: 6, message: () => t('password.minLength'), trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: () => t('password.confirmPassword'), trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' }
  ]
}

function resetForm() {
  form.value = {
    old_password: '',
    new_password: '',
    confirm_password: ''
  }
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await api.post('/admin/change-password', {
        old_password: form.value.old_password,
        new_password: form.value.new_password
      })
      ElMessage.success(t('password.success'))
      adminStore.logout()
      router.push('/login')
    } catch (e) {
      console.error('change password failed', e)
      ElMessage.error(e.response?.data?.message || t('common.failed'))
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.change-password-page {
  max-width: 1400px;
}

.main-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.card-header h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.form-wrapper {
  padding: 8px 0;

}

.password-form {
  max-width: 420px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 32px;
}

.form-actions .el-button {
  flex: 1;
}

:deep(.el-input__prefix) {
  color: #9ca3af;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
}
</style>
