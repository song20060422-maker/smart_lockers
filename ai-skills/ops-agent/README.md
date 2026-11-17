---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_97dcbf2de27c11f19919525400638852
    ReservedCode1: JIgOM6eLZeY9yYjL4VrZby7zda9i/QbmuFfVL7oBBWcYbPP1MPBEaqg5icNV0xbvo5OsEzCgIEfEVwMp56dpk4YC4M9LejvqzpReRwVyHPNTTdsxSA/7n4e6l9RyCd1GmeQjjS4EuTKzUVHSqTGHhzE4sWjuGKVWFa3b6zPNPOwwMtaEg+9qlWcfgOE=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_97dcbf2de27c11f19919525400638852
    ReservedCode2: JIgOM6eLZeY9yYjL4VrZby7zda9i/QbmuFfVL7oBBWcYbPP1MPBEaqg5icNV0xbvo5OsEzCgIEfEVwMp56dpk4YC4M9LejvqzpReRwVyHPNTTdsxSA/7n4e6l9RyCd1GmeQjjS4EuTKzUVHSqTGHhzE4sWjuGKVWFa3b6zPNPOwwMtaEg+9qlWcfgOE=
---

<!-- @ticket 10：设备运维智能体运行/使用文档 -->

# 设备运维智能体（工单 10）

运维智能体 + Skills 能力定义与工程化配置。**仅定义能力、流程编排与配置，不含 SSH/SNMP/厂商 API 底层实现。**

## 目录结构
```
ai-skills/ops-agent/
├── AGENT.md                        智能体定义（识别→计划→调用→汇总闭环）
├── SKILLS.md                       Skills 清单与规格索引（13 个）
├── skills/01..13-*.md              每个 Skill 的完整规格
├── output-schema.md                统一输出结构
├── triggers/dialogue.md            对话触发配置
├── triggers/event.md               事件触发配置
├── audit/confirmation-and-audit.md 风险确认与审计机制
├── MAINTENANCE.md                  维护与演进说明
└── README.md                       本文档
```

## 工作流程
对话/事件触发 → 识别上下文 → 计划并编排 Skills → 调用 Skills → 汇总输出统一结构结论。
高风险动作（重启/隔离/回滚）进入确认流程，无确认不执行，有记录可追溯。

## Skills 覆盖（13 个，工单要求≥10）
- 查询类：设备状态、指标、日志、配置
- 诊断类：告警归因、关联分析、健康检查、变更影响评估
- 报告类：巡检报告、事件总结
- 风险动作类：重启、隔离、回滚（均需确认+审计的占位技能）

## 使用示例
1. 对话触发：“检查设备 0203445006 CPU 高的原因并给建议” → 编排 02+03+05 → 统一输出。
2. 事件触发：CPU 告警 → 05 归因 → 若建议重启，仅产出 `pending_confirmation`，等待确认并记录审计。

## 验收对照
| 验收项 | 实现方式 |
|---|---|
| Skills 可被识别并触发 | SKILLS.md + 各 skills 规格 + triggers/ |
| 双触发（对话+事件） | triggers/dialogue.md + triggers/event.md |
| 覆盖查询/诊断/巡检/风险控制/输出沉淀 | 13 个 Skills 四类覆盖 |
| 输出结构统一、工程可用 | output-schema.md 强制统一字段 |
| 不依赖个人经验即可复用 | 显式 Skill 规格 + 闭环流程 |
| 高风险动作确认与审计 | audit/confirmation-and-audit.md（无确认不执行、可追溯） |
| 可持续维护调整 | MAINTENANCE.md 新增/下线/升级规范 |
*（内容由AI生成，仅供参考）*
