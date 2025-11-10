// pages/my/my.js
const app = getApp()

Page({
  data: {
    userInfo: {}
  },

  onShow() {
    this.setData({
      userInfo: app.globalData.userInfo || {}
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
  }
})

