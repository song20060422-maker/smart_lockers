<template>
  <!-- 页面根容器 div，包裹充值对账页面的全部内容 -->

  <div>
    <!-- 页面头部容器，class 为 page-head，内含二级标题 h2 显示“充值对账” -->

    <div class="page-head"><h2>充值对账</h2>
      <!-- “执行对账”按钮，当认证对象具备 reconcile:run 权限时渲染，class 为 btn primary，点击将 showRun 置为 true 打开弹窗 -->

      <button v-if="auth.can('reconcile:run')" class="btn primary" @click="showRun = true">+ 执行对账</button>
    <!-- 闭合页面头部 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，包裹对账数据表格 -->

    <div class="card">
      <!-- 表格元素，用于展示对账记录列表 -->

      <table>
        <!-- 表头区域，一行 tr 内依次定义 ID、批次号、网点ID、应收(元)、实收(元)、差额(元)、状态、幂等键、创建时间等列名 th 单元格 -->

        <thead><tr><th>ID</th><th>批次号</th><th>网点ID</th><th>应收(元)</th><th>实收(元)</th><th>差额(元)</th><th>状态</th><th>幂等键</th><th>创建时间</th></tr></thead>
        <!-- 表体区域，用于放置对账数据行 -->

        <tbody>
          <!-- 基于 rows 列表循环生成表格行，key 取 r.id 作为唯一标识 -->

          <tr v-for="r in rows" :key="r.id">
            <!-- 一行内的三个单元格，依次插值显示对账记录 ID、批次号 r.batch_no 与网点ID r.site_id -->

            <td>{{ r.id }}</td><td>{{ r.batch_no }}</td><td>{{ r.site_id }}</td>
            <!-- 两个单元格，分别通过 fen 函数格式化显示应收金额 r.expected_amount 与实收金额 r.actual_amount -->

            <td>{{ fen(r.expected_amount) }}</td><td>{{ fen(r.actual_amount) }}</td>
            <!-- 差额单元格，当 r.diff_amount 不为 0 时动态添加 diff 样式类，插值经 fen 格式化的差额金额 -->

            <td :class="{ 'diff': r.diff_amount !== 0 }">{{ fen(r.diff_amount) }}</td>
            <!-- 状态单元格，内含标签 span，当 r.status 等于 1 时样式类为 ok、文本“已平”，否则样式类为 warn、文本“未平” -->

            <td><span class="tag" :class="r.status === 1 ? 'ok' : 'warn'">{{ r.status === 1 ? '已平' : '未平' }}</span></td>
            <!-- 两个单元格，依次插值显示幂等键 r.idempotency_key（class 为 muted）与创建时间 r.created_at -->

            <td class="muted">{{ r.idempotency_key }}</td><td>{{ r.created_at }}</td>
          <!-- 闭合当前数据行 tr 元素 -->

          </tr>
          <!-- 空数据占位行，当 rows 长度为 0 时渲染，单元格 colspan 为 9 横跨整行，显示“暂无数据” -->

          <tr v-if="!rows.length"><td colspan="9" class="muted">暂无数据</td></tr>
        <!-- 闭合表体 tbody 元素 -->

        </tbody>
      <!-- 闭合表格 table 元素 -->

      </table>
    <!-- 闭合卡片 div 元素 -->

    </div>

    <!-- 执行对账弹窗遮罩容器，v-if 依据 showRun 是否有值决定是否渲染，class 为 modal -->

    <div v-if="showRun" class="modal">
      <!-- 对账表单，class 为 card，提交事件被 prevent 修饰阻止默认行为并调用 onRun 方法 -->

      <form class="card" @submit.prevent="onRun">
        <!-- 三级标题，显示“执行对账（仅平台运营）” -->

        <h3>执行对账（仅平台运营）</h3>
        <!-- 表单标签，包含文本“批次号”与必填输入框，v-model 双向绑定 form.batch_no -->

        <label>批次号<input v-model="form.batch_no" required /></label>
        <!-- 表单标签，包含“应收金额（元）”与必填数字输入框，v-model.number 双向绑定 form.expected，步长 0.01 -->

        <label>应收金额（元）<input v-model.number="form.expected" type="number" step="0.01" required /></label>
        <!-- 表单标签，包含“实收金额（元）”与必填数字输入框，v-model.number 双向绑定 form.actual，步长 0.01 -->

        <label>实收金额（元）<input v-model.number="form.actual" type="number" step="0.01" required /></label>
        <!-- 操作按钮区容器，class 为 actions，包裹执行与取消按钮 -->

        <div class="actions">
          <!-- 提交按钮，class 为 btn primary，type 为 submit，文本“执行” -->

          <button class="btn primary" type="submit">执行</button>
          <!-- 取消按钮，class 为 btn，type 为 button，点击将 showRun 置为 false 关闭弹窗，文本“取消” -->

          <button class="btn" type="button" @click="showRun = false">取消</button>
        <!-- 闭合操作按钮区 div 元素 -->

        </div>
      <!-- 闭合对账表单 form 元素 -->

      </form>
    <!-- 闭合弹窗遮罩 div 元素 -->

    </div>
    <!-- 全局轻提示条，v-if 依据 message 是否有值决定是否渲染，class 为 toast，插值 message 显示提示内容 -->

    <div v-if="message" class="toast">{{ message }}</div>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'
import type { ReconcileRecord } from '@/types/api'

// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 赋值更新：const { message, notify } = useNotify()

const { message, notify } = useNotify()
// 声明变量 rows = ref<ReconcileRecord[]>([])，用于保存数据或函数引用

const rows = ref<ReconcileRecord[]>([])
// 声明变量 showRun = ref(false)，用于保存数据或函数引用

const showRun = ref(false)
// 声明变量 form = ref({ batch_no: '', expected: 0, actual: 0 })，用于保存数据或函数引用

const form = ref({ batch_no: '', expected: 0, actual: 0 })
// 定义函数/方法，封装可复用的处理逻辑：const fen = (v: number) => (v / 100).toFixed(2)

const fen = (v: number) => (v / 100).toFixed(2)
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try { rows.value = (await api.listReconciles()).data || [] } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function onRun() {

async function onRun() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 res = await api.runReconcile({，用于保存数据或函数引用

    const res = await api.runReconcile({
      // 定义对象属性 batch_no，取值为 form.value.batch_no.trim()

      batch_no: form.value.batch_no.trim(),
      // 定义对象属性 expected_amount，取值为 Math.round(form.value.expected * 100)

      expected_amount: Math.round(form.value.expected * 100),
      // 定义对象属性 actual_amount，取值为 Math.round(form.value.actual * 100)

      actual_amount: Math.round(form.value.actual * 100),
    // 执行 })，实现对应的前端逻辑

    })
    // 调用 notify(res.data?.diff_amount ? `对账完成，差额 ${fen(res.data.diff_amount)} 元` : '对账完成，账目平衡')，触发对应功能

    notify(res.data?.diff_amount ? `对账完成，差额 ${fen(res.data.diff_amount)} 元` : '对账完成，账目平衡')
    // 赋值更新：showRun.value = false

    showRun.value = false
    // 赋值更新：form.value = { batch_no: '', expected: 0, actual: 0 }

    form.value = { batch_no: '', expected: 0, actual: 0 }
    // 等待异步操作完成：await load()

    await load()
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 调用 onMounted(load)，触发对应功能

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
