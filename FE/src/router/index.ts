import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layouts/Layout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL), // 使用Vite的BASE_URL环境变量
  routes: [
    {
      path: '/',
      component: Layout,
      redirect: '/home',
      children: [
        {
          path: 'home',
          name: 'Home',
          component: () => import('@/pages/Home.vue'),
          meta: { title: '饰品数据', icon: 'DataAnalysis' }
        },
        {
          path: 'admin',
          name: 'Admin',
          component: () => import('@/pages/Admin.vue'),
          meta: { title: '管理中心', icon: 'Setting' }
        }
      ]
    }
  ]
})

export default router