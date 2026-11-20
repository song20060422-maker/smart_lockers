import { defineConfig, loadEnv } from 'vite' // @ticket B端-前端：Vite 构建配置
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// @ticket B端-前端：开发态把 /api 代理到 B 端 Gin 后端 :18081
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiTarget = env.VITE_API_TARGET || 'http://localhost:18081'
  return {
    plugins: [vue()],
    resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
    server: {
      port: 5175,
      proxy: { '/api': { target: apiTarget, changeOrigin: true } },
    },
    build: { outDir: 'dist', sourcemap: false },
  }
})
