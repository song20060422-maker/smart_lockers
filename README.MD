
# 智能快寄柜系统（Smart Lockers）

基于 10 份任务工单（V1.x，2026-01）实现的一套智能快寄柜全栈系统：Go + Gin 后端、Vue3 + Vite 前端，另含 Coze 数字人、Go 性能调优工具链、IDE Skills、设备运维智能体等专项工单。

## 一、仓库信息

- 远程仓库：<https://github.com/song20060422-maker/smart_lockers.git>
- 需求来源：`D:\Users\Lenovo\Desktop\智能快寄柜`（10 份 PDF 工单）


## 二、目录结构

```
smart_lockers/                     # 本目录（工程配置与 AI 规则层）
├── AGENT.MD                       # AI 协作入口规则
├── README.MD                      # 本说明
├── .gitignore                     # 忽略规则
├── .agent/                        # AI 分规则集（详见 .agent/README.md）
├── smart_lockers/                 # 后端：Go + Gin（REST API，:18080）
├── smart_lockers_frontend/        # 前端：Vue3 + Vite
├── smart_lockers_admin/           # B 端运营后端：Go 1.23 + Gin（REST API，:18081；RBAC/审计/对账/告警）
├── smart_lockers_admin_frontend/  # B 端运营前端：Vue3.5 + Vite6 + TS + Pinia（登录守卫/按钮级 RBAC/四模块）
├── digital_human/                 # 工单07 数字人独立子模块（React + TS + Live2D）
└── ai-skills/                     # 工单09/10 Skills 相关
```

> 后端与前端为同级子目录；本层的 `AGENT.MD`/`README.MD`/`.gitignore`/`.agent/` 归属于整个仓库，不放入任一子项目内部。`smart_lockers`、`smart_lockers_frontend`、`smart_lockers_admin`、`smart_lockers_admin_frontend`、`digital_human` 五个代码子目录各自带独立 `AGENT.MD`、`README.MD` 与 `.agent/` 分规则集，按各自技术栈定制；`ai-skills/` 为 Skill 文件夹，不含 .agent 规则。

## 三、技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.20+、Gin、关系型数据库（MySQL 系）、缓存（Redis 系）、GORM/sqlx |
| 前端 | Vue3、Vite、TypeScript、Pinia、Vue Router、Axios |
| B 端运营后端（新增） | Go 1.23、Gin v1.10、JWT v5、YAML v3（:18081）；handler→service→store 分层、三角色 RBAC、审计/对账/告警、幂等键 |
| B 端运营前端（新增） | Vue 3.5、Vite 6、TypeScript 5.7、Pinia 2.3、Vue Router 4.5、Axios 1.7；登录守卫 + 按钮级 RBAC + 统一请求封装（业务码/幂等键） |
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

各工单对应模块、数据表、接口与业务细节见 `.agent/07-tickets.md`。

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

### B 端运营后台启动（新增）

```bash
# B 端运营后端（:18081，独立于用户端 :18080）
cd smart_lockers_admin
# 复制 config/config.example.yaml 为 config/config.yaml 并填写占位符后：
go run cmd/server/main.go
# 单测：go test ./...

# B 端运营前端
cd smart_lockers_admin_frontend
npm install
npm run dev          # 或 npm run build / npm run type-check
```

B 端后端凭据仅用 `config/config.example.yaml` 占位（demo_admins / jwt.secret / policy 均为示例值），禁止提交真实账号或密钥；B 端前端 API 目标见 `.env.example`。代码注释按 B 端约定使用 `// @ticket B端-…` 标注。

约定：后端配置写在本项目内 config 文件（本地明文，便于 `go run` 直接启动），不依赖临时注入环境变量；真实密码、密钥由使用者自行填写，仓库不提交敏感凭据。

### 本地配置

仓库出于安全考虑**不提交任何真实凭据**，仅提供配置模板。clone 后请按以下步骤填写本机配置：

1. **后端**：复制 `smart_lockers/config/config.example.yaml` 为 `smart_lockers/config/config.yaml`，然后编辑 `config.yaml`：
   - `database.dsn`：把 `<DB_USER>`、`<DB_PASSWORD>`、`<DB_HOST>`、`<DB_PORT>`、`<DB_NAME>` 替换为本机 MySQL 的账号、密码、地址、端口与库名（如 `root:your_pwd@tcp(127.0.0.1:3306)/smart_lockers?charset=utf8mb4&parseTime=True&loc=Local`）；建库建表可执行 `smart_lockers/scripts/init.sql`。
   - `redis.addr` / `redis.password`：填入本机 Redis 地址与密码（无密码则留空 `""`）。
   - `jwt.secret`：替换为自行生成的随机密钥（如 `openssl rand -hex 32`），切勿使用模板占位符。
2. **前端**：如需自定义 API 目标，复制 `smart_lockers_frontend/.env.example` 为 `.env.development`（或 `.env.production`），按需修改 `VITE_API_TARGET`（默认 `/api/v1` 经 Vite 代理转发到后端 `:18080`）。
3. **数字人子模块（工单07）**：复制 `digital_human/.env.example` 为 `.env.local`，填入真实的 Coze Bot ID 与 Token。

> 说明：`config.yaml`、`.env`、`.env.*` 已被 `.gitignore` 排除，仅 `*.example` 模板入库。模板中均为占位符，不含任何真实密码或密钥。

## 六、开发规范

开发前先读 [`AGENT.MD`](AGENT.MD) 与 `.agent/` 分规则。硬性要求包括：代码注释标注工单编号、单测覆盖率 ≥70%、统一错误码与日志、接口压测报告、资金/状态变更接口的事务与幂等、敏感信息脱敏与审计。

## 七、当前状态

- [x] 顶层工程配置：AGENT.MD、README.MD、.gitignore、.agent 分规则
- [x] B 端运营后台（新增）：smart_lockers_admin（Go + Gin :18081）+ smart_lockers_admin_frontend（Vue3.5 + Vite6），各带独立 AGENT.MD/.agent 分规则集
- [ ] 后端 Go + Gin 骨架搭建
- [ ] 工单 01–10 功能实现

