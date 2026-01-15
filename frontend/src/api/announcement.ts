import { request } from '@/utils/request'
import type {
  ApiResponse,
  PaginateData,
  Announcement,
  AnnouncementListParams,
  CreateAnnouncementRequest,
  UpdateAnnouncementRequest
} from '@/types'

/**
 * 公告管理相关 API
 */
export const announcementApi = {
  /**
   * 获取公告列表
   */
  getAnnouncements(params: AnnouncementListParams): Promise<ApiResponse<PaginateData<Announcement>>> {
    return request.get('/announcements', { params })
  },

  /**
   * 获取启用的公告列表（用于首页公告栏）
   */
  getActiveAnnouncements(limit: number = 5): Promise<ApiResponse<Announcement[]>> {
    return request.get('/announcements/active', { params: { limit } })
  },

  /**
   * 获取公告详情
   */
  getAnnouncement(id: number): Promise<ApiResponse<Announcement>> {
    return request.get(`/announcements/${id}`)
  },

  /**
   * 创建公告
   */
  createAnnouncement(data: CreateAnnouncementRequest): Promise<ApiResponse<Announcement>> {
    return request.post('/announcements', data)
  },

  /**
   * 更新公告
   */
  updateAnnouncement(id: number, data: UpdateAnnouncementRequest): Promise<ApiResponse<Announcement>> {
    return request.put(`/announcements/${id}`, data)
  },

  /**
   * 删除公告
   */
  deleteAnnouncement(id: number): Promise<ApiResponse<null>> {
    return request.delete(`/announcements/${id}`)
  }
}
