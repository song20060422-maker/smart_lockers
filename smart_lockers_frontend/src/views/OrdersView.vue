<script setup lang="ts">
// @ticket 05-订单查看与取出：进行中聚合 + 游标分页订单列表
// @proto 原型对齐：v-orders 存取订单；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { orderList, ongoingOrders } from '@/api/order'
import type { OrderVO, OngoingSummary } from '@/types'

// 声明变量 router = useRouter()，用于保存数据或函数引用

const router = useRouter()
// 声明变量 list = ref<OrderVO[]>([])，用于保存数据或函数引用

const list = ref<OrderVO[]>([])
// 声明变量 ongoing = ref<OngoingSummary | null>(null)，用于保存数据或函数引用

const ongoing = ref<OngoingSummary | null>(null)
// 声明变量 cursor = ref('')，用于保存数据或函数引用

const cursor = ref('')
// 声明变量 hasMore = ref(false)，用于保存数据或函数引用

const hasMore = ref(false)
// 声明变量 activeType = ref('')，用于保存数据或函数引用

const activeType = ref('')
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')
// 声明变量 loading = ref(false)，用于保存数据或函数引用

const loading = ref(false)

// 声明变量 statusText: Record<number, string> = { 0: '待支付', 1: '进行中', 2: '已完成', 3: '已超时' }，用于保存数据或函数引用

const statusText: Record<number, string> = { 0: '待支付', 1: '进行中', 2: '已完成', 3: '已超时' }

// 调用 onMounted(() => { loadOngoing(); loadList() })，触发对应功能

onMounted(() => { loadOngoing(); loadList() })

// 定义函数/方法，封装可复用的处理逻辑：async function loadOngoing() {

async function loadOngoing() {
  // 等待异步操作完成：try { ongoing.value = await ongoingOrders() }

  try { ongoing.value = await ongoingOrders() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (_) { /* 非关键 */ }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：async function loadList(reset = false) {

async function loadList(reset = false) {
  // 条件判断：if (reset) { cursor.value = ''; list.value = [] } 时执行下面逻辑

  if (reset) { cursor.value = ''; list.value = [] }
  // 赋值更新：loading.value = true

  loading.value = true
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 r = await orderList({ type: activeType.value || undefined, cursor: cursor.value, limit: 10 })，用于保存数据或函数引用

    const r = await orderList({ type: activeType.value || undefined, cursor: cursor.value, limit: 10 })
    // 赋值更新：list.value = list.value.concat(r.list)

    list.value = list.value.concat(r.list)
    // 赋值更新：cursor.value = r.next_cursor

    cursor.value = r.next_cursor
    // 赋值更新：hasMore.value = r.has_more

    hasMore.value = r.has_more
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
  // 无论成功失败都会执行的收尾逻辑

  finally { loading.value = false }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function switchType(t: string) {

function switchType(t: string) {
  // 赋值更新：activeType.value = t

  activeType.value = t
  // 调用 loadList(true)，触发对应功能

  loadList(true)
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function openDetail(o: OrderVO) {

function openDetail(o: OrderVO) {
  // 执行 router.push(`/orders/${o.order_no}`)，实现对应的前端逻辑

  router.push(`/orders/${o.order_no}`)
// 当前代码块结束

}
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列进行中统计与订单列表卡片 -->

  <div class="col">
    <!-- 进行中统计卡片，class 为 card，v-if 依据 ongoing 是否有值决定是否渲染 -->

    <div v-if="ongoing" class="card">
      <!-- 二级标题，显示“进行中” -->

      <h2>进行中</h2>
      <!-- 段落元素，显示进行中订单数 ongoing.ongoing_count 与累计金额（ongoing.total_amount 除以 100 保留两位小数） -->

      <p>进行中订单：{{ ongoing.ongoing_count }} 单，累计金额 {{ (ongoing.total_amount / 100).toFixed(2) }} 元</p>
    <!-- 闭合进行中统计卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示订单列表与筛选按钮 -->

    <div class="card">
      <!-- 二级标题，显示“订单列表” -->

      <h2>订单列表</h2>
      <!-- 行容器，class 为 row，横向排列订单类型筛选按钮（本行与下一行同属一个标签） -->

      <div class="row">
        <!-- 基于类型二维数组循环生成筛选按钮，key 取 t[0]，逐项对应全部/进行中/已完成；动态绑定 class，activeType 不等于 t[0] 时添加 btn-ghost，点击触发 switchType(t[0]) -->

        <button v-for="t in [['', '全部'], ['ongoing', '进行中'], ['done', '已完成']]" :key="t[0]"
          :class="{ 'btn-ghost': activeType !== t[0] }" @click="switchType(t[0])">{{ t[1] }}</button>
      <!-- 闭合筛选按钮行容器 div 元素 -->

      </div>
      <!-- 基于 list 数组循环生成订单条目，key 取 o.order_no，样式类为 list-item，点击触发 openDetail(o) 查看详情 -->

      <div v-for="o in list" :key="o.order_no" class="list-item" @click="openDetail(o)">
        <!-- 行容器，class 为 row，横向排列订单号与状态标签 -->

        <div class="row">
          <!-- 加粗元素，插值 o.order_no 显示订单号 -->

          <strong>{{ o.order_no }}</strong>
          <!-- 行内标签元素，class 为 tag，插值 statusText[o.status]（缺省 o.status）显示订单状态 -->

          <span class="tag">{{ statusText[o.status] || o.status }}</span>
        <!-- 闭合订单行容器 div 元素 -->

        </div>
        <!-- 段落元素，class 为 muted，显示柜门 o.door_id、金额与押金（均除以 100 保留两位小数） -->

        <p class="muted">柜门 {{ o.door_id }} / 金额 {{ (o.amount / 100).toFixed(2) }}元 / 押金 {{ (o.deposit / 100).toFixed(2) }}元</p>
      <!-- 闭合当前订单条目 div 元素 -->

      </div>
      <!-- 空状态提示段落，当 list 为空且非加载中时渲染，class 为 muted，显示“暂无订单” -->

      <p v-if="!list.length && !loading" class="muted">暂无订单</p>
      <!-- 加载更多按钮，class 为 btn-ghost，当 hasMore 为真时渲染，loading 为真时禁用，点击触发 loadList() -->

      <button v-if="hasMore" class="btn-ghost" :disabled="loading" @click="loadList()">加载更多</button>
    <!-- 闭合订单列表卡片 div 元素 -->

    </div>
    <!-- 错误提示段落，当 err 有值时渲染，class 为 err，插值 err 显示错误信息 -->

    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
