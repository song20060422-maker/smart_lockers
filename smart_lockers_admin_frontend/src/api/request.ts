// @ticket B端-前端：axios 封装——baseURL、Bearer 注入、业务码拦截、幂等键 @ticket B端-幂等
import axios, { type InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'

// 声明变量 raw = axios.create({，用于保存数据或函数引用

const raw = axios.create({
  // 定义对象属性 baseURL，取值为 import.meta.env.VITE_API_BASE || '/api/v1'

  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  // 定义对象属性 timeout，取值为 15000

  timeout: 15000,
// 执行 })，实现对应的前端逻辑

})

// 声明变量 token: string | null = localStorage.getItem('admin_token')，用于保存数据或函数引用

let token: string | null = localStorage.getItem('admin_token')
// 导出 function setToken(t: string | null) {，供其他模块引用

export function setToken(t: string | null) {
  // 赋值更新：token = t

  token = t
  // 条件判断：if (t) localStorage.setItem('admin_token', t) 时执行下面逻辑

  if (t) localStorage.setItem('admin_token', t)
  // 执行 else localStorage.removeItem('admin_token')，实现对应的前端逻辑

  else localStorage.removeItem('admin_token')
// 当前代码块结束

}
// 导出 function getToken() { return token }，供其他模块引用

export function getToken() { return token }

// 请求拦截：注入 Bearer token 与幂等键（写操作可覆盖）
// 执行 raw.interceptors.request.use((cfg: InternalAxiosRequestConfig) => {，实现对应的前端逻辑

raw.interceptors.request.use((cfg: InternalAxiosRequestConfig) => {
  // 条件判断：if (token) cfg.headers.set('Authorization', `Bearer ${token}`) 时执行下面逻辑

  if (token) cfg.headers.set('Authorization', `Bearer ${token}`)
  // 返回 cfg 给调用方

  return cfg
// 执行 })，实现对应的前端逻辑

})

// 响应拦截：后端统一 HTTP 200 + 业务 code，非 0 视为业务错误
// 执行 raw.interceptors.response.use(，实现对应的前端逻辑

raw.interceptors.response.use(
  // 执行 (resp) => {，实现对应的前端逻辑

  (resp) => {
    // 声明变量 body = resp.data as ApiResponse，用于保存数据或函数引用

    const body = resp.data as ApiResponse
    // 条件判断：if (body && typeof body.code === 'number' && body.code !== 0) 时执行下面逻辑

    if (body && typeof body.code === 'number' && body.code !== 0) {
      // 返回 Promise.reject(new Error(body.message || `业务错误(${body.code})`)) 给调用方

      return Promise.reject(new Error(body.message || `业务错误(${body.code})`))
    // 当前代码块结束

    }
    // 返回 body as any 给调用方

    return body as any
  // 定义结束

  },
  // 执行 (err) => Promise.reject(err),，实现对应的前端逻辑

  (err) => Promise.reject(err),
// 结束参数/表达式列表

)

// 拦截器已将响应解包为 ApiResponse，声明真实返回类型 @ticket B端-前端
// 导出 interface Http {，供其他模块引用

export interface Http {
  // 赋值更新：get<T = unknown>(url: string, config?: any): Promise<ApiResponse<T>>

  get<T = unknown>(url: string, config?: any): Promise<ApiResponse<T>>
  // 赋值更新：post<T = unknown>(url: string, body?: any, config?: any): Promise<ApiResponse<T>>

  post<T = unknown>(url: string, body?: any, config?: any): Promise<ApiResponse<T>>
// 当前代码块结束

}
// 声明变量 http = raw as unknown as Http，用于保存数据或函数引用

const http = raw as unknown as Http

// 为写操作附带幂等键：重试不会重复执行资金/状态变更
// 导出 function idem(key: string) {，供其他模块引用

export function idem(key: string) {
  // 返回 { headers: { 'Idempotency-Key': key } } 给调用方

  return { headers: { 'Idempotency-Key': key } }
// 当前代码块结束

}

// 导出 http，供其他模块引用

export default http
