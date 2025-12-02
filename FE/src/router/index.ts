import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layouts/Layout.vue'
import { useUserStore } from '@/stores/user'
import { showMessage } from '@/utils/message'

const router = createRouter({
  history: createWebHistory('/'), // 去掉csgo前缀
  routes: [
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
      path: '/',
      component: Layout,
      redirect: '/home',
      meta: { requiresAuth: true },
      children: [
        {
          path: 'home',
          name: 'Home',
          component: () => import('@/pages/Home.vue'),
          meta: { 
            title: '饰品数据', 
            icon: 'DataAnalysis',
            requiresAuth: true,
            requiresVip: true,  // 需要VIP权限
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
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 动态设置页面标题
  const title = to.meta.title as string
  document.title = title ? `${title} - CS Goods` : 'CS Goods'
  
  // 如果已登录且访问登录/注册页，重定向到首页
  if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
    next('/home')
    return
  }
  
  // 公开页面允许所有人访问（不需要登录）
  const publicPages = ['/reset-password', '/privacy-policy', '/user-agreement']
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
      showMessage.warning('需要VIP权限才能访问')
      next('/settings')
      return
    }
    
    // 检查管理员权限
    if (to.meta.requiresAdmin && !userStore.isAdmin) {
      showMessage.error('需要管理员权限才能访问')
      next('/home')
      return
    }
  }
  
  next()
})

export default router