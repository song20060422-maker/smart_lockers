<script setup lang="ts">
// @ticket 05-订单查看与取出：选择取出类型（结束订单全部取出/临时开柜继续寄存）
// @proto 原型对齐：04-一期用户界面-14 选择取出类型（结束订单/临时开门继续寄存）
import { ref } from 'vue'
import { useRouter } from 'vue-router'
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 selected，用于保存数据或函数引用

const selected = ref<'finish' | 'temporary'>('finish')
// 定义函数/方法，封装可复用的处理逻辑：function confirm() {

function confirm() {
  // 条件判断：if (selected.value === 'finish') { 时执行下面逻辑

  if (selected.value === 'finish') {
    // 执行 router.push({ path: '/door-exception', query: { scene: 'opened' } })，实现对应的前端逻辑

    router.push({ path: '/door-exception', query: { scene: 'opened' } })
  // 否则再次判断/执行对应逻辑

  } else {
    // 执行 router.push({ path: '/pay-success', query: { door: '18', scene: 'temporary' } })，实现对应的前端逻辑

    router.push({ path: '/pay-success', query: { door: '18', scene: 'temporary' } })
  // 当前代码块结束

  }
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 title -->

    <div class="title">请选择取出类型</div>
    <!-- div 元素，class 为 opt -->

    <div class="opt" :class="{ on: selected === 'finish' }" @click="selected = 'finish'">
      <!-- div 元素，class 为 radio -->

      <div class="radio"></div>
      <!-- div 元素，class 为 txt -->

      <div class="txt"><div class="h">结束订单，全部取出</div><div class="d">取出全部物品并结束寄存，押金原路退还</div></div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- div 元素，class 为 opt -->

    <div class="opt" :class="{ on: selected === 'temporary' }" @click="selected = 'temporary'">
      <!-- div 元素，class 为 radio -->

      <div class="radio"></div>
      <!-- div 元素，class 为 txt -->

      <div class="txt"><div class="h">临时开柜，继续寄存</div><div class="d">仅临时取物，订单继续计费，稍后自动关闭</div></div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- button 元素，class 为 primary -->

    <button class="primary" @click="confirm">确认开柜</button>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 14px; }
/* .title 样式定义 */

.title { font-size: 17px; font-weight: 700; color: #1f2937; text-align: center; }
/* .opt 样式定义 */

.opt { display: flex; align-items: center; gap: 12px; background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; padding: 16px; }
/* .opt.on 样式定义 */

.opt.on { border-color: #2563eb; background: #eff6ff; }
/* .radio 样式定义 */

.radio { width: 20px; height: 20px; border-radius: 50%; border: 2px solid #cbd5e1; flex-shrink: 0; }
/* .opt.on .radio 样式定义 */

.opt.on .radio { border-color: #2563eb; background: #2563eb; box-shadow: inset 0 0 0 4px #fff; }
/* .h 样式定义 */

.h { font-size: 15px; color: #1f2937; font-weight: 600; }
/* .d 样式定义 */

.d { font-size: 12px; color: #9ca3af; margin-top: 4px; line-height: 1.5; }
/* .primary 样式定义 */

.primary { height: 44px; border: 0; border-radius: 22px; background: #2563eb; color: #fff; font-size: 16px; margin-top: 6px; }
</style>

