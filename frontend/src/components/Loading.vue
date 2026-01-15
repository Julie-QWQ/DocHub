<template>
  <div v-loading="true" class="app-loading" element-loading-text="加载中...">
    <div class="loading-placeholder"></div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  fullScreen?: boolean
}

withDefaults(defineProps<Props>(), {
  fullScreen: false
})
</script>

<style scoped lang="scss">
.app-loading {
  width: 100%;
  height: 100%;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-placeholder {
  width: 100%;
  height: 100%;
}
</style>
