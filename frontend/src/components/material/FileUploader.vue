<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUpload } from '@/composables/useUpload'
import { useUploadConfig } from '@/composables/useUploadConfig'
import { materialApi } from '@/api/material'

/**
 * 根据文件扩展名获取 MIME 类型
 * 注意: 必须与后端 backend/internal/pkg/oss/minio.go 中的 mimeMap 保持一致
 */
const getMimeTypeByExtension = (fileName: string): string => {
  const ext = fileName.substring(fileName.lastIndexOf('.')).toLowerCase()

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

  return mimeTypes[ext] || ''
}

interface FileInfo {
  file: File
  fileKey: string
  progress: number
  status: 'pending' | 'uploading' | 'success' | 'error'
  error?: string
  mimeType?: string // 推断的 MIME 类型
}

const props = defineProps<{
  maxFileSize?: number // 字节，默认使用系统配置
  maxFiles?: number
  accept?: string // 文件类型，如 '.pdf,.doc,.docx'，默认使用系统配置
}>()

const emit = defineEmits<{
  (e: 'file-selected', files: FileInfo[]): void
  (e: 'upload-complete', fileInfos: FileInfo[]): void
}>()

const { uploading, progress, formatFileSize, uploadFile } = useUpload()
const { allowedExtensions, maxSize, accept, loadConfig, validateFileType } = useUploadConfig()

const fileList = ref<FileInfo[]>([])
const uploadingIndex = ref(-1)
const uploadRef = ref()

// 加载上传配置
onMounted(async () => {
  await loadConfig()
})

// 允许的文件类型（从系统配置获取）
const allowedTypes = computed(() => allowedExtensions.value)

// 最大文件大小（从系统配置获取，或使用 props 传入的值）
const maxAllowedFileSize = computed(() => props.maxFileSize || maxSize.value)

// accept 属性（从系统配置获取，或使用 props 传入的值）
const acceptAttr = computed(() => props.accept || accept.value)

// 本地验证文件大小
const validateFileSize = (file: File): boolean => {
  return file.size <= maxAllowedFileSize.value
}

// 处理文件选择
const handleFileChange = (uploadFile: any, uploadFiles: any[]) => {
  // 只处理当前选择的文件,不使用 uploadFiles(它会累积所有选择过的文件)
  const file = uploadFile.raw

  if (!file) {
    return
  }

  // 验证文件数量（检查现有文件列表数量 + 1）
  const currentFileCount = fileList.value.length

  if (props.maxFiles && currentFileCount >= props.maxFiles) {
    ElMessage.warning(`最多只能上传 ${props.maxFiles} 个文件`)
    // 清空 el-upload 组件内部的文件列表,防止累积
    uploadRef.value?.clearFiles()
    return
  }

  // 验证文件类型（基于扩展名）
  if (!validateFileType(file)) {
    ElMessage.error(`文件类型不支持: ${file.name}`)
    uploadRef.value?.clearFiles()
    return
  }

  // 验证文件大小
  if (!validateFileSize(file)) {
    ElMessage.error(`文件大小超出限制 (${formatFileSize(maxAllowedFileSize.value)}): ${file.name}`)
    uploadRef.value?.clearFiles()
    return
  }

  // 检查是否已存在
  const exists = fileList.value.some((f) => f.file.name === file.name && f.file.size === file.size)
  if (exists) {
    ElMessage.warning(`文件已存在: ${file.name}`)
    uploadRef.value?.clearFiles()
    return
  }

  // 推断 MIME 类型(优先使用浏览器提供的类型,如果为空则根据扩展名推断)
  let mimeType = file.type
  if (!mimeType) {
    mimeType = getMimeTypeByExtension(file.name)
  }

  const newFileInfo: FileInfo = {
    file,
    fileKey: '',
    progress: 0,
    status: 'pending',
    mimeType // 保存推断的 MIME 类型
  }

  fileList.value = [...fileList.value, newFileInfo]
  emit('file-selected', [newFileInfo])

  // 清空 el-upload 组件内部的文件列表,防止累积
  uploadRef.value?.clearFiles()
}

// 移除文件
const removeFile = async (index: number) => {
  const fileInfo = fileList.value[index]

  // 如果文件已经上传成功,需要删除OSS上的文件
  if (fileInfo.status === 'success' && fileInfo.fileKey) {
    try {
      await materialApi.deleteUploadedFile(fileInfo.fileKey)
    } catch (error: any) {
      console.error('删除OSS文件失败:', error)
      // 即使删除失败,也从列表中移除
    }
  }

  // 从列表中移除
  fileList.value.splice(index, 1)
}

// 上传单个文件
const uploadSingleFile = async (fileInfo: FileInfo): Promise<boolean> => {
  fileInfo.status = 'uploading'

  const fileKey = await uploadFile(fileInfo.file, (percent) => {
    fileInfo.progress = percent
  })

  if (fileKey) {
    fileInfo.status = 'success'
    fileInfo.fileKey = fileKey // 保存 fileKey
  } else {
    fileInfo.status = 'error'
    fileInfo.error = '上传失败'
  }

  return !!fileKey
}

// 上传所有文件
const uploadAllFiles = async (): Promise<FileInfo[]> => {
  const results: FileInfo[] = []

  for (let i = 0; i < fileList.value.length; i++) {
    const fileInfo = fileList.value[i]

    if (fileInfo.status === 'success') {
      results.push(fileInfo)
      continue
    }

    uploadingIndex.value = i
    const success = await uploadSingleFile(fileInfo)

    if (success) {
      results.push(fileInfo)
    }

    uploadingIndex.value = -1
  }

  emit('upload-complete', results)
  return results
}

// 清空文件列表
const clearFiles = () => {
  fileList.value = []
  // 清除 el-upload 组件内部的文件列表
  uploadRef.value?.clearFiles()
}

// 获取已上传的文件信息
const getUploadedFiles = computed(() => {
  return fileList.value.filter((f) => f.status === 'success')
})

// 是否有文件正在上传
const isUploading = computed(() => {
  return fileList.value.some((f) => f.status === 'uploading')
})

// 所有文件是否都已上传完成
const allUploaded = computed(() => {
  return fileList.value.length > 0 && fileList.value.every((f) => f.status === 'success')
})

// 是否有文件(不管是否上传)
const hasFiles = computed(() => {
  return fileList.value.length > 0
})

// 暴露方法给父组件
defineExpose({
  uploadAllFiles,
  clearFiles,
  getUploadedFiles,
  isUploading,
  allUploaded,
  hasFiles
})
</script>

<template>
  <div class="file-uploader">
    <el-upload
      ref="uploadRef"
      :auto-upload="false"
      :show-file-list="false"
      :on-change="handleFileChange"
      :accept="acceptAttr"
      multiple
      drag
    >
      <div class="upload-area">
        <el-icon class="upload-icon"><UploadFilled /></el-icon>
        <div class="upload-text">
          <p>拖拽文件到此处，或<em>点击上传</em></p>
          <p class="upload-tip">
            支持 {{ acceptAttr }} 等格式，单个文件最大 {{ formatFileSize(maxAllowedFileSize) }}
          </p>
        </div>
      </div>
    </el-upload>

    <!-- 文件列表 -->
    <div v-if="fileList.length > 0" class="file-list">
      <div v-for="(fileInfo, index) in fileList" :key="index" class="file-item">
        <div class="file-info">
          <el-icon class="file-icon"><Document /></el-icon>
          <div class="file-details">
            <div class="file-name">{{ fileInfo.file.name }}</div>
            <div class="file-meta">
              <span>{{ formatFileSize(fileInfo.file.size) }}</span>
              <span v-if="fileInfo.status === 'uploading'" class="upload-progress">
                {{ fileInfo.progress }}%
              </span>
              <span v-else-if="fileInfo.status === 'success'" class="upload-success">上传成功</span>
              <span v-else-if="fileInfo.status === 'error'" class="upload-error">
                {{ fileInfo.error || '上传失败' }}
              </span>
            </div>
          </div>
        </div>

        <!-- 进度条 -->
        <div v-if="fileInfo.status === 'uploading'" class="file-progress">
          <el-progress :percentage="fileInfo.progress" :show-text="false" />
        </div>

        <!-- 删除按钮 -->
        <el-button
          :icon="'Close'"
          circle
          size="small"
          :disabled="fileInfo.status === 'uploading'"
          @click="removeFile(index)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.file-uploader {
  .upload-area {
    padding: 40px 0;
    text-align: center;

    .upload-icon {
      font-size: 48px;
      color: #1a1a1a;
      margin-bottom: 16px;
    }

    .upload-text {
      p {
        margin: 8px 0;
        color: #1a1a1a;

        em {
          color: #FF6B35;
          font-style: normal;
          font-weight: 500;
        }
      }

      .upload-tip {
        font-size: 13px;
        color: #666;
      }
    }
  }

  .file-list {
    margin-top: 16px;
  }

  .file-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: #fafafa;
    border: 1px solid #e5e5e5;
    border-radius: 8px;
    margin-bottom: 8px;
    transition: all 0.2s;

    &:hover {
      background: #ffffff;
      border-color: #1a1a1a;
    }

    .file-info {
      flex: 1;
      display: flex;
      align-items: center;
      gap: 12px;
      min-width: 0;

      .file-icon {
        font-size: 24px;
        color: #666;
        flex-shrink: 0;
      }

      .file-details {
        flex: 1;
        min-width: 0;

        .file-name {
          font-size: 14px;
          color: #1a1a1a;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          font-weight: 500;
        }

        .file-meta {
          display: flex;
          align-items: center;
          gap: 12px;
          margin-top: 4px;
          font-size: 12px;
          color: #666;

          .upload-progress {
            color: #FF6B35;
            font-weight: 500;
          }

          .upload-success {
            color: #22c55e;
            font-weight: 500;
          }

          .upload-error {
            color: #dc2626;
            font-weight: 500;
          }
        }
      }
    }

    .file-progress {
      flex: 1;
      max-width: 200px;
    }
  }
}

:deep(.el-upload-dragger) {
  background: #fafafa;
  border: 2px dashed #e5e5e5;
  border-radius: 8px;
  transition: all 0.2s;

  &:hover {
    background: #ffffff;
    border-color: #1a1a1a;
  }
}

:deep(.el-progress-bar__outer) {
  background-color: #e5e5e5;
}

:deep(.el-progress-bar__inner) {
  background-color: #FF6B35;
}
</style>
