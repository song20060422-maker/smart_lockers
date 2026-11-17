/** @ticket 07-数字人Coze：编排 Hook——串联 WebSocket 回复 → TTS 播报 → Live2D 表情动作联动 */
import { useCallback, useEffect, useRef } from 'react'
import { useDispatch } from 'react-redux'
import { WebSocketClient } from '@/ws/WebSocketClient'
import { SpeechEngine } from '@/tts/speech'
import { CONFIG, getRole } from '@/config'
import { appendMessage, setConnected, setThinking, setAction, setError } from '@/store'
import type { ChatMessage, AgentReplyPayload } from '@/types'
import type { AppDispatch } from '@/store'

/** 情绪 → Live2D 动作映射（工单：智能体回复触发点头/微笑/摇头等） */
const emotionAction = (emotion: string): AgentReplyPayload['action'] => {
  const map: Record<string, 'nod' | 'smile' | 'shake' | 'blink'> = {
    happy: 'smile', positive: 'nod', sad: 'shake', neutral: 'blink',
  }
  return { action: map[emotion] ?? 'blink', emotion, duration_ms: 1200 }
}

export function useDigitalHuman() {
  const dispatch = useDispatch<AppDispatch>()
  const wsRef = useRef<WebSocketClient | null>(null)
  const speechRef = useRef<SpeechEngine | null>(null)

  useEffect(() => {
    const speech = new SpeechEngine(getRole().voice, { onRecognized: (text) => send(text, 'voice') })
    speechRef.current = speech

    const ws = new WebSocketClient(CONFIG.wsUrl, CONFIG.token, {
      onStatus: (c) => dispatch(setConnected(c)),
      onReply: (p: AgentReplyPayload) => {
        dispatch(setThinking(false))
        const msg: ChatMessage = { role: 'assistant', content: p.text, type: 'text', timestamp: Date.now(), meta: { emotion: p.emotion, scene: p.scene } }
        dispatch(appendMessage(msg))
        speech.speak(p.text)                       // 语音播报
        dispatch(setAction(p.action ?? emotionAction(p.emotion) ?? null)) // 表情/动画联动
      },
    })
    ws.connect()
    wsRef.current = ws
    return () => { ws.close(); speech.stop() }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dispatch])

  const send = useCallback((content: string, type: 'text' | 'voice' = 'text') => {
    if (!content.trim()) return
    const ok = wsRef.current?.send(content, type)
    if (!ok) { dispatch(setError('连接未就绪，无法发送')); return }
    dispatch(appendMessage({ role: 'user', content, type, timestamp: Date.now() }))
    dispatch(setThinking(true))
    dispatch(setError(null))
  }, [dispatch])

  const listen = useCallback(() => speechRef.current?.startListen(), [])

  return { send, listen }
}
