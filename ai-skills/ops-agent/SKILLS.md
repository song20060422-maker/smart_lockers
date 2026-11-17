---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_954cebcbe27c11f19919525400638852
    ReservedCode1: mcDFl8D/VRC85rCFirEx2FjtAXj8jgJ33j8qsFQ784FuN3HdS4J4OB7XcC7AAl136IW6FTpViYZqF3IMsvxEpKI62fenkb0UEPLNKe3fP8TRUwZkJKK5fR2FDjlEcb+LPIPMTa0zCJ4rpdM6dTjOLji5suZCrI7gbtLOyEQoeb2PzT5OMLm6QuzSWHM=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_954cebcbe27c11f19919525400638852
    ReservedCode2: mcDFl8D/VRC85rCFirEx2FjtAXj8jgJ33j8qsFQ784FuN3HdS4J4OB7XcC7AAl136IW6FTpViYZqF3IMsvxEpKI62fenkb0UEPLNKe3fP8TRUwZkJKK5fR2FDjlEcb+LPIPMTa0zCJ4rpdM6dTjOLji5suZCrI7gbtLOyEQoeb2PzT5OMLm6QuzSWHM=
---

<!-- @ticket 10：Skills 清单与规格索引（13 个，覆盖查询/诊断/报告/风险动作四类，工单要求≥10） -->

# 设备运维 Skills 清单（工单 10）

每个 Skill 均定义：名称、用途、适用场景、输入字段、输出字段、风险等级、是否自动执行、失败/超时返回语义。

## 查询类（低风险，可自动执行）
| 编号 | 名称 | 风险 | 自动执行 | 规格 |
|---|---|---|---|---|
| 01 | 设备状态查询 | 低 | 是 | [skills/01-device-status.md](./skills/01-device-status.md) |
| 02 | 指标查询 | 低 | 是 | [skills/02-metrics-query.md](./skills/02-metrics-query.md) |
| 03 | 日志查询 | 低 | 是 | [skills/03-log-query.md](./skills/03-log-query.md) |
| 04 | 配置查询 | 低 | 是 | [skills/04-config-query.md](./skills/04-config-query.md) |

## 诊断类（低~中风险）
| 编号 | 名称 | 风险 | 自动执行 | 规格 |
|---|---|---|---|---|
| 05 | 告警归因分析 | 中 | 是 | [skills/05-alert-attribution.md](./skills/05-alert-attribution.md) |
| 06 | 关联分析 | 中 | 是 | [skills/06-correlation.md](./skills/06-correlation.md) |
| 07 | 健康检查 | 低 | 是 | [skills/07-health-check.md](./skills/07-health-check.md) |
| 08 | 变更影响评估（轻量，仅提示） | 中 | 是 | [skills/08-change-impact.md](./skills/08-change-impact.md) |

## 报告类（低风险）
| 编号 | 名称 | 风险 | 自动执行 | 规格 |
|---|---|---|---|---|
| 09 | 巡检报告生成 | 低 | 是 | [skills/09-inspection-report.md](./skills/09-inspection-report.md) |
| 10 | 事件总结 | 低 | 是 | [skills/10-event-summary.md](./skills/10-event-summary.md) |

## 风险动作类（高风险，禁止自动执行，需确认+审计；均为占位技能）
| 编号 | 名称 | 风险 | 自动执行 | 规格 |
|---|---|---|---|---|
| 11 | 设备重启 | 高 | 否 | [skills/11-restart.md](./skills/11-restart.md) |
| 12 | 设备隔离 | 高 | 否 | [skills/12-isolate.md](./skills/12-isolate.md) |
| 13 | 配置回滚 | 高 | 否 | [skills/13-rollback.md](./skills/13-rollback.md) |

> 风险动作类仅定义“确认/审计”语义与占位，不提供底层执行实现（工单明确不要求 SSH/SNMP/厂商 API 适配细节）。
*（内容由AI生成，仅供参考）*
