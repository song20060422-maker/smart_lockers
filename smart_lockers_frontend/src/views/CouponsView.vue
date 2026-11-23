<script setup lang="ts">
// @ticket 06-个人中心：我的卡券列表（未使用 / 已使用/已失效）
// @proto 原型对齐：01-C端优惠券-1 我的优惠券页（tab 文案、券类型与限网点/不可用态演示数据）
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import CouponCard from '@/components/CouponCard.vue'
import { myCoupons } from '@/api/user'
// 引入统一券数据工具：与寄存页共用同一套归一化真实券数据
import { normalizeCoupons, toCardCoupon } from '@/utils/coupon'
// 执行 interface CouponItem { id?: number; code?: string; title?: string; amount?: number; discount?: number; hours?: number; expire_at?: string; scope?: string; expired?: boolean; usable?: boolean }，实现对应的前端逻辑

interface CouponItem { id?: number; code?: string; title?: string; amount?: number; discount?: number; hours?: number; expire_at?: string; scope?: string; expired?: boolean; usable?: boolean }
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 tab，用于保存数据或函数引用

const tab = ref<'available' | 'expired'>('available')
// 声明变量 list，用于保存数据或函数引用

const list = ref<CouponItem[]>([])
// 声明变量 tip，用于保存数据或函数引用

const tip = ref('')
// 演示券数据：完全对照原型 01-C端优惠券-1 的四张券（1元/1小时/8折限网点/1小时本网点不可用）

const demoList: CouponItem[] = [
  // 执行 { id: 1, title: '1元优惠券', amount: 1, expire_at: '23:23:23', scope: '中原区二七万达广场兑换' },，实现对应的前端逻辑

  { id: 1, title: '1元优惠券', amount: 1, expire_at: '23:23:23', scope: '中原区二七万达广场兑换', usable: true },
  // 执行 { id: 2, title: '1小时免费寄存券', hours: 1, expire_at: '23:23:23', scope: '限【中原区二七万达广场】网点使用' },，实现对应的前端逻辑

  { id: 2, title: '1小时免费寄存券', hours: 1, expire_at: '23:23:23', scope: '限【中原区二七万达广场】网点使用', usable: true },
  // 执行 { id: 3, title: '8折优惠券', discount: 8, expire_at: '23:23:23', scope: '中原区二七万达广场兑换' },，实现对应的前端逻辑

  { id: 3, title: '8折优惠券', discount: 8, expire_at: '23:23:23', scope: '中原区二七万达广场兑换', usable: false },
  // 执行 { id: 4, title: '1小时免费寄存券', hours: 1, expire_at: '23:23:23', scope: '限【中原区二七万达广场】网点使用' },，实现对应的前端逻辑

  { id: 4, title: '1小时免费寄存券', hours: 1, expire_at: '23:23:23', scope: '限【中原区二七万达广场】网点使用', usable: false },
// 执行 ]，实现对应的前端逻辑

]
// 定义函数/方法，封装可复用的处理逻辑：function shown(): CouponItem[] {
// 未使用 tab：可用或不可用（未过期）的券；已失效 tab：expired 的券

function shown(): CouponItem[] {
  // 按当前 tab 过滤
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
    // 归一化后端 {available,expired} 对象为统一券数组，逐字段空值防御

    const normalized = normalizeCoupons(res)
    // 赋值更新：真实券数据优先（与寄存页同源），为空时保留原型演示券兜底，不删改既有数据

    list.value = normalized.length > 0 ? normalized.map((c) => toCardCoupon(c) as CouponItem) : demoList
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
  <!-- 页面根容器 -->
  <div class="col">
    <!-- 顶部 tab：未使用 / 已使用/已失效（对照原型 01-C端优惠券-1） -->
    <div class="tabs">
      <!-- 未使用 tab -->
      <button :class="{ on: tab === 'available' }" @click="tab = 'available'">未使用</button>
      <!-- 已使用/已失效 tab -->
      <button :class="{ on: tab === 'expired' }" @click="tab = 'expired'">已使用/已失效</button>
    <!-- 闭合 tab -->
    </div>
    <!-- 演示数据提示 -->
    <p v-if="tip" class="tip">{{ tip }}</p>
    <!-- 券卡片列表（含限网点/不可用态） -->
    <CouponCard v-for="c in shown()" :key="c.id || c.code" :coupon="c" />
    <!-- 空状态 -->
    <div v-if="shown().length === 0" class="empty">暂无优惠券</div>
    <!-- 去寄存按钮 -->
    <button class="primary" @click="goDeposit">去寄存</button>
  <!-- 闭合根容器 -->
  </div>
</template>

<style scoped>
/* 根容器纵向排布 */
.col { display: flex; flex-direction: column; gap: 12px; }
/* 顶部 tab：白底横条，蓝色下划线高亮（对照原型） */
.tabs { display: flex; background: #fff; border-radius: 10px; overflow: hidden; }
/* tab 按钮 */
.tabs button { flex: 1; height: 44px; border: none; background: transparent; color: #6b7280; font-size: 15px; position: relative; cursor: pointer; }
/* 激活 tab：蓝字 + 底部蓝色短下划线 */
.tabs button.on { color: #2563eb; font-weight: 600; }
/* 激活下划线 */
.tabs button.on::after { content: ''; position: absolute; left: 50%; transform: translateX(-50%); bottom: 6px; width: 28px; height: 3px; border-radius: 2px; background: #2563eb; }
/* 提示灰字 */
.tip { font-size: 12px; color: #9ca3af; }
/* 空状态 */
.empty { text-align: center; color: #9ca3af; padding: 40px 0; font-size: 14px; }
/* 去寄存主按钮 */
.primary { height: 44px; border: 0; border-radius: 22px; background: #2563eb; color: #fff; font-size: 15px; }
</style>
