<script setup lang="ts">
// @ticket 04-物品寄存 + 05-订单查看与取出：异常换柜/重新开门/取消订单
// @proto 原型对齐：02-开门异常提示优化-1~7 柜门已开/重新开门/已解决/取消订单退费/换柜/剩余次数/联系方式
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 route，用于保存数据或函数引用

const route = useRoute()
// 声明变量 scene，用于保存数据或函数引用

const scene = ref<string>(route.query.scene as string || 'opened')
// 声明变量 remain，用于保存数据或函数引用

const remain = ref<number>(Number(route.query.remain) || 3)
// 声明变量 titleMap:，用于保存数据或函数引用

const titleMap: Record<string, string> = {
  // 执行 opened: '柜门已打开，请及时存取',，实现对应的前端逻辑

  opened: '柜门已打开，请及时存取',
  // 执行 notclosed: '柜门未关闭，请重新开门',，实现对应的前端逻辑

  notclosed: '柜门未关闭，请重新开门',
  // 执行 failed: '开门失败，请稍后重试',，实现对应的前端逻辑

  failed: '开门失败，请稍后重试',
  // 执行 occupied: '该柜门已被占用，请更换柜门',，实现对应的前端逻辑

  occupied: '该柜门已被占用，请更换柜门',
// 当前代码块结束

}
// 声明变量 descMap:，用于保存数据或函数引用

const descMap: Record<string, string> = {
  // 执行 opened: '为避免物品丢失，请在 30 秒内完成存取并关好柜门。',，实现对应的前端逻辑

  opened: '为避免物品丢失，请在 30 秒内完成存取并关好柜门。',
  // 执行 notclosed: '检测到柜门仍处于开启状态，请重新扫码开柜并确认关闭。',，实现对应的前端逻辑

  notclosed: '检测到柜门仍处于开启状态，请重新扫码开柜并确认关闭。',
  // 执行 failed: '设备暂时无响应，请检查网络后重试，或联系客服。',，实现对应的前端逻辑

  failed: '设备暂时无响应，请检查网络后重试，或联系客服。',
  // 执行 occupied: '当前柜门已被其他用户使用，可更换柜门继续寄存。',，实现对应的前端逻辑

  occupied: '当前柜门已被其他用户使用，可更换柜门继续寄存。',
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function reopen() {

function reopen() {
  // 赋值更新：scene.value = 'opened'

  scene.value = 'opened'
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function resolved() {

function resolved() {
  // 执行 router.push('/orders')，实现对应的前端逻辑

  router.push('/orders')
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function cancelOrder() {

function cancelOrder() {
  // 执行 router.push({ path: '/orders', query: { action: 'refund' } })，实现对应的前端逻辑

  router.push({ path: '/orders', query: { action: 'refund' } })
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function changeDoor() {

function changeDoor() {
  // 执行 router.push('/deposit')，实现对应的前端逻辑

  router.push('/deposit')
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 icon -->

    <div class="icon">!</div>
    <!-- div 元素，class 为 title -->

    <div class="title">{{ titleMap[scene] || '操作异常' }}</div>
    <!-- div 元素，class 为 desc -->

    <div class="desc">{{ descMap[scene] || '请稍后重试或联系客服。' }}</div>
    <!-- div 元素，class 为 remain -->

    <div class="remain" v-if="scene === 'failed'">今日剩余重试次数：{{ remain }} 次</div>
    <!-- div 元素，class 为 actions -->

    <div class="actions">
      <!-- button 元素，class 为 primary -->

      <button v-if="scene === 'notclosed' || scene === 'failed'" class="primary" @click="reopen">重新开门</button>
      <!-- button 元素，class 为 primary -->

      <button v-if="scene === 'occupied'" class="primary" @click="changeDoor">更换柜门</button>
      <!-- button 元素，class 为 primary -->

      <button v-if="scene === 'opened'" class="primary" @click="resolved">我已解决</button>
      <!-- button 元素，class 为 ghost -->

      <button class="ghost" @click="cancelOrder">取消订单并退费</button>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- div 元素，class 为 service -->

    <div class="service">仍有问题？<span @click="router.push('/service')">联系在线客服 ›</span></div>
    <!-- div 元素，class 为 tel -->

    <div class="tel">客服电话：400-XXX-XXXX</div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; align-items: stretch; gap: 12px; }
/* .icon 样式定义 */

.icon { width: 56px; height: 56px; margin: 6px auto 0; border-radius: 50%; background: #fef3c7; color: #d97706; font-size: 32px; font-weight: 700; display: flex; align-items: center; justify-content: center; }
/* .title 样式定义 */

.title { text-align: center; font-size: 18px; font-weight: 700; color: #1f2937; }
/* .desc 样式定义 */

.desc { text-align: center; font-size: 13px; color: #6b7280; line-height: 1.6; }
/* .remain 样式定义 */

.remain { text-align: center; font-size: 12px; color: #ef4444; }
/* .actions 样式定义 */

.actions { display: flex; flex-direction: column; gap: 10px; margin-top: 8px; }
/* .actions button 样式定义 */

.actions button { height: 42px; border-radius: 21px; font-size: 15px; }
/* .primary 样式定义 */

.primary { border: 0; background: #2563eb; color: #fff; }
/* .ghost 样式定义 */

.ghost { border: 1px solid #d1d5db; background: #fff; color: #6b7280; }
/* .service 样式定义 */

.service { text-align: center; font-size: 13px; color: #6b7280; margin-top: 8px; }
/* .service span 样式定义 */

.service span { color: #2563eb; }
/* .tel 样式定义 */

.tel { text-align: center; font-size: 12px; color: #9ca3af; }
</style>

