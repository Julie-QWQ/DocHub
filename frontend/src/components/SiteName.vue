<template>
  <span>{{ siteName }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()
const siteName = computed(() => systemStore.getConfig('site_name', 'Study-UPC'))
</script>
