---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_23ae5e67e4ab11f1b820525400de85a5
    ReservedCode1: LeHQefOs90HE5FqGoegAp4kLnkBy7llJmpDG+FQF0DuRVi8K21FQPfEsieKlLuC/4mHlGOZMNlopZzvL06snI6Q3wzvqyWBLYj596tRPUcbWPYBLlmuTlbnRICzJQEAZRr22UBQ1yOR6YxQI7BT0iVXNNyCYkHxL+a6KnXcIlQaCV+fbPyrWdQGOxHU=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_23ae5e67e4ab11f1b820525400de85a5
    ReservedCode2: LeHQefOs90HE5FqGoegAp4kLnkBy7llJmpDG+FQF0DuRVi8K21FQPfEsieKlLuC/4mHlGOZMNlopZzvL06snI6Q3wzvqyWBLYj596tRPUcbWPYBLlmuTlbnRICzJQEAZRr22UBQ1yOR6YxQI7BT0iVXNNyCYkHxL+a6KnXcIlQaCV+fbPyrWdQGOxHU=
---

# AGENT.MD — smart_lockers_frontend 用户端前端 AI 协作入口

> 本文件是 `smart_lockers_frontend/`（用户侧 C 端 Vue3 前端主站）内的 AI 协作细则。
> 仓库规则优先级：根 `../AGENT.MD` > 根 `../.agent/*.md` > **本文件** > 本目录 `.agent/*.md` > 个人习惯。冲突时以更高优先级为准；本文件只细化用户端前端，不得与根铁律冲突。
> 任何 AI 在本目录产出代码/配置/文档前，必须先读根 `../AGENT.MD` 与本文件，再按本目录 `.agent/` 分规则执行。

---

## 1. 子项目身份

| 项 | 内容 |
|---|---|
| 名称 | 智能快寄柜 · 用户端前端主站（C 端） |
| 包名 | `smart-lockers-frontend` |
| 职责 | 面向用户的 8 个业务视图：首页/通知/支付余额/寄存/订单/个人中心/运维观测 |
| 对应后端 | `../smart_lockers/`（Gin，:18080），开发代理默认 `/api/v1` → `:18080` |
| 覆盖工单 | 01/02/03/04/05/06（用户侧视图）、08（ops 观测视图）；数字人前端在 `../digital_human/`（独立 React 模块，不并入本站） |
| 边界 | 不含 B 端运营后台（在 `smart_lockers_admin_frontend/`）；不含后端与 Live2D |

## 2. 技术栈（真实依赖，禁止臆造）

| 层 | 选型 | 版本/依据 |
|---|---|---|
| 框架 | Vue 3 | ^3.5.13（Composition API + `<script setup>`） |
| 构建 | Vite | ^6.0.5 |
| 语言 | TypeScript | ^5.7.2（vue-tsc 类型检查） |
| 状态 | Pinia | ^2.3.0 |
| 路由 | Vue Router | ^4.5.0 |
| 请求 | Axios | ^1.7.9 |
| UI 库 | 无第三方 UI 依赖（沿用自研样式 styles.css） |

> 注意：本前端依赖里**没有** react/live2d/element-plus 等，禁止照搬数字人或 B 端前端的依赖假设。

## 3. 目录与代码地图（现状，不得假设为空）

```
smart_lockers_frontend/
├── AGENT.MD / README.MD / .agent/     <- 本子项目规则层（本次新增）
├── index.html / vite.config.ts / tsconfig*.json
├── .env.example                        <- VITE_API_TARGET 等模板
└── src/
    ├── main.ts / App.vue               <- 入口与根组件
    ├── api/                            <- 请求封装：request + auth/home/deposit/order/pay/user/notify/ops
    ├── views/                          <- 8 视图：Home/Notify/Pay/Deposit/Orders/OrderDetail/User/Ops
    ├── stores/                         <- Pinia：auth / location
    ├── router/index.ts                 <- 路由注册
    ├── types/index.ts                  <- TS 类型（与后端报文对齐）
    └── styles.css                      <- 全局样式
```

## 4. 本端开发铁律（继承根 §4 8 条，另加前端特有）

1. **工单编号注释**：TS/JS 关键逻辑标注 `// @ticket 0X-<工单名>：<说明>`。
2. **类型检查**：提交前 `npm run type-check`（vue-tsc）通过，禁止大量 `any` 逃逸。
3. **统一请求层**：所有接口经 `src/api/request.ts`（baseURL、超时、token 注入、错误码拦截、重复请求取消、统一提示），禁止视图内散落 axios 调用。
4. **状态管理**：跨页/全局状态用 Pinia（auth/location），局部状态用 ref/reactive。
5. **交互安全**：敏感操作（支付、取出、临时开柜、退出登录）二次确认 + 幂等防抖；手机号等个人信息脱敏展示。
6. **原型对齐视图标识**：涉及原型视图的组件标注 `// @proto 原型对齐：<视图标识>`（如 v-orders/v-fee），映射见根 `../AGENT.MD` 第 2.1 节。
7. **可配置化**：API 目标走 `VITE_API_TARGET`（`.env.example` 模板），禁止硬编码后端地址到代码。
8. **引用真实来源**：字段/接口与后端 `types/index.ts` 及根工单对齐，禁止臆造。

## 5. 分规则索引

| 文件 | 适用范围 |
|---|---|
| `.agent/README.md` | 分规则加载方式与维护约定 |
| `.agent/00-project.md` | 前端技术栈边界、目录、命名、通用质量铁律 |
| `.agent/11-frontend.md` | Vue3 组件、Pinia、路由、请求层、交互与安全 |
| `.agent/13-api.md` | 统一报文、错误码、8 类接口在前端的落点 |
| `.agent/14-testing.md` | type-check/build/联调、验收口径 |
| `.agent/15-security.md` | token 处理、脱敏、敏感操作防抖与二次确认 |
| `.agent/16-tickets.md` | 用户侧工单 → 视图/API/Store 速查 |

## 6. 执行约定

- 保证本目录 `npm install && npm run dev` 可启动，经 Vite 代理连后端 :18080。
- 不改动 `.git`、`.gitignore` 及根目录已有文件；不提交 `.env`、密钥。
- 一次只推进一个工单视图，做完即交付并标注工单编号；变更可追溯。
