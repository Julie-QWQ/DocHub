<template>
  <div class="status-timeline">
    <el-timeline>
      <el-timeline-item
        v-for="(item, index) in timelineItems"
        :key="index"
        :timestamp="formatTimestamp(item.timestamp)"
        :type="item.type"
        :icon="item.icon"
        :color="item.color"
      >
        <div class="timeline-content">
          <div class="timeline-title">{{ item.title }}</div>
          <div v-if="item.description" class="timeline-description">
            {{ item.description }}
          </div>
          <div v-if="item.operator" class="timeline-operator">
            操作人: {{ item.operator }}
          </div>
        </div>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CircleCheck, Clock, CircleClose, Warning } from '@element-plus/icons-vue'
import { format } from '@/utils/format'

interface TimelineItem {
  timestamp?: string
  title: string
  description?: string
  operator?: string
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  icon?: any
  color?: string
}

interface Props {
  status: string
  createdAt: string
  updatedAt?: string
  reviewedAt?: string
  reviewer?: string
  reviewComment?: string
}

const props = defineProps<Props>()

const timelineItems = computed<TimelineItem[]>(() => {
  const items: TimelineItem[] = [
    {
      timestamp: props.createdAt,
      title: '提交申请',
      description: '申请已提交',
      type: 'primary',
      icon: Clock,
      color: '#409eff'
    }
  ]

  if (props.status === 'approved' && props.reviewedAt) {
    items.push({
      timestamp: props.reviewedAt,
      title: '审核通过',
      description: props.reviewComment || '申请已通过审核',
      operator: props.reviewer,
      type: 'success',
      icon: CircleCheck,
      color: '#67c23a'
    })
  } else if (props.status === 'rejected' && props.reviewedAt) {
    items.push({
      timestamp: props.reviewedAt,
      title: '审核拒绝',
      description: props.reviewComment || '申请未通过审核',
      operator: props.reviewer,
      type: 'danger',
      icon: CircleClose,
      color: '#f56c6c'
    })
  } else if (props.status === 'cancelled') {
    items.push({
      timestamp: props.updatedAt,
      title: '已取消',
      description: '申请人已取消申请',
      type: 'warning',
      icon: Warning,
      color: '#e6a23c'
    })
  } else if (props.status === 'pending') {
    items.push({
      title: '待审核',
      description: '申请正在审核中',
      type: 'info',
      icon: Clock,
      color: '#909399'
    })
  }

  return items
})

const formatTimestamp = (timestamp?: string): string => {
  if (!timestamp) return ''
  return format.datetime(timestamp)
}
</script>

<style scoped lang="scss">
.status-timeline {
  .timeline-content {
    .timeline-title {
      font-weight: 500;
      margin-bottom: 4px;
    }

    .timeline-description {
      font-size: 13px;
      color: var(--el-text-color-secondary);
      margin-bottom: 4px;
    }

    .timeline-operator {
      font-size: 12px;
      color: var(--el-text-color-placeholder);
    }
  }
}
</style>
