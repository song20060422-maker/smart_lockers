// @ticket 06-个人中心：登录态 Pinia store（token + user_id）
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as apiLogin, logout as apiLogout } from '@/api/auth'
import { getToken, setToken, clearToken } from '@/api/request'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(getToken())
  const userId = ref<number>(0)
  const logged = ref<boolean>(!!token.value)

  // @ticket 06-个人中心：登录并持久化 token
  async function login(uid: number) {
    const res = await apiLogin(uid)
    setToken(res.token)
    token.value = res.token
    userId.value = uid
    logged.value = true
    return res
  }

  // @ticket 06-个人中心：退出登录清除 token
  async function logout() {
    try { await apiLogout() } catch (_) { /* 忽略退出接口异常 */ }
    clearToken()
    token.value = ''
    logged.value = false
    userId.value = 0
  }

  return { token, userId, logged, login, logout }
})
