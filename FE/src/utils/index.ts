import dayjs from 'dayjs'

// 格式化时间
export const formatTime = (time: string, format = 'YYYY-MM-DD HH:mm:ss'): string => {
  if (!time) return '--'
  return dayjs(time).format(format)
}

// 格式化价格
export const formatPrice = (price: number): string => {
  if (!price && price !== 0) return '--'
  return price.toFixed(2)
}

// 格式化百分比
export const formatPercent = (rate: number): string => {
  if (!rate && rate !== 0) return '--'
  return `${(rate * 100).toFixed(2)}%`
}

// 获取利润率颜色
export const getProfitColor = (rate: number): string => {
  if (rate > 0.2) return '#52c41a'  // 绿色
  if (rate > 0.1) return '#faad14'  // 橙色
  if (rate > 0.05) return '#1890ff' // 蓝色
  return '#f5222d' // 红色
}

// 防抖函数
export const debounce = (fn: Function, delay: number) => {
  let timeoutId: NodeJS.Timeout
  return (...args: any[]) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn.apply(null, args), delay)
  }
}

// 导出错误码相关工具
export { getErrorMessage, translateErrorMessage, ERROR_CODE_MAP, ERROR_MSG_MAP } from './errorCode'

// 导出美化的消息提示
export { showMessage } from './message'