// pages/settings/settings.js
const api = require('../../utils/api.js')
const app = getApp()

Page({
  data: {
    userInfo: {},
    settings: {
      min_sell_num: 0,
      min_diff: 0,
      min_sell_price: 0,
      max_sell_price: 10000
    },
    saving: false
  },

  onLoad() {
    this.loadUserInfo()
    this.loadSettings()
  },

  onShow() {
    this.loadUserInfo()
  },

  // 加载用户信息
  loadUserInfo() {
    const userInfo = app.globalData.userInfo || {}
    userInfo.isVip = userInfo.role === 1 || userInfo.role === 2
    this.setData({ userInfo })
  },

  // 加载设置
  async loadSettings() {
    if (!app.globalData.token) return

    try {
      const res = await api.getSettings()
      if (res.code === 1 && res.data) {
        this.setData({
          settings: res.data
        })
      }
    } catch (error) {
      console.error('加载设置失败:', error)
    }
  },

  onMinSellNumInput(e) {
    this.setData({
      'settings.min_sell_num': Number(e.detail.value)
    })
  },

  onMinDiffInput(e) {
    this.setData({
      'settings.min_diff': Number(e.detail.value)
    })
  },

  onMinPriceInput(e) {
    this.setData({
      'settings.min_sell_price': Number(e.detail.value)
    })
  },

  onMaxPriceInput(e) {
    this.setData({
      'settings.max_sell_price': Number(e.detail.value)
    })
  },

  // 保存设置
  async handleSave() {
    this.setData({ saving: true })

    try {
      const res = await api.updateSettings(this.data.settings)
      
      if (res.code === 1) {
        wx.showToast({
          title: '保存成功',
          icon: 'success'
        })
      }
    } catch (error) {
      console.error('保存失败:', error)
    } finally {
      this.setData({ saving: false })
    }
  },

  // 重置
  handleReset() {
    const that = this
    wx.showModal({
      title: '确认重置',
      content: '确定要恢复默认设置吗？',
      success(res) {
        if (res.confirm) {
          that.setData({
            settings: {
              min_sell_num: 200,
              min_diff: 1,
              min_sell_price: 0,
              max_sell_price: 10000
            }
          })
          wx.showToast({
            title: '已恢复默认设置',
            icon: 'none'
          })
        }
      }
    })
  }
})



