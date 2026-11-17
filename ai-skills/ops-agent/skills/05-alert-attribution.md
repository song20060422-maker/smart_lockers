<!-- @ticket 10：设备运维 Skill 规格——告警归因分析 -->
# 告警归因分析（工单 10）

- **名称**：告警归因分析
- **用途描述**：对告警分级、关联相关指标/日志，给出初步归因与建议动作
- **适用场景**：监控告警触发后的第一步分析；输出问题类型与风险等级
- **风险等级**：中
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| alert_id / alert_detail(告警详情，必填) | |
| device_id(必填) | |
| time_window(必填) | |
| related_metrics(可选) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| problem_type(告警类型) | |
| risk_level(低/中/高) | |
| root_cause_candidates(带置信度) | |
| evidence(依据现象) | |
| suggested_actions | |

## 失败/超时返回语义
信息不足时输出 status=partial_success 并列出缺失上下文；不臆断根因。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
