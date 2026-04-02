<template>
  <div class="config-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <h3>
            <Settings :size="18" />
            网站配置
          </h3>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本设置" name="basic">
          <el-form :model="config" label-position="top">
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="网站名称">
                  <el-input v-model="config.site_name" placeholder="请输入网站名称">
                    <template #prefix>
                      <Globe :size="16" />
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="网站描述">
                  <el-input v-model="config.site_description" placeholder="请输入网站描述" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="网站 Logo">
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
                            更换图片
                          </el-button>
                        </div>
                      </div>
                      <div v-else class="upload-placeholder">
                        <el-icon class="upload-icon"><Plus /></el-icon>
                        <div class="upload-text">点击上传 Logo</div>
                        <div class="upload-hint">建议尺寸：800×200px</div>
                      </div>
                    </el-upload>
                    <div v-if="uploadingLogo" class="upload-loading">
                      <el-icon class="is-loading"><Loading /></el-icon>
                      <span>正在上传...</span>
                    </div>
                    <div v-if="config.site_logo" class="url-display">
                      <el-input v-model="config.site_logo" placeholder="Logo URL" readonly>
                        <template #prepend>URL</template>
                      </el-input>
                    </div>
                  </div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="网站 Icon">
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
                            更换图片
                          </el-button>
                        </div>
                      </div>
                      <div v-else class="upload-placeholder icon-placeholder">
                        <el-icon class="upload-icon"><Plus /></el-icon>
                        <div class="upload-text">点击上传 Icon</div>
                        <div class="upload-hint">建议尺寸：256×256px</div>
                      </div>
                    </el-upload>
                    <div v-if="uploadingIcon" class="upload-loading">
                      <el-icon class="is-loading"><Loading /></el-icon>
                      <span>正在上传...</span>
                    </div>
                    <div v-if="config.site_icon" class="url-display">
                      <el-input v-model="config.site_icon" placeholder="Icon URL" readonly>
                        <template #prepend>URL</template>
                      </el-input>
                    </div>
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="邮件配置" name="email">
          <el-alert
            title="邮件服务说明"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px"
          >
            <template #default>
              <p style="margin: 0">邮件服务用于用户注册时发送邮箱验证码，确保用户邮箱真实有效。</p>
              <p style="margin: 8px 0 0 0; color: #909399; font-size: 13px">
                启用邮件服务后，用户注册时需要输入邮箱验证码才能完成注册。
              </p>
            </template>
          </el-alert>

          <el-form :model="config" label-position="top">
            <el-form-item label="启用邮件服务">
              <el-switch v-model="config.email_enabled" />
              <div style="color: #909399; font-size: 12px; margin-top: 4px">
                启用后，用户注册时需要邮箱验证码；禁用则跳过邮箱验证
              </div>
            </el-form-item>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="SMTP 服务器">
                  <el-input v-model="config.email_host" placeholder="smtp.example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="SMTP 端口">
                  <el-input v-model="config.email_port" placeholder="465" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="邮箱账号">
                  <el-input v-model="config.email_user" placeholder="noreply@example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="邮箱密码">
                  <el-input v-model="config.email_password" type="password" placeholder="请输入邮箱密码" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="发件人地址">
                  <el-input v-model="config.email_from" placeholder="noreply@example.com" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="发件人名称">
                  <el-input v-model="config.email_from_name" placeholder="网站名称" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="图床配置" name="storage">
          <el-form :model="config" label-position="top">
            <el-form-item label="存储类型">
              <el-select v-model="config.storage_type" placeholder="请选择存储类型" class="w-full">
                <el-option label="本地存储" value="local" />
                <el-option label="七牛云存储" value="qiniu" />
                <el-option label="阿里云OSS" value="aliyun" />
                <el-option label="腾讯云COS" value="tencent" />
              </el-select>
            </el-form-item>

            <template v-if="config.storage_type === 'local'">
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="存储路径">
                    <el-input v-model="config.local_storage_path" placeholder="./uploads" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="访问URL">
                    <el-input v-model="config.local_storage_base_url" placeholder="/uploads" />
                  </el-form-item>
                </el-col>
              </el-row>
            </template>

            <template v-if="config.storage_type === 'qiniu'">
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="Access Key">
                    <el-input v-model="config.qiniu_access_key" placeholder="请输入 Access Key" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Secret Key">
                    <el-input v-model="config.qiniu_secret_key" type="password" placeholder="请输入 Secret Key" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="Bucket 名称">
                    <el-input v-model="config.qiniu_bucket" placeholder="请输入 Bucket 名称" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="CDN 域名">
                    <el-input v-model="config.qiniu_domain" placeholder="cdn.example.com" />
                  </el-form-item>
                </el-col>
              </el-row>
            </template>

            <template v-if="config.storage_type === 'aliyun'">
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="AccessKey ID">
                    <el-input v-model="config.aliyun_access_key_id" placeholder="请输入 AccessKey ID" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="AccessKey Secret">
                    <el-input v-model="config.aliyun_access_key_secret" type="password" placeholder="请输入 AccessKey Secret" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="Bucket">
                    <el-input v-model="config.aliyun_bucket" placeholder="请输入 Bucket" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Endpoint">
                    <el-input v-model="config.aliyun_endpoint" placeholder="oss-cn-hangzhou.aliyuncs.com" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="CDN 域名">
                    <el-input v-model="config.aliyun_domain" placeholder="cdn.example.com" />
                  </el-form-item>
                </el-col>
              </el-row>
            </template>

            <template v-if="config.storage_type === 'tencent'">
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="Secret ID">
                    <el-input v-model="config.tencent_secret_id" placeholder="请输入 Secret ID" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Secret Key">
                    <el-input v-model="config.tencent_secret_key" type="password" placeholder="请输入 Secret Key" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="Bucket">
                    <el-input v-model="config.tencent_bucket" placeholder="请输入 Bucket" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Region">
                    <el-input v-model="config.tencent_region" placeholder="ap-guangzhou" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="CDN 域名">
                    <el-input v-model="config.tencent_domain" placeholder="cdn.example.com" />
                  </el-form-item>
                </el-col>
              </el-row>
            </template>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="安全设置" name="security">
          <el-form :model="config" label-position="top">
            <el-row :gutter="24">
              <el-col :span="12">
                <el-form-item label="JWT Secret">
                  <el-input v-model="config.jwt_secret" type="password" placeholder="请输入 JWT Secret" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Token 过期天数">
                  <el-input v-model="config.jwt_expire_days" placeholder="7" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div class="form-footer">
        <el-button @click="loadConfig">重置</el-button>
        <el-button type="primary" @click="saveConfig" :loading="saving">
          <Save :size="16" />
          保存配置
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Loading } from '@element-plus/icons-vue'
import api from '@/api'
import { Settings, Globe, Save } from 'lucide-vue-next'

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
  local_storage_path: '',
  local_storage_base_url: '',
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

const uploadUrl = computed(() => {
  return `/api/v1/upload`
})

const uploadHeaders = computed(() => {
  const token = localStorage.getItem('admin_token')
  return {
    Authorization: `Bearer ${token}`
  }
})

function compressImage(file, maxWidth, maxHeight, type = 'image/png') {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        let width = img.width
        let height = img.height
        
        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width = width * ratio
          height = height * ratio
        }
        
        canvas.width = width
        canvas.height = height
        
        const ctx = canvas.getContext('2d')
        ctx.drawImage(img, 0, 0, width, height)
        
        canvas.toBlob((blob) => {
          resolve(new File([blob], file.name, { type: type }))
        }, type, 0.9)
      }
      img.onerror = reject
      img.src = e.target.result
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

async function uploadFile(file, type) {
  const formData = new FormData()
  formData.append('file', file)
  
  try {
    const response = await fetch(uploadUrl.value, {
      method: 'POST',
      headers: uploadHeaders.value,
      body: formData
    })
    
    const result = await response.json()
    return result
  } catch (error) {
    console.error('Upload error:', error)
    throw error
  }
}

async function handleLogoUpload(options) {
  uploadingLogo.value = true
  
  try {
    const compressedFile = await compressImage(options.file, 800, 200, 'image/png')
    const result = await uploadFile(compressedFile, 'logo')
    
    if (result.code === 0 && result.data?.url) {
      config.value.site_logo = result.data.url
      ElMessage.success('Logo 上传成功')
    } else {
      ElMessage.error('上传失败')
    }
  } catch (error) {
    console.error('Logo upload error:', error)
    ElMessage.error('Logo 上传失败')
  } finally {
    uploadingLogo.value = false
  }
}

async function handleIconUpload(options) {
  uploadingIcon.value = true
  
  try {
    const compressedFile = await compressImage(options.file, 256, 256, 'image/png')
    const result = await uploadFile(compressedFile, 'icon')
    
    if (result.code === 0 && result.data?.url) {
      config.value.site_icon = result.data.url
      ElMessage.success('Icon 上传成功')
    } else {
      ElMessage.error('上传失败')
    }
  } catch (error) {
    console.error('Icon upload error:', error)
    ElMessage.error('Icon 上传失败')
  } finally {
    uploadingIcon.value = false
  }
}

function beforeUpload(file) {
  const isImage = file.type.startsWith('image/')
  
  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
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
        local_storage_path: res.local_storage_path || '',
        local_storage_base_url: res.local_storage_base_url || '',
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
    console.error('加载配置失败', e)
    ElMessage.error('加载配置失败')
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
    ElMessage.success('保存成功')
    loadConfig()
  } catch (e) {
    console.error('保存失败', e)
    ElMessage.error('保存失败')
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
