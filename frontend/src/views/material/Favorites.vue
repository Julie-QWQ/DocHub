<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useMaterialStore } from '@/stores/material'
import MaterialCard from '@/components/material/MaterialCard.vue'
import type { MaterialListParams } from '@/types'

const materialStore = useMaterialStore()

// 当前筛选条件
const filters = ref<MaterialListParams>({
  page: 1,
  size: 20
})

// 加载收藏列表
const loadFavorites = async () => {
  await materialStore.fetchFavorites(filters.value)
}

// 处理页码变化
const handlePageChange = async (page: number) => {
  filters.value.page = page
  await loadFavorites()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  loadFavorites()
})
</script>

<template>
  <div class="favorites-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题 -->
        <header class="stream-header">
          <h1 class="page-title">我的收藏</h1>
        </header>

        <!-- 加载状态 -->
        <div v-if="materialStore.loading" class="loading-state">
          <div v-for="i in 6" :key="i" class="skeleton-card"></div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="materialStore.materials.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"/>
          </svg>
          <h2>暂无收藏</h2>
          <p>收藏的资料会显示在这里</p>
        </div>

        <!-- 资料列表 -->
        <div v-else class="materials-list">
          <MaterialCard
            v-for="material in materialStore.materials"
            :key="material.id"
            :material="material"
          />
        </div>

        <!-- 分页 -->
        <div v-if="materialStore.total > 0" class="pagination-section">
          <el-pagination
            :current-page="materialStore.page"
            :page-size="materialStore.size"
            :total="materialStore.total"
            layout="prev, pager, next"
            @current-change="handlePageChange"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.favorites-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

// 中间内容流
.content-stream {
  max-width: 100%;
}

.stream-header {
  margin-bottom: 24px;

  .page-title {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0;
  }
}

// 资料列表
.materials-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

// 加载状态
.loading-state {
  display: flex;
  flex-direction: column;
  gap: 0;

  .skeleton-card {
    height: 200px;
    background: linear-gradient(90deg, #f5f5f5 25%, #e8e8e8 50%, #f5f5f5 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
    border-bottom: 1px solid #f2f2f2;

    &:last-child {
      border-bottom: none;
    }
  }
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

// 空状态
.empty-state {
  text-align: center;
  padding: 80px 40px;
  background: #f9fafb;
  border-radius: 12px;

  svg {
    color: #d1d5db;
    margin-bottom: 24px;
  }

  h2 {
    font-size: 22px;
    font-weight: 600;
    color: #111827;
    margin: 0 0 8px 0;
  }

  p {
    font-size: 15px;
    color: #9ca3af;
    margin: 0;
  }
}

// 分页
.pagination-section {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination) {
  .btn-prev,
  .btn-next,
  .el-pager li {
    background: none;
    border: none;
    font-weight: 500;
    color: #6b7280;

    &:hover {
      color: #111827;
    }

    &.active {
      color: #111827;
      font-weight: 600;
    }
  }
}

// 响应式
@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    .page-title {
      font-size: 24px;
    }
  }
}
</style>
