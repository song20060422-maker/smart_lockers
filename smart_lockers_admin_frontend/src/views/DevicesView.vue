<template>
  <div>
    <div class="page-head"><h2>设备监控</h2>
      <div class="stats">
        <span>总计 <b>{{ monitor.total }}</b></span>
        <span class="tag ok">在线 {{ monitor.online }}</span>
        <span class="tag warn">离线 {{ monitor.offline }}</span>
        <span class="tag err">故障 {{ monitor.fault }}</span>
      </div>
    </div>
    <div class="card">
      <table>
        <thead><tr><th>ID</th><th>编码</th><th>网点ID</th><th>类型</th><th>状态</th><th>最后心跳</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="d in rows" :key="d.id">
            <td>{{ d.id }}</td><td>{{ d.code }}</td><td>{{ d.site_id }}</td><td>{{ d.type }}</td>
            <td><span class="tag" :class="statusClass(d.status)">{{ statusText(d.status) }}</span></td>
            <td>{{ d.last_heartbeat }}</td>
            <td>
              <select v-if="auth.can('device:ctl')" :value="d.status" @change="onControl(d, $event)">
                <option :value="0">在线</option><option :value="1">离线</option><option :value="2">故障</option>
              </select>
              <span v-else class="muted">—</span>
            </td>
          </tr>
          <tr v-if="!rows.length"><td colspan="7" class="muted">暂无数据</td></tr>
        </tbody>
      </table>
    </div>
    <div v-if="message" class="toast">{{ message }}</div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'
import type { Device, DeviceMonitor } from '@/types/api'

const auth = useAuthStore()
const { message, notify } = useNotify()
const rows = ref<Device[]>([])
const monitor = ref<DeviceMonitor>({ total: 0, online: 0, offline: 0, fault: 0, items: [] })
const statusText = (s: number) => ({ 0: '在线', 1: '离线', 2: '故障' }[s] || '未知')
const statusClass = (s: number) => s === 0 ? 'ok' : s === 1 ? 'warn' : 'err'
async function load() {
  try {
    rows.value = (await api.listDevices()).data || []
    monitor.value = (await api.monitor()).data || monitor.value
  } catch (e: any) { notify(e.message) }
}
async function onControl(d: Device, e: Event) {
  const status = Number((e.target as HTMLSelectElement).value)
  try { await api.controlDevice(d.id, status); notify('控制指令已下发'); await load() }
  catch (err: any) { notify(err.message); await load() }
}
onMounted(load)
</script>
<style scoped>
.stats { display: flex; gap: 12px; align-items: center; }
.stats b { font-size: 16px; }
.muted { color: var(--muted); }
</style>
