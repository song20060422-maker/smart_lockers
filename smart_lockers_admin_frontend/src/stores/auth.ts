// @ticket B端-前端：登录态与角色权限状态（Pinia）
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api'
import { setToken, getToken } from '@/api/request'
import type { AdminUser, AdminRole } from '@/types/api'
import { PERMISSIONS } from '@/types/api'

// 导出 const useAuthStore = defineStore('auth', () => {，供其他模块引用

export const useAuthStore = defineStore('auth', () => {
  // 声明变量 user = ref<AdminUser | null>(null)，用于保存数据或函数引用

  const user = ref<AdminUser | null>(null)
  // 声明变量 token = ref<string | null>(getToken())，用于保存数据或函数引用

  const token = ref<string | null>(getToken())

  // 声明变量 role = computed<AdminRole>(() => user.value?.role || 'user')，用于保存数据或函数引用

  const role = computed<AdminRole>(() => user.value?.role || 'user')
  // 声明变量 permissions = computed(() => PERMISSIONS[role.value] || [])，用于保存数据或函数引用

  const permissions = computed(() => PERMISSIONS[role.value] || [])
  // 声明变量 loggedIn = computed(() => !!token.value)，用于保存数据或函数引用

  const loggedIn = computed(() => !!token.value)

  // 定义函数/方法，封装可复用的处理逻辑：function can(perm: string) { return permissions.value.includes(perm) }

  function can(perm: string) { return permissions.value.includes(perm) }

  // 定义函数/方法，封装可复用的处理逻辑：async function login(username: string, password: string) {

  async function login(username: string, password: string) {
    // 声明变量 res = await api.login(username, password)，用于保存数据或函数引用

    const res = await api.login(username, password)
    // 赋值更新：token.value = res.data?.token || null

    token.value = res.data?.token || null
    // 调用 setToken(token.value)，触发对应功能

    setToken(token.value)
    // 赋值更新：user.value = res.data?.user || null

    user.value = res.data?.user || null
    // 条件判断：if (!user.value) await fetchMe() 时执行下面逻辑

    if (!user.value) await fetchMe()
    // 返回 user.value 给调用方

    return user.value
  // 当前代码块结束

  }

  // 定义函数/方法，封装可复用的处理逻辑：async function fetchMe() {

  async function fetchMe() {
    // 条件判断：if (!token.value) return null 时执行下面逻辑

    if (!token.value) return null
    // 声明变量 res = await api.me()，用于保存数据或函数引用

    const res = await api.me()
    // 赋值更新：user.value = res.data || null

    user.value = res.data || null
    // 返回 user.value 给调用方

    return user.value
  // 当前代码块结束

  }

  // 定义函数/方法，封装可复用的处理逻辑：function logout() {

  function logout() {
    // 赋值更新：token.value = null

    token.value = null
    // 赋值更新：user.value = null

    user.value = null
    // 调用 setToken(null)，触发对应功能

    setToken(null)
  // 当前代码块结束

  }

  // 返回 { user, token, role, permissions, loggedIn, can, login, fetchMe, logout } 给调用方

  return { user, token, role, permissions, loggedIn, can, login, fetchMe, logout }
// 执行 })，实现对应的前端逻辑

})
