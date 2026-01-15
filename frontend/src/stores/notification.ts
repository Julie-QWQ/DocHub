import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { notificationApi } from '@/api/notification'
import type { Notification, NotificationListParams } from '@/types'

export const useNotificationStore = defineStore('notification', () => {
  // ==================== 状态 ====================
  const notifications = ref<Notification[]>([])
  const unreadNotifications = ref<Notification[]>([])
  const unreadCount = ref(0)

  const total = ref(0)
  const page = ref(1)
  const size = ref(20)
  const loading = ref(false)

  // 当前查询参数
  const currentParams = ref<NotificationListParams>({
    page: 1,
    size: 20
  })

  // ==================== 计算属性 ====================
  const hasMore = computed(() => {
    return page.value * size.value < total.value
  })

  const totalPages = computed(() => {
    return Math.ceil(total.value / size.value)
  })

  // ==================== 方法 ====================

  /**
   * 获取通知列表
   */
  const fetchNotifications = async (params?: NotificationListParams) => {
    loading.value = true
    try {
      const queryParams = {
        ...currentParams.value,
        ...params
      }

      const response = await notificationApi.getNotifications(queryParams)

      if (response.code === 0 && response.data) {
        notifications.value = response.data.list || []
        total.value = response.data.total
        page.value = response.data.page
        size.value = response.data.size
        currentParams.value = queryParams
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取未读通知
   */
  const fetchUnreadNotifications = async () => {
    loading.value = true
    try {
      const response = await notificationApi.getUnreadNotifications()

      if (response.code === 0 && response.data) {
        unreadNotifications.value = response.data || []
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取未读数量
   */
  const fetchUnreadCount = async () => {
    const response = await notificationApi.getUnreadCount()

    if (response.code === 0 && typeof response.data === 'number') {
      unreadCount.value = response.data
    }

    return response
  }

  /**
   * 标记单个通知为已读
   */
  const markAsRead = async (id: number) => {
    const response = await notificationApi.markAsRead(id)

    if (response.code === 0) {
      // 更新通知列表中的状态
      const notification = notifications.value.find(n => n.id === id)
      if (notification) {
        notification.status = 'read'
      }

      // 更新未读通知列表
      unreadNotifications.value = unreadNotifications.value.filter(n => n.id !== id)

      // 更新未读数量
      if (unreadCount.value > 0) {
        unreadCount.value--
      }
    }

    return response
  }

  /**
   * 标记所有通知为已读
   */
  const markAllAsRead = async () => {
    const response = await notificationApi.markAllAsRead()

    if (response.code === 0) {
      // 更新所有通知为已读
      notifications.value.forEach(n => {
        n.status = 'read'
      })

      // 清空未读通知
      unreadNotifications.value = []

      // 重置未读数量
      unreadCount.value = 0
    }

    return response
  }

  /**
   * 删除通知
   */
  const deleteNotification = async (id: number) => {
    const response = await notificationApi.deleteNotification(id)

    if (response.code === 0) {
      // 从列表中移除
      notifications.value = notifications.value.filter(n => n.id !== id)
      unreadNotifications.value = unreadNotifications.value.filter(n => n.id !== id)

      // 如果未读则减少计数
      const notification = notifications.value.find(n => n.id === id)
      if (notification && notification.status === 'unread' && unreadCount.value > 0) {
        unreadCount.value--
      }
    }

    return response
  }

  /**
   * 重置状态
   */
  const reset = () => {
    notifications.value = []
    unreadNotifications.value = []
    unreadCount.value = 0
    total.value = 0
    page.value = 1
    size.value = 20
    loading.value = false
    currentParams.value = {
      page: 1,
      size: 20
    }
  }

  return {
    // 状态
    notifications,
    unreadNotifications,
    unreadCount,
    total,
    page,
    size,
    loading,
    currentParams,

    // 计算属性
    hasMore,
    totalPages,

    // 方法
    fetchNotifications,
    fetchUnreadNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    reset
  }
})
