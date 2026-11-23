<script setup lang="ts">
// @ticket 03-首页：定位/城市选择/附近寄存点/寄存指南
// @proto 原型对齐：04-一期用户界面-1 首页（蓝色头部/城市/搜索/快捷标签/查询按钮/四入口/交易保障/附近寄存点列表）
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { locate, searchCities, nearbySites, getGuide } from '@/api/home'
import { useLocationStore } from '@/stores/location'
import { useAuthStore } from '@/stores/auth'
import type { City, LockerSite, GuideResult } from '@/types'
// 引入全国城市搜索函数，作为后端 cities 表数据不全时的本地兜底
import { searchNationwide } from '@/data/cities'
// 引入地图选点组件，提供无需 API Key 的地点选择
import LocationPickerMap from '@/components/LocationPickerMap.vue'

// 声明变量 router = useRouter()，用于保存数据或函数引用
const router = useRouter()
// 声明变量 locStore = useLocationStore()，用于保存数据或函数引用
const locStore = useLocationStore()
// 声明变量 auth = useAuthStore()，用于保存数据或函数引用
const auth = useAuthStore()
// 声明变量 keyword = ref('')，用于保存数据或函数引用
const keyword = ref('')
// 声明变量 cities = ref<City[]>([])，用于保存数据或函数引用
const cities = ref<City[]>([])
// 声明变量 sites = ref<LockerSite[]>([])，用于保存数据或函数引用
const sites = ref<LockerSite[]>([])
// 声明变量 guide = ref<GuideResult | null>(null)，用于保存数据或函数引用
const guide = ref<GuideResult | null>(null)
// 声明变量 err = ref('')，用于保存数据或函数引用
const err = ref('')
// 声明变量 loading = ref(false)，用于保存数据或函数引用
const loading = ref(false)
// 声明变量 showMap = ref(false)，控制地图选点面板的显示与隐藏
const showMap = ref(false)
// 声明变量 pickedAddress = ref('')，保存地图选点回填后的详细地址文案
const pickedAddress = ref('')
// 快捷搜索标签（对照原型：郑州站/郑州东站/二七广场/中原福塔）
// 声明变量 quickTags，用于保存数据或函数引用
const quickTags = ['郑州站', '郑州东站', '二七广场', '中原福塔']

// @ticket 00-project：演示兜底数据——后端不可用时保证首页正常展示，避免暴露"系统内部错误"
// 声明常量 demoSites = [...]（附近寄存点演示数据，对照原型"中交锦兰荟南门处丰巢柜"）
const demoSites = [
  // 执行 { id: 1, name: '中交锦兰荟南门处丰巢柜', address: '中交锦兰荟南门处丰巢柜中交锦兰荟南门荟…', lat: 34.75, lng: 113.62, status: 1, distance: 320 },，实现对应的前端逻辑
  { id: 1, name: '中交锦兰荟南门处丰巢柜', address: '中交锦兰荟南门处丰巢柜中交锦兰荟南门荟…', lat: 34.75, lng: 113.62, status: 1, distance: 320 },
  // 执行 { id: 2, name: '嵩山北路新途径教育一楼前台', address: '嵩山北路新途径教育一楼前台寄存点', lat: 34.76, lng: 113.61, status: 1, distance: 850 },，实现对应的前端逻辑
  { id: 2, name: '嵩山北路新途径教育一楼前台', address: '嵩山北路新途径教育一楼前台寄存点', lat: 34.76, lng: 113.61, status: 1, distance: 850 },
// 结束切片/数组字面量

]
// 声明常量 demoCities = [...]（城市搜索演示数据）
const demoCities = [
  // 执行 { id: 1, name: '北京', pinyin: 'beijing', initial: 'B', sort: 1 },，实现对应的前端逻辑
  { id: 1, name: '北京', pinyin: 'beijing', initial: 'B', sort: 1 },
  // 执行 { id: 2, name: '上海', pinyin: 'shanghai', initial: 'S', sort: 2 },，实现对应的前端逻辑
  { id: 2, name: '上海', pinyin: 'shanghai', initial: 'S', sort: 2 },
  // 执行 { id: 3, name: '广州', pinyin: 'guangzhou', initial: 'G', sort: 3 },，实现对应的前端逻辑
  { id: 3, name: '广州', pinyin: 'guangzhou', initial: 'G', sort: 3 },
  // 执行 { id: 4, name: '深圳', pinyin: 'shenzhen', initial: 'S', sort: 4 },，实现对应的前端逻辑
  { id: 4, name: '深圳', pinyin: 'shenzhen', initial: 'S', sort: 4 },
// 结束切片/数组字面量

]

// @ticket 03-首页：初始化定位并加载附近寄存点与指南
// 调用 onMounted(async () => {，触发对应功能
onMounted(async () => {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理
  try {
    // 声明变量 r = await locate({ fallback_city: '北京' })，用于保存数据或函数引用
    const r = await locate({ fallback_city: '北京' })
    // 执行 locStore.applyLocate(r)，实现对应的前端逻辑
    locStore.applyLocate(r)
  // 捕获上述逻辑抛出的异常并进行错误处理
  } catch (_) { /* 定位失败用默认城市 */ }
  // 等待异步操作完成：await Promise.all([loadSites(), loadGuide()])
  await Promise.all([loadSites(), loadGuide()])
  // 条件判断：if (!auth.logged) auth.login(10001).catch(() => {}) 时执行下面逻辑
  if (!auth.logged) auth.login(10001).catch(() => {})
// 执行 })，实现对应的前端逻辑

})

// 定义函数/方法，封装可复用的处理逻辑：async function loadSites() {
async function loadSites() {
  // 赋值更新：loading.value = true
  loading.value = true
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理
  try {
    // 声明变量 r = await nearbySites({ lat: locStore.lat, lng: locStore.lng, page: 1, page_size: 20 })，用于保存数据或函数引用
    const r = await nearbySites({ lat: locStore.lat, lng: locStore.lng, page: 1, page_size: 20 })
    // 赋值更新：sites.value = r.list
    sites.value = r.list
  // 捕获上述逻辑抛出的异常并进行错误处理
  } catch (_) {
    // 后端不可用时回退到演示网点，避免向用户显示"系统内部错误"
    sites.value = demoSites
    // 以友好提示替代原始错误文本
    err.value = '当前展示演示网点数据'
  }
  // 无论成功失败都会执行的收尾逻辑
  finally { loading.value = false }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：async function doSearchCity() {
async function doSearchCity() {
  // 基于前端内置全国城市库检索，覆盖后端 cities 表仅少量城市的问题，且无需 API Key
  cities.value = searchNationwide(keyword.value)
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function pickCity(c: City) {
function pickCity(c: City) {
  // 执行 locStore.setCity(c)，实现对应的前端逻辑
  locStore.setCity(c)
  // 赋值更新：cities.value = []
  cities.value = []
  // 赋值更新：keyword.value = ''
  keyword.value = ''
  // 调用 loadSites()，触发对应功能
  loadSites()
// 当前代码块结束

}

// 快捷标签点击：回填搜索词并刷新寄存点（对照原型快捷标签交互）
// 定义函数/方法 pickTag，封装可复用的处理逻辑
function pickTag(t: string) {
  // 回填关键词
  keyword.value = t
  // 刷新附近寄存点列表
  loadSites()
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function onMapPicked(v) —— 地图选点回填
function onMapPicked(v: { lat: number; lng: number; address: string; city: string }) {
  // 写入经纬度、城市名与来源到定位 store，来源标记为 map_picker
  locStore.applyPicked(v.lat, v.lng, v.city, 'map_picker')
  // 保存详细地址文案用于展示
  pickedAddress.value = v.address
  // 关闭地图选点面板
  showMap.value = false
  // 按新坐标刷新附近寄存点
  loadSites()
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：async function loadGuide() {
async function loadGuide() {
  // 等待异步操作完成：try { guide.value = await getGuide() }
  try { guide.value = await getGuide() }
  // 捕获上述逻辑抛出的异常并进行错误处理
  catch (_) { /* 指南非关键路径 */ }
// 当前代码块结束

}

// 进入寄存流程：携带网点 id/名称/地址（柜子信息页头部展示）
// 定义函数/方法 goDeposit，封装可复用的处理逻辑
function goDeposit(site: LockerSite) {
  // 执行 router.push，携带 site_id/site_name/site_addr，实现对应的前端逻辑
  router.push({ path: '/deposit', query: { site_id: String(site.id), site_name: site.name, site_addr: site.address } })
// 当前代码块结束

}
</script>

<template>
  <!-- 页面根容器：首页（对照原型蓝色头部 + 白色面板 + 附近寄存点） -->
  <div class="home">
    <!-- 蓝色渐变头部：品牌标语 -->
    <div class="hero">
      <!-- 主标语 -->
      <div class="hero-title">放心寄存</div>
      <!-- 副标语 -->
      <div class="hero-sub">平台安全保障</div>
      <!-- 装饰行李箱（CSS） -->
      <div class="hero-bag"><i></i></div>
    <!-- 闭合头部 -->
    </div>

    <!-- 白色操作面板：城市/搜索/查询/入口/保障 -->
    <div class="panel">
      <!-- 城市行：定位城市 + 我的附近 -->
      <div class="city-row">
        <!-- 城市定位（点击展开城市搜索） -->
        <span class="city" @click="doSearchCity">📍 {{ locStore.city || '郑州市' }} ▾</span>
        <!-- 我的附近（地图选点入口） -->
        <span class="near" @click="showMap = !showMap">➤ 我的附近</span>
      <!-- 闭合城市行 -->
      </div>
      <!-- 搜索框：火车站/地铁站/景点 -->
      <div class="search-box">
        <!-- 搜索输入 -->
        <input v-model="keyword" placeholder="搜索火车站/地铁站/景点" @keyup.enter="doSearchCity" />
        <!-- 搜索图标 -->
        <span class="s-ico">🔍</span>
      <!-- 闭合搜索框 -->
      </div>
      <!-- 城市搜索结果（ chips 横排） -->
      <div v-if="cities.length" class="chips">
        <!-- 城市chip -->
        <span v-for="c in cities" :key="c.id" class="chip" @click="pickCity(c)">{{ c.name }}</span>
      <!-- 闭合结果 -->
      </div>
      <!-- 快捷标签行 -->
      <div class="chips">
        <!-- 快捷标签 -->
        <span v-for="t in quickTags" :key="t" class="chip" @click="pickTag(t)">{{ t }}</span>
      <!-- 闭合快捷标签 -->
      </div>
      <!-- 查询寄存点大按钮 -->
      <button class="query-btn" @click="loadSites">查询寄存点</button>

      <!-- 四个功能入口：我的订单/在线客服/寄存指南/优惠券（开发中角标） -->
      <div class="entries">
        <!-- 我的订单入口 -->
        <div class="entry" @click="router.push('/orders')">
          <!-- 图标 -->
          <span class="e-ico" style="background:#e8f1ff">📋</span>
          <!-- 名称 -->
          <span class="e-name">我的订单</span>
        <!-- 闭合入口 -->
        </div>
        <!-- 在线客服入口 -->
        <div class="entry" @click="router.push('/service')">
          <!-- 图标 -->
          <span class="e-ico" style="background:#ffeef2">💗</span>
          <!-- 名称 -->
          <span class="e-name">在线客服</span>
        <!-- 闭合入口 -->
        </div>
        <!-- 寄存指南入口 -->
        <div class="entry" @click="router.push('/guide')">
          <!-- 图标 -->
          <span class="e-ico" style="background:#fff3e0">📒</span>
          <!-- 名称 -->
          <span class="e-name">寄存指南</span>
        <!-- 闭合入口 -->
        </div>
        <!-- 优惠券入口（开发中角标） -->
        <div class="entry" @click="router.push('/coupons')">
          <!-- 开发中角标 -->
          <span class="dev-badge">开发中</span>
          <!-- 图标 -->
          <span class="e-ico" style="background:#ffe9ec">🎟️</span>
          <!-- 名称 -->
          <span class="e-name">优惠券</span>
        <!-- 闭合入口 -->
        </div>
      <!-- 闭合入口区 -->
      </div>

      <!-- 小程序交易保障条 -->
      <div class="guard-bar">
        <!-- 保障图标 -->
        <span class="g-ico">🛡️</span>
        <!-- 保障文案 -->
        <span class="g-text">小程序交易保障</span>
        <!-- 说明 -->
        <span class="g-sub">先行赔付 · 消费者权益</span>
      <!-- 闭合保障条 -->
      </div>
    <!-- 闭合操作面板 -->
    </div>

    <!-- 地图选点组件（展开时） -->
    <LocationPickerMap
      v-if="showMap"
      :initial-lat="locStore.lat"
      :initial-lng="locStore.lng"
      :initial-city="locStore.city"
      @pick="onMapPicked"
      @close="showMap = false"
    />
    <!-- 已选地点提示 -->
    <p v-if="pickedAddress" class="picked">已选地点：{{ pickedAddress }}</p>

    <!-- 附近寄存点区 -->
    <div class="nearby-sec">
      <!-- 区标题 -->
      <div class="nearby-title">附近寄存点</div>
      <!-- 加载提示 -->
      <p v-if="loading" class="muted">加载中...</p>
      <!-- 寄存点卡片列表 -->
      <div v-for="s in sites" :key="s.id" class="site-item" @click="goDeposit(s)">
        <!-- 柜机缩略图（CSS） -->
        <div class="s-thumb"><i></i></div>
        <!-- 信息区 -->
        <div class="s-info">
          <!-- 网点名 -->
          <div class="s-name">{{ s.name }}</div>
          <!-- 可用格口行（演示文案对照原型：大8|中8|小8） -->
          <div class="s-stock">可用格口： 大8 | 中8 | 小8</div>
          <!-- 地址 -->
          <div class="s-addr">📍 {{ s.address }}</div>
        <!-- 闭合信息区 -->
        </div>
      <!-- 闭合寄存点卡片 -->
      </div>
      <!-- 空状态 -->
      <p v-if="!sites.length && !loading" class="muted">暂无附近寄存点</p>
    <!-- 闭合附近区 -->
    </div>

    <!-- 演示提示 -->
    <p v-if="err" class="muted tip">{{ err }}</p>
  <!-- 闭合根容器 -->
  </div>
</template>

<style scoped>
/* 首页根容器 */
.home { margin: -12px -12px 0; }
/* 蓝色渐变头部 */
.hero { position: relative; background: linear-gradient(150deg, #2f7bff, #1f5fe8); color: #fff; padding: 22px 20px 46px; }
/* 主标语 */
.hero-title { font-size: 24px; font-weight: 800; letter-spacing: 2px; }
/* 副标语 */
.hero-sub { font-size: 13px; opacity: .9; margin-top: 6px; }
/* 行李箱装饰 */
.hero-bag { position: absolute; right: 22px; top: 18px; width: 46px; height: 56px; background: #4d8dff; border-radius: 8px; opacity: .9; }
/* 行李箱拉杆 */
.hero-bag i { position: absolute; top: -10px; left: 14px; width: 16px; height: 10px; border: 3px solid #cfe0ff; border-bottom: none; border-radius: 4px 4px 0 0; }
/* 白色操作面板：覆盖头部底部 */
.panel { position: relative; margin: -34px 12px 0; background: #fff; border-radius: 14px; padding: 14px; box-shadow: 0 6px 18px rgba(31,95,232,.10); }
/* 城市行 */
.city-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
/* 城市 */
.city { font-size: 15px; font-weight: 700; color: #111827; cursor: pointer; }
/* 我的附近 */
.near { font-size: 13px; color: #2563eb; cursor: pointer; }
/* 搜索框 */
.search-box { position: relative; margin-bottom: 10px; }
/* 搜索输入 */
.search-box input { width: 100%; height: 42px; border: none; background: #f3f5f9; border-radius: 21px; padding: 0 40px 0 16px; font-size: 14px; box-sizing: border-box; }
/* 搜索图标 */
.s-ico { position: absolute; right: 14px; top: 10px; font-size: 15px; opacity: .6; }
/* chips 横排 */
.chips { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 10px; }
/* 单个 chip */
.chip { background: #f3f5f9; color: #374151; font-size: 12px; border-radius: 14px; padding: 5px 12px; cursor: pointer; }
/* 查询大按钮 */
.query-btn { width: 100%; height: 46px; border: none; border-radius: 23px; background: #2563eb; color: #fff; font-size: 16px; font-weight: 600; margin: 2px 0 14px; cursor: pointer; }
/* 四入口区 */
.entries { display: flex; justify-content: space-between; margin-bottom: 14px; }
/* 单个入口 */
.entry { position: relative; display: flex; flex-direction: column; align-items: center; gap: 6px; cursor: pointer; flex: 1; }
/* 入口图标底 */
.e-ico { width: 44px; height: 44px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 20px; }
/* 入口名称 */
.e-name { font-size: 12px; color: #374151; }
/* 开发中角标 */
.dev-badge { position: absolute; top: -8px; right: 4px; background: #ff5a1f; color: #fff; font-size: 10px; border-radius: 8px 8px 8px 0; padding: 1px 6px; }
/* 交易保障条 */
.guard-bar { display: flex; align-items: center; gap: 8px; background: #f6f8fb; border-radius: 10px; padding: 10px 12px; }
/* 保障图标 */
.g-ico { font-size: 16px; }
/* 保障文案 */
.g-text { font-size: 13px; font-weight: 600; color: #111827; }
/* 保障说明 */
.g-sub { font-size: 12px; color: #9ca3af; margin-left: auto; }
/* 已选地点提示 */
.picked { font-size: 12px; color: #6b7280; padding: 8px 16px 0; }
/* 附近寄存点区 */
.nearby-sec { padding: 14px 16px 10px; }
/* 区标题 */
.nearby-title { font-size: 16px; font-weight: 700; color: #111827; margin-bottom: 10px; }
/* 寄存点卡片 */
.site-item { display: flex; gap: 12px; background: #fff; border-radius: 12px; padding: 12px; margin-bottom: 10px; cursor: pointer; }
/* 柜机缩略图 */
.s-thumb { width: 56px; height: 56px; border-radius: 10px; background: #eef3fb; flex-shrink: 0; position: relative; }
/* 缩略图柜格造型 */
.s-thumb i { position: absolute; inset: 10px; background: linear-gradient(160deg, #9dbdf2, #5f8ef0); border-radius: 6px; box-shadow: inset 0 -10px 0 rgba(255,255,255,.4); }
/* 信息区 */
.s-info { flex: 1; min-width: 0; }
/* 网点名 */
.s-name { font-size: 15px; font-weight: 600; color: #111827; }
/* 可用格口 */
.s-stock { font-size: 12px; color: #6b7280; margin: 4px 0; }
/* 地址 */
.s-addr { font-size: 12px; color: #9ca3af; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
/* 通用灰字 */
.muted { font-size: 13px; color: #9ca3af; }
/* 提示 */
.tip { padding: 0 16px 12px; }
</style>
