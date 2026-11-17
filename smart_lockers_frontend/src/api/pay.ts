// @ticket 02-支付宝支付：余额/明细/充值/记录接口封装
import { get, post } from './request'
import type { BalanceResult, BalanceLog, RechargeOrder, CreateRechargeResult, ListResult, PageQuery } from '@/types'

// @ticket 02-支付宝支付：GET /api/v1/pay/balance
export function getBalance(): Promise<BalanceResult> {
  return get('/pay/balance')
}

// @ticket 02-支付宝支付：GET /api/v1/pay/balance/logs
export function balanceLogs(params?: PageQuery): Promise<ListResult<BalanceLog>> {
  return get('/pay/balance/logs', params)
}

// @ticket 02-支付宝支付：POST /api/v1/pay/recharge（鉴权 + 幂等）
export function createRecharge(amount: number, idem: string): Promise<CreateRechargeResult> {
  return post('/pay/recharge', { amount }, idem)
}

// @ticket 02-支付宝支付：GET /api/v1/pay/recharge/records
export function rechargeRecords(params?: PageQuery): Promise<ListResult<RechargeOrder>> {
  return get('/pay/recharge/records', params)
}
