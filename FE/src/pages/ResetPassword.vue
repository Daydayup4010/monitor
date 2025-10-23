<template>
  <div class="reset-password-page">
    <div class="reset-password-container">
      <div class="reset-password-card">
        <div class="reset-password-header">
          <h1 class="reset-password-title">🔑 重置密码</h1>
          <p class="reset-password-subtitle">通过邮箱验证码重置您的密码</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          class="reset-password-form"
          @submit.prevent="handleResetPassword"
        >
          <el-form-item prop="email">
            <el-input
              v-model="form.email"
              placeholder="请输入注册时的邮箱地址"
              size="large"
              prefix-icon="Message"
            />
          </el-form-item>

          <el-form-item prop="code">
            <div class="code-input-group">
              <el-input
                v-model="form.code"
                placeholder="请输入验证码"
                size="large"
                prefix-icon="Key"
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

          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="设置新密码"
              size="large"
              prefix-icon="Lock"
              show-password
            >
              <template #append>
                <span class="input-hint">至少6个字符</span>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="再次输入新密码"
              size="large"
              prefix-icon="Lock"
              show-password
              @keyup.enter="handleResetPassword"
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleResetPassword"
              class="reset-btn"
            >
              重置密码
            </el-button>
          </el-form-item>
        </el-form>

        <div class="reset-password-footer">
          想起密码了？
          <router-link to="/login" class="login-link">
            返回登录
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api'
import { showMessage } from '@/utils/message'
import type { FormInstance, FormRules } from 'element-plus'
import type { ResetPasswordForm } from '@/types'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: '',
  code: '',
  password: '',
  confirmPassword: '',
})

// 自定义验证器
const validateConfirmPassword = (rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

// 验证码相关
const sendingCode = ref(false)
const countdown = ref(0)
let countdownTimer: number | null = null

// 发送验证码
const handleSendCode = async () => {
  if (!form.email) {
    showMessage.warning('请先输入邮箱地址')
    return
  }

  // 验证邮箱格式
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.email)) {
    showMessage.warning('请输入正确的邮箱格式')
    return
  }

  sendingCode.value = true
  try {
    const response = await authApi.sendEmailCode({ email: form.email })
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

// 重置密码
const handleResetPassword = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true
    
    const { confirmPassword, ...resetData } = form
    const response = await authApi.resetPassword(resetData)
    
    if (response.code === 1) {
      showMessage.success('密码重置成功，请登录')
      setTimeout(() => {
        router.push('/login')
      }, 1500)
    }
  } catch (error) {
    console.error('重置密码失败:', error)
  } finally {
    loading.value = false
  }
}

// 组件卸载时清除定时器
onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>

<style scoped>
.reset-password-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.reset-password-container {
  width: 100%;
  max-width: 500px;
}

.reset-password-card {
  background: rgba(255, 255, 255, 0.98);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
}

.reset-password-header {
  text-align: center;
  margin-bottom: 32px;
}

.reset-password-title {
  font-size: 32px;
  color: #1890ff;
  margin-bottom: 8px;
  font-weight: bold;
}

.reset-password-subtitle {
  font-size: 16px;
  color: #666;
}

.reset-password-form {
  margin-top: 20px;
}

.input-hint {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
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

.reset-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border: none;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
  margin-top: 8px;
}

.reset-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.4);
}

.reset-password-footer {
  text-align: center;
  margin-top: 24px;
  color: #666;
  font-size: 14px;
}

.login-link {
  color: #1890ff;
  text-decoration: none;
  font-weight: 600;
  margin-left: 4px;
}

.login-link:hover {
  text-decoration: underline;
}

/* 响应式 */
@media (max-width: 768px) {
  .reset-password-card {
    padding: 24px;
  }

  .reset-password-title {
    font-size: 24px;
  }

  .code-input-group {
    flex-direction: column;
  }

  .code-btn {
    width: 100%;
  }
}
</style>

