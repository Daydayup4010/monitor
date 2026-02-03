// pages/vip/vip.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    plans: [],
    selectedPlan: null,
    loading: false,
    paying: false,
    userInfo: null,
    lowestPrice: '14.2',  // 默认最低月付价格
    formatExpiry: ''      // 格式化的到期时间
  },

  onLoad() {
    this.loadVipPlans()
    this.checkUserInfo()
  },

  onShow() {
    this.checkUserInfo()
    
    // 处理嵌入式小程序支付回调
    const payResult = app.globalData.payResult
    if (payResult) {
      // 清除支付结果，避免重复处理
      app.globalData.payResult = null
      
      if (payResult.success) {
        wx.showToast({
          title: '支付成功',
          icon: 'success'
        })
        // 刷新用户信息
        this.refreshUserInfo()
      }
    }
  },

  // 检查用户信息
  checkUserInfo() {
    const userInfo = app.globalData.userInfo
    let formatExpiry = ''
    if (userInfo && userInfo.vip_expiry) {
      const date = new Date(userInfo.vip_expiry)
      formatExpiry = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    }
    this.setData({ userInfo, formatExpiry })
  },

  // 加载VIP套餐
  async loadVipPlans() {
    this.setData({ loading: true })
    
    try {
      const res = await api.getVipPlans()
      if (res.code === 1 && res.data && res.data.plans) {
        // 将 plans 对象转换为数组并排序
        const plansObj = res.data.plans
        const plans = Object.values(plansObj).sort((a, b) => a.months - b.months)
        
        // 添加显示信息
        const plansWithInfo = plans.map(plan => ({
          ...plan,
          title: plan.months === 1 ? '月卡' : plan.months === 12 ? '年卡' : `${plan.months}个月`,
          pricePerMonth: (plan.price / plan.months).toFixed(1),
          savings: plan.months > 1 ? Math.round((1 - plan.price / (19.9 * plan.months)) * 100) : 0,
          recommend: plan.months === 12  // 12个月为推荐
        }))
        
        // 计算最低月付价格
        const lowestPrice = Math.min(...plansWithInfo.map(p => p.price / p.months)).toFixed(1)
        
        // 默认选中12个月（推荐）
        const defaultPlan = plansWithInfo.find(p => p.months === 12) || plansWithInfo[0]
        
        this.setData({ 
          plans: plansWithInfo,
          selectedPlan: defaultPlan,
          lowestPrice
        })
      }
    } catch (error) {
      console.error('加载套餐失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 选择套餐
  selectPlan(e) {
    const months = e.currentTarget.dataset.months
    const plan = this.data.plans.find(p => p.months === months)
    if (plan) {
      this.setData({ selectedPlan: plan })
    }
  },

  // 发起支付
  async handlePay() {
    if (!this.data.selectedPlan) {
      wx.showToast({
        title: '请选择套餐',
        icon: 'none'
      })
      return
    }

    // 检查登录状态
    if (!app.globalData.token) {
      wx.showModal({
        title: '提示',
        content: '请先登录',
        confirmText: '去登录',
        success: (res) => {
          if (res.confirm) {
            wx.navigateTo({
              url: '/pages/login/login'
            })
          }
        }
      })
      return
    }

    this.setData({ paying: true })

    try {
      // 创建支付订单（获取嵌入式小程序支付参数）
      const res = await api.createMinAppPay({
        months: this.data.selectedPlan.months
      })

      if (res.code === 1 && res.data) {
        // 使用嵌入式小程序支付（YunGouOS收银台）
        const payParams = res.data
        wx.openEmbeddedMiniProgram({
          appId: 'wxd9634afb01b983c0', // YunGouOS收银台小程序AppID（固定值）
          path: '/pages/pay/pay',      // 支付页面路径（固定值）
          extraData: {
            out_trade_no: payParams.out_trade_no,
            total_fee: payParams.total_fee,
            mch_id: payParams.mch_id,
            body: payParams.body,
            notify_url: payParams.notify_url,
            attach: payParams.attach,
            sign: payParams.sign,
            title: 'CS搬砖助手'
          },
          success: () => {
            console.log('打开收银台成功')
          },
          fail: (err) => {
            console.error('打开收银台失败:', err)
            wx.showToast({
              title: '打开支付失败',
              icon: 'none'
            })
          }
        })
      } else {
        wx.showToast({
          title: res.msg || '创建订单失败',
          icon: 'none'
        })
      }
    } catch (error) {
      console.error('支付错误:', error)
      wx.showToast({
        title: '支付失败',
        icon: 'none'
      })
    } finally {
      this.setData({ paying: false })
    }
  },

  // 刷新用户信息
  async refreshUserInfo() {
    try {
      const res = await api.getUserInfo()
      if (res.code === 1 && res.data) {
        app.globalData.userInfo = res.data
        wx.setStorageSync('userInfo', res.data)
        this.setData({ userInfo: res.data })
      }
    } catch (error) {
      console.error('刷新用户信息失败:', error)
    }
  },

  // 查看VIP记录
  goVipRecords() {
    wx.navigateTo({
      url: '/pages/vip-records/vip-records'
    })
  }
})
