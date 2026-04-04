# BBSGo Template Customization System (Plan A: Component-Level Themes)

## Table of Contents
- [I. System Architecture Design](#i-system-architecture-design)
- [II. Core Concepts and Division of Labor](#ii-core-concepts-and-division-of-labor)
- [III. Detailed Design of Recommended Plan](#iii-detailed-design-of-recommended-plan)
- [IV. Theme Management and Customization (Admin Backend)](#iv-theme-management-and-customization-admin-backend)
- [V. Theme Application (Site Frontend)](#v-theme-application-site-frontend)
- [VI. Complete Layout Customization Plan](#vi-complete-layout-customization-plan)
- [VII. Implementation Roadmap](#vii-implementation-roadmap)
- [VIII. Key Technical Points](#viii-key-technical-points)

---

## I. System Architecture Design

### 1.1 Theme Directory Structure

```
bbsgo/
├── site/
│   ├── src/
│   │   ├── themes/              # Theme directory
│   │   │   ├── default/       # Default theme
│   │   │   │   ├── components/
│   │   │   │   ├── views/
│   │   │   │   ├── styles/
│   │   │   │   ├── theme.json      # Theme config
│   │   │   │   ├── layout.js       # Layout config
│   │   │   │   └── index.js       # Theme entry
│   │   │   ├── dark/         # Dark theme
│   │   │   └── minimalist/  # Minimalist theme
│   │   └── ...
│   ├── ...
└── ...
```

### 1.2 Theme Config File (theme.json)

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

## II. Core Concepts and Division of Labor

### 2.1 Two Levels of Responsibilities

| Level | Location | Function | User Role |
|-------|----------|----------|-----------|
| **Theme Management & Customization** | `admin/src/` | Install, activate, delete themes, customize site-wide colors, layouts | Site owner/admin |
| **Theme Application** | `site/src/` | Apply admin-configured theme, unified site style | Regular users (read-only) |

### 2.2 Data Flow

```
Site owner (admin backend)
    ↓ Select theme + customize
Theme config (database)
    ↓ Site-wide application
Frontend rendering (site/src/)
```

---

## III. Detailed Design of Recommended Plan

### 3.1 Core Implementation Steps

#### Step 1: Create Theme Store (Site Frontend)

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

#### Step 2: Dynamic Component Loading

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

#### Step 3: CSS Variable System

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

## IV. Theme Management and Customization (Admin Backend)

### 4.1 Theme List and Selection Page

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
    
    <!-- Theme list -->
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
    
    <!-- Upload dialog -->
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

### 4.2 Theme Customizer (Colors + Layout)

```vue
<!-- admin/src/components/ThemeCustomizer.vue -->
<template>
  <el-dialog v-model="visible" :title="t('theme.customizeTheme')" width="800px" :close-on-click-modal="false">
    <el-tabs v-model="activeTab">
      <!-- Color customization -->
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
      
      <!-- Layout customization -->
      <el-tab-pane :label="t('theme.layout')" name="layout">
        <div class="space-y-6">
          <!-- Layout preview -->
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
          
          <!-- Layout type selection -->
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
          
          <!-- Sidebar toggles -->
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
          
          <!-- Sidebar width -->
          <div>
            <label class="block text-sm font-medium mb-2">{{ t('theme.sidebarWidth') }}</label>
            <input type="range" v-model.number="sidebarWidthNum" min="200" max="350" class="w-full">
            <div class="text-xs text-gray-500 mt-1">{{ preview.sidebarWidth }}</div>
          </div>
          
          <!-- Feature toggles -->
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

### 4.3 Backend Model (Go)

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

### 4.4 Backend Handler (Go)

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

## V. Theme Application (Site Frontend)

### 5.1 Layout Rendering Engine

```vue
<!-- site/src/components/LayoutEngine.vue -->
<template>
  <div class="layout-engine" :class="layoutClass">
    <!-- Left sidebar -->
    <aside v-if="showSidebarLeft" 
      class="sidebar sidebar-left"
      :style="{ width: config.sidebarWidth }">
      <slot name="sidebar-left"></slot>
      <template v-for="slotName in leftSlotComponents" :key="slotName">
        <component :is="getComponent(slotName)" />
      </template>
    </aside>
    
    <!-- Main content area -->
    <main class="content-area flex-1 min-w-0">
      <slot></slot>
    </main>
    
    <!-- Right sidebar -->
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
  /* Three column layout */
}

.layout-two-column {
  /* Two column layout */
}

.layout-single-column {
  /* Single column layout */
}
</style>
```

### 5.2 Usage in Pages

```vue
<!-- site/src/views/Home.vue -->
<template>
  <LayoutEngine page-type="home">
    <template #sidebar-left>
      <SidebarTags v-if="showSidebarTags" />
    </template>
    
    <!-- Main content -->
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

## VI. Complete Layout Customization Plan

### 6.1 Layout System Concept

```
Layout hierarchy:
┌─────────────────────────────────────┐
│         Page Layout                 │  (3-column/2-column/1-column - site-wide)
│  ┌─────────┬───────────┬─────────┐ │
│  │ Sidebar │  Content  │ Sidebar │ │  (Can be hidden/swap positions)
│  │  Left   │   Area    │  Right  │ │
│  └─────────┴───────────┴─────────┘ │
└─────────────────────────────────────┘
         ↓
    Component slot system
```

### 6.2 Layout Config Example

```javascript
// site/src/themes/default/layout.js
export default {
  name: 'Default Layout',
  
  // Page layout config (theme default)
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
  
  // Component mapping
  components: {
    Header: 'components/Header.vue',
    Footer: 'components/Footer.vue',
    TopicCard: 'components/TopicCard.vue',
    SidebarHotTopics: 'components/SidebarHotTopics.vue'
  },
  
  // Area slot config
  slots: {
    'home:sidebar-left': ['SidebarTags', 'SidebarUserInfo'],
    'home:sidebar-right': ['SidebarHotTopics', 'SidebarCreditUsers'],
    'profile:sidebar-left': ['SidebarUserInfo', 'SidebarBadges']
  },
  
  // Responsive breakpoints
  breakpoints: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px'
  }
}
```

---

## VII. Implementation Roadmap

| Phase | Content | Location | Estimated Effort |
|-------|---------|----------|-----------------|
| **Phase 1** | CSS theme system + theme config loading | site | 2-3 days |
| **Phase 2** | Layout engine core | site | 3-4 days |
| **Phase 3** | Theme list and activation | admin | 2-3 days |
| **Phase 4** | Color customizer | admin | 2-3 days |
| **Phase 5** | Layout customizer | admin | 3-4 days |
| **Phase 6** | Backend API and data models | server | 2-3 days |
| **Phase 7** | Multi-theme examples (dark/minimalist) | site/themes | 3-5 days |

---

## VIII. Key Technical Points

### 8.1 Tech Stack
- **Theme Management**: Pinia Store
- **CSS Solution**: CSS Variables + Tailwind CSS
- **Dynamic Loading**: Vite dynamic import
- **Config Storage**: Database (site-wide config)
- **UI Framework**: Vue 3 + Element Plus

### 8.2 Core Technical Points

1. **Dynamic Component Loading**: Vite `import()` + component mapping table
2. **Slot System**: Vue named slots + dynamic slots
3. **Responsive Breakpoints**: Tailwind CSS breakpoints + CSS media queries
4. **Data Persistence**: Database stores site-wide config
5. **Type Safety**: TypeScript interfaces define config structure

---

## IX. Feature Checklist

### 9.1 Site Owner/Administrator (Admin Backend)
- ✅ View all available themes
- ✅ Upload new themes
- ✅ Activate/switch themes
- ✅ Delete themes (except default theme)
- ✅ Customize theme colors (site-wide effect)
- ✅ Customize theme layout (site-wide effect)
- ✅ Customize feature toggles (site-wide effect)
- ✅ Real-time preview effects
- ✅ Reset to default config

### 9.2 Regular Users (Site Frontend)
- ✅ Use admin-configured theme
- ✅ Enjoy unified site style
- ✅ No customization permissions (read-only)

---

## X. Backend API Design

### 10.1 Theme Management API (Admin)

```
GET    /api/admin/themes              # Get all themes
POST   /api/admin/themes              # Upload new theme
POST   /api/admin/themes/:name/activate  # Activate theme
POST   /api/admin/themes/:name/config    # Save theme config
DELETE /api/admin/themes/:name        # Delete theme
```

### 10.2 Site Theme API (Site)

```
GET    /api/site/theme          # Get current site theme config
```

---

## XI. File Structure Summary

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

## XII. Future Expansion Suggestions

1. **Theme Market**: Allow site owners to upload, share, download themes
2. **Theme Import/Export**: Backup and restore theme configurations
3. **Child Theme System**: Create child themes based on parent themes, only overwrite what needs to be changed
4. **Theme Preview**: Real-time preview of full site effects before switching
5. **Theme Template Generator**: Quickly create new themes based on existing themes

---

*Document version: v3.0 | Last updated: 2026-04-04*
