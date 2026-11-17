import { defineConfig, loadEnv } from 'vite' // @ticket 00-project：Vite 构建配置
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// @ticket 00-project：Vite 配置——插件、路径别名、开发代理（/api -> 后端 :18080）
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiTarget = env.VITE_API_TARGET || 'http://localhost:18080'
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      port: 5173,
      proxy: {
        // @ticket 00-project：开发态把 /api 代理到 Gin 后端，避免跨域
        '/api': {
          target: apiTarget,
          changeOrigin: true,
        },
      },
    },
    build: {
      outDir: 'dist',
      sourcemap: false,
    },
  }
})
