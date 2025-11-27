// 错误码到中文消息的映射
export const ERROR_CODE_MAP: Record<number, string> = {
  // 通用错误
  0: '参数错误',
  1: '成功',
  
  // 用户模块错误 (1000-1999)
  1001: '用户名已被使用',
  1002: '密码错误',
  1003: '用户不存在',
  1004: '未找到认证Token',
  1005: '认证Token已过期',
  1006: '无效的认证Token',
  1007: '认证Token格式错误',
  1008: '权限不足',
  1009: '生成Token失败',
  1010: '验证码错误或已过期',
  1011: '邮箱已被注册',
  1012: '注册失败',
  1013: '获取用户列表失败',
  1014: '生成验证码失败',
  1015: '发送验证码失败',
  1016: '更新用户信息失败',
  1017: '删除用户失败',
  1018: '微信登录失败',
  1019: '微信绑定失败',
  1020: '参数无效',
  1021: '邮箱已存在，需要合并账号',
  1022: '账号合并失败',
  1023: '该邮箱已绑定其他微信账号',
  1024: '生成图形验证码失败',
  1025: '图形验证码错误或已过期',
  
  // 系统模块错误 (2000-2999)
  2000: '全量更新正在运行中',
  2001: '获取设置失败',
  2002: '获取饰品总数失败',
  2003: '获取饰品数据失败',
  2004: '更新设置失败',
  2005: '更新UU Token失败',
  2006: '更新Buff Token失败',
  2007: '获取Token状态失败',
  2008: '创建默认设置失败',
  2009: '获取分类失败',
  2010: '请求过于频繁，请稍后再试',
}

// 获取错误消息（中文）
export function getErrorMessage(code: number, defaultMsg?: string): string {
  return ERROR_CODE_MAP[code] || defaultMsg || '未知错误'
}

// 英文错误消息到中文的映射（用于未使用错误码的旧接口）
export const ERROR_MSG_MAP: Record<string, string> = {
  'Invalid Parameter': '参数错误',
  'The provided password is incorrect': '密码错误',
  'The requested username is already in use': '用户名已被使用',
  'The requested email is already in use': '邮箱已被注册',
  'The provided email code is incorrect': '验证码错误或已过期',
  'User account not found': '用户不存在',
  'Authentication token not found': '未找到认证Token',
  'Authentication token has expired': '认证Token已过期',
  'Invalid authentication token': '无效的认证Token',
  'Malformed authentication token': '认证Token格式错误',
  'Insufficient permissions for this operation': '权限不足',
  'Generate token error': '生成Token失败',
  'Register user error': '注册失败',
  'Get user list error': '获取用户列表失败',
  'Email code generate error': '生成验证码失败',
  'Send email code error': '发送验证码失败',
  'Update user error': '更新用户信息失败',
  'Delete user error': '删除用户失败',
  'Authorization header required': '需要登录',
  'Invalid authorization format': '认证格式错误',
  'invalid token': 'Token无效',
  'vip access required': '需要VIP权限',
  'Admin access required': '需要管理员权限',
  'Full update running': '全量更新正在运行中',
  'Get settings error': '获取设置失败',
  'Get goods total error': '获取商品总数失败',
  'Get goods data error': '获取商品数据失败',
  'Update setting error': '更新设置失败',
  'Update UU token error': '更新UU Token失败',
  'Update buff token error': '更新Buff Token失败',
  'Get token expired error': '获取Token状态失败',
  'create default setting error': '创建默认设置失败',
  'success': '操作成功',
  'Generate captcha error': '生成图形验证码失败',
  'Invalid captcha or captcha expired': '图形验证码错误或已过期',
}

// 翻译错误消息为中文
export function translateErrorMessage(msg: string): string {
  // 先尝试精确匹配
  if (ERROR_MSG_MAP[msg]) {
    return ERROR_MSG_MAP[msg]
  }
  
  // 尝试模糊匹配（包含关系）
  for (const [enMsg, cnMsg] of Object.entries(ERROR_MSG_MAP)) {
    if (msg.toLowerCase().includes(enMsg.toLowerCase())) {
      return cnMsg
    }
  }
  
  // 如果没有匹配到，返回原消息
  return msg
}

