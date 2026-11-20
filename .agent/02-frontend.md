---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_877d8901e27c11f18c5752540024e231
    ReservedCode1: NOvD+4XZqezNkRXK3YHUz0lEM1rwrbkJUZ7Q28DLnM4G/+fu+3tL5wgV9AkJGCfRw6H7WphbxxfB4fGnH4S66xOGocuII1eFLKhKsUhqJ+jigJ4Gx17uucCsu6auZJPeKvOpGWjZU3MEQhO8ZB8TJl+7t1PB07uE8HX2zGj9d+VtKbNWqAlrYBe+SRQ=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_877d8901e27c11f18c5752540024e231
    ReservedCode2: NOvD+4XZqezNkRXK3YHUz0lEM1rwrbkJUZ7Q28DLnM4G/+fu+3tL5wgV9AkJGCfRw6H7WphbxxfB4fGnH4S66xOGocuII1eFLKhKsUhqJ+jigJ4Gx17uucCsu6auZJPeKvOpGWjZU3MEQhO8ZB8TJl+7t1PB07uE8HX2zGj9d+VtKbNWqAlrYBe+SRQ=
---

# 02 · 前端开发规范（Vue3 + Vite）

> 适用：前端主站开发。前置：`00-project.md`、`04-api.md`。数字人模块另见 `08-ai-skills.md`。

## 1. 技术栈

Vue3（Composition API + `<script setup>`）、Vite、TypeScript、Pinia（状态）、Vue Router（路由）、Axios（请求）。UI 库按项目约定统一选型后全站一致。

## 2. 目录结构建议

```
src/
├── views/        # 页面：首页(03)、订单(05)、个人中心(06)、物品寄存(04)、支付/余额(02) 等
├── components/   # 通用组件
├── api/          # 接口封装（按业务域分文件）
├── stores/       # Pinia store
├── router/
├── utils/        # 请求封装、鉴权、埋点
├── types/        # TS 类型定义（与后端报文对齐）
└── assets/
```

## 3. 开发与交互约定

- 页面与工单对应：首页(03：定位/选城/附近柜/指南)、订单(05：当前/历史/详情/取出/补交)、个人中心(06：订单/头像/服务/卡券/退出)、寄存(04：选柜/时长/优惠/押金/支付)、余额充值(02)。
- 请求层统一封装：baseURL、超时、token 注入、错误码拦截、重复请求取消、loading 与错误提示统一。
- 状态类实时刷新（如进行中订单、柜门状态）用轮询或 WebSocket，注意节流与页面卸载时清理定时器/连接。
- 敏感操作（支付、取出、临时开柜、退出登录）需二次确认与幂等防抖，防止重复提交。
- 手机号等个人信息展示需脱敏（如 `185****6789`）。

## 4. 数字人模块特例（工单 07）

- 基于 `awesome-digital-human-live2d`，原生 React + TypeScript + Vite + Live2D + Redux Toolkit。
- 与 Vue3 主前端**独立**：单独子目录、独立 package.json 与依赖，禁止把 React/Live2D 依赖装进主前端。
- 通信：WebSocket（实时对话/状态）+ RESTful（配置/知识库/状态）；消息结构 `{role, content, type, timestamp}`。
- 语音：SpeechRecognition 输入 + SpeechSynthesis/第三方 TTS 输出；回复驱动 Live2D 表情/动作。
- 指标：响应 <2s、知识库准确率 ≥90%、帧率 ≥55fps、音画同步误差 <0.3s、断线重连 100%。
- 集成时通过 iframe/独立路由嵌入主站，API 网关统一鉴权。

## 5. 代码注释

关键逻辑标注工单编号：`// @ticket 03-首页：<说明>`（TS/JS 用 `//` 或 `/* */`）。

