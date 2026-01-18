<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/composables/useAuth'
import { useRegisterForm } from '@/composables/useValidator'
import { useEmailVerification } from '@/composables/useEmailVerification'
import type { RegisterRequest } from '@/types'

const emit = defineEmits<{
  success: []
}>()

const { register } = useAuth()
const {
  sendVerificationCode,
  registerWithCode,
  countdown,
  isSending,
  error: verificationError
} = useEmailVerification()

const {
  username,
  email,
  password,
  confirmPassword,
  real_name,
  major,
  className,
  errors,
  isSubmitting,
  handleSubmit
} = useRegisterForm()

// 验证码
const verificationCode = ref('')

// 发送验证码
const handleSendCode = async () => {
  if (!email.value) {
    errors.value.email = '请先输入邮箱'
    ElMessage.error(errors.value.email)
    return
  }

  // 简单的邮箱格式验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email.value)) {
    errors.value.email = '请输入有效的邮箱地址'
    ElMessage.error(errors.value.email)
    return
  }

  const success = await sendVerificationCode({
    email: email.value,
    purpose: 'register'
  })

  if (success) {
    ElMessage.success('验证码已发送，请查收邮箱')
  } else if (verificationError.value) {
    ElMessage.error(verificationError.value)
  } else {
    ElMessage.error('发送验证码失败')
  }
}

const onSubmit = handleSubmit.withControlled(async () => {
  // 添加验证码验证
  if (!verificationCode.value) {
    errors.value.verificationCode = '请输入验证码'
    return
  }

  if (verificationCode.value.length !== 6) {
    errors.value.verificationCode = '验证码为6位数字'
    return
  }

  // 使用验证码注册
  const success = await registerWithCode({
    username: username.value,
    email: email.value,
    password: password.value,
    code: verificationCode.value
  })

  if (success) {
    emit('success')
  }
})
</script>

<template>
  <form @submit="onSubmit" class="auth-form auth-form--register">
    <div class="form-grid">
      <div class="form-group">
        <label for="username">用户名</label>
        <input
          id="username"
          v-model="username"
          type="text"
          placeholder="输入用户名（字母、数字）"
          autocomplete="username"
          :disabled="isSubmitting"
        />
        <p v-if="errors.username" class="error-message">{{ errors.username }}</p>
      </div>

      <div class="form-group">
        <label for="email">邮箱</label>
        <input
          id="email"
          v-model="email"
          type="email"
          placeholder="输入邮箱地址"
          autocomplete="email"
          :disabled="isSubmitting"
        />
        <p v-if="errors.email" class="error-message">{{ errors.email }}</p>
      </div>

      <div class="form-group">
        <label for="real_name">真实姓名</label>
        <input
          id="real_name"
          v-model="real_name"
          type="text"
          placeholder="输入真实姓名"
          autocomplete="name"
          :disabled="isSubmitting"
        />
        <p v-if="errors.real_name" class="error-message">{{ errors.real_name }}</p>
      </div>

      <div class="form-group">
        <label for="major">专业</label>
        <input
          id="major"
          v-model="major"
          type="text"
          placeholder="输入专业"
          :disabled="isSubmitting"
        />
        <p v-if="errors.major" class="error-message">{{ errors.major }}</p>
      </div>

      <div class="form-group">
        <label for="class">班级</label>
        <input
          id="class"
          v-model="className"
          type="text"
          placeholder="输入班级"
          :disabled="isSubmitting"
        />
        <p v-if="errors.class" class="error-message">{{ errors.class }}</p>
      </div>

      <div class="form-group full-width">
        <label for="password">密码</label>
        <input
          id="password"
          v-model="password"
          type="password"
          placeholder="输入密码（至少6位）"
          autocomplete="new-password"
          :disabled="isSubmitting"
        />
        <p v-if="errors.password" class="error-message">{{ errors.password }}</p>
      </div>

      <div class="form-group full-width">
        <label for="confirmPassword">确认密码</label>
        <input
          id="confirmPassword"
          v-model="confirmPassword"
          type="password"
          placeholder="再次输入密码"
          autocomplete="new-password"
          :disabled="isSubmitting"
        />
        <p v-if="errors.confirmPassword" class="error-message">{{ errors.confirmPassword }}</p>
      </div>

      <div class="form-group full-width">
        <label for="verificationCode">邮箱验证码</label>
        <div class="verification-input-group">
          <input
            id="verificationCode"
            v-model="verificationCode"
            type="text"
            placeholder="输入6位验证码"
            maxlength="6"
            :disabled="isSubmitting"
          />
          <button
            type="button"
            class="send-code-button"
            :disabled="isSending || countdown > 0"
            @click="handleSendCode"
          >
            {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
          </button>
        </div>
        <p v-if="errors.verificationCode" class="error-message">{{ errors.verificationCode }}</p>
        <p v-if="verificationError" class="error-message">{{ verificationError }}</p>
      </div>
    </div>

    <button type="submit" class="submit-button" :disabled="isSubmitting">
      {{ isSubmitting ? '注册中...' : '注册' }}
    </button>

  </form>
</template>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';

.verification-input-group {
  display: flex;
  gap: 10px;

  input {
    flex: 1;
  }

  .send-code-button {
    position: relative;
    flex-shrink: 0;
    padding: 0 20px;
    height: 44px;
    // 使用与背景协调的青绿色渐变
    background: linear-gradient(
      135deg,
      rgba(15, 118, 110, 0.9) 0%,
      rgba(20, 184, 166, 0.9) 100%
    );
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    color: white;
    border-radius: 8px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    white-space: nowrap;
    transition: all 0.3s ease;
    overflow: hidden;
    box-shadow:
      0 4px 12px rgba(15, 118, 110, 0.3),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);

    // 磨砂玻璃反光效果
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(
        90deg,
        transparent,
        rgba(255, 255, 255, 0.4),
        transparent
      );
      transition: left 0.5s ease;
    }

    // 鼠标悬停时的反光效果
    &:hover:not(:disabled) {
      background: linear-gradient(
        135deg,
        rgba(15, 118, 110, 1) 0%,
        rgba(20, 184, 166, 1) 100%
      );
      border-color: rgba(255, 255, 255, 0.4);
      box-shadow:
        0 6px 20px rgba(15, 118, 110, 0.4),
        inset 0 1px 0 rgba(255, 255, 255, 0.3);
      transform: translateY(-2px);

      &::before {
        left: 100%;
      }
    }

    // 点击效果
    &:active:not(:disabled) {
      transform: translateY(0);
      box-shadow:
        0 2px 8px rgba(15, 118, 110, 0.3),
        inset 0 1px 0 rgba(255, 255, 255, 0.2);
    }

    // 禁用状态
    &:disabled {
      background: rgba(144, 147, 153, 0.3);
      border-color: rgba(144, 147, 153, 0.2);
      color: rgba(255, 255, 255, 0.4);
      cursor: not-allowed;
      box-shadow: none;
      transform: none;

      &::before {
        display: none;
      }
    }

    // 倒计时状态
    &:not(:disabled):hover {
      animation: shimmer 1.5s infinite;
    }
  }
}

@keyframes shimmer {
  0% {
    box-shadow:
      0 4px 12px rgba(15, 118, 110, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
  50% {
    box-shadow:
      0 6px 20px rgba(15, 118, 110, 0.4),
      inset 0 1px 0 rgba(255, 255, 255, 0.3);
  }
  100% {
    box-shadow:
      0 4px 12px rgba(15, 118, 110, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }
}
</style>
