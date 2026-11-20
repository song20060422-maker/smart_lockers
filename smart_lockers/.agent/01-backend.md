---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_9ba72055e4aa11f1b820525400de85a5
    ReservedCode1: LcsSMVa2MODDqqbLOF8/W1iF0UMrw0V4frogw3tLBzFHbFzqFp/ujbumPdoM9miO5xWa/bWzVj1FaK2Xr1pvKmySIuJrVLeKrKPXteQWcsE9j+bE7Ml9gz8pfqiRicstZu5CryGJcAw48lJpp8BYXpjLZWjFhNbUpNEI/o4oIPVGzLRpgDEVrQDwPTQ=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_9ba72055e4aa11f1b820525400de85a5
    ReservedCode2: LcsSMVa2MODDqqbLOF8/W1iF0UMrw0V4frogw3tLBzFHbFzqFp/ujbumPdoM9miO5xWa/bWzVj1FaK2Xr1pvKmySIuJrVLeKrKPXteQWcsE9j+bE7Ml9gz8pfqiRicstZu5CryGJcAw48lJpp8BYXpjLZWjFhNbUpNEI/o4oIPVGzLRpgDEVrQDwPTQ=
---

# 01 · 用户端后端开发规范（Go + Gin）

> 适用：本服务所有后端工单（01–06、08）。前置：`00-project.md`、`03-database.md`、`04-api.md`、`06-security.md`。

## 1. 分层架构

```
handler(HTTP 入参/出参) -> service(业务编排/事务边界) -> model/pkg(数据访问/数据)
                      \-> middleware(JWT鉴权/限流/日志/Recovery)
```

- handler 不写业务逻辑，只做参数校验与响应封装；统一走 `internal/pkg/response.go`，禁止 handler 内散落 `c.JSON`。
- service 是事务与幂等的落点；跨层用明确 DTO/VO，禁止把 entity 直接透传到 HTTP 响应（需脱敏）。
- `internal/pkg/` 提供 db / redis / log / response / errcode / alipay 公共能力。

## 2. Gin 工程约定

- 路由在 `internal/router/router.go` 按业务域分组：`/api/v1/home`、`/api/v1/deposit`、`/api/v1/order`、`/api/v1/pay`、`/api/v1/user`、`/api/v1/notify`、`/api/v1/ops`。
- 中间件链：Recovery → RequestID → AccessLog → Auth(JWT) → RateLimit；受保护接口走 Auth。
- 配置加载集中在 `internal/config/config.go`，启动失败快速失败（fail-fast）并输出清晰错误。

## 3. 事务与幂等（工单 02/04/05 强约束）

- 资金变动（充值、支付、补缴）、柜门状态变更、订单状态流转必须在事务内完成。
- 所有变更类接口支持幂等键（`Idempotency-Key` 头或业务唯一键 + 唯一索引），防止重复扣费、重复落单、重复补缴。
- 支付回调（`internal/pkg/alipay.go` + pay handler）需验签、防重放、状态机校验；异常状态（未支付成功、重复回调、状态不一致）定义明确处理语义。

## 4. 缓存策略（工单 03/05）

- 定位结果短时缓存；订单详情缓存优先 + 精准失效；热点寄存点查询加缓存。
- 缓存与 DB 一致性：更新走「更新 DB 后删缓存」或延迟双删；金额、状态类字段以 DB 为准。

## 5. 错误码与日志

- 统一错误码表见 `04-api.md`，按业务域分段；错误信息对用户友好、对内有上下文。
- 结构化日志必带：request_id、工单编号、用户/设备标识（脱敏）、耗时；禁止打印手机号、密码、token 全量。

## 6. 性能观测（工单 08）

- `ops.go` 提供性能观测接口；pprof/trace/expvar/metrics 类接口须鉴权或限制内网访问，防止敏感数据泄露。
- 采样开销目标 <3%；慢查询结合 SQL 日志定位。

## 7. 各工单后端落点速查

| 工单 | 本目录重点 | 主要文件 |
|---|---|---|
| 01 | 通知事件、模板渲染、去重频控、重试升级、审计监控 | handler/notify.go、service/notify.go |
| 02 | 余额/明细/充值/记录，支付宝回调幂等入账 | handler/pay.go、service/pay.go、pkg/alipay.go |
| 03 | 定位+IP回退、城市搜索、寄存点距离排序分页、指南 | handler/home.go、service/home.go |
| 04 | 柜子查询、收费试算、支付落单、异常换柜、事务 | handler/deposit.go、service/deposit.go、service/billing.go |
| 05 | 列表游标分页、详情聚合、临时开柜限时限次、超时补交 | handler/order.go、service/order.go |
| 06 | 订单查询、头像上传校验、服务配置、卡券校验、退出会话失效 | handler/user.go、service/user.go |
| 08 | 性能分析接口、火焰图/trace、指标 | handler/ops.go |

## 8. 代码注释

每个文件/关键函数标注工单编号：`// @ticket XX-<工单名>：<说明>`。
