<template>
  <div class="search-results">
    <div v-if="showSummary" class="search-summary">
      <span class="summary-text">
        找到 <strong>{{ total }}</strong> 条结果
        <span v-if="keyword" class="keyword">
          关键词: <el-tag size="small">{{ keyword }}</el-tag>
        </span>
      </span>
      <div class="sort-options">
        <el-radio-group v-model="localSortBy" size="small" @change="handleSortChange">
          <el-radio-button label="created_at">最新</el-radio-button>
          <el-radio-button label="download_count">下载量</el-radio-button>
          <el-radio-button label="favorite_count">收藏量</el-radio-button>
          <el-radio-button label="view_count">浏览量</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <div v-loading="loading" class="results-list">
      <div v-if="results.length === 0 && !loading" class="empty">
        <el-empty :description="emptyText">
          <el-button type="primary" @click="handleReset">重置搜索</el-button>
        </el-empty>
      </div>

      <el-card
        v-for="item in results"
        :key="item.id"
        class="result-card"
        shadow="hover"
        @click="handleItemClick(item)"
      >
        <div class="result-header">
          <el-tag
            :type="getCategoryType(item.category)"
            size="small"
            class="category-tag"
          >
            {{ getCategoryLabel(item.category) }}
          </el-tag>
          <h3 class="result-title">{{ item.title }}</h3>
        </div>

        <div class="result-description">
          {{ item.description }}
        </div>

        <div class="result-meta">
          <div class="meta-left">
            <span v-if="item.course_name" class="course-name">
              <el-icon><Reading /></el-icon>
              {{ item.course_name }}
            </span>
            <span class="uploader">
              <el-icon><User /></el-icon>
              {{ item.uploader_name }}
            </span>
            <span class="upload-time">
              <el-icon><Clock /></el-icon>
              {{ formatDate(item.created_at) }}
            </span>
          </div>
          <div class="meta-right">
            <span class="stat-item">
              <el-icon><Download /></el-icon>
              {{ item.download_count }}
            </span>
            <span class="stat-item">
              <el-icon><Star /></el-icon>
              {{ item.favorite_count }}
            </span>
            <span class="stat-item">
              <el-icon><View /></el-icon>
              {{ item.view_count }}
            </span>
          </div>
        </div>

        <div v-if="item.tags && item.tags.length > 0" class="result-tags">
          <el-tag
            v-for="(tag, index) in item.tags.slice(0, 5)"
            :key="index"
            size="small"
            class="tag-item"
          >
            {{ tag }}
          </el-tag>
        </div>
      </el-card>
    </div>

    <div v-if="totalPages > 1" class="pagination-wrapper">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
        @size-change="handlePageSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Reading, User, Clock, Download, Star, View } from '@element-plus/icons-vue'
import type { SearchResult, MaterialCategory } from '@/api/search'
import { useMaterialCategoryStore } from '@/stores/materialCategory'
import { useRouter } from 'vue-router'

interface Props {
  results: SearchResult[]
  total: number
  currentPage: number
  pageSize: number
  totalPages: number
  loading?: boolean
  keyword?: string
  sortBy?: string
  showSummary?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  keyword: '',
  sortBy: 'created_at',
  showSummary: true
})

const emit = defineEmits<{
  'page-change': [page: number]
  'page-size-change': [size: number]
  'sort-change': [sortBy: string]
  'item-click': [item: SearchResult]
  'reset': []
}>()

const router = useRouter()
const materialCategoryStore = useMaterialCategoryStore()
const localSortBy = ref(props.sortBy)

const emptyText = computed(() => {
  return props.keyword ? `未找到与 "${props.keyword}" 相关的结果` : '暂无数据'
})

// 从 store 获取资料类型名称
const getCategoryLabel = (category: MaterialCategory) => {
  return materialCategoryStore.getCategoryName(category)
}

const getCategoryType = (category: MaterialCategory) => {
  // 根据类别代码返回对应的标签类型
  const typeMap: Record<string, any> = {
    courseware: 'primary',
    textbook: 'success',
    reference: 'info',
    exam_paper: 'danger',
    note: 'warning',
    exercise: '',
    experiment: 'success',
    thesis: 'warning',
    other: 'info'
  }
  return typeMap[category] || 'default'
}

// 加载资料类型
onMounted(async () => {
  try {
    await materialCategoryStore.fetchActiveCategories()
  } catch (error: any) {
    console.error('加载资料类型失败:', error)
  }
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days} 天前`
  if (days < 30) return `${Math.floor(days / 7)} 周前`
  if (days < 365) return `${Math.floor(days / 30)} 月前`
  return `${Math.floor(days / 365)} 年前`
}

const handlePageChange = (page: number) => {
  emit('page-change', page)
}

const handlePageSizeChange = (size: number) => {
  emit('page-size-change', size)
}

const handleSortChange = (value: string) => {
  emit('sort-change', value)
}

const handleItemClick = (item: SearchResult) => {
  emit('item-click', item)
  router.push(`/materials/${item.id}`)
}

const handleReset = () => {
  emit('reset')
}
</script>

<style scoped lang="scss">
.search-results {
  .search-summary {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    margin-bottom: 16px;
    background: var(--el-fill-color-light);
    border-radius: 4px;

    .summary-text {
      font-size: 14px;
      color: var(--el-text-color-regular);

      strong {
        color: var(--el-color-primary);
        font-size: 18px;
        margin: 0 4px;
      }

      .keyword {
        margin-left: 12px;
      }
    }
  }

  .results-list {
    min-height: 400px;

    .empty {
      padding: 60px 20px;
    }

    .result-card {
      margin-bottom: 16px;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }

      .result-header {
        display: flex;
        align-items: flex-start;
        gap: 8px;
        margin-bottom: 12px;

        .category-tag {
          flex-shrink: 0;
          margin-top: 2px;
        }

        .result-title {
          flex: 1;
          margin: 0;
          font-size: 16px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }
      }

      .result-description {
        font-size: 14px;
        color: var(--el-text-color-regular);
        line-height: 1.6;
        margin-bottom: 12px;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }

      .result-meta {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;

        .meta-left,
        .meta-right {
          display: flex;
          gap: 16px;
        }

        .course-name,
        .uploader,
        .upload-time,
        .stat-item {
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

      .result-tags {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;

        .tag-item {
          cursor: pointer;
        }
      }
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 20px 0;
  }
}
</style>
