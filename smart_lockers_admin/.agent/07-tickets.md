---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_ad6c6a36e4aa11f1b820525400de85a5
    ReservedCode1: CBhRvv2nKX2Wk2nTj6bGPitbFiWdKEtSP+gvHnLE1ivAo6HczGpiFw8DhKtQiJUWAq0NDP1j+eZ01EO7eR6ebNOrbwdGq8n+SratxqQMIFlT9I0pL93a9YeXhZZmCRg+XX/eV0QHtSuTz5B/SSK2epj9O9zcUVpdo8JRxdPW80N256U2Yv4ZMNvBqbw=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_ad6c6a36e4aa11f1b820525400de85a5
    ReservedCode2: CBhRvv2nKX2Wk2nTj6bGPitbFiWdKEtSP+gvHnLE1ivAo6HczGpiFw8DhKtQiJUWAq0NDP1j+eZ01EO7eR6ebNOrbwdGq8n+SratxqQMIFlT9I0pL93a9YeXhZZmCRg+XX/eV0QHtSuTz5B/SSK2epj9O9zcUVpdo8JRxdPW80N256U2Yv4ZMNvBqbw=
---

# 07 · 模块映射速查 — B 端管理后端

> 面向具体模块开发时查本文件；业务细节以真实代码与根工单原文为准。本表只列 B 端后端落点。

## 鉴权（handler/auth.go、middleware/auth.go、username.go）
- 登录签发 JWT、当前用户、退出；`username.go` 提供操作人上下文。
- 账号来源：`config.yaml` 的 `demo_admins`（占位演示）；生产走数据库/密管。

## 网点管理（handler/site.go、service/site.go）
- 网点列表/详情/新增/编辑/停用；权限 `site:read` / `site:write`。
- 数据域：网点管理员仅本网点（`SiteScope`），平台运营跨网点。

## 设备管理（handler/device.go、service/device.go）
- 设备列表/详情/远程指令；权限 `device:read` / `device:control`。
- 远程控制为风险动作：二次确认 + 审计记录。

## 对账（handler/reconcile.go、service/reconcile.go、store/time.go / time2.go）
- 对账任务、执行对账、记录查询；权限 `reconcile:read` / `reconcile:run`（run 仅平台运营）。
- 容差由 `policy.reconcile_tolerance_fen` 配置；时间窗逻辑在 store/time*.go。

## 告警（handler/alarm.go、service/alarm.go）
- 告警列表/详情/处置；权限 `alarm:read` / `alarm:handle`。
- 级别枚举由 `policy.alarm_levels` 配置；处置为风险动作：确认 + 审计。

## 审计（pkg/audit.go、handler 各模块）
- 敏感操作留痕：操作人/角色/网点/动作/对象/结果/时间；查询权限 `audit:read`（平台运营）。
- 保留条数由 `policy.audit_max_records` 配置。

## 通用（router/wire.go、pkg/*）
- `wire.go` 依赖装配；`router.go` 路由注册。
- `errcode.go` 错误码、`response.go` 响应、`log.go` 日志、`mask.go` 脱敏、`audit.go` 审计。

