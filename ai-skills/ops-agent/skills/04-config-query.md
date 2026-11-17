<!-- @ticket 10：设备运维 Skill 规格——配置查询 -->
# 配置查询（工单 10）

- **名称**：配置查询
- **用途描述**：查询设备/服务的当前配置与最近变更记录（只读）
- **适用场景**：变更影响评估、确认配置是否符合预期
- **风险等级**：低
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_id(必填) | |
| service(可选) | |
| config_keys(配置项列表，可选) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| configs[{key,value,version}] | |
| last_changed(最近变更时间/人) | |
| diff(与基线差异，若有) | |

## 失败/超时返回语义
配置不可读返回 status=failed 并说明缺口；无变更记录时明确标注而非臆断。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
