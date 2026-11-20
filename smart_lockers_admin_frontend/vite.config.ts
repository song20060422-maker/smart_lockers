import { defineConfig, loadEnv } from 'vite' // @ticket B端-前端：Vite 构建配置
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// @ticket B端-前端：开发态把 /api 代理到 B 端 Gin 后端 :18081
// 导出 defineConfig(({ mode }) => {，供其他模块引用

export default defineConfig(({ mode }) => {
  // 声明变量 env = loadEnv(mode, process.cwd(), '')，用于保存数据或函数引用

  const env = loadEnv(mode, process.cwd(), '')
  // 声明变量 apiTarget = env.VITE_API_TARGET || 'http://localhost:18081'，用于保存数据或函数引用

  const apiTarget = env.VITE_API_TARGET || 'http://localhost:18081'
  // 返回 { 给调用方

  return {
    // 定义对象属性 plugins，取值为 [vue()]

    plugins: [vue()],
    // 定义对象属性 resolve，取值为 { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } }

    resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
    // 定义对象属性 server，取值为 {

    server: {
      // 定义对象属性 port，取值为 5175

      port: 5175,
      // 定义对象属性 proxy，取值为 { '/api': { target: apiTarget, changeOrigin: true } }

      proxy: { '/api': { target: apiTarget, changeOrigin: true } },
    // 定义结束

    },
    // 定义对象属性 build，取值为 { outDir: 'dist', sourcemap: false }

    build: { outDir: 'dist', sourcemap: false },
  // 当前代码块结束

  }
// 执行 })，实现对应的前端逻辑

})
