<template>
  <div class="auth-layout">
    <div class="auth-container">
      <div class="auth-header">
        <h1><SiteName /></h1>
        <p>{{ siteDescription }}</p>
      </div>
      <div class="auth-content">
        <router-view />
      </div>
      <div class="auth-footer">
        <p>&copy; {{ currentYear }} <SiteName />. All rights reserved.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SiteName from '@/components/SiteName.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学习资料托管平台'))
const currentYear = computed(() => new Date().getFullYear())
</script>

<style scoped lang="scss">
.auth-layout {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.auth-container {
  width: 400px;
  padding: 40px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.auth-header {
  text-align: center;
  margin-bottom: 30px;

  h1 {
    font-size: 28px;
    color: #303133;
    margin-bottom: 8px;
  }

  p {
    font-size: 14px;
    color: #909399;
  }
}

.auth-content {
  margin-bottom: 20px;
}

.auth-footer {
  text-align: center;

  p {
    font-size: 12px;
    color: #c0c4cc;
  }
}
</style>
