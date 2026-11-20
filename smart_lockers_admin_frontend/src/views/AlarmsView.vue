<template>
  <!-- 页面根容器 div，包裹告警处理页面的全部内容 -->

  <div>
    <!-- 页面头部容器，class 为 page-head，内含二级标题 h2 显示“告警处理” -->

    <div class="page-head"><h2>告警处理</h2>
      <!-- 行内文本元素，class 为 muted，显示“待处理”并插值 pendingCount 得到待处理告警条数 -->

      <span class="muted">待处理 {{ pendingCount }} 条</span>
    <!-- 闭合页面头部 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，包裹告警数据表格 -->

    <div class="card">
      <!-- 表格元素，用于展示告警记录列表 -->

      <table>
        <!-- 表头区域，一行 tr 内依次定义 ID、级别、网点ID、设备ID、内容、状态、处理人、处理备注、操作等列名 th 单元格 -->

        <thead><tr><th>ID</th><th>级别</th><th>网点ID</th><th>设备ID</th><th>内容</th><th>状态</th><th>处理人</th><th>处理备注</th><th>操作</th></tr></thead>
        <!-- 表体区域，用于放置告警数据行 -->

        <tbody>
          <!-- 基于 rows 列表循环生成表格行，key 取 a.id 作为唯一标识 -->

          <tr v-for="a in rows" :key="a.id">
            <!-- 单元格，插值 a.id 显示告警记录 ID -->

            <td>{{ a.id }}</td>
            <!-- 单元格，内含标签 span，class 为 tag 并通过 levelClass(a.level) 动态计算级别样式类，插值 a.level 显示告警级别文本 -->

            <td><span class="tag" :class="levelClass(a.level)">{{ a.level }}</span></td>
            <!-- 一行内的三个单元格，依次插值显示网点ID（a.site_id）、设备ID（a.device_id）与告警内容（a.content） -->

            <td>{{ a.site_id }}</td><td>{{ a.device_id }}</td><td>{{ a.content }}</td>
            <!-- 状态单元格，内含标签 span，当 a.status 等于 1 时样式类为 ok、文本“已处理”，否则样式类为 err、文本“待处理” -->

            <td><span class="tag" :class="a.status === 1 ? 'ok' : 'err'">{{ a.status === 1 ? '已处理' : '待处理' }}</span></td>
            <!-- 一行内的两个单元格，依次插值显示处理人 a.handler 和处理备注 a.handle_note，为空时显示破折号“—” -->

            <td>{{ a.handler || '—' }}</td><td>{{ a.handle_note || '—' }}</td>
            <!-- 操作列单元格起始，内部放置“处理”按钮或占位文本 -->

            <td>
              <!-- “处理”按钮，当告警状态为 0（待处理）且认证对象具备 alarm:handle 权限时才渲染，class 为 btn primary，点击触发 openHandle(a) 打开处理弹窗 -->

              <button v-if="a.status === 0 && auth.can('alarm:handle')" class="btn primary" @click="openHandle(a)">处理</button>
              <!-- v-if 的备选分支，当不满足处理条件时显示占位文本“—”，class 为 muted -->

              <span v-else class="muted">—</span>
            <!-- 闭合操作列单元格 -->

            </td>
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

    <!-- 处理告警弹窗遮罩容器，v-if 依据 target 是否有值决定是否渲染，class 为 modal -->

    <div v-if="target" class="modal">
      <!-- 处理表单，class 为 card，提交事件被 prevent 修饰阻止默认行为并调用 onHandle 方法 -->

      <form class="card" @submit.prevent="onHandle">
        <!-- 三级标题，显示“处理告警 #”并插值 target.id 得到当前告警编号 -->

        <h3>处理告警 #{{ target.id }}</h3>
        <!-- 段落元素，class 为 muted，插值 target.content 显示待处理告警的内容文本 -->

        <p class="muted">{{ target.content }}</p>
        <!-- 表单标签，包含文本“处理备注”与多行文本框 textarea，v-model 双向绑定 note 变量，rows 为 3 -->

        <label>处理备注<textarea v-model="note" rows="3" /></label>
        <!-- 操作按钮区容器，class 为 actions，包裹确认与取消按钮 -->

        <div class="actions">
          <!-- 提交按钮，class 为 btn primary，type 为 submit，文本“确认处理” -->

          <button class="btn primary" type="submit">确认处理</button>
          <!-- 取消按钮，class 为 btn，type 为 button，点击将 target 置为 null 关闭弹窗，文本“取消” -->

          <button class="btn" type="button" @click="target = null">取消</button>
        <!-- 闭合操作按钮区 div 元素 -->

        </div>
      <!-- 闭合处理表单 form 元素 -->

      </form>
    <!-- 闭合弹窗遮罩 div 元素 -->

    </div>
    <!-- 全局轻提示条，v-if 依据 message 是否有值决定是否渲染，class 为 toast，插值 message 显示提示内容 -->

    <div v-if="message" class="toast">{{ message }}</div>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { api } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'
import type { Alarm } from '@/types/api'

// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 赋值更新：const { message, notify } = useNotify()

const { message, notify } = useNotify()
// 声明变量 rows = ref<Alarm[]>([])，用于保存数据或函数引用

const rows = ref<Alarm[]>([])
// 声明变量 target = ref<Alarm | null>(null)，用于保存数据或函数引用

const target = ref<Alarm | null>(null)
// 声明变量 note = ref('')，用于保存数据或函数引用

const note = ref('')
// 声明变量 pendingCount = computed(() => rows.value.filter((a) => a.status === 0).length)，用于保存数据或函数引用

const pendingCount = computed(() => rows.value.filter((a) => a.status === 0).length)
// 定义函数/方法，封装可复用的处理逻辑：const levelClass = (l: string) => l === 'critical' ? 'err' : l === 'warn' ? 'warn' : 'ok'

const levelClass = (l: string) => l === 'critical' ? 'err' : l === 'warn' ? 'warn' : 'ok'
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try { rows.value = (await api.listAlarms()).data || [] } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function openHandle(a: Alarm) { target.value = a; note.value = '' }

function openHandle(a: Alarm) { target.value = a; note.value = '' }
// 定义函数/方法，封装可复用的处理逻辑：async function onHandle() {

async function onHandle() {
  // 条件判断：if (!target.value) return 时执行下面逻辑

  if (!target.value) return
  // 等待异步操作完成：try { await api.handleAlarm(target.value.id, note.value); notify('处理成功'); target.value = null; await load() }

  try { await api.handleAlarm(target.value.id, note.value); notify('处理成功'); target.value = null; await load() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 调用 onMounted(load)，触发对应功能

onMounted(load)
</script>
<style scoped>
.muted { color: var(--muted); }
.modal { position: fixed; inset: 0; background: #0006; display: flex; align-items: center; justify-content: center; z-index: 50; }
.modal form { width: 400px; display: flex; flex-direction: column; gap: 12px; }
.modal label { display: flex; flex-direction: column; gap: 5px; font-size: 13px; color: var(--muted); }
textarea { padding: 8px; border: 1px solid var(--border); border-radius: 6px; resize: vertical; }
.actions { display: flex; gap: 10px; justify-content: flex-end; }
</style>
