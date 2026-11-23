<script setup lang="ts">
// @ticket 04-物品寄存：柜子信息页——格口选择/寄存时长/优惠券/预估费用/押金与预付/收费说明/确认寄存/更换格口
// @proto 原型对齐：04-一期用户界面-柜子信息系列、01-C端优惠券-2/3/4、03-预支付相关提示-1~4、02-开门异常-5 更换格口弹窗
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { queryCabinets, quoteFee, createDeposit, changeCabinet } from '@/api/deposit'
import { myCoupons } from '@/api/user'
import type { Cabinet, FeeResult } from '@/types'
// 新增统一券数据工具：归一化后端 {available,expired} 对象并对空值做防御，杜绝 -NaN元
import { normalizeCoupons, couponOptionText, type NormalizedCoupon } from '@/utils/coupon'

// 声明变量 route = useRoute()，用于保存数据或函数引用
const route = useRoute()
// 声明变量 router = useRouter()，用于保存数据或函数引用
const router = useRouter()
// 网点信息：从首页/附近寄存跳转 query 携带（对照原型顶部网点名与地址行）
// 声明变量 siteName = ref(...)，用于保存数据或函数引用
const siteName = ref<string>((route.query.site_name as string) || '嵩山北路新途径教育一楼前台')
// 声明变量 siteAddr = ref(...)，用于保存数据或函数引用
const siteAddr = ref<string>((route.query.site_addr as string) || '中交锦兰荟南门处丰巢柜中交锦兰荟南门荟…')
// 声明变量 siteId = Number(route.query.site_id || 0)，用于保存数据或函数引用
const siteId = Number(route.query.site_id || 0)
// 声明变量 cabinets = ref<Cabinet[]>([])，用于保存数据或函数引用
const cabinets = ref<Cabinet[]>([])
// 声明变量 coupons = ref<NormalizedCoupon[]>([])，统一券口径
const coupons = ref<NormalizedCoupon[]>([])
// 声明变量 doorId = ref<number>(0)，用于保存数据或函数引用
const doorId = ref<number>(0)
// 声明变量 hours = ref<number>(3)，用于保存数据或函数引用（原型默认 3 元/小时示例）
const hours = ref<number>(3)
// 声明变量 couponId = ref<number | undefined>(undefined)，用于保存数据或函数引用
const couponId = ref<number | undefined>(undefined)
// 声明变量 fee = ref<FeeResult | null>(null)，用于保存数据或函数引用
const fee = ref<FeeResult | null>(null)
// 计费模式演示切换：hourly 按小时 / daily 按天（对照原型两种价格与收费说明文案）
// 声明变量 mode = ref<'hourly' | 'daily'>('hourly')，用于保存数据或函数引用
const mode = ref<'hourly' | 'daily'>('hourly')
// 声明变量 orderNo = ref('')，用于保存数据或函数引用
const orderNo = ref('')
// 声明变量 err = ref('')，用于保存数据或函数引用
const err = ref('')
// 声明变量 msg = ref('')，用于保存数据或函数引用
const msg = ref('')
// 更换格口弹层显隐（对照原型 02-5 底部弹窗）
// 声明变量 showChange = ref(false)，用于保存数据或函数引用
const showChange = ref(false)

// 格口三档定义（对照原型：小格口-可放文件-400*280*600mm、中格口-可放背包-400*600*600mm、大格口-可放20寸行李箱-400*840*600mm）
// 声明变量 SIZES，用于保存数据或函数引用
const SIZES: Record<string, { label: string; desc: string; dim: string; priceH: string; priceD: string }> = {
  // 小格口
  S: { label: '小格口', desc: '可放文件', dim: '400*280*600mm', priceH: '1元/小时', priceD: '10/天' },
  // 中格口
  M: { label: '中格口', desc: '可放背包', dim: '400*600*600mm', priceH: '2元/小时', priceD: '20元/天' },
  // 大格口
  L: { label: '大格口', desc: '可放20寸行李箱', dim: '400*840*600mm', priceH: '3元/小时', priceD: '30元/天' },
// 当前代码块结束
}
// 按尺寸聚合格口余量（同尺寸多柜合并余量，展示"小格口（余2）"）
// 声明变量 sizeStock，用于保存数据或函数引用
const sizeStock = computed<Record<string, number>>(() => {
  // 聚合结果
  const out: Record<string, number> = {}
  // 循环处理：for of cabinets，逐轮迭代
  for (const c of cabinets.value) {
    // 累加各尺寸可用数
    out[c.size] = (out[c.size] || 0) + c.available
  // 当前代码块结束
  }
  // 返回聚合结果
  return out
// 当前代码块结束
})
// 当前选中格口的尺寸（用于计费试算）
// 声明变量 selSize，用于保存数据或函数引用
const selSize = computed<string>(() => {
  // 从已选柜门反查尺寸
  const c = cabinets.value.find((x) => x.id === doorId.value)
  // 返回尺寸（默认 S）
  return c?.size || 'S'
// 当前代码块结束
})
// 预估费用文案（原型红字"预估费用：¥3"）
// 声明变量 estText，用于保存数据或函数引用
const estText = computed<string>(() => {
  // 已试算则显示应付金额（不含押金：payable 已含押金时展示押金拆分在底部栏）
  if (fee.value) return `¥${(fee.value.payable / 100).toFixed(0)}`
  // 未试算按时长×单价估算（演示）
  return `¥${hours.value}`
// 当前代码块结束
})
// 底部实付/预支付金额（含押金）
// 声明变量 payText，用于保存数据或函数引用
const payText = computed<string>(() => {
  // 已试算用后端应付
  if (fee.value) return (fee.value.payable / 100).toFixed(2).replace(/\.00$/, '')
  // 兜底演示值
  return mode.value === 'hourly' ? '102' : '109'
// 当前代码块结束
})
// 收费说明文案（对照原型：按时/按天两套）
// 声明变量 feeNotes，用于保存数据或函数引用
const feeNotes = computed<string[]>(() => {
  // 按小时模式
  if (mode.value === 'hourly') {
    // 返回按时文案（对照原型）
    return ['租金1元/小时', '寄存后每24小时封顶金额20元', '中途开门不限次数', '寄存后前10分钟内结束订单不收取费用']
  // 当前代码块结束
  }
  // 按天模式文案（对照原型）
  return ['租金20元/天', '当天23:59超时后按每1天收费20元', '中途开门不限次数', '寄存后前2个小时内结束订单不收取费用']
// 当前代码块结束
})

// 调用 onMounted(async () => {，触发对应功能
onMounted(async () => {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理
  try {
    // 等待异步操作完成：拉取柜子列表与优惠券
    cabinets.value = await queryCabinets({ site_id: siteId })
    // 等待异步操作完成：统一券口径并空值防御
    coupons.value = normalizeCoupons(await myCoupons())
  // 捕获上述逻辑抛出的异常并进行错误处理
  } catch (e) { err.value = (e as Error).message }
// 执行 })，实现对应的前端逻辑

})

// 选择格口：记录 doorId 并立即试算（原型点选格口卡即带出预估费用）
// 定义函数/方法 pickSize，封装可复用的处理逻辑
function pickSize(id: number) {
  // 记录选中柜门
  doorId.value = id
  // 触发费用试算
  doQuote()
// 当前代码块结束
}

// 定义函数/方法，封装可复用的处理逻辑：async function doQuote() {
async function doQuote() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理
  try {
    // 等待异步操作完成：按选中尺寸+时长+券试算
    fee.value = await quoteFee({ size: selSize.value, hours: hours.value, coupon: couponId.value })
  // 捕获上述逻辑抛出的异常并进行错误处理
  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}

// @ticket 04-物品寄存：确认寄存（幂等键防重复提交），成功后进入支付确认页（原型流程）
// 定义函数/方法，封装可复用的处理逻辑：async function submitOrder() {
async function submitOrder() {
  // 条件判断：未选格口时提示
  if (!doorId.value) { err.value = '请先选择格口'; return }
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理
  try {
    // 声明变量 idem，幂等键
    const idem = `deposit-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    // 声明变量 r，落单结果
    const r = await createDeposit({
      // 定义对象属性 door_id，取值为 doorId.value, site_id: siteId, hours: hours.value, coupon: couponId.value
      door_id: doorId.value, site_id: siteId, hours: hours.value, coupon: couponId.value,
    // 执行 }, idem)，实现对应的前端逻辑
    }, idem)
    // 赋值更新：orderNo.value = r.order_no
    orderNo.value = r.order_no
    // 跳转支付确认页（对照原型 03-预支付提示后的确认流程）
    router.push({ path: '/pay-confirm', query: { order_no: r.order_no } })
  // 捕获上述逻辑抛出的异常并进行错误处理
  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}

// 更换格口（对照原型 02-5 弹窗"确认"）：调用换柜接口换到同尺寸新柜门
// 定义函数/方法，封装可复用的处理逻辑：async function confirmChange() {
async function confirmChange() {
  // 条件判断：暂无可换柜订单时提示
  if (!orderNo.value) { err.value = '暂无可换柜的订单'; showChange.value = false; return }
  // 找一个同尺寸的可用柜门
  const target = cabinets.value.find((c) => c.size === selSize.value && c.available > 0 && c.id !== doorId.value)
  // 条件判断：无可换目标
  if (!target) { err.value = '暂无可更换的格口'; showChange.value = false; return }
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理
  try {
    // 调用换柜接口
    const r = await changeCabinet(orderNo.value, target.id)
    // 更新选中柜门
    doorId.value = r.door_id
    // 关闭弹层并提示
    showChange.value = false
    // 赋值更新：msg.value
    msg.value = `已为您打开新的格口 ${r.door_id}`
  // 捕获上述逻辑抛出的异常并进行错误处理
  } catch (e) { err.value = (e as Error).message; showChange.value = false }
// 当前代码块结束

}
// 分转元展示
// 定义函数/方法 fen，封装可复用的处理逻辑
function fen(v: number): string { return (v / 100).toFixed(2) }
</script>

<template>
  <!-- 页面根容器：柜子信息页（对照原型） -->
  <div class="dep">
    <!-- 网点信息头部卡 -->
    <div class="site-card">
      <!-- 网点名称 -->
      <div class="site-name">{{ siteName }}</div>
      <!-- 网点地址（单行截断） -->
      <div class="site-addr">{{ siteAddr }}</div>
    <!-- 闭合头部卡 -->
    </div>

    <!-- 选择格口区 -->
    <div class="card sec">
      <!-- 区标题行：选择格口 + 收费规则入口 -->
      <div class="sec-head">
        <!-- 标题 -->
        <span class="sec-title">选择格口</span>
        <!-- 收费规则入口（点击切换计费模式演示按小时/按天两套规则） -->
        <span class="fee-rule" @click="mode = mode === 'hourly' ? 'daily' : 'hourly'">收费规则 ›</span>
      <!-- 闭合标题行 -->
      </div>
      <!-- 格口三卡横排 -->
      <div class="size-row">
        <!-- 循环 S/M/L 三档 -->
        <div v-for="(meta, s) in SIZES" :key="s" class="size-card" :class="{ on: selSize === s }" @click="pickSize(cabinets.find(c => c.size === s)?.id || 0)">
          <!-- 选中角标 -->
          <span v-if="selSize === s" class="check">✓</span>
          <!-- 格口图标（CSS 柜格造型，S/M/L 尺寸渐增） -->
          <div class="box-icon" :class="`b-${s}`"><i></i></div>
          <!-- 格口名称（含余量） -->
          <div class="size-name">{{ meta.label }}（余{{ sizeStock[s] ?? 0 }}）</div>
          <!-- 可放物品 -->
          <div class="size-desc">{{ meta.desc }}</div>
          <!-- 尺寸 -->
          <div class="size-dim">{{ meta.dim }}</div>
          <!-- 价格（按小时/按天两种展示，对照原型红字） -->
          <div class="size-price">{{ mode === 'hourly' ? meta.priceH : meta.priceD }}</div>
        <!-- 闭合格口卡 -->
        </div>
      <!-- 闭合格口横排 -->
      </div>
    <!-- 闭合选择格口区 -->
    </div>

    <!-- 选择寄存区（按小时模式；按天模式原型无此区） -->
    <div v-if="mode === 'hourly'" class="card sec">
      <!-- 区标题 + 免费寄存角标 -->
      <div class="sec-head">
        <!-- 标题 -->
        <span class="sec-title">选择寄存</span>
        <!-- 免费寄存角标 -->
        <span class="free-badge">😊 10分钟内免费寄存</span>
      <!-- 闭合标题行 -->
      </div>
      <!-- 寄存时长行 -->
      <div class="form-line">
        <!-- 标签 -->
        <span class="fl-label">寄存时长</span>
        <!-- 时长输入 -->
        <input v-model.number="hours" type="number" min="1" class="fl-input" @change="doQuote" />
        <!-- 单位 -->
        <span class="fl-unit">元/小时</span>
      <!-- 闭合时长行 -->
      </div>
      <!-- 优惠券行 -->
      <div class="form-line">
        <!-- 标签 -->
        <span class="fl-label">优惠券</span>
        <!-- 优惠券下拉（样式化为胶囊行） -->
        <select v-model.number="couponId" class="fl-select" @change="doQuote">
          <!-- 不使用选项 -->
          <option :value="undefined">不使用</option>
          <!-- 券选项（统一面值文案） -->
          <option v-for="c in coupons" :key="c.id" :value="c.id">{{ couponOptionText(c) }}</option>
        <!-- 闭合下拉 -->
        </select>
        <!-- 箭头 -->
        <span class="fl-arrow">›</span>
      <!-- 闭合优惠券行 -->
      </div>
      <!-- 结算提示 + 预估费用 -->
      <div class="est-row">
        <!-- 优惠结算提示红字 -->
        <span class="est-tip">优惠将在结束订单时结算</span>
        <!-- 预估费用红字 -->
        <span class="est-fee">预估费用：{{ estText }}</span>
      <!-- 闭合预估行 -->
      </div>
    <!-- 闭合选择寄存区 -->
    </div>

    <!-- 收费说明区 -->
    <div class="card sec">
      <!-- 区标题 -->
      <div class="sec-title solo">收费说明</div>
      <!-- 说明条目 -->
      <p v-for="(n, i) in feeNotes" :key="i" class="note">{{ i + 1 }}.{{ n }}</p>
    <!-- 闭合收费说明区 -->
    </div>

    <!-- 资金模式条：押金 / 预付（对照原型 03-预支付相关提示；按 fee.deposit 判定） -->
    <div class="fund-bar">
      <!-- 押金模式（deposit>0）：按时取件全退 | 超时自动扣费 -->
      <template v-if="fee && fee.deposit > 0">
        <!-- 押金条 -->
        <span class="fund-ico">💰</span>
        <!-- 押金文案 -->
        <span class="fund-label">寄存押金 ¥{{ fen(fee.deposit) }}</span>
        <!-- 规则说明 -->
        <span class="fund-rule">按时取件全退 | 超时自动扣费</span>
      <!-- 当前代码块结束 -->
      </template>
      <!-- 预付模式：按时计费 | 结束订单时多退少补 -->
      <template v-else>
        <!-- 预付图标 -->
        <span class="fund-ico">💰</span>
        <!-- 预付文案 -->
        <span class="fund-label">预付金额 ¥30</span>
        <!-- 规则说明 -->
        <span class="fund-rule">按时计费 | 结束订单时多退少补</span>
      <!-- 当前代码块结束 -->
      </template>
    <!-- 闭合资金模式条 -->
    </div>

    <!-- 更换格口底部弹层（对照原型 02-5） -->
    <div v-if="showChange" class="mask" @click.self="showChange = false">
      <!-- 弹层主体 -->
      <div class="sheet">
        <!-- 弹层标题 + 关闭 -->
        <div class="sheet-head">
          <!-- 标题 -->
          <span class="sheet-title">更换格口</span>
          <!-- 关闭按钮 -->
          <span class="sheet-close" @click="showChange = false">✕</span>
        <!-- 闭合标题行 -->
        </div>
        <!-- 柜门编号卡 -->
        <div class="sheet-door">
          <!-- 标签 -->
          <span class="sd-label">柜门编号</span>
          <!-- 编号大字 -->
          <span class="sd-no">{{ doorId || '—' }}号</span>
        <!-- 闭合编号卡 -->
        </div>
        <!-- 即将开门提示橙字 -->
        <p class="sheet-tip">即将为您打开新的格口</p>
        <!-- 底部操作行：取消订单链接 + 确认按钮 -->
        <div class="sheet-foot">
          <!-- 取消订单链接 -->
          <span class="cancel-link">不想用了，取消订单</span>
          <!-- 确认按钮 -->
          <button class="confirm-btn" @click="confirmChange">确认</button>
        <!-- 闭合操作行 -->
        </div>
      <!-- 闭合弹层主体 -->
      </div>
    <!-- 闭合遮罩 -->
    </div>

    <!-- 底部固定操作栏：实付/预支付 + 确认寄存 -->
    <div class="bottom-bar">
      <!-- 左侧金额区 -->
      <div class="pay-info">
        <!-- 金额标签（押金模式"实付" / 预付模式"预支付"） -->
        <span class="pay-label">{{ fee && fee.deposit > 0 ? '实付' : '预支付' }}：</span>
        <!-- 金额大字 -->
        <span class="pay-amount">¥{{ payText }}</span>
        <!-- 含押金/多退少补小字 -->
        <span class="pay-sub">{{ fee && fee.deposit > 0 ? `含押金¥${fen(fee.deposit)}` : '订单结束时多退少补' }}</span>
      <!-- 闭合金额区 -->
      </div>
      <!-- 确认寄存按钮 -->
      <button class="submit-btn" @click="submitOrder">确认寄存</button>
    <!-- 闭合底部栏 -->
    </div>

    <!-- 提示与错误 -->
    <p v-if="msg" class="ok">{{ msg }}</p>
    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合根容器 -->
  </div>
</template>

<style scoped>
/* 根容器：底部栏预留空间 */
.dep { padding-bottom: 76px; }
/* 网点头部卡 */
.site-card { background: #fff; border-radius: 12px; padding: 16px; margin-bottom: 12px; }
/* 网点名称 */
.site-name { font-size: 17px; font-weight: 700; color: #111827; }
/* 网点地址 */
.site-addr { font-size: 13px; color: #9ca3af; margin-top: 6px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
/* 通用区块卡 */
.card.sec { background: #fff; border-radius: 12px; padding: 14px; margin-bottom: 12px; }
/* 区标题行 */
.sec-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
/* 区标题 */
.sec-title { font-size: 15px; font-weight: 700; color: #111827; }
/* 独占标题（收费说明） */
.solo { margin-bottom: 8px; }
/* 收费规则入口灰字 */
.fee-rule { font-size: 13px; color: #6b7280; cursor: pointer; }
/* 格口三卡横排 */
.size-row { display: flex; gap: 10px; }
/* 格口卡 */
.size-card { flex: 1; position: relative; border: 1px solid #e5e7eb; border-radius: 10px; padding: 12px 8px; text-align: center; cursor: pointer; }
/* 选中态：蓝框 */
.size-card.on { border-color: #2563eb; background: #f5f9ff; }
/* 选中角标 */
.check { position: absolute; top: -6px; right: -6px; width: 20px; height: 20px; border-radius: 50%; background: #2563eb; color: #fff; font-size: 12px; line-height: 20px; }
/* 格口图标底座 */
.box-icon { height: 34px; margin-bottom: 6px; position: relative; }
/* 柜格造型：圆角矩形+格子线 */
.box-icon i { position: absolute; left: 50%; top: 6px; transform: translateX(-50%); width: 38px; height: 22px; background: linear-gradient(135deg, #93b8f5, #5f8ef0); border-radius: 5px; box-shadow: inset 0 -6px 0 rgba(255,255,255,.35); }
/* 大格口图标放大 */
.b-L i { width: 46px; height: 26px; }
/* 格口名称 */
.size-name { font-size: 13px; font-weight: 600; color: #111827; }
/* 可放物品 */
.size-desc { font-size: 11px; color: #6b7280; margin-top: 2px; }
/* 尺寸 */
.size-dim { font-size: 10px; color: #9ca3af; margin-top: 1px; }
/* 价格红字 */
.size-price { font-size: 13px; color: #ff4d4f; font-weight: 600; margin-top: 6px; }
/* 免费寄存角标 */
.free-badge { font-size: 12px; color: #ff7a45; background: #fff3ec; border-radius: 10px; padding: 2px 8px; }
/* 表单行 */
.form-line { display: flex; align-items: center; gap: 10px; padding: 8px 0; border-bottom: 1px solid #f3f4f6; }
/* 表单标签 */
.fl-label { width: 64px; font-size: 14px; color: #374151; }
/* 时长输入 */
.fl-input { flex: 1; border: none; background: #f5f7fa; border-radius: 8px; height: 38px; padding: 0 12px; font-size: 15px; text-align: left; }
/* 单位灰字 */
.fl-unit { font-size: 13px; color: #9ca3af; }
/* 优惠券下拉（胶囊行样式） */
.fl-select { flex: 1; border: none; background: #f5f7fa; border-radius: 8px; height: 38px; padding: 0 12px; font-size: 14px; color: #111827; }
/* 箭头 */
.fl-arrow { color: #c9ced6; font-size: 18px; }
/* 结算提示行 */
.est-row { display: flex; align-items: center; justify-content: space-between; margin-top: 10px; }
/* 优惠结算提示 */
.est-tip { font-size: 12px; color: #ff4d4f; }
/* 预估费用 */
.est-fee { font-size: 15px; color: #ff4d4f; font-weight: 700; }
/* 收费说明条目 */
.note { font-size: 13px; color: #374151; line-height: 2; margin: 0; }
/* 资金模式条：浅蓝底 */
.fund-bar { display: flex; align-items: center; gap: 8px; background: #eef4ff; border-radius: 10px; padding: 10px 14px; margin-bottom: 12px; }
/* 图标 */
.fund-ico { font-size: 16px; }
/* 资金标签蓝字 */
.fund-label { font-size: 13px; color: #2563eb; font-weight: 600; }
/* 规则说明 */
.fund-rule { font-size: 12px; color: #2563eb; margin-left: auto; }
/* 遮罩 */
.mask { position: fixed; inset: 0; background: rgba(0,0,0,.45); z-index: 50; display: flex; align-items: flex-end; justify-content: center; }
/* 底部弹层 */
.sheet { width: 100%; max-width: 430px; background: #fff; border-radius: 16px 16px 0 0; padding: 16px 16px 22px; }
/* 弹层标题行 */
.sheet-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
/* 弹层标题 */
.sheet-title { font-size: 16px; font-weight: 700; }
/* 关闭按钮 */
.sheet-close { color: #6b7280; font-size: 16px; cursor: pointer; }
/* 柜门编号卡 */
.sheet-door { background: #f5f7fa; border-radius: 10px; padding: 14px; display: flex; align-items: center; justify-content: space-between; }
/* 编号标签 */
.sd-label { font-size: 12px; color: #9ca3af; }
/* 编号大字 */
.sd-no { font-size: 30px; font-weight: 800; color: #111827; }
/* 即将开门提示 */
.sheet-tip { font-size: 13px; color: #ff7a45; margin: 12px 0 4px; }
/* 弹层底部操作行 */
.sheet-foot { display: flex; align-items: center; justify-content: space-between; margin-top: 10px; }
/* 取消订单链接 */
.cancel-link { color: #2563eb; font-size: 13px; text-decoration: underline; cursor: pointer; }
/* 确认按钮 */
.confirm-btn { background: #2563eb; color: #fff; border: none; border-radius: 20px; padding: 10px 34px; font-size: 15px; cursor: pointer; }
/* 底部固定操作栏 */
.bottom-bar { position: fixed; bottom: 0; left: 50%; transform: translateX(-50%); width: 100%; max-width: 430px; background: #fff; border-top: 1px solid #e5e7eb; display: flex; align-items: center; justify-content: space-between; padding: 10px 14px; z-index: 30; }
/* 金额区 */
.pay-info { display: flex; align-items: baseline; gap: 4px; }
/* 金额标签 */
.pay-label { font-size: 13px; color: #374151; }
/* 金额大字橙红 */
.pay-amount { font-size: 22px; font-weight: 800; color: #ff5a1f; }
/* 含押金/多退少补小字 */
.pay-sub { font-size: 11px; color: #9ca3af; }
/* 确认寄存按钮 */
.submit-btn { background: #2563eb; color: #fff; border: none; border-radius: 22px; padding: 12px 30px; font-size: 16px; font-weight: 600; cursor: pointer; }
/* 成功提示 */
.ok { color: #16a34a; font-size: 13px; }
/* 错误提示 */
.err { color: #ff4d4f; font-size: 13px; }
</style>
