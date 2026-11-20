---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_25c1efb9e4ab11f18663525400cd780f
    ReservedCode1: haVS0Hfl3KT/xD8SdnbbdOfzaMK74LYHJY41Jxaainic9KrBOmvo32MNqYqAy3Vx6wG+yxyhnW0FmS8NWnqXEBmJwJI1VwuOFJSTtte245JNN8U/aRdjCOQ9KLMVvcuK3rnu+YWEFO86BrLK8pu6o2gV/iWtcbQaw/mUQpt6Lnf8F15YAjkGYH7mIJQ=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_25c1efb9e4ab11f18663525400cd780f
    ReservedCode2: haVS0Hfl3KT/xD8SdnbbdOfzaMK74LYHJY41Jxaainic9KrBOmvo32MNqYqAy3Vx6wG+yxyhnW0FmS8NWnqXEBmJwJI1VwuOFJSTtte245JNN8U/aRdjCOQ9KLMVvcuK3rnu+YWEFO86BrLK8pu6o2gV/iWtcbQaw/mUQpt6Lnf8F15YAjkGYH7mIJQ=
---

# 智能快寄柜 · 用户端前端（smart_lockers_frontend）

面向用户的 Vue3 + Vite + TypeScript 前端主站，对接用户端后端 `:18080`。工程规则见 [AGENT.MD](AGENT.MD) 与 [.agent/](.agent/README.md)。

## 一、技术栈

| 项 | 选型 |
|---|---|
| 框架 | Vue 3（Composition API / `<script setup>`） |
| 构建 | Vite 6 |
| 语言 | TypeScript 5.7（vue-tsc） |
| 状态 | Pinia 2 |
| 路由 | Vue Router 4 |
| 请求 | Axios 1.7 |

## 二、目录结构

```
src/
├── main.ts / App.vue          入口与根组件
├── api/                       request + auth/home/deposit/order/pay/user/notify/ops
├── views/                     8 视图：Home / Notify / Pay / Deposit / Orders / OrderDetail / User / Ops
├── stores/                    Pinia：auth / location
├── router/index.ts            路由
├── types/index.ts             TS 类型（与后端报文对齐）
└── styles.css                 全局样式
```

## 三、视图与工单

| 视图 | 工单 | 说明 |
|---|---|---|
| HomeView | 03 | 定位/选城/附近寄存点/指南 |
| NotifyView | 01 | 消息通知 |
| PayView | 02 | 余额/充值 |
| DepositView | 04 | 选柜/时长/优惠/押金/支付 |
| OrdersView / OrderDetailView | 05 | 订单列表/详情/取出/补交 |
| UserView | 06 | 订单/头像/服务/卡券/退出 |
| OpsView | 08 | 性能观测 |

## 四、快速开始

```bash
cp .env.example .env.development   # 按需配置 VITE_API_TARGET（默认 /api/v1）
npm install
npm run dev                        # Vite 开发服务器，/api 代理到后端 :18080
npm run type-check                 # vue-tsc 类型检查
npm run build                      # 生产构建
```

> 数字人模块是独立 React 子项目（`../digital_human/`，端口 5174），不并入本站依赖；运营后台前端在 `../smart_lockers_admin_frontend/`。

## 五、开发规范

遵循根 `../AGENT.MD` 与本目录 `AGENT.MD`：工单编号注释、type-check 通过、统一请求层、Pinia 管理全局状态、敏感操作二次确认与防抖、个人信息脱敏、原型视图标识对齐。

## 六、当前状态

- [x] 8 视图 + 8 个 api 模块 + 2 个 Pinia store + 路由/类型
- [ ] 按工单持续迭代、补齐交互与联调

