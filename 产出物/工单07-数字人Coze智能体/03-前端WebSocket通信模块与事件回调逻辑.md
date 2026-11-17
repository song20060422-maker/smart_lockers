# 工单07 数字人Coze智能体 — 前端 WebSocket 通信模块与事件回调逻辑
> **工单编号：07｜数字人Coze智能体**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 1. WebSocket 通信模块（src/ws/WebSocketClient.ts）
- **连接**：`new WebSocket(url?token=...)`；onopen 重置重试计数并回传状态。
- **断线重连**：指数退避 `delay=min(1000*2^retry, 30000)`，maxRetry=10；`manualClose` 区分主动关闭，避免误重连（目标重连成功率 100%）。
- **消息处理**：仅处理 `role=assistant` 帧，解析失败静默忽略（工单要求无控制台错误）；回复附带 emotion/scene/action 触发表情动画联动。
- **发送**：`send(content, type)` 校验 readyState，构造 `{role:'user',content,type,timestamp}`。
```ts
private handleMessage(raw: string) {
  const msg = JSON.parse(raw)
  if (msg.role === 'assistant') this.handlers.onReply(msg.payload ?? {
    text: msg.content, emotion: msg.meta?.emotion ?? 'neutral', scene: msg.meta?.scene ?? 'chat' })
}
private scheduleReconnect() {
  if (this.retry >= this.maxRetry) return
  const delay = Math.min(1000 * Math.pow(2, this.retry), 30000)
  this.retry += 1; this.timer = setTimeout(() => this.connect(), delay)
}
```

## 2. 编排 Hook / 事件回调（src/hooks/useDigitalHuman.ts）
串联 **WebSocket 回复 → TTS 播报 → Live2D 表情动作联动**：
- onStatus → `dispatch(setConnected)` 更新连接态。
- onReply → 落库消息 + `setThinking(false)` + `speech.speak(text)` 播报 + `setAction(action ?? emotionAction(emotion))` 动画联动。
- **情绪→动作映射**：`happy→smile, positive→nod, sad→shake, neutral→blink`，duration_ms=1200。
- send → 校验空值、发送失败提示"连接未就绪"、更新 thinking/error 状态；listen → 启动语音识别。
```ts
const emotionAction = (emotion) => ({
  action: ({ happy:'smile', positive:'nod', sad:'shake', neutral:'blink' }[emotion] ?? 'blink'),
  emotion, duration_ms: 1200 })
```

## 3. 语音模块（src/tts/speech.ts）
- 输入：SpeechRecognition（lang=zh-CN，interimResults=false）→ onRecognized。
- 输出：SpeechSynthesis TTS，voice 按角色 voiceName 匹配，onend 回调；可切换 Azure/OpenAI。
