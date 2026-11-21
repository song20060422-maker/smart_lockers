// @ticket 01-消息通知：发送/记录/监控/场景接口封装
import { get, post } from './request'
import type { NotifySendResult, NotifyRecordsResult, MonitorResult, Scene } from '@/types'

// @ticket 01-消息通知：POST /api/v1/notify/send（鉴权 + 幂等）
// 导出 function sendNotify(payload: RecordOfUnknown, idem: string): Promise<NotifySendResult> {，供其他模块引用

export function sendNotify(payload: RecordOfUnknown, idem: string): Promise<NotifySendResult> {
  // 返回 post('/notify/send', payload, idem) 给调用方

  return post('/notify/send', payload, idem)
// 当前代码块结束

}

// @ticket 01-消息通知：GET /api/v1/notify/records（鉴权）
// 导出 function notifyRecords(params: RecordOfUnknown): Promise<NotifyRecordsResult> {，供其他模块引用

export function notifyRecords(params: RecordOfUnknown): Promise<NotifyRecordsResult> {
  // 返回 get('/notify/records', params) 给调用方

  return get('/notify/records', params)
// 当前代码块结束

}

// @ticket 01-消息通知：GET /api/v1/notify/monitor（鉴权）
// 导出 function notifyMonitor(): Promise<MonitorResult> {，供其他模块引用

export function notifyMonitor(): Promise<MonitorResult> {
  // 返回 get('/notify/monitor') 给调用方

  return get('/notify/monitor')
// 当前代码块结束

}

// @ticket 01-消息通知：POST /api/v1/notify/scene（鉴权，按场景模板发送）
// 导出 function sendScene(payload: RecordOfUnknown, idem: string): Promise<NotifySendResult> {，供其他模块引用

export function sendScene(payload: RecordOfUnknown, idem: string): Promise<NotifySendResult> {
  // 返回 post('/notify/scene', payload, idem) 给调用方

  return post('/notify/scene', payload, idem)
// 当前代码块结束

}

// @ticket 01-消息通知：GET /api/v1/notify/scenes（鉴权）
// 导出 function notifyScenes(): Promise<Scene[]> {，供其他模块引用

export function notifyScenes(): Promise<Scene[]> {
  // 返回 get('/notify/scenes') 给调用方

  return get('/notify/scenes')
// 当前代码块结束

}

// 定义类型别名，描述数据的形状：type RecordOfUnknown = Record<string, unknown>

type RecordOfUnknown = Record<string, unknown>
