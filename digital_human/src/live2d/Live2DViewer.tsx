/** @ticket 07-数字人Coze：Live2D 渲染桥接层——模型路径配置 + 表情/姿态动画驱动（对接 awesome-digital-human-live2d） */
import { useEffect, useRef } from 'react'
import type { Live2DAction } from '@/types'

interface Props {
  modelPath: string
  action: Live2DAction | null
  onError?: (e: string) => void
}

/**
 * 该组件为 Live2D 渲染的桥接抽象：
 * - 生产环境将 canvas 交给 awesome-digital-human-live2d 的 PIXI/Live2D Core SDK 托管（商业 SDK 需自行授权，见 public/models/README.md）。
 * - 此处保留标准挂载点、模型路径配置与动作驱动接口，保证依赖隔离、可替换为任意渲染后端。
 */
export default function Live2DViewer({ modelPath, action, onError }: Props) {
  const hostRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!modelPath) { onError?.('未配置 Live2D 模型路径'); return }
    // 实际 SDK 挂载：mountLive2D(hostRef.current, modelPath)，由部署方注入
    hostRef.current?.setAttribute('data-model', modelPath)
  }, [modelPath, onError])

  useEffect(() => {
    if (!action || !hostRef.current) return
    // 动作驱动：节点微笑/点头/摇头/眨眼等；由渲染后端消费 data-action
    hostRef.current.setAttribute('data-action', action.action)
    hostRef.current.setAttribute('data-emotion', action.emotion ?? 'neutral')
    const t = setTimeout(() => hostRef.current?.removeAttribute('data-action'), action.duration_ms ?? 1200)
    return () => clearTimeout(t)
  }, [action])

  return (
    <div ref={hostRef} className="live2d-host" aria-label="Live2D 数字人渲染区（ticket 07）">
      <canvas id="live2d-canvas" />
    </div>
  )
}
