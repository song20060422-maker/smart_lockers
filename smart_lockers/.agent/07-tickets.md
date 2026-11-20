---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_a2149c13e4aa11f1b820525400de85a5
    ReservedCode1: Z8SNyEPF+hI3U9onJY/sLy1GOkY2tFLfW4b3lPzsnGcV5+C4gXQqN9kg8aCC6xwgFXN10ODxYOHn0q9Li6kBUEUkvoCJiE85KGdMbuEMUXNhRhpA7NzdeRQOkjTxhiVFaDSNTYbCnRgmdrlRPLBnFCVfTwqENjM4kH9j5yCj0ApQMAKwpUVF5uov7xY=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_a2149c13e4aa11f1b820525400de85a5
    ReservedCode2: Z8SNyEPF+hI3U9onJY/sLy1GOkY2tFLfW4b3lPzsnGcV5+C4gXQqN9kg8aCC6xwgFXN10ODxYOHn0q9Li6kBUEUkvoCJiE85KGdMbuEMUXNhRhpA7NzdeRQOkjTxhiVFaDSNTYbCnRgmdrlRPLBnFCVfTwqENjM4kH9j5yCj0ApQMAKwpUVF5uov7xY=
---

# 07 · 工单映射速查 — 用户端后端

> 面向具体工单开发时查本文件；业务细节以根目录工单原文为准。本表只列用户侧后端落点。

## 工单 01 · 消息通知（handler/notify.go、service/notify.go）
- 多渠道（站内信/短信/邮件/回调）、模板+变量渲染、去重频控（指纹+窗口）、失败重试升级、审计监控。
- 8 类场景：操作失败→管理员、设备异常→网点管理员、订单异常→用户/网点、取件提醒、订单完成、长时间未关门(>1min)、营业前超时预警。

## 工单 02 · 快寄员支付宝支付（handler/pay.go、service/pay.go、pkg/alipay.go）
- 四类能力：余额查询、余额明细、支付宝充值、充值记录。
- 任意余额变动必有明细；回调验签 + 以交易号幂等入账；资金可追溯、可扩展。
- 配置：`config.yaml` alipay 段（沙箱 `sandbox: true`）。

## 工单 03 · 首页（handler/home.go、service/home.go）
- 定位 GPS 优先 + IP 回退、数据不持久化、短时缓存、≤200ms。
- 城市按拼音/中文搜索；附近寄存点距离排序分页（地理索引）、P95≤300ms；寄存指南版本/多语言后台可配。

## 工单 04 · 物品寄存（handler/deposit.go、service/deposit.go、service/billing.go）
- 柜子查询 P95≤200ms；时间配置与分时定价；优惠券+押金规则；收费计算 P95≤250ms。
- 支付完成落单（状态"已占用"+流水+幂等）；异常释放/换柜事务保护。

## 工单 05 · 订单查看与取出（handler/order.go、service/order.go）
- 当前/历史列表游标分页、详情（寄存地点+位置）、临时开柜（限时限次+审计）、超时补交（幂等+解锁/续期）、已结束详情、进行中聚合。
- 状态机 + 动作权限校验、读写分离 + 热点缓存精准失效、幂等键 + 短期锁、审计轨迹、历史归档分区。

## 工单 06 · 个人中心（handler/user.go、service/user.go）
- 订单查询 P95≤300ms；头像上传格式/大小校验 P95≤500ms；我的服务按角色动态配置；卡券可用/过期实时校验 + 下单幂等抵扣；退出登录清会话即时失效。

## 工单 08 · 性能调优观测（handler/ops.go）
- pprof/trace/expvar/metrics/block profile 等分析接口；采样开销 <3%；分析接口鉴权/内网限制；火焰图与慢查询定位。

## 工单 07 · 数字人（本服务仅后端支撑）
- 数字人前端在 `../digital_human/`（React/Live2D）；本服务可提供 RESTful 配置/知识库接口与 WebSocket 对话转发支撑，消息结构 `{role, content, type, timestamp}`；指标响应 <2s。
