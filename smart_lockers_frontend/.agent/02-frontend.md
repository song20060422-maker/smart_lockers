---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_294ac28ee4ab11f1b820525400de85a5
    ReservedCode1: azzWA+xVmGH/2CqMCiPLSOh3N/9J2GdXlbDyErFVoPOcD59LYSmYXyni0HY1kQarD0SrEZQFuAHnfvrCFuZPjJ6MxbaE45Ko1qIcso3GotpJzTDLhOoicX/I7g+DqnNVOzA+IGSC9FF4fRe60zJD/VRw/ZaURpMYX+8ZwMghynWlid1yszSvatiRpHo=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_294ac28ee4ab11f1b820525400de85a5
    ReservedCode2: azzWA+xVmGH/2CqMCiPLSOh3N/9J2GdXlbDyErFVoPOcD59LYSmYXyni0HY1kQarD0SrEZQFuAHnfvrCFuZPjJ6MxbaE45Ko1qIcso3GotpJzTDLhOoicX/I7g+DqnNVOzA+IGSC9FF4fRe60zJD/VRw/ZaURpMYX+8ZwMghynWlid1yszSvatiRpHo=
---

# 02 · 前端开发规范 — 用户端（Vue3 + Vite）

> 适用：本站所有页面/组件开发。前置：`00-project.md`、`04-api.md`。

## 1. 组件约定

- Vue3 Composition API + `<script setup>`；组件大驼峰命名，props/emits 显式声明类型。
- 页面放 `src/views/`，通用组件按需抽到 `components/`（如新增）；复杂逻辑抽 composable 或 Pinia，避免巨型组件。

## 2. Pinia 状态

- `src/stores/auth.ts`（登录态/token/用户信息）、`src/stores/location.ts`（定位/城市）。
- 跨页共享状态入 store，页面私有状态用 ref/reactive；退出登录须清空 auth 状态并使令牌失效。

## 3. 路由

- `src/router/index.ts` 统一管理；受保护页面做登录态守卫，未登录跳登录/授权入口。
- 视图与工单：Home(03)、Notify(01)、Pay(02)、Deposit(04)、Orders/OrderDetail(05)、User(06)、Ops(08)。

## 4. 请求层（src/api）

- `request.ts` 统一：baseURL（取自 `VITE_API_TARGET`）、超时、token 注入、错误码拦截、重复请求取消、loading 与错误提示。
- 业务 API 分文件：`auth/home/deposit/order/pay/user/notify/ops.ts`；视图只调 api，不直接 axios。
- 类型与后端报文对齐（错误码 unified、分页 `{list,total,page,page_size}`，订单历史游标分页 `{list,next_cursor,has_more}`）。

## 5. 交互与安全（详见 06-security.md）

- 实时状态（进行中订单、柜门）用轮询，页面卸载时清理定时器，注意节流。
- 敏感操作二次确认 + 幂等防抖，防止重复提交；手机号等脱敏展示（如 `185****6789`）。

## 6. 原型对齐

- 涉及原型的视图标注 `// @proto 原型对齐：<视图标识>`，映射见根 `../AGENT.MD` 第 2.1 节（如 v-orders/v-fee/v-sites）；运营中台专属视图无对应时如实标注，禁止臆造。

## 7. 代码注释

关键逻辑标注 `// @ticket 0X-<工单名>：<说明>`。

