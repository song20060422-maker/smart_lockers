// @ticket 05-订单查看与取出：列表/详情/开柜/补交/完成接口封装
import { get, post } from './request'
import type {
  OrderListResult, OrderDetail, OngoingResult, OpenAuthResult, MakeupResult, OkResult,
} from '@/types'

// @ticket 05-订单查看与取出：GET /api/v1/orders（游标分页）
// 导出 function orderList(params: { type?: string; cursor?: string; limit?: number }): Promise<OrderListResult> {，供其他模块引用

export function orderList(params: { type?: string; cursor?: string; limit?: number }): Promise<OrderListResult> {
  // 返回 get('/orders', params) 给调用方

  return get('/orders', params)
// 当前代码块结束

}

// @ticket 05-订单查看与取出：GET /api/v1/orders/:order_no
// 导出 function orderDetail(orderNo: string): Promise<OrderDetail> {，供其他模块引用

export function orderDetail(orderNo: string): Promise<OrderDetail> {
  // 返回 get(`/orders/${orderNo}`) 给调用方

  return get(`/orders/${orderNo}`)
// 当前代码块结束

}

// @ticket 05-订单查看与取出：GET /api/v1/orders/summary/ongoing
// 导出 function ongoingOrders(): Promise<OngoingResult> {，供其他模块引用

export function ongoingOrders(): Promise<OngoingResult> {
  // 返回 get('/orders/summary/ongoing') 给调用方

  return get('/orders/summary/ongoing')
// 当前代码块结束

}

// @ticket 05-订单查看与取出：POST /api/v1/orders/open-auth/:token/use（鉴权）
// 导出 function useOpenAuth(token: string, doorId: number, idem: string): Promise<OkResult> {，供其他模块引用

export function useOpenAuth(token: string, doorId: number, idem: string): Promise<OkResult> {
  // 返回 post(`/orders/open-auth/${token}/use`, { door_id: doorId }, idem) 给调用方

  return post(`/orders/open-auth/${token}/use`, { door_id: doorId }, idem)
// 当前代码块结束

}

// @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/open-auth（鉴权）
// 导出 function openAuth(orderNo: string, data: { max_times?: number; ttl_sec?: number }): Promise<OpenAuthResult> {，供其他模块引用

export function openAuth(orderNo: string, data: { max_times?: number; ttl_sec?: number }): Promise<OpenAuthResult> {
  // 返回 post(`/orders/${orderNo}/open-auth`, data) 给调用方

  return post(`/orders/${orderNo}/open-auth`, data)
// 当前代码块结束

}

// @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/makeup（鉴权 + 幂等）
// 导出 function makeup(orderNo: string, amount: number, idem: string): Promise<MakeupResult> {，供其他模块引用

export function makeup(orderNo: string, amount: number, idem: string): Promise<MakeupResult> {
  // 返回 post(`/orders/${orderNo}/makeup`, { amount }, idem) 给调用方

  return post(`/orders/${orderNo}/makeup`, { amount }, idem)
// 当前代码块结束

}

// @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/complete（鉴权 + 幂等）
// 导出 function completeOrder(orderNo: string, idem: string): Promise<OkResult> {，供其他模块引用

export function completeOrder(orderNo: string, idem: string): Promise<OkResult> {
  // 返回 post(`/orders/${orderNo}/complete`, {}, idem) 给调用方

  return post(`/orders/${orderNo}/complete`, {}, idem)
// 当前代码块结束

}
