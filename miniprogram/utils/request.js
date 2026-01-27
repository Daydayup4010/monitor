// 网络请求封装
const app = getApp()

const request = (url, options = {}) => {
  return new Promise((resolve, reject) => {
    const token = app.globalData.token || wx.getStorageSync('token')
    
    // 判断是否是公开接口（不需要登录）
    const isPublicApi = url.startsWith('/public/') || url === '/wechat/login'
    
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
          } else if ((data.code === 1005 || data.code === 1006) && !isPublicApi) {
            // Token过期或无效（公开接口不处理）
            app.clearLoginInfo()
            wx.navigateTo({
              url: '/pages/login/login'
            })
            reject(data)
          } else if (data.code === 1026) {
            // 账号已在其他设备登录（单设备登录限制）
            wx.showToast({
              title: '账号已在其他设备登录',
              icon: 'none',
              duration: 2000
            })
            app.clearLoginInfo()
            setTimeout(() => {
              wx.reLaunch({
                url: '/pages/login/login'
              })
            }, 1500)
            reject(data)
          } else if (data.code === 1021) {
            // 邮箱已存在需要合并，不显示toast，让页面自己处理
            resolve(data)
          } else {
            wx.showToast({
              title: data.msg || '请求失败',
              icon: 'none'
            })
            reject(data)
          }
        } else if (res.statusCode === 403) {
          // 403 权限不足，提示开通VIP
          wx.showModal({
            title: '提示',
            content: '暂无权限访问，请先开通VIP',
            confirmText: '去开通',
            cancelText: '取消',
            success(modal) {
              if (modal.confirm) {
                wx.navigateTo({
                  url: '/pages/vip/vip'
                })
              }
            }
          })
          reject({
            code: 403,
            msg: '暂无权限访问，请先开通VIP',
            statusCode: res.statusCode
          })
        } else if (res.statusCode === 401) {
          // 401 未授权
          app.clearLoginInfo()
          wx.navigateTo({
            url: '/pages/login/login'
          })
          reject({
            code: 401,
            msg: '请先登录',
            statusCode: res.statusCode
          })
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

