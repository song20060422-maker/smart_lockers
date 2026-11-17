/** @ticket 07-数字人Coze：运行时配置与角色切换（配置化，符合工单角色切换要求） */
export const CONFIG = {
  cozeApiBase: import.meta.env.VITE_COZE_API_BASE ?? '/coze/v1',
  botId: import.meta.env.VITE_COZE_BOT_ID ?? '',
  token: import.meta.env.VITE_COZE_TOKEN ?? '',
  wsUrl: import.meta.env.VITE_WS_URL ?? '/ws/digital-human',
  ttsProvider: (import.meta.env.VITE_TTS_PROVIDER ?? 'web-speech') as 'web-speech' | 'azure' | 'openai',
  modelPath: import.meta.env.VITE_LIVE2D_MODEL_PATH ?? '/models/haru/haru.model3.json',
}

export interface RoleConfig { id: string; name: string; modelPath: string; voice: string }

/** 可切换的数字人角色（品牌代言人 / 虚拟讲解员） */
export const ROLES: Record<string, RoleConfig> = {
  presenter: { id: 'presenter', name: '品牌代言人', modelPath: '/models/haru/haru.model3.json', voice: 'zh-CN-Xiaoxiao' },
  guide: { id: 'guide', name: '虚拟讲解员', modelPath: '/models/shizuku/shizuku.model3.json', voice: 'zh-CN-Yunxi' },
}

let currentRole: RoleConfig = ROLES.presenter
export const setRole = (id: string) => { if (ROLES[id]) currentRole = ROLES[id]; return currentRole }
export const getRole = () => currentRole
