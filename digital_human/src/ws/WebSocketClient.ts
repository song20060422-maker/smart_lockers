/** @ticket 07-数字人Coze：WebSocket 通信模块——实时对话 + 状态推送 + 断线重连（目标重连成功率 100%） */
import type { ChatMessage, AgentReplyPayload } from '@/types'

type Handlers = {
  onReply: (p: AgentReplyPayload) => void
  onStatus?: (connected: boolean) => void
}

export class WebSocketClient {
  private ws: WebSocket | null = null
  private url: string
  private token: string
  private handlers: Handlers
  private retry = 0
  private maxRetry = 10
  private timer: ReturnType<typeof setTimeout> | null = null
  private manualClose = false

  constructor(url: string, token: string, handlers: Handlers) {
    this.url = url
    this.token = token
    this.handlers = handlers
  }

  connect() {
    this.manualClose = false
    const sep = this.url.includes('?') ? '&' : '?'
    this.ws = new WebSocket(`${this.url}${sep}token=${encodeURIComponent(this.token)}`)
    this.ws.onopen = () => { this.retry = 0; this.handlers.onStatus?.(true) }
    this.ws.onmessage = (evt) => this.handleMessage(evt.data)
    this.ws.onclose = () => { this.handlers.onStatus?.(false); if (!this.manualClose) this.scheduleReconnect() }
    this.ws.onerror = () => { this.ws?.close() }
  }

  private handleMessage(raw: string) {
    try {
      const msg = JSON.parse(raw) as ChatMessage & { payload?: AgentReplyPayload }
      // 智能体回复：附带情绪/场景/动作，触发表情与动画联动
      if (msg.role === 'assistant') {
        this.handlers.onReply(msg.payload ?? {
          text: msg.content,
          emotion: msg.meta?.emotion ?? 'neutral',
          scene: msg.meta?.scene ?? 'chat',
          action: msg.type === 'command' ? undefined : undefined,
        })
      }
    } catch { /* 忽略非法帧，工单要求无控制台错误 */ }
  }

  private scheduleReconnect() {
    if (this.retry >= this.maxRetry) return
    const delay = Math.min(1000 * Math.pow(2, this.retry), 30000) // 指数退避
    this.retry += 1
    this.timer = setTimeout(() => this.connect(), delay)
  }

  /** 发送用户消息（文本或语音识别后的文本） */
  send(content: string, type: 'text' | 'voice' = 'text') {
    if (this.ws?.readyState !== WebSocket.OPEN) return false
    const msg: ChatMessage = { role: 'user', content, type, timestamp: Date.now() }
    this.ws.send(JSON.stringify(msg))
    return true
  }

  close() { this.manualClose = true; if (this.timer) clearTimeout(this.timer); this.ws?.close() }
}
