// @ticket 03-首页：定位与当前城市 Pinia store
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { City, LocateResult } from '@/types'

export const useLocationStore = defineStore('location', () => {
  const lat = ref<number>(0)
  const lng = ref<number>(0)
  const city = ref<string>('')
  const source = ref<string>('')

  // @ticket 03-首页：写入定位结果
  function applyLocate(r: LocateResult) {
    lat.value = r.lat
    lng.value = r.lng
    city.value = r.city || city.value
    source.value = r.source
  }

  // @ticket 03-首页：切换城市
  function setCity(c: City) {
    city.value = c.name
  }

  return { lat, lng, city, source, applyLocate, setCity }
})
