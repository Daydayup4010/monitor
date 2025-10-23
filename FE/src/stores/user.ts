import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import { showMessage } from '@/utils/message'
import type { UserInfo, LoginForm, EmailLoginForm, RegisterForm } from '@/types'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)
  const token = ref<string>('')
  const loading = ref(false)

  // 用户角色常量
  const ROLE_NORMAL = 0
  const ROLE_VIP = 1
  const ROLE_ADMIN = 2

  // 是否已登录
  const isLoggedIn = computed(() => !!token.value && !!userInfo.value)

  // 是否是VIP
  const isVip = computed(() => {
    if (!userInfo.value) return false
    if (userInfo.value.role === ROLE_ADMIN) return true  // 管理员也算VIP
    if (userInfo.value.role === ROLE_VIP && userInfo.value.vip_expiry) {
      return new Date(userInfo.value.vip_expiry) > new Date()
    }
    return false
  })

  // 是否是管理员
  const isAdmin = computed(() => {
    return userInfo.value?.role === ROLE_ADMIN
  })

  // 用户类型标签
  const userTypeLabel = computed(() => {
    if (!userInfo.value) return '未登录'
    if (userInfo.value.role === ROLE_ADMIN) return '管理员'
    if (isVip.value) return 'VIP会员'
    return '普通用户'
  })

  // 从本地存储加载用户信息
  const loadFromStorage = () => {
    const storedToken = localStorage.getItem('token')
    const storedUserInfo = localStorage.getItem('userInfo')
    
    if (storedToken && storedUserInfo) {
      try {
        token.value = storedToken
        userInfo.value = JSON.parse(storedUserInfo)
      } catch (e) {
        console.error('解析用户信息失败:', e)
        logout()
      }
    }
  }

  // 保存到本地存储
  const saveToStorage = (newToken: string, newUserInfo: UserInfo) => {
    token.value = newToken
    userInfo.value = newUserInfo
    localStorage.setItem('token', newToken)
    localStorage.setItem('userInfo', JSON.stringify(newUserInfo))
    
    // 开发环境下输出调试信息
    if (import.meta.env.DEV) {
      console.log('保存用户信息:', {
        用户名: newUserInfo.username,
        角色: newUserInfo.role,
        VIP到期: newUserInfo.vip_expiry,
        是否VIP: isVip.value,
        是否管理员: isAdmin.value
      })
    }
  }

  // 账户密码登录
  const login = async (form: LoginForm) => {
    loading.value = true
    try {
      const response = await authApi.login(form)
      if (response.code === 1 && response.token && response.data) {
        saveToStorage(response.token, response.data)
        showMessage.success('登录成功')
        return true
      }
      return false
    } catch (error: any) {
      // 错误已在拦截器中处理，这里不再重复提示
      return false
    } finally {
      loading.value = false
    }
  }

  // 邮箱验证码登录
  const emailLogin = async (form: EmailLoginForm) => {
    loading.value = true
    try {
      const response = await authApi.emailLogin(form)
      if (response.code === 1 && response.token && response.data) {
        saveToStorage(response.token, response.data)
        showMessage.success('登录成功')
        return true
      }
      return false
    } catch (error: any) {
      // 错误已在拦截器中处理，这里不再重复提示
      return false
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (form: RegisterForm) => {
    loading.value = true
    try {
      const { confirmPassword, ...registerData } = form
      const response = await authApi.register(registerData)
      if (response.code === 1) {
        showMessage.success('注册成功，请登录')
        return true
      }
      return false
    } catch (error: any) {
      // 错误已在拦截器中处理，这里不再重复提示
      return false
    } finally {
      loading.value = false
    }
  }

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const response = await authApi.getSelfInfo()
      if (response.code === 1 && response.data) {
        userInfo.value = response.data
        localStorage.setItem('userInfo', JSON.stringify(response.data))
        return true
      }
      return false
    } catch (error) {
      console.error('获取用户信息失败:', error)
      return false
    }
  }

  // 刷新Token
  const refreshToken = async () => {
    try {
      const response = await authApi.refreshToken()
      if (response.code === 1 && response.token) {
        const userData = response.data
        saveToStorage(response.token, userData)
        return true
      }
      return false
    } catch (error) {
      console.error('刷新Token失败:', error)
      return false
    }
  }

  // 登出
  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    router.push('/login')
    showMessage.info('已退出登录')
  }

  // 检查VIP状态
  const checkVipStatus = (): boolean => {
    if (!isLoggedIn.value) {
      showMessage.warning('请先登录')
      router.push('/login')
      return false
    }
    
    if (!isVip.value) {
      showMessage.warning('需要VIP权限才能访问')
      return false
    }
    
    return true
  }

  // 检查管理员权限
  const checkAdminPermission = (): boolean => {
    if (!isLoggedIn.value) {
      showMessage.warning('请先登录')
      router.push('/login')
      return false
    }
    
    if (!isAdmin.value) {
      showMessage.error('需要管理员权限')
      return false
    }
    
    return true
  }

  return {
    // 状态
    userInfo,
    token,
    loading,
    // 计算属性
    isLoggedIn,
    isVip,
    isAdmin,
    userTypeLabel,
    // 方法
    loadFromStorage,
    login,
    emailLogin,
    register,
    getUserInfo,
    refreshToken,
    logout,
    checkVipStatus,
    checkAdminPermission,
  }
})

