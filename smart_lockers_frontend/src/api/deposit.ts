// @ticket 04-物品寄存：柜子/试算/落单/换柜接口封装
import { get, post } from './request'
import type { Cabinet, FeeResult, CreateDepositResult, ChangeCabinetResult } from '@/types'

// @ticket 04-物品寄存：GET /api/v1/deposit/cabinets
export function queryCabinets(params: { site_id?: number; size?: string }): Promise<Cabinet[]> {
  return get('/deposit/cabinets', params)
}

// @ticket 04-物品寄存：GET /api/v1/deposit/quote
export function quoteFee(params: { size: string; hours: number; coupon?: number }): Promise<FeeResult> {
  return get('/deposit/quote', params)
}

// @ticket 04-物品寄存：POST /api/v1/deposit/orders（鉴权 + 幂等）
export function createDeposit(data: CreateDepositPayload, idem: string): Promise<CreateDepositResult> {
  return post('/deposit/orders', data, idem)
}

// @ticket 04-物品寄存：POST /api/v1/deposit/orders/:order_no/change-cabinet
export function changeCabinet(orderNo: string, newDoorId: number): Promise<ChangeCabinetResult> {
  return post(`/deposit/orders/${orderNo}/change-cabinet`, { new_door_id: newDoorId })
}

interface CreateDepositPayload {
  door_id: number
  site_id: number
  hours: number
  coupon?: number
}
