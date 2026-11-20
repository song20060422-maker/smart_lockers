---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_9d675d69e4aa11f1b820525400de85a5
    ReservedCode1: zyC5gUWInVZjAq/lwyp4dBDbFT1rG+71cO9Ua+4NK0W4l2XUBGHUkAYdBQX8TcPypgIuSatMCojdtV8kNGjpESQjADvHQJRWhFufREK6Ldfdt6yoRySCq3d7oH/C1rIwQX0uRl5N3f1e0LrABvsRUxfSd9JqP57QLCdg7Npvw2CZiTvba8GU8hTVB/g=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_9d675d69e4aa11f1b820525400de85a5
    ReservedCode2: zyC5gUWInVZjAq/lwyp4dBDbFT1rG+71cO9Ua+4NK0W4l2XUBGHUkAYdBQX8TcPypgIuSatMCojdtV8kNGjpESQjADvHQJRWhFufREK6Ldfdt6yoRySCq3d7oH/C1rIwQX0uRl5N3f1e0LrABvsRUxfSd9JqP57QLCdg7Npvw2CZiTvba8GU8hTVB/g=
---

# 03 · 数据库规范 — 用户端后端（GORM + MySQL / Redis）

> 适用：本服务涉及数据表、索引、缓存的开发。前置：`00-project.md`、`01-backend.md`。

## 1. 选型与连接

- ORM：GORM v1.25.12 + MySQL 驱动 v1.5.7；连接配置在 `config/config.yaml` 的 `database` 段（dsn、max_open_conns、max_idle_conns、conn_max_lifetime_sec）。
- 缓存：go-redis/v9 v9.7.3，配置在 `redis` 段（addr、password、db）。
- 公共连接封装在 `internal/pkg/db.go`、`internal/pkg/redis.go`；禁止在 handler/service 中散落新建连接。

## 2. 表设计约定

- 数据结构定义在 `internal/model/entity.go`，表名与字段风格全服务统一。
- 主键、创建/更新时间、软删除字段按 GORM 约定或显式声明，保持一致。
- 建表与种子数据脚本：`scripts/init.sql`（已填充 17 表种子数据），结构变更须同步更新脚本。
- 核心业务表（依工单）：账户余额/余额明细/充值订单/充值记录（02）、寄存订单（04）、订单与主详情（05）、用户/卡券（06）、通知模板/发送记录（01）。

## 3. 索引与约束

- 查询字段（城市、寄存点坐标、订单状态、用户 ID、创建时间）建合适索引；附近寄存点用地理索引支撑距离排序（03）。
- 订单历史大数据量用游标分页（cursor + limit），对 `(用户, 状态, 时间)` 建联合索引。
- 幂等用业务唯一键 + 唯一索引兜底（充值单号、支付流水号）。

## 4. 缓存策略

- 定位结果短时缓存、订单详情缓存 + 精准失效、热点寄存点缓存。
- 一致性：先更新 DB 再删缓存 / 延迟双删；金额、状态以 DB 为准。
- 缓存 key 命名统一前缀 + 业务域 + 标识，设置合理 TTL，避免永久 key。

## 5. 事务边界（工单 02/04/05）

- 资金变动、寄存落单、订单状态流转、异常换柜在 service 层事务内完成；事务内不调用外部慢 IO（如 TTS/通知发送），改为事务后异步。
- 支付回调幂等入账：以充值单号/交易号为幂等键，状态机校验后入账，重复回调返回同一结果。

## 6. 红线

- 配置里的真实 DSN/密码只在被 `.gitignore` 排除的 `config.yaml`，禁止入库到 `config.example.yaml` 或代码。
- 不臆造表名、字段或索引；一切以 `entity.go` 与 `scripts/init.sql` 为准。

