<template>
  <div class="auth-container">
    <div class="auth-logo">
      <h1><SiteName /></h1>
      <p class="subtitle">{{ siteDescription }}</p>
    </div>

    <div class="auth-card">
      <h2 class="auth-title">注册</h2>
      <p class="auth-description">创建您的账号，开始使用学习资料平台</p>

      <RegisterForm @success="handleRegisterSuccess" />

      <div class="auth-footer">
        已有账号？
        <router-link to="/login">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import RegisterForm from '@/components/RegisterForm.vue'
import SiteName from '@/components/SiteName.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学院学习资料托管平台'))

const handleRegisterSuccess = () => {
  // 注册成功后由 composable 处理跳转
}
</script>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';
</style>
