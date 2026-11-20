---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_a6295337e4aa11f1b820525400de85a5
    ReservedCode1: eC5K8GLPlZQBJjbSen1j7gviDZOfSiJ3yfPCVCwOdfCdYlt7PyvlkIGT0fV+fDQa5igpZ+BZJ0IUzjqci6CHX+2m5QFhKAWed6hm8R88wY0BYwyMu9/uZSg/M/7pGR1YyxG0YPWOHviMcX4AIAXolkvduJ/ItQh5lwrCKcSr4LFJ13O1Tk1d8Q9q8qs=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_a6295337e4aa11f1b820525400de85a5
    ReservedCode2: eC5K8GLPlZQBJjbSen1j7gviDZOfSiJ3yfPCVCwOdfCdYlt7PyvlkIGT0fV+fDQa5igpZ+BZJ0IUzjqci6CHX+2m5QFhKAWed6hm8R88wY0BYwyMu9/uZSg/M/7pGR1YyxG0YPWOHviMcX4AIAXolkvduJ/ItQh5lwrCKcSr4LFJ13O1Tk1d8Q9q8qs=
---

# .agent 分规则集 — B 端管理后端（smart_lockers_admin）

本目录是 B 端运营后台的 AI 编码规则集，配合本子项目 `../AGENT.MD` 与仓库根 `../../AGENT.MD` 使用。与用户端后端规则不同：本服务以 RBAC、数据域隔离、审计为重心，且不涉及 GORM/Redis。

## 加载约定

1. **任何任务**：先读根 `../../AGENT.MD` → 本子项目 `../AGENT.MD` → 本文件 → `00-project.md`。
2. **后端开发**：追加 `01-backend.md`、`04-api.md`、`06-security.md`、`05-testing.md`。
3. **权限/数据相关**：必读 `06-security.md`（RBAC 权限矩阵 + `SiteScope` 数据域）。
4. **具体模块**：追加 `07-tickets.md` 查 B 端模块/权限点速查。

## 规则文件清单

| 文件 | 主题 | 关键内容 |
|---|---|---|
| `00-project.md` | 项目全局 | Go+Gin 技术栈、目录、命名、质量铁律（B 端） |
| `01-backend.md` | 后端规范 | Gin 分层、store 存储抽象、RBAC 装配、配置、错误码日志 |
| `04-api.md` | 接口规范 | RESTful、统一报文、错误码、网点/设备/对账/告警接口 |
| `05-testing.md` | 测试规范 | go test 单测维持高覆盖率、联调、压测报告 |
| `06-security.md` | 安全规范 | RBAC 权限矩阵、数据域隔离、审计、脱敏、风险动作确认 |
| `07-tickets.md` | 模块映射 | B 端职责 → 模块/接口/权限点速查 |

## 维护约定

- 规则随本服务演进更新；新增规则须在本表登记。
- 内容必须来自真实代码或根工单原文，禁止臆造接口、字段或指标。
- 每条规则变更注明来源；与 `../AGENT.MD` 或根 `../../AGENT.MD` 冲突时，以更高优先级为准。

