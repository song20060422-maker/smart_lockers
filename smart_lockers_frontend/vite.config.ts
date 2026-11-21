import { defineConfig, loadEnv } from 'vite' // @ticket 00-project：Vite 构建配置
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// @ticket 00-project：Vite 配置——插件、路径别名、开发代理（/api -> 后端 :18080）
// 导出 defineConfig(({ mode }) => {，供其他模块引用

export default defineConfig(({ mode }) => {
  // 声明变量 env = loadEnv(mode, process.cwd(), '')，用于保存数据或函数引用

  const env = loadEnv(mode, process.cwd(), '')
  // 声明变量 apiTarget = env.VITE_API_TARGET || 'http://localhost:18080'，用于保存数据或函数引用

  const apiTarget = env.VITE_API_TARGET || 'http://localhost:18080'
  // 返回 { 给调用方

  return {
    // 定义对象属性 plugins，取值为 [vue()]

    plugins: [vue()],
    // 定义对象属性 resolve，取值为 {

    resolve: {
      // 定义对象属性 alias，取值为 {

      alias: {
        // 执行 '@': fileURLToPath(new URL('./src', import.meta.url)),，实现对应的前端逻辑

        '@': fileURLToPath(new URL('./src', import.meta.url)),
      // 定义结束

      },
    // 定义结束

    },
    // 定义对象属性 server，取值为 {

    server: {
      // 定义对象属性 port，取值为 5173

      port: 5173,
      // 定义对象属性 proxy，取值为 {

      proxy: {
        // @ticket 00-project：开发态把 /api 代理到 Gin 后端，避免跨域
        // 执行 '/api': {，实现对应的前端逻辑

        '/api': {
          // 定义对象属性 target，取值为 apiTarget

          target: apiTarget,
          // 定义对象属性 changeOrigin，取值为 true

          changeOrigin: true,
        // 定义结束

        },
      // 定义结束

      },
    // 定义结束

    },
    // 定义对象属性 build，取值为 {

    build: {
      // 定义对象属性 outDir，取值为 'dist'

      outDir: 'dist',
      // 定义对象属性 sourcemap，取值为 false

      sourcemap: false,
    // 定义结束

    },
  // 当前代码块结束

  }
// 执行 })，实现对应的前端逻辑

})
