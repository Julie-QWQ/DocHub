<template>
  <div class="announcements-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>公告管理</h1>
      <div class="header-actions">
        <span class="total-count">共 {{ total }} 条公告</span>
        <button class="create-btn" @click="handleCreate">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          创建公告
        </button>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <select v-model="filters.priority" class="filter-select" @change="handleSearch">
        <option value="">全部优先级</option>
        <option value="high">重要</option>
        <option value="normal">普通</option>
      </select>

      <select v-model="filters.is_active" class="filter-select" @change="handleSearch">
        <option value="">全部状态</option>
        <option value="true">启用</option>
        <option value="false">禁用</option>
      </select>

      <button v-if="hasFilter" class="reset-btn" @click="handleReset">
        重置筛选
      </button>
    </div>

    <!-- 公告列表 -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else class="announcements-list">
      <div
        v-for="announcement in announcements"
        :key="announcement.id"
        class="announcement-card"
        :class="{ inactive: !announcement.is_active }"
      >
        <div class="announcement-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M2 12h5"/>
            <path d="M4 12v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-6"/>
            <path d="M5 12V7a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v5"/>
            <path d="M15 12v-2"/>
            <path d="M15 12h5"/>
            <path d="M18 12v-2"/>
          </svg>
        </div>

        <div class="announcement-content">
          <div class="announcement-title">
            {{ announcement.title }}
            <span v-if="announcement.priority === 'high'" class="priority-badge high">重要</span>
          </div>
          <div class="announcement-preview">{{ announcement.content }}</div>
          <div class="announcement-meta">
            <span class="meta-item">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="16" y1="2" x2="16" y2="6"></line>
                <line x1="8" y1="2" x2="8" y2="6"></line>
                <line x1="3" y1="10" x2="21" y2="10"></line>
              </svg>
              {{ formatDate(announcement.created_at) }}
            </span>
            <span v-if="announcement.published_at" class="meta-item">
              发布: {{ formatDate(announcement.published_at) }}
            </span>
            <span v-if="announcement.expires_at" class="meta-item">
              过期: {{ formatDate(announcement.expires_at) }}
            </span>
          </div>
        </div>

        <div class="announcement-status">
          <span :class="['status-badge', announcement.is_active ? 'active' : 'inactive']">
            {{ announcement.is_active ? '启用' : '禁用' }}
          </span>
        </div>

        <div class="announcement-actions">
          <button class="action-btn" @click="handleEdit(announcement)">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
            </svg>
            编辑
          </button>
          <button class="action-btn danger" @click="handleDelete(announcement)">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            删除
          </button>
        </div>
      </div>

      <div v-if="!loading && announcements.length === 0" class="empty-state">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M2 12h5"/>
          <path d="M4 12v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-6"/>
          <path d="M5 12V7a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v5"/>
          <path d="M15 12v-2"/>
          <path d="M15 12h5"/>
          <path d="M18 12v-2"/>
        </svg>
        <p>{{ hasFilter ? '没有符合条件的公告' : '暂无公告' }}</p>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > 0" class="pagination">
      <span class="pagination-info">
        显示 {{ (pagination.page - 1) * pagination.page_size + 1 }}-
        {{ Math.min(pagination.page * pagination.page_size, total) }}
        / 共 {{ total }} 条
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
          :disabled="pagination.page * pagination.page_size >= total"
          @click="changePage(pagination.page + 1)"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 创建/编辑对话框 -->
    <div v-if="dialogVisible" class="dialog-overlay" @click.self="handleCancel">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ dialogTitle }}</h3>
          <button class="close-btn" @click="handleCancel">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <div class="form-group">
            <label for="title">标题 <span class="required">*</span></label>
            <input
              id="title"
              v-model="formData.title"
              type="text"
              class="form-input"
              placeholder="请输入公告标题（至少2个字符）"
              maxlength="200"
            />
            <span class="char-count" :class="{ invalid: formData.title.trim().length < 2 && formData.title.trim().length > 0 }">
              有效字符: {{ formData.title.trim().length }}/200 (最少2个)
            </span>
          </div>

          <div class="form-group">
            <label for="content">内容 <span class="required">*</span></label>
            <textarea
              id="content"
              v-model="formData.content"
              class="form-textarea"
              rows="6"
              placeholder="请输入公告内容（至少10个字符）"
              maxlength="10000"
            ></textarea>
            <span class="char-count" :class="{ invalid: formData.content.trim().length < 10 && formData.content.trim().length > 0 }">
              有效字符: {{ formData.content.trim().length }}/10000 (最少10个)
            </span>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>优先级</label>
              <div class="radio-group">
                <label class="radio-label">
                  <input
                    v-model="formData.priority"
                    type="radio"
                    value="normal"
                  />
                  <span>普通</span>
                </label>
                <label class="radio-label">
                  <input
                    v-model="formData.priority"
                    type="radio"
                    value="high"
                  />
                  <span>重要</span>
                </label>
              </div>
            </div>

            <div class="form-group">
              <label>状态</label>
              <div class="radio-group">
                <label class="radio-label">
                  <input
                    v-model="formData.is_active"
                    type="radio"
                    :value="true"
                  />
                  <span>启用</span>
                </label>
                <label class="radio-label">
                  <input
                    v-model="formData.is_active"
                    type="radio"
                    :value="false"
                  />
                  <span>禁用</span>
                </label>
              </div>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="published_at">发布时间</label>
              <input
                id="published_at"
                v-model="formData.published_at"
                type="datetime-local"
                class="form-input"
              />
              <span class="form-hint">留空则立即发布</span>
            </div>

            <div class="form-group">
              <label for="expires_at">过期时间</label>
              <input
                id="expires_at"
                v-model="formData.expires_at"
                type="datetime-local"
                class="form-input"
              />
              <span class="form-hint">留空则永不过期</span>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="handleCancel">取消</button>
          <button class="btn btn-primary" @click="handleSubmit">{{ formData.id ? '保存' : '创建' }}</button>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <div v-if="deleteDialogVisible" class="dialog-overlay" @click.self="deleteDialogVisible = false">
      <div class="dialog dialog-small">
        <div class="dialog-header">
          <h3>确认删除</h3>
          <button class="close-btn" @click="deleteDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <div class="delete-confirmation">
            <div class="warning-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21.73 18l-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
            </div>
            <p class="delete-message">
              确定要删除公告 <strong>{{ currentAnnouncement?.title }}</strong> 吗?
            </p>
            <p class="delete-warning">此操作不可恢复,公告将被永久删除。</p>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="deleteDialogVisible = false">取消</button>
          <button class="btn btn-danger" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { announcementApi } from '@/api/announcement'
import type { Announcement, CreateAnnouncementRequest, UpdateAnnouncementRequest } from '@/types'

const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('创建公告')
const deleteDialogVisible = ref(false)
const currentAnnouncement = ref<Announcement>()

const announcements = ref<Announcement[]>([])
const total = ref(0)

const pagination = reactive({
  page: 1,
  page_size: 20
})

const filters = reactive({
  priority: '',
  is_active: ''
})

const formData = reactive<CreateAnnouncementRequest & { id?: number; published_at?: string; expires_at?: string }>({
  title: '',
  content: '',
  priority: 'normal',
  is_active: true,
  published_at: undefined,
  expires_at: undefined
})

const hasFilter = computed(() => {
  return filters.priority !== '' || filters.is_active !== ''
})

// 加载公告列表
const loadAnnouncements = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }

    if (filters.priority) params.priority = filters.priority
    if (filters.is_active) params.is_active = filters.is_active === 'true'

    const response = await announcementApi.getAnnouncements(params)
    if (response.code === 0) {
      announcements.value = response.data.list
      total.value = response.data.total
    }
  } catch (error: any) {
    console.error('加载公告失败:', error)
    // 静默失败,不显示错误提示
  } finally {
    loading.value = false
  }
}

// 打开创建对话框
const handleCreate = () => {
  dialogTitle.value = '创建公告'
  Object.assign(formData, {
    title: '',
    content: '',
    priority: 'normal',
    is_active: true,
    published_at: undefined,
    expires_at: undefined
  })
  dialogVisible.value = true
}

// 打开编辑对话框
const handleEdit = (announcement: Announcement) => {
  dialogTitle.value = '编辑公告'
  Object.assign(formData, {
    id: announcement.id,
    title: announcement.title,
    content: announcement.content,
    priority: announcement.priority,
    is_active: announcement.is_active,
    published_at: announcement.published_at ? announcement.published_at.slice(0, 16) : '',
    expires_at: announcement.expires_at ? announcement.expires_at.slice(0, 16) : ''
  })
  dialogVisible.value = true
}

// 删除公告
const handleDelete = (announcement: Announcement) => {
  currentAnnouncement.value = announcement
  deleteDialogVisible.value = true
}

// 确认删除
const confirmDelete = async () => {
  if (!currentAnnouncement.value) return

  try {
    await announcementApi.deleteAnnouncement(currentAnnouncement.value.id)
    deleteDialogVisible.value = false
    await loadAnnouncements()
  } catch (error: any) {
    console.error('删除失败:', error)
  }
}

// 提交表单
const handleSubmit = async () => {
  // 验证必填字段
  const trimmedTitle = formData.title.trim()
  const trimmedContent = formData.content.trim()

  if (!trimmedTitle) {
    alert('请输入标题')
    return
  }
  if (trimmedTitle.length < 2) {
    alert('标题至少需要2个字符')
    return
  }
  if (!trimmedContent) {
    alert('请输入内容')
    return
  }
  if (trimmedContent.length < 10) {
    alert('内容至少需要10个字符')
    return
  }

  // 验证过期时间必须晚于发布时间
  if (formData.published_at && formData.expires_at) {
    const publishedTime = new Date(formData.published_at).getTime()
    const expiresTime = new Date(formData.expires_at).getTime()

    if (expiresTime <= publishedTime) {
      alert('过期时间必须晚于发布时间')
      return
    }
  }

  try {
    const data: any = {
      title: trimmedTitle,
      content: trimmedContent,
      priority: formData.priority,
      is_active: formData.is_active
    }

    // 处理时间字段
    if (formData.published_at) {
      data.published_at = new Date(formData.published_at).toISOString()
    }
    if (formData.expires_at) {
      data.expires_at = new Date(formData.expires_at).toISOString()
    }

    if (formData.id) {
      // 更新
      await announcementApi.updateAnnouncement(formData.id, data)
    } else {
      // 创建
      await announcementApi.createAnnouncement(data)
    }

    dialogVisible.value = false
    await loadAnnouncements()
  } catch (error: any) {
    console.error('操作失败:', error)
  }
}

// 取消
const handleCancel = () => {
  dialogVisible.value = false
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadAnnouncements()
}

// 重置筛选
const handleReset = () => {
  filters.priority = ''
  filters.is_active = ''
  handleSearch()
}

// 页码变化
const changePage = (page: number) => {
  pagination.page = page
  loadAnnouncements()
}

// 格式化时间
const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()

  // 检查日期是否有效
  if (isNaN(date.getTime())) return '-'

  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  // 如果是未来时间,显示具体日期
  if (days < 0) {
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  loadAnnouncements()
})
</script>

<style scoped lang="scss">
.announcements-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h1 {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 16px;
  }
}

.total-count {
  font-size: 14px;
  color: #6b7280;
}

.create-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  background: #111827;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #000000;
  }
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.filter-select {
  padding: 10px 14px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #111827;
  font-size: 14px;
  cursor: pointer;
  outline: none;
  transition: all 0.2s;

  &:focus {
    border-color: #111827;
  }
}

.reset-btn {
  padding: 10px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  color: #6b7280;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #9ca3af;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f4f6;
  border-top-color: #111827;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.announcements-list {
  min-height: 300px;
}

.announcement-card {
  display: grid;
  grid-template-columns: auto 1fr auto auto;
  gap: 16px;
  align-items: center;
  padding: 20px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  margin-bottom: 12px;
  transition: all 0.2s;

  &:hover {
    border-color: #d1d5db;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  }

  &.inactive {
    opacity: 0.6;
  }

  @media (max-width: 768px) {
    grid-template-columns: auto 1fr;
    gap: 12px;

    .announcement-status {
      display: none;
    }
  }
}

.announcement-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  background: #f3f4f6;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.announcement-content {
  min-width: 0;
  flex: 1;
}

.announcement-title {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.priority-badge {
  padding: 2px 8px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 4px;

  &.high {
    background: #fef3c7;
    color: #b45309;
  }
}

.announcement-preview {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
}

.announcement-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #9ca3af;
  flex-wrap: wrap;

  .meta-item {
    display: flex;
    align-items: center;
    gap: 4px;

    svg {
      flex-shrink: 0;
    }
  }
}

.status-badge {
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 6px;
  white-space: nowrap;

  &.active {
    background: #dcfce7;
    color: #15803d;
  }

  &.inactive {
    background: #e5e7eb;
    color: #6b7280;
  }
}

.announcement-actions {
  display: flex;
  gap: 8px;

  @media (max-width: 768px) {
    grid-column: 2;
    justify-content: flex-end;
  }
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  background: #f3f4f6;
  color: #374151;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #e5e7eb;
  }

  &.danger {
    color: #b91c1c;
    background: #fef2f2;

    &:hover {
      background: #fee2e2;
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #9ca3af;

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
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.pagination-info {
  font-size: 13px;
  color: #6b7280;
}

.pagination-controls {
  display: flex;
  gap: 8px;
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: #ffffff;
  color: #111827;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    background: #f9fafb;
    border-color: #d1d5db;
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
}

.dialog {
  background: #ffffff;
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.dialog-small {
  max-width: 420px;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;

  h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #111827;
  }

  .close-btn {
    background: none;
    border: none;
    padding: 4px;
    cursor: pointer;
    color: #6b7280;
    transition: color 0.2s;

    &:hover {
      color: #111827;
    }
  }
}

.dialog-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
  position: relative;

  &:last-child {
    margin-bottom: 0;
  }

  label {
    display: block;
    font-size: 13px;
    font-weight: 600;
    color: #374151;
    margin-bottom: 8px;
  }

  .required {
    color: #b91c1c;
  }
}

.char-count {
  position: absolute;
  bottom: -20px;
  right: 0;
  font-size: 12px;
  color: #9ca3af;

  &.invalid {
    color: #b91c1c;
    font-weight: 500;
  }
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  color: #111827;
  outline: none;
  transition: all 0.2s;
  font-family: inherit;

  &:focus {
    border-color: #111827;
  }

  &::placeholder {
    color: #9ca3af;
  }
}

.form-textarea {
  resize: vertical;
  min-height: 120px;
}

.form-hint {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: #9ca3af;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;

  @media (max-width: 600px) {
    grid-template-columns: 1fr;
  }
}

.radio-group {
  display: flex;
  gap: 16px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #111827;

  input[type="radio"] {
    cursor: pointer;
  }
}

.radio-group {
  display: flex;
  gap: 16px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #111827;

  input[type="radio"] {
    cursor: pointer;
    accent-color: #111827;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &.btn-secondary {
    background: #f3f4f6;
    color: #111827;

    &:hover {
      background: #e5e7eb;
    }
  }

  &.btn-primary {
    background: #111827;
    color: #ffffff;

    &:hover {
      background: #000000;
    }
  }

  &.btn-danger {
    background: #b91c1c;
    color: #ffffff;

    &:hover {
      background: #991b1b;
    }
  }
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
  color: #111827;
  margin: 0 0 12px 0;
  line-height: 1.5;

  strong {
    font-weight: 600;
    color: #b91c1c;
  }
}

.delete-warning {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}
</style>
