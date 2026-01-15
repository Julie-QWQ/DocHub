<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import MaterialUpload from '@/components/material/MaterialUpload.vue'
import { useAuthStore } from '@/stores/auth'
import { onMounted } from 'vue'

const router = useRouter()
const authStore = useAuthStore()

// 权限检查
onMounted(() => {
  if (!authStore.isCommittee) {
    ElMessage.warning('您没有权限访问此页面')
    router.push('/materials')
  }
})

// 上传成功
const handleSuccess = () => {
  ElMessage.success('资料上传成功,等待管理员审核')
  router.push('/materials/my')
}

// 取消上传
const handleCancel = () => {
  router.back()
}
</script>

<template>
  <div class="upload-page">
    <div class="page-container">
      <main class="content-stream">
        <!-- 页面标题和返回 -->
        <header class="stream-header">
          <button class="back-btn" @click="router.back()">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            返回
          </button>
          <h1 class="page-title">上传资料</h1>
        </header>

        <!-- 上传表单 -->
        <div class="upload-content">
          <MaterialUpload @success="handleSuccess" @cancel="handleCancel" />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped lang="scss">
.upload-page {
  width: 100%;
  background: #ffffff;
  min-height: 100vh;
}

.page-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 24px;
}

.content-stream {
  max-width: 100%;
}

.stream-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f2f2f2;

  .back-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    background: none;
    border: none;
    color: #666;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: color 0.15s;

    &:hover {
      color: #1a1a1a;
    }

    svg {
      width: 20px;
      height: 20px;
    }
  }

  .page-title {
    font-size: 32px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0;
    letter-spacing: -0.02em;
  }
}

.upload-content {
  // 内容样式在 MaterialUpload 组件中定义
}

@media (max-width: 768px) {
  .page-container {
    padding: 24px 16px;
  }

  .stream-header {
    flex-direction: column;
    align-items: flex-start;

    .page-title {
      font-size: 26px;
    }
  }
}
</style>
