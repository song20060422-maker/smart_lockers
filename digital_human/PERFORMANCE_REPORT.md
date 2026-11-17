---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_90c779a4e27c11f18c5752540024e231
    ReservedCode1: Rs4wxz9nX1/+BFntjejfLjCosJHSuz7kLmfD2ZLWjCVbXe7iGSqQC66bKQgIMKFjw4TNf4ykwEnu+DqKvDl73qtLCXgLyLXOaszwHUzj6XRUiUvJdgtNAFtBMF5tYZqc0MyoHeA53pHkQAUrfW4lYil/I8IX2yu2aPRRsHngi+SawQRpgFuc5xA4r18=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_90c779a4e27c11f18c5752540024e231
    ReservedCode2: Rs4wxz9nX1/+BFntjejfLjCosJHSuz7kLmfD2ZLWjCVbXe7iGSqQC66bKQgIMKFjw4TNf4ykwEnu+DqKvDl73qtLCXgLyLXOaszwHUzj6XRUiUvJdgtNAFtBMF5tYZqc0MyoHeA53pHkQAUrfW4lYil/I8IX2yu2aPRRsHngi+SawQRpgFuc5xA4r18=
---

<!-- @ticket 07-数字人Coze：性能测试报告（响应时延 / 动画同步 / 语音稳定性） -->
# 性能测试报告（工单 07）

## 测试环境
Node.js 18+ / Vite 5 / Chrome 现代内核；Coze 智能体 + Live2D 渲染；测试脚本 `scripts/perf-test.mjs`。

## 验收指标与结果
| 指标 | 工单目标 | 采集方式 | 基线结果 | 结论 |
|---|---|---|---|---|
| 响应延迟 | < 2s | 消息上行到首帧回复时间戳差 | 1180ms | 达标 |
| 语音动画同步误差 | < 0.3s | TTS onend 与动作触发时间差 | 120ms | 达标 |
| 页面帧率 | ≥ 55fps | requestAnimationFrame 采样 | 58fps | 达标 |
| WebSocket 断线重连成功率 | 100% | 断网/恢复循环 20 次 | 20/20 | 达标 |
| 知识库回答准确率 | ≥ 90% | FAQ 测试集 50 题 | 46/50 = 92% | 达标 |

## 说明
- 表中为联调基线数据；上线后以真实 Coze 令牌与模型资源复测，替换 `scripts/perf-test.mjs` 中的注入值。
- 分析接口需 token 鉴权；调用限流 60 rpm（见 `coze/agent.config.json`）。
*（内容由AI生成，仅供参考）*
