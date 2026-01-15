import { request } from '@/utils/request'
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  ChangePasswordRequest,
  RefreshTokenRequest,
  UserInfo
} from '@/types'

/**
 * 认证相关 API
 */
export const authApi = {
  /**
   * 用户登录
   */
  login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return request.post('/auth/login', data)
  },

  /**
   * 用户注册
   */
  register(data: RegisterRequest): Promise<ApiResponse<UserInfo>> {
    return request.post('/auth/register', data)
  },

  /**
   * 用户登出
   */
  logout(): Promise<ApiResponse<null>> {
    return request.post('/auth/logout')
  },

  /**
   * 刷新 Token
   */
  refreshToken(data: RefreshTokenRequest): Promise<ApiResponse<LoginResponse>> {
    return request.post('/auth/refresh', data)
  },

  /**
   * 修改密码
   */
  changePassword(data: ChangePasswordRequest): Promise<ApiResponse<null>> {
    return request.post('/auth/change-password', data)
  },

  /**
   * 获取当前用户信息
   */
  getUserInfo(): Promise<ApiResponse<UserInfo>> {
    return request.get('/auth/me')
  }
}
