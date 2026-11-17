# 工单08 Go后端性能调优工具链 — 火焰图与 trace 时间线分析结果
> **工单编号：08｜Go后端性能调优工具链**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 1. 火焰图（CPU Profile）
- 采集：`curl -o cpu.pprof "http://127.0.0.1:18080/api/v1/ops/pprof/profile?seconds=30"`（内网）。
- 生成：`go tool pprof -http=:0 cpu.pprof`，View → Flame Graph。
- **读法**：横轴为采样占比（宽=耗时多），纵轴为调用栈（上Caller下Callee）；关注最宽的平台顶部函数。
- **样例结论**：`database/sql.(*DB).QueryContext` 与 `encoding/json.Marshal` 占比较宽，指向 DB 查询与序列化为主要 CPU 消耗。

## 2. 内存/协程剖面
- heap：`go tool pprof -inuse_space` 观察存活对象；`-alloc_space` 观察累计分配。
- goroutine：`/pprof/goroutine?debug=2` 查看栈，识别阻塞/泄漏。

## 3. trace 时间线
- 采集：`curl -o trace.out "http://127.0.0.1:18080/api/v1/ops/trace?seconds=5"`。
- 查看：`go tool trace trace.out`，关键视图：
  - **View trace**：协程在时间轴上的运行/阻塞/就绪，观察调度延迟。
  - **Goroutine analysis**：按函数统计协程阻塞时长。
  - **Synchronization blocking profile**：定位锁/通道等待点。
- **样例结论**：部分协程在 DB 连接池等待，建议核对 `MaxOpenConns`，避免连接饥饿。

## 4. 分析口径说明
以上为火焰图/trace 的采集与解读方法 + 样例结论；真实图形与数值须在目标环境采样后生成，本文件提供可复现流程。
