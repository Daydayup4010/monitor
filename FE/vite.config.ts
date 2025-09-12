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
    hmr: {
      port: 3001, // HMR WebSocket 端口
      // 确保HMR使用正确的主机地址
      host: 'localhost'
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
  // 只在构建时使用base路径，开发时不使用
  base: process.env.NODE_ENV === 'production' ? '/csgo/' : '/',
  define: {
    // 定义 Vue 3 特性标志
    __VUE_OPTIONS_API__: true,
    __VUE_PROD_DEVTOOLS__: false,
    __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
  }
})