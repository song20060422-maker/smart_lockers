<script setup lang="ts">
// @proto 原型对齐：01-C端优惠券-5、06-13/14 优惠券卡片（可使用/已过期）
// @ticket 06-个人中心：优惠卡券（抵扣金额/时长/有效期/适用范围/使用规则）
import { computed } from 'vue'
// 声明变量 props，用于保存数据或函数引用

const props = defineProps<{ coupon: { id?: number; code?: string; title?: string; amount?: number; discount?: number; hours?: number; expire_at?: string; scope?: string; expired?: boolean } }>()
// 声明变量 amountText，用于保存数据或函数引用

const amountText = computed<string>(() => {
  // 条件判断：if (props.coupon.amount != null) { 时执行下面逻辑

  if (props.coupon.amount != null) {
    // 返回 '¥' + props.coupon.amount 给调用方

    return '¥' + props.coupon.amount
  // 当前代码块结束

  }
  // 条件判断：if (props.coupon.discount != null) { 时执行下面逻辑

  if (props.coupon.discount != null) {
    // 返回 props.coupon.discount + '折' 给调用方

    return props.coupon.discount + '折'
  // 当前代码块结束

  }
  // 条件判断：if (props.coupon.hours != null) { 时执行下面逻辑

  if (props.coupon.hours != null) {
    // 返回 props.coupon.hours + '小时' 给调用方

    return props.coupon.hours + '小时'
  // 当前代码块结束

  }
  // 返回 '券' 给调用方

  return '券'
// 当前代码块结束

})
</script>

<template>
  <!-- div 元素，class 为 coupon -->

  <div class="coupon" :class="{ gray: coupon.expired }">
    <!-- div 元素，class 为 left -->

    <div class="left"><span class="amt">{{ amountText }}</span><span class="unit">{{ coupon.title || '寄存优惠券' }}</span></div>
    <!-- div 元素，class 为 right -->

    <div class="right">
      <!-- div 元素，class 为 row -->

      <div class="row">券码：{{ coupon.code || '—' }}</div>
      <!-- div 元素，class 为 row -->

      <div class="row">有效期至：{{ coupon.expire_at || '长期有效' }}</div>
      <!-- div 元素，class 为 row -->

      <div class="row">适用范围：{{ coupon.scope || '全场寄存柜' }}</div>
      <!-- div 元素，class 为 row rule -->

      <div class="row rule">使用规则：到场出示券码，限寄存服务使用</div>
    <!-- 闭合 div 元素 -->

    </div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .coupon 样式定义 */

.coupon { display: flex; background: #fff; border-radius: 10px; overflow: hidden; margin-bottom: 10px; border: 1px solid #dbeafe; }
/* .coupon.gray 样式定义 */

.coupon.gray { filter: grayscale(1); opacity: .7; }
/* .left 样式定义 */

.left { width: 96px; flex-shrink: 0; background: #eff6ff; display: flex; flex-direction: column; align-items: center; justify-content: center; color: #2563eb; }
/* .amt 样式定义 */

.amt { font-size: 22px; font-weight: 800; }
/* .unit 样式定义 */

.unit { font-size: 11px; margin-top: 2px; }
/* .right 样式定义 */

.right { flex: 1; padding: 10px 12px; }
/* .row 样式定义 */

.row { font-size: 12px; color: #4b5563; line-height: 1.8; }
/* .rule 样式定义 */

.rule { color: #9ca3af; }
</style>

