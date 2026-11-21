<script setup lang="ts">
// @proto 原型对齐：04-一期用户界面-13/17、06-用户界面规划-8/15 柜门位置网格图
// @ticket 05-订单查看与取出：订单详情柜门平面图/位置指示
import { computed } from 'vue'
// 声明变量 props，用于保存数据或函数引用

const props = defineProps<{ activeDoor?: number; cols?: number; rows?: number }>()
// 声明变量 cols，用于保存数据或函数引用

const cols = computed<number>(() => props.cols || 6)
// 声明变量 rows，用于保存数据或函数引用

const rows = computed<number>(() => props.rows || 4)
// 声明变量 doors，用于保存数据或函数引用

const doors = computed<number[]>(() => {
  // 声明变量 total，用于保存数据或函数引用

  const total = cols.value * rows.value
  // 声明变量 list:，用于保存数据或函数引用

  const list: number[] = []
  // 执行 for (let i = 1; i <= total; i++) {，实现对应的前端逻辑

  for (let i = 1; i <= total; i++) {
    // 执行 list.push(i)，实现对应的前端逻辑

    list.push(i)
  // 当前代码块结束

  }
  // 返回 list 给调用方

  return list
// 当前代码块结束

})
// 定义函数/方法，封装可复用的处理逻辑：function isActive(no: number): boolean {

function isActive(no: number): boolean {
  // 返回 no === props.activeDoor 给调用方

  return no === props.activeDoor
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 door-map -->

  <div class="door-map">
    <!-- div 元素，class 为 map-title -->

    <div class="map-title">柜门位置示意图</div>
    <!-- div 元素，class 为 grid -->

    <div class="grid" :style="{ gridTemplateColumns: 'repeat(' + cols + ', 1fr)' }">
      <!-- div 元素，class 为 door -->

      <div v-for="no in doors" :key="no" class="door" :class="{ active: isActive(no) }">{{ no }}</div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- div 元素，class 为 legend -->

    <div v-if="activeDoor" class="legend">当前柜门：{{ activeDoor }} 号</div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .door-map 样式定义 */

.door-map { background: #fff; border-radius: 10px; padding: 12px; }
/* .map-title 样式定义 */

.map-title { font-size: 14px; color: #1f2937; font-weight: 600; margin-bottom: 10px; }
/* .grid 样式定义 */

.grid { display: grid; gap: 6px; }
/* .door 样式定义 */

.door { aspect-ratio: 1; display: flex; align-items: center; justify-content: center; background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 6px; font-size: 12px; color: #2563eb; }
/* .door.active 样式定义 */

.door.active { background: #2563eb; color: #fff; font-weight: 700; box-shadow: 0 0 0 2px #93c5fd; }
/* .legend 样式定义 */

.legend { margin-top: 10px; font-size: 13px; color: #2563eb; }
</style>

