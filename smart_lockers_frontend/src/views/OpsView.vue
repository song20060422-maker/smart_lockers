<script setup lang="ts">
// @ticket 08-性能调优：pprof/trace/expvar 运维观测入口（需鉴权 + 内网）
// @proto 原型对齐：ops 远程运维(性能采样观测)；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { health, opsPprof, opsTrace, opsExpvar } from '@/api/ops'

const healthState = ref('检测中')
const pprof = ref<unknown>(null)
const trace = ref<unknown>(null)
const expvar = ref<unknown>(null)
const err = ref('')

onMounted(async () => {
  try { const h = await health(); healthState.value = h.status }
  catch (e) { healthState.value = '异常'; err.value = (e as Error).message }
})
async function loadPprof() { try { pprof.value = await opsPprof() } catch (e) { err.value = (e as Error).message } }
async function loadTrace() { try { trace.value = await opsTrace() } catch (e) { err.value = (e as Error).message } }
async function loadExpvar() { try { expvar.value = await opsExpvar() } catch (e) { err.value = (e as Error).message } }
</script>

<template>
  <div class="col">
    <div class="card">
      <h2>服务健康</h2>
      <p>状态：<span class="tag">{{ healthState }}</span></p>
    </div>
    <div class="card">
      <h2>性能采样（需鉴权 / 内网）</h2>
      <div class="row">
        <button class="btn-ghost" @click="loadPprof">pprof</button>
        <button class="btn-ghost" @click="loadTrace">trace</button>
        <button class="btn-ghost" @click="loadExpvar">expvar</button>
      </div>
      <pre v-if="pprof" class="muted">{{ JSON.stringify(pprof, null, 2) }}</pre>
      <pre v-if="trace" class="muted">{{ JSON.stringify(trace, null, 2) }}</pre>
      <pre v-if="expvar" class="muted">{{ JSON.stringify(expvar, null, 2) }}</pre>
    </div>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
