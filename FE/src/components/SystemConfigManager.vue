<template>
  <div class="config-manager">
    <div class="config-card">
      <h3 class="card-title">小程序配置</h3>
      <div class="config-item">
        <div class="config-info">
          <div class="config-label">VIP开通入口</div>
          <div class="config-desc">控制小程序端是否显示VIP开通入口，关闭后用户无法在小程序内开通VIP</div>
        </div>
        <div class="config-action">
          <label class="switch">
            <input 
              type="checkbox" 
              :checked="minappVipEnabled" 
              @change="toggleMinappVip"
              :disabled="loading"
            >
            <span class="slider"></span>
          </label>
          <span class="status-text" :class="{ enabled: minappVipEnabled }">
            {{ minappVipEnabled ? '已开启' : '已关闭' }}
          </span>
        </div>
      </div>
    </div>

    <div class="config-card">
      <h3 class="card-title">说明</h3>
      <div class="config-tips">
        <p>• 小程序审核期间建议关闭VIP开通入口</p>
        <p>• 审核通过后可重新开启</p>
        <p>• 此设置仅影响小程序端，不影响Web端</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { systemConfigApi } from '@/api'
import { showMessage } from '@/utils/message'

const loading = ref(false)
const minappVipEnabled = ref(false)

// 获取系统配置
const fetchConfig = async () => {
  try {
    const res = await systemConfigApi.getConfig()
    if (res.code === 1) {
      minappVipEnabled.value = res.data.minapp_vip_enabled || false
    }
  } catch (error) {
    console.error('获取配置失败:', error)
  }
}

// 切换小程序VIP开关
const toggleMinappVip = async (e: Event) => {
  const target = e.target as HTMLInputElement
  const newValue = target.checked

  loading.value = true
  try {
    await systemConfigApi.setMinappVipEnabled(newValue)
    minappVipEnabled.value = newValue
    showMessage.success(newValue ? 'VIP入口已开启' : 'VIP入口已关闭')
  } catch (error) {
    // 恢复原状态
    target.checked = !newValue
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped>
.config-manager {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.config-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0 0 20px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
}

.config-info {
  flex: 1;
}

.config-label {
  font-size: 15px;
  font-weight: 500;
  color: #262626;
  margin-bottom: 6px;
}

.config-desc {
  font-size: 13px;
  color: #8c8c8c;
}

.config-action {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-text {
  font-size: 13px;
  color: #8c8c8c;
  min-width: 50px;
}

.status-text.enabled {
  color: #52c41a;
}

/* Switch 样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 22px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: 0.3s;
  border-radius: 22px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #1890ff;
}

input:disabled + .slider {
  opacity: 0.6;
  cursor: not-allowed;
}

input:checked + .slider:before {
  transform: translateX(22px);
}

/* 说明卡片 */
.config-tips {
  background: #fafafa;
  border-radius: 8px;
  padding: 16px;
}

.config-tips p {
  font-size: 13px;
  color: #595959;
  margin: 0 0 8px 0;
  line-height: 1.6;
}

.config-tips p:last-child {
  margin-bottom: 0;
}
</style>
