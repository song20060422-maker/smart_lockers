// @ticket 08-性能调优：运维观测接口封装（鉴权 + 内网限制）
import { get, request } from './request'
import type { HealthResult, OpsPprofResult, OpsTraceResult, OpsExpvarResult } from '@/types'

// @ticket 08-性能调优：GET /api/v1/ops/pprof/（鉴权）
export function opsPprof(): Promise<OpsPprofResult> {
  return get('/ops/pprof/')
}

// @ticket 08-性能调优：GET /api/v1/ops/pprof/{profile}（鉴权，profile=profile|heap|goroutine|block）
export function opsPprofProfile(profile: 'profile' | 'heap' | 'goroutine' | 'block'): Promise<OpsPprofResult> {
  return get(`/ops/pprof/${profile}`)
}

// @ticket 08-性能调优：GET /api/v1/ops/trace（鉴权）
export function opsTrace(): Promise<OpsTraceResult> {
  return get('/ops/trace')
}

// @ticket 08-性能调优：GET /api/v1/ops/expvar（鉴权）
export function opsExpvar(): Promise<OpsExpvarResult> {
  return get('/ops/expvar')
}

// @ticket 00-project：GET /health（健康检查，后端注册在 v1 组外，用绝对路径绕过 baseURL）
export function health(): Promise<HealthResult> {
  return request<HealthResult>({ url: '/health', method: 'GET', baseURL: '' })
}
