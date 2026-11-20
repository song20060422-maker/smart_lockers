---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_26cdfe3ae4ab11f1b820525400de85a5
    ReservedCode1: IwsUSYZbh96wDSQ+3z4MNFrLhjTrTHuJrTPyDHGCoT0T9fa9eotpgwA8AK6ZDuaA7hLGABjbq+csjKTn6NOzawan2A09j0Ywe2sMmM2oiPkrBFuQRmfRGHrZo75Z4SH/WWrrNIj/AFGhx9EmOkc1HucEGskIjWEDRhZmFl/fhFA99XQPIRfHPMfQzyk=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_26cdfe3ae4ab11f1b820525400de85a5
    ReservedCode2: IwsUSYZbh96wDSQ+3z4MNFrLhjTrTHuJrTPyDHGCoT0T9fa9eotpgwA8AK6ZDuaA7hLGABjbq+csjKTn6NOzawan2A09j0Ywe2sMmM2oiPkrBFuQRmfRGHrZo75Z4SH/WWrrNIj/AFGhx9EmOkc1HucEGskIjWEDRhZmFl/fhFA99XQPIRfHPMfQzyk=
---

# .agent 分规则集 — 用户端前端（smart_lockers_frontend）

本目录是用户端 Vue3 前端主站的 AI 编码规则集，配合本子项目 `../AGENT.MD` 与仓库根 `../../AGENT.MD` 使用。与数字人（React/Live2D）和 B 端前端规则不同：本站以 Vue3 + Pinia + 统一请求层为重心。

## 加载约定

1. **任何任务**：先读根 `../../AGENT.MD` → 本子项目 `../AGENT.MD` → 本文件 → `00-project.md`。
2. **页面/组件开发**：追加 `02-frontend.md`、`04-api.md`。
3. **涉及登录/支付/取出/个人中心**：必读 `06-security.md`。
4. **具体工单**：追加 `07-tickets.md` 查视图/API/Store 速查。

## 规则文件清单

| 文件 | 主题 | 关键内容 |
|---|---|---|
| `00-project.md` | 项目全局 | Vue3 技术栈边界、目录、命名、质量铁律 |
| `02-frontend.md` | 前端规范 | Vue3 组件、Pinia、路由、请求层、交互安全、原型对齐 |
| `04-api.md` | 接口规范 | 统一报文、错误码、8 类 API 在前端落点 |
| `05-testing.md` | 测试规范 | type-check/build、联调、验收口径 |
| `06-security.md` | 安全规范 | token、脱敏、敏感操作防抖与二次确认 |
| `07-tickets.md` | 工单映射 | 用户侧工单 → 视图/API/Store 速查 |

## 维护约定

- 内容必须来自真实代码或根工单原文，禁止臆造接口、字段或指标。
- 每条规则变更注明来源；与 `../AGENT.MD` 或根 `../../AGENT.MD` 冲突时以更高优先级为准。

