/**
 * 统一的上传工具模块
 * 提供文件上传、秒传、图片压缩等功能
 * 所有图片上传前都会强制压缩到 500KB 以下
 */

import axios from 'axios'
import { ElMessage } from 'element-plus'

// 默认压缩到 500KB
const DEFAULT_MAX_SIZE = 500 * 1024

// 获取 token
function getToken() {
  return localStorage.getItem('admin_token')
}

// 创建 axios 实例
const uploadApi = axios.create({
  baseURL: '/api/v1',
  timeout: 60000
})

uploadApi.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

uploadApi.interceptors.response.use((response) => {
  const res = response.data
  if (res && res.code === 0 && res.data && res.data.url) {
    return res.data
  }
  return response
})

/**
 * 计算文件的 SHA-256 哈希值
 * @param {File} file 文件对象
 * @returns {Promise<string>} 十六进制哈希字符串
 */
export async function calculateFileHash(file) {
  const buffer = await file.arrayBuffer()
  const hashBuffer = await crypto.subtle.digest('SHA-256', buffer)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
}

/**
 * 压缩图片到指定大小以下（使用 canvas）
 * @param {File} file 原始文件
 * @param {number} maxSize 最大文件大小（字节），默认 500KB
 * @param {number} maxWidth 最大宽度
 * @param {number} maxHeight 最大高度
 * @returns {Promise<File>} 压缩后的文件
 */
export function compressImageToSize(file, maxSize = DEFAULT_MAX_SIZE, maxWidth = 1920, maxHeight = 1080) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        // 计算初始尺寸
        let width = img.width
        let height = img.height

        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width = Math.round(width * ratio)
          height = Math.round(height * ratio)
        }

        // 创建 canvas
        const canvas = document.createElement('canvas')
        canvas.width = width
        canvas.height = height

        const ctx = canvas.getContext('2d')
        ctx.drawImage(img, 0, 0, width, height)

        // 递归压缩直到文件大小符合要求
        const compress = (quality) => {
          if (quality < 0.1) {
            // 质量已经很低了，直接返回
            const blob = dataURLtoBlob(canvas.toDataURL('image/jpeg', 0.1))
            resolve(new File([blob], file.name.replace(/\.[^.]+$/, '.jpg'), { type: 'image/jpeg' }))
            return
          }

          const dataURL = canvas.toDataURL('image/jpeg', quality)
          const size = Math.round((dataURL.length - dataURL.indexOf(',') - 1) * 0.75)

          if (size <= maxSize) {
            const blob = dataURLtoBlob(dataURL)
            resolve(new File([blob], file.name.replace(/\.[^.]+$/, '.jpg'), { type: 'image/jpeg' }))
          } else {
            // 继续压缩，每次减少 10% 质量
            compress(quality - 0.1)
          }
        }

        // 从 0.9 开始压缩
        compress(0.9)
      }
      img.onerror = reject
      img.src = e.target.result
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

/**
 * 将 dataURL 转换为 Blob
 */
function dataURLtoBlob(dataURL) {
  const arr = dataURL.split(',')
  const mime = arr[0].match(/:(.*?);/)[1]
  const bstr = atob(arr[1])
  let n = bstr.length
  const u8arr = new Uint8Array(n)
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n)
  }
  return new Blob([u8arr], { type: mime })
}

/**
 * 检查文件是否已存在（秒传）
 * @param {string} filename 文件名
 * @param {string} contentHash 文件内容哈希
 * @returns {Promise<{exists: boolean, url?: string, key?: string}>}
 */
export async function checkFileExists(filename, contentHash) {
  try {
    const res = await uploadApi.get('/upload/check', {
      params: { filename, content_hash: contentHash }
    })
    return res
  } catch (e) {
    console.error('检查文件存在失败:', e)
    return { exists: false }
  }
}

/**
 * 上传文件
 * @param {File} file 文件对象
 * @param {Object} options 配置选项
 * @param {string} options.dir 存储目录
 * @param {string} options.contentHash 文件内容哈希（用于秒传）
 * @returns {Promise<string>} 文件URL
 */
export async function uploadFile(file, options = {}) {
  const { dir = '', contentHash = '' } = options
  const formData = new FormData()
  formData.append('file', file)

  // 构建查询参数
  const params = {}
  if (dir) params.dir = dir
  if (contentHash) params.content_hash = contentHash

  const res = await uploadApi.post('/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    params
  })

  if (res && res.url) {
    return res.url
  }
  throw new Error('上传失败')
}

/**
 * 通用图片上传（自动压缩到 500KB 以下，带秒传）
 * @param {File} file 原始文件
 * @param {Object} options 配置选项
 * @param {string} options.dir 存储目录
 * @param {number} options.maxSize 最大文件大小（字节）
 * @param {Function} options.onInstant 秒传成功回调
 * @returns {Promise<string>} 文件URL
 */
export async function uploadImage(file, options = {}) {
  const { dir = '', maxSize = DEFAULT_MAX_SIZE, onInstant } = options

  // 1. 压缩图片到指定大小
  const compressedFile = await compressImageToSize(file, maxSize)
  console.log(`图片压缩完成: ${file.name} (${file.size}) -> ${compressedFile.name} (${compressedFile.size})`)

  // 2. 计算压缩后文件的 hash
  const contentHash = await calculateFileHash(compressedFile)

  // 3. 检查秒传
  const checkRes = await checkFileExists(compressedFile.name, contentHash)
  if (checkRes.exists && checkRes.url) {
    console.log('秒传成功:', file.name)
    onInstant?.(checkRes.url)
    return checkRes.url
  }

  // 4. 上传
  return uploadFile(compressedFile, { dir, contentHash })
}

/**
 * 压缩图片并上传（带秒传，用于需要指定尺寸的场合）
 * @param {File} file 原始文件
 * @param {number} maxWidth 最大宽度
 * @param {number} maxHeight 最大高度
 * @param {Object} options 配置选项
 * @param {string} options.dir 存储目录
 * @param {string} options.mimeType 输出格式
 * @param {number} options.maxSize 最大文件大小
 * @param {Function} options.onInstant 秒传成功回调
 * @returns {Promise<string>} 文件URL
 */
export async function compressAndUpload(file, maxWidth, maxHeight, options = {}) {
  const { dir = '', mimeType = 'image/png', maxSize = DEFAULT_MAX_SIZE, onInstant } = options

  // 1. 先压缩到指定尺寸
  const compressedFile = await compressImageToSize(file, maxSize, maxWidth, maxHeight)
  console.log(`图片压缩完成: ${file.name} -> ${compressedFile.name} (${compressedFile.size})`)

  // 2. 计算 hash
  const contentHash = await calculateFileHash(compressedFile)

  // 3. 检查秒传
  const checkRes = await checkFileExists(compressedFile.name, contentHash)
  if (checkRes.exists && checkRes.url) {
    console.log('秒传成功:', file.name)
    onInstant?.(checkRes.url)
    return checkRes.url
  }

  // 4. 上传
  return uploadFile(compressedFile, { dir, contentHash })
}
