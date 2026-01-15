import request from '@/utils/request'

/**
 * 系统配置响应接口
 */
export interface SystemConfigResponse {
  config_key: string
  config_value: string
  description: string
  category: string
  created_at: string
  updated_at: string
}

/**
 * 上传配置接口
 */
export interface UploadConfig {
  max_size: number
  allowed_types: string[]
  max_size_mb: number
  accept: string
}

/**
 * 获取系统配置（公共接口）
 */
export function getPublicSystemConfig(key: string) {
  return request<string>({
    url: `/system/configs/${key}`,
    method: 'get'
  })
}

/**
 * 批量获取系统配置（公共接口）
 */
export function getPublicSystemConfigs(keys: string[]) {
  return request<Record<string, string>>({
    url: '/system/configs',
    method: 'get',
    params: { keys: keys.join(',') }
  })
}

/**
 * 获取上传配置（公共接口）
 */
export function getUploadConfig() {
  return request<UploadConfig>({
    url: '/system/upload-config',
    method: 'get'
  })
}
