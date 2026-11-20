---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_ac5bf17ee4aa11f18663525400cd780f
    ReservedCode1: pXgCk3qgtiOOzIrwyP11yO11gomEvKHgMvAOg62i8DTuu0qKLy3cvphaZIfsRpkHifPijA8utfUclF8czXPSgLj/WYC4Fw6iMiwZEEndB20jfv18LfLdILGpcnqCQGooQuIt1ehen6HuuRhZ57NzLXR84zksK2kQPbImibMOOYKpnZ7vhy2SPElW5Cg=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_ac5bf17ee4aa11f18663525400cd780f
    ReservedCode2: pXgCk3qgtiOOzIrwyP11yO11gomEvKHgMvAOg62i8DTuu0qKLy3cvphaZIfsRpkHifPijA8utfUclF8czXPSgLj/WYC4Fw6iMiwZEEndB20jfv18LfLdILGpcnqCQGooQuIt1ehen6HuuRhZ57NzLXR84zksK2kQPbImibMOOYKpnZ7vhy2SPElW5Cg=
---

# 06 · 安全规范 — B 端管理后端

> 适用：本服务所有涉及权限、数据、审计的开发。前置：`00-project.md`、`01-backend.md`。本文件为 B 端安全核心。

## 1. RBAC 权限矩阵（internal/middleware/rbac.go）

三角色权限（来源：真实 `rolePermissions` 映射）：

| 权限点 | 普通用户 user | 网点管理员 site_admin | 平台运营 operator |
|---|:--:|:--:|:--:|
| site:read | ✓ | ✓ | ✓ |
| site:write | — | ✓ | ✓ |
| device:read | — | ✓ | ✓ |
| device:control | — | ✓ | ✓ |
| reconcile:read | — | ✓ | ✓ |
| reconcile:run | — | — | ✓ |
| alarm:read | — | ✓ | ✓ |
| alarm:handle | — | ✓ | ✓ |
| audit:read | — | — | ✓ |

- 受保护接口挂 `RequirePerm(perm)`；无权限统一 `ErrForb`。新增权限点须同步本表与 `07-tickets.md`。

## 2. 数据域隔离（SiteScope）

- `SiteScope(u)`：网点管理员（`RoleSiteAdmin`）仅返回其 `SiteID`；平台运营返回 0（全部）。
- 所有列表/详情/写操作查询必须套用数据域过滤，禁止越权访问其他网点；单测须覆盖隔离效果。

## 3. 鉴权（JWT）

- 登录经 `handler/auth.go` 签发 JWT（`golang-jwt/v5`），密钥取 `config.yaml` 的 `jwt.secret`；受保护接口走 `middleware/auth.go`。
- `demo_admins` 密码为占位符，仅本地演示；真实账号走数据库/密管，禁止入库真实密码。

## 4. 审计（pkg/audit.go）

- 敏感操作（设备控制、对账执行、告警处置、网点变更、登录）写审计记录：操作人、角色、网点、动作、对象、结果、时间。
- 审计保留条数由 `policy.audit_max_records` 配置；审计查询仅 `audit:read`（平台运营）。

## 5. 脱敏（pkg/mask.go）

- 手机号、账号、设备标识等敏感字段按需脱敏展示；日志禁止打印敏感信息全量。

## 6. 风险动作确认

- 设备远程控制、对账执行等风险动作须二次确认；高风险动作无确认不得执行，并留审计记录。

## 7. 红线

- 不提交真实密码、密钥、Token；`config.yaml` 不得入库。
- 对账容差、告警级别以 `policy` 配置为准，禁止硬编码绕过。
