<template>
  <!-- 页面根容器 div，包裹网点管理页面的全部内容 -->

  <div>
    <!-- 页面头部容器，class 为 page-head，包裹标题与新建按钮 -->

    <div class="page-head">
      <!-- 二级标题，显示“网点管理” -->

      <h2>网点管理</h2>
      <!-- “新建网点”按钮，当认证对象具备 site:write 权限时渲染，class 为 btn primary，点击将 showCreate 置为 true 打开弹窗 -->

      <button v-if="auth.can('site:write')" class="btn primary" @click="showCreate = true">+ 新建网点</button>
    <!-- 闭合页面头部 div 元素 -->

    </div>
    <!-- 卡片容器，class 为 card，包裹网点数据表格 -->

    <div class="card">
      <!-- 表格元素，用于展示网点记录列表 -->

      <table>
        <!-- 表头区域，一行 tr 内依次定义 ID、编码、名称、城市、地址、负责人、柜门、故障、状态、操作等列名 th 单元格 -->

        <thead><tr><th>ID</th><th>编码</th><th>名称</th><th>城市</th><th>地址</th><th>负责人</th><th>柜门</th><th>故障</th><th>状态</th><th>操作</th></tr></thead>
        <!-- 表体区域，用于放置网点数据行 -->

        <tbody>
          <!-- 基于 rows 列表循环生成表格行，key 取 s.id 作为唯一标识 -->

          <tr v-for="s in rows" :key="s.id">
            <!-- 一行内的三个单元格，依次插值显示网点 ID、编码 s.code 与名称 s.name -->

            <td>{{ s.id }}</td><td>{{ s.code }}</td><td>{{ s.name }}</td>
            <!-- 一行内的三个单元格，依次插值显示城市 s.city、地址 s.address 与负责人电话 s.manager_phone -->

            <td>{{ s.city }}</td><td>{{ s.address }}</td><td>{{ s.manager_phone }}</td>
            <!-- 两个单元格，依次插值显示柜门总数 s.door_total 与故障柜门数 s.door_fault -->

            <td>{{ s.door_total }}</td><td>{{ s.door_fault }}</td>
            <!-- 状态单元格，内含标签 span，当 s.status 等于 1 时样式类为 ok、文本“启用”，否则样式类为 err、文本“停用” -->

            <td><span class="tag" :class="s.status === 1 ? 'ok' : 'err'">{{ s.status === 1 ? '启用' : '停用' }}</span></td>
            <!-- 操作列单元格起始，内部放置启停切换按钮或占位文本 -->

            <td>
              <!-- 启停切换按钮，当认证对象具备 site:write 权限时渲染，class 为 btn，点击触发 toggleStatus(s) 切换网点状态（本行与后续行共同构成按钮） -->

              <button v-if="auth.can('site:write')" class="btn" @click="toggleStatus(s)">
                <!-- 按钮文本插值，当 s.status 等于 1 时显示“停用”，否则显示“启用” -->

                {{ s.status === 1 ? '停用' : '启用' }}
              <!-- 闭合启停切换按钮 button 元素 -->

              </button>
              <!-- v-if 的备选分支，无编辑权限时显示占位文本“—”，class 为 muted -->

              <span v-else class="muted">—</span>
            <!-- 闭合操作列单元格 -->

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

    <!-- 新建网点弹窗遮罩容器，v-if 依据 showCreate 是否有值决定是否渲染，class 为 modal -->

    <div v-if="showCreate" class="modal">
      <!-- 新建表单，class 为 card，提交事件被 prevent 修饰阻止默认行为并调用 onCreate 方法 -->

      <form class="card" @submit.prevent="onCreate">
        <!-- 三级标题，显示“新建网点” -->

        <h3>新建网点</h3>
        <!-- 表单标签，包含文本“编码”与必填输入框，v-model 双向绑定 form.code -->

        <label>编码<input v-model="form.code" required /></label>
        <!-- 表单标签，包含文本“名称”与必填输入框，v-model 双向绑定 form.name -->

        <label>名称<input v-model="form.name" required /></label>
        <!-- 表单标签，包含文本“城市”与输入框，v-model 双向绑定 form.city -->

        <label>城市<input v-model="form.city" /></label>
        <!-- 表单标签，包含文本“地址”与输入框，v-model 双向绑定 form.address -->

        <label>地址<input v-model="form.address" /></label>
        <!-- 表单标签，包含文本“负责人电话”与输入框，v-model 双向绑定 form.manager_phone -->

        <label>负责人电话<input v-model="form.manager_phone" /></label>
        <!-- 操作按钮区容器，class 为 actions，包裹提交与取消按钮 -->

        <div class="actions">
          <!-- 提交按钮，class 为 btn primary，type 为 submit，文本“提交” -->

          <button class="btn primary" type="submit">提交</button>
          <!-- 取消按钮，class 为 btn，type 为 button，点击将 showCreate 置为 false 关闭弹窗，文本“取消” -->

          <button class="btn" type="button" @click="showCreate = false">取消</button>
        <!-- 闭合操作按钮区 div 元素 -->

        </div>
      <!-- 闭合新建表单 form 元素 -->

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
import type { Site } from '@/types/api'

// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 赋值更新：const { message, notify } = useNotify()

const { message, notify } = useNotify()
// 声明变量 rows = ref<Site[]>([])，用于保存数据或函数引用

const rows = ref<Site[]>([])
// 声明变量 showCreate = ref(false)，用于保存数据或函数引用

const showCreate = ref(false)
// 声明变量 form = ref({ code: '', name: '', city: '', address: '', manager_phone: '' })，用于保存数据或函数引用

const form = ref({ code: '', name: '', city: '', address: '', manager_phone: '' })
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try { rows.value = (await api.listSites()).data || [] } catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function toggleStatus(s: Site) {

async function toggleStatus(s: Site) {
  // 等待异步操作完成：try { await api.updateSiteStatus(s.id, s.status === 1 ? 0 : 1); notify('操作成功'); await load() }

  try { await api.updateSiteStatus(s.id, s.status === 1 ? 0 : 1); notify('操作成功'); await load() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e: any) { notify(e.message) }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function onCreate() {

async function onCreate() {
  // 等待异步操作完成：try { await api.createSite({ ...form.value }); notify('创建成功'); showCreate.value = false; form.value = { code: '', name: '', city: '', address: '', manager_ph...

  try { await api.createSite({ ...form.value }); notify('创建成功'); showCreate.value = false; form.value = { code: '', name: '', city: '', address: '', manager_phone: '' }; await load() }
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
.modal form { width: 380px; display: flex; flex-direction: column; gap: 12px; }
.modal label { display: flex; flex-direction: column; gap: 5px; font-size: 13px; color: var(--muted); }
.actions { display: flex; gap: 10px; justify-content: flex-end; }
</style>
