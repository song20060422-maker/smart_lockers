---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_ba730f02e4aa11f1b820525400de85a5
    ReservedCode1: 7RMkMkFoLz37c/R2i5SWuvNr+kjxafJjodAmoZRvpw//8fYmgWNqGjfFpoBU9MHOaQymVOk30bpgjiJyjc7l1qRm6qporYj0/MVoJYO0Er87US8yjQG2BtzuFiBMa4ehqSUmqJYSK60uv0zAUfSpJ57Cr78+5PFOatpjNZlg4dpIzkES3ZlnRfdTtC8=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_ba730f02e4aa11f1b820525400de85a5
    ReservedCode2: 7RMkMkFoLz37c/R2i5SWuvNr+kjxafJjodAmoZRvpw//8fYmgWNqGjfFpoBU9MHOaQymVOk30bpgjiJyjc7l1qRm6qporYj0/MVoJYO0Er87US8yjQG2BtzuFiBMa4ehqSUmqJYSK60uv0zAUfSpJ57Cr78+5PFOatpjNZlg4dpIzkES3ZlnRfdTtC8=
---

# 06 · 安全规范 — B 端后台前端

> 适用：本模块涉及令牌、权限、输入、配置的处理。前置：`00-project.md`、`04-api.md`。

## 1. 令牌安全

- token 仅存 localStorage（`admin_token`），由 `setToken/getToken` 统一管理；禁止写入日志、URL、页面明文。
- 请求统一由 `request.ts` 注入 `Authorization: Bearer`；退出登录 `logout()` 清理 token 与用户态；过期 token 由守卫/服务端拦截跳登录。

## 2. RBAC 权限

- 三角色（admin/operator/user）权限由 `src/types/api.ts` 的 `PERMISSIONS` 派生，`useAuthStore().can(perm)` 判断。
- **前端权限仅用于展示控制**，服务端必须对每个写操作做二次鉴权；禁止仅依赖前端隐藏就认为安全。

## 3. 输入与内容安全

- 用户输入按文本处理并转义，避免 XSS/注入；接口返回内容按文本渲染，不执行脚本。
- 文件中止上传/导出等内容按类型白名单校验。

## 4. 配置安全

- `VITE_API_BASE` 走 `.env.local`（gitignore 排除），禁止写死域名/端口；生产环境使用 HTTPS。
- 不提交任何密钥/令牌；`.env.example` 仅含占位值。

## 5. 红线

- 不臆造接口、字段、指标；接口与权限标识以 `src/api/`、`src/types/api.ts` 真实代码与根工单原文为准。

