<template>
  <div class="app-container">
    <!-- 顶部导航栏 -->
    <header class="top-header">
      <div class="header-content">
        <div class="brand">
          <div class="brand-icon">
            <el-icon size="28" color="white"><TrendCharts /></el-icon>
          </div>
          <h1 class="brand-title">饰品差价监控</h1>
        </div>
        
        <nav class="nav-menu">
          <el-menu
            :default-active="$route.path"
            mode="horizontal"
            router
            class="header-menu"
          >
            <el-menu-item
              v-for="route in menuRoutes"
              :key="route.path"
              :index="route.path"
              class="nav-item"
            >
              <el-icon><component :is="route.meta.icon" /></el-icon>
              <span>{{ route.meta.title }}</span>
            </el-menu-item>
          </el-menu>
        </nav>
        
        <div class="header-actions">
          <el-space size="large">
            <TokenStatus />
            <el-button 
              circle
              type="text"
              class="settings-btn"
              @click="goToAdmin"
            >
              <el-icon size="18"><Setting /></el-icon>
            </el-button>
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
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TokenStatus from '@/components/TokenStatus.vue'
import { useTokenStore } from '@/stores/token'

const router = useRouter()
const tokenStore = useTokenStore()

// 菜单路由
const menuRoutes = computed(() => {
  return router.getRoutes()
    .filter(route => route.meta?.title && route.path !== '/')
    .map(route => ({
      path: route.path,
      meta: route.meta
    }))
})

// 跳转到管理中心
const goToAdmin = () => {
  router.push({ name: 'Admin' })
}

// 初始化时验证Token状态
onMounted(() => {
  tokenStore.verifyTokens()
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  position: relative;
  overflow-x: hidden;
}

.top-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(24, 144, 255, 0.2);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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
  max-width: 600px;
}

.header-menu {
  border-bottom: none;
  background: transparent;
}

.nav-item {
  border-radius: 8px;
  margin: 0 4px;
}

.nav-item:hover {
  background: rgba(24, 144, 255, 0.1);
}

.nav-item.is-active {
  background: rgba(24, 144, 255, 0.15);
  color: #1890ff;
}

.header-actions {
  display: flex;
  align-items: center;
}

.settings-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.05);
}

.settings-btn:hover {
  background: rgba(24, 144, 255, 0.1);
}

.main-content {
  min-height: calc(100vh - 70px);
  padding: 20px 40px;
}

.content-wrapper {
  width: 100%;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  overflow: hidden;
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
    min-height: calc(100vh - 60px);
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
    min-height: calc(100vh - 55px);
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