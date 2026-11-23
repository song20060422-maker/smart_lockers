// @ticket B端-前端：路由与登录守卫（三类角色共用页面，按钮级由权限控制）
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
// 引入鉴权状态仓库，守卫据演示会话判定登录态
import { useAuthStore } from '@/stores/auth'
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
      // IoT 中台总览：设备状态分布/EMQX健康/Topic与ACL规范 @ticket B端-IoT中台
      // 执行 { path: 'iot', name: 'iot', component: () => import('@/views/IotOverviewView.vue'), meta: { title: 'IoT 总览', perm: 'device:read' } },，实现对应的前端逻辑

      { path: 'iot', name: 'iot', component: () => import('@/views/IotOverviewView.vue'), meta: { title: 'IoT 总览', perm: 'device:read' } },
      // OTA 固件升级 @ticket B端-IoT中台
      // 执行 { path: 'ota', name: 'ota', component: () => import('@/views/OtaView.vue'), meta: { title: 'OTA 升级', perm: 'device:read' } },，实现对应的前端逻辑

      { path: 'ota', name: 'ota', component: () => import('@/views/OtaView.vue'), meta: { title: 'OTA 升级', perm: 'device:read' } },
    // 执行 ],，实现对应的前端逻辑

    ],
  // 定义结束

  },
// 结束切片/数组字面量

]

// 声明变量 router = createRouter({ history: createWebHistory(), routes })，用于保存数据或函数引用

const router = createRouter({ history: createWebHistory(), routes })

// 执行 router.beforeEach((to) => {，注册全局登录守卫

// 登录守卫：根路径与所有受保护后台路由统一分流，未建立演示会话一律跳转登录页
router.beforeEach((to) => {
  // 取出鉴权仓库实例，基于本地演示会话判定登录态
  const auth = useAuthStore()
  // 登录页：已登录用户直接回后台，避免重复登录；未登录则放行展示登录页
  if (to.path === '/login') {
    // 已有演示会话时重定向到后台首页
    if (auth.loggedIn) return { path: '/sites' }
    // 未登录时放行进入登录页
    return true
  }
  // 受保护路由未登录：重定向到登录页，并携带原目标以便登录后回跳
  if (!auth.loggedIn) return { path: '/login', query: { redirect: to.fullPath } }
  // 权限拦截：路由声明 meta.perm 且当前角色无该权限时，重定向到第一个有权访问的页面 @ticket B端-RBAC
  if (typeof to.meta.perm === 'string' && !auth.can(to.meta.perm)) {
    // 依次尝试各页面权限，找到当前角色首个可访问页面
    const fallback = ['/sites', '/devices', '/iot', '/ota', '/reconcile', '/alarms'].find((p) => {
      // 从路由表解析目标页面声明的权限标识
      const rec = router.getRoutes().find((r) => r.path === p)
      // 该页面无权限或当前角色拥有权限时可作为兜底页
      return !rec?.meta?.perm || auth.can(rec.meta.perm as string)
    // 执行 })，实现对应的前端逻辑

    })
    // 有兜底页面则跳转，避免越权页面白屏
    return { path: fallback || '/login' }
  }
  // 已登录且具备权限则正常放行
  return true
// 执行 })，实现对应的前端逻辑

})

// 导出 router，供其他模块引用

export default router
