---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_9a7bc5eae4aa11f1b820525400de85a5
    ReservedCode1: CK4twEr1iUXtTfoYbKOz2T4CdXH0S1HZE15vqopDzGEoChztAnTLEfFpOFLoWywPx6SVAUFXa1JUIkyuel5rxx05hsDPNj13pGrfMa/QZjwoTSGIcmjfa2dkpdlOZZLmhSo1Mp+AXYgrwdtJmaEOepNHJQ0sLOVSAXD446+1X9YhDsCGvilJ6kn3/JQ=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_9a7bc5eae4aa11f1b820525400de85a5
    ReservedCode2: CK4twEr1iUXtTfoYbKOz2T4CdXH0S1HZE15vqopDzGEoChztAnTLEfFpOFLoWywPx6SVAUFXa1JUIkyuel5rxx05hsDPNj13pGrfMa/QZjwoTSGIcmjfa2dkpdlOZZLmhSo1Mp+AXYgrwdtJmaEOepNHJQ0sLOVSAXD446+1X9YhDsCGvilJ6kn3/JQ=
---

# 00 · 用户端后端全局规则

> 适用：本目录所有任务。加载顺序：根 `../../AGENT.MD` → `../AGENT.MD` → 本文件 → 其他分规则。

## 1. 技术栈边界（不得混用）

- 本目录 = 用户侧后端：Go 1.23 + Gin v1.10.0 + GORM/MySQL + Redis(v9) + JWT(v5) + 支付宝 SDK v3。
- 前端在 `../smart_lockers_frontend/`（Vue3）；B 端在 `../smart_lockers_admin/`（独立 Gin 服务，:18081）；数字人在 `../digital_human/`（React/Live2D）。
- 本服务只提供 HTTP/JSON API 与必要的 WebSocket 支撑，禁止在本目录引入 React/Vue 前端依赖。

## 2. 目录与命名

- 分层：`cmd/server`（启动）、`config`（配置）、`internal/{config,handler,service,model,middleware,router,pkg}`、`scripts`（SQL）。
- 命名：Go 包/文件小写或短名，导出符号大驼峰；路由 `/api/v1/<业务域>/<资源>`，字段风格全服务统一。
- 新增业务域须同时更新 `internal/router/router.go` 与本目录 `.agent/04-api.md` 的接口清单。

## 3. 工单编号注释（强制）

- 每个新增/修改文件关键逻辑处标注 `// @ticket 0X-<工单名>：<说明>`，如 `// @ticket 04-物品寄存：收费试算`。
- 用户侧工单 01/02/03/04/05/06 与性能 08 均强制；缺失视为不合格。

## 4. 通用质量铁律

1. 单测覆盖率 ≥70%（service / pkg / middleware 必须配套 `*_test.go`）。
2. 错误码只用 `internal/pkg/errcode.go`，日志只用 `internal/pkg/log.go` 结构化输出，禁止自定义散落格式。
3. 关键接口可压测并提供报告与解读。
4. 资金、状态变更、柜门动作：事务 + 幂等 + 鉴权 + 脱敏 + 审计，缺一不可。
5. 收费、通知渠道与频控、JWT 过期等可配置、可动态生效。
6. 变更可追溯，命名规范，文档清晰。

## 5. 标准产出物清单

《设计说明书》、数据字典、ER 图、状态转换图、接口报文设计说明书、建表/索引脚本、接口文档、联调与压测报告、日志与审计规范、上线说明与回滚预案。

## 6. 可启动性与配置安全

- 保证在本目录 `go run ./cmd/server` 即可启动（:18080）；配置集中于 `config/config.yaml`。
- `config.yaml` 含明文敏感值（数据库密码、JWT 密钥、支付宝私钥/公钥），必须被 `.gitignore` 排除，禁止以任何形式入库；`config.example.yaml` 仅放占位符。
- 不修改 `.git`、`.gitignore` 与根目录已有文件；不提交 `.env`/私钥等敏感文件。

## 7. 红线

- 不臆造接口、字段、指标；结论来自真实代码或根工单原文。
- 一次只推进一个工单，做完即交付并标注工单编号。
- 不把真实私钥、Token 写入任何入库文件。
