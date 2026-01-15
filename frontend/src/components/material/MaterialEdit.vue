<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { materialApi } from '@/api/material'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import type { Material, MaterialCategory } from '@/types'

interface Props {
  materialId: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'cancel'): void
}>()

const materialCategoryStore = useMaterialCategoryStore()

const loading = ref(false)
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
  await loadMaterial()
})

// 加载资料详情
const loadMaterial = async () => {
  loading.value = true
  try {
    const response = await materialApi.getMaterial(props.materialId)
    if (response.code === 0 && response.data) {
      const material = response.data
      formData.title = material.title || ''
      formData.description = material.description || ''
      formData.category = material.category || '' as MaterialCategory
      formData.course_name = material.course_name || ''
      formData.tags = material.tags || []
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载资料失败')
  } finally {
    loading.value = false
  }
}

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
    await materialApi.updateMaterial(props.materialId, {
      title: formData.title,
      description: formData.description,
      category: formData.category,
      course_name: formData.course_name,
      tags: formData.tags
    })

    ElMessage.success('资料更新成功')
    emit('success')
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  emit('cancel')
}

// 组件挂载时加载资料
onMounted(() => {
  loadMaterial()
})
</script>

<template>
  <div class="edit-form">
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <form v-else @submit.prevent="handleSubmit">
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
            <option v-for="option in categoryOptions" :key="option.value" :value="option.value">
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
            placeholder="如：高等数学"
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
          maxlength="2000"
        ></textarea>
        <span class="hint">最多 2000 个字符</span>
      </div>

      <!-- 标签 -->
      <div class="form-group">
        <label class="form-label">标签</label>
        <div class="tags-input-container">
          <div class="tags-list">
            <span
              v-for="tag in formData.tags"
              :key="tag"
              class="tag-item"
            >
              {{ tag }}
              <button
                type="button"
                class="tag-remove"
                @click="handleRemoveTag(tag)"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </span>
          </div>
          <div class="tag-input-row">
            <input
              v-model="tagInput"
              type="text"
              class="tag-input"
              placeholder="输入标签后按回车添加"
              maxlength="50"
              @keydown.enter.prevent="handleAddTag"
            />
            <button
              type="button"
              class="tag-add-btn"
              @click="handleAddTag"
            >
              添加
            </button>
          </div>
        </div>
        <span class="hint">最多添加 10 个标签，每个标签最多 50 个字符</span>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <button
          type="button"
          class="btn btn-secondary"
          @click="handleCancel"
          :disabled="submitting"
        >
          取消
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          :disabled="submitting"
        >
          <span v-if="submitting">保存中...</span>
          <span v-else>保存修改</span>
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped lang="scss">
.edit-form {
  max-width: 700px;
  margin: 0 auto;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #666;

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 3px solid #f3f4f6;
    border-top-color: #1a1a1a;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin-bottom: 16px;
  }

  p {
    margin: 0;
    font-size: 14px;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.form-group {
  margin-bottom: 24px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.half {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;

  .required {
    color: #dc2626;
  }
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  font-size: 14px;
  line-height: 1.5;
  color: #1a1a1a;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  transition: all 0.15s;

  &::placeholder {
    color: #9ca3af;
  }

  &:focus {
    outline: none;
    border-color: #1a1a1a;
    box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.1);
  }

  &.error {
    border-color: #dc2626;
  }
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.error-message {
  display: block;
  margin-top: 6px;
  font-size: 13px;
  color: #dc2626;
}

.hint {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: #9ca3af;
}

.tags-input-container {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
  transition: all 0.15s;

  &:focus-within {
    border-color: #1a1a1a;
    box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.1);
  }
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 13px;
  color: #374151;
}

.tag-remove {
  display: flex;
  align-items: center;
  padding: 2px;
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.15s;

  &:hover {
    color: #1a1a1a;
    background: rgba(0, 0, 0, 0.05);
  }

  svg {
    width: 12px;
    height: 12px;
  }
}

.tag-input-row {
  display: flex;
  gap: 8px;
}

.tag-input {
  flex: 1;
  padding: 8px 10px;
  font-size: 14px;
  border: none;
  background: transparent;
  outline: none;

  &::placeholder {
    color: #9ca3af;
  }
}

.tag-add-btn {
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 500;
  color: #ffffff;
  background: #1a1a1a;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;

  &:hover {
    background: #000000;
    transform: translateY(-1px);
  }

  &:active {
    transform: translateY(0);
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #f2f2f2;
}

.btn {
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &.btn-secondary {
    background: #f3f4f6;
    color: #1a1a1a;

    &:hover:not(:disabled) {
      background: #e5e7eb;
    }
  }

  &.btn-primary {
    background: #1a1a1a;
    color: #ffffff;

    &:hover:not(:disabled) {
      background: #000000;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }

    &:active:not(:disabled) {
      transform: translateY(0);
    }
  }
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column-reverse;

    .btn {
      width: 100%;
    }
  }
}
</style>
