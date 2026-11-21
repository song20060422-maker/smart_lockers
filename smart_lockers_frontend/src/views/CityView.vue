<script setup lang="ts">
// @ticket 03-首页：城市切换
// @proto 原型对齐：04-一期用户界面-2、06-用户界面规划-2 选择城市（城市网格/拼音搜索/规划中）
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { searchCities } from '@/api/home'
import { useLocationStore } from '@/stores/location'
import type { City } from '@/types'
// 执行 interface CityItem { name: string; pinyin?: string; planning?: boolean }，实现对应的前端逻辑

interface CityItem { name: string; pinyin?: string; planning?: boolean }
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 location，用于保存数据或函数引用

const location = useLocationStore()
// 声明变量 keyword，用于保存数据或函数引用

const keyword = ref('')
// 声明变量 cities，用于保存数据或函数引用

const cities = ref<CityItem[]>([])
// 声明变量 demoCities:，用于保存数据或函数引用

const demoCities: CityItem[] = [
  // 执行 { name: '北京', pinyin: 'beijing' }, { name: '上海', pinyin: 'shanghai' },，实现对应的前端逻辑

  { name: '北京', pinyin: 'beijing' }, { name: '上海', pinyin: 'shanghai' },
  // 执行 { name: '广州', pinyin: 'guangzhou' }, { name: '深圳', pinyin: 'shenzhen' },，实现对应的前端逻辑

  { name: '广州', pinyin: 'guangzhou' }, { name: '深圳', pinyin: 'shenzhen' },
  // 执行 { name: '杭州', pinyin: 'hangzhou' }, { name: '成都', pinyin: 'chengdu' },，实现对应的前端逻辑

  { name: '杭州', pinyin: 'hangzhou' }, { name: '成都', pinyin: 'chengdu' },
  // 执行 { name: '武汉', pinyin: 'wuhan' }, { name: '西安', pinyin: 'xian' },，实现对应的前端逻辑

  { name: '武汉', pinyin: 'wuhan' }, { name: '西安', pinyin: 'xian' },
  // 执行 { name: '南京', pinyin: 'nanjing', planning: true },，实现对应的前端逻辑

  { name: '南京', pinyin: 'nanjing', planning: true },
// 执行 ]，实现对应的前端逻辑

]
// 声明变量 shown，用于保存数据或函数引用

const shown = computed<CityItem[]>(() => {
  // 声明变量 k，用于保存数据或函数引用

  const k = keyword.value.trim().toLowerCase()
  // 条件判断：if (!k) return cities.value 时执行下面逻辑

  if (!k) return cities.value
  // 返回 cities.value.filter((c) => c.name.includes(k) || (c.pinyin || '').includes(k)) 给调用方

  return cities.value.filter((c) => c.name.includes(k) || (c.pinyin || '').includes(k))
// 当前代码块结束

})
// 定义函数/方法，封装可复用的处理逻辑：function choose(c: CityItem) {

function choose(c: CityItem) {
  // 条件判断：if (c.planning) return 时执行下面逻辑

  if (c.planning) return
  // 声明变量 full:，用于保存数据或函数引用

  const full: City = { id: 0, name: c.name, pinyin: c.pinyin || '', initial: '', sort: 0 }
  // 执行 location.setCity(full)，实现对应的前端逻辑

  location.setCity(full)
  // 执行 router.back()，实现对应的前端逻辑

  router.back()
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 执行 try {，实现对应的前端逻辑

  try {
    // 声明变量 res，用于保存数据或函数引用

    const res = await searchCities('')
    // 赋值更新：cities.value = Array.isArray(res) && res.length ? res.map((c) => ({ name: c.name, pinyin: c.pinyin })) : demoCities

    cities.value = Array.isArray(res) && res.length ? res.map((c) => ({ name: c.name, pinyin: c.pinyin })) : demoCities
  // 执行 } catch (e) {，实现对应的前端逻辑

  } catch (e) {
    // 赋值更新：cities.value = demoCities

    cities.value = demoCities
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

    <div class="search"><input v-model="keyword" placeholder="输入城市名或拼音搜索" /></div>
    <!-- div 元素，class 为 section-title -->

    <div class="section-title">当前定位 / 热门城市</div>
    <!-- div 元素，class 为 grid -->

    <div class="grid">
      <!-- button 元素，class 为 city -->

      <button v-for="c in shown" :key="c.name" class="city" :class="{ plan: c.planning }" @click="choose(c)">
        <!-- 插值文本，显示动态内容 -->

        {{ c.name }}<span v-if="c.planning" class="badge">2025规划中</span>
      <!-- 闭合 button 元素 -->

      </button>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- p 元素，class 为 empty -->

    <p v-if="shown.length === 0" class="empty">未找到相关城市</p>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 12px; }
/* .search input 样式定义 */

.search input { width: 100%; height: 40px; box-sizing: border-box; border: 1px solid #dbeafe; border-radius: 20px; padding: 0 16px; font-size: 14px; outline: none; }
/* .section-title 样式定义 */

.section-title { font-size: 14px; color: #6b7280; }
/* .grid 样式定义 */

.grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
/* .city 样式定义 */

.city { position: relative; height: 44px; border: 1px solid #bfdbfe; border-radius: 8px; background: #fff; color: #1f2937; font-size: 14px; }
/* .city.plan 样式定义 */

.city.plan { background: #f3f4f6; color: #9ca3af; }
/* .badge 样式定义 */

.badge { position: absolute; right: 2px; top: 2px; font-size: 9px; color: #9ca3af; }
/* .empty 样式定义 */

.empty { text-align: center; color: #9ca3af; padding: 30px 0; }
</style>

