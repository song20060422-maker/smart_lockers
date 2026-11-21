// @ticket 00-project：Axios 请求封装——统一 baseURL、token 注入、响应外壳解包、错误码拦截、幂等键、重复请求取消
import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from 'axios'
import type { ApiError, ApiResponse } from '@/types'

// @ticket 00-project：业务错误类，承载后端 code/message/request_id
// 导出 class BizError extends Error {，供其他模块引用

export class BizError extends Error {
  // 定义对象属性 code，取值为 number

  code: number
  // 执行 requestId?: string，实现对应的前端逻辑

  requestId?: string
  // 调用 constructor(code: number, message: string, requestId?: string) {，触发对应功能

  constructor(code: number, message: string, requestId?: string) {
    // 调用 super(message)，触发对应功能

    super(message)
    // 赋值更新：this.name = 'BizError'

    this.name = 'BizError'
    // 赋值更新：this.code = code

    this.code = code
    // 赋值更新：this.requestId = requestId

    this.requestId = requestId
  // 当前代码块结束

  }
// 当前代码块结束

}

// @ticket 06-个人中心：token 读写（localStorage 持久化，退出登录即时清除）
// 声明变量 TOKEN_KEY = 'sl_token'，用于保存数据或函数引用

const TOKEN_KEY = 'sl_token'
// 导出 function getToken(): string {，供其他模块引用

export function getToken(): string {
  // 返回 localStorage.getItem(TOKEN_KEY) || '' 给调用方

  return localStorage.getItem(TOKEN_KEY) || ''
// 当前代码块结束

}
// 导出 function setToken(token: string): void {，供其他模块引用

export function setToken(token: string): void {
  // 执行 localStorage.setItem(TOKEN_KEY, token)，实现对应的前端逻辑

  localStorage.setItem(TOKEN_KEY, token)
// 当前代码块结束

}
// 导出 function clearToken(): void {，供其他模块引用

export function clearToken(): void {
  // 执行 localStorage.removeItem(TOKEN_KEY)，实现对应的前端逻辑

  localStorage.removeItem(TOKEN_KEY)
// 当前代码块结束

}

// @ticket 00-project：运行中请求 Map，用于取消重复请求（同 method+url+params+data）
// 声明变量 pending = new Map<string, AbortController>()，用于保存数据或函数引用

const pending = new Map<string, AbortController>()

// 定义函数/方法，封装可复用的处理逻辑：function requestKey(config: AxiosRequestConfig): string {

function requestKey(config: AxiosRequestConfig): string {
  // 返回 [config.method, config.url, JSON.stringify(config.params), JSON.stringify(config.data)].join('|') 给调用方

  return [config.method, config.url, JSON.stringify(config.params), JSON.stringify(config.data)].join('|')
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function addPending(config: InternalAxiosRequestConfig): void {

function addPending(config: InternalAxiosRequestConfig): void {
  // 声明变量 key = requestKey(config)，用于保存数据或函数引用

  const key = requestKey(config)
  // 条件判断：if (!key) return 时执行下面逻辑

  if (!key) return
  // 声明变量 controller = new AbortController()，用于保存数据或函数引用

  const controller = new AbortController()
  // 赋值更新：config.signal = controller.signal

  config.signal = controller.signal
  // 执行 pending.set(key, controller)，实现对应的前端逻辑

  pending.set(key, controller)
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function removePending(config: AxiosRequestConfig): void {

function removePending(config: AxiosRequestConfig): void {
  // 声明变量 key = requestKey(config)，用于保存数据或函数引用

  const key = requestKey(config)
  // 条件判断：if (key && pending.has(key)) 时执行下面逻辑

  if (key && pending.has(key)) {
    // 执行 pending.delete(key)，实现对应的前端逻辑

    pending.delete(key)
  // 当前代码块结束

  }
// 当前代码块结束

}

// @ticket 00-project：创建 axios 实例，baseURL 指向后端 /api/v1（见 .env）
// 声明变量 service: AxiosInstance = axios.create({，用于保存数据或函数引用

const service: AxiosInstance = axios.create({
  // 定义对象属性 baseURL，取值为 import.meta.env.VITE_API_BASE || '/api/v1'

  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
  // 定义对象属性 timeout，取值为 15000

  timeout: 15000,
  // 定义对象属性 headers，取值为 { 'Content-Type': 'application/json' }

  headers: { 'Content-Type': 'application/json' },
// 执行 })，实现对应的前端逻辑

})

// @ticket 00-project + 06：请求拦截器——注入 Bearer token、幂等键、取消重复请求
// 执行 service.interceptors.request.use(，实现对应的前端逻辑

service.interceptors.request.use(
  // 执行 (config) => {，实现对应的前端逻辑

  (config) => {
    // 调用 removePending(config)，触发对应功能

    removePending(config)
    // 调用 addPending(config)，触发对应功能

    addPending(config)
    // 声明变量 token = getToken()，用于保存数据或函数引用

    const token = getToken()
    // 条件判断：if (token) 时执行下面逻辑

    if (token) {
      // 赋值更新：config.headers.Authorization = `Bearer ${token}`

      config.headers.Authorization = `Bearer ${token}`
    // 当前代码块结束

    }
    // @ticket 00-project：写接口幂等键，从配置扩展字段透传为 Idempotency-Key 头
    // 定义函数/方法，封装可复用的处理逻辑：const idem = (config as AxiosRequestConfig & { idempotencyKey?: string }).idempotencyKey

    const idem = (config as AxiosRequestConfig & { idempotencyKey?: string }).idempotencyKey
    // 条件判断：if (idem) 时执行下面逻辑

    if (idem) {
      // 赋值更新：config.headers['Idempotency-Key'] = idem

      config.headers['Idempotency-Key'] = idem
    // 当前代码块结束

    }
    // 返回 config 给调用方

    return config
  // 定义结束

  },
  // 执行 (error) => Promise.reject(error),，实现对应的前端逻辑

  (error) => Promise.reject(error),
// 结束参数/表达式列表

)

// @ticket 00-project：响应拦截器——解包 {code,message,data,request_id} 外壳，code!==0 走错误处理
// 执行 service.interceptors.response.use(，实现对应的前端逻辑

service.interceptors.response.use(
  // 执行 (response) => {，实现对应的前端逻辑

  (response) => {
    // 调用 removePending(response.config)，触发对应功能

    removePending(response.config)
    // 声明变量 body = response.data as ApiResponse，用于保存数据或函数引用

    const body = response.data as ApiResponse
    // 条件判断：if (!body || typeof body.code !== 'number') 时执行下面逻辑

    if (!body || typeof body.code !== 'number') {
      // 返回 Promise.reject(new BizError(-1, '响应格式异常')) 给调用方

      return Promise.reject(new BizError(-1, '响应格式异常'))
    // 当前代码块结束

    }
    // 条件判断：if (body.code === 0) 时执行下面逻辑

    if (body.code === 0) {
      // 返回 body.data as never 给调用方

      return body.data as never
    // 当前代码块结束

    }
    // @ticket 06：登录失效（10002）清除本地 token
    // 条件判断：if (body.code === 10002) 时执行下面逻辑

    if (body.code === 10002) {
      // 调用 clearToken()，触发对应功能

      clearToken()
    // 当前代码块结束

    }
    // 声明变量 err: ApiError = { code: body.code, message: body.message, request_id: body.request_id }，用于保存数据或函数引用

    const err: ApiError = { code: body.code, message: body.message, request_id: body.request_id }
    // 返回 Promise.reject(new BizError(err.code, err.message, err.request_id)) 给调用方

    return Promise.reject(new BizError(err.code, err.message, err.request_id))
  // 定义结束

  },
  // 执行 (error) => {，实现对应的前端逻辑

  (error) => {
    // 条件判断：if (error?.config) removePending(error.config) 时执行下面逻辑

    if (error?.config) removePending(error.config)
    // 条件判断：if (axios.isCancel(error)) 时执行下面逻辑

    if (axios.isCancel(error)) {
      // 返回 Promise.reject(new BizError(-2, '请求已取消（重复提交）')) 给调用方

      return Promise.reject(new BizError(-2, '请求已取消（重复提交）'))
    // 当前代码块结束

    }
    // 声明变量 status = error?.response?.status，用于保存数据或函数引用

    const status = error?.response?.status
    // 声明变量 message = error?.message || '网络异常，请稍后重试'，用于保存数据或函数引用

    let message = error?.message || '网络异常，请稍后重试'
    // 条件判断：if (status === 403) message = '禁止访问' 时执行下面逻辑

    if (status === 403) message = '禁止访问'
    // 否则再次判断：(status === 404) message = '资源不存在' 时执行

    else if (status === 404) message = '资源不存在'
    // 否则再次判断：(status && status >= 500) message = '系统内部错误' 时执行

    else if (status && status >= 500) message = '系统内部错误'
    // 返回 Promise.reject(new BizError(status ? -status : -1, message)) 给调用方

    return Promise.reject(new BizError(status ? -status : -1, message))
  // 定义结束

  },
// 结束参数/表达式列表

)

// @ticket 00-project：通用请求方法；支持通过 config.idempotencyKey 设置幂等键
// 导出 function request<T>(config: AxiosRequestConfig & { idempotencyKey?: string }): Promise<T> {，供其他模块引用

export function request<T>(config: AxiosRequestConfig & { idempotencyKey?: string }): Promise<T> {
  // 返回 service.request(config) as unknown as Promise<T> 给调用方

  return service.request(config) as unknown as Promise<T>
// 当前代码块结束

}

// @ticket 00-project：便捷方法，降低各业务 API 模块的书写复杂度
// 导出 function get<T>(url: string, params?: object): Promise<T> {，供其他模块引用

export function get<T>(url: string, params?: object): Promise<T> {
  // 返回 request<T>({ url, method: 'GET', params }) 给调用方

  return request<T>({ url, method: 'GET', params })
// 当前代码块结束

}

// @ticket 00-project：POST 便捷方法，第三个参数透传幂等键
// 导出 function post<T>(url: string, data?: unknown, idempotencyKey?: string): Promise<T> {，供其他模块引用

export function post<T>(url: string, data?: unknown, idempotencyKey?: string): Promise<T> {
  // 返回 request<T>({ url, method: 'POST', data, idempotencyKey }) 给调用方

  return request<T>({ url, method: 'POST', data, idempotencyKey })
// 当前代码块结束

}

// 导出 service，供其他模块引用

export default service
