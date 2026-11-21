<script setup lang="ts">
// @ticket 04-物品寄存 + 05-订单查看与取出：支付完成/柜门已开提示
// @proto 原型对齐：04-一期用户界面-10/16/17/18、06-用户界面规划-15/16/20 支付完成与柜门已开（柜号大字/网格图/查看订单/继续寄存）
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import CabinetMap from '@/components/CabinetMap.vue'
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 route，用于保存数据或函数引用

const route = useRoute()
// 声明变量 doorNo，用于保存数据或函数引用

const doorNo = ref<number>(Number(route.query.door) || 18)
// 声明变量 countdown，用于保存数据或函数引用

const countdown = ref(30)
// 声明变量 doorOpened，用于保存数据或函数引用

const doorOpened = ref(true)
// 声明变量 timeText，用于保存数据或函数引用

const timeText = computed<string>(() => {
  // 声明变量 m，用于保存数据或函数引用

  const m = String(Math.floor(countdown.value / 60)).padStart(2, '0')
  // 声明变量 s，用于保存数据或函数引用

  const s = String(countdown.value % 60).padStart(2, '0')
  // 返回 m + ':' + s 给调用方

  return m + ':' + s
// 当前代码块结束

})
// 定义函数/方法，封装可复用的处理逻辑：function viewOrders() {

function viewOrders() {
  // 执行 router.push('/orders')，实现对应的前端逻辑

  router.push('/orders')
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function again() {

function again() {
  // 执行 router.push('/deposit')，实现对应的前端逻辑

  router.push('/deposit')
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 success-icon -->

    <div class="success-icon">✓</div>
    <!-- div 元素，class 为 success-t -->

    <div class="success-t">支付成功</div>
    <!-- div 元素，class 为 open-box -->

    <div v-if="doorOpened" class="open-box">
      <!-- div 元素，class 为 open-title -->

      <div class="open-title">柜门已开启，请放入物品</div>
      <!-- div 元素，class 为 door-no -->

      <div class="door-no">{{ doorNo }}<span>号柜门</span></div>
      <!-- div 元素，class 为 cd -->

      <div class="cd">柜门保持开启 {{ timeText }}</div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- CabinetMap 元素 -->

    <CabinetMap :active-door="doorNo" />
    <!-- div 元素，class 为 btns -->

    <div class="btns">
      <!-- button 元素，class 为 ghost -->

      <button class="ghost" @click="viewOrders">查看订单</button>
      <!-- button 元素，class 为 primary -->

      <button class="primary" @click="again">继续寄存</button>
    <!-- 闭合 div 元素 -->

    </div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 12px; align-items: stretch; }
/* .success-icon 样式定义 */

.success-icon { width: 56px; height: 56px; margin: 6px auto 0; border-radius: 50%; background: #22c55e; color: #fff; font-size: 32px; display: flex; align-items: center; justify-content: center; }
/* .success-t 样式定义 */

.success-t { text-align: center; font-size: 18px; font-weight: 700; color: #1f2937; }
/* .open-box 样式定义 */

.open-box { background: #eff6ff; border-radius: 10px; padding: 16px; text-align: center; }
/* .open-title 样式定义 */

.open-title { font-size: 14px; color: #2563eb; }
/* .door-no 样式定义 */

.door-no { font-size: 48px; font-weight: 800; color: #1d4ed8; line-height: 1.2; }
/* .door-no span 样式定义 */

.door-no span { font-size: 16px; margin-left: 4px; }
/* .cd 样式定义 */

.cd { font-size: 12px; color: #6b7280; }
/* .btns 样式定义 */

.btns { display: flex; gap: 10px; }
/* .btns button 样式定义 */

.btns button { flex: 1; height: 42px; border-radius: 21px; font-size: 15px; }
/* .ghost 样式定义 */

.ghost { border: 1px solid #2563eb; background: #fff; color: #2563eb; }
/* .primary 样式定义 */

.primary { border: 0; background: #2563eb; color: #fff; }
</style>

