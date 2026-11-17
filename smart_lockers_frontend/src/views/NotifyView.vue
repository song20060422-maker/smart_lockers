<script setup lang="ts">
// @ticket 01-消息通知：发送/记录/监控/场景
import { ref, onMounted } from 'vue'
import { sendNotify, sendScene, notifyRecords, notifyMonitor, notifyScenes } from '@/api/notify'
import type { NotifyRecordsResult, MonitorResult, Scene } from '@/types'

const records = ref<NotifyRecordsResult | null>(null)
const monitor = ref<MonitorResult | null>(null)
const scenes = ref<Scene[]>([])
const sendResult = ref('')
const err = ref('')

const form = ref({ tpl_code: 'notify_test', receiver_type: 'user', receiver_id: 10001, channel: 'station', content: '测试通知' })

onMounted(() => { loadRecords(); loadMonitor(); loadScenes() })
async function loadRecords() {
  try { records.value = await notifyRecords({ page: 1, page_size: 20 }) }
  catch (e) { err.value = (e as Error).message }
}
async function loadMonitor() {
  try { monitor.value = await notifyMonitor() }
  catch (_) {}
}
async function loadScenes() {
  try { scenes.value = await notifyScenes() }
  catch (_) {}
}
// @ticket 01-消息通知：发送走幂等键防止重复
async function doSend() {
  try {
    const idem = `notify-${Date.now()}`
    await sendNotify(form.value, idem)
    sendResult.value = '发送成功'
    await loadRecords(); await loadMonitor()
  } catch (e) { err.value = (e as Error).message }
}
// @ticket 01-消息通知：按场景模板发送
async function doSendScene() {
  const code = window.prompt('输入场景编码（如 order_paid）')
  if (!code) return
  try {
    await sendScene({ scene_code: code, receiver_id: form.value.receiver_id }, `scene-${Date.now()}`)
    sendResult.value = `场景 ${code} 已发送`
    await loadRecords()
  } catch (e) { err.value = (e as Error).message }
}
</script>

<template>
  <div class="col">
    <div v-if="monitor" class="card">
      <h2>通知监控</h2>
      <p>待发 {{ monitor.pending }} / 成功 {{ monitor.success }} / 失败 {{ monitor.failed }} / 重试 {{ monitor.retrying }}</p>
    </div>
    <div class="card">
      <h2>发送通知</h2>
      <div class="row">
        <input v-model="form.tpl_code" placeholder="模板编码" />
        <input v-model.number="form.receiver_id" type="number" placeholder="接收者 ID" />
        <input v-model="form.content" placeholder="内容" />
        <button @click="doSend">发送</button>
        <button class="btn-ghost" @click="doSendScene">按场景发送</button>
      </div>
      <p v-if="sendResult" class="muted">{{ sendResult }}</p>
    </div>
    <div class="card">
      <h2>通知场景</h2>
      <span v-for="s in scenes" :key="s.code" class="tag">{{ s.name }}</span>
    </div>
    <div class="card">
      <h2>通知记录</h2>
      <div v-for="n in records?.list || []" :key="n.id" class="list-item">
        <p>{{ n.channel }} / {{ n.content }}</p>
        <p class="muted">状态 {{ n.status }} / 重试 {{ n.retry_count }} 次 / {{ n.created_at }}</p>
      </div>
      <p v-if="!records?.list?.length" class="muted">暂无记录</p>
    </div>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
