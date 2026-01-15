<script setup lang="ts">
import MaterialDetail from '@/components/material/MaterialDetail.vue'

const props = defineProps<{
  id: string
}>()
</script>

<template>
  <div class="material-detail-page">
    <MaterialDetail :material-id="props.id" />
  </div>
</template>

<style scoped lang="scss">
.material-detail-page {
  min-height: 100vh;
  background: #fafafa;
}
</style>
