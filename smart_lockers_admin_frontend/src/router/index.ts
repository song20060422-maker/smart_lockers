// @ticket B端-前端：路由与登录守卫（三类角色共用页面，按钮级由权限控制）
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { getToken } from '@/api/request'
import Layout from '@/layouts/MainLayout.vue'

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue') },
  {
    path: '/',
    component: Layout,
    redirect: '/sites',
    children: [
      { path: 'sites', name: 'sites', component: () => import('@/views/SitesView.vue'), meta: { title: '网点管理', perm: 'site:read' } },
      { path: 'devices', name: 'devices', component: () => import('@/views/DevicesView.vue'), meta: { title: '设备监控', perm: 'device:read' } },
      { path: 'reconcile', name: 'reconcile', component: () => import('@/views/ReconcileView.vue'), meta: { title: '充值对账', perm: 'reconcile:read' } },
      { path: 'alarms', name: 'alarms', component: () => import('@/views/AlarmsView.vue'), meta: { title: '告警处理', perm: 'alarm:read' } },
    ],
  },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to) => {
  if (to.path === '/login') return true
  if (!getToken()) return { path: '/login' }
  return true
})

export default router
