<template>
  <div class="hot-rank">
    <el-card class="rank-card">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Trophy /></el-icon>
          <span class="header-title">{{ title }}</span>
        </div>
      </template>

      <div v-loading="loading" class="rank-list">
        <div v-if="items.length === 0 && !loading" class="empty">
          <el-empty description="暂无数据" />
        </div>

        <div
          v-for="(item, index) in items"
          :key="item.id"
          class="rank-item"
          @click="handleItemClick(item)"
        >
          <div class="rank-number" :class="`rank-${index + 1}`">
            {{ index + 1 }}
          </div>
          <div class="item-content">
            <div class="item-title">{{ item.title }}</div>
            <div class="item-meta">
              <span class="meta-item">
                <el-icon><Download /></el-icon>
                {{ item.download_count }}
              </span>
              <span class="meta-item">
                <el-icon><Star /></el-icon>
                {{ item.favorite_count }}
              </span>
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ item.view_count }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { Trophy, Download, Star, View } from '@element-plus/icons-vue'

interface HotItem {
  id: number
  title: string
  download_count: number
  favorite_count: number
  view_count: number
}

interface Props {
  title?: string
  items: HotItem[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  title: '热门排行榜',
  loading: false
})

const emit = defineEmits<{
  'item-click': [item: HotItem]
}>()

const handleItemClick = (item: HotItem) => {
  emit('item-click', item)
}
</script>

<style scoped lang="scss">
.hot-rank {
  .rank-card {
    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;

      .header-icon {
        font-size: 20px;
        color: #f59e0b;
      }

      .header-title {
        font-size: 16px;
        font-weight: 600;
      }
    }

    .rank-list {
      .empty {
        padding: 20px 0;
      }

      .rank-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 12px;
        margin-bottom: 8px;
        background: var(--el-fill-color-light);
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
          background: var(--el-fill-color);
          transform: translateX(4px);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }

        .rank-number {
          flex-shrink: 0;
          width: 32px;
          height: 32px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: var(--el-border-color);
          border-radius: 4px;
          font-weight: 600;
          font-size: 16px;
          color: var(--el-text-color-secondary);

          &.rank-1 {
            background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
            color: #92400e;
            box-shadow: 0 2px 8px rgba(255, 215, 0, 0.3);
          }

          &.rank-2 {
            background: linear-gradient(135deg, #c0c0c0 0%, #e5e7eb 100%);
            color: #4b5563;
            box-shadow: 0 2px 8px rgba(192, 192, 192, 0.3);
          }

          &.rank-3 {
            background: linear-gradient(135deg, #cd7f32 0%, #daa06d 100%);
            color: #78350f;
            box-shadow: 0 2px 8px rgba(205, 127, 50, 0.3);
          }
        }

        .item-content {
          flex: 1;
          min-width: 0;

          .item-title {
            font-size: 14px;
            color: var(--el-text-color-primary);
            font-weight: 500;
            margin-bottom: 8px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .item-meta {
            display: flex;
            gap: 16px;

            .meta-item {
              display: flex;
              align-items: center;
              gap: 4px;
              font-size: 12px;
              color: var(--el-text-color-secondary);

              .el-icon {
                font-size: 14px;
              }
            }
          }
        }
      }
    }
  }
}
</style>
