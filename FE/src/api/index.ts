import axios from 'axios'
import { ElMessage } from 'element-plus'
import type { 
  UUToken, 
  BuffToken, 
  Settings, 
  PaginationParams, 
  ApiResponse,
  TokenStatus,
  SkinItem
} from '@/types'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const { data } = response
    if (data.code === 0) {
      ElMessage.error(data.msg || '请求失败')
      return Promise.reject(new Error(data.msg || '请求失败'))
    }
    return data
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

// Token相关API
export const tokenApi = {
  // 更新UU Token
  updateUUToken: (data: UUToken): Promise<ApiResponse> => 
    api.post('/tokens/uu', data),
  
  // 更新Buff Token
  updateBuffToken: (data: BuffToken): Promise<ApiResponse> => 
    api.post('/tokens/buff', data),
  
  // 验证Token状态（获取状态用）
  verifyTokens: (): Promise<ApiResponse<TokenStatus>> => 
    api.get('/tokens/verify'),
    
  // 手动验证所有Token（验证按钮用）
  manualVerifyTokens: (): Promise<ApiResponse<TokenStatus>> => 
    api.post('/tokens/verify'),
}

// 设置相关API
export const settingsApi = {
  // 获取设置
  getSettings: (): Promise<ApiResponse<Settings>> => 
    api.get('/settings'),
  
  // 更新设置
  updateSettings: (data: Settings): Promise<ApiResponse> => 
    api.put('/settings', data),
}

// 数据相关API
export const dataApi = {
  // 获取饰品数据
  getSkinItems: (params: PaginationParams): Promise<ApiResponse<SkinItem[]>> => 
    api.get('/data', { params }),
}

export default api
