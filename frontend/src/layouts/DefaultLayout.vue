<template>
  <div class="default-layout">
    <AppSidebar />
    <div class="layout-main">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import AppSidebar from '@/components/layout/AppSidebar.vue'
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
