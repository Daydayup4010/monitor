// API接口定义
const request = require('./request.js')

module.exports = {
  // 微信登录
  wechatLogin: (code) => request.post('/wechat/login', { code }),
  
  // 绑定邮箱
  bindEmail: (email, code, password) => request.post('/wechat/bind-email', { email, code, password }),
  
  // 合并账号
  mergeAccount: (email, code) => request.post('/wechat/merge-account', { email, code }),
  
  // 发送邮箱验证码
  sendEmailCode: (email) => request.post('/wechat/send-email-code', { email }),
  
  // 验证邮箱验证码
  verifyEmailCode: (email, code) => request.post('/user/verify-email-code', { email, code }),
  
  // 获取饰品数据（VIP）
  getSkinItems: (params) => request.get('/vip/goods/data', params),
  
  // 获取公开首页数据（无需登录）
  getPublicHomeData: () => request.get('/public/home'),
  
  // 获取分类列表
  getCategories: () => request.get('/vip/goods/category'),
  
  // 获取个人设置
  getSettings: () => request.get('/vip/settings'),
  
  // 更新个人设置
  updateSettings: (data) => request.put('/vip/settings', data),
  
  // 获取用户信息
  getUserInfo: () => request.get('/user/self'),
  
  // 获取饰品详情
  getGoodsDetail: (params) => request.get('/vip/goods/detail', params),
  
  // VIP支付相关
  getVipPlans: () => request.get('/payment/vip-price'),
  createMinAppPay: (data) => request.post('/payment/minapp', data),
  queryPayOrder: (orderNo) => request.get('/payment/query', { order_no: orderNo }),
  getVipRecords: (params) => request.get('/payment/vip-records', params),
}

