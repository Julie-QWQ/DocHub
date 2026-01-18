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
        // 动态代码分割 - 更细粒度的分包
        manualChunks(id) {
          // Vue 核心库
          if (id.includes('node_modules/vue/') || id.includes('node_modules/@vue/') || id.includes('node_modules/pinia/')) {
            return 'vue-vendor'
          }
          // Element Plus UI 库
          if (id.includes('node_modules/element-plus/') || id.includes('node_modules/@element-plus/')) {
            return 'element-plus'
          }
          // 其他第三方库
          if (id.includes('node_modules/')) {
            return 'vendor'
          }
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
