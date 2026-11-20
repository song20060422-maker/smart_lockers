---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_a7e8b83ee4aa11f18663525400cd780f
    ReservedCode1: N1HG4SIMYuCpOpRxCoITwtBs8U7pIbM9SE/w/uFD0cMnXcbo0/JfAqZ6ZKCDXTjQofh/cfvyhDMlVJAJweykd8O0zu5vySkeLgHlyrQMXTGKPgbqZEDXs50N79w/Yr7IHSeJ18wc+Mlf/TMsXND+UfoUGR/9EAclNAE18ECuBMTagfx8ayZUhsaVE1Q=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_a7e8b83ee4aa11f18663525400cd780f
    ReservedCode2: N1HG4SIMYuCpOpRxCoITwtBs8U7pIbM9SE/w/uFD0cMnXcbo0/JfAqZ6ZKCDXTjQofh/cfvyhDMlVJAJweykd8O0zu5vySkeLgHlyrQMXTGKPgbqZEDXs50N79w/Yr7IHSeJ18wc+Mlf/TMsXND+UfoUGR/9EAclNAE18ECuBMTagfx8ayZUhsaVE1Q=
---

# 00 · B 端管理后端全局规则

> 适用：本目录所有任务。加载顺序：根 `../../AGENT.MD` → `../AGENT.MD` → 本文件 → 其他分规则。

## 1. 技术栈边界（不得混用）

- 本目录 = B 端运营后台：Go 1.23 + Gin v1.10.0 + JWT(v5) + 自研 `internal/store/` 存储抽象 + YAML 配置。
- **无 GORM、无 Redis、无支付宝 SDK**：不要照搬用户端后端 `../smart_lockers/` 的实现假设。
- 前端在 `../smart_lockers_admin_frontend/`（Vue3）；用户端后端在 `../smart_lockers/`（:18080）；数字人在 `../digital_human/`。

## 2. 目录与命名

- 分层：`cmd/server`、`config`、`internal/{config,handler,service,model,middleware,router,store,pkg}`。
- 新增模块须同时更新 `internal/router/router.go`（或 wire.go）、`internal/middleware/rbac.go` 权限点，以及本目录 `.agent/04-api.md` 与 `07-tickets.md`。
- 命名风格全服务统一；权限点采用 `<资源>:<动作>`（如 `device:control`）。

## 3. 工单编号注释（强制）

- 每个新增/修改文件关键逻辑处标注 `// @ticket B端-<模块>：<说明>`，如 `// @ticket B端-RBAC：数据域约束`。

## 4. 通用质量铁律

1. 单测覆盖率维持 ≥70%（现状约 90.3%），新增逻辑补齐 `*_test.go`。
2. 错误码只用 `pkg/errcode.go`，日志只用 `pkg/log.go`，响应统一走 `pkg/response.go`。
3. 关键接口可压测并提供报告。
4. 写操作：RBAC 鉴权 + 数据域隔离 + 审计 + 脱敏，缺一不可；风险动作需确认。
5. 告警级别、对账容差、审计条数走 `config.yaml` 的 `policy` 配置，可动态生效。
6. 变更可追溯，命名规范，文档清晰。

## 5. 标准产出物清单

《设计说明书》、数据字典、接口报文设计说明书、接口文档、联调与压测报告、权限矩阵说明、审计与脱敏规范、上线说明与回滚预案。

## 6. 可启动性与配置安全

- 保证本目录 `go run ./cmd/server` 即可启动（:18081）；配置集中 `config/config.yaml`（由 `config.example.yaml` 复制）。
- `config.yaml` 含 JWT 密钥与演示账号，必须被 `.gitignore` 排除；`demo_admins` 密码为占位符，真实账号走数据库/密管。
- 不修改 `.git`、`.gitignore` 与根目录已有文件。

## 7. 红线

- 不臆造接口、字段、指标；结论来自真实代码或根工单原文。
- 一次只推进一个模块，做完即交付并标注 B 端工单编号。
- 不提交真实密码、密钥、Token。
