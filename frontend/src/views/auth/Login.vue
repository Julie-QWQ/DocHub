<template>
  <div class="auth-container">
    <div class="auth-logo">
      <h1><SiteName /></h1>
      <p class="subtitle">{{ siteDescription }}</p>
    </div>

    <div class="auth-card">
      <h2 class="auth-title">登录</h2>
      <p class="auth-description">欢迎回来！请登录您的账号</p>

      <LoginForm @success="handleLoginSuccess" />

      <div class="auth-footer">
        还没有账号？
        <router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import LoginForm from '@/components/LoginForm.vue'
import SiteName from '@/components/SiteName.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学院学习资料托管平台'))

const handleLoginSuccess = () => {
  // 登录成功后由 composable 处理跳转
}
</script>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';
</style>
