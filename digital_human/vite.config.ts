// @ticket 07-数字人Coze：Vite 构建配置（独立子模块）
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath, URL } from 'node:url'

// 导出 defineConfig({，供其他模块引用

export default defineConfig({
  // 定义对象属性 plugins，取值为 [react()]

  plugins: [react()],
  // 定义对象属性 resolve，取值为 { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } }

  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  // 定义对象属性 server，取值为 {

  server: {
    // 定义对象属性 port，取值为 5174, // 与主 Vue3 前端 5173 隔离

    port: 5174, // 与主 Vue3 前端 5173 隔离
    // 定义对象属性 proxy，取值为 {

    proxy: {
      // 执行 '/coze': { target: 'https://api.coze.cn', changeOrigin: true, rewrite: (p) => p.replace(/^\/coze/, '') },，实现对应的前端逻辑

      '/coze': { target: 'https://api.coze.cn', changeOrigin: true, rewrite: (p) => p.replace(/^\/coze/, '') },
      // 执行 '/ws': { target: 'ws://localhost:18080', ws: true },，实现对应的前端逻辑

      '/ws': { target: 'ws://localhost:18080', ws: true },
    // 定义结束

    },
  // 定义结束

  },
  // 定义对象属性 build，取值为 { outDir: 'dist', sourcemap: false }

  build: { outDir: 'dist', sourcemap: false },
// 执行 })，实现对应的前端逻辑

})
