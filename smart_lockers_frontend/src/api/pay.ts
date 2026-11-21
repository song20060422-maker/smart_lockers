// @ticket 02-支付宝支付：余额/明细/充值/记录接口封装
import { get, post } from './request'
import type { BalanceResult, BalanceLog, RechargeOrder, CreateRechargeResult, ListResult, PageQuery } from '@/types'

// @ticket 02-支付宝支付：GET /api/v1/pay/balance
// 导出 function getBalance(): Promise<BalanceResult> {，供其他模块引用

export function getBalance(): Promise<BalanceResult> {
  // 返回 get('/pay/balance') 给调用方

  return get('/pay/balance')
// 当前代码块结束

}

// @ticket 02-支付宝支付：GET /api/v1/pay/balance/logs
// 导出 function balanceLogs(params?: PageQuery): Promise<ListResult<BalanceLog>> {，供其他模块引用

export function balanceLogs(params?: PageQuery): Promise<ListResult<BalanceLog>> {
  // 返回 get('/pay/balance/logs', params) 给调用方

  return get('/pay/balance/logs', params)
// 当前代码块结束

}

// @ticket 02-支付宝支付：POST /api/v1/pay/recharge（鉴权 + 幂等）
// 导出 function createRecharge(amount: number, idem: string): Promise<CreateRechargeResult> {，供其他模块引用

export function createRecharge(amount: number, idem: string): Promise<CreateRechargeResult> {
  // 返回 post('/pay/recharge', { amount }, idem) 给调用方

  return post('/pay/recharge', { amount }, idem)
// 当前代码块结束

}

// @ticket 02-支付宝支付：GET /api/v1/pay/recharge/records
// 导出 function rechargeRecords(params?: PageQuery): Promise<ListResult<RechargeOrder>> {，供其他模块引用

export function rechargeRecords(params?: PageQuery): Promise<ListResult<RechargeOrder>> {
  // 返回 get('/pay/recharge/records', params) 给调用方

  return get('/pay/recharge/records', params)
// 当前代码块结束

}
