---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_8f545a80e4aa11f1b820525400de85a5
    ReservedCode1: cfWdJKP4TtB5K6XJ8xQ91MARpI1KcGyNm4n17jnSqCrhVdi7epcWYe8AHK492VAMEKQUK5kZD4wo4qgO+mHYrfwxtNbtKS6/xDBchp0M5faaz6qbM6f7Rx2BA0vizv2n2i4N5Biv0Mh8onueiPROzWq9AMnxbDCB0OjOxIm4ixdFFuiaKf+pa7KTSCo=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_8f545a80e4aa11f1b820525400de85a5
    ReservedCode2: cfWdJKP4TtB5K6XJ8xQ91MARpI1KcGyNm4n17jnSqCrhVdi7epcWYe8AHK492VAMEKQUK5kZD4wo4qgO+mHYrfwxtNbtKS6/xDBchp0M5faaz6qbM6f7Rx2BA0vizv2n2i4N5Biv0Mh8onueiPROzWq9AMnxbDCB0OjOxIm4ixdFFuiaKf+pa7KTSCo=
---

# .agent 分规则集 — 数字人模块（digital_human）

本目录是数字人模块的 AI 编码规则集，配合本子项目 `../AGENT.MD` 与仓库根 `../../AGENT.MD` 使用。与两个 Vue3 前端规则不同：本模块以 React + Redux + Coze + Live2D + WebSocket + 语音为重心。

## 加载约定

1. **任何任务**：先读根 `../../AGENT.MD` → 本子项目 `../AGENT.MD` → 本文件 → `00-project.md`。
2. **组件/交互开发**：追加 `02-frontend.md`（React/Live2D/语音/Hook）。
3. **对话/连接开发**：追加 `04-api.md`（Coze + WebSocket 消息 schema）。
4. **涉及凭证/配置**：必读 `06-security.md`。
5. **工单 07**：追加 `07-tickets.md` 查速查。

## 规则文件清单

| 文件 | 主题 | 关键内容 |
|---|---|---|
| `00-project.md` | 项目全局 | React 技术栈边界、目录、命名、质量铁律 |
| `02-frontend.md` | 前端规范 | React 组件、Redux、Live2D、语音、useDigitalHuman 编排 |
| `04-api.md` | 接口规范 | Coze 对话、WebSocket 消息 schema、后端支撑约定 |
| `05-testing.md` | 测试规范 | 类型检查/构建、联调/性能脚本、验收口径 |
| `06-security.md` | 安全规范 | Coze 凭证、输入内容、连接安全 |
| `07-tickets.md` | 工单映射 | 07 数字人 → 组件/API/Hook/config 速查 |

## 维护约定

- 内容必须来自真实代码或根工单原文，禁止臆造接口、字段或指标。
- 每条规则变更注明来源；与 `../AGENT.MD` 或根 `../../AGENT.MD` 冲突时以更高优先级为准。
