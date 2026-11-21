<script setup lang="ts">
// @ticket 00-project：根组件——手机APP外壳（顶部标题栏 + 内容区 + 底部 TabBar）
import { useRoute, useRouter } from 'vue-router'
import { computed } from 'vue'

// 声明变量 route = useRoute()，用于保存数据或函数引用

const route = useRoute()
// 声明变量 router = useRouter()，用于保存数据或函数引用

const router = useRouter()
// 声明变量 title = computed<string>(() => (route.meta.title as string) || '智能快寄柜')，用于保存数据或函数引用

const title = computed<string>(() => (route.meta.title as string) || '智能快寄柜')
// 声明变量 isLogin = computed<boolean>(() => route.path === '/login')，用于保存数据或函数引用

const isLogin = computed<boolean>(() => route.path === '/login')
// 声明变量 canBack = computed<boolean>(() => route.path !== '/home' && route.path !== '/login')，用于保存数据或函数引用

const canBack = computed<boolean>(() => route.path !== '/home' && route.path !== '/login')

// 底部 TabBar 主导航项（首页 / 寄存 / 订单 / 我的）
// 声明变量 tabs = [，用于保存数据或函数引用

const tabs = [
  // 执行 { path: '/home', label: '首页', icon: '⌂' },，实现对应的前端逻辑

  { path: '/home', label: '首页', icon: '⌂' },
  // 执行 { path: '/deposit', label: '寄存', icon: '▣' },，实现对应的前端逻辑

  { path: '/deposit', label: '寄存', icon: '▣' },
  // 执行 { path: '/orders', label: '订单', icon: '▤' },，实现对应的前端逻辑

  { path: '/orders', label: '订单', icon: '▤' },
  // 执行 { path: '/user', label: '我的', icon: '◉' },，实现对应的前端逻辑

  { path: '/user', label: '我的', icon: '◉' },
// 结束切片/数组字面量

]

// 返回上一页
// 定义函数/方法，封装可复用的处理逻辑：function goBack() {

function goBack() {
  // 执行 router.back()，实现对应的前端逻辑

  router.back()
// 当前代码块结束

}
</script>

<template>
  <!-- 手机模拟外壳容器，class 为 phone-shell，限定手机宽度并居中 -->
  <div class="phone-shell">
    <!-- 登录页视图，v-if 依据 isLogin 决定是否渲染，独立渲染不含顶部栏与底部 TabBar -->
    <router-view v-if="isLogin" />
    <!-- 主体布局模板，v-else 在登录页以外渲染，包含顶部栏、内容区与底部导航 -->
    <template v-else>
      <!-- 顶部标题栏元素，class 为 app-header，展示返回按钮与当前页面标题 -->
      <header class="app-header">
        <!-- 返回按钮，v-if 依据 canBack 决定是否渲染，点击触发 goBack -->
        <button v-if="canBack" class="header-back" @click="goBack">‹</button>
        <!-- 标题文本元素，class 为 app-title，插值 title 显示当前页面标题 -->
        <span class="app-title">{{ title }}</span>
      <!-- 闭合顶部标题栏 header 元素 -->

      </header>
      <!-- 内容滚动区，class 为 app-body，承载路由视图 -->
      <main class="app-body">
        <!-- 路由出口，渲染当前页面组件 -->
        <router-view />
      <!-- 闭合内容滚动区 main 元素 -->

      </main>
      <!-- 底部标签导航栏元素，class 为 tabbar，基于 tabs 循环生成主导航项 -->
      <nav class="tabbar">
        <!-- 基于 tabs 数组循环生成标签链接，key 取 t.path，跳转目标为 t.path，样式类为 tab -->
        <router-link v-for="t in tabs" :key="t.path" :to="t.path" class="tab">
          <!-- 标签图标元素，class 为 tab-icon，插值 t.icon 显示图标 -->
          <span class="tab-icon">{{ t.icon }}</span>
          <!-- 标签文字元素，class 为 tab-label，插值 t.label 显示名称 -->
          <span class="tab-label">{{ t.label }}</span>
        <!-- 闭合标签链接 router-link 元素 -->

        </router-link>
      <!-- 闭合底部标签导航栏 nav 元素 -->

      </nav>
    </template>
  <!-- 闭合手机模拟外壳 div 元素 -->

  </div>
</template>

<style scoped>
/* 手机模拟外壳：限定宽度、居中、模拟机身阴影 */
.phone-shell { position: relative; width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; background: #f5f7fa; box-shadow: 0 0 24px rgba(0, 0, 0, .12); overflow: hidden; }
/* 顶部标题栏：粘性定位在壳内顶部 */
.app-header { position: sticky; top: 0; z-index: 20; display: flex; align-items: center; justify-content: center; height: 50px; background: #2563eb; color: #fff; }
/* 返回按钮：绝对定位在标题栏左侧 */
.header-back { position: absolute; left: 8px; background: transparent; color: #fff; font-size: 26px; padding: 0 10px; }
/* 标题文本样式 */
.app-title { font-size: 17px; font-weight: 600; }
/* 内容滚动区：预留底部 TabBar 高度 */
.app-body { min-height: calc(100vh - 100px); padding: 12px 12px 72px; }
/* 底部标签导航栏：固定在壳内底部并水平居中 */
.tabbar { position: fixed; bottom: 0; left: 50%; transform: translateX(-50%); width: 100%; max-width: 430px; z-index: 20; display: flex; background: #fff; border-top: 1px solid #e5e7eb; }
/* 单个标签项：均分宽度、图标文字纵向排列 */
.tab { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 8px 0; color: #6b7280; text-decoration: none; font-size: 12px; }
/* 当前激活标签高亮 */
.tab.router-link-active { color: #2563eb; }
/* 标签图标样式 */
.tab-icon { font-size: 20px; line-height: 1; }
</style>
