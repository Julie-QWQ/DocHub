import { ref, computed, onMounted } from 'vue'
import { getUploadConfig } from '@/api/system'
import type { UploadConfig } from '@/api/system'

// 全局上传配置状态
const uploadConfig = ref<UploadConfig>({
  max_size: 52428800, // 默认 50MB
  allowed_types: ['pdf', 'docx', 'doc', 'pptx', 'ppt', 'zip', 'rar'],
  max_size_mb: 50,
  accept: '.pdf,.docx,.doc,.pptx,.ppt,.zip,.rar'
})

const loading = ref(false)
const loaded = ref(false)

/**
 * 上传配置管理
 */
export function useUploadConfig() {
  /**
   * 加载上传配置
   */
  const loadConfig = async () => {
    if (loaded.value) return

    loading.value = true
    try {
      const response = await getUploadConfig()
      if (response.code === 0 && response.data) {
        uploadConfig.value = response.data
        loaded.value = true
      }
    } catch (error) {
      console.error('加载上传配置失败:', error)
      // 使用默认值
    } finally {
      loading.value = false
    }
  }

  /**
   * 刷新配置
   */
  const refreshConfig = async () => {
    loaded.value = false
    await loadConfig()
  }

  /**
   * 获取允许的文件扩展名列表
   */
  const allowedExtensions = computed(() => {
    return uploadConfig.value.allowed_types.map(type => type.toLowerCase())
  })

  /**
   * 获取允许的 MIME 类型列表
   */
  const allowedMimeTypes = computed(() => {
    const mimeMap: Record<string, string> = {
      pdf: 'application/pdf',
      doc: 'application/msword',
      docx: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      xls: 'application/vnd.ms-excel',
      xlsx: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      ppt: 'application/vnd.ms-powerpoint',
      pptx: 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
      txt: 'text/plain',
      md: 'text/markdown',
      csv: 'text/csv',
      zip: 'application/zip',
      rar: 'application/x-rar-compressed',
      '7z': 'application/x-7z-compressed',
      tar: 'application/x-tar',
      jpg: 'image/jpeg',
      jpeg: 'image/jpeg',
      png: 'image/png',
      gif: 'image/gif',
      bmp: 'image/bmp',
      webp: 'image/webp'
    }

    return uploadConfig.value.allowed_types
      .map(type => mimeMap[type])
      .filter(Boolean)
  })

  /**
   * 验证文件类型（基于扩展名）
   */
  const validateFileType = (file: File): boolean => {
    // 获取文件扩展名
    const fileName = file.name.toLowerCase()
    const extIndex = fileName.lastIndexOf('.')

    if (extIndex === -1) {
      return false // 没有扩展名
    }

    const extension = fileName.slice(extIndex + 1)
    return allowedExtensions.value.includes(extension)
  }

  /**
   * 验证文件大小
   */
  const validateFileSize = (file: File): boolean => {
    return file.size <= uploadConfig.value.max_size
  }

  return {
    // 状态
    uploadConfig: computed(() => uploadConfig.value),
    loading: computed(() => loading.value),
    loaded: computed(() => loaded.value),

    // 计算属性
    allowedExtensions,
    allowedMimeTypes,
    maxSize: computed(() => uploadConfig.value.max_size),
    maxSizeMB: computed(() => uploadConfig.value.max_size_mb),
    accept: computed(() => uploadConfig.value.accept),

    // 方法
    loadConfig,
    refreshConfig,
    validateFileType,
    validateFileSize
  }
}

/**
 * 初始化上传配置（在应用启动时调用）
 */
export async function initUploadConfig() {
  const { loadConfig } = useUploadConfig()
  await loadConfig()
}
