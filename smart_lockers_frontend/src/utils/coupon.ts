// @ticket 06-个人中心 / 04-物品寄存：优惠券字段口径统一与空值防御（只新增，不改动既有代码）
// @proto 原型对齐：01-C端优惠券-1~5、06-13/14；券字段以后端 GET /api/v1/user/coupons 的 { available, expired } 为准

// 后端 Coupon 结构：{ id, code, discount(int64 分), scope, expire_at, used }
// 前端展示兼容 amount(元)/hours(时长) 两类演示券，统一归一到本结构
export interface NormalizedCoupon {
  // 券唯一标识（后端 id，恒为 number；缺失时兜底 0）
  id: number
  // 券码
  code: string
  // 券名称（后端无 title 时按面值推断）
  title: string
  // 抵扣金额，单位元（由 discount 分换算；金额券）
  amount: number
  // 抵扣时长，单位小时（时长券，后端 discount=0 且无金额时为 0）
  hours: number
  // 折扣文案（如「8折」，非折扣券为空）
  discountText: string
  // 有效期字符串
  expire_at: string
  // 适用范围
  scope: string
  // 是否已使用
  used: boolean
  // 是否已过期（来自 expired 分组）
  expired: boolean
}

// 将后端优惠券响应（对象/数组/空值）归一为统一数组；空值一律防御，绝不产生 NaN
export function normalizeCoupons(raw: unknown): NormalizedCoupon[] {
  // 后端返回 { available, expired } 对象：先拆成两个分组再合并
  if (raw && typeof raw === 'object' && !Array.isArray(raw)) {
    // available 分组默认空数组
    const available = (raw as Record<string, unknown>).available
    // expired 分组默认空数组
    const expired = (raw as Record<string, unknown>).expired
    // 合并两组：可用在前、过期在后，并打上 expired 标记
    return [
      // 归一化可用分组
      ...(Array.isArray(available) ? available.map((c) => normalizeOne(c, false)) : []),
      // 归一化过期分组
      ...(Array.isArray(expired) ? expired.map((c) => normalizeOne(c, true)) : []),
    ]
  }
  // 兼容旧版数组返回：全部按可用处理
  if (Array.isArray(raw)) {
    // 逐条归一化
    return raw.map((c) => normalizeOne(c, false))
  }
  // 空值（null/undefined/其他）防御：返回空数组
  return []
}

// 归一化单张券；对缺失字段做空值防御，金额换算不做除法直到确认是分值
export function normalizeOne(raw: unknown, expired: boolean): NormalizedCoupon {
  // 非对象输入防御：返回占位券，避免渲染报错
  if (!raw || typeof raw !== 'object') {
    // 返回全空占位券
    return {
      id: 0,
      code: '',
      title: '',
      amount: 0,
      hours: 0,
      discountText: '',
      expire_at: '',
      scope: '',
      used: false,
      expired,
    }
  }
  // 原始券对象
  const c = raw as Record<string, unknown>
  // 券码兜底为空串
  const code = typeof c.code === 'string' ? c.code : ''
  // 数值字段安全读取：非有限数一律 0，杜绝 NaN
  const readNum = (v: unknown): number => {
    // number 且有限才保留，否则 0
    return typeof v === 'number' && Number.isFinite(v) ? v : 0
  }
  // 后端 discount 单位分；amount 单位元。优先 amount，其次 discount/100
  const amountFromAmount = readNum(c.amount)
  // 由分换算的金额
  const amountFromDiscount = readNum(c.discount) / 100
  // 最终金额：amount 字段优先，否则用 discount 换算
  const amount = amountFromAmount > 0 ? amountFromAmount : amountFromDiscount
  // 时长券：hours>0 且无金额时视为时长券金额清零
  const hours = readNum(c.hours)
  // 折扣文案：仅接受形如「8折」的字符串
  const discountText = typeof c.discountText === 'string' ? c.discountText : ''
  // 名称兜底
  const title = typeof c.title === 'string' && c.title ? c.title : inferTitle(amount, hours, discountText)
  // 返回统一结构
  return {
    id: typeof c.id === 'number' && Number.isFinite(c.id) ? c.id : 0,
    code,
    title,
    amount,
    hours,
    discountText,
    expire_at: typeof c.expire_at === 'string' ? c.expire_at : '',
    scope: typeof c.scope === 'string' ? c.scope : '',
    used: c.used === true,
    expired,
  }
}

// 按面值推断券名称（仅演示/后端无 title 时使用）
export function inferTitle(amount: number, hours: number, discountText: string): string {
  // 折扣券
  if (discountText) return '折扣券'
  // 时长券
  if (hours > 0) return '时长赠送券'
  // 金额券
  if (amount > 0) return '立减券'
  // 无面值兜底
  return '寄存优惠券'
}

// 生成券面值展示文案（元/时长/折扣/通用），供卡片与下拉共用，空值防御
export function couponFaceText(c: NormalizedCoupon): string {
  // 折扣券
  if (c.discountText) return c.discountText
  // 时长券
  if (c.hours > 0) return `${c.hours}小时`
  // 金额券
  if (c.amount > 0) return `¥${c.amount}`
  // 兜底
  return '券'
}

// 生成寄存页下拉展示文案：券码 + 面值，杜绝 -NaN元
export function couponOptionText(c: NormalizedCoupon): string {
  // 有券码拼接券码
  const prefix = c.code ? `${c.code} ` : ''
  // 返回「券码 面值」
  return `${prefix}${couponFaceText(c)}`
}

// 卡片组件 CouponCard 所需的券结构（与其 prop 类型一致）
export interface CardCoupon {
  // 券 id
  id?: number
  // 券码
  code?: string
  // 券名称
  title?: string
  // 抵扣金额（元）
  amount?: number
  // 抵扣时长（小时）
  hours?: number
  // 有效期
  expire_at?: string
  // 适用范围
  scope?: string
  // 是否过期
  expired?: boolean
}

// 生成卡片展示用的结构（保留 CouponCard 现有 amount/discount/hours 口径）
export function toCardCoupon(c: NormalizedCoupon): CardCoupon {
  // 返回卡片所需字段；折扣文案用 discount 字段承载，其余原样
  return {
    id: c.id,
    code: c.code,
    // 折扣文案并入名称展示，避免 CouponCard 的 discount「N折」口径误读金额分
    title: c.discountText ? `${c.title}（${c.discountText}）` : c.title,
    // 抵扣金额（元），无金额时不传
    amount: c.amount > 0 ? c.amount : undefined,
    // 抵扣时长（小时），无时长时不传
    hours: c.hours > 0 ? c.hours : undefined,
    expire_at: c.expire_at,
    scope: c.scope,
    expired: c.expired,
  }
}
