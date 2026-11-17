<!-- @ticket 10：设备运维 Skill 规格——日志查询 -->
# 日志查询（工单 10）

- **名称**：日志查询
- **用途描述**：按设备、服务、关键字、时间窗口检索并汇总日志
- **适用场景**：排查错误、关联告警时间线、复盘事件
- **风险等级**：低
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_id(必填) | |
| service(服务名，可选) | |
| keywords(关键字，可选) | |
| time_window(必填) | |
| limit(条数上限) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| logs[{timestamp,level,message}] | |
| matched_count | |
| highlights(错误/异常级别条目) | |

## 失败/超时返回语义
无匹配返回 status=success + 空列表；查询超时 status=partial_success 并注明已覆盖时间范围。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
