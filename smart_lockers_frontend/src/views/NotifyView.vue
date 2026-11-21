<script setup lang="ts">
// @ticket 01-消息通知：发送/记录/监控/场景
// @proto 原型对齐：notify 通知中心 / alarm 告警中心；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { sendNotify, sendScene, notifyRecords, notifyMonitor, notifyScenes } from '@/api/notify'
import type { NotifyRecordsResult, MonitorResult, Scene } from '@/types'

// 声明变量 records = ref<NotifyRecordsResult | null>(null)，用于保存数据或函数引用

const records = ref<NotifyRecordsResult | null>(null)
// 声明变量 monitor = ref<MonitorResult | null>(null)，用于保存数据或函数引用

const monitor = ref<MonitorResult | null>(null)
// 声明变量 scenes = ref<Scene[]>([])，用于保存数据或函数引用

const scenes = ref<Scene[]>([])
// 声明变量 sendResult = ref('')，用于保存数据或函数引用

const sendResult = ref('')
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')

// 声明变量 form = ref({ tpl_code: 'notify_test', receiver_type: 'user', receiver_id: 10001, channel: 'station', content: '测试通知' })，用于保存数据或函数引用

const form = ref({ tpl_code: 'notify_test', receiver_type: 'user', receiver_id: 10001, channel: 'station', content: '测试通知' })

// 调用 onMounted(() => { loadRecords(); loadMonitor(); loadScenes() })，触发对应功能

onMounted(() => { loadRecords(); loadMonitor(); loadScenes() })
// 定义函数/方法，封装可复用的处理逻辑：async function loadRecords() {

async function loadRecords() {
  // 等待异步操作完成：try { records.value = await notifyRecords({ page: 1, page_size: 20 }) }

  try { records.value = await notifyRecords({ page: 1, page_size: 20 }) }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function loadMonitor() {

async function loadMonitor() {
  // 等待异步操作完成：try { monitor.value = await notifyMonitor() }

  try { monitor.value = await notifyMonitor() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (_) {}
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function loadScenes() {

async function loadScenes() {
  // 等待异步操作完成：try { scenes.value = await notifyScenes() }

  try { scenes.value = await notifyScenes() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (_) {}
// 当前代码块结束

}
// @ticket 01-消息通知：发送走幂等键防止重复
// 定义函数/方法，封装可复用的处理逻辑：async function doSend() {

async function doSend() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 idem = `notify-${Date.now()}`，用于保存数据或函数引用

    const idem = `notify-${Date.now()}`
    // 等待异步操作完成：await sendNotify(form.value, idem)

    await sendNotify(form.value, idem)
    // 赋值更新：sendResult.value = '发送成功'

    sendResult.value = '发送成功'
    // 等待异步操作完成：await loadRecords(); await loadMonitor()

    await loadRecords(); await loadMonitor()
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
// @ticket 01-消息通知：按场景模板发送
// 定义函数/方法，封装可复用的处理逻辑：async function doSendScene() {

async function doSendScene() {
  // 声明变量 code = window.prompt('输入场景编码（如 order_paid）')，用于保存数据或函数引用

  const code = window.prompt('输入场景编码（如 order_paid）')
  // 条件判断：if (!code) return 时执行下面逻辑

  if (!code) return
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 等待异步操作完成：await sendScene({ scene_code: code, receiver_id: form.value.receiver_id }, `scene-${Date.now()}`)

    await sendScene({ scene_code: code, receiver_id: form.value.receiver_id }, `scene-${Date.now()}`)
    // 赋值更新：sendResult.value = `场景 ${code} 已发送`

    sendResult.value = `场景 ${code} 已发送`
    // 等待异步操作完成：await loadRecords()

    await loadRecords()
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列通知相关卡片 -->

  <div class="col">
    <!-- 通知监控卡片，class 为 card，v-if 依据 monitor 是否有值决定是否渲染 -->

    <div v-if="monitor" class="card">
      <!-- 二级标题，显示“通知监控” -->

      <h2>通知监控</h2>
      <!-- 段落元素，依次插值 monitor.pending、success、failed、retrying 显示各状态通知数量 -->

      <p>待发 {{ monitor.pending }} / 成功 {{ monitor.success }} / 失败 {{ monitor.failed }} / 重试 {{ monitor.retrying }}</p>
    <!-- 闭合通知监控卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示发送通知表单 -->

    <div class="card">
      <!-- 二级标题，显示“发送通知” -->

      <h2>发送通知</h2>
      <!-- 行容器，class 为 row，横向排列模板编码、接收者、内容与发送按钮 -->

      <div class="row">
        <!-- 模板编码输入框，v-model 双向绑定 form.tpl_code，占位提示“模板编码” -->

        <input v-model="form.tpl_code" placeholder="模板编码" />
        <!-- 接收者 ID 输入框，type 为 number，v-model.number 双向绑定 form.receiver_id，占位提示“接收者 ID” -->

        <input v-model.number="form.receiver_id" type="number" placeholder="接收者 ID" />
        <!-- 内容输入框，v-model 双向绑定 form.content，占位提示“内容” -->

        <input v-model="form.content" placeholder="内容" />
        <!-- 发送按钮，点击触发 doSend 发送通知 -->

        <button @click="doSend">发送</button>
        <!-- 按场景发送按钮，class 为 btn-ghost，点击触发 doSendScene -->

        <button class="btn-ghost" @click="doSendScene">按场景发送</button>
      <!-- 闭合发送表单行容器 div 元素 -->

      </div>
      <!-- 发送结果段落，当 sendResult 有值时渲染，class 为 muted，插值 sendResult 显示结果文本 -->

      <p v-if="sendResult" class="muted">{{ sendResult }}</p>
    <!-- 闭合发送通知卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示可用通知场景 -->

    <div class="card">
      <!-- 二级标题，显示“通知场景” -->

      <h2>通知场景</h2>
      <!-- 基于 scenes 数组循环生成场景标签，key 取 s.code，样式类为 tag，插值 s.name 显示场景名称 -->

      <span v-for="s in scenes" :key="s.code" class="tag">{{ s.name }}</span>
    <!-- 闭合通知场景卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示通知发送记录 -->

    <div class="card">
      <!-- 二级标题，显示“通知记录” -->

      <h2>通知记录</h2>
      <!-- 基于 records?.list 数组循环生成记录条目，key 取 n.id，样式类为 list-item -->

      <div v-for="n in records?.list || []" :key="n.id" class="list-item">
        <!-- 段落元素，插值 n.channel 与 n.content 显示通知渠道与内容 -->

        <p>{{ n.channel }} / {{ n.content }}</p>
        <!-- 段落元素，class 为 muted，显示状态 n.status、重试次数 n.retry_count 与创建时间 n.created_at -->

        <p class="muted">状态 {{ n.status }} / 重试 {{ n.retry_count }} 次 / {{ n.created_at }}</p>
      <!-- 闭合当前通知记录条目 div 元素 -->

      </div>
      <!-- 空状态提示段落，当 records?.list 为空时渲染，class 为 muted，显示“暂无记录” -->

      <p v-if="!records?.list?.length" class="muted">暂无记录</p>
    <!-- 闭合通知记录卡片 div 元素 -->

    </div>
    <!-- 错误提示段落，当 err 有值时渲染，class 为 err，插值 err 显示错误信息 -->

    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
