import axios from 'axios'
import type { ApiResponse } from '@/types'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

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

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => Promise.reject(error)
)

// 通知相关类型
export interface NotificationItem {
  id: string
  title: string
  content: string
  image_url?: string  // 图片URL（可选）
  is_read: boolean
  created_at: string
}

// 通知相关API
export const notificationApi = {
  // 获取通知列表
  getList: (params: { page_num: number; page_size: number }): Promise<ApiResponse<NotificationItem[]>> =>
    api.get('/notification/list', { params }),
  
  // 获取未读数量
  getUnreadCount: (): Promise<ApiResponse<{ unread_count: number }>> =>
    api.get('/notification/unread-count'),
  
  // 标记单条已读
  markAsRead: (data: { notification_id: string }): Promise<ApiResponse> =>
    api.post('/notification/read', data),
  
  // 标记全部已读
  markAllAsRead: (): Promise<ApiResponse> =>
    api.post('/notification/read-all'),
}

// 通知管理API（管理员）
export const notificationAdminApi = {
  // 创建通知
  create: (data: { title: string; content: string; image_url?: string }): Promise<ApiResponse> =>
    api.post('/admin/notification', data),
  
  // 获取所有通知
  getAll: (params: { page_num: number; page_size: number }): Promise<ApiResponse<NotificationItem[]>> =>
    api.get('/admin/notifications', { params }),
  
  // 删除通知
  delete: (notificationId: string): Promise<ApiResponse> =>
    api.delete('/admin/notification', { params: { notification_id: notificationId } }),
  
  // 上传图片
  uploadImage: (file: File): Promise<ApiResponse<{ url: string }>> => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/admin/upload-image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
}
