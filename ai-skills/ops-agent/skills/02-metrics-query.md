<!-- @ticket 10：设备运维 Skill 规格——指标查询 -->
# 指标查询（工单 10）

- **名称**：指标查询
- **用途描述**：按设备、指标项、时间窗口拉取 CPU/内存/磁盘/网络等指标并汇总
- **适用场景**：告警归因、健康检查需要指标依据；工单示例“检查设备 A CPU 高”
- **风险等级**：低
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_id(必填) | |
| metrics(指标名列表，如 cpu/mem/disk/net) | |
| time_window(起止时间，必填) | |
| aggregation(avg/max/p95) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| series[{metric,points[],agg_value}] | |
| anomalies(超阈值点列表) | |
| thresholds(采用的阈值) | |

## 失败/超时返回语义
超时返回 status=partial_success 并标明已获取的时间片段；完全失败 status=failed。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
