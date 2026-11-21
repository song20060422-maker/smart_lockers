// @ticket 06-个人中心：头像/服务/卡券接口封装（登录登出见 auth.ts）
import { get, post } from './request'
import type { AvatarResult, MyServicesResult, Coupon, OrderSummaryResult } from '@/types'

// @ticket 06-个人中心：POST /api/v1/user/avatar（鉴权）
// 导出 function updateAvatar(payload: AvatarPayload): Promise<AvatarResult> {，供其他模块引用

export function updateAvatar(payload: AvatarPayload): Promise<AvatarResult> {
  // 返回 post('/user/avatar', payload) 给调用方

  return post('/user/avatar', payload)
// 当前代码块结束

}

// @ticket 06-个人中心：GET /api/v1/user/services（鉴权）
// 导出 function myServices(): Promise<MyServicesResult> {，供其他模块引用

export function myServices(): Promise<MyServicesResult> {
  // 返回 get('/user/services') 给调用方

  return get('/user/services')
// 当前代码块结束

}

// @ticket 06-个人中心：GET /api/v1/user/coupons（鉴权）
// 导出 function myCoupons(): Promise<Coupon[]> {，供其他模块引用

export function myCoupons(): Promise<Coupon[]> {
  // 返回 get('/user/coupons') 给调用方

  return get('/user/coupons')
// 当前代码块结束

}

// @ticket 06-个人中心：POST /api/v1/user/coupons/:code/validate（鉴权）
// 导出 function validateCoupon(code: string): Promise<Coupon> {，供其他模块引用

export function validateCoupon(code: string): Promise<Coupon> {
  // 返回 post(`/user/coupons/${code}/validate`, {}) 给调用方

  return post(`/user/coupons/${code}/validate`, {})
// 当前代码块结束

}

// 定义接口，约束对象的结构：interface AvatarPayload {

interface AvatarPayload {
  // 定义对象属性 url，取值为 string

  url: string
  // 执行 content_type?: string，实现对应的前端逻辑

  content_type?: string
  // 执行 size?: number，实现对应的前端逻辑

  size?: number
// 当前代码块结束

}


// @ticket 06-个人中心：GET /api/v1/user/order-summary（鉴权，订单汇总）
// 导出 function orderSummary(): Promise<OrderSummaryResult> {，供其他模块引用

export function orderSummary(): Promise<OrderSummaryResult> {
  // 返回 get('/user/order-summary') 给调用方

  return get('/user/order-summary')
// 当前代码块结束

}
