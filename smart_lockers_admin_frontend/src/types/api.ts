// @ticket B端-前端：统一类型定义，字段与后端 model/entity.go + pkg/response.go 对齐

// 统一响应结构（后端 HTTP 200 + 业务 code）
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
  request_id?: string
}

export type AdminRole = 'user' | 'site_admin' | 'operator'

export interface AdminUser {
  id: number
  username: string
  role: AdminRole
  site_id: number
  status: number
}

// 网点管理 @ticket B端-网点管理
export interface Site {
  id: number
  code: string
  name: string
  city: string
  address: string
  manager_phone: string
  status: number // 1=启用 0=停用
  door_total: number
  door_fault: number
  created_at: string
  updated_at: string
}

// 设备监控 @ticket B端-设备监控
export interface Device {
  id: number
  site_id: number
  code: string
  type: string
  status: number // 0=在线 1=离线 2=故障
  last_heartbeat: string
  updated_at: string
}

export interface DeviceMonitor {
  total: number
  online: number
  offline: number
  fault: number
  items: Device[]
}

// 充值对账 @ticket B端-充值对账
export interface ReconcileRecord {
  id: number
  batch_no: string
  site_id: number
  expected_amount: number // 分
  actual_amount: number   // 分
  diff_amount: number     // 分
  status: number
  idempotency_key: string
  operator_id: number
  created_at: string
}

// 告警处理 @ticket B端-告警处理
export interface Alarm {
  id: number
  site_id: number
  device_id: number
  level: string
  content: string
  status: number // 0=待处理 1=已处理
  handler: string
  handle_note: string
  idempotency_key: string
  created_at: string
  handled_at: string
}

// 权限点（与后端 rbac.go 对齐，用于按钮级显隐）
export const PERMISSIONS: Record<AdminRole, string[]> = {
  user: ['site:read'],
  site_admin: ['site:read', 'site:write', 'device:read', 'device:ctl', 'reconcile:read', 'alarm:read', 'alarm:handle'],
  operator: ['site:read', 'site:write', 'device:read', 'device:ctl', 'reconcile:read', 'reconcile:run', 'alarm:read', 'alarm:handle', 'audit:read'],
}
