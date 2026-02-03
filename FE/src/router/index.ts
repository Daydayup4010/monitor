import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layouts/Layout.vue'
import { useUserStore } from '@/stores/user'
import { showMessage } from '@/utils/message'

const router = createRouter({
  history: createWebHistory('/'), // 去掉csgo前缀
  routes: [
    {
      path: '/',
      name: 'Landing',
      component: () => import('@/pages/Landing.vue'),
      meta: { title: 'CS2饰品搬砖平台', requiresAuth: false }
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/pages/Login.vue'),
      meta: { title: '登录', requiresAuth: false }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/pages/Register.vue'),
      meta: { title: '注册', requiresAuth: false }
    },
    {
      path: '/reset-password',
      name: 'ResetPassword',
      component: () => import('@/pages/ResetPassword.vue'),
      meta: { title: '重置密码', requiresAuth: false, hideInMenu: true }
    },
    {
      path: '/privacy-policy',
      name: 'PrivacyPolicy',
      component: () => import('@/pages/PrivacyPolicy.vue'),
      meta: { title: '隐私政策', requiresAuth: false, hideInMenu: true }
    },
    {
      path: '/user-agreement',
      name: 'UserAgreement',
      component: () => import('@/pages/UserAgreement.vue'),
      meta: { title: '用户协议', requiresAuth: false, hideInMenu: true }
    },
    {
      path: '/app',
      component: Layout,
      redirect: '/app/dashboard',
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/pages/Dashboard.vue'),
          meta: { 
            title: '首页', 
            icon: 'HomeFilled',
            requiresAuth: true,
            requiresVip: true
          }
        },
        {
          path: 'ranking',
          name: 'Ranking',
          component: () => import('@/pages/Ranking.vue'),
          meta: { 
            title: '饰品榜单', 
            icon: 'TrendCharts',
            requiresAuth: true,
            requiresVip: true
          }
        },
        {
          path: 'home',
          name: 'Home',
          component: () => import('@/pages/Home.vue'),
          meta: { 
            title: '挂刀/搬砖', 
            icon: 'DataAnalysis',
            requiresAuth: true,
            requiresVip: true
          }
        },
        {
          path: 'big-item-bidding',
          name: 'BigItemBidding',
          component: () => import('@/pages/BigItemBidding.vue'),
          meta: { 
            title: '大件求购', 
            icon: 'Coin',
            requiresAuth: true,
            requiresVip: true,
            isNew: true
          }
        },
        {
          path: 'detail',
          name: 'Detail',
          component: () => import('@/pages/Detail.vue'),
          meta: { 
            title: '饰品详情', 
            icon: 'InfoFilled',
            requiresAuth: true,
            hideInMenu: true  // 不在导航菜单中显示
          }
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('@/pages/Settings.vue'),
          meta: { 
            title: '个人设置', 
            icon: 'User',
            requiresAuth: true,
            hideInMenu: true  // 不在导航菜单中显示
          }
        },
        {
          path: 'notifications',
          name: 'NotificationCenter',
          component: () => import('@/pages/NotificationCenter.vue'),
          meta: { 
            title: '消息中心', 
            icon: 'Bell',
            requiresAuth: true,
            hideInMenu: true  // 不在导航菜单中显示
          }
        },
        {
          path: 'vip',
          name: 'VipService',
          component: () => import('@/pages/Settings.vue'),
          meta: { 
            title: 'VIP服务', 
            icon: 'Medal',
            requiresAuth: true,
            defaultTab: 'vip'  // 用于设置默认tab
            // 不需要requiresVip，让非VIP用户也能访问
          }
        },
        {
          path: 'admin',
          name: 'Admin',
          component: () => import('@/pages/Admin.vue'),
          meta: { 
            title: '管理中心', 
            icon: 'Setting',
            requiresAuth: true,
            requiresAdmin: true  // 需要管理员权限
          }
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  
  // 动态设置页面标题
  const title = to.meta.title as string
  document.title = title ? `${title} - CS Goods` : 'CS Goods'
  
  // 如果已登录且访问登录/注册页面，根据VIP状态跳转
  if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
    // VIP用户跳转到dashboard，非VIP用户跳转到公开首页
    next(userStore.isVip ? '/app/dashboard' : '/')
    return
  }
  
  // 公开页面允许所有人访问（不需要登录）
  const publicPages = ['/', '/reset-password', '/privacy-policy', '/user-agreement']
  if (publicPages.includes(to.path)) {
    next()
    return
  }
  
  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    if (!userStore.isLoggedIn) {
      showMessage.warning('请先登录')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }
    
    // 检查VIP权限
    if (to.meta.requiresVip && !userStore.isVip) {
      showMessage.warning('需要VIP权限才能访问该功能')
      next('/')  // 非VIP用户重定向到公开首页
      return
    }
    
    // 检查管理员权限
    if (to.meta.requiresAdmin && !userStore.isAdmin) {
      showMessage.error('需要管理员权限才能访问')
      next('/')
      return
    }
  }
  
  next()
})

export default router