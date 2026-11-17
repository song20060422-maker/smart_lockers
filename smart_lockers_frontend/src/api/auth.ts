// @ticket 06-个人中心：登录/退出接口封装
import { request } from './request'
import type { LoginResult, OkResult } from '@/types'

// @ticket 06-个人中心：POST /api/v1/auth/login，返回 token
export function login(userId: number): Promise<LoginResult> {
  return request({ url: '/auth/login', method: 'POST', data: { user_id: userId } })
}

// @ticket 06-个人中心：POST /api/v1/auth/logout（令牌即时失效、幂等）
export function logout(): Promise<OkResult> {
  return request({
    url: '/auth/logout',
    method: 'POST',
    idempotencyKey: `logout-${Date.now()}`,
  })
}
