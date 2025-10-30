// pages/home/home.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    skinList: [],
    platforms: ['Buff', 'UU'],
    sourceIndex: 0,  // Buff
    targetIndex: 1,  // UU
    
    // 类别筛选
    categoryOptions: ['全部类别'],  // 从后端动态加载
    categories: [],  // 后端返回的分类列表
    categoryIndex: 0,
    categoryValue: '',
    
    // 排序选项
    sortOptions: ['默认排序', '价格差 ↑', '价格差 ↓', '利润率 ↑', '利润率 ↓'],
    sortIndex: 0,
    sortField: '',
    sortDesc: false,
    
    pageNum: 1,
    pageSize: 20,
    hasMore: true,
    loading: false
  },

  onLoad() {
    this.checkLogin()
    this.loadCategories()  // 加载分类列表
  },

  onShow() {
    if (app.globalData.token) {
      this.loadData()
    }
  },

  // 加载分类列表
  async loadCategories() {
    try {
      const res = await api.getCategories()
      if (res.code === 1 && res.data) {
        const categoryOptions = ['全部类别', ...res.data]
        this.setData({
          categories: res.data,
          categoryOptions: categoryOptions
        })
      }
    } catch (error) {
      console.error('加载分类失败:', error)
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

    const params = {
      page_num: this.data.pageNum,
      page_size: this.data.pageSize,
      source: source,
      target: target
    }

    // 添加类别筛选
    if (this.data.categoryValue) {
      params.category = this.data.categoryValue
    }

    // 添加排序
    if (this.data.sortField) {
      params.sort = this.data.sortField
      params.desc = this.data.sortDesc
    }

    try {
      const res = await api.getSkinItems(params)

      if (res.code === 1) {
        const newList = this.data.pageNum === 1 ? res.data : [...this.data.skinList, ...res.data]
        this.setData({
          skinList: newList || [],
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
        targetIndex: index == 0 ? 1 : 0,
        pageNum: 1
      }, () => {
        this.loadData()
      })
    } else {
      this.setData({ 
        sourceIndex: index,
        pageNum: 1
      }, () => {
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
        sourceIndex: index == 0 ? 1 : 0,
        pageNum: 1
      }, () => {
        this.loadData()
      })
    } else {
      this.setData({ 
        targetIndex: index,
        pageNum: 1
      }, () => {
        this.loadData()
      })
    }
  },

  // 类别筛选
  onCategoryChange(e) {
    const index = e.detail.value
    // index=0是"全部类别"，对应空字符串
    // index>0是具体分类
    const categoryValue = index === 0 ? '' : this.data.categories[index - 1]
    
    this.setData({
      categoryIndex: index,
      categoryValue: categoryValue,
      pageNum: 1
    }, () => {
      this.loadData()
    })
  },

  // 排序筛选
  onSortChange(e) {
    const index = e.detail.value
    const sortMap = [
      { field: '', desc: false },          // 默认
      { field: 'price_diff', desc: false }, // 价格差 ↑
      { field: 'price_diff', desc: true },  // 价格差 ↓
      { field: 'profit_rate', desc: false },// 利润率 ↑
      { field: 'profit_rate', desc: true }, // 利润率 ↓
    ]
    
    const sort = sortMap[index]
    
    this.setData({
      sortIndex: index,
      sortField: sort.field,
      sortDesc: sort.desc,
      pageNum: 1
    }, () => {
      this.loadData()
    })
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
    this.setData({ 
      pageNum: 1,
      skinList: []
    }, () => {
      this.loadData()
      wx.stopPullDownRefresh()
    })
  }
})
