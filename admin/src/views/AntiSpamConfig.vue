<template>
  <div class="antispam-page">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
        </div>
      </template>

    <el-tabs v-model="activeTab" class="mb-6">
      <el-tab-pane :label="t('antispam.rateLimit')" name="rate">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">{{ t('antispam.rateLimitTitle') }}</h3>
          <el-form :model="rateConfig" label-width="200px">
            <el-form-item :label="t('antispam.topicMinInterval')">
              <el-input-number v-model="rateConfig.topic_min_interval" :min="10" :max="600" />
            </el-form-item>
            <el-form-item :label="t('antispam.commentMinInterval')">
              <el-input-number v-model="rateConfig.comment_min_interval" :min="5" :max="300" />
            </el-form-item>
            <el-form-item :label="t('antispam.maxTopicsPerDay')">
              <el-input-number v-model="rateConfig.max_topics_per_day" :min="1" :max="100" />
            </el-form-item>
            <el-form-item :label="t('antispam.maxCommentsPerDay')">
              <el-input-number v-model="rateConfig.max_comments_per_day" :min="1" :max="200" />
            </el-form-item>
            <el-divider />
            <h4 class="text-md font-medium mb-3">{{ t('antispam.newUserLimit') }}</h4>
            <el-form-item :label="t('antispam.newUserHours2')">
              <el-input-number v-model="rateConfig.new_user_hours" :min="1" :max="168" />
            </el-form-item>
            <el-form-item :label="t('antispam.newUserTopicLimit')">
              <el-input-number v-model="rateConfig.new_user_max_topics_per_day" :min="1" :max="20" />
            </el-form-item>
            <el-form-item :label="t('antispam.newUserCommentLimit')">
              <el-input-number v-model="rateConfig.new_user_max_comments_per_day" :min="1" :max="50" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveRateConfig">{{ t('antispam.saveConfig') }}</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="t('antispam.contentQuality')" name="quality">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">{{ t('antispam.qualityDetect') }}</h3>
          <el-form :model="qualityConfig" label-width="200px">
            <el-form-item :label="t('antispam.minContentLength')">
              <el-input-number v-model="qualityConfig.min_content_length" :min="5" :max="100" />
              <span class="ml-2 text-gray-500">{{ t('antispam.charCount') }}</span>
            </el-form-item>
            <el-form-item :label="t('antispam.repeatCharThreshold')">
              <el-input-number v-model="qualityConfig.repeat_char_threshold" :min="3" :max="20" />
              <span class="ml-2 text-gray-500">{{ t('antispam.repeatCharDesc') }}</span>
            </el-form-item>
            <el-form-item :label="t('antispam.similarityThreshold')">
              <el-slider v-model="qualityConfig.similarity_threshold" :min="0.5" :max="1" :step="0.1" show-input />
              <span class="ml-2 text-gray-500">{{ t('antispam.similarityDesc') }}</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveQualityConfig">{{ t('antispam.saveConfig') }}</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="t('antispam.reputation')" name="reputation">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">{{ t('antispam.reputationConfig') }}</h3>
          <el-form :model="reputationConfig" label-width="200px">
            <el-form-item :label="t('antispam.lowReputationThreshold')">
              <el-input-number v-model="reputationConfig.low_reputation_threshold" :min="20" :max="80" />
              <span class="ml-2 text-gray-500">{{ t('antispam.lowRepDesc') }}</span>
            </el-form-item>
            <el-form-item :label="t('antispam.banReputationThreshold')">
              <el-input-number v-model="reputationConfig.ban_reputation_threshold" :min="0" :max="50" />
              <span class="ml-2 text-gray-500">{{ t('antispam.banRepDesc') }}</span>
            </el-form-item>
            <el-form-item :label="t('antispam.banLowReputation')">
              <el-switch v-model="reputationConfig.ban_low_reputation" />
            </el-form-item>
            <el-divider />
            <h4 class="text-md font-medium mb-3">{{ t('antispam.hotWeight') }}</h4>
            <el-form-item :label="t('antispam.lowQualityHot')">
              <el-slider v-model="reputationConfig.low_quality_hot_multiplier" :min="0" :max="1" :step="0.1" show-input />
            </el-form-item>
            <el-form-item :label="t('antispam.lowRepHot')">
              <el-slider v-model="reputationConfig.low_reputation_hot_multiplier" :min="0" :max="1" :step="0.1" show-input />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveReputationConfig">{{ t('antispam.saveConfig') }}</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="t('antispam.spamKeywords')" name="keywords">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">{{ t('antispam.spamKeywordsTitle') }}</h3>
          <div class="flex gap-4 mb-4">
            <el-input
              v-model="newKeyword"
              :placeholder="t('antispam.addKeywordPlaceholder')"
              style="width: 300px"
              @keyup.enter="addKeyword"
            />
            <el-button type="primary" @click="addKeyword">{{ t('antispam.addKeyword') }}</el-button>
          </div>
          <el-alert
            :title="t('antispam.keywordTip')"
            type="info"
            :closable="false"
            class="mb-4"
          />
          <el-table :data="keywordsDisplay" style="width: 100%">
            <el-table-column prop="index" label="ID" width="80" />
            <el-table-column prop="keyword" :label="t('antispam.keyword')" />
            <el-table-column :label="t('common.actions')" width="120" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="deleteKeyword(row)">
                  {{ t('common.delete') }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="t('antispam.report')" name="report">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-medium mb-4">{{ t('antispam.reportAutoHandle') }}</h3>
          <el-form :model="reportConfig" label-width="200px">
            <el-form-item :label="t('antispam.autoHide')">
              <el-input-number v-model="reportConfig.report_threshold" :min="1" :max="10" />
              <span class="ml-2 text-gray-500">{{ t('antispam.autoHideDesc') }}</span>
            </el-form-item>
            <el-form-item :label="t('antispam.autoBan')">
              <el-input-number v-model="reportConfig.report_ban_threshold" :min="1" :max="20" />
              <span class="ml-2 text-gray-500">{{ t('antispam.autoBanDesc') }}</span>
            </el-form-item>
            <el-form-item :label="t('antispam.autoBanDays')">
              <el-input-number v-model="reportConfig.report_ban_days" :min="1" :max="30" />
            </el-form-item>
            <el-form-item :label="t('antispam.maxReportPerDay')">
              <el-input-number v-model="reportConfig.max_reports_per_day" :min="1" :max="50" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveReportConfig">{{ t('antispam.saveConfig') }}</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <el-tab-pane :label="t('antispam.userReputation')" name="users">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium">{{ t('antispam.userRepManage') }}</h3>
            <el-input v-model="userSearch" :placeholder="t('antispam.searchUser')" style="width: 300px" clearable>
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>

          <el-table :data="filteredUsers" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="username" :label="t('antispam.username')" width="150" />
            <el-table-column prop="nickname" :label="t('antispam.nickname')" width="150" />
            <el-table-column :label="t('antispam.reputationScore')" width="200">
              <template #default="{ row }">
                <el-progress
                  :percentage="row.reputation"
                  :color="getReputationColor(row.reputation)"
                  :stroke-width="20"
                  :text-inside="true"
                />
              </template>
            </el-table-column>
            <el-table-column :label="t('antispam.reputationLevel')" width="120">
              <template #default="{ row }">
                <el-tag :type="getReputationTagType(row.reputation)">
                  {{ getReputationLevel(row.reputation) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="t('antispam.registerTime')" width="180" />
            <el-table-column :label="t('common.actions')" fixed="right" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="showReputationDialog(row)">{{ t('antispam.adjustScore') }}</el-button>
                <el-button size="small" type="danger" @click="banUser(row)" v-if="!isBanned(row)">{{ t('antispam.ban') }}</el-button>
                <el-button size="small" type="success" @click="unbanUser(row)" v-else>{{ t('antispam.unban') }}</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="flex justify-center mt-4">
            <el-pagination
              v-model:current-page="userPage"
              :page-size="20"
              :total="userTotal"
              layout="prev, pager, next"
              @current-change="loadUsers"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="reputationDialogVisible" :title="t('antispam.adjustScore')" width="400px">
      <el-form :model="reputationForm" label-width="100px">
        <el-form-item :label="t('antispam.currentScore')">
          <el-input :value="currentUser?.reputation" disabled />
        </el-form-item>
        <el-form-item :label="t('antispam.adjustChange')">
          <el-input-number v-model="reputationForm.change" :min="-100" :max="100" />
        </el-form-item>
        <el-form-item :label="t('antispam.adjustReason')">
          <el-input v-model="reputationForm.reason" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reputationDialogVisible = false">{{ t('antispam.cancel') }}</el-button>
        <el-button type="primary" @click="adjustReputation">{{ t('antispam.confirm') }}</el-button>
      </template>
    </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { Shield } from 'lucide-vue-next'
import api from '@/api'

const { t } = useI18n()
const activeTab = ref('rate')
const userSearch = ref('')
const userPage = ref(1)
const userTotal = ref(0)
const users = ref([])
const reputationDialogVisible = ref(false)
const currentUser = ref(null)

const rateConfig = ref({
  topic_min_interval: 60,
  comment_min_interval: 30,
  max_topics_per_day: 10,
  max_comments_per_day: 50,
  new_user_hours: 24,
  new_user_max_topics_per_day: 3,
  new_user_max_comments_per_day: 10
})

const qualityConfig = ref({
  min_content_length: 10,
  repeat_char_threshold: 5,
  similarity_threshold: 0.8
})

const reputationConfig = ref({
  low_reputation_threshold: 60,
  ban_reputation_threshold: 20,
  ban_low_reputation: true,
  low_quality_hot_multiplier: 0.3,
  low_reputation_hot_multiplier: 0.5
})

const reportConfig = ref({
  report_threshold: 3,
  report_ban_threshold: 5,
  report_ban_days: 3,
  max_reports_per_day: 10
})

const reputationForm = ref({
  change: 0,
  reason: ''
})

const keywords = ref([])
const newKeyword = ref('')

const keywordsDisplay = computed(() => {
  return keywords.value.map((kw, idx) => ({
    index: idx + 1,
    keyword: kw
  }))
})

const filteredUsers = computed(() => {
  if (!userSearch.value) return users.value
  const search = userSearch.value.toLowerCase()
  return users.value.filter(u =>
    u.username?.toLowerCase().includes(search) ||
    u.nickname?.toLowerCase().includes(search)
  )
})

onMounted(() => {
  loadConfigs()
  loadUsers()
  loadKeywords()
})

async function loadConfigs() {
  try {
    const res = await api.get('/admin/antispam/config')
    if (res) {
      Object.keys(res).forEach(key => {
        if (rateConfig.value.hasOwnProperty(key)) {
          rateConfig.value[key] = parseInt(res[key])
        } else if (qualityConfig.value.hasOwnProperty(key)) {
          const val = res[key]
          qualityConfig.value[key] = key === 'similarity_threshold' ? parseFloat(val) : parseInt(val)
        } else if (reputationConfig.value.hasOwnProperty(key)) {
          const val = res[key]
          if (key === 'ban_low_reputation') {
            reputationConfig.value[key] = val === 'true'
          } else if (key.includes('multiplier')) {
            reputationConfig.value[key] = parseFloat(val)
          } else {
            reputationConfig.value[key] = parseInt(val)
          }
        } else if (reportConfig.value.hasOwnProperty(key)) {
          reportConfig.value[key] = parseInt(res[key])
        }
      })
    }
  } catch (e) {
    console.error('Load config failed', e)
  }
}

async function saveRateConfig() {
  await saveConfig(rateConfig.value)
}

async function saveQualityConfig() {
  await saveConfig(qualityConfig.value)
}

async function saveReputationConfig() {
  await saveConfig(reputationConfig.value)
}

async function saveReportConfig() {
  await saveConfig(reportConfig.value)
}

async function saveConfig(config) {
  try {
    await api.post('/admin/antispam/config', config)
    ElMessage.success(t('antispam.saveSuccess'))
  } catch (e) {
    ElMessage.error(t('antispam.saveFailed'))
    console.error(e)
  }
}

async function loadKeywords() {
  try {
    const res = await api.get('/admin/antispam/keywords')
    keywords.value = res.keywords || []
  } catch (e) {
    console.error('Load keywords failed', e)
  }
}

async function addKeyword() {
  if (!newKeyword.value.trim()) {
    ElMessage.warning(t('antispam.pleaseEnterKeyword'))
    return
  }
  try {
    await api.post('/admin/antispam/keywords', { keyword: newKeyword.value.trim() })
    ElMessage.success(t('antispam.addKeywordSuccess'))
    newKeyword.value = ''
    loadKeywords()
  } catch (e) {
    ElMessage.error(t('antispam.addKeywordFailed'))
    console.error(e)
  }
}

async function deleteKeyword(row) {
  try {
    await ElMessageBox.confirm(t('antispam.deleteKeywordConfirm'), t('common.warning'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await api.delete('/admin/antispam/keywords', { data: { keyword: row.keyword } })
    ElMessage.success(t('antispam.deleteKeywordSuccess'))
    loadKeywords()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(t('antispam.deleteKeywordFailed'))
      console.error(e)
    }
  }
}

async function loadUsers() {
  try {
    const res = await api.get('/admin/users', {
      params: { page: userPage.value }
    })
    users.value = res.list || []
    userTotal.value = res.total || 0
  } catch (e) {
    console.error('Load users failed', e)
  }
}

function getReputationColor(reputation) {
  if (reputation >= 80) return '#67c23a'
  if (reputation >= 60) return '#409eff'
  if (reputation >= 40) return '#e6a23c'
  if (reputation >= 20) return '#f56c6c'
  return '#909399'
}

function getReputationTagType(reputation) {
  if (reputation >= 80) return 'success'
  if (reputation >= 60) return ''
  if (reputation >= 40) return 'warning'
  return 'danger'
}

function getReputationLevel(reputation) {
  if (reputation >= 80) return t('antispam.levels.normal')
  if (reputation >= 60) return t('antispam.levels.verify')
  if (reputation >= 40) return t('antispam.levels.restricted')
  if (reputation >= 20) return t('antispam.levels.serious')
  return t('antispam.levels.banned')
}

function isBanned(user) {
  return user.reputation < 20
}

function showReputationDialog(user) {
  currentUser.value = user
  reputationForm.value = { change: 0, reason: '' }
  reputationDialogVisible.value = true
}

async function adjustReputation() {
  try {
    await api.post(`/admin/users/${currentUser.value.id}/reputation`, reputationForm.value)
    ElMessage.success(t('antispam.adjustSuccess'))
    reputationDialogVisible.value = false
    loadUsers()
  } catch (e) {
    ElMessage.error(t('antispam.adjustFailed'))
    console.error(e)
  }
}

async function banUser(user) {
  try {
    await ElMessageBox.confirm(t('antispam.banConfirm'), t('antispam.ban'), {
      confirmButtonText: t('antispam.confirm'),
      cancelButtonText: t('antispam.cancel'),
      type: 'warning'
    })
    await api.post(`/admin/users/${user.id}/ban`, { reason: t('antispam.banReason') })
    ElMessage.success(t('antispam.banSuccess'))
    loadUsers()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(t('antispam.operationFailed'))
      console.error(e)
    }
  }
}

async function unbanUser(user) {
  try {
    await api.post(`/admin/users/${user.id}/unban`)
    ElMessage.success(t('antispam.unbanSuccess'))
    loadUsers()
  } catch (e) {
    ElMessage.error(t('antispam.operationFailed'))
    console.error(e)
  }
}
</script>

<style scoped>
.antispam-page {
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
</style>
