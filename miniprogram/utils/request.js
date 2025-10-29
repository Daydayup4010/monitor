// 网络请求封装
const app = getApp()

const request = (url, options = {}) => {
  return new Promise((resolve, reject) => {
    const token = app.globalData.token || wx.getStorageSync('token')
    
    wx.request({
      url: `${app.globalData.baseURL}${url}`,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      },
      success(res) {
        if (res.statusCode === 200) {
          const data = res.data
          
          // 统一错误处理
          if (data.code === 1) {
            resolve(data)
          } else if (data.code === 1005 || data.code === 1006) {
            // Token过期或无效
            app.clearLoginInfo()
            wx.reLaunch({
              url: '/pages/login/login'
            })
            reject(data)
          } else {
            wx.showToast({
              title: data.msg || '请求失败',
              icon: 'none'
            })
            reject(data)
          }
        } else {
          wx.showToast({
            title: '网络错误',
            icon: 'none'
          })
          reject(res)
        }
      },
      fail(err) {
        wx.showToast({
          title: '网络请求失败',
          icon: 'none'
        })
        reject(err)
      }
    })
  })
}

module.exports = {
  get: (url, data) => request(url, { method: 'GET', data }),
  post: (url, data) => request(url, { method: 'POST', data }),
  put: (url, data) => request(url, { method: 'PUT', data }),
  delete: (url, data) => request(url, { method: 'DELETE', data }),
}

