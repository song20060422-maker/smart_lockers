<script setup lang="ts">
// @ticket 08-性能调优：pprof/trace/expvar 运维观测入口（需鉴权 + 内网）
// @proto 原型对齐：ops 远程运维(性能采样观测)；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { health, opsPprof, opsTrace, opsExpvar } from '@/api/ops'

// 声明变量 healthState = ref('检测中')，用于保存数据或函数引用

const healthState = ref('检测中')
// 声明变量 pprof = ref<unknown>(null)，用于保存数据或函数引用

const pprof = ref<unknown>(null)
// 声明变量 trace = ref<unknown>(null)，用于保存数据或函数引用

const trace = ref<unknown>(null)
// 声明变量 expvar = ref<unknown>(null)，用于保存数据或函数引用

const expvar = ref<unknown>(null)
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')

// 调用 onMounted(async () => {，触发对应功能

onMounted(async () => {
  // 等待异步操作完成：try { const h = await health(); healthState.value = h.status }

  try { const h = await health(); healthState.value = h.status }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e) { healthState.value = '异常'; err.value = (e as Error).message }
// 执行 })，实现对应的前端逻辑

})
// 定义函数/方法，封装可复用的处理逻辑：async function loadPprof() { try { pprof.value = await opsPprof() } catch (e) { err.value = (e as Error).messa

async function loadPprof() { try { pprof.value = await opsPprof() } catch (e) { err.value = (e as Error).message } }
// 定义函数/方法，封装可复用的处理逻辑：async function loadTrace() { try { trace.value = await opsTrace() } catch (e) { err.value = (e as Error).messa

async function loadTrace() { try { trace.value = await opsTrace() } catch (e) { err.value = (e as Error).message } }
// 定义函数/方法，封装可复用的处理逻辑：async function loadExpvar() { try { expvar.value = await opsExpvar() } catch (e) { err.value = (e as Error).me

async function loadExpvar() { try { expvar.value = await opsExpvar() } catch (e) { err.value = (e as Error).message } }
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列运维相关卡片 -->

  <div class="col">
    <!-- 卡片容器，class 为 card，展示服务健康状态 -->

    <div class="card">
      <!-- 二级标题，显示“服务健康” -->

      <h2>服务健康</h2>
      <!-- 段落元素，显示“状态：”并以内联标签 span（class 为 tag）插值 healthState 显示健康状态值 -->

      <p>状态：<span class="tag">{{ healthState }}</span></p>
    <!-- 闭合服务健康卡片 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，展示性能采样数据（需鉴权 / 内网） -->

    <div class="card">
      <!-- 二级标题，显示“性能采样（需鉴权 / 内网）” -->

      <h2>性能采样（需鉴权 / 内网）</h2>
      <!-- 行容器，class 为 row，横向排列采样类型按钮 -->

      <div class="row">
        <!-- pprof 采样按钮，class 为 btn-ghost，点击触发 loadPprof -->

        <button class="btn-ghost" @click="loadPprof">pprof</button>
        <!-- trace 采样按钮，class 为 btn-ghost，点击触发 loadTrace -->

        <button class="btn-ghost" @click="loadTrace">trace</button>
        <!-- expvar 采样按钮，class 为 btn-ghost，点击触发 loadExpvar -->

        <button class="btn-ghost" @click="loadExpvar">expvar</button>
      <!-- 闭合采样按钮行容器 div 元素 -->

      </div>
      <!-- 预格式化元素，当 pprof 有值时渲染，class 为 muted，通过 JSON.stringify(pprof, null, 2) 格式化显示采样数据 -->

      <pre v-if="pprof" class="muted">{{ JSON.stringify(pprof, null, 2) }}</pre>
      <!-- 预格式化元素，当 trace 有值时渲染，class 为 muted，通过 JSON.stringify(trace, null, 2) 格式化显示采样数据 -->

      <pre v-if="trace" class="muted">{{ JSON.stringify(trace, null, 2) }}</pre>
      <!-- 预格式化元素，当 expvar 有值时渲染，class 为 muted，通过 JSON.stringify(expvar, null, 2) 格式化显示采样数据 -->

      <pre v-if="expvar" class="muted">{{ JSON.stringify(expvar, null, 2) }}</pre>
    <!-- 闭合性能采样卡片 div 元素 -->

    </div>
    <!-- 错误提示段落，当 err 有值时渲染，class 为 err，插值 err 显示错误信息 -->

    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
