import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    chunkSizeWarningLimit: 1500,
    // 代码分割优化
    rollupOptions: {
      output: {
        // 手动分包优化
        manualChunks: {
          // Vue 核心库
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          // Element Plus UI 库
          'element-plus': ['element-plus', '@element-plus/icons-vue'],
          // 其他第三方库
          'vendor': ['axios', 'dayjs', 'echarts', 'vee-validate', 'yup']
        },
        // chunk 文件命名
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: '[ext]/[name]-[hash].[ext]'
      }
    },
    // 启用 CSS 代码分割
    cssCodeSplit: true,
    // 预加载优化
    modulePreload: {
      polyfill: true
    }
  },
  // 定义全局常量
  define: {
    __USE_MOCK__: process.env.USE_MOCK === 'true'
  }
})
