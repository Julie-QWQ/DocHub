<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import type { MaterialCategory, MaterialStatus } from '@/types'

interface Filters {
  category?: MaterialCategory
  status?: MaterialStatus
  course_name?: string
  sort_by?: string
  order?: 'asc' | 'desc'
}

const props = defineProps<{
  modelValue: Filters
  showStatus?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Filters): void
  (e: 'search'): void
  (e: 'reset'): void
}>()

const materialCategoryStore = useMaterialCategoryStore()

// 资料分类选项（从 store 动态获取）
const categoryOptions = computed(() => {
  const allOption = { label: '全部分类', value: '' as MaterialCategory | '' }
  const activeOptions = materialCategoryStore.activeCategories.map(cat => ({
    label: cat.name,
    value: cat.code as MaterialCategory
  }))
  return [allOption, ...activeOptions]
})

// 加载资料类型
onMounted(async () => {
  try {
    await materialCategoryStore.fetchActiveCategories()
  } catch (error: any) {
    console.error('加载资料类型失败:', error)
  }
})

// 资料状态选项
const statusOptions = [
  { label: '全部状态', value: '' },
  { label: '待审核', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已拒绝', value: 'rejected' }
]

// 排序选项
const sortOptions = [
  { label: '最新发布', value: 'created_at' },
  { label: '最近更新', value: 'updated_at' },
  { label: '下载量', value: 'download_count' },
  { label: '浏览量', value: 'view_count' },
  { label: '收藏量', value: 'favorite_count' }
]

// 本地表单数据
const formData = ref<Filters>({ ...props.modelValue })

// 监听外部变化
watch(
  () => props.modelValue,
  (newVal) => {
    formData.value = { ...newVal }
  },
  { deep: true }
)

// 更新筛选条件
const updateFilters = () => {
  const filters: Filters = {}

  if (formData.value.category) filters.category = formData.value.category
  if (formData.value.status) filters.status = formData.value.status
  if (formData.value.course_name) filters.course_name = formData.value.course_name
  if (formData.value.sort_by) filters.sort_by = formData.value.sort_by as any
  if (formData.value.order) filters.order = formData.value.order

  emit('update:modelValue', filters)
  emit('search')
}

// 重置筛选条件
const resetFilters = () => {
  formData.value = {
    sort_by: 'created_at',
    order: 'desc'
  }
  emit('update:modelValue', formData.value)
  emit('reset')
}

// 切换排序方向
const toggleOrder = () => {
  formData.value.order = formData.value.order === 'asc' ? 'desc' : 'asc'
  updateFilters()
}
</script>

<template>
  <div class="material-filters">
    <el-form :model="formData" inline>
      <!-- 分类筛选 -->
      <el-form-item label="分类">
        <el-select
          v-model="formData.category"
          placeholder="请选择分类"
          clearable
          @change="updateFilters"
        >
          <el-option
            v-for="option in categoryOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value || undefined"
          />
        </el-select>
      </el-form-item>

      <!-- 状态筛选（可选显示） -->
      <el-form-item v-if="showStatus" label="状态">
        <el-select
          v-model="formData.status"
          placeholder="请选择状态"
          clearable
          @change="updateFilters"
        >
          <el-option
            v-for="option in statusOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value || undefined"
          />
        </el-select>
      </el-form-item>

      <!-- 课程名称筛选 -->
      <el-form-item label="课程">
        <el-input
          v-model="formData.course_name"
          placeholder="请输入课程名称"
          clearable
          @keyup.enter="updateFilters"
        />
      </el-form-item>

      <!-- 排序方式 -->
      <el-form-item label="排序">
        <el-select
          v-model="formData.sort_by"
          placeholder="请选择排序方式"
          @change="updateFilters"
        >
          <el-option
            v-for="option in sortOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
      </el-form-item>

      <!-- 排序方向 -->
      <el-form-item>
        <el-button circle @click="toggleOrder">
          <el-icon>
            <component :is="formData.order === 'asc' ? 'SortDescending' : 'SortAscending'" />
          </el-icon>
        </el-button>
      </el-form-item>

      <!-- 搜索按钮 -->
      <el-form-item>
        <el-button type="primary" :icon="'Search'" @click="updateFilters">
          搜索
        </el-button>
      </el-form-item>

      <!-- 重置按钮 -->
      <el-form-item>
        <el-button :icon="'RefreshLeft'" @click="resetFilters">
          重置
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<style scoped lang="scss">
.material-filters {
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;

  :deep(.el-form) {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
  }

  :deep(.el-form-item) {
    margin-bottom: 0;
  }
}

@media (max-width: 768px) {
  .material-filters {
    :deep(.el-form) {
      flex-direction: column;
    }

    :deep(.el-form-item) {
      width: 100%;
    }

    :deep(.el-select),
    :deep(.el-input) {
      width: 100%;
    }
  }
}
</style>
