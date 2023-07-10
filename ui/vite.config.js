import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import WindiCSS from "vite-plugin-windicss"
import path from 'path'
const resolve = (dir) => path.join(__dirname, dir)

export default defineConfig({
  plugins: [vue(), WindiCSS()],
  resolve: {
    alias: {
      '@': resolve('src')
    }
  },
  define: {
    'process.env': {
      VITE_APP_API_BASE_URL: 'http://127.0.0.1:8083',
      VITE_APP_BASE_SOCKET: 'ws://127.0.0.1:8083/kael/chat'
    }
  },
  server: {
    cors: true,
    open: true,
    proxy: {
      '/kael': {
        target: process.env.VITE_APP_API_BASE_URL,
        changeOrigin: true
      }
    }
  }
})
