---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_f3cd8207e4ad11f1b820525400de85a5
    ReservedCode1: zqtDMvFX6AVU3gJMqfyGIwU63adqYH6/DL9dR+qVPIAQ9BANg91s+tNqfC43S9VZmyOpEG1SCieuXW863wcaSa8kr7UWj0HnvYM+4+0ZD6Kl9abFLGNZs8AM7AN4WCKdi1c73S8iWscqF0hUvBd+XvPmzze6n48vyiIubCB6dARjUSpVOllQyREqTeg=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_f3cd8207e4ad11f1b820525400de85a5
    ReservedCode2: zqtDMvFX6AVU3gJMqfyGIwU63adqYH6/DL9dR+qVPIAQ9BANg91s+tNqfC43S9VZmyOpEG1SCieuXW863wcaSa8kr7UWj0HnvYM+4+0ZD6Kl9abFLGNZs8AM7AN4WCKdi1c73S8iWscqF0hUvBd+XvPmzze6n48vyiIubCB6dARjUSpVOllQyREqTeg=
---

# 09 · Git 提交与分支规范

> 适用：所有代码变更的分支创建与 commit 提交。前置：`00-project.md`。

## 1. 分支开发规范

- 所有代码变更必须先创建分支再开发提交，禁止直接在主分支/默认分支上写代码。
- 分支命名基于工单/需求编号，示例：
  - `feat/T-134`（新功能，对应需求/工单编号）
  - `fix/xxx`（Bug 修复，对应 Bug 编号）
- 基本流程：分支 → 提交 → 合并。
  1. 从最新主分支创建工作分支（按上述命名约定）。
  2. 在分支上开发，并按本文件规范撰写 commit。
  3. 变更完成后发起合并回主分支。

## 2. 提交格式

```
[类型] [版本号] [需求/Bug编号] 描述
```

## 3. 类型

| 类型 | 用途 |
|------|------|
| feat | 新增功能需求 |
| fix | 修复 Bug |
| hotfix | 线上紧急修复（来不及提 Bug 单，事后补） |
| opt | 项目优化（性能、代码质量、重构等） |
| docs | 文档变更 |
| test | 测试相关 |

## 4. 必填规则

| 类型 | 版本号 | 编号 | 说明 |
|------|--------|------|------|
| feat | 必填 | 必填 | 新功能必须关联需求编号 |
| fix | 必填 | 必填 | Bug 修复必须关联 Bug 编号 |
| hotfix | 必填 | 可选 | 线上紧急修复，编号可事后补 |
| opt | 必填 | 可选 | 优化可以不关联编号 |
| docs | 可选 | 可选 | 文档变更可以都不带 |
| test | 可选 | 可选 | 测试相关可以都不带 |

## 5. 示例

- `[feat] [v2.1.0] [#6974499249] 新增用户注册接口`
- `[fix] [v2.1.0] [#6922369604] 修复详情接口缺少 appId 过滤`
- `[hotfix] [v2.1.0] 修复线上用户无法登录问题`
- `[hotfix] [v2.1.0] [#6974499251] 修复线上支付回调超时（事后补编号）`
- `[opt] [v2.1.0] 优化首页列表查询性能，增加复合索引`
- `[opt] [v2.1.0] [#6974499250] 重构支付回调逻辑`
- `[docs] 更新 API 接口文档`
- `[test] [v2.1.0] 补充订单模块单元测试`

## 6. 规则

- 描述使用中文，简洁明了
- 版本号格式：v{主版本}.{次版本}.{修订号}，如 v2.1.0
- 编号是项目的工作项 ID（需求或缺陷），格式 #数字
- 一次 commit 只包含一个类型，混合变更拆分为多次 commit
- feat 和 fix 提交时，AI 会通过 meegle CLI 校验编号是否真实存在

## 7. 编号校验

提交 feat 或 fix 类型时，AI 会自动执行：

```bash
meegle workitem get --work-item-id {编号} --project-key r17te9
```

如果编号不存在或查询失败，会提醒开发者确认编号是否正确。
