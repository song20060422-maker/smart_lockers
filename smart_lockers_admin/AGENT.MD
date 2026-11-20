---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_a33c3f76e4aa11f18663525400cd780f
    ReservedCode1: rOMzeOf2Vpl20daaeyUv59SAM5gYwn23sj2nbjBUYJziH1YxQvvt0Vdl2kInhsw57GhMXEFyijNHLdJybFmiGa0EsSMQ9w8nQDz4LqNmM/jaeHURh4LqQ1hfeIBwH4uftKNKLg2iWB0d1yVMXCwP+WbbMI11np1PRmNmIgJ9Mtvx4VBq1Qn0XAxYcD0=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_a33c3f76e4aa11f18663525400cd780f
    ReservedCode2: rOMzeOf2Vpl20daaeyUv59SAM5gYwn23sj2nbjBUYJziH1YxQvvt0Vdl2kInhsw57GhMXEFyijNHLdJybFmiGa0EsSMQ9w8nQDz4LqNmM/jaeHURh4LqQ1hfeIBwH4uftKNKLg2iWB0d1yVMXCwP+WbbMI11np1PRmNmIgJ9Mtvx4VBq1Qn0XAxYcD0=
---

# AGENT.MD — smart_lockers_admin B 端管理后端 AI 协作入口

> 本文件是 `smart_lockers_admin/`（B 端运营后台服务）内的 AI 协作细则。
> 仓库规则优先级：根 `../AGENT.MD` > 根 `../.agent/*.md` > **本文件** > 本目录 `.agent/*.md` > 个人习惯。冲突时以更高优先级为准；本文件只细化 B 端后端，不得与根铁律冲突。
> 任何 AI 在本目录产出代码/配置/文档前，必须先读根 `../AGENT.MD` 与本文件，再按本目录 `.agent/` 分规则执行。

---

## 1. 子项目身份

| 项 | 内容 |
|---|---|
| 名称 | 智能快寄柜 · B 端运营后台后端服务 |
| 模块名 | `smart_lockers_admin`（go.mod `module smart_lockers_admin`） |
| 职责 | 网点管理、设备管理、对账、告警处置、RBAC 鉴权、审计与脱敏 |
| 监听端口 | `:18081`（`config/config.example.yaml` → `server.addr`，避开主站 :18080） |
| 角色模型 | 三角色 RBAC：平台运营 `operator`、网点管理员 `site_admin`、普通用户 `user`（`internal/middleware/rbac.go`） |
| 边界 | 独立于用户端后端 `smart_lockers/`；不含前端页面（在 `smart_lockers_admin_frontend/`）；不使用 GORM/Redis（数据经 `internal/store/` 抽象） |

## 2. 技术栈（真实依赖，禁止臆造）

| 层 | 选型 | 版本/依据 |
|---|---|---|
| 语言 | Go | 1.23（go.mod） |
| Web 框架 | github.com/gin-gonic/gin | v1.10.0 |
| 鉴权 | github.com/golang-jwt/jwt/v5 | v5.2.2 |
| 唯一标识 | github.com/google/uuid | v1.6.0 |
| 配置 | gopkg.in/yaml.v3 | v3.0.1 |
| 存储 | 自研 `internal/store/` 抽象（含 time.go/time2.go），**无 GORM、无 Redis 依赖** |

> 注意：本服务 go.mod 中不含 gorm/redis/alipay 依赖，禁止照搬用户端后端的实现假设。

## 3. 目录与代码地图（现状，不得假设为空）

```
smart_lockers_admin/
├── AGENT.MD / README.MD / .agent/     <- 本子项目规则层（本次新增）
├── go.mod / go.sum
├── cmd/server/main.go                 <- 启动入口（:18081）
├── config/config.example.yaml          <- 配置模板（demo_admins / policy）
└── internal/
    ├── config/                         <- 配置加载（含 config_test.go）
    ├── handler/                        <- HTTP 层：auth/site/device/reconcile/alarm/helper（含 handler_test.go）
    ├── middleware/                     <- auth / rbac / username / 通用中间件（均含单测）
    ├── model/entity.go                 <- 数据结构（AdminRole 等）
    ├── router/                         <- router.go / wire.go（含 router_test.go）
    ├── service/                        <- 业务层：site/device/reconcile/alarm/helper（含 service_test.go）
    ├── store/                          <- 数据存储抽象：store.go / time.go / time2.go（含 store_test.go）
    └── pkg/                            <- audit/errcode/log/mask/response（均含单测）
```

## 4. 本端开发铁律（继承根 §4 8 条，另加 B 端特有）

1. **工单编号注释**：标注 `// @ticket B端-<模块>：<说明>`（如 `// @ticket B端-RBAC：数据域约束`）。
2. **单测覆盖率 ≥70%**：现状 go test 覆盖率约 90.3%，新增逻辑须维持并补齐 `*_test.go`。
3. **统一错误码与日志**：只用 `internal/pkg/errcode.go` 与 `log.go`；统一响应走 `response.go`。
4. **接口可压测**：变更接口提供压测口径与报告。
5. **RBAC + 数据域隔离**：所有受保护接口经 RBAC 权限点校验（`rbac.go`），并按 `SiteScope` 限制数据域——网点管理员仅限本网点，平台运营可跨网点（site_id=0 表示全部）。
6. **审计与脱敏**：敏感操作留审计记录（`pkg/audit.go`）；个人信息脱敏（`pkg/mask.go`）；对账容差、告警级别、审计条数走 `policy` 配置。
7. **配置安全**：demo 账号密码为占位符，真实账号走数据库/密管；JWT 密钥由使用者填写，禁止入库真实凭据。
8. **引用真实来源**：结论基于本目录真实代码与根工单原文，禁止臆造接口、字段、指标。

## 5. 分规则索引

| 文件 | 适用范围 |
|---|---|
| `.agent/README.md` | 分规则加载方式与维护约定 |
| `.agent/00-project.md` | B 端技术栈边界、目录、命名、通用质量铁律 |
| `.agent/10-backend.md` | Gin 分层、store 存储抽象、RBAC、配置、错误码日志 |
| `.agent/13-api.md` | RESTful、统一报文、错误码、网点/设备/对账/告警接口 |
| `.agent/14-testing.md` | go test 单测（维持高覆盖率）、联调、压测报告 |
| `.agent/15-security.md` | RBAC 权限矩阵、数据域隔离、审计、脱敏、风险动作确认 |
| `.agent/16-tickets.md` | B 端职责 → 模块/接口/权限点速查 |

## 6. 执行约定

- 保证在本目录直接 `go run ./cmd/server` 即可启动（:18081）；配置集中在 `config/config.yaml`（由模板复制）。
- 不改动 `.git`、`.gitignore` 及根目录已有文件；不提交敏感凭据。
- 一次只推进一个模块，做完即交付并标注 B 端工单编号；所有变更可追溯。
