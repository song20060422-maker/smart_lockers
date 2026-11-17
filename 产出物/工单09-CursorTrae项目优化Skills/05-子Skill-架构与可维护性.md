# 工单09 Cursor/Trae项目优化Skills — 子 Skill：架构与可维护性（architecture.md）
> **工单编号：09｜CursorTrae项目优化Skills**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 维度：架构与可维护性
聚焦模块职责边界与演进风险，输出遵循统一输出结构（见 01）。

## 关注点
- **职责越界**：handler 内写业务逻辑、model 混入 HTTP 逻辑。
- **循环依赖 / 包耦合**：internal 包边界模糊、接口与实现未分离。
- **重复代码**：多处相同逻辑未抽象。
- **可测试性**：强耦合导致难以单测（未依赖接口）。

## 分析依据来源
- 目录分层（cmd/internal/pkg）、import 关系、函数长度与职责。

## 输出示例
```markdown
## 优化结论 #1
- **问题类型**：架构
- **风险等级**：中（职责越界，后续演进与维护成本上升）
- **工程背景**：internal/handler/order.go 直接拼接 SQL 并操作 DB
- **优化建议**：下沉到 service 层，handler 仅做参数校验与响应封装
- **依据**：handler.CreateOrder 内出现 db.Exec(...)
```
