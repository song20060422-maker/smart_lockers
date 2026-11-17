<script setup lang="ts">
// @ticket 05-订单查看与取出：进行中聚合 + 游标分页订单列表
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { orderList, ongoingOrders } from '@/api/order'
import type { OrderVO, OngoingSummary } from '@/types'

const router = useRouter()
const list = ref<OrderVO[]>([])
const ongoing = ref<OngoingSummary | null>(null)
const cursor = ref('')
const hasMore = ref(false)
const activeType = ref('')
const err = ref('')
const loading = ref(false)

const statusText: Record<number, string> = { 0: '待支付', 1: '进行中', 2: '已完成', 3: '已超时' }

onMounted(() => { loadOngoing(); loadList() })

async function loadOngoing() {
  try { ongoing.value = await ongoingOrders() }
  catch (_) { /* 非关键 */ }
}

async function loadList(reset = false) {
  if (reset) { cursor.value = ''; list.value = [] }
  loading.value = true
  try {
    const r = await orderList({ type: activeType.value || undefined, cursor: cursor.value, limit: 10 })
    list.value = list.value.concat(r.list)
    cursor.value = r.next_cursor
    hasMore.value = r.has_more
  } catch (e) { err.value = (e as Error).message }
  finally { loading.value = false }
}

function switchType(t: string) {
  activeType.value = t
  loadList(true)
}

function openDetail(o: OrderVO) {
  router.push(`/orders/${o.order_no}`)
}
</script>

<template>
  <div class="col">
    <div v-if="ongoing" class="card">
      <h2>进行中</h2>
      <p>进行中订单：{{ ongoing.ongoing_count }} 单，累计金额 {{ (ongoing.total_amount / 100).toFixed(2) }} 元</p>
    </div>
    <div class="card">
      <h2>订单列表</h2>
      <div class="row">
        <button v-for="t in [['', '全部'], ['ongoing', '进行中'], ['done', '已完成']]" :key="t[0]"
          :class="{ 'btn-ghost': activeType !== t[0] }" @click="switchType(t[0])">{{ t[1] }}</button>
      </div>
      <div v-for="o in list" :key="o.order_no" class="list-item" @click="openDetail(o)">
        <div class="row">
          <strong>{{ o.order_no }}</strong>
          <span class="tag">{{ statusText[o.status] || o.status }}</span>
        </div>
        <p class="muted">柜门 {{ o.door_id }} / 金额 {{ (o.amount / 100).toFixed(2) }}元 / 押金 {{ (o.deposit / 100).toFixed(2) }}元</p>
      </div>
      <p v-if="!list.length && !loading" class="muted">暂无订单</p>
      <button v-if="hasMore" class="btn-ghost" :disabled="loading" @click="loadList()">加载更多</button>
    </div>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
