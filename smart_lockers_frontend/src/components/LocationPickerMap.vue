// @ticket 03-首页：地图选点补缺——Leaflet + OpenStreetMap 实现无需 API Key 的地点选择
// @proto 原型对齐：首页定位页新增“在地图上选点”，点选后回填当前城市与经纬度
<script setup lang="ts">
// 引入 Vue 组合式 API：ref/onMounted/onBeforeUnmount 管理组件生命周期与状态
import { ref, onMounted, onBeforeUnmount } from 'vue'
// 引入 Leaflet 命名空间，用于创建地图、图层与标记
import * as L from 'leaflet'
// 引入 Leaflet 基础样式，保证瓦片与控件正常显示
import 'leaflet/dist/leaflet.css'
// 引入全国城市数据，支持按城市名/拼音跳转到地图中心
import { nationwideCities } from '@/data/cities'

// 定义父组件传入的初始定位属性：经纬度与城市名
const props = defineProps<{
  // 初始纬度，0 表示未定位
  initialLat: number
  // 初始经度，0 表示未定位
  initialLng: number
  // 初始城市名，仅用于提示
  initialCity?: string
}>()

// 声明组件对外抛出的事件：pick 回填选点、close 关闭地图
// 选点结果结构 PickResult 由父组件通过事件参数接收
// 定义选点结果结构，含经纬度、详细地址与城市名
interface PickResult {
  // 所选地点纬度
  lat: number
  // 所选地点经度
  lng: number
  // 反向地理编码得到的详细地址
  address: string
  // 反向地理编码得到的城市名，可能为空
  city: string
}

const emit = defineEmits<{
  // 用户在地图上点选并确认后触发
  (e: 'pick', value: PickResult): void
  // 用户点击关闭地图时触发
  (e: 'close'): void
}>()

// 地图容器 DOM 引用，挂载后用于初始化 Leaflet
const mapEl = ref<HTMLDivElement | null>(null)
// 地图实例引用，onBeforeUnmount 时销毁以释放监听
let map: L.Map | null = null
// 选点标记引用，重复点击时移动而非新建
let marker: L.CircleMarker | null = null
// 当前选中点纬度，0 表示尚未选点
const pickedLat = ref(0)
// 当前选中点经度，0 表示尚未选点
const pickedLng = ref(0)
// 反向地理编码得到的地址文案
const addressText = ref('')
// 反向地理编码得到的城市名
const cityText = ref('')
// 地址解析中的标志，用于禁用确认按钮
const resolving = ref(false)
// 城市跳转输入框绑定值
const cityQuery = ref('')
// 记录最近一次地址解析时间戳，用于丢弃过期请求结果
let lastReverseAt = 0

// 判断经纬度是否有效：非 0 且落在合理范围内
function valid(lat: number, lng: number): boolean {
  // 纬度需在 -90 到 90 之间
  if (lat <= -90 || lat >= 90) return false
  // 经度需在 -180 到 180 之间
  if (lng <= -180 || lng >= 180) return false
  // 两者均非 0 才视为有效定位
  return lat !== 0 && lng !== 0
}

// 调用 Nominatim 反向地理编码，将经纬度解析为地址与城市（无需 API Key）
async function reverseGeocode(lat: number, lng: number): Promise<{ city: string; address: string }> {
  // 构造 OpenStreetMap Nominatim reverse 接口地址，指定中文与 jsonv2 格式
  const url = `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${lat}&lon=${lng}&accept-language=zh-CN&zoom=16`
  // 网络异常时返回空结果，由调用方用坐标兜底
  try {
    // 发起跨域 GET 请求，Nominatim 允许浏览器直接访问
    const res = await fetch(url, { headers: { Accept: 'application/json' } })
    // 非 2xx 响应视为失败
    if (!res.ok) return { city: '', address: '' }
    // 解析返回 JSON
    const data = await res.json()
    // 取地址对象，可能为空
    const a = data.address || {}
    // 城市名优先取 city，其次 state/province
    const city = a.city || a.province || a.state || ''
    // 详细地址取 display_name
    const address = data.display_name || ''
    // 返回解析结果
    return { city, address }
  // 捕获网络或解析异常
  } catch (_) {
    // 失败返回空结果
    return { city: '', address: '' }
  }
}

// 选中给定坐标：更新标记、发起地址解析
async function selectAt(lat: number, lng: number): Promise<void> {
  // 写入选中坐标
  pickedLat.value = lat
  // 写入选中经度
  pickedLng.value = lng
  // 标记已存在则移动到新位置
  if (marker) marker.setLatLng([lat, lng])
  // 标记不存在则创建一个蓝色圆点标记
  else if (map) marker = L.circleMarker([lat, lng], { radius: 8, color: '#2563eb', fillColor: '#2563eb', fillOpacity: 0.6 }).addTo(map)
  // 进入解析中状态
  resolving.value = true
  // 先用坐标文案兜底，解析成功后再替换
  addressText.value = `坐标 ${lat.toFixed(4)}, ${lng.toFixed(4)}`
  // 清空旧城市名
  cityText.value = ''
  // 记录本次解析时间戳，用于丢弃过期请求结果
  const stamp = Date.now()
  // 最近解析时间更新为当前戳
  lastReverseAt = stamp
  // 调用反向地理编码
  const r = await reverseGeocode(lat, lng)
  // 若在等待期间用户又点了新点，丢弃旧结果
  if (lastReverseAt !== stamp) return
  // 有详细地址则覆盖坐标兜底文案
  if (r.address) addressText.value = r.address
  // 保存城市名
  cityText.value = r.city
  // 结束解析中状态
  resolving.value = false
}

// 地图点击事件处理函数，参数为 Leaflet 鼠标事件
function onMapClick(ev: L.LeafletMouseEvent): void {
  // 以点击位置执行选点流程
  void selectAt(ev.latlng.lat, ev.latlng.lng)
}

// 按输入关键字在全国城市中查找并跳转到该城市中心
function jumpToCity(): void {
  // 去除首尾空白
  const k = cityQuery.value.trim()
  // 空输入直接返回
  if (!k) return
  // 转小写便于拼音匹配
  const lower = k.toLowerCase()
  // 在内置全国数据中按中文名/全拼/首字母精确匹配
  const found = nationwideCities.find((c) => c.name === k || c.pinyin === lower || c.initial.toLowerCase() === lower)
  // 未匹配到任何城市时给出提示
  if (!found) {
    // 提示用户当前支持的匹配方式
    addressText.value = `未找到“${k}”，请输入城市中文名或拼音`
    // 结束
    return
  }
  // 地图存在则将视图移动到该城市中心
  if (map) map.setView([found.lat, found.lng], 12)
  // 以该城市中心点执行选点
  void selectAt(found.lat, found.lng)
}

// 确认回填：将选点结果抛给父组件
function confirmPick(): void {
  // 未选点则不允许确认
  if (!pickedLat.value || !pickedLng.value) return
  // 解析中禁用确认，避免地址尚未就绪
  if (resolving.value) return
  // 抛出选点结果
  emit('pick', { lat: pickedLat.value, lng: pickedLng.value, address: addressText.value, city: cityText.value })
}

// 组件挂载后初始化地图
onMounted(() => {
  // 容器未就绪则跳过初始化
  if (!mapEl.value) return
  // 计算初始中心：入参有效则用它，否则回落到北京
  const lat = valid(props.initialLat, props.initialLng) ? props.initialLat : 39.9042
  // 计算初始经度
  const lng = valid(props.initialLat, props.initialLng) ? props.initialLng : 116.4074
  // 创建地图实例，设置中心与缩放级别
  map = L.map(mapEl.value, { center: [lat, lng], zoom: 12 })
  // 添加 OpenStreetMap 瓦片图层，无需 API Key
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    // 遵循 OSM 署名要求
    attribution: '&copy; OpenStreetMap contributors',
    // 最大缩放级别
    maxZoom: 19,
  }).addTo(map)
  // 绑定地图点击事件
  map.on('click', onMapClick)
  // 若已有有效定位，则在初始位置放置标记并解析地址
  if (valid(props.initialLat, props.initialLng)) void selectAt(lat, lng)
  // 容器可能是动态渲染，下一帧刷新尺寸确保瓦片正确铺满
  setTimeout(() => map && map.invalidateSize(), 0)
})

// 组件卸载前销毁地图，避免内存泄漏与重复容器错误
onBeforeUnmount(() => {
  // 存在地图实例则移除
  if (map) map.remove()
  // 清空引用
  map = null
  // 清空标记引用
  marker = null
})
</script>

<template>
  <!-- 地图选点面板根容器，纵向排列城市跳转、地图与操作区 -->
  <div class="map-picker">
    <!-- 城市跳转行：输入城市名或拼音，快速定位到该城市中心 -->
    <div class="map-row">
      <!-- 城市输入框，回车触发 jumpToCity -->
      <input v-model="cityQuery" placeholder="输入城市名/拼音定位" @keyup.enter="jumpToCity" />
      <!-- 定位到城市按钮 -->
      <button class="btn-ghost" @click="jumpToCity">定位城市</button>
      <!-- 关闭地图按钮 -->
      <button class="btn-ghost" @click="emit('close')">关闭地图</button>
    </div>
    <!-- Leaflet 地图容器，ref 用于初始化 -->
    <div ref="mapEl" class="map-canvas"></div>
    <!-- 地址提示段落：解析中显示“地址解析中”，否则显示所选地址或操作引导 -->
    <p class="muted">{{ resolving ? '地址解析中...' : (addressText || '请在地图上点击选择具体地点') }}</p>
    <!-- 操作行：确认回填与重新提示 -->
    <div class="map-row">
      <!-- 确认回填按钮，未选点或解析中禁用 -->
      <button :disabled="!pickedLat || resolving" @click="confirmPick">确认回填</button>
      <!-- 说明文案，提示数据来源无需 Key -->
      <span class="muted">地图数据来自 OpenStreetMap，无需 API Key</span>
    </div>
  </div>
</template>

<style scoped>
/* 地图选点面板：浅色边框卡片，与其他卡片风格一致 */
.map-picker { display: flex; flex-direction: column; gap: 10px; padding: 12px; border: 1px solid #dbeafe; border-radius: 10px; background: #fff; }
/* 操作行：横向排列，元素间留 8px 间距 */
.map-row { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
/* 城市输入框占满剩余宽度，最小宽度 140px */
.map-row input { flex: 1; min-width: 140px; height: 36px; box-sizing: border-box; border: 1px solid #d1d5db; border-radius: 6px; padding: 0 10px; font-size: 14px; }
/* 地图画布固定高度，保证可见面积 */
.map-canvas { width: 100%; height: 300px; border-radius: 8px; overflow: hidden; background: #e5e7eb; }
</style>
