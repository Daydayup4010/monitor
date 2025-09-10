<template>
  <div class="token-status">
    <el-space>
      <el-tooltip content="UU平台Token状态" placement="bottom">
        <div class="token-indicator">
          <el-icon 
            size="16" 
            :color="tokenStore.isTokenValid('uu') ? '#52c41a' : '#f5222d'"
          >
            <Key />
          </el-icon>
          <span class="token-label">UU</span>
        </div>
      </el-tooltip>
      
      <el-tooltip content="Buff平台Token状态" placement="bottom">
        <div class="token-indicator">
          <el-icon 
            size="16" 
            :color="tokenStore.isTokenValid('buff') ? '#52c41a' : '#f5222d'"
          >
            <Key />
          </el-icon>
          <span class="token-label">Buff</span>
        </div>
      </el-tooltip>
    </el-space>
  </div>
</template>

<script setup lang="ts">
import { useTokenStore } from '@/stores/token'

const tokenStore = useTokenStore()
</script>

<style scoped>
.token-status {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  backdrop-filter: blur(10px);
}

.token-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.8);
  transition: all 0.3s ease;
}

.token-indicator:hover {
  background: rgba(255, 255, 255, 0.9);
  transform: scale(1.05);
}

.token-label {
  font-size: 12px;
  font-weight: 600;
  color: #666;
}
</style>
