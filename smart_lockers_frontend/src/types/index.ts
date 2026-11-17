// @ticket 00-project：全局 TypeScript 类型定义，字段与后端统一报文/各 service 返回对齐

// @ticket 00-project：统一响应外壳，后端 {code,message,data,request_id}，code=0 成功
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
  request_id?: string
}

// @ticket 00-project：业务错误（响应拦截器 reject 后抛出）
export interface ApiError {
  code: number
  message: string
  request_id?: string
}

// @ticket 00-project：分页参数
export interface PageQuery {
  page?: number
  page_size?: number
}

// ================= 工单 06 登录/个人中心 =================
// @ticket 06-个人中心：登录响应
export interface LoginResult {
  token: string
}
// @ticket 06-个人中心：我的服务项（对应 model.MyService）
export interface MyService {
  id: number
  title: string
  icon?: string
  path?: string
  roles?: string
  enabled?: boolean
}
// @ticket 06-个人中心：卡券（对应 model.Coupon）
export interface Coupon {
  id: number
  code: string
  discount: number
  expire_at: string
  used: boolean
}
// @ticket 06-个人中心：卡券列表响应
export interface CouponList {
  available: Coupon[]
  expired: Coupon[]
}

// ================= 工单 03 首页 =================
// @ticket 03-首页：定位结果
export interface LocateResult {
  source: string
  city?: string
  lat: number
  lng: number
}
// @ticket 03-首页：城市（对应 model.City）
export interface City {
  id: number
  name: string
  pinyin: string
  initial: string
  sort: number
}
// @ticket 03-首页：附近寄存点（LockerSite + distance）
export interface LockerSite {
  id: number
  name: string
  address: string
  lat: number
  lng: number
  status: number
  distance?: number
}
// @ticket 03-首页：附近寄存点列表响应
export interface NearbySitesResult {
  list: LockerSite[]
  total: number
  page: number
  page_size: number
}
// @ticket 03-首页：寄存指南
export interface GuideResult {
  version: string
  lang: string
  content: string
}

// ================= 工单 04 物品寄存 =================
// @ticket 04-物品寄存：柜子分组（CabinetQuery）
export interface Cabinet {
  id: number
  site_id: number
  code: string
  size: string
  total: number
  available: number
}
// @ticket 04-物品寄存：费用试算结果（FeeResult，金额单位：分）
export interface FeeResult {
  base_fee: number
  discount: number
  deposit: number
  payable: number
}
// @ticket 04-物品寄存：创建寄存订单响应
export interface CreateDepositResult {
  order_no: string
  payable?: number
  idempotent: boolean
}
// @ticket 04-物品寄存：换柜响应
export interface ChangeCabinetResult {
  order_no: string
  door_id: number
}

// ================= 工单 02 支付宝支付 =================
// @ticket 02-支付宝支付：余额
export interface BalanceResult {
  user_id: number
  balance: number
  currency: string
}
// @ticket 02-支付宝支付：余额明细（对应 model.BalanceLog）
export interface BalanceLog {
  id: number
  user_id: number
  change: number
  balance_after: number
  biz_type: string
  biz_no: string
  created_at: string
}
// @ticket 02-支付宝支付：充值单（对应 model.RechargeOrder）
export interface RechargeOrder {
  id: number
  charge_no: string
  user_id: number
  amount: number
  status: number
  trade_no?: string
  paid_at?: string
  created_at: string
}
// @ticket 02-支付宝支付：充值创建响应
export interface CreateRechargeResult {
  charge_no: string
  amount: number
  status: number
}
// @ticket 02-支付宝支付：分页列表通用结构
export interface ListResult<T> {
  list: T[]
  page: number
  page_size: number
  total?: number
}

// ================= 工单 05 订单查看与取出 =================
// @ticket 05-订单查看与取出：订单视图（OrderVO）
export interface OrderVO {
  order_no: string
  site_id: number
  door_id: number
  start_time: string
  end_time: string
  amount: number
  deposit: number
  status: number
}
// @ticket 05-订单查看与取出：游标分页订单列表
export interface CursorOrderResult {
  list: OrderVO[]
  next_cursor: string
  has_more: boolean
}
// @ticket 05-订单查看与取出：订单详情
export interface OrderDetail extends OrderVO {
  site?: {
    name: string
    address: string
    lat: number
    lng: number
  }
}
// @ticket 05-订单查看与取出：进行中聚合
export interface OngoingSummary {
  ongoing_count: number
  total_amount: number
}
// @ticket 05-订单查看与取出：临时开柜授权
export interface OpenAuthResult {
  token: string
  expire_at: string
  max_times: number
}

// ================= 工单 01 消息通知 =================
// @ticket 01-消息通知：发送结果
export interface NotifySendResult {
  message_id?: number
  content?: string
  deduped?: boolean
}
// @ticket 01-消息通知：通知消息（对应 model.NotifyMessage）
export interface NotifyMessage {
  id: number
  tpl_code: string
  receiver_type: string
  receiver_id: number
  channel: string
  content: string
  status: number
  dedupe_key: string
  retry_count: number
  created_at: string
}
// @ticket 01-消息通知：监控统计
export interface NotifyMonitor {
  pending: number
  success: number
  failed: number
  retrying: number
}

// @ticket 00-project：通用操作结果
export interface OkResult {
  logout?: boolean
  ok?: boolean
  idempotent?: boolean
  [key: string]: unknown
}


// ================= 别名与补充类型（对齐各 API 模块引用） =================
// @ticket 05-订单查看与取出：订单列表结果别名
export type OrderListResult = CursorOrderResult
// @ticket 05-订单查看与取出：进行中聚合别名
export type OngoingResult = OngoingSummary
// @ticket 05-订单查看与取出：补缴结果
export interface MakeupResult {
  order_no: string
  payable: number
  idempotent?: boolean
}

// @ticket 06-个人中心：我的服务列表响应
export interface MyServicesResult {
  list: MyService[]
}

// @ticket 01-消息通知：通知记录分页响应
export interface NotifyRecordsResult {
  list: NotifyMessage[]
  total: number
  page: number
  page_size: number
}
// @ticket 01-消息通知：监控统计别名
export type MonitorResult = NotifyMonitor
// @ticket 01-消息通知：通知场景
export interface Scene {
  code: string
  name: string
  channels?: string[]
}

// @ticket 08-性能调优：健康检查响应
export interface HealthResult {
  status: string
  time?: string
}
// @ticket 08-性能调优：pprof 采样结果
export interface OpsPprofResult {
  file?: string
  url?: string
  [key: string]: unknown
}
// @ticket 08-性能调优：trace 采样结果
export interface OpsTraceResult {
  file?: string
  url?: string
  [key: string]: unknown
}
// @ticket 08-性能调优：expvar 指标
export type OpsExpvarResult = Record<string, unknown>


// @ticket 06-个人中心：头像更新结果
export interface AvatarResult {
  url: string
  ok?: boolean
}


// @ticket 06-个人中心：订单汇总（对应 handler.OrderSummary）
export interface OrderSummaryResult {
  total?: number
  ongoing?: number
  done?: number
  [key: string]: unknown
}
