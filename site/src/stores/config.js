import { defineStore } from 'pinia'
import { reactive } from 'vue'
import api from '@/api'

export const useConfigStore = defineStore('config', () => {
  const state = reactive({
    loading: true,
    allow_register: true,
    allow_post: true,
    allow_comment: true,
    allow_poll: true,
    credit_topic: 20,
    credit_post: 5,
    credit_signin: 10,
    site_name: '彩虹BBS',
    site_logo: '',
    site_icon: ''
  })

  function parseBool(value, defaultValue) {
    if (typeof value === 'boolean') return value
    if (value === 'true' || value === '1') return true
    if (value === 'false' || value === '0') return false
    return defaultValue
  }

  async function loadConfig() {
    state.loading = true
    try {
      const res = await api.get('/config')
      console.log('✅ 配置加载成功:', res)
      
      if (res) {
        state.allow_register = parseBool(res.allow_register, true)
        state.allow_post = parseBool(res.allow_post, true)
        state.allow_comment = parseBool(res.allow_comment, true)
        state.allow_poll = parseBool(res.allow_poll, true)
        state.credit_topic = parseInt(res.credit_topic) || 20
        state.credit_post = parseInt(res.credit_post) || 5
        state.credit_signin = parseInt(res.credit_signin) || 10
        
        if (res.site_name) state.site_name = res.site_name
        if (res.site_logo) state.site_logo = res.site_logo
        if (res.site_icon) state.site_icon = res.site_icon
        
        console.log('✅ 解析后的配置:', {
          allow_post: state.allow_post,
          allow_comment: state.allow_comment
        })
      }
    } catch (e) {
      console.error('❌ 配置加载失败:', e)
    } finally {
      state.loading = false
    }
  }

  return {
    state,
    loadConfig
  }
})
