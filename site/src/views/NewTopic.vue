<template>
  <div class="max-w-5xl mx-auto px-4 py-8">
    <div v-if="!configStore.state.allow_post" class="bg-white rounded-2xl shadow-sm p-12 text-center">
      <svg class="w-20 h-20 text-gray-300 mx-auto mb-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>
      </svg>
      <h2 class="text-2xl font-bold text-gray-900 mb-3">发帖功能已关闭</h2>
      <p class="text-gray-500 text-lg">管理员暂时关闭了发帖功能，请稍后再试。</p>
      <router-link to="/"
        class="inline-block mt-6 px-6 py-3 bg-blue-500 text-white rounded-xl hover:bg-blue-600 transition-colors">
        返回首页
      </router-link>
    </div>
    <div v-else>
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-2">发表新帖</h1>
        <p class="text-gray-500">分享你的想法和见解</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
          <div class="mb-6">
            <label class="block text-gray-800 text-sm font-semibold mb-2">标题</label>
            <input type="text" v-model="form.title"
              class="w-full px-5 py-2 text-[1rem]  border-2 border-gray-200 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-50 transition-all"
              placeholder="给你的帖子起一个吸引人的标题..." required>
          </div>

          <div class="mb-6">
            <label class="block text-gray-800 text-sm font-semibold mb-2">选择版块</label>
            <select v-model="form.forum_id"
              class="w-full px-5 py-2 text-[1rem] border-2 border-gray-200 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-50 transition-all bg-white"
              required>
              <option value="" disabled selected>请选择版块</option>
              <option v-for="forum in forums" :key="forum.id" :value="forum.id">{{ forum.name }}</option>
            </select>
          </div>

          <div class="mb-6">
            <label class="block text-gray-800 text-sm font-semibold mb-3">话题标签 <span
                class="text-gray-400 font-normal">（可选，最多3个）</span></label>
            <div class="relative">
              <div class="flex flex-wrap gap-2 mb-3">
                <span v-for="(tag, index) in selectedTags" :key="index"
                  class="inline-flex items-center px-4 py-2 bg-blue-50 text-blue-700 rounded-full text-sm font-medium">
                  #{{ tag }}
                  <button type="button" @click="removeTag(index)"
                    class="ml-2 text-blue-400 hover:text-blue-600 transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12">
                      </path>
                    </svg>
                  </button>
                </span>
              </div>
              <input type="text" v-model="tagInput" @input="searchTags" @keydown.enter.prevent="addTag"
                @keydown.down="navigateSuggestion(1)" @keydown.up="navigateSuggestion(-1)"
                @keydown.escape="showSuggestions = false"
                class="w-full px-5 py-2 border-2 border-gray-200 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-50 transition-all"
                placeholder="输入话题名称，按回车添加" :disabled="selectedTags.length >= 3">
              <div v-if="showSuggestions && suggestions.length > 0"
                class="absolute z-20 w-full mt-2 bg-white border-2 border-gray-200 rounded-xl shadow-xl max-h-56 overflow-y-auto">
                <div v-for="(suggestion, index) in suggestions" :key="suggestion.id"
                  @click="selectSuggestion(suggestion.name)"
                  :class="['px-5 py-3 cursor-pointer transition-colors', index === suggestionIndex ? 'bg-blue-50' : 'hover:bg-gray-50']">
                  <span class="font-medium text-gray-800">#{{ suggestion.name }}</span>
                  <span class="text-xs text-gray-400 ml-3">{{ suggestion.usage_count }}次使用</span>
                </div>
              </div>
            </div>
            <p class="text-xs text-gray-400 mt-2">话题由用户自由创建，2-20个字符</p>
          </div>
        </div>

        <div v-if="configStore.state.allow_poll" class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
          <div class="flex items-center justify-between mb-4">
            <label class="text-gray-800 text-sm font-semibold">添加投票</label>
            <button type="button" @click="showPoll = !showPoll"
              :class="['relative inline-flex h-6 w-11 items-center rounded-full transition-colors', showPoll ? 'bg-blue-500' : 'bg-gray-200']">
              <span :class="['inline-block h-4 w-4 transform rounded-full bg-white transition-transform', showPoll ? 'translate-x-6' : 'translate-x-1']"></span>
            </button>
          </div>

          <div v-if="showPoll" class="space-y-4 pt-4 border-t border-gray-100">
            <div>
              <label class="block text-gray-700 text-sm font-medium mb-2">投票标题 <span class="text-gray-400 font-normal">（可选，默认使用帖子标题）</span></label>
              <input type="text" v-model="pollForm.title"
                class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:border-blue-500"
                placeholder="输入投票标题">
            </div>

            <div>
              <label class="block text-gray-700 text-sm font-medium mb-2">投票选项 <span class="text-gray-400 font-normal">（2-10个）</span></label>
              <div class="space-y-2">
                <div v-for="(option, index) in pollForm.options" :key="index" class="flex gap-2">
                  <input type="text" v-model="pollForm.options[index]"
                    class="flex-1 px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:border-blue-500"
                    :placeholder="'选项 ' + (index + 1)">
                  <button type="button" @click="removePollOption(index)" v-if="pollForm.options.length > 2"
                    class="px-3 py-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                  </button>
                </div>
              </div>
              <button type="button" @click="addPollOption" v-if="pollForm.options.length < 10"
                class="mt-2 px-4 py-2 text-blue-500 hover:bg-blue-50 rounded-lg transition-colors text-sm font-medium">
                + 添加选项
              </button>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">投票类型</label>
                <select v-model="pollForm.poll_type"
                  class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:border-blue-500 bg-white">
                  <option value="single">单选</option>
                  <option value="multiple">多选</option>
                </select>
              </div>
              <div v-if="pollForm.poll_type === 'multiple'">
                <label class="block text-gray-700 text-sm font-medium mb-2">最多可选</label>
                <input type="number" v-model.number="pollForm.max_choices" min="2" :max="pollForm.options.length"
                  class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:border-blue-500">
              </div>
            </div>

            <div>
              <label class="block text-gray-700 text-sm font-medium mb-2">截止时间 <span class="text-gray-400 font-normal">（可选，不设置则永久有效）</span></label>
              <input type="datetime-local" v-model="pollForm.end_time"
                class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:border-blue-500">
            </div>
          </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
          <div class="px-6 py-4 border-b border-gray-100 bg-gray-50">
            <div class="flex items-center justify-between">
              <label class="text-gray-800 text-sm font-semibold">帖子内容</label>
              <div class="flex flex-wrap items-center gap-3 mt-2">
                <div class="text-xs text-gray-500 flex items-center gap-4">
                  <span class="inline-flex items-center gap-1">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z">
                      </path>
                    </svg>
                    支持图片和视频
                  </span>
                </div>
                <button type="button" @click="openVideoUpload"
                  class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-white bg-gradient-to-r from-purple-500 to-purple-600 rounded-lg hover:from-purple-600 hover:to-purple-700 transition-all shadow-sm">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z">
                    </path>
                  </svg>
                  上传视频
                </button>
              </div>
            </div>
          </div>

          <input type="file" ref="videoInputRef" accept="video/*" class="hidden" @change="handleVideoFileSelect" />

          <Editor ref="editorRef" :value="form.content" @change="handleEditorChange" :plugins="plugins" :locale="zhHans"
            placeholder="开始编写你的精彩内容..." :upload-images="handleUploadImage" :upload-files="handleUploadVideo"
            :mode="editorMode" :sanitize="sanitizeSchema" class="bytemd-editor" />

          <div class="px-6 py-4 bg-gray-50 border-t border-gray-100">
            <div class="flex items-center gap-4 text-sm text-gray-500">
              <div class="flex items-center gap-1.5">
                <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                支持 Markdown 语法
              </div>
              <div class="w-px h-4 bg-gray-300"></div>
              <div class="flex items-center gap-1.5">
                <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
                </svg>
                拖拽/粘贴图片上传
              </div>
              <div class="w-px h-4 bg-gray-300"></div>
              <div class="flex items-center gap-1.5">
                <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z">
                  </path>
                </svg>
                右侧实时预览
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-4 pt-4">
          <button type="button" @click="$router.back()"
            class="px-8 py-3.5 border-2 border-gray-200 text-gray-700 rounded-xl hover:bg-gray-50 hover:border-gray-300 font-semibold transition-all">
            取消
          </button>
          <button type="submit" :disabled="submitting"
            class="px-10 py-3.5 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-xl hover:from-blue-600 hover:to-blue-700 disabled:opacity-50 disabled:cursor-not-allowed font-semibold shadow-lg shadow-blue-500/25 transition-all hover:shadow-xl hover:shadow-blue-500/30">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none">
                </circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              发布中...
            </span>
            <span v-else>立即发布</span>
          </button>
        </div>
      </form>
    </div>

    <!-- 视频上传进度对话框 -->
    <el-dialog v-model="videoUploading" title="视频上传中" width="400px" :close-on-click-modal="false" :show-close="false">
      <div class="py-4">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm text-gray-600">上传进度</span>
          <span class="text-sm font-medium text-blue-600">{{ videoUploadProgress }}%</span>
        </div>
        <el-progress :percentage="videoUploadProgress" :stroke-width="12" />
        <p class="text-xs text-gray-400 mt-3 text-center">请耐心等待，上传完成后将自动插入视频</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useConfigStore } from '@/stores/config'
import { Editor } from '@bytemd/vue-next'
import gfm from '@bytemd/plugin-gfm'
import highlight from '@bytemd/plugin-highlight'
import mediumZoom from '@bytemd/plugin-medium-zoom'
import math from '@bytemd/plugin-math'
import 'bytemd/dist/index.css'
import 'highlight.js/styles/github.css'
import 'katex/dist/katex.css'
import api, { pollApi } from '@/api'
import { ElMessage } from 'element-plus'
import { uploadImage, uploadVideo } from '@/utils/upload'

const plugins = [
  gfm(),
  highlight(),
  mediumZoom(),
  math(),
  createResizePlugin()
]

function createResizePlugin() {
  return {
    remark: (processor) => processor,
    rehype: (processor) => processor,
    viewerEffect: ({ markdownBody }) => {
      const setupResize = () => {
        const elements = markdownBody.querySelectorAll('img, video')
        elements.forEach(element => {
          if (!element.parentNode.classList?.contains('resize-container')) {
            makeResizable(element)
          }
        })
      }

      const makeResizable = (element) => {
        // 创建容器来包装元素和手柄
        const container = document.createElement('div')
        container.className = 'resize-container'
        container.style.position = 'relative'
        container.style.display = 'inline-block'
        container.style.border = '2px solid #1890ff'
        container.style.borderRadius = '8px'
        container.style.padding = '2px'
        container.style.transition = 'all 0.2s ease'
        
        // 移动元素到容器中
        if (element.parentNode) {
          element.parentNode.insertBefore(container, element)
          container.appendChild(element)
        }
        
        // 创建尺寸显示
        const sizeDisplay = document.createElement('div')
        sizeDisplay.className = 'size-display'
        sizeDisplay.style.position = 'absolute'
        sizeDisplay.style.top = '50%'
        sizeDisplay.style.left = '50%'
        sizeDisplay.style.transform = 'translate(-50%, -50%)'
        sizeDisplay.style.background = 'rgba(0, 0, 0, 0.7)'
        sizeDisplay.style.color = 'white'
        sizeDisplay.style.padding = '4px 8px'
        sizeDisplay.style.borderRadius = '4px'
        sizeDisplay.style.fontSize = '12px'
        sizeDisplay.style.zIndex = '20'
        sizeDisplay.style.pointerEvents = 'none'
        sizeDisplay.textContent = `${element.offsetWidth} x ${element.offsetHeight}`
        container.appendChild(sizeDisplay)
        
        // 创建四个角的调整点
        const corners = ['top-left', 'top-right', 'bottom-left', 'bottom-right']
        const cursors = ['nwse-resize', 'nesw-resize', 'nesw-resize', 'nwse-resize']
        
        corners.forEach((corner, index) => {
          const handle = document.createElement('div')
          handle.className = `resize-handle ${corner}`
          handle.style.position = 'absolute'
          handle.style.width = '12px'
          handle.style.height = '12px'
          handle.style.background = '#1890ff'
          handle.style.border = '2px solid white'
          handle.style.borderRadius = '50%'
          handle.style.cursor = cursors[index]
          handle.style.zIndex = '20'
          handle.style.boxShadow = '0 2px 4px rgba(0, 0, 0, 0.3)'
          
          // 设置位置
          if (corner.includes('top')) {
            handle.style.top = '-6px'
          } else {
            handle.style.bottom = '-6px'
          }
          if (corner.includes('left')) {
            handle.style.left = '-6px'
          } else {
            handle.style.right = '-6px'
          }
          
          container.appendChild(handle)
          
          // 添加拖拽功能
          let isResizing = false
          let startX, startY, startWidth, startHeight
          
          handle.addEventListener('mousedown', (e) => {
            e.stopPropagation()
            isResizing = true
            startX = e.clientX
            startY = e.clientY
            startWidth = element.offsetWidth
            startHeight = element.offsetHeight
            
            document.body.style.cursor = cursors[index]
            document.body.style.userSelect = 'none'
          })
          
          document.addEventListener('mousemove', (e) => {
            if (!isResizing) return
            
            const width = startWidth + (e.clientX - startX) * (corner.includes('left') ? -1 : 1)
            const height = startHeight + (e.clientY - startY) * (corner.includes('top') ? -1 : 1)
            
            if (width > 100 && height > 100) {
              element.style.width = `${width}px`
              element.style.height = `${height}px`
              sizeDisplay.textContent = `${Math.round(width)} x ${Math.round(height)}`
            }
          })
          
          document.addEventListener('mouseup', () => {
            if (isResizing) {
              isResizing = false
              document.body.style.cursor = ''
              document.body.style.userSelect = ''
              
              updateMarkdownSize(element)
            }
          })
        })
      }

      const updateMarkdownSize = (element) => {
        const width = Math.round(element.offsetWidth)
        const height = Math.round(element.offsetHeight)

        let currentContent = form.value.content
        let updatedContent = currentContent

        if (element.tagName === 'IMG') {
          const src = element.getAttribute('src')
          const alt = element.getAttribute('alt') || ''
          
          console.log('调整图片大小:', { src, alt, width, height })

          // 使用正则表达式匹配 Markdown 图片语法 ![alt](src) 或 ![alt](src "title")
          const mdRegex = new RegExp(`!\\[([^\\]]*)\\]\\(${src.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}(?:\\s+"[^"]*")?\\)`, 'i')
          const mdMatch = currentContent.match(mdRegex)
          
          if (mdMatch) {
            console.log('找到 Markdown 图片:', mdMatch[0])
            // 替换 Markdown 为 HTML img 标签
            const imgTag = `<img src="${src}" alt="${alt}" width="${width}" height="${height}" style="max-width: 100%; border-radius: 8px; display: block; margin: 1rem 0;">`
            updatedContent = currentContent.replace(mdRegex, imgTag)
          } else {
            // 如果没有找到 Markdown，尝试匹配现有的 <img> 标签
            const imgRegex = new RegExp(`<img\\s+[^>]*src=["']${src.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}["'][^>]*>`, 'i')
            const imgMatch = currentContent.match(imgRegex)
            
            if (imgMatch) {
              console.log('找到 img 标签:', imgMatch[0])
              const imgTag = `<img src="${src}" alt="${alt}" width="${width}" height="${height}" style="max-width: 100%; border-radius: 8px; display: block; margin: 1rem 0;">`
              updatedContent = currentContent.replace(imgRegex, imgTag)
            } else {
              console.warn('未找到匹配的图片')
            }
          }
        } else if (element.tagName === 'VIDEO') {
          const src = element.getAttribute('src')
          
          console.log('调整视频大小:', { src, width, height })

          // 匹配 <video> 标签
          const videoRegex = new RegExp(`<video\\s+[^>]*src=["']${src.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}["'][^>]*>\\s*</video>`, 'i')
          const videoMatch = currentContent.match(videoRegex)
          
          if (videoMatch) {
            console.log('找到 video 标签:', videoMatch[0])
            const videoTag = `<video src="${src}" width="${width}" height="${height}" controls style="max-width: 100%; border-radius: 8px; display: block; margin: 1rem 0;"></video>`
            updatedContent = currentContent.replace(videoRegex, videoTag)
          } else {
            console.warn('未找到匹配的视频')
          }
        }

        if (updatedContent !== currentContent) {
          form.value.content = updatedContent
          console.log('✅ Markdown 已更新')
          console.log('新内容:', updatedContent)
        } else {
          console.warn('❌ 未更新内容')
        }
      }

      setupResize()

      const observer = new MutationObserver(setupResize)
      observer.observe(markdownBody, { childList: true, subtree: true })

      return () => observer.disconnect()
    }
  }
}

function sanitizeSchema(schema) {
  const newSchema = { ...schema }

  newSchema.tagNames = [...(newSchema.tagNames || []), 'video']
  newSchema.attributes = {
    ...(newSchema.attributes || {}),
    video: ['src', 'controls', 'style', 'width', 'height', 'autoplay', 'loop', 'muted', 'poster']
  }

  return newSchema
}

const zhHans = {
  'placeholder': '开始编写你的精彩内容...',
  'uploadError': '上传失败',
  'bold': '加粗',
  'italic': '斜体',
  'strike': '删除线',
  'link': '链接',
  'quote': '引用',
  'code': '代码',
  'image': '图片',
  'file': '文件',
  'table': '表格',
  'ordered-list': '有序列表',
  'unordered-list': '无序列表',
  'task-list': '任务列表',
  'heading-1': '一级标题',
  'heading-2': '二级标题',
  'heading-3': '三级标题',
  'heading-4': '四级标题',
  'heading-5': '五级标题',
  'heading-6': '六级标题',
}

const router = useRouter()
const configStore = useConfigStore()

const form = ref({
  title: '',
  content: '',
  forum_id: '',
  tag_names: []
})
const forums = ref([])
const selectedTags = ref([])
const tagInput = ref('')
const suggestions = ref([])
const showSuggestions = ref(false)
const suggestionIndex = ref(-1)
const submitting = ref(false)
const editorMode = ref('split')
const editorRef = ref(null)
const videoInputRef = ref(null)
let searchTimeout = null

const showPoll = ref(false)
const videoUploading = ref(false)
const videoUploadProgress = ref(0)
const pollForm = ref({
  title: '',
  options: ['', ''],
  poll_type: 'single',
  max_choices: 1,
  end_time: ''
})

function addPollOption() {
  if (pollForm.value.options.length < 10) {
    pollForm.value.options.push('')
  }
}

function removePollOption(index) {
  if (pollForm.value.options.length > 2) {
    pollForm.value.options.splice(index, 1)
  }
}

function openVideoUpload() {
  videoInputRef.value?.click()
}

async function handleVideoFileSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  // 检查文件大小（50MB）
  if (file.size > 50 * 1024 * 1024) {
    ElMessage.warning('视频大小超过50MB，建议自行压缩后上传')
    event.target.value = ''
    return
  }

  videoUploading.value = true
  videoUploadProgress.value = 0

  try {
    const url = await uploadVideo(file, {
      onInstant: () => {
        ElMessage.success('视频上传成功（秒传）')
        videoUploadProgress.value = 100
      },
      onProgress: (percent) => {
        videoUploadProgress.value = percent
      }
    })

    ElMessage.success('视频上传成功')
    const videoMarkdown = `\n<video src="${url}" controls style="max-width: 100%; border-radius: 8px; display: block; margin: 1rem 0;"></video>\n`
    form.value.content += videoMarkdown
  } catch (error) {
    console.error('Video upload error:', error)
    if (error.message && error.message.startsWith('FILE_TOO_LARGE')) {
      ElMessage.warning('视频大小超过50MB，建议自行压缩后上传')
    } else {
      ElMessage.error('视频上传失败')
    }
  } finally {
    videoUploading.value = false
    event.target.value = ''
  }
}

function handleEditorChange(value) {
  form.value.content = value
}

async function handleUploadImage(files) {
  const file = files[0]
  if (!file) return []

  try {
    const url = await uploadImage(file, {
      dir: 'images',
      onInstant: () => ElMessage.success('图片上传成功（秒传）')
    })

    ElMessage.success('图片上传成功')
    
    // 直接插入 HTML img 标签而不是 Markdown 语法
    const imgTag = `<img src="${url}" alt="${file.name}" style="max-width: 100%; border-radius: 8px; display: block; margin: 1rem 0;">`
    form.value.content += `\n${imgTag}\n`
    
    // 返回空数组，因为我们已经手动插入了内容
    return []
  } catch (error) {
    console.error('Image upload error:', error)
    ElMessage.error('图片上传失败')
    return []
  }
}

async function handleUploadVideo(files) {
  const file = files[0]
  if (!file) return []

  // 检查文件大小（50MB）
  if (file.size > 50 * 1024 * 1024) {
    ElMessage.warning('视频大小超过50MB，建议自行压缩后上传')
    return []
  }

  videoUploading.value = true
  videoUploadProgress.value = 0

  try {
    const url = await uploadVideo(file, {
      onInstant: () => {
        ElMessage.success('视频上传成功（秒传）')
        videoUploadProgress.value = 100
      },
      onProgress: (percent) => {
        videoUploadProgress.value = percent
      }
    })

    ElMessage.success('视频上传成功')
    return [{
      url: url,
      alt: file.name,
      title: file.name
    }]
  } catch (error) {
    console.error('Video upload error:', error)
    if (error.message && error.message.startsWith('FILE_TOO_LARGE')) {
      ElMessage.warning('视频大小超过50MB，建议自行压缩后上传')
    } else {
      ElMessage.error('视频上传失败')
    }
    return []
  } finally {
    videoUploading.value = false
  }
}

async function loadForums() {
  try {
    const res = await api.get('/forums')
    forums.value = res || []
  } catch (e) {
    console.error(e)
  }
}

function searchTags() {
  if (searchTimeout) clearTimeout(searchTimeout)

  if (!tagInput.value.trim()) {
    suggestions.value = []
    showSuggestions.value = false
    return
  }

  searchTimeout = setTimeout(async () => {
    try {
      const res = await api.get('/tags/search', { params: { q: tagInput.value.trim() } })
      suggestions.value = res || []
      showSuggestions.value = suggestions.value.length > 0
      suggestionIndex.value = -1
    } catch (e) {
      console.error(e)
    }
  }, 300)
}

function addTag() {
  const tagName = tagInput.value.trim()
  if (!tagName) return

  if (tagName.length < 2 || tagName.length > 20) {
    return
  }

  if (selectedTags.value.includes(tagName)) {
    tagInput.value = ''
    showSuggestions.value = false
    return
  }

  if (selectedTags.value.length >= 3) {
    return
  }

  selectedTags.value.push(tagName)
  form.value.tag_names = selectedTags.value
  tagInput.value = ''
  showSuggestions.value = false
}

function selectSuggestion(name) {
  if (selectedTags.value.includes(name)) {
    tagInput.value = ''
    showSuggestions.value = false
    return
  }

  if (selectedTags.value.length >= 3) return

  selectedTags.value.push(name)
  form.value.tag_names = selectedTags.value
  tagInput.value = ''
  showSuggestions.value = false
}

function removeTag(index) {
  selectedTags.value.splice(index, 1)
  form.value.tag_names = selectedTags.value
}

function navigateSuggestion(direction) {
  if (!showSuggestions.value || suggestions.value.length === 0) return

  const newIndex = suggestionIndex.value + direction
  if (newIndex >= 0 && newIndex < suggestions.value.length) {
    suggestionIndex.value = newIndex
  }
}

async function handleSubmit() {
  if (!form.value.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }

  if (!form.value.forum_id) {
    ElMessage.warning('请选择版块')
    return
  }

  if (!form.value.content || form.value.content.trim() === '') {
    console.log('content:', form.value.content)
    ElMessage.warning('请输入内容')
    return
  }

  if (showPoll.value) {
    const validOptions = pollForm.value.options.filter(opt => opt.trim() !== '')
    if (validOptions.length < 2) {
      ElMessage.warning('请至少填写2个投票选项')
      return
    }
  }

  submitting.value = true

  try {
    const res = await api.post('/topics', form.value)
    
    if (showPoll.value) {
      try {
        const validOptions = pollForm.value.options.filter(opt => opt.trim() !== '')
        const pollData = {
          topic_id: res.id,
          title: pollForm.value.title || form.value.title,
          poll_type: pollForm.value.poll_type,
          max_choices: pollForm.value.poll_type === 'multiple' ? pollForm.value.max_choices : 1,
          options: validOptions.map(text => ({ text })),
        }
        
        if (pollForm.value.end_time) {
          pollData.end_time = new Date(pollForm.value.end_time).toISOString()
        }
        
        await pollApi.createPoll(pollData)
      } catch (pollError) {
        console.error('创建投票失败:', pollError)
        ElMessage.warning('帖子发布成功，但投票创建失败')
      }
    }
    
    ElMessage.success('发布成功')
    setTimeout(() => {
      router.push(`/topic/${res.id}`)
    }, 1500)
  } catch (error) {
    console.error(error)
    ElMessage.error(error.response?.data?.message || '发布失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadForums()
  // 强制设置编辑器模式为双栏
  setTimeout(() => {
    editorMode.value = 'split'
  }, 100)
})

// 调试 content 变化
watch(() => form.value.content, (val) => {
  console.log('content changed:', val)
})
</script>

<style scoped>
.bytemd-editor {
  min-height: 650px;
  width: 100%;
}

.bytemd-editor :deep(.bytemd) {
  min-height: 650px;
  border: none;
}

.bytemd-editor :deep(.bytemd-toolbar) {
  border-bottom: 1px solid #e5e7eb;
}

.bytemd-editor :deep(.bytemd-editor-left) {
  min-height: 580px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  line-height: 1.6;
  color: #111827;
}

.bytemd-editor :deep(.bytemd-editor-right) {
  min-height: 580px;
}

.bytemd-editor :deep(.bytemd-preview) {
  padding: 1rem 1.5rem;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: 16px;
  line-height: 1.7;
  color: #111827;
}

.bytemd-editor :deep(.markdown-body > *:first-child) {
  margin-top: -8px !important;
}

.bytemd-editor :deep(.bytemd-preview h1),
.bytemd-editor :deep(.bytemd-preview h2),
.bytemd-editor :deep(.bytemd-preview h3),
.bytemd-editor :deep(.bytemd-preview h4),
.bytemd-editor :deep(.bytemd-preview h5),
.bytemd-editor :deep(.bytemd-preview h6) {
  color: #111827;
  font-weight: 700;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
}

.bytemd-editor :deep(.bytemd-preview p) {
  color: #111827;
  margin: 0.75rem 0;
}

.bytemd-editor :deep(.bytemd-preview strong) {
  color: #111827;
  font-weight: 700;
}

.bytemd-editor :deep(.bytemd-toolbar-btn[title="GitHub"]),
.bytemd-editor :deep(.bytemd-toolbar-btn[title="帮助"]),
.bytemd-editor :deep(.bytemd-toolbar-btn[title="Help"]) {
  display: none;
}

.bytemd-editor :deep(.bytemd-preview video) {
  max-width: 100%;
  border-radius: 8px;
  margin: 1rem 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.bytemd-editor :deep(.resizable) {
  position: relative;
  display: inline-block;
}

.bytemd-editor :deep(.resize-handle) {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 16px;
  height: 16px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 0 0 8px 0;
  cursor: se-resize;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.bytemd-editor :deep(.resize-handle:hover) {
  opacity: 1;
  background: rgba(0, 0, 0, 0.8);
}

.bytemd-editor :deep(.resizable:hover .resize-handle) {
  opacity: 0.8;
}
</style>
