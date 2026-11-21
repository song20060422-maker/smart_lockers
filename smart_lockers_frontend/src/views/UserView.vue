<script setup lang="ts">
// @ticket 06-个人中心：登录/头像/我的服务/卡券/退出
// @proto 原型对齐：v-marketing 营销奖促中心(卡券)；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { myServices, myCoupons, updateAvatar, orderSummary } from '@/api/user'
import type { MyService, Coupon } from '@/types'

// 声明变量 router = useRouter()，用于保存数据或函数引用

const router = useRouter()
// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 声明变量 services = ref<MyService[]>([])，用于保存数据或函数引用

const services = ref<MyService[]>([])
// 声明变量 summary = ref<Record<string, unknown> | null>(null)，用于保存数据或函数引用

const summary = ref<Record<string, unknown> | null>(null)
// 声明变量 coupons = ref<Coupon[]>([])，用于保存数据或函数引用

const coupons = ref<Coupon[]>([])
// 声明变量 loginId = ref<number>(10001)，用于保存数据或函数引用

const loginId = ref<number>(10001)
// 声明变量 avatarUrl = ref('')，用于保存数据或函数引用

const avatarUrl = ref('')
// 声明变量 msg = ref('')，用于保存数据或函数引用

const msg = ref('')
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')

// 调用 onMounted(() => {，触发对应功能

onMounted(() => {
  // 条件判断：if (auth.logged) loadData() 时执行下面逻辑

  if (auth.logged) loadData()
// 执行 })，实现对应的前端逻辑

})

// 定义函数/方法，封装可复用的处理逻辑：async function loadData() {

async function loadData() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 等待异步操作完成：const [s, c, sum] = await Promise.all([myServices(), myCoupons(), orderSummary()])

    const [s, c, sum] = await Promise.all([myServices(), myCoupons(), orderSummary()])
    // 赋值更新：services.value = s.list

    services.value = s.list
    // 赋值更新：coupons.value = c

    coupons.value = c
    // 赋值更新：summary.value = sum

    summary.value = sum
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function doLogin() {

async function doLogin() {
  // 等待异步操作完成：try { await auth.login(loginId.value); msg.value = '登录成功'; await loadData() }

  try { await auth.login(loginId.value); msg.value = '登录成功'; await loadData() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function doLogout() {

async function doLogout() {
  // 条件判断：if (!window.confirm('确认退出登录？')) return 时执行下面逻辑

  if (!window.confirm('确认退出登录？')) return
  // 等待异步操作完成：await auth.logout()

  await auth.logout()
  // 赋值更新：services.value = []; coupons.value = []

  services.value = []; coupons.value = []
  // 赋值更新：msg.value = '已退出'

  msg.value = '已退出'
  // 跳转登录页，形成“退出→重新登录”闭环

  router.push('/login')
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function doAvatar() {

async function doAvatar() {
  // 条件判断：if (!avatarUrl.value) { err.value = '请输入头像 URL'; return } 时执行下面逻辑

  if (!avatarUrl.value) { err.value = '请输入头像 URL'; return }
  // 等待异步操作完成：try { await updateAvatar({ url: avatarUrl.value }); msg.value = '头像已更新' }

  try { await updateAvatar({ url: avatarUrl.value }); msg.value = '头像已更新' }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function fen(v: number): string { return (v / 100).toFixed(2) }

function fen(v: number): string { return (v / 100).toFixed(2) }
// 定义函数/方法，封装可复用的处理逻辑：function mask(code: string): string { return code.length > 4 ? code.slice(0, 2) + '****' + code.slice(-2) : co

function mask(code: string): string { return code.length > 4 ? code.slice(0, 2) + '****' + code.slice(-2) : code }
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列个人中心与服务卡片 -->

  <div class="col">
    <!-- 卡片容器，class 为 card，展示登录状态与个人操作 -->

    <div class="card">
      <!-- 二级标题，显示“个人中心” -->

      <h2>个人中心</h2>
      <!-- 已登录分支容器，v-if 依据 auth.logged 为真时渲染，包裹登录用户信息与操作 -->

      <template v-if="auth.logged">
        <!-- 段落元素，显示“已登录（用户”并插值 auth.userId 后闭合括号 -->

        <p>已登录（用户 {{ auth.userId }}）</p>
        <!-- 行容器，class 为 row，横向排列头像输入框与操作按钮 -->

        <div class="row">
          <!-- 头像 URL 输入框，v-model 双向绑定 avatarUrl，占位提示“头像 URL” -->

          <input v-model="avatarUrl" placeholder="头像 URL" />
          <!-- 更新头像按钮，点击触发 doAvatar 更新用户头像 -->

          <button @click="doAvatar">更新头像</button>
          <!-- 退出登录按钮，class 为 btn-ghost，点击触发 doLogout 退出当前账号 -->

          <button class="btn-ghost" @click="doLogout">退出登录</button>
        <!-- 闭合已登录操作行容器 div 元素 -->

        </div>
      <!-- 闭合已登录分支 template 元素 -->

      </template>
      <!-- 未登录分支容器，v-else 在 auth.logged 为假时渲染，包裹登录表单 -->

      <template v-else>
        <!-- 行容器，class 为 row，横向排列用户 ID 输入框与登录按钮 -->

        <div class="row">
          <!-- 用户 ID 输入框，type 为 number，v-model.number 双向绑定 loginId，占位提示“用户 ID” -->

          <input v-model.number="loginId" type="number" placeholder="用户 ID" />
          <!-- 登录按钮，点击触发 doLogin 执行登录 -->

          <button @click="doLogin">登录</button>
        <!-- 闭合未登录操作行容器 div 元素 -->

        </div>
      <!-- 闭合未登录分支 template 元素 -->

      </template>
      <!-- 提示段落，当 msg 有值时渲染，class 为 muted，插值 msg 显示提示内容 -->

      <p v-if="msg" class="muted">{{ msg }}</p>
    <!-- 闭合个人中心卡片 div 元素 -->

    </div>
    <!-- 订单汇总卡片，class 为 card，当 auth.logged 为真且 summary 有值时渲染 -->

    <div v-if="auth.logged && summary" class="card">
      <!-- 二级标题，显示“订单汇总” -->

      <h2>订单汇总</h2>
      <!-- 段落元素，class 为 muted，通过 JSON.stringify(summary) 显示订单汇总数据 -->

      <p class="muted">{{ JSON.stringify(summary) }}</p>
    <!-- 闭合订单汇总卡片 div 元素 -->

    </div>
    <!-- 我的服务卡片，class 为 card，当 auth.logged 为真时渲染 -->

    <div v-if="auth.logged" class="card">
      <!-- 二级标题，显示“我的服务” -->

      <h2>我的服务</h2>
      <!-- 行容器，class 为 row，横向排列服务标签 -->

      <div class="row">
        <!-- 基于 services 数组循环生成服务标签，key 取 s.id，样式类为 tag，插值 s.title 显示服务名称 -->

        <span v-for="s in services" :key="s.id" class="tag">{{ s.title }}</span>
      <!-- 闭合服务标签行容器 div 元素 -->

      </div>
      <!-- 空状态提示段落，当 services 为空时渲染，class 为 muted，显示“暂无服务” -->

      <p v-if="!services.length" class="muted">暂无服务</p>
    <!-- 闭合我的服务卡片 div 元素 -->

    </div>
    <!-- 我的卡券卡片，class 为 card，当 auth.logged 为真时渲染 -->

    <div v-if="auth.logged" class="card">
      <!-- 二级标题，显示“我的卡券” -->

      <h2>我的卡券</h2>
      <!-- 基于 coupons 数组循环生成卡券条目，key 取 c.id，样式类为 list-item -->

      <div v-for="c in coupons" :key="c.id" class="list-item">
        <!-- 行内文本元素，显示经 mask 遮罩的卡券编码 c.code、经 fen 格式化的优惠金额 c.discount，以及 c.used 为真时显示“已用”否则“可用” -->

        <span>{{ mask(c.code) }} / 优惠 {{ fen(c.discount) }}元 / {{ c.used ? '已用' : '可用' }}</span>
      <!-- 闭合当前卡券条目 div 元素 -->

      </div>
      <!-- 空状态提示段落，当 coupons 为空时渲染，class 为 muted，显示“暂无卡券” -->

      <p v-if="!coupons.length" class="muted">暂无卡券</p>
    <!-- 闭合我的卡券卡片 div 元素 -->

    </div>
    <!-- ===== 补齐入口：更多服务跳转（只新增，不改动上方原有卡片） ===== -->
    <!-- 更多服务卡片容器，class 为 card，集中放置原型图/工单对应页面的跳转入口 -->
    <div class="card">
      <!-- 二级标题，显示“更多服务” -->
      <h2>更多服务</h2>
      <!-- 行容器，class 为 row，横向排列各功能入口链接 -->
      <div class="row">
        <!-- 链接：我的优惠券，跳转 /coupons -->
        <router-link to="/coupons" class="tag" style="text-decoration:none">我的优惠券</router-link>
        <!-- 链接：在线客服，跳转 /service -->
        <router-link to="/service" class="tag" style="text-decoration:none">在线客服</router-link>
        <!-- 链接：人脸识别，跳转 /face -->
        <router-link to="/face" class="tag" style="text-decoration:none">人脸识别</router-link>
        <!-- 链接：寄存指南，跳转 /guide -->
        <router-link to="/guide" class="tag" style="text-decoration:none">寄存指南</router-link>
        <!-- 链接：消息通知，跳转 /notify -->
        <router-link to="/notify" class="tag" style="text-decoration:none">消息通知</router-link>
      <!-- 闭合入口链接行容器 div 元素 -->
      </div>
    <!-- 闭合更多服务卡片 div 元素 -->
    </div>

    <!-- 错误提示段落，当 err 有值时渲染，class 为 err，插值 err 显示错误信息 -->

    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
