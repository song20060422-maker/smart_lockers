<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="brand">智能快寄柜 · B端</div>
      <nav>
        <router-link v-for="m in menus" :key="m.path" :to="m.path" class="nav-item"
          :class="{ active: $route.path === m.path }">{{ m.title }}</router-link>
      </nav>
    </aside>
    <section class="main">
      <header class="topbar">
        <span>{{ currentTitle }}</span>
        <div class="user">
          <span>{{ auth.user?.username || '未登录' }}</span>
          <span class="role">{{ roleLabel }}</span>
          <button class="btn" @click="onLogout">退出</button>
        </div>
      </header>
      <main class="content"><router-view /></main>
    </section>
    <div v-if="message" class="toast">{{ message }}</div>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const { message } = useNotify()
const menus = [
  { path: '/sites', title: '网点管理' },
  { path: '/devices', title: '设备监控' },
  { path: '/reconcile', title: '充值对账' },
  { path: '/alarms', title: '告警处理' },
]
const currentTitle = computed(() => (route.meta.title as string) || '运营后台')
const roleLabel = computed(() => ({ user: '普通用户', site_admin: '网点管理员', operator: '平台运营' }[auth.role]))
function onLogout() { auth.logout(); router.replace('/login') }
</script>
<style scoped>
.layout { display: flex; min-height: 100vh; }
.sidebar { width: 200px; background: #1d2430; color: #c9d1dc; flex-shrink: 0; }
.brand { padding: 20px 18px; font-size: 16px; font-weight: 600; color: #fff; border-bottom: 1px solid #2c3543; }
.nav-item { display: block; padding: 13px 18px; }
.nav-item:hover { background: #263041; }
.nav-item.active { background: var(--primary); color: #fff; }
.main { flex: 1; min-width: 0; }
.topbar { height: 56px; background: #fff; border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 22px; }
.user { display: flex; align-items: center; gap: 12px; }
.role { color: var(--primary); font-size: 12px; background: #eaf2ff; padding: 2px 8px; border-radius: 10px; }
.content { padding: 22px; }
</style>
