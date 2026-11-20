---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_b6d4d173e4aa11f1b820525400de85a5
    ReservedCode1: WQiclXuXmxD5KQSX8xHsLWPVuibNT869O4qxSObRnp1OJoFpiL8//x8WsuILXUSvlqO1BNQM9sk6Bzg+RYBV8mprRBm2bUmZc6/3IC5L3gC+6yLHqOF85HyS6mYU1CYMmDjdAwOgvAkJEoU5Ucc9o5UWzQrbrssMEcVuS9IkZTi4AA1Htbqr5l6VHZg=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_b6d4d173e4aa11f1b820525400de85a5
    ReservedCode2: WQiclXuXmxD5KQSX8xHsLWPVuibNT869O4qxSObRnp1OJoFpiL8//x8WsuILXUSvlqO1BNQM9sk6Bzg+RYBV8mprRBm2bUmZc6/3IC5L3gC+6yLHqOF85HyS6mYU1CYMmDjdAwOgvAkJEoU5Ucc9o5UWzQrbrssMEcVuS9IkZTi4AA1Htbqr5l6VHZg=
---

# 05 · 测试规范 — B 端后台前端

> 适用：本模块所有测试活动。前置：`00-project.md`、`04-api.md`。

## 1. 类型检查与构建（必过）

- 类型检查：`npm run type-check`（vue-tsc --noEmit）通过，避免 any 逃逸。
- 生产构建：`npm run build`（vite build）通过；`npm run preview` 可预览产物。

## 2. 联调测试要点

- 登录：账号密码登录 → token 写入 localStorage → `fetchMe` 回填用户与角色 → 跳转首页；错误密码/过期 token 处理正确。
- 权限：三角色分别登录，验证菜单与按钮按 `Permission`/`can(perm)` 展示正确；越权路由被守卫拦截或服务端拒绝。
- 四大模块：网点管理、设备监控、充值对账、告警处理的列表/详情/写操作主流程走通。
- 幂等：对账确认、告警处理等写操作带 `Idempotency-Key`，重试不重复执行。
- 业务码：非 0 `code` 由拦截器统一抛错并提示，不出现静默失败。

## 3. 边界与健壮性验证

- token 失效/过期时跳转登录；退出登录清理令牌与用户态。
- 接口超时/网络异常有提示；空列表/加载态展示正确。
- 敏感信息（token）不出现在日志、URL、页面明文。

## 4. 验收口径

- 工单交付前：type-check/build 通过、三角色权限联调通过、写操作幂等验证通过、业务码错误提示正确、无凭证泄露。

## 5. 注释

验证用例标注 `// @ticket B端-前端：<说明>`。

