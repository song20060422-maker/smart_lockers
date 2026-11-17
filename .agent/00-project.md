---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_8574eeede27c11f19919525400638852
    ReservedCode1: nF3AhavtcWkktboyURDWXTYYHPbzcqSsXlguXGoX85bt5+pzztLTN4yl9fXplcylxb6d8zmUShFZCX/U5FttDt6sjo4EfL8Ut7dddKp1ZI1rG9YkcMkKgJdnObuoN2I9QD5zwM+lQoxMwRxdK2RwsHkYahwfISC8bgs6EMUljFQRoVKeyvsshUGsVMQ=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_8574eeede27c11f19919525400638852
    ReservedCode2: nF3AhavtcWkktboyURDWXTYYHPbzcqSsXlguXGoX85bt5+pzztLTN4yl9fXplcylxb6d8zmUShFZCX/U5FttDt6sjo4EfL8Ut7dddKp1ZI1rG9YkcMkKgJdnObuoN2I9QD5zwM+lQoxMwRxdK2RwsHkYahwfISC8bgs6EMUljFQRoVKeyvsshUGsVMQ=
---

# 00 · 项目全局规则

> 适用：所有任务。加载顺序：本文件在 `../AGENT.MD` 之后、其他分规则之前。

## 1. 技术栈边界（不得混用）

- 后端：`smart_lockers/` = Go 1.20+ + Gin；接口 RESTful + JSON；关系型库（MySQL 系）+ 缓存（Redis 系）。
- 前端：`smart_lockers_frontend/` = Vue3 + Vite + TypeScript + Pinia + Vue Router + Axios。
- 工单 07 数字人模块：基于 `awesome-digital-human-live2d`，原生 React + TS + Vite + Live2D；作为**独立子模块**集成，依赖与主 Vue3 前端隔离，禁止混装。

## 2. 目录与命名

- 后端采用分层结构建议：`cmd/`、`configs/`、`internal/`（handler / service / dao / model / middleware / pkg）、`docs/`、`scripts/`、`test/`。
- 前端建议：`src/`（views / components / api / stores / router / utils / types）。
- 命名：Go 包/文件小写下划线或短名，导出符号大驼峰；前端组件大驼峰、目录小写中划线。接口路径、字段名全项目统一风格（snake_case 或 camelCase 二选一，全局一致）。

## 3. 工单编号注释（强制）

- 每个工单新增/修改的源文件，关键逻辑处必须有注释标注工单编号，格式建议：`// @ticket 04-物品寄存：<说明>` / `// @ticket 05-订单查看与取出：<说明>`。
- 10 个工单备注均明确要求"代码注释中需要写明工单编号"，缺失视为不合格。

## 4. 通用质量铁律（源自 10 工单共性要求）

1. 单元测试覆盖率 ≥70%（工单 03/04/06 等明确）。
2. 错误码与日志格式全项目统一，禁止各模块自定义格式。
3. 关键接口须可压测并提供压测报告与解读（10 工单产出物共性）。
4. 涉及资金、状态变更、柜门动作：事务保护 + 幂等控制 + 鉴权 + 脱敏 + 审计，缺一不可。
5. 业务规则（收费、通知渠道、服务入口、指南文案、频控阈值）需可配置、可动态生效。
6. 所有数据结构/接口命名规范，文档清晰，变更可追溯。

## 5. 标准产出物清单（每个业务工单交付时对齐）

《设计说明书》、数据字典、ER 图、状态转换图、接口报文设计说明书、数据库建表及索引脚本、接口文档（请求/响应/错误码）、接口联调与测试报告、压测报告、日志与审计规范、上线说明与回滚预案。

## 6. 可启动性约定

- 后端须保证在 `smart_lockers/` 根目录 `go run` 即可启动；配置集中于 `configs/` 下的配置文件（本地明文），不依赖运行时手动注入环境变量。
- 任何真实密码、密钥、Token 不得写死提交到仓库；配置文件中仅放置本地开发占位值，生产凭据由部署环境或用户自行提供。

## 7. 红线

- 不臆造接口、字段、指标；一切结论来自工单原文或真实代码。
- 不修改 `.git`、不提交 `.env`/私钥等敏感文件（见 `.gitignore`）。
- 一次只推进一个工单，做完即交付并标注工单编号。
*（内容由AI生成，仅供参考）*
