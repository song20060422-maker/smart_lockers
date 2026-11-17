<!-- @ticket 10：设备运维 Skill 规格——关联分析 -->
# 关联分析（工单 10）

- **名称**：关联分析
- **用途描述**：跨设备/服务/时间关联告警与变更，发现共性根因
- **适用场景**：多设备同时异常、疑似共性变更引起
- **风险等级**：中
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_ids(设备列表) | |
| signals(告警/指标异常列表) | |
| time_window(必填) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| correlations[{signals,time_overlap,confidence}] | |
| common_factors(共性因素) | |
| hypotheses(假设列表) | |

## 失败/超时返回语义
数据不足 status=partial_success；无关联时明确返回“未发现显著关联”而非空泛结论。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
