<template>
  <div class="search-bar">
    <el-input
      v-model="searchKeyword"
      :placeholder="placeholder"
      :prefix-icon="Search"
      clearable
      size="large"
      class="search-input"
      @keyup.enter="handleSearch"
      @clear="handleClear"
    >
      <template #append>
        <el-button :icon="Search" @click="handleSearch" />
      </template>
    </el-input>

    <!-- 热门搜索词 -->
    <div v-if="showHotKeywords && hotKeywords.length > 0" class="hot-keywords">
      <span class="label">热门搜索:</span>
      <el-tag
        v-for="(keyword, index) in hotKeywords.slice(0, 8)"
        :key="index"
        class="keyword-tag"
        @click="handleHotKeywordClick(keyword.keyword)"
      >
        {{ keyword.keyword }}
      </el-tag>
    </div>

    <!-- 搜索历史 -->
    <div v-if="showSearchHistory && searchHistory.length > 0 && !searchKeyword" class="search-history">
      <div class="history-header">
        <span class="label">搜索历史</span>
        <el-button link type="danger" @click="handleClearHistory">清空</el-button>
      </div>
      <div class="history-list">
        <div
          v-for="(item, index) in searchHistory"
          :key="index"
          class="history-item"
          @click="handleHistoryClick(item.keyword)"
        >
          <el-icon class="history-icon"><Clock /></el-icon>
          <span class="history-keyword">{{ item.keyword }}</span>
          <span class="history-count">{{ item.result_count }} 条结果</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search, Clock } from '@element-plus/icons-vue'
import { useSearchStore } from '@/stores/search'
import { ElMessage } from 'element-plus'

interface Props {
  placeholder?: string
  showHotKeywords?: boolean
  showSearchHistory?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请输入搜索关键词',
  showHotKeywords: true,
  showSearchHistory: true
})

const emit = defineEmits<{
  search: [keyword: string]
}>()

const searchStore = useSearchStore()
const searchKeyword = ref('')
const hotKeywords = ref<any[]>([])
const searchHistory = ref<any[]>([])

onMounted(async () => {
  if (props.showHotKeywords) {
    try {
      await searchStore.fetchHotKeywords(8)
      hotKeywords.value = searchStore.hotKeywords
    } catch (error) {
      console.error('获取热门搜索词失败:', error)
    }
  }

  if (props.showSearchHistory) {
    try {
      await searchStore.fetchSearchHistory(10)
      searchHistory.value = searchStore.searchHistory
    } catch (error) {
      console.error('获取搜索历史失败:', error)
    }
  }
})

const handleSearch = () => {
  if (!searchKeyword.value.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }
  emit('search', searchKeyword.value.trim())
}

const handleClear = () => {
  searchKeyword.value = ''
}

const handleHotKeywordClick = (keyword: string) => {
  searchKeyword.value = keyword
  emit('search', keyword)
}

const handleHistoryClick = (keyword: string) => {
  searchKeyword.value = keyword
  emit('search', keyword)
}

const handleClearHistory = async () => {
  try {
    await searchStore.clearSearchHistoryApi()
    searchHistory.value = []
    ElMessage.success('搜索历史已清空')
  } catch (error) {
    ElMessage.error('清空搜索历史失败')
  }
}
</script>

<style scoped lang="scss">
.search-bar {
  width: 100%;

  .search-input {
    margin-bottom: 16px;
  }

  .hot-keywords {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
    padding: 12px;
    background: var(--el-fill-color-light);
    border-radius: 4px;

    .label {
      font-size: 14px;
      color: var(--el-text-color-secondary);
      font-weight: 500;
    }

    .keyword-tag {
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      }
    }
  }

  .search-history {
    margin-top: 16px;
    padding: 12px;
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color);
    border-radius: 4px;

    .history-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .label {
        font-size: 14px;
        color: var(--el-text-color-secondary);
        font-weight: 500;
      }
    }

    .history-list {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }

    .history-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px;
      cursor: pointer;
      border-radius: 4px;
      transition: background 0.2s;

      &:hover {
        background: var(--el-fill-color-light);
      }

      .history-icon {
        color: var(--el-text-color-secondary);
      }

      .history-keyword {
        flex: 1;
        font-size: 14px;
        color: var(--el-text-color-primary);
      }

      .history-count {
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }
  }
}
</style>
