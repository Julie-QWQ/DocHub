import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

/**
 * 格式化日期为相对时间
 */
export const formatDate = (date: string | Date): string => {
  const now = dayjs()
  const target = dayjs(date)
  const diff = now.diff(target, 'day')

  if (diff === 0) return '今天'
  if (diff === 1) return '昨天'
  if (diff < 7) return `${diff} 天前`
  if (diff < 30) return `${Math.floor(diff / 7)} 周前`
  if (diff < 365) return `${Math.floor(diff / 30)} 月前`
  return `${Math.floor(diff / 365)} 年前`
}

/**
 * 格式化工具函数
 */
export const format = {
  /**
   * 格式化日期时间
   */
  datetime(date: string | Date, format = 'YYYY-MM-DD HH:mm:ss'): string {
    return dayjs(date).format(format)
  },

  /**
   * 格式化日期
   */
  date(date: string | Date): string {
    return dayjs(date).format('YYYY-MM-DD')
  },

  /**
   * 格式化时间
   */
  time(date: string | Date): string {
    return dayjs(date).format('HH:mm:ss')
  },

  /**
   * 格式化相对时间
   */
  relative(date: string | Date): string {
    return dayjs(date).fromNow()
  },

  /**
   * 格式化文件大小
   */
  fileSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  },

  /**
   * 格式化数字（千分位）
   */
  number(num: number): string {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  },

  /**
   * 格式化百分比
   */
  percent(value: number, total: number, decimals = 2): string {
    if (total === 0) return '0%'
    return ((value / total) * 100).toFixed(decimals) + '%'
  },

  /**
   * 截断文本
   */
  truncate(text: string, length = 50, suffix = '...'): string {
    if (text.length <= length) return text
    return text.substring(0, length) + suffix
  }
}
