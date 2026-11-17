<!-- @ticket 10：设备运维 Skill 规格——变更影响评估（轻量） -->
# 变更影响评估（轻量）（工单 10）

- **名称**：变更影响评估（轻量）
- **用途描述**：结合变更范围评估潜在影响点与回滚建议（只提示，不做实现）
- **适用场景**：设备/配置变更前后的风险评估；工单要求仅做提示
- **风险等级**：中
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| change_scope(变更范围：设备/配置/服务，必填) | |
| device_ids(可选) | |
| time_window(可选) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| impact_points(潜在影响点，带依据) | |
| risk_level | |
| rollback_suggestion(提示性，需确认) | |
| uncertainties(不确定项) | |

## 失败/超时返回语义
变更信息不足 status=partial_success 并列出缺口；禁止无依据的过度影响推断。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
