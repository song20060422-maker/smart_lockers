// @ticket 03-首页：定位/城市/寄存点/指南接口封装
import { request } from './request'
import type { LocateResult, City, NearbySitesResult, GuideResult } from '@/types'

// @ticket 03-首页：POST /api/v1/home/locate（GPS 优先，IP/默认回退）
// 导出 function locate(params: {，供其他模块引用

export function locate(params: {
  // 执行 lat?: number，实现对应的前端逻辑

  lat?: number
  // 执行 lng?: number，实现对应的前端逻辑

  lng?: number
  // 执行 ip?: string，实现对应的前端逻辑

  ip?: string
  // 执行 fallback_city?: string，实现对应的前端逻辑

  fallback_city?: string
// 执行 }): Promise<LocateResult> {，实现对应的前端逻辑

}): Promise<LocateResult> {
  // 返回 request({ url: '/home/locate', method: 'POST', data: params }) 给调用方

  return request({ url: '/home/locate', method: 'POST', data: params })
// 当前代码块结束

}

// @ticket 03-首页：GET /api/v1/home/cities?keyword=
// 导出 function searchCities(keyword: string): Promise<City[]> {，供其他模块引用

export function searchCities(keyword: string): Promise<City[]> {
  // 返回 request({ url: '/home/cities', method: 'GET', params: { keyword } }) 给调用方

  return request({ url: '/home/cities', method: 'GET', params: { keyword } })
// 当前代码块结束

}

// @ticket 03-首页：GET /api/v1/home/sites?lat=&lng=&page=&page_size=
// 导出 function nearbySites(params: {，供其他模块引用

export function nearbySites(params: {
  // 执行 lat?: number，实现对应的前端逻辑

  lat?: number
  // 执行 lng?: number，实现对应的前端逻辑

  lng?: number
  // 执行 page?: number，实现对应的前端逻辑

  page?: number
  // 执行 page_size?: number，实现对应的前端逻辑

  page_size?: number
// 执行 }): Promise<NearbySitesResult> {，实现对应的前端逻辑

}): Promise<NearbySitesResult> {
  // 返回 request({ url: '/home/sites', method: 'GET', params }) 给调用方

  return request({ url: '/home/sites', method: 'GET', params })
// 当前代码块结束

}

// @ticket 03-首页：GET /api/v1/home/guide?version=&lang=
// 导出 function getGuide(params?: { version?: string; lang?: string }): Promise<GuideResult> {，供其他模块引用

export function getGuide(params?: { version?: string; lang?: string }): Promise<GuideResult> {
  // 返回 request({ url: '/home/guide', method: 'GET', params }) 给调用方

  return request({ url: '/home/guide', method: 'GET', params })
// 当前代码块结束

}
