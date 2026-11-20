---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_88d6a8ede27c11f18c5752540024e231
    ReservedCode1: Ljkhxas0SVHpQ++HU4njiFGU39nJJseiqNFmf7aQqAESEz1UnJZBsO/6zoSobUHBt4gsIsBol5q9HVAh2hDcLuRBMTddxyrphfzHx5qQP7Ok4V6pPsOaDlkWjla8/EcCaRKAP2Wj86AnW/lrMuEhYVqleOlP/cu1ZN6C2AO1NqG/Exl/+OFBiFDVrbs=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_88d6a8ede27c11f18c5752540024e231
    ReservedCode2: Ljkhxas0SVHpQ++HU4njiFGU39nJJseiqNFmf7aQqAESEz1UnJZBsO/6zoSobUHBt4gsIsBol5q9HVAh2hDcLuRBMTddxyrphfzHx5qQP7Ok4V6pPsOaDlkWjla8/EcCaRKAP2Wj86AnW/lrMuEhYVqleOlP/cu1ZN6C2AO1NqG/Exl/+OFBiFDVrbs=
---

# 03 · 数据库与缓存规范

> 适用：所有涉及数据持久化的工单。前置：`00-project.md`。

## 1. 建表通用规范

- 主键统一 `id`（bigint，自增或雪花）；每张表含 `created_at`、`updated_at`，逻辑删除用 `deleted_at`（索引需注意）。
- 字段类型精确匹配业务：金额用 `decimal`/整数分（禁用 float），时间用 datetime，状态用 smallint + 状态机约束。
- 表名/字段名风格全库统一（建议 snake_case），注释完整，中文语义明确。
- 唯一业务键建唯一索引（幂等基础），如订单号、支付流水号、卡券实例码。

## 2. 核心表清单（按工单归集，建表时按需细化）

| 工单 | 表/实体（建议） |
|---|---|
| 01 | 通知模板、通知消息、通知发送记录、频控指纹记录 |
| 02 | 账户余额、余额明细、充值订单、充值记录 |
| 03 | 城市、寄存点、寄存指南（含版本/多语言字段） |
| 04 | 柜子/柜门、收费规则、寄存订单、支付流水 |
| 05 | 订单主表、订单明细、临时开柜授权、补缴记录 |
| 06 | 用户资料、优惠卡券（可用/过期）、我的服务配置、会话/令牌 |
| 07 | 知识库条目、对话会话记录（如落库） |
| 10 | 设备信息、巡检报告、运维 Skills 定义、审计记录 |

> 表名仅为规划建议，落库以《设计说明书》与 ER 图为准；每个字段须写清来源工单编号。

## 3. 索引与查询

- 高频查询字段建索引：城市、寄存点坐标（空间/地理索引）、订单状态+用户+时间、卡券状态+有效期。
- 订单列表按 (用户, 时间/游标) 建复合索引支撑游标分页；避免深分页。
- 慢查询结合 SQL 日志与火焰图（工单 08）定位；避免 `SELECT *`、避免无索引排序。

## 4. 归档与分区（工单 05）

- 订单历史数据按时间归档/分区（如按月），进行中订单与历史订单分表或分区，保证查询性能。
- 归档策略需可回滚，归档过程有审计记录。

## 5. 缓存策略

- 缓存键命名规范：`业务:实体:标识[:条件]`，设置合理 TTL 与随机抖动防雪崩。
- 缓存优先 + 精准失效：订单详情、城市列表、热点寄存点；金额/状态类以 DB 为准。
- 频控（工单 01）、幂等键（工单 02/04/05）可借助缓存实现，但持久化与一致性需兜底。
- 缓存不可用需降级，不得导致核心交易流程不可用。

## 6. 数据字典与脚本

每个工单需产出：数据字典、ER 图、状态转换图、建表及索引更新脚本；脚本文件名含工单编号与版本，可重复执行且向后兼容。

## 7. 代码注释

建表脚本与模型定义处标注工单编号：`-- @ticket 04-物品寄存：柜门占用状态` / `// @ticket 04-物品寄存：订单实体`。
