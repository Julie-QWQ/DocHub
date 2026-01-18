<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useCommitteeStore } from '@/stores/committee'
import { format } from '@/utils/format'

const router = useRouter()
const committeeStore = useCommitteeStore()

const loading = ref(false)
const filterStatus = ref('')
const page = ref(1)
const size = ref(20)

const applications = computed(() => committeeStore.applications)
const total = computed(() => committeeStore.total)

const loadApplications = async () => {
  loading.value = true
  try {
    await committeeStore.fetchMyApplications({
      page: page.value,
      size: size.value,
      status: filterStatus.value || undefined
    })
  } finally {
    loading.value = false
  }
}

const handleApply = () => {
  router.push('/committee/apply')
}

const handleFilterChange = (status: string) => {
  filterStatus.value = status
  page.value = 1
  loadApplications()
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadApplications()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleCancel = async (id: number) => {
  await committeeStore.cancelApplication(id)
  ElMessage.success('申请已取消')
  await loadApplications()
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    cancelled: '已取消'
  }
  return textMap[status] || status
}

const getStatusClass = (status: string) => {
  const classMap: Record<string, string> = {
    pending: 'status-pending',
    approved: 'status-approved',
    rejected: 'status-rejected',
    cancelled: 'status-cancelled'
  }
  return classMap[status] || ''
}

onMounted(() => {
  loadApplications()
})
</script>

<template>
  <div class="applications-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题 -->
        <header class="stream-header">
          <h1 class="page-title">我的申请记录</h1>
          <button class="apply-btn" @click="handleApply">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14M5 12h14"/>
            </svg>
            新建申请
          </button>
        </header>

        <!-- 筛选器 -->
        <div class="filters-section">
          <button
            :class="['filter-btn', { active: filterStatus === '' }]"
            @click="handleFilterChange('')"
          >
            全部
          </button>
          <button
            :class="['filter-btn', { active: filterStatus === 'pending' }]"
            @click="handleFilterChange('pending')"
          >
            待审核
          </button>
          <button
            :class="['filter-btn', { active: filterStatus === 'approved' }]"
            @click="handleFilterChange('approved')"
          >
            已通过
          </button>
          <button
            :class="['filter-btn', { active: filterStatus === 'rejected' }]"
            @click="handleFilterChange('rejected')"
          >
            已拒绝
          </button>
          <button
            :class="['filter-btn', { active: filterStatus === 'cancelled' }]"
            @click="handleFilterChange('cancelled')"
          >
            已取消
          </button>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div v-for="i in 5" :key="i" class="skeleton-card"></div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="applications.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect width="18" height="18" x="3" y="3" rx="2"/>
            <path d="M3 9h18"/>
            <path d="M9 21V9"/>
          </svg>
          <h2>暂无申请记录</h2>
          <p>您还没有提交过学习委员申请</p>
          <button class="apply-btn-large" @click="handleApply">
            立即申请
          </button>
        </div>

        <!-- 申请列表 -->
        <div v-else class="applications-list">
          <div
            v-for="application in applications"
            :key="application.id"
            class="application-card"
          >
            <div class="card-header">
              <span :class="['status-badge', getStatusClass(application.status)]">
                {{ getStatusText(application.status) }}
              </span>
              <button
                v-if="application.status === 'pending'"
                class="cancel-btn"
                @click="handleCancel(application.id)"
              >
                取消申请
              </button>
            </div>

            <div class="card-body">
              <div class="application-reason">
                <label>申请理由</label>
                <p>{{ application.reason }}</p>
              </div>

              <div class="application-meta">
                <div class="meta-item">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <polyline points="12 6 12 12 16 14"/>
                  </svg>
                  <span>申请时间: {{ format.datetime(application.created_at) }}</span>
                </div>

                <div v-if="application.reviewed_at" class="meta-item">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <polyline points="12 6 12 12 16 14"/>
                  </svg>
                  <span>审核时间: {{ format.datetime(application.reviewed_at) }}</span>
                </div>

                <div v-if="application.reviewer" class="meta-item">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                  <span>审核人: {{ application.reviewer.real_name || application.reviewer.username }}</span>
                </div>
              </div>

              <div v-if="application.review_comment" class="application-comment">
                <label>审核意见</label>
                <div class="comment-content">{{ application.review_comment }}</div>
              </div>
            </div>
          </div>
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
.applications-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 900px;
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

  .apply-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 10px 20px;
    background: #1a1a1a;
    color: #ffffff;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #333;
    }

    svg {
      width: 18px;
      height: 18px;
    }
  }
}

.filters-section {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  overflow-x: auto;
  padding-bottom: 4px;

  .filter-btn {
    padding: 8px 16px;
    background: #ffffff;
    border: 1px solid #e5e5e5;
    border-radius: 6px;
    font-size: 14px;
    color: #666;
    cursor: pointer;
    transition: all 0.15s;
    white-space: nowrap;

    &:hover {
      border-color: #1a1a1a;
      color: #1a1a1a;
    }

    &.active {
      background: #1a1a1a;
      border-color: #1a1a1a;
      color: #ffffff;
      font-weight: 500;
    }
  }
}

.loading-state {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .skeleton-card {
    height: 200px;
    background: linear-gradient(90deg, #f5f5f5 25%, #e8e8e8 50%, #f5f5f5 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
    border-radius: 8px;
  }
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.empty-state {
  text-align: center;
  padding: 80px 40px;
  background: #fafafa;
  border-radius: 12px;

  svg {
    color: #ccc;
    margin-bottom: 24px;
  }

  h2 {
    font-size: 22px;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0 0 8px 0;
  }

  p {
    font-size: 15px;
    color: #999;
    margin: 0 0 24px 0;
  }

  .apply-btn-large {
    display: inline-block;
    padding: 14px 32px;
    background: #1a1a1a;
    color: #ffffff;
    border: none;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

    &:hover {
      background: #333;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }

    &:active {
      transform: translateY(0);
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
  }
}

.applications-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.application-card {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s;

  &:hover {
    border-color: #d1d5db;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #f2f2f2;
    background: #fafafa;

    .status-badge {
      padding: 4px 10px;
      border-radius: 6px;
      font-size: 13px;
      font-weight: 500;

      &.status-pending {
        background: #fef3c7;
        color: #d97706;
      }

      &.status-approved {
        background: #dcfce7;
        color: #16a34a;
      }

      &.status-rejected {
        background: #fee2e2;
        color: #dc2626;
      }

      &.status-cancelled {
        background: #f3f4f6;
        color: #6b7280;
      }
    }

    .cancel-btn {
      padding: 6px 14px;
      background: none;
      border: 1px solid #dc2626;
      color: #dc2626;
      border-radius: 6px;
      font-size: 13px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.15s;

      &:hover {
        background: #dc2626;
        color: #ffffff;
      }
    }
  }

  .card-body {
    padding: 20px;

    .application-reason {
      margin-bottom: 20px;

      label {
        display: block;
        font-size: 13px;
        font-weight: 600;
        color: #666;
        margin-bottom: 8px;
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }

      p {
        font-size: 15px;
        line-height: 1.6;
        color: #1a1a1a;
        margin: 0;
      }
    }

    .application-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 20px;
      margin-bottom: 16px;
      padding-bottom: 16px;
      border-bottom: 1px solid #f2f2f2;

      .meta-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 13px;
        color: #666;

        svg {
          width: 16px;
          height: 16px;
          color: #999;
        }
      }
    }

    .application-comment {
      label {
        display: block;
        font-size: 13px;
        font-weight: 600;
        color: #666;
        margin-bottom: 8px;
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }

      .comment-content {
        padding: 12px 16px;
        background: #fafafa;
        border-radius: 6px;
        font-size: 14px;
        line-height: 1.6;
        color: #1a1a1a;
      }
    }
  }
}

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
    color: #666;

    &:hover {
      color: #1a1a1a;
    }

    &.active {
      color: #1a1a1a;
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
    align-items: flex-start !important;
    gap: 12px;

    .page-title {
      font-size: 26px;
    }

    .apply-btn {
      width: 100%;
      justify-content: center;
    }
  }

  .filters-section {
    .filter-btn {
      flex: 1;
      min-width: calc(50% - 4px);
      justify-content: center;
    }
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start !important;
    gap: 12px;
  }

  .application-meta {
    flex-direction: column;
    gap: 8px !important;
  }
}
</style>
