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
        <!-- 表头区域：新增 SN/固件列与 MQTT 指令操作列 @ticket B端-IoT中台 -->

        <thead><tr><th>ID</th><th>编码</th><th>设备SN</th><th>固件</th><th>网点ID</th><th>类型</th><th>状态</th><th>最后心跳</th><th>操作</th><th>MQTT指令</th></tr></thead>
        <!-- 表体区域，用于放置设备数据行 -->

        <tbody>
          <!-- 基于 rows 列表循环生成表格行，key 取 d.id 作为唯一标识 -->

          <tr v-for="d in rows" :key="d.id">
            <!-- 一行内的基础信息单元格：ID/编码/SN/固件/网点ID/类型 -->

            <td>{{ d.id }}</td><td>{{ d.code }}</td><td class="muted">{{ d.device_sn || '—' }}</td><td class="muted">{{ d.firmware || '—' }}</td><td>{{ d.site_id }}</td><td>{{ d.type }}</td>
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
            <!-- MQTT 指令列：开门/关门/重启快捷指令（对齐原型 devices 视图"下发控制指令"）@ticket B端-IoT中台 -->

            <td>
              <!-- 指令按钮组，ota:manage 权限控制显隐 -->

              <span v-if="auth.can('ota:manage')" class="cmds">
                <!-- 开门指令按钮 -->

                <button class="btn" @click="onCommand(d, 'open_door')">开门</button>
                <!-- 关门指令按钮 -->

                <button class="btn" @click="onCommand(d, 'close_door')">关门</button>
                <!-- 重启指令按钮 -->

                <button class="btn" @click="onCommand(d, 'reboot')">重启</button>
              <!-- 闭合按钮组 -->

              </span>
              <!-- 无权限占位 -->

              <span v-else class="muted">—</span>
            <!-- 闭合 MQTT 指令列 -->

            </td>
          <!-- 闭合当前数据行 tr 元素 -->

          </tr>
          <!-- 空数据占位行，当 rows 长度为 0 时渲染，单元格 colspan 为 10 横跨整行，显示“暂无数据” -->

          <tr v-if="!rows.length"><td colspan="10" class="muted">暂无数据</td></tr>
        <!-- 闭合表体 tbody 元素 -->

        </tbody>
      <!-- 闭合表格 table 元素 -->

      </table>
    <!-- 闭合卡片 div 元素 -->

    </div>
    <!-- 指令响应记录卡片（对齐原型 devices 视图"指令响应记录 device/response/{id}/…"）@ticket B端-IoT中台 -->

    <div class="card">
      <!-- 卡片标题与刷新按钮 -->

      <div class="sec-h"><h3>指令响应记录</h3><button class="btn" @click="loadCmds">刷新</button></div>
      <!-- 指令记录表格 -->

      <table>
        <!-- 表头：指令号/设备/类型/柜格/状态/响应/触发人/时间 -->

        <thead><tr><th>指令号</th><th>设备</th><th>类型</th><th>柜格</th><th>状态</th><th>设备响应</th><th>触发人</th><th>时间</th></tr></thead>
        <!-- 表体：循环指令记录 -->

        <tbody>
          <tr v-for="c in cmds" :key="c.id">
            <!-- 指令号（等宽展示） -->

            <td><code>{{ c.command_id }}</code></td>
            <!-- 设备 SN -->

            <td>{{ c.device_sn || c.device_id }}</td>
            <!-- 指令类型 -->

            <td>{{ c.cmd_type }}</td>
            <!-- 柜格编号 -->

            <td>{{ c.locker_no || '—' }}</td>
            <!-- 指令状态标签 -->

            <td><span class="tag" :class="cmdClass(c.status)">{{ cmdText(c.status) }}</span></td>
            <!-- 设备响应报文（截断展示） -->

            <td class="muted resp">{{ c.response || '—' }}</td>
            <!-- 触发人 -->

            <td>{{ c.operator }}</td>
            <!-- 下发时间 -->

            <td class="muted">{{ c.created_at }}</td>
          <!-- 闭合行 -->

          </tr>
          <!-- 空数据占位 -->

          <tr v-if="!cmds.length"><td colspan="8" class="muted">暂无指令记录</td></tr>
        <!-- 闭合表体 -->

        </tbody>
      <!-- 闭合表格 -->

      </table>
    <!-- 闭合卡片 -->

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
import type { Device, DeviceMonitor, DeviceCommand } from '@/types/api'

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

// ============ IoT MQTT 指令（@ticket B端-IoT中台） ============
// 声明变量 cmds = ref<DeviceCommand[]>([])，指令响应记录列表

const cmds = ref<DeviceCommand[]>([])
// 指令状态文案：0已下发 1成功 2失败 3超时/未发布

const cmdText = (s: number) => ({ 0: '已下发', 1: '成功', 2: '失败', 3: '超时' }[s] || '未知')
// 指令状态样式

const cmdClass = (s: number) => ({ 0: 'warn', 1: 'ok', 2: 'err', 3: 'err' }[s] || 'muted')
// 定义函数/方法 loadCmds：刷新指令响应记录

async function loadCmds() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 赋值更新：cmds.value = (await api.iotCommands()).data || []

    cmds.value = (await api.iotCommands()).data || []
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 定义函数/方法 onCommand：下发 MQTT 控制指令（open_door/close_door/reboot）

async function onCommand(d: Device, cmdType: string) {
  // 等待异步操作完成：try { await api.iotDispatch(d.id, cmdType) }

  try {
    // 调用下发接口

    const res: any = await api.iotDispatch(d.id, cmdType)
    // 成功提示：published=false 时提示演示模式（Broker 未连接，指令仅落库）

    notify(res.data?.published ? `指令已下发：${res.data.command_id}` : `指令已记录（Broker 未连接）：${res.data?.command_id || ''}`)
    // 刷新指令记录

    await loadCmds()
  // 等待异步操作完成：catch (err: any) { notify(err.message) }

  } catch (err: any) { notify(err.message) }
// 当前代码块结束

}
// 初始加载指令记录（与设备列表并行）
// 调用 onMounted(loadCmds)，触发对应功能

onMounted(loadCmds)
</script>
<style scoped>
.stats { display: flex; gap: 12px; align-items: center; }
.stats b { font-size: 16px; }
.muted { color: var(--muted); }
.sec-h { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.sec-h h3 { margin: 0; font-size: 14px; }
.cmds { display: inline-flex; gap: 6px; }
.cmds .btn { padding: 3px 10px; font-size: 12px; }
code { background: #f4f6fa; padding: 1px 6px; border-radius: 4px; font-size: 12px; }
.resp { max-width: 260px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
