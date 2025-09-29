// API响应类型
export interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
  total?: number
}

// 饰品数据类型
export interface SkinItem {
  id: number
  name: string
  category: string
  image_url: string
  u_price: number
  buff_price: number
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
}

// Token类型
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
