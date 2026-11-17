// @ticket 01-消息通知：发送/记录/监控/场景接口封装
import { get, post } from './request'
import type { NotifySendResult, NotifyRecordsResult, MonitorResult, Scene } from '@/types'

// @ticket 01-消息通知：POST /api/v1/notify/send（鉴权 + 幂等）
export function sendNotify(payload: RecordOfUnknown, idem: string): Promise<NotifySendResult> {
  return post('/notify/send', payload, idem)
}

// @ticket 01-消息通知：GET /api/v1/notify/records（鉴权）
export function notifyRecords(params: RecordOfUnknown): Promise<NotifyRecordsResult> {
  return get('/notify/records', params)
}

// @ticket 01-消息通知：GET /api/v1/notify/monitor（鉴权）
export function notifyMonitor(): Promise<MonitorResult> {
  return get('/notify/monitor')
}

// @ticket 01-消息通知：POST /api/v1/notify/scene（鉴权，按场景模板发送）
export function sendScene(payload: RecordOfUnknown, idem: string): Promise<NotifySendResult> {
  return post('/notify/scene', payload, idem)
}

// @ticket 01-消息通知：GET /api/v1/notify/scenes（鉴权）
export function notifyScenes(): Promise<Scene[]> {
  return get('/notify/scenes')
}

type RecordOfUnknown = Record<string, unknown>
