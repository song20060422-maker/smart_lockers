<script setup lang="ts">
// @ticket 06-个人中心：我的卡券列表（可使用/已过期）
// @proto 原型对齐：01-C端优惠券-1~5、06-用户界面规划-13/14 优惠券列表与规则
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import CouponCard from '@/components/CouponCard.vue'
import { myCoupons } from '@/api/user'
// 执行 interface CouponItem { id?: number; code?: string; title?: string; amount?: number; discount?: number; hours?: number; expire_at?: string; scope?: string; expired?: boolean }，实现对应的前端逻辑

interface CouponItem { id?: number; code?: string; title?: string; amount?: number; discount?: number; hours?: number; expire_at?: string; scope?: string; expired?: boolean }
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 tab，用于保存数据或函数引用

const tab = ref<'available' | 'expired'>('available')
// 声明变量 list，用于保存数据或函数引用

const list = ref<CouponItem[]>([])
// 声明变量 tip，用于保存数据或函数引用

const tip = ref('')
// 声明变量 demoList:，用于保存数据或函数引用

const demoList: CouponItem[] = [
  // 执行 { id: 1, code: 'NEW2026', title: '新客立减券', amount: 5, expire_at: '2026-12-31', scope: '全场寄存柜' },，实现对应的前端逻辑

  { id: 1, code: 'NEW2026', title: '新客立减券', amount: 5, expire_at: '2026-12-31', scope: '全场寄存柜' },
  // 执行 { id: 2, code: 'HOUR3', title: '时长赠送券', hours: 3, expire_at: '2026-11-30', scope: '按小时计费' },，实现对应的前端逻辑

  { id: 2, code: 'HOUR3', title: '时长赠送券', hours: 3, expire_at: '2026-11-30', scope: '按小时计费' },
// 执行 ]，实现对应的前端逻辑

]
// 定义函数/方法，封装可复用的处理逻辑：function shown(): CouponItem[] {

function shown(): CouponItem[] {
  // 返回 list.value.filter((c) => (tab.value === 'expired') === !!c.expired) 给调用方

  return list.value.filter((c) => (tab.value === 'expired') === !!c.expired)
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function goDeposit() {

function goDeposit() {
  // 执行 router.push('/deposit')，实现对应的前端逻辑

  router.push('/deposit')
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 执行 try {，实现对应的前端逻辑

  try {
    // 声明变量 res，用于保存数据或函数引用

    const res = await myCoupons()
    // 赋值更新：list.value = Array.isArray(res) ? res as CouponItem[] : demoList

    list.value = Array.isArray(res) ? res as CouponItem[] : demoList
    // 赋值更新：tip.value = ''

    tip.value = ''
  // 执行 } catch (e) {，实现对应的前端逻辑

  } catch (e) {
    // 赋值更新：list.value = demoList

    list.value = demoList
    // 赋值更新：tip.value = '当前展示演示数据'

    tip.value = '当前展示演示数据'
  // 当前代码块结束

  }
// 当前代码块结束

}
// 执行 onMounted(load)，实现对应的前端逻辑

onMounted(load)
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 tabs -->

    <div class="tabs">
      <!-- button 元素，class 为 { on: tab === 'available' } -->

      <button :class="{ on: tab === 'available' }" @click="tab = 'available'">可使用</button>
      <!-- button 元素，class 为 { on: tab === 'expired' } -->

      <button :class="{ on: tab === 'expired' }" @click="tab = 'expired'">已过期</button>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- p 元素，class 为 tip -->

    <p v-if="tip" class="tip">{{ tip }}</p>
    <!-- CouponCard 元素 -->

    <CouponCard v-for="c in shown()" :key="c.id || c.code" :coupon="c" />
    <!-- div 元素，class 为 empty -->

    <div v-if="shown().length === 0" class="empty">暂无优惠券</div>
    <!-- button 元素，class 为 primary -->

    <button class="primary" @click="goDeposit">去寄存</button>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 10px; }
/* .tabs 样式定义 */

.tabs { display: flex; gap: 8px; }
/* .tabs button 样式定义 */

.tabs button { flex: 1; height: 36px; border: 1px solid #bfdbfe; border-radius: 18px; background: #fff; color: #2563eb; font-size: 14px; }
/* .tabs button.on 样式定义 */

.tabs button.on { background: #2563eb; color: #fff; }
/* .tip 样式定义 */

.tip { font-size: 12px; color: #9ca3af; }
/* .empty 样式定义 */

.empty { text-align: center; color: #9ca3af; padding: 30px 0; font-size: 14px; }
/* .primary 样式定义 */

.primary { margin-top: 8px; height: 42px; border: 0; border-radius: 21px; background: #2563eb; color: #fff; font-size: 15px; }
</style>

