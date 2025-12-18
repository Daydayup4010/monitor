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
    <div class="user-banner" :class="{ 'user-banner-normal': !userStore.isVip && !userStore.isAdmin }">
      <div 
        class="user-avatar" 
        :style="{ backgroundImage: getUserAvatarBg(), backgroundSize: 'contain', backgroundPosition: 'center', backgroundRepeat: 'no-repeat' }"
      ></div>
      <div class="user-details">
        <div class="user-badge"><span v-if="badgeIcon">{{ badgeIcon }} </span>{{ userStore.userTypeLabel }}</div>
        <h3>{{ userStore.userInfo?.username }}</h3>
        <p :style="{ opacity: userStore.isVip || userStore.isAdmin ? 0.9 : 1 }">{{ userStore.userInfo?.email }}</p>
        <p v-if="userStore.isVip && userStore.userInfo?.vip_expiry" style="opacity: 0.8; font-size: 13px; margin-top: 4px;">
          VIP到期：{{ formatDate(userStore.userInfo.vip_expiry) }}
        </p>
      </div>
    </div>

    <!-- VIP会员卡片 -->
    <div class="vip-card">
      <div class="vip-header">
        <div class="vip-title">
          <span class="vip-icon">👑</span>
          <span>VIP会员</span>
        </div>
        <div class="vip-price">
          <span class="currency">￥</span>
          <span class="amount">19.9</span>
          <span class="period">/月</span>
        </div>
      </div>
      
      <div class="vip-features">
        <div class="feature-group">
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>查看饰品最大7天的价格走势变化</span>
          </div>
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>查看各大平台实时价格对比</span>
          </div>
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>查看今日、7日、15日、30日涨跌榜单</span>
          </div>
        </div>
        
        <div class="feature-divider"></div>
        
        <div class="feature-group">
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>搬砖比价功能使用特权</span>
          </div>
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>挂刀比价功能使用特权</span>
          </div>
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>大件求购功能使用特权</span>
          </div>
        </div>
        
        <div class="feature-divider"></div>
        
        <div class="feature-group">
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>尊贵VIP会员身份标识</span>
          </div>
          <div class="feature-item">
            <el-icon class="feature-icon"><Check /></el-icon>
            <span>优先获得新功能体验资格</span>
          </div>
        </div>
      </div>
      
      <div class="vip-action">
        <button 
          class="vip-btn" 
          @click="showVipSelectDialog = true"
          :disabled="isCreatingOrder"
        >
          <span v-if="userStore.isVip">续费会员</span>
          <span v-else>立即开通</span>
        </button>
      </div>
    </div>

    <!-- 选择购买时长弹窗 -->
    <el-dialog
      v-model="showVipSelectDialog"
      :title="userStore.isVip ? '续费会员' : '开通会员'"
      width="480px"
      :close-on-click-modal="false"
    >
      <div class="vip-select-content">
        <p class="select-tip">请选择购买时长</p>
        <div class="month-options">
          <div 
            v-for="option in monthOptions" 
            :key="option.months"
            class="month-option"
            :class="{ 'active': selectedMonths === option.months, 'recommend': option.recommend }"
            @click="selectedMonths = option.months"
          >
            <div class="option-tag" v-if="option.recommend">推荐</div>
            <div class="option-months">{{ option.months }}个月</div>
            <div class="option-price">
              <span class="current-price">￥{{ option.price }}</span>
            </div>
            <div class="option-unit">￥{{ (option.price / option.months).toFixed(1) }}/月</div>
          </div>
        </div>
        <div class="selected-info">
          <span>已选择：{{ selectedMonths }}个月</span>
          <span class="total-price">合计：<em>￥{{ selectedPrice }}</em></span>
        </div>
      </div>
      <template #footer>
        <el-button @click="showVipSelectDialog = false">取消</el-button>
        <el-button type="primary" :loading="isCreatingOrder" @click="handleCreateOrder">
          立即支付
        </el-button>
      </template>
    </el-dialog>

    <!-- 修改密码 -->
    <div class="card">
      <div class="card-title">🔐 账号安全</div>
      <div style="padding: 24px;">
        <p style="font-size: 14px; color: #666; margin-bottom: 16px;">
          通过邮箱验证码修改您的登录密码
        </p>
        <button class="btn btn-primary" @click="showPasswordDialog = true">
          修改密码
        </button>
      </div>
    </div>

    <!-- 修改密码弹窗 -->
    <el-dialog
      v-model="showPasswordDialog"
      :title="passwordStep === 1 ? '验证邮箱' : '设置新密码'"
      width="420px"
      :close-on-click-modal="false"
      @close="resetPasswordForm"
    >
      <!-- 步骤指示器 -->
      <div class="password-steps">
        <div class="step" :class="{ active: passwordStep >= 1, done: passwordStep > 1 }">
          <div class="step-num">{{ passwordStep > 1 ? '✓' : '1' }}</div>
          <div class="step-text">验证邮箱</div>
        </div>
        <div class="step-line" :class="{ active: passwordStep > 1 }"></div>
        <div class="step" :class="{ active: passwordStep >= 2 }">
          <div class="step-num">2</div>
          <div class="step-text">设置密码</div>
        </div>
      </div>

      <!-- 步骤1：验证邮箱 -->
      <el-form
        v-if="passwordStep === 1"
        ref="verifyFormRef"
        :model="passwordForm"
        :rules="verifyRules"
        label-position="top"
      >
        <el-form-item label="邮箱地址">
          <el-input
            :value="userStore.userInfo?.email"
            disabled
          />
        </el-form-item>

        <el-form-item label="邮箱验证码" prop="code">
          <div style="display: flex; gap: 12px;">
            <el-input
              v-model="passwordForm.code"
              placeholder="请输入6位验证码"
              style="flex: 1;"
              maxlength="6"
            />
            <el-button
              type="success"
              :disabled="passwordCountdown > 0 || sendingPasswordCode"
              @click="handleSendPasswordCode"
            >
              {{ passwordCountdown > 0 ? `${passwordCountdown}秒后重试` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>

      <!-- 步骤2：设置新密码 -->
      <el-form
        v-else
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-position="top"
      >
        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="passwordForm.password"
            type="password"
            placeholder="至少6个字符"
            show-password
          />
        </el-form-item>

        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button 
          v-if="passwordStep === 1" 
          type="primary" 
          :loading="isVerifyingCode" 
          @click="handleVerifyCode"
        >
          下一步
        </el-button>
        <template v-else>
          <el-button @click="passwordStep = 1">上一步</el-button>
          <el-button type="primary" :loading="isResettingPassword" @click="handleResetPassword">
            确认修改
          </el-button>
        </template>
      </template>
    </el-dialog>

    <!-- 支付二维码弹窗 -->
    <el-dialog
      v-model="showPayDialog"
      title="微信扫码支付"
      width="400px"
      :close-on-click-modal="false"
      @close="handleClosePayDialog"
    >
      <div class="pay-dialog-content">
        <div class="pay-amount">
          <span class="label">支付金额</span>
          <span class="price">￥{{ currentOrder?.amount?.toFixed(2) || '19.90' }}</span>
        </div>
        
        <div class="qrcode-container" v-if="currentOrder?.qrcode_img">
          <img :src="'data:image/png;base64,' + currentOrder.qrcode_img" alt="支付二维码" />
        </div>
        <div class="qrcode-loading" v-else>
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>正在加载二维码...</span>
        </div>
        
        <div class="pay-tips">
          <p>请使用微信扫描二维码完成支付</p>
          <p class="order-no" v-if="currentOrder?.order_no">订单号：{{ currentOrder.order_no }}</p>
        </div>
        
        <div class="pay-status" v-if="isPolling">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>等待支付中...</span>
        </div>
      </div>
    </el-dialog>

    <!-- 支付成功弹窗 -->
    <el-dialog
      v-model="showSuccessDialog"
      title="支付成功"
      width="360px"
    >
      <div class="success-dialog-content">
        <div class="success-icon">🎉</div>
        <h3>恭喜您成为VIP会员！</h3>
        <p>您的VIP权益已生效</p>
        <button class="btn btn-primary" @click="handleSuccessConfirm" style="margin-top: 20px;">
          开始体验
        </button>
      </div>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { paymentApi, authApi, type PaymentOrder } from '@/api'
import { showMessage } from '@/utils/message'
import { DataAnalysis, Check, Loading } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import loginIcon from '@/assets/icons/login.png'
import registerIcon from '@/assets/icons/register.png'

const router = useRouter()
const userStore = useUserStore()

const badgeIcon = computed(() => {
  if (userStore.isAdmin) return '👨‍💼'
  if (userStore.isVip) return '👑'
  return ''  // 普通用户不显示图标
})

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const goToHome = () => {
  router.push('/app/ranking')
}

// 获取用户头像背景图
const getUserAvatarBg = () => {
  if (userStore.isVip || userStore.isAdmin) {
    return `url(${loginIcon})`
  } else {
    return `url(${registerIcon})`
  }
}

// 支付相关状态
const showVipSelectDialog = ref(false)
const showPayDialog = ref(false)
const showSuccessDialog = ref(false)
const isCreatingOrder = ref(false)
const isPolling = ref(false)
const currentOrder = ref<PaymentOrder | null>(null)
const selectedMonths = ref(3) // 默认选择3个月
let pollingTimer: ReturnType<typeof setInterval> | null = null

// 月份选项
const monthOptions = [
  { months: 1, price: 19.9, recommend: false },
  { months: 3, price: 49.9, recommend: false },
  { months: 6, price: 89.9, recommend: false },
  { months: 12, price: 169.9, recommend: true },
]

// 计算选中的价格
const selectedPrice = computed(() => {
  const option = monthOptions.find(o => o.months === selectedMonths.value)
  return option ? option.price : 19.9
})

// 创建订单
const handleCreateOrder = async () => {
  isCreatingOrder.value = true
  try {
    const res = await paymentApi.createOrder(selectedMonths.value)
    if (res.code === 1 && res.data) {
      currentOrder.value = res.data
      showVipSelectDialog.value = false
      showPayDialog.value = true
      // 开始轮询订单状态
      startPolling(res.data.order_no)
    } else {
      showMessage.error(res.msg || '创建订单失败')
    }
  } catch (error) {
    showMessage.error('创建订单失败，请稍后重试')
  } finally {
    isCreatingOrder.value = false
  }
}

// 开始轮询订单状态
const startPolling = (orderNo: string) => {
  isPolling.value = true
  // 每3秒查询一次订单状态
  pollingTimer = setInterval(async () => {
    try {
      const res = await paymentApi.queryOrder(orderNo)
      if (res.code === 1 && res.data) {
        if (res.data.status === 1) {
          // 支付成功
          stopPolling()
          showPayDialog.value = false
          showSuccessDialog.value = true
          // 刷新用户信息
          await userStore.getUserInfo()
        }
      }
    } catch (error) {
      console.error('查询订单状态失败:', error)
    }
  }, 3000)
}

// 停止轮询
const stopPolling = () => {
  isPolling.value = false
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

// 关闭支付弹窗
const handleClosePayDialog = () => {
  stopPolling()
  currentOrder.value = null
}

// 支付成功确认
const handleSuccessConfirm = () => {
  showSuccessDialog.value = false
  // 如果之前不是VIP，跳转到首页
  if (userStore.isVip) {
    router.push('/app/dashboard')
  }
}

// 修改密码相关状态
const showPasswordDialog = ref(false)
const passwordStep = ref(1) // 1: 验证邮箱, 2: 设置密码
const verifyFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const isVerifyingCode = ref(false)
const isResettingPassword = ref(false)
const sendingPasswordCode = ref(false)
const passwordCountdown = ref(0)
let passwordCountdownTimer: ReturnType<typeof setInterval> | null = null

const passwordForm = reactive({
  code: '',
  password: '',
  confirmPassword: '',
})

// 步骤1验证规则
const verifyRules: FormRules = {
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' },
  ],
}

// 密码验证规则
const validatePassword = (_rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请输入新密码'))
  } else if (value.length < 6) {
    callback(new Error('密码长度不能少于6位'))
  } else if (/\s/.test(value)) {
    callback(new Error('密码不能包含空格'))
  } else {
    callback()
  }
}

const validateConfirmPassword = (_rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// 步骤2密码规则
const passwordRules: FormRules = {
  password: [
    { required: true, validator: validatePassword, trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

// 发送验证码
const handleSendPasswordCode = async () => {
  const email = userStore.userInfo?.email
  if (!email) {
    showMessage.warning('无法获取邮箱地址')
    return
  }

  sendingPasswordCode.value = true
  try {
    const response = await authApi.sendEmailCode({ email })
    if (response.code === 1) {
      showMessage.success('验证码已发送，请查收邮件')
      passwordCountdown.value = 60
      passwordCountdownTimer = setInterval(() => {
        passwordCountdown.value--
        if (passwordCountdown.value <= 0 && passwordCountdownTimer) {
          clearInterval(passwordCountdownTimer)
          passwordCountdownTimer = null
        }
      }, 1000)
    }
  } catch (error: any) {
    // 错误已在拦截器中处理
  } finally {
    sendingPasswordCode.value = false
  }
}

// 验证验证码（步骤1 -> 步骤2）
const handleVerifyCode = async () => {
  if (!verifyFormRef.value) return

  try {
    await verifyFormRef.value.validate()
    
    const email = userStore.userInfo?.email
    if (!email) {
      showMessage.warning('无法获取邮箱地址')
      return
    }
    
    isVerifyingCode.value = true
    
    // 调用验证接口验证验证码
    const response = await authApi.verifyEmailCode({ 
      email, 
      code: passwordForm.code 
    })
    
    if (response.code === 1) {
      // 验证通过，进入步骤2
      passwordStep.value = 2
    } else {
      showMessage.error(response.msg || '验证码错误')
    }
  } catch (error) {
    console.error('验证码验证失败:', error)
  } finally {
    isVerifyingCode.value = false
  }
}

// 重置密码
const handleResetPassword = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    
    const email = userStore.userInfo?.email
    if (!email) {
      showMessage.warning('无法获取邮箱地址')
      return
    }
    
    isResettingPassword.value = true
    
    const response = await authApi.resetPassword({
      email,
      code: passwordForm.code,
      password: passwordForm.password,
    })
    
    if (response.code === 1) {
      showMessage.success('密码修改成功')
      showPasswordDialog.value = false
      resetPasswordForm()
    }
  } catch (error) {
    console.error('修改密码失败:', error)
  } finally {
    isResettingPassword.value = false
  }
}

// 重置表单
const resetPasswordForm = () => {
  passwordStep.value = 1
  passwordForm.code = ''
  passwordForm.password = ''
  passwordForm.confirmPassword = ''
  if (passwordCountdownTimer) {
    clearInterval(passwordCountdownTimer)
    passwordCountdownTimer = null
  }
  passwordCountdown.value = 0
}

// 组件卸载时停止轮询
onUnmounted(() => {
  stopPolling()
  if (passwordCountdownTimer) {
    clearInterval(passwordCountdownTimer)
  }
})
</script>

<style scoped>
.settings-page {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

/* 密码修改步骤指示器 */
.password-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.step-num {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #e8e8e8;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
}

.step.active .step-num {
  background: #1890ff;
  color: #fff;
}

.step.done .step-num {
  background: #52c41a;
  color: #fff;
}

.step-text {
  font-size: 13px;
  color: #999;
  transition: all 0.3s;
}

.step.active .step-text {
  color: #333;
  font-weight: 500;
}

.step-line {
  width: 60px;
  height: 2px;
  background: #e8e8e8;
  margin: 0 16px;
  margin-bottom: 20px;
  transition: all 0.3s;
}

.step-line.active {
  background: #52c41a;
}

/* 普通用户信息横幅样式 */
.user-banner-normal {
  background: #fff !important;
  color: #333 !important;
  border: 1px solid #e8e8e8;
}

.user-banner-normal .user-badge {
  color: #666;
}

.user-banner-normal h3 {
  color: #333;
}

.user-banner-normal p {
  color: #666;
}

/* VIP卡片样式 */
.vip-card {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  color: #fff;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  position: relative;
  overflow: hidden;
}

.vip-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255, 215, 0, 0.1) 0%, transparent 70%);
  pointer-events: none;
}

.vip-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
  position: relative;
  z-index: 1;
}

.vip-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 700;
}

.vip-icon {
  font-size: 36px;
}

.vip-price {
  text-align: right;
}

.vip-price .currency {
  font-size: 20px;
  color: #ffd700;
}

.vip-price .amount {
  font-size: 48px;
  font-weight: 700;
  color: #ffd700;
  line-height: 1;
}

.vip-price .period {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
}

.vip-price .original-price {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  text-decoration: line-through;
  margin-top: 4px;
}

.vip-features {
  position: relative;
  z-index: 1;
}

.feature-group {
  margin-bottom: 16px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  font-size: 15px;
  color: rgba(255, 255, 255, 0.9);
}

.feature-icon {
  color: #4ade80;
  font-size: 18px;
  flex-shrink: 0;
}

.feature-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  margin: 16px 0;
}

.vip-action {
  margin-top: 28px;
  text-align: center;
  position: relative;
  z-index: 1;
}

.vip-btn {
  width: auto;
  min-width: 140px;
  padding: 14px 32px;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
  background: linear-gradient(135deg, #ffd700 0%, #ffed4a 100%);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(255, 215, 0, 0.4);
}

.vip-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 215, 0, 0.5);
}

.vip-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.vip-tip {
  margin-top: 12px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

/* 支付弹窗样式 */
.pay-dialog-content {
  text-align: center;
  padding: 20px;
}

.pay-amount {
  margin-bottom: 24px;
}

.pay-amount .label {
  font-size: 14px;
  color: #666;
  margin-right: 8px;
}

.pay-amount .price {
  font-size: 32px;
  font-weight: 700;
  color: #e6a23c;
}

.qrcode-container {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.qrcode-container img {
  width: 200px;
  height: 200px;
  border: 1px solid #eee;
  border-radius: 8px;
}

.qrcode-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #999;
  gap: 12px;
}

.qrcode-loading .el-icon {
  font-size: 32px;
}

.pay-tips {
  color: #666;
  font-size: 14px;
}

.pay-tips .order-no {
  margin-top: 8px;
  font-size: 12px;
  color: #999;
}

.pay-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  color: #409eff;
  font-size: 14px;
}

/* 成功弹窗样式 */
.success-dialog-content {
  text-align: center;
  padding: 20px;
}

.success-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.success-dialog-content h3 {
  font-size: 20px;
  color: #333;
  margin-bottom: 8px;
}

.success-dialog-content p {
  color: #666;
  font-size: 14px;
}

/* 旋转动画 */
.is-loading {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* VIP选择弹窗样式 */
.vip-select-content {
  padding: 10px 0;
}

.select-tip {
  text-align: center;
  color: #666;
  margin-bottom: 20px;
  font-size: 14px;
}

.month-options {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.month-option {
  position: relative;
  padding: 20px 16px;
  border: 2px solid #e8e8e8;
  border-radius: 12px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
}

.month-option:hover {
  border-color: #1890ff;
  background: #f0f7ff;
}

.month-option.active {
  border-color: #1890ff;
  background: #e6f4ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.month-option.recommend {
  border-color: #ff6b00;
}

.month-option.recommend.active {
  border-color: #ff6b00;
  background: #fff7e6;
  box-shadow: 0 0 0 2px rgba(255, 107, 0, 0.2);
}

.option-tag {
  position: absolute;
  top: -10px;
  right: 10px;
  background: linear-gradient(135deg, #ff6b00, #ff9500);
  color: #fff;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.option-months {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.option-price {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 8px;
  margin-bottom: 4px;
}

.option-price .current-price {
  font-size: 24px;
  font-weight: 700;
  color: #ff6b00;
}

.option-price .original-price {
  font-size: 14px;
  color: #999;
  text-decoration: line-through;
}

.option-unit {
  font-size: 12px;
  color: #999;
}

.selected-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  font-size: 14px;
  color: #666;
}

.selected-info .total-price {
  font-size: 14px;
}

.selected-info .total-price em {
  font-size: 20px;
  font-weight: 700;
  color: #ff6b00;
  font-style: normal;
}
</style>
