---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_2e3ac21fe4ab11f18663525400cd780f
    ReservedCode1: wmFBMhrYrcWkLY2CxPWuEsVBO0CvvocXGWugXL4FXhQ4phFyr906qos595RVT7nLtFqaVOHlCD5roP3mTfEfd9KJIWXHK8DAwSzbwD95oPJfORcN8Cu/Nar+fCpy+Fw3zVWJ5Uy1AfjRazUTeKoFMDVMJd0LZ9O6qKAq/KEShMuRqczt/Ed3x5XgA+A=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_2e3ac21fe4ab11f18663525400cd780f
    ReservedCode2: wmFBMhrYrcWkLY2CxPWuEsVBO0CvvocXGWugXL4FXhQ4phFyr906qos595RVT7nLtFqaVOHlCD5roP3mTfEfd9KJIWXHK8DAwSzbwD95oPJfORcN8Cu/Nar+fCpy+Fw3zVWJ5Uy1AfjRazUTeKoFMDVMJd0LZ9O6qKAq/KEShMuRqczt/Ed3x5XgA+A=
---

# 07 · 工单映射速查 — 用户端前端

> 面向具体工单开发时查本文件；业务细节以根工单原文为准。本表只列用户端前端视图/API/Store 落点。

## 工单 01 · 消息通知（views/NotifyView.vue、api/notify.ts）
- 通知列表/详情/已读；去重频控提示；异常/取件提醒/完成通知按类型展示。

## 工单 02 · 支付宝支付（views/PayView.vue、api/pay.ts、stores/auth.ts）
- 余额查询/明细、支付宝充值、充值记录；充值支付二次确认 + 防抖；回调后刷新余额与明细。

## 工单 03 · 首页（views/HomeView.vue、api/home.ts、stores/location.ts）
- GPS 定位/IP 回退、城市拼音/中文搜索、附近寄存点距离排序分页、指南版本/多语言。

## 工单 04 · 物品寄存（views/DepositView.vue、api/deposit.ts）
- 柜子查询、时长配置与分时定价、优惠券/押金规则、收费试算展示、支付落单、异常换柜。

## 工单 05 · 订单查看与取出（views/OrdersView.vue、OrderDetailView.vue、api/order.ts）
- 当前/历史订单（游标分页）、详情（地点+位置）、临时开柜（限时限次+审计）、超时补交（幂等）、进行中聚合。

## 工单 06 · 个人中心（views/UserView.vue、api/user.ts、api/auth.ts）
- 订单摘要、头像上传（格式/大小校验）、我的服务按角色、卡券可用/过期、退出登录即时失效。

## 工单 08 · 性能观测（views/OpsView.vue、api/ops.ts，如需要）
- 观测指标展示；分析接口鉴权访问。

## 原型对齐
- 视图涉及原型时标注 `// @proto 原型对齐：<视图标识>`，映射见根 `../AGENT.MD` 第 2.1 节；运营中台专属视图无对应时如实标注。

## 数字人（工单 07）
- 数字人前端是独立 React 子项目 `../digital_human/`（端口 5174），不并入本站；本站如需入口仅提供跳转/嵌入链接，Live2D 与 Coze 逻辑不在此实现。

