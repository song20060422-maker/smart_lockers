// @ticket 07-数字人Coze：Vite 构建配置（独立子模块）
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [react()],
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  server: {
    port: 5174, // 与主 Vue3 前端 5173 隔离
    proxy: {
      '/coze': { target: 'https://api.coze.cn', changeOrigin: true, rewrite: (p) => p.replace(/^\/coze/, '') },
      '/ws': { target: 'ws://localhost:18080', ws: true },
    },
  },
  build: { outDir: 'dist', sourcemap: false },
})
