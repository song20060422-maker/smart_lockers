<script setup lang="ts">
// @ticket 02-支付宝支付：余额/明细/充值/充值记录
// @proto 原型对齐：无直接对应(用户侧支付能力，原型未单列视图)；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { getBalance, balanceLogs, createRecharge, rechargeRecords } from '@/api/pay'
import type { BalanceResult, BalanceLog, RechargeOrder } from '@/types'

// 声明变量 balance = ref<BalanceResult | null>(null)，用于保存数据或函数引用

const balance = ref<BalanceResult | null>(null)
// 声明变量 logs = ref<BalanceLog[]>([])，用于保存数据或函数引用

const logs = ref<BalanceLog[]>([])
// 声明变量 records = ref<RechargeOrder[]>([])，用于保存数据或函数引用

const records = ref<RechargeOrder[]>([])
// 声明变量 amount = ref<number>(0)，用于保存数据或函数引用

const amount = ref<number>(0)
// 声明变量 msg = ref('')，用于保存数据或函数引用

const msg = ref('')
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')

// 调用 onMounted(() => { loadBalance(); loadLogs(); loadRecords() })，触发对应功能

onMounted(() => { loadBalance(); loadLogs(); loadRecords() })

// 定义函数/方法，封装可复用的处理逻辑：async function loadBalance() {

async function loadBalance() {
  // 等待异步操作完成：try { balance.value = await getBalance() }

  try { balance.value = await getBalance() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function loadLogs() {

async function loadLogs() {
  // 等待异步操作完成：try { const r = await balanceLogs({ page: 1, page_size: 20 }); logs.value = r.list }

  try { const r = await balanceLogs({ page: 1, page_size: 20 }); logs.value = r.list }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (_) {}
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function loadRecords() {

async function loadRecords() {
  // 等待异步操作完成：try { const r = await rechargeRecords({ page: 1, page_size: 20 }); records.value = r.list }

  try { const r = await rechargeRecords({ page: 1, page_size: 20 }); records.value = r.list }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (_) {}
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function doRecharge() {

async function doRecharge() {
  // 条件判断：if (!amount.value || amount.value <= 0) { err.value = '请输入有效金额（分）'; return } 时执行下面逻辑

  if (!amount.value || amount.value <= 0) { err.value = '请输入有效金额（分）'; return }
  // 条件判断：if (!window.confirm(`确认充值 ${(amount.value / 100).toFixed(2)} 元？`)) return 时执行下面逻辑

  if (!window.confirm(`确认充值 ${(amount.value / 100).toFixed(2)} 元？`)) return
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 idem = `recharge-${Date.now()}`，用于保存数据或函数引用

    const idem = `recharge-${Date.now()}`
    // 等待异步操作完成：await createRecharge(amount.value, idem)

    await createRecharge(amount.value, idem)
    // 赋值更新：msg.value = '充值单已创建'

    msg.value = '充值单已创建'
    // 等待异步操作完成：await loadBalance(); await loadRecords()

    await loadBalance(); await loadRecords()
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function fen(v: number): string { return (v / 100).toFixed(2) }

function fen(v: number): string { return (v / 100).toFixed(2) }
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列余额、明细与充值记录卡片 -->

  <div class="col">
    <!-- 卡片容器，class 为 card，展示账户余额与充值操作 -->

    <div class="card">
      <!-- 二级标题，显示“账户余额” -->

      <h2>账户余额</h2>
      <!-- 段落元素，当 balance 有值时渲染，显示余额并以加粗元素展示经 fen 格式化的金额与币种 balance.currency -->

      <p v-if="balance">余额：<strong>{{ fen(balance.balance) }} {{ balance.currency }}</strong></p>
      <!-- 加载提示段落，当 balance 无值时渲染，class 为 muted，显示“加载中...” -->

      <p v-else class="muted">加载中...</p>
      <!-- 行容器，class 为 row，横向排列充值金额输入框与按钮 -->

      <div class="row">
        <!-- 充值金额输入框，type 为 number，v-model.number 双向绑定 amount，最小值 1，占位提示“充值金额(分)” -->

        <input v-model.number="amount" type="number" min="1" placeholder="充值金额(分)" />
        <!-- 充值按钮，点击触发 doRecharge 执行充值 -->

        <button @click="doRecharge">充值</button>
      <!-- 闭合充值操作行容器 div 元素 -->

      </div>
      <!-- 提示段落，当 msg 有值时渲染，class 为 muted，插值 msg 显示提示内容 -->

      <p v-if="msg" class="muted">{{ msg }}</p>
    <!-- 闭合账户余额卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示余额变动明细表格 -->

    <div class="card">
      <!-- 二级标题，显示“余额明细” -->

      <h2>余额明细</h2>
      <!-- 表格元素，用于展示余额变动记录 -->

      <table>
        <!-- 表头区域，一行 tr 内依次定义时间、类型、变动、余额列名 th 单元格 -->

        <thead><tr><th>时间</th><th>类型</th><th>变动</th><th>余额</th></tr></thead>
        <!-- 表体区域，用于放置余额变动数据行 -->

        <tbody>
          <!-- 基于 logs 数组循环生成明细行，key 取 l.id -->

          <tr v-for="l in logs" :key="l.id">
            <!-- 两个单元格，依次插值显示创建时间 l.created_at 与业务类型 l.biz_type -->

            <td>{{ l.created_at }}</td><td>{{ l.biz_type }}</td>
            <!-- 两个单元格，分别通过 fen 格式化显示变动金额 l.change 与变动后余额 l.balance_after -->

            <td>{{ fen(l.change) }}</td><td>{{ fen(l.balance_after) }}</td>
          <!-- 闭合当前明细行 tr 元素 -->

          </tr>
        <!-- 闭合表体 tbody 元素 -->

        </tbody>
      <!-- 闭合表格 table 元素 -->

      </table>
      <!-- 空状态提示段落，当 logs 为空时渲染，class 为 muted，显示“暂无明细” -->

      <p v-if="!logs.length" class="muted">暂无明细</p>
    <!-- 闭合余额明细卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示充值记录列表 -->

    <div class="card">
      <!-- 二级标题，显示“充值记录” -->

      <h2>充值记录</h2>
      <!-- 基于 records 数组循环生成充值记录条目，key 取 r.id，样式类为 list-item -->

      <div v-for="r in records" :key="r.id" class="list-item">
        <!-- 行内文本元素，插值显示充值单号 r.charge_no、经 fen 格式化的金额 r.amount 与状态 r.status -->

        <span>{{ r.charge_no }} / {{ fen(r.amount) }}元 / 状态 {{ r.status }}</span>
      <!-- 闭合当前充值记录条目 div 元素 -->

      </div>
      <!-- 空状态提示段落，当 records 为空时渲染，class 为 muted，显示“暂无充值记录” -->

      <p v-if="!records.length" class="muted">暂无充值记录</p>
    <!-- 闭合充值记录卡片 div 元素 -->

    </div>
    <!-- 错误提示段落，当 err 有值时渲染，class 为 err，插值 err 显示错误信息 -->

    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
