import { request } from '@/utils/request'
import type {
  ApiResponse,
  PaginateData,
  Notification,
  NotificationListParams
} from '@/types'

/**
 * 通知相关 API
 */
export const notificationApi = {
  /**
   * 获取通知列表
   */
  getNotifications(params?: NotificationListParams): Promise<ApiResponse<PaginateData<Notification>>> {
    return request.get('/notifications', { params })
  },

  /**
   * 获取未读通知
   */
  getUnreadNotifications(): Promise<ApiResponse<Notification[]>> {
    return request.get('/notifications/unread')
  },

  /**
   * 获取未读数量
   */
  getUnreadCount(): Promise<ApiResponse<number>> {
    return request.get('/notifications/unread/count')
  },

  /**
   * 标记单个通知为已读
   */
  markAsRead(id: number): Promise<ApiResponse<null>> {
    return request.post(`/notifications/${id}/read`)
  },

  /**
   * 标记所有通知为已读
   */
  markAllAsRead(): Promise<ApiResponse<null>> {
    return request.post('/notifications/read-all')
  },

  /**
   * 删除通知
   */
  deleteNotification(id: number): Promise<ApiResponse<null>> {
    return request.delete(`/notifications/${id}`)
  }
}
