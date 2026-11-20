---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_919ce731e4aa11f1b820525400de85a5
    ReservedCode1: wN7i+OMAEKgtSjV7oWz4NNUiFkSXSLSUpjCtnYTqKh7yuXugw8ORaYIXdBNK1f0RwsY1jSbtwioktoxH+wd66Sd+mUQy5ofPHLHWyg3WjKGO4IAwFcq96l+B5gBeaIOemLdITEQaP1XdoFspGQsfTCOO0/wOTwrmwo7xzGnI3oDdUVjAZwxnUvCh0GA=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_919ce731e4aa11f1b820525400de85a5
    ReservedCode2: wN7i+OMAEKgtSjV7oWz4NNUiFkSXSLSUpjCtnYTqKh7yuXugw8ORaYIXdBNK1f0RwsY1jSbtwioktoxH+wd66Sd+mUQy5ofPHLHWyg3WjKGO4IAwFcq96l+B5gBeaIOemLdITEQaP1XdoFspGQsfTCOO0/wOTwrmwo7xzGnI3oDdUVjAZwxnUvCh0GA=
---

# 02 · 前端规范 — 数字人（React + Redux + Live2D + 语音）

> 适用：本模块所有组件/交互开发。前置：`00-project.md`、`04-api.md`。

## 1. React 组件

- 函数组件 + Hooks；`App.tsx` 编排整体，`live2d/Live2DViewer.tsx` 负责形象渲染。
- props/state 显式类型，避免 any；组件职责单一，复杂逻辑下沉到 hook（`hooks/useDigitalHuman.ts`）或 store。

## 2. Redux 状态（src/store/index.ts）

- 统一在 Redux store 管理对话消息、连接状态、当前角色、UI 状态；组件通过 react-redux 读写，禁止散落局部全局态。
- 角色切换经 `src/config/index.ts` 的 `setRole` / `getRole`，切换时同步更新模型路径（`modelPath`）与音色（`voice`）。

## 3. Live2D 形象（live2d/Live2DViewer.tsx）

- 模型资源走配置（`VITE_LIVE2D_MODEL_PATH` / `ROLES[].modelPath`），资源在 `public/models/`（商业模型不入库）；加载失败给占位/降级提示。
- 形象动作由 `Live2DAction { action, emotion?, duration_ms? }` 驱动（`action`: nod/smile/shake/blink/transform/idle），与 Coze 回复的 `AgentReplyPayload.action` 联动，禁止散落 setTimeout 直接操作 DOM。

## 4. 语音（tts/speech.ts）

- 封装 TTS，提供方由 `VITE_TTS_PROVIDER` 配置（web-speech / azure / openai）；按 `AgentReplyPayload.text` / `emotion` 播报。
- 浏览器不支持或授权失败时降级到文本展示，不阻塞核心对话；音频资源及时释放。

## 5. 核心编排（hooks/useDigitalHuman.ts）

- 统一管理：连接状态、消息列表、输入、发送、形象动作、角色；对外暴露清晰接口供组件消费。
- 发送流程：用户输入 → 经 WebSocket/Coze 发送 → 接收 `AgentReplyPayload` → 更新消息列表 → 驱动 TTS 播报与 Live2D 动作。
- 页面/组件卸载时清理 WebSocket、心跳定时器、语音资源。

## 6. 消息类型（types/index.ts）

- `ChatMessage { role: 'user'|'assistant', content, type: 'text'|'voice'|'command', timestamp, meta?:{ emotion?, scene?, task? } }`；`AgentReplyPayload { text, emotion, scene, action? }`；枚举明确，禁止隐式结构。

## 7. 代码注释

关键逻辑标注 `// @ticket 07-数字人Coze：<说明>`。
