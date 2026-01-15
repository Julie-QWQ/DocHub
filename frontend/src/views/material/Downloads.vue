<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { materialApi } from '@/api/material'
import type { Material } from '@/types'

const router = useRouter()
const loading = ref(false)
const downloadRecords = ref<Material[]>([])
const total = ref(0)
const page = ref(1)
const size = ref(20)

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  if (days < 30) return `${Math.floor(days / 7)} 周前`
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

// 加载下载记录
const loadDownloadRecords = async () => {
  loading.value = true
  try {
    const response = await materialApi.getDownloadRecords({
      page: page.value,
      page_size: size.value
    })

    if (response.code === 0 && response.data) {
      downloadRecords.value = response.data.materials || response.data.list || []
      total.value = response.data.total
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载下载记录失败')
  } finally {
    loading.value = false
  }
}

// 处理页码变化
const handlePageChange = async (newPage: number) => {
  page.value = newPage
  await loadDownloadRecords()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 跳转到资料详情
const goToDetail = (materialId: number) => {
  router.push(`/materials/${materialId}`)
}

onMounted(() => {
  loadDownloadRecords()
})
</script>

<template>
  <div class="downloads-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题 -->
        <header class="stream-header">
          <h1 class="page-title">下载历史</h1>
        </header>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div v-for="i in 6" :key="i" class="skeleton-card"></div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="downloadRecords.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          <h2>暂无下载记录</h2>
          <p>下载过的资料会显示在这里</p>
        </div>

        <!-- 下载记录列表 -->
        <div v-else class="downloads-list">
          <article
            v-for="record in downloadRecords"
            :key="record.id"
            class="download-item"
            @click="goToDetail(record.id)"
          >
            <div class="material-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
              </svg>
            </div>

            <div class="material-info">
              <div class="material-title">{{ record.title }}</div>
              <div class="material-meta">
                <span v-if="record.course_name" class="course-name">{{ record.course_name }}</span>
                <span v-if="record.file_name" class="separator">·</span>
                <span v-if="record.file_name" class="file-name">{{ record.file_name }}</span>
              </div>
              <div class="material-details">
                <span>{{ formatFileSize(record.file_size) }}</span>
                <span class="separator">·</span>
                <span>{{ formatTime(record.created_at) }}</span>
              </div>
            </div>

            <div class="material-actions">
              <button class="action-btn primary">查看详情</button>
            </div>
          </article>
        </div>

        <!-- 分页 -->
        <div v-if="total > 0" class="pagination-section">
          <el-pagination
            :current-page="page"
            :page-size="size"
            :total="total"
            layout="prev, pager, next"
            @current-change="handlePageChange"
          />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.downloads-page {
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
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f2f2f2;

  .page-title {
    font-size: 32px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0;
    letter-spacing: -0.02em;
  }
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

// 下载记录列表
.downloads-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.download-item {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 16px;
  align-items: center;
  padding: 20px 0;
  background: transparent;
  border: none;
  border-bottom: 1px solid #f2f2f2;
  border-radius: 0;
  transition: all 0.2s;
  cursor: pointer;

  &:hover {
    background: #f9fafb;
  }

  &:first-child {
    border-top: 1px solid #f2f2f2;
  }

  @media (max-width: 768px) {
    grid-template-columns: auto 1fr;
    gap: 12px;

    .material-actions {
      grid-column: 2;
      justify-content: flex-end;
    }
  }
}

.material-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  background: #f3f4f6;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.material-info {
  min-width: 0;
}

.material-title {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 6px;
}

.material-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
  font-size: 13px;
  color: #6b7280;
}

.separator {
  color: #d1d5db;
}

.file-name {
  font-size: 13px;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.material-details {
  font-size: 13px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  gap: 6px;
}

.material-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 8px 14px;
  border: none;
  border-radius: 6px;
  background: #f3f4f6;
  color: #374151;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #e5e7eb;
  }

  &.primary {
    background: #111827;
    color: #ffffff;

    &:hover {
      background: #000000;
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
      font-size: 26px;
    }
  }
}

// 分页
.pagination-section {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #f2f2f2;
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
</style>
