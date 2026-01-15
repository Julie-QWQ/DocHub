import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { materialApi } from '@/api/material'
import type { UploadSignatureRequest } from '@/types'

/**
 * 根据文件扩展名获取 MIME 类型
 * 注意: 必须与后端允许的类型保持一致
 */
const getMimeTypeByExtension = (ext: string): string => {
  // 与后端 backend/internal/pkg/oss/minio.go 中的 mimeMap 保持一致
  const mimeTypes: Record<string, string> = {
    '.pdf':  'application/pdf',
    '.doc':  'application/msword',
    '.docx': 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    '.xls':  'application/vnd.ms-excel',
    '.xlsx': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    '.ppt':  'application/vnd.ms-powerpoint',
    '.pptx': 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    '.txt':  'text/plain',
    '.md':   'text/markdown',
    '.csv':  'text/csv',
    '.zip':  'application/zip',
    '.rar':  'application/x-rar-compressed',
    '.7z':   'application/x-7z-compressed',
    '.tar':  'application/x-tar',
    '.jpg':  'image/jpeg',
    '.jpeg': 'image/jpeg',
    '.png':  'image/png',
    '.gif':  'image/gif',
    '.bmp':  'image/bmp',
    '.webp': 'image/webp',
  }

  const mimeType = mimeTypes[ext.toLowerCase()]

  // 如果找不到对应的 MIME 类型,返回空字符串
  if (!mimeType) {
    console.warn(`不支持的文件类型: ${ext},后端仅支持: pdf, docx, doc, pptx, ppt, txt, md, zip, rar 等格式`)
    return ''
  }

  return mimeType
}

/**
 * 文件上传相关的组合式函数
 */
export function useUpload() {
  const uploading = ref(false)
  const progress = ref(0)

  /**
   * 格式化文件大小
   */
  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  /**
   * 验证文件类型
   */
  const validateFileType = (file: File, allowedTypes: string[] = []): boolean => {
    if (allowedTypes.length === 0) {
      // 默认允许的类型
      allowedTypes = [
        'application/pdf',
        'application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        'application/vnd.ms-excel',
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        'application/vnd.ms-powerpoint',
        'application/vnd.openxmlformats-officedocument.presentationml.presentation',
        'application/zip',
        'application/x-rar-compressed',
        'application/x-7z-compressed',
        'image/jpeg',
        'image/png',
        'image/gif'
      ]
    }

    return allowedTypes.includes(file.type)
  }

  /**
   * 验证文件大小
   */
  const validateFileSize = (file: File, maxSize: number = 536870912): boolean => {
    // 默认最大 512MB
    return file.size <= maxSize
  }

  /**
   * 获取上传签名
   */
  const getUploadSignature = async (file: File): Promise<{ uploadUrl: string; fileKey: string } | null> => {
    try {
      // 如果 file.type 为空,根据文件扩展名推断 mime_type
      let mimeType = file.type
      if (!mimeType && file.name) {
        const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
        mimeType = getMimeTypeByExtension(ext)
      }

      // 如果 MIME 类型仍然为空,说明文件类型不支持
      if (!mimeType) {
        const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
        ElMessage.error(`不支持的文件类型: ${ext || '未知'}。仅支持: PDF, DOC, DOCX, PPT, PPTX, TXT, MD, ZIP, RAR 等格式`)
        return null
      }

      const data: UploadSignatureRequest = {
        file_name: file.name,
        file_size: file.size,
        mime_type: mimeType
      }

      const response = await materialApi.getUploadSignature(data)

      if (response.code === 0 && response.data) {
        return {
          uploadUrl: response.data.upload_url,
          fileKey: response.data.file_key
        }
      } else {
        ElMessage.error(response.message || '获取上传签名失败')
        return null
      }
    } catch (error: any) {
      ElMessage.error(error.message || '获取上传签名失败')
      return null
    }
  }

  /**
   * 上传文件到 OSS
   * @returns 成功返回 fileKey,失败返回 null
   */
  const uploadFile = async (file: File, onProgress?: (progress: number) => void): Promise<string | null> => {
    // 1. 获取上传签名
    const signature = await getUploadSignature(file)
    if (!signature) {
      return null
    }

    uploading.value = true
    progress.value = 0

    try {
      // 2. 使用 XMLHttpRequest 上传文件（支持进度监控）
      await new Promise<void>((resolve, reject) => {
        const xhr = new XMLHttpRequest()

        // 监听上传进度
        xhr.upload.addEventListener('progress', (e) => {
          if (e.lengthComputable) {
            progress.value = Math.round((e.loaded / e.total) * 100)
            onProgress?.(progress.value)
          }
        })

        // 监听上传完成
        xhr.addEventListener('load', () => {
          if (xhr.status === 200) {
            resolve()
          } else {
            reject(new Error('上传失败'))
          }
        })

        // 监听上传错误
        xhr.addEventListener('error', () => {
          reject(new Error('网络错误'))
        })

        // 发送 PUT 请求
        xhr.open('PUT', signature.uploadUrl)
        xhr.setRequestHeader('Content-Type', file.type)
        xhr.send(file)
      })

      // 上传成功,返回 fileKey
      return signature.fileKey
    } catch (error: any) {
      ElMessage.error(error.message || '文件上传失败')
      return null
    } finally {
      uploading.value = false
      progress.value = 0
    }
  }

  /**
   * 批量上传文件
   */
  const uploadFiles = async (
    files: File[],
    onProgress?: (index: number, progress: number) => void
  ): Promise<(boolean | string)[]> => {
    const results: (boolean | string)[] = []

    for (let i = 0; i < files.length; i++) {
      const success = await uploadFile(files[i], (progress) => {
        onProgress?.(i, progress)
      })

      results.push(success ? files[i].name : false)
    }

    return results
  }

  /**
   * 从文件对象读取文件信息
   */
  const getFileInfo = (file: File) => {
    return {
      name: file.name,
      size: file.size,
      type: file.type,
      lastModified: new Date(file.lastModified)
    }
  }

  /**
   * 检查是否是图片文件
   */
  const isImage = (file: File): boolean => {
    return file.type.startsWith('image/')
  }

  /**
   * 检查是否是 PDF 文件
   */
  const isPDF = (file: File): boolean => {
    return file.type === 'application/pdf'
  }

  /**
   * 生成文件预览 URL
   */
  const createPreviewUrl = (file: File): string | null => {
    if (isImage(file)) {
      return URL.createObjectURL(file)
    }
    return null
  }

  /**
   * 释放预览 URL
   */
  const revokePreviewUrl = (url: string) => {
    URL.revokeObjectURL(url)
  }

  return {
    // 状态
    uploading,
    progress,

    // 工具方法
    formatFileSize,
    validateFileType,
    validateFileSize,
    getFileInfo,
    isImage,
    isPDF,
    createPreviewUrl,
    revokePreviewUrl,

    // 上传方法
    getUploadSignature,
    uploadFile,
    uploadFiles
  }
}
