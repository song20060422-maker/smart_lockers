---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_9929e62ce4aa11f18663525400cd780f
    ReservedCode1: BzauQ0HwhV9P0M+GGbIya46XFg3DfNuZizg7t6mKP3mKd/656WseR+rWjcXqC8AYtE+wvo97Qr6rch+TYBS61tiPljx4HLZvuryu7AsdsjWjSV08G8oDpUJuVNvCz23LSyjha/mSj3rTBuOTDsrz3D563T3+boUhfdWDFzdjzeCtq5QCJcP224bMnlE=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_9929e62ce4aa11f18663525400cd780f
    ReservedCode2: BzauQ0HwhV9P0M+GGbIya46XFg3DfNuZizg7t6mKP3mKd/656WseR+rWjcXqC8AYtE+wvo97Qr6rch+TYBS61tiPljx4HLZvuryu7AsdsjWjSV08G8oDpUJuVNvCz23LSyjha/mSj3rTBuOTDsrz3D563T3+boUhfdWDFzdjzeCtq5QCJcP224bMnlE=
---

# .agent 分规则集 — 用户端后端（smart_lockers）

本目录是用户端后端服务的 AI 编码规则集，配合本子项目 `../AGENT.MD` 与仓库根 `../../AGENT.MD` 使用。规则按领域拆分，AI 进入对应开发任务前必须加载对应规则文件。

## 加载约定

1. **任何任务**：先读根 `../../AGENT.MD` → 本子项目 `../AGENT.MD` → 本文件 → `00-project.md`。
2. **后端开发**：追加 `01-backend.md`、`03-database.md`、`04-api.md`、`06-security.md`、`05-testing.md`。
3. **具体工单**：追加 `07-tickets.md` 查用户侧工单的业务细节与指标。
4. **涉及支付/资金**：必读 `06-security.md` 的支付宝回调验签与幂等入账条款。

## 规则文件清单

| 文件 | 主题 | 关键内容 |
|---|---|---|
| `00-project.md` | 项目全局 | Go+Gin 技术栈边界、目录、命名、工单编号注释、质量铁律 |
| `01-backend.md` | 后端规范 | Gin 分层、路由分组、事务/幂等、支付宝、缓存、错误码日志 |
| `03-database.md` | 数据库规范 | GORM 表设计、MySQL 索引、Redis 缓存、init.sql 种子数据 |
| `04-api.md` | 接口规范 | RESTful、统一报文、错误码分段、8 类接口落点 |
| `05-testing.md` | 测试规范 | go test 单测≥70%、联调、压测报告 |
| `06-security.md` | 安全规范 | JWT 鉴权、幂等防重放、支付宝回调验签、脱敏、审计 |
| `07-tickets.md` | 工单映射 | 用户侧 01–06/08 → 模块/接口/指标速查 |

## 维护约定

- 规则随本服务演进更新；新增规则须在本表登记并指明适用工单。
- 规则内容必须来自真实代码或根工单原文，禁止臆造接口、字段或指标。
- 每条规则变更注明来源工单编号；与 `../AGENT.MD` 或根 `../../AGENT.MD` 冲突时，以更高优先级为准。

