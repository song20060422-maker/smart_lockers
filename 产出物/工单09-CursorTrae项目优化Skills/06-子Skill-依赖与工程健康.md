# 工单09 Cursor/Trae项目优化Skills — 子 Skill：依赖与工程健康度（dependency.md）
> **工单编号：09｜CursorTrae项目优化Skills**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 维度：依赖与工程健康度
聚焦依赖合理性与工程结构健康，输出遵循统一输出结构（见 01）。

## 关注点
- **依赖膨胀**：无用依赖、功能重复库、未使用 require。
- **版本风险**：过时/存在漏洞版本、主版本混用。
- **工程规范**：缺少 go.sum 校验、无统一 lint、无 CI 检查。
- **配置管理**：明文凭据入库风险（如 config.yaml 应 gitignore，提供 example 模板）。

## 分析依据来源
- go.mod/go.sum 条目、.gitignore、CI 配置。

## 输出示例
```markdown
## 优化结论 #1
- **问题类型**：依赖 / 工程健康
- **风险等级**：中（凭据泄露风险）
- **工程背景**：config/config.yaml 含明文 DB 密码与 JWT secret
- **优化建议**：纳入 .gitignore，提供 config.example.yaml，真实配置走环境变量/密钥管理
- **依据**：config.yaml 中 dsn: root:root@tcp(...)
```
