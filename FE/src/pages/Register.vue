<template>
  <div class="auth-page-wrapper">
    <div class="auth-page">
      <div class="auth-card">
      <div class="auth-header">
        <h2 class="auth-title">
          <img src="@/assets/icons/register.png" style="height: 36px; width: auto; vertical-align: middle; margin-right: 8px; object-fit: contain;" alt="注册" />
          注册
        </h2>
        <p class="auth-subtitle">创建新账户，开始使用系统</p>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="auth-form"
        @submit.prevent="handleRegister"
      >
        <div class="form-item">
          <label class="form-label">用户名</label>
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="3-20个字符"
            />
          </el-form-item>
        </div>

        <div class="form-item">
          <label class="form-label">邮箱地址</label>
          <el-form-item prop="email">
            <el-input
              v-model="form.email"
              type="email"
              placeholder="用于接收验证码"
              @blur="checkEmailAvailability"
            >
              <template #suffix>
                <span v-if="emailChecking" style="color: #1890ff; font-size: 12px; padding-right: 8px;">检查中...</span>
                <span v-else-if="emailCheckResult === 'available'" style="color: #52c41a; font-size: 12px; padding-right: 8px;">✓ 可用</span>
                <span v-else-if="emailCheckResult === 'taken'" style="color: #ff4d4f; font-size: 12px; padding-right: 8px;">已注册</span>
              </template>
            </el-input>
          </el-form-item>
        </div>

        <div class="form-item">
          <label class="form-label">邮箱验证码</label>
          <el-form-item prop="code">
            <div style="display: flex; gap: 12px;">
              <el-input
                v-model="form.code"
                placeholder="请输入验证码"
                style="flex: 1;"
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

        <div class="form-item">
          <label class="form-label">设置密码</label>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="至少6个字符"
              show-password
            />
          </el-form-item>
        </div>

        <div class="form-item">
          <label class="form-label">确认密码</label>
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="再次输入密码"
              show-password
              @keyup.enter="handleRegister"
            />
          </el-form-item>
        </div>

        <button type="submit" class="btn btn-primary" :disabled="userStore.loading">
          {{ userStore.loading ? '注册中...' : '立即注册' }}
        </button>
      </el-form>

      <div class="auth-footer">
        已有账户？<router-link to="/login">立即登录</router-link>
      </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { authApi } from '@/api'
import { showMessage, debounce } from '@/utils'
import type { FormInstance, FormRules } from 'element-plus'
import type { RegisterForm } from '@/types'
import Footer from '@/components/Footer.vue'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const form = reactive<RegisterForm>({
  username: '',
  email: '',
  code: '',
  password: '',
  confirmPassword: '',
})

const emailChecking = ref(false)
const emailCheckResult = ref<'available' | 'taken' | ''>('')
const lastCheckedEmail = ref('')

const validateUsername = (rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请输入用户名'))
  } else if (value.length < 3 || value.length > 20) {
    callback(new Error('用户名长度应在3-20个字符之间'))
  } else if (!/^[a-zA-Z0-9_]+$/.test(value)) {
    callback(new Error('用户名只能包含字母、数字和下划线'))
  } else {
    callback()
  }
}

const validatePassword = (rule: any, value: any, callback: any) => {
  if (!value) {
    callback(new Error('请输入密码'))
  } else if (value.length < 6) {
    callback(new Error('密码长度不能少于6位'))
  } else if (/\s/.test(value)) {
    callback(new Error('密码不能包含空格'))
  } else if (!/^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?`~]+$/.test(value)) {
    callback(new Error('密码包含违规字符，只允许字母、数字和常见符号'))
  } else {
    callback()
  }
}

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
  username: [
    { required: true, validator: validateUsername, trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码长度为6位', trigger: 'blur' },
  ],
  password: [
    { required: true, validator: validatePassword, trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

const checkEmailAvailability = async () => {
  if (!form.email) {
    emailCheckResult.value = ''
    lastCheckedEmail.value = ''
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.email)) {
    emailCheckResult.value = ''
    return
  }

  if (form.email === lastCheckedEmail.value && emailCheckResult.value) {
    return
  }

  emailChecking.value = true
  emailCheckResult.value = ''
  
  try {
    const response = await authApi.checkEmailExist({ email: form.email })
    if (response.code === 1) {
      emailCheckResult.value = 'available'
      lastCheckedEmail.value = form.email
    } else if (response.code === 1011) {
      emailCheckResult.value = 'taken'
      lastCheckedEmail.value = form.email
    }
  } catch (error: any) {
    emailCheckResult.value = ''
  } finally {
    emailChecking.value = false
  }
}

const sendingCode = ref(false)
const countdown = ref(0)
let countdownTimer: number | null = null

const handleSendCode = async () => {
  if (!form.email) {
    showMessage.warning('请先输入邮箱地址')
    return
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.email)) {
    showMessage.warning('请输入正确的邮箱格式')
    return
  }

  if (emailCheckResult.value === 'taken') {
    showMessage.warning('该邮箱已被注册，无法注册')
    return
  }

  if (emailCheckResult.value !== 'available') {
    await checkEmailAvailability()
    if (emailCheckResult.value !== 'available') {
      return
    }
  }

  sendingCode.value = true
  try {
    const response = await authApi.sendEmailCode({ email: form.email })
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

const handleRegister = async () => {
  if (!formRef.value) return

  // 检查两次密码是否一致
  if (form.password !== form.confirmPassword) {
    showMessage.error('两次输入的密码不一致')
    return
  }

  try {
    await formRef.value.validate()
    
    if (emailCheckResult.value === 'taken') {
      showMessage.warning('该邮箱已被注册，请使用其他邮箱')
      return
    }
    
    if (emailCheckResult.value !== 'available') {
      await checkEmailAvailability()
      if (emailCheckResult.value !== 'available') {
        return
      }
    }
    
    const success = await userStore.register(form)
    if (success) {
      setTimeout(() => {
        router.push('/login')
      }, 1500)
    }
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const debouncedCheckEmail = debounce(checkEmailAvailability, 800)
watch(() => form.email, (newEmail, oldEmail) => {
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

