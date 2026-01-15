import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest, RegisterRequest, ChangePasswordRequest } from '@/types'

/**
 * 认证相关的组合式函数
 */
export function useAuth() {
  const authStore = useAuthStore()
  const router = useRouter()
  const route = useRoute()

  /**
   * 登录
   */
  const login = async (credentials: LoginRequest) => {
    const success = await authStore.login(credentials)
    if (success) {
      // 登录成功后跳转
      const redirect = route.query.redirect as string
      router.push(redirect || '/materials')
    }
    return success
  }

  /**
   * 注册
   */
  const register = async (data: RegisterRequest) => {
    const success = await authStore.register(data)
    if (success) {
      // 注册成功后跳转到登录页
      router.push('/login')
    }
    return success
  }

  /**
   * 登出
   */
  const logout = async () => {
    // 立即清除本地认证状态
    authStore.clearAuth()

    // 显示成功消息
    const { ElMessage } = await import('element-plus')
    ElMessage.success('已退出登录')

    // 异步调用 logout API(不等待,避免阻塞跳转)
    authApi.logout().catch(error => {
      console.error('Logout API failed:', error)
    })

    // 立即跳转到登录页
    await router.push('/login')
  }

  /**
   * 修改密码
   */
  const changePassword = async (data: ChangePasswordRequest) => {
    return await authStore.changePassword(data)
  }

  /**
   * 刷新用户信息
   */
  const refreshUserInfo = async () => {
    return await authStore.fetchUserInfo()
  }

  return {
    // 状态
    user: computed(() => authStore.user),
    isLoggedIn: computed(() => authStore.isLoggedIn),
    userRole: computed(() => authStore.userRole),
    userName: computed(() => authStore.userName),
    userAvatar: computed(() => authStore.userAvatar),
    isAdmin: computed(() => authStore.isAdmin),
    isCommittee: computed(() => authStore.isCommittee),

    // 方法
    login,
    register,
    logout,
    changePassword,
    refreshUserInfo,
    hasRole: authStore.hasRole
  }
}
