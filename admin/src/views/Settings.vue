<template>
  <div class="settings-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <h3>
            <Sliders :size="18" />
            系统设置
          </h3>
        </div>
      </template>

      <el-form ref="formRef" :model="settings" label-position="top">
        <div class="settings-section">
          <h4 class="section-title">
            <ToggleLeft :size="16" />
            功能开关
          </h4>
          <div class="switch-grid">
            <div class="switch-item">
              <div class="switch-info">
                <span class="switch-label">允许用户注册</span>
                <span class="switch-desc">控制是否开放新用户注册</span>
              </div>
              <el-switch v-model="settings.allow_register" />
            </div>
            <div class="switch-item">
              <div class="switch-info">
                <span class="switch-label">允许发帖</span>
                <span class="switch-desc">控制是否允许用户发布帖子</span>
              </div>
              <el-switch v-model="settings.allow_post" />
            </div>
            <div class="switch-item">
              <div class="switch-info">
                <span class="switch-label">允许评论</span>
                <span class="switch-desc">控制是否允许用户发表评论</span>
              </div>
              <el-switch v-model="settings.allow_comment" />
            </div>
            <div class="switch-item">
              <div class="switch-info">
                <span class="switch-label">允许发起投票</span>
                <span class="switch-desc">控制是否允许用户发布投票</span>
              </div>
              <el-switch v-model="settings.allow_poll" />
            </div>
          </div>
        </div>

        <div class="settings-section">
          <h4 class="section-title">
            <Coins :size="16" />
            积分规则
          </h4>
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="发帖积分">
                <el-input-number v-model="settings.credits_topic" :min="0" :max="999" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="评论积分">
                <el-input-number v-model="settings.credits_post" :min="0" :max="999" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="签到积分">
                <el-input-number v-model="settings.credits_signin" :min="0" :max="999" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <div class="form-footer">
          <el-button @click="resetSettings">重置</el-button>
          <el-button type="primary" @click="saveSettings" :loading="saving">
            <Save :size="16" />
            保存设置
          </el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Sliders, ToggleLeft, Coins, Save } from 'lucide-vue-next'
import api from '@/api'

const settings = ref({
  allow_register: true,
  allow_post: true,
  allow_comment: true,
  allow_poll: true,
  credits_topic: 10,
  credits_post: 5,
  credits_signin: 10
})

const saving = ref(false)
const formRef = ref(null)

async function loadSettings() {
  try {
    const config = await api.get('/config')
    settings.value = {
      allow_register: config.allow_register !== 'false',
      allow_post: config.allow_post !== 'false',
      allow_comment: config.allow_comment !== 'false',
      allow_poll: config.allow_poll !== 'false',
      credits_topic: parseInt(config.credit_topic) || 10,
      credits_post: parseInt(config.credit_post) || 5,
      credits_signin: parseInt(config.credit_signin) || 10
    }
  } catch (e) {
    console.error('加载设置失败', e)
  }
}

function resetSettings() {
  loadSettings()
}

async function saveSettings() {
  saving.value = true
  try {
    await api.put('/admin/config', {
      allow_register: settings.value.allow_register ? 'true' : 'false',
      allow_post: settings.value.allow_post ? 'true' : 'false',
      allow_comment: settings.value.allow_comment ? 'true' : 'false',
      allow_poll: settings.value.allow_poll ? 'true' : 'false',
      credit_topic: String(settings.value.credits_topic),
      credit_post: String(settings.value.credits_post),
      credit_signin: String(settings.value.credits_signin)
    })
    ElMessage.success('设置已保存')
  } catch (e) {
    console.error('保存设置失败', e)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings-page {
  max-width: 1000px;
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

.settings-section {
  margin-bottom: 32px;
}

.settings-section:last-of-type {
  margin-bottom: 0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f3f4f6;
}

.switch-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.switch-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
}

.switch-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.switch-label {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
}

.switch-desc {
  font-size: 12px;
  color: #9ca3af;
}

.form-footer {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #f3f4f6;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
}
</style>
