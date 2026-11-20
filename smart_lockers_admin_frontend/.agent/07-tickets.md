# 07 · 工单映射速查 — B 端后台前端

> 面向 B 端前端工单开发时查本文件；业务细节以根工单原文与真实代码为准。注释前缀 `B端-前端`，幂等 `B端-幂等`。

## 工单 → 落点速查

| 能力 | 落点 | 说明 |
|---|---|---|
| 登录/鉴权 | `src/stores/auth.ts` | login / fetchMe / logout / can(perm)，token 存 localStorage |
| 路由守卫 | `src/router/index.ts` | beforeEach 校验 getToken()，未登录跳 /login；meta.perm 标注权限 |
| 主布局 | `src/layouts/MainLayout.vue` | 三角色共用框架 |
| 请求封装 | `src/api/request.ts` | Bearer 注入、业务码拦截、幂等键（idem） |
| 接口定义 | `src/api/index.ts` | 登录/用户/网点/设备/对账/告警接口集中定义 |
| 类型与权限 | `src/types/api.ts` | ApiResponse / AdminUser / AdminRole / PERMISSIONS |
| 通知 | `src/composables/useNotify.ts` | 成功/失败统一提示 |
| 网点管理 | `src/views/SitesView.vue` | perm: site:read |
| 设备监控 | `src/views/DevicesView.vue` | perm: device:read |
| 充值对账 | `src/views/ReconcileView.vue` | perm: reconcile:read；写操作带幂等键 |
| 告警处理 | `src/views/AlarmsView.vue` | perm: alarm:read；写操作带幂等键 |
| 配置 | `.env.example` → `.env.local` | VITE_API_BASE（默认 /api/v1） |

## 三角色 RBAC

- 角色：admin（管理员）/ operator（运营）/ user（普通），由 `AdminUser.role` 决定。
- 权限标识：`site:read` / `device:read` / `reconcile:read` / `alarm:read`，映射见 `PERMISSIONS`。

## 后端对接

- 对接 `../smart_lockers_admin/`（Gin :18081），路径前缀 `/api/v1`；响应 `ApiResponse`，业务 `code === 0` 成功。

