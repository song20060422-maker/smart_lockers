// @ticket 06-个人中心：登录/退出接口封装
import { request } from './request'
import type { LoginResult, OkResult, DemoLoginResult } from '@/types'

// @ticket 06-个人中心：POST /api/v1/auth/login，返回 token
// 导出 function login(userId: number): Promise<LoginResult> {，供其他模块引用

export function login(userId: number): Promise<LoginResult> {
  // 返回 request({ url: '/auth/login', method: 'POST', data: { user_id: userId } }) 给调用方

  return request({ url: '/auth/login', method: 'POST', data: { user_id: userId } })
// 当前代码块结束

}

// 修复偏差4：POST /api/v1/auth/demo-login，手机号+验证码换取后端真实 JWT @ticket 00-project
// 导出 function demoLogin：调后端演示登录接口，返回 token 与 user_id，禁止前端本地生成 token

export function demoLogin(phone: string, code: string): Promise<DemoLoginResult> {
  // 返回 request：向后端发起真实 HTTP 登录请求

  return request({ url: '/auth/demo-login', method: 'POST', data: { phone, code } })
// 当前代码块结束

}

// @ticket 06-个人中心：POST /api/v1/auth/logout（令牌即时失效、幂等）
// 导出 function logout(): Promise<OkResult> {，供其他模块引用

export function logout(): Promise<OkResult> {
  // 返回 request({ 给调用方

  return request({
    // 定义对象属性 url，取值为 '/auth/logout'

    url: '/auth/logout',
    // 定义对象属性 method，取值为 'POST'

    method: 'POST',
    // 定义对象属性 idempotencyKey，取值为 `logout-${Date.now()}`

    idempotencyKey: `logout-${Date.now()}`,
  // 执行 })，实现对应的前端逻辑

  })
// 当前代码块结束

}
