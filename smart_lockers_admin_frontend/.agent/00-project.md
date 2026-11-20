---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_b31722a4e4aa11f18663525400cd780f
    ReservedCode1: 69u0jZmCoYA9aWnf1YEDcx7yQ/y+bLNN1QzYoaheJx9bzXlHDnRlwm3tHyp1Fhb5QjQhYVpndYdlFihmEAx4qjqIkAlaQnWyQ6PO3qkv62cjBYxhS4kEVWO+Z9RkcIoHIP64UzKAFbui9k5Lh/vOoitQlh1mIbx98vXXuQKncr6uDgA7LCStueUiP6U=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_b31722a4e4aa11f18663525400cd780f
    ReservedCode2: 69u0jZmCoYA9aWnf1YEDcx7yQ/y+bLNN1QzYoaheJx9bzXlHDnRlwm3tHyp1Fhb5QjQhYVpndYdlFihmEAx4qjqIkAlaQnWyQ6PO3qkv62cjBYxhS4kEVWO+Z9RkcIoHIP64UzKAFbui9k5Lh/vOoitQlh1mIbx98vXXuQKncr6uDgA7LCStueUiP6U=
---

# 00 · B 端运营后台前端全局规则

> 适用：本目录所有任务。加载顺序：根 `../../AGENT.MD` → `../AGENT.MD` → 本文件 → 其他分规则。

## 1. 技术栈边界（不得混用）

- 本模块 = B 端前端：Vue 3.5 + Vite 6 + TypeScript 5.7 + vue-router 4.5 + Pinia 2.3 + axios 1.7。
- 与用户端前端（同为 Vue3，但业务域不同）、数字人模块（React）、两个 Gin 后端（Go）结构不同；禁止混用依赖与目录假设。
- 接口来自 `../smart_lockers_admin/`（Gin :18081），默认 `/api/v1`；本模块只做展示与交互，权限须服务端二次校验。

## 2. 目录与命名

- `src/{api,composables,layouts,router,stores,styles,types,views}`；视图 `<名称>View.vue`，布局 `MainLayout.vue`。
- 接口集中在 `src/api/`（`index.ts` 定义、`request.ts` 封装），禁止在组件内散落 axios 调用。
- 新增页面/权限须同步 `src/router/index.ts`（`meta.perm`）、`src/types/api.ts`（`PERMISSIONS`）、本目录 `07-tickets.md`。

## 3. 工单编号注释（强制）

- 关键逻辑标注 `// @ticket B端-前端：<说明>`，幂等处 `// @ticket B端-幂等`。

## 4. 通用质量铁律

1. 提交前 `npm run type-check`（vue-tsc --noEmit）通过，避免 any 逃逸。
2. 所有请求走 `request.ts`，响应 `ApiResponse`，业务 `code !== 0` 统一报错。
3. 鉴权：`router.beforeEach` 校验 token；菜单/按钮用 `useAuthStore().can(perm)` 控制。
4. 写操作通过 `idem(key)` 携带 `Idempotency-Key`。
5. token 存 localStorage（`admin_token`），禁止入日志/URL；退出登录清理。
6. 配置驱动：`VITE_API_BASE`（默认 `/api/v1`），禁止写死域名/端口。

## 5. 标准产出物清单

页面/组件说明、接口对接说明、RBAC 权限矩阵、联调记录、类型检查与构建结果、上线说明与回滚预案。

## 6. 可启动性与配置安全

- 保证 `npm install && npm run dev` 可启动；`cp .env.example .env.local` 按需配置 `VITE_API_BASE`。
- 不修改 `.git`、`.gitignore`、根目录已有文件及业务代码；不提交 token/密钥。

## 7. 红线

- 不臆造接口、字段、指标；结论来自真实代码或根工单原文。
- 一次只推进一个能力，做完即交付并标注 `B端-前端`。
