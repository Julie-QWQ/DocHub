<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useMaterialStore } from '@/stores/material'
import MaterialCard from '@/components/material/MaterialCard.vue'
import type { MaterialListParams } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const materialStore = useMaterialStore()

const currentTab = ref('all')

// 标签页配置
const tabs = [
  { label: '全部', name: 'all', status: undefined },
  { label: '待审核', name: 'pending', status: 'pending' as const },
  { label: '已通过', name: 'approved', status: 'approved' as const },
  { label: '已拒绝', name: 'rejected', status: 'rejected' as const }
]

// 判断用户是否可以上传资料
const canUploadMaterial = computed(() => authStore.isCommittee)

// 当前筛选条件
const filters = ref<MaterialListParams>({
  page: 1,
  size: 20,
  sort_by: 'created_at',
  order: 'desc'
})

// 加载资料列表
const loadMaterials = async () => {
  await materialStore.fetchMaterials({
    ...filters.value,
    status: tabs.find((t) => t.name === currentTab.value)?.status,
    uploader_id: authStore.user?.id // 传递当前用户ID作为上传者筛选
  })
}

// 处理标签页变化
const handleTabChange = async () => {
  filters.value.page = 1
  await loadMaterials()
}

// 处理页码变化
const handlePageChange = async (page: number) => {
  filters.value.page = page
  await loadMaterials()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  loadMaterials()
})
</script>

<template>
  <div class="my-materials-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题和操作 -->
        <header class="stream-header">
          <h1 class="page-title">我的资料</h1>
          <el-button
            v-if="canUploadMaterial"
            type="primary"
            @click="router.push('/materials/upload')"
            class="upload-btn"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14M5 12h14"/>
            </svg>
            上传资料
          </el-button>
        </header>

        <!-- 标签页 -->
        <div class="tabs-section">
          <button
            v-for="tab in tabs"
            :key="tab.name"
            :class="['tab-button', { active: currentTab === tab.name }]"
            @click="currentTab = tab.name; handleTabChange()"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- 加载状态 -->
        <div v-if="materialStore.loading" class="loading-state">
          <div v-for="i in 6" :key="i" class="skeleton-card"></div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="materialStore.materials.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
          </svg>
          <h2>暂无资料</h2>
          <p>你上传的资料会显示在这里</p>
          <button v-if="canUploadMaterial" @click="router.push('/materials/upload')" class="upload-first-btn">
            上传第一个资料
          </button>
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
.my-materials-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.content-stream {
  max-width: 100%;
}

.stream-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  .page-title {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0;
  }

  .upload-btn {
    background: #111827;
    color: #ffffff;
    border: none;
    padding: 10px 20px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 6px;

    &:hover {
      background: #000000;
    }

    svg {
      width: 16px;
      height: 16px;
    }
  }
}

.tabs-section {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;

  .tab-button {
    padding: 8px 16px;
    background: none;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #f3f4f6;
      color: #111827;
    }

    &.active {
      background: #111827;
      color: #ffffff;
      font-weight: 500;
    }
  }
}

.materials-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

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
    margin: 0 0 24px 0;
  }

  .upload-first-btn {
    background: #111827;
    color: #ffffff;
    border: none;
    padding: 12px 24px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #000000;
    }
  }
}

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

@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;

    .page-title {
      font-size: 24px;
    }
  }

  .tabs-section {
    flex-wrap: wrap;
  }
}
</style>
