<template>
  <aside class="app-sidebar">
    <div class="sidebar-header">
      <h2>管理后台</h2>
    </div>
    <el-menu
      :default-active="activeMenu"
      class="sidebar-menu"
      :collapse="isCollapse"
      router
    >
      <el-menu-item index="/admin/dashboard">
        <el-icon><DataLine /></el-icon>
        <template #title>概览</template>
      </el-menu-item>
      <el-menu-item index="/admin/materials">
        <el-icon><Document /></el-icon>
        <template #title>资料管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/applications">
        <el-icon><UserFilled /></el-icon>
        <template #title>申请审核</template>
      </el-menu-item>
      <el-menu-item index="/admin/reports">
        <el-icon><Warning /></el-icon>
        <template #title>举报处理</template>
      </el-menu-item>
      <el-menu-item index="/admin/users">
        <el-icon><User /></el-icon>
        <template #title>用户管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/settings">
        <el-icon><Setting /></el-icon>
        <template #title>系统设置</template>
      </el-menu-item>
    </el-menu>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  DataLine,
  Document,
  UserFilled,
  Warning,
  User,
  Setting
} from '@element-plus/icons-vue'

const route = useRoute()

const activeMenu = computed(() => route.path)
const isCollapse = computed(() => false)
</script>

<style scoped lang="scss">
.app-sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 200px;
  background-color: #304156;
  z-index: 1000;
}

.sidebar-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);

  h2 {
    color: #ffffff;
    font-size: 18px;
  }
}

.sidebar-menu {
  border-right: none;
  background-color: transparent;

  :deep(.el-menu-item) {
    color: #bfcbd9;

    &:hover {
      background-color: rgba(255, 255, 255, 0.05);
    }

    &.is-active {
      color: #409eff;
      background-color: rgba(64, 158, 255, 0.1);
    }
  }

  :deep(.el-icon) {
    color: inherit;
  }
}
</style>
