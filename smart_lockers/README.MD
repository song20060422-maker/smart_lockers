---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_98173dcbe4aa11f18663525400cd780f
    ReservedCode1: MhzwlzsJJjsWi3lQ0O1C+eWUvI0NiuzMtozlfFdmpG1QuDiwruk/FN2M9mpf6kqWhfjzGHsID6zs+zlecQ4FZ/UKxK9M0NFnHFQXhHlADGAuVTszd44+M/UAoQjBy5zyYjdBKzlR6nsKxJpqdaz+dFXQDSjJCQgwAGOjK//CsBhB40X5vnrT39WHCAU=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_98173dcbe4aa11f18663525400cd780f
    ReservedCode2: MhzwlzsJJjsWi3lQ0O1C+eWUvI0NiuzMtozlfFdmpG1QuDiwruk/FN2M9mpf6kqWhfjzGHsID6zs+zlecQ4FZ/UKxK9M0NFnHFQXhHlADGAuVTszd44+M/UAoQjBy5zyYjdBKzlR6nsKxJpqdaz+dFXQDSjJCQgwAGOjK//CsBhB40X5vnrT39WHCAU=
---

# 智能快寄柜 · 用户端后端（smart_lockers）

面向用户侧（C 端）的 Go + Gin 后端服务，监听 `:18080`，为小程序/用户 App 提供通知、支付、首页、寄存、订单、个人中心等 API。工程规则见 [AGENT.MD](AGENT.MD) 与 [.agent/](.agent/README.md)。

## 一、技术栈

| 项 | 选型 |
|---|---|
| Go | 1.23 |
| Web | Gin v1.10.0 |
| ORM / DB | GORM + MySQL 驱动 |
| 缓存 | Redis（go-redis/v9） |
| 鉴权 | JWT（golang-jwt/v5） |
| 支付 | 支付宝 SDK v3（沙箱，`internal/pkg/alipay.go`） |
| 配置 | YAML（gopkg.in/yaml.v3） |

## 二、目录结构

```
cmd/server/main.go          启动入口（:18080）
config/                     config.example.yaml（模板）/ config.yaml（本地明文，被忽略）
internal/config/            配置加载
internal/handler/           HTTP 层：deposit / home / notify / ops / order / pay / user / helper
internal/service/           业务层：billing / deposit / home / notify / order / pay / user
internal/model/             数据结构（entity.go）
internal/middleware/        鉴权与通用中间件（含单测）
internal/router/            路由注册（/api/v1/...）
internal/pkg/               db / redis / log / response / errcode / alipay（含单测）
scripts/init.sql            建表与种子数据
```

## 三、接口分组（对应 handler）

| 业务域 | 路由前缀 | 工单 | handler |
|---|---|---|---|
| 消息通知 | `/api/v1/notify` | 01 | notify.go |
| 支付/充值 | `/api/v1/pay` | 02 | pay.go |
| 首页 | `/api/v1/home` | 03 | home.go |
| 物品寄存 | `/api/v1/deposit` | 04 | deposit.go |
| 订单/取出 | `/api/v1/order` | 05 | order.go |
| 个人中心 | `/api/v1/user` | 06 | user.go |
| 性能观测 | `/api/v1/ops` | 08 | ops.go |

## 四、快速开始

```bash
# 1. 配置：复制模板并填写本机值
cp config/config.example.yaml config/config.yaml
#   填写 database.dsn / redis.addr / jwt.secret；支付宝参数见模板 alipay 段

# 2. 初始化数据库（MySQL）
mysql -u <user> -p < scripts/init.sql

# 3. 启动（:18080）
go run ./cmd/server

# 4. 测试
go test ./...
```

> 安全：`config.yaml` 含明文凭据（数据库密码、JWT 密钥、支付宝私钥），已被 `.gitignore` 排除，禁止提交仓库。模板中均为占位符。

## 五、开发规范

遵循根 `../AGENT.MD` 与本目录 `AGENT.MD`：代码注释标注工单编号、单测覆盖率 ≥70%、统一错误码与日志、关键接口可压测、资金/状态变更接口事务与幂等、支付回调验签与防重放、敏感信息脱敏与审计。

## 六、当前状态

- [x] 8 handler + 7 service + 中间件 + pkg + 路由骨架
- [x] 支付宝沙箱支付链路（充值单、回调验签、幂等入账）
- [x] 种子数据脚本 scripts/init.sql
- [ ] 按工单持续迭代并补齐压测报告

