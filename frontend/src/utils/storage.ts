import { STORAGE_KEYS } from './constants'

/**
 * 本地存储工具
 */
export const storage = {
  /**
   * 设置 Token
   */
  setToken(token: string): void {
    localStorage.setItem(STORAGE_KEYS.TOKEN, token)
  },

  /**
   * 获取 Token
   */
  getToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.TOKEN)
  },

  /**
   * 设置刷新 Token
   */
  setRefreshToken(token: string): void {
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, token)
  },

  /**
   * 获取刷新 Token
   */
  getRefreshToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
  },

  /**
   * 设置 Token 过期时间
   */
  setTokenExpireTime(timestamp: number): void {
    localStorage.setItem(STORAGE_KEYS.TOKEN_EXPIRE_TIME, String(timestamp))
  },

  /**
   * 获取 Token 过期时间
   */
  getTokenExpireTime(): number | null {
    const time = localStorage.getItem(STORAGE_KEYS.TOKEN_EXPIRE_TIME)
    return time ? Number(time) : null
  },

  /**
   * 移除 Token
   */
  removeToken(): void {
    localStorage.removeItem(STORAGE_KEYS.TOKEN)
  },

  /**
   * 移除刷新 Token
   */
  removeRefreshToken(): void {
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)
  },

  /**
   * 移除 Token 过期时间
   */
  removeTokenExpireTime(): void {
    localStorage.removeItem(STORAGE_KEYS.TOKEN_EXPIRE_TIME)
  },

  /**
   * 设置用户信息
   */
  setUser(userInfo: any): void {
    localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(userInfo))
  },

  /**
   * 获取用户信息
   */
  getUser(): any | null {
    const info = localStorage.getItem(STORAGE_KEYS.USER_INFO)
    return info ? JSON.parse(info) : null
  },

  /**
   * 移除用户信息
   */
  removeUser(): void {
    localStorage.removeItem(STORAGE_KEYS.USER_INFO)
  },

  /**
   * 清除所有认证信息
   */
  clearAuth(): void {
    this.removeToken()
    this.removeRefreshToken()
    this.removeTokenExpireTime()
    this.removeUser()
  },

  /**
   * 设置主题
   */
  setTheme(theme: string): void {
    localStorage.setItem(STORAGE_KEYS.THEME, theme)
  },

  /**
   * 获取主题
   */
  getTheme(): string | null {
    return localStorage.getItem(STORAGE_KEYS.THEME)
  },

  /**
   * 设置会话存储
   */
  setSession(key: string, value: any): void {
    sessionStorage.setItem(key, JSON.stringify(value))
  },

  /**
   * 获取会话存储
   */
  getSession<T = any>(key: string): T | null {
    const value = sessionStorage.getItem(key)
    return value ? JSON.parse(value) : null
  },

  /**
   * 移除会话存储
   */
  removeSession(key: string): void {
    sessionStorage.removeItem(key)
  }
}
