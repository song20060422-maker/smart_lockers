<!-- @ticket 10：设备运维 Skill 规格——事件总结 -->
# 事件总结（工单 10）

- **名称**：事件总结
- **用途描述**：对一次运维事件/告警处置过程生成可复盘总结
- **适用场景**：事件闭环后的复盘材料、工单处理记录
- **风险等级**：低
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| event_id(必填) | |
| time_window(必填) | |
| related_skills(调用过的 Skills) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| timeline(事件时间线) | |
| root_cause_summary | |
| actions_taken(已执行动作) | |
| lessons(经验沉淀要点) | |

## 失败/超时返回语义
记录不全时标注缺失环节；总结须基于已调用 Skill 的真实结论。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
