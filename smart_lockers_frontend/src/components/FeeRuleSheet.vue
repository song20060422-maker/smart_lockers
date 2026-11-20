<script setup lang="ts">
// @proto 原型对齐：03-预支付相关提示-3/4、04-一期用户界面-8 收费规则弹窗
// @ticket 04-物品寄存：收费规则（租柜费/封顶/营业截止）可配置化展示
// 执行 defineProps<{ visible: boolean }>()，实现对应的前端逻辑

defineProps<{ visible: boolean }>()
// 声明变量 emit，用于保存数据或函数引用

const emit = defineEmits<{ (e: 'close'): void }>()
// 定义函数/方法，封装可复用的处理逻辑：function close() {

function close() {
  // 执行 emit('close')，实现对应的前端逻辑

  emit('close')
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 mask -->

  <div v-if="visible" class="mask" @click.self="close">
    <!-- div 元素，class 为 sheet -->

    <div class="sheet">
      <!-- div 元素，class 为 sheet-title -->

      <div class="sheet-title">收费规则</div>
      <!-- div 元素，class 为 rule-row -->

      <div class="rule-row"><span>计费方式</span><span>按小时 / 按天可选</span></div>
      <!-- div 元素，class 为 rule-row -->

      <div class="rule-row"><span>小格租柜费</span><span>1 元/小时，10 元/天</span></div>
      <!-- div 元素，class 为 rule-row -->

      <div class="rule-row"><span>中格租柜费</span><span>2 元/小时，20 元/天</span></div>
      <!-- div 元素，class 为 rule-row -->

      <div class="rule-row"><span>大格租柜费</span><span>3 元/小时，30 元/天</span></div>
      <!-- div 元素，class 为 rule-row -->

      <div class="rule-row"><span>封顶说明</span><span>单日按柜型封顶 10/20/30 元</span></div>
      <!-- div 元素，class 为 rule-row -->

      <div class="rule-row"><span>押金</span><span>99 元，结单后原路退还</span></div>
      <!-- div 元素，class 为 rule-row -->

      <div class="rule-row"><span>营业截止</span><span>请在营业结束前取出物品</span></div>
      <!-- button 元素，class 为 primary -->

      <button class="primary" @click="close">我知道了</button>
    <!-- 闭合 div 元素 -->

    </div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .mask 样式定义 */

.mask { position: fixed; inset: 0; background: rgba(0,0,0,.4); display: flex; align-items: flex-end; justify-content: center; z-index: 60; }
/* .sheet 样式定义 */

.sheet { width: 100%; max-width: 430px; background: #fff; border-radius: 16px 16px 0 0; padding: 20px 18px 16px; }
/* .sheet-title 样式定义 */

.sheet-title { font-size: 17px; font-weight: 700; color: #1f2937; text-align: center; margin-bottom: 16px; }
/* .rule-row 样式定义 */

.rule-row { display: flex; justify-content: space-between; font-size: 13px; color: #4b5563; padding: 8px 0; border-bottom: 1px solid #f3f4f6; }
/* .rule-row span:last-child 样式定义 */

.rule-row span:last-child { color: #1f2937; }
/* .primary 样式定义 */

.primary { margin-top: 16px; width: 100%; height: 42px; border: 0; border-radius: 21px; background: #2563eb; color: #fff; font-size: 15px; }
</style>

