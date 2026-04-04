# BBSGo Theme Style Design (Fully Differentiated)

## Table of Contents
- [Theme Style Comparison](#theme-style-comparison)
- [Style 1: Modern Card (Material Design)](#style-1-modern-card-material-design)
- [Style 2: News Portal](#style-2-news-portal)
- [Style 3: Social Feed](#style-3-social-feed)
- [Style 4: Classic Forum](#style-4-classic-forum)
- [Style 5: Minimal Reading](#style-5-minimal-reading)

---

## Theme Style Comparison

| Feature | Modern Card | News Portal | Social Feed | Classic Forum | Minimal Reading |
|---------|-------------|-------------|-------------|---------------|-----------------|
| **Layout** | Three-column cards | Two-column magazine | Single-column waterfall | Classic forum | Centered single-column |
| **Cards** | Large rounded shadows | Image-text mix | Left-right aligned | Table rows | Borderless |
| **Navigation** | Fixed top | Top menu | Bottom tabs | Side navigation | Minimal top |
| **Interaction** | Hover animations | Large image titles | Likes/comments | Floor replies | Focus on content |
| **Best For** | General community | Content portal | Youth community | Tech forum | Blog/reading |

---

## Style 1: Modern Card (Material Design)

### Design Philosophy
- Google Material Design style
- Large rounded corners, soft shadows
- Clear hierarchy
- Smooth interactive animations

### Homepage Layout
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]  [Search]                [Nav]  [Avatar]  [Notif] │ ← Top navbar
├──────────┬──────────────────────────────────────────┬───────────┤
│          │  [Tags]  All | Hot | Latest | Featured     │           │
│  [Hot]   │  ┌────────────────────────────────────┐  │  [Active] │
│  Topics  │  │ 🟦  [Pinned]  Post big title        │  │  Users    │
│         │  │  Summary text...                       │  │           │
│  [Hot]   │  │  [Image] [Image] [Image]            │  │  [Hot]    │
│  Users   │  └────────────────────────────────────┘  │  Topics   │
│         │  ┌────────────────────────────────────┐  │           │
│         │  │ 🟦  Post big title                   │  │           │
│         │  │  Summary text...                       │  │           │
│         │  └────────────────────────────────────┘  │           │
└──────────┴──────────────────────────────────────────┴───────────┘
```

### Card Style
```vue
<!-- Large rounded corners, large shadows -->
<div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-300">
  <div class="flex items-center space-x-4 mb-4">
    <img class="w-12 h-12 rounded-full" src="avatar">
    <div>
      <div class="font-semibold">Username</div>
      <div class="text-sm text-gray-500">2 hours ago</div>
    </div>
  </div>
  <h2 class="text-xl font-bold mb-3">Post big title</h2>
  <p class="text-gray-600 mb-4">Summary content...</p>
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

### Color Scheme
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

## Style 2: News Portal

### Design Philosophy
- Magazine/news portal style
- Large image + title visual impact
- High information density
- Multi-column layout

### Homepage Layout
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]          [Home] [Tech] [Life] [Entertainment] [More]     [Login] │
├─────────────────────────────────────────────────────────────────┤
│  📰 [Carousel]  Top headline 1              [Sidebar Hot]  Title 1        │
│     Title big title                    Title 2                 │
│                                       Title 3                 │
│  ┌──────────────┬───────────────────────────────────────────┐  │
│  │ [Thumb] Title 2 │ [Thumb] Title 5       [Thumb] Title 8       │  │
│  │  Summary...     │  Summary...           Summary...           │  │
│  ├──────────────┼───────────────────────────────────────────┤  │
│  │ [Thumb] Title 3 │ [Thumb] Title 6       [Thumb] Title 9       │  │
│  │  Summary...     │  Summary...           Summary...           │  │
│  ├──────────────┼───────────────────────────────────────────┤  │
│  │ [Thumb] Title 4 │ [Thumb] Title 7       [Thumb] Title 10      │  │
│  │  Summary...     │  Summary...           Summary...           │  │
│  └──────────────┴───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### Article Card Style
```vue
<!-- Headline big card -->
<div class="relative overflow-hidden rounded-lg">
  <img class="w-full h-64 object-cover" src="cover-image">
  <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-6">
    <span class="inline-block bg-red-600 text-white px-3 py-1 rounded text-sm mb-2">Headline</span>
    <h1 class="text-2xl font-bold text-white">This is the big headline</h1>
    <p class="text-gray-200 mt-2">Summary content...</p>
  </div>
</div>

<!-- List small card -->
<div class="flex space-x-4 p-4 border-b hover:bg-gray-50">
  <img class="w-24 h-24 object-cover rounded" src="thumbnail">
  <div class="flex-1">
    <h3 class="font-bold hover:text-blue-600">Article title</h3>
    <p class="text-sm text-gray-500 mt-1 line-clamp-2">Summary content...</p>
    <div class="flex items-center space-x-3 mt-2 text-xs text-gray-400">
      <span>Source</span>
      <span>·</span>
      <span>2 hours ago</span>
      <span>·</span>
      <span>1.2k views</span>
    </div>
  </div>
</div>
```

### Color Scheme
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

## Style 3: Social Feed

### Design Philosophy
- Timeline similar to Weibo, Twitter
- Left-right columns, content centered
- Emphasize interaction (likes, reposts, comments)
- Mobile-friendly

### Homepage Layout
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]  [🏠 Home] [🔍 Explore] [✏️ Post]  [🔔] [👤]        │ ← Top
├──────────┬──────────────────────────┬───────────────────────────┤
│          │  [Post box] Share thoughts...     │  [My Profile]              │
│  [Menu]  │  ┌─────────────────────┐  │  [Avatar] Username        │
│  Home    │  │ 🟦  [Avatar] Username  │  │  Following 128 | Followers 256 │
│  Explore │  │  Feed content...        │  │                         │
│  Messages│  │  [Image] [Image]     │  │  [Recommended follows]            │
│  Notifs  │  │  ❤️ 128  🔄 32  💬  │  │  [Avatar] Username  [+]   │
│          │  └─────────────────────┘  │  [Avatar] Username  [+]   │
│          │  ┌─────────────────────┐  │  [Avatar] Username  [+]   │
│          │  │ 🟦  [Avatar] Username  │  │                         │
│          │  │  Feed content...        │  │  [Hot topics]            │
│          │  │  ❤️ 256  🔄 64  💬  │  │  #Topic1  12.5k       │
│          │  └─────────────────────┘  │  #Topic2   8.2k       │
└──────────┴──────────────────────────┴───────────────────────────┘
```

### Feed Card Style
```vue
<div class="border-b border-gray-200 p-4 hover:bg-gray-50">
  <div class="flex space-x-3">
    <img class="w-12 h-12 rounded-full" src="avatar">
    <div class="flex-1 min-w-0">
      <div class="flex items-center space-x-2">
        <span class="font-bold">Username</span>
        <span class="text-gray-500">@username</span>
        <span class="text-gray-400">·</span>
        <span class="text-gray-500">2h</span>
      </div>
      <p class="mt-1 text-gray-900 whitespace-pre-wrap">
        Feed content, can be multi-line, supports line breaks
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

### Color Scheme
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

## Style 4: Classic Forum

### Design Philosophy
- Traditional forum (Discourse, phpBB) style
- Tabular list, high information density
- Floor replies, clear discussion tree
- Sidebar navigation

### Homepage Layout
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]                    [User] [Messages] [Settings] [Logout]        │
├──────────┬──────────────────────────────────────────────────────┤
│          │  ┌─────────────────────────────────────────────┐    │
│  [Forums]│  │ [Icon]  Forum name                Topics | Posts │    │
│  Tech    │  ├─────────────────────────────────────────────┤    │
│  Discuss │  │ [Pinned] [Featured]  Post title                │    │
│  Share   │  │  [Avatar]  Author · 2h ago · 32 replies · 1.2k  │    │
│  Feedback│  ├─────────────────────────────────────────────┤    │
│          │  │ [Featured]  Post title                        │    │
│  [Online]│  │  [Avatar]  Author · 4h ago · 64 replies · 2.5k  │    │
│  Users   │  ├─────────────────────────────────────────────┤    │
│  [Avatar]│  │  Post title                               │    │
│  [Avatar]│  │  [Avatar]  Author · 6h ago · 8 replies · 320    │    │
│          │  └─────────────────────────────────────────────┘    │
│          │  [Pagination]  1 2 3 4 5 ...  Next                      │
└──────────┴──────────────────────────────────────────────────────┘
```

### Post List Style
```vue
<!-- Tabular list -->
<table class="w-full">
  <thead class="bg-gray-100">
    <tr>
      <th class="px-4 py-3 text-left">Topic</th>
      <th class="px-4 py-3 text-center w-24">Replies</th>
      <th class="px-4 py-3 text-center w-24">Views</th>
      <th class="px-4 py-3 text-left w-48">Last Update</th>
    </tr>
  </thead>
  <tbody>
    <tr class="border-b hover:bg-gray-50">
      <td class="px-4 py-4">
        <div class="flex items-start space-x-3">
          <img class="w-10 h-10 rounded-full" src="avatar">
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2">
              <span v-if="topic.isPinned" class="bg-red-100 text-red-700 px-2 py-0.5 rounded text-xs">Pinned</span>
              <span v-if="topic.isEssence" class="bg-yellow-100 text-yellow-700 px-2 py-0.5 rounded text-xs">Featured</span>
              <a href="#" class="font-medium hover:text-blue-600">Post title</a>
            </div>
            <div class="text-sm text-gray-500 mt-1">
              Author · Posted 2 hours ago
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
            <div>Last replier</div>
            <div class="text-gray-500">30 minutes ago</div>
          </div>
        </div>
      </td>
    </tr>
  </tbody>
</table>
```

### Color Scheme
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

## Style 5: Minimal Reading

### Design Philosophy
- Reading experience similar to Medium, Jianshu
- Large whitespace, focus on content
- Centered narrow column, comfortable reading
- No distracting elements

### Homepage Layout
```
┌─────────────────────────────────────────────────────────────────┐
│  [Logo]                              [Search] [User]         │ ← Minimal
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│                    ┌─────────────────────┐                      │
│                    │  Featured articles           │                      │
│                    ├─────────────────────┤                      │
│                    │  Article big title         │                      │
│                    │                     │                      │
│                    │  Summary content...        │                      │
│                    │                     │                      │
│                    │  [Author Avatar]  Author name  │                      │
│                    │  2026-04-04 · 5 min read │                    │
│                    └─────────────────────┘                      │
│                    ┌─────────────────────┐                      │
│                    │  Article big title         │                      │
│                    │                     │                      │
│                    │  Summary content...        │                      │
│                    │                     │                      │
│                    │  [Author Avatar]  Author name  │                      │
│                    └─────────────────────┘                      │
│                    ┌─────────────────────┐                      │
│                    │  Article big title         │                      │
│                    └─────────────────────┘                      │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Article Card Style
```vue
<article class="max-w-2xl mx-auto py-12">
  <header class="mb-8">
    <div class="flex items-center space-x-3 mb-4">
      <img class="w-12 h-12 rounded-full" src="authorAvatar">
      <div>
        <div class="font-medium">Author name</div>
        <div class="text-sm text-gray-500">April 4, 2026 · 5 min read</div>
      </div>
    </div>
    <h1 class="text-3xl font-bold leading-tight">Article big title</h1>
  </header>
  
  <div v-if="coverImage" class="mb-8">
    <img :src="coverImage" class="w-full rounded-lg">
  </div>
  
  <div class="prose prose-lg prose-gray max-w-none">
    <p>First paragraph content...</p>
    <p>Second paragraph content...</p>
    <blockquote>Quote content</blockquote>
    <pre><code>Code block</code></pre>
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
        🔗 Share
      </button>
    </div>
  </footer>
</article>
```

### Color Scheme
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

## Theme Implementation Points

### 1. Theme Selector
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
  { id: 'material', name: 'Modern Card', icon: '🃏', description: 'Material Design style, large rounded cards' },
  { id: 'news', name: 'News Portal', icon: '📰', description: 'Magazine portal style, large image titles' },
  { id: 'social', name: 'Social Feed', icon: '💬', description: 'Weibo/Twitter style timeline' },
  { id: 'forum', name: 'Classic Forum', icon: '💭', description: 'Traditional forum style, information dense' },
  { id: 'minimal', name: 'Minimal Reading', icon: '📖', description: 'Medium style, focus on reading' }
]
</script>
```

### 2. Dynamic Component Switching
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
      // ... other styles
    }
  }
  
  return { currentStyle, getStyleComponents }
})
```

---

## Implementation Suggestions

1. **Implement 2-3 core styles first**
   - Modern Card (default)
   - Social Feed
   - Minimal Reading

2. **Each style in independent directory**
   ```
   site/src/themes/
   ├── material/
   ├── news/
   ├── social/
   ├── forum/
   └── minimal/
   ```

3. **Share base components**
   - User avatars, badges, etc. can be reused
   - Unified API call logic

---

*Document version: v1.0 | Last updated: 2026-04-04*
