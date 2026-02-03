// app.js
const api = require('./utils/api.js')

App({
  onLaunch() {
    // 小程序启动时执行
    console.log('小程序启动')
    
    // 检查登录状态
    this.checkLoginStatus()
    
    // 获取小程序配置
    this.fetchMinAppConfig()
  },

  onShow(options) {
    // 处理嵌入式小程序支付回调
    console.log('App.onShow options:', options)
    
    if (options && options.referrerInfo && options.referrerInfo.extraData) {
      const extraData = options.referrerInfo.extraData
      console.log('收到支付回调 extraData:', JSON.stringify(extraData))
      
      // 保存支付结果到 globalData
      // YunGouOS 返回: { code: 0, data: { orderNo: "xxx" }, msg: "支付成功" }
      this.globalData.payResult = {
        success: extraData.code === 0,
        orderNo: extraData.data?.orderNo || '',
        msg: extraData.msg || ''
      }
      
      console.log('设置 payResult:', this.globalData.payResult)
      
      // 如果支付失败，显示错误提示
      if (extraData.code !== 0) {
        wx.showToast({
          title: extraData.msg || '支付失败',
          icon: 'none',
          duration: 3000
        })
      }
    }
  },

  globalData: {
    userInfo: null,
    token: null,
    // baseURL: 'http://localhost:3100/api/v1',  // 后端API地址（本地开发）
    baseURL: 'https://www.csgoods.com.cn/api/v1',  // 后端API地址（线上）
    // 小程序配置
    minAppConfig: {
      vipEnabled: false  // VIP开通入口开关，默认关闭
    },
    configLoaded: false  // 配置是否已加载
  },

  // 获取小程序配置
  async fetchMinAppConfig() {
    try {
      const res = await api.getMinAppConfig()
      if (res.code === 1 && res.data) {
        this.globalData.minAppConfig = {
          vipEnabled: res.data.vip_enabled || false
        }
      }
    } catch (error) {
      console.error('获取小程序配置失败:', error)
    } finally {
      this.globalData.configLoaded = true
    }
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    const userInfo = wx.getStorageSync('userInfo')
    
    if (token && userInfo) {
      this.globalData.token = token
      this.globalData.userInfo = userInfo
    }
  },

  // 保存登录信息
  saveLoginInfo(token, userInfo) {
    this.globalData.token = token
    this.globalData.userInfo = userInfo
    wx.setStorageSync('token', token)
    wx.setStorageSync('userInfo', userInfo)
  },

  // 清除登录信息
  clearLoginInfo() {
    this.globalData.token = null
    this.globalData.userInfo = null
    wx.removeStorageSync('token')
    wx.removeStorageSync('userInfo')
  }
})

