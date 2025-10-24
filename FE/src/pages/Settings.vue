<template>
  <div class="settings-page">
    <!-- 返回按钮 -->
    <div v-if="userStore.isVip" style="text-align: right; margin-bottom: 16px;">
      <button class="btn btn-primary" @click="goToHome" style="display: inline-flex; align-items: center; gap: 6px;">
        <el-icon><DataAnalysis /></el-icon>
        <span>返回饰品数据</span>
      </button>
    </div>

    <!-- 用户信息横幅 -->
    <div class="user-banner">
      <div class="user-avatar" :style="{ backgroundImage: getUserAvatarBg(), backgroundSize: 'contain', backgroundPosition: 'center', backgroundRepeat: 'no-repeat' }">
      </div>
      <div class="user-details">
        <div class="user-badge">{{ badgeIcon }} {{ userStore.userTypeLabel }}</div>
        <h3>{{ userStore.userInfo?.username }}</h3>
        <p style="opacity: 0.9;">{{ userStore.userInfo?.email }}</p>
        <p v-if="userStore.isVip && userStore.userInfo?.vip_expiry" style="opacity: 0.8; font-size: 13px; margin-top: 4px;">
          VIP到期：{{ formatDate(userStore.userInfo.vip_expiry) }}
        </p>
      </div>
    </div>

    <!-- 筛选参数设置 -->
    <div class="card">
      <div class="card-title">⚙️ 筛选参数配置</div>

      <el-form ref="formRef" :model="form" :rules="rules">
        <div class="two-cols">
          <div class="form-item">
            <label class="form-label">最小销售量</label>
            <el-form-item prop="min_sell_num">
              <el-input-number
                v-model="form.min_sell_num"
                :min="0"
                :max="10000"
                controls-position="right"
                placeholder="低于此值将被过滤"
                style="width: 100%"
              />
            </el-form-item>
          </div>

          <div class="form-item">
            <label class="form-label">最小价格差（元）</label>
            <el-form-item prop="min_diff">
              <el-input-number
                v-model="form.min_diff"
                :min="0"
                :max="1000"
                :step="0.1"
                :precision="2"
                controls-position="right"
                placeholder="UU与Buff的价格差"
                style="width: 100%"
              />
            </el-form-item>
          </div>

          <div class="form-item">
            <label class="form-label">最低价格（元）</label>
            <el-form-item prop="min_sell_price">
              <el-input-number
                v-model="form.min_sell_price"
                :min="0"
                :max="100000"
                :precision="2"
                controls-position="right"
                style="width: 100%"
              />
            </el-form-item>
          </div>

          <div class="form-item">
            <label class="form-label">最高价格（元）</label>
            <el-form-item prop="max_sell_price">
              <el-input-number
                v-model="form.max_sell_price"
                :min="0"
                :max="100000"
                :precision="2"
                controls-position="right"
                style="width: 100%"
              />
            </el-form-item>
          </div>
        </div>

        <div style="text-align: center; margin-top: 24px;">
          <button type="button" class="btn btn-primary" style="min-width: 160px;" :disabled="loading" @click="handleSave">
            {{ loading ? '保存中...' : '保存设置' }}
          </button>
          <button type="button" class="btn btn-secondary" style="min-width: 120px; margin-left: 12px;" @click="handleReset">
            重置
          </button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useSettingsStore } from '@/stores/settings'
import { showMessage } from '@/utils/message'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'

const router = useRouter()
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

const badgeIcon = computed(() => {
  if (userStore.isAdmin) return '👨‍💼'
  if (userStore.isVip) return '👑'
  return '👤'
})

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const loadSettings = async () => {
  if (!userStore.isVip) return
  
  await settingsStore.getSettings()
  Object.assign(form, settingsStore.settings)
}

const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true
    await settingsStore.updateSettings(form)
  } catch (error) {
    console.error('保存设置失败:', error)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  form.min_sell_num = 200
  form.min_diff = 1
  form.min_sell_price = 0
  form.max_sell_price = 10000
  showMessage.info('已恢复默认设置')
}

const goToHome = () => {
  router.push('/home')
}

// 获取用户头像背景图
const getUserAvatarBg = () => {
  // VIP或管理员用login.png，普通用户用register.png
  if (userStore.isVip || userStore.isAdmin) {
    return `url(/src/assets/icons/login.png)`
  } else {
    return `url(/src/assets/icons/register.png)`
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
/* 所有样式在unified.css中 */
.settings-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}
</style>

