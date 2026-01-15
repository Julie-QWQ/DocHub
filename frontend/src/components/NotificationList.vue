<template>
  <div class="notification-list">
    <div class="notification-header">
      <h3>通知</h3>
      <el-button
        v-if="unreadCount > 0"
        link
        type="primary"
        @click="handleMarkAllRead"
      >
        全部标为已读
      </el-button>
    </div>

    <el-scrollbar max-height="400px">
      <el-empty v-if="notifications.length === 0" description="暂无通知" />

      <div v-else class="notification-items">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-item"
          :class="{ unread: notification.status === 'unread' }"
          @click="handleClick(notification)"
        >
          <el-icon class="notification-icon" :size="20">
            <component :is="getIcon(notification.type)" />
          </el-icon>

          <div class="notification-content">
            <div class="notification-title">{{ notification.title }}</div>
            <div class="notification-desc">{{ notification.content }}</div>
            <div class="notification-time">
              {{ formatNotificationTime(notification.created_at) }}
            </div>
          </div>

          <div v-if="notification.status === 'unread'" class="unread-dot" />
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { Bell, Document, User, Warning } from '@element-plus/icons-vue'
import { useNotification } from '@/composables/useNotification'
import type { Notification } from '@/types'

const props = defineProps<{
  showUnreadOnly?: boolean
}>()

const emit = defineEmits<{
  click: [notification: Notification]
}>()

const {
  notifications,
  unreadNotifications,
  unreadCount,
  fetchNotifications,
  fetchUnreadNotifications,
  markAllAsRead,
  formatNotificationTime
} = useNotification()

// 根据属性决定显示哪些通知
const displayNotifications = computed(() => {
  return props.showUnreadOnly ? unreadNotifications.value : notifications.value
})

const getIcon = (type: string) => {
  const iconMap: Record<string, any> = {
    system: Bell,
    material: Document,
    committee: User,
    report: Warning
  }
  return iconMap[type] || Bell
}

const handleClick = async (notification: Notification) => {
  emit('click', notification)

  // 如果是未读通知,标记为已读
  if (notification.status === 'unread') {
    await markAsRead(notification.id)
  }
}

const handleMarkAllRead = async () => {
  await markAllAsRead()
}

onMounted(async () => {
  if (props.showUnreadOnly) {
    await fetchUnreadNotifications()
  } else {
    await fetchNotifications()
  }
})
</script>

<style scoped lang="scss">
.notification-list {
  width: 380px;

  .notification-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    border-bottom: 1px solid var(--el-border-color-light);

    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
    }
  }

  .notification-items {
    .notification-item {
      display: flex;
      align-items: flex-start;
      padding: 12px 16px;
      cursor: pointer;
      transition: background-color 0.2s;
      position: relative;

      &:hover {
        background-color: var(--el-fill-color-light);
      }

      &.unread {
        background-color: var(--el-color-info-light-9);
      }

      .notification-icon {
        margin-right: 12px;
        margin-top: 2px;
        color: var(--el-color-primary);
      }

      .notification-content {
        flex: 1;
        min-width: 0;

        .notification-title {
          font-size: 14px;
          font-weight: 500;
          margin-bottom: 4px;
          color: var(--el-text-color-primary);
        }

        .notification-desc {
          font-size: 13px;
          color: var(--el-text-color-regular);
          margin-bottom: 4px;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }

        .notification-time {
          font-size: 12px;
          color: var(--el-text-color-secondary);
        }
      }

      .unread-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background-color: var(--el-color-danger);
        margin-left: 8px;
        margin-top: 6px;
        flex-shrink: 0;
      }
    }
  }
}
</style>
