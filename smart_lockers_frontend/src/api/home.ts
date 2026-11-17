// @ticket 03-首页：定位/城市/寄存点/指南接口封装
import { request } from './request'
import type { LocateResult, City, NearbySitesResult, GuideResult } from '@/types'

// @ticket 03-首页：POST /api/v1/home/locate（GPS 优先，IP/默认回退）
export function locate(params: {
  lat?: number
  lng?: number
  ip?: string
  fallback_city?: string
}): Promise<LocateResult> {
  return request({ url: '/home/locate', method: 'POST', data: params })
}

// @ticket 03-首页：GET /api/v1/home/cities?keyword=
export function searchCities(keyword: string): Promise<City[]> {
  return request({ url: '/home/cities', method: 'GET', params: { keyword } })
}

// @ticket 03-首页：GET /api/v1/home/sites?lat=&lng=&page=&page_size=
export function nearbySites(params: {
  lat?: number
  lng?: number
  page?: number
  page_size?: number
}): Promise<NearbySitesResult> {
  return request({ url: '/home/sites', method: 'GET', params })
}

// @ticket 03-首页：GET /api/v1/home/guide?version=&lang=
export function getGuide(params?: { version?: string; lang?: string }): Promise<GuideResult> {
  return request({ url: '/home/guide', method: 'GET', params })
}
