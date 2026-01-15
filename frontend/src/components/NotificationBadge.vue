<template>
  <el-badge :value="unreadCount" :hidden="unreadCount === 0" :max="99">
    <el-button :icon="Bell" circle @click="handleClick" />
  </el-badge>
</template>

<script setup lang="ts">
import { Bell } from '@element-plus/icons-vue'
import { useNotification } from '@/composables/useNotification'

const { unreadCount } = useNotification()

const emit = defineEmits<{
  click: []
}>()

const handleClick = () => {
  emit('click')
}
</script>

<style scoped lang="scss">
.el-badge {
  :deep(.el-badge__content) {
    transform: translateY(-50%) translateX(100%);
  }
}
</style>
