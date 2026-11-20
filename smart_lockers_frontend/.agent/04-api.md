---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_2a87c8cce4ab11f18663525400cd780f
    ReservedCode1: 5KxyoKSi3cghJyGM8l3mHkrLtXFfUUL0fRydcM0kKMdH8xEACgD8rHzmSx3Rgu7EryX29asJjb3wVLb6xr1K7A99xM4xeePi1ak35gGOIqNGU+2gDwHssusE23jF/27kWgJCyGXMW045tvtKCWhRoy7m7mi2hIh9Mfa4IRYxwgS4zeXZ8WekWGq1ZAk=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_2a87c8cce4ab11f18663525400cd780f
    ReservedCode2: 5KxyoKSi3cghJyGM8l3mHkrLtXFfUUL0fRydcM0kKMdH8xEACgD8rHzmSx3Rgu7EryX29asJjb3wVLb6xr1K7A99xM4xeePi1ak35gGOIqNGU+2gDwHssusE23jF/27kWgJCyGXMW045tvtKCWhRoy7m7mi2hIh9Mfa4IRYxwgS4zeXZ8WekWGq1ZAk=
---

# 04 · 接口规范 — 用户端前端对接

> 适用：本站所有接口调用。前置：`00-project.md`、`02-frontend.md`。

## 1. 统一报文

- 响应 `{ code, message, data, request_id, trace_id }`；`code=0` 成功，非 0 由 `request.ts` 拦截器统一处理。
- 分页 `{ list, total, page, page_size }`；订单历史游标分页 `{ list, next_cursor, has_more }`。

## 2. 错误码处理

- 前端只按统一错误码表分支处理，禁止视图内散落魔法数字；通用错误（参数/未授权/限流/系统）走统一提示，业务错误按业务域细分。
- 401/未授权：清登录态并引导重新登录；幂等冲突类错误给用户明确提示，禁止重复提交。

## 3. API 模块与工单落点

| 工单 | API 文件 | 主要能力 |
|---|---|---|
| 01 | notify.ts | 通知列表/详情/已读 |
| 02 | pay.ts | 余额/明细/充值/充值记录 |
| 03 | home.ts | 定位/城市/寄存点/指南 |
| 04 | deposit.ts | 柜子/收费试算/落单/换柜 |
| 05 | order.ts | 订单列表/详情/临时开柜/补交 |
| 06 | user.ts、auth.ts | 订单摘要/头像/服务/卡券/退出、登录 |
| 08 | ops.ts | 性能观测指标（如需要） |

## 4. 请求约定

- baseURL 取 `VITE_API_TARGET`（默认 `/api/v1`，Vite 代理到 :18080）；超时、token 注入、重复取消、loading 与错误提示统一在 `request.ts`。
- 写操作（支付、落单、开柜、补交、退出）前端做幂等防抖，避免连点重复提交。

## 5. 注释

接口封装处标注 `// @ticket 0X-<工单名>：<说明>`。

