<!-- @ticket 10：设备运维 Skill 规格——设备隔离（占位·需确认） -->
# 设备隔离（占位·需确认）（工单 10）

- **名称**：设备隔离（占位·需确认）
- **用途描述**：高风险动作：隔离设备（暂停纳管/流量）。仅定义确认/审计语义与占位
- **适用场景**：设备疑似故障扩散、安全事件时；需严格控制
- **风险等级**：高
- **是否允许自动执行**：否（高风险，须确认，见 audit/）

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_id(必填) | |
| reason(必填) | |
| scope(隔离范围，必填) | |
| evidence(必填) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| action_status(pending_confirmation/confirmed/rejected) | |
| confirmation(确认人/时间/依据) | |
| audit_record_id | |

## 失败/超时返回语义
未确认 status=pending_confirmation 且不执行；超时按拒绝处理；不提供底层隔离实现。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
