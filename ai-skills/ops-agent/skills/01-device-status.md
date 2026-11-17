<!-- @ticket 10：设备运维 Skill 规格——设备状态查询 -->
# 设备状态查询（工单 10）

- **名称**：设备状态查询
- **用途描述**：查询指定设备/批次的在线状态、运行状态与基础信息
- **适用场景**：运维人员想了解某设备当前是否在线、运行状态；作为其他诊断的前置步骤
- **风险等级**：低
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_id(设备标识，必填) | |
| scope(single/batch，默认 single) | |
| filters(区域/类型，可选) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| devices[{device_id,online,status,last_seen}] | |
| summary(在线数/离线数) | |
| context_window(查询时间窗口) | |

## 失败/超时返回语义
失败返回 status=failed 并附错误原因；部分设备不可达时 status=partial_success 并列出失败设备。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
