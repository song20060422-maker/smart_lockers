// @ticket 05-订单查看与取出：列表/详情/开柜/补交/完成接口封装
import { get, post } from './request'
import type {
  OrderListResult, OrderDetail, OngoingResult, OpenAuthResult, MakeupResult, OkResult,
} from '@/types'

// @ticket 05-订单查看与取出：GET /api/v1/orders（游标分页）
export function orderList(params: { type?: string; cursor?: string; limit?: number }): Promise<OrderListResult> {
  return get('/orders', params)
}

// @ticket 05-订单查看与取出：GET /api/v1/orders/:order_no
export function orderDetail(orderNo: string): Promise<OrderDetail> {
  return get(`/orders/${orderNo}`)
}

// @ticket 05-订单查看与取出：GET /api/v1/orders/summary/ongoing
export function ongoingOrders(): Promise<OngoingResult> {
  return get('/orders/summary/ongoing')
}

// @ticket 05-订单查看与取出：POST /api/v1/orders/open-auth/:token/use（鉴权）
export function useOpenAuth(token: string, doorId: number, idem: string): Promise<OkResult> {
  return post(`/orders/open-auth/${token}/use`, { door_id: doorId }, idem)
}

// @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/open-auth（鉴权）
export function openAuth(orderNo: string, data: { max_times?: number; ttl_sec?: number }): Promise<OpenAuthResult> {
  return post(`/orders/${orderNo}/open-auth`, data)
}

// @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/makeup（鉴权 + 幂等）
export function makeup(orderNo: string, amount: number, idem: string): Promise<MakeupResult> {
  return post(`/orders/${orderNo}/makeup`, { amount }, idem)
}

// @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/complete（鉴权 + 幂等）
export function completeOrder(orderNo: string, idem: string): Promise<OkResult> {
  return post(`/orders/${orderNo}/complete`, {}, idem)
}
