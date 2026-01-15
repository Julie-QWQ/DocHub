import request from '@/utils/request'

// 统计相关类型定义

export interface OverviewStatistics {
  users: UserStatistics
  materials: MaterialStatistics
  downloads: DownloadStatistics
  applications: ApplicationStatistics
  visits: VisitStatistics
}

export interface UserStatistics {
  total: number
  today: number
  week: number
  month: number
  active: number
  by_role: Record<string, number>
}

export interface MaterialStatistics {
  total: number
  approved: number
  pending: number
  rejected: number
  offline: number
  today: number
  week: number
  by_category: Record<string, number>
}

export interface DownloadStatistics {
  total: number
  today: number
  week: number
  month: number
}

export interface ApplicationStatistics {
  total: number
  pending: number
  approved: number
  rejected: number
}

export interface VisitStatistics {
  total: number
  today: number
  week: number
  month: number
  unique: number
}

export interface TrendData {
  date: string
  count: number
  value: number
}

/**
 * 获取概览统计
 */
export function getOverviewStatistics() {
  return request<OverviewStatistics>({
    url: '/admin/statistics/overview',
    method: 'get'
  })
}

/**
 * 获取用户统计
 */
export function getUserStatistics() {
  return request<UserStatistics>({
    url: '/admin/statistics/users',
    method: 'get'
  })
}

/**
 * 获取用户趋势
 */
export function getUserTrend(days: number = 30) {
  return request<TrendData[]>({
    url: '/admin/statistics/users/trend',
    method: 'get',
    params: { days }
  })
}

/**
 * 获取资料统计
 */
export function getMaterialStatistics() {
  return request<MaterialStatistics>({
    url: '/admin/statistics/materials',
    method: 'get'
  })
}

/**
 * 获取资料趋势
 */
export function getMaterialTrend(days: number = 30) {
  return request<TrendData[]>({
    url: '/admin/statistics/materials/trend',
    method: 'get',
    params: { days }
  })
}

/**
 * 获取下载统计
 */
export function getDownloadStatistics() {
  return request<DownloadStatistics>({
    url: '/admin/statistics/downloads',
    method: 'get'
  })
}

/**
 * 获取下载趋势
 */
export function getDownloadTrend(days: number = 30) {
  return request<TrendData[]>({
    url: '/admin/statistics/downloads/trend',
    method: 'get',
    params: { days }
  })
}

/**
 * 获取申请统计
 */
export function getApplicationStatistics() {
  return request<ApplicationStatistics>({
    url: '/admin/statistics/applications',
    method: 'get'
  })
}

/**
 * 获取访问统计
 */
export function getVisitStatistics() {
  return request<VisitStatistics>({
    url: '/admin/statistics/visits',
    method: 'get'
  })
}

/**
 * 获取访问趋势
 */
export function getVisitTrend(days: number = 30) {
  return request<TrendData[]>({
    url: '/admin/statistics/visits/trend',
    method: 'get',
    params: { days }
  })
}

/**
 * 记录页面浏览
 */
export function recordPageView(path: string, referer?: string) {
  return request({
    url: '/statistics/page-view',
    method: 'post',
    data: { path, referer }
  })
}
