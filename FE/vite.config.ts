import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3000,
    host: '0.0.0.0', // 允许外部访问
    // 允许所有主机访问（开发环境）
    allowedHosts: 'all',
    hmr: {
      port: 3001, // HMR WebSocket 端口
      // 设置HMR主机为当前域名
      host: 'www.2333tv.top'
    },
    proxy: {
      '/api': {
        target: 'http://localhost:3601',
        changeOrigin: true,
      }
    }
  },
  // 生产环境使用base路径
  base: '/csgo/',
  define: {
    // 定义 Vue 3 特性标志
    __VUE_OPTIONS_API__: true,
    __VUE_PROD_DEVTOOLS__: false,
    __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
  }
})