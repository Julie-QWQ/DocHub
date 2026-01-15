<template>
  <el-card class="statistics-card" :body-style="{ padding: '20px' }">
    <div class="card-content">
      <div class="card-icon" :style="{ backgroundColor: iconBg }">
        <el-icon :size="24">
          <component :is="icon" />
        </el-icon>
      </div>
      <div class="card-info">
        <div class="card-title">{{ title }}</div>
        <div class="card-value">{{ value }}</div>
        <div class="card-trend" v-if="trend !== undefined">
          <span :class="trend >= 0 ? 'trend-up' : 'trend-down'">
            <el-icon><Top v-if="trend >= 0" /><Bottom v-else /></el-icon>
            {{ Math.abs(trend) }}%
          </span>
          <span class="trend-label">{{ trendLabel }}</span>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { Top, Bottom } from '@element-plus/icons-vue'

interface Props {
  title: string
  value: string | number
  icon: any
  iconBg?: string
  trend?: number
  trendLabel?: string
}

withDefaults(defineProps<Props>(), {
  iconBg: '#409EFF',
  trendLabel: '较昨日'
})
</script>

<style scoped lang="scss">
.statistics-card {
  transition: all 0.3s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
}

.card-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.card-info {
  flex: 1;
  min-width: 0;
}

.card-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.card-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.card-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;

  .trend-up {
    color: #67c23a;
  }

  .trend-down {
    color: #f56c6c;
  }

  .trend-label {
    color: #909399;
  }
}
</style>
