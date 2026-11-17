/** @ticket 07-数字人Coze：统一消息结构（与 Coze 插件契约一致） */
export type MessageRole = 'user' | 'assistant'
export type MessageType = 'text' | 'voice' | 'command'

export interface ChatMessage {
  role: MessageRole
  content: string
  type: MessageType
  timestamp: number
  /** 情绪/场景/任务状态，由 Coze 工作流输出附带 */
  meta?: { emotion?: string; scene?: string; task?: string }
}

export interface Live2DAction {
  action: 'nod' | 'smile' | 'shake' | 'blink' | 'transform' | 'idle'
  emotion?: string
  duration_ms?: number
}

export interface AgentReplyPayload {
  text: string
  emotion: string
  scene: string
  action?: Live2DAction
}
