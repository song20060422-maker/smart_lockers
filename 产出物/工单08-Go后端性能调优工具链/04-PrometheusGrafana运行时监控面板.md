# 工单08 Go后端性能调优工具链 — Prometheus/Grafana 运行时监控面板
> **工单编号：08｜Go后端性能调优工具链**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 1. 指标来源
- 应用 runtime/metrics 通过 `/api/v1/ops/expvar` 暴露（Go expvar + runtime/metrics）。
- 建议用 `prometheus/client_golang` 暴露 `/metrics`，采集以下维度。

## 2. 核心指标
| 指标 | 类型 | 说明 |
|---|---|---|
| go_goroutines | Gauge | 当前协程数（泄漏监控） |
| go_memstats_heap_inuse_bytes | Gauge | 堆使用量 |
| go_gc_duration_seconds | Summary | GC 耗时分布 |
| http_request_duration_seconds | Histogram | 接口时延（按 method/path/status） |
| http_requests_total | Counter | 请求量/QPS |
| db_conn_wait_total | Counter | DB 连接池等待 |

## 3. Grafana 面板设计
- **概览行**：QPS、P95 时延、错误率、协程数、堆内存（Stat 面板）。
- **时延面板**：http_request_duration_seconds 的 P50/P95/P99 折线（按接口分组）。
- **运行时面板**：协程数趋势、GC 耗时与频率、堆 inuse/alloc 双轴。
- **告警**：P95>500ms 持续 5min、协程数异常增长、错误率>1% 触发告警。

## 4. 部署要点
- Prometheus 定时拉取 `/metrics`；面板 JSON 可版本化管理，纳入 `scripts/monitoring/`。
- 监控接口同样需内网+鉴权保护，避免指标泄露。
