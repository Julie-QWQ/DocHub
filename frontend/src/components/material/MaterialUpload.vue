<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMaterialStore } from '@/stores/material'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import FileUploader from './FileUploader.vue'
import type { MaterialCategory } from '@/types'

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'cancel'): void
}>()

const materialStore = useMaterialStore()
const materialCategoryStore = useMaterialCategoryStore()

const fileUploaderRef = ref()
const submitting = ref(false)

// 表单数据
const formData = reactive({
  title: '',
  description: '',
  category: '' as MaterialCategory,
  course_name: '',
  tags: [] as string[]
})

// 输入的标签
const tagInput = ref('')

// 表单错误提示
const errors = ref<Record<string, string>>({})

// 资料分类选项（从 store 动态获取）
const categoryOptions = computed(() => {
  return materialCategoryStore.activeCategories.map(cat => ({
    label: cat.name,
    value: cat.code as MaterialCategory
  }))
})

// 加载资料类型
onMounted(async () => {
  try {
    await materialCategoryStore.fetchActiveCategories()
  } catch (error: any) {
    ElMessage.error('加载资料类型失败: ' + error.message)
  }
})

// 验证表单
const validateForm = (): boolean => {
  errors.value = {}

  if (!formData.title.trim()) {
    errors.value.title = '请输入资料标题'
  }

  if (!formData.category) {
    errors.value.category = '请选择资料分类'
  }

  if (!formData.course_name.trim()) {
    errors.value.course_name = '请输入课程名称'
  }

  // 检查文件
  if (!fileUploaderRef.value) {
    errors.value.file = '请选择文件'
  } else {
    if (!fileUploaderRef.value.hasFiles) {
      errors.value.file = '请选择文件'
    }
  }

  return Object.keys(errors.value).length === 0
}

// 添加标签
const handleAddTag = () => {
  const tag = tagInput.value.trim()
  if (!tag) return

  if (formData.tags.includes(tag)) {
    ElMessage.warning('标签已存在')
    return
  }

  if (formData.tags.length >= 10) {
    ElMessage.warning('最多添加 10 个标签')
    return
  }

  formData.tags.push(tag)
  tagInput.value = ''
}

// 删除标签
const handleRemoveTag = (tag: string) => {
  const index = formData.tags.indexOf(tag)
  if (index !== -1) {
    formData.tags.splice(index, 1)
  }
}

// 提交表单
const handleSubmit = async () => {
  // 清除之前的错误
  errors.value = {}

  // 验证表单
  if (!validateForm()) {
    return
  }

  submitting.value = true

  try {
    // 上传文件（如果还未上传）
    if (!fileUploaderRef.value.allUploaded) {
      await fileUploaderRef.value.uploadAllFiles()
    }

    // 获取已上传的文件
    const files = fileUploaderRef.value.getUploadedFiles
    if (!files || files.length === 0) {
      ElMessage.error('文件上传失败')
      return
    }

    const file = files[0]

    // 创建资料(使用推断的 MIME 类型,确保 .md 等文件类型能正确识别)
    const mimeType = file.mimeType || file.file.type
    if (!mimeType) {
      ElMessage.error('无法识别文件类型')
      return
    }

    await materialStore.createMaterial({
      title: formData.title,
      description: formData.description,
      category: formData.category,
      course_name: formData.course_name,
      file_name: file.file.name,
      file_size: file.file.size,
      file_key: file.fileKey,
      mime_type: mimeType,
      tags: formData.tags
    })

    ElMessage.success('资料上传成功')
    emit('success')
  } catch (error: any) {
    ElMessage.error(error.message || '上传失败')
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  emit('cancel')
}

// 重置表单
const resetForm = () => {
  formData.title = ''
  formData.description = ''
  formData.category = '' as MaterialCategory
  formData.course_name = ''
  formData.tags = []
  tagInput.value = ''
  errors.value = {}
  fileUploaderRef.value?.clearFiles()
}
</script>

<template>
  <div class="upload-form">
    <form @submit.prevent="handleSubmit">
      <!-- 文件上传 -->
      <div class="form-group">
        <label class="form-label">
          上传文件 <span class="required">*</span>
        </label>
        <FileUploader ref="fileUploaderRef" :max-files="1" />
        <span v-if="errors.file" class="error-message">{{ errors.file }}</span>
      </div>

      <!-- 资料标题 -->
      <div class="form-group">
        <label class="form-label">
          资料标题 <span class="required">*</span>
        </label>
        <input
          v-model="formData.title"
          type="text"
          class="form-input"
          placeholder="请输入资料标题"
          maxlength="200"
          :class="{ error: errors.title }"
        />
        <span v-if="errors.title" class="error-message">{{ errors.title }}</span>
      </div>

      <!-- 资料分类和课程名称 -->
      <div class="form-row">
        <div class="form-group half">
          <label class="form-label">
            资料分类 <span class="required">*</span>
          </label>
          <select
            v-model="formData.category"
            class="form-select"
            :class="{ error: errors.category }"
          >
            <option value="">请选择分类</option>
            <option
              v-for="option in categoryOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </option>
          </select>
          <span v-if="errors.category" class="error-message">{{ errors.category }}</span>
        </div>

        <div class="form-group half">
          <label class="form-label">
            课程名称 <span class="required">*</span>
          </label>
          <input
            v-model="formData.course_name"
            type="text"
            class="form-input"
            placeholder="请输入课程名称"
            maxlength="100"
            :class="{ error: errors.course_name }"
          />
          <span v-if="errors.course_name" class="error-message">{{ errors.course_name }}</span>
        </div>
      </div>

      <!-- 资料描述 -->
      <div class="form-group">
        <label class="form-label">资料描述</label>
        <textarea
          v-model="formData.description"
          class="form-textarea"
          placeholder="请输入资料描述（选填）"
          rows="4"
          maxlength="1000"
        ></textarea>
        <span class="char-count">{{ formData.description.length }}/1000</span>
      </div>

      <!-- 标签 -->
      <div class="form-group">
        <label class="form-label">标签</label>
        <div class="tags-input">
          <span
            v-for="tag in formData.tags"
            :key="tag"
            class="tag-item"
          >
            {{ tag }}
            <button type="button" class="tag-close" @click="handleRemoveTag(tag)">
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </span>
          <input
            v-if="formData.tags.length < 10"
            v-model="tagInput"
            type="text"
            class="tag-input-field"
            placeholder="输入标签后按回车添加"
            @keyup.enter="handleAddTag"
          />
          <span v-else class="tags-limit">最多添加 10 个标签</span>
        </div>
        <span class="help-text">按回车键添加标签,最多添加 10 个</span>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button type="submit" class="btn btn-primary" :disabled="submitting">
          <span v-if="submitting">提交中...</span>
          <span v-else>提交审核</span>
        </button>
        <button type="button" class="btn btn-secondary" @click="resetForm">
          重置
        </button>
        <button type="button" class="btn btn-text" @click="handleCancel">
          取消
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped lang="scss">
.upload-form {
  background: #ffffff;
}

.form-group {
  margin-bottom: 24px;

  &.half {
    flex: 1;
    min-width: 0;
  }

  .form-label {
    display: block;
    font-size: 14px;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 8px;

    .required {
      color: #dc2626;
      margin-left: 2px;
    }
  }

  .form-input,
  .form-select,
  .form-textarea {
    width: 100%;
    padding: 12px 16px;
    background: #ffffff;
    border: 1px solid #e5e5e5;
    border-radius: 8px;
    font-size: 15px;
    color: #1a1a1a;
    outline: none;
    transition: all 0.2s;
    font-family: inherit;

    &::placeholder {
      color: #999;
    }

    &:focus {
      border-color: #1a1a1a;
      box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.1);
    }

    &.error {
      border-color: #dc2626;

      &:focus {
        box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.1);
      }
    }
  }

  .form-textarea {
    resize: vertical;
    min-height: 100px;
  }

  .error-message {
    display: block;
    font-size: 13px;
    color: #dc2626;
    margin-top: 6px;
  }

  .char-count {
    display: block;
    text-align: right;
    font-size: 12px;
    color: #999;
    margin-top: 6px;
  }

  .help-text {
    display: block;
    font-size: 13px;
    color: #999;
    margin-top: 6px;
  }

  .tags-input {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    padding: 12px;
    background: #ffffff;
    border: 1px solid #e5e5e5;
    border-radius: 8px;
    min-height: 48px;
    transition: all 0.2s;

    &:focus-within {
      border-color: #1a1a1a;
      box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.1);
    }

    .tag-item {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 6px 12px;
      background: #1a1a1a;
      color: #ffffff;
      border-radius: 6px;
      font-size: 13px;
      font-weight: 500;

      .tag-close {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 16px;
        height: 16px;
        background: rgba(255, 255, 255, 0.2);
        border: none;
        border-radius: 3px;
        cursor: pointer;
        transition: background 0.15s;

        &:hover {
          background: rgba(255, 255, 255, 0.3);
        }

        svg {
          width: 10px;
          height: 10px;
        }
      }
    }

    .tag-input-field {
      flex: 1;
      min-width: 150px;
      border: none;
      background: transparent;
      font-size: 14px;
      color: #1a1a1a;
      outline: none;

      &::placeholder {
        color: #999;
      }
    }

    .tags-limit {
      font-size: 13px;
      color: #999;
    }
  }
}

.form-row {
  display: flex;
  gap: 16px;

  .form-group {
    margin-bottom: 0;
  }
}

.form-actions {
  display: flex;
  gap: 12px;
  padding-top: 8px;
  border-top: 1px solid #f2f2f2;
  margin-top: 32px;

  .btn {
    padding: 12px 24px;
    font-size: 15px;
    font-weight: 500;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    border: none;

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    &-primary {
      background: #1a1a1a;
      color: #ffffff;

      &:hover:not(:disabled) {
        background: #333;
      }
    }

    &-secondary {
      background: #ffffff;
      color: #1a1a1a;
      border: 1px solid #e5e5e5;

      &:hover:not(:disabled) {
        background: #fafafa;
        border-color: #1a1a1a;
      }
    }

    &-text {
      background: transparent;
      color: #666;

      &:hover:not(:disabled) {
        color: #1a1a1a;
        background: #fafafa;
      }
    }
  }
}

@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
    gap: 24px;
  }

  .form-actions {
    flex-direction: column;

    .btn {
      width: 100%;
    }
  }
}
</style>
