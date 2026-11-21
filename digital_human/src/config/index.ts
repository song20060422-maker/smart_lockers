/** @ticket 07-数字人Coze：运行时配置与角色切换（配置化，符合工单角色切换要求） */
// 导出 const CONFIG = {，供其他模块引用

export const CONFIG = {
  // 定义对象属性 cozeApiBase，取值为 import.meta.env.VITE_COZE_API_BASE ?? '/coze/v1'

  cozeApiBase: import.meta.env.VITE_COZE_API_BASE ?? '/coze/v1',
  // 定义对象属性 botId，取值为 import.meta.env.VITE_COZE_BOT_ID ?? ''

  botId: import.meta.env.VITE_COZE_BOT_ID ?? '',
  // 定义对象属性 token，取值为 import.meta.env.VITE_COZE_TOKEN ?? ''

  token: import.meta.env.VITE_COZE_TOKEN ?? '',
  // 定义对象属性 wsUrl，取值为 import.meta.env.VITE_WS_URL ?? '/ws/digital-human'

  wsUrl: import.meta.env.VITE_WS_URL ?? '/ws/digital-human',
  // 定义对象属性 ttsProvider，取值为 (import.meta.env.VITE_TTS_PROVIDER ?? 'web-speech') as 'web-speech' | 'azure' | 'openai'

  ttsProvider: (import.meta.env.VITE_TTS_PROVIDER ?? 'web-speech') as 'web-speech' | 'azure' | 'openai',
  // 定义对象属性 modelPath，取值为 import.meta.env.VITE_LIVE2D_MODEL_PATH ?? '/models/haru/haru.model3.json'

  modelPath: import.meta.env.VITE_LIVE2D_MODEL_PATH ?? '/models/haru/haru.model3.json',
// 当前代码块结束

}

// 导出 interface RoleConfig { id: string; name: string; modelPath: string; voice: string }，供其他模块引用

export interface RoleConfig { id: string; name: string; modelPath: string; voice: string }

/** 可切换的数字人角色（品牌代言人 / 虚拟讲解员） */
// 导出 const ROLES: Record<string, RoleConfig> = {，供其他模块引用

export const ROLES: Record<string, RoleConfig> = {
  // 定义对象属性 presenter，取值为 { id: 'presenter', name: '品牌代言人', modelPath: '/models/haru/haru.model3.json', voice: 'zh-CN-Xiaoxiao' }

  presenter: { id: 'presenter', name: '品牌代言人', modelPath: '/models/haru/haru.model3.json', voice: 'zh-CN-Xiaoxiao' },
  // 定义对象属性 guide，取值为 { id: 'guide', name: '虚拟讲解员', modelPath: '/models/shizuku/shizuku.model3.json', voice: 'zh-CN-Yunxi' }

  guide: { id: 'guide', name: '虚拟讲解员', modelPath: '/models/shizuku/shizuku.model3.json', voice: 'zh-CN-Yunxi' },
// 当前代码块结束

}

// 声明变量 currentRole: RoleConfig = ROLES.presenter，用于保存数据或函数引用

let currentRole: RoleConfig = ROLES.presenter
// 导出 const setRole = (id: string) => { if (ROLES[id]) currentRole = ROLES[id]; return currentRole }，供其他模块引用

export const setRole = (id: string) => { if (ROLES[id]) currentRole = ROLES[id]; return currentRole }
// 导出 const getRole = () => currentRole，供其他模块引用

export const getRole = () => currentRole
