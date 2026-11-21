<script setup lang="ts">
// @proto 原型对齐：04-19/20/21、06-6/9/10/11 订单多维状态
// @ticket 05-订单查看与取出：订单状态（进行中/已超时/待补缴/已完成/已退费/异常）
import { computed } from 'vue'
// 声明变量 props，用于保存数据或函数引用

const props = defineProps<{ status: string }>()
// 声明变量 label，用于保存数据或函数引用

const label = computed<string>(() => {
  // 声明变量 map:，用于保存数据或函数引用

  const map: Record<string, string> = {
    // 执行 PENDING_PAY: '待支付', ONGOING: '进行中', OVERTIME: '已超时',，实现对应的前端逻辑

    PENDING_PAY: '待支付', ONGOING: '进行中', OVERTIME: '已超时',
    // 执行 NEED_SUPPLEMENT: '待补缴', COMPLETED: '已完成', REFUNDED: '已退费',，实现对应的前端逻辑

    NEED_SUPPLEMENT: '待补缴', COMPLETED: '已完成', REFUNDED: '已退费',
    // 执行 SUPPLEMENTED: '已补交', ABNORMAL: '异常', CANCELLED: '已取消',，实现对应的前端逻辑

    SUPPLEMENTED: '已补交', ABNORMAL: '异常', CANCELLED: '已取消',
  // 当前代码块结束

  }
  // 返回 map[props.status] || props.status || '未知' 给调用方

  return map[props.status] || props.status || '未知'
// 当前代码块结束

})
// 声明变量 cls，用于保存数据或函数引用

const cls = computed<string>(() => {
  // 声明变量 map:，用于保存数据或函数引用

  const map: Record<string, string> = {
    // 执行 ONGOING: 'tag-blue', PENDING_PAY: 'tag-orange', OVERTIME: 'tag-red',，实现对应的前端逻辑

    ONGOING: 'tag-blue', PENDING_PAY: 'tag-orange', OVERTIME: 'tag-red',
    // 执行 NEED_SUPPLEMENT: 'tag-red', COMPLETED: 'tag-green', REFUNDED: 'tag-gray',，实现对应的前端逻辑

    NEED_SUPPLEMENT: 'tag-red', COMPLETED: 'tag-green', REFUNDED: 'tag-gray',
    // 执行 SUPPLEMENTED: 'tag-green', ABNORMAL: 'tag-red', CANCELLED: 'tag-gray',，实现对应的前端逻辑

    SUPPLEMENTED: 'tag-green', ABNORMAL: 'tag-red', CANCELLED: 'tag-gray',
  // 当前代码块结束

  }
  // 返回 map[props.status] || 'tag-gray' 给调用方

  return map[props.status] || 'tag-gray'
// 当前代码块结束

})
</script>

<template>
  <!-- span 元素，class 为 tag -->

  <span class="tag" :class="cls">{{ label }}</span>
</template>

<style scoped>
/* .tag 样式定义 */

.tag { display: inline-block; padding: 2px 8px; border-radius: 10px; font-size: 12px; }
/* .tag-blue 样式定义 */

.tag-blue { background: #dbeafe; color: #1d4ed8; }
/* .tag-orange 样式定义 */

.tag-orange { background: #fef3c7; color: #b45309; }
/* .tag-red 样式定义 */

.tag-red { background: #fee2e2; color: #b91c1c; }
/* .tag-green 样式定义 */

.tag-green { background: #dcfce7; color: #15803d; }
/* .tag-gray 样式定义 */

.tag-gray { background: #f3f4f6; color: #6b7280; }
</style>

