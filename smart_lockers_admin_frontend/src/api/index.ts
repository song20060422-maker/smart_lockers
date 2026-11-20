// @ticket B端-前端：对齐后端 router.go 全部接口
import http, { idem } from './request'
import type {
  AdminUser, Site, Device, DeviceMonitor, ReconcileRecord, Alarm,
} from '@/types/api'

export const api = {
  // 登录鉴权
  login: (username: string, password: string) =>
    http.post<{ token: string; user: AdminUser }>('/auth/login', { username, password }),
  me: () => http.get<AdminUser>('/auth/me'),

  // 网点管理
  listSites: () => http.get<Site[]>('/admin/sites'),
  createSite: (body: Partial<Site>) => http.post<Site>('/admin/sites', body, idem(`site-${Date.now()}`)),
  updateSiteStatus: (id: number, status: number) =>
    http.post<Site>(`/admin/sites/${id}/status`, { status }, idem(`site-status-${id}-${Date.now()}`)),

  // 设备监控
  listDevices: () => http.get<Device[]>('/admin/devices'),
  monitor: () => http.get<DeviceMonitor>('/admin/devices/monitor'),
  controlDevice: (id: number, status: number) =>
    http.post<Device>(`/admin/devices/${id}/control`, { status }, idem(`device-${id}-${Date.now()}`)),

  // 充值对账
  listReconciles: () => http.get<ReconcileRecord[]>('/admin/reconcile'),
  runReconcile: (body: { batch_no: string; expected_amount: number; actual_amount: number }) =>
    http.post<ReconcileRecord>('/admin/reconcile/run', body, idem(`recon-${body.batch_no}`)),

  // 告警处理
  listAlarms: () => http.get<Alarm[]>('/admin/alarms'),
  handleAlarm: (id: number, note: string) =>
    http.post<Alarm>(`/admin/alarms/${id}/handle`, { note }, idem(`alarm-${id}-${Date.now()}`)),
}
