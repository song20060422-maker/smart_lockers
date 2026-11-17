<!-- @ticket 10：设备运维 Skill 规格——巡检报告生成 -->
# 巡检报告生成（工单 10）

- **名称**：巡检报告生成
- **用途描述**：汇总巡检/健康检查结果为结构化报告
- **适用场景**：周期性巡检、运维质量评估材料
- **风险等级**：低
- **是否允许自动执行**：是

## 输入字段（参数）
| 字段 | 说明 |
|---|---|
| device_scope(设备范围) | |
| time_window(必填) | |
| include(包含的检查项) | |

## 输出字段（结构）
| 字段 | 说明 |
|---|---|
| report(结构化巡检报告) | |
| overview(总体健康度) | |
| findings(问题清单，分级) | |
| recommendations | |

## 失败/超时返回语义
数据源缺失时报告标注“数据不完整”章节；不伪造指标。

> 规格定义遵循 [AGENT.md](../AGENT.md) 闭环与 [output-schema.md](../output-schema.md) 统一输出；风险动作确认审计见 [audit/confirmation-and-audit.md](../audit/confirmation-and-audit.md)。
