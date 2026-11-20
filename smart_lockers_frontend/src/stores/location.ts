// @ticket 03-首页：定位与当前城市 Pinia store
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { City, LocateResult } from '@/types'

// 导出 const useLocationStore = defineStore('location', () => {，供其他模块引用

export const useLocationStore = defineStore('location', () => {
  // 声明变量 lat = ref<number>(0)，用于保存数据或函数引用

  const lat = ref<number>(0)
  // 声明变量 lng = ref<number>(0)，用于保存数据或函数引用

  const lng = ref<number>(0)
  // 声明变量 city = ref<string>('')，用于保存数据或函数引用

  const city = ref<string>('')
  // 声明变量 source = ref<string>('')，用于保存数据或函数引用

  const source = ref<string>('')

  // @ticket 03-首页：写入定位结果
  // 定义函数/方法，封装可复用的处理逻辑：function applyLocate(r: LocateResult) {

  function applyLocate(r: LocateResult) {
    // 赋值更新：lat.value = r.lat

    lat.value = r.lat
    // 赋值更新：lng.value = r.lng

    lng.value = r.lng
    // 赋值更新：city.value = r.city || city.value

    city.value = r.city || city.value
    // 赋值更新：source.value = r.source

    source.value = r.source
  // 当前代码块结束

  }

  // @ticket 03-首页：切换城市
  // 定义函数/方法，封装可复用的处理逻辑：function setCity(c: City) {

  function setCity(c: City) {
    // 赋值更新：city.value = c.name

    city.value = c.name
  // 当前代码块结束

  }

  // 返回 { lat, lng, city, source, applyLocate, setCity } 给调用方

  return { lat, lng, city, source, applyLocate, setCity }
// 执行 })，实现对应的前端逻辑

})
