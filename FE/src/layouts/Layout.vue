<template>
  <div class="app-container">
    <!-- 顶部导航栏 -->
    <header class="top-header">
      <div class="header-content">
        <div class="brand" @click="goHome" style="cursor: pointer;">
          <div class="brand-icon">
            <el-icon size="28" color="white"><TrendCharts /></el-icon>
          </div>
          <h1 class="brand-title">CS Goods</h1>
        </div>
        
        <nav class="nav-menu">
          <el-menu
            :default-active="$route.path"
            mode="horizontal"
            router
            class="header-menu"
            :ellipsis="false"
          >
            <el-menu-item
              v-for="route in menuRoutes"
              :key="route.path"
              :index="route.path"
              class="nav-item"
            >
              <span>{{ route.meta.title }}</span>
            </el-menu-item>
          </el-menu>
        </nav>
        
        <div class="header-actions">
          <el-space size="large">
            <!-- 用户信息下拉菜单 -->
            <el-dropdown trigger="click" @command="handleCommand">
              <div class="user-dropdown">
                <div class="user-avatar" :style="{ backgroundImage: getUserAvatarBg() }">
                </div>
                <div class="user-info">
                  <div class="user-name">{{ userStore.userInfo?.username }}</div>
                  <div class="user-type">{{ userStore.userTypeLabel }}</div>
                </div>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="dashboard">
                    <el-icon><HomeFilled /></el-icon>
                    <span>首页</span>
                  </el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isVip" command="ranking">
                    <el-icon><TrendCharts /></el-icon>
                    <span>饰品榜单</span>
                  </el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isVip" command="home">
                    <el-icon><DataAnalysis /></el-icon>
                    <span>挂刀/搬砖</span>
                  </el-dropdown-item>
                  <el-dropdown-item command="settings">
                    <el-icon><User /></el-icon>
                    <span>个人设置</span>
                  </el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isAdmin" command="admin">
                    <el-icon><Setting /></el-icon>
                    <span>管理中心</span>
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    <span>退出登录</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </el-space>
        </div>
      </div>
    </header>
    
    <!-- 主要内容区域 -->
    <main class="main-content">
      <div class="content-wrapper">
        <router-view />
      </div>
    </main>

    <!-- 页脚 -->
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import loginIcon from '@/assets/icons/login.png'
import registerIcon from '@/assets/icons/register.png'
import Footer from '@/components/Footer.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 菜单路由 - 根据权限过滤
const menuRoutes = computed(() => {
  return router.getRoutes()
    .filter(route => {
      if (!route.meta?.title || route.path === '/') return false
      
      // 过滤掉登录注册页
      if (route.path === '/login' || route.path === '/register') return false
      
      // 过滤掉hideInMenu的路由
      if (route.meta.hideInMenu) return false
      
      // 检查VIP权限
      if (route.meta.requiresVip && !userStore.isVip) return false
      
      // 检查管理员权限
      if (route.meta.requiresAdmin && !userStore.isAdmin) return false
      
      return true
    })
    .map(route => ({
      path: route.path,
      meta: route.meta
    }))
})

// 获取用户头像背景图
const getUserAvatarBg = () => {
  // VIP或管理员用login.png，普通用户用register.png
  if (userStore.isVip || userStore.isAdmin) {
    return `url(${loginIcon})`
  } else {
    return `url(${registerIcon})`
  }
}

// 点击logo回到首页
const goHome = () => {
  router.push('/app/dashboard')
}

// 处理下拉菜单命令
const handleCommand = (command: string) => {
  switch (command) {
    case 'dashboard':
      router.push('/app/dashboard')
      break
    case 'ranking':
      router.push('/app/ranking')
      break
    case 'home':
      router.push('/app/home')
      break
    case 'settings':
      router.push('/app/settings')
      break
    case 'admin':
      router.push('/app/admin')
      break
    case 'logout':
      userStore.logout()
      break
  }
}

// 初始化时加载用户信息
onMounted(() => {
  userStore.loadFromStorage()
  if (userStore.isLoggedIn) {
    userStore.getUserInfo()
  }
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background: #f5f7fa;
  position: relative;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
}

.top-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  height: 70px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 45px;
  height: 45px;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
}

.brand-title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #1890ff, #722ed1);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-menu {
  flex: 1;
  display: flex;
  justify-content: center;
}

.header-menu {
  border-bottom: none;
  background: transparent;
}

:deep(.el-menu--horizontal) {
  border-bottom: none;
}

:deep(.el-menu--horizontal .el-menu-item) {
  height: 69px;
  line-height: 69px;
  padding: 0 24px;
  border-radius: 0;
  margin: 0;
  margin-bottom: -1px;
  border-bottom: 2px solid transparent;
  font-size: 18px;
  color: #333;
}

:deep(.el-menu--horizontal .el-menu-item:hover) {
  background: transparent;
  color: #1890ff;
}

:deep(.el-menu--horizontal .el-menu-item.is-active) {
  background: transparent;
  color: #1890ff;
  border-bottom: 2px solid #1890ff;
}

:deep(.el-menu--horizontal .el-menu-item span) {
  font-size: 18px;
  font-weight: normal;
}

.nav-item {
  border-radius: 0;
  margin: 0;
}

.nav-item:hover {
  background: transparent;
}

.nav-item.is-active {
  background: transparent;
  color: #1890ff;
}

.header-actions {
  display: flex;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-dropdown:hover {
  background: rgba(24, 144, 255, 0.1);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.user-name {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  line-height: 1.2;
}

.user-type {
  font-size: 12px;
  color: #999;
  line-height: 1.2;
}

.main-content {
  flex: 1;
  padding: 20px;
}

.content-wrapper {
  width: 100%;
  background: transparent;
  overflow: visible;
}

/* 全局美化 */
:deep(.el-card) {
  border: none;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-table) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-table th) {
  background: linear-gradient(135deg, #f8faff, #e6f7ff);
  font-weight: 600;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .header-content {
    padding: 0 24px;
  }
  
  .main-content {
    padding: 20px 24px;
  }
}

@media (max-width: 768px) {
  .header-content {
    padding: 0 16px;
    height: 60px;
  }
  
  .brand-icon {
    width: 35px;
    height: 35px;
  }
  
  .brand-title {
    font-size: 16px;
  }
  
  .nav-menu {
    display: none;
  }
  
  .main-content {
    padding: 12px;
  }
  
  .content-wrapper {
    border-radius: 16px;
  }
}

@media (max-width: 480px) {
  .header-content {
    padding: 0 12px;
    height: 55px;
  }
  
  .brand {
    gap: 8px;
  }
  
  .brand-icon {
    width: 30px;
    height: 30px;
    border-radius: 8px;
  }
  
  .brand-title {
    font-size: 14px;
  }
  
  .main-content {
    padding: 8px;
  }
  
  .content-wrapper {
    border-radius: 12px;
  }
  
  .settings-btn {
    width: 35px;
    height: 35px;
  }
}
</style>