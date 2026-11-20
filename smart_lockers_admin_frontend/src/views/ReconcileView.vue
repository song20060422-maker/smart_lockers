<template>
  <div>
    <div class="page-head"><h2>充值对账</h2>
      <button v-if="auth.can('reconcile:run')" class="btn primary" @click="showRun = true">+ 执行对账</button>
    </div>
    <div class="card">
      <table>
        <thead><tr><th>ID</th><th>批次号</th><th>网点ID</th><th>应收(元)</th><th>实收(元)</th><th>差额(元)</th><th>状态</th><th>幂等键</th><th>创建时间</th></tr></thead>
        <tbody>
          <tr v-for="r in rows" :key="r.id">
            <td>{{ r.id }}</td><td>{{ r.batch_no }}</td><td>{{ r.site_id }}</td>
            <td>{{ fen(r.expected_amount) }}</td><td>{{ fen(r.actual_amount) }}</td>
            <td :class="{ 'diff': r.diff_amount !== 0 }">{{ fen(r.diff_amount) }}</td>
            <td><span class="tag" :class="r.status === 1 ? 'ok' : 'warn'">{{ r.status === 1 ? '已平' : '未平' }}</span></td>
            <td class="muted">{{ r.idempotency_key }}</td><td>{{ r.created_at }}</td>
          </tr>
          <tr v-if="!rows.length"><td colspan="9" class="muted">暂无数据</td></tr>
        </tbody>
      </table>
    </div>

    <div v-if="showRun" class="modal">
      <form class="card" @submit.prevent="onRun">
        <h3>执行对账（仅平台运营）</h3>
        <label>批次号<input v-model="form.batch_no" required /></label>
        <label>应收金额（元）<input v-model.number="form.expected" type="number" step="0.01" required /></label>
        <label>实收金额（元）<input v-model.number="form.actual" type="number" step="0.01" required /></label>
        <div class="actions">
          <button class="btn primary" type="submit">执行</button>
          <button class="btn" type="button" @click="showRun = false">取消</button>
        </div>
      </form>
    </div>
    <div v-if="message" class="toast">{{ message }}</div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'
import type { ReconcileRecord } from '@/types/api'

const auth = useAuthStore()
const { message, notify } = useNotify()
const rows = ref<ReconcileRecord[]>([])
const showRun = ref(false)
const form = ref({ batch_no: '', expected: 0, actual: 0 })
const fen = (v: number) => (v / 100).toFixed(2)
async function load() {
  try { rows.value = (await api.listReconciles()).data || [] } catch (e: any) { notify(e.message) }
}
async function onRun() {
  try {
    const res = await api.runReconcile({
      batch_no: form.value.batch_no.trim(),
      expected_amount: Math.round(form.value.expected * 100),
      actual_amount: Math.round(form.value.actual * 100),
    })
    notify(res.data?.diff_amount ? `对账完成，差额 ${fen(res.data.diff_amount)} 元` : '对账完成，账目平衡')
    showRun.value = false
    form.value = { batch_no: '', expected: 0, actual: 0 }
    await load()
  } catch (e: any) { notify(e.message) }
}
onMounted(load)
</script>
<style scoped>
.muted { color: var(--muted); font-size: 12px; }
.diff { color: var(--danger); font-weight: 600; }
.modal { position: fixed; inset: 0; background: #0006; display: flex; align-items: center; justify-content: center; z-index: 50; }
.modal form { width: 380px; display: flex; flex-direction: column; gap: 12px; }
.modal label { display: flex; flex-direction: column; gap: 5px; font-size: 13px; color: var(--muted); }
.actions { display: flex; gap: 10px; justify-content: flex-end; }
</style>
