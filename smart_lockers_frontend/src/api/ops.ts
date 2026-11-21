// @ticket 08-性能调优：运维观测接口封装（鉴权 + 内网限制）
import { get, request } from './request'
import type { HealthResult, OpsPprofResult, OpsTraceResult, OpsExpvarResult } from '@/types'

// @ticket 08-性能调优：GET /api/v1/ops/pprof/（鉴权）
// 导出 function opsPprof(): Promise<OpsPprofResult> {，供其他模块引用

export function opsPprof(): Promise<OpsPprofResult> {
  // 返回 get('/ops/pprof/') 给调用方

  return get('/ops/pprof/')
// 当前代码块结束

}

// @ticket 08-性能调优：GET /api/v1/ops/pprof/{profile}（鉴权，profile=profile|heap|goroutine|block）
// 导出 function opsPprofProfile(profile: 'profile' | 'heap' | 'goroutine' | 'block'): Promise<OpsPprofResult> {，供其他模块引用

export function opsPprofProfile(profile: 'profile' | 'heap' | 'goroutine' | 'block'): Promise<OpsPprofResult> {
  // 返回 get(`/ops/pprof/${profile}`) 给调用方

  return get(`/ops/pprof/${profile}`)
// 当前代码块结束

}

// @ticket 08-性能调优：GET /api/v1/ops/trace（鉴权）
// 导出 function opsTrace(): Promise<OpsTraceResult> {，供其他模块引用

export function opsTrace(): Promise<OpsTraceResult> {
  // 返回 get('/ops/trace') 给调用方

  return get('/ops/trace')
// 当前代码块结束

}

// @ticket 08-性能调优：GET /api/v1/ops/expvar（鉴权）
// 导出 function opsExpvar(): Promise<OpsExpvarResult> {，供其他模块引用

export function opsExpvar(): Promise<OpsExpvarResult> {
  // 返回 get('/ops/expvar') 给调用方

  return get('/ops/expvar')
// 当前代码块结束

}

// @ticket 00-project：GET /health（健康检查，后端注册在 v1 组外，用绝对路径绕过 baseURL）
// 导出 function health(): Promise<HealthResult> {，供其他模块引用

export function health(): Promise<HealthResult> {
  // 返回 request<HealthResult>({ url: '/health', method: 'GET', baseURL: '' }) 给调用方

  return request<HealthResult>({ url: '/health', method: 'GET', baseURL: '' })
// 当前代码块结束

}
