---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_9ed1149fe4aa11f18663525400cd780f
    ReservedCode1: gg23JV6WNtjVqfFusKpayCcSWbvO4zfI/oOLSwCG/j7TFWNwn2raHRONhhhLDPH1UDmnItVl3UF6ru6rzOlzpfTyhYuUsHTfcYsC+hmoDia14qg0a0VsWARCp+3oSrE4sktpBk879EKNB57oacVrS+0qfDDyGhq5wuL3Aj5lk7t/W4i+/yHoziP/+eM=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_9ed1149fe4aa11f18663525400cd780f
    ReservedCode2: gg23JV6WNtjVqfFusKpayCcSWbvO4zfI/oOLSwCG/j7TFWNwn2raHRONhhhLDPH1UDmnItVl3UF6ru6rzOlzpfTyhYuUsHTfcYsC+hmoDia14qg0a0VsWARCp+3oSrE4sktpBk879EKNB57oacVrS+0qfDDyGhq5wuL3Aj5lk7t/W4i+/yHoziP/+eM=
---

# 04 · 接口规范 — 用户端后端（RESTful + 统一报文）

> 适用：本服务所有对外接口。前置：`00-project.md`、`01-backend.md`。

## 1. RESTful 约定

- 路径：`/api/v1/<业务域>/<资源>`，名词复数；方法 GET 查、POST 建、PUT/PATCH 改、DELETE 删。
- 业务域：`notify`(01)、`pay`(02)、`home`(03)、`deposit`(04)、`order`(05)、`user`(06)、`ops`(08)。
- 变更类接口需幂等（见 `06-security.md`）。

## 2. 统一响应报文（internal/pkg/response.go）

```json
{ "code": 0, "message": "ok", "data": {}, "request_id": "xxx", "trace_id": "xxx" }
```

- `code=0` 成功，非 0 为错误码；分页返回 `{ list, total, page, page_size }`，历史订单用游标分页 `{ list, next_cursor, has_more }`。
- 所有 handler 经统一响应封装返回，禁止各自拼装 JSON 结构。

## 3. 错误码分段（内部以 errcode.go 为准）

| 段 | 业务域 |
|---|---|
| 10000 | 通用（参数、未授权、限流、系统错误） |
| 11xxx | 首页/定位（03） |
| 12xxx | 账户/支付（02） |
| 13xxx | 寄存/收费（04） |
| 14xxx | 订单/取出（05） |
| 15xxx | 个人中心/卡券（06） |
| 16xxx | 通知（01） |
| 18xxx | 性能/运维接口（08） |

- 每个错误码唯一、可读、对内详细、有处理建议；禁止魔法数字。

## 4. 接口清单（按 handler）

| 工单 | 主要接口 |
|---|---|
| 01 | 模板管理、发送、发送记录查询、监控统计 |
| 02 | 余额查询、余额明细、支付宝充值、充值记录、回调通知 |
| 03 | 定位、城市列表/搜索、附近寄存点、寄存指南 |
| 04 | 柜子查询、收费试算、创建寄存订单、异常换柜 |
| 05 | 当前/历史订单列表、订单详情、临时开柜、超时补交、进行中聚合 |
| 06 | 订单摘要、头像上传、我的服务、卡券列表、退出登录 |
| 08 | 性能观测指标查询（鉴权/内网限制） |

## 5. 分页、限流、鉴权

- 列表统一分页；订单历史用游标分页。
- 写接口按用户/IP/设备维度限流（登录、支付、充值、开柜、头像上传）。
- 受保护接口走 JWT 中间件；退出登录后令牌立即失效；敏感字段按需脱敏；传输 HTTPS。

## 6. 文档与注释

- 接口变更同步更新接口文档（OpenAPI），含请求参数、响应示例、错误码。
- 路由/接口定义处标注工单编号：`// @ticket 04-物品寄存：收费试算接口`。

