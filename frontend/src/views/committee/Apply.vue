<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useCommitteeStore } from '@/stores/committee'
import { useAuthStore } from '@/stores/auth'
import { format } from '@/utils/format'

const router = useRouter()
const committeeStore = useCommitteeStore()
const authStore = useAuthStore()

const loading = ref(false)
const loadingHistory = ref(false)

const formData = reactive({
  reason: ''
})

const errors = ref<Record<string, string>>({})

// 验证表单
const validateForm = (): boolean => {
  errors.value = {}

  if (!formData.reason.trim()) {
    errors.value.reason = '请填写申请理由'
    return false
  }

  if (formData.reason.trim().length < 20) {
    errors.value.reason = '申请理由至少 20 个字'
    return false
  }

  if (formData.reason.length > 500) {
    errors.value.reason = '申请理由不能超过 500 个字'
    return false
  }

  return true
}

const recentApplications = computed(() => {
  return committeeStore.applications.slice(0, 3)
})

const hasApplications = computed(() => {
  return committeeStore.applications.length > 0
})

const handleSubmit = async () => {
  // 清除之前的错误
  errors.value = {}

  // 验证表单
  if (!validateForm()) {
    return
  }

  // 检查是否已有待审核申请
  const pendingApplication = committeeStore.applications.find(
    app => app.status === 'pending'
  )

  if (pendingApplication) {
    ElMessage.warning('您已有待审核的申请，请等待审核完成')
    return
  }

  loading.value = true
  try {
    await committeeStore.applyForCommittee({
      reason: formData.reason
    })

    ElMessage.success('申请提交成功，请等待管理员审核')
    formData.reason = ''
    await loadApplications()
  } catch (error: any) {
    ElMessage.error(error.message || '申请提交失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  formData.reason = ''
  errors.value = {}
}

const handleViewAll = () => {
  router.push('/committee/applications')
}

const loadApplications = async () => {
  loadingHistory.value = true
  try {
    await committeeStore.fetchMyApplications({ page: 1, size: 10 })
  } finally {
    loadingHistory.value = false
  }
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
  <div class="committee-apply-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题 -->
        <header class="stream-header">
          <h1 class="page-title">申请成为学习委员</h1>
        </header>

        <!-- 申请表单 -->
        <div class="apply-section">
          <div class="section-intro">
            <p>提交申请后，管理员将对您的申请进行审核</p>
          </div>

          <form @submit.prevent="handleSubmit" class="apply-form">
            <div class="form-group">
              <label class="form-label">
                申请理由 <span class="required">*</span>
              </label>
              <textarea
                v-model="formData.reason"
                class="form-textarea"
                :class="{ error: errors.reason }"
                placeholder="请详细说明您想成为学习委员的理由（至少 20 个字）"
                rows="10"
                maxlength="500"
              ></textarea>
              <span class="char-count">{{ formData.reason.length }}/500</span>
              <span v-if="errors.reason" class="error-message">{{ errors.reason }}</span>
              <div class="form-tip">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M12 16v-4"/>
                  <path d="M12 8h.01"/>
                </svg>
                <span>请认真填写申请理由，这将是管理员审核的重要依据</span>
              </div>
            </div>

            <div class="form-actions">
              <button type="submit" class="btn btn-primary" :disabled="loading">
                <span v-if="loading">提交中...</span>
                <span v-else>提交申请</span>
              </button>
              <button type="button" class="btn btn-secondary" @click="handleReset">
                重置
              </button>
            </div>
          </form>
        </div>

        <!-- 我的申请记录 -->
        <div v-if="hasApplications" class="history-section">
          <div class="section-header">
            <h2 class="section-title">我的申请记录</h2>
            <button class="view-all-btn" @click="handleViewAll">
              查看全部
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M5 12h14"/>
                <path d="m12 5 7 7-7 7"/>
              </svg>
            </button>
          </div>

          <div v-if="loadingHistory" class="loading-state">
            <div v-for="i in 3" :key="i" class="skeleton-item"></div>
          </div>

          <div v-else class="application-list">
            <div
              v-for="application in recentApplications"
              :key="application.id"
              class="application-item"
            >
              <div class="application-header">
                <span :class="['status-badge', getStatusClass(application.status)]">
                  {{ getStatusText(application.status) }}
                </span>
                <span class="application-time">
                  {{ format.datetime(application.created_at) }}
                </span>
              </div>

              <div class="application-content">
                <p class="application-reason">{{ application.reason }}</p>

                <div v-if="application.review_comment" class="application-comment">
                  <div class="comment-label">审核意见</div>
                  <div class="comment-text">{{ application.review_comment }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.committee-apply-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 800px;
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

.apply-section {
  margin-bottom: 48px;

  .section-intro {
    margin-bottom: 24px;

    p {
      font-size: 15px;
      color: #666;
      margin: 0;
    }
  }

  .apply-form {
    .form-group {
      margin-bottom: 24px;

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

      .form-textarea {
        width: 100%;
        padding: 16px;
        background: #ffffff;
        border: 1px solid #e5e5e5;
        border-radius: 8px;
        font-size: 15px;
        color: #1a1a1a;
        outline: none;
        transition: all 0.2s;
        font-family: inherit;
        line-height: 1.6;
        resize: vertical;

        &::placeholder {
          color: #999;
        }

        &:focus {
          border-color: #1a1a1a;
          box-shadow: 0 0 0 3px rgba(26, 26, 26, 0.1);
        }

        &.error {
          border-color: #dc2626;

          &:focus {
            box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.1);
          }
        }
      }

      .char-count {
        display: block;
        text-align: right;
        font-size: 12px;
        color: #999;
        margin-top: 8px;
      }

      .error-message {
        display: block;
        font-size: 13px;
        color: #dc2626;
        margin-top: 6px;
      }

      .form-tip {
        display: flex;
        align-items: center;
        gap: 6px;
        margin-top: 12px;
        font-size: 13px;
        color: #666;

        svg {
          color: #FF6B35;
          flex-shrink: 0;
        }
      }
    }

    .form-actions {
      display: flex;
      gap: 12px;
      padding-top: 8px;

      .btn {
        padding: 12px 24px;
        font-size: 15px;
        font-weight: 500;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s;
        border: none;

        &:disabled {
          opacity: 0.5;
          cursor: not-allowed;
        }

        &-primary {
          background: #1a1a1a;
          color: #ffffff;

          &:hover:not(:disabled) {
            background: #333;
          }
        }

        &-secondary {
          background: #ffffff;
          color: #1a1a1a;
          border: 1px solid #e5e5e5;

          &:hover:not(:disabled) {
            background: #fafafa;
            border-color: #1a1a1a;
          }
        }
      }
    }
  }
}

.history-section {
  padding-top: 32px;
  border-top: 1px solid #f2f2f2;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    .section-title {
      font-size: 24px;
      font-weight: 600;
      color: #1a1a1a;
      margin: 0;
    }

    .view-all-btn {
      display: flex;
      align-items: center;
      gap: 4px;
      background: none;
      border: none;
      color: #666;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: color 0.15s;

      &:hover {
        color: #1a1a1a;
      }

      svg {
        width: 16px;
        height: 16px;
      }
    }
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    gap: 16px;

    .skeleton-item {
      height: 120px;
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

  .application-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .application-item {
    padding: 20px;
    background: #fafafa;
    border: 1px solid #e5e5e5;
    border-radius: 8px;
    transition: all 0.2s;

    &:hover {
      background: #ffffff;
      border-color: #1a1a1a;
    }

    .application-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .status-badge {
        padding: 4px 12px;
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

      .application-time {
        font-size: 13px;
        color: #999;
      }
    }

    .application-content {
      .application-reason {
        font-size: 14px;
        line-height: 1.6;
        color: #1a1a1a;
        margin: 0 0 12px 0;
      }

      .application-comment {
        padding: 12px;
        background: #ffffff;
        border-radius: 6px;

        .comment-label {
          font-size: 13px;
          font-weight: 600;
          color: #1a1a1a;
          margin-bottom: 6px;
        }

        .comment-text {
          font-size: 14px;
          line-height: 1.6;
          color: #666;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    .page-title {
      font-size: 26px;
    }
  }

  .form-actions {
    flex-direction: column;

    .btn {
      width: 100%;
    }
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
