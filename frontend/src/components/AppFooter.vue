<template>
  <footer class="app-footer">
    <div class="footer-content">
      <p>&copy; {{ currentYear }} <SiteName />. All rights reserved.</p>
      <p class="footer-links">
        <router-link to="/about">关于我们</router-link>
        <span>·</span>
        <router-link to="/privacy">隐私政策</router-link>
        <span>·</span>
        <router-link to="/terms">服务条款</router-link>
      </p>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SiteName from '@/components/SiteName.vue'

const currentYear = computed(() => new Date().getFullYear())
</script>

<style scoped lang="scss">
.app-footer {
  padding: 20px;
  background-color: #ffffff;
  border-top: 1px solid #e4e7ed;
}

.footer-content {
  text-align: center;
  color: #909399;
  font-size: 14px;

  p {
    margin: 5px 0;
  }
}

.footer-links {
  a {
    color: #909399;
    margin: 0 5px;

    &:hover {
      color: #409eff;
    }
  }

  span {
    margin: 0 5px;
  }
}
</style>
