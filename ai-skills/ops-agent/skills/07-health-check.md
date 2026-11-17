<!-- @ticket 10：设备运维 Skill 规格——健康检查 -->
# 健康检查（工单 10）

- **名称**：健康检查
- **用途描述**：对单设备/批量设备执行健康检查并给出状态结论
- **适用场景**：日常巡检、故障前健康评估
- **风险等级**：低
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_ids(必填) | |
| checks(检查项：连通性/资源/服务，默认全部) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| results[{device_id,checks[],health(healthy/warning/critical)}] | |
| summary(各状态设备数) | |

## 失败/超时返回语义
部分设备不可达 status=partial_success 并列出；单设备失败不改变整体状态汇总。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
