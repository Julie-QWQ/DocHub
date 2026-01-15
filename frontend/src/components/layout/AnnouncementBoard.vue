<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { announcementApi } from '@/api/announcement'
import type { Announcement } from '@/types'

// 公告数据
const announcements = ref<Announcement[]>([])
const loading = ref(false)
const showDetailDialog = ref(false)
const selectedAnnouncement = ref<Announcement | null>(null)

// 加载公告
const loadAnnouncements = async () => {
  loading.value = true
  try {
    const response = await announcementApi.getActiveAnnouncements(5)
    if (response.code === 0 && response.data) {
      announcements.value = response.data
    }
  } catch (error: any) {
    // 静默失败，不影响用户体验
    console.error('加载公告失败:', error)
  } finally {
    loading.value = false
  }
}

const formatTime = (dateStr: string) => {
  if (!dateStr) return ''

  const date = new Date(dateStr)
  const now = new Date()

  // 检查日期是否有效
  if (isNaN(date.getTime())) return ''

  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  // 如果是未来时间,显示具体日期
  if (days < 0) {
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

// 点击公告查看详情
const handleAnnouncementClick = (announcement: Announcement) => {
  selectedAnnouncement.value = announcement
  showDetailDialog.value = true
}

// 关闭详情弹窗
const closeDetailDialog = () => {
  showDetailDialog.value = false
  selectedAnnouncement.value = null
}

// 组件挂载时加载公告
onMounted(() => {
  loadAnnouncements()
})
</script>

<template>
  <aside class="announcement-board">
    <div class="board-header">
      <h3 class="board-title">公告栏</h3>
    </div>

    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="announcements.length === 0" class="empty-container">
      <p>暂无公告</p>
    </div>

    <div v-else class="announcement-list">
      <article
        v-for="announcement in announcements"
        :key="announcement.id"
        class="announcement-item"
        @click="handleAnnouncementClick(announcement)"
      >
        <div class="announcement-header">
          <h4 class="announcement-title">
            <span
              v-if="announcement.priority === 'high'"
              class="priority-badge"
            >
              重要
            </span>
            {{ announcement.title }}
          </h4>
          <span class="announcement-time">{{ formatTime(announcement.published_at || announcement.created_at) }}</span>
        </div>
        <p class="announcement-content">{{ announcement.content }}</p>
      </article>
    </div>

    <div class="board-footer">
      <router-link to="/announcements" class="view-all-link">查看全部公告 →</router-link>
    </div>
  </aside>

  <!-- 公告详情弹窗 -->
  <Transition name="dialog">
    <div v-if="showDetailDialog" class="dialog-overlay" @click.self="closeDetailDialog">
      <div class="dialog-content" @click.stop>
        <div class="dialog-header">
          <h3 class="dialog-title">
            <span v-if="selectedAnnouncement?.priority === 'high'" class="priority-badge">重要</span>
            {{ selectedAnnouncement?.title }}
          </h3>
          <button class="close-btn" @click="closeDetailDialog">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>
        <div class="dialog-body">
          <div class="announcement-meta">
            <span class="meta-item">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
              {{ formatTime(selectedAnnouncement?.published_at || selectedAnnouncement?.created_at || '') }}
            </span>
            <span v-if="selectedAnnouncement?.author" class="meta-item">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              {{ selectedAnnouncement.author.real_name || selectedAnnouncement.author.username }}
            </span>
          </div>
          <div class="announcement-full-content" v-html="selectedAnnouncement?.content"></div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped lang="scss">
.announcement-board {
  background: #fafafa;
  border-radius: 12px;
  padding: 20px;
  position: sticky;
  top: 24px;
}

.board-header {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e5e5e5;

  .board-title {
    font-size: 16px;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0;
  }
}

.loading-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #999;
  min-height: 200px;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #f3f4f6;
  border-top-color: #1a1a1a;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 16px;
}

.announcement-item {
  background: #ffffff;
  border: 1px solid #f2f2f2;
  border-radius: 8px;
  padding: 16px;
  transition: all 0.2s;
  cursor: pointer;

  &:hover {
    border-color: #667eea;
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
    transform: translateY(-2px);
  }
}

.announcement-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
  gap: 8px;
}

.announcement-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
}

.priority-badge {
  padding: 2px 6px;
  background: #fef3c7;
  color: #d97706;
  font-size: 11px;
  font-weight: 600;
  border-radius: 4px;
  white-space: nowrap;
}

.announcement-time {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
}

.announcement-content {
  font-size: 13px;
  color: #666;
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.board-footer {
  padding-top: 12px;
  border-top: 1px solid #e5e5e5;
  text-align: center;
}

.view-all-link {
  font-size: 13px;
  color: #666;
  text-decoration: none;
  transition: color 0.2s;

  &:hover {
    color: #1a1a1a;
  }
}

// 响应式
@media (max-width: 1200px) {
  .announcement-board {
    display: none;
  }
}

// 弹窗样式
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  padding: 20px;
  animation: fadeIn 0.3s ease;

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
}

.dialog-content {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9ff 100%);
  border-radius: 20px;
  box-shadow:
    0 20px 60px rgba(0, 0, 0, 0.15),
    0 8px 24px rgba(102, 126, 234, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
  max-width: 700px;
  width: 100%;
  max-height: 85vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  animation: dialogSlideIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  position: relative;

  // 顶部装饰条
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  }

  @keyframes dialogSlideIn {
    0% {
      opacity: 0;
      transform: translateY(-40px) scale(0.9);
    }
    100% {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }
}

.dialog-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 28px 32px 24px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08) 0%, rgba(118, 75, 162, 0.06) 100%);
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  gap: 20px;
  position: relative;

  // 装饰性背景图案
  &::before {
    content: '';
    position: absolute;
    top: -50%;
    right: -20%;
    width: 200px;
    height: 200px;
    background: radial-gradient(circle, rgba(102, 126, 234, 0.06) 0%, transparent 70%);
    border-radius: 50%;
    pointer-events: none;
  }
}

.dialog-title {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  line-height: 1.3;
  color: #1a1a2e;
  position: relative;
  z-index: 1;

  .priority-badge {
    padding: 4px 12px;
    background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
    color: #ffffff;
    font-size: 12px;
    font-weight: 700;
    border-radius: 20px;
    white-space: nowrap;
    box-shadow: 0 2px 8px rgba(238, 90, 111, 0.3);
    letter-spacing: 0.5px;
  }
}

.close-btn {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.08) 100%);
  border: 1px solid rgba(102, 126, 234, 0.15);
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #667eea;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
  position: relative;
  z-index: 1;

  svg {
    width: 20px;
    height: 20px;
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  &:hover {
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.15) 100%);
    border-color: rgba(102, 126, 234, 0.3);
    transform: scale(1.05);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);

    svg {
      transform: rotate(90deg);
    }
  }

  &:active {
    transform: scale(0.95);
  }
}

.dialog-body {
  padding: 28px 32px;
  overflow-y: auto;
  flex: 1;

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-track {
    background: rgba(102, 126, 234, 0.05);
    border-radius: 4px;
    margin: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 4px;
    transition: background 0.3s;

    &:hover {
      background: linear-gradient(135deg, #5568d3 0%, #6a4a8f 100%);
    }
  }
}

.announcement-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 24px;
  padding: 18px 20px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.04) 0%, rgba(118, 75, 162, 0.03) 100%);
  border: 1px solid rgba(102, 126, 234, 0.08);
  border-radius: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #667eea;
  background: rgba(255, 255, 255, 0.8);
  padding: 8px 14px;
  border-radius: 8px;
  transition: all 0.2s;

  svg {
    flex-shrink: 0;
    width: 16px;
    height: 16px;
  }

  &:hover {
    background: rgba(102, 126, 234, 0.1);
    transform: translateY(-1px);
  }
}

.announcement-full-content {
  font-size: 15px;
  color: #2d3748;
  line-height: 1.8;
  font-weight: 400;

  :deep p {
    margin-bottom: 16px;
    position: relative;
    padding-left: 16px;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 10px;
      bottom: 10px;
      width: 3px;
      background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
      border-radius: 2px;
    }

    &:last-child {
      margin-bottom: 0;
    }
  }

  :deep strong {
    color: #667eea;
    font-weight: 600;
  }

  :deep a {
    color: #667eea;
    text-decoration: none;
    border-bottom: 1px dashed #667eea;
    transition: all 0.2s;

    &:hover {
      color: #764ba2;
      border-bottom-style: solid;
    }
  }

  :deep ul, ::v-deep ol {
    margin: 16px 0;
    padding-left: 28px;
  }

  :deep li {
    margin-bottom: 10px;
    line-height: 1.8;
    position: relative;

    &::marker {
      color: #667eea;
      font-weight: 600;
    }
  }

  :deep code {
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 0.9em;
    font-family: 'Courier New', monospace;
  }

  :deep blockquote {
    margin: 16px 0;
    padding: 16px 20px;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.06) 0%, rgba(118, 75, 162, 0.04) 100%);
    border-left: 4px solid #667eea;
    border-radius: 0 8px 8px 0;
    font-style: italic;
  }
}

// 弹窗过渡动画
.dialog-enter-active,
.dialog-leave-active {
  transition: all 0.3s ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-from .dialog-content,
.dialog-leave-to .dialog-content {
  transform: translateY(-20px) scale(0.95);
}
</style>
