# BBSGo 主题风格设计（完全差异化）

## 目录
- [主题风格对比](#主题风格对比)
- [风格 1：现代卡片式（Material Design）](#风格-1现代卡片式material-design)
- [风格 2：新闻门户式](#风格-2新闻门户式)
- [风格 3：社交动态流](#风格-3社交动态流)
- [风格 4：论坛经典式](#风格-4论坛经典式)
- [风格 5：极简阅读式](#风格-5极简阅读式)

---

## 主题风格对比

| 特性 | 现代卡片式 | 新闻门户式 | 社交动态流 | 论坛经典式 | 极简阅读式 |
|------|-----------|-----------|-----------|-----------|-----------|
| **布局** | 三列卡片 | 两列杂志 | 单列瀑布 | 经典论坛 | 居中单列 |
| **卡片** | 大圆角阴影 | 图文混排 | 左右对齐 | 表格行 | 无边框 |
| **导航** | 顶部固定 | 顶部菜单 | 底部 Tab | 侧边导航 | 顶部极简 |
| **互动** | Hover 动效 | 大图标题 | 点赞评论 | 楼层回复 | 专注内容 |
| **适用** | 通用社区 | 内容门户 | 年轻人社区 | 技术论坛 | 博客/阅读 |

---

## 风格 1：现代卡片式（Material Design）

### 设计理念
- Google Material Design 风格
- 大圆角、柔和阴影
- 清晰的层次结构
- 流畅的交互动画

### 首页布局
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]  [搜索]                [导航]  [用户头像]  [通知] │ ← 顶部导航栏
├──────────┬──────────────────────────────────────────┬───────────┤
│          │  [标签栏] 全部 | 热门 | 最新 | 精华     │           │
│  [热门]  │  ┌────────────────────────────────────┐  │  [活跃]  │
│  话题   │  │ 🟦  [置顶标记] 帖子大标题        │  │  用户    │
│         │  │  摘要文字...                       │  │           │
│  [热门]  │  │  [图片] [图片] [图片]            │  │  [热门]  │
│  用户   │  └────────────────────────────────────┘  │  话题    │
│         │  ┌────────────────────────────────────┐  │           │
│         │  │ 🟦  帖子大标题                   │  │           │
│         │  │  摘要文字...                       │  │           │
│         │  └────────────────────────────────────┘  │           │
└──────────┴──────────────────────────────────────────┴───────────┘
```

### 卡片样式
```vue
<!-- 大圆角、大阴影 -->
<div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-300">
  <div class="flex items-center space-x-4 mb-4">
    <img class="w-12 h-12 rounded-full" src="avatar">
    <div>
      <div class="font-semibold">用户名</div>
      <div class="text-sm text-gray-500">2小时前</div>
    </div>
  </div>
  <h2 class="text-xl font-bold mb-3">帖子大标题</h2>
  <p class="text-gray-600 mb-4">摘要内容...</p>
  <div class="grid grid-cols-3 gap-2 mb-4">
    <img class="rounded-lg" src="img1">
    <img class="rounded-lg" src="img2">
    <img class="rounded-lg" src="img3">
  </div>
  <div class="flex items-center justify-between pt-4 border-t">
    <button class="flex items-center space-x-2 px-4 py-2 rounded-full hover:bg-red-50">
      ❤️ 128
    </button>
    <button class="flex items-center space-x-2 px-4 py-2 rounded-full hover:bg-blue-50">
      💬 32
    </button>
  </div>
</div>
```

### 颜色方案
```css
:root {
  --color-primary: #6366F1;
  --color-primary-dark: #4F46E5;
  --color-surface: #FFFFFF;
  --color-surface-variant: #F8FAFC;
  --color-on-surface: #1F2937;
  --shadow-1: 0 1px 3px rgba(0,0,0,0.12);
  --shadow-2: 0 4px 6px rgba(0,0,0,0.1);
  --shadow-3: 0 10px 15px rgba(0,0,0,0.1);
  --radius-xl: 16px;
  --radius-2xl: 24px;
}
```

---

## 风格 2：新闻门户式

### 设计理念
- 杂志/新闻门户风格
- 大图+标题的视觉冲击
- 信息密度高
- 多栏目布局

### 首页布局
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]          [首页] [科技] [生活] [娱乐] [更多]     [登录] │
├─────────────────────────────────────────────────────────────────┤
│  📰 [轮播大图] 置顶头条1              [边栏热门] 标题1        │
│     标题文字大标题                    标题2                 │
│                                       标题3                 │
│  ┌──────────────┬───────────────────────────────────────────┐  │
│  │ [小图] 标题2 │ [小图] 标题5       [小图] 标题8       │  │
│  │  摘要...     │  摘要...           摘要...           │  │
│  ├──────────────┼───────────────────────────────────────────┤  │
│  │ [小图] 标题3 │ [小图] 标题6       [小图] 标题9       │  │
│  │  摘要...     │  摘要...           摘要...           │  │
│  ├──────────────┼───────────────────────────────────────────┤  │
│  │ [小图] 标题4 │ [小图] 标题7       [小图] 标题10      │  │
│  │  摘要...     │  摘要...           摘要...           │  │
│  └──────────────┴───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 文章卡片样式
```vue
<!-- 头条大卡片 -->
<div class="relative overflow-hidden rounded-lg">
  <img class="w-full h-64 object-cover" src="cover-image">
  <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-6">
    <span class="inline-block bg-red-600 text-white px-3 py-1 rounded text-sm mb-2">头条</span>
    <h1 class="text-2xl font-bold text-white">这是头条大标题</h1>
    <p class="text-gray-200 mt-2">摘要内容...</p>
  </div>
</div>

<!-- 列表小卡片 -->
<div class="flex space-x-4 p-4 border-b hover:bg-gray-50">
  <img class="w-24 h-24 object-cover rounded" src="thumbnail">
  <div class="flex-1">
    <h3 class="font-bold hover:text-blue-600">文章标题</h3>
    <p class="text-sm text-gray-500 mt-1 line-clamp-2">摘要内容...</p>
    <div class="flex items-center space-x-3 mt-2 text-xs text-gray-400">
      <span>来源</span>
      <span>·</span>
      <span>2小时前</span>
      <span>·</span>
      <span>1.2k 阅读</span>
    </div>
  </div>
</div>
```

### 颜色方案
```css
:root {
  --color-primary: #DC2626;
  --color-accent: #2563EB;
  --color-background: #FFFFFF;
  --color-border: #E5E7EB;
  --color-text: #111827;
  --color-text-secondary: #6B7280;
  --font-heading: 'Noto Serif SC', serif;
  --font-body: 'Noto Sans SC', sans-serif;
}
```

---

## 风格 3：社交动态流

### 设计理念
- 类似微博、Twitter 的时间线
- 左右分栏，内容居中
- 强调互动（点赞、转发、评论）
- 移动端友好

### 首页布局
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]  [🏠 首页] [🔍 探索] [✏️ 发帖]  [🔔] [👤]        │ ← 顶部
├──────────┬──────────────────────────┬───────────────────────────┤
│          │  [发布框] 分享想法...     │  [我的资料]              │
│  [菜单]  │  ┌─────────────────────┐  │  [头像] 用户名        │
│  首页    │  │ 🟦  [头像] 用户名  │  │  关注 128 | 粉丝 256 │
│  探索    │  │  动态内容...        │  │                         │
│  消息    │  │  [图片] [图片]     │  │  [推荐关注]            │
│  通知    │  │  ❤️ 128  🔄 32  💬  │  │  [头像] 用户名  [+]   │
│          │  └─────────────────────┘  │  [头像] 用户名  [+]   │
│          │  ┌─────────────────────┐  │  [头像] 用户名  [+]   │
│          │  │ 🟦  [头像] 用户名  │  │                         │
│          │  │  动态内容...        │  │  [热门话题]            │
│          │  │  ❤️ 256  🔄 64  💬  │  │  #话题1  12.5k       │
│          │  └─────────────────────┘  │  #话题2   8.2k       │
└──────────┴──────────────────────────┴───────────────────────────┘
```

### 动态卡片样式
```vue
<div class="border-b border-gray-200 p-4 hover:bg-gray-50">
  <div class="flex space-x-3">
    <img class="w-12 h-12 rounded-full" src="avatar">
    <div class="flex-1 min-w-0">
      <div class="flex items-center space-x-2">
        <span class="font-bold">用户名</span>
        <span class="text-gray-500">@username</span>
        <span class="text-gray-400">·</span>
        <span class="text-gray-500">2h</span>
      </div>
      <p class="mt-1 text-gray-900 whitespace-pre-wrap">
        动态内容，可以是多行文字，支持换行
      </p>
      <div v-if="images" class="mt-3 grid grid-cols-2 gap-2 rounded-lg overflow-hidden">
        <img v-for="img in images" :key="img" :src="img" class="w-full h-48 object-cover">
      </div>
      <div class="flex items-center justify-between mt-4 max-w-md text-gray-500">
        <button class="flex items-center space-x-2 hover:text-blue-500">
          <svg class="w-5 h-5">...</svg>
          <span>128</span>
        </button>
        <button class="flex items-center space-x-2 hover:text-green-500">
          <svg class="w-5 h-5">...</svg>
          <span>32</span>
        </button>
        <button class="flex items-center space-x-2 hover:text-red-500">
          <svg class="w-5 h-5">...</svg>
          <span>256</span>
        </button>
        <button class="flex items-center space-x-2 hover:text-blue-500">
          <svg class="w-5 h-5">...</svg>
        </button>
      </div>
    </div>
  </div>
</div>
```

### 颜色方案
```css
:root {
  --color-primary: #1D9BF0;
  --color-secondary: #0F1419;
  --color-background: #FFFFFF;
  --color-border: #EFF3F4;
  --color-text: #0F1419;
  --color-text-secondary: #536471;
  --hover-blue: rgba(29, 155, 240, 0.1);
  --hover-red: rgba(244, 33, 46, 0.1);
  --hover-green: rgba(0, 186, 124, 0.1);
}
```

---

## 风格 4：论坛经典式

### 设计理念
- 传统论坛（Discourse、phpBB）风格
- 表格式列表，信息密度高
- 楼层回复，清晰的讨论树
- 侧边栏导航

### 首页布局
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]                    [用户] [消息] [设置] [退出]        │
├──────────┬──────────────────────────────────────────────────────┤
│          │  ┌─────────────────────────────────────────────┐    │
│  [版块]  │  │ [图标] 版块名称                主题 | 帖子 │    │
│  技术区  │  ├─────────────────────────────────────────────┤    │
│  讨论区  │  │ [置顶] [精华] 帖子标题                │    │
│  分享区  │  │  [头像] 作者 · 2小时前 · 32回复 · 1.2k  │    │
│  反馈区  │  ├─────────────────────────────────────────────┤    │
│          │  │ [精华] 帖子标题                        │    │
│  [在线]  │  │  [头像] 作者 · 4小时前 · 64回复 · 2.5k  │    │
│  用户    │  ├─────────────────────────────────────────────┤    │
│  [头像]  │  │ 帖子标题                               │    │
│  [头像]  │  │  [头像] 作者 · 6小时前 · 8回复 · 320    │    │
│          │  └─────────────────────────────────────────────┘    │
│          │  [分页]  1 2 3 4 5 ... 下一页                      │
└──────────┴──────────────────────────────────────────────────────┘
```

### 帖子列表样式
```vue
<!-- 表格式列表 -->
<table class="w-full">
  <thead class="bg-gray-100">
    <tr>
      <th class="px-4 py-3 text-left">主题</th>
      <th class="px-4 py-3 text-center w-24">回复</th>
      <th class="px-4 py-3 text-center w-24">浏览</th>
      <th class="px-4 py-3 text-left w-48">最后更新</th>
    </tr>
  </thead>
  <tbody>
    <tr class="border-b hover:bg-gray-50">
      <td class="px-4 py-4">
        <div class="flex items-start space-x-3">
          <img class="w-10 h-10 rounded-full" src="avatar">
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2">
              <span v-if="topic.isPinned" class="bg-red-100 text-red-700 px-2 py-0.5 rounded text-xs">置顶</span>
              <span v-if="topic.isEssence" class="bg-yellow-100 text-yellow-700 px-2 py-0.5 rounded text-xs">精华</span>
              <a href="#" class="font-medium hover:text-blue-600">帖子标题</a>
            </div>
            <div class="text-sm text-gray-500 mt-1">
              作者 · 发布于 2小时前
            </div>
          </div>
        </div>
      </td>
      <td class="px-4 py-4 text-center font-medium">32</td>
      <td class="px-4 py-4 text-center text-gray-500">1.2k</td>
      <td class="px-4 py-4">
        <div class="flex items-center space-x-2">
          <img class="w-6 h-6 rounded-full" src="lastAvatar">
          <div class="text-sm">
            <div>最后回复者</div>
            <div class="text-gray-500">30分钟前</div>
          </div>
        </div>
      </td>
    </tr>
  </tbody>
</table>
```

### 颜色方案
```css
:root {
  --color-primary: #2563EB;
  --color-success: #16A34A;
  --color-warning: #D97706;
  --color-danger: #DC2626;
  --color-background: #FFFFFF;
  --color-header: #F9FAFB;
  --color-border: #E5E7EB;
  --color-text: #111827;
  --color-text-secondary: #6B7280;
}
```

---

## 风格 5：极简阅读式

### 设计理念
- 类似 Medium、简书的阅读体验
- 大留白，专注内容
- 居中窄栏，舒适阅读
- 无干扰元素

### 首页布局
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]                              [搜索] [用户]         │ ← 极简
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│                    ┌─────────────────────┐                      │
│                    │  精选文章           │                      │
│                    ├─────────────────────┤                      │
│                    │  文章大标题         │                      │
│                    │                     │                      │
│                    │  摘要内容...        │                      │
│                    │                     │                      │
│                    │  [作者头像] 作者名  │                      │
│                    │  2026-04-04 · 5分钟读 │                    │
│                    └─────────────────────┘                      │
│                    ┌─────────────────────┐                      │
│                    │  文章大标题         │                      │
│                    │                     │                      │
│                    │  摘要内容...        │                      │
│                    │                     │                      │
│                    │  [作者头像] 作者名  │                      │
│                    └─────────────────────┘                      │
│                    ┌─────────────────────┐                      │
│                    │  文章大标题         │                      │
│                    └─────────────────────┘                      │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 文章卡片样式
```vue
<article class="max-w-2xl mx-auto py-12">
  <header class="mb-8">
    <div class="flex items-center space-x-3 mb-4">
      <img class="w-12 h-12 rounded-full" src="authorAvatar">
      <div>
        <div class="font-medium">作者姓名</div>
        <div class="text-sm text-gray-500">2026年4月4日 · 5分钟阅读</div>
      </div>
    </div>
    <h1 class="text-3xl font-bold leading-tight">文章大标题</h1>
  </header>
  
  <div v-if="coverImage" class="mb-8">
    <img :src="coverImage" class="w-full rounded-lg">
  </div>
  
  <div class="prose prose-lg prose-gray max-w-none">
    <p>第一段内容...</p>
    <p>第二段内容...</p>
    <blockquote>引用内容</blockquote>
    <pre><code>代码块</code></pre>
  </div>
  
  <footer class="mt-12 pt-8 border-t">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <button class="flex items-center space-x-2 text-gray-600 hover:text-gray-900">
          <span>❤️</span>
          <span>128</span>
        </button>
        <button class="flex items-center space-x-2 text-gray-600 hover:text-gray-900">
          <span>💬</span>
          <span>32</span>
        </button>
      </div>
      <button class="text-gray-600 hover:text-gray-900">
        🔗 分享
      </button>
    </div>
  </footer>
</article>
```

### 颜色方案
```css
:root {
  --color-text: #1A1A1A;
  --color-text-secondary: #757575;
  --color-background: #FFFFFF;
  --color-accent: #1A8917;
  --font-serif: 'Noto Serif SC', 'Source Han Serif SC', serif;
  --font-sans: 'Noto Sans SC', 'Source Han Sans SC', sans-serif;
  --content-max-width: 680px;
  --line-height: 1.8;
}

.prose {
  font-family: var(--font-serif);
  font-size: 1.125rem;
  line-height: var(--line-height);
}
```

---

## 主题实现要点

### 1. 主题选择器
```vue
<!-- admin/src/components/ThemeStyleSelector.vue -->
<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    <div v-for="style in styles" :key="style.id"
      @click="selectStyle(style)"
      :class="['p-4 border-2 rounded-xl cursor-pointer transition-all',
        selectedStyle === style.id ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-gray-300']">
      <div class="aspect-video bg-gray-100 rounded-lg mb-3 flex items-center justify-center">
        <span class="text-4xl">{{ style.icon }}</span>
      </div>
      <h3 class="font-semibold">{{ style.name }}</h3>
      <p class="text-sm text-gray-500 mt-1">{{ style.description }}</p>
    </div>
  </div>
</template>

<script setup>
const styles = [
  { id: 'material', name: '现代卡片式', icon: '🃏', description: 'Material Design 风格，大圆角卡片' },
  { id: 'news', name: '新闻门户式', icon: '📰', description: '杂志门户风格，大图标题' },
  { id: 'social', name: '社交动态流', icon: '💬', description: '微博/Twitter 风格时间线' },
  { id: 'forum', name: '论坛经典式', icon: '💭', description: '传统论坛风格，信息密集' },
  { id: 'minimal', name: '极简阅读式', icon: '📖', description: 'Medium 风格，专注阅读' }
]
</script>
```

### 2. 动态组件切换
```javascript
// site/src/stores/theme.js
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const currentStyle = ref('material')
  
  function getStyleComponents(style) {
    switch (style) {
      case 'material':
        return {
          TopicCard: () => import('@/themes/material/components/TopicCard.vue'),
          Home: () => import('@/themes/material/views/Home.vue')
        }
      case 'news':
        return {
          TopicCard: () => import('@/themes/news/components/TopicCard.vue'),
          Home: () => import('@/themes/news/views/Home.vue')
        }
      // ... 其他风格
    }
  }
  
  return { currentStyle, getStyleComponents }
})
```

---

## 实施建议

1. **先实现 2-3 种核心风格**
   - 现代卡片式（默认）
   - 社交动态流
   - 极简阅读式

2. **每个风格独立目录**
   ```
   site/src/themes/
   ├── material/
   ├── news/
   ├── social/
   ├── forum/
   └── minimal/
   ```

3. **共享基础组件**
   - 用户头像、徽章等可复用
   - API 调用逻辑统一

---

*文档版本：v1.0 | 最后更新：2026-04-04*
