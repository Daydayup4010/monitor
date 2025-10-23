<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <h1 class="login-title">🎮 CSGO饰品系统</h1>
          <p class="login-subtitle">登录以访问完整功能</p>
        </div>

        <!-- 登录方式切换 -->
        <div class="login-tabs">
          <div 
            class="tab-item" 
            :class="{ active: loginType === 'password' }"
            @click="loginType = 'password'"
          >
            账户密码登录
          </div>
          <div 
            class="tab-item" 
            :class="{ active: loginType === 'email' }"
            @click="loginType = 'email'"
          >
            邮箱验证码登录
          </div>
        </div>

        <!-- 账户密码登录表单 -->
        <el-form
          v-if="loginType === 'password'"
          ref="passwordFormRef"
          :model="passwordForm"
          :rules="passwordRules"
          class="login-form"
          @submit.prevent="handlePasswordLogin"
        >
          <el-form-item prop="email">
            <el-input
              v-model="passwordForm.email"
              placeholder="请输入邮箱地址"
              size="large"
              prefix-icon="Message"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="passwordForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              prefix-icon="Lock"
              show-password
              @keyup.enter="handlePasswordLogin"
            />
          </el-form-item>

          <div class="form-actions">
            <router-link to="/reset-password" class="forgot-link">
              忘记密码？
            </router-link>
          </div>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="userStore.loading"
              @click="handlePasswordLogin"
              class="login-btn"
            >
              立即登录
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 邮箱验证码登录表单 -->
        <el-form
          v-else
          ref="emailFormRef"
          :model="emailForm"
          :rules="emailRules"
          class="login-form"
          @submit.prevent="handleEmailLogin"
        >
          <el-form-item prop="email">
            <el-input
              v-model="emailForm.email"
              placeholder="请输入邮箱地址"
              size="large"
              prefix-icon="Message"
            />
          </el-form-item>

          <el-form-item prop="code">
            <div class="code-input-group">
              <el-input
                v-model="emailForm.code"
                placeholder="请输入验证码"
                size="large"
                prefix-icon="Key"
                @keyup.enter="handleEmailLogin"
              />
              <el-button
                type="primary"
                size="large"
                :disabled="countdown > 0"
                :loading="sendingCode"
                @click="handleSendCode"
                class="code-btn"
              >
                {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
              </el-button>
            </div>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="userStore.loading"
              @click="handleEmailLogin"
              class="login-btn"
            >
              立即登录
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          还没有账户？
          <router-link to="/register" class="register-link">
            立即注册
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { authApi } from '@/api'
import { showMessage } from '@/utils/message'
import type { FormInstance, FormRules } from 'element-plus'
import type { LoginForm, EmailLoginForm } from '@/types'

const router = useRouter()
const userStore = useUserStore()

// 登录方式
const loginType = ref<'password' | 'email'>('password')

// 账户密码登录表单
const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive<LoginForm>({
  email: '',
  password: '',
})

const passwordRules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
}

// 邮箱验证码登录表单
const emailFormRef = ref<FormInstance>()
const emailForm = reactive<EmailLoginForm>({
  email: '',
  code: '',
})

const emailRules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' },
  ],
}

// 验证码相关
const sendingCode = ref(false)
const countdown = ref(0)
let countdownTimer: number | null = null

// 发送验证码
const handleSendCode = async () => {
  if (!emailForm.email) {
    showMessage.warning('请先输入邮箱地址')
    return
  }

  // 验证邮箱格式
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(emailForm.email)) {
    showMessage.warning('请输入正确的邮箱格式')
    return
  }

  sendingCode.value = true
  try {
    const response = await authApi.sendEmailCode({ email: emailForm.email })
    if (response.code === 1) {
      showMessage.success('验证码已发送，请查收邮件')
      // 开始倒计时
      countdown.value = 60
      countdownTimer = window.setInterval(() => {
        countdown.value--
        if (countdown.value <= 0 && countdownTimer) {
          clearInterval(countdownTimer)
          countdownTimer = null
        }
      }, 1000)
    }
  } catch (error: any) {
    // 错误已在拦截器中处理
  } finally {
    sendingCode.value = false
  }
}

// 账户密码登录
const handlePasswordLogin = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    const success = await userStore.login(passwordForm)
    if (success) {
      // 登录成功后，等待一下让状态更新
      await new Promise(resolve => setTimeout(resolve, 100))
      // 根据用户权限跳转
      if (userStore.isVip || userStore.isAdmin) {
        router.push('/home')
      } else {
        router.push('/settings')
      }
    }
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 邮箱验证码登录
const handleEmailLogin = async () => {
  if (!emailFormRef.value) return

  try {
    await emailFormRef.value.validate()
    const success = await userStore.emailLogin(emailForm)
    if (success) {
      // 登录成功后，等待一下让状态更新
      await new Promise(resolve => setTimeout(resolve, 100))
      // 根据用户权限跳转
      if (userStore.isVip || userStore.isAdmin) {
        router.push('/home')
      } else {
        router.push('/settings')
      }
    }
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 组件卸载时清除定时器
import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 450px;
}

.login-card {
  background: rgba(255, 255, 255, 0.98);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 32px;
  color: #1890ff;
  margin-bottom: 8px;
  font-weight: bold;
}

.login-subtitle {
  font-size: 16px;
  color: #666;
}

.login-tabs {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.tab-item {
  flex: 1;
  padding: 12px;
  text-align: center;
  border: 2px solid #e0e0e0;
  border-radius: 10px;
  background: white;
  cursor: pointer;
  font-weight: 600;
  color: #666;
  transition: all 0.3s ease;
}

.tab-item:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.tab-item.active {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  color: white;
  border-color: #1890ff;
}

.login-form {
  margin-top: 20px;
}

.form-actions {
  text-align: right;
  margin-bottom: 16px;
}

.forgot-link {
  color: #1890ff;
  text-decoration: none;
  font-size: 14px;
}

.forgot-link:hover {
  text-decoration: underline;
}

.code-input-group {
  display: flex;
  gap: 12px;
  width: 100%;
}

.code-input-group :deep(.el-input) {
  flex: 1;
}

.code-btn {
  white-space: nowrap;
  min-width: 120px;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border: none;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.4);
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  color: #666;
  font-size: 14px;
}

.register-link {
  color: #1890ff;
  text-decoration: none;
  font-weight: 600;
  margin-left: 4px;
}

.register-link:hover {
  text-decoration: underline;
}

/* 响应式 */
@media (max-width: 768px) {
  .login-card {
    padding: 24px;
  }

  .login-title {
    font-size: 24px;
  }

  .tab-item {
    font-size: 14px;
    padding: 10px;
  }
}
</style>

