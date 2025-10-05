import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    // vvvvvvvvvvvv 添加这行 vvvvvvvvvvvv
    host: true, 
    // ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
    proxy: {
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/static': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      }
    }
  }
})