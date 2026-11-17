# 工单07 数字人Coze智能体 — Live2D 数字人配置文件与渲染实现代码
> **工单编号：07｜数字人Coze智能体**　需求来源：智能快寄柜 10 份任务工单 V1.x　创建人：邹世军　出品：八维文化与产业研究院（内容由 AI 生成，仅供参考）

## 1. 模型与角色配置（src/config/index.ts）
运行时通过环境变量配置，支持角色切换：`CONFIG` 含 `cozeApiBase / botId / token / wsUrl / ttsProvider / modelPath`。
- ttsProvider 可切换 `web-speech | azure | openai`；modelPath 默认 `/models/haru/haru.model3.json`。
- 角色表 `ROLES`：`presenter`（品牌代言人，voice=zh-CN-Xiaoxiao）、`guide`（虚拟讲解员，voice=zh-CN-Yunxi）。
- `setRole(id)/getRole()` 实现运行时角色切换。

## 2. 渲染桥接组件（src/live2d/Live2DViewer.tsx）
生产环境将 canvas 交给 awesome-digital-human-live2d 的 PIXI/Live2D Core SDK 托管（商业 SDK 需自行授权，见 public/models/README.md）；本层保留标准挂载点与动作驱动接口，依赖隔离、可替换任意渲染后端。
```tsx
interface Props { modelPath: string; action: Live2DAction | null; onError?: (e: string) => void }
export default function Live2DViewer({ modelPath, action, onError }: Props) {
  const hostRef = useRef<HTMLDivElement>(null)
  useEffect(() => {                        // 模型挂载点
    if (!modelPath) { onError?.('未配置 Live2D 模型路径'); return }
    hostRef.current?.setAttribute('data-model', modelPath)   // mountLive2D(host, modelPath) 由部署方注入
  }, [modelPath, onError])
  useEffect(() => {                        // 动作/表情驱动
    if (!action || !hostRef.current) return
    hostRef.current.setAttribute('data-action', action.action)   // nod/smile/shake/blink
    hostRef.current.setAttribute('data-emotion', action.emotion ?? 'neutral')
    const t = setTimeout(() => hostRef.current?.removeAttribute('data-action'), action.duration_ms ?? 1200)
    return () => clearTimeout(t)
  }, [action])
  return <div ref={hostRef} className="live2d-host" aria-label="Live2D 数字人渲染区（ticket 07）"><canvas id="live2d-canvas" /></div>
}
```
> 模型路径 `data-model`、动作 `data-action`、情绪 `data-emotion`、持续 `duration_ms` 由渲染后端消费，前端只负责下发，做到渲染与业务解耦。
