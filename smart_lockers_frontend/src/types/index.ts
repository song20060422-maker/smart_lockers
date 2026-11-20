// @ticket 00-project：全局 TypeScript 类型定义，字段与后端统一报文/各 service 返回对齐

// @ticket 00-project：统一响应外壳，后端 {code,message,data,request_id}，code=0 成功
// 导出 interface ApiResponse<T = unknown> {，供其他模块引用

export interface ApiResponse<T = unknown> {
  // 定义对象属性 code，取值为 number

  code: number
  // 定义对象属性 message，取值为 string

  message: string
  // 执行 data?: T，实现对应的前端逻辑

  data?: T
  // 执行 request_id?: string，实现对应的前端逻辑

  request_id?: string
// 当前代码块结束

}

// @ticket 00-project：业务错误（响应拦截器 reject 后抛出）
// 导出 interface ApiError {，供其他模块引用

export interface ApiError {
  // 定义对象属性 code，取值为 number

  code: number
  // 定义对象属性 message，取值为 string

  message: string
  // 执行 request_id?: string，实现对应的前端逻辑

  request_id?: string
// 当前代码块结束

}

// @ticket 00-project：分页参数
// 导出 interface PageQuery {，供其他模块引用

export interface PageQuery {
  // 执行 page?: number，实现对应的前端逻辑

  page?: number
  // 执行 page_size?: number，实现对应的前端逻辑

  page_size?: number
// 当前代码块结束

}

// ================= 工单 06 登录/个人中心 =================
// @ticket 06-个人中心：登录响应
// 导出 interface LoginResult {，供其他模块引用

export interface LoginResult {
  // 定义对象属性 token，取值为 string

  token: string
// 当前代码块结束

}

// 修复偏差4：演示登录响应，含后端签发的真实 token 与 user_id @ticket 00-project
// 导出 interface DemoLoginResult：与 POST /auth/demo-login 后端返回字段一致

export interface DemoLoginResult {
  // 后端签发的 JWT token

  token: string
  // 登录用户 ID

  user_id: number
  // 登录手机号

  phone?: string
// 当前代码块结束

}
// @ticket 06-个人中心：我的服务项（对应 model.MyService）
// 导出 interface MyService {，供其他模块引用

export interface MyService {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 title，取值为 string

  title: string
  // 执行 icon?: string，实现对应的前端逻辑

  icon?: string
  // 执行 path?: string，实现对应的前端逻辑

  path?: string
  // 执行 roles?: string，实现对应的前端逻辑

  roles?: string
  // 执行 enabled?: boolean，实现对应的前端逻辑

  enabled?: boolean
// 当前代码块结束

}
// @ticket 06-个人中心：卡券（对应 model.Coupon）
// 导出 interface Coupon {，供其他模块引用

export interface Coupon {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 code，取值为 string

  code: string
  // 定义对象属性 discount，取值为 number

  discount: number
  // 定义对象属性 expire_at，取值为 string

  expire_at: string
  // 定义对象属性 used，取值为 boolean

  used: boolean
// 当前代码块结束

}
// @ticket 06-个人中心：卡券列表响应
// 导出 interface CouponList {，供其他模块引用

export interface CouponList {
  // 定义对象属性 available，取值为 Coupon[]

  available: Coupon[]
  // 定义对象属性 expired，取值为 Coupon[]

  expired: Coupon[]
// 当前代码块结束

}

// ================= 工单 03 首页 =================
// @ticket 03-首页：定位结果
// 导出 interface LocateResult {，供其他模块引用

export interface LocateResult {
  // 定义对象属性 source，取值为 string

  source: string
  // 执行 city?: string，实现对应的前端逻辑

  city?: string
  // 定义对象属性 lat，取值为 number

  lat: number
  // 定义对象属性 lng，取值为 number

  lng: number
// 当前代码块结束

}
// @ticket 03-首页：城市（对应 model.City）
// 导出 interface City {，供其他模块引用

export interface City {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 name，取值为 string

  name: string
  // 定义对象属性 pinyin，取值为 string

  pinyin: string
  // 定义对象属性 initial，取值为 string

  initial: string
  // 定义对象属性 sort，取值为 number

  sort: number
// 当前代码块结束

}
// @ticket 03-首页：附近寄存点（LockerSite + distance）
// 导出 interface LockerSite {，供其他模块引用

export interface LockerSite {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 name，取值为 string

  name: string
  // 定义对象属性 address，取值为 string

  address: string
  // 定义对象属性 lat，取值为 number

  lat: number
  // 定义对象属性 lng，取值为 number

  lng: number
  // 定义对象属性 status，取值为 number

  status: number
  // 执行 distance?: number，实现对应的前端逻辑

  distance?: number
// 当前代码块结束

}
// @ticket 03-首页：附近寄存点列表响应
// 导出 interface NearbySitesResult {，供其他模块引用

export interface NearbySitesResult {
  // 定义对象属性 list，取值为 LockerSite[]

  list: LockerSite[]
  // 定义对象属性 total，取值为 number

  total: number
  // 定义对象属性 page，取值为 number

  page: number
  // 定义对象属性 page_size，取值为 number

  page_size: number
// 当前代码块结束

}
// @ticket 03-首页：寄存指南
// 导出 interface GuideResult {，供其他模块引用

export interface GuideResult {
  // 定义对象属性 version，取值为 string

  version: string
  // 定义对象属性 lang，取值为 string

  lang: string
  // 定义对象属性 content，取值为 string

  content: string
// 当前代码块结束

}

// ================= 工单 04 物品寄存 =================
// @ticket 04-物品寄存：柜子分组（CabinetQuery）
// 导出 interface Cabinet {，供其他模块引用

export interface Cabinet {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 site_id，取值为 number

  site_id: number
  // 定义对象属性 code，取值为 string

  code: string
  // 定义对象属性 size，取值为 string

  size: string
  // 定义对象属性 total，取值为 number

  total: number
  // 定义对象属性 available，取值为 number

  available: number
// 当前代码块结束

}
// @ticket 04-物品寄存：费用试算结果（FeeResult，金额单位：分）
// 导出 interface FeeResult {，供其他模块引用

export interface FeeResult {
  // 定义对象属性 base_fee，取值为 number

  base_fee: number
  // 定义对象属性 discount，取值为 number

  discount: number
  // 定义对象属性 deposit，取值为 number

  deposit: number
  // 定义对象属性 payable，取值为 number

  payable: number
// 当前代码块结束

}
// @ticket 04-物品寄存：创建寄存订单响应
// 导出 interface CreateDepositResult {，供其他模块引用

export interface CreateDepositResult {
  // 定义对象属性 order_no，取值为 string

  order_no: string
  // 执行 payable?: number，实现对应的前端逻辑

  payable?: number
  // 定义对象属性 idempotent，取值为 boolean

  idempotent: boolean
// 当前代码块结束

}
// @ticket 04-物品寄存：换柜响应
// 导出 interface ChangeCabinetResult {，供其他模块引用

export interface ChangeCabinetResult {
  // 定义对象属性 order_no，取值为 string

  order_no: string
  // 定义对象属性 door_id，取值为 number

  door_id: number
// 当前代码块结束

}

// ================= 工单 02 支付宝支付 =================
// @ticket 02-支付宝支付：余额
// 导出 interface BalanceResult {，供其他模块引用

export interface BalanceResult {
  // 定义对象属性 user_id，取值为 number

  user_id: number
  // 定义对象属性 balance，取值为 number

  balance: number
  // 定义对象属性 currency，取值为 string

  currency: string
// 当前代码块结束

}
// @ticket 02-支付宝支付：余额明细（对应 model.BalanceLog）
// 导出 interface BalanceLog {，供其他模块引用

export interface BalanceLog {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 user_id，取值为 number

  user_id: number
  // 定义对象属性 change，取值为 number

  change: number
  // 定义对象属性 balance_after，取值为 number

  balance_after: number
  // 定义对象属性 biz_type，取值为 string

  biz_type: string
  // 定义对象属性 biz_no，取值为 string

  biz_no: string
  // 定义对象属性 created_at，取值为 string

  created_at: string
// 当前代码块结束

}
// @ticket 02-支付宝支付：充值单（对应 model.RechargeOrder）
// 导出 interface RechargeOrder {，供其他模块引用

export interface RechargeOrder {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 charge_no，取值为 string

  charge_no: string
  // 定义对象属性 user_id，取值为 number

  user_id: number
  // 定义对象属性 amount，取值为 number

  amount: number
  // 定义对象属性 status，取值为 number

  status: number
  // 执行 trade_no?: string，实现对应的前端逻辑

  trade_no?: string
  // 执行 paid_at?: string，实现对应的前端逻辑

  paid_at?: string
  // 定义对象属性 created_at，取值为 string

  created_at: string
// 当前代码块结束

}
// @ticket 02-支付宝支付：充值创建响应
// 导出 interface CreateRechargeResult {，供其他模块引用

export interface CreateRechargeResult {
  // 定义对象属性 charge_no，取值为 string

  charge_no: string
  // 定义对象属性 amount，取值为 number

  amount: number
  // 定义对象属性 status，取值为 number

  status: number
// 当前代码块结束

}
// @ticket 02-支付宝支付：分页列表通用结构
// 导出 interface ListResult<T> {，供其他模块引用

export interface ListResult<T> {
  // 定义对象属性 list，取值为 T[]

  list: T[]
  // 定义对象属性 page，取值为 number

  page: number
  // 定义对象属性 page_size，取值为 number

  page_size: number
  // 执行 total?: number，实现对应的前端逻辑

  total?: number
// 当前代码块结束

}

// ================= 工单 05 订单查看与取出 =================
// @ticket 05-订单查看与取出：订单视图（OrderVO）
// 导出 interface OrderVO {，供其他模块引用

export interface OrderVO {
  // 定义对象属性 order_no，取值为 string

  order_no: string
  // 定义对象属性 site_id，取值为 number

  site_id: number
  // 定义对象属性 door_id，取值为 number

  door_id: number
  // 定义对象属性 start_time，取值为 string

  start_time: string
  // 定义对象属性 end_time，取值为 string

  end_time: string
  // 定义对象属性 amount，取值为 number

  amount: number
  // 定义对象属性 deposit，取值为 number

  deposit: number
  // 定义对象属性 status，取值为 number

  status: number
// 当前代码块结束

}
// @ticket 05-订单查看与取出：游标分页订单列表
// 导出 interface CursorOrderResult {，供其他模块引用

export interface CursorOrderResult {
  // 定义对象属性 list，取值为 OrderVO[]

  list: OrderVO[]
  // 定义对象属性 next_cursor，取值为 string

  next_cursor: string
  // 定义对象属性 has_more，取值为 boolean

  has_more: boolean
// 当前代码块结束

}
// @ticket 05-订单查看与取出：订单详情
// 导出 interface OrderDetail extends OrderVO {，供其他模块引用

export interface OrderDetail extends OrderVO {
  // 执行 site?: {，实现对应的前端逻辑

  site?: {
    // 定义对象属性 name，取值为 string

    name: string
    // 定义对象属性 address，取值为 string

    address: string
    // 定义对象属性 lat，取值为 number

    lat: number
    // 定义对象属性 lng，取值为 number

    lng: number
  // 当前代码块结束

  }
// 当前代码块结束

}
// @ticket 05-订单查看与取出：进行中聚合
// 导出 interface OngoingSummary {，供其他模块引用

export interface OngoingSummary {
  // 定义对象属性 ongoing_count，取值为 number

  ongoing_count: number
  // 定义对象属性 total_amount，取值为 number

  total_amount: number
// 当前代码块结束

}
// @ticket 05-订单查看与取出：临时开柜授权
// 导出 interface OpenAuthResult {，供其他模块引用

export interface OpenAuthResult {
  // 定义对象属性 token，取值为 string

  token: string
  // 定义对象属性 expire_at，取值为 string

  expire_at: string
  // 定义对象属性 max_times，取值为 number

  max_times: number
// 当前代码块结束

}

// ================= 工单 01 消息通知 =================
// @ticket 01-消息通知：发送结果
// 导出 interface NotifySendResult {，供其他模块引用

export interface NotifySendResult {
  // 执行 message_id?: number，实现对应的前端逻辑

  message_id?: number
  // 执行 content?: string，实现对应的前端逻辑

  content?: string
  // 执行 deduped?: boolean，实现对应的前端逻辑

  deduped?: boolean
// 当前代码块结束

}
// @ticket 01-消息通知：通知消息（对应 model.NotifyMessage）
// 导出 interface NotifyMessage {，供其他模块引用

export interface NotifyMessage {
  // 定义对象属性 id，取值为 number

  id: number
  // 定义对象属性 tpl_code，取值为 string

  tpl_code: string
  // 定义对象属性 receiver_type，取值为 string

  receiver_type: string
  // 定义对象属性 receiver_id，取值为 number

  receiver_id: number
  // 定义对象属性 channel，取值为 string

  channel: string
  // 定义对象属性 content，取值为 string

  content: string
  // 定义对象属性 status，取值为 number

  status: number
  // 定义对象属性 dedupe_key，取值为 string

  dedupe_key: string
  // 定义对象属性 retry_count，取值为 number

  retry_count: number
  // 定义对象属性 created_at，取值为 string

  created_at: string
// 当前代码块结束

}
// @ticket 01-消息通知：监控统计
// 导出 interface NotifyMonitor {，供其他模块引用

export interface NotifyMonitor {
  // 定义对象属性 pending，取值为 number

  pending: number
  // 定义对象属性 success，取值为 number

  success: number
  // 定义对象属性 failed，取值为 number

  failed: number
  // 定义对象属性 retrying，取值为 number

  retrying: number
// 当前代码块结束

}

// @ticket 00-project：通用操作结果
// 导出 interface OkResult {，供其他模块引用

export interface OkResult {
  // 执行 logout?: boolean，实现对应的前端逻辑

  logout?: boolean
  // 执行 ok?: boolean，实现对应的前端逻辑

  ok?: boolean
  // 执行 idempotent?: boolean，实现对应的前端逻辑

  idempotent?: boolean
  // 执行 [key: string]: unknown，实现对应的前端逻辑

  [key: string]: unknown
// 当前代码块结束

}


// ================= 别名与补充类型（对齐各 API 模块引用） =================
// @ticket 05-订单查看与取出：订单列表结果别名
// 导出 type OrderListResult = CursorOrderResult，供其他模块引用

export type OrderListResult = CursorOrderResult
// @ticket 05-订单查看与取出：进行中聚合别名
// 导出 type OngoingResult = OngoingSummary，供其他模块引用

export type OngoingResult = OngoingSummary
// @ticket 05-订单查看与取出：补缴结果
// 导出 interface MakeupResult {，供其他模块引用

export interface MakeupResult {
  // 定义对象属性 order_no，取值为 string

  order_no: string
  // 定义对象属性 payable，取值为 number

  payable: number
  // 执行 idempotent?: boolean，实现对应的前端逻辑

  idempotent?: boolean
// 当前代码块结束

}

// @ticket 06-个人中心：我的服务列表响应
// 导出 interface MyServicesResult {，供其他模块引用

export interface MyServicesResult {
  // 定义对象属性 list，取值为 MyService[]

  list: MyService[]
// 当前代码块结束

}

// @ticket 01-消息通知：通知记录分页响应
// 导出 interface NotifyRecordsResult {，供其他模块引用

export interface NotifyRecordsResult {
  // 定义对象属性 list，取值为 NotifyMessage[]

  list: NotifyMessage[]
  // 定义对象属性 total，取值为 number

  total: number
  // 定义对象属性 page，取值为 number

  page: number
  // 定义对象属性 page_size，取值为 number

  page_size: number
// 当前代码块结束

}
// @ticket 01-消息通知：监控统计别名
// 导出 type MonitorResult = NotifyMonitor，供其他模块引用

export type MonitorResult = NotifyMonitor
// @ticket 01-消息通知：通知场景
// 导出 interface Scene {，供其他模块引用

export interface Scene {
  // 定义对象属性 code，取值为 string

  code: string
  // 定义对象属性 name，取值为 string

  name: string
  // 执行 channels?: string[]，实现对应的前端逻辑

  channels?: string[]
// 当前代码块结束

}

// @ticket 08-性能调优：健康检查响应
// 导出 interface HealthResult {，供其他模块引用

export interface HealthResult {
  // 定义对象属性 status，取值为 string

  status: string
  // 执行 time?: string，实现对应的前端逻辑

  time?: string
// 当前代码块结束

}
// @ticket 08-性能调优：pprof 采样结果
// 导出 interface OpsPprofResult {，供其他模块引用

export interface OpsPprofResult {
  // 执行 file?: string，实现对应的前端逻辑

  file?: string
  // 执行 url?: string，实现对应的前端逻辑

  url?: string
  // 执行 [key: string]: unknown，实现对应的前端逻辑

  [key: string]: unknown
// 当前代码块结束

}
// @ticket 08-性能调优：trace 采样结果
// 导出 interface OpsTraceResult {，供其他模块引用

export interface OpsTraceResult {
  // 执行 file?: string，实现对应的前端逻辑

  file?: string
  // 执行 url?: string，实现对应的前端逻辑

  url?: string
  // 执行 [key: string]: unknown，实现对应的前端逻辑

  [key: string]: unknown
// 当前代码块结束

}
// @ticket 08-性能调优：expvar 指标
// 导出 type OpsExpvarResult = Record<string, unknown>，供其他模块引用

export type OpsExpvarResult = Record<string, unknown>


// @ticket 06-个人中心：头像更新结果
// 导出 interface AvatarResult {，供其他模块引用

export interface AvatarResult {
  // 定义对象属性 url，取值为 string

  url: string
  // 执行 ok?: boolean，实现对应的前端逻辑

  ok?: boolean
// 当前代码块结束

}


// @ticket 06-个人中心：订单汇总（对应 handler.OrderSummary）
// 导出 interface OrderSummaryResult {，供其他模块引用

export interface OrderSummaryResult {
  // 执行 total?: number，实现对应的前端逻辑

  total?: number
  // 执行 ongoing?: number，实现对应的前端逻辑

  ongoing?: number
  // 执行 done?: number，实现对应的前端逻辑

  done?: number
  // 执行 [key: string]: unknown，实现对应的前端逻辑

  [key: string]: unknown
// 当前代码块结束

}
