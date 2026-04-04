# BBSGo 模板定制系统方案（方案A：组件级主题）

## 目录
- [一、系统架构设计](#一系统架构设计)
- [二、核心概念与分工](#二核心概念与分工)
- [三、推荐方案详细设计](#三推荐方案详细设计)
- [四、主题管理与定制（admin 后台）](#四主题管理与定制admin-后台)
- [五、主题应用（site 前端）](#五主题应用site-前端)
- [六、完整布局定制方案](#六完整布局定制方案)
- [七、实施路线图](#七实施路线图)
- [八、关键技术点](#八关键技术点)

---

## 一、系统架构设计

### 1.1 主题目录结构

```
bbsgo/
├── site/
│   ├── src/
│   │   ├── themes/              # 主题目录
│   │   │   ├── default/       # 默认主题
│   │   │   │   ├── components/
│   │   │   │   ├── views/
│   │   │   │   ├── styles/
│   │   │   │   ├── theme.json      # 主题配置
│   │   │   │   ├── layout.js       # 布局配置
│   │   │   │   └── index.js       # 主题入口
│   │   │   ├── dark/         # 深色主题
│   │   │   └── minimalist/  # 简约主题
│   │   └── ...
│   ├── ...
└── ...
```

### 1.2 主题配置文件 (theme.json)

```json
{
  "name": "Default Theme",
  "description": "The default theme for BBSGo",
  "version": "1.0.0",
  "author": "BBSGo Team",
  "screenshot": "screenshot.png",
  "colors": {
    "primary": "#3B82F6",
    "secondary": "#64748B",
    "background": "#FFFFFF",
    "text": "#1F2937"
  },
  "fonts": {
    "heading": "Inter, sans-serif",
    "body": "Inter, sans-serif"
  },
  "layouts": {
    "sidebar": "left",
    "sidebarWidth": "256px"
  },
  "features": {
    "showAvatars": true,
    "showBadges": true,
    "compactMode": false
  },
  "configurable": {
    "colors": ["primary", "secondary", "background", "text"],
    "layouts": ["type", "sidebarLeft", "sidebarRight", "sidebarWidth"],
    "features": ["showAvatars", "showBadges", "compactMode"]
  }
}
```

---

## 二、核心概念与分工

### 2.1 两个层面的职责

| 层面 | 位置 | 功能 | 用户角色 |
|------|------|------|----------|
| **主题管理与定制** | `admin/src/` | 安装、激活、删除主题，定制全站颜色、布局 | 站主/管理员 |
| **主题应用** | `site/src/` | 应用管理员设置的主题，全站统一风格 | 普通用户（只读） |

### 2.2 数据流转

```
站主（admin后台）
    ↓ 选择主题 + 定制
主题配置（数据库）
    ↓ 全站应用
前端渲染（site/src/）
```

---

## 三、推荐方案详细设计

### 3.1 核心实现步骤

#### 步骤 1：创建主题 Store（site 前端）

```javascript
// site/src/stores/theme.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useThemeStore = defineStore('theme', () => {
  const currentTheme = ref('default')
  const themeConfig = ref(null)
  const layoutConfig = ref(null)
  const themeComponents = ref({})
  const siteConfig = ref(null)

  async function loadSiteTheme() {
    try {
      const res = await api.get('/site/theme')
      if (res) {
        currentTheme.value = res.theme_name || 'default'
        siteConfig.value = res.config || {}
      }
    } catch (e) {
      console.error('Load site theme failed:', e)
    }
    await loadTheme(currentTheme.value)
  }

  async function loadTheme(themeName) {
    try {
      const [themeModule, layoutModule] = await Promise.all([
        import(`@/themes/${themeName}/index.js`),
        import(`@/themes/${themeName}/layout.js`)
      ])
      
      themeConfig.value = { ...themeModule.default, ...siteConfig.value }
      layoutConfig.value = layoutModule.default
      
      await loadThemeComponents(themeName)
      applyTheme(themeConfig.value)
    } catch (e) {
      console.error('Load theme failed:', e)
      if (themeName !== 'default') {
        await loadTheme('default')
      }
    }
  }

  async function loadThemeComponents(themeName) {
    const componentMap = layoutConfig.value.components || {}
    
    for (const [name, path] of Object.entries(componentMap)) {
      try {
        const module = await import(`@/themes/${themeName}/${path}`)
        themeComponents.value[name] = module.default
      } catch (e) {
        console.warn(`Failed to load component ${name}, using default`)
        try {
          const defaultModule = await import(`@/themes/default/${path}`)
          themeComponents.value[name] = defaultModule.default
        } catch (e2) {
          console.error('Failed to load default component:', e2)
        }
      }
    }
  }

  function applyTheme(config) {
    const root = document.documentElement
    Object.entries(config.colors || {}).forEach(([key, value]) => {
      root.style.setProperty(`--color-${key}`, value)
    })
  }

  return {
    currentTheme,
    themeConfig,
    layoutConfig,
    themeComponents,
    siteConfig,
    loadSiteTheme,
    loadTheme
  }
})
```

#### 步骤 2：动态组件加载

```javascript
// site/src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useThemeStore } from '@/stores/theme'

function getComponent(themeName, componentName) {
  return () => import(`@/themes/${themeName}/views/${componentName}.vue`)
    .catch(() => import(`@/themes/default/views/${componentName}.vue`))
}

const routes = [
  {
    path: '/',
    component: () => {
      const themeStore = useThemeStore()
      return getComponent(themeStore.currentTheme, 'Home')
    }
  }
]
```

#### 步骤 3：CSS 变量系统

```css
/* site/src/style.css */
:root {
  --color-primary: #3B82F6;
  --color-secondary: #64748B;
  --color-background: #FFFFFF;
  --color-text: #1F2937;
}

.btn-primary {
  background-color: var(--color-primary);
}
```

---

## 四、主题管理与定制（admin 后台）

### 4.1 主题列表与选择页面

```vue
<!-- admin/src/views/Themes.vue -->
<template>
  <div class="themes-management">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-xl font-semibold">{{ t('theme.themes') }}</h2>
      <button @click="showUploadDialog = true" class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600">
        {{ t('theme.uploadTheme') }}
      </button>
    </div>
    
    <!-- 主题列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="theme in themes" :key="theme.name"
        :class="['theme-card bg-white rounded-lg shadow-sm overflow-hidden',
          theme.isActive ? 'ring-2 ring-blue-500' : '']">
        <img v-if="theme.screenshot" :src="theme.screenshot" class="w-full h-40 object-cover">
        <div class="p-4">
          <div class="flex justify-between items-start mb-2">
            <div>
              <h3 class="font-medium text-lg">{{ theme.displayName || theme.name }}</h3>
              <p class="text-sm text-gray-500">v{{ theme.version }}</p>
            </div>
            <span v-if="theme.isActive" class="bg-green-100 text-green-700 text-xs px-2 py-1 rounded">
              {{ t('theme.active') }}
            </span>
          </div>
          <p class="text-sm text-gray-600 mb-3">{{ theme.description }}</p>
          <p class="text-xs text-gray-400 mb-3">{{ t('theme.author') }}: {{ theme.author }}</p>
          <div class="flex space-x-2">
            <button v-if="!theme.isActive" @click="activateTheme(theme)"
              class="flex-1 bg-blue-500 text-white py-1.5 rounded text-sm hover:bg-blue-600">
              {{ t('theme.activate') }}
            </button>
            <button @click="openCustomizer(theme)"
              class="flex-1 bg-gray-100 text-gray-700 py-1.5 rounded text-sm hover:bg-gray-200">
              {{ t('theme.customize') }}
            </button>
            <button v-if="theme.name !== 'default'" @click="deleteTheme(theme)"
              class="px-3 py-1.5 border border-red-300 text-red-600 rounded text-sm hover:bg-red-50">
              {{ t('theme.delete') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 上传对话框 -->
    <el-dialog v-model="showUploadDialog" :title="t('theme.uploadTheme')" width="500px">
      <el-upload
        class="theme-uploader"
        drag
        action="/api/admin/themes/upload"
        :on-success="handleUploadSuccess"
        :before-upload="beforeUpload"
        accept=".zip">
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          {{ t('theme.dragOrClick') }}
        </div>
        <template #tip>
          <div class="el-upload__tip text-xs text-gray-500">
            {{ t('theme.uploadTip') }}
          </div>
        </template>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const { t } = useI18n()

const themes = ref([])
const showUploadDialog = ref(false)
const showCustomizer = ref(false)
const currentCustomizeTheme = ref(null)

async function loadThemes() {
  try {
    const res = await api.get('/admin/themes')
    themes.value = res || []
  } catch (e) {
    console.error(e)
  }
}

async function activateTheme(theme) {
  try {
    await ElMessageBox.confirm(
      t('theme.confirmActivate', { name: theme.displayName || theme.name }),
      t('theme.activate'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    await api.post(`/admin/themes/${theme.name}/activate`)
    ElMessage.success(t('theme.activated'))
    loadThemes()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function openCustomizer(theme) {
  currentCustomizeTheme.value = theme
  showCustomizer.value = true
}

async function deleteTheme(theme) {
  try {
    await ElMessageBox.confirm(
      t('theme.confirmDelete', { name: theme.displayName || theme.name }),
      t('theme.delete'),
      { confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel'), type: 'warning' }
    )
    await api.delete(`/admin/themes/${theme.name}`)
    ElMessage.success(t('theme.deleted'))
    loadThemes()
  } catch (e) {
    if (e !== 'cancel') {
      console.error(e)
    }
  }
}

function beforeUpload(file) {
  const isZip = file.type === 'application/zip' || file.name.endsWith('.zip')
  if (!isZip) {
    ElMessage.error(t('theme.onlyZipAllowed'))
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error(t('theme.fileTooLarge'))
    return false
  }
  return true
}

function handleUploadSuccess() {
  ElMessage.success(t('theme.uploaded'))
  showUploadDialog.value = false
  loadThemes()
}

onMounted(() => {
  loadThemes()
})
</script>
```

### 4.2 主题定制器（颜色 + 布局）

```vue
<!-- admin/src/components/ThemeCustomizer.vue -->
<template>
  <el-dialog v-model="visible" :title="t('theme.customizeTheme')" width="800px" :close-on-click-modal="false">
    <el-tabs v-model="activeTab">
      <!-- 颜色定制 -->
      <el-tab-pane :label="t('theme.colors')" name="colors">
        <div class="space-y-4">
          <div v-for="colorKey in configurableColors" :key="colorKey">
            <label class="block text-sm font-medium mb-2">{{ getColorLabel(colorKey) }}</label>
            <div class="flex items-center space-x-2">
              <input type="color" v-model="customColors[colorKey]" class="w-12 h-10 rounded">
              <span class="text-sm text-gray-600 font-mono">{{ customColors[colorKey] }}</span>
              <button @click="resetColor(colorKey)" class="text-sm text-gray-500 hover:text-gray-700">
                {{ t('theme.reset') }}
              </button>
            </div>
          </div>
        </div>
      </el-tab-pane>
      
      <!-- 布局定制 -->
      <el-tab-pane :label="t('theme.layout')" name="layout">
        <div class="space-y-6">
          <!-- 布局预览 -->
          <div class="layout-preview bg-gray-100 p-4 rounded-lg">
            <div class="flex gap-2" :class="previewLayoutClass">
              <div v-if="preview.sidebarLeft" 
                class="sidebar-preview bg-blue-200 rounded p-2 text-center text-sm"
                :style="{ width: preview.sidebarWidth }">
                {{ t('theme.leftSidebar') }}
              </div>
              <div class="content-preview flex-1 bg-green-200 rounded p-2 text-center text-sm">
                {{ t('theme.content') }}
              </div>
              <div v-if="preview.sidebarRight" 
                class="sidebar-preview bg-purple-200 rounded p-2 text-center text-sm"
                :style="{ width: preview.sidebarWidth }">
                {{ t('theme.rightSidebar') }}
              </div>
            </div>
          </div>
          
          <!-- 布局类型选择 -->
          <div>
            <label class="block text-sm font-medium mb-2">{{ t('theme.layoutType') }}</label>
            <div class="grid grid-cols-3 gap-2">
              <button v-for="type in layoutTypes" :key="type.value"
                @click="setLayoutType(type.value)"
                :class="['p-3 border rounded-lg text-center',
                  preview.type === type.value ? 'border-blue-500 bg-blue-50' : 'border-gray-200']">
                <span class="text-2xl">{{ type.icon }}</span>
                <div class="text-xs mt-1">{{ type.label }}</div>
              </button>
            </div>
          </div>
          
          <!-- 侧边栏开关 -->
          <div v-if="preview.type !== 'single-column'">
            <label class="flex items-center space-x-2">
              <input type="checkbox" v-model="preview.sidebarLeft">
              <span class="text-sm">{{ t('theme.showLeftSidebar') }}</span>
            </label>
            <label v-if="preview.type === 'three-column'" class="flex items-center space-x-2 mt-2">
              <input type="checkbox" v-model="preview.sidebarRight">
              <span class="text-sm">{{ t('theme.showRightSidebar') }}</span>
            </label>
          </div>
          
          <!-- 侧边栏宽度 -->
          <div>
            <label class="block text-sm font-medium mb-2">{{ t('theme.sidebarWidth') }}</label>
            <input type="range" v-model.number="sidebarWidthNum" min="200" max="350" class="w-full">
            <div class="text-xs text-gray-500 mt-1">{{ preview.sidebarWidth }}</div>
          </div>
          
          <!-- 功能开关 -->
          <div>
            <label class="block text-sm font-medium mb-2">{{ t('theme.features') }}</label>
            <div class="space-y-2">
              <label v-for="feature in configurableFeatures" :key="feature"
                class="flex items-center space-x-2">
                <input type="checkbox" v-model="customFeatures[feature]">
                <span class="text-sm">{{ getFeatureLabel(feature) }}</span>
              </label>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    
    <template #footer>
      <button @click="resetAll" class="text-gray-500 hover:text-gray-700">
        {{ t('theme.resetAll') }}
      </button>
      <el-button @click="visible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" @click="saveConfig" :loading="saving">
        {{ t('common.save') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  theme: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const activeTab = ref('colors')
const saving = ref(false)

const layoutTypes = [
  { value: 'three-column', icon: '🟦🟩🟪', label: t('theme.threeColumn') },
  { value: 'two-column', icon: '🟦🟩', label: t('theme.twoColumn') },
  { value: 'single-column', icon: '🟩', label: t('theme.singleColumn') }
]

const customColors = ref({})
const customFeatures = ref({})
const preview = ref({
  type: 'three-column',
  sidebarLeft: true,
  sidebarRight: true,
  sidebarWidth: '256px'
})

const configurableColors = computed(() => {
  return props.theme?.config?.configurable?.colors || []
})

const configurableFeatures = computed(() => {
  return props.theme?.config?.configurable?.features || []
})

const sidebarWidthNum = computed({
  get: () => parseInt(preview.value.sidebarWidth),
  set: (val) => preview.value.sidebarWidth = `${val}px`
})

const previewLayoutClass = computed(() => {
  if (preview.value.type === 'three-column') return 'flex-row'
  if (preview.value.type === 'two-column') return 'flex-row'
  return 'flex-col'
})

function getColorLabel(key) {
  const labels = {
    primary: t('theme.primaryColor'),
    secondary: t('theme.secondaryColor'),
    background: t('theme.backgroundColor'),
    text: t('theme.textColor')
  }
  return labels[key] || key
}

function getFeatureLabel(key) {
  const labels = {
    showAvatars: t('theme.showAvatars'),
    showBadges: t('theme.showBadges'),
    compactMode: t('theme.compactMode')
  }
  return labels[key] || key
}

function setLayoutType(type) {
  preview.value.type = type
  if (type === 'single-column') {
    preview.value.sidebarLeft = false
    preview.value.sidebarRight = false
  } else if (type === 'two-column') {
    preview.value.sidebarLeft = true
    preview.value.sidebarRight = false
  } else {
    preview.value.sidebarLeft = true
    preview.value.sidebarRight = true
  }
}

function resetColor(key) {
  if (props.theme?.config?.colors?.[key]) {
    customColors.value[key] = props.theme.config.colors[key]
  }
}

function resetAll() {
  if (props.theme?.config) {
    customColors.value = { ...props.theme.config.colors }
    customFeatures.value = { ...props.theme.config.features }
    preview.value = {
      type: props.theme.config.layouts?.type || 'three-column',
      sidebarLeft: props.theme.config.layouts?.sidebarLeft !== false,
      sidebarRight: props.theme.config.layouts?.sidebarRight !== false,
      sidebarWidth: props.theme.config.layouts?.sidebarWidth || '256px'
    }
  }
}

async function saveConfig() {
  saving.value = true
  try {
    await api.post(`/admin/themes/${props.theme.name}/config`, {
      colors: customColors.value,
      features: customFeatures.value,
      layouts: preview.value
    })
    ElMessage.success(t('theme.configSaved'))
    visible.value = false
  } catch (e) {
    console.error(e)
    ElMessage.error(t('theme.configSaveFailed'))
  } finally {
    saving.value = false
  }
}

watch(() => props.theme, (theme) => {
  if (theme) {
    resetAll()
  }
}, { immediate: true })
</script>
```

### 4.3 后端模型（Go）

```go
// server/models/theme.go
package models

import (
	"gorm.io/gorm"
	"time"
)

type Theme struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null;uniqueIndex" json:"name"`
	DisplayName string    `gorm:"size:100" json:"display_name"`
	Description string    `gorm:"type:text" json:"description"`
	Version     string    `gorm:"size:20" json:"version"`
	Author      string    `gorm:"size:100" json:"author"`
	Screenshot  string    `gorm:"size:255" json:"screenshot"`
	IsActive    bool      `gorm:"default:false" json:"is_active"`
	Config      JSON      `gorm:"type:json" json:"config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SiteThemeConfig struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ThemeName string    `gorm:"size:50;default:'default'" json:"theme_name"`
	Config    JSON      `gorm:"type:json" json:"config"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

### 4.4 后端 Handler（Go）

```go
// server/handlers/admin_theme.go
package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"bbsgo/models"
)

func GetThemes(c *gin.Context) {
	var themes []models.Theme
	models.DB.Find(&themes)
	c.JSON(http.StatusOK, themes)
}

func ActivateTheme(c *gin.Context) {
	themeName := c.Param("name")
	
	models.DB.Model(&models.Theme{}).Where("is_active = ?", true).Update("is_active", false)
	models.DB.Model(&models.Theme{}).Where("name = ?", themeName).Update("is_active", true)
	
	var siteConfig models.SiteThemeConfig
	models.DB.FirstOrCreate(&siteConfig, models.SiteThemeConfig{})
	siteConfig.ThemeName = themeName
	models.DB.Save(&siteConfig)
	
	c.JSON(http.StatusOK, gin.H{"message": "Theme activated"})
}

func SaveThemeConfig(c *gin.Context) {
	themeName := c.Param("name")
	
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var theme models.Theme
	models.DB.Where("name = ?", themeName).First(&theme)
	
	if theme.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}
	
	theme.Config = config
	models.DB.Save(&theme)
	
	var siteConfig models.SiteThemeConfig
	models.DB.FirstOrCreate(&siteConfig, models.SiteThemeConfig{})
	siteConfig.Config = config
	models.DB.Save(&siteConfig)
	
	c.JSON(http.StatusOK, gin.H{"message": "Config saved"})
}

func DeleteTheme(c *gin.Context) {
	themeName := c.Param("name")
	
	if themeName == "default" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete default theme"})
		return
	}
	
	models.DB.Where("name = ?", themeName).Delete(&models.Theme{})
	c.JSON(http.StatusOK, gin.H{"message": "Theme deleted"})
}

func UploadTheme(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Theme uploaded"})
}
```

---

## 五、主题应用（site 前端）

### 5.1 布局渲染引擎

```vue
<!-- site/src/components/LayoutEngine.vue -->
<template>
  <div class="layout-engine" :class="layoutClass">
    <!-- 左栏 -->
    <aside v-if="showSidebarLeft" 
      class="sidebar sidebar-left"
      :style="{ width: config.sidebarWidth }">
      <slot name="sidebar-left"></slot>
      <template v-for="slotName in leftSlotComponents" :key="slotName">
        <component :is="getComponent(slotName)" />
      </template>
    </aside>
    
    <!-- 主内容区 -->
    <main class="content-area flex-1 min-w-0">
      <slot></slot>
    </main>
    
    <!-- 右栏 -->
    <aside v-if="showSidebarRight" 
      class="sidebar sidebar-right"
      :style="{ width: config.sidebarWidth }">
      <slot name="sidebar-right"></slot>
      <template v-for="slotName in rightSlotComponents" :key="slotName">
        <component :is="getComponent(slotName)" />
      </template>
    </aside>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

const props = defineProps({
  pageType: {
    type: String,
    default: 'home'
  }
})

const themeStore = useThemeStore()

const config = computed(() => {
  const pageLayouts = themeStore.layoutConfig?.pageLayouts || {}
  const siteLayouts = themeStore.siteConfig?.layouts || {}
  return { ...pageLayouts[props.pageType], ...siteLayouts }
})

const showSidebarLeft = computed(() => config.value.sidebarLeft !== false)
const showSidebarRight = computed(() => config.value.sidebarRight !== false)

const layoutClass = computed(() => {
  return [
    `layout-${config.value.type || 'three-column'}`,
    'flex flex-col lg:flex-row gap-4'
  ].join(' ')
})

function getComponent(name) {
  return themeStore.themeComponents[name]
}
</script>

<style scoped>
.layout-three-column {
  /* 三列布局 */
}

.layout-two-column {
  /* 两列布局 */
}

.layout-single-column {
  /* 单列布局 */
}
</style>
```

### 5.2 在页面中使用

```vue
<!-- site/src/views/Home.vue -->
<template>
  <LayoutEngine page-type="home">
    <template #sidebar-left>
      <SidebarTags v-if="showSidebarTags" />
    </template>
    
    <!-- 主内容 -->
    <div class="space-y-4">
      <div v-for="topic in topics" :key="topic.id">
        <TopicCard :topic="topic" />
      </div>
    </div>
    
    <template #sidebar-right>
      <SidebarHotTopics v-if="showSidebarHotTopics" />
      <SidebarCreditUsers v-if="showSidebarCreditUsers" />
    </template>
  </LayoutEngine>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import LayoutEngine from '@/components/LayoutEngine.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const layout = computed(() => {
  const pageLayouts = themeStore.layoutConfig?.pageLayouts || {}
  return pageLayouts.home || {}
})

const showSidebarTags = computed(() => layout.value.slots?.['sidebar-left']?.includes('SidebarTags'))
const showSidebarHotTopics = computed(() => layout.value.slots?.['sidebar-right']?.includes('SidebarHotTopics'))
const showSidebarCreditUsers = computed(() => layout.value.slots?.['sidebar-right']?.includes('SidebarCreditUsers'))

onMounted(() => {
  themeStore.loadSiteTheme()
})
</script>
```

---

## 六、完整布局定制方案

### 6.1 布局系统概念

```
布局层级结构：
┌─────────────────────────────────────┐
│         Page Layout                 │  (3列/2列/1列 - 全站统一)
│  ┌─────────┬───────────┬─────────┐ │
│  │ Sidebar │  Content  │ Sidebar │ │  (可隐藏/交换位置)
│  │  Left   │   Area    │  Right  │ │
│  └─────────┴───────────┴─────────┘ │
└─────────────────────────────────────┘
         ↓
    组件插槽系统
```

### 6.2 布局配置示例

```javascript
// site/src/themes/default/layout.js
export default {
  name: 'Default Layout',
  
  // 页面布局配置（主题默认）
  pageLayouts: {
    home: {
      type: 'three-column',
      sidebarLeft: true,
      sidebarRight: true,
      sidebarWidth: '256px',
      contentOrder: ['sidebar-left', 'content', 'sidebar-right']
    },
    profile: {
      type: 'two-column',
      sidebarLeft: true,
      sidebarRight: false,
      sidebarWidth: '280px'
    },
    topic: {
      type: 'single-column'
    }
  },
  
  // 组件映射
  components: {
    Header: 'components/Header.vue',
    Footer: 'components/Footer.vue',
    TopicCard: 'components/TopicCard.vue',
    SidebarHotTopics: 'components/SidebarHotTopics.vue'
  },
  
  // 区域插槽配置
  slots: {
    'home:sidebar-left': ['SidebarTags', 'SidebarUserInfo'],
    'home:sidebar-right': ['SidebarHotTopics', 'SidebarCreditUsers'],
    'profile:sidebar-left': ['SidebarUserInfo', 'SidebarBadges']
  },
  
  // 响应式断点
  breakpoints: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px'
  }
}
```

---

## 七、实施路线图

| 阶段 | 内容 | 位置 | 预估工作量 |
|------|------|------|-----------|
| **Phase 1** | CSS 主题系统 + 主题配置加载 | site | 2-3 天 |
| **Phase 2** | 布局引擎核心 | site | 3-4 天 |
| **Phase 3** | 主题列表与激活 | admin | 2-3 天 |
| **Phase 4** | 颜色定制器 | admin | 2-3 天 |
| **Phase 5** | 布局定制器 | admin | 3-4 天 |
| **Phase 6** | 后端 API 和数据模型 | server | 2-3 天 |
| **Phase 7** | 多主题示例（深色/简约） | site/themes | 3-5 天 |

---

## 八、关键技术点

### 8.1 技术栈
- **主题管理**：Pinia Store
- **CSS 方案**：CSS Variables + Tailwind CSS
- **动态加载**：Vite 动态 import
- **配置存储**：数据库（全站配置）
- **UI 框架**：Vue 3 + Element Plus

### 8.2 核心技术点

1. **动态组件加载**：Vite `import()` + 组件映射表
2. **插槽系统**：Vue 的具名插槽 + 动态插槽
3. **响应式断点**：Tailwind CSS 断点 + CSS 媒体查询
4. **数据持久化**：数据库存储全站配置
5. **类型安全**：TypeScript 接口定义配置结构

---

## 九、功能清单

### 9.1 站主/管理员（admin 后台）
- ✅ 查看所有可用主题
- ✅ 上传新主题
- ✅ 激活/切换主题
- ✅ 删除主题（默认主题除外）
- ✅ 定制主题颜色（全站生效）
- ✅ 定制主题布局（全站生效）
- ✅ 定制功能开关（全站生效）
- ✅ 实时预览效果
- ✅ 重置为默认配置

### 9.2 普通用户（site 前端）
- ✅ 使用管理员设置的主题
- ✅ 享受全站统一的风格
- ✅ 无定制权限（只读）

---

## 十、后端 API 设计

### 10.1 主题管理 API（admin）

```
GET    /api/admin/themes              # 获取所有主题
POST   /api/admin/themes              # 上传新主题
POST   /api/admin/themes/:name/activate  # 激活主题
POST   /api/admin/themes/:name/config    # 保存主题配置
DELETE /api/admin/themes/:name        # 删除主题
```

### 10.2 站点主题 API（site）

```
GET    /api/site/theme          # 获取当前站点主题配置
```

---

## 十一、文件结构总结

```
bbsgo/
├── site/
│   ├── src/
│   │   ├── themes/
│   │   │   ├── default/
│   │   │   │   ├── components/
│   │   │   │   │   ├── Header.vue
│   │   │   │   │   ├── Footer.vue
│   │   │   │   │   ├── TopicCard.vue
│   │   │   │   │   └── ...
│   │   │   │   ├── views/
│   │   │   │   │   ├── Home.vue
│   │   │   │   │   ├── Profile.vue
│   │   │   │   │   └── ...
│   │   │   │   ├── styles/
│   │   │   │   │   └── index.css
│   │   │   │   ├── theme.json
│   │   │   │   ├── layout.js
│   │   │   │   └── index.js
│   │   │   ├── dark/
│   │   │   └── minimalist/
│   │   ├── components/
│   │   │   └── LayoutEngine.vue
│   │   ├── stores/
│   │   │   └── theme.js
│   │   └── ...
│   └── ...
├── admin/
│   ├── src/
│   │   ├── views/
│   │   │   └── Themes.vue
│   │   ├── components/
│   │   │   └── ThemeCustomizer.vue
│   │   └── ...
│   └── ...
├── server/
│   ├── models/
│   │   └── theme.go
│   ├── handlers/
│   │   └── admin_theme.go
│   └── ...
└── THEME-CUSTOMIZATION.md
```

---

## 十二、后续扩展建议

1. **主题市场**：允许站主上传、分享、下载主题
2. **主题导入/导出**：备份和恢复主题配置
3. **子主题系统**：基于父主题创建子主题，只覆盖需要修改的部分
4. **主题预览**：切换前实时预览全站效果
5. **主题模板生成器**：基于现有主题快速创建新主题

---

*文档版本：v3.0 | 最后更新：2026-04-04*
