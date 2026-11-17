<!-- @ticket 10：设备运维 Skill 规格——配置回滚（占位·需确认） -->
# 配置回滚（占位·需确认）（工单 10）

- **名称**：配置回滚（占位·需确认）
- **用途描述**：高风险动作：回滚配置/版本。仅定义确认/审计语义与占位
- **适用场景**：变更引发问题时；与 08 变更影响评估联动
- **风险等级**：高
- **是否允许自动执行**：否（高风险，须确认，见 audit/）

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| target(回滚目标：配置项/版本，必填) | |
| reason(必填) | |
| evidence(必填) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| action_status(pending_confirmation/confirmed/rejected) | |
| confirmation(确认人/时间/依据) | |
| audit_record_id | |
| rollback_suggestion(来自08) | |

## 失败/超时返回语义
未确认不执行；确认超时视为拒绝；回滚执行实现不在本 Skill 范围，仅保留确认与审计占位。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
