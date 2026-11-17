<script setup lang="ts">
// @ticket 02-支付宝支付：余额/明细/充值/充值记录
import { ref, onMounted } from 'vue'
import { getBalance, balanceLogs, createRecharge, rechargeRecords } from '@/api/pay'
import type { BalanceResult, BalanceLog, RechargeOrder } from '@/types'

const balance = ref<BalanceResult | null>(null)
const logs = ref<BalanceLog[]>([])
const records = ref<RechargeOrder[]>([])
const amount = ref<number>(0)
const msg = ref('')
const err = ref('')

onMounted(() => { loadBalance(); loadLogs(); loadRecords() })

async function loadBalance() {
  try { balance.value = await getBalance() }
  catch (e) { err.value = (e as Error).message }
}
async function loadLogs() {
  try { const r = await balanceLogs({ page: 1, page_size: 20 }); logs.value = r.list }
  catch (_) {}
}
async function loadRecords() {
  try { const r = await rechargeRecords({ page: 1, page_size: 20 }); records.value = r.list }
  catch (_) {}
}
async function doRecharge() {
  if (!amount.value || amount.value <= 0) { err.value = '请输入有效金额（分）'; return }
  if (!window.confirm(`确认充值 ${(amount.value / 100).toFixed(2)} 元？`)) return
  try {
    const idem = `recharge-${Date.now()}`
    await createRecharge(amount.value, idem)
    msg.value = '充值单已创建'
    await loadBalance(); await loadRecords()
  } catch (e) { err.value = (e as Error).message }
}
function fen(v: number): string { return (v / 100).toFixed(2) }
</script>

<template>
  <div class="col">
    <div class="card">
      <h2>账户余额</h2>
      <p v-if="balance">余额：<strong>{{ fen(balance.balance) }} {{ balance.currency }}</strong></p>
      <p v-else class="muted">加载中...</p>
      <div class="row">
        <input v-model.number="amount" type="number" min="1" placeholder="充值金额(分)" />
        <button @click="doRecharge">充值</button>
      </div>
      <p v-if="msg" class="muted">{{ msg }}</p>
    </div>
    <div class="card">
      <h2>余额明细</h2>
      <table>
        <thead><tr><th>时间</th><th>类型</th><th>变动</th><th>余额</th></tr></thead>
        <tbody>
          <tr v-for="l in logs" :key="l.id">
            <td>{{ l.created_at }}</td><td>{{ l.biz_type }}</td>
            <td>{{ fen(l.change) }}</td><td>{{ fen(l.balance_after) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-if="!logs.length" class="muted">暂无明细</p>
    </div>
    <div class="card">
      <h2>充值记录</h2>
      <div v-for="r in records" :key="r.id" class="list-item">
        <span>{{ r.charge_no }} / {{ fen(r.amount) }}元 / 状态 {{ r.status }}</span>
      </div>
      <p v-if="!records.length" class="muted">暂无充值记录</p>
    </div>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
