<script setup lang="ts">
// @ticket 04-物品寄存：选柜/时长/优惠券/费用试算/落单/换柜
// @proto 原型对齐：v-fee 网点收费规则 / v-opendoor 异常换柜；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { queryCabinets, quoteFee, createDeposit, changeCabinet } from '@/api/deposit'
import { myCoupons } from '@/api/user'
import type { Cabinet, FeeResult, Coupon } from '@/types'

const route = useRoute()
const router = useRouter()
const siteId = Number(route.query.site_id || 0)
const cabinets = ref<Cabinet[]>([])
const coupons = ref<Coupon[]>([])
const doorId = ref<number>(0)
const hours = ref<number>(1)
const couponId = ref<number | undefined>(undefined)
const fee = ref<FeeResult | null>(null)
const orderNo = ref('')
const err = ref('')
const msg = ref('')

onMounted(async () => {
  try {
    cabinets.value = await queryCabinets({ site_id: siteId })
    coupons.value = await myCoupons()
  } catch (e) { err.value = (e as Error).message }
})

async function doQuote() {
  try {
    fee.value = await quoteFee({ size: 'M', hours: hours.value, coupon: couponId.value })
  } catch (e) { err.value = (e as Error).message }
}

// @ticket 04-物品寄存：敏感操作二次确认 + 幂等键防重复提交
async function submitOrder() {
  if (!doorId.value) { err.value = '请先选择柜门'; return }
  if (!window.confirm('确认按当前费用下单并支付？')) return
  try {
    const idem = `deposit-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    const r = await createDeposit({
      door_id: doorId.value, site_id: siteId, hours: hours.value, coupon: couponId.value,
    }, idem)
    orderNo.value = r.order_no
    msg.value = '下单成功'
  } catch (e) { err.value = (e as Error).message }
}

async function doChangeCabinet() {
  if (!orderNo.value) { err.value = '暂无可换柜的订单'; return }
  const newDoor = window.prompt('输入新柜门 ID')
  if (!newDoor) return
  try {
    const r = await changeCabinet(orderNo.value, Number(newDoor))
    doorId.value = r.door_id
    msg.value = `已换柜至 ${r.door_id}`
  } catch (e) { err.value = (e as Error).message }
}

function fen(v: number): string { return (v / 100).toFixed(2) }
</script>

<template>
  <div class="col">
    <div class="card">
      <h2>选择柜子</h2>
      <div v-for="c in cabinets" :key="c.id" class="list-item">
        <div class="row">
          <strong>{{ c.code }}</strong>
          <span class="tag">{{ c.size }}</span>
          <span class="muted">可用 {{ c.available }}/{{ c.total }}</span>
          <button :class="{ 'btn-ghost': doorId !== c.id }" @click="doorId = c.id; doQuote()">
            {{ doorId === c.id ? '已选' : '选择' }}
          </button>
        </div>
      </div>
      <p v-if="!cabinets.length" class="muted">该寄存点暂无可用柜子</p>
    </div>

    <div class="card">
      <h2>寄存设置</h2>
      <div class="row">
        <label>时长(小时) <input v-model.number="hours" type="number" min="1" @change="doQuote" /></label>
        <label>优惠券
          <select v-model.number="couponId" @change="doQuote">
            <option :value="undefined">不使用</option>
            <option v-for="c in coupons" :key="c.id" :value="c.id">{{ c.code }} -{{ fen(c.discount) }}元</option>
          </select>
        </label>
        <button @click="doQuote">试算费用</button>
      </div>
      <div v-if="fee" class="list-item">
        <p>基础费：{{ fen(fee.base_fee) }}元 / 优惠：{{ fen(fee.discount) }}元</p>
        <p>押金：{{ fen(fee.deposit) }}元 / 应付：<strong>{{ fen(fee.payable) }}元</strong></p>
      </div>
      <div class="row">
        <button :disabled="!fee" @click="submitOrder">下单并支付</button>
        <button class="btn-ghost" @click="doChangeCabinet">异常换柜</button>
      </div>
      <p v-if="msg" class="muted">{{ msg }}</p>
    </div>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
