<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useLoginForm } from '@/composables/useValidator'
import type { LoginRequest } from '@/types'

const emit = defineEmits<{
  success: []
}>()

const { login } = useAuth()
const router = useRouter()

const {
  username,
  password,
  errors,
  isSubmitting,
  handleSubmit
} = useLoginForm()

const onSubmit = handleSubmit.withControlled(async () => {
  const credentials: LoginRequest = {
    username: username.value,
    password: password.value
  }

  const success = await login(credentials)
  if (success) {
    emit('success')
  }
})
</script>

<template>
  <form @submit="onSubmit" class="auth-form">
    <div class="form-group">
      <label for="username">用户名或邮箱</label>
      <input
        id="username"
        v-model="username"
        type="text"
        placeholder="输入用户名或邮箱"
        autocomplete="username"
        :disabled="isSubmitting"
      />
      <p v-if="errors.username" class="error-message">{{ errors.username }}</p>
    </div>

    <div class="form-group">
      <label for="password">密码</label>
      <input
        id="password"
        v-model="password"
        type="password"
        placeholder="输入密码"
        autocomplete="current-password"
        :disabled="isSubmitting"
      />
      <p v-if="errors.password" class="error-message">{{ errors.password }}</p>
    </div>

    <div class="form-actions">
      <div class="forgot-password">
        <router-link to="/forgot-password">忘记密码？</router-link>
      </div>
    </div>

    <button type="submit" class="submit-button" :disabled="isSubmitting">
      {{ isSubmitting ? '登录中...' : '登录' }}
    </button>

    <div class="divider">
      <span>或</span>
    </div>

    <div class="social-login">
      <button type="button" class="social-button" @click="router.push('/register')">
        注册新账号
      </button>
    </div>
  </form>
</template>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';
</style>
