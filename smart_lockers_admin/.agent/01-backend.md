---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_a8fa186de4aa11f18663525400cd780f
    ReservedCode1: AiZK7WYTe9lsLupL1fqvWqA1g3/uE534CDRUj23Gs+wZZ5mn0zskzJBOqFMB1nXUk2CqhqTKEmZyNlEfBavQQOQEI00WQCVvT3TYWhKT9axrZXda9a51+bl918GPfa4QO50zrZuLuKMwWKU0bdkD4oMcPTiveO12LumNtQMhpsZmfnOZrjLzWgURxb4=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_a8fa186de4aa11f18663525400cd780f
    ReservedCode2: AiZK7WYTe9lsLupL1fqvWqA1g3/uE534CDRUj23Gs+wZZ5mn0zskzJBOqFMB1nXUk2CqhqTKEmZyNlEfBavQQOQEI00WQCVvT3TYWhKT9axrZXda9a51+bl918GPfa4QO50zrZuLuKMwWKU0bdkD4oMcPTiveO12LumNtQMhpsZmfnOZrjLzWgURxb4=
---

# 01 · B 端后端开发规范（Go + Gin）

> 适用：本服务所有后端模块。前置：`00-project.md`、`04-api.md`、`06-security.md`。

## 1. 分层架构

```
handler(HTTP 入参/出参) -> service(业务编排) -> store(数据存取)
                      \-> middleware(JWT鉴权 / RBAC / username / 通用)
                      \-> pkg(audit / mask / errcode / response / log)
```

- handler 不写业务逻辑，统一走 `internal/pkg/response.go`；`helper.go` 放公共辅助。
- service 承载业务规则；store 是数据存取抽象（`store.go` + `time.go`/`time2.go`），屏蔽底层存储细节。
- 依赖装配在 `internal/router/wire.go`；路由注册在 `router.go`。

## 2. Gin 工程约定

- 路由按业务域分组：`/api/v1/auth`、`/api/v1/site`、`/api/v1/device`、`/api/v1/reconcile`、`/api/v1/alarm`、`/api/v1/audit`。
- 中间件链：Recovery → RequestID → AccessLog → Auth(JWT) → Username → RBAC(RequirePerm)；受保护接口按权限点挂 `RequirePerm`。
- 配置加载集中在 `internal/config/config.go`，失败快速失败。

## 3. RBAC 装配（本服务核心，工单 B端-RBAC）

- 权限点定义在 `internal/middleware/rbac.go`（`site:read/write`、`device:read/control`、`reconcile:read/run`、`alarm:read/handle`、`audit:read`）。
- 角色→权限映射 `rolePermissions`：`RoleNormalUser` / `RoleSiteAdmin` / `RolePlatformOp`。
- 数据域 `SiteScope`：网点管理员仅本网点，平台运营跨网点（site_id=0 为全部）。所有查询须套用数据域过滤，禁止越权读取其他网点数据。

## 4. 存储抽象（internal/store）

- 数据读写统一经 `store` 接口；`time.go`/`time2.go` 处理时间相关逻辑（对账时间窗、告警时间等）。
- 本服务无 GORM/Redis；如需持久化改造，须在 store 层替换实现，不侵入 handler/service。

## 5. 配置与可配置化

- `config/config.example.yaml` 提供 `demo_admins`（占位）与 `policy`：
  - `alarm_levels`：告警级别枚举（info/warn/critical）。
  - `reconcile_tolerance_fen`：对账容差（分）。
  - `audit_max_records`：审计保留条数。
- 业务逻辑读取 policy，禁止硬编码阈值。

## 6. 错误码与日志

- 错误码 `pkg/errcode.go`，结构化日志 `pkg/log.go`（带 request_id、操作人、网点标识、耗时），禁止打印敏感信息。
- 敏感操作写审计 `pkg/audit.go`；个人信息脱敏 `pkg/mask.go`。

## 7. 各模块落点速查

| 模块 | 重点 | 主要文件 |
|---|---|---|
| 鉴权 | 登录签发 JWT、当前用户 | handler/auth.go、middleware/auth.go |
| 网点 | 网点 CRUD + 数据域 | handler/site.go、service/site.go |
| 设备 | 设备查询/控制（风险动作确认） | handler/device.go、service/device.go |
| 对账 | 对账执行（容差配置）、记录 | handler/reconcile.go、service/reconcile.go |
| 告警 | 告警查询/处置 | handler/alarm.go、service/alarm.go |
| 审计 | 审计记录查询（平台运营） | pkg/audit.go |

## 8. 代码注释

每个文件/关键函数标注 `// @ticket B端-<模块>：<说明>`。
