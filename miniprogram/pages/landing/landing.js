// pages/landing/landing.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    searchKeyword: '',
    searchResults: [],
    showResults: false,
    loading: false,
    isLoggedIn: false,
    needLogin: false,
    isTouchingResults: false  // 是否正在触摸搜索结果列表
  },

  onLoad() {
    this.checkLoginStatus()
  },

  onShow() {
    this.checkLoginStatus()
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    this.setData({
      isLoggedIn: !!token
    })
  },

  // 搜索输入
  onSearchInput(e) {
    const keyword = e.detail.value
    this.setData({ searchKeyword: keyword })
    
    // 未登录提示
    if (!this.data.isLoggedIn) {
      if (keyword.trim()) {
        this.setData({
          showResults: true,
          searchResults: [],
          needLogin: true
        })
      } else {
        this.setData({ showResults: false, needLogin: false })
      }
      return
    }
    
    // 防抖处理
    if (this.searchTimer) {
      clearTimeout(this.searchTimer)
    }
    
    if (!keyword.trim()) {
      this.setData({
        searchResults: [],
        showResults: false
      })
      return
    }
    
    this.searchTimer = setTimeout(() => {
      this.doSearch()
    }, 300)
  },

  // 执行搜索
  async doSearch() {
    if (!this.data.searchKeyword.trim()) return
    
    this.setData({ loading: true, showResults: true })
    
    try {
      const res = await api.searchGoods({
        keyword: this.data.searchKeyword.trim(),
        limit: 30
      })
      
      if (res.code === 1 && res.data) {
        this.setData({
          searchResults: res.data
        })
      }
    } catch (error) {
      console.error('搜索失败:', error)
    } finally {
      this.setData({ loading: false })
    }
  },

  // 清除搜索
  clearSearch() {
    this.setData({
      searchKeyword: '',
      searchResults: [],
      showResults: false
    })
  },

  // 点击搜索结果
  goDetail(e) {
    const { item } = e.currentTarget.dataset
    
    if (!this.data.isLoggedIn) {
      wx.showModal({
        title: '提示',
        content: '请先登录后查看饰品详情',
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
    
    wx.navigateTo({
      url: `/pages/detail/detail?market_hash_name=${encodeURIComponent(item.marketHashName)}`
    })
  },

  // 隐藏搜索结果
  hideResults() {
    // 如果正在触摸搜索结果列表，不隐藏
    if (this.data.isTouchingResults) {
      return
    }
    // 延迟隐藏，以便点击事件能够触发
    setTimeout(() => {
      if (!this.data.isTouchingResults) {
        this.setData({ showResults: false })
      }
    }, 200)
  },

  // 开始触摸搜索结果
  onResultsTouchStart() {
    this.setData({ isTouchingResults: true })
  },

  // 结束触摸搜索结果
  onResultsTouchEnd() {
    this.setData({ isTouchingResults: false })
  },

  // 阻止滑动事件冒泡
  onResultsTouchMove() {
    // 阻止事件冒泡，让 scroll-view 处理滚动
  },

  // 图片加载失败
  onImageError(e) {
    const index = e.currentTarget.dataset.index
    const key = `searchResults[${index}].iconUrl`
    this.setData({
      [key]: 'https://steamcommunity-a.akamaihd.net/economy/image/-9a81dlWLwJ2UUGcVs_nsVtzdOEdtWwKGZZLQHTxDZ7I56KU0Zwwo4NUX4oFJZEHLbXQ9QVcJY8gulRPQV6CF7b9mMnYZh9SHY27gZKBl_JbMKyJI24H65S1xtXZwKb2YOqHxj4F68Nz2L2Y9oj2jQDm_RY4am-mctWXdFc5NQuDqAHqx-fmg5_v7oOJlyU1fQmQdw/360fx360f'
    })
  },

  // 去登录
  goLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    })
  },

  // 去搬砖榜
  goHome() {
    wx.switchTab({
      url: '/pages/home/home'
    })
  }
})
