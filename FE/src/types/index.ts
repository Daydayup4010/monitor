// API响应类型
export interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
  total?: number
  token?: string
  user?: UserInfo
}

// 用户相关类型
export interface UserInfo {
  id: string
  username: string
  email: string
  role: number  // 0: 普通用户, 1: VIP, 2: 管理员
  vip_expiry?: string
  last_login?: string
  created_at?: string
}

export interface LoginForm {
  email: string
  password: string
}

export interface EmailLoginForm {
  email: string
  code: string
}

export interface RegisterForm {
  username: string
  email: string
  password: string
  confirmPassword: string
  code: string
}

export interface SendEmailCodeRequest {
  email: string
}

export interface ResetPasswordForm {
  email: string
  password: string
  code: string
}

// 用户列表相关
export interface UserListParams {
  page_num: number
  page_size: number
  search?: string
}

export interface UserListItem {
  id: string
  user_name: string
  email: string
  role: number
  vip_expiry?: string
  last_login?: string
}

export interface RenewVipRequest {
  user_id: string
  days: number
}

// 饰品数据类型
export interface SkinItem {
  id: number
  name: string
  category: string
  image_url: string
  source_price: number    // 买入价格（根据source平台）
  target_price: number    // 卖出价格（根据target平台）
  u_price?: number        // 兼容旧数据
  buff_price?: number     // 兼容旧数据
  price_diff: number
  profit_rate: number
  updated_at: string
}

// 分页参数
export interface PaginationParams {
  page_num: number
  page_size: number
  sort?: string  // 排序字段
  desc?: boolean // 是否降序
  user_id?: string
  category?: string  // 分类筛选
  search?: string    // 搜索关键词
  source?: string    // 买入平台
  target?: string    // 卖出平台
}

// Platform Token类型
export interface UUToken {
  authorization: string
  uk: string
}

export interface BuffToken {
  session: string
  csrf_token: string
}

export interface TokenStatus {
  uu: 'yes' | 'no'  // no表示有效，yes表示无效
  buff: 'yes' | 'no'
}

// 设置类型
export interface Settings {
  user_id?: string
  min_sell_num: number
  min_diff: number
  max_sell_price: number
  min_sell_price: number
}

// 统计数据类型
export interface Statistics {
  totalItems: number
  avgProfitRate: number
  maxPriceDiff: number
  lastUpdateTime: string
}
