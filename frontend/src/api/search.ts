import request from '@/utils/request'
import type { ApiResponse, PageResponse } from '@/types/api'

/**
 * 资料分类枚举
 */
export enum MaterialCategory {
  COURSEWARE = 'courseware',
  EXAM = 'exam',
  EXPERIMENT = 'experiment',
  EXERCISE = 'exercise',
  REFERENCE = 'reference',
  OTHER = 'other'
}

/**
 * 资料状态枚举
 */
export enum MaterialStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected'
}

/**
 * 搜索请求参数
 */
export interface SearchRequest {
  keyword?: string
  category?: MaterialCategory
  course_name?: string
  tags?: string[]
  start_date?: string
  end_date?: string
  sort_by?: string
  sort_order?: string
  page?: number
  page_size?: number
}

/**
 * 搜索结果
 */
export interface SearchResult {
  id: number
  title: string
  description: string
  category: MaterialCategory
  course_name: string
  tags: string[]
  file_size: number
  download_count: number
  favorite_count: number
  view_count: number
  uploader_id: number
  uploader_name: string
  created_at: string
  updated_at: string
}

/**
 * 搜索响应
 */
export interface SearchResponse {
  results: SearchResult[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

/**
 * 热门搜索词
 */
export interface HotKeyword {
  id: number
  keyword: string
  search_count: number
  last_searched_at: string
}

/**
 * 搜索历史
 */
export interface SearchHistory {
  id: number
  keyword: string
  result_count: number
  created_at: string
}

/**
 * 推荐类型
 */
export enum RecommendationType {
  HOT = 'hot',
  PERSONALIZED = 'personalized',
  RELATED = 'related',
  DOWNLOADED = 'downloaded'
}

/**
 * 推荐请求参数
 */
export interface RecommendationRequest {
  type?: RecommendationType
  material_id?: number
  limit?: number
}

/**
 * 推荐结果
 */
export interface RecommendationResult {
  material: SearchResult
  reason: string
  score: number
}

/**
 * 搜索 API
 */
export const searchApi = {
  /**
   * 搜索资料
   */
  search: (params: SearchRequest) => {
    return request.get<ApiResponse<SearchResponse>>('/search', { params })
  },

  /**
   * 获取热门搜索词
   */
  getHotKeywords: (limit: number = 20) => {
    return request.get<ApiResponse<HotKeyword[]>>('/search/hot-keywords', {
      params: { limit }
    })
  },

  /**
   * 获取搜索历史
   */
  getSearchHistory: (limit: number = 20) => {
    return request.get<ApiResponse<SearchHistory[]>>('/search/history', {
      params: { limit }
    })
  },

  /**
   * 清空搜索历史
   */
  clearSearchHistory: () => {
    return request.delete<ApiResponse<void>>('/search/history')
  },

  /**
   * 获取热门资料
   */
  getHotMaterials: (limit: number = 20) => {
    return request.get<ApiResponse<SearchResult[]>>('/materials/hot', {
      params: { limit }
    })
  },

  /**
   * 获取推荐资料
   */
  getRecommendations: (params: RecommendationRequest) => {
    return request.get<ApiResponse<RecommendationResult[]>>('/materials/recommend', {
      params
    })
  }
}
