---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_b43423bee4aa11f18663525400cd780f
    ReservedCode1: sSnKpUQulv5zlqoUXfccGy2x5KVpdWUJipV61DUYfLO7ZF6zTCUWWst4pIXwZ0yj2EtO6CWilA1gy5bPzmDCf938hQEhbTgoC8s15Jd3VLp4o1TX0p62UDi7xHeW3augIYmF8YNEJ3qA+jURt3rcugpnXGdqciJdug4YFwpTMywl7uq70oSkshIefXQ=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_b43423bee4aa11f18663525400cd780f
    ReservedCode2: sSnKpUQulv5zlqoUXfccGy2x5KVpdWUJipV61DUYfLO7ZF6zTCUWWst4pIXwZ0yj2EtO6CWilA1gy5bPzmDCf938hQEhbTgoC8s15Jd3VLp4o1TX0p62UDi7xHeW3augIYmF8YNEJ3qA+jURt3rcugpnXGdqciJdug4YFwpTMywl7uq70oSkshIefXQ=
---

# 02 · 前端规范 — B 端后台（Vue3 + Pinia + RBAC）

> 适用：本模块所有页面/组件开发。前置：`00-project.md`、`04-api.md`。

## 1. Vue 组件

- 统一 `<script setup lang="ts">` 组合式 API；视图放 `src/views/<名称>View.vue`，通用布局放 `src/layouts/`。
- props/emits/响应式显式类型，避免 any；列表/表单/弹窗等复杂交互拆分子组件，逻辑下沉到 composables。

## 2. 状态（Pinia）

- 登录态/角色/权限统一在 `src/stores/auth.ts`（`useAuthStore`）：`user`、`token`、`role`、`permissions`、`loggedIn`、`can(perm)`、`login`、`fetchMe`、`logout`。
- 组件通过 store 读写，禁止散落全局态；角色由 `AdminUser.role` 计算，权限由 `PERMISSIONS[role]` 派生。

## 3. 路由与权限

- `src/router/index.ts` 采用 history 模式；`beforeEach` 校验 `getToken()`，未登录跳 `/login`。
- 业务路由挂在 `MainLayout.vue` 下，通过 `meta.perm` 标注所需权限（site:read / device:read / reconcile:read / alarm:read）；菜单与按钮级用 `can(perm)` 控制，前端权限仅展示，服务端须二次校验。

## 4. 通知（composables/useNotify.ts）

- 成功/失败提示统一走 `useNotify`，禁止组件各自实现 alert；业务错误由 `request.ts` 拦截器抛出后统一提示。

## 5. 四大业务模块

- 网点管理 `SitesView.vue`、设备监控 `DevicesView.vue`、充值对账 `ReconcileView.vue`、告警处理 `AlarmsView.vue`；写操作（对账确认、告警处理）须带幂等键。

## 6. 代码注释

关键逻辑标注 `// @ticket B端-前端：<说明>`。
