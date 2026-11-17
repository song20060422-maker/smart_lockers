# 工单07 数字人Coze智能体 — Coze 智能体配置文件（插件/工作流/知识库）
> **工单编号：07｜数字人Coze智能体**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 说明
数字人对话智能体"智能快寄柜数字人助手"的 Coze 侧配置文件共四份，均位于 `digital_human/coze/`，配置化驱动、可替换。以下为原文要点。

## 1. 主配置 agent.config.json
- **智能体**：name=`智能快寄柜数字人助手`，model=`coze-default`，locale=`zh-CN`；面向寄存场景，支持问候/FAQ/闲聊/情绪响应。
- **通信协议**：WebSocket（实时对话）+ RESTful；消息结构 `{role: user|assistant, content: string, type: text|voice|command, timestamp: number}`。
- **安全**：`auth=token`、`cors=allowlist`、`rate_limit_rpm=60`。
- **关联**：plugins=`./plugins.json`、workflow=`./workflow.json`、knowledge_base=`./knowledge_base.json`。

## 2. 插件 plugins.json（4 个）
| 插件 | 方向 | 参数 | 用途 |
|---|---|---|---|
| message_forward | 双向 | role/content/type/timestamp | 用户消息上行与回复下行的统一通道 |
| tts_output | 出站 | text/voice/emotion | 文本转语音，驱动前端 SpeechSynthesis/第三方 TTS |
| api_call | 出站 | endpoint/method/payload | 调用外部只读 API（寄存点查询、订单状态） |
| live2d_action_trigger | 出站 | action/emotion/duration_ms | 下发 Live2D 动作（nod/smile/shake/transform） |

## 3. 工作流 workflow.json
主流程 `digital_human_main_flow`：
```
greeting(问候/start) → intent(意图识别)
  ├─ faq → knowledge_retrieval(kb=locker_faq) → response
  ├─ chitchat → llm → response
  └─ emotion → emotion_detect → response
response(输出, attach=[emotion,scene,task], memory=context_window) → action(live2d_action_trigger)
```
- 上下文记忆：sliding_window，max_turns=10。
- 降级策略：出错 `friendly_retry`，低置信度 `ask_clarification`。

## 4. 知识库 knowledge_base.json
- id=`locker_faq`，来源：`faq.json`(weight 1.0) + `brand.md`(weight 0.8)。
- 检索：top_k=3，min_score=0.7，strategy=`semantic+keyword`。
- 标签：品牌/寄存/收费/取件/故障/景区。
> 以上四份配置均已标注 `@ticket 07-数字人Coze`。
