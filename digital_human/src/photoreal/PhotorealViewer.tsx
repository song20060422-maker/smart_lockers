/** @ticket 07-数字人Coze：写实真人数字人渲染组件（接口与 Live2DViewer 对齐，可按配置切换） */
import { useEffect, useRef } from 'react'
import type { Live2DAction } from '@/types'
import { PhotorealAvatarEngine } from '@/photoreal/PhotorealAvatarEngine'

// 定义组件 Props：与 Live2DViewer 保持一致的驱动接口
interface Props {
  // portraitSrc：当前角色写实人像图地址
  portraitSrc: string
  // action：智能体下发的动作/情绪指令
  action: Live2DAction | null
  // speaking：是否正在语音播报
  speaking: boolean
  // speakPulse：语音边界计数，驱动口型脉冲
  speakPulse: number
  // onError：错误回调
  onError?: (e: string) => void
}

// 导出默认组件 PhotorealViewer：写实数字人舞台
export default function PhotorealViewer({ portraitSrc, action, speaking, speakPulse, onError }: Props) {
  // 声明 hostRef：舞台宿主
  const hostRef = useRef<HTMLDivElement>(null)
  // 声明 engineRef：引擎实例复用
  const engineRef = useRef<PhotorealAvatarEngine | null>(null)

  // 挂载 effect：创建引擎并启动动画循环，卸载时释放
  useEffect(() => {
    // 读取宿主节点
    const host = hostRef.current
    // 宿主缺失时跳过
    if (!host) return
    // 创建写实引擎
    const engine = new PhotorealAvatarEngine(host, portraitSrc)
    // 保存实例
    engineRef.current = engine
    // 启动动画
    engine.start()
    // 窗口尺寸变化时重建图层定位（区域标定按像素换算）
    const onResize = () => { engine.stop(); engine.start() }
    // 监听窗口尺寸
    window.addEventListener('resize', onResize)
    // 返回清理函数
    return () => { engine.stop(); window.removeEventListener('resize', onResize); engineRef.current = null }
  // portraitSrc 变化时重建引擎以应用新人像
  }, [portraitSrc])

  // 人像地址同步 effect（引擎随依赖重建，无需手动切换）
  useEffect(() => {
    // 地址为空时报错提示
    if (!portraitSrc) onError?.('未配置写实数字人形象')
  }, [portraitSrc, onError])

  // 动作驱动 effect：nod/shake 等交给引擎
  useEffect(() => {
    // 无动作或引擎未就绪时跳过
    if (!action || !engineRef.current) return
    // 调用引擎动作接口
    engineRef.current.setAction(action)
  }, [action])

  // 播报状态 effect：同步口型开关与光环
  useEffect(() => {
    // 调用引擎播报状态
    engineRef.current?.setSpeaking(speaking)
  }, [speaking])

  // 语音边界 effect：注入一次张口脉冲
  useEffect(() => {
    // 计数为 0 时跳过（初始值）
    if (speakPulse === 0) return
    // 触发口型脉冲
    engineRef.current?.pulseMouth()
  }, [speakPulse])

  // 渲染写实数字人舞台
  return (
    <div ref={hostRef} className="ph-stage" aria-label="写实数字人渲染区（ticket 07）" />
  )
// 当前代码块结束
}
