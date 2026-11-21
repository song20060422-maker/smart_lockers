<template>
  <!-- 后台布局根容器，class 为 layout，包裹侧边栏与主内容区 -->

  <div class="layout">
    <!-- 侧边栏容器，class 为 sidebar，用于放置品牌标识与导航菜单 -->

    <aside class="sidebar">
      <!-- 品牌标题文本块，显示“智能快寄柜 · B端”字样 -->

      <div class="brand">智能快寄柜 · B端</div>
      <!-- 导航区域元素，包裹顶部导航链接列表 -->

      <nav>
        <!-- 基于 menus 菜单数组循环生成路由链接，key 取 m.path，跳转目标为 m.path，样式类为 nav-item（本行与下一行同属一个标签，属性换行书写） -->

        <router-link v-for="m in menus" :key="m.path" :to="m.path" class="nav-item"
          :class="{ active: $route.path === m.path }">{{ m.title }}</router-link>
      <!-- 闭合导航区域 nav 元素 -->

      </nav>
    <!-- 闭合侧边栏 aside 元素 -->

    </aside>
    <!-- 主内容区容器，class 为 main，包裹顶栏与页面主体 -->

    <section class="main">
      <!-- 顶部栏元素，class 为 topbar，展示当前标题与用户信息 -->

      <header class="topbar">
        <!-- 行内文本元素，插值 currentTitle 显示当前页面标题 -->

        <span>{{ currentTitle }}</span>
        <!-- 用户信息区容器，class 为 user，包裹用户名、角色与退出按钮 -->

        <div class="user">
          <!-- 行内文本元素，插值显示认证用户对象的 username，若为空则显示“未登录” -->

          <span>{{ auth.user?.username || '未登录' }}</span>
          <!-- 行内文本元素，class 为 role，插值 roleLabel 显示当前用户角色名称 -->

          <span class="role">{{ roleLabel }}</span>
          <!-- 退出按钮，class 为 btn，点击事件绑定 onLogout 方法执行登出，按钮文本为“退出” -->

          <button class="btn" @click="onLogout">退出</button>
        <!-- 闭合用户信息区 div 元素 -->

        </div>
      <!-- 闭合顶部栏 header 元素 -->

      </header>
      <!-- 页面主体元素，class 为 content，内部放置路由视图出口，用于渲染匹配的子路由组件 -->

      <main class="content"><router-view /></main>
    <!-- 闭合主内容区 section 元素 -->

    </section>
    <!-- 全局轻提示条，v-if 依据 message 是否有值决定是否渲染，class 为 toast，插值 message 显示提示内容 -->

    <div v-if="message" class="toast">{{ message }}</div>
  <!-- 闭合布局根容器 div 元素 -->

  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'

// 声明变量 route = useRoute()，用于保存数据或函数引用

const route = useRoute()
// 声明变量 router = useRouter()，用于保存数据或函数引用

const router = useRouter()
// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 赋值更新：const { message } = useNotify()

const { message } = useNotify()
// 声明变量 menus = [，用于保存数据或函数引用

const menus = [
  // 执行 { path: '/sites', title: '网点管理' },，实现对应的前端逻辑

  { path: '/sites', title: '网点管理' },
  // 执行 { path: '/devices', title: '设备监控' },，实现对应的前端逻辑

  { path: '/devices', title: '设备监控' },
  // 执行 { path: '/reconcile', title: '充值对账' },，实现对应的前端逻辑

  { path: '/reconcile', title: '充值对账' },
  // 执行 { path: '/alarms', title: '告警处理' },，实现对应的前端逻辑

  { path: '/alarms', title: '告警处理' },
// 结束切片/数组字面量

]
// 声明变量 currentTitle = computed(() => (route.meta.title as string) || '运营后台')，用于保存数据或函数引用

const currentTitle = computed(() => (route.meta.title as string) || '运营后台')
// 声明变量 roleLabel = computed(() => ({ user: '普通用户', site_admin: '网点管理员', operator: '平台运营' }[auth.role]))，用于保存数据或函数引用

const roleLabel = computed(() => ({ user: '普通用户', site_admin: '网点管理员', operator: '平台运营' }[auth.role]))
// 定义函数/方法，封装可复用的处理逻辑：function onLogout() { auth.logout(); router.replace('/login') }

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
