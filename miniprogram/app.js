// app.js
App({
  onLaunch() {
    // 小程序启动时执行
    console.log('小程序启动')
    
    // 检查登录状态
    this.checkLoginStatus()
  },

  globalData: {
    userInfo: null,
    token: null,
    baseURL: 'http://localhost:3601/api/v1'  // 后端API地址
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

