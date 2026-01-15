import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { committeeApi } from '@/api/committee'
import type {
  CommitteeApplication,
  CommitteeApplicationListParams,
  CreateCommitteeApplicationRequest,
  ReviewCommitteeApplicationRequest
} from '@/types'

export const useCommitteeStore = defineStore('committee', () => {
  // ==================== 状态 ====================
  const applications = ref<CommitteeApplication[]>([])
  const currentApplication = ref<CommitteeApplication | null>(null)
  const total = ref(0)
  const page = ref(1)
  const size = ref(20)
  const loading = ref(false)
  const pendingCount = ref(0)

  // 当前查询参数
  const currentParams = ref<CommitteeApplicationListParams>({
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
   * 申请成为学委
   */
  const applyForCommittee = async (data: CreateCommitteeApplicationRequest) => {
    const response = await committeeApi.applyForCommittee(data)
    if (response.code === 0) {
      // 刷新列表
      await fetchMyApplications()
    }
    return response
  }

  /**
   * 获取我的申请列表
   */
  const fetchMyApplications = async (params?: CommitteeApplicationListParams) => {
    loading.value = true
    try {
      const queryParams = {
        ...currentParams.value,
        ...params
      }

      const response = await committeeApi.getMyApplications(queryParams)

      if (response.code === 0 && response.data) {
        applications.value = response.data.list || []
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
   * 获取申请详情
   */
  const fetchApplicationDetail = async (id: number) => {
    loading.value = true
    try {
      const response = await committeeApi.getApplicationDetail(id)

      if (response.code === 0 && response.data) {
        currentApplication.value = response.data
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 取消申请
   */
  const cancelApplication = async (id: number) => {
    const response = await committeeApi.cancelApplication(id)

    if (response.code === 0) {
      // 从列表中移除
      applications.value = applications.value.filter(app => app.id !== id)

      // 清空当前申请
      if (currentApplication.value?.id === id) {
        currentApplication.value = null
      }
    }

    return response
  }

  /**
   * 管理员获取所有申请列表
   */
  const fetchAllApplications = async (params?: CommitteeApplicationListParams) => {
    loading.value = true
    try {
      const queryParams = {
        ...currentParams.value,
        ...params
      }

      const response = await committeeApi.getAllApplications(queryParams)

      if (response.code === 0 && response.data) {
        applications.value = response.data.list || []
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
   * 管理员审核学委申请
   */
  const reviewApplication = async (id: number, data: ReviewCommitteeApplicationRequest) => {
    const response = await committeeApi.reviewApplication(id, data)

    if (response.code === 0 && response.data) {
      // 更新列表中的数据
      const index = applications.value.findIndex(app => app.id === id)
      if (index !== -1) {
        applications.value[index] = response.data
      }

      // 更新当前申请
      if (currentApplication.value?.id === id) {
        currentApplication.value = response.data
      }

      // 更新待审核数量
      await fetchPendingCount()
    }

    return response
  }

  /**
   * 获取待审核数量
   */
  const fetchPendingCount = async () => {
    const response = await committeeApi.getPendingCount()

    if (response.code === 0 && response.data) {
      pendingCount.value = response.data.count
    }

    return response
  }

  /**
   * 重置状态
   */
  const reset = () => {
    applications.value = []
    currentApplication.value = null
    total.value = 0
    page.value = 1
    size.value = 20
    loading.value = false
    pendingCount.value = 0
    currentParams.value = {
      page: 1,
      size: 20
    }
  }

  return {
    // 状态
    applications,
    currentApplication,
    total,
    page,
    size,
    loading,
    pendingCount,
    currentParams,

    // 计算属性
    hasMore,
    totalPages,

    // 方法
    applyForCommittee,
    fetchMyApplications,
    fetchApplicationDetail,
    cancelApplication,
    fetchAllApplications,
    reviewApplication,
    fetchPendingCount,
    reset
  }
})
