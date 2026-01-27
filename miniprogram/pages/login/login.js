// pages/login/login.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    loading: false,
    agreed: false  // 是否同意协议
  },

  onLoad() {
    // 检查是否已登录
    if (app.globalData.token) {
      wx.switchTab({
        url: '/pages/home/home'
      })
    }
  },

  // 先逛逛（跳过登录）
  skip() {
    wx.switchTab({
      url: '/pages/home/home'
    })
  },

  // 跳转用户协议
  goUserAgreement() {
    wx.navigateTo({
      url: '/pages/user-agreement/user-agreement'
    })
  },

  // 跳转隐私政策
  goPrivacyPolicy() {
    wx.navigateTo({
      url: '/pages/privacy-policy/privacy-policy'
    })
  },

  // 切换同意状态
  toggleAgree() {
    this.setData({
      agreed: !this.data.agreed
    })
  },

  // 显示协议弹窗
  showAgreementModal() {
    const that = this
    wx.showModal({
      title: '温馨提示',
      content: '请阅读并同意《用户协议》和《隐私政策》',
      cancelText: '不同意',
      confirmText: '同意',
      success(res) {
        if (res.confirm) {
          // 用户点击同意
          that.setData({ agreed: true })
          // 继续登录
          that.doLogin()
        }
      }
    })
  },

  // 微信登录
  async handleWechatLogin() {
    if (this.data.loading) return

    // 检查是否同意协议
    if (!this.data.agreed) {
      this.showAgreementModal()
      return
    }

    this.doLogin()
  },

  // 执行登录
  async doLogin() {
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
        
        // 重置引导状态，让首页显示引导弹窗
        app.globalData.hasShownGuide = false
        
        wx.showToast({
          title: '登录成功',
          icon: 'success'
        })

        // 跳转到首页，首页会处理VIP/邮箱绑定引导
        setTimeout(() => {
          wx.switchTab({
            url: '/pages/home/home'
          })
        }, 1000)
      }
    } catch (error) {
      console.error('登录失败:', error)
    } finally {
      this.setData({ loading: false })
    }
  }
})

