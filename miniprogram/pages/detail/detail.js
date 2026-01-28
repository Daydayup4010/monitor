// pages/detail/detail.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    loading: true,
    goodsDetail: null,
    market_hash_name: '',
    // 平台图标映射
    platformIconMap: {
      'YOUPIN': '/images/uu.png',
      'BUFF': '/images/buff.png',
      'C5': '/images/c5.png',
      'STEAM': '/images/steam.png'
    },
    // 走势图相关
    platforms: [
      { value: 'YOUPIN', label: '悠悠' },
      { value: 'BUFF', label: 'BUFF' },
      { value: 'C5', label: 'C5GAME' },
      { value: 'STEAM', label: 'Steam' }
    ],
    timeRanges: [
      { value: 30, label: '近1个月' },
      { value: 90, label: '近3个月' },
      { value: 180, label: '近6个月' },
      { value: 365, label: '近1年' }
    ],
    selectedPlatformIndex: 0,
    selectedTimeIndex: 0,
    chartData: null,
    // 磨损/品质相关
    relatedWears: null,
    currentQuality: '',
    currentQualityWears: [],
    hasWearOptions: false,
    hasOtherQuality: false,
    otherQualityLabel: '',
    otherQualityMinPrice: 0,
    // 磨损等级中文映射
    wearLabelMap: {
      'Factory New': '崭新出厂',
      'Minimal Wear': '略有磨损',
      'Field-Tested': '久经沙场',
      'Well-Worn': '破损不堪',
      'Battle-Scarred': '战痕累累',
      'NO_WEAR': '无磨损'
    },
    // 贴纸变体中文映射
    stickerVariantMap: {
      '普通': '普通',
      'Holo': '全息',
      'Gold': '金色',
      'Foil': '闪亮'
    },
    // 图表tooltip相关
    showTooltip: false,
    tooltipX: 0,
    tooltipY: 0,
    tooltipData: {
      date: '',
      price: '',
      count: ''
    }
  },

  onLoad(options) {
    if (options.market_hash_name) {
      this.setData({
        market_hash_name: decodeURIComponent(options.market_hash_name)
      })
      this.fetchDetail()
      this.fetchRelatedWears()
    }
  },

  onReady() {
    if (this.data.goodsDetail) {
      this.initChartAndDraw()
    }
  },

  // 获取详情数据
  async fetchDetail() {
    this.setData({ loading: true })
    
    try {
      const days = this.data.timeRanges[this.data.selectedTimeIndex].value
      const res = await api.getGoodsDetail({
        market_hash_name: this.data.market_hash_name,
        days: days
      })
      
      if (res.code === 1 && res.data) {
        const detail = this.processDetail(res.data)
        this.setData({
          goodsDetail: detail
        }, () => {
          setTimeout(() => {
            this.initChartAndDraw()
          }, 100)
        })
        wx.setNavigationBarTitle({
          title: detail.name || '饰品详情'
        })
      }
    } catch (error) {
      console.error('获取详情失败:', error)
      wx.showToast({
        title: '加载失败',
        icon: 'none'
      })
    } finally {
      this.setData({ loading: false })
    }
  },

  // 获取关联磨损数据
  async fetchRelatedWears() {
    try {
      const res = await api.getRelatedWears({
        market_hash_name: this.data.market_hash_name
      })
      
      if (res.code === 1 && res.data) {
        this.setData({ relatedWears: res.data })
        this.processRelatedWears()
      }
    } catch (error) {
      console.error('获取关联磨损失败:', error)
    }
  },

  // 处理关联磨损数据
  processRelatedWears() {
    const relatedWears = this.data.relatedWears
    if (!relatedWears) return

    const marketHashName = this.data.market_hash_name
    let currentQuality = ''
    
    // 找到当前饰品的品质
    for (const quality of relatedWears.qualities || []) {
      const wears = relatedWears.wears[quality] || []
      if (wears.some(w => w.hash_name === marketHashName)) {
        currentQuality = quality
        break
      }
    }
    
    if (!currentQuality) {
      currentQuality = relatedWears.current_quality || ''
    }

    // 当前品质的磨损列表
    let currentQualityWears = relatedWears.wears[currentQuality] || []
    
    // 添加显示标签和价格
    currentQualityWears = currentQualityWears.map(item => ({
      ...item,
      label: this.getWearLabel(item.wear),
      priceText: item.price > 0 ? item.price.toFixed(2) : '-',
      isActive: item.hash_name === marketHashName
    }))

    // 是否有磨损选项
    const hasWearOptions = currentQualityWears.length > 1 || 
      (currentQualityWears.length === 1 && currentQualityWears[0].wear !== 'NO_WEAR')

    // 是否有其他品质
    const hasOtherQuality = (relatedWears.qualities || []).length > 1
    
    // 其他品质标签
    let otherQualityLabel = ''
    let otherQualityMinPrice = 0
    
    if (hasOtherQuality) {
      const otherQ = (relatedWears.qualities || []).find(q => q !== currentQuality)
      if (otherQ) {
        const otherWears = relatedWears.wears[otherQ] || []
        
        // 根据名称判断显示标签
        if (otherWears.length > 0) {
          const otherHashName = otherWears[0].hash_name
          if (otherHashName.includes('StatTrak')) {
            otherQualityLabel = 'StatTrak™'
          } else if (otherHashName.startsWith('Souvenir ')) {
            otherQualityLabel = '纪念品'
          } else if (otherHashName.startsWith('★ ')) {
            otherQualityLabel = '★'
          } else {
            otherQualityLabel = '普通'
          }
        }
        
        // 找到对应磨损的价格
        const currentItem = currentQualityWears.find(w => w.isActive)
        const currentWear = currentItem?.wear || ''
        const matchingItem = otherWears.find(w => w.wear === currentWear)
        if (matchingItem && matchingItem.price > 0) {
          otherQualityMinPrice = matchingItem.price
        } else {
          const firstWithPrice = otherWears.find(w => w.price > 0)
          otherQualityMinPrice = firstWithPrice?.price || 0
        }
      }
    }

    this.setData({
      currentQuality,
      currentQualityWears,
      hasWearOptions,
      hasOtherQuality,
      otherQualityLabel,
      otherQualityMinPrice
    })
  },

  // 获取磨损标签
  getWearLabel(wear) {
    if (this.data.wearLabelMap[wear]) {
      return this.data.wearLabelMap[wear]
    }
    if (this.data.stickerVariantMap[wear]) {
      return this.data.stickerVariantMap[wear]
    }
    return wear
  },

  // 切换磨损
  switchWear(e) {
    const hashName = e.currentTarget.dataset.hashname
    if (!hashName || hashName === this.data.market_hash_name) return
    
    // 导航到新的详情页
    wx.redirectTo({
      url: `/pages/detail/detail?market_hash_name=${encodeURIComponent(hashName)}`
    })
  },

  // 切换品质
  toggleQuality() {
    const relatedWears = this.data.relatedWears
    if (!relatedWears) return
    
    const currentQ = this.data.currentQuality
    const otherQ = (relatedWears.qualities || []).find(q => q !== currentQ)
    if (!otherQ) return
    
    const otherWears = relatedWears.wears[otherQ] || []
    const currentItem = this.data.currentQualityWears.find(w => w.isActive)
    const currentWear = currentItem?.wear || ''
    
    // 找到对应磨损的饰品
    const targetItem = otherWears.find(w => w.wear === currentWear) || otherWears[0]
    
    if (targetItem) {
      wx.redirectTo({
        url: `/pages/detail/detail?market_hash_name=${encodeURIComponent(targetItem.hash_name)}`
      })
    }
  },

  // 初始化图表并绘制
  initChartAndDraw() {
    if (this.canvas) {
      this.drawChart()
      return
    }
    this.initChart()
  },

  // 处理详情数据
  processDetail(data) {
    let platformList = data.platformList || []
    if (platformList.length > 0) {
      const prices = platformList.filter(p => p.sellPrice > 0).map(p => p.sellPrice)
      const minPrice = prices.length > 0 ? Math.min(...prices) : 0
      
      platformList = platformList.map(p => ({
        ...p,
        sellPriceText: Number(p.sellPrice || 0).toFixed(2),
        biddingPriceText: Number(p.biddingPrice || 0).toFixed(2),
        isLowest: p.sellPrice === minPrice && p.sellPrice > 0,
        updateTimeText: this.formatUpdateTime(p.updateTime),
        icon: this.data.platformIconMap[p.platform] || ''
      }))
    }

    let priceChange = data.priceChange || []
    if (priceChange.length > 0) {
      priceChange = priceChange.map(item => ({
        ...item,
        priceDiffText: Math.abs(item.priceDiff).toFixed(2),
        changeRateText: (item.changeRate >= 0 ? '+' : '') + item.changeRate.toFixed(2) + '%'
      }))
    }

    return {
      ...data,
      platformList,
      priceChange,
      priceHistory: data.priceHistory || {}
    }
  },

  // 初始化图表
  initChart() {
    const query = wx.createSelectorQuery()
    query.select('#priceChart')
      .fields({ node: true, size: true })
      .exec((res) => {
        if (res[0] && res[0].node) {
          const canvas = res[0].node
          const ctx = canvas.getContext('2d')
          
          let width = res[0].width
          let height = res[0].height
          
          if (!width || !height) {
            const systemInfo = wx.getSystemInfoSync()
            width = systemInfo.windowWidth - 56
            height = 200
          }
          
          const dpr = wx.getSystemInfoSync().pixelRatio
          canvas.width = width * dpr
          canvas.height = height * dpr
          ctx.scale(dpr, dpr)
          
          this.canvas = canvas
          this.ctx = ctx
          this.canvasWidth = width
          this.canvasHeight = height
          
          if (this.data.goodsDetail) {
            this.drawChart()
          }
        } else {
          setTimeout(() => {
            this.initChart()
          }, 200)
        }
      })
  },

  // 绘制图表
  drawChart() {
    if (!this.ctx || !this.data.goodsDetail) return

    const ctx = this.ctx
    const width = this.canvasWidth
    const height = this.canvasHeight
    const padding = { top: 30, right: 20, bottom: 40, left: 60 }

    const platform = this.data.platforms[this.data.selectedPlatformIndex].value
    const priceHistory = this.data.goodsDetail.priceHistory || {}
    const historyData = priceHistory[platform] || []

    ctx.clearRect(0, 0, width, height)

    if (historyData.length === 0) {
      ctx.fillStyle = '#999'
      ctx.font = '14px sans-serif'
      ctx.textAlign = 'center'
      ctx.fillText('暂无该平台数据', width / 2, height / 2)
      this.chartInfo = null
      return
    }

    const dates = historyData.map(item => item.date)
    const prices = historyData.map(item => item.sellPrice)
    const counts = historyData.map(item => item.sellCount)

    const minPrice = Math.min(...prices) * 0.95
    const maxPrice = Math.max(...prices) * 1.05
    const maxCount = Math.max(...counts) * 1.1

    const chartWidth = width - padding.left - padding.right
    const chartHeight = height - padding.top - padding.bottom
    
    // 保存图表信息供触摸事件使用
    this.chartInfo = {
      historyData,
      padding,
      chartWidth,
      chartHeight,
      minPrice,
      maxPrice,
      maxCount
    }

    // 绘制背景网格
    ctx.strokeStyle = '#f0f0f0'
    ctx.lineWidth = 1
    for (let i = 0; i <= 4; i++) {
      const y = padding.top + (chartHeight / 4) * i
      ctx.beginPath()
      ctx.moveTo(padding.left, y)
      ctx.lineTo(width - padding.right, y)
      ctx.stroke()
    }

    // Y轴标签
    ctx.fillStyle = '#999'
    ctx.font = '11px sans-serif'
    ctx.textAlign = 'right'
    for (let i = 0; i <= 4; i++) {
      const price = maxPrice - ((maxPrice - minPrice) / 4) * i
      const y = padding.top + (chartHeight / 4) * i
      ctx.fillText('¥' + price.toFixed(0), padding.left - 8, y + 4)
    }

    // X轴标签
    ctx.textAlign = 'center'
    const labelCount = Math.min(5, dates.length)
    const step = Math.floor(dates.length / labelCount)
    for (let i = 0; i < labelCount; i++) {
      const index = i * step
      if (index < dates.length) {
        const x = padding.left + (chartWidth / (dates.length - 1)) * index
        const date = dates[index].slice(5)
        ctx.fillText(date, x, height - 10)
      }
    }

    // 价格折线
    ctx.strokeStyle = '#f97316'
    ctx.lineWidth = 2
    ctx.beginPath()
    prices.forEach((price, index) => {
      const x = padding.left + (chartWidth / (prices.length - 1)) * index
      const y = padding.top + chartHeight - ((price - minPrice) / (maxPrice - minPrice)) * chartHeight
      if (index === 0) {
        ctx.moveTo(x, y)
      } else {
        ctx.lineTo(x, y)
      }
    })
    ctx.stroke()

    // 价格区域填充
    ctx.fillStyle = 'rgba(249, 115, 22, 0.1)'
    ctx.beginPath()
    prices.forEach((price, index) => {
      const x = padding.left + (chartWidth / (prices.length - 1)) * index
      const y = padding.top + chartHeight - ((price - minPrice) / (maxPrice - minPrice)) * chartHeight
      if (index === 0) {
        ctx.moveTo(x, y)
      } else {
        ctx.lineTo(x, y)
      }
    })
    ctx.lineTo(padding.left + chartWidth, padding.top + chartHeight)
    ctx.lineTo(padding.left, padding.top + chartHeight)
    ctx.closePath()
    ctx.fill()

    // 在售数量折线
    ctx.strokeStyle = '#1890ff'
    ctx.lineWidth = 2
    ctx.beginPath()
    counts.forEach((count, index) => {
      const x = padding.left + (chartWidth / (counts.length - 1)) * index
      const y = padding.top + chartHeight - (count / maxCount) * chartHeight
      if (index === 0) {
        ctx.moveTo(x, y)
      } else {
        ctx.lineTo(x, y)
      }
    })
    ctx.stroke()

    // 图例
    const legendY = 12
    ctx.fillStyle = '#f97316'
    ctx.fillRect(padding.left, legendY, 20, 3)
    ctx.fillStyle = '#666'
    ctx.font = '12px sans-serif'
    ctx.textAlign = 'left'
    ctx.fillText('价格', padding.left + 25, legendY + 4)
    
    ctx.fillStyle = '#1890ff'
    ctx.fillRect(padding.left + 70, legendY, 20, 3)
    ctx.fillStyle = '#666'
    ctx.fillText('在售数量', padding.left + 95, legendY + 4)
  },

  // 切换平台
  onPlatformChange(e) {
    this.setData({
      selectedPlatformIndex: e.detail.value
    })
    this.drawChart()
  },

  // 切换时间范围
  async onTimeChange(e) {
    const index = parseInt(e.detail.value)
    if (index === this.data.selectedTimeIndex) return
    
    const days = this.data.timeRanges[index].value
    
    this.setData({ selectedTimeIndex: index, loading: true })
    
    try {
      const res = await api.getGoodsDetail({
        market_hash_name: this.data.market_hash_name,
        days: days
      })
      
      if (res.code === 1 && res.data) {
        const detail = this.processDetail(res.data)
        console.log('处理后detail.priceHistory:', detail.priceHistory)
        
        // loading: false 后 canvas 会重新创建，必须重新初始化
        // 先清除旧的 canvas 引用
        this.canvas = null
        this.ctx = null
        
        this.setData({ goodsDetail: detail, loading: false }, () => {
          // 等待 DOM 渲染完成后重新初始化 canvas
          setTimeout(() => {
            this.initChart()
          }, 100)
        })
      } else {
        this.setData({ loading: false })
      }
    } catch (error) {
      console.error('切换时间范围失败:', error)
      this.setData({ loading: false })
      wx.showToast({ title: '加载失败', icon: 'none' })
    }
  },

  // 格式化更新时间
  formatUpdateTime(timestamp) {
    if (!timestamp) return ''
    const now = Date.now()
    const diff = now - timestamp * 1000
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    
    if (minutes < 1) return '刚刚'
    if (minutes < 60) return `${minutes}分钟前`
    if (hours < 24) return `${hours}小时前`
    return `${Math.floor(hours / 24)}天前`
  },

  // 获取品质颜色
  getRarityColor(rarity) {
    const colorMap = {
      '隐秘': '#EB4B4B',
      '保密': '#D32CE6',
      '受限': '#8B008B',
      '军规级': '#4B69FF',
      '工业级': '#5E98D9',
      '消费级': '#808080'
    }
    return colorMap[rarity] || '#333333'
  },

  // 返回上一页
  goBack() {
    wx.navigateBack()
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.fetchDetail().then(() => {
      wx.stopPullDownRefresh()
    })
  },

  // 图表触摸事件
  onChartTouch(e) {
    if (!this.chartInfo || !this.chartInfo.historyData || this.chartInfo.historyData.length === 0) return

    const touch = e.touches[0]
    const { historyData, padding, chartWidth, chartHeight } = this.chartInfo
    
    // 获取canvas相对于页面的位置
    const query = wx.createSelectorQuery()
    query.select('#priceChart').boundingClientRect((rect) => {
      if (!rect) return
      
      const x = touch.clientX - rect.left
      const y = touch.clientY - rect.top
      
      // 判断是否在图表区域内
      if (x < padding.left || x > padding.left + chartWidth) {
        this.setData({ showTooltip: false })
        return
      }
      
      // 计算对应的数据索引
      const dataIndex = Math.round((x - padding.left) / chartWidth * (historyData.length - 1))
      const clampedIndex = Math.max(0, Math.min(dataIndex, historyData.length - 1))
      
      const dataPoint = historyData[clampedIndex]
      if (!dataPoint) return
      
      // 计算tooltip位置
      let tooltipX = x
      let tooltipY = 10
      
      // 边界处理
      if (tooltipX < 80) tooltipX = 80
      if (tooltipX > rect.width - 80) tooltipX = rect.width - 80
      
      this.setData({
        showTooltip: true,
        tooltipX: tooltipX,
        tooltipY: tooltipY,
        tooltipData: {
          date: dataPoint.date,
          price: dataPoint.sellPrice.toFixed(2),
          count: dataPoint.sellCount
        }
      })
      
      // 绘制选中点
      this.drawSelectedPoint(clampedIndex)
    }).exec()
  },

  // 触摸结束隐藏tooltip
  onChartTouchEnd() {
    setTimeout(() => {
      this.setData({ showTooltip: false })
      // 重新绘制图表（移除选中点）
      this.drawChart()
    }, 1500)
  },

  // 绘制选中点
  drawSelectedPoint(index) {
    if (!this.ctx || !this.chartInfo) return
    
    const { historyData, padding, chartWidth, chartHeight, minPrice, maxPrice, maxCount } = this.chartInfo
    const dataPoint = historyData[index]
    if (!dataPoint) return
    
    // 先重绘图表
    this.drawChart()
    
    const ctx = this.ctx
    const x = padding.left + (chartWidth / (historyData.length - 1)) * index
    const priceY = padding.top + chartHeight - ((dataPoint.sellPrice - minPrice) / (maxPrice - minPrice)) * chartHeight
    const countY = padding.top + chartHeight - (dataPoint.sellCount / maxCount) * chartHeight
    
    // 绘制垂直线
    ctx.strokeStyle = 'rgba(0, 0, 0, 0.2)'
    ctx.lineWidth = 1
    ctx.setLineDash([4, 4])
    ctx.beginPath()
    ctx.moveTo(x, padding.top)
    ctx.lineTo(x, padding.top + chartHeight)
    ctx.stroke()
    ctx.setLineDash([])
    
    // 绘制价格点
    ctx.fillStyle = '#f97316'
    ctx.beginPath()
    ctx.arc(x, priceY, 5, 0, Math.PI * 2)
    ctx.fill()
    ctx.strokeStyle = '#fff'
    ctx.lineWidth = 2
    ctx.stroke()
    
    // 绘制数量点
    ctx.fillStyle = '#1890ff'
    ctx.beginPath()
    ctx.arc(x, countY, 5, 0, Math.PI * 2)
    ctx.fill()
    ctx.strokeStyle = '#fff'
    ctx.lineWidth = 2
    ctx.stroke()
  }
})
