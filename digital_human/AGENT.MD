---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_8dc90666e4aa11f1b820525400de85a5
    ReservedCode1: ILDgMYRfWrbsGtejzUz8aCjNFSiUGN6m7VS+6oU4gCZPS4Ycu1TwKeiI+2y5wNE0cYDcOQi40fmNZSbU2IvhXRqpM04nHJFcL4PTa1JBTpFtlUOzBf+gF99i3Jf4lXFXfcO74k3ziG0il8UsHeR7F3F+ABxjJ9CiR2TLCrcSIiWA97olC7NvHaB/DZo=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_8dc90666e4aa11f1b820525400de85a5
    ReservedCode2: ILDgMYRfWrbsGtejzUz8aCjNFSiUGN6m7VS+6oU4gCZPS4Ycu1TwKeiI+2y5wNE0cYDcOQi40fmNZSbU2IvhXRqpM04nHJFcL4PTa1JBTpFtlUOzBf+gF99i3Jf4lXFXfcO74k3ziG0il8UsHeR7F3F+ABxjJ9CiR2TLCrcSIiWA97olC7NvHaB/DZo=
---

# AGENT.MD — digital_human 数字人 AI 协作入口

> 本文件是 `digital_human/`（智能快寄柜 · 数字人交互模块）内的 AI 协作细则。
> 仓库规则优先级：根 `../AGENT.MD` > 根 `../.agent/*.md` > **本文件** > 本目录 `.agent/*.md` > 个人习惯。冲突时以更高优先级为准；本文件只细化数字人模块，不得与根铁律冲突。
> 任何 AI 在本目录产出代码/配置/文档前，必须先读根 `../AGENT.MD` 与本文件，再按本目录 `.agent/` 分规则执行。

---

## 1. 子项目身份

| 项 | 内容 |
|---|---|
| 名称 | 智能快寄柜 · 数字人（Digital Human）交互模块 |
| 包名 | `smart-lockers-digital-human`（`package.json`，`"type": "module"`） |
| 职责 | Live2D 形象 + 语音交互 + Coze 智能体对话，提供沉浸式智能客服 |
| 端口 | 5174（Vite；与主 Vue3 前端 5173 隔离） |
| 覆盖工单 | 07 · 数字人（Coze 智能体 + Live2D） |
| 边界 | **独立 React 子模块**，依赖与主前端、B 端前端、后端完全隔离；独立 `package.json` |

## 2. 技术栈（真实依赖，禁止臆造）

| 层 | 选型 | 版本/依据 |
|---|---|---|
| 框架 | React | react ^18.3.1 / react-dom ^18.3.1 / @vitejs/plugin-react ^4.3.1 |
| 语言 | TypeScript | ^5.5.4（`tsc --noEmit`） |
| 构建 | Vite | ^5.4.2 |
| 状态 | Redux Toolkit + react-redux | ^2.2.7 / ^9.1.2（`src/store/`） |
| 智能体 | Coze | 配置见 `coze/` 与 `.env.example`（bot/token/ws） |
| 实时通信 | WebSocket（`src/ws/WebSocketClient.ts`） | `/ws/digital-human` |
| 语音 | TTS（web-speech / azure / openai，可配） | `src/tts/speech.ts` |
| 形象 | Live2D（`src/live2d/Live2DViewer.tsx`） | 商业模型资源不随仓库提交 |

> 注意：本模块是 React + Redux（非 Vue/Pinia），端口 5174（非 5173），禁止照搬其他前端的依赖假设。

## 3. 目录与代码地图（现状，不得假设为空）

```
digital_human/
├── AGENT.MD / README.md / PERFORMANCE_REPORT.md / .agent/   <- 规则与既有文档（README.md 为原运行文档，保留）
├── index.html / package.json / tsconfig.json / vite.config.ts
├── .env.example                        <- 环境变量模板（复制为 .env.local，勿提交）
├── coze/                               <- Coze 智能体配置：agent.config.json / plugins.json /
│   └── knowledge/                         workflow.json / knowledge_base.json / faq.json / brand.md
├── public/models/README.md             <- Live2D 模型资源目录（模型不入库）
├── scripts/                            <- integration-test.mjs / perf-test.mjs
└── src/
    ├── main.tsx / App.tsx / styles.css / vite-env.d.ts
    ├── config/index.ts                 <- 运行时配置与角色切换（CONFIG / ROLES / setRole）
    ├── store/index.ts                  <- Redux store
    ├── hooks/useDigitalHuman.ts        <- 核心状态与流程编排
    ├── live2d/Live2DViewer.tsx         <- Live2D 形象渲染
    ├── tts/speech.ts                   <- 语音合成
    ├── ws/WebSocketClient.ts           <- WebSocket 客户端
    └── types/index.ts                  <- 统一消息契约（ChatMessage / Live2DAction / AgentReplyPayload）
```

## 4. 本端开发铁律（继承根 §4 8 条，另加数字人特有）

1. **工单编号注释**：关键逻辑标注 `// @ticket 07-数字人Coze：<说明>`（与既有代码一致）。
2. **类型检查**：提交前 `npm run typecheck`（tsc --noEmit）通过，避免 any 逃逸。
3. **消息 schema 对齐**：对话消息统一 `ChatMessage { role, content, type, timestamp, meta? }`（见 `src/types/index.ts`），前端、WebSocket、Coze 三方一致，禁止各自扩展字段。
4. **连接健壮性**：WebSocket 断线重连 + 指数退避 + 心跳，页面卸载时清理连接/定时器，避免泄漏。
5. **凭证安全**：Coze token 只经 `.env.local`（由 `.env.example` 复制，gitignore 排除）注入，禁止入库/硬编码；WebSocket 代理到后端 :18080。
6. **语音降级**：TTS 提供方由 `VITE_TTS_PROVIDER` 配置（web-speech/azure/openai），不可用时文本兜底，不阻塞核心对话。
7. **可配置化**：端点、botId、模型路径、角色（`src/config/index.ts` 的 `ROLES`）全走配置，禁止写死。
8. **引用真实来源**：接口、字段、消息结构以本模块真实代码与根工单原文为准，禁止臆造。

## 5. 分规则索引

| 文件 | 适用范围 |
|---|---|
| `.agent/README.md` | 分规则加载方式与维护约定 |
| `.agent/00-project.md` | 数字人技术栈边界、目录、命名、质量铁律 |
| `.agent/11-frontend.md` | React 组件、Redux、Live2D、语音、useDigitalHuman 编排 |
| `.agent/13-api.md` | Coze 对话、WebSocket 消息 schema、后端支撑约定 |
| `.agent/14-testing.md` | 类型检查/构建、联调/性能脚本、验收口径 |
| `.agent/15-security.md` | Coze 凭证、输入内容、连接安全 |
| `.agent/16-tickets.md` | 工单 07 → 组件/API/Hook/config 速查 |

## 6. 执行约定

- 保证本目录 `npm install && npm run dev` 可启动（:5174）；`.env.local` 从 `.env.example` 复制并填写。
- 不改动 `.git`、`.gitignore`、`README.md`、`PERFORMANCE_REPORT.md` 及根目录已有文件；不提交 token/密钥。
- 一次只推进一个能力（形象/语音/对话/连接），做完即交付并标注 `07-数字人Coze`；变更可追溯。

