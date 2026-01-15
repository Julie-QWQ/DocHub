import { defineStore } from 'pinia'
import { ref } from 'vue'
import { materialCategoryApi } from '@/api/materialCategory'
import type { MaterialCategoryConfig, MaterialCategoryConfigRequest } from '@/types'

export const useMaterialCategoryStore = defineStore('materialCategory', () => {
  // ==================== 状态 ====================
  const categories = ref<MaterialCategoryConfig[]>([])
  const activeCategories = ref<MaterialCategoryConfig[]>([])
  const loading = ref(false)

  // ==================== 方法 ====================

  /**
   * 获取所有资料类型（包括禁用的）
   */
  const fetchAllCategories = async () => {
    loading.value = true
    try {
      const response = await materialCategoryApi.list(false)

      if (response.code === 0 && response.data) {
        categories.value = response.data
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取启用的资料类型
   */
  const fetchActiveCategories = async () => {
    loading.value = true
    try {
      const response = await materialCategoryApi.list(true)

      if (response.code === 0 && response.data) {
        activeCategories.value = response.data
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 创建资料类型
   */
  const createCategory = async (data: MaterialCategoryConfigRequest) => {
    const response = await materialCategoryApi.create(data)

    if (response.code === 0 && response.data) {
      // 添加到列表
      categories.value.push(response.data)
      if (response.data.is_active) {
        activeCategories.value.push(response.data)
      }
    }

    return response
  }

  /**
   * 更新资料类型
   */
  const updateCategory = async (id: number, data: MaterialCategoryConfigRequest) => {
    const response = await materialCategoryApi.update(id, data)

    if (response.code === 0 && response.data) {
      // 更新列表中的数据
      const index = categories.value.findIndex(c => c.id === id)
      if (index !== -1) {
        categories.value[index] = response.data
      }

      // 更新活跃列表
      const activeIndex = activeCategories.value.findIndex(c => c.id === id)
      if (response.data.is_active) {
        if (activeIndex !== -1) {
          activeCategories.value[activeIndex] = response.data
        } else {
          activeCategories.value.push(response.data)
        }
      } else if (activeIndex !== -1) {
        activeCategories.value.splice(activeIndex, 1)
      }
    }

    return response
  }

  /**
   * 删除资料类型
   */
  const deleteCategory = async (id: number) => {
    const response = await materialCategoryApi.delete(id)

    if (response.code === 0) {
      // 从列表中移除
      categories.value = categories.value.filter(c => c.id !== id)
      activeCategories.value = activeCategories.value.filter(c => c.id !== id)
    }

    return response
  }

  /**
   * 切换启用状态
   */
  const toggleStatus = async (id: number) => {
    const response = await materialCategoryApi.toggleStatus(id)

    if (response.code === 0 && response.data) {
      // 更新列表中的数据
      const index = categories.value.findIndex(c => c.id === id)
      if (index !== -1) {
        categories.value[index] = response.data
      }

      // 更新活跃列表
      const activeIndex = activeCategories.value.findIndex(c => c.id === id)
      if (response.data.is_active) {
        if (activeIndex === -1) {
          activeCategories.value.push(response.data)
        }
      } else if (activeIndex !== -1) {
        activeCategories.value.splice(activeIndex, 1)
      }
    }

    return response
  }

  /**
   * 根据代码获取资料类型
   */
  const getCategoryByCode = (code: string): MaterialCategoryConfig | undefined => {
    return activeCategories.value.find(c => c.code === code)
  }

  /**
   * 根据代码获取资料类型名称
   */
  const getCategoryName = (code: string): string => {
    const category = getCategoryByCode(code)
    return category?.name || code
  }

  /**
   * 重置状态
   */
  const reset = () => {
    categories.value = []
    activeCategories.value = []
    loading.value = false
  }

  return {
    // 状态
    categories,
    activeCategories,
    loading,

    // 方法
    fetchAllCategories,
    fetchActiveCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    toggleStatus,
    getCategoryByCode,
    getCategoryName,
    reset
  }
})
