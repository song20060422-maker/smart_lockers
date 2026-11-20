---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_8c4ee4f0e27c11f19919525400638852
    ReservedCode1: 9X0oo5a3t4Ly2Jrz9ALpdtfWTZKwLRVAo4Z4qeRr/7rTfaqY8R0Ms0nbehwmcmmWCUJo6GYMuYPxES3WpnL+ZLbSmsXI1lL0+UoD6G89hDJ+iZxAZpIxwKJhHVHsNdB+vl9W7sdaAWP+HDm6QtimbH1C1ir+vNrCDvwY8Slbl+IDdwnND/olB6oS6qc=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_8c4ee4f0e27c11f19919525400638852
    ReservedCode2: 9X0oo5a3t4Ly2Jrz9ALpdtfWTZKwLRVAo4Z4qeRr/7rTfaqY8R0Ms0nbehwmcmmWCUJo6GYMuYPxES3WpnL+ZLbSmsXI1lL0+UoD6G89hDJ+iZxAZpIxwKJhHVHsNdB+vl9W7sdaAWP+HDm6QtimbH1C1ir+vNrCDvwY8Slbl+IDdwnND/olB6oS6qc=
---

# 06 · 安全规范

> 适用：所有涉及鉴权、资金、状态变更、风险的工单（尤其 01/02/04/05/06/10）。前置：`00-project.md`。

## 1. 鉴权与会话

- 受保护接口统一鉴权；按角色区分权限：普通用户、快寄员、网点管理员、运营/系统管理员。
- 退出登录（06）：立即清除会话/令牌并记日志；令牌加入失效名单，退出后无法访问受保护资源；多次调用幂等。
- 临时开柜（05）：限时限次授权，权限校验 + 并发锁，操作留审计。

## 2. 敏感信息与脱敏

- 手机号（`185****6789`）、姓名、地址、余额、柜门位置等按场景脱敏；日志与监控禁止打印完整敏感值。
- 头像上传（06）：校验格式与大小，防恶意文件注入；存储路径不可执行；返回路径可访问且全局一致。
- 定位数据不持久化（03）；分析接口需鉴权或内网限制（08）。

## 3. 幂等与防重放

- 资金与状态变更接口强制幂等（唯一键/幂等键）；支付回调校验签名、防重放、状态机校验（02/04）。
- 前端二次确认 + 后端防抖双重防护（支付、补缴、开柜、退出）。
- 频控：通知发送去重频控（01）按接收人+场景+指纹+时间窗；登录/支付/上传按维度限流。

## 4. 审计

- 关键动作必须审计：资金变动、订单状态流转、柜门开关/换柜、临时开柜、补缴、退出登录、通知发送结果、高风险运维动作。
- 审计记录含：操作人/设备、时间、动作、对象、前后状态、结果、request_id；只追加不可篡改。
- 通知全链路审计（01）：模板、渠道、送达状态、重试升级、失败原因可查。

## 5. 高风险动作确认（工单 10 核心）

- 运维动作分级：低/中/高；高风险（重启、隔离、限流、回滚）必须"无确认不得执行"，有审计可追溯。
- Skills 需定义：风险等级、是否允许自动执行、失败/超时返回语义；高风险动作仅提示与占位，需人工确认。
- 运维智能体编排遵循：识别→计划→调用 Skills→汇总；结论可解释、基于真实上下文，信息不足须说明缺口。

## 6. 凭据红线

- 仓库与文档禁止出现真实密码、密钥、Token；配置文件仅放本地占位值。
- 不修改 `.git`、不提交 `.env`/私钥（见 `.gitignore`）。

## 7. 代码注释

安全关键逻辑标注工单编号：`// @ticket 05-订单查看与取出：临时开柜限时限次校验`。
