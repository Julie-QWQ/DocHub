import { request } from '@/utils/request'
import type {
  ApiResponse,
  PaginateData,
  ReviewRecord,
  ReviewHistoryParams,
  HandleReportRequest,
  Material
} from '@/types'

/**
 * 审核相关 API
 */
export const reviewApi = {
  /**
   * 获取待审核资料列表
   */
  getPendingMaterials(params?: { page?: number; size?: number }): Promise<ApiResponse<PaginateData<Material>>> {
    return request.get('/admin/materials/pending', { params })
  },

  /**
   * 审核资料
   */
  reviewMaterial(id: number, approved: boolean, comment?: string): Promise<ApiResponse<Material>> {
    return request.post(`/admin/materials/${id}/review`, {
      status: approved ? 'approved' : 'rejected',
      rejection_reason: comment
    })
  },

  /**
   * 处理举报
   */
  handleReport(id: number, data: HandleReportRequest): Promise<ApiResponse<null>> {
    return request.post(`/admin/reports/${id}/handle`, data)
  },

  /**
   * 获取审核历史
   */
  getReviewHistory(params?: ReviewHistoryParams): Promise<ApiResponse<PaginateData<ReviewRecord>>> {
    return request.get('/admin/review/history', { params })
  },

  /**
   * 获取审核人统计
   */
  getReviewerStatistics(reviewerId: number): Promise<ApiResponse<{
    total_reviews: number
    approved_count: number
    rejected_count: number
    material_count: number
    committee_count: number
    report_count: number
  }>> {
    return request.get(`/admin/reviewers/${reviewerId}/statistics`)
  }
}
