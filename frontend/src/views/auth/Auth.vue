<template>
  <div class="auth-split">
    <div class="auth-split-grid">
      <!-- 左侧介绍文字 -->
      <section class="auth-intro">
        <transition-group
          :name="isLogin ? 'slide-left' : 'slide-right'"
          tag="div"
          class="intro-wrapper"
        >
          <div
            :key="currentMode"
            class="intro-content"
          >
            <div>
              <h1 class="auth-brand-title interactive-text"><SiteName /></h1>
              <p class="auth-brand-subtitle interactive-text">{{ siteDescription }}</p>
            </div>
            <div>
              <h2 class="auth-intro-heading interactive-text">{{ introTitle }}</h2>
              <p class="auth-intro-text interactive-text">
                {{ introDescription }}
              </p>
              <ul class="auth-intro-list">
                <li v-for="(item, index) in introFeatures" :key="index" class="interactive-text">
                  {{ item }}
                </li>
              </ul>
            </div>
          </div>
        </transition-group>
      </section>

      <!-- 右侧表单卡片 -->
      <section class="auth-panel">
        <transition-group
          name="slide-up"
          tag="div"
          class="card-wrapper"
        >
          <div :key="currentMode" class="auth-card">
            <h2 class="auth-title">{{ cardTitle }}</h2>
            <p class="auth-description">{{ cardDescription }}</p>

            <!-- 登录表单 -->
            <LoginForm v-if="isLogin" @success="handleLoginSuccess" />

            <!-- 注册表单 -->
            <RegisterForm v-else @success="handleRegisterSuccess" />

            <div class="auth-footer">
              {{ footerText }}
              <router-link :to="footerLink" @click.prevent="toggleMode">
                {{ footerLinkText }}
              </router-link>
            </div>
          </div>
        </transition-group>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import LoginForm from '@/components/LoginForm.vue'
import RegisterForm from '@/components/RegisterForm.vue'
import SiteName from '@/components/SiteName.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteDescription = computed(() => systemStore.getConfig('site_description', '学院学习资料托管平台'))
const router = useRouter()
const route = useRoute()

// 当前模式: login 或 register
const currentMode = ref<'login' | 'register'>((route.name as string) === 'register' ? 'register' : 'login')

const isLogin = computed(() => currentMode.value === 'login')

// 登录页面的左侧内容
const loginIntro = {
  title: '把资料整理成你的学习基地',
  description: '课程资料、试卷与实验一站管理,快速搜索与分享。',
  features: [
    '搜索课件、试卷与实验资料',
    '收藏常用资源并追踪下载',
    '与同学共享高质量资料'
  ]
}

// 注册页面的左侧内容
const registerIntro = {
  title: '创建你的学习账号',
  description: '上传、整理并共享学习资料,让学习更高效。',
  features: [
    '快速上传与分类资料',
    '随时查看下载与收藏记录',
    '获取公告与审核通知'
  ]
}

const introTitle = computed(() => isLogin.value ? loginIntro.title : registerIntro.title)
const introDescription = computed(() => isLogin.value ? loginIntro.description : registerIntro.description)
const introFeatures = computed(() => isLogin.value ? loginIntro.features : registerIntro.features)

const cardTitle = computed(() => isLogin.value ? '登录' : '注册')
const cardDescription = computed(() => isLogin.value ? '欢迎回来,请登录你的账号' : '创建账号,开始使用学习资料平台')
const footerText = computed(() => isLogin.value ? '还没有账号?' : '已有账号?')
const footerLinkText = computed(() => isLogin.value ? '立即注册' : '立即登录')
const footerLink = computed(() => isLogin.value ? '/register' : '/login')

const toggleMode = () => {
  currentMode.value = isLogin.value ? 'register' : 'login'
  // 更新 URL
  const newPath = isLogin.value ? '/register' : '/login'
  window.history.pushState({}, '', newPath)
}

const handleLoginSuccess = () => {
  const redirect = route.query.redirect as string
  router.push(redirect || '/materials')
}

const handleRegisterSuccess = () => {
  currentMode.value = 'login'
  window.history.pushState({}, '', '/login')
}
</script>

<style scoped lang="scss">
@import '@/assets/styles/auth.scss';

.intro-wrapper,
.card-wrapper {
  position: relative;
  width: 100%;
}

.intro-content {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

// 右侧卡片向上滑动动画 - 使用 move 确保同步
.slide-up-move,
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateY(80px) scale(0.95);
  filter: blur(8px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateY(-80px) scale(0.95);
  filter: blur(8px);
}

// 确保离开的元素脱离文档流
.slide-up-leave-active {
  position: absolute;
  width: 100%;
}

// 左侧文字向左滑出动画
.slide-left-move,
.slide-left-enter-active,
.slide-left-leave-active {
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.slide-left-enter-from {
  opacity: 0;
  transform: translateX(80px) scale(0.95);
  filter: blur(4px);
}

.slide-left-leave-to {
  opacity: 0;
  transform: translateX(-80px) scale(0.95);
  filter: blur(4px);
}

// 确保离开的元素脱离文档流
.slide-left-leave-active {
  position: absolute;
  width: 100%;
}

// 左侧文字向右滑入动画
.slide-right-move,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.slide-right-enter-from {
  opacity: 0;
  transform: translateX(-80px) scale(0.95);
  filter: blur(4px);
}

.slide-right-leave-to {
  opacity: 0;
  transform: translateX(80px) scale(0.95);
  filter: blur(4px);
}

// 确保离开的元素脱离文档流
.slide-right-leave-active {
  position: absolute;
  width: 100%;
}
</style>
