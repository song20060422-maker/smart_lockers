---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_9fef497ae4aa11f1b820525400de85a5
    ReservedCode1: C4sUYxvOWnXirCepmfBjnZh5W5Cl19M+sOQeQPcEp9Klgkw/nHDVNI5q3fmfuqRC4UjWcdmcxjjsXh9AB0hUSXk7DcxVT1kZCykGH2xu+IMMZQfVkPGyim6Kr7CsX9gLZV67sXtAsGJtFJDE5bdBiBuIk5EtLdiwbJdYRkrKV2Yw4Ro0AMI//RB74OA=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_9fef497ae4aa11f1b820525400de85a5
    ReservedCode2: C4sUYxvOWnXirCepmfBjnZh5W5Cl19M+sOQeQPcEp9Klgkw/nHDVNI5q3fmfuqRC4UjWcdmcxjjsXh9AB0hUSXk7DcxVT1kZCykGH2xu+IMMZQfVkPGyim6Kr7CsX9gLZV67sXtAsGJtFJDE5bdBiBuIk5EtLdiwbJdYRkrKV2Yw4Ro0AMI//RB74OA=
---

# 05 · 测试规范 — 用户端后端

> 适用：本服务所有测试活动。前置：`00-project.md`、`01-backend.md`、`04-api.md`。

## 1. 单元测试（覆盖率 ≥70%）

- 命令：`go test ./...`；service、pkg、middleware 必须配套 `*_test.go`（现有：`auth_test.go`、`errcode_test.go`、`billing_test.go`、`notify_test.go` 等，需持续补齐）。
- 表驱动测试为主，覆盖正常路径、边界、错误路径；外部依赖（DB/Redis/支付宝）用接口 mock/桩，不连真实环境。
- 资金与状态流转逻辑（收费试算、支付幂等、订单状态机）为单测重点。

## 2. 集成/联调测试

- handler → service → DB/Redis 的链路用测试库或容器；支付宝联调用沙箱（`config.yaml` 中 `sandbox: true`），不触碰真实资金。
- 联调覆盖：充值单生成 → 回调验签 → 幂等入账 → 余额/明细一致性。

## 3. 接口压测与报告

- 关键接口（寄存、订单、支付、附近寄存点）提供压测能力与报告，含 QPS、P95/P99、错误率、资源占用与解读。
- 性能目标参考根 `07-tickets.md`：定位 ≤200ms、收费计算 P95≤250ms、订单查询 P95≤300ms、城市/寄存点 P95≤300ms。

## 4. 验收口径

- 工单交付前：`go test ./...` 全绿、覆盖率达标、关键接口联调通过、压测报告齐备、敏感信息无泄露。
- 缺陷须有对应用例回归。

## 5. 注释

测试文件关键用例标注工单编号：`// @ticket 05-订单查看与取出：超时补交幂等用例`。

