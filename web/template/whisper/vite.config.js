import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
// https://juejin.cn/post/7049665818612269070
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    //设置 server.hmr.overlay 为 false 可以禁用开发服务器错误的屏蔽
    // hmr: { overlay: false },
    host: '0.0.0.0', // 服务器监听的IP地址
    port: 80, // 服务器监听的端口号
  },
})
