// @ticket B端-前端：axios 封装——baseURL、Bearer 注入、业务码拦截、幂等键 @ticket B端-幂等
import axios, { type InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'

const raw = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 15000,
})

let token: string | null = localStorage.getItem('admin_token')
export function setToken(t: string | null) {
  token = t
  if (t) localStorage.setItem('admin_token', t)
  else localStorage.removeItem('admin_token')
}
export function getToken() { return token }

// 请求拦截：注入 Bearer token 与幂等键（写操作可覆盖）
raw.interceptors.request.use((cfg: InternalAxiosRequestConfig) => {
  if (token) cfg.headers.set('Authorization', `Bearer ${token}`)
  return cfg
})

// 响应拦截：后端统一 HTTP 200 + 业务 code，非 0 视为业务错误
raw.interceptors.response.use(
  (resp) => {
    const body = resp.data as ApiResponse
    if (body && typeof body.code === 'number' && body.code !== 0) {
      return Promise.reject(new Error(body.message || `业务错误(${body.code})`))
    }
    return body as any
  },
  (err) => Promise.reject(err),
)

// 拦截器已将响应解包为 ApiResponse，声明真实返回类型 @ticket B端-前端
export interface Http {
  get<T = unknown>(url: string, config?: any): Promise<ApiResponse<T>>
  post<T = unknown>(url: string, body?: any, config?: any): Promise<ApiResponse<T>>
}
const http = raw as unknown as Http

// 为写操作附带幂等键：重试不会重复执行资金/状态变更
export function idem(key: string) {
  return { headers: { 'Idempotency-Key': key } }
}

export default http
