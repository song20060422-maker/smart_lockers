<script setup lang="ts">
// @ticket 04-物品寄存 + 02-支付宝支付：费用支付确认（须知/押金/总额/确认支付）
// @proto 原型对齐：03-预支付相关提示-1~4、04-一期用户界面-9 费用支付确认弹层
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import FeeRuleSheet from '@/components/FeeRuleSheet.vue'
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 route，用于保存数据或函数引用

const route = useRoute()
// 声明变量 showRule，用于保存数据或函数引用

const showRule = ref(false)
// 声明变量 agreed，用于保存数据或函数引用

const agreed = ref(false)
// 声明变量 hours，用于保存数据或函数引用

const hours = ref<number>(Number(route.query.hours) || 2)
// 声明变量 rate，用于保存数据或函数引用

const rate = ref<number>(Number(route.query.rate) || 2)
// 声明变量 rent，用于保存数据或函数引用

const rent = computed<number>(() => hours.value * rate.value)
// 声明变量 deposit，用于保存数据或函数引用

const deposit = ref<number>(99)
// 声明变量 total，用于保存数据或函数引用

const total = computed<number>(() => rent.value + deposit.value)
// 定义函数/方法，封装可复用的处理逻辑：function confirmPay() {

function confirmPay() {
  // 条件判断：if (!agreed.value) return 时执行下面逻辑

  if (!agreed.value) return
  // 执行 router.push({ path: '/pay-success', query: { hours: String(hours.value) } })，实现对应的前端逻辑

  router.push({ path: '/pay-success', query: { hours: String(hours.value) } })
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 card -->

    <div class="card">
      <!-- div 元素，class 为 h -->

      <div class="h">支付须知</div>
      <!-- ul 元素 -->

      <ul>
        <!-- li 元素 -->

        <li>租柜费按柜型与时长计算，押金结单后原路退还</li>
        <!-- li 元素 -->

        <li>请在营业结束前取出物品，超时按规则续费</li>
        <!-- li 元素 -->

        <li>支付成功即开启对应柜门</li>
      <!-- 闭合 ul 元素 -->

      </ul>
      <!-- button 元素，class 为 link -->

      <button class="link" @click="showRule = true">查看收费规则 ›</button>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- div 元素，class 为 card -->

    <div class="card">
      <!-- div 元素，class 为 row -->

      <div class="row"><span>寄存时长</span><span>{{ hours }} 小时</span></div>
      <!-- div 元素，class 为 row -->

      <div class="row"><span>租柜费</span><span>¥{{ rent }}</span></div>
      <!-- div 元素，class 为 row -->

      <div class="row"><span>押金</span><span>¥{{ deposit }}</span></div>
      <!-- div 元素，class 为 row total -->

      <div class="row total"><span>订单总额</span><span>¥{{ total }}</span></div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- label 元素，class 为 agree -->

    <label class="agree"><input type="checkbox" v-model="agreed" /> 我已阅读并同意《寄存服务协议》</label>
    <!-- button 元素，class 为 primary -->

    <button class="primary" :disabled="!agreed" @click="confirmPay">确认支付 ¥{{ total }}</button>
    <!-- FeeRuleSheet 元素 -->

    <FeeRuleSheet :visible="showRule" @close="showRule = false" />
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 12px; }
/* .card 样式定义 */

.card { background: #fff; border-radius: 10px; padding: 14px; border: 1px solid #e5e7eb; }
/* .h 样式定义 */

.h { font-size: 15px; font-weight: 700; color: #1f2937; margin-bottom: 8px; }
/* ul 样式定义 */

ul { margin: 0; padding-left: 18px; color: #6b7280; font-size: 13px; line-height: 1.9; }
/* .link 样式定义 */

.link { border: 0; background: none; color: #2563eb; font-size: 13px; padding: 8px 0 0; }
/* .row 样式定义 */

.row { display: flex; justify-content: space-between; font-size: 14px; color: #4b5563; padding: 6px 0; }
/* .total 样式定义 */

.total { border-top: 1px solid #f3f4f6; margin-top: 6px; font-size: 16px; font-weight: 700; color: #dc2626; }
/* .agree 样式定义 */

.agree { font-size: 12px; color: #6b7280; }
/* .primary 样式定义 */

.primary { height: 44px; border: 0; border-radius: 22px; background: #2563eb; color: #fff; font-size: 16px; }
/* .primary:disabled 样式定义 */

.primary:disabled { background: #93c5fd; }
</style>

