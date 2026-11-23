<template>
  <!-- OTA 固件升级页根容器（对齐原型 ota 视图"OTA 固件升级"）@ticket B端-IoT中台 -->

  <div>
    <!-- 页面头部：标题 + 新建任务按钮（ota:manage 权限控制显隐） -->

    <div class="page-head"><h2>OTA 升级</h2>
      <!-- 新建升级任务按钮，点击展开表单 -->

      <button v-if="auth.can('ota:manage')" class="btn primary" @click="showForm = !showForm">{{ showForm ? '收起' : '新建升级任务' }}</button>
    <!-- 闭合页头 -->

    </div>
    <!-- 新建任务表单卡片（v-if 展示） -->

    <div v-if="showForm" class="card form">
      <!-- 目标设备下拉框 -->

      <label>目标设备 <select v-model.number="form.device_id"><option v-for="d in devices" :key="d.id" :value="d.id">{{ d.code }}（{{ d.device_sn || '未绑定SN' }}）</option></select></label>
      <!-- 固件版本输入框 -->

      <label>固件版本 <input v-model="form.firmware_version" placeholder="v1.3.0" /></label>
      <!-- 下载地址输入框 -->

      <label>下载地址 <input v-model="form.download_url" placeholder="https://ota.example.com/firmware/v1.3.0.bin" /></label>
      <!-- MD5 校验输入框 -->

      <label>MD5 <input v-model="form.md5" placeholder="固件MD5（可选）" /></label>
      <!-- 文件大小输入框 -->

      <label>大小(字节) <input v-model.number="form.file_size" type="number" placeholder="0" /></label>
      <!-- 强制升级开关 -->

      <label class="chk"><input type="checkbox" v-model="form.force_upgrade" /> 强制升级</label>
      <!-- 提交按钮 -->

      <button class="btn primary" @click="onCreate">下发升级通知</button>
    <!-- 闭合表单卡片 -->

    </div>
    <!-- OTA 任务列表卡片 -->

    <div class="card">
      <!-- 任务表格 -->

      <table>
        <!-- 表头：任务号/设备/目标版本/状态/进度/说明/操作人/创建时间 -->

        <thead><tr><th>任务号</th><th>设备</th><th>目标版本</th><th>状态</th><th>进度</th><th>说明</th><th>操作人</th><th>创建时间</th></tr></thead>
        <!-- 表体：循环任务列表 -->

        <tbody>
          <tr v-for="t in rows" :key="t.id">
            <!-- 任务号（等宽展示） -->

            <td><code>{{ t.ota_id }}</code></td>
            <!-- 设备 SN -->

            <td>{{ t.device_sn || t.device_id }}</td>
            <!-- 目标固件版本 -->

            <td>{{ t.firmware_version }}</td>
            <!-- 状态标签 -->

            <td><span class="tag" :class="statusClass(t.status)">{{ statusText(t.status) }}</span></td>
            <!-- 进度条单元格 -->

            <td>
              <!-- 进度条容器 -->

              <div class="bar"><div class="fill" :style="{ width: t.progress + '%' }"></div></div>
              <!-- 进度文本 -->

              <span class="muted">{{ t.progress }}%</span>
            <!-- 闭合进度单元格 -->

            </td>
            <!-- 设备上报说明 -->

            <td class="muted">{{ t.message || '—' }}</td>
            <!-- 操作人 -->

            <td>{{ t.operator }}</td>
            <!-- 创建时间 -->

            <td class="muted">{{ t.created_at }}</td>
          <!-- 闭合行 -->

          </tr>
          <!-- 空数据占位 -->

          <tr v-if="!rows.length"><td colspan="8" class="muted">暂无任务</td></tr>
        <!-- 闭合表体 -->

        </tbody>
      <!-- 闭合表格 -->

      </table>
    <!-- 闭合卡片 -->

    </div>
    <!-- 全局轻提示条 -->

    <div v-if="message" class="toast">{{ message }}</div>
  <!-- 闭合根容器 -->

  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'
import type { Device, OtaTask } from '@/types/api'

// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 赋值更新：const { message, notify } = useNotify()

const { message, notify } = useNotify()
// 声明变量 rows = ref<OtaTask[]>([])，OTA 任务列表

const rows = ref<OtaTask[]>([])
// 声明变量 devices = ref<Device[]>([])，设备下拉数据

const devices = ref<Device[]>([])
// 声明变量 showForm = ref(false)，表单显隐

const showForm = ref(false)
// 声明变量 form = reactive({...})，新建任务表单模型

const form = reactive({ device_id: 0, firmware_version: '', download_url: '', md5: '', file_size: 0, force_upgrade: false })
// 状态文案映射：0已下发 1下载中 2升级中 3成功 4失败

const statusText = (s: number) => ({ 0: '已下发', 1: '下载中', 2: '升级中', 3: '成功', 4: '失败' }[s] || '未知')
// 状态样式映射

const statusClass = (s: number) => ({ 0: 'warn', 1: 'warn', 2: 'warn', 3: 'ok', 4: 'err' }[s] || 'muted')
// 定义函数/方法 load：拉取设备与任务列表

async function load() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 赋值更新：devices.value = (await api.listDevices()).data || []

    devices.value = (await api.listDevices()).data || []
    // 默认选中第一台设备

    if (!form.device_id && devices.value.length) form.device_id = devices.value[0].id
    // 赋值更新：rows.value = (await api.iotOtaTasks()).data || []

    rows.value = (await api.iotOtaTasks()).data || []
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 定义函数/方法 onCreate：提交新建 OTA 任务

async function onCreate() {
  // 参数校验：设备与版本、地址必填

  if (!form.device_id || !form.firmware_version || !form.download_url) { notify('请填写设备、版本与下载地址'); return }
  // 等待异步操作完成：try { await api.iotOtaCreate(form.device_id, { ...form }) }

  try {
    // 调用创建接口

    const res: any = await api.iotOtaCreate(form.device_id, { ...form })
    // 成功提示：published=false 时提示演示模式（Broker 未连接）

    notify(res.data?.published ? `任务已下发：${res.data.ota_id}` : `任务已记录（Broker 未连接）：${res.data?.ota_id || ''}`)
    // 收起表单并刷新列表

    showForm.value = false
    // 刷新任务列表

    await load()
  // 等待异步操作完成：catch (e: any) { notify(e.message) }

  } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 调用 onMounted(load)，触发对应功能

onMounted(load)
</script>
<style scoped>
.form { display: flex; flex-wrap: wrap; gap: 12px; align-items: center; margin-bottom: 16px; }
.form label { display: flex; align-items: center; gap: 6px; font-size: 13px; }
.form input, .form select { padding: 6px 10px; border: 1px solid var(--border); border-radius: 6px; }
.form .chk input { margin-right: 4px; }
.bar { display: inline-block; width: 90px; height: 6px; background: #eef1f6; border-radius: 3px; margin-right: 6px; vertical-align: middle; }
.fill { height: 100%; background: var(--primary); border-radius: 3px; }
code { background: #f4f6fa; padding: 1px 6px; border-radius: 4px; font-size: 12px; }
.muted { color: var(--muted); }
</style>
