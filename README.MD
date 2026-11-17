---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: b2aa3b4cd48b832ce0392d027319fcdd_830c25a0e27c11f18c5752540024e231
    ReservedCode1: tmPcG3iVx9YMJj9D3W4LUxOl7LvC5OCVRLK1CkQJAcgxJWdOPV7MYwQWHApdHYro0dWyKbdQkZ0aP1JLV/Xro6bixHV0dqYC8gGLaA7Bz9g3Y0XkM/Q2vfdQ8wZmMMzpL+0cDuEJglvsG4/h34ZnYcqBfNlxFPQW7SQhdnP8QDEusJRcK2KWaif3Pt8=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: b2aa3b4cd48b832ce0392d027319fcdd_830c25a0e27c11f18c5752540024e231
    ReservedCode2: tmPcG3iVx9YMJj9D3W4LUxOl7LvC5OCVRLK1CkQJAcgxJWdOPV7MYwQWHApdHYro0dWyKbdQkZ0aP1JLV/Xro6bixHV0dqYC8gGLaA7Bz9g3Y0XkM/Q2vfdQ8wZmMMzpL+0cDuEJglvsG4/h34ZnYcqBfNlxFPQW7SQhdnP8QDEusJRcK2KWaif3Pt8=
---

# 智能快寄柜系统（Smart Lockers）

基于 10 份任务工单（V1.x，2026-01）实现的一套智能快寄柜全栈系统：Go + Gin 后端、Vue3 + Vite 前端，另含 Coze 数字人、Go 性能调优工具链、IDE Skills、设备运维智能体等专项工单。

## 一、仓库信息

- 远程仓库：<https://github.com/song20060422-maker/smart_lockers.git>
- 需求来源：`D:\Users\Lenovo\Desktop\智能快寄柜`（10 份 PDF 工单）
- 原型地址：<https://modao.cc/proto/3fxWzFGJsp37iCk5fS2h/sharing?view_mode=read_only>
- 工单创建人：邹世军；出品：八维文化与产业研究院

## 二、目录结构

```
smart_lockers/                     # 本目录（工程配置与 AI 规则层）
├── AGENT.MD                       # AI 协作入口规则
├── README.MD                      # 本说明
├── .gitignore                     # 忽略规则
├── .agent/                        # AI 分规则集（详见 .agent/README.md）
├── smart_lockers/                 # 后端：Go + Gin（REST API）
└── smart_lockers_frontend/        # 前端：Vue3 + Vite
```

> 后端与前端为同级子目录；本层的 `AGENT.MD`/`README.MD`/`.gitignore`/`.agent/` 归属于整个仓库，不放入任一子项目内部。

## 三、技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.20+、Gin、关系型数据库（MySQL 系）、缓存（Redis 系）、GORM/sqlx |
| 前端 | Vue3、Vite、TypeScript、Pinia、Vue Router、Axios |
| 数字人（工单07，独立模块） | Coze（插件/工作流/知识库）、awesome-digital-human-live2d、React + TS + Vite、Live2D、WebSocket、Web Speech/TTS |
| 性能调优（工单08） | pprof、trace、expvar、runtime/metrics、gops、pyroscope/Parca、Prometheus/Grafana、火焰图 |
| AI 能力（工单09/10） | IDE Skills（Cursor/Trae）、运维 Agent + Skills |

## 四、工单一览

| 编号 | 工单 | 核心功能 |
|---|---|---|
| 01 | 消息通知 | 多渠道通知告警中心、模板渲染、去重频控、重试升级、审计监控 |
| 02 | 快寄员支付宝支付 | 余额/明细/充值/充值记录，资金可追溯 |
| 03 | 首页 | 定位、城市选择、附近寄存点搜索、寄存指南 |
| 04 | 物品寄存 | 柜子选择、收费规则、支付落单、异常换柜 |
| 05 | 订单查看与取出 | 订单列表/详情、临时开柜、超时补交、进行中聚合 |
| 06 | 个人中心 | 订单、头像、我的服务、优惠卡券、退出登录 |
| 07 | 数字人 Coze 智能体 | Coze 智能体 + Live2D 数字人 + 语音联动 |
| 08 | Go 后端性能调优工具链 | pprof/trace/持续分析/火焰图/监控面板 |
| 09 | Cursor/Trae 项目优化 Skills | Go 项目优化 Skills 能力配置 |
| 10 | 结合 Skills 开发设备运维智能体 | 运维 Agent + ≥10 运维 Skills + 风险确认审计 |

各工单对应模块、数据表、接口与业务细节见 `.agent/16-tickets.md`。

## 五、快速开始

> 后端子项目当前为空仓库（仅 `.git`，无提交）；前端目录当前为空。以下为搭建完成后的预期启动方式。

```bash
# 后端
cd smart_lockers
go mod init <module>
# 依赖安装后，配置文件集中在 config/ 下，直接运行：
go run main.go

# 前端
cd smart_lockers_frontend
npm install
npm run dev
```

约定：后端配置写在本项目内 config 文件（本地明文，便于 `go run` 直接启动），不依赖临时注入环境变量；真实密码、密钥由使用者自行填写，仓库不提交敏感凭据。

## 六、开发规范

开发前先读 [`AGENT.MD`](AGENT.MD) 与 `.agent/` 分规则。硬性要求包括：代码注释标注工单编号、单测覆盖率 ≥70%、统一错误码与日志、接口压测报告、资金/状态变更接口的事务与幂等、敏感信息脱敏与审计。

## 七、当前状态

- [x] 顶层工程配置：AGENT.MD、README.MD、.gitignore、.agent 分规则
- [ ] 后端 Go + Gin 骨架搭建
- [ ] 工单 01–10 功能实现
*（内容由AI生成，仅供参考）*
