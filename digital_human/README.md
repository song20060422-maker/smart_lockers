---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_8f77c617e27c11f19919525400638852
    ReservedCode1: 3kqC7m8WKqcha85XaRAVEXJxpXqbVMNEdAWYzFL/VHymmNuCbZeq4zQ/FA8E96zOAOxUQNfcBhTejaZCya79Gv2QzRyDb56qiB1gndjjjhOV9xVCP0c6bF7QkLtOSrgiKKjiZhSrgGbclnmB2WaBHd/kqGOPhWKirJZT4Yc5AuSQNnJhE0bU+hTUDwU=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_8f77c617e27c11f19919525400638852
    ReservedCode2: 3kqC7m8WKqcha85XaRAVEXJxpXqbVMNEdAWYzFL/VHymmNuCbZeq4zQ/FA8E96zOAOxUQNfcBhTejaZCya79Gv2QzRyDb56qiB1gndjjjhOV9xVCP0c6bF7QkLtOSrgiKKjiZhSrgGbclnmB2WaBHd/kqGOPhWKirJZT4Yc5AuSQNnJhE0bU+hTUDwU=
---

<!-- @ticket 07-数字人Coze：项目运行文档（环境搭建 / 配置步骤 / 测试用例） -->
# 智能快寄柜数字人（工单 07）

基于 Coze 智能体 + `awesome-digital-human-live2d` 的 Live2D 数字人独立子模块。**依赖与主 Vue3 前端（smart_lockers_frontend）完全隔离**，独立 `package.json`，独立端口 5174。

## 一、环境搭建
- Node.js 18+；`npm install`
- Live2D Core SDK 为商业授权组件，模型资源见 `public/models/README.md`（不随仓库提交）

## 二、配置步骤
1. `cp .env.example .env.local`，填写 Coze 的 `VITE_COZE_BOT_ID` / `VITE_COZE_TOKEN`、WebSocket 地址与 TTS 提供方。
2. 在 Coze 平台按 `coze/agent.config.json`、`plugins.json`、`workflow.json`、`knowledge_base.json` 创建智能体：启用插件、主工作流 + 分支节点、导入 `coze/knowledge/faq.json`。
3. 放置 Live2D 模型到 `public/models/<角色>/`，配置 `VITE_LIVE2D_MODEL_PATH`；角色切换见 `src/config/index.ts` 的 `ROLES`。
4. `npm run dev` 启动（端口 5174，与前端 5173 隔离）；WebSocket 代理到后端 `:18080`。

## 三、交互链路
用户语音/文本 → WebSocket 上行 → Coze 智能体（工作流 + 知识库）→ 回复（含 emotion/scene/action）→ TTS 播报 + Live2D 表情动作联动。

## 四、测试用例
| 编号 | 场景 | 预期 |
|---|---|---|
| T1 | 输入“如何寄存物品” | FAQ 命中，回复正确，触发点头 |
| T2 | 语音输入 | SpeechRecognition 识别后转文本上行 |
| T3 | 知识库问答 | 匹配度 ≥ 阈值，回答准确率 ≥ 90% |
| T4 | 多轮对话 | 上下文记忆生效，指代理解正确 |
| T5 | 断网重连 | 指数退避重连，成功率 100% |
| T6 | 角色切换 | 模型路径与音色随之切换 |

## 五、验证命令
`npm run typecheck`（类型检查）、`npm run build`（构建）、`npm run test:integration`（联调桩）、`npm run test:perf`（性能桩）。
