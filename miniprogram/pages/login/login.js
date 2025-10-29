// pages/login/login.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    loading: false
  },

  onLoad() {
    // 检查是否已登录
    if (app.globalData.token) {
      wx.switchTab({
        url: '/pages/home/home'
      })
    }
  },

  // 微信登录
  async handleWechatLogin() {
    if (this.data.loading) return
    
    this.setData({ loading: true })

    try {
      // 1. 获取微信登录code
      const loginRes = await wx.login()
      
      if (!loginRes.code) {
        wx.showToast({
          title: '获取登录凭证失败',
          icon: 'none'
        })
        return
      }

      // 2. 调用后端API登录
      const res = await api.wechatLogin(loginRes.code)
      
      if (res.code === 1) {
        // 保存登录信息
        app.saveLoginInfo(res.token, res.data)
        
        wx.showToast({
          title: '登录成功',
          icon: 'success'
        })

        // 3. 检查是否需要绑定邮箱
        if (!res.data.has_email) {
          setTimeout(() => {
            wx.showModal({
              title: '提示',
              content: '建议绑定邮箱，可以在Web端登录查看更多数据',
              confirmText: '去绑定',
              cancelText: '暂不绑定',
              success(modal) {
                if (modal.confirm) {
                  wx.navigateTo({
                    url: '/pages/bind-email/bind-email'
                  })
                } else {
                  wx.switchTab({
                    url: '/pages/home/home'
                  })
                }
              }
            })
          }, 1000)
        } else {
          // 4. 跳转到首页
          setTimeout(() => {
            wx.switchTab({
              url: '/pages/home/home'
            })
          }, 1000)
        }
      }
    } catch (error) {
      console.error('登录失败:', error)
    } finally {
      this.setData({ loading: false })
    }
  }
})

