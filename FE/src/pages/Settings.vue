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

    <!-- 提示信息 -->
    <div class="card">
      <div class="card-title">⚙️ 筛选参数配置</div>
      <div style="padding: 24px; text-align: center;">
        <p style="font-size: 15px; color: #595959; margin-bottom: 16px;">
          筛选参数（价格范围、成交量等）已移至"饰品数据"页面，方便您直接调整和查看结果。
        </p>
        <button class="btn btn-primary" @click="goToHome" style="padding: 12px 32px;">
          前往饰品数据页面
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const badgeIcon = computed(() => {
  if (userStore.isAdmin) return '👨‍💼'
  if (userStore.isVip) return '👑'
  return '👤'
})

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
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
</script>

<style scoped>
/* 所有样式在unified.css中 */
.settings-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}
</style>

