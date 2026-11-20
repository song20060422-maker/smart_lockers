// @ticket B端-前端：路由与登录守卫（三类角色共用页面，按钮级由权限控制）
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { getToken } from '@/api/request'
import Layout from '@/layouts/MainLayout.vue'

// 声明变量 routes: RouteRecordRaw[] = [，用于保存数据或函数引用

const routes: RouteRecordRaw[] = [
  // 执行 { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue') },，实现对应的前端逻辑

  { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue') },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/'

    path: '/',
    // 定义对象属性 component，取值为 Layout

    component: Layout,
    // 定义对象属性 redirect，取值为 '/sites'

    redirect: '/sites',
    // 定义对象属性 children，取值为 [

    children: [
      // 执行 { path: 'sites', name: 'sites', component: () => import('@/views/SitesView.vue'), meta: { title: '网点管理', perm: 'site:read' } },，实现对应的前端逻辑

      { path: 'sites', name: 'sites', component: () => import('@/views/SitesView.vue'), meta: { title: '网点管理', perm: 'site:read' } },
      // 执行 { path: 'devices', name: 'devices', component: () => import('@/views/DevicesView.vue'), meta: { title: '设备监控', perm: 'device:read' } },，实现对应的前端逻辑

      { path: 'devices', name: 'devices', component: () => import('@/views/DevicesView.vue'), meta: { title: '设备监控', perm: 'device:read' } },
      // 执行 { path: 'reconcile', name: 'reconcile', component: () => import('@/views/ReconcileView.vue'), meta: { title: '充值对账', perm: 'reconcile:read' } },，实现对应的前端逻辑

      { path: 'reconcile', name: 'reconcile', component: () => import('@/views/ReconcileView.vue'), meta: { title: '充值对账', perm: 'reconcile:read' } },
      // 执行 { path: 'alarms', name: 'alarms', component: () => import('@/views/AlarmsView.vue'), meta: { title: '告警处理', perm: 'alarm:read' } },，实现对应的前端逻辑

      { path: 'alarms', name: 'alarms', component: () => import('@/views/AlarmsView.vue'), meta: { title: '告警处理', perm: 'alarm:read' } },
    // 执行 ],，实现对应的前端逻辑

    ],
  // 定义结束

  },
// 结束切片/数组字面量

]

// 声明变量 router = createRouter({ history: createWebHistory(), routes })，用于保存数据或函数引用

const router = createRouter({ history: createWebHistory(), routes })

// 执行 router.beforeEach((to) => {，实现对应的前端逻辑

router.beforeEach((to) => {
  // 条件判断：if (to.path === '/login') return true 时执行下面逻辑

  if (to.path === '/login') return true
  // 条件判断：if (!getToken()) return { path: '/login' } 时执行下面逻辑

  if (!getToken()) return { path: '/login' }
  // 返回 true 给调用方

  return true
// 执行 })，实现对应的前端逻辑

})

// 导出 router，供其他模块引用

export default router
