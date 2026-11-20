---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_a47e1bb2e4aa11f1b820525400de85a5
    ReservedCode1: Sy+Umn8pHBl6C0QdaSc+RwrqRP+fDVUSa03QxTemMypybMinmNWzyi+YPiSK6umkPw4Yb89ykDctwfREidbb/UPynzzd4kRN74poW/rR6pSiVT02GNGYz3ri1N5xk6e+nClcjegU9mjjjf1NaKu8CiNBideDy2EcHmLf7WXx5TlgadeqKbU3WGOsjQc=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_a47e1bb2e4aa11f1b820525400de85a5
    ReservedCode2: Sy+Umn8pHBl6C0QdaSc+RwrqRP+fDVUSa03QxTemMypybMinmNWzyi+YPiSK6umkPw4Yb89ykDctwfREidbb/UPynzzd4kRN74poW/rR6pSiVT02GNGYz3ri1N5xk6e+nClcjegU9mjjjf1NaKu8CiNBideDy2EcHmLf7WXx5TlgadeqKbU3WGOsjQc=
---

# 智能快寄柜 · B 端运营后台（smart_lockers_admin）

面向运营人员的 Go + Gin 后端服务，监听 `:18081`，提供网点、设备、对账、告警管理与三角色 RBAC。工程规则见 [AGENT.MD](AGENT.MD) 与 [.agent/](.agent/README.md)。

## 一、技术栈

| 项 | 选型 |
|---|---|
| Go | 1.23 |
| Web | Gin v1.10.0 |
| 鉴权 | JWT（golang-jwt/v5）+ 自研 RBAC |
| 存储 | 自研 `internal/store/` 抽象（无 GORM / 无 Redis） |
| 配置 | YAML（gopkg.in/yaml.v3） |

## 二、目录结构

```
cmd/server/main.go          启动入口（:18081）
config/config.example.yaml  配置模板（demo_admins / policy）
internal/config/            配置加载
internal/handler/           HTTP 层：auth / site / device / reconcile / alarm / helper
internal/middleware/        auth / rbac（权限矩阵+数据域）/ username / 通用
internal/model/             数据结构（AdminRole 等）
internal/router/            路由注册 / wire 依赖装配
internal/service/           业务层：site / device / reconcile / alarm / helper
internal/store/             数据存储抽象（store.go / time.go / time2.go）
internal/pkg/               audit / errcode / log / mask / response
```

## 三、业务模块与角色

| 模块 | 路由前缀 | 主要权限点 |
|---|---|---|
| 鉴权 | `/api/v1/auth` | 登录、当前用户 |
| 网点 | `/api/v1/site` | site:read / site:write |
| 设备 | `/api/v1/device` | device:read / device:control |
| 对账 | `/api/v1/reconcile` | reconcile:read / reconcile:run |
| 告警 | `/api/v1/alarm` | alarm:read / alarm:handle |
| 审计 | `/api/v1/audit` | audit:read |

三角色：`operator`（平台运营，跨网点+审计）、`site_admin`（网点管理员，限本网点）、`user`（普通用户，仅网点只读）。数据域由 `SiteScope` 约束（site_id=0 表示全部）。

## 四、快速开始

```bash
# 1. 配置
cp config/config.example.yaml config/config.yaml
#   填写 jwt.secret；demo_admins 密码为占位符，仅本地演示

# 2. 启动（:18081）
go run ./cmd/server

# 3. 测试
go test ./...
```

> 安全：`config.yaml` 含 JWT 密钥与演示账号，禁止提交仓库；真实账号走数据库/密管。可配置项：`policy.alarm_levels`、`reconcile_tolerance_fen`、`audit_max_records`。

## 五、开发规范

遵循根 `../AGENT.MD` 与本目录 `AGENT.MD`：B 端工单编号注释、单测覆盖率维持 ≥70%（现状约 90.3%）、统一错误码与日志、RBAC 权限点 + 数据域隔离、审计与脱敏、风险动作确认。

## 六、当前状态

- [x] 网点/设备/对账/告警 4 大模块 + 鉴权
- [x] 三角色 RBAC + 数据域隔离 + 审计 + 脱敏
- [x] 单测覆盖率约 90.3%
- [ ] 按运营需求持续迭代并补齐压测报告
