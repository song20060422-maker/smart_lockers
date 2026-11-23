/** @ticket 07-数字人Coze：编排 Hook——串联 WebSocket 回复 → TTS 播报 → Live2D 表情动作联动 */
import { useCallback, useEffect, useRef } from 'react'
import { useDispatch } from 'react-redux'
import { WebSocketClient } from '@/ws/WebSocketClient'
import { SpeechEngine } from '@/tts/speech'
import { CONFIG, getRole } from '@/config'
import { appendMessage, setConnected, setThinking, setAction, setSpeaking, bumpSpeakPulse, setError, clearConversation } from '@/store'
import type { ChatMessage, AgentReplyPayload } from '@/types'
import type { AppDispatch } from '@/store'

/** 情绪 → Live2D 动作映射（工单：智能体回复触发点头/微笑/摇头等） */
// 定义函数/方法，封装可复用的处理逻辑：const emotionAction = (emotion: string): AgentReplyPayload['action'] => {

const emotionAction = (emotion: string): AgentReplyPayload['action'] => {
  // 声明变量 map: Record<string, 'nod' | 'smile' | 'shake' | 'blink'> = {，用于保存数据或函数引用

  const map: Record<string, 'nod' | 'smile' | 'shake' | 'blink'> = {
    // 定义对象属性 happy，取值为 'smile', positive: 'nod', sad: 'shake', neutral: 'blink'

    happy: 'smile', positive: 'nod', sad: 'shake', neutral: 'blink',
  // 当前代码块结束

  }
  // 返回 { action: map[emotion] ?? 'blink', emotion, duration_ms: 1200 } 给调用方

  return { action: map[emotion] ?? 'blink', emotion, duration_ms: 1200 }
// 当前代码块结束

}

// 导出 function useDigitalHuman() {，供其他模块引用

export function useDigitalHuman(role: string | null) {
  // 声明变量 dispatch = useDispatch<AppDispatch>()，用于保存数据或函数引用

  const dispatch = useDispatch<AppDispatch>()
  // 声明变量 wsRef = useRef<WebSocketClient | null>(null)，用于保存数据或函数引用

  const wsRef = useRef<WebSocketClient | null>(null)
  // 声明变量 speechRef = useRef<SpeechEngine | null>(null)，用于保存数据或函数引用

  const speechRef = useRef<SpeechEngine | null>(null)

  // 会话代次 ref：账号切换/重连时 +1，用于丢弃旧连接迟到的欢迎语与回复，杜绝串号 @ticket 07
  // 声明变量 sessionRef，保存运行过程中的自增会话编号

  const sessionRef = useRef(0)

  // 调用 useEffect(() => {，触发对应功能

  useEffect(() => {
    // 会话代次自增：本次 effect 建立的连接绑定唯一编号，旧连接回调因编号过期被忽略 @ticket 07
    const session = ++sessionRef.current
    // 账号角色变化（含首次登录/切换账号）先清空上一账号的对话与瞬时状态，保证角色内容完全隔离 @ticket 07
    dispatch(clearConversation())
    // 未登录或无业务角色时不建立连接，等待账号登录判定角色后再接入 @ticket 07
    if (!role) return
    // 声明变量 speech = new SpeechEngine(getRole().voice, { onRecognized: (text) => send(text, 'voice') })，用于保存数据或函数引用

    const speech = new SpeechEngine(getRole().voice, {
      // 语音识别结果回调：把识别文本作为用户消息发送

      onRecognized: (text) => send(text, 'voice'),
      // 播报开始回调：开启数字人说话口型状态

      onSpeakStart: () => dispatch(setSpeaking(true)),
      // 语音边界回调：递增计数以触发一次张口脉冲，贴合真实发音

      onSpeakBoundary: () => dispatch(bumpSpeakPulse()),
      // 播报结束回调：关闭数字人说话口型状态

      onSpeakEnd: () => dispatch(setSpeaking(false)),
    })
    // 赋值更新：speechRef.current = speech

    speechRef.current = speech

    // 声明变量 ws = new WebSocketClient，接入账号判定的业务角色（由父组件传入）@ticket 07

    const ws = new WebSocketClient(CONFIG.wsUrl, CONFIG.token, role, {
      // 定义对象属性 onStatus，仅处理当前会话代次的状态，避免旧连接迟到事件改状态 @ticket 07

      onStatus: (c) => {
        // 丢弃过期连接的状态回调 @ticket 07

        if (sessionRef.current !== session) return
        // 同步连接态到 Redux，驱动 UI 连接指示 @ticket 07
        dispatch(setConnected(c))
      },
      // 定义对象属性 onReply，仅处理当前会话代次的消息，杜绝切换账号后旧角色欢迎语/回复串入 @ticket 07

      onReply: (p: AgentReplyPayload) => {
        // 丢弃过期连接的欢迎语与回复，这是串号与重复欢迎语的关键拦截点 @ticket 07

        if (sessionRef.current !== session) return
        // 调用 dispatch(setThinking(false))，触发对应功能

        dispatch(setThinking(false))
        // 声明变量 msg: ChatMessage = { role: 'assistant', content: p.text, type: 'text', timestamp: Date.now(), meta: { emotion: p.emotion, scene: p.scene } }，用于保存数据或函数引用

        const msg: ChatMessage = { role: 'assistant', content: p.text, type: 'text', timestamp: Date.now(), meta: { emotion: p.emotion, scene: p.scene } }
        // 调用 dispatch(appendMessage(msg))，触发对应功能

        dispatch(appendMessage(msg))
        // 执行 speech.speak(p.text)                       // 语音播报，实现对应的前端逻辑

        speech.speak(p.text)                       // 语音播报
        // 调用 dispatch(setAction(p.action ?? emotionAction(p.emotion) ?? null)) // 表情/动画联动，触发对应功能

        dispatch(setAction(p.action ?? emotionAction(p.emotion) ?? null)) // 表情/动画联动
      // 定义结束

      },
    // 执行 })，实现对应的前端逻辑

    })
    // 执行 ws.connect()，实现对应的前端逻辑

    ws.connect()
    // 赋值更新：wsRef.current = ws

    wsRef.current = ws
    // 返回 () => { ws.close(); speech.stop() } 给调用方，并使当前会话代次立即失效以丢弃旧连接迟到回调 @ticket 07

    return () => { sessionRef.current += 1; ws.close(); speech.stop() }
  // role 变化时重连以切换到对应角色内容 @ticket 07
  }, [dispatch, role])

  // 声明变量 send = useCallback((content: string, type: 'text' | 'voice' = 'text') => {，用于保存数据或函数引用

  const send = useCallback((content: string, type: 'text' | 'voice' = 'text') => {
    // 条件判断：if (!content.trim()) return 时执行下面逻辑

    if (!content.trim()) return
    // 声明变量 ok = wsRef.current?.send(content, type)，用于保存数据或函数引用

    const ok = wsRef.current?.send(content, type)
    // 条件判断：if (!ok) { dispatch(setError('连接未就绪，无法发送')); return } 时执行下面逻辑

    if (!ok) { dispatch(setError('连接未就绪，无法发送')); return }
    // 调用 dispatch(appendMessage({ role: 'user', content, type, timestamp: Date.now() }))，触发对应功能

    dispatch(appendMessage({ role: 'user', content, type, timestamp: Date.now() }))
    // 调用 dispatch(setThinking(true))，触发对应功能

    dispatch(setThinking(true))
    // 调用 dispatch(setError(null))，触发对应功能

    dispatch(setError(null))
  // 执行 }, [dispatch])，实现对应的前端逻辑

  }, [dispatch])

  // 声明变量 listen = useCallback(() => speechRef.current?.startListen(), [])，用于保存数据或函数引用

  const listen = useCallback(() => speechRef.current?.startListen(), [])

  // 返回 { send, listen } 给调用方

  return { send, listen }
// 当前代码块结束

}
