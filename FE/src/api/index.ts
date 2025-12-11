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
      // 检查是否是被踢出（单设备登录限制）
      if (errorData?.code === 1026) {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        router.push('/login')
        showMessage.warning('您的账号已在其他设备登录，请重新登录')
        return Promise.reject(error)
      }

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

    // 429 请求限流
    if (error.response?.status === 429) {
      const retryAfter = errorData?.retryAfter
      const limitMsg = retryAfter 
        ? `请求过于频繁，请 ${retryAfter} 秒后再试`
        : errorMessage
      showMessage.warning(limitMsg)
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

// 验证码响应类型
export interface CaptchaResponse {
  code: number
  captcha_id: string
  captcha_img: string
}

// 验证码API
export const captchaApi = {
  // 获取图形验证码
  getCaptcha: (): Promise<CaptchaResponse> =>
    api.get('/captcha'),
}

// 用户认证相关API
export const authApi = {
  // 账户密码登录（需要验证码）
  login: (data: LoginForm & { captcha_id: string; captcha_code: string }): Promise<ApiResponse> => 
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
  
  // 登出
  logout: (): Promise<ApiResponse> => 
    api.post('/user/logout'),
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

// 涨跌幅数据类型
export interface PriceIncreaseItem {
  marketHashName: string
  name: string
  iconUrl: string
  platform: string
  todayPrice: number
  yesterdayPrice: number
  price3DaysAgo: number | null
  price7DaysAgo: number | null
  price15DaysAgo: number | null
  price30DaysAgo: number | null
  priceChange: number
  increaseRate1D: number
  increaseRate3D: number | null
  increaseRate7D: number | null
  increaseRate15D: number | null
  increaseRate30D: number | null
  // 在售数相关字段
  todaySellCount: number
  yesterdaySellCount: number
  sellCount3DaysAgo: number | null
  sellCount7DaysAgo: number | null
  sellCount15DaysAgo: number | null
  sellCount30DaysAgo: number | null
  sellCountChange: number
  sellCountRate1D: number
  sellCountRate3D: number | null
  sellCountRate7D: number | null
  sellCountRate15D: number | null
  sellCountRate30D: number | null
}

// 价格历史数据项
export interface PriceHistoryItem {
  date: string
  sellPrice: number
  sellCount: number
}

// 平台在售信息
export interface GoodsPlatformInfo {
  platform: string
  platformName: string
  sellPrice: number
  sellCount: number
  biddingPrice: number
  biddingCount: number
  link: string
}

// 商品详情响应
// 价格涨幅项
export interface PriceChangeItem {
  label: string      // 今日、本周、本月
  priceDiff: number  // 价格差
  changeRate: number // 涨跌幅百分比
  isUp: boolean      // 是否上涨
}

export interface GoodsDetailResponse {
  marketHashName: string
  name: string
  iconUrl: string
  rarityName: string    // 品质名称
  qualityName: string   // 类别名称
  priceHistory: Record<string, PriceHistoryItem[]>  // key: 平台名，value: 历史数据数组
  platformList: GoodsPlatformInfo[]
  priceChange: PriceChangeItem[]  // 悠悠平台的涨幅信息
}

// 搬砖数据项类型
export interface BrickMovingItem {
  id: number
  market_hash_name: string
  name: string
  image_url: string
  target_price: number
  source_price: number
  price_diff: number
  profit_rate: number
  sell_count: number
  turn_over: number
  bidding_count: number
  bidding_price: number
  target_update_time: string
  source_update_time: string
  platform_list: Array<{
    platform: string
    platformName: string
    sellPrice: number
    sellCount: number
    biddingPrice: number
    biddingCount: number
    link: string
    priceDiff: number
  }>
}

// 公开首页数据响应
export interface PublicHomeData {
  rankingList: PriceIncreaseItem[]
  brickMoving: BrickMovingItem[]
}

// 公开API（无需登录）
export const publicApi = {
  // 获取公开首页数据
  getHomeData: (): Promise<ApiResponse<PublicHomeData>> =>
    api.get('/public/home'),
}

// 数据相关API（VIP用户）
export const dataApi = {
  // 获取饰品数据
  getSkinItems: (params: PaginationParams): Promise<ApiResponse<SkinItem[]>> => 
    api.get('/vip/goods/data', { params }),
  
  // 获取分类列表
  getCategories: (): Promise<ApiResponse<string[]>> => 
    api.get('/vip/goods/category'),
  
  // 获取涨跌幅排行
  getPriceIncrease: (params: { is_desc: boolean; limit: number }): Promise<ApiResponse<PriceIncreaseItem[]>> =>
    api.get('/vip/goods/price-increase', { params }),
  
  // 获取商品详情
  getGoodsDetail: (params: { market_hash_name: string; days?: number }): Promise<ApiResponse<GoodsDetailResponse>> =>
    api.get('/vip/goods/detail', { params }),
  
  // 搜索商品
  searchGoods: (params: { keyword: string; limit?: number }): Promise<ApiResponse<SearchResult[]>> =>
    api.get('/vip/goods/search', { params }),
}

// 搜索结果类型
export interface SearchResult {
  name: string
  marketHashName: string
  iconUrl: string
}

export default api
