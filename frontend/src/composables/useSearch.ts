import { ref, watch } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import type { MaterialListParams } from '@/types'

/**
 * 搜索相关的组合式函数
 */
export function useSearch(initialKeyword = '') {
  const keyword = ref(initialKeyword)
  const searching = ref(false)
  const searchHistory = ref<string[]>([])

  // 从 localStorage 加载搜索历史
  const loadSearchHistory = () => {
    try {
      const history = localStorage.getItem('search_history')
      if (history) {
        searchHistory.value = JSON.parse(history)
      }
    } catch (error) {
      console.error('加载搜索历史失败:', error)
    }
  }

  // 保存搜索历史到 localStorage
  const saveSearchHistory = () => {
    try {
      localStorage.setItem('search_history', JSON.stringify(searchHistory.value))
    } catch (error) {
      console.error('保存搜索历史失败:', error)
    }
  }

  /**
   * 添加搜索历史
   */
  const addSearchHistory = (word: string) => {
    if (!word || word.trim() === '') return

    // 移除已存在的相同关键词
    const index = searchHistory.value.indexOf(word)
    if (index !== -1) {
      searchHistory.value.splice(index, 1)
    }

    // 添加到开头
    searchHistory.value.unshift(word)

    // 限制历史记录数量（最多 20 条）
    if (searchHistory.value.length > 20) {
      searchHistory.value = searchHistory.value.slice(0, 20)
    }

    saveSearchHistory()
  }

  /**
   * 清除搜索历史
   */
  const clearSearchHistory = () => {
    searchHistory.value = []
    saveSearchHistory()
  }

  /**
   * 删除指定的搜索历史
   */
  const removeSearchHistory = (word: string) => {
    const index = searchHistory.value.indexOf(word)
    if (index !== -1) {
      searchHistory.value.splice(index, 1)
      saveSearchHistory()
    }
  }

  /**
   * 执行搜索
   */
  const performSearch = async (
    searchFn: (params: MaterialListParams) => Promise<void>,
    params?: Partial<MaterialListParams>
  ) => {
    if (keyword.value.trim() === '') {
      return
    }

    searching.value = true
    try {
      await searchFn({
        keyword: keyword.value.trim(),
        ...params
      })

      // 添加到搜索历史
      addSearchHistory(keyword.value.trim())
    } finally {
      searching.value = false
    }
  }

  /**
   * 防抖搜索
   */
  const debouncedSearch = useDebounceFn(async (
    searchFn: (params: MaterialListParams) => Promise<void>,
    params?: Partial<MaterialListParams>
  ) => {
    await performSearch(searchFn, params)
  }, 500)

  /**
   * 重置搜索
   */
  const resetSearch = async (resetFn: () => Promise<void>) => {
    keyword.value = ''
    await resetFn()
  }

  /**
   * 快速搜索（从历史记录）
   */
  const quickSearch = (word: string) => {
    keyword.value = word
  }

  /**
   * 获取热门搜索关键词（示例数据）
   */
  const getHotKeywords = (): string[] => {
    return ['高等数学', '线性代数', 'C语言', '数据结构', '计算机网络', '操作系统']
  }

  // 监听关键词变化
  watch(keyword, (newKeyword) => {
    if (newKeyword.trim() !== '') {
      // 可以在这里触发自动搜索
    }
  })

  // 初始化时加载搜索历史
  loadSearchHistory()

  return {
    // 状态
    keyword,
    searching,
    searchHistory,

    // 方法
    performSearch,
    debouncedSearch,
    resetSearch,
    quickSearch,
    addSearchHistory,
    removeSearchHistory,
    clearSearchHistory,
    getHotKeywords
  }
}
