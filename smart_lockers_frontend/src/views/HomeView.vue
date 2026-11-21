<script setup lang="ts">
// @ticket 03-首页：定位/城市选择/附近寄存点/寄存指南
// @proto 原型对齐：v-sites 网点柜组监控；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { locate, searchCities, nearbySites, getGuide } from '@/api/home'
import { useLocationStore } from '@/stores/location'
import { useAuthStore } from '@/stores/auth'
import type { City, LockerSite, GuideResult } from '@/types'

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

// @ticket 00-project：演示兜底数据——后端不可用时保证首页正常展示，避免暴露“系统内部错误”
// 声明常量 demoSites = [...]（附近寄存点演示数据）

const demoSites = [
  // 执行 { id: 1, name: '中关村寄存点', address: '海淀区中关村大街 1 号', lat: 39.98, lng: 116.31, status: 1, distance: 320 },，实现对应的前端逻辑

  { id: 1, name: '中关村寄存点', address: '海淀区中关村大街 1 号', lat: 39.98, lng: 116.31, status: 1, distance: 320 },
  // 执行 { id: 2, name: '国贸寄存点', address: '朝阳区建国路 88 号', lat: 39.91, lng: 116.46, status: 1, distance: 850 },，实现对应的前端逻辑

  { id: 2, name: '国贸寄存点', address: '朝阳区建国路 88 号', lat: 39.91, lng: 116.46, status: 1, distance: 850 },
  // 执行 { id: 3, name: '三里屯寄存点', address: '朝阳区三里屯路 19 号', lat: 39.94, lng: 116.45, status: 1, distance: 1200 },，实现对应的前端逻辑

  { id: 3, name: '三里屯寄存点', address: '朝阳区三里屯路 19 号', lat: 39.94, lng: 116.45, status: 1, distance: 1200 },
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
    // 后端不可用时回退到演示网点，避免向用户显示“系统内部错误”
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
  // 等待异步操作完成：try { cities.value = await searchCities(keyword.value) }

  try { cities.value = await searchCities(keyword.value) }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (_) {
    // 后端不可用时用演示城市按关键字模糊过滤
    cities.value = demoCities.filter((c) => !keyword.value || c.name.includes(keyword.value))
    // 以友好提示替代原始错误文本

    err.value = '当前展示演示城市数据'
  }
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

// 定义函数/方法，封装可复用的处理逻辑：async function loadGuide() {

async function loadGuide() {
  // 等待异步操作完成：try { guide.value = await getGuide() }

  try { guide.value = await getGuide() }
  // 捕获上述逻辑抛出的异常并进行错误处理

  catch (_) { /* 指南非关键路径 */ }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function goDeposit(site: LockerSite) {

function goDeposit(site: LockerSite) {
  // 执行 router.push({ path: '/deposit', query: { site_id: String(site.id) } })，实现对应的前端逻辑

  router.push({ path: '/deposit', query: { site_id: String(site.id) } })
// 当前代码块结束

}
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列定位、寄存点与指南卡片 -->

  <div class="col">
    <!-- 卡片容器，class 为 card，展示首页定位与城市搜索 -->

    <div class="card">
      <!-- 二级标题，显示“首页定位” -->

      <h2>首页定位</h2>
      <!-- 段落元素，class 为 muted，显示当前城市 locStore.city（缺省“定位中”）与来源 locStore.source（缺省“-”） -->

      <p class="muted">当前城市：{{ locStore.city || '定位中' }}（来源：{{ locStore.source || '-' }}）</p>
      <!-- 行容器，class 为 row，横向排列城市搜索输入框与按钮 -->

      <div class="row">
        <!-- 城市搜索输入框，v-model 双向绑定 keyword，占位提示“搜索城市”，回车键触发 doSearchCity -->

        <input v-model="keyword" placeholder="搜索城市" @keyup.enter="doSearchCity" />
        <!-- 搜索城市按钮，点击触发 doSearchCity -->

        <button @click="doSearchCity">搜索城市</button>
      <!-- 闭合城市搜索行容器 div 元素 -->

      </div>
      <!-- 城市结果容器，class 为 list-item，当 cities 长度大于 0 时渲染 -->

      <div v-if="cities.length" class="list-item">
        <!-- 基于 cities 数组循环生成城市按钮，key 取 c.id，样式类为 btn-ghost，点击触发 pickCity(c)，文本插值 c.name -->

        <button v-for="c in cities" :key="c.id" class="btn-ghost" @click="pickCity(c)">{{ c.name }}</button>
      <!-- 闭合城市结果容器 div 元素 -->

      </div>
    <!-- 闭合首页定位卡片 div 元素 -->

    </div>

    <!-- 卡片容器，class 为 card，展示附近寄存点列表 -->

    <div class="card">
      <!-- 二级标题，显示“附近寄存点” -->

      <h2>附近寄存点</h2>
      <!-- 加载提示段落，当 loading 为真时渲染，class 为 muted，显示“加载中...” -->

      <p v-if="loading" class="muted">加载中...</p>
      <!-- 基于 sites 数组循环生成寄存点条目，key 取 s.id，样式类为 list-item -->

      <div v-for="s in sites" :key="s.id" class="list-item">
        <!-- 行容器，class 为 row，横向排列寄存点名称与距离标签 -->

        <div class="row">
          <!-- 加粗元素，插值 s.name 显示寄存点名称 -->

          <strong>{{ s.name }}</strong>
          <!-- 行内标签元素，class 为 tag，显示距离：s.distance 非空时拼接“m”，否则显示“附近” -->

          <span class="tag">{{ s.distance != null ? s.distance + 'm' : '附近' }}</span>
        <!-- 闭合寄存点行容器 div 元素 -->

        </div>
        <!-- 段落元素，class 为 muted，插值 s.address 显示寄存点地址 -->

        <p class="muted">{{ s.address }}</p>
        <!-- 选择寄存按钮，点击触发 goDeposit(s) 进入该寄存点的寄存流程 -->

        <button @click="goDeposit(s)">选择寄存</button>
      <!-- 闭合当前寄存点条目 div 元素 -->

      </div>
      <!-- 空状态提示段落，当 sites 为空且非加载中时渲染，class 为 muted，显示“暂无附近寄存点” -->

      <p v-if="!sites.length && !loading" class="muted">暂无附近寄存点</p>
    <!-- 闭合附近寄存点卡片 div 元素 -->

    </div>

    <!-- 寄存指南卡片，class 为 card，v-if 依据 guide 是否有值决定是否渲染 -->

    <div v-if="guide" class="card">
      <!-- 二级标题，显示“寄存指南” -->

      <h2>寄存指南</h2>
      <!-- 段落元素，插值 guide.content 显示指南正文 -->

      <p>{{ guide.content }}</p>
      <!-- 段落元素，class 为 muted，显示指南版本 guide.version 与语言 guide.lang -->

      <p class="muted">版本 {{ guide.version }} / {{ guide.lang }}</p>
    <!-- 闭合寄存指南卡片 div 元素 -->

    </div>
    <!-- ===== 补齐入口：更多功能跳转（只新增，不改动上方原有卡片） ===== -->
    <!-- 更多功能卡片容器，class 为 card，放置首页相关原型界面的跳转入口 -->
    <div class="card">
      <!-- 二级标题，显示“更多功能” -->
      <h2>更多功能</h2>
      <!-- 行容器，class 为 row，横向排列各功能入口链接 -->
      <div class="row">
        <!-- 链接：选择城市，跳转 /city -->
        <router-link to="/city" class="tag" style="text-decoration:none">选择城市</router-link>
        <!-- 链接：附近寄存点，跳转 /nearby -->
        <router-link to="/nearby" class="tag" style="text-decoration:none">附近寄存点</router-link>
        <!-- 链接：寄存指南，跳转 /guide -->
        <router-link to="/guide" class="tag" style="text-decoration:none">寄存指南</router-link>
      <!-- 闭合入口链接行容器 div 元素 -->
      </div>
    <!-- 闭合更多功能卡片 div 元素 -->
    </div>

    <!-- 演示提示段落，当 err 有值时渲染，class 为 muted，插值 err 显示友好提示信息 -->

    <p v-if="err" class="muted">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
