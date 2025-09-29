import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layouts/Layout.vue'

const router = createRouter({
  history: createWebHistory('/csgo/'), // 设置基础路径，匹配 nginx 配置
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