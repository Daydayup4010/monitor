<template>
  <div class="admin-page">
    <el-tabs v-model="activeTab" type="border-card" class="admin-tabs">
      <!-- 用户管理 -->
      <el-tab-pane label="用户管理" name="users">
        <div class="tab-content">
          <UserManager />
        </div>
      </el-tab-pane>
      
      <!-- Token管理 -->
      <el-tab-pane label="Token管理" name="tokens">
        <div class="tab-content">
          <TokenManager />
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import UserManager from '@/components/UserManager.vue'
import TokenManager from '@/components/TokenManager.vue'

const activeTab = ref('users')
</script>

<style scoped>
.admin-page {
  padding: 20px;
  min-height: 100%;
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

/* 响应式 */
@media (max-width: 768px) {
  .admin-page {
    padding: 12px;
  }

  .tab-content {
    padding: 20px 16px;
    min-height: 500px;
  }

  :deep(.admin-tabs > .el-tabs__header .el-tabs__item) {
    font-size: 14px;
    padding: 0 16px;
    margin-right: 16px;
  }
}
</style>
