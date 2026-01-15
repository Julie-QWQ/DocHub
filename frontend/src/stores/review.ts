import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { reviewApi } from '@/api/review'
import type {
  Material,
  ReviewRecord,
  ReviewHistoryParams,
  HandleReportRequest
} from '@/types'

export const useReviewStore = defineStore('review', () => {
  // ==================== 状态 ====================
  const pendingMaterials = ref<Material[]>([])
  const reviewHistory = ref<ReviewRecord[]>([])
  const currentStatistics = ref<{
    total_reviews: number
    approved_count: number
    rejected_count: number
    material_count: number
    committee_count: number
    report_count: number
  } | null>(null)

  const total = ref(0)
  const page = ref(1)
  const size = ref(20)
  const loading = ref(false)

  // 当前查询参数
  const currentParams = ref<ReviewHistoryParams>({
    page: 1,
    size: 20
  })

  // ==================== 计算属性 ====================
  const hasMore = computed(() => {
    return page.value * size.value < total.value
  })

  const totalPages = computed(() => {
    return Math.ceil(total.value / size.value)
  })

  // ==================== 方法 ====================

  /**
   * 获取待审核资料列表
   */
  const fetchPendingMaterials = async (params?: { page?: number; size?: number }) => {
    loading.value = true
    try {
      const response = await reviewApi.getPendingMaterials(params || { page: 1, size: 20 })

      if (response.code === 0 && response.data) {
        pendingMaterials.value = response.data.materials || response.data.list || []
        total.value = response.data.total
        page.value = response.data.page
        size.value = response.data.page_size || response.data.size || 20
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 审核资料
   */
  const reviewMaterial = async (id: number, approved: boolean, comment?: string) => {
    const response = await reviewApi.reviewMaterial(id, approved, comment)

    if (response.code === 0) {
      // 从待审核列表中移除
      pendingMaterials.value = pendingMaterials.value.filter(m => m.id !== id)
    }

    return response
  }

  /**
   * 处理举报
   */
  const handleReport = async (id: number, data: HandleReportRequest) => {
    return await reviewApi.handleReport(id, data)
  }

  /**
   * 获取审核历史
   */
  const fetchReviewHistory = async (params?: ReviewHistoryParams) => {
    loading.value = true
    try {
      const queryParams = {
        ...currentParams.value,
        ...params
      }

      const response = await reviewApi.getReviewHistory(queryParams)

      if (response.code === 0 && response.data) {
        reviewHistory.value = response.data.list || []
        total.value = response.data.total
        page.value = response.data.page
        size.value = response.data.size
        currentParams.value = queryParams
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取审核人统计
   */
  const fetchReviewerStatistics = async (reviewerId: number) => {
    const response = await reviewApi.getReviewerStatistics(reviewerId)

    if (response.code === 0 && response.data) {
      currentStatistics.value = response.data
    }

    return response
  }

  /**
   * 重置状态
   */
  const reset = () => {
    pendingMaterials.value = []
    reviewHistory.value = []
    currentStatistics.value = null
    total.value = 0
    page.value = 1
    size.value = 20
    loading.value = false
    currentParams.value = {
      page: 1,
      size: 20
    }
  }

  return {
    // 状态
    pendingMaterials,
    reviewHistory,
    currentStatistics,
    total,
    page,
    size,
    loading,
    currentParams,

    // 计算属性
    hasMore,
    totalPages,

    // 方法
    fetchPendingMaterials,
    reviewMaterial,
    handleReport,
    fetchReviewHistory,
    fetchReviewerStatistics,
    reset
  }
})
