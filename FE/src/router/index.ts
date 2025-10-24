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
            title: '设置', 
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
  
  // 如果已登录且访问登录/注册页，重定向到首页
  if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
    next('/home')
    return
  }
  
  // 重置密码页面允许所有人访问
  if (to.path === '/reset-password') {
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