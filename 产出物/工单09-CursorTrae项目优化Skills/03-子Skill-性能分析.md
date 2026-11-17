# 工单09 Cursor/Trae项目优化Skills — 子 Skill：性能分析（performance.md）
> **工单编号：09｜CursorTrae项目优化Skills**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 维度：性能
聚焦 Go 项目运行时效率与资源使用，输出遵循统一输出结构（见 01）。

## 关注点
- **数据库访问**：N+1 查询、缺失索引、全表扫描、未用覆盖索引、连接未释放。
- **序列化/分配**：高频 JSON Marshal/Unmarshal、字符串拼接、临时对象导致 GC 压力。
- **缓存使用**：热点数据未缓存、缓存粒度不当、缓存穿透/击穿。
- **I/O**：同步阻塞 I/O、未批量处理、无超时控制。

## 分析依据来源
- pprof CPU/heap 剖面（工单08工具链）、具体文件/函数（如 handler → service → model 调用链）。
- 引用示例：`internal/service/order.go` 中列表查询循环内单条查询 → N+1。

## 输出示例（遵循统一结构）
```markdown
## 优化结论 #1
- **问题类型**：性能
- **风险等级**：中（列表接口 P95 偏高，影响用户体验）
- **工程背景**：internal/service/order.go ListOrders 在循环内按 order_id 逐条查询
- **优化建议**：改为 JOIN 或批量 IN 查询 + 覆盖索引 (user_id,status)
- **依据**：for 循环内 db.QueryRow("SELECT ... WHERE id=?", id)
```
