---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_a10c5505e4aa11f18663525400cd780f
    ReservedCode1: 2aKDsjL3jQDKIWzxP6kPZBAqdOW2x86HgA8pvRNp75zK8BBgZ4ErpP8aVd9/DTNZ88nD+A9XT8yYcJ7rolUcryGOMm7ipmdrp4/FEXNw2/0hjm7yTYpRIIjjYpQ6/eyWc8AO+HtFZarnPoloeJ+4c005SssMl734KVnOvPC5HP/cSmPW1Oggy2fjpsI=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_a10c5505e4aa11f18663525400cd780f
    ReservedCode2: 2aKDsjL3jQDKIWzxP6kPZBAqdOW2x86HgA8pvRNp75zK8BBgZ4ErpP8aVd9/DTNZ88nD+A9XT8yYcJ7rolUcryGOMm7ipmdrp4/FEXNw2/0hjm7yTYpRIIjjYpQ6/eyWc8AO+HtFZarnPoloeJ+4c005SssMl734KVnOvPC5HP/cSmPW1Oggy2fjpsI=
---

# 06 · 安全规范 — 用户端后端

> 适用：本服务所有涉及鉴权、资金、敏感数据的开发。前置：`00-project.md`、`01-backend.md`。

## 1. 鉴权（JWT）

- 受保护接口走 `internal/middleware/auth.go` 的 JWT 校验；令牌自 `golang-jwt/jwt/v5` 签发，密钥取自 `config.yaml` 的 `jwt.secret`（高强度随机串，禁止占位符入库）。
- 退出登录（工单 06）须清会话并使令牌即时失效；越权返回统一 403 错误码。

## 2. 幂等与防重放

- 写接口支持幂等键；同一键重复请求返回同一结果，不产生重复副作用（重复扣费、重复落单、重复补缴）。
- 回调类接口防重放：校验时间戳/随机数 + 签名，重复回调按幂等键短路。

## 3. 支付宝回调验签（工单 02，重点）

- 异步通知（`notify_url`）必须先验签再处理业务；验签使用支付宝公钥（`config.yaml` 的 `alipay.public_key`），禁止跳过验签直接入账。
- 校验交易号/充值单号一致性、金额一致性、交易状态；以交易号为幂等键，状态机校验后入账。
- 私钥/公钥仅存被 `.gitignore` 排除的 `config.yaml`；沙箱网关由 SDK 选择，禁止硬编码真实私钥或写死生产网关。

## 4. 脱敏

- 手机号、余额、地址、头像 URL 等个人信息按需脱敏（如 `185****6789`）；日志不得打印敏感字段全量。

## 5. 审计

- 资金变动、柜门动作、订单状态流转、临时开柜、退出登录须留审计记录（谁/何时/何操作/结果），可追溯。

## 6. 红线

- 不在任何入库文件写入真实密码、私钥、Token；`config.yaml` 不得提交。
- 分析/压测接口鉴权或内网限制访问（工单 08）。
