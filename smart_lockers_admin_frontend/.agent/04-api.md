---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_b55b130be4aa11f1b820525400de85a5
    ReservedCode1: A9XzVda6ES6dhR30zPrWiVHu759apnr4ubv/u/Ehbd178HSBvrT95WCQeaFEs2BkpFLOhbipUlGaWTggGR+4FnRjpWL05cZM7c12BRS8wM/K6i48re15E7ituPRBjXRHIOk5ZOxHHCDEPVGCpQ64zZbBOQ2m3v6S1wihVfw8oN85y2Ovt5ZJRj08gmE=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_b55b130be4aa11f1b820525400de85a5
    ReservedCode2: A9XzVda6ES6dhR30zPrWiVHu759apnr4ubv/u/Ehbd178HSBvrT95WCQeaFEs2BkpFLOhbipUlGaWTggGR+4FnRjpWL05cZM7c12BRS8wM/K6i48re15E7ituPRBjXRHIOk5ZOxHHCDEPVGCpQ64zZbBOQ2m3v6S1wihVfw8oN85y2Ovt5ZJRj08gmE=
---

# 04 · 接口规范 — B 端后台（axios + 业务码 + 幂等）

> 适用：本模块所有接口与请求开发。前置：`00-project.md`、`02-frontend.md`。

## 1. 请求封装（src/api/request.ts，唯一出口）

- 所有请求走 `src/api/request.ts`，禁止组件内自建 axios 实例。
- `baseURL` 取 `VITE_API_BASE || '/api/v1'`，超时 15000ms；接口定义集中在 `src/api/index.ts`。
- 请求拦截注入 `Authorization: Bearer ${token}`（token 来自 localStorage `admin_token`）。
- 响应拦截：后端统一 HTTP 200 + 业务 `code`，`code !== 0` 视为业务错误并 `Promise.reject(new Error(message))`；拦截器已解包为 `ApiResponse`，对外返回类型按 `ApiResponse<T>` 声明。

## 2. 响应结构（src/types/api.ts）

- 统一 `ApiResponse<T> { code: number; message: string; data?: T }`；业务码 0 为成功，非 0 走错误提示。
- 新增接口须同步在 `src/types/api.ts` 定义请求/响应类型，禁止 any。

## 3. 幂等键（写操作强制）

- 资金/状态变更类写操作（充值对账确认、告警处理等）通过 `idem(key)` 附带 `Idempotency-Key` 头，防重试重复执行：`api.xxx(body, idem(key))`（`// @ticket B端-幂等`）。
- 幂等键由业务唯一标识生成（如单号 + 操作类型），保证同一操作重试键一致。

## 4. 后端对接约定

- 接口由 `../smart_lockers_admin/`（Gin :18081）提供，路径前缀 `/api/v1`；登录 `/auth/login`、当前用户 `/auth/me` 等以 `src/api/index.ts` 真实定义为准。
- 权限标识（site:read / device:read / reconcile:read / alarm:read）与后端 RBAC 一致；接口指标响应 <2s（以根工单为准）。

## 5. 注释

接口定义与拦截逻辑处标注 `// @ticket B端-前端：<说明>`，幂等处标注 `// @ticket B端-幂等`。
