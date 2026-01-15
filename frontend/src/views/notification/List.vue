<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useNotification } from '@/composables/useNotification'
import type { Notification } from '@/types'

const router = useRouter()

const {
  notifications,
  unreadCount,
  loading,
  total,
  page,
  size,
  fetchNotifications,
  markAsRead,
  deleteNotification,
  markAllAsRead,
  getNotificationTypeText,
  formatNotificationTime
} = useNotification()

// 当前选中的通知类型
const selectedType = ref<string>('')
const onlyUnread = ref(false)

// 弹窗状态
const deleteDialogVisible = ref(false)
const markAllDialogVisible = ref(false)
const currentNotification = ref<Notification>()

// 过滤后的通知列表
const filteredNotifications = computed(() => {
  let filtered = notifications.value

  if (selectedType.value) {
    filtered = filtered.filter(n => n.type === selectedType.value)
  }

  if (onlyUnread.value) {
    filtered = filtered.filter(n => n.status === 'unread')
  }

  return filtered
})

// 加载通知列表
const loadNotifications = async () => {
  await fetchNotifications()
}

// 处理通知点击
const handleNotificationClick = async (notification: Notification) => {
  if (notification.status === 'unread') {
    await markAsRead(notification.id)
  }

  // 如果有链接,跳转到对应页面
  if (notification.link) {
    router.push(notification.link)
  }
}

// 处理删除通知
const handleDelete = (notification: Notification) => {
  currentNotification.value = notification
  deleteDialogVisible.value = true
}

// 确认删除
const confirmDelete = async () => {
  if (!currentNotification.value) return

  try {
    await deleteNotification(currentNotification.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  }
}

// 处理全部标记为已读
const handleMarkAllRead = () => {
  if (unreadCount.value === 0) {
    ElMessage.info('没有未读通知')
    return
  }

  markAllDialogVisible.value = true
}

// 确认全部标记为已读
const confirmMarkAllRead = async () => {
  try {
    await markAllAsRead()
    ElMessage.success('操作成功')
    markAllDialogVisible.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 处理页码变化
const handlePageChange = async (newPage: number) => {
  await fetchNotifications({ page: newPage })
}

// 处理每页大小变化
const handleSizeChange = async (newSize: number) => {
  await fetchNotifications({ page: 1, size: newSize })
}

onMounted(() => {
  loadNotifications()
})
</script>

<template>
  <div class="notifications-page">
    <div class="page-header">
      <h1>通知中心</h1>
      <div class="header-actions">
        <el-button
          v-if="unreadCount > 0"
          type="primary"
          @click="handleMarkAllRead"
        >
          全部标为已读
        </el-button>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-radio-group v-model="selectedType">
        <el-radio-button label="">全部</el-radio-button>
        <el-radio-button label="system">系统通知</el-radio-button>
        <el-radio-button label="material">资料审核</el-radio-button>
        <el-radio-button label="committee">学委申请</el-radio-button>
        <el-radio-button label="report">举报处理</el-radio-button>
      </el-radio-group>

      <el-checkbox v-model="onlyUnread" style="margin-left: 20px">
        仅显示未读
      </el-checkbox>
    </div>

    <!-- 通知列表 -->
    <div v-loading="loading" class="notifications-container">
      <el-empty v-if="filteredNotifications.length === 0" description="暂无通知" />

      <div v-else class="notification-list">
        <div
          v-for="notification in filteredNotifications"
          :key="notification.id"
          class="notification-item"
          :class="{ unread: notification.status === 'unread' }"
        >
          <el-icon class="notification-icon" :size="24">
            <Bell />
          </el-icon>

          <div class="notification-content">
            <div class="notification-header">
              <span class="notification-type">
                {{ getNotificationTypeText(notification.type) }}
              </span>
              <span class="notification-time">
                {{ formatNotificationTime(notification.created_at) }}
              </span>
            </div>

            <div class="notification-title">{{ notification.title }}</div>
            <div class="notification-desc">{{ notification.content }}</div>

            <div v-if="notification.status === 'unread'" class="unread-badge">
              未读
            </div>
          </div>

          <div class="notification-actions">
            <el-button
              type="primary"
              link
              @click.stop="handleNotificationClick(notification)"
            >
              查看详情
            </el-button>
            <el-button
              type="danger"
              link
              @click.stop="handleDelete(notification)"
            >
              删除
            </el-button>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="total > 0" class="pagination-container">
        <el-pagination
          :current-page="page"
          :page-size="size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
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
              确定要删除这条通知吗？
            </p>
            <p class="delete-warning">此操作不可恢复，该通知将被永久删除。</p>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="deleteDialogVisible = false">取消</button>
          <button class="btn btn-danger" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>

    <!-- 全部标记为已读对话框 -->
    <div v-if="markAllDialogVisible" class="dialog-overlay" @click.self="markAllDialogVisible = false">
      <div class="dialog dialog-small">
        <div class="dialog-header">
          <h3>确认操作</h3>
          <button class="close-btn" @click="markAllDialogVisible = false">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 12 12"></path>
            </svg>
          </button>
        </div>

        <div class="dialog-body">
          <div class="delete-confirmation">
            <div class="info-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="16" x2="12" y2="12"></line>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
              </svg>
            </div>
            <p class="delete-message">
              确定要将所有通知标记为已读吗？
            </p>
            <p class="delete-warning">此操作将标记 {{ unreadCount }} 条未读通知为已读状态。</p>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="markAllDialogVisible = false">取消</button>
          <button class="btn btn-primary" @click="confirmMarkAllRead">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Bell } from '@element-plus/icons-vue'
</script>

<style scoped lang="scss">
.notifications-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    h1 {
      font-size: 24px;
      font-weight: 600;
      margin: 0;
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .filter-bar {
    display: flex;
    align-items: center;
    padding: 16px;
    background-color: var(--el-fill-color-blank);
    border-radius: 4px;
    margin-bottom: 20px;
  }

  .notifications-container {
    min-height: 400px;

    .notification-list {
      .notification-item {
        display: flex;
        align-items: flex-start;
        padding: 20px;
        background-color: var(--el-fill-color-blank);
        border: 1px solid var(--el-border-color-light);
        border-radius: 4px;
        margin-bottom: 12px;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          border-color: var(--el-color-primary);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }

        &.unread {
          background-color: var(--el-color-info-light-9);
          border-left: 3px solid var(--el-color-primary);
        }

        .notification-icon {
          margin-right: 16px;
          margin-top: 2px;
          color: var(--el-color-primary);
          flex-shrink: 0;
        }

        .notification-content {
          flex: 1;
          min-width: 0;

          .notification-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;

            .notification-type {
              font-size: 12px;
              color: var(--el-color-primary);
              background-color: var(--el-color-primary-light-9);
              padding: 2px 8px;
              border-radius: 4px;
            }

            .notification-time {
              font-size: 12px;
              color: var(--el-text-color-secondary);
            }
          }

          .notification-title {
            font-size: 16px;
            font-weight: 500;
            color: var(--el-text-color-primary);
            margin-bottom: 8px;
          }

          .notification-desc {
            font-size: 14px;
            color: var(--el-text-color-regular);
            line-height: 1.6;
          }

          .unread-badge {
            display: inline-block;
            margin-top: 8px;
            padding: 2px 8px;
            font-size: 12px;
            color: var(--el-color-danger);
            background-color: var(--el-color-danger-light-9);
            border-radius: 4px;
          }
        }

        .notification-actions {
          display: flex;
          flex-direction: column;
          gap: 8px;
          margin-left: 16px;
          flex-shrink: 0;
        }
      }
    }

    .pagination-container {
      display: flex;
      justify-content: center;
      margin-top: 24px;
    }
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
  max-height: 90vh;
  overflow-y: auto;

  &.dialog-small {
    max-width: 420px;
  }
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

.delete-confirmation {
  text-align: center;
  padding: 12px 0;
}

.warning-icon,
.info-icon {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  margin-bottom: 20px;
}

.warning-icon {
  background: #fef3c7;
  color: #b45309;
}

.info-icon {
  background: #dbeafe;
  color: #1d4ed8;
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
</style>
