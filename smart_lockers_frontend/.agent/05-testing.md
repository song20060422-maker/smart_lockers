---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_2bb446dce4ab11f1b820525400de85a5
    ReservedCode1: OQCdu7kMSpScLTa1DDOKCyhwU/xmU4/SRqI+5I+eyE0W3k95b3z4uToqOmXezjywziSLKe9mOpDPCn/AkS7uZ1mX1Dw90jQTeZ0tN8oBRyFpOlIkarkupp72nZg0vL2oEg5ojBSon0ShUUtDY+FmvRyEFVbUzOeIPN7JNCI1eDmygewkODWdxzFIrew=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_2bb446dce4ab11f1b820525400de85a5
    ReservedCode2: OQCdu7kMSpScLTa1DDOKCyhwU/xmU4/SRqI+5I+eyE0W3k95b3z4uToqOmXezjywziSLKe9mOpDPCn/AkS7uZ1mX1Dw90jQTeZ0tN8oBRyFpOlIkarkupp72nZg0vL2oEg5ojBSon0ShUUtDY+FmvRyEFVbUzOeIPN7JNCI1eDmygewkODWdxzFIrew=
---

# 05 · 测试规范 — 用户端前端

> 适用：本站所有测试活动。前置：`00-project.md`、`02-frontend.md`、`04-api.md`。

## 1. 静态检查与构建（必过）

- 类型检查：`npm run type-check`（vue-tsc），零类型错误、避免 any 逃逸。
- 生产构建：`npm run build`，构建通过、产物无致命告警。
- 提交前两项都必须通过。

## 2. 联调测试

- 经 Vite 代理（`/api` → 后端 :18080）做端到端联调，覆盖 8 个视图主流程。
- 关键链路：登录/授权 → 定位与附近寄存点 → 选柜与收费试算 → 支付落单 → 订单列表/详情 → 临时开柜/补交 → 个人中心 → 退出登录。
- 错误分支：接口非 0 码、401 未授权、支付重复提交（幂等）、支付回调异步刷新余额与明细。

## 3. 交互与边界验证

- 敏感操作二次确认与防抖：连点支付/落单/开柜不产生重复请求。
- 定时轮询在页面卸载时正确清理；空列表/加载失败/超时等边界 UI 状态正确。

## 4. 验收口径

- 工单视图交付前：type-check 通过、build 通过、主流程联调通过、关键错误分支验证、敏感信息无泄露。

## 5. 注释

验证用例关联工单：`// @ticket 02-支付宝支付：重复提交幂等验证`。

