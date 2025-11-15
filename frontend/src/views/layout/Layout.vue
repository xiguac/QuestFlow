<template>
  <el-container class="layout-container">
    <!-- 左侧菜单 -->
    <el-aside width="200px" class="aside">
      <div class="logo-container">
        🌊 QuestFlow
      </div>
      <el-menu
        :default-active="activeMenu"
        class="aside-menu"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>我的表单</span>
        </el-menu-item>
        <el-menu-item index="/editor/new">
          <el-icon><Plus /></el-icon>
          <span>创建表单</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 右侧主内容区 -->
    <el-container class="main-container">
      <!-- 顶部 Header -->
      <el-header class="header">
        <div><!-- 预留给面包屑导航 --></div>
        <div class="user-info">
          <el-dropdown @command="handleCommand">
            <span class="el-dropdown-link">
              <el-avatar size="small" :src="avatarUrl" />
              <span class="username">{{ userStore.userInfo.username }}</span>
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 内容区域 Main -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  DataAnalysis,
  Plus,
  ArrowDown,
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const avatarUrl = computed(() => `https://i.pravatar.cc/150?u=${userStore.userInfo.username}`)

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    handleLogout()
  } else if (command === 'profile') {
    router.push('/profile')
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
}
.aside {
  background-color: #f5f7fa;
  border-right: 1px solid #e6e6e6;
  .logo-container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 60px;
    font-size: 20px;
    font-weight: bold;
    color: #409eff;
  }
  .aside-menu {
    border-right: none;
  }
}
.main-container {
  display: flex;
  flex-direction: column;
  .header {
    background-color: #ffffff;
    border-bottom: 1px solid #e6e6e6;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    flex-shrink: 0;
    .user-info .el-dropdown-link {
      cursor: pointer;
      display: flex;
      align-items: center;
      .username {
        margin-left: 8px;
        margin-right: 4px;
      }
    }
  }
  .main-content {
    background-color: #f0f2f5;
    padding: 20px;
    flex-grow: 1;
    overflow-y: auto;
  }
}
</style>
