/** @ticket 07-数字人Coze：动漫全身数字人渲染组件（接口与 Live2DViewer/PhotorealViewer 对齐） */
import { useEffect, useRef } from 'react'
import type { Live2DAction } from '@/types'
import { AnimeAvatarEngine } from '@/anime/AnimeAvatarEngine'

// 定义组件 Props：与其他渲染器保持一致的驱动接口
interface Props {
  // poseDir：姿势集目录（角色切换预留，默认 /anime）
  poseDir?: string
  // action：智能体下发的动作/情绪指令
  action: Live2DAction | null
  // speaking：是否正在语音播报
  speaking: boolean
  // speakPulse：语音边界计数，驱动全身弹跳与小动作
  speakPulse: number
  // onError：错误回调
  onError?: (e: string) => void
}

// 导出默认组件 AnimeViewer：动漫全身数字人舞台
export default function AnimeViewer({ poseDir = '/anime', action, speaking, speakPulse, onError }: Props) {
  // 声明 hostRef：舞台宿主
  const hostRef = useRef<HTMLDivElement>(null)
  // 声明 engineRef：引擎实例复用
  const engineRef = useRef<AnimeAvatarEngine | null>(null)

  // 挂载 effect：创建引擎并启动动画循环，卸载时释放
  useEffect(() => {
    // 读取宿主节点
    const host = hostRef.current
    // 宿主缺失时跳过
    if (!host) return
    // 创建动漫全身引擎
    const engine = new AnimeAvatarEngine(host, poseDir)
    // 保存实例
    engineRef.current = engine
    // 启动动画
    engine.start()
    // 返回清理函数
    return () => { engine.stop(); engineRef.current = null }
  // poseDir 变化时重建引擎（切换角色姿势集）
  }, [poseDir])

  // 姿势目录同步 effect
  useEffect(() => {
    // 目录为空时报错提示
    if (!poseDir) onError?.('未配置动漫数字人姿势集')
  }, [poseDir, onError])

  // 动作驱动 effect：nod/shake/wave/point/happy 等映射姿势切换
  useEffect(() => {
    // 无动作或引擎未就绪时跳过
    if (!action || !engineRef.current) return
    // 调用引擎动作接口
    engineRef.current.setAction(action)
  }, [action])

  // 播报状态 effect：同步声波条与轻弹循环
  useEffect(() => {
    // 调用引擎播报状态
    engineRef.current?.setSpeaking(speaking)
  }, [speaking])

  // 语音边界 effect：注入全身弹跳脉冲（偶发切换小动作）
  useEffect(() => {
    // 计数为 0 时跳过（初始值）
    if (speakPulse === 0) return
    // 触发弹跳脉冲
    engineRef.current?.pulseMouth()
  }, [speakPulse])

  // 渲染动漫全身数字人舞台
  return <div ref={hostRef} className="an-stage" aria-label="动漫数字人渲染区（ticket 07）" />
// 当前代码块结束
}
