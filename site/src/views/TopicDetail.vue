<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 bg-white rounded-lg shadow-sm">
    <div v-if="topic" class="mb-6">
      <div class="flex items-start justify-between mb-4">
        <h1 class="text-2xl font-bold text-gray-900">{{ topic.title }}</h1>
        <button v-if="canDeleteTopic" @click="handleDeleteTopic"
          class="flex items-center space-x-1 px-3 py-1.5 text-sm text-red-600 hover:text-white hover:bg-red-500 border border-red-300 rounded-lg transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
            </path>
          </svg>
          <span>删除</span>
        </button>
      </div>
      <div v-if="topic.tags && topic.tags.length > 0" class="flex items-center flex-wrap gap-2 mb-4">
        <router-link v-for="tag in topic.tags" :key="tag.id" :to="`/?tag=${tag.id}`"
          class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded-full hover:bg-blue-200">
          #{{ tag.name }}
        </router-link>
      </div>
      <div class="flex items-center space-x-4 mb-6 pb-6 border-b">
        <router-link :to="`/user/${topic.user_id}`">
          <img :src="getUserAvatar(topic.user)" class="w-12 h-12 rounded-full">
        </router-link>
        <div class="flex-1">
          <div class="flex items-center gap-2">
            <router-link :to="`/user/${topic.user_id}`" class="font-medium text-gray-900 hover:text-blue-500">{{
              getUserDisplayName(topic.user) }}</router-link>
            <div v-if="displayAuthorBadges.length > 0" class="flex items-center gap-1">
              <SvgBadge v-for="badge in displayAuthorBadges" :key="badge.id"
                :type="badge.icon" :size="24" :title="badge.name" />
            </div>
          </div>
          <div class="text-sm text-gray-500">{{ formatTime(topic.created_at) }} · {{ topic.view_count }} 浏览</div>
        </div>
      </div>
      <div class="prose max-w-none mb-6 topic-content" v-html="renderMarkdown(topic.content)"></div>

      <div v-if="poll && configStore.state.allow_poll"
        class="mb-6 p-6 bg-gradient-to-r from-blue-50 to-purple-50 rounded-xl border border-blue-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-bold text-gray-900">{{ poll.title || topic.title }}</h3>
          <span v-if="isPollEnded"
            class="px-3 py-1 text-xs font-medium bg-gray-200 text-gray-600 rounded-full">已结束</span>
          <span v-else-if="poll.end_time" class="text-sm text-gray-500">
            剩余 {{ getRemainingTime(poll.end_time) }}
          </span>
        </div>

        <div v-if="!hasVoted && !isPollEnded" class="space-y-3">
          <div v-for="option in poll.options" :key="option.id"
            @click="poll.poll_type === 'single' && selectOption(option.id)"
            :class="['p-4 rounded-lg border-2 cursor-pointer transition-all',
              selectedOptions.includes(option.id) ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-blue-300 bg-white']">
            <div class="flex items-center">
              <input v-if="poll.poll_type === 'single'" type="radio" :checked="selectedOptions.includes(option.id)"
                class="w-4 h-4 text-blue-600" @click.stop>
              <input v-else type="checkbox" :checked="selectedOptions.includes(option.id)"
                @change="toggleOption(option.id)" @click.stop class="w-4 h-4 text-blue-600 rounded">
              <span class="ml-3 text-gray-700">{{ option.text }}</span>
            </div>
          </div>

          <div class="flex items-center justify-between pt-4">
            <span class="text-sm text-gray-500">
              {{ poll.poll_type === 'single' ? '单选' : `多选，最多选${poll.max_choices}项` }}
            </span>
            <button @click="submitVote" :disabled="selectedOptions.length === 0 || submittingVote"
              class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
              {{ submittingVote ? '提交中...' : '提交投票' }}
            </button>
          </div>
        </div>

        <div v-else class="space-y-3">
          <div v-for="option in poll.options" :key="option.id" class="relative">
            <div :class="['p-4 rounded-lg border-2 transition-all',
              votedOptionIds.includes(option.id) ? 'border-blue-500 bg-blue-50' : 'border-gray-200 bg-white']">
              <div class="flex items-center justify-between mb-2">
                <span class="text-gray-700 font-medium">{{ option.text }}</span>
                <span class="text-sm font-medium text-gray-600">
                  {{ getPercentage(option.vote_count) }}%
                  <span v-if="votedOptionIds.includes(option.id)" class="ml-2 text-blue-500">✓ 已选</span>
                </span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                <div class="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full transition-all duration-500"
                  :style="{ width: getPercentage(option.vote_count) + '%' }"></div>
              </div>
              <div class="text-xs text-gray-500 mt-1">{{ option.vote_count }} 票</div>
            </div>
          </div>

          <div class="text-center text-sm text-gray-500 pt-2">
            共 {{ poll.total_votes }} 票
            <span v-if="hasVoted" class="ml-2 text-blue-500">· 您已投票</span>
          </div>
        </div>
      </div>

      <div class="flex items-center space-x-4 pt-4 border-t">
        <button @click="toggleLike"
          :class="['flex items-center space-x-2 transition-colors', liked ? 'text-red-500' : 'text-gray-500 hover:text-red-500']">
          <svg class="w-5 h-5" :fill="liked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z">
            </path>
          </svg>
          <span>{{ topic.like_count }}</span>
        </button>
        <button @click="toggleFavorite"
          :class="['flex items-center space-x-2 transition-colors', favorited ? 'text-yellow-500' : 'text-gray-500 hover:text-yellow-500']">
          <svg class="w-5 h-5" :fill="favorited ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"></path>
          </svg>
          <span>{{ favorited ? '已收藏' : '收藏' }}</span>
        </button>
        <button @click="shareTopic"
          class="flex items-center space-x-2 text-gray-500 hover:text-green-500 transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z">
            </path>
          </svg>
          <span>分享</span>
        </button>
        <button @click="openReportDialog('topic', topic.id, topic.title)"
          class="flex items-center space-x-2 text-gray-500 hover:text-red-500 transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9">
            </path>
          </svg>
          <span>举报</span>
        </button>
      </div>
    </div>
    <div class="mt-8">
      <h3 class="text-lg font-medium text-gray-900 mb-4">{{ posts.length }} 条评论</h3>
      <div v-if="userStore.isLoggedIn && configStore.state.allow_comment && topic?.allow_comment" class="mb-6">
        <div v-if="replyTo" class="mb-2 text-sm text-gray-600 flex items-center">
          <span>回复 @{{ replyToUser }}</span>
          <button @click="cancelReply" class="ml-2 text-red-500 hover:text-red-600">取消</button>
        </div>
        <textarea v-model="newPost" rows="3"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
          placeholder="写下你的评论..."></textarea>
        <div class="flex justify-end mt-2">
          <button @click="submitPost"
            class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600">发表评论</button>
        </div>
      </div>
      <div v-else-if="!configStore.state.allow_comment"
        class="mb-6 p-4 bg-gray-100 rounded-lg text-center text-gray-500">
        评论功能已关闭
      </div>
      <div v-else-if="topic && !topic.allow_comment" class="mb-6 p-4 bg-gray-100 rounded-lg text-center text-gray-500">
        本话题已关闭评论
      </div>
      <div class="space-y-4">
        <div v-for="post in sortedPosts" :key="post.id" :id="'post-' + post.id"
          :class="[
            'flex space-x-4 p-4 rounded-lg transition-all',
            post.is_best ? 'bg-gradient-to-r from-yellow-50 to-orange-50 border-2 border-yellow-300 shadow-md' : 'bg-gray-50'
          ]">
          <img :src="getUserAvatar(post.user)" class="w-10 h-10 rounded-full">
          <div class="flex-1">
            <div class="flex items-center justify-between mb-1">
              <div class="flex items-center space-x-2">
                <span v-if="post.is_best" class="px-2 py-0.5 text-xs font-bold bg-yellow-500 text-white rounded-full animate-pulse">最佳</span>
                <span v-if="post.is_pinned" class="text-xs text-red-500 font-medium">置顶</span>
                <SvgBadge v-if="post.is_best" type="gold-comment" :size="20" title="最佳评论" />
                <span class="font-medium text-gray-900">{{ getUserDisplayName(post.user) }}</span>
                <span v-if="post.reply_user" class="text-gray-400 text-sm">回复 @{{ post.reply_user.nickname || post.reply_user.username }}</span>
                <span v-if="post.user_id === topic?.user_id" class="px-1.5 py-0.5 text-xs bg-red-500 text-white rounded">楼主</span>
                <div v-if="getCommentAuthorTopBadge(post)" class="flex items-center gap-0.5 ml-1">
                  <SvgBadge :type="getCommentAuthorTopBadge(post).icon" :size="16" :title="getCommentAuthorTopBadge(post).name" />
                </div>
                <span class="text-sm text-gray-500">{{ formatTime(post.created_at) }}</span>
              </div>
              <div class="flex gap-2">
                <button v-if="canBestComment(post)" @click="toggleCommentBest(post)"
                  :class="['text-xs transition-colors', post.is_best ? 'text-yellow-500 hover:text-yellow-600' : 'text-gray-400 hover:text-yellow-500']">
                  {{ post.is_best ? '取消最佳' : '标记最佳' }}
                </button>
                <button v-if="canPinComment(post)" @click="toggleCommentPin(post)"
                  :class="['text-xs transition-colors', post.is_pinned ? 'text-red-500 hover:text-red-600' : 'text-gray-400 hover:text-red-500']">
                  {{ post.is_pinned ? '取消置顶' : '置顶' }}
                </button>
                <button v-if="canDeletePost(post)" @click="handleDeletePost(post)"
                  class="text-xs text-gray-400 hover:text-red-500 transition-colors ml-2">删除</button>
                <button v-if="canReportPost(post)"
                  @click="openReportDialog('comment', post.id, post.content.substring(0, 50))"
                  class="text-xs text-gray-400 hover:text-red-500 transition-colors ml-2">举报</button>
              </div>
            </div>

            <p class="text-gray-700">{{ post.content }}</p>
            <div class="flex items-center space-x-4 mt-2 text-sm">
              <button @click="togglePostLike(post)"
                :class="['transition-colors', getPostLiked(post.id) ? 'text-red-500' : 'text-gray-500 hover:text-red-500']">
                {{ getPostLiked(post.id) ? '❤️' : '🤍' }} {{ post.like_count }}
              </button>
              <button @click="openReply(post)" class="text-gray-500 hover:text-blue-500">回复</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片查看器 -->
    <div v-if="showLightbox" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-90"
      @click="closeLightbox">
      <button @click="closeLightbox"
        class="absolute top-4 right-4 text-white text-4xl hover:text-gray-300 transition-colors z-10">
        ×
      </button>
      <img :src="lightboxImage" class="max-w-full max-h-full object-contain" @click.stop>
    </div>

    <!-- 分享对话框 -->
    <el-dialog v-model="shareDialogVisible" title="分享" width="480px" :close-on-click-modal="true">
      <div class="share-dialog">
        <!-- 链接显示和复制 -->
        <div class="share-link-section mb-6">
          <div class="flex items-center bg-gray-50 rounded-lg p-4">
            <div class="flex-1 mr-4">
              <p class="text-sm text-gray-500 mb-1">分享链接</p>
              <p class="text-sm text-gray-800 truncate">{{ shareUrl }}</p>
            </div>
            <button @click="copyShareLink"
              class="flex items-center px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3">
                </path>
              </svg>
              {{ copied ? '已复制' : '复制' }}
            </button>
          </div>
        </div>

        <!-- 分享方式 -->
        <div class="share-methods">
          <p class="text-sm text-gray-500 mb-4">选择分享方式</p>
          <div class="grid grid-cols-4 gap-4">
            <!-- 微信 -->
            <div @click="shareToWechat"
              class="share-item cursor-pointer flex flex-col items-center p-4 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center mb-2">
                <svg class="w-7 h-7 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path
                    d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178A1.17 1.17 0 0 1 4.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 0 1-1.162 1.178 1.17 1.17 0 0 1-1.162-1.178c0-.651.52-1.18 1.162-1.18zm5.34 2.867c-1.797-.052-3.746.512-5.28 1.786-1.72 1.428-2.687 3.72-1.78 6.22.942 2.453 3.666 4.229 6.884 4.229.826 0 1.622-.12 2.361-.336a.722.722 0 0 1 .598.082l1.584.926a.272.272 0 0 0 .14.047c.134 0 .24-.111.24-.247 0-.06-.023-.12-.038-.177l-.327-1.233a.582.582 0 0 1-.023-.156.49.49 0 0 1 .201-.398C23.024 18.48 24 16.82 24 14.98c0-3.21-2.931-5.837-6.656-6.088V8.89c-.135-.015-.264-.032-.406-.032zm-2.53 3.274c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.97-.982zm4.844 0c.535 0 .969.44.969.982a.976.976 0 0 1-.969.983.976.976 0 0 1-.969-.983c0-.542.434-.982.969-.982z" />
                </svg>
              </div>
              <span class="text-sm text-gray-700">微信</span>
            </div>

            <!-- 微信朋友圈 -->
            <div @click="shareToMoments"
              class="share-item cursor-pointer flex flex-col items-center p-4 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="w-12 h-12 bg-green-600 rounded-full flex items-center justify-center mb-2">
                <svg class="w-7 h-7 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path
                    d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" />
                </svg>
              </div>
              <span class="text-sm text-gray-700">朋友圈</span>
            </div>

            <!-- 微博 -->
            <div @click="shareToWeibo"
              class="share-item cursor-pointer flex flex-col items-center p-4 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="w-12 h-12 bg-red-500 rounded-full flex items-center justify-center mb-2">
                <svg class="w-7 h-7 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path
                    d="M10.098 20.323c-3.977.391-7.414-1.406-7.672-4.02-.259-2.609 2.759-5.047 6.74-5.441 3.979-.394 7.413 1.404 7.671 4.018.259 2.6-2.759 5.049-6.739 5.443zm-.521-1.262c2.706-.271 4.631-1.663 4.809-3.118.178-1.457-1.592-2.904-4.301-2.636-2.71.269-4.633 1.665-4.812 3.124-.178 1.456 1.59 2.901 4.304 2.63zm.51-3.012c-1.23-.122-2.429.37-2.682 1.096-.25.723.463 1.417 1.69 1.536 1.225.12 2.425-.371 2.676-1.095.254-.727-.46-1.419-1.684-1.537zm.082-.971c-.465-.046-.914.138-1.013.412-.098.275.174.54.642.588.466.048.916-.137 1.015-.414.1-.276-.172-.54-.644-.586zm.379.493c-.146-.015-.287.045-.32.133-.032.088.053.173.2.187.147.016.288-.043.321-.133.032-.086-.054-.172-.201-.187z" />
                </svg>
              </div>
              <span class="text-sm text-gray-700">微博</span>
            </div>

            <!-- QQ -->
            <div @click="shareToQQ"
              class="share-item cursor-pointer flex flex-col items-center p-4 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center mb-2">
                <svg class="w-7 h-7 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path
                    d="M12.003 2c-2.265 0-6.29 1.364-6.29 7.325v1.195S3.55 14.96 3.55 17.474c0 .665.17 1.025.28 1.025.114 0 .85-.365.85-.365s-.062.48.307.928c.37.445 1.035.742 1.035.742s-.426.19-.426.633c0 .44.38.82.846.82.468 0 2.216-.285 3.558-.532 1.345-.25 2.61-.25 3.95 0 1.34.247 3.088.532 3.555.532.467 0 .846-.38.846-.82 0-.443-.427-.633-.427-.633s.664-.297 1.035-.742c.369-.448.307-.928.307-.928s.736.365.85.365c.11 0 .28-.36.28-1.025 0-2.514-2.164-6.954-2.164-6.954V9.325C18.293 3.364 14.268 2 12.003 2z" />
                </svg>
              </div>
              <span class="text-sm text-gray-700">QQ</span>
            </div>

            <!-- 二维码 -->
            <div @click="showQRCode"
              class="share-item cursor-pointer flex flex-col items-center p-4 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="w-12 h-12 bg-gray-800 rounded-full flex items-center justify-center mb-2">
                <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z">
                  </path>
                </svg>
              </div>
              <span class="text-sm text-gray-700">二维码</span>
            </div>
          </div>
        </div>

        <!-- 二维码显示 -->
        <div v-if="qrcodeVisible" class="qrcode-section mt-6 pt-6 border-t">
          <p class="text-sm text-gray-500 mb-4 text-center">扫描二维码分享</p>
          <div class="flex justify-center">
            <div class="bg-white p-4 rounded-lg shadow-md">
              <div id="qrcode" class="w-48 h-48"></div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 举报对话框 -->
    <el-dialog v-model="reportDialogVisible" title="举报内容" width="420px" :close-on-click-modal="false">
      <div class="report-dialog">
        <p class="report-tip">请选择举报原因：</p>
        <el-radio-group v-model="selectedReportReason" class="report-reasons">
          <el-radio value="spam">🚫 垃圾广告</el-radio>
          <el-radio value="illegal">🔞 违规内容</el-radio>
          <el-radio value="attack">😡 人身攻击</el-radio>
          <el-radio value="rumor">📰 谣言虚假信息</el-radio>
          <el-radio value="other">➕ 其他</el-radio>
        </el-radio-group>
        <div class="report-detail">
          <p class="report-tip">补充说明（可选）：</p>
          <el-input v-model="reportDetail" type="textarea" :rows="3" maxlength="500" show-word-limit
            placeholder="请详细描述问题..." />
        </div>
      </div>
      <template #footer>
        <el-button @click="closeReportDialog">取消</el-button>
        <el-button type="danger" @click="submitReport" :loading="reportSubmitting" :disabled="!selectedReportReason">
          提交举报
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useConfigStore } from '@/stores/config'
import api, { pollApi, topicApi, commentApi, commentPinApi, reportApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserAvatar, getUserDisplayName } from '@/utils/user'
import { renderMarkdown } from '@/utils/markdown'
import { getDisplayBadges } from '@/utils/badge'
import QRCode from 'qrcode'
import SvgBadge from '@/components/SvgBadge.vue'

const route = useRoute()
const userStore = useUserStore()
const configStore = useConfigStore()
const topic = ref(null)
const posts = ref([])
const newPost = ref('')
const liked = ref(false)
const favorited = ref(false)
const postLikes = ref({})
const showLightbox = ref(false)
const lightboxImage = ref('')
const authorBadges = ref([])

const poll = ref(null)
const selectedOptions = ref([])
const submittingVote = ref(false)
const votedOptionIds = ref([])
const hasVotedFromServer = ref(false)

// 举报相关
const reportDialogVisible = ref(false)
const reportTargetType = ref('')
const reportTargetId = ref(null)
const reportTargetTitle = ref('')
const selectedReportReason = ref('')
const reportDetail = ref('')
const reportSubmitting = ref(false)

// 回复相关
const replyTo = ref(null)
const replyToUser = ref(null)

// 分享相关
const shareDialogVisible = ref(false)
const shareUrl = ref('')
const copied = ref(false)
const qrcodeVisible = ref(false)

const hasVoted = computed(() => hasVotedFromServer.value || votedOptionIds.value.length > 0)
const isPollEnded = computed(() => {
  if (!poll.value?.end_time) return false
  return new Date(poll.value.end_time) < new Date()
})

// 作者勋章展示（取前3个，按优先级排序）
const displayAuthorBadges = computed(() => {
  return getDisplayBadges(authorBadges.value, 'author-card')
})

const canDeleteTopic = computed(() => {
  if (!userStore.isLoggedIn || !topic.value) return false
  const isAuthor = topic.value.user_id === userStore.user?.id
  const isAdmin = userStore.user?.role === 2
  return isAuthor || isAdmin
})

const sortedPosts = computed(() => {
  return [...posts.value].sort((a, b) => {
    if (a.is_best && !b.is_best) return -1
    if (!a.is_best && b.is_best) return 1
    if (a.is_pinned && !b.is_pinned) return -1
    if (!a.is_pinned && b.is_pinned) return 1
    return new Date(b.created_at) - new Date(a.created_at)
  })
})

function canDeletePost(post) {
  if (!userStore.isLoggedIn) return false
  const isAuthor = post.user_id === userStore.user?.id
  const isAdmin = userStore.user?.role === 2
  return isAuthor || isAdmin
}

function canReportPost(post) {
  if (!userStore.isLoggedIn) return false
  // 不能举报自己
  if (post.user_id === userStore.user?.id) return false
  return true
}

function canPinComment(post) {
  // 只有帖子作者可以置顶评论
  if (!userStore.isLoggedIn || !topic.value) return false
  return topic.value.user_id === userStore.user?.id
}

function canBestComment(post) {
  // 只有帖子作者可以标记最佳评论
  if (!userStore.isLoggedIn || !topic.value) return false
  return topic.value.user_id === userStore.user?.id
}

// 获取评论作者的最高优先级勋章（用于显示）
function getCommentAuthorTopBadge(post) {
  if (!post.user?.badges || post.user.badges.length === 0) return null
  const topBadges = getDisplayBadges(post.user.badges, 'comment')
  return topBadges.length > 0 ? topBadges[0] : null
}

async function toggleCommentBest(post) {
  try {
    await ElMessageBox.confirm(
      post.is_best ? '确定要取消最佳评论吗？' : '确定要标记为最佳评论吗？',
      post.is_best ? '取消最佳' : '标记最佳',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await commentPinApi.bestComment(topic.value.id, post.id, !post.is_best)
    post.is_best = !post.is_best
    ElMessage.success(post.is_best ? '已标记为最佳评论' : '已取消最佳评论')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('操作失败', e)
      ElMessage.error('操作失败')
    }
  }
}

async function toggleCommentPin(post) {
  try {
    await ElMessageBox.confirm(
      post.is_pinned ? '确定要取消置顶这条评论吗？' : '确定要置顶这条评论吗？',
      post.is_pinned ? '取消置顶' : '置顶评论',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await commentPinApi.pinComment(topic.value.id, post.id, !post.is_pinned)

    // 记录原始索引
    const originalIndex = posts.value.findIndex(p => p.id === post.id)
    if (originalIndex === -1) return

    // 记录是否原本已置顶
    const wasPinned = post.is_pinned

    // 更新置顶状态
    post.is_pinned = !post.is_pinned

    // 从原位置移除
    posts.value.splice(originalIndex, 1)

    if (post.is_pinned) {
      // 置顶：插入到置顶评论之后（最前面）
      let insertIndex = 0
      for (let i = 0; i < posts.value.length; i++) {
        if (!posts.value[i].is_pinned) {
          insertIndex = i
          break
        }
        insertIndex = i + 1
      }
      posts.value.splice(insertIndex, 0, post)
    } else {
      // 取消置顶：移到非置顶评论之后
      let insertIndex = posts.value.length
      for (let i = 0; i < posts.value.length; i++) {
        if (!posts.value[i].is_pinned) {
          insertIndex = i
          break
        }
      }
      posts.value.splice(insertIndex, 0, post)
    }

    ElMessage.success(post.is_pinned ? '评论已置顶' : '评论已取消置顶')
  } catch (e) {
    if (e !== 'cancel') {
      console.error('操作失败', e)
      ElMessage.error('操作失败')
    }
  }
}

function openReportDialog(targetType, targetId, targetTitle) {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }
  reportTargetType.value = targetType
  reportTargetId.value = targetId
  reportTargetTitle.value = targetTitle
  selectedReportReason.value = ''
  reportDetail.value = ''
  reportDialogVisible.value = true
}

function closeReportDialog() {
  reportDialogVisible.value = false
}

async function submitReport() {
  if (!selectedReportReason.value) {
    ElMessage.warning('请选择举报原因')
    return
  }

  reportSubmitting.value = true
  try {
    await reportApi.createReport({
      target_type: reportTargetType.value,
      target_id: reportTargetId.value,
      reason: selectedReportReason.value,
      detail: reportDetail.value
    })
    ElMessage.success('举报已提交，感谢您的反馈')
    closeReportDialog()
  } catch (e) {
    console.error('举报失败', e)
    ElMessage.error(e.response?.data?.message || '举报失败')
  } finally {
    reportSubmitting.value = false
  }
}

async function handleDeletePost(post) {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '删除评论', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
  } catch {
    return
  }

  try {
    await commentApi.deleteComment(post.id)
    posts.value = posts.value.filter(p => p.id !== post.id)
    ElMessage.success('评论已删除')
  } catch (e) {
    console.error('删除失败', e)
    ElMessage.error(e.response?.data?.message || '删除失败')
  }
}

function getPercentage(voteCount) {
  if (!poll.value || poll.value.total_votes === 0) return 0
  return Math.round((voteCount / poll.value.total_votes) * 100)
}

function getRemainingTime(endTime) {
  const end = new Date(endTime)
  const now = new Date()
  const diff = end - now

  if (diff <= 0) return '已结束'

  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))

  if (days > 0) return `${days}天${hours}小时`
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  if (hours > 0) return `${hours}小时${minutes}分钟`
  return `${minutes}分钟`
}

function selectOption(optionId) {
  selectedOptions.value = [optionId]
}

function toggleOption(optionId) {
  const index = selectedOptions.value.indexOf(optionId)
  if (index > -1) {
    selectedOptions.value.splice(index, 1)
  } else {
    if (selectedOptions.value.length < poll.value.max_choices) {
      selectedOptions.value.push(optionId)
    } else {
      ElMessage.warning(`最多只能选择${poll.value.max_choices}项`)
    }
  }
}

async function loadPoll() {
  try {
    const res = await pollApi.getPollByTopic(route.params.id)
    console.log('loadPoll response:', res)
    if (res && res.poll) {
      poll.value = res.poll
      hasVotedFromServer.value = res.has_voted || false
      votedOptionIds.value = res.voted_option_ids || []
      console.log('hasVotedFromServer:', hasVotedFromServer.value, 'votedOptionIds:', votedOptionIds.value)
    }
  } catch (e) {
    console.log('No poll for this topic')
  }
}

async function handleDeleteTopic() {
  try {
    await ElMessageBox.confirm('确定要删除这篇帖子吗？删除后无法恢复。', '删除帖子', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
  } catch {
    return
  }

  try {
    await topicApi.deleteTopic(topic.value.id)
    ElMessage.success('帖子已删除')
    window.location.href = '/'
  } catch (e) {
    console.error('删除失败', e)
    ElMessage.error(e.response?.data?.message || '删除失败')
  }
}

async function submitVote() {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  if (selectedOptions.value.length === 0) {
    ElMessage.warning('请选择投票选项')
    return
  }

  submittingVote.value = true
  try {
    const res = await pollApi.submitVote({
      poll_id: poll.value.id,
      option_ids: selectedOptions.value
    })
    console.log('submitVote response:', res)

    // 立即更新状态
    hasVotedFromServer.value = true
    votedOptionIds.value = [...selectedOptions.value]

    // 再刷新获取最新数据
    await loadPoll()
    ElMessage.success('投票成功')
  } catch (e) {
    console.error(e)
    ElMessage.error(e.response?.data?.message || '投票失败')
  } finally {
    submittingVote.value = false
  }
}

function formatTime(time) {
  const date = new Date(time)
  const now = new Date()
  const diff = now - date
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  return Math.floor(diff / 86400000) + '天前'
}

async function loadTopic() {
  try {
    const id = route.params.id
    topic.value = await api.get(`/topics/${id}`)
    const postsRes = await api.get(`/topics/${id}/comments`)
    posts.value = postsRes.list || postsRes || []

    if (userStore.isLoggedIn) {
      await checkLikeStatus()
      await checkFavoriteStatus()
    }

    await loadPoll()
    
    // 加载作者勋章
    if (topic.value?.user_id) {
      await loadAuthorBadges()
    }

    // 设置媒体查看器
    setupMediaViewers()
  } catch (e) {
    console.error(e)
  }
}

async function loadAuthorBadges() {
  try {
    if (!topic.value?.user_id) return
    const res = await api.get(`/users/${topic.value.user_id}/badges`)
    authorBadges.value = res || []
  } catch (e) {
    console.error('加载作者勋章失败', e)
  }
}

async function checkLikeStatus() {
  try {
    const res = await api.post('/likes/check', {
      target_type: 'topic',
      target_id: topic.value.id
    })
    liked.value = res.liked
  } catch (e) {
    console.error(e)
  }
}

async function checkFavoriteStatus() {
  try {
    const res = await api.post('/favorites/check', {
      topic_id: topic.value.id
    })
    favorited.value = res.favorited
  } catch (e) {
    console.error(e)
  }
}

async function toggleLike() {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (liked.value) {
      await api.delete(`/likes?target_type=topic&target_id=${topic.value.id}`)
      topic.value.like_count--
    } else {
      await api.post('/likes', {
        target_type: 'topic',
        target_id: topic.value.id
      })
      topic.value.like_count++
    }
    liked.value = !liked.value
  } catch (e) {
    console.error(e)
    ElMessage.error('操作失败')
  }
}

async function toggleFavorite() {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (favorited.value) {
      await api.delete(`/favorites?topic_id=${topic.value.id}`)
    } else {
      await api.post('/favorites', {
        topic_id: topic.value.id
      })
    }
    favorited.value = !favorited.value
  } catch (e) {
    console.error(e)
    ElMessage.error('操作失败')
  }
}

function shareTopic() {
  shareUrl.value = window.location.href
  copied.value = false
  qrcodeVisible.value = false
  shareDialogVisible.value = true
}

async function copyShareLink() {
  try {
    await navigator.clipboard.writeText(shareUrl.value)
    copied.value = true
    ElMessage.success('链接已复制到剪贴板')
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (e) {
    console.error(e)
    ElMessage.error('复制失败，请手动复制')
  }
}

function isWechatBrowser() {
  const ua = navigator.userAgent.toLowerCase()
  return ua.includes('micromessenger')
}

function isMobile() {
  return /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(navigator.userAgent)
}

function shareToWechat() {
  if (isWechatBrowser()) {
    ElMessage.info('请点击右上角菜单按钮分享给好友')
  } else if (isMobile()) {
    copyShareLink()
    setTimeout(() => {
      ElMessage.success('链接已复制，正在打开微信...')
      setTimeout(() => {
        window.location.href = 'weixin://'
      }, 500)
    }, 300)
  } else {
    ElMessage.info('请使用微信扫描二维码分享')
    showQRCode()
  }
}

function shareToMoments() {
  if (isWechatBrowser()) {
    ElMessage.info('请点击右上角菜单按钮分享到朋友圈')
  } else if (isMobile()) {
    copyShareLink()
    setTimeout(() => {
      ElMessage.success('链接已复制，正在打开微信...')
      setTimeout(() => {
        window.location.href = 'weixin://'
      }, 500)
    }, 300)
  } else {
    ElMessage.info('请使用微信扫描二维码分享到朋友圈')
    showQRCode()
  }
}

function shareToWeibo() {
  const url = `https://service.weibo.com/share/share.php?url=${encodeURIComponent(shareUrl.value)}&title=${encodeURIComponent(topic.value?.title || '')}`
  window.open(url, '_blank', 'width=600,height=400')
}

function shareToQQ() {
  const url = `https://connect.qq.com/widget/shareqq/index.html?url=${encodeURIComponent(shareUrl.value)}&title=${encodeURIComponent(topic.value?.title || '')}`
  window.open(url, '_blank', 'width=600,height=400')
}

function showQRCode() {
  qrcodeVisible.value = true
  nextTick(() => {
    const qrcodeEl = document.getElementById('qrcode')
    if (qrcodeEl) {
      qrcodeEl.innerHTML = ''
      try {
        const canvas = document.createElement('canvas')
        qrcodeEl.appendChild(canvas)
        QRCode.toCanvas(canvas, shareUrl.value, {
          width: 192,
          height: 192,
          margin: 1,
          color: {
            dark: '#000000',
            light: '#ffffff'
          }
        }, function (error) {
          if (error) {
            console.error(error)
            qrcodeEl.innerHTML = '<div class="flex items-center justify-center h-full text-gray-500 text-sm">二维码生成失败</div>'
          }
        })
      } catch (e) {
        console.error(e)
        qrcodeEl.innerHTML = '<div class="flex items-center justify-center h-full text-gray-500 text-sm">二维码生成失败</div>'
      }
    }
  })
}

function getPostLiked(postId) {
  return postLikes.value[postId] || false
}

async function togglePostLike(post) {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (getPostLiked(post.id)) {
      await api.delete(`/likes?target_type=comment&target_id=${post.id}`)
      post.like_count--
    } else {
      await api.post('/likes', {
        target_type: 'comment',
        target_id: post.id
      })
      post.like_count++
    }
    postLikes.value[post.id] = !getPostLiked(post.id)
  } catch (e) {
    console.error(e)
    ElMessage.error('操作失败')
  }
}

async function submitPost() {
  if (!newPost.value.trim()) return
  try {
    const payload = { content: newPost.value }
    if (replyTo.value) {
      payload.reply_to_id = replyTo.value
    }
    await api.post(`/topics/${route.params.id}/comments`, payload)
    newPost.value = ''
    replyTo.value = null
    replyToUser.value = null
    loadTopic()
  } catch (e) {
    console.error(e)
  }
}

function openReply(post) {
  replyTo.value = post.id
  replyToUser.value = getUserDisplayName(post.user)
  // 滚动到评论输入框
  document.querySelector('.topic-detail-container')?.scrollIntoView({ behavior: 'smooth' })
}

function cancelReply() {
  replyTo.value = null
  replyToUser.value = null
}

function openLightbox(src) {
  lightboxImage.value = src
  showLightbox.value = true
  document.body.style.overflow = 'hidden'
}

function closeLightbox() {
  showLightbox.value = false
  document.body.style.overflow = ''
}

function setupMediaViewers() {
  setTimeout(() => {
    const content = document.querySelector('.topic-content')
    if (!content) return

    // 为图片添加点击放大功能
    const images = content.querySelectorAll('img')
    images.forEach(img => {
      img.style.cursor = 'pointer'
      img.style.transition = 'transform 0.2s, box-shadow 0.2s'

      img.addEventListener('mouseenter', () => {
        img.style.transform = 'scale(1.02)'
        img.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)'
      })

      img.addEventListener('mouseleave', () => {
        img.style.transform = ''
        img.style.boxShadow = ''
      })

      img.addEventListener('click', () => {
        openLightbox(img.src)
      })
    })

    // 为视频添加全屏播放功能
    const videos = content.querySelectorAll('video')
    videos.forEach(video => {
      video.style.cursor = 'pointer'
      video.style.transition = 'transform 0.2s, box-shadow 0.2s'

      video.addEventListener('mouseenter', () => {
        video.style.transform = 'scale(1.02)'
        video.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)'
      })

      video.addEventListener('mouseleave', () => {
        video.style.transform = ''
        video.style.boxShadow = ''
      })

      video.addEventListener('click', () => {
        if (video.requestFullscreen) {
          video.requestFullscreen()
        } else if (video.webkitRequestFullscreen) {
          video.webkitRequestFullscreen()
        } else if (video.msRequestFullscreen) {
          video.msRequestFullscreen()
        }
      })
    })
  }, 100)
}

onMounted(() => {
  loadTopic()
})
</script>

<style scoped>
.topic-content img,
.topic-content video {
  border-radius: 8px;
  margin: 1rem 0;
  max-width: 100%;
  height: auto;
}

.topic-content img:hover,
.topic-content video:hover {
  position: relative;
}

.topic-content img::after,
.topic-content video::after {
  content: '🔍';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 24px;
  opacity: 0;
  transition: opacity 0.2s;
}

.topic-content img:hover::after,
.topic-content video:hover::after {
  opacity: 1;
}
</style>
