<script setup lang="ts">
// @ticket 06-个人中心：登录/头像/我的服务/卡券/退出
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { myServices, myCoupons, updateAvatar, orderSummary } from '@/api/user'
import type { MyService, Coupon } from '@/types'

const router = useRouter()
const auth = useAuthStore()
const services = ref<MyService[]>([])
const summary = ref<Record<string, unknown> | null>(null)
const coupons = ref<Coupon[]>([])
const loginId = ref<number>(10001)
const avatarUrl = ref('')
const msg = ref('')
const err = ref('')

onMounted(() => {
  if (auth.logged) loadData()
})

async function loadData() {
  try {
    const [s, c, sum] = await Promise.all([myServices(), myCoupons(), orderSummary()])
    services.value = s.list
    coupons.value = c
    summary.value = sum
  } catch (e) { err.value = (e as Error).message }
}
async function doLogin() {
  try { await auth.login(loginId.value); msg.value = '登录成功'; await loadData() }
  catch (e) { err.value = (e as Error).message }
}
async function doLogout() {
  if (!window.confirm('确认退出登录？')) return
  await auth.logout()
  services.value = []; coupons.value = []
  msg.value = '已退出'
}
async function doAvatar() {
  if (!avatarUrl.value) { err.value = '请输入头像 URL'; return }
  try { await updateAvatar({ url: avatarUrl.value }); msg.value = '头像已更新' }
  catch (e) { err.value = (e as Error).message }
}
function fen(v: number): string { return (v / 100).toFixed(2) }
function mask(code: string): string { return code.length > 4 ? code.slice(0, 2) + '****' + code.slice(-2) : code }
</script>

<template>
  <div class="col">
    <div class="card">
      <h2>个人中心</h2>
      <template v-if="auth.logged">
        <p>已登录（用户 {{ auth.userId }}）</p>
        <div class="row">
          <input v-model="avatarUrl" placeholder="头像 URL" />
          <button @click="doAvatar">更新头像</button>
          <button class="btn-ghost" @click="doLogout">退出登录</button>
        </div>
      </template>
      <template v-else>
        <div class="row">
          <input v-model.number="loginId" type="number" placeholder="用户 ID" />
          <button @click="doLogin">登录</button>
        </div>
      </template>
      <p v-if="msg" class="muted">{{ msg }}</p>
    </div>
    <div v-if="auth.logged && summary" class="card">
      <h2>订单汇总</h2>
      <p class="muted">{{ JSON.stringify(summary) }}</p>
    </div>
    <div v-if="auth.logged" class="card">
      <h2>我的服务</h2>
      <div class="row">
        <span v-for="s in services" :key="s.id" class="tag">{{ s.title }}</span>
      </div>
      <p v-if="!services.length" class="muted">暂无服务</p>
    </div>
    <div v-if="auth.logged" class="card">
      <h2>我的卡券</h2>
      <div v-for="c in coupons" :key="c.id" class="list-item">
        <span>{{ mask(c.code) }} / 优惠 {{ fen(c.discount) }}元 / {{ c.used ? '已用' : '可用' }}</span>
      </div>
      <p v-if="!coupons.length" class="muted">暂无卡券</p>
    </div>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
