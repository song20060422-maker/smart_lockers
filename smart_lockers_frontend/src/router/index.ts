// @ticket 00-project：Vue Router 路由表——按工单组织页面
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/home' },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
    meta: { title: '首页' },
  },
  {
    path: '/deposit',
    name: 'Deposit',
    component: () => import('@/views/DepositView.vue'),
    meta: { title: '物品寄存' },
  },
  {
    path: '/orders',
    name: 'Orders',
    component: () => import('@/views/OrdersView.vue'),
    meta: { title: '我的订单' },
  },
  {
    path: '/orders/:orderNo',
    name: 'OrderDetail',
    component: () => import('@/views/OrderDetailView.vue'),
    meta: { title: '订单详情' },
  },
  {
    path: '/pay',
    name: 'Pay',
    component: () => import('@/views/PayView.vue'),
    meta: { title: '余额充值' },
  },
  {
    path: '/user',
    name: 'User',
    component: () => import('@/views/UserView.vue'),
    meta: { title: '个人中心' },
  },
  {
    path: '/notify',
    name: 'Notify',
    component: () => import('@/views/NotifyView.vue'),
    meta: { title: '消息通知' },
  },
  {
    path: '/ops',
    name: 'Ops',
    component: () => import('@/views/OpsView.vue'),
    meta: { title: '性能运维' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// @ticket 00-project：切换页面更新标题
router.afterEach((to) => {
  const t = to.meta.title as string | undefined
  document.title = t ? `${t} - 智能快寄柜` : '智能快寄柜'
})

export default router
