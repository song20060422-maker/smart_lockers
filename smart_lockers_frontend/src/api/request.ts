// @ticket 00-project：Axios 请求封装——统一 baseURL、token 注入、响应外壳解包、错误码拦截、幂等键、重复请求取消
import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from 'axios'
import type { ApiError, ApiResponse } from '@/types'

// @ticket 00-project：业务错误类，承载后端 code/message/request_id
export class BizError extends Error {
  code: number
  requestId?: string
  constructor(code: number, message: string, requestId?: string) {
    super(message)
    this.name = 'BizError'
    this.code = code
    this.requestId = requestId
  }
}

// @ticket 06-个人中心：token 读写（localStorage 持久化，退出登录即时清除）
const TOKEN_KEY = 'sl_token'
export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}
export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}
export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

// @ticket 00-project：运行中请求 Map，用于取消重复请求（同 method+url+params+data）
const pending = new Map<string, AbortController>()

function requestKey(config: AxiosRequestConfig): string {
  return [config.method, config.url, JSON.stringify(config.params), JSON.stringify(config.data)].join('|')
}

function addPending(config: InternalAxiosRequestConfig): void {
  const key = requestKey(config)
  if (!key) return
  const controller = new AbortController()
  config.signal = controller.signal
  pending.set(key, controller)
}

function removePending(config: AxiosRequestConfig): void {
  const key = requestKey(config)
  if (key && pending.has(key)) {
    pending.delete(key)
  }
}

// @ticket 00-project：创建 axios 实例，baseURL 指向后端 /api/v1（见 .env）
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

// @ticket 00-project + 06：请求拦截器——注入 Bearer token、幂等键、取消重复请求
service.interceptors.request.use(
  (config) => {
    removePending(config)
    addPending(config)
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    // @ticket 00-project：写接口幂等键，从配置扩展字段透传为 Idempotency-Key 头
    const idem = (config as AxiosRequestConfig & { idempotencyKey?: string }).idempotencyKey
    if (idem) {
      config.headers['Idempotency-Key'] = idem
    }
    return config
  },
  (error) => Promise.reject(error),
)

// @ticket 00-project：响应拦截器——解包 {code,message,data,request_id} 外壳，code!==0 走错误处理
service.interceptors.response.use(
  (response) => {
    removePending(response.config)
    const body = response.data as ApiResponse
    if (!body || typeof body.code !== 'number') {
      return Promise.reject(new BizError(-1, '响应格式异常'))
    }
    if (body.code === 0) {
      return body.data as never
    }
    // @ticket 06：登录失效（10002）清除本地 token
    if (body.code === 10002) {
      clearToken()
    }
    const err: ApiError = { code: body.code, message: body.message, request_id: body.request_id }
    return Promise.reject(new BizError(err.code, err.message, err.request_id))
  },
  (error) => {
    if (error?.config) removePending(error.config)
    if (axios.isCancel(error)) {
      return Promise.reject(new BizError(-2, '请求已取消（重复提交）'))
    }
    const status = error?.response?.status
    let message = error?.message || '网络异常，请稍后重试'
    if (status === 403) message = '禁止访问'
    else if (status === 404) message = '资源不存在'
    else if (status && status >= 500) message = '系统内部错误'
    return Promise.reject(new BizError(status ? -status : -1, message))
  },
)

// @ticket 00-project：通用请求方法；支持通过 config.idempotencyKey 设置幂等键
export function request<T>(config: AxiosRequestConfig & { idempotencyKey?: string }): Promise<T> {
  return service.request(config) as unknown as Promise<T>
}

// @ticket 00-project：便捷方法，降低各业务 API 模块的书写复杂度
export function get<T>(url: string, params?: object): Promise<T> {
  return request<T>({ url, method: 'GET', params })
}

// @ticket 00-project：POST 便捷方法，第三个参数透传幂等键
export function post<T>(url: string, data?: unknown, idempotencyKey?: string): Promise<T> {
  return request<T>({ url, method: 'POST', data, idempotencyKey })
}

export default service
