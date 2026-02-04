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
    
    // 购买/出售方案
    // buyType: sell(在售价购买) / bidding(求购价购买)
    // sellType: sell(挂底价) / bidding(丢求购)
    buyType: 'sell',
    sellType: 'sell',
    
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
    showFilter: false,
    
    // 筛选描述
    filterDescription: ''
  },

  onLoad() {
    this.checkLoginStatus()
    if (this.data.isLoggedIn) {
      this.loadSettings().then(() => {
        this.updateFilterDescription()
        this.loadData()
      })
    } else {
      this.updateFilterDescription()
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
    
    // 如果VIP开关关闭，不引导开通VIP
    const vipEnabled = app.globalData.minAppConfig?.vipEnabled || false

    const userInfo = app.globalData.userInfo
    if (!userInfo) return

    // 判断VIP状态：role为1且vip_expiry未过期
    let isVip = false
    if (userInfo.role === 1 && userInfo.vip_expiry) {
      const vipExpiry = new Date(userInfo.vip_expiry)
      isVip = vipExpiry > new Date()
    }
    const hasEmail = userInfo.email && userInfo.email !== ''

    // 优先提醒绑定邮箱（未绑定邮箱的用户每次登录都提醒）
    if (!hasEmail) {
      app.globalData.hasShownGuide = true
      this.showBindEmailGuide()
      return
    }

    // 如果不是VIP且VIP开关开启，引导开通VIP
    if (!isVip && vipEnabled) {
      app.globalData.hasShownGuide = true
      this.showVipGuide()
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
      content: '绑定邮箱后可在网页端登录使用更多功能',
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

  // 类别标签点击
  onCategoryTap(e) {
    const index = e.currentTarget.dataset.index
    if (index === this.data.categoryIndex) return
    this.setData({ 
      categoryIndex: index,
      pageNum: 1,
      skinList: []
    }, () => {
      this.updateFilterDescription()
      this.loadData()
    })
  },

  // 排序标签点击
  onSortTap(e) {
    const index = e.currentTarget.dataset.index
    if (index === this.data.sortIndex) return
    this.setData({ 
      sortIndex: index,
      pageNum: 1,
      skinList: []
    }, () => {
      this.loadData()
    })
  },

  // 类别选择（picker备用）
  onCategoryChange(e) {
    this.onCategoryTap({ currentTarget: { dataset: { index: parseInt(e.detail.value) } } })
  },

  // 排序选择（picker备用）
  onSortChange(e) {
    this.onSortTap({ currentTarget: { dataset: { index: parseInt(e.detail.value) } } })
  },

  // 买入平台切换
  onSourceChange(e) {
    const index = parseInt(e.detail.value)
    if (index === this.data.sourceIndex) return
    
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
      this.updateFilterDescription()
      this.loadData()
    })
  },

  // 卖出平台切换
  onTargetChange(e) {
    const index = parseInt(e.detail.value)
    if (index === this.data.targetIndex) return
    
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
      this.updateFilterDescription()
      this.loadData()
    })
  },

  // 购买方案切换
  onBuyTypeChange(e) {
    const type = e.currentTarget.dataset.type
    if (type === this.data.buyType) return
    this.setData({
      buyType: type,
      pageNum: 1,
      skinList: []
    }, () => {
      this.updateFilterDescription()
      this.loadData()
    })
  },

  // 出售方案切换
  onSellTypeChange(e) {
    const type = e.currentTarget.dataset.type
    if (type === this.data.sellType) return
    this.setData({
      sellType: type,
      pageNum: 1,
      skinList: []
    }, () => {
      this.updateFilterDescription()
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
      this.updateFilterDescription()
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
      target: target,
      buy_type: this.data.buyType,
      sell_type: this.data.sellType
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
    // 平台名称和图标映射
    const platformIconMap = {
      'BUFF': '/images/buff.png',
      '悠悠': '/images/uu.png',
      'C5GAME': '/images/c5.png',
      'Steam': '/images/steam.png'
    }
    
    // 当前卖出平台名称
    const targetPlatformName = this.data.platforms[this.data.targetIndex]
    
    return list.map(item => {
      const rate = item.profit_rate || 0
      const percent = rate * 100
      
      // 处理平台列表中的价格
      let platformList = item.platform_list || []
      if (platformList.length > 0) {
        platformList = platformList.map(p => ({
          ...p,
          sellPriceNum: Number(p.sellPrice || 0),
          sellPrice: Number(p.sellPrice || 0).toFixed(2),
          biddingPrice: Number(p.biddingPrice || 0).toFixed(2),
          icon: platformIconMap[p.platformName] || '/images/buff.png'
        }))
      }
      
      // 分离卖出平台数据和其他平台数据
      let targetPlatformData = null
      let otherPlatforms = []
      
      if (platformList.length > 0) {
        targetPlatformData = platformList.find(p => p.platformName === targetPlatformName)
        otherPlatforms = platformList.filter(p => p.platformName !== targetPlatformName)
      }
      
      // 判断是否有平台价格低于买入平台价格
      const sourcePrice = Number(item.source_price || 0)
      const hasLowerPrice = platformList.some(p => 
        p.sellPriceNum > 0 && p.sellPriceNum < sourcePrice
      )
      
      return {
        ...item,
        sourcePriceText: Number(item.source_price || 0).toFixed(2),
        targetPriceText: Number(item.target_price || 0).toFixed(2),
        priceDiffText: Number(item.price_diff || 0).toFixed(2),
        profitRateText: percent.toFixed(1) + '%',
        profitClass: this.getProfitClass(rate),
        platform_list: platformList,
        targetPlatformData: targetPlatformData,
        otherPlatforms: otherPlatforms,
        showPlatform: false,
        hasLowerPrice: hasLowerPrice
      }
    })
  },

  loadMore() {
    if (this.data.loading || !this.data.hasMore || !this.data.isLoggedIn) return
    this.setData({ pageNum: this.data.pageNum + 1 }, () => {
      this.loadData()
    })
  },

  // 页面触底自动加载更多
  onReachBottom() {
    this.loadMore()
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

  // 更新筛选描述
  updateFilterDescription() {
    const sourceName = this.data.platforms[this.data.sourceIndex]
    const targetName = this.data.platforms[this.data.targetIndex]
    
    // 购买方案描述
    const buyTypeDesc = this.data.buyType === 'bidding' ? '求购价购买' : '在售价购买'
    // 出售方案描述
    const sellTypeDesc = this.data.sellType === 'bidding' ? '丢求购' : '挂底价'
    
    const minPrice = this.data.minPrice || 0
    const maxPrice = this.data.maxPrice || 10000
    const priceRange = `价格${minPrice}-${maxPrice}元`
    
    let conditions = []
    if (this.data.minSellNum) {
      conditions.push(`在售量大于${this.data.minSellNum}件`)
    }
    if (this.data.minDiff) {
      conditions.push(`价格差大于${this.data.minDiff}元`)
    }
    if (this.data.categoryIndex > 0) {
      conditions.push(`类别：${this.data.categoryOptions[this.data.categoryIndex]}`)
    }
    
    const conditionStr = conditions.length > 0 ? '，' + conditions.join('，') : ''
    
    this.setData({
      filterDescription: `从${sourceName}（${buyTypeDesc}）→ ${targetName}（${sellTypeDesc}）（${priceRange}${conditionStr}）`
    })
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
