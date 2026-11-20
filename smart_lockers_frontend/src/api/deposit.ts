// @ticket 04-物品寄存：柜子/试算/落单/换柜接口封装
import { get, post } from './request'
import type { Cabinet, FeeResult, CreateDepositResult, ChangeCabinetResult } from '@/types'

// @ticket 04-物品寄存：GET /api/v1/deposit/cabinets
// 导出 function queryCabinets(params: { site_id?: number; size?: string }): Promise<Cabinet[]> {，供其他模块引用

export function queryCabinets(params: { site_id?: number; size?: string }): Promise<Cabinet[]> {
  // 返回 get('/deposit/cabinets', params) 给调用方

  return get('/deposit/cabinets', params)
// 当前代码块结束

}

// @ticket 04-物品寄存：GET /api/v1/deposit/quote
// 导出 function quoteFee(params: { size: string; hours: number; coupon?: number }): Promise<FeeResult> {，供其他模块引用

export function quoteFee(params: { size: string; hours: number; coupon?: number }): Promise<FeeResult> {
  // 返回 get('/deposit/quote', params) 给调用方

  return get('/deposit/quote', params)
// 当前代码块结束

}

// @ticket 04-物品寄存：POST /api/v1/deposit/orders（鉴权 + 幂等）
// 导出 function createDeposit(data: CreateDepositPayload, idem: string): Promise<CreateDepositResult> {，供其他模块引用

export function createDeposit(data: CreateDepositPayload, idem: string): Promise<CreateDepositResult> {
  // 返回 post('/deposit/orders', data, idem) 给调用方

  return post('/deposit/orders', data, idem)
// 当前代码块结束

}

// @ticket 04-物品寄存：POST /api/v1/deposit/orders/:order_no/change-cabinet
// 导出 function changeCabinet(orderNo: string, newDoorId: number): Promise<ChangeCabinetResult> {，供其他模块引用

export function changeCabinet(orderNo: string, newDoorId: number): Promise<ChangeCabinetResult> {
  // 返回 post(`/deposit/orders/${orderNo}/change-cabinet`, { new_door_id: newDoorId }) 给调用方

  return post(`/deposit/orders/${orderNo}/change-cabinet`, { new_door_id: newDoorId })
// 当前代码块结束

}

// 定义接口，约束对象的结构：interface CreateDepositPayload {

interface CreateDepositPayload {
  // 定义对象属性 door_id，取值为 number

  door_id: number
  // 定义对象属性 site_id，取值为 number

  site_id: number
  // 定义对象属性 hours，取值为 number

  hours: number
  // 执行 coupon?: number，实现对应的前端逻辑

  coupon?: number
// 当前代码块结束

}
