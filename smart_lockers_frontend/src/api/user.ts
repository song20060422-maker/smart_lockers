// @ticket 06-个人中心：头像/服务/卡券接口封装（登录登出见 auth.ts）
import { get, post } from './request'
import type { AvatarResult, MyServicesResult, Coupon, OrderSummaryResult } from '@/types'

// @ticket 06-个人中心：POST /api/v1/user/avatar（鉴权）
export function updateAvatar(payload: AvatarPayload): Promise<AvatarResult> {
  return post('/user/avatar', payload)
}

// @ticket 06-个人中心：GET /api/v1/user/services（鉴权）
export function myServices(): Promise<MyServicesResult> {
  return get('/user/services')
}

// @ticket 06-个人中心：GET /api/v1/user/coupons（鉴权）
export function myCoupons(): Promise<Coupon[]> {
  return get('/user/coupons')
}

// @ticket 06-个人中心：POST /api/v1/user/coupons/:code/validate（鉴权）
export function validateCoupon(code: string): Promise<Coupon> {
  return post(`/user/coupons/${code}/validate`, {})
}

interface AvatarPayload {
  url: string
  content_type?: string
  size?: number
}


// @ticket 06-个人中心：GET /api/v1/user/order-summary（鉴权，订单汇总）
export function orderSummary(): Promise<OrderSummaryResult> {
  return get('/user/order-summary')
}
