import { computed } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import type { NotificationListParams } from '@/types'

/**
 * 通知相关的组合式函数
 */
export function useNotification() {
  const notificationStore = useNotificationStore()

  /**
   * 获取通知列表
   */
  const fetchNotifications = async (params?: NotificationListParams) => {
    return await notificationStore.fetchNotifications(params)
  }

  /**
   * 获取未读通知
   */
  const fetchUnreadNotifications = async () => {
    return await notificationStore.fetchUnreadNotifications()
  }

  /**
   * 获取未读数量
   */
  const fetchUnreadCount = async () => {
    return await notificationStore.fetchUnreadCount()
  }

  /**
   * 标记通知为已读
   */
  const markAsRead = async (id: number) => {
    return await notificationStore.markAsRead(id)
  }

  /**
   * 标记所有通知为已读
   */
  const markAllAsRead = async () => {
    return await notificationStore.markAllAsRead()
  }

  /**
   * 删除通知
   */
  const deleteNotification = async (id: number) => {
    return await notificationStore.deleteNotification(id)
  }

  /**
   * 格式化通知类型
   */
  const getNotificationTypeText = (type: string) => {
    const typeMap: Record<string, string> = {
      system: '系统通知',
      material: '资料审核',
      committee: '学委申请',
      report: '举报处理'
    }
    return typeMap[type] || type
  }

  /**
   * 获取通知类型对应的图标
   */
  const getNotificationIcon = (type: string) => {
    const iconMap: Record<string, string> = {
      system: 'Bell',
      material: 'Document',
      committee: 'User',
      report: 'Warning'
    }
    return iconMap[type] || 'Bell'
  }

  /**
   * 格式化通知时间
   */
  const formatNotificationTime = (createdAt: string) => {
    const now = new Date()
    const created = new Date(createdAt)
    const diff = now.getTime() - created.getTime()

    const minute = 60 * 1000
    const hour = 60 * minute
    const day = 24 * hour

    if (diff < minute) {
      return '刚刚'
    } else if (diff < hour) {
      return `${Math.floor(diff / minute)}分钟前`
    } else if (diff < day) {
      return `${Math.floor(diff / hour)}小时前`
    } else if (diff < 7 * day) {
      return `${Math.floor(diff / day)}天前`
    } else {
      return created.toLocaleDateString('zh-CN')
    }
  }

  return {
    // 状态
    notifications: computed(() => notificationStore.notifications),
    unreadNotifications: computed(() => notificationStore.unreadNotifications),
    unreadCount: computed(() => notificationStore.unreadCount),
    loading: computed(() => notificationStore.loading),
    hasMore: computed(() => notificationStore.hasMore),
    totalPages: computed(() => notificationStore.totalPages),

    // 方法
    fetchNotifications,
    fetchUnreadNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    deleteNotification,

    // 工具方法
    getNotificationTypeText,
    getNotificationIcon,
    formatNotificationTime
  }
}
