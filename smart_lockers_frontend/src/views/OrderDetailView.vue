<script setup lang="ts">
// @ticket 05-订单查看与取出：详情/临时开柜/超时补交/完成
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { orderDetail, openAuth, makeup, completeOrder } from '@/api/order'
import type { OrderDetail, OpenAuthResult } from '@/types'

const route = useRoute()
const orderNo = String(route.params.orderNo)
const detail = ref<OrderDetail | null>(null)
const openResult = ref<OpenAuthResult | null>(null)
const makeupAmount = ref<number>(0)
const msg = ref('')
const err = ref('')

const statusText: Record<number, string> = { 0: '待支付', 1: '进行中', 2: '已完成', 3: '已超时' }

onMounted(async () => {
  try { detail.value = await orderDetail(orderNo) }
  catch (e) { err.value = (e as Error).message }
})

// @ticket 05-订单查看与取出：临时开柜需二次确认
async function doOpenAuth() {
  if (!window.confirm('确认临时开柜？')) return
  try {
    openResult.value = await openAuth(orderNo, { max_times: 3, ttl_sec: 300 })
    msg.value = '开柜令牌已生成'
  } catch (e) { err.value = (e as Error).message }
}

async function doMakeup() {
  if (!makeupAmount.value) { err.value = '请输入补交金额（分）'; return }
  if (!window.confirm('确认补交？')) return
  try {
    const idem = `makeup-${Date.now()}`
    await makeup(orderNo, makeupAmount.value, idem)
    msg.value = '补交成功'
    detail.value = await orderDetail(orderNo)
  } catch (e) { err.value = (e as Error).message }
}

async function doComplete() {
  if (!window.confirm('确认完成订单？')) return
  try {
    const idem = `complete-${Date.now()}`
    await completeOrder(orderNo, idem)
    msg.value = '订单已完成'
    detail.value = await orderDetail(orderNo)
  } catch (e) { err.value = (e as Error).message }
}
</script>

<template>
  <div class="col">
    <div v-if="detail" class="card">
      <h2>订单 {{ detail.order_no }}</h2>
      <p>状态：<span class="tag">{{ statusText[detail.status] || detail.status }}</span></p>
      <p>柜门：{{ detail.door_id }} / 金额：{{ (detail.amount / 100).toFixed(2) }}元 / 押金：{{ (detail.deposit / 100).toFixed(2) }}元</p>
      <p class="muted">开始 {{ detail.start_time }} / 结束 {{ detail.end_time }}</p>
      <div v-if="detail.site" class="list-item">
        <p>{{ detail.site.name }}</p>
        <p class="muted">{{ detail.site.address }}</p>
      </div>
      <div class="row">
        <button @click="doOpenAuth">临时开柜</button>
        <button @click="doComplete">完成订单</button>
      </div>
    </div>
    <div class="card">
      <h2>超时补交</h2>
      <div class="row">
        <input v-model.number="makeupAmount" type="number" min="0" placeholder="金额(分)" />
        <button @click="doMakeup">补交</button>
      </div>
    </div>
    <div v-if="openResult" class="card">
      <h2>开柜令牌</h2>
      <p>令牌：{{ openResult.token }}</p>
      <p class="muted">有效期至 {{ openResult.expire_at }} / 最多 {{ openResult.max_times }} 次</p>
    </div>
    <p v-if="msg" class="muted">{{ msg }}</p>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
