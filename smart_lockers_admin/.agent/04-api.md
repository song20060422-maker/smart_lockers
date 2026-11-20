---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_aa15d360e4aa11f18663525400cd780f
    ReservedCode1: Us8A3Wd9ukLCb0u+nN0PvHSGIwRdCzm4VZF3vjeAGF8fn3oSCJ+snyim97UpimzsiaMLTOGCJuBcL0bamHM3y9vthlHI/zcfv60+KuR90ucjl5EgtPojfa667eOIZ18sBuLG80TynaOTveZ2Z8bis73UrypbTEl1lLt9Amnja+tFTccYxYSohrhgUFA=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_aa15d360e4aa11f18663525400cd780f
    ReservedCode2: Us8A3Wd9ukLCb0u+nN0PvHSGIwRdCzm4VZF3vjeAGF8fn3oSCJ+snyim97UpimzsiaMLTOGCJuBcL0bamHM3y9vthlHI/zcfv60+KuR90ucjl5EgtPojfa667eOIZ18sBuLG80TynaOTveZ2Z8bis73UrypbTEl1lLt9Amnja+tFTccYxYSohrhgUFA=
---

# 04 · 接口规范 — B 端管理后端（RESTful + 统一报文）

> 适用：本服务所有对外接口。前置：`00-project.md`、`01-backend.md`、`06-security.md`。

## 1. RESTful 约定

- 路径：`/api/v1/<业务域>/<资源>`；方法 GET 查、POST 建/动作、PUT/PATCH 改、DELETE 删。
- 业务域：`auth`、`site`、`device`、`reconcile`、`alarm`、`audit`。
- 每个接口须声明所需权限点（如 `device:control`），并在路由上挂 `RequirePerm`。

## 2. 统一响应报文（pkg/response.go）

```json
{ "code": 0, "message": "ok", "data": {}, "request_id": "xxx" }
```

- `code=0` 成功；分页返回 `{ list, total, page, page_size }`。
- 所有 handler 经统一响应封装，禁止各自拼装 JSON。

## 3. 错误码分段（内部以 errcode.go 为准）

| 段 | 业务域 |
|---|---|
| 10000 | 通用（参数、未授权、限流、系统错误） |
| 20xxx | 鉴权/账号（auth） |
| 21xxx | 网点（site） |
| 22xxx | 设备（device） |
| 23xxx | 对账（reconcile） |
| 24xxx | 告警（alarm） |
| 25xxx | 审计（audit） |

- 权限不足统一 `ErrForb`；每个错误码唯一、可读、有处理建议。

## 4. 接口清单

| 模块 | 主要接口 | 权限点 |
|---|---|---|
| 鉴权 | 登录、当前用户、退出 | 公开 / 已登录 |
| 网点 | 网点列表、详情、新增/编辑/停用 | site:read / site:write |
| 设备 | 设备列表、详情、远程指令 | device:read / device:control |
| 对账 | 对账任务、执行对账、记录查询 | reconcile:read / reconcile:run |
| 告警 | 告警列表、详情、处置 | alarm:read / alarm:handle |
| 审计 | 审计记录查询 | audit:read |

## 5. 数据域与分页

- 所有列表接口须套用 `SiteScope` 数据域过滤：网点管理员仅本网点，平台运营可全部（site_id=0）。
- 列表统一分页参数；写接口按操作人/IP 限流；受保护接口走 JWT。

## 6. 文档与注释

- 接口变更同步更新接口文档（请求参数、响应示例、错误码、所需权限点）。
- 路由/接口定义处标注 `// @ticket B端-<模块>：<说明>`。
