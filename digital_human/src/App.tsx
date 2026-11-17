/** @ticket 07-数字人Coze：数字人交互主界面（输入 → 智能体 → 语音+表情联动） */
import { useEffect, useRef, useState } from 'react'
import { useSelector } from 'react-redux'
import Live2DViewer from '@/live2d/Live2DViewer'
import { useDigitalHuman } from '@/hooks/useDigitalHuman'
import { CONFIG, ROLES, setRole, getRole } from '@/config'
import type { RootState } from '@/store'
import './styles.css'

export default function App() {
  const { messages, connected, thinking, action, error } = useSelector((s: RootState) => s.hub)
  const { send, listen } = useDigitalHuman()
  const [input, setInput] = useState('')
  const [modelPath, setModelPath] = useState(CONFIG.modelPath)
  const endRef = useRef<HTMLDivElement>(null)

  useEffect(() => { endRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  const switchRole = (id: string) => { setModelPath(setRole(id).modelPath) }

  return (
    <div className="app">
      <header className="topbar">
        <span className="title">智能快寄柜数字人 <small>@ticket 07</small></span>
        <span className={`conn ${connected ? 'on' : 'off'}`}>{connected ? '已连接' : '连接中…'}</span>
        <select defaultValue={getRole().id} onChange={e => switchRole(e.target.value)}>
          {Object.values(ROLES).map(r => <option key={r.id} value={r.id}>{r.name}</option>)}
        </select>
      </header>

      <div className="stage">
        <Live2DViewer modelPath={modelPath} action={action} />
        <div className="chat">
          <div className="messages">
            {messages.map((m, i) => (
              <div key={i} className={`bubble ${m.role}`}>
                <span className="who">{m.role === 'user' ? '我' : '数字人'}</span>
                <span>{m.content}</span>
              </div>
            ))}
            {thinking && <div className="bubble assistant">思考中…</div>}
            {error && <div className="error">{error}</div>}
            <div ref={endRef} />
          </div>
          <div className="composer">
            <button onClick={listen} title="语音输入">🎤</button>
            <input value={input} onChange={e => setInput(e.target.value)}
              onKeyDown={e => { if (e.key === 'Enter') { send(input); setInput('') } }}
              placeholder="输入问题，回车发送" />
            <button onClick={() => { send(input); setInput('') }}>发送</button>
          </div>
        </div>
      </div>
    </div>
  )
}
