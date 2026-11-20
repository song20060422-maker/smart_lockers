<script setup lang="ts">
// @ticket 05-订单查看与取出：详情/临时开柜/超时补交/完成
// @proto 原型对齐：v-orders 存取订单详情 / v-opendoor 临时开柜；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { orderDetail, openAuth, makeup, completeOrder } from '@/api/order'
import type { OrderDetail, OpenAuthResult } from '@/types'

// 声明变量 route = useRoute()，用于保存数据或函数引用

const route = useRoute()
// 声明变量 orderNo = String(route.params.orderNo)，用于保存数据或函数引用

const orderNo = String(route.params.orderNo)
// 声明变量 detail = ref<OrderDetail | null>(null)，用于保存数据或函数引用

const detail = ref<OrderDetail | null>(null)
// 声明变量 openResult = ref<OpenAuthResult | null>(null)，用于保存数据或函数引用

const openResult = ref<OpenAuthResult | null>(null)
// 声明变量 makeupAmount = ref<number>(0)，用于保存数据或函数引用

const makeupAmount = ref<number>(0)
// 声明变量 msg = ref('')，用于保存数据或函数引用

const msg = ref('')
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')

// 声明变量 statusText: Record<number, string> = { 0: '待支付', 1: '进行中', 2: '已完成', 3: '已超时' }，用于保存数据或函数引用

const statusText: Record<number, string> = { 0: '待支付', 1: '进行中', 2: '已完成', 3: '已超时' }

// 调用 onMounted(async () => {，触发对应功能

onMounted(async () => {
  // 等待异步操作完成：try { detail.value = await orderDetail(orderNo) }

  try { detail.value = await orderDetail(orderNo) }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e) { err.value = (e as Error).message }
// 执行 })，实现对应的前端逻辑

})

// @ticket 05-订单查看与取出：临时开柜需二次确认
// 定义函数/方法，封装可复用的处理逻辑：async function doOpenAuth() {

async function doOpenAuth() {
  // 条件判断：if (!window.confirm('确认临时开柜？')) return 时执行下面逻辑

  if (!window.confirm('确认临时开柜？')) return
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 等待异步操作完成：openResult.value = await openAuth(orderNo, { max_times: 3, ttl_sec: 300 })

    openResult.value = await openAuth(orderNo, { max_times: 3, ttl_sec: 300 })
    // 赋值更新：msg.value = '开柜令牌已生成'

    msg.value = '开柜令牌已生成'
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：async function doMakeup() {

async function doMakeup() {
  // 条件判断：if (!makeupAmount.value) { err.value = '请输入补交金额（分）'; return } 时执行下面逻辑

  if (!makeupAmount.value) { err.value = '请输入补交金额（分）'; return }
  // 条件判断：if (!window.confirm('确认补交？')) return 时执行下面逻辑

  if (!window.confirm('确认补交？')) return
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 idem = `makeup-${Date.now()}`，用于保存数据或函数引用

    const idem = `makeup-${Date.now()}`
    // 等待异步操作完成：await makeup(orderNo, makeupAmount.value, idem)

    await makeup(orderNo, makeupAmount.value, idem)
    // 赋值更新：msg.value = '补交成功'

    msg.value = '补交成功'
    // 等待异步操作完成：detail.value = await orderDetail(orderNo)

    detail.value = await orderDetail(orderNo)
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：async function doComplete() {

async function doComplete() {
  // 条件判断：if (!window.confirm('确认完成订单？')) return 时执行下面逻辑

  if (!window.confirm('确认完成订单？')) return
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 idem = `complete-${Date.now()}`，用于保存数据或函数引用

    const idem = `complete-${Date.now()}`
    // 等待异步操作完成：await completeOrder(orderNo, idem)

    await completeOrder(orderNo, idem)
    // 赋值更新：msg.value = '订单已完成'

    msg.value = '订单已完成'
    // 等待异步操作完成：detail.value = await orderDetail(orderNo)

    detail.value = await orderDetail(orderNo)
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列订单详情相关卡片 -->

  <div class="col">
    <!-- 订单详情卡片，class 为 card，v-if 依据 detail 是否有值决定是否渲染 -->

    <div v-if="detail" class="card">
      <!-- 二级标题，显示“订单”并插值 detail.order_no 显示订单号 -->

      <h2>订单 {{ detail.order_no }}</h2>
      <!-- 段落元素，显示状态：以内联标签 span 插值 statusText[detail.status]（缺省 detail.status）显示状态文本 -->

      <p>状态：<span class="tag">{{ statusText[detail.status] || detail.status }}</span></p>
      <!-- 段落元素，显示柜门 detail.door_id、金额与押金（金额除以 100 后保留两位小数） -->

      <p>柜门：{{ detail.door_id }} / 金额：{{ (detail.amount / 100).toFixed(2) }}元 / 押金：{{ (detail.deposit / 100).toFixed(2) }}元</p>
      <!-- 段落元素，class 为 muted，显示开始时间 detail.start_time 与结束时间 detail.end_time -->

      <p class="muted">开始 {{ detail.start_time }} / 结束 {{ detail.end_time }}</p>
      <!-- 网点信息容器，class 为 list-item，v-if 依据 detail.site 是否有值决定是否渲染 -->

      <div v-if="detail.site" class="list-item">
        <!-- 段落元素，插值 detail.site.name 显示网点名称 -->

        <p>{{ detail.site.name }}</p>
        <!-- 段落元素，class 为 muted，插值 detail.site.address 显示网点地址 -->

        <p class="muted">{{ detail.site.address }}</p>
      <!-- 闭合网点信息容器 div 元素 -->

      </div>
      <!-- 行容器，class 为 row，横向排列订单操作按钮 -->

      <div class="row">
        <!-- 临时开柜按钮，点击触发 doOpenAuth 获取开柜授权 -->

        <button @click="doOpenAuth">临时开柜</button>
        <!-- 完成订单按钮，点击触发 doComplete 完成当前订单 -->

        <button @click="doComplete">完成订单</button>
      <!-- 闭合订单操作行容器 div 元素 -->

      </div>
    <!-- 闭合订单详情卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示超时补交表单 -->

    <div class="card">
      <!-- 二级标题，显示“超时补交” -->

      <h2>超时补交</h2>
      <!-- 行容器，class 为 row，横向排列补交金额输入框与按钮 -->

      <div class="row">
        <!-- 补交金额输入框，type 为 number，v-model.number 双向绑定 makeupAmount，最小值 0，占位提示“金额(分)” -->

        <input v-model.number="makeupAmount" type="number" min="0" placeholder="金额(分)" />
        <!-- 补交按钮，点击触发 doMakeup 执行补交 -->

        <button @click="doMakeup">补交</button>
      <!-- 闭合补交表单行容器 div 元素 -->

      </div>
    <!-- 闭合超时补交卡片 div 元素 -->

    </div>
    <!-- 开柜令牌卡片，class 为 card，v-if 依据 openResult 是否有值决定是否渲染 -->

    <div v-if="openResult" class="card">
      <!-- 二级标题，显示“开柜令牌” -->

      <h2>开柜令牌</h2>
      <!-- 段落元素，显示“令牌：”并插值 openResult.token 显示令牌值 -->

      <p>令牌：{{ openResult.token }}</p>
      <!-- 段落元素，class 为 muted，显示令牌有效期 openResult.expire_at 与最多使用次数 openResult.max_times -->

      <p class="muted">有效期至 {{ openResult.expire_at }} / 最多 {{ openResult.max_times }} 次</p>
    <!-- 闭合开柜令牌卡片 div 元素 -->

    </div>
    <!-- ===== 补齐入口：订单相关原型界面跳转（只新增，不改动上方原有内容） ===== -->
    <!-- 更多操作卡片容器，class 为 card，放置取出/异常/补交等界面的跳转入口 -->
    <div class="card">
      <!-- 二级标题，显示“更多操作” -->
      <h2>更多操作</h2>
      <!-- 行容器，class 为 row，横向排列各操作入口链接 -->
      <div class="row">
        <!-- 链接：取出类型选择，跳转 /pickup -->
        <router-link to="/pickup" class="tag" style="text-decoration:none">取出类型</router-link>
        <!-- 链接：开门异常提示，跳转 /door-exception -->
        <router-link to="/door-exception" class="tag" style="text-decoration:none">开门异常</router-link>
        <!-- 链接：超时补交，跳转 /supplement -->
        <router-link to="/supplement" class="tag" style="text-decoration:none">超时补交</router-link>
      <!-- 闭合入口链接行容器 div 元素 -->
      </div>
    <!-- 闭合更多操作卡片 div 元素 -->
    </div>

    <!-- 提示段落，当 msg 有值时渲染，class 为 muted，插值 msg 显示提示内容 -->

    <p v-if="msg" class="muted">{{ msg }}</p>
    <!-- 错误提示段落，当 err 有值时渲染，class 为 err，插值 err 显示错误信息 -->

    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
