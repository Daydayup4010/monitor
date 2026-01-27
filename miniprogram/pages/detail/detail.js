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
    }
  },

  onLoad(options) {
    if (options.market_hash_name) {
      this.setData({
        market_hash_name: decodeURIComponent(options.market_hash_name)
      })
      this.fetchDetail()
    }
  },

  // 获取详情数据
  async fetchDetail() {
    this.setData({ loading: true })
    
    try {
      const res = await api.getGoodsDetail({
        market_hash_name: this.data.market_hash_name,
        days: 30
      })
      
      if (res.code === 1 && res.data) {
        const detail = this.processDetail(res.data)
        this.setData({
          goodsDetail: detail
        })
        // 设置页面标题
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

  // 处理详情数据
  processDetail(data) {
    // 处理平台列表
    let platformList = data.platformList || []
    if (platformList.length > 0) {
      // 找到最低价
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

    // 处理价格涨幅
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
      priceChange
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

  // 跳转到平台购买
  openPlatformLink(e) {
    const link = e.currentTarget.dataset.link
    if (link) {
      // 复制链接到剪贴板
      wx.setClipboardData({
        data: link,
        success: () => {
          wx.showToast({
            title: '链接已复制',
            icon: 'success'
          })
        }
      })
    }
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
  }
})
