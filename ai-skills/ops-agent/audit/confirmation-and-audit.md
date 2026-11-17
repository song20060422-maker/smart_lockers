---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_96821e30e27c11f19919525400638852
    ReservedCode1: EWWpAbx/lfBcAXmTQIkbpCNR+HzS6ECzcasN3nPgT1oVOegRyNsJS8ZybAYT52sMiJDEwkYWbbNIhBvNEA7Ei3YO1A4GNDhca+Qz0irsflHcFhBVhgfhFBThyj1TKLSXjmuY3vSLQqQuq3o6FeQ1MNvkNT+kVlldOx7lvC91nlxofno378q3b0A8d2s=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_96821e30e27c11f19919525400638852
    ReservedCode2: EWWpAbx/lfBcAXmTQIkbpCNR+HzS6ECzcasN3nPgT1oVOegRyNsJS8ZybAYT52sMiJDEwkYWbbNIhBvNEA7Ei3YO1A4GNDhca+Qz0irsflHcFhBVhgfhFBThyj1TKLSXjmuY3vSLQqQuq3o6FeQ1MNvkNT+kVlldOx7lvC91nlxofno378q3b0A8d2s=
---

<!-- @ticket 10：风险动作确认与审计机制（无确认不得执行、有记录可追溯） -->

# 风险动作确认与审计（工单 10）

## 适用范围
高风险动作类 Skills（11 重启 / 12 隔离 / 13 回滚），风险等级=高，是否自动执行=否。

## 强制规则
1. **无确认不得执行**：即使由事件自动触发，高风险动作也只能产出 `pending_confirmation` 状态，不得进入执行。
2. **确认要素**：确认人（工号/账号）、确认时间、确认依据（关联分析结论/工单号）、动作与目标设备。
3. **超时语义**：确认超时（建议 5 分钟）视为拒绝，状态置 `rejected`，不执行。
4. **可追溯**：每次确认/拒绝均生成审计记录，包含动作全生命周期状态。

## 审计记录结构（audit-log）
```json
{
  "audit_id": "AUD-<时间戳>-<序号>",
  "ticket": "10",
  "action": "restart|isolate|rollback",
  "device_id": "0203445006",
  "status": "pending_confirmation|confirmed|rejected|executed|failed",
  "trigger": "dialogue|event",
  "requester": "触发来源（人/系统）",
  "evidence": ["05 告警归因结论", "工单号"],
  "confirm": { "operator": "确认人", "time": "ISO8601", "basis": "确认依据" },
  "created_at": "ISO8601",
  "updated_at": "ISO8601"
}
```

## 状态机
`pending_confirmation → confirmed → executed(由底层系统执行，非本 Skill 实现) → 记录归档`
任一环节失败写入 `status=failed` 并保留上下文，便于复盘。

## 使用价值
- 满足工单验收项：无确认不得执行、有记录可追溯。
- 审计记录可用于事件复盘、运维质量评估与责任追溯。
*（内容由AI生成，仅供参考）*
