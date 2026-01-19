<template>
  <div class="profile-page">
    <div class="page-header">
      <div>
        <h1>个人中心</h1>
      </div>
      <div class="header-actions">
        <button v-if="canManageMaterials" class="action-btn" @click="router.push('/materials/my')">我的资料</button>
        <button class="action-btn" @click="router.push('/favorites')">我的收藏</button>
        <button class="action-btn" @click="router.push('/downloads')">下载历史</button>
      </div>
    </div>

    <section class="profile-card">
      <div class="avatar">
        <img v-if="avatar" :src="avatar" :alt="displayName" />
        <span v-else class="avatar-fallback">{{ defaultAvatar }}</span>
      </div>
      <div class="user-main">
        <div class="user-name">
          <h2>{{ displayName || '-' }}</h2>
          <span class="role-badge" :class="roleClass">{{ roleText }}</span>
          <span class="status-badge" :class="statusClass">{{ statusText }}</span>
        </div>
        <div class="user-meta">
          <span v-if="user?.username">@{{ user.username }}</span>
          <span v-if="user?.email">· {{ user.email }}</span>
        </div>
      </div>
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-label">上传资料</div>
          <div class="stat-value">{{ statValue(stats.uploads) }}</div>
        </div>
        <div class="stat-item">
          <div class="stat-label">收藏数</div>
          <div class="stat-value">{{ statValue(stats.favorites) }}</div>
        </div>
      <div class="stat-item">
        <div class="stat-label">下载次数</div>
        <div class="stat-value">{{ statValue(stats.downloads) }}</div>
      </div>
      <div class="stat-item">
        <div class="stat-label">今日剩余下载</div>
        <div class="stat-value">{{ dailyRemainingText }}</div>
      </div>
    </div>
  </section>

    <section class="info-grid">
      <div class="info-section">
        <h3>基本信息</h3>
        <div class="info-row">
          <span class="info-label">用户名</span>
          <span class="info-value">{{ user?.username || '-' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">真实姓名</span>
          <span class="info-value">{{ user?.real_name || '-' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">角色</span>
          <span class="info-value">{{ roleText }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">状态</span>
          <span class="info-value">{{ statusText }}</span>
        </div>
      </div>

      <div class="info-section">
        <h3>联系方式</h3>
        <div class="info-row">
          <span class="info-label">邮箱</span>
          <span class="info-value">{{ user?.email || '-' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">手机号</span>
          <span class="info-value">{{ user?.phone || '-' }}</span>
        </div>
      </div>

      <div class="info-section">
        <h3>学籍信息</h3>
        <div class="info-row">
          <span class="info-label">专业</span>
          <span class="info-value">{{ user?.major || '-' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label">班级</span>
          <span class="info-value">{{ user?.class || '-' }}</span>
        </div>
      </div>

      <div class="info-section">
        <h3>账号信息</h3>
        <div class="info-row">
          <span class="info-label">注册时间</span>
          <span class="info-value">{{ formatDate(user?.created_at) }}</span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { materialApi } from '@/api/material'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const displayName = computed(() => authStore.userName)
const avatar = computed(() => authStore.userAvatar)
const defaultAvatar = computed(() => (displayName.value ? displayName.value.charAt(0).toUpperCase() : 'U'))
const canManageMaterials = computed(() => authStore.isCommittee)

const statsLoading = ref(false)
const stats = ref({
  uploads: 0,
  favorites: 0,
  downloads: 0,
  dailyRemaining: 0,
  dailyUnlimited: false
})

const roleText = computed(() => {
  const role = user.value?.role
  const map: Record<string, string> = {
    admin: '管理员',
    committee: '学委',
    student: '学生'
  }
  return map[role || ''] || '-'
})

const roleClass = computed(() => {
  const role = user.value?.role
  if (role === 'admin') return 'role-admin'
  if (role === 'committee') return 'role-committee'
  return 'role-student'
})

const statusText = computed(() => {
  const status = user.value?.status
  const map: Record<string, string> = {
    active: '正常',
    banned: '已封禁'
  }
  return map[status || ''] || '-'
})

const statusClass = computed(() => {
  const status = user.value?.status
  return status === 'banned' ? 'status-banned' : 'status-active'
})

const statValue = (value: number) => {
  return statsLoading.value ? '-' : value
}

const dailyRemainingText = computed(() => {
  if (statsLoading.value) return '-'
  return stats.value.dailyUnlimited ? '不限' : `${stats.value.dailyRemaining}`
})

const formatDate = (date?: string) => {
  if (!date) return '-'
  const d = new Date(date)
  if (Number.isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const loadStats = async () => {
  if (!user.value?.id) return
  statsLoading.value = true
  try {
    const [uploadRes, favoriteRes, downloadRes, quotaRes] = await Promise.all([
      materialApi.getMaterials({
        page: 1,
        size: 1,
        uploader_id: user.value.id
      }),
      materialApi.getFavorites({
        page: 1,
        size: 1
      }),
      materialApi.getDownloadRecords({
        page: 1,
        page_size: 1
      }),
      materialApi.getDownloadQuota()
    ])

    if (uploadRes?.data) {
      stats.value.uploads = uploadRes.data.total || 0
    }
    if (favoriteRes?.data) {
      stats.value.favorites = favoriteRes.data.total || 0
    }
    if (downloadRes?.data) {
      stats.value.downloads = downloadRes.data.total || 0
    }
    if (quotaRes?.data) {
      stats.value.dailyUnlimited = quotaRes.data.unlimited
      stats.value.dailyRemaining = quotaRes.data.remaining
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取个人统计失败')
  } finally {
    statsLoading.value = false
  }
}

onMounted(async () => {
  if (!user.value) {
    try {
      await authStore.fetchUserInfo()
    } catch {
      return
    }
  }
  loadStats()
})
</script>

<style scoped lang="scss">
.profile-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px 48px;
  min-height: 100vh;
  background: #ffffff;
}

.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;

  h1 {
    font-size: 28px;
    font-weight: 700;
    color: #111827;
    margin: 0 0 6px 0;
  }

  .page-subtitle {
    margin: 0;
    color: #6b7280;
    font-size: 14px;
  }
}

.header-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.action-btn {
  padding: 8px 14px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  background: #ffffff;
  color: #111827;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    border-color: #d1d5db;
    background: #f9fafb;
  }
}

.profile-card {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 24px;
  align-items: center;
  padding: 24px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  margin-bottom: 24px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.04);
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  overflow: hidden;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.avatar-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #111827;
  color: #ffffff;
  font-size: 24px;
  font-weight: 600;
}

.user-main {
  min-width: 0;
}

.user-name {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;

  h2 {
    margin: 0;
    font-size: 22px;
    font-weight: 600;
    color: #111827;
  }
}

.user-meta {
  margin-top: 6px;
  font-size: 14px;
  color: #6b7280;
}

.role-badge,
.status-badge {
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 600;
  border-radius: 999px;
}

.role-admin {
  background: #fee2e2;
  color: #b91c1c;
}

.role-committee {
  background: #fef3c7;
  color: #b45309;
}

.role-student {
  background: #e5e7eb;
  color: #6b7280;
}

.status-active {
  background: #dcfce7;
  color: #15803d;
}

.status-banned {
  background: #fee2e2;
  color: #b91c1c;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  min-width: 260px;
}

.stat-item {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 12px 14px;
  background: #fafafa;
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 6px;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.info-section {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 16px;
  background: #fafafa;

  h3 {
    margin: 0 0 12px 0;
    font-size: 15px;
    font-weight: 600;
    color: #111827;
  }
}

.info-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
  color: #4b5563;
  padding: 6px 0;
}

.info-label {
  color: #6b7280;
  flex-shrink: 0;
}

.info-value {
  color: #111827;
  text-align: right;
  word-break: break-all;
}

@media (max-width: 900px) {
  .profile-card {
    grid-template-columns: 1fr;
    text-align: left;
  }

  .stats-grid {
    width: 100%;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .profile-page {
    padding: 24px 16px 32px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
