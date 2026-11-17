/// <reference types="vite/client" /> // @ticket 00-project：Vite 客户端类型声明

// @ticket 00-project：环境变量类型定义
interface ImportMetaEnv {
  readonly VITE_API_BASE: string
  readonly VITE_API_TARGET: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, unknown>
  export default component
}
