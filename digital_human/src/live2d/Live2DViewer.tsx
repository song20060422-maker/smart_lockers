/** @ticket 07-数字人Coze：Live2D 渲染桥接层——模型路径配置 + 表情/姿态动画驱动（对接 awesome-digital-human-live2d） */
import { useEffect, useRef } from 'react'
import type { Live2DAction } from '@/types'
import { CartoonAvatarEngine } from '@/live2d/CartoonAvatarEngine'

interface Props {
  modelPath: string
  action: Live2DAction | null
  // 声明属性 speaking，标识当前是否正在语音播报

  speaking: boolean
  // 声明属性 speakPulse，语音边界计数，变化时触发一次张口脉冲

  speakPulse: number
  onError?: (e: string) => void
}

/**
 * 该组件为 Live2D 渲染的桥接抽象：
 * - 生产环境将 canvas 交给 awesome-digital-human-live2d 的 PIXI/Live2D Core SDK 托管（商业 SDK 需自行授权，见 public/models/README.md）。
 * - 此处保留标准挂载点、模型路径配置与动作驱动接口，保证依赖隔离、可替换为任意渲染后端。
 */
export default function Live2DViewer({ modelPath, action, speaking, speakPulse, onError }: Props) {
  const hostRef = useRef<HTMLDivElement>(null)
  // 声明 ref engineRef，保存自绘引擎实例，供动作 effect 复用同一实例

  const engineRef = useRef<CartoonAvatarEngine | null>(null)

  // 挂载 effect：创建自绘引擎并启动动画循环，卸载时停止释放监听器

  useEffect(() => {
    // 声明变量 host，读取渲染容器节点

    const host = hostRef.current
    // 声明变量 canvas，读取容器内的 canvas 节点

    const canvas = host?.querySelector('canvas')
    // 条件判断：容器与 canvas 均存在时才初始化引擎

    if (!host || !canvas) return
    // 创建引擎实例，传入 canvas 与宿主容器

    const engine = new CartoonAvatarEngine(canvas, host)
    // 赋值更新：保存引擎实例供其他 effect 使用

    engineRef.current = engine
    // 调用方法：启动动画循环

    engine.start()
    // 返回清理函数：卸载时停止引擎并清空引用

    return () => { engine.stop(); engineRef.current = null }
  // 依赖数组为空，引擎生命周期与组件一致

  }, [])

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
    // 调用方法：将现有 payload 的动作/情绪交给自绘引擎，驱动形象的表情与动作

    engineRef.current?.setAction(action)
    const t = setTimeout(() => hostRef.current?.removeAttribute('data-action'), action.duration_ms ?? 1200)
    return () => clearTimeout(t)
  }, [action])

  // 播报状态 effect：把说话状态同步给引擎，开启或关闭口型与声波光环

  useEffect(() => {
    engineRef.current?.setSpeaking(speaking)
  }, [speaking])

  // 语音边界 effect：计数变化时注入一次张口脉冲，使口型贴合真实发音

  useEffect(() => {
    if (speakPulse === 0) return
    engineRef.current?.pulseMouth()
  }, [speakPulse])

  return (
    <div ref={hostRef} className="live2d-host" aria-label="Live2D 数字人渲染区（ticket 07）">
      <canvas id="live2d-canvas" />
    </div>
  )
}
