<template>
  <div class="admin-page">
    <el-tabs v-model="activeTab" type="border-card" class="admin-tabs">
      <el-tab-pane label="Token管理" name="tokens">
        <div class="tab-content">
          <TokenManager />
        </div>
      </el-tab-pane>
      
      <el-tab-pane label="系统设置" name="settings">
        <div class="tab-content">
          <SystemSettings />
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import TokenManager from '@/components/TokenManager.vue'
import SystemSettings from '@/components/SystemSettings.vue'

const activeTab = ref('tokens')
</script>

<style scoped>
.admin-page {
  padding: 32px 48px;
  min-height: 100vh;
  animation: fadeInUp 0.6s ease-out;
}

:deep(.admin-tabs) {
  border: none;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 16px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

:deep(.admin-tabs > .el-tabs__header) {
  background: linear-gradient(135deg, #f8faff 0%, #e6f7ff 100%);
  border-bottom: 2px solid #1890ff;
  margin: 0;
}

:deep(.admin-tabs > .el-tabs__header .el-tabs__nav-wrap) {
  padding: 0 24px;
}

:deep(.admin-tabs > .el-tabs__header .el-tabs__item) {
  border: none;
  height: 60px;
  line-height: 60px;
  font-size: 16px;
  font-weight: 600;
  color: #666;
  margin-right: 32px;
  padding: 0 24px;
  border-radius: 12px 12px 0 0;
  transition: all 0.3s ease;
}

:deep(.admin-tabs > .el-tabs__header .el-tabs__item:hover) {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
}

:deep(.admin-tabs > .el-tabs__header .el-tabs__item.is-active) {
  background: rgba(24, 144, 255, 0.15);
  color: #1890ff;
  border-bottom-color: transparent;
}

:deep(.admin-tabs > .el-tabs__content) {
  padding: 0;
  background: transparent;
}

.tab-content {
  padding: 32px;
  min-height: 600px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 16px;
}

/* 移动端响应式样式 */
@media (max-width: 1200px) {
  .admin-page {
    padding: 24px 32px;
  }
  
  .tab-content {
    padding: 24px;
  }
}

@media (max-width: 768px) {
  .admin-page {
    padding: 16px;
  }
  
  :deep(.admin-tabs) {
    border-radius: 12px;
  }
  
  :deep(.admin-tabs > .el-tabs__header .el-tabs__nav-wrap) {
    padding: 0 16px;
  }
  
  :deep(.admin-tabs > .el-tabs__header .el-tabs__item) {
    height: 50px;
    line-height: 50px;
    font-size: 14px;
    margin-right: 16px;
    padding: 0 16px;
  }
  
  .tab-content {
    padding: 20px 16px;
    min-height: 500px;
    border-radius: 12px;
  }
}

@media (max-width: 480px) {
  .admin-page {
    padding: 12px;
  }
  
  :deep(.admin-tabs) {
    border-radius: 8px;
  }
  
  :deep(.admin-tabs > .el-tabs__header .el-tabs__nav-wrap) {
    padding: 0 12px;
  }
  
  :deep(.admin-tabs > .el-tabs__header .el-tabs__item) {
    height: 45px;
    line-height: 45px;
    font-size: 13px;
    margin-right: 8px;
    padding: 0 12px;
  }
  
  .tab-content {
    padding: 16px 12px;
    min-height: 400px;
    border-radius: 8px;
  }
}
</style>