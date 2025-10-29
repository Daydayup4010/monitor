// API接口定义
const request = require('./request.js')

module.exports = {
  // 微信登录
  wechatLogin: (code) => request.post('/wechat/login', { code }),
  
  // 绑定邮箱
  bindEmail: (email, code, password) => request.post('/wechat/bind-email', { email, code, password }),
  
  // 发送邮箱验证码
  sendEmailCode: (email) => request.post('/wechat/send-email-code', { email }),
  
  // 获取饰品数据
  getSkinItems: (params) => request.get('/vip/goods/data', params),
  
  // 获取分类列表
  getCategories: () => request.get('/vip/goods/category'),
  
  // 获取个人设置
  getSettings: () => request.get('/vip/settings'),
  
  // 更新个人设置
  updateSettings: (data) => request.put('/vip/settings', data),
  
  // 获取用户信息
  getUserInfo: () => request.get('/user/self'),
}

