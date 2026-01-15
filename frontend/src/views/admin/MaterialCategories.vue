<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import type { MaterialCategoryConfig, MaterialCategoryConfigRequest } from '@/types'

const materialCategoryStore = useMaterialCategoryStore()

const loading = ref(false)
const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const deleteDialogVisible = ref(false)
const categoryToDelete = ref<MaterialCategoryConfig | null>(null)
const formData = ref<MaterialCategoryConfigRequest>({
  code: '',
  name: '',
  description: '',
  icon: '',
  sort_order: 0,
  is_active: true
})
const currentEditId = ref<number | null>(null)

// 常见图标列表
const iconOptions = [
  { label: '文档 (document)', value: 'document' },
  { label: '书籍 (book)', value: 'book' },
  { label: '参考 (reference)', value: 'reference' },
  { label: '考试 (exam)', value: 'exam' },
  { label: '练习 (exercise)', value: 'exercise' },
  { label: '实验 (experiment)', value: 'experiment' },
  { label: '笔记 (note)', value: 'note' },
  { label: '论文 (thesis)', value: 'thesis' },
  { label: '其他 (other)', value: 'other' }
]

// 获取资料类型列表
const fetchCategories = async () => {
  loading.value = true
  try {
    await materialCategoryStore.fetchAllCategories()
  } finally {
    loading.value = false
  }
}

// 排序后的列表
const sortedCategories = computed(() => {
  return [...materialCategoryStore.categories].sort((a, b) => a.sort_order - b.sort_order)
})

// 打开创建对话框
const openCreateDialog = () => {
  dialogMode.value = 'create'
  currentEditId.value = null
  formData.value = {
    code: '',
    name: '',
    description: '',
    icon: '',
    sort_order: sortedCategories.value.length + 1,
    is_active: true
  }
  dialogVisible.value = true
}

// 打开编辑对话框
const openEditDialog = (category: MaterialCategoryConfig) => {
  dialogMode.value = 'edit'
  currentEditId.value = category.id
  formData.value = {
    code: category.code,
    name: category.name,
    description: category.description || '',
    icon: category.icon || '',
    sort_order: category.sort_order,
    is_active: category.is_active
  }
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  // 验证
  if (!formData.value.code.trim()) {
    ElMessage.error('请输入类型代码')
    return
  }
  if (!formData.value.name.trim()) {
    ElMessage.error('请输入类型名称')
    return
  }

  loading.value = true
  try {
    if (dialogMode.value === 'create') {
      await materialCategoryStore.createCategory(formData.value)
      ElMessage.success('创建成功')
    } else {
      if (currentEditId.value) {
        await materialCategoryStore.updateCategory(currentEditId.value, formData.value)
        ElMessage.success('更新成功')
      }
    }
    dialogVisible.value = false
    await fetchCategories()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

// 删除资料类型
const handleDelete = (category: MaterialCategoryConfig) => {
  categoryToDelete.value = category
  deleteDialogVisible.value = true
}

// 确认删除
const confirmDelete = async () => {
  if (!categoryToDelete.value) return

  loading.value = true
  try {
    await materialCategoryStore.deleteCategory(categoryToDelete.value.id)
    deleteDialogVisible.value = false
    categoryToDelete.value = null
    await fetchCategories()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  } finally {
    loading.value = false
  }
}

// 取消删除
const cancelDelete = () => {
  deleteDialogVisible.value = false
  categoryToDelete.value = null
}

// 切换启用状态
const handleToggleStatus = async (category: MaterialCategoryConfig) => {
  loading.value = true
  try {
    await materialCategoryStore.toggleStatus(category.id)
    ElMessage.success(category.is_active ? '已禁用' : '已启用')
    await fetchCategories()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

// 上移
const moveUp = async (index: number) => {
  if (index === 0) return

  const categories = sortedCategories.value
  const current = categories[index]
  const prev = categories[index - 1]

  // 交换排序
  const tempSort = current.sort_order
  current.sort_order = prev.sort_order
  prev.sort_order = tempSort

  // 更新
  loading.value = true
  try {
    await materialCategoryStore.updateCategory(current.id, {
      code: current.code,
      name: current.name,
      description: current.description,
      icon: current.icon,
      sort_order: current.sort_order,
      is_active: current.is_active
    })
    await materialCategoryStore.updateCategory(prev.id, {
      code: prev.code,
      name: prev.name,
      description: prev.description,
      icon: prev.icon,
      sort_order: prev.sort_order,
      is_active: prev.is_active
    })
    await fetchCategories()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

// 下移
const moveDown = async (index: number) => {
  if (index === sortedCategories.value.length - 1) return

  const categories = sortedCategories.value
  const current = categories[index]
  const next = categories[index + 1]

  // 交换排序
  const tempSort = current.sort_order
  current.sort_order = next.sort_order
  next.sort_order = tempSort

  // 更新
  loading.value = true
  try {
    await materialCategoryStore.updateCategory(current.id, {
      code: current.code,
      name: current.name,
      description: current.description,
      icon: current.icon,
      sort_order: current.sort_order,
      is_active: current.is_active
    })
    await materialCategoryStore.updateCategory(next.id, {
      code: next.code,
      name: next.name,
      description: next.description,
      icon: next.icon,
      sort_order: next.sort_order,
      is_active: next.is_active
    })
    await fetchCategories()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<template>
  <div class="material-categories-page">
    <div class="page-header">
      <h1 class="page-title">资料类型管理</h1>
      <button class="btn btn-primary" @click="openCreateDialog">
        新增类型
      </button>
    </div>

    <div class="categories-list" v-loading="loading">
      <div
        v-for="(category, index) in sortedCategories"
        :key="category.id"
        class="category-item"
        :class="{ disabled: !category.is_active }"
      >
        <div class="category-info">
          <div class="category-sort">
            <span class="sort-number">{{ category.sort_order }}</span>
            <div class="sort-controls">
              <button
                class="sort-btn"
                :disabled="index === 0"
                @click="moveUp(index)"
              >
                ↑
              </button>
              <button
                class="sort-btn"
                :disabled="index === sortedCategories.length - 1"
                @click="moveDown(index)"
              >
                ↓
              </button>
            </div>
          </div>
          <div class="category-details">
            <div class="category-name">{{ category.name }}</div>
            <div class="category-meta">
              <span class="category-code">代码: {{ category.code }}</span>
              <span v-if="category.description" class="category-description">
                {{ category.description }}
              </span>
            </div>
          </div>
        </div>
        <div class="category-actions">
          <button
            class="btn btn-sm"
            :class="category.is_active ? 'btn-warning' : 'btn-success'"
            @click="handleToggleStatus(category)"
          >
            {{ category.is_active ? '禁用' : '启用' }}
          </button>
          <button
            class="btn btn-sm btn-secondary"
            @click="openEditDialog(category)"
          >
            编辑
          </button>
          <button
            class="btn btn-sm btn-danger"
            @click="handleDelete(category)"
          >
            删除
          </button>
        </div>
      </div>

      <div v-if="sortedCategories.length === 0 && !loading" class="empty-state">
        <p>暂无资料类型</p>
      </div>
    </div>

    <!-- 创建/编辑对话框 -->
    <div v-if="dialogVisible" class="dialog-overlay" @click.self="dialogVisible = false">
      <div class="dialog">
        <div class="dialog-header">
          <h2>{{ dialogMode === 'create' ? '新增资料类型' : '编辑资料类型' }}</h2>
          <button class="dialog-close" @click="dialogVisible = false">×</button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label class="form-label">类型代码 <span class="required">*</span></label>
            <input
              v-model="formData.code"
              type="text"
              class="form-input"
              placeholder="如: courseware, exam_paper"
              :disabled="dialogMode === 'edit'"
            />
            <span class="help-text">英文代码,用于系统标识,创建后不可修改</span>
          </div>

          <div class="form-group">
            <label class="form-label">类型名称 <span class="required">*</span></label>
            <input
              v-model="formData.name"
              type="text"
              class="form-input"
              placeholder="如: 课件、试卷"
            />
          </div>

          <div class="form-group">
            <label class="form-label">描述</label>
            <textarea
              v-model="formData.description"
              class="form-textarea"
              rows="3"
              placeholder="资料类型描述"
            ></textarea>
          </div>

          <div class="form-group">
            <label class="form-label">图标</label>
            <select v-model="formData.icon" class="form-select">
              <option value="">请选择图标</option>
              <option v-for="icon in iconOptions" :key="icon.value" :value="icon.value">
                {{ icon.label }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">排序</label>
            <input
              v-model.number="formData.sort_order"
              type="number"
              class="form-input"
              min="1"
            />
            <span class="help-text">数字越小越靠前</span>
          </div>

          <div class="form-group">
            <label class="form-checkbox">
              <input type="checkbox" v-model="formData.is_active" />
              <span>启用此类型</span>
            </label>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="dialogVisible = false">
            取消
          </button>
          <button class="btn btn-primary" @click="handleSubmit" :disabled="loading">
            {{ dialogMode === 'create' ? '创建' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <div v-if="deleteDialogVisible" class="dialog-overlay" @click.self="cancelDelete">
      <div class="dialog dialog-delete">
        <div class="dialog-header">
          <h2>删除确认</h2>
          <button class="dialog-close" @click="cancelDelete">×</button>
        </div>
        <div class="dialog-body">
          <div class="delete-warning">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <p>确定要删除资料类型"<strong>{{ categoryToDelete?.name }}</strong>"吗?</p>
            <p class="warning-text">此操作不可恢复,请谨慎操作。</p>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="cancelDelete">
            取消
          </button>
          <button class="btn btn-danger" @click="confirmDelete" :disabled="loading">
            确定删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.material-categories-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0;
}

.categories-list {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f2f2f2;
  transition: background 0.2s;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: #fafafa;
  }

  &.disabled {
    opacity: 0.5;
  }
}

.category-info {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
}

.category-sort {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sort-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: #f5f5f5;
  border-radius: 6px;
  font-weight: 600;
  color: #666;
}

.sort-controls {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.sort-btn {
  padding: 2px 6px;
  font-size: 12px;
  background: transparent;
  border: 1px solid #e5e5e5;
  border-radius: 3px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    background: #f5f5f5;
    border-color: #1a1a1a;
  }

  &:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }
}

.category-details {
  flex: 1;
}

.category-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.category-meta {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: #999;
}

.category-code {
  font-family: 'Courier New', monospace;
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
}

.category-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
  color: #999;
}

// 对话框样式
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: #ffffff;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f2f2f2;

  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
  }
}

.dialog-close {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  transition: color 0.2s;

  &:hover {
    color: #1a1a1a;
  }
}

.dialog-body {
  padding: 24px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #f2f2f2;
}

.form-group {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
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
  padding: 10px 12px;
  border: 1px solid #e5e5e5;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;

  &:focus {
    outline: none;
    border-color: #1a1a1a;
  }

  &:disabled {
    background: #f5f5f5;
    cursor: not-allowed;
  }
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.form-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;

  input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: pointer;
  }

  span {
    font-size: 14px;
    color: #1a1a1a;
  }
}

.help-text {
  display: block;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &.btn-primary {
    background: #1a1a1a;
    color: #ffffff;

    &:hover:not(:disabled) {
      background: #333;
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  &.btn-secondary {
    background: #ffffff;
    border: 1px solid #e5e5e5;
    color: #1a1a1a;

    &:hover {
      background: #fafafa;
      border-color: #1a1a1a;
    }
  }

  &.btn-success {
    background: #059669;
    color: #ffffff;

    &:hover {
      background: #047857;
    }
  }

  &.btn-warning {
    background: #d97706;
    color: #ffffff;

    &:hover {
      background: #b45309;
    }
  }

  &.btn-danger {
    background: #dc2626;
    color: #ffffff;

    &:hover {
      background: #b91c1c;
    }
  }

  &.btn-sm {
    padding: 6px 12px;
    font-size: 13px;
  }
}

// 删除确认对话框样式
.dialog-delete {
  max-width: 380px;

  .delete-warning {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 16px 0;

    svg {
      color: #f59e0b;
      margin-bottom: 16px;
      width: 40px;
      height: 40px;
    }

    p {
      font-size: 14px;
      line-height: 1.5;
      color: #1a1a1a;
      margin: 6px 0;

      strong {
        font-weight: 600;
        color: #dc2626;
      }
    }

    .warning-text {
      font-size: 12px;
      color: #999;
      margin-top: 12px;
    }
  }
}
</style>
