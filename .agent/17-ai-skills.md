---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_8e608b5be27c11f19919525400638852
    ReservedCode1: O26FR7fbw0GmAEdyLFCoXU66dzff7fjVLRQiwPnSreKxQhysYYhd1u8PvY1s+LSUS8zbP2GtZ5txU/QD6EHrk9/+a2puYnFUBW2VfDdSCd0f3ri7x7RW+Nlxent+pQQ6F7+enMwmz02UhLkmccKjU7SltdOS3UjmiAqmoxKZ4i/Y2GfEXMapzcH0pcU=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_8e608b5be27c11f19919525400638852
    ReservedCode2: O26FR7fbw0GmAEdyLFCoXU66dzff7fjVLRQiwPnSreKxQhysYYhd1u8PvY1s+LSUS8zbP2GtZ5txU/QD6EHrk9/+a2puYnFUBW2VfDdSCd0f3ri7x7RW+Nlxent+pQQ6F7+enMwmz02UhLkmccKjU7SltdOS3UjmiAqmoxKZ4i/Y2GfEXMapzcH0pcU=
---

# 17 · AI 能力规范：Skills / Agent / 数字人

> 适用：工单 09（Cursor/Trae 优化 Skills）、工单 10（设备运维智能体）、工单 07（数字人 Coze）。前置：`00-project.md`。

## 1. 工单 09 · Go 项目优化 Skills（Cursor/Trae）

- 目标：把 Go 项目优化经验抽象为 IDE 内可触发的 Skills，沉淀为团队可复用资产。
- 分析维度：性能优化/瓶颈识别/低效模式、并发风险/稳定性、架构职责与边界/演进风险、依赖合理性/工程健康。
- 触发：在 Cursor/Trae 中通过 Skills/Agent 能力稳定可复用地触发；必须基于当前 Go 项目真实上下文（结构、模块关系、变更范围），禁止模板化泛化建议。
- 原则：聚焦工程价值、避免无意义过度优化、区分风险等级（不一刀切）、结论可解释。
- 统一输出结构：①问题类型（性能/并发/架构/依赖）②风险等级/影响程度 ③工程背景 ④可执行优化方向 ⑤清晰结构便于 PR/MR 审查与沉淀。
- 产出：Skills 配置/定义文件（含触发说明、上下文要求、输出模板）、运行文档、可复用维护说明。

## 2. 工单 10 · 设备运维智能体 + Skills

### 智能体闭环
`识别（对话/告警/工单事件）→ 计划（基于设备/时间窗/指标/变更范围选 Skill 并编排）→ 调用 Skills → 统一汇总输出`。
- 双触发：对话触发（如"检查设备 A CPU 高原因并给建议"）+ 事件触发（监控告警、巡检计划、工单状态变化）。
- 上下文感知：必须基于真实设备标识、指标/告警、时间窗、相关服务、历史事件；信息不足须明确缺口与补充建议，禁止模板化结论。

### Skills 规范（每个 Skill 必备字段）
| 字段 | 说明 |
|---|---|
| 名称 / 用途 / 适用场景 | 清晰定义 |
| 输入字段（参数）定义 | 结构化 |
| 输出字段（结构）定义 | 结构化 |
| 风险等级 | 低 / 中 / 高 |
| 是否允许自动执行 | 是 / 否 |
| 失败/超时返回语义 | 成功 / 失败 / 部分成功 |

### Skills 清单（≥10，覆盖四类）
- 查询类：状态查询、指标查询、日志查询、配置查询。
- 诊断类：告警归因、关联分析、健康检查。
- 报告类：巡检报告、事件总结。
- 风险动作类：重启/隔离/限流/回滚的确认占位 Skill（明确"确认/审计"语义，仅提示不落地实现）。

### 统一输出结构（必备字段）
①问题类型（告警/性能/网络/存储/服务/配置/变更影响）②风险等级（低/中/高）③工程/运维背景（基于上下文）④可执行下一步动作（建议/排查/需补充信息）⑤执行记录摘要（调用了哪些 Skill、各自结论）。
- 高风险动作：无确认不得执行，且需审计记录可追溯。
- Skills 可新增/下线/升级，不破坏整体规范。

## 3. 工单 07 · 数字人 Coze 智能体（前端独立模块）

- 平台：Coze（插件 Plugin / 工作流 Workflow / 知识库 Knowledge Base）。
- 渲染：`awesome-digital-human-live2d`（React + TS + Vite + Live2D + Redux Toolkit），独立子模块，依赖与主 Vue3 前端隔离。
- 通信：WebSocket（对话/状态/指令）+ RESTful（配置/知识库/状态）；消息结构 `{role, content, type: text|voice|command, timestamp}`。
- 语音：SpeechRecognition 输入、SpeechSynthesis 或第三方 TTS（Azure/OpenAI）输出，驱动 Live2D 表情/动作。
- 工作流：问候→识别意图→调用知识库→响应输出，多轮上下文记忆；知识库条目设标签/权重。
- 指标：响应<2s、知识库准确率≥90%、帧率≥55fps、音画同步误差<0.3s、断线重连100%；token 鉴权、跨域控制、调用限流。
- 产出：Coze 配置、Live2D 配置与渲染代码、WebSocket 模块与回调、联调脚本、运行文档、性能报告。

## 4. 通用要求

- 所有 AI 能力代码/配置注释标注工单编号：`// @ticket 10-设备运维智能体：<说明>`。
- 配置随项目演进可维护；结论基于真实上下文，可解释、可复查、可沉淀。
*（内容由AI生成，仅供参考）*
