import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      "/remove": {
        target: "http://127.0.0.1:3000/",
        changeOrigin: true,
        secure: false,
      },
      "/log/": {
        target: "http://127.0.0.1:3000/",
        changeOrigin: true,
        secure: false,
      },
      "/webfinger/": {
        target: "http://127.0.0.1:3000/",
        changeOrigin: true,
        secure: false,
      }
    },
  },
})
