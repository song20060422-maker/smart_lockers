---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_2d0efdf7e4ab11f18663525400cd780f
    ReservedCode1: vq9bTtdtpV+QlNnAiIPcWWkBogKZVz7woL0/EkQCIBTsPvjeRUR7qkAHijxYW2FLdIeU5zXCKhJX1L+b77GUohJpy+dDvtzXUnOo0ezw4TjboOSr4k4TVpsFViLWuCW5arQT5rLr8ipR3H1odhyQx5+kozTsb/UtNWB7Rie87j96B8jVc7+ZzAUsHzo=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_2d0efdf7e4ab11f18663525400cd780f
    ReservedCode2: vq9bTtdtpV+QlNnAiIPcWWkBogKZVz7woL0/EkQCIBTsPvjeRUR7qkAHijxYW2FLdIeU5zXCKhJX1L+b77GUohJpy+dDvtzXUnOo0ezw4TjboOSr4k4TVpsFViLWuCW5arQT5rLr8ipR3H1odhyQx5+kozTsb/UtNWB7Rie87j96B8jVc7+ZzAUsHzo=
---

# 06 · 安全规范 — 用户端前端

> 适用：本站涉及登录、支付、个人信息的开发。前置：`00-project.md`、`02-frontend.md`。

## 1. Token 处理

- 登录成功后 token 由 `src/stores/auth.ts` 管理；`request.ts` 统一注入 Authorization，禁止各视图手动拼 header。
- 退出登录（工单 06）：清本地登录态并使令牌即时失效；401 时统一清态并引导重新登录。
- token/密钥不写入源码、不打印到 console、不落地到可被 XSS 读取的位置（按项目既定方案）。

## 2. 脱敏

- 手机号、账号、地址等个人信息在页面按需脱敏展示（如 `185****6789`）；禁止明文展示全量敏感字段。

## 3. 敏感操作防抖与二次确认

- 支付充值、寄存落单、临时开柜、超时补交、退出登录等：前端做幂等防抖（提交中禁用按钮）+ 二次确认弹窗，防止连点重复提交。
- 重复提交类错误按统一错误码提示，不盲目重试。

## 4. 输入与展示安全

- 用户输入按类型校验（时长、金额、城市搜索）；富文本/服务端返回内容按文本渲染，避免 XSS。
- API 目标走 `VITE_API_TARGET`，禁止硬编码后端地址；`.env` 与密钥不得提交。

## 5. 红线

- 不臆造接口、字段、指标；字段/报文以后端 `types/index.ts` 与真实接口为准。
