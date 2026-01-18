import { ref } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

export interface SendCodeParams {
  email: string
  purpose: 'register' | 'login' | 'reset_password'
}

export interface VerifyCodeParams {
  email: string
  code: string
  purpose: 'register' | 'login' | 'reset_password'
}

export interface RegisterWithCodeParams {
  username: string
  email: string
  password: string
  code: string
}

export interface LoginWithCodeParams {
  email: string
  password?: string
  code: string
}

export function useEmailVerification() {
  const isSending = ref(false)
  const countdown = ref(0)
  const error = ref<string | null>(null)
  const authStore = useAuthStore()

  // 发送验证码
  const sendVerificationCode = async (params: SendCodeParams): Promise<boolean> => {
    if (isSending.value) return false

    // 清理邮箱地址(去除首尾空格)
    const cleanedEmail = params.email.trim()
    if (!cleanedEmail) {
      error.value = '请输入邮箱地址'
      return false
    }

    isSending.value = true
    error.value = null

    // 调试日志
    console.log('发送验证码请求:', {
      email: cleanedEmail,
      purpose: params.purpose,
      emailLength: cleanedEmail.length
    })

    try {
      const response = await axios.post('/api/v1/verification/send', {
        email: cleanedEmail,
        purpose: params.purpose
      })

      console.log('发送验证码响应:', response.data)

      // 开始倒计时(60秒)
      countdown.value = 60
      const timer = setInterval(() => {
        countdown.value--
        if (countdown.value <= 0) {
          clearInterval(timer)
        }
      }, 1000)

      return true
    } catch (err: any) {
      console.error('发送验证码错误:', err.response?.data)
      error.value = err.response?.data?.message || '发送验证码失败'
      return false
    } finally {
      isSending.value = false
    }
  }

  // 验证验证码
  const verifyCode = async (params: VerifyCodeParams): Promise<boolean> => {
    error.value = null

    try {
      await axios.post('/api/v1/verification/verify', params)
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || '验证码验证失败'
      return false
    }
  }

  // 使用验证码注册
  const registerWithCode = async (params: RegisterWithCodeParams): Promise<boolean> => {
    error.value = null

    try {
      await axios.post('/api/v1/auth/register', params)
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || '注册失败'
      return false
    }
  }

  // 使用验证码登录
  const loginWithCode = async (params: LoginWithCodeParams): Promise<boolean> => {
    error.value = null

    try {
      const response = await axios.post('/api/v1/auth/login', params)
      if (response?.data?.data) {
        authStore.applyLoginResponse(response.data.data)
      }
      return true
    } catch (err: any) {
      error.value = err.response?.data?.message || '登录失败'
      ElMessage.error(error.value)
      return false
    }
  }

  return {
    isSending,
    countdown,
    error,
    sendVerificationCode,
    verifyCode,
    registerWithCode,
    loginWithCode
  }
}
