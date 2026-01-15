<template>
  <div class="app-pagination">
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="currentPageSize"
      :page-sizes="pageSizes"
      :total="total"
      :background="true"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  page?: number
  size?: number
  total: number
  pageSizes?: number[]
}

const props = withDefaults(defineProps<Props>(), {
  page: 1,
  size: 20,
  pageSizes: () => [10, 20, 50, 100]
})

const emit = defineEmits<{
  (e: 'update:page', page: number): void
  (e: 'update:size', size: number): void
  (e: 'change', page: number, size: number): void
}>()

const currentPage = computed({
  get: () => props.page,
  set: (val) => emit('update:page', val)
})

const currentPageSize = computed({
  get: () => props.size,
  set: (val) => emit('update:size', val)
})

const handleSizeChange = (size: number) => {
  emit('update:size', size)
  emit('change', currentPage.value, size)
}

const handleCurrentChange = (page: number) => {
  emit('update:page', page)
  emit('change', page, currentPageSize.value)
}
</script>

<style scoped lang="scss">
.app-pagination {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}
</style>
