<template>
  <div class="login-wrap">
    <form class="card login-card" @submit.prevent="onLogin">
      <h2>智能快寄柜 · B端运营后台</h2>
      <p class="sub">请使用分配的运营账号登录</p>
      <label>账号</label>
      <input v-model="username" placeholder="operator" autocomplete="username" />
      <label>密码</label>
      <input v-model="password" type="password" placeholder="••••••" autocomplete="current-password" />
      <button class="btn primary" type="submit" :disabled="loading">{{ loading ? '登录中…' : '登录' }}</button>
      <p v-if="error" class="err">{{ error }}</p>
    </form>
    <div v-if="message" class="toast">{{ message }}</div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'

const router = useRouter()
const auth = useAuthStore()
const { message, notify } = useNotify()
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
async function onLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(username.value.trim(), password.value)
    notify('登录成功')
    router.replace('/sites')
  } catch (e: any) {
    error.value = e?.message || '登录失败，请检查账号密码'
  } finally { loading.value = false }
}
</script>
<style scoped>
.login-wrap { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #1f6feb22, #1d243022); }
.login-card { width: 360px; display: flex; flex-direction: column; gap: 12px; }
.login-card h2 { font-size: 18px; }
.sub, .err { color: var(--muted); font-size: 12px; }
.err { color: var(--danger); }
label { font-size: 13px; color: var(--muted); }
</style>
