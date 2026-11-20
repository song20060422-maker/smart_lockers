---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_92b9a222e4aa11f18663525400cd780f
    ReservedCode1: LtOaXUHidbKIdzVA5Lwvc7TqIJ1R5xKpIx1smWVd2kmSjqdA0ErMSU/Fx6iD+Puo30KvoHu2JX22tYSMxihw+3lkRyB2TYy/XkmCs4NCVQuAQlxdtH0u6fsHZOeWxfSTWkREjIuobhDszxoIlmc1Mizd+p2RCYaqO0Bak1AZnGVug1WYMuhEWlRvx/s=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_92b9a222e4aa11f18663525400cd780f
    ReservedCode2: LtOaXUHidbKIdzVA5Lwvc7TqIJ1R5xKpIx1smWVd2kmSjqdA0ErMSU/Fx6iD+Puo30KvoHu2JX22tYSMxihw+3lkRyB2TYy/XkmCs4NCVQuAQlxdtH0u6fsHZOeWxfSTWkREjIuobhDszxoIlmc1Mizd+p2RCYaqO0Bak1AZnGVug1WYMuhEWlRvx/s=
---

# 04 · 接口规范 — 数字人对话（Coze + WebSocket）

> 适用：本模块所有对话与连接开发。前置：`00-project.md`、`02-frontend.md`。

## 1. 消息 schema（统一，三方对齐，以 src/types/index.ts 为准）

```ts
ChatMessage { role: 'user'|'assistant'; content: string; type: 'text'|'voice'|'command'; timestamp: number;
              meta?: { emotion?: string; scene?: string; task?: string } }
AgentReplyPayload { text: string; emotion: string; scene: string; action?: Live2DAction }
Live2DAction { action: 'nod'|'smile'|'shake'|'blink'|'transform'|'idle'; emotion?: string; duration_ms?: number }
```

- 前端、WebSocket、Coze 三方使用同一契约，禁止各自扩展字段；新增字段须同步 `types/index.ts` 与 Coze 插件契约。

## 2. Coze 对话与配置

- 配置集中 `src/config/index.ts` 的 `CONFIG`：`cozeApiBase`（默认 `/coze/v1`）、`botId`、`token`、`wsUrl`（默认 `/ws/digital-human`）、`ttsProvider`、`modelPath`，均来自 `.env.local` 注入，禁止硬编码。
- Coze 智能体定义见 `coze/agent.config.json`、`plugins.json`、`workflow.json`、`knowledge_base.json` 及 `knowledge/faq.json`、`brand.md`；调整智能体须同步更新这些配置。
- 鉴权 token 由环境变量 `VITE_COZE_TOKEN` 注入，禁止入库；响应解析为 `AgentReplyPayload`，错误按统一提示，失败可重试（带退避）。

## 3. WebSocket（src/ws/WebSocketClient.ts）

- 连接地址 `VITE_WS_URL`（默认 `/ws/digital-human`），经代理转发到后端 `:18080`。
- 建立后发心跳，断线按指数退避重连（设最大重试次数），重连成功后恢复未确认消息。
- 页面卸载/组件销毁时 close 连接并清理心跳定时器，避免泄漏。
- 消息收发遵循 `ChatMessage` 契约；接收后交 `useDigitalHuman`/Redux 更新状态，禁止组件内各自处理 socket。

## 4. 后端支撑约定（可选，非本模块实现）

- 业务知识库/配置可由 `../smart_lockers/`（Gin :18080）提供 RESTful 接口或 WebSocket 对话转发；支撑接口指标响应 <2s，消息结构保持一致。

## 5. 注释

接口/消息处理处标注 `// @ticket 07-数字人Coze：<说明>`。
