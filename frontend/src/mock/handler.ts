/**
 * Mock 请求处理器
 * 当 USE_MOCK=true 时，使用 Mock 数据代替真实 API
 */

import {
  mockMaterials,
  mockNotifications,
  mockDownloadRecords,
  mockFavorites,
  mockStatistics,
  mockPaginate,
  mockSuccess
} from './index'

// 模拟延迟
const delay = (ms: number = 500) => new Promise(resolve => setTimeout(resolve, ms))

// Mock 处理器映射
const mockHandlers = {
  // 获取资料列表
  'GET /api/v1/materials': async (params: any) => {
    await delay()
    const page = params.page || 1
    const size = params.size || 20
    return mockPaginate(mockMaterials, page, size)
  },

  // 获取资料详情
  'GET /api/v1/materials/:id': async (params: any) => {
    await delay()
    const material = mockMaterials.find(m => m.id === parseInt(params.id))
    if (material) {
      return mockSuccess(material)
    }
    return { code: 10004, message: '资料不存在', data: null }
  },

  // 获取通知列表
  'GET /api/v1/notifications': async () => {
    await delay()
    return mockSuccess(mockNotifications)
  },

  // 获取下载记录
  'GET /api/v1/downloads': async () => {
    await delay()
    return mockSuccess(mockDownloadRecords)
  },

  // 获取收藏列表
  'GET /api/v1/favorites': async () => {
    await delay()
    return mockSuccess(mockFavorites)
  },

  // 获取统计数据
  'GET /api/v1/statistics/overview': async () => {
    await delay()
    return mockSuccess(mockStatistics.overview)
  },

  // 获取资料统计
  'GET /api/v1/statistics/materials': async () => {
    await delay()
    return mockSuccess(mockStatistics.materials)
  },

  // 获取下载统计
  'GET /api/v1/statistics/downloads': async () => {
    await delay()
    return mockSuccess(mockStatistics.downloads)
  }
}

/**
 * 处理 Mock 请求
 */
export async function handleMockRequest(
  method: string,
  url: string,
  params?: any
): Promise<any> {
  // 移除查询参数
  const path = url.split('?')[0]

  // 构建处理器 key
  const key = `${method} ${path}`

  // 查找处理器
  const handler = mockHandlers[key as keyof typeof mockHandlers]

  if (handler) {
    return await handler(params)
  }

  // 未找到处理器
  console.warn(`Mock handler not found for: ${key}`)
  return { code: 10004, message: '接口未定义', data: null }
}

/**
 * 检查是否使用 Mock
 */
export function shouldUseMock(): boolean {
  return import.meta.env.USE_MOCK === 'true'
}
