/** @ticket 07-数字人Coze：WebSocket 通信模块——实时对话 + 状态推送 + 断线重连（目标重连成功率 100%） */
import type { ChatMessage, AgentReplyPayload } from '@/types'

// 定义类型别名，描述数据的形状：type Handlers = {

type Handlers = {
  // 定义对象属性 onReply，取值为 (p: AgentReplyPayload) => void

  onReply: (p: AgentReplyPayload) => void
  // 执行 onStatus?: (connected: boolean) => void，实现对应的前端逻辑

  onStatus?: (connected: boolean) => void
// 当前代码块结束

}

// 导出 class WebSocketClient {，供其他模块引用

export class WebSocketClient {
  // 赋值更新：private ws: WebSocket | null = null

  private ws: WebSocket | null = null
  // 执行 private url: string，实现对应的前端逻辑

  private url: string
  // 执行 private token: string，实现对应的前端逻辑

  private token: string
  // 执行 private handlers: Handlers，实现对应的前端逻辑

  private handlers: Handlers
  // 赋值更新：private retry = 0

  private retry = 0
  // 赋值更新：private maxRetry = 10

  private maxRetry = 10
  // 赋值更新：private timer: ReturnType<typeof setTimeout> | null = null

  private timer: ReturnType<typeof setTimeout> | null = null
  // 赋值更新：private manualClose = false

  private manualClose = false

  // 调用 constructor(url: string, token: string, handlers: Handlers) {，触发对应功能

  constructor(url: string, token: string, handlers: Handlers) {
    // 赋值更新：this.url = url

    this.url = url
    // 赋值更新：this.token = token

    this.token = token
    // 赋值更新：this.handlers = handlers

    this.handlers = handlers
  // 当前代码块结束

  }

  // 调用 connect() {，触发对应功能

  connect() {
    // 赋值更新：this.manualClose = false

    this.manualClose = false
    // 声明变量 sep = this.url.includes('?') ? '&' : '?'，用于保存数据或函数引用

    const sep = this.url.includes('?') ? '&' : '?'
    // 赋值更新：this.ws = new WebSocket(`${this.url}${sep}token=${encodeURIComponent(this.token)}`)

    this.ws = new WebSocket(`${this.url}${sep}token=${encodeURIComponent(this.token)}`)
    // 赋值更新：this.ws.onopen = () => { this.retry = 0; this.handlers.onStatus?.(true) }

    this.ws.onopen = () => { this.retry = 0; this.handlers.onStatus?.(true) }
    // 赋值更新：this.ws.onmessage = (evt) => this.handleMessage(evt.data)

    this.ws.onmessage = (evt) => this.handleMessage(evt.data)
    // 赋值更新：this.ws.onclose = () => { this.handlers.onStatus?.(false); if (!this.manualClose) this.scheduleReconnect() }

    this.ws.onclose = () => { this.handlers.onStatus?.(false); if (!this.manualClose) this.scheduleReconnect() }
    // 赋值更新：this.ws.onerror = () => { this.ws?.close() }

    this.ws.onerror = () => { this.ws?.close() }
  // 当前代码块结束

  }

  // 执行 private handleMessage(raw: string) {，实现对应的前端逻辑

  private handleMessage(raw: string) {
    // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

    try {
      // 声明变量 msg = JSON.parse(raw) as ChatMessage & { payload?: AgentReplyPayload }，用于保存数据或函数引用

      const msg = JSON.parse(raw) as ChatMessage & { payload?: AgentReplyPayload }
      // 智能体回复：附带情绪/场景/动作，触发表情与动画联动
      // 条件判断：if (msg.role === 'assistant') 时执行下面逻辑

      if (msg.role === 'assistant') {
        // 执行 this.handlers.onReply(msg.payload ?? {，实现对应的前端逻辑

        this.handlers.onReply(msg.payload ?? {
          // 定义对象属性 text，取值为 msg.content

          text: msg.content,
          // 定义对象属性 emotion，取值为 msg.meta?.emotion ?? 'neutral'

          emotion: msg.meta?.emotion ?? 'neutral',
          // 定义对象属性 scene，取值为 msg.meta?.scene ?? 'chat'

          scene: msg.meta?.scene ?? 'chat',
          // 定义对象属性 action，取值为 msg.type === 'command' ? undefined : undefined

          action: msg.type === 'command' ? undefined : undefined,
        // 执行 })，实现对应的前端逻辑

        })
      // 当前代码块结束

      }
    // 捕获上述逻辑抛出的异常并进行错误处理

    } catch { /* 忽略非法帧，工单要求无控制台错误 */ }
  // 当前代码块结束

  }

  // 执行 private scheduleReconnect() {，实现对应的前端逻辑

  private scheduleReconnect() {
    // 条件判断：if (this.retry >= this.maxRetry) return 时执行下面逻辑

    if (this.retry >= this.maxRetry) return
    // 声明变量 delay = Math.min(1000 * Math.pow(2, this.retry), 30000) // 指数退避，用于保存数据或函数引用

    const delay = Math.min(1000 * Math.pow(2, this.retry), 30000) // 指数退避
    // 执行 this.retry += 1，实现对应的前端逻辑

    this.retry += 1
    // 赋值更新：this.timer = setTimeout(() => this.connect(), delay)

    this.timer = setTimeout(() => this.connect(), delay)
  // 当前代码块结束

  }

  /** 发送用户消息（文本或语音识别后的文本） */
  // 赋值更新：send(content: string, type: 'text' | 'voice' = 'text') {

  send(content: string, type: 'text' | 'voice' = 'text') {
    // 条件判断：if (this.ws?.readyState !== WebSocket.OPEN) return false 时执行下面逻辑

    if (this.ws?.readyState !== WebSocket.OPEN) return false
    // 声明变量 msg: ChatMessage = { role: 'user', content, type, timestamp: Date.now() }，用于保存数据或函数引用

    const msg: ChatMessage = { role: 'user', content, type, timestamp: Date.now() }
    // 执行 this.ws.send(JSON.stringify(msg))，实现对应的前端逻辑

    this.ws.send(JSON.stringify(msg))
    // 返回 true 给调用方

    return true
  // 当前代码块结束

  }

  // 赋值更新：close() { this.manualClose = true; if (this.timer) clearTimeout(this.timer); this.ws?.close() }

  close() { this.manualClose = true; if (this.timer) clearTimeout(this.timer); this.ws?.close() }
// 当前代码块结束

}
