<template>
  <header class="app-header">
    <div class="header-left">
      <h1 class="logo"><SiteName /></h1>
    </div>
    <nav class="header-nav">
      <router-link to="/">首页</router-link>
      <router-link to="/materials">资料</router-link>
      <router-link to="/about">关于</router-link>
    </nav>
    <div class="header-right">
      <template v-if="isLoggedIn">
        <!-- 公告中心 -->
        <div class="announcement-trigger" @click="showAnnouncementPanel = true">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M2 12h5"/>
            <path d="M4 12v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-6"/>
            <path d="M5 12V7a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v5"/>
            <path d="M15 12v-2"/>
            <path d="M15 12h5"/>
            <path d="M18 12v-2"/>
          </svg>
          <span v-if="activeAnnouncementCount > 0" class="badge">{{ activeAnnouncementCount }}</span>
        </div>

        <!-- 通知中心 -->
        <el-popover
          placement="bottom"
          :width="380"
          trigger="click"
          @show="handleNotificationShow"
        >
          <template #reference>
            <NotificationBadge @click="handleNotificationClick" />
          </template>
          <NotificationList
            :show-unread-only="true"
            @click="handleNotificationItemClick"
          />
        </el-popover>

        <el-dropdown>
          <span class="user-name">
            <el-icon><User /></el-icon>
            {{ userInfo?.username }}
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleNotifications">通知中心</el-dropdown-item>
              <el-dropdown-item @click="handleProfile">个人中心</el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
      <template v-else>
        <router-link to="/login" class="btn-login">登录</router-link>
        <router-link to="/register" class="btn-register">注册</router-link>
      </template>
    </div>

    <!-- 公告面板弹窗 -->
    <Transition name="announcement-panel">
      <div v-if="showAnnouncementPanel" class="announcement-panel-overlay" @click.self="showAnnouncementPanel = false">
        <div class="announcement-panel">
          <div class="announcement-panel-header">
            <div class="header-left">
              <svg class="announcement-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M2 12h5"/>
                <path d="M4 12v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-6"/>
                <path d="M5 12V7a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v5"/>
                <path d="M15 12v-2"/>
                <path d="M15 12h5"/>
                <path d="M18 12v-2"/>
              </svg>
              <h3>系统公告</h3>
            </div>
            <button class="close-btn" @click="showAnnouncementPanel = false">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6 6 18"></path>
                <path d="m6 6 12 12"></path>
              </svg>
            </button>
          </div>

          <div class="announcement-panel-body">
            <div v-if="loadingAnnouncements" class="loading-state">
              <div v-for="i in 3" :key="i" class="announcement-skeleton"></div>
            </div>

            <div v-else-if="announcements.length === 0" class="empty-state">
              <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M2 12h5"/>
                <path d="M4 12v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-6"/>
                <path d="M5 12V7a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v5"/>
              </svg>
              <p>暂无公告</p>
            </div>

            <div v-else class="announcements-list">
              <div
                v-for="announcement in announcements"
                :key="announcement.id"
                class="announcement-item"
                @click="viewAnnouncementDetail(announcement)"
              >
                <div class="announcement-item-header">
                  <div class="announcement-title-row">
                    <h4 class="announcement-title">{{ announcement.title }}</h4>
                    <span v-if="announcement.priority === 'high'" class="important-badge">重要</span>
                  </div>
                  <span class="announcement-date">{{ formatAnnouncementDate(announcement.published_at || announcement.created_at) }}</span>
                </div>
                <p class="announcement-preview">{{ getPreviewText(announcement.content) }}</p>
              </div>
            </div>
          </div>

          <div class="announcement-panel-footer">
            <span class="announcement-count">共 {{ announcements.length }} 条公告</span>
          </div>
        </div>
      </div>
    </Transition>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { storage } from '@/utils'
import { useNotification } from '@/composables/useNotification'
import { announcementApi } from '@/api/announcement'
import NotificationBadge from '@/components/NotificationBadge.vue'
import NotificationList from '@/components/NotificationList.vue'
import SiteName from '@/components/SiteName.vue'
import type { Notification } from '@/types'
import type { Announcement } from '@/types'

const router = useRouter()
const { fetchUnreadCount } = useNotification()

const isLoggedIn = computed(() => !!storage.getToken())
const userInfo = computed(() => {
  const userStr = localStorage.getItem('user')
  return userStr ? JSON.parse(userStr) : null
})

// 公告相关状态
const showAnnouncementPanel = ref(false)
const loadingAnnouncements = ref(false)
const announcements = ref<Announcement[]>([])
const activeAnnouncementCount = ref(0)

const handleProfile = () => {
  router.push('/profile')
}

const handleLogout = () => {
  storage.clearAuth()
  ElMessage.success('退出登录成功')
  router.push('/login')
}

const handleNotifications = () => {
  router.push('/notifications')
}

const handleNotificationClick = () => {
  // 徽章点击事件,由 popover 处理
}

const handleNotificationShow = async () => {
  // 显示时刷新未读通知
  await fetchUnreadCount()
}

const handleNotificationItemClick = (notification: Notification) => {
  // 如果有链接,跳转到对应页面
  if (notification.link) {
    router.push(notification.link)
  }
}

// 加载公告列表
const loadAnnouncements = async () => {
  loadingAnnouncements.value = true
  try {
    const response = await announcementApi.getActiveAnnouncements(10)
    console.log('公告API响应:', response)
    announcements.value = response.data || []
    activeAnnouncementCount.value = announcements.value.length
  } catch (error) {
    console.error('加载公告失败:', error)
    ElMessage.error('加载公告失败')
  } finally {
    loadingAnnouncements.value = false
  }
}

// 格式化公告日期
const formatAnnouncementDate = (date: string) => {
  if (!date) return '-'

  const d = new Date(date)
  const now = new Date()

  // 检查日期是否有效
  if (isNaN(d.getTime())) return '-'

  const diff = now.getTime() - d.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  // 如果是未来时间,显示具体日期
  if (days < 0) {
    return d.toLocaleDateString('zh-CN', {
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
  if (days < 30) return `${Math.floor(days / 7)}周前`
  return d.toLocaleDateString('zh-CN')
}

// 获取公告预览文本
const getPreviewText = (content: string) => {
  const plainText = content.replace(/<[^>]*>/g, '')
  return plainText.length > 60 ? plainText.substring(0, 60) + '...' : plainText
}

// 查看公告详情
const viewAnnouncementDetail = (announcement: Announcement) => {
  showAnnouncementPanel.value = false
  ElMessage.info(announcement.content)
}

// 初始化时获取未读通知数量和公告列表
onMounted(async () => {
  if (isLoggedIn.value) {
    await fetchUnreadCount()
    await loadAnnouncements()
  }
})
</script>

<style scoped lang="scss">
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 20px;
  background-color: #ffffff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-left {
  .logo {
    font-size: 24px;
    font-weight: bold;
    color: #409eff;
  }
}

.header-nav {
  display: flex;
  gap: 20px;

  a {
    color: #606266;
    font-size: 14px;

    &:hover,
    &.router-link-active {
      color: #409eff;
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;

  .user-name {
    display: flex;
    align-items: center;
    gap: 5px;
    cursor: pointer;
    color: #606266;
  }

  .btn-login,
  .btn-register {
    padding: 6px 16px;
    border-radius: 4px;
    font-size: 14px;
  }

  .btn-login {
    color: #409eff;
  }

  .btn-register {
    background-color: #409eff;
    color: #ffffff;
  }
}

// 公告触发器
.announcement-trigger {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;
  color: #606266;

  &:hover {
    background: #f5f7fa;
    color: #409eff;
  }

  .badge {
    position: absolute;
    top: 4px;
    right: 4px;
    min-width: 18px;
    height: 18px;
    padding: 0 5px;
    background: #f56c6c;
    color: #ffffff;
    font-size: 12px;
    font-weight: 600;
    border-radius: 9px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid #ffffff;
  }
}

// 公告面板弹窗
.announcement-panel-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(2px);
}

.announcement-panel {
  width: 420px;
  max-height: 600px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  animation: slideIn 0.3s ease-out;

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
}

.announcement-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #ffffff;

  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;

    .announcement-icon {
      flex-shrink: 0;
    }

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #ffffff;
    }
  }

  .close-btn {
    background: rgba(255, 255, 255, 0.2);
    border: none;
    width: 32px;
    height: 32px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: #ffffff;
    transition: all 0.2s;
    flex-shrink: 0;

    &:hover {
      background: rgba(255, 255, 255, 0.3);
    }
  }
}

.announcement-panel-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: #dcdfe6;
    border-radius: 3px;

    &:hover {
      background: #c0c4cc;
    }
  }
}

.loading-state {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.announcement-skeleton {
  height: 80px;
  background: linear-gradient(90deg, #f5f5f5 25%, #e8e8e8 50%, #f5f5f5 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 8px;

  @keyframes shimmer {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #909399;

  svg {
    margin-bottom: 12px;
    opacity: 0.5;
  }

  p {
    margin: 0;
    font-size: 14px;
  }
}

.announcements-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.announcement-item {
  padding: 16px;
  background: #fafafa;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #ffffff;
    border-color: #667eea;
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
    transform: translateY(-1px);
  }
}

.announcement-item-header {
  margin-bottom: 8px;
}

.announcement-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;

  .announcement-title {
    margin: 0;
    font-size: 15px;
    font-weight: 600;
    color: #111827;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .important-badge {
    flex-shrink: 0;
    padding: 2px 8px;
    background: #fef3c7;
    color: #b45309;
    font-size: 12px;
    font-weight: 600;
    border-radius: 4px;
  }
}

.announcement-date {
  font-size: 12px;
  color: #9ca3af;
}

.announcement-preview {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.announcement-panel-footer {
  padding: 12px 20px;
  background: #fafafa;
  border-top: 1px solid #e5e7eb;

  .announcement-count {
    font-size: 13px;
    color: #909399;
  }
}

// 弹窗过渡动画
.announcement-panel-enter-active,
.announcement-panel-leave-active {
  transition: all 0.3s ease;
}

.announcement-panel-enter-from,
.announcement-panel-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
