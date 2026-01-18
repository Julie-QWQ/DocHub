<template>
  <div class="default-layout">
    <AppSidebar />
    <div class="layout-main">
      <!-- 页面过渡动画 + 加载状态 -->
      <router-view v-slot="{ Component, route }">
        <transition :name="'page-fade'" mode="out-in">
          <component :is="Component" :key="route.path" />
        </transition>
      </router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useNavigationStore } from '@/stores/navigation'
import AppSidebar from '@/components/layout/AppSidebar.vue'

const route = useRoute()
const navigationStore = useNavigationStore()

// 当前页面是否正在加载
const isPageLoading = computed(() => {
  return navigationStore.isNavigating || navigationStore.isPageLoading(route.path)
})
</script>

<style scoped lang="scss">
.default-layout {
  display: flex;
  min-height: 100vh;
  padding-left: 260px;
}

.layout-main {
  flex: 1;
  width: 100%;
  position: relative;
}

// 页面切换动画
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.15s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}

.page-fade-enter-to,
.page-fade-leave-from {
  opacity: 1;
}

@media (max-width: 1024px) {
  .default-layout {
    padding-left: 80px;
  }
}

@media (max-width: 768px) {
  .default-layout {
    padding-left: 0;
  }
}
</style>
