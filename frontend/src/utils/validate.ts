/**
 * 验证工具函数
 */
export const validate = {
  /**
   * 验证邮箱
   */
  email(email: string): boolean {
    const reg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
    return reg.test(email)
  },

  /**
   * 验证手机号（中国大陆）
   */
  phone(phone: string): boolean {
    const reg = /^1[3-9]\d{9}$/
    return reg.test(phone)
  },

  /**
   * 验证用户名（字母开头，允许字母数字下划线，4-20位）
   */
  username(username: string): boolean {
    const reg = /^[a-zA-Z][a-zA-Z0-9_]{3,19}$/
    return reg.test(username)
  },

  /**
   * 验证密码（至少包含字母和数字，8-20位）
   */
  password(password: string): boolean {
    const reg = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]{8,20}$/
    return reg.test(password)
  },

  /**
   * 验证 URL
   */
  url(url: string): boolean {
    try {
      new URL(url)
      return true
    } catch {
      return false
    }
  },

  /**
   * 验证文件类型
   */
  fileType(fileName: string, allowedTypes: string[]): boolean {
    const ext = fileName.substring(fileName.lastIndexOf('.')).toLowerCase()
    return allowedTypes.some(type => {
      const typeExt = type.substring(type.lastIndexOf('/')).toLowerCase()
      return ext === typeExt || type.includes(ext)
    })
  },

  /**
   * 验证文件大小
   */
  fileSize(size: number, maxSize: number): boolean {
    return size <= maxSize
  },

  /**
   * 验证是否为空
   */
  empty(value: any): boolean {
    if (value === null || value === undefined) return true
    if (typeof value === 'string') return value.trim().length === 0
    if (Array.isArray(value)) return value.length === 0
    if (typeof value === 'object') return Object.keys(value).length === 0
    return false
  }
}
