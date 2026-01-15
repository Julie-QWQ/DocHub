import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { searchApi, type SearchRequest, type SearchResult, type SearchResponse, type HotKeyword, type SearchHistory, type RecommendationResult } from '@/api/search'

export const useSearchStore = defineStore('search', () => {
  // 状态
  const searchResults = ref<SearchResult[]>([])
  const searchTotal = ref(0)
  const searchPage = ref(1)
  const searchPageSize = ref(20)
  const searchTotalPages = ref(0)
  const currentKeyword = ref('')
  const loading = ref(false)

  const hotKeywords = ref<HotKeyword[]>([])
  const hotKeywordsLoading = ref(false)

  const searchHistory = ref<SearchHistory[]>([])
  const searchHistoryLoading = ref(false)

  const hotMaterials = ref<SearchResult[]>([])
  const hotMaterialsLoading = ref(false)

  const recommendations = ref<RecommendationResult[]>([])
  const recommendationsLoading = ref(false)

  // 计算属性
  const hasSearchResults = computed(() => searchResults.value.length > 0)
  const hasHotKeywords = computed(() => hotKeywords.value.length > 0)
  const hasSearchHistory = computed(() => searchHistory.value.length > 0)
  const hasHotMaterials = computed(() => hotMaterials.value.length > 0)
  const hasRecommendations = computed(() => recommendations.value.length > 0)

  // 方法
  const search = async (params: SearchRequest) => {
    loading.value = true
    try {
      const response = await searchApi.search(params)
      if (response.data.code === 0) {
        const data = response.data.data
        searchResults.value = data.results
        searchTotal.value = data.total
        searchPage.value = data.page
        searchPageSize.value = data.page_size
        searchTotalPages.value = data.total_pages
        if (params.keyword) {
          currentKeyword.value = params.keyword
        }
      }
      return response.data
    } catch (error) {
      console.error('搜索失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const clearSearch = () => {
    searchResults.value = []
    searchTotal.value = 0
    searchPage.value = 1
    searchTotalPages.value = 0
    currentKeyword.value = ''
  }

  const fetchHotKeywords = async (limit: number = 20) => {
    hotKeywordsLoading.value = true
    try {
      const response = await searchApi.getHotKeywords(limit)
      if (response.data.code === 0) {
        hotKeywords.value = response.data.data
      }
      return response.data
    } catch (error) {
      console.error('获取热门搜索词失败:', error)
      throw error
    } finally {
      hotKeywordsLoading.value = false
    }
  }

  const fetchSearchHistory = async (limit: number = 20) => {
    searchHistoryLoading.value = true
    try {
      const response = await searchApi.getSearchHistory(limit)
      if (response.data.code === 0) {
        searchHistory.value = response.data.data
      }
      return response.data
    } catch (error) {
      console.error('获取搜索历史失败:', error)
      throw error
    } finally {
      searchHistoryLoading.value = false
    }
  }

  const clearSearchHistoryApi = async () => {
    try {
      const response = await searchApi.clearSearchHistory()
      if (response.data.code === 0) {
        searchHistory.value = []
      }
      return response.data
    } catch (error) {
      console.error('清空搜索历史失败:', error)
      throw error
    }
  }

  const fetchHotMaterials = async (limit: number = 20) => {
    hotMaterialsLoading.value = true
    try {
      const response = await searchApi.getHotMaterials(limit)
      if (response.data.code === 0) {
        hotMaterials.value = response.data.data
      }
      return response.data
    } catch (error) {
      console.error('获取热门资料失败:', error)
      throw error
    } finally {
      hotMaterialsLoading.value = false
    }
  }

  const fetchRecommendations = async (params: { type?: string; material_id?: number; limit?: number }) => {
    recommendationsLoading.value = true
    try {
      const response = await searchApi.getRecommendations(params)
      if (response.data.code === 0) {
        recommendations.value = response.data.data
      }
      return response.data
    } catch (error) {
      console.error('获取推荐资料失败:', error)
      throw error
    } finally {
      recommendationsLoading.value = false
    }
  }

  return {
    // 状态
    searchResults,
    searchTotal,
    searchPage,
    searchPageSize,
    searchTotalPages,
    currentKeyword,
    loading,
    hotKeywords,
    hotKeywordsLoading,
    searchHistory,
    searchHistoryLoading,
    hotMaterials,
    hotMaterialsLoading,
    recommendations,
    recommendationsLoading,

    // 计算属性
    hasSearchResults,
    hasHotKeywords,
    hasSearchHistory,
    hasHotMaterials,
    hasRecommendations,

    // 方法
    search,
    clearSearch,
    fetchHotKeywords,
    fetchSearchHistory,
    clearSearchHistoryApi,
    fetchHotMaterials,
    fetchRecommendations
  }
})
