<template>
  <!-- IoT 中台总览页根容器 @ticket B端-IoT中台 -->

  <div>
    <!-- 页面头部：标题与 EMQX 接入层健康状态 -->

    <div class="page-head"><h2>IoT 总览</h2>
      <!-- 设备状态分布标签区（对齐原型 overview"设备状态分布"，数据来自 device/status Topic） -->

      <div class="stats">
        <!-- 总设备数 -->

        <span>设备 <b>{{ ov.devices.total }}</b></span>
        <!-- 在线数 -->

        <span class="tag ok">在线 {{ ov.devices.online }}</span>
        <!-- 离线数 -->

        <span class="tag warn">离线 {{ ov.devices.offline }}</span>
        <!-- 故障数 -->

        <span class="tag err">故障 {{ ov.devices.fault }}</span>
        <!-- EMQX 健康状态（对齐原型"EMQX 集群健康"） -->

        <span class="tag" :class="ov.emqx.status === 'connected' ? 'ok' : 'warn'">EMQX {{ emqxText }}</span>
      <!-- 闭合统计区 -->

      </div>
    <!-- 闭合页头 -->

    </div>
    <!-- 指令/OTA 汇总卡片行 -->

    <div class="cards">
      <!-- 指令总数卡片 -->

      <div class="card mini"><div class="k">累计下发指令</div><div class="v">{{ ov.commands_total }}</div></div>
      <!-- OTA 任务总数卡片 -->

      <div class="card mini"><div class="k">OTA 任务</div><div class="v">{{ ov.ota_total }}</div></div>
      <!-- Broker 地址卡片 -->

      <div class="card mini"><div class="k">Broker</div><div class="v small">{{ ov.emqx.broker || '未配置' }}</div></div>
    <!-- 闭合卡片行 -->

    </div>
    <!-- Topic 设计规范表（对齐原型 mqtt 视图"Topic 设计规范"，格式 {业务域}/{消息类型}/{设备标识}/{具体功能}） -->

    <div class="card">
      <h3>Topic 设计规范</h3>
      <!-- 规范条目表格 -->

      <table>
        <!-- 表头：Topic / 方向 / QoS / 说明 -->

        <thead><tr><th>Topic</th><th>方向</th><th>QoS</th><th>说明</th></tr></thead>
        <!-- 表体：循环 topic_spec -->

        <tbody>
          <tr v-for="(t, i) in ov.topic_spec" :key="i">
            <!-- Topic 模板（等宽字体展示） -->

            <td><code>{{ t.topic }}</code></td>
            <!-- 方向：上行=设备→云端，下行=云端→设备 -->

            <td>{{ t.direction }}</td>
            <!-- QoS 等级 -->

            <td>{{ t.qos }}</td>
            <!-- 说明文案 -->

            <td>{{ t.desc }}</td>
          <!-- 闭合行 -->

          </tr>
        <!-- 闭合表体 -->

        </tbody>
      <!-- 闭合表格 -->

      </table>
    <!-- 闭合卡片 -->

    </div>
    <!-- ACL 权限规则表（对齐原型 mqtt 视图"ACL 权限规则"） -->

    <div class="card">
      <h3>ACL 权限规则</h3>
      <!-- ACL 条目表格 -->

      <table>
        <!-- 表头：角色 / 权限 / 规则 / 说明 -->

        <thead><tr><th>角色</th><th>权限</th><th>规则</th><th>说明</th></tr></thead>
        <!-- 表体：循环 acl_spec -->

        <tbody>
          <tr v-for="(a, i) in ov.acl_spec" :key="i">
            <!-- 角色名 -->

            <td>{{ a.role }}</td>
            <!-- 权限动作 -->

            <td><span class="tag" :class="a.perm === 'deny' ? 'err' : 'ok'">{{ a.perm }}</span></td>
            <!-- ACL 规则（等宽展示） -->

            <td><code>{{ a.rule }}</code></td>
            <!-- 说明文案 -->

            <td>{{ a.desc }}</td>
          <!-- 闭合行 -->

          </tr>
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
import { ref, computed, onMounted } from 'vue'
import { api } from '@/api'
import { useNotify } from '@/composables/useNotify'
import type { IotOverview } from '@/types/api'

// 赋值更新：const { message, notify } = useNotify()

const { message, notify } = useNotify()
// 声明变量 ov = ref<IotOverview>，初始为空态占位数据

const ov = ref<IotOverview>({
  // 设备统计占位

  devices: { total: 0, online: 0, offline: 0, fault: 0 },
  // EMQX 占位

  emqx: { status: 'standby', broker: '' },
  // 指令/OTA 占位

  commands_total: 0, ota_total: 0,
  // 规范表占位

  topic_spec: [], acl_spec: [],
// 当前代码块结束

})
// 计算属性 emqxText：connected→已连接，standby→演示模式

const emqxText = computed(() => ov.value.emqx.status === 'connected' ? '已连接' : '演示模式')
// 定义函数/方法 load：拉取总览数据

async function load() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 赋值更新：ov.value = (await api.iotOverview()).data || ov.value

    ov.value = (await api.iotOverview()).data || ov.value
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 调用 onMounted(load)，触发对应功能

onMounted(load)
</script>
<style scoped>
.cards { display: flex; gap: 14px; margin-bottom: 16px; }
.card.mini { flex: 1; padding: 14px 18px; }
.card .k { color: var(--muted); font-size: 12px; }
.card .v { font-size: 22px; font-weight: 600; margin-top: 4px; }
.card .v.small { font-size: 13px; font-weight: 400; word-break: break-all; }
code { background: #f4f6fa; padding: 1px 6px; border-radius: 4px; font-size: 12px; }
h3 { margin: 0 0 10px; font-size: 14px; }
.card { margin-bottom: 16px; }
</style>
