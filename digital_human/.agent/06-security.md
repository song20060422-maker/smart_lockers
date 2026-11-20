---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_94ed5b2de4aa11f18663525400cd780f
    ReservedCode1: FDYYNESGDXH6S78WJpPPOhSbNBt4oxl7hqtKQTMPqCNVxDTkzNOAt9UXvwKLGZybw951wGBCAb6RzVKc7G6uOs1CKE89eGpjtiocbdLZMn132MV4lLLE5KApjdmH8XMm9HdTI3mspZXsBSTLMgE3k4WXFJ+J8V3eRYlhOdGeoxZJ5Jr2Nd8ajz+KfMk=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_94ed5b2de4aa11f18663525400cd780f
    ReservedCode2: FDYYNESGDXH6S78WJpPPOhSbNBt4oxl7hqtKQTMPqCNVxDTkzNOAt9UXvwKLGZybw951wGBCAb6RzVKc7G6uOs1CKE89eGpjtiocbdLZMn132MV4lLLE5KApjdmH8XMm9HdTI3mspZXsBSTLMgE3k4WXFJ+J8V3eRYlhOdGeoxZJ5Jr2Nd8ajz+KfMk=
---

# 06 · 安全规范 — 数字人模块

> 适用：本模块涉及凭证、对话内容、连接的处理。前置：`00-project.md`、`04-api.md`。

## 1. Coze 凭证安全

- token 由 `VITE_COZE_TOKEN` 从 `.env.local`（由 `.env.example` 复制）注入，必须被 gitignore 排除，禁止提交仓库。
- 禁止在源码、console、日志中硬编码或打印 token/API key；Vite 环境变量仅注入构建期需要的最小值，敏感调用优先经后端/代理转发，避免前端直连暴露凭证。

## 2. 对话内容安全

- 用户输入按文本处理并转义，避免 XSS/注入；服务端/Coze 返回内容按文本渲染，不执行 HTML/脚本。
- 敏感对话内容不落日志；如需持久化按脱敏处理。

## 3. 连接安全

- WebSocket 生产环境使用 wss（`VITE_WS_URL`），校验服务端来源；心跳与鉴权按 Coze/服务端要求处理。
- 断线重连设最大次数与退避上限，避免无限重连导致资源耗尽；卸载时彻底清理连接与定时器。

## 4. 配置安全

- `cozeApiBase`、`wsUrl`、`botId`、`modelPath`、`ttsProvider` 走 `src/config/index.ts` 与环境变量，禁止写死；`.env.local`、`coze/` 中如含密钥须谨慎，禁止提交真实凭据。

## 5. 红线

- 不臆造接口、字段、指标；消息结构与接口以真实代码（`types/index.ts`、`config/index.ts`）与根工单原文为准。

