---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_ab488b55e4aa11f1b820525400de85a5
    ReservedCode1: NrvwNWKI21B3NlcdqE7FbfSXtm8PuMyhmEoyY2TQ/hHwbtUnDlM7nH+YCgG4e8cR0DyDjkExvBTnC76QHVrtmTeHNd2vVFQy82/D3OvTv1Xew6FU41bahVFKw5Dt9iB9rTitaRjkq0Z9Dk408vaFqdPcyp01GYDn80S1+iEike0IVSpHrjjRgkdbc9s=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_ab488b55e4aa11f1b820525400de85a5
    ReservedCode2: NrvwNWKI21B3NlcdqE7FbfSXtm8PuMyhmEoyY2TQ/hHwbtUnDlM7nH+YCgG4e8cR0DyDjkExvBTnC76QHVrtmTeHNd2vVFQy82/D3OvTv1Xew6FU41bahVFKw5Dt9iB9rTitaRjkq0Z9Dk408vaFqdPcyp01GYDn80S1+iEike0IVSpHrjjRgkdbc9s=
---

# 05 · 测试规范 — B 端管理后端

> 适用：本服务所有测试活动。前置：`00-project.md`、`01-backend.md`、`06-security.md`。

## 1. 单元测试（覆盖率 ≥70%，现状约 90.3%）

- 命令：`go test ./...`；本服务各层已有单测：`config_test.go`、`handler_test.go`、`auth_test.go`、`rbac_test.go`、`middleware_test.go`、`audit_test.go`、`mask_test.go`、`errcode_test.go`、`response_test.go`、`router_test.go`、`service_test.go`、`store_test.go`。
- 新增逻辑必须补齐 `*_test.go`，维持高覆盖率；表驱动，覆盖正常/边界/错误路径。
- **RBAC 重点**：权限矩阵（三角色×权限点）与数据域（`SiteScope`）须有专门用例，验证越权被拒、网点隔离生效。
- **脱敏重点**：`pkg/mask.go` 各类敏感字段的脱敏结果须覆盖。

## 2. 集成/联调测试

- handler → service → store 链路测试；用测试 store 实例，不污染真实数据。
- 鉴权链路：登录签发 JWT → 携带令牌访问受保护接口 → 权限点校验 → 数据域过滤，全链路验证。

## 3. 接口压测与报告

- 关键接口（网点列表、设备控制、对账执行、告警处置）提供压测报告：QPS、P95/P99、错误率、资源占用与解读。
- 对账执行类长任务关注超时与并发控制；采样/分析接口鉴权访问。

## 4. 验收口径

- 工单交付前：`go test ./...` 全绿、覆盖率不下降、RBAC/数据域用例通过、审计记录齐全、敏感信息无泄露。

## 5. 注释

测试文件关键用例标注 `// @ticket B端-<模块>：<说明>`。

