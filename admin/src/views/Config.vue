<template>
  <div class="config-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane :label="t('config.basic')" name="basic">
          <el-form :model="config" label-position="top">
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item :label="t('config.siteName')">
                  <el-input v-model="config.site_name" :placeholder="t('config.siteNamePlaceholder')">
                    <template #prefix>
                      <Globe :size="16" />
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('config.siteDescription')">
                  <el-input v-model="config.site_description" :placeholder="t('config.siteDescPlaceholder')" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item :label="t('config.siteLogo')">
                  <div class="upload-container">
                    <el-upload
                      class="image-uploader"
                      :show-file-list="false"
                      :http-request="handleLogoUpload"
                      :before-upload="beforeUpload"
                      accept="image/*"
                    >
                      <div v-if="config.site_logo" class="image-preview">
                        <img :src="config.site_logo" class="preview-img" />
                        <div class="image-actions">
                          <el-button type="primary" size="small" @click.stop="config.site_logo = ''">
                            {{ t('config.changeImage') }}
                          </el-button>
                        </div>
                      </div>
                      <div v-else class="upload-placeholder">
                        <el-icon class="upload-icon"><Plus /></el-icon>
                        <div class="upload-text">{{ t('config.uploadLogo') }}</div>
                        <div class="upload-hint">{{ t('config.logoSize') }}</div>
                      </div>
                    </el-upload>
                    <div v-if="uploadingLogo" class="upload-loading">
                      <el-icon class="is-loading"><Loading /></el-icon>
                      <span>{{ t('config.uploading') }}</span>
                    </div>
                    <div v-if="config.site_logo" class="url-display">
                      <el-input v-model="config.site_logo" placeholder="Logo URL" readonly>
                        <template #prepend>{{ t('config.logoUrl') }}</template>
                      </el-input>
                    </div>
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('config.siteIcon')">
                  <div class="upload-container">
                    <el-upload
                      class="image-uploader icon-uploader"
                      :show-file-list="false"
                      :http-request="handleIconUpload"
                      :before-upload="beforeUpload"
                      accept="image/*"
                    >
                      <div v-if="config.site_icon" class="image-preview icon-preview">
                        <img :src="config.site_icon" class="preview-img icon-img" />
                        <div class="image-actions">
                          <el-button type="primary" size="small" @click.stop="config.site_icon = ''">
                            {{ t('config.changeImage') }}
                          </el-button>
                        </div>
                      </div>
                      <div v-else class="upload-placeholder icon-placeholder">
                        <el-icon class="upload-icon"><Plus /></el-icon>
                        <div class="upload-text">{{ t('config.uploadIcon') }}</div>
                        <div class="upload-hint">{{ t('config.iconSize') }}</div>
                      </div>
                    </el-upload>
                    <div v-if="uploadingIcon" class="upload-loading">
                      <el-icon class="is-loading"><Loading /></el-icon>
                      <span>{{ t('config.uploading') }}</span>
                    </div>
                    <div v-if="config.site_icon" class="url-display">
                      <el-input v-model="config.site_icon" placeholder="Icon URL" readonly>
                        <template #prepend>{{ t('config.iconUrl') }}</template>
                      </el-input>
                    </div>
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="t('config.email')" name="email">
          <el-alert
            :title="t('config.emailTip')"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px"
          >
            <template #default>
              <p style="margin: 0">{{ t('config.emailTip2') }}</p>
            </template>
          </el-alert>

          <el-form :model="config" label-position="top">
            <el-form-item :label="t('config.emailEnabled')">
              <el-switch v-model="config.email_enabled" />
              <div style="color: #909399; font-size: 12px; margin-top: 4px">
                {{ t('config.emailEnabledDesc') }}
              </div>
            </el-form-item>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item :label="t('config.emailHost')">
                  <el-input v-model="config.email_host" placeholder="smtp.example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('config.emailPort')">
                  <el-input v-model="config.email_port" placeholder="465" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item :label="t('config.emailUser')">
                  <el-input v-model="config.email_user" placeholder="noreply@example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('config.emailPassword')">
                  <el-input v-model="config.email_password" type="password" :placeholder="t('config.emailPassword')" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item :label="t('config.emailFrom')">
                  <el-input v-model="config.email_from" placeholder="noreply@example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('config.emailFromName')">
                  <el-input v-model="config.email_from_name" :placeholder="t('config.siteName')" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="t('config.storage')" name="storage">
          <el-form :model="config" label-position="top">
            <el-form-item :label="t('config.storageType')">
              <el-select v-model="config.storage_type" :placeholder="t('config.selectStorageType')" class="w-full">
                <el-option :label="t('config.storageLocal')" value="local" />
                <el-option :label="t('config.storageQiniu')" value="qiniu" />
                <el-option :label="t('config.storageAliyun')" value="aliyun" />
                <el-option :label="t('config.storageTencent')" value="tencent" />
              </el-select>
            </el-form-item>

            <template v-if="config.storage_type === 'local'">
              <el-alert type="info" :closable="false" show-icon>
                <template #title>
                  {{ t('config.localStorageTip') }}
                </template>
                <template #default>
                  {{ t('config.localStorageDesc') }}
                </template>
              </el-alert>
            </template>

            <template v-if="config.storage_type === 'qiniu'">
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.qiniuAccessKey')">
                    <el-input v-model="config.qiniu_access_key" :placeholder="t('config.qiniuAccessKey')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="t('config.qiniuSecretKey')">
                    <el-input v-model="config.qiniu_secret_key" type="password" :placeholder="t('config.qiniuSecretKey')" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.qiniuBucket')">
                    <el-input v-model="config.qiniu_bucket" :placeholder="t('config.qiniuBucket')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="t('config.qiniuDomain')">
                    <el-input v-model="config.qiniu_domain" placeholder="cdn.example.com" />
                  </el-form-item>
                </el-col>
              </el-row>
            </template>

            <template v-if="config.storage_type === 'aliyun'">
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.aliyunAccessKeyId')">
                    <el-input v-model="config.aliyun_access_key_id" :placeholder="t('config.aliyunAccessKeyId')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="t('config.aliyunAccessKeySecret')">
                    <el-input v-model="config.aliyun_access_key_secret" type="password" :placeholder="t('config.aliyunAccessKeySecret')" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.aliyunBucket')">
                    <el-input v-model="config.aliyun_bucket" :placeholder="t('config.aliyunBucket')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="t('config.aliyunEndpoint')">
                    <el-input v-model="config.aliyun_endpoint" placeholder="oss-cn-hangzhou.aliyuncs.com" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.aliyunDomain')">
                    <el-input v-model="config.aliyun_domain" placeholder="cdn.example.com" />
                  </el-form-item>
                </el-col>
              </el-row>
            </template>

            <template v-if="config.storage_type === 'tencent'">
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.tencentSecretId')">
                    <el-input v-model="config.tencent_secret_id" :placeholder="t('config.tencentSecretId')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="t('config.tencentSecretKey')">
                    <el-input v-model="config.tencent_secret_key" type="password" :placeholder="t('config.tencentSecretKey')" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.tencentBucket')">
                    <el-input v-model="config.tencent_bucket" :placeholder="t('config.tencentBucket')" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item :label="t('config.tencentRegion')">
                    <el-input v-model="config.tencent_region" placeholder="ap-guangzhou" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item :label="t('config.tencentDomain')">
                    <el-input v-model="config.tencent_domain" placeholder="cdn.example.com" />
                  </el-form-item>
                </el-col>
              </el-row>
            </template>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="t('config.security')" name="security">
          <el-form :model="config" label-position="top">
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item :label="t('config.jwtSecret')">
                  <el-input v-model="config.jwt_secret" type="password" :placeholder="t('config.jwtSecretPlaceholder')" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="t('config.jwtExpireDays')">
                  <el-input v-model="config.jwt_expire_days" placeholder="7" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div class="form-footer">
        <el-button @click="loadConfig">{{ t('config.reset') }}</el-button>
        <el-button type="primary" @click="saveConfig" :loading="saving">
          <Save :size="16" />
          {{ t('config.saveConfig') }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, Loading } from '@element-plus/icons-vue'
import api from '@/api'
import { Settings, Globe, Save } from 'lucide-vue-next'
import { compressAndUpload } from '@/utils/upload'

const { t } = useI18n()
const activeTab = ref('basic')
const config = ref({
  site_name: '',
  site_logo: '',
  site_icon: '',
  site_description: '',
  email_enabled: false,
  email_host: '',
  email_port: '465',
  email_user: '',
  email_password: '',
  email_from: '',
  email_from_name: '',
  storage_type: 'local',
  qiniu_access_key: '',
  qiniu_secret_key: '',
  qiniu_bucket: '',
  qiniu_domain: '',
  aliyun_access_key_id: '',
  aliyun_access_key_secret: '',
  aliyun_bucket: '',
  aliyun_endpoint: '',
  aliyun_domain: '',
  tencent_secret_id: '',
  tencent_secret_key: '',
  tencent_bucket: '',
  tencent_region: '',
  tencent_domain: '',
  jwt_secret: '',
  jwt_expire_days: '7'
})

const saving = ref(false)
const uploadingLogo = ref(false)
const uploadingIcon = ref(false)

async function handleLogoUpload(options) {
  uploadingLogo.value = true

  try {
    const url = await compressAndUpload(options.file, 800, 200, {
      dir: 'logo',
      mimeType: 'image/png',
      onInstant: () => ElMessage.success(t('config.logoUploadSuccess') + '（' + t('config.instantUpload') + '）')
    })
    config.value.site_logo = url
    ElMessage.success(t('config.logoUploadSuccess'))
  } catch (error) {
    console.error('Logo upload error:', error)
    ElMessage.error(t('config.logoUploadFailed'))
  } finally {
    uploadingLogo.value = false
  }
}

async function handleIconUpload(options) {
  uploadingIcon.value = true

  try {
    const url = await compressAndUpload(options.file, 256, 256, {
      dir: 'icon',
      mimeType: 'image/png',
      onInstant: () => ElMessage.success(t('config.iconUploadSuccess') + '（' + t('config.instantUpload') + '）')
    })
    config.value.site_icon = url
    ElMessage.success(t('config.iconUploadSuccess'))
  } catch (error) {
    console.error('Icon upload error:', error)
    ElMessage.error(t('config.iconUploadFailed'))
  } finally {
    uploadingIcon.value = false
  }
}

function beforeUpload(file) {
  const isImage = file.type.startsWith('image/')

  if (!isImage) {
    ElMessage.error(t('config.onlyImage'))
    return false
  }

  return true
}

async function loadConfig() {
  try {
    const res = await api.get('/config')
    if (res) {
      config.value = {
        site_name: res.site_name || '',
        site_logo: res.site_logo || '',
        site_icon: res.site_icon || '',
        site_description: res.site_description || '',
        email_enabled: res.email_enabled === 'true',
        email_host: res.email_host || '',
        email_port: res.email_port || '465',
        email_user: res.email_user || '',
        email_password: res.email_password || '',
        email_from: res.email_from || '',
        email_from_name: res.email_from_name || '',
        storage_type: res.storage_type || 'local',
        qiniu_access_key: res.qiniu_access_key || '',
        qiniu_secret_key: res.qiniu_secret_key || '',
        qiniu_bucket: res.qiniu_bucket || '',
        qiniu_domain: res.qiniu_domain || '',
        aliyun_access_key_id: res.aliyun_access_key_id || '',
        aliyun_access_key_secret: res.aliyun_access_key_secret || '',
        aliyun_bucket: res.aliyun_bucket || '',
        aliyun_endpoint: res.aliyun_endpoint || '',
        aliyun_domain: res.aliyun_domain || '',
        tencent_secret_id: res.tencent_secret_id || '',
        tencent_secret_key: res.tencent_secret_key || '',
        tencent_bucket: res.tencent_bucket || '',
        tencent_region: res.tencent_region || '',
        tencent_domain: res.tencent_domain || '',
        jwt_secret: res.jwt_secret || '',
        jwt_expire_days: res.jwt_expire_days || '7'
      }
    }
  } catch (e) {
    console.error('Load config failed', e)
    ElMessage.error(t('config.loadFailed'))
  }
}

async function saveConfig() {
  saving.value = true
  try {
    const data = {
      ...config.value,
      email_enabled: config.value.email_enabled ? 'true' : 'false'
    }
    await api.put('/admin/config', data)
    ElMessage.success(t('config.saveSuccess'))
    loadConfig()
  } catch (e) {
    console.error('Save config failed', e)
    ElMessage.error(t('config.saveFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.config-page {
  max-width: 1000px;
}

.w-full {
  width: 100%;
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

.upload-container {
  display: inline-block;
}

.image-uploader {
  width: 200px;
  height: 200px;
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s;
}

.image-uploader:hover {
  border-color: #409eff;
}

.upload-placeholder {
  width: 200px;
  height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
}

.upload-icon {
  font-size: 48px;
  color: #8c939d;
  margin-bottom: 8px;
}

.upload-text {
  color: #8c939d;
  font-size: 14px;
}

.upload-hint {
  color: #c0c4cc;
  font-size: 12px;
  margin-top: 4px;
}

.upload-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  color: #409eff;
  font-size: 14px;
  
  .el-icon {
    margin-right: 8px;
    font-size: 16px;
  }
  
  .is-loading {
    animation: rotating 2s linear infinite;
  }
}

@keyframes rotating {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.image-preview {
  width: 200px;
  height: 200px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}

.preview-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.icon-img {
  width: 128px;
  height: 128px;
  object-fit: contain;
}

.image-actions {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
}

.url-display {
  margin-top: 12px;
}

.form-footer {
  margin-top: 24px;
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

:deep(.el-input__prefix) {
  color: #9ca3af;
}

:deep(.el-tabs__item) {
  font-weight: 500;
}

:deep(.el-upload) {
  width: 100%;
  height: 100%;
}

:deep(.el-upload-dragger) {
  width: 100%;
  height: 100%;
  border: none;
  background: transparent;
}
</style>
