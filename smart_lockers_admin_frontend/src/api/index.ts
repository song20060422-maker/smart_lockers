// @ticket B端-前端：对齐后端 router.go 全部接口
import http, { idem } from './request'
import type {
  AdminUser, Site, Device, DeviceMonitor, ReconcileRecord, Alarm,
} from '@/types/api'

// 导出 const api = {，供其他模块引用

export const api = {
  // 登录鉴权
  // 定义对象属性 login，取值为 (username: string, password: string) =>

  login: (username: string, password: string) =>
    // 执行 http.post<{ token: string; user: AdminUser }>('/auth/login', { username, password }),，实现对应的前端逻辑

    http.post<{ token: string; user: AdminUser }>('/auth/login', { username, password }),
  // 定义对象属性 me，取值为 () => http.get<AdminUser>('/auth/me')

  me: () => http.get<AdminUser>('/auth/me'),

  // 网点管理
  // 定义对象属性 listSites，取值为 () => http.get<Site[]>('/admin/sites')

  listSites: () => http.get<Site[]>('/admin/sites'),
  // 定义对象属性 createSite，取值为 (body: Partial<Site>) => http.post<Site>('/admin/sites', body, idem(`site-${Date.now()}`))

  createSite: (body: Partial<Site>) => http.post<Site>('/admin/sites', body, idem(`site-${Date.now()}`)),
  // 定义对象属性 updateSiteStatus，取值为 (id: number, status: number) =>

  updateSiteStatus: (id: number, status: number) =>
    // 执行 http.post<Site>(`/admin/sites/${id}/status`, { status }, idem(`site-status-${id}-${Date.now()}`)),，实现对应的前端逻辑

    http.post<Site>(`/admin/sites/${id}/status`, { status }, idem(`site-status-${id}-${Date.now()}`)),

  // 设备监控
  // 定义对象属性 listDevices，取值为 () => http.get<Device[]>('/admin/devices')

  listDevices: () => http.get<Device[]>('/admin/devices'),
  // 定义对象属性 monitor，取值为 () => http.get<DeviceMonitor>('/admin/devices/monitor')

  monitor: () => http.get<DeviceMonitor>('/admin/devices/monitor'),
  // 定义对象属性 controlDevice，取值为 (id: number, status: number) =>

  controlDevice: (id: number, status: number) =>
    // 执行 http.post<Device>(`/admin/devices/${id}/control`, { status }, idem(`device-${id}-${Date.now()}`)),，实现对应的前端逻辑

    http.post<Device>(`/admin/devices/${id}/control`, { status }, idem(`device-${id}-${Date.now()}`)),

  // 充值对账
  // 定义对象属性 listReconciles，取值为 () => http.get<ReconcileRecord[]>('/admin/reconcile')

  listReconciles: () => http.get<ReconcileRecord[]>('/admin/reconcile'),
  // 定义对象属性 runReconcile，取值为 (body: { batch_no: string; expected_amount: number; actual_amount: number }) =>

  runReconcile: (body: { batch_no: string; expected_amount: number; actual_amount: number }) =>
    // 执行 http.post<ReconcileRecord>('/admin/reconcile/run', body, idem(`recon-${body.batch_no}`)),，实现对应的前端逻辑

    http.post<ReconcileRecord>('/admin/reconcile/run', body, idem(`recon-${body.batch_no}`)),

  // 告警处理
  // 定义对象属性 listAlarms，取值为 () => http.get<Alarm[]>('/admin/alarms')

  listAlarms: () => http.get<Alarm[]>('/admin/alarms'),
  // 定义对象属性 handleAlarm，取值为 (id: number, note: string) =>

  handleAlarm: (id: number, note: string) =>
    // 执行 http.post<Alarm>(`/admin/alarms/${id}/handle`, { note }, idem(`alarm-${id}-${Date.now()}`)),，实现对应的前端逻辑

    http.post<Alarm>(`/admin/alarms/${id}/handle`, { note }, idem(`alarm-${id}-${Date.now()}`)),
// 当前代码块结束

}
