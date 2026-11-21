<template>
  <!-- 页面根容器 div，包裹设备监控页面的全部内容 -->

  <div>
    <!-- 页面头部容器，class 为 page-head，内含二级标题 h2 显示“设备监控” -->

    <div class="page-head"><h2>设备监控</h2>
      <!-- 统计信息区容器，class 为 stats，包裹各状态设备数量 -->

      <div class="stats">
        <!-- 行内元素，显示“总计”并通过加粗 b 元素插值 monitor.total 得到设备总数 -->

        <span>总计 <b>{{ monitor.total }}</b></span>
        <!-- 行内标签元素，class 为 tag ok，显示“在线”并插值 monitor.online 得到在线设备数 -->

        <span class="tag ok">在线 {{ monitor.online }}</span>
        <!-- 行内标签元素，class 为 tag warn，显示“离线”并插值 monitor.offline 得到离线设备数 -->

        <span class="tag warn">离线 {{ monitor.offline }}</span>
        <!-- 行内标签元素，class 为 tag err，显示“故障”并插值 monitor.fault 得到故障设备数 -->

        <span class="tag err">故障 {{ monitor.fault }}</span>
      <!-- 闭合统计信息区 div 元素 -->

      </div>
    <!-- 闭合页面头部 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，包裹设备数据表格 -->

    <div class="card">
      <!-- 表格元素，用于展示设备记录列表 -->

      <table>
        <!-- 表头区域，一行 tr 内依次定义 ID、编码、网点ID、类型、状态、最后心跳、操作等列名 th 单元格 -->

        <thead><tr><th>ID</th><th>编码</th><th>网点ID</th><th>类型</th><th>状态</th><th>最后心跳</th><th>操作</th></tr></thead>
        <!-- 表体区域，用于放置设备数据行 -->

        <tbody>
          <!-- 基于 rows 列表循环生成表格行，key 取 d.id 作为唯一标识 -->

          <tr v-for="d in rows" :key="d.id">
            <!-- 一行内的四个单元格，依次插值显示设备 ID、编码 d.code、网点ID d.site_id 与类型 d.type -->

            <td>{{ d.id }}</td><td>{{ d.code }}</td><td>{{ d.site_id }}</td><td>{{ d.type }}</td>
            <!-- 状态单元格，内含标签 span，class 为 tag 并通过 statusClass(d.status) 动态计算样式类，插值 statusText(d.status) 显示状态文本 -->

            <td><span class="tag" :class="statusClass(d.status)">{{ statusText(d.status) }}</span></td>
            <!-- 单元格，插值 d.last_heartbeat 显示设备最后心跳时间 -->

            <td>{{ d.last_heartbeat }}</td>
            <!-- 操作列单元格起始，内部放置状态控制下拉框或占位文本 -->

            <td>
              <!-- 修复偏差5：权限点与后端 device:control 对齐，原 device:ctl 永远为 false 导致控制框不渲染 @ticket B端-前端 -->

              <select v-if="auth.can('device:control')" :value="d.status" @change="onControl(d, $event)">
                <!-- 三个下拉选项，值分别为 0（在线）、1（离线）、2（故障） -->

                <option :value="0">在线</option><option :value="1">离线</option><option :value="2">故障</option>
              <!-- 闭合状态控制下拉框 select 元素 -->

              </select>
              <!-- v-if 的备选分支，无控制权限时显示占位文本“—”，class 为 muted -->

              <span v-else class="muted">—</span>
            <!-- 闭合操作列单元格 -->

            </td>
          <!-- 闭合当前数据行 tr 元素 -->

          </tr>
          <!-- 空数据占位行，当 rows 长度为 0 时渲染，单元格 colspan 为 7 横跨整行，显示“暂无数据” -->

          <tr v-if="!rows.length"><td colspan="7" class="muted">暂无数据</td></tr>
        <!-- 闭合表体 tbody 元素 -->

        </tbody>
      <!-- 闭合表格 table 元素 -->

      </table>
    <!-- 闭合卡片 div 元素 -->

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
import type { Device, DeviceMonitor } from '@/types/api'

// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 赋值更新：const { message, notify } = useNotify()

const { message, notify } = useNotify()
// 声明变量 rows = ref<Device[]>([])，用于保存数据或函数引用

const rows = ref<Device[]>([])
// 声明变量 monitor = ref<DeviceMonitor>({ total: 0, online: 0, offline: 0, fault: 0, items: [] })，用于保存数据或函数引用

const monitor = ref<DeviceMonitor>({ total: 0, online: 0, offline: 0, fault: 0, items: [] })
// 定义函数/方法，封装可复用的处理逻辑：const statusText = (s: number) => ({ 0: '在线', 1: '离线', 2: '故障' }[s] || '未知')

const statusText = (s: number) => ({ 0: '在线', 1: '离线', 2: '故障' }[s] || '未知')
// 定义函数/方法，封装可复用的处理逻辑：const statusClass = (s: number) => s === 0 ? 'ok' : s === 1 ? 'warn' : 'err'

const statusClass = (s: number) => s === 0 ? 'ok' : s === 1 ? 'warn' : 'err'
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 赋值更新：rows.value = (await api.listDevices()).data || []

    rows.value = (await api.listDevices()).data || []
    // 赋值更新：monitor.value = (await api.monitor()).data || monitor.value

    monitor.value = (await api.monitor()).data || monitor.value
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function onControl(d: Device, e: Event) {

async function onControl(d: Device, e: Event) {
  // 声明变量 status = Number((e.target as HTMLSelectElement).value)，用于保存数据或函数引用

  const status = Number((e.target as HTMLSelectElement).value)
  // 等待异步操作完成：try { await api.controlDevice(d.id, status); notify('控制指令已下发'); await load() }

  try { await api.controlDevice(d.id, status); notify('控制指令已下发'); await load() }
  // 等待异步操作完成：catch (err: any) { notify(err.message); await load() }

  catch (err: any) { notify(err.message); await load() }
// 当前代码块结束

}
// 调用 onMounted(load)，触发对应功能

onMounted(load)
</script>
<style scoped>
.stats { display: flex; gap: 12px; align-items: center; }
.stats b { font-size: 16px; }
.muted { color: var(--muted); }
</style>
