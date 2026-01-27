// pages/bind-email/bind-email.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    step: 1,  // 当前步骤：1=输入邮箱验证码，2=设置密码
    email: '',
    code: '',
    password: '',
    confirmPassword: '',
    countdown: 0,
    sending: false,
    submitting: false,
    emailExists: false  // 邮箱是否已存在（决定是新绑定还是合并）
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

        // 保存邮箱是否已存在的状态
        this.setData({ 
          emailExists: res.email_exists || false
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

  // 第一步：验证邮箱
  async handleStep1() {
    if (!this.data.email || !this.data.code) {
      wx.showToast({
        title: '请填写邮箱和验证码',
        icon: 'none'
      })
      return
    }

    this.setData({ submitting: true })

    try {
      // 先验证验证码
      const res = await api.verifyEmailCode(this.data.email, this.data.code)
      
      if (res.code === 1) {
        // 验证成功，判断下一步
        if (this.data.emailExists) {
          // 邮箱已存在，直接合并
          this.mergeAccount()
        } else {
          // 新绑定，进入设置密码步骤
          this.setData({ step: 2 })
        }
      }
    } catch (error) {
      console.error('验证失败:', error)
    } finally {
      this.setData({ submitting: false })
    }
  },

  // 第二步：设置密码并绑定
  async handleStep2() {
    if (!this.data.password || !this.data.confirmPassword) {
      wx.showToast({
        title: '请设置密码',
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
        this.handleBindSuccess(res)
      } else if (res.code === 1021) {
        // 邮箱已存在（兜底处理）
        this.mergeAccount()
      }
    } catch (error) {
      console.error('绑定失败:', error)
    } finally {
      this.setData({ submitting: false })
    }
  },

  // 返回上一步
  goBack() {
    this.setData({ step: 1 })
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
  handleBindSuccess(res) {
    wx.showToast({
      title: '绑定成功，获得2天VIP试用',
      icon: 'none',
      duration: 2000
    })

    // 更新用户信息，包括VIP状态
    const userInfo = app.globalData.userInfo
    userInfo.email = this.data.email
    userInfo.has_email = true
    userInfo.role = 1  // VIP角色
    if (res && res.data && res.data.vip_expiry) {
      userInfo.vip_expiry = res.data.vip_expiry
    }
    app.saveLoginInfo(app.globalData.token, userInfo)

    // 跳转到首页
    setTimeout(() => {
      wx.switchTab({
        url: '/pages/home/home'
      })
    }, 2000)
  },

  // 跳过绑定
  skip() {
    wx.switchTab({
      url: '/pages/home/home'
    })
  }
})
