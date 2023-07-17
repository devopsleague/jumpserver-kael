import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import WindiCSS from "vite-plugin-windicss"
import path from 'path'
const resolve = (dir) => path.join(__dirname, dir)

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())

  return {
    base: '/kael/',
    plugins: [vue(), WindiCSS()],
    resolve: {
      alias: {
        '@': resolve('src')
      }
    },
    define: {
      'process.env': env
    },
    server: {
      cors: true,
      open: true,
      proxy: {
        '/kael/interrupt_current_ask': {
          target: env.VITE_APP_BASE_URL,
          changeOrigin: true
        }
      }
    }
  }
})
