import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

/**
 * 导航状态管理 Store
 * 用于实现立即反馈的侧边栏高亮和全局加载状态
 */
export const useNavigationStore = defineStore('navigation', () => {
  // 当前激活的路由路径（用于立即更新侧边栏高亮）
  const currentPath = ref<string>('')

  // 当前路由是否正在加载
  const isNavigating = ref(false)

  // 页面加载状态映射
  const pageLoadingStates = ref<Record<string, boolean>>({})

  // 设置当前路径（立即调用，不等待路由切换）
  const setCurrentPath = (path: string) => {
    currentPath.value = path
  }

  // 标记页面开始加载
  const setPageLoading = (path: string, loading: boolean) => {
    pageLoadingStates.value[path] = loading
  }

  // 检查页面是否正在加载
  const isPageLoading = (path: string) => {
    return pageLoadingStates.value[path] || false
  }

  // 标记导航开始
  const startNavigation = () => {
    isNavigating.value = true
  }

  // 标记导航结束
  const endNavigation = () => {
    isNavigating.value = false
  }

  // 当前路由是否激活（用于侧边栏高亮）
  const isPathActive = (path: string) => {
    // 精确匹配
    if (currentPath.value === path) return true

    // 前缀匹配（用于管理后台路由）
    if (path.startsWith('/admin/') && currentPath.value.startsWith(path)) {
      return true
    }

    return false
  }

  return {
    currentPath,
    isNavigating,
    pageLoadingStates,
    setCurrentPath,
    setPageLoading,
    isPageLoading,
    startNavigation,
    endNavigation,
    isPathActive
  }
})
