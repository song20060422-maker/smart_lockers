---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_b1d55d67e4aa11f1b820525400de85a5
    ReservedCode1: XdzpNxDlVkYZkEMB9QKx40t7L21VE2sC02UFlfFwGIDMZDrWyYlA8lPSC6fpI+qghGaAmx5EsfrTatzyxuM6YnRksA4hu3ud3FX0PqVO+YhlZG6iGG0cRliElUS0hzVF8WyN9T2z6mTTpKLMphMwYewWthcsGKGhYxurP7TGkaXdNj5BZHmMCSKeQXY=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_b1d55d67e4aa11f1b820525400de85a5
    ReservedCode2: XdzpNxDlVkYZkEMB9QKx40t7L21VE2sC02UFlfFwGIDMZDrWyYlA8lPSC6fpI+qghGaAmx5EsfrTatzyxuM6YnRksA4hu3ud3FX0PqVO+YhlZG6iGG0cRliElUS0hzVF8WyN9T2z6mTTpKLMphMwYewWthcsGKGhYxurP7TGkaXdNj5BZHmMCSKeQXY=
---

# .agent 分规则集 — B 端运营后台前端（smart_lockers_admin_frontend）

本目录是 B 端前端的 AI 编码规则集，配合本子项目 `../AGENT.MD` 与仓库根 `../../AGENT.MD` 使用。重心：Vue3 + Pinia + vue-router + axios 请求封装 + 三角色 RBAC + 幂等键。

## 加载约定

1. **任何任务**：先读根 `../../AGENT.MD` → 本子项目 `../AGENT.MD` → 本文件 → `00-project.md`。
2. **页面/组件开发**：追加 `02-frontend.md`（Vue 组件/Pinia/路由守卫/RBAC）。
3. **接口/请求开发**：追加 `04-api.md`（axios 封装/业务码/幂等键/后端对接）。
4. **涉及权限/令牌**：必读 `06-security.md`。
5. **B 端工单**：追加 `07-tickets.md` 查速查。

## 规则文件清单

| 文件 | 主题 | 关键内容 |
|---|---|---|
| `00-project.md` | 项目全局 | Vue3 技术栈边界、目录、命名、质量铁律 |
| `02-frontend.md` | 前端规范 | Vue 组件、Pinia auth、路由守卫、RBAC、通知 |
| `04-api.md` | 接口规范 | axios 封装、业务码、幂等键、后端对接约定 |
| `05-testing.md` | 测试规范 | 类型检查/构建、联调、验收口径 |
| `06-security.md` | 安全规范 | 令牌、RBAC、输入、配置安全 |
| `07-tickets.md` | 工单映射 | B 端工单 → 视图/接口/store 速查 |

## 维护约定

- 内容必须来自真实代码或根工单原文，禁止臆造接口、字段或指标。
- 每条规则变更注明来源；与 `../AGENT.MD` 或根 `../../AGENT.MD` 冲突时以更高优先级为准。
