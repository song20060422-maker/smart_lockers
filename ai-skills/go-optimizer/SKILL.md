---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_91f2b811e27c11f19919525400638852
    ReservedCode1: /LGI6cIDskBRtZGbfSu2Bk6JaJejTa6Qb9ACDsVtkdmMarvji04aRJOgptDHlU/bqionEvlxjGdufwa3Pwoiw0IppQ7SPnfKB2aSuPWXGZd78PNgbLgyiPyyni0Id2qXvyvmb8eT88qvaac57EV8praD+8lGcpTAe28WnKUgZkIC2IFf5M+hcKjrDVw=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_91f2b811e27c11f19919525400638852
    ReservedCode2: /LGI6cIDskBRtZGbfSu2Bk6JaJejTa6Qb9ACDsVtkdmMarvji04aRJOgptDHlU/bqionEvlxjGdufwa3Pwoiw0IppQ7SPnfKB2aSuPWXGZd78PNgbLgyiPyyni0Id2qXvyvmb8eT88qvaac57EV8praD+8lGcpTAe28WnKUgZkIC2IFf5M+hcKjrDVw=
---

---
name: go-project-optimizer
description: "@ticket 09：Go 项目优化分析主 Skill。在 Cursor/Trae 中触发，基于当前 Go 项目真实上下文，输出统一结构的优化结论。覆盖性能、并发稳定性、架构可维护性、依赖工程健康度。"
version: 1.0.0
ticket: "全栈开发-智能快寄柜项目-Cursor或者Trae配置项目优化Skills能力工单"
triggers:
  - "分析当前 Go 项目的优化点"
  - "审查这个项目/模块的性能与并发风险"
  - "对本次变更做工程优化评估"
---

<!-- @ticket 09：Go 项目优化 Skills 主定义（能力定义与工程化配置，非具体算法实现） -->

# Go 项目优化分析 Skill（工单 09）

## 一、能力目标
将常见 Go 项目优化经验抽象为 IDE 内可触发的标准化能力，使优化分析成为**结构化、可复用、团队资产化**的工程能力，而非依赖个人的一次性判断。

## 二、触发方式（Cursor / Trae）
- 在 Cursor：本 Skill 通过 `.cursor/rules/go-optimizer.mdc` 注册，输入上方 triggers 任一语义即可触发。
- 在 Trae：通过 `.trae/rules/go-optimizer.md` 注册。
- 触发方式稳定、可复用，不绑定具体实现。

## 三、上下文感知（强制）
分析必须基于当前 Go 项目真实上下文，**禁止脱离上下文的泛化/模板化建议**：
1. 识别项目结构（go.mod 模块名、目录分层、main 包位置）。
2. 识别模块关系与调用边界（internal/pkg、接口与实现）。
3. 识别变更范围（git diff / 当前打开文件 / 指定路径）。
4. 结论须引用具体文件、函数或代码片段作为依据。

## 四、分析维度（调用子 Skill）
| 子 Skill | 维度 | 关注点 |
|---|---|---|
| [performance.md](./skills/performance.md) | 性能 | 性能瓶颈、低效实现模式 |
| [concurrency.md](./skills/concurrency.md) | 并发与稳定性 | 并发风险、稳定性隐患 |
| [architecture.md](./skills/architecture.md) | 架构与可维护性 | 模块职责边界、演进风险 |
| [dependency.md](./skills/dependency.md) | 依赖与工程健康度 | 依赖合理性、工程结构健康 |

## 五、优化分析原则
- **聚焦工程价值**：避免无意义或过度优化；无充分依据不出“升级配置/重构”类空泛建议。
- **区分风险等级**：高/中/低分级，避免“一刀切”。
- **可解释**：每条结论说明依据来自哪段上下文。
- **信息不足时**：明确指出缺口与需要补充的信息，不臆断。

## 六、统一输出结构（必须）
所有优化分析输出必须遵循 [output-template.md](./output-template.md)，逐条包含：
1. 问题类型（性能 / 并发 / 架构 / 依赖 / 工程健康）
2. 风险等级或影响程度（高 / 中 / 低）
3. 工程背景（基于上下文的简述，含文件/函数依据）
4. 可执行的优化方向或建议
5. 结构清晰、可复查、可沉淀（可用于 PR/MR 审查、优化记录、团队质量评估）

## 七、复用与维护
见 [MAINTENANCE.md](./MAINTENANCE.md)：子 Skill 可新增/下线/升级，统一输出结构不变。
*（内容由AI生成，仅供参考）*
