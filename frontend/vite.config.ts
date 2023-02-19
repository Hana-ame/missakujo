import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    sourcemap: true,
  },
  server: {
    proxy: {
      "/api-missakujo/delete": {
        target: "http://127.0.0.1:3000/",
        changeOrigin: true,
        secure: false,
      },
      "/api-missakujo/log/": {
        target: "http://127.0.0.1:3000/",
        changeOrigin: true,
        secure: false,
      },
      "/api-missakujo/webfinger/": {
        target: "http://127.0.0.1:3000/",
        changeOrigin: true,
        secure: false,
      }
    },
  },
})
