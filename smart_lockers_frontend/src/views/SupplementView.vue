<script setup lang="ts">
// @ticket 05-订单查看与取出：超时补交（补交并取出/8折优惠/支付后自动开柜）
// @proto 原型对齐：04-一期用户界面-15/20、06-用户界面规划-9 超时补交与优惠
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 hours，用于保存数据或函数引用

const hours = ref<number>(2)
// 声明变量 rate，用于保存数据或函数引用

const rate = ref<number>(2)
// 声明变量 raw，用于保存数据或函数引用

const raw = computed<number>(() => hours.value * rate.value)
// 声明变量 payable，用于保存数据或函数引用

const payable = computed<number>(() => Math.round(raw.value * 0.8 * 100) / 100)
// 定义函数/方法，封装可复用的处理逻辑：function payAndOpen() {

function payAndOpen() {
  // 执行 router.push({ path: '/pay-success', query: { door: '18', scene: 'pickup' } })，实现对应的前端逻辑

  router.push({ path: '/pay-success', query: { door: '18', scene: 'pickup' } })
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 notice -->

    <div class="notice">您的订单已超时，请补交超时费用后取出物品。</div>
    <!-- div 元素，class 为 card -->

    <div class="card">
      <!-- div 元素，class 为 row -->

      <div class="row"><span>超时时长</span><span>{{ hours }} 小时</span></div>
      <!-- div 元素，class 为 row -->

      <div class="row"><span>应付费用</span><span>¥{{ raw }}</span></div>
      <!-- div 元素，class 为 row discount -->

      <div class="row discount"><span>限时 8 折优惠</span><span>-¥{{ (raw - payable).toFixed(2) }}</span></div>
      <!-- div 元素，class 为 row total -->

      <div class="row total"><span>实付金额</span><span>¥{{ payable }}</span></div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- label 元素，class 为 agree -->

    <label class="agree">支付完成后柜门将自动开启，请及时取出物品</label>
    <!-- button 元素，class 为 primary -->

    <button class="primary" @click="payAndOpen">补交 ¥{{ payable }} 并取出</button>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 12px; }
/* .notice 样式定义 */

.notice { background: #fef3c7; color: #92400e; font-size: 13px; padding: 12px; border-radius: 8px; line-height: 1.6; }
/* .card 样式定义 */

.card { background: #fff; border-radius: 10px; padding: 14px; border: 1px solid #e5e7eb; }
/* .row 样式定义 */

.row { display: flex; justify-content: space-between; font-size: 14px; color: #4b5563; padding: 6px 0; }
/* .discount 样式定义 */

.discount { color: #16a34a; }
/* .total 样式定义 */

.total { border-top: 1px solid #f3f4f6; margin-top: 6px; font-size: 16px; font-weight: 700; color: #dc2626; }
/* .agree 样式定义 */

.agree { font-size: 12px; color: #9ca3af; }
/* .primary 样式定义 */

.primary { height: 44px; border: 0; border-radius: 22px; background: #2563eb; color: #fff; font-size: 16px; }
</style>

