<template>
  <div class="auth-page-wrapper">
    <div class="auth-page">
      <div class="auth-card">
      <div class="auth-header">
        <h2 class="auth-title">
          <img src="@/assets/icons/login.png" style="height: 36px; width: auto; vertical-align: middle; margin-right: 8px; object-fit: contain;" alt="登录" />
          登录
        </h2>
        <p class="auth-subtitle">选择登录方式进入系统</p>
      </div>

      <!-- Tab切换 -->
      <div class="tabs">
        <div class="tab" :class="{ active: loginType === 'password' }" @click="loginType = 'password'">
          密码登录
        </div>
        <div class="tab" :class="{ active: loginType === 'email' }" @click="loginType = 'email'">
          验证码登录
        </div>
      </div>

      <!-- 密码登录表单 -->
      <el-form
        v-show="loginType === 'password'"
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        class="auth-form"
        @submit.prevent="handlePasswordLogin"
      >
        <div class="form-item">
          <label class="form-label">邮箱地址</label>
          <el-form-item prop="email">
            <el-input
              v-model="passwordForm.email"
              type="email"
              placeholder="请输入邮箱地址"
            />
          </el-form-item>
        </div>

        <div class="form-item">
          <label class="form-label">密码</label>
          <el-form-item prop="password">
            <el-input
              v-model="passwordForm.password"
              type="password"
              placeholder="请输入密码"
              show-password
              @keyup.enter="handlePasswordLogin"
            />
          </el-form-item>
        </div>

        <div class="form-item">
          <label class="form-label">图形验证码</label>
          <el-form-item prop="captchaCode">
            <div style="display: flex; gap: 12px; align-items: center;">
              <el-input
                v-model="passwordForm.captchaCode"
                placeholder="请输入验证码"
                style="flex: 1;"
                @keyup.enter="handlePasswordLogin"
              />
              <img 
                v-if="captchaImg" 
                :src="captchaImg" 
                @click="refreshCaptcha"
                style="height: 40px; cursor: pointer; border-radius: 4px; border: 1px solid #dcdfe6;"
                title="点击刷新验证码"
              />
              <span v-else style="color: #999; font-size: 12px;">加载中...</span>
            </div>
          </el-form-item>
        </div>

        <div style="text-align: right; margin-bottom: 12px;">
          <router-link to="/reset-password" style="color: #1890ff; font-size: 13px; text-decoration: none;">
            忘记密码？
          </router-link>
        </div>

        <button type="submit" class="btn btn-primary" :disabled="userStore.loading || !captchaId || isSubmitting">
          {{ (userStore.loading || isSubmitting) ? '登录中...' : '立即登录' }}
        </button>
      </el-form>

      <!-- 验证码登录表单 -->
      <el-form
        v-show="loginType === 'email'"
        ref="emailFormRef"
        :model="emailForm"
        :rules="emailRules"
        class="auth-form"
        @submit.prevent="handleEmailLogin"
      >
        <div class="form-item">
          <label class="form-label">邮箱地址</label>
          <el-form-item prop="email">
            <el-input
              v-model="emailForm.email"
              type="email"
              placeholder="请输入邮箱地址"
              @blur="checkEmailExist"
            >
              <template #suffix>
                <span v-if="emailChecking" style="color: #1890ff; font-size: 12px; padding-right: 8px;">检查中...</span>
                <span v-else-if="emailCheckResult === 'exist'" style="color: #52c41a; font-size: 12px; padding-right: 8px;">✓</span>
                <span v-else-if="emailCheckResult === 'notexist'" style="color: #ff4d4f; font-size: 12px; padding-right: 8px;">不存在</span>
              </template>
            </el-input>
          </el-form-item>
        </div>

        <div class="form-item">
          <label class="form-label">邮箱验证码</label>
          <el-form-item prop="code">
            <div style="display: flex; gap: 12px;">
              <el-input
                v-model="emailForm.code"
                placeholder="请输入验证码"
                style="flex: 1;"
                @keyup.enter="handleEmailLogin"
              />
              <button
                type="button"
                class="btn btn-success"
                style="white-space: nowrap;"
                :disabled="countdown > 0 || sendingCode"
                @click="handleSendCode"
              >
                {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
              </button>
            </div>
          </el-form-item>
        </div>

        <button type="submit" class="btn btn-primary" :disabled="userStore.loading">
          {{ userStore.loading ? '登录中...' : '立即登录' }}
        </button>
      </el-form>

      <div class="auth-footer">
        还没有账户？<router-link to="/register">立即注册</router-link>
      </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { authApi, captchaApi } from '@/api'
import { showMessage, debounce } from '@/utils'
import type { FormInstance, FormRules } from 'element-plus'
import type { LoginForm, EmailLoginForm } from '@/types'
import Footer from '@/components/Footer.vue'

const router = useRouter()
const userStore = useUserStore()

const loginType = ref<'password' | 'email'>('password')
const isSubmitting = ref(false) // 防止重复提交

// 图形验证码
const captchaId = ref('')
const captchaImg = ref('')

const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive<LoginForm & { captchaCode: string }>({
  email: '',
  password: '',
  captchaCode: '',
})

// 获取图形验证码
const refreshCaptcha = async () => {
  try {
    const response = await captchaApi.getCaptcha()
    if (response.code === 1) {
      captchaId.value = response.captcha_id
      captchaImg.value = response.captcha_img
    }
  } catch (error) {
    console.error('获取验证码失败:', error)
  }
}

// 页面加载时获取验证码
onMounted(() => {
  refreshCaptcha()
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
  captchaCode: [
    { required: true, message: '请输入图形验证码', trigger: 'blur' },
  ],
}

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

const sendingCode = ref(false)
const countdown = ref(0)
let countdownTimer: number | null = null

const emailChecking = ref(false)
const emailCheckResult = ref<'exist' | 'notexist' | ''>('')
const lastCheckedEmail = ref('')

const checkEmailExist = async () => {
  if (!emailForm.email) {
    emailCheckResult.value = ''
    lastCheckedEmail.value = ''
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(emailForm.email)) {
    emailCheckResult.value = ''
    return
  }

  if (emailForm.email === lastCheckedEmail.value && emailCheckResult.value) {
    return
  }

  emailChecking.value = true
  emailCheckResult.value = ''
  
  try {
    const response = await authApi.checkEmailExist({ email: emailForm.email })
    if (response.code === 1) {
      emailCheckResult.value = 'notexist'
      lastCheckedEmail.value = emailForm.email
    } else if (response.code === 1011) {
      emailCheckResult.value = 'exist'
      lastCheckedEmail.value = emailForm.email
    }
  } catch (error: any) {
    emailCheckResult.value = ''
  } finally {
    emailChecking.value = false
  }
}

const handleSendCode = async () => {
  if (!emailForm.email) {
    showMessage.warning('请先输入邮箱地址')
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(emailForm.email)) {
    showMessage.warning('请输入正确的邮箱格式')
    return
  }

  if (emailCheckResult.value === 'notexist') {
    showMessage.warning('该邮箱未注册，无法登录')
    return
  }

  if (emailCheckResult.value !== 'exist') {
    await checkEmailExist()
    if (emailCheckResult.value !== 'exist') {
      return
    }
  }

  sendingCode.value = true
  try {
    const response = await authApi.sendEmailCode({ email: emailForm.email })
    if (response.code === 1) {
      showMessage.success('验证码已发送，请查收邮件')
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

const handlePasswordLogin = async () => {
  if (!passwordFormRef.value || isSubmitting.value) return

  try {
    await passwordFormRef.value.validate()
    
    // 防止重复提交
    isSubmitting.value = true
    
    const success = await userStore.login({
      email: passwordForm.email,
      password: passwordForm.password,
      captcha_id: captchaId.value,
      captcha_code: passwordForm.captchaCode,
    })
    if (success) {
      await new Promise(resolve => setTimeout(resolve, 100))
      if (userStore.isVip || userStore.isAdmin) {
        router.push('/app/ranking')
      } else {
        router.push('/')  // 非VIP用户跳转到公开首页
      }
    } else {
      // 登录失败，刷新验证码
      refreshCaptcha()
      passwordForm.captchaCode = ''
    }
  } catch (error) {
    console.error('表单验证失败:', error)
    // 登录失败，刷新验证码
    refreshCaptcha()
    passwordForm.captchaCode = ''
  } finally {
    isSubmitting.value = false
  }
}

const handleEmailLogin = async () => {
  if (!emailFormRef.value) return

  try {
    await emailFormRef.value.validate()
    const success = await userStore.emailLogin(emailForm)
    if (success) {
      await new Promise(resolve => setTimeout(resolve, 100))
      if (userStore.isVip || userStore.isAdmin) {
        router.push('/app/ranking')
      } else {
        router.push('/')  // 非VIP用户跳转到公开首页
      }
    }
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const debouncedCheckEmail = debounce(checkEmailExist, 800)
watch(() => emailForm.email, (newEmail, oldEmail) => {
  if (newEmail !== oldEmail) {
    emailCheckResult.value = ''
    if (newEmail !== lastCheckedEmail.value) {
      lastCheckedEmail.value = ''
    }
  }
  
  if (newEmail && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(newEmail)) {
    debouncedCheckEmail()
  }
})

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>

<style scoped>
/* 样式在 unified.css 中 */
</style>

