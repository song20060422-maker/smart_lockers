<script setup lang="ts">
// @ticket 06-个人中心：在线客服/使用帮助
// @proto 原型对齐：06-用户界面规划-1/19 在线客服（常见问题展开/联系客服）
import { ref } from 'vue'
// 声明变量 faqs，用于保存数据或函数引用

const faqs = ref<string[]>([
  // 执行 '如何寄存物品？',，实现对应的前端逻辑

  '如何寄存物品？',
  // 执行 '忘记取件码怎么办？',，实现对应的前端逻辑

  '忘记取件码怎么办？',
  // 执行 '超时未取如何收费？',，实现对应的前端逻辑

  '超时未取如何收费？',
  // 执行 '押金什么时候退还？',，实现对应的前端逻辑

  '押金什么时候退还？',
  // 执行 '柜门打不开怎么办？',，实现对应的前端逻辑

  '柜门打不开怎么办？',
// 执行 ])，实现对应的前端逻辑

])
// 声明变量 openIdx，用于保存数据或函数引用

const openIdx = ref<number>(-1)
// 声明变量 showContact，用于保存数据或函数引用

const showContact = ref(false)
// 声明变量 answers:，用于保存数据或函数引用

const answers: Record<number, string> = {
  // 执行 0: '在首页查询附近网点，到柜机扫码选择空闲柜门，支付后放入物品即可。',，实现对应的前端逻辑

  0: '在首页查询附近网点，到柜机扫码选择空闲柜门，支付后放入物品即可。',
  // 执行 1: '可在订单详情页查看取件码，或使用人脸识别/扫码开柜。',，实现对应的前端逻辑

  1: '可在订单详情页查看取件码，或使用人脸识别/扫码开柜。',
  // 执行 2: '超时按柜型时长规则续费，可在订单页补交（享限时优惠）。',，实现对应的前端逻辑

  2: '超时按柜型时长规则续费，可在订单页补交（享限时优惠）。',
  // 执行 3: '结单后押金将原路退还，通常 1-3 个工作日到账。',，实现对应的前端逻辑

  3: '结单后押金将原路退还，通常 1-3 个工作日到账。',
  // 执行 4: '请在开门异常页点击重新开门，或联系客服协助处理。',，实现对应的前端逻辑

  4: '请在开门异常页点击重新开门，或联系客服协助处理。',
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function toggle(i: number) {

function toggle(i: number) {
  // 赋值更新：openIdx.value = openIdx.value === i ? -1 : i

  openIdx.value = openIdx.value === i ? -1 : i
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 head -->

    <div class="head">常见问题</div>
    <!-- div 元素，class 为 faq -->

    <div v-for="(f, i) in faqs" :key="i" class="faq">
      <!-- div 元素，class 为 q -->

      <div class="q" @click="toggle(i)">{{ f }}<span>{{ openIdx === i ? '−' : '+' }}</span></div>
      <!-- div 元素，class 为 a -->

      <div v-if="openIdx === i" class="a">{{ answers[i] }}</div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- button 元素，class 为 primary -->

    <button class="primary" @click="showContact = true">联系在线客服</button>
    <!-- div 元素，class 为 contact -->

    <div v-if="showContact" class="contact">客服电话：400-XXX-XXXX<br />服务时间：08:00 - 22:00</div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 10px; }
/* .head 样式定义 */

.head { font-size: 16px; font-weight: 700; color: #1f2937; }
/* .faq 样式定义 */

.faq { background: #fff; border-radius: 8px; border: 1px solid #e5e7eb; overflow: hidden; }
/* .q 样式定义 */

.q { display: flex; justify-content: space-between; padding: 12px; font-size: 14px; color: #1f2937; }
/* .q span 样式定义 */

.q span { color: #2563eb; }
/* .a 样式定义 */

.a { padding: 0 12px 12px; font-size: 13px; color: #6b7280; line-height: 1.6; }
/* .primary 样式定义 */

.primary { height: 42px; border: 0; border-radius: 21px; background: #2563eb; color: #fff; font-size: 15px; }
/* .contact 样式定义 */

.contact { background: #eff6ff; color: #2563eb; font-size: 13px; padding: 12px; border-radius: 8px; line-height: 1.8; text-align: center; }
</style>

