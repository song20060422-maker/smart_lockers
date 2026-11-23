/** @ticket 07-数字人Coze：账号密码登录页——调用后端真实查库接口校验（不再前端写死演示账号），成功后按账号角色直接进入对应数字人对话页（无独立角色选择页） */

import { useState } from 'react'

import type { FormEvent } from 'react'

// 后端登录接口路径：Vite 代理 /api → http://localhost:18080 @ticket 07-数字人Coze
// 声明常量 LOGIN_API，用于保存运行过程中的数据

const LOGIN_API = '/api/v1/digital-human/login'

// 初始种子账号列表（数据库真实存在的账号，仅供联调填充，校验完全由后端查库完成）@ticket 07
// 三个业务角色账号登录后直接进入对应角色，无需再选角色
const SEED_ACCOUNTS: { label: string; account: string; password: string }[] = [
  { label: '寄件引导员', account: 'sender', password: '123456' },
  { label: '取件指引员', account: 'receiver', password: '123456' },
  { label: '智能客服专员', account: 'service', password: '123456' },
]

// 定义组件 Props：onSuccess 为校验通过后的回调，回传后端签发的 token 与账号判定角色
// 声明接口 Props，用于保存运行过程中的数据类型

interface Props {
  // onSuccess：校验通过后通知父组件直接进入对应角色页面，回传 token 与 role（账号即角色，跳过选择页）
  onSuccess: (token: string, role: string) => void
// 当前代码块结束

}

// 定义后端统一响应结构：与 pkg.Response 字段契约一致 @ticket 07-数字人Coze
// 声明接口 LoginResp，用于保存运行过程中的数据类型

interface LoginResp {
  // code：业务码，0 表示成功
  code: number
  // message：业务提示文案
  message: string
  // data：成功时携带 token/username/role
  data?: { token: string; username: string; role: string }
}

// 导出默认组件 AccountLoginPage：账号密码登录页
// 声明函数 AccountLoginPage，用于渲染账号密码登录界面

export default function AccountLoginPage({ onSuccess }: Props) {
  // 声明状态 account，保存账号输入
  const [account, setAccount] = useState('')

  // 声明状态 password，保存密码输入
  const [password, setPassword] = useState('')

  // 声明状态 error，保存校验失败提示
  const [error, setError] = useState('')

  // 声明状态 loading，防止重复提交
  const [loading, setLoading] = useState(false)

  // 处理表单提交：调用后端真实查库登录接口
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    // 阻止浏览器默认提交行为
    e.preventDefault()

    // 校验账号非空（去除首尾空白）
    if (!account.trim()) {
      // 设置账号为空提示
      setError('请输入账号')
      // 终止提交
      return
    }

    // 校验密码非空
    if (!password) {
      // 设置密码为空提示
      setError('请输入密码')
      // 终止提交
      return
    }

    // 防止重复提交
    if (loading) return

    // 标记提交中
    setLoading(true)
    // 清空上一次的错误提示
    setError('')

    try {
      // 调用后端账号密码查库校验接口
      const resp = await fetch(LOGIN_API, {
        // POST 方法
        method: 'POST',
        // JSON 请求头
        headers: { 'Content-Type': 'application/json' },
        // 字段契约：username/password
        body: JSON.stringify({ username: account.trim(), password }),
      })
      // 解析统一响应结构
      const json = (await resp.json()) as LoginResp
      // 条件判断：HTTP 与业务码均成功时直接进入账号对应角色页面（跳过角色选择）
      if (resp.ok && json.code === 0 && json.data?.token) {
        // 校验通过：回传 token 与账号判定角色，由父组件直接进入对应角色
        onSuccess(json.data.token, json.data.role)
      } else {
        // 校验失败：使用后端友好文案，不暴露原始报错；缺失时兜底
        setError(json.message || '登录失败，请稍后重试')
      }
    } catch {
      // 网络/服务不可用：友好提示，不暴露系统内部错误文本
      setError('服务暂不可用，请稍后重试')
    } finally {
      // 恢复提交状态
      setLoading(false)
    }
  }

  // 渲染账号密码登录页
  return (
    // 登录页根容器
    <div className="dh-account-login">
      {/* 登录卡片 */}
      <div className="dh-account-card">
        {/* 品牌标题 */}
        <h1 className="dh-account-title">智能快寄柜 · 数字人</h1>
        {/* 副标题 */}
        <p className="dh-account-sub">请使用数据库账号密码登录</p>
        {/* 登录表单：回车即提交 */}
        <form className="dh-account-form" onSubmit={handleSubmit}>
          {/* 账号标签 */}
          <label className="dh-account-label" htmlFor="dh-account">账号</label>
          {/* 账号输入框 */}
          <input
            id="dh-account"
            className="dh-account-input"
            value={account}
            onChange={e => setAccount(e.target.value)}
            placeholder="请输入账号"
            autoComplete="username"
          />
          {/* 密码标签 */}
          <label className="dh-account-label" htmlFor="dh-password">密码</label>
          {/* 密码输入框 */}
          <input
            id="dh-password"
            className="dh-account-input"
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            placeholder="请输入密码"
            autoComplete="current-password"
          />
          {/* 校验失败提示 */}
          {error && <p className="dh-account-error">{error}</p>}
          {/* 登录提交按钮 */}
          <button className="dh-account-submit" type="submit" disabled={loading}>
            {loading ? '登录中…' : '登录'}
          </button>
        </form>
        {/* 数据库初始账号快速填充：三角色账号，登录后直接进入对应角色，校验仍由后端查库完成 @ticket 07 */}
        <div className="dh-account-seeds">
          {/* 遍历三角色种子账号，各自提供一键填充 */}
          {SEED_ACCOUNTS.map(s => (
            <button
              key={s.account}
              className="dh-account-demo"
              type="button"
              onClick={() => { setAccount(s.account); setPassword(s.password); setError('') }}
            >
              {s.label}：{s.account} / {s.password}（点击填充）
            </button>
          ))}
        </div>
      </div>
    </div>
  )
// 当前代码块结束

}
