import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getPublicSystemConfigs } from '@/api/system'

export const useSystemStore = defineStore('system', () => {
  // 系统配置
  const configs = ref<Record<string, string>>({
    site_name: 'Study-UPC',
    site_description: '学院学习资料托管平台',
    maintenance_mode: 'false'
  })

  // 是否已加载
  const loaded = ref(false)

  /**
   * 加载系统配置
   */
  async function loadConfigs() {
    if (loaded.value) return

    try {
      const { data } = await getPublicSystemConfigs([
        'site_name',
        'site_description',
        'maintenance_mode'
      ])
      configs.value = { ...configs.value, ...data }
      loaded.value = true

      // 更新页面标题
      updatePageTitle()
    } catch (error) {
      console.error('加载系统配置失败:', error)
    }
  }

  /**
   * 更新页面标题
   */
  function updatePageTitle() {
    const siteName = configs.value.site_name || 'Study-UPC'
    document.title = `${siteName} - 学习资料托管平台`
  }

  /**
   * 获取配置值
   */
  function getConfig(key: string, defaultValue = '') {
    return configs.value[key] || defaultValue
  }

  /**
   * 设置配置值（用于更新后立即生效）
   */
  function setConfig(key: string, value: string) {
    configs.value[key] = value
    if (key === 'site_name') {
      updatePageTitle()
    }
  }

  return {
    configs,
    loaded,
    loadConfigs,
    getConfig,
    setConfig
  }
})
