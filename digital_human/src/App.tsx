/** @ticket 07-数字人Coze：数字人交互主界面（输入 → 智能体 → 语音+表情联动）
 * 登录门：先账号密码登录（后端真实查库校验，成功后签发 token），通过后再按场景角色进入对应数字人页面 @ticket 07
 * 动漫全身模式默认启用（anime：从头到脚全身展示 + 姿势动作库），VITE_AVATAR_MODE 可切 photoreal/cartoon
 */
import { useEffect, useRef, useState } from 'react'
import { useSelector } from 'react-redux'
import Live2DViewer from '@/live2d/Live2DViewer'
import PhotorealViewer from '@/photoreal/PhotorealViewer'
import AnimeViewer from '@/anime/AnimeViewer'
// 角色选择页 LoginPage 已移除：账号登录后按账号角色直接进入对话，不再独立选角色 @ticket 07

// 账号密码登录页：未通过账号校验时优先展示 @ticket 07
// 导入默认组件 AccountLoginPage，用于账号密码登录门

import AccountLoginPage from '@/pages/AccountLoginPage'
import type { DhRole } from '@/pages/types'
import { useDigitalHuman } from '@/hooks/useDigitalHuman'
import { CONFIG, setRole, getRole } from '@/config'
import type { RootState } from '@/store'
import './styles.css'

// 登录会话持久化键名（localStorage）
// 声明变量 ROLE_KEY，用于保存运行过程中的数据

const ROLE_KEY = 'dh_login_role'

// 账号登录态持久化键名（localStorage）：'1' 表示账号密码已校验通过
// 声明变量 ACCOUNT_KEY，用于保存运行过程中的数据

const ACCOUNT_KEY = 'dh_account_authed'

// 后端签发的登录令牌键名（localStorage）：真实查库登录成功后保存，供后续需要鉴权的接口使用 @ticket 07-数字人Coze
// 声明变量 TOKEN_KEY，用于保存运行过程中的数据

const TOKEN_KEY = 'dh_account_token'

// 场景角色 → 数字人形象映射：寄件引导=品牌代言人，取件指引/智能客服=虚拟讲解员 @ticket 07
// 声明常量 ROLE_AVATAR，用于保存运行过程中的数据

const ROLE_AVATAR: Record<DhRole, string> = {
  // 寄件引导员 → 品牌代言人
  sender: 'presenter',
  // 取件指引员 → 虚拟讲解员
  receiver: 'guide',
  // 智能客服专员 → 品牌代言人
  service: 'presenter',
// 当前代码块结束

}

// 场景角色中文名（顶栏徽标展示）
// 声明常量 ROLE_LABEL，用于保存运行过程中的数据

const ROLE_LABEL: Record<DhRole, string> = {
  // 寄件引导员
  sender: '寄件引导员',
  // 取件指引员
  receiver: '取件指引员',
  // 智能客服专员
  service: '智能客服专员',
// 当前代码块结束

}

// 后端返回的账号角色 → 业务角色白名单归一化：operator 等旧账号兜底为智能客服，避免脏值 @ticket 07
// 定义函数 normalizeRole，用于把任意后端角色字符串收敛为合法 DhRole

const normalizeRole = (raw: string | undefined | null): DhRole => {
  // 仅三业务角色合法，其余（含 operator/空值）统一兜底为 service
  if (raw === 'sender' || raw === 'receiver' || raw === 'service') return raw
  // 兜底角色
  return 'service'
}

export default function App() {
  const { messages, connected, thinking, action, speaking, speakPulse, error } = useSelector((s: RootState) => s.hub)
  // 登录角色：初始从 localStorage 恢复会话（null 表示未登录）。须在联网 hook 之前声明 @ticket 07
  const [loginRole, setLoginRole] = useState<DhRole | null>(() => {
    // 读取持久化会话并做合法性归一化（兼容旧 operator 会话）
    const saved = localStorage.getItem(ROLE_KEY)
    return saved ? normalizeRole(saved) : null
  })
  // 传入账号判定的业务角色：未登录时为 null，联网 hook 内不建立连接 @ticket 07
  const { send, listen } = useDigitalHuman(loginRole)
  const [input, setInput] = useState('')
  const [role, setRoleState] = useState(getRole())
  const endRef = useRef<HTMLDivElement>(null)

  // 账号登录态：初始从 localStorage 恢复（false 表示未通过账号密码校验，先展示账号密码登录页）
  const [accountAuthed, setAccountAuthed] = useState<boolean>(() => {
    // 读取账号密码校验标记
    return localStorage.getItem(ACCOUNT_KEY) === '1'
  })

  useEffect(() => { endRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  // 账号密码校验通过：直接按账号判定角色进入对应对话，跳过独立角色选择页 @ticket 07
  const handleAccountLogin = (token: string, rawRole: string) => {
    // 后端按账号返回业务角色，归一化后作为本次会话角色
    const r = normalizeRole(rawRole)
    // 写入账号登录态标记
    localStorage.setItem(ACCOUNT_KEY, '1')
    // 保存后端签发的登录令牌，供后续接口鉴权使用
    localStorage.setItem(TOKEN_KEY, token)
    // 持久化账号判定角色会话
    localStorage.setItem(ROLE_KEY, r)
    // 更新账号登录态状态
    setAccountAuthed(true)
    // 直接进入账号对应角色（无独立选择页）
    setLoginRole(r)
    // 按映射切换数字人角色配置（寄件引导=代言人，取件指引=讲解员，客服=代言人）
    setRoleState(setRole(ROLE_AVATAR[r]))
  }

  // 退出登录：清空账号登录态与角色会话，回到账号密码登录页
  const handleLogout = () => {
    // 清除账号登录态标记
    localStorage.removeItem(ACCOUNT_KEY)
    // 清除后端签发的登录令牌
    localStorage.removeItem(TOKEN_KEY)
    // 清除持久化角色会话
    localStorage.removeItem(ROLE_KEY)
    // 重置账号登录态
    setAccountAuthed(false)
    // 重置登录状态
    setLoginRole(null)
  }

  // 未通过账号密码校验：展示账号密码登录页（不渲染对话界面）
  if (!accountAuthed) {
    // 渲染账号密码登录页
    return <AccountLoginPage onSuccess={handleAccountLogin} />
  }

  // 账号已登录但角色缺失（兼容旧会话/异常态）：兜底按智能客服进入，绝不再回到角色选择页 @ticket 07
  const activeRole: DhRole = loginRole ?? 'service'

  return (
    <div className="app">
      <header className="topbar">
        <span className="title">智能快寄柜数字人</span>
        {/* 当前登录场景角色徽标（登录时决定，页面内不再提供形象切换）@ticket 07 */}
        <span className="role-badge">{ROLE_LABEL[activeRole]}</span>
        {/* 当前数字人形象（随场景角色自动绑定，只读展示） */}
        <span className="avatar-label">{role.name}</span>
        <span className={`conn ${connected ? 'on' : 'off'}`}>{connected ? '已连接' : '连接中…'}</span>
        {/* 退出登录按钮 */}
        <button className="logout-btn" onClick={handleLogout}>退出</button>
      </header>

      <div className="stage">
        {/* 动漫全身模式（默认）：从头到脚全身展示，姿势动作库驱动丰富动作 @ticket 07 */}
        {CONFIG.avatarMode === 'anime' ? (
          <AnimeViewer action={action} speaking={speaking} speakPulse={speakPulse} />
        ) : CONFIG.avatarMode === 'photoreal' ? (
          /* 写实真人模式（保留可切换）：照片级分层绑骨渲染 */
          <PhotorealViewer portraitSrc={role.portrait} action={action} speaking={speaking} speakPulse={speakPulse} />
        ) : (
          /* 卡通模式（兼容保留）：Live2D 模型渲染 */
          <Live2DViewer modelPath={role.modelPath} action={action} speaking={speaking} speakPulse={speakPulse} />
        )}
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
