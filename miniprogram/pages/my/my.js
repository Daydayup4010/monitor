// pages/my/my.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    userInfo: {},
    isLoggedIn: false,
    vipEnabled: false,  // VIP入口开关
    vipExpiryText: '',  // VIP到期时间文本
    isVipValid: false   // VIP是否有效
  },

  async onShow() {
    // 从storage重新读取用户信息，确保是最新的
    app.checkLoginStatus()
    
    const isLoggedIn = !!app.globalData.token
    
    // 先用本地数据渲染
    this.updateUserDisplay()

    // 如果已登录，主动从服务器刷新用户信息
    if (isLoggedIn) {
      this.refreshUserInfo()
    }

    // 每次进入页面都重新获取配置（确保配置是最新的）
    await app.fetchMinAppConfig()
    
    // 重新获取最新的vipEnabled并更新显示
    this.updateUserDisplay()
    
    const vipEnabled = app.globalData.minAppConfig?.vipEnabled || false
    // 检查是否需要引导（仅VIP开关开启时才引导）
    if (isLoggedIn && vipEnabled) {
      this.checkUserGuide()
    }
  },

  // 更新用户信息显示
  updateUserDisplay() {
    const isLoggedIn = !!app.globalData.token
    const vipEnabled = app.globalData.minAppConfig?.vipEnabled || false
    const userInfo = app.globalData.userInfo || {}
    
    // 计算VIP状态和到期时间
    let vipExpiryText = ''
    let isVipValid = false
    if (userInfo.role === 1 && userInfo.vip_expiry) {
      const vipExpiry = new Date(userInfo.vip_expiry)
      isVipValid = vipExpiry > new Date()
      vipExpiryText = this.formatDate(userInfo.vip_expiry)
    }
    
    this.setData({
      isLoggedIn: isLoggedIn,
      userInfo: userInfo,
      vipEnabled: vipEnabled,
      vipExpiryText: vipExpiryText,
      isVipValid: isVipValid
    })
  },

  // 从服务器刷新用户信息
  async refreshUserInfo() {
    try {
      const res = await api.getUserInfo()
      if (res.code === 1 && res.data) {
        // 更新全局和本地存储
        app.globalData.userInfo = res.data
        wx.setStorageSync('userInfo', res.data)
        
        // 更新页面显示
        this.updateUserDisplay()
      }
    } catch (error) {
      console.error('刷新用户信息失败:', error)
    }
  },

  // 检查用户引导（VIP和邮箱绑定）
  checkUserGuide() {
    // 如果已经显示过引导，不再重复
    if (app.globalData.hasShownGuide) return

    const userInfo = app.globalData.userInfo
    if (!userInfo) return

    // 判断VIP状态：role为1且vip_expiry未过期
    let isVip = false
    if (userInfo.role === 1 && userInfo.vip_expiry) {
      const vipExpiry = new Date(userInfo.vip_expiry)
      isVip = vipExpiry > new Date()
    }
    const hasEmail = userInfo.email && userInfo.email !== ''

    // 如果不是VIP
    if (!isVip) {
      app.globalData.hasShownGuide = true
      if (hasEmail) {
        // 已绑定邮箱，引导开通VIP
        this.showVipGuide()
      } else {
        // 未绑定邮箱，提示先绑定邮箱
        this.showBindEmailGuide()
      }
    }
  },

  // 显示VIP引导
  showVipGuide() {
    wx.showModal({
      title: '开通VIP',
      content: '开通VIP会员，解锁全部数据和功能',
      confirmText: '去开通',
      cancelText: '稍后再说',
      success: (res) => {
        if (res.confirm) {
          wx.navigateTo({
            url: '/pages/vip/vip'
          })
        }
      }
    })
  },

  // 显示绑定邮箱引导
  showBindEmailGuide() {
    wx.showModal({
      title: '绑定邮箱',
      content: '绑定邮箱后可在网页端登录使用更多功能',
      confirmText: '去绑定',
      cancelText: '稍后再说',
      success: (res) => {
        if (res.confirm) {
          wx.navigateTo({
            url: '/pages/bind-email/bind-email'
          })
        }
      }
    })
  },

  // 跳转登录
  goLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    })
  },

  goSettings() {
    wx.navigateTo({
      url: '/pages/settings/settings'
    })
  },

  goBindEmail() {
    wx.navigateTo({
      url: '/pages/bind-email/bind-email'
    })
  },

  goVip() {
    // 检查VIP入口开关
    if (!app.globalData.minAppConfig?.vipEnabled) {
      wx.showToast({
        title: '暂不支持小程序开通',
        icon: 'none'
      })
      return
    }
    wx.navigateTo({
      url: '/pages/vip/vip'
    })
  },

  logout() {
    wx.showModal({
      title: '提示',
      content: '确定要退出登录吗？',
      success(res) {
        if (res.confirm) {
          app.clearLoginInfo()
          wx.reLaunch({
            url: '/pages/login/login'
          })
        }
      }
    })
  },

  // 格式化日期
  formatDate(dateStr) {
    if (!dateStr) return ''
    const date = new Date(dateStr)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }
})

