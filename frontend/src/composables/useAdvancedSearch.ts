import { ref, computed } from 'vue'
import { useSearchStore } from '@/stores/search'
import type { SearchRequest, MaterialCategory } from '@/api/search'
import { ElMessage } from 'element-plus'

/**
 * 高级搜索功能 composable
 */
export function useAdvancedSearch() {
  const searchStore = useSearchStore()

  const keyword = ref('')
  const category = ref<MaterialCategory | undefined>()
  const courseName = ref('')
  const selectedTags = ref<string[]>([])
  const dateRange = ref<[Date, Date] | null>(null)
  const sortBy = ref('created_at')
  const sortOrder = ref<'desc' | 'asc'>('desc')
  const currentPage = ref(1)
  const pageSize = ref(20)

  const loading = computed(() => searchStore.loading)
  const results = computed(() => searchStore.searchResults)
  const total = computed(() => searchStore.searchTotal)
  const totalPages = computed(() => searchStore.searchTotalPages)

  /**
   * 执行搜索
   */
  const performSearch = async () => {
    const params: SearchRequest = {
      page: currentPage.value,
      page_size: pageSize.value,
      sort_by: sortBy.value,
      sort_order: sortOrder.value
    }

    if (keyword.value.trim()) {
      params.keyword = keyword.value.trim()
    }

    if (category.value) {
      params.category = category.value
    }

    if (courseName.value.trim()) {
      params.course_name = courseName.value.trim()
    }

    if (selectedTags.value.length > 0) {
      params.tags = selectedTags.value
    }

    if (dateRange.value) {
      params.start_date = dateRange.value[0].toISOString().split('T')[0]
      params.end_date = dateRange.value[1].toISOString().split('T')[0]
    }

    try {
      await searchStore.search(params)
    } catch (error) {
      ElMessage.error('搜索失败,请稍后重试')
    }
  }

  /**
   * 重置搜索条件
   */
  const resetFilters = () => {
    keyword.value = ''
    category.value = undefined
    courseName.value = ''
    selectedTags.value = []
    dateRange.value = null
    sortBy.value = 'created_at'
    sortOrder.value = 'desc'
    currentPage.value = 1
  }

  /**
   * 处理页码变化
   */
  const handlePageChange = (page: number) => {
    currentPage.value = page
    performSearch()
  }

  /**
   * 处理每页数量变化
   */
  const handlePageSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    performSearch()
  }

  /**
   * 处理排序变化
   */
  const handleSortChange = (field: string) => {
    if (sortBy.value === field) {
      sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
    } else {
      sortBy.value = field
      sortOrder.value = 'desc'
    }
    performSearch()
  }

  /**
   * 快速搜索
   */
  const quickSearch = (searchKeyword: string) => {
    keyword.value = searchKeyword
    resetFilters()
    performSearch()
  }

  return {
    // 状态
    keyword,
    category,
    courseName,
    selectedTags,
    dateRange,
    sortBy,
    sortOrder,
    currentPage,
    pageSize,

    // 计算属性
    loading,
    results,
    total,
    totalPages,

    // 方法
    performSearch,
    resetFilters,
    handlePageChange,
    handlePageSizeChange,
    handleSortChange,
    quickSearch
  }
}
