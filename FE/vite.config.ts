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
    hmr: false, // 生产环境禁用HMR
    proxy: {
      '/api': {
        target: 'http://localhost:3100',
        changeOrigin: true,
      }
    }
  },
  // 基础路径
  base: '/',
  build: {
    // 生产构建优化
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    minify: 'terser',
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          elementPlus: ['element-plus']
        }
      }
    }
  },
  define: {
    // 定义 Vue 3 特性标志
    __VUE_OPTIONS_API__: true,
    __VUE_PROD_DEVTOOLS__: false,
    __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
  }
})