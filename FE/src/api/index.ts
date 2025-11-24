import axios from 'axios'
import router from '@/router'
import { getErrorMessage, translateErrorMessage } from '@/utils/errorCode'
import { showMessage } from '@/utils/message'
import type { 
  Settings, 
  PaginationParams, 
  ApiResponse,
  SkinItem,
  LoginForm,
  EmailLoginForm,
  RegisterForm,
  SendEmailCodeRequest,
  ResetPasswordForm,
  UserInfo,
  UserListParams,
  UserListItem,
  RenewVipRequest
} from '@/types'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// Token刷新标志
let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

// 添加订阅者
function subscribeTokenRefresh(cb: (token: string) => void) {
  refreshSubscribers.push(cb)
}

// 通知所有订阅者
function onRefreshed(token: string) {
  refreshSubscribers.forEach(cb => cb(token))
  refreshSubscribers = []
}

// 请求拦截器 - 添加Token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理Token过期和自动刷新
api.interceptors.response.use(
  (response) => {
    const { data } = response
    const config = response.config
    
    // 邮箱校验接口的错误不显示消息，且不reject（通过UI状态显示）
    const isEmailCheckEndpoint = config.url?.includes('/user/email-exist')
    
    // 如果code不是1（成功），根据错误码或消息显示中文错误
    if (data.code !== 1 && data.code !== undefined) {
      const errorMsg = getErrorMessage(data.code, data.msg)
      const translatedMsg = translateErrorMessage(errorMsg)
      
      // 邮箱校验接口特殊处理：不显示错误，不reject，直接返回数据
      if (isEmailCheckEndpoint) {
        return data
      }
      
      // 其他接口显示错误并reject
      showMessage.error(translatedMsg)
      return Promise.reject(new Error(translatedMsg))
    }
    return data
  },
  async (error) => {
    const originalRequest = error.config
    const errorData = error.response?.data

    // 获取错误消息（优先使用错误码映射，其次翻译英文消息）
    let errorMessage = '网络错误'
    if (errorData) {
      if (errorData.code !== undefined) {
        errorMessage = getErrorMessage(errorData.code, errorData.msg)
      } else if (errorData.msg) {
        errorMessage = errorData.msg
      } else if (errorData.error) {
        errorMessage = errorData.error
      }
      // 翻译成中文
      errorMessage = translateErrorMessage(errorMessage)
    }

    // 401错误处理
    if (error.response?.status === 401) {
      // 检查是否是登录相关的接口（登录、注册等不需要Token的接口）
      const isAuthEndpoint = originalRequest.url?.includes('/user/login') || 
                            originalRequest.url?.includes('/user/register') ||
                            originalRequest.url?.includes('/user/send-email') ||
                            originalRequest.url?.includes('/user/reset-password')
      
      // 如果是登录接口返回401，直接显示错误，不尝试刷新Token
      if (isAuthEndpoint) {
        showMessage.error(errorMessage)
        return Promise.reject(error)
      }
      
      // 其他接口的401，尝试刷新Token
      if (!originalRequest._retry) {
        if (isRefreshing) {
          // 如果正在刷新，将请求加入队列
          return new Promise(resolve => {
            subscribeTokenRefresh(token => {
              originalRequest.headers.Authorization = `Bearer ${token}`
              resolve(api(originalRequest))
            })
          })
        }

        originalRequest._retry = true
        isRefreshing = true

        try {
          // 调用刷新接口
          const { data } = await axios.post('/api/v1/user/refresh', {}, {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`
            }
          })
          
          if (data.code === 1 && data.token) {
            const newToken = data.token
            
            // 保存新Token
            localStorage.setItem('token', newToken)
            
            // 更新原请求的Token
            originalRequest.headers.Authorization = `Bearer ${newToken}`
            
            // 通知所有等待的请求
            onRefreshed(newToken)
            
            return api(originalRequest)
          } else {
            throw new Error('刷新Token失败')
          }
        } catch (err) {
          // 刷新失败，清除token并跳转到登录页
          localStorage.removeItem('token')
          localStorage.removeItem('userInfo')
          router.push('/login')
          showMessage.warning('登录已过期，请重新登录')
          return Promise.reject(err)
        } finally {
          isRefreshing = false
        }
      }
    }

    // 403权限不足
    if (error.response?.status === 403) {
      showMessage.error(errorMessage || '权限不足')
      return Promise.reject(error)
    }

    // VIP权限不足（特殊code标识）
    if (errorData?.code === 'VIP_REQUIRED') {
      showMessage.warning('需要VIP权限才能访问')
      return Promise.reject(error)
    }

    // 其他错误，只显示有效的错误消息
    if (errorMessage && errorMessage !== '网络错误') {
      showMessage.error(errorMessage)
    }
    return Promise.reject(error)
  }
)

// 用户认证相关API
export const authApi = {
  // 账户密码登录
  login: (data: LoginForm): Promise<ApiResponse> => 
    api.post('/user/login', data),
  
  // 邮箱验证码登录
  emailLogin: (data: EmailLoginForm): Promise<ApiResponse> => 
    api.post('/user/email-login', data),
  
  // 注册
  register: (data: Omit<RegisterForm, 'confirmPassword'>): Promise<ApiResponse> => 
    api.post('/user/register', data),
  
  // 发送邮箱验证码
  sendEmailCode: (data: SendEmailCodeRequest): Promise<ApiResponse> => 
    api.post('/user/send-email', data),
  
  // 重置密码
  resetPassword: (data: ResetPasswordForm): Promise<ApiResponse> => 
    api.post('/user/reset-password', data),
  
  // 检查邮箱是否存在
  checkEmailExist: (data: SendEmailCodeRequest): Promise<ApiResponse> => 
    api.post('/user/email-exist', data),
  
  // 获取当前用户信息
  getSelfInfo: (): Promise<ApiResponse<UserInfo>> => 
    api.get('/user/self'),
  
  // 刷新Token
  refreshToken: (): Promise<ApiResponse> => 
    api.post('/user/refresh'),
}

// 用户管理API（管理员）
export const userApi = {
  // 获取用户列表
  getUserList: (params: UserListParams): Promise<ApiResponse<UserListItem[]>> => 
    api.get('/admin/users', { params }),
  
  // 删除用户
  deleteUser: (userId: string): Promise<ApiResponse> => 
    api.delete('/admin/user', { params: { user_id: userId } }),
  
  // VIP续费
  renewVip: (data: RenewVipRequest): Promise<ApiResponse> => 
    api.post('/admin/vip-expiry', data),
}

// 设置相关API（VIP用户）
export const settingsApi = {
  // 获取设置
  getSettings: (): Promise<ApiResponse<Settings>> => 
    api.get('/vip/settings'),
  
  // 更新设置
  updateSettings: (data: Settings): Promise<ApiResponse> => 
    api.put('/vip/settings', data),
}

// 数据相关API（VIP用户）
export const dataApi = {
  // 获取饰品数据
  getSkinItems: (params: PaginationParams): Promise<ApiResponse<SkinItem[]>> => 
    api.get('/vip/goods/data', { params }),
  
  // 获取分类列表
  getCategories: (): Promise<ApiResponse<string[]>> => 
    api.get('/vip/goods/category'),
}

export default api
