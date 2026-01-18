import { ref, type Ref } from 'vue'
import { ElMessage } from 'element-plus'

/**
 * 乐观更新组合式函数
 * 用于需要立即反馈 UI，后台异步调用 API 的场景
 *
 * @example
 * const { mutate: addFavorite, loading } = useOptimisticMutation({
 *   mutationFn: (id: number) => materialApi.addFavorite(id),
 *   onMutate: (id) => {
 *     const material = materialStore.materials.find(m => m.id === id)
 *     if (material) {
 *       material.is_favorited = true
 *       material.favorite_count++
 *     }
 *   },
 *   onRollback: (id) => {
 *     const material = materialStore.materials.find(m => m.id === id)
 *     if (material) {
 *       material.is_favorited = false
 *       material.favorite_count--
 *     }
 *   },
 *   onSuccessMessage: '收藏成功',
 *   onErrorMessage: '收藏失败，已回滚'
 * })
 */
export function useOptimisticMutation<T = any, V = any>(
  options: {
    /** API 调用函数 */
    mutationFn: (variables: V) => Promise<T>

    /** 乐观更新：立即更新 UI */
    onMutate: (variables: V) => void

    /** 回滚：API 失败时恢复 UI */
    onRollback: (variables: V) => void

    /** 成功提示消息 */
    onSuccessMessage?: string

    /** 失败提示消息 */
    onErrorMessage?: string

    /** 成功后的回调（可选） */
    onSuccess?: (data: T, variables: V) => void

    /** 失败后的回调（可选） */
    onError?: (error: Error, variables: V) => void

    /** 是否显示成功消息（默认 true） */
    showSuccessMessage?: boolean

    /** 是否显示失败消息（默认 true） */
    showErrorMessage?: boolean
  }
) {
  const loading: Ref<boolean> = ref(false)
  const error: Ref<Error | null> = ref(null)

  const mutate = async (variables: V): Promise<T | undefined> => {
    loading.value = true
    error.value = null

    // 1. 立即执行乐观更新
    options.onMutate(variables)

    // 显示成功消息
    if (options.showSuccessMessage !== false && options.onSuccessMessage) {
      ElMessage.success(options.onSuccessMessage)
    }

    try {
      // 2. 异步调用 API
      const data = await options.mutationFn(variables)

      // 3. 成功回调
      options.onSuccess?.(data, variables)

      return data
    } catch (err) {
      const errorObj = err as Error

      // 4. 回滚状态
      options.onRollback(variables)

      error.value = errorObj
      options.onError?.(errorObj, variables)

      // 显示失败消息
      if (options.showErrorMessage !== false) {
        ElMessage.error(options.onErrorMessage || '操作失败，已回滚')
      }

      throw errorObj
    } finally {
      loading.value = false
    }
  }

  return {
    mutate,
    loading,
    error
  }
}
