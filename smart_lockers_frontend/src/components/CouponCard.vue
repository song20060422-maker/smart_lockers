<script setup lang="ts">
// @proto 原型对齐：01-C端优惠券-1~5 我的优惠券卡片（面额/失效时间/去使用/使用规则展开/限网点/不可用态）
// @ticket 06-个人中心：优惠卡券
import { computed, ref } from 'vue'
// 券卡片入参：兼容后端归一化字段 + 原型展示字段（siteLimit 限网点文案、usable 当前网点可用性）
const props = defineProps<{
  // 券唯一标识
  id?: number
  // 券码
  code?: string
  // 券标题（如"1元优惠券"）
  title?: string
  // 抵扣金额（元）
  amount?: number
  // 折扣（如 8 表示 8 折）
  discount?: number
  // 时长券（小时）
  hours?: number
  // 失效时间
  expire_at?: string
  // 适用范围/限网点文案
  scope?: string
  // 是否已使用/已失效（灰显）
  expired?: boolean
  // 当前网点是否可用（false 显示"当前网点不可用"+不可用按钮）
  usable?: boolean
}>()
// 规则展开状态（原型：使用规则∨/∧ 点击展开规则说明）
const expanded = ref(false)

// 左侧大面额文案：金额券→"1元"，时长券→"1小时"，折扣券→"8折"
const bigText = computed<string>(() => {
  // 金额券优先
  if (props.amount != null) return String(props.amount)
  // 折扣券
  if (props.discount != null) return String(props.discount)
  // 时长券
  if (props.hours != null) return String(props.hours)
  // 兜底
  return '券'
// 当前代码块结束
})
// 面额单位：元 / 折 / 小时
const unitText = computed<string>(() => {
  // 金额券单位"元"
  if (props.amount != null) return '元'
  // 折扣券单位"折"
  if (props.discount != null) return '折'
  // 时长券单位"小时"
  if (props.hours != null) return '小时'
  // 兜底空
  return ''
// 当前代码块结束
})
// 券标题（右侧首行）
const titleText = computed<string>(() => {
  // 有标题直接用
  if (props.title) return props.title
  // 金额券兜底标题
  if (props.amount != null) return `${props.amount}元优惠券`
  // 折扣券兜底标题
  if (props.discount != null) return `${props.discount}折优惠券`
  // 时长券兜底标题
  if (props.hours != null) return `${props.hours}小时免费寄存券`
  // 兜底
  return '寄存优惠券'
// 当前代码块结束
})
// 失效时间文案（原型红字："23:23:23后失效"）
const expireText = computed<string>(() => {
  // 有失效时间则展示"xx后失效"
  if (props.expire_at) return `${props.expire_at}后失效`
  // 长期有效兜底
  return '长期有效'
// 当前代码块结束
})
// 限网点底部文案（原型："中原区二七万达广场兑换" 或 "限【xx】网点使用"）
const scopeText = computed<string>(() => props.scope || '全场通用')
// 不可用判定：已失效 或 当前网点不可用
const disabled = computed<boolean>(() => !!props.expired || props.usable === false)
// 不可用原因红字（原型："当前网点不可用"/"本网点不可用"）
const disableTip = computed<string>(() => props.expired ? '' : (props.usable === false ? '本网点不可用' : ''))
</script>

<template>
  <!-- 券卡片根容器：expired 灰显 -->
  <div class="coupon" :class="{ gray: expired }">
    <!-- 卡片主体：左侧面额区 + 右侧信息区 -->
    <div class="main">
      <!-- 左侧面额区 -->
      <div class="left">
        <!-- 大面额数字 -->
        <span class="big">{{ bigText }}</span>
        <!-- 单位 -->
        <span class="unit">{{ unitText }}</span>
      <!-- 闭合左区 -->
      </div>
      <!-- 右侧信息区 -->
      <div class="right">
        <!-- 首行：券标题 + 去使用/不可用按钮 -->
        <div class="row head">
          <!-- 券标题 -->
          <span class="title">{{ titleText }}</span>
          <!-- 去使用按钮（可用态） -->
          <button v-if="!disabled" class="use-btn">去使用</button>
          <!-- 不可用按钮（灰态） -->
          <span v-else class="use-btn gray-btn">不可用</span>
        <!-- 闭合首行 -->
        </div>
        <!-- 第二行：失效时间红字 + 不可用原因 -->
        <div class="row expire-row">
          <!-- 失效时间 -->
          <span class="expire">{{ expireText }}</span>
          <!-- 当前网点不可用提示 -->
          <span v-if="disableTip" class="tip">{{ disableTip }}</span>
        <!-- 闭合第二行 -->
        </div>
        <!-- 第三行：使用规则 ∨/∧（点击展开收起） -->
        <div class="row rule-row" @click="expanded = !expanded">
          <!-- 规则入口 -->
          <span class="rule-entry">使用规则 <b>{{ expanded ? '∧' : '∨' }}</b></span>
        <!-- 闭合第三行 -->
        </div>
        <!-- 规则展开内容 -->
        <div v-if="expanded" class="rule-body">
          <!-- 规则说明文案（对照原型） -->
          规则说明：优惠券仅限本人使用，一旦使用无法恢复。下单时选择优惠券，结束订单时将按优惠券结算。
        <!-- 闭合规则内容 -->
        </div>
      <!-- 闭合右区 -->
      </div>
    <!-- 闭合卡片主体 -->
    </div>
    <!-- 底部限网点灰条 -->
    <div class="foot">{{ scopeText }}</div>
  <!-- 闭合券卡片 -->
  </div>
</template>

<style scoped>
/* 券卡片：白底圆角，底部灰条通过负边距贴合 */
.coupon { background: #fff; border-radius: 10px; overflow: hidden; margin-bottom: 12px; border: 1px solid #e5e7eb; }
/* 灰显态 */
.coupon.gray { filter: grayscale(1); opacity: .75; }
/* 主体：左右两栏 */
.main { display: flex; padding: 14px 12px 10px; }
/* 左侧面额区：大字红橙色调（原型 1元/8折 为橙红色） */
.left { width: 86px; flex-shrink: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; border-right: 1px dashed #e5e7eb; padding-right: 10px; }
/* 大面额数字 */
.big { font-size: 34px; font-weight: 800; color: #ff5a1f; line-height: 1; }
/* 面额单位 */
.unit { font-size: 12px; color: #ff5a1f; margin-top: 4px; }
/* 右侧信息区 */
.right { flex: 1; padding-left: 12px; display: flex; flex-direction: column; gap: 6px; }
/* 行容器 */
.row { display: flex; align-items: center; justify-content: space-between; }
/* 券标题 */
.title { font-size: 15px; font-weight: 600; color: #111827; }
/* 去使用按钮：蓝色胶囊 */
.use-btn { background: #2563eb; color: #fff; border: none; border-radius: 14px; padding: 4px 14px; font-size: 12px; cursor: pointer; }
/* 不可用按钮：灰色胶囊 */
.gray-btn { background: #c9ced6; cursor: default; }
/* 失效时间红字 */
.expire { color: #ff4d4f; font-size: 12px; }
/* 不可用原因红字 */
.tip { color: #ff4d4f; font-size: 12px; }
/* 规则入口灰字 */
.rule-entry { color: #6b7280; font-size: 12px; cursor: pointer; }
/* 规则展开内容 */
.rule-body { font-size: 12px; color: #9ca3af; line-height: 1.7; background: #f9fafb; border-radius: 6px; padding: 8px 10px; }
/* 底部限网点灰条 */
.foot { background: #f3f4f6; color: #6b7280; font-size: 12px; padding: 6px 14px; }
</style>
