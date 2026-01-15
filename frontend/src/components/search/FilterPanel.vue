<template>
  <div class="filter-panel">
    <el-card class="filter-card">
      <template #header>
        <div class="card-header">
          <span>筛选条件</span>
          <el-button link type="primary" @click="handleReset">重置</el-button>
        </div>
      </template>

      <el-form :model="filters" label-width="80px" label-position="left">
        <!-- 资料分类 -->
        <el-form-item label="资料分类">
          <el-select
            v-model="filters.category"
            placeholder="请选择分类"
            clearable
            @change="handleFilterChange"
          >
            <el-option label="课件" value="courseware" />
            <el-option label="试卷" value="exam" />
            <el-option label="实验" value="experiment" />
            <el-option label="习题" value="exercise" />
            <el-option label="参考资料" value="reference" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>

        <!-- 课程名称 -->
        <el-form-item label="课程名称">
          <el-input
            v-model="filters.courseName"
            placeholder="请输入课程名称"
            clearable
            @keyup.enter="handleFilterChange"
            @clear="handleFilterChange"
          />
        </el-form-item>

        <!-- 标签 -->
        <el-form-item label="标签">
          <el-select
            v-model="filters.tags"
            placeholder="请选择标签"
            multiple
            collapse-tags
            collapse-tags-tooltip
            @change="handleFilterChange"
          >
            <el-option label="重点" value="重点" />
            <el-option label="难点" value="难点" />
            <el-option label="易错" value="易错" />
            <el-option label="基础" value="基础" />
            <el-option label="进阶" value="进阶" />
          </el-select>
        </el-form-item>

        <!-- 时间范围 -->
        <el-form-item label="上传时间">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="handleFilterChange"
          />
        </el-form-item>

        <!-- 排序方式 -->
        <el-form-item label="排序方式">
          <el-select v-model="sortBy" @change="handleSortChange">
            <el-option label="上传时间" value="created_at" />
            <el-option label="下载量" value="download_count" />
            <el-option label="收藏量" value="favorite_count" />
            <el-option label="浏览量" value="view_count" />
          </el-select>
        </el-form-item>

        <!-- 排序方向 -->
        <el-form-item label="排序方向">
          <el-radio-group v-model="sortOrder" @change="handleSortChange">
            <el-radio-button label="desc">降序</el-radio-button>
            <el-radio-button label="asc">升序</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { MaterialCategory } from '@/api/search'

interface Filters {
  category?: MaterialCategory
  courseName: string
  tags: string[]
  dateRange: [string, string] | null
}

interface Props {
  modelValue: Filters
  sortBy?: string
  sortOrder?: 'desc' | 'asc'
}

const props = withDefaults(defineProps<Props>(), {
  sortBy: 'created_at',
  sortOrder: 'desc'
})

const emit = defineEmits<{
  'update:modelValue': [value: Filters]
  'update:sortBy': [value: string]
  'update:sortOrder': [value: 'desc' | 'asc']
  'filter-change': []
  'sort-change': []
}>()

const filters = reactive<Filters>({
  category: props.modelValue.category,
  courseName: props.modelValue.courseName,
  tags: props.modelValue.tags || [],
  dateRange: props.modelValue.dateRange
})

const sortBy = ref(props.sortBy)
const sortOrder = ref<'desc' | 'asc'>(props.sortOrder)

const handleFilterChange = () => {
  emit('update:modelValue', { ...filters })
  emit('filter-change')
}

const handleSortChange = () => {
  emit('update:sortBy', sortBy.value)
  emit('update:sortOrder', sortOrder.value)
  emit('sort-change')
}

const handleReset = () => {
  filters.category = undefined
  filters.courseName = ''
  filters.tags = []
  filters.dateRange = null
  sortBy.value = 'created_at'
  sortOrder.value = 'desc'
  handleFilterChange()
  handleSortChange()
}
</script>

<style scoped lang="scss">
.filter-panel {
  .filter-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    :deep(.el-form-item) {
      margin-bottom: 16px;
    }

    :deep(.el-select),
    :deep(.el-input) {
      width: 100%;
    }
  }
}
</style>
