// pages/home/home.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    skinList: [],
    platforms: ['Buff', 'UU'],
    sourceIndex: 0,  // Buff
    targetIndex: 1,  // UU
    pageNum: 1,
    pageSize: 20,
    hasMore: true,
    loading: false
  },

  onLoad() {
    this.checkLogin()
  },

  onShow() {
    if (app.globalData.token) {
      this.loadData()
    }
  },

  // 检查登录
  checkLogin() {
    if (!app.globalData.token) {
      wx.reLaunch({
        url: '/pages/login/login'
      })
      return false
    }
    return true
  },

  // 加载数据
  async loadData() {
    if (!this.checkLogin()) return
    
    this.setData({ loading: true })

    const source = this.data.sourceIndex === 0 ? 'buff' : 'uu'
    const target = this.data.targetIndex === 0 ? 'buff' : 'uu'

    try {
      const res = await api.getSkinItems({
        page_num: this.data.pageNum,
        page_size: this.data.pageSize,
        source: source,
        target: target
      })

      if (res.code === 1) {
        this.setData({
          skinList: res.data || [],
          hasMore: res.data && res.data.length >= this.data.pageSize
        })
      }
    } catch (error) {
      console.error('加载数据失败:', error)
    } finally {
      this.setData({ loading: false })
    }
  },

  // 切换买入平台
  onSourceChange(e) {
    const index = e.detail.value
    if (index == this.data.targetIndex) {
      // 自动切换目标平台
      this.setData({
        sourceIndex: index,
        targetIndex: index == 0 ? 1 : 0
      }, () => {
        this.loadData()
      })
    } else {
      this.setData({ sourceIndex: index }, () => {
        this.loadData()
      })
    }
  },

  // 切换卖出平台
  onTargetChange(e) {
    const index = e.detail.value
    if (index == this.data.sourceIndex) {
      // 自动切换来源平台
      this.setData({
        targetIndex: index,
        sourceIndex: index == 0 ? 1 : 0
      }, () => {
        this.loadData()
      })
    } else {
      this.setData({ targetIndex: index }, () => {
        this.loadData()
      })
    }
  },

  // 加载更多
  async loadMore() {
    if (this.data.loading || !this.data.hasMore) return
    
    this.setData({ 
      pageNum: this.data.pageNum + 1 
    }, () => {
      this.loadData()
    })
  },

  // 获取利润率样式
  getProfitClass(rate) {
    if (rate >= 60) return 'danger'
    if (rate >= 40) return 'warning'
    if (rate >= 20) return 'success'
    return 'primary'
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.setData({ pageNum: 1 }, () => {
      this.loadData()
      wx.stopPullDownRefresh()
    })
  }
})

