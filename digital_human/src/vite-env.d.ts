/// <reference types="vite/client" />
/** @ticket 07-数字人Coze：环境变量类型声明 */
interface ImportMetaEnv {
  readonly VITE_COZE_API_BASE: string
  readonly VITE_COZE_BOT_ID: string
  readonly VITE_COZE_TOKEN: string
  readonly VITE_WS_URL: string
  readonly VITE_TTS_PROVIDER: string
  readonly VITE_LIVE2D_MODEL_PATH: string
}
interface ImportMeta { readonly env: ImportMetaEnv }
