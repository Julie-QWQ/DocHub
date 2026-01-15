<template>
  <div class="settings-container">
    <!-- 分类标签 -->
    <div class="tabs">
      <button
        :class="['tab-btn', { active: activeCategory === '' }]"
        @click="handleTabChange('')"
      >
        全部
      </button>
      <button
        :class="['tab-btn', { active: activeCategory === 'general' }]"
        @click="handleTabChange('general')"
      >
        通用设置
      </button>
      <button
        :class="['tab-btn', { active: activeCategory === 'upload' }]"
        @click="handleTabChange('upload')"
      >
        上传设置
      </button>
      <button
        :class="['tab-btn', { active: activeCategory === 'download' }]"
        @click="handleTabChange('download')"
      >
        下载设置
      </button>
      <button
        :class="['tab-btn', { active: activeCategory === 'auth' }]"
        @click="handleTabChange('auth')"
      >
        认证设置
      </button>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <div class="search-box">
        <input
          v-model="searchForm.keyword"
          type="text"
          placeholder="搜索配置键或说明..."
          class="search-input"
          @keyup.enter="handleSearch"
        />
      </div>
      <div class="action-buttons">
        <button class="btn btn-secondary" @click="handleReset">重置</button>
        <button class="btn btn-primary" @click="handleAdd">新增配置</button>
      </div>
    </div>

    <!-- 配置列表 -->
    <div v-loading="loading" class="config-list">
      <div v-for="config in configList" :key="config.config_key" class="config-card">
        <div class="config-header">
          <div class="config-key">{{ config.config_key }}</div>
          <span :class="['category-tag', config.category]">
            {{ getCategoryText(config.category) }}
          </span>
        </div>
        <div class="config-value">{{ config.config_value }}</div>
        <div class="config-description">{{ config.description }}</div>
        <div class="config-meta">
          <span class="update-time">更新于 {{ formatDate(config.updated_at) }}</span>
        </div>
        <div class="config-actions">
          <button class="action-btn" @click="handleEdit(config)">编辑</button>
          <button class="action-btn danger" @click="handleDelete(config)">删除</button>
        </div>
      </div>

      <div v-if="!loading && configList.length === 0" class="empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
        <p>暂无配置数据</p>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="pagination.total > 0" class="pagination">
      <span class="pagination-info">
        显示 {{ (pagination.page - 1) * pagination.page_size + 1 }}-{{ Math.min(pagination.page * pagination.page_size, pagination.total) }} / 共 {{ pagination.total }} 条
      </span>
      <div class="pagination-controls">
        <button
          class="page-btn"
          :disabled="pagination.page === 1"
          @click="changePage(pagination.page - 1)"
        >
          上一页
        </button>
        <button
          class="page-btn"
          :disabled="pagination.page * pagination.page_size >= pagination.total"
          @click="changePage(pagination.page + 1)"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 编辑/新增对话框 -->
    <div v-if="dialogVisible" class="dialog-overlay" @click.self="dialogVisible = false">
      <div class="dialog">
        <div class="dialog-header">
          <h2>{{ isEdit ? '编辑配置' : '新增配置' }}</h2>
          <button class="close-btn" @click="dialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label class="form-label">配置键</label>
              <input
                v-model="configForm.config_key"
                type="text"
                class="form-input"
                placeholder="请输入配置键(英文)"
                :disabled="isEdit"
              />
            </div>

            <div class="form-group">
              <label class="form-label">配置值</label>
              <textarea
                v-model="configForm.config_value"
                class="form-textarea"
                rows="4"
                placeholder="请输入配置值"
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">分类</label>
              <select v-model="configForm.category" class="form-select">
                <option value="general">通用</option>
                <option value="upload">上传</option>
                <option value="download">下载</option>
                <option value="auth">认证</option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">说明</label>
              <textarea
                v-model="configForm.description"
                class="form-textarea"
                rows="2"
                placeholder="请输入配置说明"
              ></textarea>
            </div>

            <div class="dialog-footer">
              <button type="button" class="btn btn-secondary" @click="dialogVisible = false">
                取消
              </button>
              <button type="submit" class="btn btn-primary" :disabled="submitting">
                {{ submitting ? '提交中...' : '确定' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <div v-if="deleteDialogVisible" class="dialog-overlay" @click.self="deleteDialogVisible = false">
      <div class="dialog dialog-small">
        <div class="dialog-header">
          <h2>确认删除</h2>
          <button class="close-btn" @click="deleteDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <div class="delete-confirmation">
            <div class="warning-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
            </div>
            <p class="delete-message">
              确定要删除配置 <strong>{{ currentConfig?.config_key }}</strong> 吗？
            </p>
            <p class="delete-warning">此操作不可恢复,该配置将被永久删除。</p>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="deleteDialogVisible = false">取消</button>
          <button class="btn btn-danger" @click="confirmDelete" :disabled="deleting">
            {{ deleting ? '删除中...' : '删除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getSystemConfigList,
  createSystemConfig,
  updateSystemConfig,
  deleteSystemConfig
} from '@/api/admin'
import { useSystemStore } from '@/stores/system'
import type { SystemConfig } from '@/api/admin'

const systemStore = useSystemStore()

const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const configList = ref<SystemConfig[]>([])

const activeCategory = ref('')
const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const dialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const isEdit = ref(false)
const currentConfig = ref<SystemConfig>()

const configForm = reactive({
  config_key: '',
  config_value: '',
  category: 'general',
  description: ''
})

// 加载配置列表
async function loadConfigList() {
  loading.value = true
  try {
    const response = await getSystemConfigList({
      page: pagination.page,
      page_size: pagination.page_size,
      category: activeCategory.value,
      keyword: searchForm.keyword
    })
    configList.value = response.data.list || []
    pagination.total = response.data.total || 0
  } catch (error: any) {
    ElMessage.error(error.message || '获取配置列表失败')
  } finally {
    loading.value = false
  }
}

// 切换标签
function handleTabChange(category: string) {
  activeCategory.value = category
  pagination.page = 1
  loadConfigList()
}

// 搜索
function handleSearch() {
  pagination.page = 1
  loadConfigList()
}

// 重置
function handleReset() {
  searchForm.keyword = ''
  handleSearch()
}

// 新增配置
function handleAdd() {
  isEdit.value = false
  Object.assign(configForm, {
    config_key: '',
    config_value: '',
    category: activeCategory.value || 'general',
    description: ''
  })
  dialogVisible.value = true
}

// 编辑配置
function handleEdit(row: SystemConfig) {
  isEdit.value = true
  Object.assign(configForm, {
    config_key: row.config_key,
    config_value: row.config_value,
    category: row.category,
    description: row.description
  })
  dialogVisible.value = true
}

// 提交表单
async function handleSubmit() {
  // 简单验证
  if (!configForm.config_key.trim()) {
    ElMessage.error('请输入配置键')
    return
  }
  if (!configForm.config_value.trim()) {
    ElMessage.error('请输入配置值')
    return
  }
  if (!configForm.category) {
    ElMessage.error('请选择分类')
    return
  }
  if (!configForm.description.trim()) {
    ElMessage.error('请输入配置说明')
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateSystemConfig({
        config_key: configForm.config_key,
        config_value: configForm.config_value
      })
      // 更新全局 store
      systemStore.setConfig(configForm.config_key, configForm.config_value)
      ElMessage.success('更新配置成功')
    } else {
      await createSystemConfig(configForm as any)
      ElMessage.success('创建配置成功')
    }
    dialogVisible.value = false
    loadConfigList()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 删除配置
function handleDelete(row: SystemConfig) {
  currentConfig.value = row
  deleteDialogVisible.value = true
}

// 确认删除
async function confirmDelete() {
  if (!currentConfig.value) return

  deleting.value = true
  try {
    await deleteSystemConfig(currentConfig.value.config_key)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    loadConfigList()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  } finally {
    deleting.value = false
  }
}

// 分页
function changePage(newPage: number) {
  pagination.page = newPage
  loadConfigList()
}

// 格式化日期
function formatDate(date: string) {
  return new Date(date).toLocaleString('zh-CN')
}

// 获取分类文本
function getCategoryText(category: string) {
  const map: Record<string, string> = {
    general: '通用',
    upload: '上传',
    download: '下载',
    auth: '认证'
  }
  return map[category] || category
}

onMounted(() => {
  loadConfigList()
})
</script>

<style scoped lang="scss">
.settings-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 1px solid #e5e5e5;
  padding-bottom: 0;
}

.tab-btn {
  position: relative;
  padding: 12px 16px;
  border: none;
  background: none;
  color: #666666;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;

  &:hover {
    color: #000000;
  }

  &.active {
    color: #000000;
    border-bottom-color: #000000;
  }
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  gap: 16px;

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
  }
}

.search-box {
  flex: 1;
  max-width: 400px;

  @media (max-width: 768px) {
    max-width: none;
  }
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cccccc;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;

  &:focus {
    outline: none;
    border-color: #000000;
  }

  &::placeholder {
    color: #999999;
  }
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.config-list {
  min-height: 300px;
  margin-bottom: 24px;
}

.config-card {
  padding: 20px;
  background: #ffffff;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  margin-bottom: 16px;
  transition: all 0.2s;

  &:hover {
    border-color: #cccccc;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  }
}

.config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.config-key {
  font-size: 16px;
  font-weight: 600;
  color: #000000;
  font-family: 'Courier New', monospace;
}

.category-tag {
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 6px;

  &.general {
    background: #e5e5e5;
    color: #666666;
  }

  &.upload {
    background: #dcfce7;
    color: #15803d;
  }

  &.download {
    background: #fef3c7;
    color: #b45309;
  }

  &.auth {
    background: #fee2e2;
    color: #b91c1c;
  }
}

.config-value {
  font-size: 14px;
  color: #333333;
  margin-bottom: 8px;
  font-family: 'Courier New', monospace;
  word-break: break-all;
  max-height: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.config-description {
  font-size: 13px;
  color: #666666;
  margin-bottom: 12px;
}

.config-meta {
  margin-bottom: 12px;
}

.update-time {
  font-size: 12px;
  color: #999999;
}

.config-actions {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.action-btn {
  padding: 6px 12px;
  border: 1px solid #cccccc;
  border-radius: 6px;
  background: #ffffff;
  color: #666666;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    border-color: #999999;
    background: #fafafa;
  }

  &.danger:hover {
    border-color: #b91c1c;
    color: #b91c1c;
    background: #fff5f5;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #999999;

  svg {
    margin-bottom: 16px;
  }

  p {
    margin: 0;
    font-size: 14px;
  }
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 24px;
  border-top: 1px solid #e5e5e5;
}

.pagination-info {
  font-size: 13px;
  color: #666666;
}

.pagination-controls {
  display: flex;
  gap: 8px;
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid #cccccc;
  border-radius: 6px;
  background: #ffffff;
  color: #000000;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    border-color: #999999;
    background: #fafafa;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.dialog {
  background: #ffffff;
  border: 1px solid #000000;
  border-radius: 8px;
  width: 100%;
  max-width: 560px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);

  &.dialog-small {
    max-width: 420px;
  }
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e5e5;

  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #000000;
  }

  .close-btn {
    background: none;
    border: none;
    padding: 4px;
    cursor: pointer;
    color: #666666;
    transition: color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      color: #000000;
    }
  }
}

.dialog-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #000000;
  margin-bottom: 8px;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #cccccc;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
  line-height: 1.5;
  transition: border-color 0.2s;

  &:focus {
    outline: none;
    border-color: #000000;
  }

  &:disabled {
    background: #f5f5f5;
    cursor: not-allowed;
  }

  &::placeholder {
    color: #999999;
  }
}

.form-textarea {
  resize: vertical;
}

.form-select {
  cursor: pointer;
  background-color: #ffffff;
}

.delete-confirmation {
  text-align: center;
  padding: 12px 0;
}

.warning-icon {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #fef3c7;
  color: #b45309;
  margin-bottom: 20px;
}

.delete-message {
  font-size: 16px;
  font-weight: 500;
  color: #000000;
  margin: 0 0 12px 0;
  line-height: 1.5;

  strong {
    font-weight: 600;
    color: #b91c1c;
  }
}

.delete-warning {
  font-size: 14px;
  color: #666666;
  margin: 0;
  line-height: 1.5;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #e5e5e5;
}

.btn {
  padding: 10px 20px;
  border: 1px solid #000000;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  background: #ffffff;
  color: #000000;

  &:hover:not(:disabled) {
    background: #f5f5f5;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &.btn-secondary {
    border-color: #cccccc;
    color: #666666;

    &:hover:not(:disabled) {
      border-color: #999999;
      background: #fafafa;
    }
  }

  &.btn-primary {
    background: #000000;
    color: #ffffff;

    &:hover:not(:disabled) {
      background: #333333;
    }
  }

  &.btn-danger {
    background: #b91c1c;
    border-color: #b91c1c;
    color: #ffffff;

    &:hover:not(:disabled) {
      background: #991b1b;
      border-color: #991b1b;
    }
  }
}
</style>
