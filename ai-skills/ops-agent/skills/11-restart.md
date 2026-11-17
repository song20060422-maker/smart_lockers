<!-- @ticket 10：设备运维 Skill 规格——设备重启（占位·需确认） -->
# 设备重启（占位·需确认）（工单 10）

- **名称**：设备重启（占位·需确认）
- **用途描述**：高风险动作：重启设备。仅定义确认/审计语义与占位，不提供执行实现
- **适用场景**：告警处置中判断需要重启时；工单要求高风险动作必须确认
- **风险等级**：高
- **是否允许自动执行**：否（高风险，须确认，见 audit/）

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_id(必填) | |
| reason(重启原因，必填) | |
| evidence(依据，必填) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| action_status(pending_confirmation/confirmed/rejected) | |
| confirmation(确认人/时间/依据) | |
| audit_record_id | |

## 失败/超时返回语义
未获确认时 status=pending_confirmation 且绝不执行；确认超时视为拒绝；执行类实现不在本 Skill 范围。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
