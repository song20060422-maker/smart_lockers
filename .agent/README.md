---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_84655cdde27c11f18c5752540024e231
    ReservedCode1: JQb4czEuSutXklgm2Pyp2Rkw/j6CojzuSvBUOmyE9j95U185AJrC9CD9poHUBZ0Nw0EjNom3ARs2ybQ3cEnXMoOqocpWIIiEAqzHXjvUSERQ6/FCSBdkGyvcj1C0DXXCeht12DShVBWf6cI5/QJscQkcq7KR5hpRDCLF1ZIM9BV3ZhL+68KEtdw75Ig=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_84655cdde27c11f18c5752540024e231
    ReservedCode2: JQb4czEuSutXklgm2Pyp2Rkw/j6CojzuSvBUOmyE9j95U185AJrC9CD9poHUBZ0Nw0EjNom3ARs2ybQ3cEnXMoOqocpWIIiEAqzHXjvUSERQ6/FCSBdkGyvcj1C0DXXCeht12DShVBWf6cI5/QJscQkcq7KR5hpRDCLF1ZIM9BV3ZhL+68KEtdw75Ig=
---

# .agent 分规则集

本目录是智能快寄柜项目的 AI 编码规则集，配合最外层 `AGENT.MD` 使用。规则按领域拆分，AI 在进入对应开发任务前必须加载对应规则文件。

## 加载约定

1. **任何任务**：先读 `../AGENT.MD`（顶层入口），再读 `00-project.md`（全局规则）。
2. **后端开发**：追加 `10-backend.md`、`12-database.md`、`13-api.md`、`15-security.md`、`14-testing.md`。
3. **前端开发**：追加 `11-frontend.md`、`13-api.md`。
4. **具体工单**：追加 `16-tickets.md` 查对应工单的业务细节与指标。
5. **AI/Skills 类工单（09/10）或数字人（07）**：追加 `17-ai-skills.md`。

## 规则文件清单

| 文件 | 主题 | 关键内容 |
|---|---|---|
| `00-project.md` | 项目全局 | 技术栈边界、目录、命名、工单编号注释、产出物、质量铁律 |
| `10-backend.md` | 后端规范 | Gin 分层、中间件、事务、幂等、错误码、日志、缓存、pprof |
| `11-frontend.md` | 前端规范 | Vue3 目录、组件、Pinia、Axios、路由、数字人模块特例 |
| `12-database.md` | 数据库规范 | 建表字段、索引、归档分区、缓存失效、核心表清单 |
| `13-api.md` | 接口规范 | RESTful、统一报文、错误码表、分页、限流、鉴权 |
| `14-testing.md` | 测试规范 | 单测≥70%、联调、压测报告、验收标准 |
| `15-security.md` | 安全规范 | 鉴权、脱敏、幂等防重放、审计、高风险动作确认 |
| `16-tickets.md` | 工单映射 | 10 工单 → 模块/表/接口/指标速查 |
| `17-ai-skills.md` | AI 能力规范 | 工单09/10 Skills 定义、工单07 数字人联调 |

## 维护约定

- 规则随项目演进持续更新；新增规则须在本表登记并指明适用工单。
- 规则内容必须来自工单原文或真实代码，禁止臆造接口、字段或指标。
- 每条规则变更需注明来源工单编号；与 `AGENT.MD` 冲突时以 `AGENT.MD` 为准。
*（内容由AI生成，仅供参考）*
