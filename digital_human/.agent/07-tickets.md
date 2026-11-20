---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_95f7a785e4aa11f18663525400cd780f
    ReservedCode1: OgA+54FpgBFjftk5pJiLJu3S/rMhKljuLDWQoX6S3lK9g2EDt4X7V1fJeyg1fGcAu1J/fLXRswyROhcKIE3Ull0elryuAl/sIAmy5mK2lne/dYuKx9u4I+LTmT8dzkjvY9pEM2Ucobsg8wWg2qOwkWz9yn0LkWodGB2mRYkpP6Om2fpU5dZw6tFpN04=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_95f7a785e4aa11f18663525400cd780f
    ReservedCode2: OgA+54FpgBFjftk5pJiLJu3S/rMhKljuLDWQoX6S3lK9g2EDt4X7V1fJeyg1fGcAu1J/fLXRswyROhcKIE3Ull0elryuAl/sIAmy5mK2lne/dYuKx9u4I+LTmT8dzkjvY9pEM2Ucobsg8wWg2qOwkWz9yn0LkWodGB2mRYkpP6Om2fpU5dZw6tFpN04=
---

# 07 · 工单映射速查 — 数字人模块

> 面向工单 07 开发时查本文件；业务细节以根工单原文与真实代码为准。本表只列数字人模块落点。

## 工单 07 · 数字人（对应根 AGENT.MD 数字人任务）

| 能力 | 落点 | 说明 |
|---|---|---|
| 形象渲染 | `src/live2d/Live2DViewer.tsx` | Live2D 模型加载、动作/表情联动；资源在 `public/models/`，不入库 |
| 状态管理 | `src/store/index.ts`（Redux Toolkit） | 对话消息/连接状态/角色/UI 状态统一管理 |
| 智能体对话 | `coze/agent.config.json`、`plugins.json`、`workflow.json`、`knowledge_base.json`、`knowledge/faq.json` | Coze 智能体/插件/工作流/知识库配置 |
| 运行时配置 | `src/config/index.ts` | `CONFIG` / `ROLES` / `setRole` / `getRole`，端点/botId/模型/角色可配 |
| 实时通信 | `src/ws/WebSocketClient.ts` | 心跳、断线重连（指数退避）、卸载清理；`/ws/digital-human` |
| 语音交互 | `src/tts/speech.ts` | TTS（web-speech/azure/openai，可配）；不可用时文本降级 |
| 核心编排 | `src/hooks/useDigitalHuman.ts` | 连接/消息/输入/发送/形象状态统一管理 |
| 消息契约 | `src/types/index.ts` | `ChatMessage` / `AgentReplyPayload` / `Live2DAction` |
| 环境变量 | `.env.example` → `.env.local` | `VITE_COZE_API_BASE` / `VITE_COZE_BOT_ID` / `VITE_COZE_TOKEN` / `VITE_WS_URL` / `VITE_TTS_PROVIDER` / `VITE_LIVE2D_MODEL_PATH` |
| 联调/性能 | `scripts/integration-test.mjs`、`scripts/perf-test.mjs`、`PERFORMANCE_REPORT.md` | `npm run test:integration` / `npm run test:perf` |

## 消息 schema

- `ChatMessage { role: 'user'|'assistant', content, type: 'text'|'voice'|'command', timestamp, meta?:{ emotion?, scene?, task? } }`
- `AgentReplyPayload { text, emotion, scene, action?: Live2DAction }`；`Live2DAction.action ∈ {nod, smile, shake, blink, transform, idle}`。

## 交互链路

用户语音/文本 → WebSocket 上行 → Coze 智能体（工作流 + 知识库）→ `AgentReplyPayload`（含 emotion/scene/action）→ TTS 播报 + Live2D 表情动作联动。

## 后端支撑（可选）

- 业务知识库/配置可由 `../smart_lockers/`（Gin :18080）提供 RESTful / WebSocket 支撑，消息结构保持一致，指标响应 <2s。
