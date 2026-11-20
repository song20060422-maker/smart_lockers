// @ticket B端-前端：统一类型定义，字段与后端 model/entity.go + pkg/response.go 对齐

// 统一响应结构（后端 HTTP 200 + 业务 code）
// 导出 interface ApiResponse<T = unknown> {，供其他模块引用

export interface ApiResponse<T = unknown> {
  // 定义对象属性 code，取值为 number

  code: number
  // 定义对象属性 message，取值为 string

  message: string
  // 执行 data?: T，实现对应的前端逻辑

  data?: T
  // 执行 request_id?: string，实现对应的前端逻辑

  request_id?: string
// 当前代码块结束

}

// 导出 type AdminRole = 'user' | 'site_admin' | 'operator'，供其他模块引用

export type AdminRole = 'user' | 'site_admin' | 'operator'

// 导出 interface AdminUser {，供其他模块引用

export interface AdminUser {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 username，取值为 string

  username: string
  // 定义对象属性 role，取值为 AdminRole

  role: AdminRole
  // 定义对象属性 site_id，取值为 number

  site_id: number
  // 定义对象属性 status，取值为 number

  status: number
// 当前代码块结束

}

// 网点管理 @ticket B端-网点管理
// 导出 interface Site {，供其他模块引用

export interface Site {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 code，取值为 string

  code: string
  // 定义对象属性 name，取值为 string

  name: string
  // 定义对象属性 city，取值为 string

  city: string
  // 定义对象属性 address，取值为 string

  address: string
  // 定义对象属性 manager_phone，取值为 string

  manager_phone: string
  // 定义对象属性 status，取值为 number // 1=启用 0=停用

  status: number // 1=启用 0=停用
  // 定义对象属性 door_total，取值为 number

  door_total: number
  // 定义对象属性 door_fault，取值为 number

  door_fault: number
  // 定义对象属性 created_at，取值为 string

  created_at: string
  // 定义对象属性 updated_at，取值为 string

  updated_at: string
// 当前代码块结束

}

// 设备监控 @ticket B端-设备监控
// 导出 interface Device {，供其他模块引用

export interface Device {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 site_id，取值为 number

  site_id: number
  // 定义对象属性 code，取值为 string

  code: string
  // 定义对象属性 type，取值为 string

  type: string
  // 定义对象属性 status，取值为 number // 0=在线 1=离线 2=故障

  status: number // 0=在线 1=离线 2=故障
  // 定义对象属性 last_heartbeat，取值为 string

  last_heartbeat: string
  // 定义对象属性 updated_at，取值为 string

  updated_at: string
// 当前代码块结束

}

// 导出 interface DeviceMonitor {，供其他模块引用

export interface DeviceMonitor {
  // 定义对象属性 total，取值为 number

  total: number
  // 定义对象属性 online，取值为 number

  online: number
  // 定义对象属性 offline，取值为 number

  offline: number
  // 定义对象属性 fault，取值为 number

  fault: number
  // 定义对象属性 items，取值为 Device[]

  items: Device[]
// 当前代码块结束

}

// 充值对账 @ticket B端-充值对账
// 导出 interface ReconcileRecord {，供其他模块引用

export interface ReconcileRecord {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 batch_no，取值为 string

  batch_no: string
  // 定义对象属性 site_id，取值为 number

  site_id: number
  // 定义对象属性 expected_amount，取值为 number // 分

  expected_amount: number // 分
  // 定义对象属性 actual_amount，取值为 number   // 分

  actual_amount: number   // 分
  // 定义对象属性 diff_amount，取值为 number     // 分

  diff_amount: number     // 分
  // 定义对象属性 status，取值为 number

  status: number
  // 定义对象属性 idempotency_key，取值为 string

  idempotency_key: string
  // 定义对象属性 operator_id，取值为 number

  operator_id: number
  // 定义对象属性 created_at，取值为 string

  created_at: string
// 当前代码块结束

}

// 告警处理 @ticket B端-告警处理
// 导出 interface Alarm {，供其他模块引用

export interface Alarm {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 site_id，取值为 number

  site_id: number
  // 定义对象属性 device_id，取值为 number

  device_id: number
  // 定义对象属性 level，取值为 string

  level: string
  // 定义对象属性 content，取值为 string

  content: string
  // 定义对象属性 status，取值为 number // 0=待处理 1=已处理

  status: number // 0=待处理 1=已处理
  // 定义对象属性 handler，取值为 string

  handler: string
  // 定义对象属性 handle_note，取值为 string

  handle_note: string
  // 定义对象属性 idempotency_key，取值为 string

  idempotency_key: string
  // 定义对象属性 created_at，取值为 string

  created_at: string
  // 定义对象属性 handled_at，取值为 string

  handled_at: string
// 当前代码块结束

}

// 权限点（与后端 rbac.go 对齐，用于按钮级显隐）
// 导出 const PERMISSIONS: Record<AdminRole, string[]> = {，供其他模块引用

export const PERMISSIONS: Record<AdminRole, string[]> = {
  // 定义对象属性 user，取值为 ['site:read']

  user: ['site:read'],
  // 修复偏差5：权限点统一为后端 rbac.go 定义的 device:control，原 device:ctl 前端无法命中后端鉴权 @ticket B端-前端

  site_admin: ['site:read', 'site:write', 'device:read', 'device:control', 'reconcile:read', 'alarm:read', 'alarm:handle',],
  // 定义对象属性 operator，取值为含 device:control 的权限列表

  operator: ['site:read', 'site:write', 'device:read', 'device:control', 'reconcile:read', 'reconcile:run', 'alarm:read', 'alarm:handle', 'audit:read'],
// 当前代码块结束

}
