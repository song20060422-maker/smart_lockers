<template>
  <div>
    <div class="page-head"><h2>告警处理</h2>
      <span class="muted">待处理 {{ pendingCount }} 条</span>
    </div>
    <div class="card">
      <table>
        <thead><tr><th>ID</th><th>级别</th><th>网点ID</th><th>设备ID</th><th>内容</th><th>状态</th><th>处理人</th><th>处理备注</th><th>操作</th></tr></thead>
        <tbody>
          <tr v-for="a in rows" :key="a.id">
            <td>{{ a.id }}</td>
            <td><span class="tag" :class="levelClass(a.level)">{{ a.level }}</span></td>
            <td>{{ a.site_id }}</td><td>{{ a.device_id }}</td><td>{{ a.content }}</td>
            <td><span class="tag" :class="a.status === 1 ? 'ok' : 'err'">{{ a.status === 1 ? '已处理' : '待处理' }}</span></td>
            <td>{{ a.handler || '—' }}</td><td>{{ a.handle_note || '—' }}</td>
            <td>
              <button v-if="a.status === 0 && auth.can('alarm:handle')" class="btn primary" @click="openHandle(a)">处理</button>
              <span v-else class="muted">—</span>
            </td>
          </tr>
          <tr v-if="!rows.length"><td colspan="9" class="muted">暂无数据</td></tr>
        </tbody>
      </table>
    </div>

    <div v-if="target" class="modal">
      <form class="card" @submit.prevent="onHandle">
        <h3>处理告警 #{{ target.id }}</h3>
        <p class="muted">{{ target.content }}</p>
        <label>处理备注<textarea v-model="note" rows="3" /></label>
        <div class="actions">
          <button class="btn primary" type="submit">确认处理</button>
          <button class="btn" type="button" @click="target = null">取消</button>
        </div>
      </form>
    </div>
    <div v-if="message" class="toast">{{ message }}</div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'
import type { Alarm } from '@/types/api'

const auth = useAuthStore()
const { message, notify } = useNotify()
const rows = ref<Alarm[]>([])
const target = ref<Alarm | null>(null)
const note = ref('')
const pendingCount = computed(() => rows.value.filter((a) => a.status === 0).length)
const levelClass = (l: string) => l === 'critical' ? 'err' : l === 'warn' ? 'warn' : 'ok'
async function load() {
  try { rows.value = (await api.listAlarms()).data || [] } catch (e: any) { notify(e.message) }
}
function openHandle(a: Alarm) { target.value = a; note.value = '' }
async function onHandle() {
  if (!target.value) return
  try { await api.handleAlarm(target.value.id, note.value); notify('处理成功'); target.value = null; await load() }
  catch (e: any) { notify(e.message) }
}
onMounted(load)
</script>
<style scoped>
.muted { color: var(--muted); }
.modal { position: fixed; inset: 0; background: #0006; display: flex; align-items: center; justify-content: center; z-index: 50; }
.modal form { width: 400px; display: flex; flex-direction: column; gap: 12px; }
.modal label { display: flex; flex-direction: column; gap: 5px; font-size: 13px; color: var(--muted); }
textarea { padding: 8px; border: 1px solid var(--border); border-radius: 6px; resize: vertical; }
.actions { display: flex; gap: 10px; justify-content: flex-end; }
</style>
