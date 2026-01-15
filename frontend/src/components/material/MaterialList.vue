<script setup lang="ts">
import { computed, onMounted } from 'vue'
import MaterialCard from './MaterialCard.vue'
import Pagination from '@/components/Pagination.vue'
import type { Material } from '@/types'

interface Props {
  materials: Material[]
  loading: boolean
  total: number
  page: number
  size: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'page-change', page: number): void
  (e: 'size-change', size: number): void
}>()

// 是否有数据
const hasData = computed(() => props.materials.length > 0)

// 处理页码变化
const handlePageChange = (page: number) => {
  emit('page-change', page)
}

// 处理每页大小变化
const handleSizeChange = (size: number) => {
  emit('size-change', size)
}
</script>

<template>
  <div class="material-list">
    <!-- 加载状态 -->
    <div v-if="loading" class="material-list__loading">
      <el-skeleton :animated="true" :count="6">
        <template #template>
          <el-skeleton-item variant="rect" style="height: 200px; margin-bottom: 16px" />
        </template>
      </el-skeleton>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!hasData" class="material-list__empty">
      <el-empty description="暂无资料" />
    </div>

    <!-- 资料列表 -->
    <template v-else>
      <div class="material-list__grid">
        <MaterialCard
          v-for="material in materials"
          :key="material.id"
          :material="material"
        />
      </div>

      <!-- 分页 -->
      <div class="material-list__pagination">
        <Pagination
          :total="total"
          :page="page"
          :size="size"
          @page-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </template>
  </div>
</template>

<style scoped lang="scss">
.material-list {
  &__grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }

  &__pagination {
    display: flex;
    justify-content: center;
    padding-top: 16px;
    border-top: 1px solid var(--el-border-color-lighter);
  }

  &__empty,
  &__loading {
    padding: 60px 0;
  }
}

// 响应式布局
@media (max-width: 768px) {
  .material-list {
    &__grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
