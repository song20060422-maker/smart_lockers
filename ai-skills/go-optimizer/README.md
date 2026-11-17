---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_92d2f413e27c11f18c5752540024e231
    ReservedCode1: zVxnI2miLUXusyl5KsDbEE5RVF1Ifm7acqtox747jOXhC5cPJk2N8E8aZqsLEt1q/wI2YFuuMNMLbd7fOQ7s0KrL42xP8fIT1gOj6fXWPKT2uogXqeujJIxLpIt6PyuH6Ks16fIRV7SBcoJbAjAKaE6iLHvH/x754a0pAaSWBwsYS0q7sOB919TK+PU=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_92d2f413e27c11f18c5752540024e231
    ReservedCode2: zVxnI2miLUXusyl5KsDbEE5RVF1Ifm7acqtox747jOXhC5cPJk2N8E8aZqsLEt1q/wI2YFuuMNMLbd7fOQ7s0KrL42xP8fIT1gOj6fXWPKT2uogXqeujJIxLpIt6PyuH6Ks16fIRV7SBcoJbAjAKaE6iLHvH/x754a0pAaSWBwsYS0q7sOB919TK+PU=
---

<!-- @ticket 09：运行/使用文档——Cursor / Trae 配置项目优化 Skills 能力 -->

# Go 项目优化 Skills（工单 09）

将 Go 优化经验沉淀为 IDE 内可触发的标准化 Skills 能力。**仅定义能力与工程化配置，不含具体算法/工具实现。**

## 目录结构
```
ai-skills/go-optimizer/
├── SKILL.md                 主 Skill（触发/上下文/原则/统一输出）
├── skills/                  四个分析维度子 Skill
│   ├── performance.md       性能
│   ├── concurrency.md       并发与稳定性
│   ├── architecture.md      架构与可维护性
│   └── dependency.md        依赖与工程健康度
├── output-template.md       统一输出结构模板
├── .cursor/rules/*.mdc      Cursor 触发配置
├── .trae/rules/*.md         Trae 触发配置
├── MAINTENANCE.md           维护与迭代说明
└── README.md                本文档
```

## 接入步骤
1. **Cursor**：将 `.cursor/rules/go-optimizer.mdc` 复制到目标 Go 项目的 `.cursor/rules/` 目录；`skills/`、`SKILL.md` 等一并放入 `ai-skills/go-optimizer/`。
2. **Trae**：将 `.trae/rules/go-optimizer.md` 内容加入 Trae 项目规则。
3. 在对话中输入触发语义（见 SKILL.md frontmatter triggers），即可触发优化分析。

## 覆盖维度（工单要求）
性能优化、潜在瓶颈识别、低效实现提示；并发风险与稳定性问题；模块职责边界与演进风险；依赖合理性与工程健康度。

## 统一输出
问题类型 + 风险等级 + 工程背景 + 可执行建议 + 依据，结构清晰可复查（详见 output-template.md）。

## 验收对照
| 验收项 | 实现方式 |
|---|---|
| Skills 可被识别并触发 | .cursor/rules 与 .trae/rules 注册 |
| 覆盖约定能力范围 | 四个子 Skill 覆盖四维度 |
| 输出结构统一、工程可用 | output-template.md 强制统一结构 |
| 不依赖个人经验即可复用 | SKILL.md + 子 Skill 显式化经验 |
| 可随项目持续维护 | MAINTENANCE.md 扩展/升级/下线规范 |
*（内容由AI生成，仅供参考）*
