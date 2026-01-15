<template>
  <div class="auth-container">
    <div class="auth-logo">
      <h1><SiteName /></h1>
      <p class="subtitle">{{ siteDescription }}</p>
    </div>

    <div class="auth-card">
      <h2 class="auth-title">忘记密码</h2>
      <p class="auth-description">输入您的邮箱地址，我们将发送密码重置链接</p>

      <form @submit="handleSubmit" class="auth-form">
        <div class="form-group">
          <label for="email">邮箱地址</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="输入注册时使用的邮箱"
            :disabled="isSubmitting"
          />
        </div>

        <button type="submit" class="submit-button" :disabled="isSubmitting || !email">
          {{ isSubmitting ? '发送中...' : '发送重置链接' }}
        </button>

        <div class="auth-footer">
          记起密码了？
          <router-link to="/login">返回登录</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import SiteName from '@/components/SiteName.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学院学习资料托管平台'))

const email = ref('')
const isSubmitting = ref(false)

const handleSubmit = async (e: Event) => {
  e.preventDefault()

  if (!email.value) {
    ElMessage.warning('请输入邮箱地址')
    return
  }

  isSubmitting.value = true

  // 模拟发送请求
  setTimeout(() => {
    ElMessage.success('如果该邮箱已注册，您将收到密码重置链接')
    isSubmitting.value = false
    email.value = ''
  }, 1500)
}
</script>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';
</style>
