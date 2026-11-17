---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_80c26d2be27c11f18c5752540024e231
    ReservedCode1: Glwo1+SWAvLBtN8upgyz3z/GsrvzykxeziX06o37QDZDyihCLhEjHBfp0VOEHENw0QFPCs1bYA6chRC8QdN8T7S2qfF/a0vCohlZQlExEIswgXNhTGAhUIIHr7aLqw+bpjYOgfwO9YWwEw0u+YE+d5ZOZCb7pxlqJOIULTNdJKo9e4yIMmixDOgp4rA=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_80c26d2be27c11f18c5752540024e231
    ReservedCode2: Glwo1+SWAvLBtN8upgyz3z/GsrvzykxeziX06o37QDZDyihCLhEjHBfp0VOEHENw0QFPCs1bYA6chRC8QdN8T7S2qfF/a0vCohlZQlExEIswgXNhTGAhUIIHr7aLqw+bpjYOgfwO9YWwEw0u+YE+d5ZOZCb7pxlqJOIULTNdJKo9e4yIMmixDOgp4rA=
---

# AGENT.MD — 智能快寄柜项目 AI 协作入口

> 本文件是本项目所有 AI 编码助手（Cursor / Trae / Claude Code / Marvis 等）的**顶层必读规则**。
> 任何 AI 在本仓库产出代码、文档或配置前，必须先读本文件，并按 `.agent/` 下的分规则执行。
> 规则优先级：`AGENT.MD` > `.agent/*.md` > 个人习惯。冲突时以更高优先级为准。

---

## 1. 项目身份

| 项 | 内容 |
|---|---|
| 项目名称 | 智能快寄柜系统（Smart Lockers） |
| 需求来源 | `D:\Users\Lenovo\Desktop\智能快寄柜` 下 10 份任务工单 V1.x（2026-01，创建人：邹世军，出品：八维文化与产业研究院） |
| 代码仓库 | https://github.com/song20060422-maker/smart_lockers.git |
| 全局目标 | 按 10 个工单逐项实现功能：先补齐工程配置（本文件、README、.gitignore、.agent），再按工单开发后端与前端 |
| 原型地址 | https://modao.cc/proto/3fxWzFGJsp37iCk5fS2h/sharing?view_mode=read_only |

## 2. 技术栈与目录结构

```
smart_lockers/                          <- 本文件所在的最外层（配置与规则层）
├── AGENT.MD                            <- AI 协作入口（本文件）
├── README.MD                           <- 项目说明（人读）
├── .gitignore                          <- 忽略规则（Go + Node + IDE）
├── .agent/                             <- 分规则集（AI 读，见第 5 节索引）
├── smart_lockers/                      <- 后端：Go + Gin（Gin 框架 REST API）【空仓库，待搭建】
└── smart_lockers_frontend/             <- 前端：Vue3 + Vite【空目录，待搭建】
```

- **后端**：Go 1.20+，Gin 框架；关系型库 + 缓存层；接口RESTful + JSON。
- **前端**：Vue3 + Vite（主站/小程序管理端风格的 Web 前端）。
- **特例**：工单 07 数字人模块基于 `awesome-digital-human-live2d`，其原生技术栈为 React + TypeScript + Vite + Live2D；作为独立子模块集成，不得与主前端 Vue3 技术栈混用依赖。
- 仓库现状：`smart_lockers/` 仅有 `.git`（无任何提交、未配置 remote）；`smart_lockers_frontend/` 为空目录。两个子项目均为**从零搭建**，AI 不得假设已有代码结构。

## 3. 十工单一览（必读，写每行代码都要能对应到工单）

| 编号 | 工单名称 | 核心能力 | 工时 | 关键指标 / 约束 |
|---|---|---|---|---|
| 01 | 消息通知 | 统一通知告警中心：站内信/短信/邮件/回调多渠道；模板+变量渲染；去重频控；失败重试升级；审计+监控看板；8 类场景 | 3 人日 | 投递分位时延达标、失败率可控、模板变量安全转义 |
| 02 | 快寄员支付宝支付 | 余额查询、余额明细、支付宝充值、充值记录查询四类能力；资金可追溯 | 3 人日 | 任意余额变动必有明细；充值与入账关系清晰；可复用可扩展 |
| 03 | 首页 | 定位当前城市(GPS/IP 回退)、城市选择(拼音/中文搜索)、附近寄存点搜索(距离排序分页)、寄存指南(版本/多语言) | 3 人日 | 定位≤200ms；城市搜索/寄存点查询 P95≤300ms；千级并发；定位数据不持久化 |
| 04 | 物品寄存 | 柜子选择、时间配置、优惠券+押金规则、收费规则计算、支付完成落单、支付异常释放/换柜 | 3 人日 | 柜子查询 P95≤200ms；收费计算 P95≤250ms；全流程事务保护；支付幂等 |
| 05 | 订单查看与取出 | 当前/历史订单列表、订单详情(寄存地点+位置)、临时开柜、超时补交、已结束详情、进行中聚合 | 5.5 人日 | 游标分页、幂等键+短期锁、读写分离+热点缓存精准失效、审计轨迹 |
| 06 | 个人中心 | 用户订单查询、头像上传管理、我的服务聚合、优惠卡券(可用/过期)、退出登录 | 2 人日 | 订单查询 P95≤300ms；头像上传 P95≤500ms；上传文件校验；会话即时失效 |
| 07 | 数字人 Coze 智能体 | Coze 插件/工作流/知识库 + Live2D 数字人渲染 + WebSocket/语音联动 | 5.5 人日 | 响应<2s；知识库准确率≥90%；帧率≥55fps；同步误差<0.3s；断线重连100% |
| 08 | Go 后端性能调优工具链 | pprof/trace/expvar/runtime-metrics/block profile/gops/pyroscope/火焰图/Prometheus | 中 | 采样对业务性能影响<3%；分析接口需鉴权限访问 |
| 09 | Cursor/Trae 项目优化 Skills | 将 Go 优化经验沉淀为 IDE 内可触发的 Skills（性能/并发/架构/依赖） | 1 人日 | 基于真实上下文、结论可解释、分级、结构统一、可复用 |
| 10 | 结合 Skills 开发设备运维智能体 | 运维智能体(识别→计划→调用 Skills→汇总) + ≥10 个运维 Skills + 风险动作确认审计 | 5.5 人日 | 对话+事件双触发；风险分级；高风险动作无确认不得执行；输出可追溯 |

## 4. 全局开发铁律（10 个工单的共性硬约束，违反即返工）

1. **工单编号注释**：每个工单新增/修改的代码文件，关键逻辑处注释必须写明工单编号（如 `// @ticket 04-物品寄存：支付幂等校验`）。10 个工单备注均强制要求。
2. **单测覆盖率 ≥70%**：核心接口/流程必须配套单元测试，工单 01/03/04/05/06 明确要求。
3. **错误码与日志格式统一**：全项目统一错误码表与日志规范，禁止散落自定义格式。
4. **接口压力测试**：每个接口需提供压测能力与压测报告解读（10 工单产出物共性要求）。
5. **事务 / 幂等 / 鉴权 / 脱敏 / 审计**：涉及资金、状态变更、柜门动作的接口，必须有事务保护、幂等控制、鉴权、敏感信息脱敏与审计记录。
6. **标准产出物**：每个工单需交付《设计说明书》、数据字典、ER 图、状态转换图、接口报文设计说明书、建表/索引脚本、接口文档、联调/压测报告、上线说明与回滚预案。
7. **可配置化**：收费规则、渠道策略、服务入口、指南文案、频控阈值等需可配置、可动态生效。
8. **引用真实来源**：结论必须基于实际代码/数据/工单原文，禁止臆造接口、字段或指标。

## 5. 分规则索引（进入对应开发任务前必读对应文件）

| 文件 | 适用范围 |
|---|---|
| `.agent/README.md` | 分规则加载方式与维护约定 |
| `.agent/00-project.md` | 项目全局规则：技术栈、目录、命名、通用质量铁律 |
| `.agent/10-backend.md` | 后端 Go + Gin 分层、中间件、事务、幂等、错误码、日志 |
| `.agent/11-frontend.md` | 前端 Vue3 + Vite 目录、组件、状态管理、请求层、数字人特例 |
| `.agent/12-database.md` | 表设计、索引、字段规范、归档分区、缓存策略 |
| `.agent/13-api.md` | RESTful 接口、统一报文、错误码表、分页、限流 |
| `.agent/14-testing.md` | 单测、集成/联调、压测报告、验收口径 |
| `.agent/15-security.md` | 鉴权、脱敏、幂等防重放、审计、风险动作确认 |
| `.agent/16-tickets.md` | 10 工单 → 模块/表/接口 映射与业务细节速查 |
| `.agent/17-ai-skills.md` | 工单 09/10 的 Skills/Agent 能力定义规范（含工单 07 数字人） |

## 6. 执行约定

- **先配置后功能**：工程配置文件优先完成；功能代码按工单编号顺序推进，一次只做一个工单，做完即交付并标注工单编号。
- **本地可启动**：后端须保证在项目根目录直接 `go run` 即可启动，配置集中在 config 文件，不依赖手动注入环境变量。
- **不越界**：不在 AI 规则/文档中写入任何真实密码、密钥、Token；凭据由用户在配置中自行填写。
- **变更可追溯**：所有规则、代码变更需在提交信息或注释中标注来源工单编号与变更原因。
*（内容由AI生成，仅供参考）*
