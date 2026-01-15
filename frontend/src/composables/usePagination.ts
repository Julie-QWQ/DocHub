import { ref, computed } from 'vue'
import type { PaginationParams } from '@/types'

/**
 * 分页逻辑组合式函数
 */
export function usePagination(initialPageSize = 20) {
  const page = ref(1)
  const size = ref(initialPageSize)
  const total = ref(0)

  /**
   * 总页数
   */
  const totalPages = computed(() => {
    return Math.ceil(total.value / size.value) || 1
  })

  /**
   * 是否有更多数据
   */
  const hasMore = computed(() => {
    return page.value * size.value < total.value
  })

  /**
   * 是否是第一页
   */
  const isFirstPage = computed(() => {
    return page.value === 1
  })

  /**
   * 是否是最后一页
   */
  const isLastPage = computed(() => {
    return page.value >= totalPages.value
  })

  /**
   * 重置分页
   */
  const reset = () => {
    page.value = 1
    total.value = 0
  }

  /**
   * 上一页
   */
  const prevPage = () => {
    if (!isFirstPage.value) {
      page.value--
    }
  }

  /**
   * 下一页
   */
  const nextPage = () => {
    if (!isLastPage.value) {
      page.value++
    }
  }

  /**
   * 跳转到指定页
   */
  const goToPage = (targetPage: number) => {
    if (targetPage >= 1 && targetPage <= totalPages.value) {
      page.value = targetPage
    }
  }

  /**
   * 改变每页大小
   */
  const changePageSize = (newSize: number) => {
    size.value = newSize
    page.value = 1 // 重置到第一页
  }

  /**
   * 设置总数
   */
  const setTotal = (newTotal: number) => {
    total.value = newTotal
  }

  /**
   * 获取查询参数
   */
  const getParams = (): PaginationParams => {
    return {
      page: page.value,
      size: size.value
    }
  }

  return {
    // 状态
    page,
    size,
    total,

    // 计算属性
    totalPages,
    hasMore,
    isFirstPage,
    isLastPage,

    // 方法
    reset,
    prevPage,
    nextPage,
    goToPage,
    changePageSize,
    setTotal,
    getParams
  }
}
