# 工单09 Cursor/Trae项目优化Skills — 子 Skill：并发与稳定性（concurrency.md）
> **工单编号：09｜CursorTrae项目优化Skills**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 维度：并发与稳定性
聚焦并发安全、资源生命周期与故障恢复，输出遵循统一输出结构（见 01）。

## 关注点
- **竞态**：共享变量未加锁/未用 sync 原语；建议 `go test -race` 验证。
- **协程泄漏**：goroutine 无退出条件、无 context 取消、channel 未关闭。
- **锁竞争**：粗粒度锁、锁内耗时操作（可用 pprof block 定位）。
- **超时与取消**：外部调用无 context timeout、无重试退避。
- **资源池**：DB 连接池 MaxOpenConns/MaxIdleConns 配置不当。

## 分析依据来源
- `-race` 检测报告、pprof goroutine/block、具体同步代码块。

## 输出示例
```markdown
## 优化结论 #1
- **问题类型**：并发
- **风险等级**：高（存在数据竞态，可能导致脏数据）
- **工程背景**：internal/service/cache.go 全局 map 被多请求并发读写无锁
- **优化建议**：改用 sync.Map 或 RWMutex，补充 -race 单测
- **依据**：var cache = map[string]Order{} 直接 get/set
```
