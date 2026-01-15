import { request } from '@/utils/request'
import type {
  ApiResponse,
  MaterialCategoryConfig,
  MaterialCategoryConfigRequest
} from '@/types'

/**
 * 资料类型配置相关 API
 */
export const materialCategoryApi = {
  /**
   * 获取资料类型列表
   */
  list(activeOnly = false): Promise<ApiResponse<MaterialCategoryConfig[]>> {
    return request.get('/material-categories', {
      params: { active_only: activeOnly }
    })
  },

  /**
   * 获取资料类型详情
   */
  get(id: number): Promise<ApiResponse<MaterialCategoryConfig>> {
    return request.get(`/material-categories/${id}`)
  },

  /**
   * 创建资料类型（管理员权限）
   */
  create(data: MaterialCategoryConfigRequest): Promise<ApiResponse<MaterialCategoryConfig>> {
    return request.post('/admin/material-categories', data)
  },

  /**
   * 更新资料类型（管理员权限）
   */
  update(id: number, data: MaterialCategoryConfigRequest): Promise<ApiResponse<MaterialCategoryConfig>> {
    return request.put(`/admin/material-categories/${id}`, data)
  },

  /**
   * 删除资料类型（管理员权限）
   */
  delete(id: number): Promise<ApiResponse<null>> {
    return request.delete(`/admin/material-categories/${id}`)
  },

  /**
   * 切换资料类型启用状态（管理员权限）
   */
  toggleStatus(id: number): Promise<ApiResponse<MaterialCategoryConfig>> {
    return request.post(`/admin/material-categories/${id}/toggle`)
  }
}
