<template>
  <div class="admin-layout">
    <AppSidebar />
    <div class="admin-main">
      <AppHeader />
      <div class="admin-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import AppSidebar from '@/components/AppSidebar.vue'
</script>

<style scoped lang="scss">
.admin-layout {
  display: flex;
  min-height: 100vh;
}

.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 200px;
  transition: margin-left 0.3s;

  &.collapsed {
    margin-left: 64px;
  }
}

.admin-content {
  flex: 1;
  padding: 20px;
  background-color: #f2f3f5;
}
</style>
