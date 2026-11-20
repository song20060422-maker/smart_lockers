---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_af0d2bc7e4aa11f18663525400cd780f
    ReservedCode1: gcT5fTs3VtSk896oAZ7CSH7hsBeFuUr+os4SMWQGQ3ZE0c/J8UMP1klzrGolF11KzuQFUWQLZaD95m9MAd3tvIG+Cgp0Vv1836bzrjpD1aI37Pz3pVLFiYgOWvp/PX6crinJMBGciD1GZ+gs7b9rIlDfRaO6l0FhXS8nz5kxt76sWf74CtZLg04olXc=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_af0d2bc7e4aa11f18663525400cd780f
    ReservedCode2: gcT5fTs3VtSk896oAZ7CSH7hsBeFuUr+os4SMWQGQ3ZE0c/J8UMP1klzrGolF11KzuQFUWQLZaD95m9MAd3tvIG+Cgp0Vv1836bzrjpD1aI37Pz3pVLFiYgOWvp/PX6crinJMBGciD1GZ+gs7b9rIlDfRaO6l0FhXS8nz5kxt76sWf74CtZLg04olXc=
---

# AGENT.MD — smart_lockers_admin_frontend AI 协作入口

> 本文件是 `smart_lockers_admin_frontend/`（智能快寄柜 · B 端运营后台前端）内的 AI 协作细则。
> 仓库规则优先级：根 `../AGENT.MD` > 根 `../.agent/*.md` > **本文件** > 本目录 `.agent/*.md` > 个人习惯。冲突时以更高优先级为准；本文件只细化 B 端前端，不得与根铁律冲突。
> 任何 AI 在本目录产出代码/配置/文档前，必须先读根 `../AGENT.MD` 与本文件，再按本目录 `.agent/` 分规则执行。

---

## 1. 子项目身份

| 项 | 内容 |
|---|---|
| 名称 | 智能快寄柜 · B 端运营后台前端 |
| 包名 | `smart-lockers-admin-frontend`（`package.json`，`"type": "module"`） |
| 职责 | 网点管理、设备监控、充值对账、告警处理；登录鉴权与三角色 RBAC 按钮级权限控制 |
| 覆盖工单 | B 端（前端侧），注释前缀 `@ticket B端-前端` / `@ticket B端-幂等` |
| 对接后端 | `../smart_lockers_admin/`（Gin，端口 :18081）；API baseURL 默认 `/api/v1` |
| 边界 | 独立前端工程，Vite dev server 默认端口 5173；与用户端前端 `../smart_lockers_frontend/` 分离 |

## 2. 技术栈（真实依赖，禁止臆造）

| 层 | 选型 | 版本/依据 |
|---|---|---|
| 框架 | Vue 3 | vue ^3.5.13（`<script setup>` + 组合式 API） |
| 构建 | Vite | ^6.0.5 + @vitejs/plugin-vue ^5.2.1 |
| 语言 | TypeScript | ^5.7.2 + vue-tsc ^2.2.0（`npm run type-check`） |
| 路由 | vue-router | ^4.5.0（history 模式 + 登录守卫） |
| 状态 | Pinia | ^2.3.0（`src/stores/auth.ts` 登录态/角色/权限） |
| HTTP | axios | ^1.7.9（`src/api/request.ts` 统一封装） |

> 本模块为 Vue3 + Vite6 + TS + Pinia，非 React（数字人）、非 Go 后端；禁止照搬其他子项目依赖假设。

## 3. 目录与代码地图（现状，不得假设为空）

```
smart_lockers_admin_frontend/
├── AGENT.MD / README.MD / .agent/            <- 规则与文档
├── index.html / package.json / tsconfig.json / tsconfig.node.json / vite.config.ts
├── .env.example                              <- VITE_API_BASE 等环境变量模板
└── src/
    ├── main.tsx? -> main.ts / App.vue / env.d.ts
    ├── api/{index.ts, request.ts}            <- 接口集中定义 + axios 封装（Bearer/业务码/幂等键）
    ├── composables/useNotify.ts              <- 通知组合式函数
    ├── layouts/MainLayout.vue                <- 主布局（三角色共用）
    ├── router/index.ts                       <- 路由 + 登录守卫 + meta.perm
    ├── stores/auth.ts                        <- 登录态/角色/权限（can(perm)）
    ├── styles/global.css
    ├── types/api.ts                          <- ApiResponse / AdminUser / AdminRole / PERMISSIONS
    └── views/{LoginView,SitesView,DevicesView,ReconcileView,AlarmsView}.vue
```

## 4. 本端开发铁律（继承根 §4 8 条，另加 B 端前端特有）

1. **工单编号注释**：关键逻辑标注 `// @ticket B端-前端：<说明>`（幂等处 `@ticket B端-幂等`）。
2. **类型检查**：提交前 `npm run type-check`（vue-tsc --noEmit）通过，避免 any 逃逸。
3. **统一请求封装**：所有接口走 `src/api/request.ts`，禁止散落 axios 实例；响应解包为 `ApiResponse`，业务 `code !== 0` 统一报错。
4. **鉴权与 RBAC**：路由守卫校验 token（`getToken`）；菜单/按钮权限用 `useAuthStore().can(perm)` + 路由 `meta.perm`，三角色（admin/operator/user）按钮级控制，前端权限仅作展示，服务端须二次校验。
5. **幂等键**：写操作（充值对账状态变更、告警处理等资金/状态变更）通过 `idem(key)` 附带 `Idempotency-Key`，防重试重复执行。
6. **令牌安全**：token 存 localStorage（`admin_token`），禁止写入日志/URL；退出登录清理 token 与用户态。
7. **可配置化**：API baseURL 走 `VITE_API_BASE`（默认 `/api/v1`），禁止写死域名/端口。
8. **引用真实来源**：接口、字段、权限标识以 `src/` 真实代码与根工单原文为准，禁止臆造。

## 5. 分规则索引

| 文件 | 适用范围 |
|---|---|
| `.agent/README.md` | 分规则加载方式与维护约定 |
| `.agent/00-project.md` | B 端前端技术栈边界、目录、命名、质量铁律 |
| `.agent/11-frontend.md` | Vue 组件、Pinia、路由守卫、RBAC、通知 |
| `.agent/13-api.md` | axios 封装、业务码、幂等键、后端对接 |
| `.agent/14-testing.md` | 类型检查、构建、联调、验收口径 |
| `.agent/15-security.md` | 令牌、RBAC、输入、配置安全 |
| `.agent/16-tickets.md` | B 端工单 → 视图/接口/store 速查 |

## 6. 执行约定

- 保证 `npm install && npm run dev` 可启动；`npm run type-check` 通过；接口对接 `../smart_lockers_admin/`（:18081）。
- 不改动 `.git`、`.gitignore`、根目录已有文件及任何业务代码；不提交 token/密钥。
- 一次只推进一个能力，做完即交付并标注 `B端-前端`；变更可追溯。
