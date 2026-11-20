---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_93dbd251e4aa11f18663525400cd780f
    ReservedCode1: NH+3ClA5DlGF+nWWSjk2KlM6RzYlXusNRC3x+CDJ+gYCUosbP57xpFRhOCClr8sGYmEcg61bdpXNN4z78HkxANd6uvjo2Chm/EQ102X6OvLQ2sK/7l/JR6I94T8+RswnPteQHOOOZIgjOfPmHouoj6OzMnIFiJeHKt3TuFKK45qMamX5JPpIbmJK7Dg=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_93dbd251e4aa11f18663525400cd780f
    ReservedCode2: NH+3ClA5DlGF+nWWSjk2KlM6RzYlXusNRC3x+CDJ+gYCUosbP57xpFRhOCClr8sGYmEcg61bdpXNN4z78HkxANd6uvjo2Chm/EQ102X6OvLQ2sK/7l/JR6I94T8+RswnPteQHOOOZIgjOfPmHouoj6OzMnIFiJeHKt3TuFKK45qMamX5JPpIbmJK7Dg=
---

# 05 · 测试规范 — 数字人模块

> 适用：本模块所有测试活动。前置：`00-project.md`、`02-frontend.md`、`04-api.md`。

## 1. 类型检查与构建（必过）

- 类型检查：`npm run typecheck`（tsc --noEmit）通过，避免 any 逃逸。
- 生产构建：`npm run build`（tsc --noEmit && vite build）通过，Live2D 资源正常打包。

## 2. 联调与性能脚本（仓库既有）

- 联调：`npm run test:integration`（`scripts/integration-test.mjs`）；覆盖 FAQ 命中、语音输入、知识库问答、多轮对话、断网重连、角色切换等主流程。
- 性能：`npm run test:perf`（`scripts/perf-test.mjs`）；结果与解读见 `PERFORMANCE_REPORT.md`，迭代时同步更新。

## 3. 联调测试要点

- Coze 对话：填好 `.env.local` 后，验证 输入 → WS/Coze → `AgentReplyPayload` → 消息更新 → TTS 播报 → Live2D 动作联动的完整链路。
- WebSocket：验证心跳、断线重连（指数退避）、重连后消息恢复、卸载后连接与定时器正确清理。
- 语音：验证 TTS 正常路径；提供方不可用/授权失败时验证文本降级可用。
- 角色切换：验证 `setRole` 后模型路径与音色随 `ROLES` 切换。

## 4. 边界与健壮性验证

- Live2D 模型资源加载失败给占位/降级提示，不崩溃。
- 空响应/Coze 超时/网络异常时消息列表状态与用户提示正确。
- 消息 schema 三方一致：前端展示、WS 传输、Coze 解析字段对齐。

## 5. 验收口径

- 工单 07 交付前：typecheck/build 通过、主流程联调通过、重连与降级验证通过、性能报告更新、凭证无泄露。

## 6. 注释

验证用例标注 `// @ticket 07-数字人Coze：<说明>`。
