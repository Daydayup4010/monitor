<template>
  <div class="admin-page">
    <div class="admin-layout">
      <!-- 左侧菜单栏 -->
      <div class="admin-sidebar">
        <div class="menu-item" :class="{ active: activeTab === 'users' }" @click="activeTab = 'users'">
          <div class="menu-icon">👥</div>
          <div class="menu-text">用户管理</div>
        </div>
        <div class="menu-item" :class="{ active: activeTab === 'orders' }" @click="activeTab = 'orders'">
          <div class="menu-icon">📋</div>
          <div class="menu-text">订单管理</div>
        </div>
        <div class="menu-item" :class="{ active: activeTab === 'notifications' }" @click="activeTab = 'notifications'">
          <div class="menu-icon">🔔</div>
          <div class="menu-text">通知管理</div>
        </div>
        <div class="menu-item" :class="{ active: activeTab === 'config' }" @click="activeTab = 'config'">
          <div class="menu-icon">⚙️</div>
          <div class="menu-text">系统配置</div>
        </div>
      </div>

      <!-- 右侧内容区 -->
      <div class="admin-content">
        <UserManager v-if="activeTab === 'users'" />
        <OrderManager v-else-if="activeTab === 'orders'" />
        <NotificationManager v-else-if="activeTab === 'notifications'" />
        <SystemConfigManager v-else-if="activeTab === 'config'" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import UserManager from '@/components/UserManager.vue'
import OrderManager from '@/components/OrderManager.vue'
import NotificationManager from '@/components/NotificationManager.vue'
import SystemConfigManager from '@/components/SystemConfigManager.vue'

const activeTab = ref('users')
</script>

<style scoped>
.admin-page {
  padding: 20px;
  max-width: 1600px;
  margin: 0 auto;
}

.admin-layout {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

/* 左侧菜单栏 */
.admin-sidebar {
  width: 200px;
  background: white;
  border-radius: 12px;
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  flex-shrink: 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 8px;
}

.menu-item:last-child {
  margin-bottom: 0;
}

.menu-item:hover {
  background: #f5f7fa;
}

.menu-item.active {
  background: #e6f7ff;
  color: #1890ff;
}

.menu-icon {
  font-size: 20px;
  line-height: 1;
}

.menu-text {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
}

.menu-item.active .menu-text {
  color: #1890ff;
  font-weight: 600;
}

/* 右侧内容区 */
.admin-content {
  flex: 1;
  min-width: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .admin-layout {
    flex-direction: column;
  }

  .admin-sidebar {
    width: 100%;
    display: flex;
    overflow-x: auto;
  }

  .menu-item {
    flex-direction: column;
    gap: 4px;
    margin-bottom: 0;
    margin-right: 8px;
    min-width: 80px;
  }

  .menu-text {
    font-size: 12px;
  }
}
</style>
