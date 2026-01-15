<template>
  <div class="applications-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>学委申请审核</h1>
      <div class="header-actions">
        <span v-if="pendingCount > 0" class="pending-badge">{{ pendingCount }} 待审核</span>
        <button class="refresh-btn" @click="handleRefresh">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
            <path d="M3 3v5h5"></path>
            <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"></path>
            <path d="M16 21h5v-5"></path>
          </svg>
          刷新
        </button>
      </div>
    </div>

    <!-- 申请列表 -->
    <div v-loading="loading" class="applications-list">
      <div v-for="app in applications" :key="app.id" class="application-card">
        <div class="applicant-avatar">
          <div class="avatar-placeholder">
            {{ (app.user?.real_name || app.user?.username || '?').charAt(0).toUpperCase() }}
          </div>
        </div>

        <div class="application-info">
          <div class="applicant-name">
            {{ app.user?.real_name || app.user?.username }}
          </div>
          <div class="application-reason">{{ app.reason }}</div>
          <div class="application-meta">
            <span>{{ format.datetime(app.created_at) }}</span>
            <span class="separator">·</span>
            <span v-if="app.reviewer">
              审核人: {{ app.reviewer.real_name || app.reviewer.username }}
            </span>
          </div>
        </div>

        <div class="application-status">
          <span :class="['status-badge', app.status]">
            {{ getStatusText(app.status) }}
          </span>
        </div>

        <div class="application-actions">
          <button
            v-if="app.status === 'pending'"
            class="action-btn primary"
            @click="handleReview(app)"
          >
            审核
          </button>
          <button class="action-btn" @click="handleViewDetail(app)">
            查看详情
          </button>
        </div>
      </div>

      <div v-if="!loading && applications.length === 0" class="empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
          <circle cx="8.5" cy="7" r="4"></circle>
          <line x1="20" y1="8" x2="20" y2="14"></line>
          <line x1="23" y1="11" x2="17" y2="11"></line>
        </svg>
        <p>暂无申请记录</p>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > 0" class="pagination">
      <span class="pagination-info">
        显示 {{ (page - 1) * size + 1 }}-{{ Math.min(page * size, total) }} / 共 {{ total }} 条
      </span>
      <div class="pagination-controls">
        <button class="page-btn" :disabled="page === 1" @click="handlePageChange(page - 1)">
          上一页
        </button>
        <button class="page-btn" :disabled="page * size >= total" @click="handlePageChange(page + 1)">
          下一页
        </button>
      </div>
    </div>

    <!-- 审核对话框 -->
    <ReviewApplication
      v-model="reviewDialogVisible"
      :application="currentApplication"
      @success="handleReviewSuccess"
    />

    <!-- 详情对话框 -->
    <div v-if="detailDialogVisible" class="dialog-overlay" @click.self="detailDialogVisible = false">
      <div class="dialog">
        <div class="dialog-header">
          <h3>申请详情</h3>
          <button class="close-btn" @click="detailDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div v-if="currentApplication" class="dialog-body">
          <div class="detail-group">
            <label>申请人</label>
            <div class="detail-value">
              {{ currentApplication.user?.real_name || currentApplication.user?.username }}
            </div>
          </div>

          <div class="detail-group">
            <label>申请理由</label>
            <div class="detail-value">{{ currentApplication.reason }}</div>
          </div>

          <div class="detail-group">
            <label>申请状态</label>
            <span :class="['status-badge', currentApplication.status]">
              {{ getStatusText(currentApplication.status) }}
            </span>
          </div>

          <div class="detail-row">
            <div class="detail-group">
              <label>申请时间</label>
              <div class="detail-value">{{ format.datetime(currentApplication.created_at) }}</div>
            </div>
            <div v-if="currentApplication.reviewed_at" class="detail-group">
              <label>审核时间</label>
              <div class="detail-value">{{ format.datetime(currentApplication.reviewed_at) }}</div>
            </div>
          </div>

          <div v-if="currentApplication.reviewer" class="detail-group">
            <label>审核人</label>
            <div class="detail-value">
              {{ currentApplication.reviewer.real_name || currentApplication.reviewer.username }}
            </div>
          </div>

          <div v-if="currentApplication.review_comment" class="detail-group">
            <label>审核意见</label>
            <div class="detail-value">{{ currentApplication.review_comment }}</div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="detailDialogVisible = false">关闭</button>
          <button
            v-if="currentApplication?.status === 'pending'"
            class="btn btn-primary"
            @click="handleReviewFromDetail"
          >
            立即审核
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useCommitteeStore } from '@/stores/committee'
import { format } from '@/utils/format'
import ReviewApplication from '@/components/ReviewApplication.vue'
import type { CommitteeApplication } from '@/types'

const committeeStore = useCommitteeStore()

const loading = ref(false)
const page = ref(1)
const size = ref(20)
const pendingCount = ref(0)

const reviewDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentApplication = ref<CommitteeApplication | null>(null)

const applications = computed(() => committeeStore.applications)
const total = computed(() => committeeStore.total)

const loadApplications = async () => {
  loading.value = true
  try {
    await committeeStore.fetchAllApplications({
      page: page.value,
      size: size.value
    })

    const pendingResponse = await committeeStore.fetchPendingCount()
    if (pendingResponse.code === 0) {
      pendingCount.value = pendingResponse.data.count
    }
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  loadApplications()
}

const handlePageChange = (newPage: number) => {
  page.value = newPage
  loadApplications()
}

const handleReview = (application: CommitteeApplication) => {
  currentApplication.value = application
  reviewDialogVisible.value = true
}

const handleViewDetail = (application: CommitteeApplication) => {
  currentApplication.value = application
  detailDialogVisible.value = true
}

const handleReviewFromDetail = () => {
  detailDialogVisible.value = false
  reviewDialogVisible.value = true
}

const handleReviewSuccess = () => {
  loadApplications()
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    cancelled: 'info'
  }
  return typeMap[status] || ''
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

onMounted(() => {
  loadApplications()
})
</script>

<style scoped lang="scss">
.applications-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h1 {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.pending-badge {
  padding: 6px 12px;
  background: #fef3c7;
  color: #b45309;
  font-size: 13px;
  font-weight: 600;
  border-radius: 6px;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #111827;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }
}

.applications-list {
  min-height: 300px;
}

.application-card {
  display: grid;
  grid-template-columns: auto 1fr auto auto;
  gap: 16px;
  align-items: center;
  padding: 20px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  margin-bottom: 12px;
  transition: all 0.2s;

  &:hover {
    border-color: #d1d5db;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  }

  @media (max-width: 768px) {
    grid-template-columns: auto 1fr;
    gap: 12px;

    .application-status {
      display: none;
    }
  }
}

.applicant-avatar {
  .avatar-placeholder {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: #f3f4f6;
    color: #6b7280;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    font-weight: 600;
  }
}

.application-info {
  min-width: 0;
}

.applicant-name {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 6px;
}

.application-reason {
  font-size: 14px;
  color: #374151;
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5;
  max-height: 3em;
}

.application-meta {
  font-size: 13px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  gap: 6px;

  .separator {
    color: #d1d5db;
  }
}

.status-badge {
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 6px;
  white-space: nowrap;

  &.pending {
    background: #fef3c7;
    color: #b45309;
  }

  &.approved {
    background: #dcfce7;
    color: #15803d;
  }

  &.rejected {
    background: #fee2e2;
    color: #b91c1c;
  }

  &.cancelled {
    background: #e5e7eb;
    color: #6b7280;
  }
}

.application-actions {
  display: flex;
  gap: 8px;

  @media (max-width: 768px) {
    grid-column: 2;
    justify-content: flex-end;
  }
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #9ca3af;

  svg {
    margin-bottom: 16px;
  }

  p {
    margin: 0;
    font-size: 14px;
  }
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.pagination-info {
  font-size: 13px;
  color: #6b7280;
}

.pagination-controls {
  display: flex;
  gap: 8px;
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: #ffffff;
  color: #111827;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    background: #f9fafb;
    border-color: #d1d5db;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

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
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;

  h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #111827;
  }

  .close-btn {
    background: none;
    border: none;
    padding: 4px;
    cursor: pointer;
    color: #6b7280;
    transition: color 0.2s;

    &:hover {
      color: #111827;
    }
  }
}

.dialog-body {
  padding: 20px;
}

.detail-group {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }

  label {
    display: block;
    font-size: 13px;
    font-weight: 600;
    color: #374151;
    margin-bottom: 8px;
  }
}

.detail-value {
  font-size: 14px;
  color: #111827;
}

.detail-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;

  @media (max-width: 600px) {
    grid-template-columns: 1fr;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &.btn-secondary {
    background: #f3f4f6;
    color: #111827;

    &:hover {
      background: #e5e7eb;
    }
  }

  &.btn-primary {
    background: #111827;
    color: #ffffff;

    &:hover {
      background: #000000;
    }
  }
}
</style>
