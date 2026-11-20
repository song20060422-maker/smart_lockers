---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_b08bb4c7e4aa11f18663525400cd780f
    ReservedCode1: U30qz4By4goQaOCNC8XP+JLJZiElaTl0pO7c7TBTtmsxcZhkfCWZuVm1iy3zJN5LtuWP7o0MPlW02qeFpOqJZO56wlQAcoDXqWersbVDxU13VWOv37HDrsuxAWg3LU9s9ZbHHUP6t39jDZWUHC8OjXpZuBBhQKms8ke4RcSCdPg+pyp6ipfbEycrYK0=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_b08bb4c7e4aa11f18663525400cd780f
    ReservedCode2: U30qz4By4goQaOCNC8XP+JLJZiElaTl0pO7c7TBTtmsxcZhkfCWZuVm1iy3zJN5LtuWP7o0MPlW02qeFpOqJZO56wlQAcoDXqWersbVDxU13VWOv37HDrsuxAWg3LU9s9ZbHHUP6t39jDZWUHC8OjXpZuBBhQKms8ke4RcSCdPg+pyp6ipfbEycrYK0=
---

# 智能快寄柜 · B 端运营后台前端（smart_lockers_admin_frontend）

## 一、项目定位

面向运营人员的后台管理前端：网点管理、设备监控、充值对账、告警处理，以及登录鉴权与三角色（admin/operator/user）RBAC 权限控制。对应工单前缀 `B端-前端`（幂等 `B端-幂等`）。

## 二、技术栈

- Vue 3.5（`<script setup>` 组合式 API）
- Vite 6 + @vitejs/plugin-vue 5
- TypeScript 5.7 + vue-tsc 2.2
- vue-router 4.5（history 模式 + 登录守卫）
- Pinia 2.3（登录态/角色/权限）
- axios 1.7（统一请求封装）

## 三、快速开始

```bash
npm install
cp .env.example .env.local   # 按需设置 VITE_API_BASE，默认 /api/v1
npm run dev                  # 开发（默认 5173）
npm run type-check           # 类型检查
npm run build                # 生产构建
```

后端接口由 `../smart_lockers_admin/`（Gin :18081）提供，默认经代理访问 `/api/v1`。

## 四、目录结构

```
src/
├── api/          接口定义(index.ts) + axios 封装(request.ts：Bearer/业务码/幂等键)
├── composables/  通知等组合式函数
├── layouts/      主布局 MainLayout.vue（三角色共用）
├── router/       路由 + 登录守卫 + meta.perm
├── stores/       Pinia auth（登录态/角色/权限 can()）
├── styles/       全局样式
├── types/        ApiResponse / AdminUser / AdminRole / PERMISSIONS
└── views/        Login / Sites / Devices / Reconcile / Alarms
```

## 五、开发约定

1. 接口统一走 `src/api/request.ts`；响应为 `ApiResponse`，业务 `code !== 0` 报错。
2. 鉴权：路由守卫校验 token；菜单/按钮权限用 `useAuthStore().can(perm)`；前端权限仅展示，服务端二次校验。
3. 写操作通过 `idem(key)` 携带 `Idempotency-Key` 防重复提交。
4. 关键逻辑注释 `// @ticket B端-前端：<说明>`；类型检查通过后再提交。
5. 配置走 `VITE_API_BASE`，禁止写死域名/端口；token 存 localStorage 且不入日志。

## 六、AI 协作规则

本目录 `.agent/` 与 `AGENT.MD` 为 AI 协作规则集，加载优先级：根 `../AGENT.MD` > 本 `AGENT.MD` > `.agent/*.md`。规则依据真实代码与根工单，禁止臆造接口与字段。

---
