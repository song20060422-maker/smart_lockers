/** @ticket 07-数字人Coze：统一消息结构（与 Coze 插件契约一致） */
// 导出 type MessageRole = 'user' | 'assistant'，供其他模块引用

export type MessageRole = 'user' | 'assistant'
// 导出 type MessageType = 'text' | 'voice' | 'command'，供其他模块引用

export type MessageType = 'text' | 'voice' | 'command'

// 导出 interface ChatMessage {，供其他模块引用

export interface ChatMessage {
  // 定义对象属性 role，取值为 MessageRole

  role: MessageRole
  // 定义对象属性 content，取值为 string

  content: string
  // 定义对象属性 type，取值为 MessageType

  type: MessageType
  // 定义对象属性 timestamp，取值为 number

  timestamp: number
  /** 情绪/场景/任务状态，由 Coze 工作流输出附带 */
  // 执行 meta?: { emotion?: string; scene?: string; task?: string }，实现对应的前端逻辑

  meta?: { emotion?: string; scene?: string; task?: string }
// 当前代码块结束

}

// 导出 interface Live2DAction {，供其他模块引用

export interface Live2DAction {
  // 定义对象属性 action，取值为 'nod' | 'smile' | 'shake' | 'blink' | 'transform' | 'idle'

  action: 'nod' | 'smile' | 'shake' | 'blink' | 'transform' | 'idle'
  // 执行 emotion?: string，实现对应的前端逻辑

  emotion?: string
  // 执行 duration_ms?: number，实现对应的前端逻辑

  duration_ms?: number
// 当前代码块结束

}

// 导出 interface AgentReplyPayload {，供其他模块引用

export interface AgentReplyPayload {
  // 定义对象属性 text，取值为 string

  text: string
  // 定义对象属性 emotion，取值为 string

  emotion: string
  // 定义对象属性 scene，取值为 string

  scene: string
  // 执行 action?: Live2DAction，实现对应的前端逻辑

  action?: Live2DAction
// 当前代码块结束

}
