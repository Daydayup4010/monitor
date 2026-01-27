// pages/home/home.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    skinList: [],
    loading: false,
    isLoggedIn: false,
    
    // 平台选择
    platforms: ['BUFF', '悠悠', 'C5GAME', 'Steam'],
    platformKeys: ['buff', 'uu', 'c5', 'steam'],
    sourceIndex: 1,  // UU
    targetIndex: 0,  // Buff
    
    // 筛选参数
    minPrice: '',      // 最低价格
    maxPrice: '',      // 最高价格
    minSellNum: '',    // 最小在售量
    minDiff: '',       // 最小价格差
    
    // 类别筛选
    categoryOptions: ['全部类别', '匕首', '手套', '步枪', '手枪', '微型冲锋枪', '霰弹枪', '机枪', '印花', '涂鸦', '探员', '挂件', '音乐盒', '武器箱', '布章'],
    categoryIndex: 0,
    
    // 排序
    sortOptions: ['默认排序', '价格差 ↑', '价格差 ↓', '利润率 ↑', '利润率 ↓'],
    sortIndex: 0,
    
    // 分页
    pageNum: 1,
    pageSize: 20,
    hasMore: true,
    total: 0,
    
    // 展开筛选面板
    showFilter: false
  },

  onLoad() {
    this.checkLoginStatus()
    if (this.data.isLoggedIn) {
      this.loadSettings().then(() => {
        this.loadData()
      })
    } else {
      this.loadData()
    }
  },

  onShow() {
    // 从storage重新读取用户信息，确保是最新的
    app.checkLoginStatus()
    this.checkLoginStatus()
    
    const isLoggedIn = !!app.globalData.token
    if (isLoggedIn !== this.data.isLoggedIn) {
      this.setData({ 
        isLoggedIn,
        pageNum: 1,
        skinList: []
      })
      if (isLoggedIn) {
        this.loadSettings().then(() => {
          this.loadData()
        })
        // 检查是否需要引导
        this.checkUserGuide()
      } else {
        this.loadData()
      }
    } else if (isLoggedIn) {
      // 已登录状态下每次显示都检查引导
      this.checkUserGuide()
    }
  },

  // 检查用户引导（VIP和邮箱绑定）
  checkUserGuide() {
    // 如果已经显示过引导，不再重复
    if (app.globalData.hasShownGuide) return

    const userInfo = app.globalData.userInfo
    if (!userInfo) return

    // 判断VIP状态：role为1且vip_expiry未过期
    let isVip = false
    if (userInfo.role === 1 && userInfo.vip_expiry) {
      const vipExpiry = new Date(userInfo.vip_expiry)
      isVip = vipExpiry > new Date()
    }
    const hasEmail = userInfo.email && userInfo.email !== ''

    // 如果不是VIP
    if (!isVip) {
      app.globalData.hasShownGuide = true
      if (hasEmail) {
        // 已绑定邮箱，引导开通VIP
        this.showVipGuide()
      } else {
        // 未绑定邮箱，提示先绑定邮箱
        this.showBindEmailGuide()
      }
    }
  },

  // 显示VIP引导
  showVipGuide() {
    wx.showModal({
      title: '开通VIP',
      content: '开通VIP会员，解锁全部数据和功能',
      confirmText: '去开通',
      cancelText: '稍后再说',
      success: (res) => {
        if (res.confirm) {
          wx.navigateTo({
            url: '/pages/vip/vip'
          })
        }
      }
    })
  },

  // 显示绑定邮箱引导
  showBindEmailGuide() {
    wx.showModal({
      title: '绑定邮箱',
      content: '绑定邮箱可免费获得2天VIP试用时间，还能在Web端登录查看更多数据',
      confirmText: '去绑定',
      cancelText: '稍后再说',
      success: (res) => {
        if (res.confirm) {
          wx.navigateTo({
            url: '/pages/bind-email/bind-email'
          })
        }
      }
    })
  },

  // 加载用户设置
  async loadSettings() {
    try {
      const res = await api.getSettings()
      if (res.code === 1 && res.data) {
        this.setData({
          minPrice: res.data.min_sell_price || '',
          maxPrice: res.data.max_sell_price || '',
          minSellNum: res.data.min_sell_num || '',
          minDiff: res.data.min_diff || ''
        })
      }
    } catch (error) {
      console.error('加载设置失败:', error)
    }
  },

  checkLoginStatus() {
    this.setData({
      isLoggedIn: !!app.globalData.token
    })
  },

  goLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    })
  },

  // 切换筛选面板
  toggleFilter() {
    this.setData({
      showFilter: !this.data.showFilter
    })
  },

  // 输入事件
  onMinPriceInput(e) {
    this.setData({ minPrice: e.detail.value })
  },
  onMaxPriceInput(e) {
    this.setData({ maxPrice: e.detail.value })
  },
  onMinSellNumInput(e) {
    this.setData({ minSellNum: e.detail.value })
  },
  onMinDiffInput(e) {
    this.setData({ minDiff: e.detail.value })
  },

  // 类别选择
  onCategoryChange(e) {
    this.setData({ categoryIndex: parseInt(e.detail.value) })
  },

  // 排序选择
  onSortChange(e) {
    this.setData({ sortIndex: parseInt(e.detail.value) })
  },

  // 平台切换
  onSourceChange(e) {
    const index = parseInt(e.detail.value)
    let targetIndex = this.data.targetIndex
    // 如果买入卖出平台相同，自动切换到下一个
    if (index === targetIndex) {
      targetIndex = (index + 1) % this.data.platforms.length
    }
    this.setData({
      sourceIndex: index,
      targetIndex: targetIndex,
      pageNum: 1,
      skinList: []
    }, () => {
      this.loadData()
    })
  },

  onTargetChange(e) {
    const index = parseInt(e.detail.value)
    let sourceIndex = this.data.sourceIndex
    // 如果买入卖出平台相同，自动切换到下一个
    if (index === sourceIndex) {
      sourceIndex = (index + 1) % this.data.platforms.length
    }
    this.setData({
      targetIndex: index,
      sourceIndex: sourceIndex,
      pageNum: 1,
      skinList: []
    }, () => {
      this.loadData()
    })
  },

  // 确定搜索（同时保存设置）
  async doSearch() {
    // 保存设置到后端
    try {
      await api.updateSettings({
        min_sell_price: Number(this.data.minPrice) || 0,
        max_sell_price: Number(this.data.maxPrice) || 10000,
        min_sell_num: Number(this.data.minSellNum) || 0,
        min_diff: Number(this.data.minDiff) || 0
      })
    } catch (error) {
      console.error('保存设置失败:', error)
    }

    this.setData({
      pageNum: 1,
      skinList: [],
      showFilter: false
    }, () => {
      this.loadData()
    })
  },

  // 重置筛选
  resetFilter() {
    this.setData({
      minPrice: '',
      maxPrice: '',
      minSellNum: '',
      minDiff: '',
      categoryIndex: 0,
      sortIndex: 0
    })
  },

  // 加载数据
  async loadData() {
    this.setData({ loading: true })

    try {
      if (this.data.isLoggedIn) {
        await this.loadVipData()
      } else {
        await this.loadPublicData()
      }
    } catch (error) {
      console.error('加载数据失败:', error)
    } finally {
      this.setData({ loading: false })
    }
  },

  async loadPublicData() {
    const res = await api.getPublicHomeData()
    if (res.code === 1 && res.data) {
      const skinList = this.processData(res.data.brickMoving || [])
      this.setData({
        skinList,
        hasMore: false,
        total: skinList.length
      })
    }
  },

  async loadVipData() {
    const source = this.data.platformKeys[this.data.sourceIndex]
    const target = this.data.platformKeys[this.data.targetIndex]

    const params = {
      page_num: this.data.pageNum,
      page_size: this.data.pageSize,
      source: source,
      target: target
    }

    // 添加筛选参数
    if (this.data.minPrice) {
      params.min_sell_price = Number(this.data.minPrice)
    }
    if (this.data.maxPrice) {
      params.max_sell_price = Number(this.data.maxPrice)
    }
    if (this.data.minSellNum) {
      params.min_sell_num = Number(this.data.minSellNum)
    }
    if (this.data.minDiff) {
      params.min_diff = Number(this.data.minDiff)
    }
    if (this.data.categoryIndex > 0) {
      params.category = this.data.categoryOptions[this.data.categoryIndex]
    }

    // 排序
    const sortMap = [
      { field: '', desc: false },
      { field: 'price_diff', desc: false },
      { field: 'price_diff', desc: true },
      { field: 'profit_rate', desc: false },
      { field: 'profit_rate', desc: true }
    ]
    const sort = sortMap[this.data.sortIndex]
    if (sort.field) {
      params.sort = sort.field
      params.desc = sort.desc
    }

    const res = await api.getSkinItems(params)

    if (res.code === 1) {
      const processedData = this.processData(res.data || [])
      const newList = this.data.pageNum === 1 ? processedData : [...this.data.skinList, ...processedData]
      this.setData({
        skinList: newList,
        hasMore: res.data && res.data.length >= this.data.pageSize,
        total: res.total || newList.length
      })
    }
  },

  processData(list) {
    return list.map(item => {
      const rate = item.profit_rate || 0
      const percent = rate * 100
      
      // 处理平台列表中的价格
      let platformList = item.platform_list || []
      if (platformList.length > 0) {
        platformList = platformList.map(p => ({
          ...p,
          sellPrice: Number(p.sellPrice || 0).toFixed(2),
          biddingPrice: Number(p.biddingPrice || 0).toFixed(2)
        }))
      }
      
      return {
        ...item,
        sourcePriceText: Number(item.source_price || 0).toFixed(2),
        targetPriceText: Number(item.target_price || 0).toFixed(2),
        priceDiffText: Number(item.price_diff || 0).toFixed(2),
        profitRateText: percent.toFixed(1) + '%',
        profitClass: this.getProfitClass(rate),
        platform_list: platformList,
        showPlatform: false
      }
    })
  },

  loadMore() {
    if (this.data.loading || !this.data.hasMore || !this.data.isLoggedIn) return
    this.setData({ pageNum: this.data.pageNum + 1 }, () => {
      this.loadData()
    })
  },

  onPullDownRefresh() {
    this.setData({ pageNum: 1, skinList: [] }, () => {
      this.loadData().then(() => {
        wx.stopPullDownRefresh()
      })
    })
  },

  getProfitClass(rate) {
    const percent = rate * 100
    if (percent >= 60) return 'danger'
    if (percent >= 40) return 'warning'
    if (percent >= 20) return 'success'
    return 'primary'
  },

  // 切换平台数据展开状态
  togglePlatformData(e) {
    const index = e.currentTarget.dataset.index
    const key = `skinList[${index}].showPlatform`
    const currentValue = this.data.skinList[index].showPlatform
    this.setData({
      [key]: !currentValue
    })
  },

  // 跳转到详情页
  goDetail(e) {
    const marketHashName = e.currentTarget.dataset.hashname
    if (!marketHashName) return
    
    // 检查登录状态
    if (!this.data.isLoggedIn) {
      wx.showModal({
        title: '提示',
        content: '登录后可查看详情',
        confirmText: '去登录',
        cancelText: '取消',
        success: (res) => {
          if (res.confirm) {
            this.goLogin()
          }
        }
      })
      return
    }
    
    wx.navigateTo({
      url: `/pages/detail/detail?market_hash_name=${encodeURIComponent(marketHashName)}`
    })
  }
})
