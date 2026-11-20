<script setup lang="ts">
// @ticket 03-首页：附近寄存点搜索
// @proto 原型对齐：04-一期用户界面-3/4/5、06-用户界面规划-3/4/5 附近寄存点（分类入口/空闲数/设备离线/搜索）
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { nearbySites } from '@/api/home'
import { useLocationStore } from '@/stores/location'
// 执行 interface SiteItem { id?: number; name: string; address?: string; distance?: string; free?: number; total?: number; offline?: boolean; category?: string }，实现对应的前端逻辑

interface SiteItem { id?: number; name: string; address?: string; distance?: string; free?: number; total?: number; offline?: boolean; category?: string }
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 location，用于保存数据或函数引用

const location = useLocationStore()
// 声明变量 keyword，用于保存数据或函数引用

const keyword = ref('')
// 声明变量 category，用于保存数据或函数引用

const category = ref('全部')
// 声明变量 cats，用于保存数据或函数引用

const cats = ['全部', '商圈', '车站', '景区', '医院', '校园']
// 声明变量 list，用于保存数据或函数引用

const list = ref<SiteItem[]>([])
// 声明变量 demo:，用于保存数据或函数引用

const demo: SiteItem[] = [
  // 执行 { id: 1, name: '国贸中心寄存点', address: '建国门外大街1号', distance: '230m', free: 8, total: 12, category: '商圈' },，实现对应的前端逻辑

  { id: 1, name: '国贸中心寄存点', address: '建国门外大街1号', distance: '230m', free: 8, total: 12, category: '商圈' },
  // 执行 { id: 2, name: '北京南站寄存点', address: '永外大街车站广场', distance: '1.2km', free: 3, total: 20, category: '车站' },，实现对应的前端逻辑

  { id: 2, name: '北京南站寄存点', address: '永外大街车站广场', distance: '1.2km', free: 3, total: 20, category: '车站' },
  // 执行 { id: 3, name: '故宫北门寄存点', address: '景山前街4号', distance: '2.5km', free: 0, total: 10, offline: true, category: '景区' },，实现对应的前端逻辑

  { id: 3, name: '故宫北门寄存点', address: '景山前街4号', distance: '2.5km', free: 0, total: 10, offline: true, category: '景区' },
// 执行 ]，实现对应的前端逻辑

]
// 定义函数/方法，封装可复用的处理逻辑：function shown(): SiteItem[] {

function shown(): SiteItem[] {
  // 声明变量 k，用于保存数据或函数引用

  const k = keyword.value.trim()
  // 返回 list.value.filter((s) => (category.value === '全部' || s.category === category.value) && (!k || s.name.includes(k) || (s.address || '').includes(k))) 给调用方

  return list.value.filter((s) => (category.value === '全部' || s.category === category.value) && (!k || s.name.includes(k) || (s.address || '').includes(k)))
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function openSite(s: SiteItem) {

function openSite(s: SiteItem) {
  // 条件判断：if (s.offline) return 时执行下面逻辑

  if (s.offline) return
  // 执行 router.push({ path: '/deposit', query: { site: String(s.id || '') } })，实现对应的前端逻辑

  router.push({ path: '/deposit', query: { site: String(s.id || '') } })
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 执行 try {，实现对应的前端逻辑

  try {
    // 声明变量 res，用于保存数据或函数引用

    const res = await nearbySites({ lat: location.lat, lng: location.lng })
    // 声明变量 arr，用于保存数据或函数引用

    const arr = (res && Array.isArray(res.list)) ? res.list.map((s) => ({
      // 执行 id: s.id, name: s.name, address: s.address,，实现对应的前端逻辑

      id: s.id, name: s.name, address: s.address,
      // 执行 distance: s.distance != null ? s.distance + 'm' : '',，实现对应的前端逻辑

      distance: s.distance != null ? s.distance + 'm' : '',
      // 执行 free: s.status === 1 ? 6 : 0, total: 12, offline: s.status !== 1, category: '商圈',，实现对应的前端逻辑

      free: s.status === 1 ? 6 : 0, total: 12, offline: s.status !== 1, category: '商圈',
    // 执行 })) : []，实现对应的前端逻辑

    })) : []
    // 赋值更新：list.value = arr.length ? arr : demo

    list.value = arr.length ? arr : demo
  // 执行 } catch (e) {，实现对应的前端逻辑

  } catch (e) {
    // 赋值更新：list.value = demo

    list.value = demo
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
    <!-- div 元素，class 为 search -->

    <div class="search"><input v-model="keyword" placeholder="搜索站点名称 / 地址" /></div>
    <!-- div 元素，class 为 cats -->

    <div class="cats">
      <!-- button 元素，class 为 { on: category === c } -->

      <button v-for="c in cats" :key="c" :class="{ on: category === c }" @click="category = c">{{ c }}</button>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- div 元素，class 为 site card -->

    <div v-for="s in shown()" :key="s.id || s.name" class="site card" :class="{ off: s.offline }" @click="openSite(s)">
      <!-- div 元素，class 为 row1 -->

      <div class="row1"><span class="name">{{ s.name }}</span><span class="dist">{{ s.distance }}</span></div>
      <!-- div 元素，class 为 addr -->

      <div class="addr">{{ s.address }}</div>
      <!-- div 元素，class 为 row2 -->

      <div class="row2">
        <!-- span 元素，class 为 off-tag -->

        <span v-if="s.offline" class="off-tag">设备离线</span>
        <!-- span 元素，class 为 free -->

        <span v-else class="free">空闲 {{ s.free }}/{{ s.total }}</span>
      <!-- 闭合 div 元素 -->

      </div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- div 元素，class 为 empty -->

    <div v-if="shown().length === 0" class="empty">附近暂无寄存点</div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 10px; }
/* .search input 样式定义 */

.search input { width: 100%; height: 40px; box-sizing: border-box; border: 1px solid #dbeafe; border-radius: 20px; padding: 0 16px; font-size: 14px; outline: none; }
/* .cats 样式定义 */

.cats { display: flex; gap: 6px; overflow-x: auto; }
/* .cats button 样式定义 */

.cats button { flex-shrink: 0; padding: 4px 12px; border: 1px solid #bfdbfe; border-radius: 16px; background: #fff; color: #2563eb; font-size: 12px; }
/* .cats button.on 样式定义 */

.cats button.on { background: #2563eb; color: #fff; }
/* .card 样式定义 */

.card { background: #fff; border-radius: 10px; padding: 12px; border: 1px solid #e5e7eb; }
/* .card.off 样式定义 */

.card.off { opacity: .6; }
/* .row1 样式定义 */

.row1 { display: flex; justify-content: space-between; }
/* .name 样式定义 */

.name { font-size: 15px; font-weight: 600; color: #1f2937; }
/* .dist 样式定义 */

.dist { font-size: 12px; color: #6b7280; }
/* .addr 样式定义 */

.addr { font-size: 12px; color: #9ca3af; margin: 4px 0; }
/* .row2 样式定义 */

.row2 { font-size: 12px; }
/* .free 样式定义 */

.free { color: #16a34a; }
/* .off-tag 样式定义 */

.off-tag { color: #ef4444; }
/* .empty 样式定义 */

.empty { text-align: center; color: #9ca3af; padding: 30px 0; }
</style>

