<script setup lang="ts">
// @ticket 03-首页：定位/城市选择/附近寄存点/寄存指南
// @proto 原型对齐：v-sites 网点柜组监控；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { locate, searchCities, nearbySites, getGuide } from '@/api/home'
import { useLocationStore } from '@/stores/location'
import { useAuthStore } from '@/stores/auth'
import type { City, LockerSite, GuideResult } from '@/types'

const router = useRouter()
const locStore = useLocationStore()
const auth = useAuthStore()
const keyword = ref('')
const cities = ref<City[]>([])
const sites = ref<LockerSite[]>([])
const guide = ref<GuideResult | null>(null)
const err = ref('')
const loading = ref(false)

// @ticket 03-首页：初始化定位并加载附近寄存点与指南
onMounted(async () => {
  try {
    const r = await locate({ fallback_city: '北京' })
    locStore.applyLocate(r)
  } catch (_) { /* 定位失败用默认城市 */ }
  await Promise.all([loadSites(), loadGuide()])
  if (!auth.logged) auth.login(10001).catch(() => {})
})

async function loadSites() {
  loading.value = true
  try {
    const r = await nearbySites({ lat: locStore.lat, lng: locStore.lng, page: 1, page_size: 20 })
    sites.value = r.list
  } catch (e) { err.value = (e as Error).message }
  finally { loading.value = false }
}

async function doSearchCity() {
  try { cities.value = await searchCities(keyword.value) }
  catch (e) { err.value = (e as Error).message }
}

function pickCity(c: City) {
  locStore.setCity(c)
  cities.value = []
  keyword.value = ''
  loadSites()
}

async function loadGuide() {
  try { guide.value = await getGuide() }
  catch (_) { /* 指南非关键路径 */ }
}

function goDeposit(site: LockerSite) {
  router.push({ path: '/deposit', query: { site_id: String(site.id) } })
}
</script>

<template>
  <div class="col">
    <div class="card">
      <h2>首页定位</h2>
      <p class="muted">当前城市：{{ locStore.city || '定位中' }}（来源：{{ locStore.source || '-' }}）</p>
      <div class="row">
        <input v-model="keyword" placeholder="搜索城市" @keyup.enter="doSearchCity" />
        <button @click="doSearchCity">搜索城市</button>
      </div>
      <div v-if="cities.length" class="list-item">
        <button v-for="c in cities" :key="c.id" class="btn-ghost" @click="pickCity(c)">{{ c.name }}</button>
      </div>
    </div>

    <div class="card">
      <h2>附近寄存点</h2>
      <p v-if="loading" class="muted">加载中...</p>
      <div v-for="s in sites" :key="s.id" class="list-item">
        <div class="row">
          <strong>{{ s.name }}</strong>
          <span class="tag">{{ s.distance != null ? s.distance + 'm' : '附近' }}</span>
        </div>
        <p class="muted">{{ s.address }}</p>
        <button @click="goDeposit(s)">选择寄存</button>
      </div>
      <p v-if="!sites.length && !loading" class="muted">暂无附近寄存点</p>
    </div>

    <div v-if="guide" class="card">
      <h2>寄存指南</h2>
      <p>{{ guide.content }}</p>
      <p class="muted">版本 {{ guide.version }} / {{ guide.lang }}</p>
    </div>
    <p v-if="err" class="err">{{ err }}</p>
  </div>
</template>
