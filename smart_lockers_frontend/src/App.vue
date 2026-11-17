<script setup lang="ts">
// @ticket 00-project：根组件——顶部导航 + 路由出口
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()
const showNav = computed<boolean>(() => route.path !== '/')
const navs = [
  { path: '/home', label: '首页' },
  { path: '/deposit', label: '寄存' },
  { path: '/orders', label: '订单' },
  { path: '/pay', label: '充值' },
  { path: '/user', label: '我的' },
  { path: '/notify', label: '通知' },
  { path: '/ops', label: '运维' },
]
</script>

<template>
  <div class="app">
    <header class="topbar">
      <span class="brand">智能快寄柜</span>
      <nav v-if="showNav" class="nav">
        <router-link v-for="n in navs" :key="n.path" :to="n.path" class="nav-link">
          {{ n.label }}
        </router-link>
      </nav>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.app { min-height: 100vh; }
.topbar { display: flex; align-items: center; gap: 24px; padding: 12px 24px; background: #fff; border-bottom: 1px solid #e5e7eb; position: sticky; top: 0; z-index: 10; }
.brand { font-weight: 700; color: #2563eb; }
.nav { display: flex; gap: 16px; flex-wrap: wrap; }
.nav-link { color: #374151; text-decoration: none; font-size: 14px; }
.nav-link.router-link-active { color: #2563eb; font-weight: 600; }
.main { max-width: 960px; margin: 0 auto; padding: 20px 16px; }
</style>
