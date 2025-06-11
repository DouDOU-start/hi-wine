import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

// 环境变量配置
const env = process.env.NODE_ENV || 'development';
const apiBaseUrl = process.env.VITE_API_BASE_URL || 'http://localhost:8000';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: apiBaseUrl,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api/v1'),
        ws: true
      },
    },
  },
  // 定义环境变量
  define: {
    'process.env': {
      NODE_ENV: JSON.stringify(env),
      VITE_API_BASE_URL: JSON.stringify(apiBaseUrl),
    }
  },
}); 