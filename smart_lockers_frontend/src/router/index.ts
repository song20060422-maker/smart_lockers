// @ticket 00-project：Vue Router 路由表——按工单组织页面
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

// 声明变量 routes: RouteRecordRaw[] = [，用于保存数据或函数引用

const routes: RouteRecordRaw[] = [
  // 执行 { path: '/', redirect: '/home' },，实现对应的前端逻辑

  // @ticket 00-project：根路径重定向到登录页
  { path: '/', redirect: '/login' },
  {
    // @ticket 00-project：登录页（手机号验证码 / 微信授权 演示登录）
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { title: '登录' },
  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/home'

    path: '/home',
    // 定义对象属性 name，取值为 'Home'

    name: 'Home',
    // 定义对象属性 component，取值为 () => import('@/views/HomeView.vue')

    component: () => import('@/views/HomeView.vue'),
    // 定义对象属性 meta，取值为 { title: '首页' }

    meta: { title: '首页' },
  // 定义结束

  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/deposit'

    path: '/deposit',
    // 定义对象属性 name，取值为 'Deposit'

    name: 'Deposit',
    // 定义对象属性 component，取值为 () => import('@/views/DepositView.vue')

    component: () => import('@/views/DepositView.vue'),
    // 定义对象属性 meta，取值为 { title: '物品寄存' }

    meta: { title: '物品寄存' },
  // 定义结束

  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/orders'

    path: '/orders',
    // 定义对象属性 name，取值为 'Orders'

    name: 'Orders',
    // 定义对象属性 component，取值为 () => import('@/views/OrdersView.vue')

    component: () => import('@/views/OrdersView.vue'),
    // 定义对象属性 meta，取值为 { title: '我的订单' }

    meta: { title: '我的订单' },
  // 定义结束

  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/orders/:orderNo'

    path: '/orders/:orderNo',
    // 定义对象属性 name，取值为 'OrderDetail'

    name: 'OrderDetail',
    // 定义对象属性 component，取值为 () => import('@/views/OrderDetailView.vue')

    component: () => import('@/views/OrderDetailView.vue'),
    // 定义对象属性 meta，取值为 { title: '订单详情' }

    meta: { title: '订单详情' },
  // 定义结束

  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/pay'

    path: '/pay',
    // 定义对象属性 name，取值为 'Pay'

    name: 'Pay',
    // 定义对象属性 component，取值为 () => import('@/views/PayView.vue')

    component: () => import('@/views/PayView.vue'),
    // 定义对象属性 meta，取值为 { title: '余额充值' }

    meta: { title: '余额充值' },
  // 定义结束

  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/user'

    path: '/user',
    // 定义对象属性 name，取值为 'User'

    name: 'User',
    // 定义对象属性 component，取值为 () => import('@/views/UserView.vue')

    component: () => import('@/views/UserView.vue'),
    // 定义对象属性 meta，取值为 { title: '个人中心' }

    meta: { title: '个人中心' },
  // 定义结束

  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/notify'

    path: '/notify',
    // 定义对象属性 name，取值为 'Notify'

    name: 'Notify',
    // 定义对象属性 component，取值为 () => import('@/views/NotifyView.vue')

    component: () => import('@/views/NotifyView.vue'),
    // 定义对象属性 meta，取值为 { title: '消息通知' }

    meta: { title: '消息通知' },
  // 定义结束

  },
  // 进入该代码块，开始执行其中的逻辑

  {
    // 定义对象属性 path，取值为 '/ops'

    path: '/ops',
    // 定义对象属性 name，取值为 'Ops'

    name: 'Ops',
    // 定义对象属性 component，取值为 () => import('@/views/OpsView.vue')

    component: () => import('@/views/OpsView.vue'),
    // 定义对象属性 meta，取值为 { title: '性能运维' }

    meta: { title: '性能运维' },
  // 定义结束

  },
  // ===== 补齐：原型图与工单对应的缺失界面（只新增路由，不改动上方原有路由） =====
  // @ticket 06-个人中心 + @proto 01-C端优惠券 / 06-用户界面规划：我的优惠券
  { path: '/coupons', name: 'Coupons', component: () => import('@/views/CouponsView.vue'), meta: { title: '我的优惠券' } },
  // @ticket 03-首页 + @proto 04-一期用户界面-2：选择城市
  { path: '/city', name: 'City', component: () => import('@/views/CityView.vue'), meta: { title: '选择城市' } },
  // @ticket 03-首页 + @proto 04-一期用户界面-3/4/5：附近寄存点
  { path: '/nearby', name: 'Nearby', component: () => import('@/views/NearbyView.vue'), meta: { title: '附近寄存点' } },
  // @ticket 03-首页 + @proto 05-寄存指南长图：寄存指南
  { path: '/guide', name: 'Guide', component: () => import('@/views/GuideView.vue'), meta: { title: '寄存指南' } },
  // @ticket 04-物品寄存 + @proto 03-预支付相关提示：费用确认
  { path: '/pay-confirm', name: 'PayConfirm', component: () => import('@/views/PayConfirmView.vue'), meta: { title: '确认支付' } },
  // @ticket 04-物品寄存 + @proto 03-预支付-4：支付完成开柜
  { path: '/pay-success', name: 'PaySuccess', component: () => import('@/views/PaySuccessView.vue'), meta: { title: '支付完成' } },
  // @ticket 04/05 + @proto 02-开门异常提示优化-1~7：开门异常
  { path: '/door-exception', name: 'DoorException', component: () => import('@/views/DoorExceptionView.vue'), meta: { title: '开门提示' } },
  // @ticket 05-订单查看与取出 + @proto 04/06：超时补交
  { path: '/supplement', name: 'Supplement', component: () => import('@/views/SupplementView.vue'), meta: { title: '超时补交' } },
  // @ticket 04-物品寄存 + @proto 06-用户界面规划-18：人脸识别
  { path: '/face', name: 'Face', component: () => import('@/views/FaceView.vue'), meta: { title: '人脸识别' } },
  // @ticket 06-个人中心 + @proto 06-用户界面规划-1/19：在线客服
  { path: '/service', name: 'Service', component: () => import('@/views/ServiceView.vue'), meta: { title: '在线客服' } },
  // @ticket 05-订单查看与取出 + @proto 04-一期用户界面-14：取出类型选择
  { path: '/pickup', name: 'Pickup', component: () => import('@/views/PickupView.vue'), meta: { title: '取出类型' } },
// 结束切片/数组字面量

]

// 声明变量 router = createRouter({，用于保存数据或函数引用

const router = createRouter({
  // 定义对象属性 history，取值为 createWebHistory()

  history: createWebHistory(),
  // 执行 routes,，实现对应的前端逻辑

  routes,
// 执行 })，实现对应的前端逻辑

})

// @ticket 00-project：切换页面更新标题
// 执行 router.afterEach((to) => {，实现对应的前端逻辑

router.afterEach((to) => {
  // 声明变量 t = to.meta.title as string | undefined，用于保存数据或函数引用

  const t = to.meta.title as string | undefined
  // 赋值更新：document.title = t ? `${t} - 智能快寄柜` : '智能快寄柜'

  document.title = t ? `${t} - 智能快寄柜` : '智能快寄柜'
// 执行 })，实现对应的前端逻辑

})

// 导出 router，供其他模块引用

export default router
