---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_86879284e27c11f19919525400638852
    ReservedCode1: REsPvNk9/k305N6OjV540eTrmBt098yIwfLE8Ouo7w7bGQhC5Daq42A5n5j1hxiROvFEeMmUHJ3bk8h9FES8G3uDMPcHgg9Kauk/MEWzOydr5TF4xc1zYYbBbiI5Po0pRI8gCy/6QAKnmMRbaRgsmS1RkOht7dWLZKzgXIxO3ZTVMFtsvQxr/dTr2kc=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_86879284e27c11f19919525400638852
    ReservedCode2: REsPvNk9/k305N6OjV540eTrmBt098yIwfLE8Ouo7w7bGQhC5Daq42A5n5j1hxiROvFEeMmUHJ3bk8h9FES8G3uDMPcHgg9Kauk/MEWzOydr5TF4xc1zYYbBbiI5Po0pRI8gCy/6QAKnmMRbaRgsmS1RkOht7dWLZKzgXIxO3ZTVMFtsvQxr/dTr2kc=
---

# 10 · 后端开发规范（Go + Gin）

> 适用：所有后端工单（01–06、08）。前置：`00-project.md`、`12-database.md`、`13-api.md`、`15-security.md`。

## 1. 分层架构

```
handler(HTTP 入参/出参) -> service(业务编排/事务边界) -> dao(数据访问) -> model/entity(数据结构)
                     \-> middleware(鉴权/限流/日志/恢复)
```

- handler 不写业务逻辑；service 是事务与幂等的落点；dao 只做数据读写；跨层用明确的 DTO/VO，禁止把 entity 直接透传到 HTTP 响应（需脱敏）。

## 2. Gin 工程约定

- 路由按工单业务域分组（`/api/v1/home`、`/api/v1/order`、`/api/v1/pay`、`/api/v1/notify` 等），统一中间件链：Recovery、RequestID、AccessLog、Auth、RateLimit。
- 统一响应封装与全局异常处理中间件，业务错误通过统一错误码返回，禁止 handler 内直接 `c.JSON` 散落返回。
- 配置加载集中在 `configs/`，启动失败快速失败（fail-fast）并输出清晰错误。

## 3. 事务与幂等（工单 02/04/05 强约束）

- 资金变动（充值、支付、补缴）、柜门状态变更、订单状态流转必须在事务内完成，保证一致性。
- 所有变更类接口支持幂等键（如 `Idempotency-Key` 或业务唯一键 + 唯一索引），防止重复扣费、重复落单、重复补缴。
- 支付回调需防重放、校验签名、状态机校验；异常状态（未支付成功、重复回调、状态不一致）需定义明确处理语义。

## 4. 缓存策略（工单 03/05）

- 城市定位结果短时缓存；订单详情缓存优先 + 精准失效；热点寄存点查询加缓存。
- 缓存与 DB 一致性：更新走"更新 DB 后删缓存"或延迟双删；金额、状态类字段以 DB 为准并做一致性校验。

## 5. 错误码与日志

- 统一错误码表（见 `13-api.md`），错误码按业务域分段；错误信息对用户友好、对内有详细上下文。
- 结构化日志必带：request_id、工单编号、用户/设备标识（脱敏）、耗时；禁止打印敏感信息（手机号、密码、token 全量）。

## 6. 性能与可观测（工单 08）

- 订单管理模块集成 pprof（`net/http/pprof`）、trace、expvar/runtime/metrics、block profile、gops、pyroscope/Parca。
- pprof/分析接口必须鉴权或限制内网访问，防止敏感数据泄露；采样开销目标 <3%。
- 关键接口 P95 目标见 `16-tickets.md`；慢查询结合 SQL 日志与火焰图定位。

## 7. 各工单后端落点速查

| 工单 | 后端重点 |
|---|---|
| 01 | 事件模型/路由、模板渲染、去重频控（指纹+窗口）、失败重试升级、定时任务扫描、审计与监控 |
| 02 | 四类支付能力（余额/明细/充值/记录）、资金一致性、充值回调幂等入账 |
| 03 | 定位+IP回退、城市列表/搜索、寄存点距离排序分页、指南版本/多语言 |
| 04 | 柜子查询、收费规则计算、支付落单、异常释放与换柜、事务 |
| 05 | 列表游标分页、详情聚合、临时开柜限时限次、超时补交、进行中聚合、幂等锁 |
| 06 | 订单查询、头像上传校验、服务动态配置、卡券状态校验、退出会话失效 |
| 08 | 性能分析工具链集成、火焰图/trace、Prometheus 指标 |

## 8. 代码注释

每个文件/关键函数标注工单编号：`// @ticket XX-<工单名>：<说明>`。
*（内容由AI生成，仅供参考）*
