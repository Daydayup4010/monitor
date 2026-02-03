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
    console.log('vip.onShow - payResult:', payResult)
    
    if (payResult) {
      // 清除支付结果，避免重复处理
      app.globalData.payResult = null
      
      if (payResult.success) {
        wx.showToast({
          title: '支付成功，正在更新...',
          icon: 'loading',
          duration: 2000
        })
        // 刷新用户信息
        this.refreshUserInfo()
      } else {
        console.log('支付失败:', payResult.msg)
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
        
        if (plans.length === 0) {
          this.setData({ loading: false })
          return
        }
        
        // 找出月均价最低的套餐（用于计算折扣和推荐）
        const monthlyPrices = plans.map(p => p.price / p.months)
        const lowestMonthlyPrice = Math.min(...monthlyPrices)
        const lowestMonthlyPriceIndex = monthlyPrices.indexOf(lowestMonthlyPrice)
        
        // 获取单月套餐价格作为折扣基准（如果没有则用第一个套餐的月均价）
        const monthlyPlan = plans.find(p => p.months === 1)
        const baseMonthlyPrice = monthlyPlan ? monthlyPlan.price : plans[0].price / plans[0].months
        
        // 添加显示信息
        const plansWithInfo = plans.map((plan, index) => {
          const pricePerMonth = plan.price / plan.months
          // 相对于单月价格的折扣
          const savings = plan.months > 1 ? Math.round((1 - pricePerMonth / baseMonthlyPrice) * 100) : 0
          
          return {
            ...plan,
            title: plan.months === 1 ? '月卡' : plan.months === 12 ? '年卡' : `${plan.months}个月`,
            pricePerMonth: pricePerMonth.toFixed(1),
            savings: savings > 0 ? savings : 0,
            recommend: index === lowestMonthlyPriceIndex  // 月均价最低的为推荐
          }
        })
        
        // 默认选中推荐套餐
        const defaultPlan = plansWithInfo.find(p => p.recommend) || plansWithInfo[0]
        
        this.setData({ 
          plans: plansWithInfo,
          selectedPlan: defaultPlan,
          lowestPrice: lowestMonthlyPrice.toFixed(1)
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
      // 等待一下，确保后端回调已处理
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      const res = await api.getUserInfo()
      if (res.code === 1 && res.data) {
        app.globalData.userInfo = res.data
        wx.setStorageSync('userInfo', res.data)
        
        // 更新 VIP 到期时间显示
        let formatExpiry = ''
        if (res.data.vip_expiry) {
          const date = new Date(res.data.vip_expiry)
          formatExpiry = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
        }
        
        this.setData({ 
          userInfo: res.data,
          formatExpiry: formatExpiry
        })
        
        console.log('用户信息已刷新:', res.data)
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
