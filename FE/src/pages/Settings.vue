<template>
  <div class="settings-page">
    <!-- 用户信息卡片 -->
    <el-card class="user-info-card">
      <div class="user-info-content">
        <div class="user-avatar">
          {{ userStore.userInfo?.username?.charAt(0).toUpperCase() || 'U' }}
        </div>
        <div class="user-details">
          <div class="user-badge" :class="badgeClass">
            <span class="badge-icon">{{ badgeIcon }}</span>
            {{ userStore.userTypeLabel }}
          </div>
          <h2 class="user-name">{{ userStore.userInfo?.username }}</h2>
          <p class="user-email">{{ userStore.userInfo?.email }}</p>
          <p v-if="userStore.isVip && userStore.userInfo?.vip_expiry" class="user-vip-expiry">
            VIP到期时间：{{ formatDate(userStore.userInfo.vip_expiry) }}
          </p>
        </div>
      </div>
    </el-card>

    <!-- VIP权限提示（普通用户） -->
    <el-alert
      v-if="!userStore.isVip"
      title="需要VIP权限"
      type="warning"
      :closable="false"
      class="vip-alert"
    >
      <template #default>
        <p>您当前是普通用户，无法查看饰品数据内容。升级为VIP会员即可解锁所有功能。</p>
        <el-button type="primary" size="small" style="margin-top: 12px">
          联系管理员开通VIP
        </el-button>
      </template>
    </el-alert>

    <!-- 筛选参数设置（VIP用户） -->
    <el-card v-if="userStore.isVip" class="settings-card">
      <template #header>
        <div class="card-header">
          <span>⚙️ 个人筛选参数设置</span>
          <span class="card-subtitle">每个用户独立的饰品筛选配置</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="140px"
        class="settings-form"
      >
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="最小销售量" prop="min_sell_num">
              <el-input-number
                v-model="form.min_sell_num"
                :min="0"
                :max="10000"
                controls-position="right"
                style="width: 100%"
              />
              <div class="form-help">低于此值的商品将被过滤</div>
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="最小价格差异" prop="min_diff">
              <el-input-number
                v-model="form.min_diff"
                :min="0"
                :max="1000"
                :step="0.1"
                :precision="2"
                controls-position="right"
                style="width: 100%"
              />
              <div class="form-help">UU与Buff的最小价格差，单位：元</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="最低销售价格" prop="min_sell_price">
              <el-input-number
                v-model="form.min_sell_price"
                :min="0"
                :max="100000"
                :precision="2"
                controls-position="right"
                style="width: 100%"
              />
              <div class="form-help">低于此价格的商品将被过滤</div>
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="最高销售价格" prop="max_sell_price">
              <el-input-number
                v-model="form.max_sell_price"
                :min="0"
                :max="100000"
                :precision="2"
                controls-position="right"
                style="width: 100%"
              />
              <div class="form-help">高于此价格的商品将被过滤</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item class="form-actions">
          <el-button type="primary" size="large" :loading="loading" @click="handleSave">
            保存设置
          </el-button>
          <el-button size="large" @click="handleReset">
            恢复默认
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useSettingsStore } from '@/stores/settings'
import { showMessage } from '@/utils/message'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'

const userStore = useUserStore()
const settingsStore = useSettingsStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  min_sell_num: 0,
  min_diff: 0,
  min_sell_price: 0,
  max_sell_price: 10000,
})

const rules: FormRules = {
  min_sell_num: [
    { required: true, message: '请输入最小销售量', trigger: 'blur' },
  ],
  min_diff: [
    { required: true, message: '请输入最小价格差异', trigger: 'blur' },
  ],
  min_sell_price: [
    { required: true, message: '请输入最低销售价格', trigger: 'blur' },
  ],
  max_sell_price: [
    { required: true, message: '请输入最高销售价格', trigger: 'blur' },
    {
      validator: (rule: any, value: number, callback: Function) => {
        if (value <= form.min_sell_price) {
          callback(new Error('最高价格必须大于最低价格'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

// 用户徽章样式
const badgeClass = computed(() => {
  if (userStore.isAdmin) return 'badge-admin'
  if (userStore.isVip) return 'badge-vip'
  return 'badge-normal'
})

const badgeIcon = computed(() => {
  if (userStore.isAdmin) return '👨‍💼'
  if (userStore.isVip) return '👑'
  return '👤'
})

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 加载设置
const loadSettings = async () => {
  if (!userStore.isVip) return
  
  await settingsStore.getSettings()
  Object.assign(form, settingsStore.settings)
}

// 保存设置
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true
    await settingsStore.updateSettings(form)
    // 成功消息已在store中显示
  } catch (error) {
    console.error('保存设置失败:', error)
  } finally {
    loading.value = false
  }
}

// 恢复默认
const handleReset = () => {
  form.min_sell_num = 200
  form.min_diff = 1
  form.min_sell_price = 0
  form.max_sell_price = 10000
  showMessage.info('已恢复默认设置')
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings-page {
  padding: 20px 32px;
  max-width: 1200px;
  margin: 0 auto;
}

.user-info-card {
  margin-bottom: 20px;
  border-radius: 16px;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.user-info-card :deep(.el-card__body) {
  padding: 30px;
}

.user-info-content {
  display: flex;
  align-items: center;
  gap: 24px;
  color: white;
}

.user-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36px;
  font-weight: bold;
  color: #667eea;
}

.user-details {
  flex: 1;
}

.user-badge {
  display: inline-block;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 8px;
}

.badge-vip {
  background: rgba(255, 215, 0, 0.3);
  border: 2px solid gold;
}

.badge-normal {
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.5);
}

.badge-admin {
  background: rgba(255, 69, 0, 0.3);
  border: 2px solid #ff4500;
}

.badge-icon {
  margin-right: 4px;
}

.user-name {
  font-size: 24px;
  margin: 8px 0;
  font-weight: bold;
}

.user-email {
  opacity: 0.9;
  font-size: 16px;
  margin: 4px 0;
}

.user-vip-expiry {
  opacity: 0.8;
  font-size: 14px;
  margin-top: 8px;
}

.vip-alert {
  margin-bottom: 20px;
  border-radius: 12px;
}

.settings-card {
  border-radius: 16px;
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.card-subtitle {
  font-size: 14px;
  color: #999;
  font-weight: normal;
}

.settings-form {
  padding: 20px 0;
}

.form-help {
  font-size: 13px;
  color: #999;
  margin-top: 4px;
  line-height: 1.5;
}

.form-actions {
  margin-top: 32px;
  text-align: center;
}

.form-actions :deep(.el-form-item__content) {
  justify-content: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .settings-page {
    padding: 12px;
  }

  .user-info-content {
    flex-direction: column;
    text-align: center;
  }

  :deep(.el-col) {
    width: 100% !important;
  }
}
</style>

