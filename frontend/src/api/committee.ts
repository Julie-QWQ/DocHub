import { request } from '@/utils/request'
import type {
  ApiResponse,
  PaginateData,
  CommitteeApplication,
  CreateCommitteeApplicationRequest,
  CommitteeApplicationListParams,
  ReviewCommitteeApplicationRequest
} from '@/types'

/**
 * 学委申请相关 API
 */
export const committeeApi = {
  /**
   * 申请成为学委
   */
  applyForCommittee(data: CreateCommitteeApplicationRequest): Promise<ApiResponse<CommitteeApplication>> {
    return request.post('/user/apply-committee', data)
  },

  /**
   * 获取我的申请列表
   */
  getMyApplications(params?: CommitteeApplicationListParams): Promise<ApiResponse<PaginateData<CommitteeApplication>>> {
    return request.get('/user/applications', { params })
  },

  /**
   * 获取申请详情
   */
  getApplicationDetail(id: number): Promise<ApiResponse<CommitteeApplication>> {
    return request.get(`/user/applications/${id}`)
  },

  /**
   * 取消申请
   */
  cancelApplication(id: number): Promise<ApiResponse<null>> {
    return request.post(`/user/applications/${id}/cancel`)
  },

  /**
   * 管理员获取所有申请列表
   */
  getAllApplications(params?: CommitteeApplicationListParams): Promise<ApiResponse<PaginateData<CommitteeApplication>>> {
    return request.get('/admin/applications', { params })
  },

  /**
   * 管理员审核学委申请
   */
  reviewApplication(id: number, data: ReviewCommitteeApplicationRequest): Promise<ApiResponse<CommitteeApplication>> {
    return request.post(`/admin/applications/${id}/review`, data)
  },

  /**
   * 获取待审核数量
   */
  getPendingCount(): Promise<ApiResponse<{ count: number }>> {
    return request.get('/admin/applications/pending/count')
  }
}
