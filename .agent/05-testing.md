---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_8b128fb7e27c11f19919525400638852
    ReservedCode1: bBj6NnPtCEB3G6ZTikKR7cNWelorayMoqdbPakc2PTUeYlrkn9dJtwnA+p8MHBJ5uNg3lxSbxGY7+6inK0lUYcOkmMBHAF5ocZ2LwdYQt+bUbGnbDFXyO8o9i9gFX6tez7VG1Tsg1ICDwWFFutcVSylY+kkXjFkmPxDod7nhZL/FU8Ma51YDgkBpftc=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_8b128fb7e27c11f19919525400638852
    ReservedCode2: bBj6NnPtCEB3G6ZTikKR7cNWelorayMoqdbPakc2PTUeYlrkn9dJtwnA+p8MHBJ5uNg3lxSbxGY7+6inK0lUYcOkmMBHAF5ocZ2LwdYQt+bUbGnbDFXyO8o9i9gFX6tez7VG1Tsg1ICDwWFFutcVSylY+kkXjFkmPxDod7nhZL/FU8Ma51YDgkBpftc=
---

# 05 · 测试与验收规范

> 适用：所有工单。前置：`00-project.md`、`04-api.md`。

## 1. 单元测试（覆盖率 ≥70%）

- 核心 service、收费计算、状态机、幂等逻辑、模板变量渲染等必须有单测；覆盖率目标 ≥70%（10 工单共性要求）。
- 表驱动测试为主；外部依赖（DB/缓存/第三方）用 mock/接口隔离。
- 关键边界：余额不足、重复回调、超时补缴、卡券过期、头像格式非法、频控触发。

## 2. 集成 / 联调测试

- 提供联调用例与测试报告；前后端按统一报文对接。
- 涉及第三方：支付宝充值（02）、Coze API（07）、短信/邮件/回调渠道（01）需 mock 与真实环境两套用例。
- 事务一致性用例：支付落单失败回滚（04）、补缴幂等（05）、充值入账一致（02）。

## 3. 压力测试（10 工单产出物共性要求）

- 每个核心接口提供压测脚本与压测报告，含 QPS、P50/P95/P99、错误率、资源占用。
- 对齐工单 P95 指标：定位 ≤200ms；城市搜索/寄存点查询/订单查询 P95 ≤300ms；柜子查询 P95 ≤200ms、收费计算 P95 ≤250ms；头像上传 P95 ≤500ms；数字人响应 <2s；采样开销 <3%。
- 高并发场景（千级并发）验证缓存命中、锁竞争、连接池；结合工单 08 火焰图/trace 分析瓶颈。
- 报告需含"压测报告解读"：瓶颈定位、优化前后对比、结论。

## 4. 验收口径

| 维度 | 通用要求 |
|---|---|
| 功能 | 覆盖工单功能概述全部条目，状态流转与权限正确 |
| 性能 | 达到该工单标注的 P95/时延/并发指标 |
| 安全 | 鉴权、脱敏、幂等、审计、高风险确认到位 |
| 代码质量 | 覆盖率 ≥70%、错误码日志统一、注释含工单编号、命名规范 |
| 文档 | 设计说明书、数据字典、ER、状态图、接口文档、压测/联调报告齐全 |

## 5. 报告产物

测试/压测报告命名含工单编号与日期，随工单交付归档，供复盘与质量评估复用。

## 6. 代码注释

测试用例标注工单编号：`// @ticket 04-物品寄存：收费计算边界用例`。
