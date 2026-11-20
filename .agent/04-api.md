---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_89de8d78e27c11f19919525400638852
    ReservedCode1: ZjujLKAbdOLt1BSW5Ch6YjjTd/mUBCHqGpWkt8c5n4kfwxtVXaFsSuUrHrCReWtW1TX1JrLW1VWRkQGj5uh6Z0wCieIZn/DBcaNBh3JXxbITfU7mgBpRznkgVathvcOtxvGp43lcZotug7cdNdjOOBaRsbQTw7yEZPOyaWvzZ8X7Gugr6Ic95w9ZA/k=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_89de8d78e27c11f19919525400638852
    ReservedCode2: ZjujLKAbdOLt1BSW5Ch6YjjTd/mUBCHqGpWkt8c5n4kfwxtVXaFsSuUrHrCReWtW1TX1JrLW1VWRkQGj5uh6Z0wCieIZn/DBcaNBh3JXxbITfU7mgBpRznkgVathvcOtxvGp43lcZotug7cdNdjOOBaRsbQTw7yEZPOyaWvzZ8X7Gugr6Ic95w9ZA/k=
---

# 04 · 接口规范（RESTful + 统一报文）

> 适用：所有前后端交互接口。前置：`00-project.md`。

## 1. RESTful 约定

- 路径：`/api/v1/<业务域>/<资源>`，名词复数，动作通过 HTTP 方法或子资源表达；版本前缀 v1。
- 方法：GET 查、POST 建、PUT/PATCH 改、DELETE 删；变更类接口需幂等（见 `06-security.md`）。
- 业务域示例：`home`（03）、`locker`/`deposit`（04）、`order`（05）、`pay`/`account`（02）、`user`/`coupon`（06）、`notify`（01）、`digital-human`（07）、`ops`（10）。

## 2. 统一响应报文

```json
{
  "code": 0,
  "message": "ok",
  "data": {},
  "request_id": "xxx",
  "trace_id": "xxx"
}
```

- `code=0` 成功，非 0 为错误码；前后端以同一份错误码表对齐。
- 分页统一返回：`{ list: [], total, page, page_size }`；游标分页返回 `{ list: [], next_cursor, has_more }`（工单 05 明确游标分页）。

## 3. 错误码分段（建议，落库以统一错误码表为准）

| 段 | 业务域 |
|---|---|
| 10000 | 通用（参数错误、未授权、限流、系统错误） |
| 11xxx | 首页/定位（03） |
| 12xxx | 账户/支付（02） |
| 13xxx | 寄存/收费（04） |
| 14xxx | 订单/取出（05） |
| 15xxx | 个人中心/头像/卡券（06） |
| 16xxx | 通知（01） |
| 17xxx | 数字人（07） |
| 18xxx | 运维智能体（10） |

- 每个错误码需：唯一、可读信息、对内详细上下文、对应处理建议；禁止散落魔法数字。

## 4. 分页、限流、幂等

- 列表统一分页参数；订单历史等大数据量用游标分页（`cursor` + `limit`）。
- 写接口限流：登录、支付、充值、开柜、头像上传等按用户/IP/设备维度限流。
- 写接口要求幂等键：`Idempotency-Key` 头或业务唯一键；同一键重复请求返回同一结果，不产生重复副作用。

## 5. 鉴权与传输

- 受保护接口走统一鉴权中间件（token/session），退出登录后令牌立即失效（工单 06）。
- 敏感字段（手机号、余额、地址）按需脱敏；传输 HTTPS；关键写操作二次确认。
- 接口文档（OpenAPI/Swagger）随代码更新，含请求参数、响应示例、错误码。

## 6. 各工单接口落点速查

| 工单 | 主要接口 |
|---|---|
| 01 | 模板管理、发送、发送记录查询、监控统计 |
| 02 | 余额查询、余额明细、支付宝充值、充值记录 |
| 03 | 定位、城市列表/搜索、附近寄存点、寄存指南 |
| 04 | 柜子查询、收费试算、创建寄存订单（支付落单）、异常换柜 |
| 05 | 当前/历史订单列表、订单详情、临时开柜、超时补交、进行中聚合 |
| 06 | 订单摘要、头像上传、我的服务、卡券列表、退出登录 |
| 07 | Coze 配置、对话 WebSocket、知识库查询、状态 |
| 10 | 对话/事件触发、Skills 调用、巡检报告、审计查询 |

## 7. 代码注释

接口定义/路由处标注工单编号：`// @ticket 04-物品寄存：收费试算接口`。
