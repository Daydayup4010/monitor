import { ElMessage } from 'element-plus'
import type { MessageOptions } from 'element-plus'

// 自定义消息配置
const messageConfig: Partial<MessageOptions> = {
  duration: 3000,
  showClose: false,
  offset: 60,
}

// 美化的消息提示
export const showMessage = {
  success: (message: string, options?: Partial<MessageOptions>) => {
    ElMessage.success({
      message,
      ...messageConfig,
      customClass: 'custom-message custom-message-success',
      ...options,
    })
  },
  
  error: (message: string, options?: Partial<MessageOptions>) => {
    // 过滤掉不应该显示的错误消息
    if (!message || shouldIgnoreError(message)) {
      return
    }
    
    ElMessage.error({
      message,
      ...messageConfig,
      customClass: 'custom-message custom-message-error',
      ...options,
    })
  },
  
  warning: (message: string, options?: Partial<MessageOptions>) => {
    ElMessage.warning({
      message,
      ...messageConfig,
      customClass: 'custom-message custom-message-warning',
      ...options,
    })
  },
  
  info: (message: string, options?: Partial<MessageOptions>) => {
    ElMessage.info({
      message,
      ...messageConfig,
      customClass: 'custom-message custom-message-info',
      ...options,
    })
  },
}

// 判断是否应该忽略的错误消息
function shouldIgnoreError(message: string): boolean {
  const ignorePatterns = [
    /Request failed with status code/i,
    /Network Error/i,
    /timeout of \d+ms exceeded/i,
  ]
  
  return ignorePatterns.some(pattern => pattern.test(message))
}

