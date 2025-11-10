// pages/bind-email/bind-email.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    email: '',
    code: '',
    password: '',
    confirmPassword: '',
    countdown: 0,
    sending: false,
    submitting: false
  },

  onEmailInput(e) {
    this.setData({ email: e.detail.value })
  },

  onCodeInput(e) {
    this.setData({ code: e.detail.value })
  },

  onPasswordInput(e) {
    this.setData({ password: e.detail.value })
  },

  onConfirmPasswordInput(e) {
    this.setData({ confirmPassword: e.detail.value })
  },

  // 发送验证码
  async sendCode() {
    if (!this.data.email) {
      wx.showToast({
        title: '请输入邮箱地址',
        icon: 'none'
      })
      return
    }

    // 验证邮箱格式
    const emailReg = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailReg.test(this.data.email)) {
      wx.showToast({
        title: '邮箱格式不正确',
        icon: 'none'
      })
      return
    }

    this.setData({ sending: true })

    try {
      const res = await api.sendEmailCode(this.data.email)
      
      if (res.code === 1) {
        wx.showToast({
          title: '验证码已发送',
          icon: 'success'
        })

        // 开始倒计时
        this.setData({ countdown: 60 })
        const timer = setInterval(() => {
          if (this.data.countdown <= 1) {
            clearInterval(timer)
            this.setData({ countdown: 0 })
          } else {
            this.setData({ countdown: this.data.countdown - 1 })
          }
        }, 1000)
      }
    } catch (error) {
      console.error('发送验证码失败:', error)
    } finally {
      this.setData({ sending: false })
    }
  },

  // 提交绑定
  async handleSubmit() {
    // 验证表单
    if (!this.data.email || !this.data.code || !this.data.password || !this.data.confirmPassword) {
      wx.showToast({
        title: '请填写完整信息',
        icon: 'none'
      })
      return
    }

    if (this.data.password !== this.data.confirmPassword) {
      wx.showToast({
        title: '两次密码不一致',
        icon: 'none'
      })
      return
    }

    if (this.data.password.length < 6) {
      wx.showToast({
        title: '密码至少6个字符',
        icon: 'none'
      })
      return
    }

    this.setData({ submitting: true })

    try {
      const res = await api.bindEmail(
        this.data.email,
        this.data.code,
        this.data.password
      )

      if (res.code === 1) {
        // 绑定成功
        this.handleBindSuccess()
      } else if (res.code === 1021) {
        // 邮箱已存在，需要合并
        this.showMergeDialog()
      }
    } catch (error) {
      console.error('绑定失败:', error)
    } finally {
      this.setData({ submitting: false })
    }
  },

  // 显示合并账号对话框
  showMergeDialog() {
    const that = this
    wx.showModal({
      title: '邮箱已注册',
      content: '此邮箱已在Web端注册，是否合并到该账号？验证码已验证，可直接合并。',
      confirmText: '确认合并',
      cancelText: '取消',
      success(modal) {
        if (modal.confirm) {
          that.mergeAccount()
        }
      }
    })
  },

  // 合并账号
  async mergeAccount() {
    wx.showLoading({ title: '合并中...' })

    try {
      const res = await api.mergeAccount(
        this.data.email,
        this.data.code
      )

      if (res.code === 1) {
        wx.hideLoading()
        wx.showToast({
          title: '合并成功',
          icon: 'success'
        })

        // 保存新token和用户信息
        app.saveLoginInfo(res.token, res.data)

        // 跳转到首页
        setTimeout(() => {
          wx.switchTab({
            url: '/pages/home/home'
          })
        }, 1500)
      } else {
        wx.hideLoading()
        wx.showToast({
          title: res.msg || '合并失败',
          icon: 'none'
        })
      }
    } catch (error) {
      wx.hideLoading()
      console.error('合并失败:', error)
    }
  },

  // 绑定成功处理
  handleBindSuccess() {
    wx.showToast({
      title: '绑定成功',
      icon: 'success'
    })

    // 更新用户信息
    const userInfo = app.globalData.userInfo
    userInfo.email = this.data.email
    userInfo.has_email = true
    app.saveLoginInfo(app.globalData.token, userInfo)

    // 跳转到首页
    setTimeout(() => {
      wx.switchTab({
        url: '/pages/home/home'
      })
    }, 1500)
  },

  // 跳过绑定
  skip() {
    wx.switchTab({
      url: '/pages/home/home'
    })
  }
})
