import errorCodes from '@/i18n/errors'
import { ElMessage } from 'element-plus'

/**
 * 根据错误码获取 i18n 翻译 key
 * @param {number} code - 错误码
 * @returns {string} - i18n key
 */
export function getErrorI18nKey(code) {
  return errorCodes[code] || 'errors.serverInternal'
}

/**
 * 处理 API 错误，自动显示 i18n 翻译
 * @param {Error|Object} error - 错误对象
 * @param {Function} t - i18n t 函数
 * @param {string} fallbackMsg - 备用消息（未找到错误码时使用）
 */
export function handleApiError(error, t, fallbackMsg = 'errors.serverInternal') {
  if (error.code) {
    const key = getErrorI18nKey(error.code)
    ElMessage.error(t(key))
  } else {
    ElMessage.error(t(fallbackMsg))
  }
}
