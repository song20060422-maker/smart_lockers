// @ticket B端-前端：登录态与角色权限状态（Pinia）
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api'
import { setToken, getToken } from '@/api/request'
import type { AdminUser, AdminRole } from '@/types/api'
import { PERMISSIONS } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AdminUser | null>(null)
  const token = ref<string | null>(getToken())

  const role = computed<AdminRole>(() => user.value?.role || 'user')
  const permissions = computed(() => PERMISSIONS[role.value] || [])
  const loggedIn = computed(() => !!token.value)

  function can(perm: string) { return permissions.value.includes(perm) }

  async function login(username: string, password: string) {
    const res = await api.login(username, password)
    token.value = res.data?.token || null
    setToken(token.value)
    user.value = res.data?.user || null
    if (!user.value) await fetchMe()
    return user.value
  }

  async function fetchMe() {
    if (!token.value) return null
    const res = await api.me()
    user.value = res.data || null
    return user.value
  }

  function logout() {
    token.value = null
    user.value = null
    setToken(null)
  }

  return { user, token, role, permissions, loggedIn, can, login, fetchMe, logout }
})
