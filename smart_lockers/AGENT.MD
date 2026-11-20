---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_9709970ae4aa11f1b820525400de85a5
    ReservedCode1: YMOZ7oPflRB4ChvFdyR5Cj5rgIXCHVcf3TBOOpcI743UvjD+KVVr5/1o/e9eizz3yMuZLU7o03zt3njnOJad8pQjc06iOfpSXASBTYXLXW98rEKljAIFLxqc0LHmwt1Liy+SKS8oqoZ8xTG30Xc+7rYm/evvopqhcgwcVQTKAhcy1paiVsfLOiBCun0=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_9709970ae4aa11f1b820525400de85a5
    ReservedCode2: YMOZ7oPflRB4ChvFdyR5Cj5rgIXCHVcf3TBOOpcI743UvjD+KVVr5/1o/e9eizz3yMuZLU7o03zt3njnOJad8pQjc06iOfpSXASBTYXLXW98rEKljAIFLxqc0LHmwt1Liy+SKS8oqoZ8xTG30Xc+7rYm/evvopqhcgwcVQTKAhcy1paiVsfLOiBCun0=
---

# AGENT.MD — smart_lockers 用户端后端 AI 协作入口

> 本文件是 `smart_lockers/`（用户侧 C 端后端服务）内的 AI 协作细则。
> 仓库规则优先级：根 `../AGENT.MD` > 根 `../.agent/*.md` > **本文件** > 本目录 `.agent/*.md` > 个人习惯。冲突时以更高优先级为准；本文件只细化用户端后端，不得与根铁律冲突。
> 任何 AI 在本目录产出代码/配置/文档前，必须先读根 `../AGENT.MD` 与本文件，再按本目录 `.agent/` 分规则执行。

---

## 1. 子项目身份

| 项 | 内容 |
|---|---|
| 名称 | 智能快寄柜 · 用户端（C 端）后端服务 |
| 模块名 | `smart_lockers`（go.mod `module smart_lockers`） |
| 职责 | 面向小程序/用户 App 的 8 类业务 API：通知、支付充值、首页、物品寄存、订单取出、个人中心、远程运维性能观测，以及数字人 WebSocket 后端支撑 |
| 监听端口 | `:18080`（`config/config.example.yaml` → `server.addr`） |
| 覆盖工单 | 01/02/03/04/05/06（用户侧主功能）、08（性能观测接口）、07（仅提供后端联调支撑，数字人前端在 `digital_human/`） |
| 边界 | 不含 B 端运营后台（在 `smart_lockers_admin/`）；不含任何前端页面；不实现 Live2D/Coze 前端 |

## 2. 技术栈（真实依赖，禁止臆造）

| 层 | 选型 | 版本/依据 |
|---|---|---|
| 语言 | Go | 1.23（go.mod） |
| Web 框架 | github.com/gin-gonic/gin | v1.10.0 |
| ORM | gorm.io/gorm + gorm.io/driver/mysql | v1.25.12 / v1.5.7 |
| 缓存 | github.com/redis/go-redis/v9 | v9.7.3 |
| 鉴权 | github.com/golang-jwt/jwt/v5 | v5.2.2 |
| 支付 | github.com/smartwalle/alipay/v3（沙箱） | v3.2.31，`internal/pkg/alipay.go` |
| 唯一标识 | github.com/google/uuid | v1.6.0 |
| 配置 | gopkg.in/yaml.v3 | v3.0.1 |

## 3. 目录与代码地图（现状，不得假设为空）

```
smart_lockers/
├── AGENT.MD / README.MD / .agent/     <- 本子项目规则层（本次新增）
├── go.mod / go.sum
├── cmd/server/main.go                 <- 启动入口（:18080）
├── config/
│   ├── config.example.yaml            <- 配置模板（入库）
│   └── config.yaml                    <- 本地明文配置（.gitignore 排除，含支付宝敏感参数）
└── internal/
    ├── config/config.go               <- 配置加载
    ├── handler/                       <- HTTP 入参/出参：deposit/home/notify/ops/order/pay/user/helper
    ├── service/                       <- 业务编排：billing/deposit/home/notify/order/pay/user
    ├── model/entity.go                <- 数据结构
    ├── middleware/                    <- auth / 通用中间件（含单测）
    ├── router/router.go               <- 路由注册（/api/v1/...）
    └── pkg/                           <- db/redis/log/response/errcode/alipay（含单测）
```

## 4. 本端开发铁律（继承根 §4 8 条，另加用户端特有）

1. **工单编号注释**：每个新增/修改文件的关键逻辑处标注 `// @ticket 0X-<工单名>：<说明>`。
2. **单测覆盖率 ≥70%**：handler 之外的 service、pkg、middleware 必须配套单测（`go test ./...`）。
3. **统一错误码与日志**：只用 `internal/pkg/errcode.go` 的错误码与 `log.go` 的结构化日志，禁止散落 `c.JSON` 自定义格式。
4. **接口可压测**：变更接口提供压测口径与报告。
5. **事务 / 幂等 / 鉴权 / 脱敏 / 审计**：资金变动（充值/支付/补缴）、柜门状态、订单状态流转必须事务保护 + 幂等键；支付回调需验签 + 防重放；手机号等脱敏。
6. **支付宝链路**：`config.yaml` 中的 `app_id`/私钥/公钥为敏感信息，仅存被 `.gitignore` 排除的本地文件；沙箱网关由 SDK 选择，禁止硬编码真实私钥入库。
7. **可配置化**：通知渠道与频控、收费规则、JWT 过期等走配置；收费实际规则以数据表为准。
8. **引用真实来源**：结论基于本目录真实代码与根工单原文，禁止臆造接口、字段、指标。

## 5. 分规则索引

| 文件 | 适用范围 |
|---|---|
| `.agent/README.md` | 分规则加载方式与维护约定 |
| `.agent/00-project.md` | 后端技术栈边界、目录、命名、通用质量铁律 |
| `.agent/10-backend.md` | Gin 分层、路由分组、事务/幂等、支付宝、缓存、错误码日志 |
| `.agent/12-database.md` | GORM 表设计、MySQL 索引、Redis 缓存、init.sql 种子数据 |
| `.agent/13-api.md` | RESTful、统一报文、错误码分段、8 类接口落点 |
| `.agent/14-testing.md` | go test 单测≥70%、联调、压测报告 |
| `.agent/15-security.md` | JWT 鉴权、幂等防重放、支付宝回调验签、脱敏、审计 |
| `.agent/16-tickets.md` | 用户侧工单 01–06/08 → 模块/接口/指标速查 |

## 6. 执行约定

- 保证在本目录直接 `go run ./cmd/server` 可启动；配置集中在 `config/config.yaml`，不依赖运行时注入环境变量。
- 不改动 `.git`、`.gitignore` 及根目录已有文件；不提交 `config.yaml`、`.env`、私钥。
- 一次只推进一个工单，做完即交付并标注工单编号；所有变更可追溯。
