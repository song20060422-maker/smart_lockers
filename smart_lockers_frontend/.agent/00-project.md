---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_27fab8d0e4ab11f1b820525400de85a5
    ReservedCode1: +E6IeDivMXPKP4qW+fMkNh0pxnKvaSxk2YUtcAwDFUOGeRXLMBfDiAtdJKh5+w9ozh0Dtc15kr3neW7iXUQizrJ3+fo5uSSFJDjW/XhV4GyOa+74XazFd/4YrIRlwtn0Fra+pLBDMiGSvLsxXRa40L815LULVlkYzG/VTCSFaCd10M3HC4BTcpS7M8Y=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_27fab8d0e4ab11f1b820525400de85a5
    ReservedCode2: +E6IeDivMXPKP4qW+fMkNh0pxnKvaSxk2YUtcAwDFUOGeRXLMBfDiAtdJKh5+w9ozh0Dtc15kr3neW7iXUQizrJ3+fo5uSSFJDjW/XhV4GyOa+74XazFd/4YrIRlwtn0Fra+pLBDMiGSvLsxXRa40L815LULVlkYzG/VTCSFaCd10M3HC4BTcpS7M8Y=
---

# 00 · 用户端前端全局规则

> 适用：本目录所有任务。加载顺序：根 `../../AGENT.MD` → `../AGENT.MD` → 本文件 → 其他分规则。

## 1. 技术栈边界（不得混用）

- 本站 = 用户端前端：Vue 3.5 + Vite 6 + TypeScript 5.7 + Pinia 2 + Vue Router 4 + Axios 1.7。
- 无第三方 UI 库、无 React/Live2D、无 Element 组件库：不要照搬数字人 `../digital_human/` 或 B 端前端的依赖假设。
- 后端在 `../smart_lockers/`（:18080）；数字人前端独立在 `../digital_human/`（React，:5174）；B 端前端在 `../smart_lockers_admin_frontend/`。

## 2. 目录与命名

- `src/{api,views,stores,router,types}` + `styles.css`；组件大驼峰、目录小写，文件名与视图/API 语义一致。
- 新增视图须同步更新 `src/router/index.ts`、`src/api/`（按需）与本目录 `07-tickets.md`。
- 字段/接口类型在 `src/types/index.ts` 与后端报文对齐，禁止散落隐式 any。

## 3. 工单编号注释（强制）

- TS/JS 关键逻辑处标注 `// @ticket 0X-<工单名>：<说明>`；涉及原型视图再标注 `// @proto 原型对齐：<视图标识>`。

## 4. 通用质量铁律

1. 提交前 `npm run type-check` 通过，避免 any 逃逸。
2. 所有接口经 `src/api/request.ts` 统一封装，禁止视图内散落 axios。
3. 全局状态用 Pinia（auth/location），局部用 ref/reactive。
4. 敏感操作（支付、取出、临时开柜、退出登录）二次确认 + 防抖；个人信息脱敏。
5. API 目标走 `VITE_API_TARGET`，禁止硬编码后端地址。
6. 变更可追溯，命名规范，文档清晰。

## 5. 标准产出物清单

页面/交互说明、接口对接文档（请求/响应/错误码）、原型对齐说明、联调记录、关键流程的验证用例、上线说明与回滚预案。

## 6. 可启动性与配置安全

- 保证 `npm install && npm run dev` 可启动，Vite 代理 `/api` → 后端 :18080。
- `.env.example` 仅放模板；`.env*` 与密钥不得提交；不修改 `.git`、`.gitignore` 与根目录已有文件。

## 7. 红线

- 不臆造接口、字段、指标；结论来自真实代码或根工单原文。
- 一次只推进一个工单视图，做完即交付并标注工单编号。

