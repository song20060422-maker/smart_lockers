---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_907fbf4ae4aa11f1b820525400de85a5
    ReservedCode1: tHY/sWIIDsNqQFIzPL3DHZCcpHENAHJFq8M7Y3eCfoTFYSCgknqGOrQlvYZQHw96LYw7bnPmc2W/EiurrfBgmLfV90EuWUMS+KNz+oQaMOkw+T7khtoDGUNUJy/7xnGp/QdSpz8eDoZH7u0fo2f03akEcj3RyzLiZ8kiIa/Wg3Y48GBb96rRZXGlVSc=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_907fbf4ae4aa11f1b820525400de85a5
    ReservedCode2: tHY/sWIIDsNqQFIzPL3DHZCcpHENAHJFq8M7Y3eCfoTFYSCgknqGOrQlvYZQHw96LYw7bnPmc2W/EiurrfBgmLfV90EuWUMS+KNz+oQaMOkw+T7khtoDGUNUJy/7xnGp/QdSpz8eDoZH7u0fo2f03akEcj3RyzLiZ8kiIa/Wg3Y48GBb96rRZXGlVSc=
---

# 00 · 数字人模块全局规则

> 适用：本目录所有任务。加载顺序：根 `../../AGENT.MD` → `../AGENT.MD` → 本文件 → 其他分规则。

## 1. 技术栈边界（不得混用）

- 本模块 = 数字人：React 18 + TypeScript 5.5 + Vite 5 + Redux Toolkit + Coze + WebSocket + TTS + Live2D。
- 独立 React 子模块，端口 5174；不要照搬 `../smart_lockers_frontend/` / `../smart_lockers_admin_frontend/`（Vue3 + Pinia）的依赖与结构，也不照搬后端的 GORM/Redis 假设。
- 后端支撑在 `../smart_lockers/`（Gin :18080）；Coze 凭证仅本地 `.env.local`。

## 2. 目录与命名

- `src/{config,store,hooks,live2d,tts,ws,types}` + `coze/` + `scripts/` + `public/models/`。
- 组件大驼峰、hook 以 `use` 开头、工具按职责分目录；新增能力同步更新 `src/types/index.ts`、`src/config/index.ts` 与本目录 `07-tickets.md`。
- 消息结构统一 `ChatMessage { role, content, type, timestamp, meta? }`，三方（前端/WS/Coze）一致，禁止各自扩展。

## 3. 工单编号注释（强制）

- 关键逻辑标注 `// @ticket 07-数字人Coze：<说明>`（与既有代码一致），如 `// @ticket 07-数字人Coze：WebSocket 断线重连`。

## 4. 通用质量铁律

1. 提交前 `npm run typecheck` / `npm run build` 通过，避免 any 逃逸。
2. 对话消息三方 schema 对齐（见 `src/types/index.ts`）。
3. WebSocket 断线重连 + 心跳 + 卸载清理，避免连接/定时器泄漏。
4. Coze 凭证仅 `.env.local`（gitignore 排除），禁止入库/硬编码。
5. TTS 提供方（`VITE_TTS_PROVIDER`，web-speech/azure/openai）不可用时有文本兜底。
6. 配置驱动：`src/config/index.ts` 的 `CONFIG` 与 `ROLES` 管理端点、botId、模型路径、角色，禁止写死。

## 5. 标准产出物清单

交互流程说明、消息 schema 设计、Coze 智能体配置说明、联调记录、连接健壮性与降级验证、性能报告（`PERFORMANCE_REPORT.md`）、上线说明与回滚预案。

## 6. 可启动性与配置安全

- 保证 `npm install && npm run dev` 可启动（:5174）；`.env.local` 由 `.env.example` 复制填写（`VITE_COZE_API_BASE`、`VITE_COZE_BOT_ID`、`VITE_COZE_TOKEN`、`VITE_WS_URL`、`VITE_TTS_PROVIDER`、`VITE_LIVE2D_MODEL_PATH`）。
- `.env.local` 含 token 必须被 gitignore 排除，不得提交；不修改 `.git`、`.gitignore`、`README.md`、`PERFORMANCE_REPORT.md` 与根目录已有文件。

## 7. 红线

- 不臆造接口、字段、指标；结论来自真实代码或根工单原文。
- 一次只推进一个能力，做完即交付并标注 `07-数字人Coze`。
