// @ticket 06-个人中心：登录态 Pinia store（token + user_id）
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as apiLogin, logout as apiLogout, demoLogin as apiDemoLogin } from '@/api/auth'
import { getToken, setToken, clearToken } from '@/api/request'

// 导出 const useAuthStore = defineStore('auth', () => {，供其他模块引用

export const useAuthStore = defineStore('auth', () => {
  // 声明变量 token = ref<string>(getToken())，用于保存数据或函数引用

  const token = ref<string>(getToken())
  // 声明变量 userId = ref<number>(0)，用于保存数据或函数引用

  const userId = ref<number>(0)
  // 声明变量 logged = ref<boolean>(!!token.value)，用于保存数据或函数引用

  const logged = ref<boolean>(!!token.value)

  // @ticket 06-个人中心：登录并持久化 token
  // 定义函数/方法，封装可复用的处理逻辑：async function login(uid: number) {

  async function login(uid: number) {
    // 声明变量 res = await apiLogin(uid)，用于保存数据或函数引用

    const res = await apiLogin(uid)
    // 调用 setToken(res.token)，触发对应功能

    setToken(res.token)
    // 赋值更新：token.value = res.token

    token.value = res.token
    // 赋值更新：userId.value = uid

    userId.value = uid
    // 赋值更新：logged.value = true

    logged.value = true
    // 返回 res 给调用方

    return res
  // 当前代码块结束

  }

  // @ticket 06-个人中心：退出登录清除 token
  // 定义函数/方法，封装可复用的处理逻辑：async function logout() {

  async function logout() {
    // 等待异步操作完成：try { await apiLogout() } catch (_) { /* 忽略退出接口异常 */ }

    try { await apiLogout() } catch (_) { /* 忽略退出接口异常 */ }
    // 调用 clearToken()，触发对应功能

    clearToken()
    // 赋值更新：token.value = ''

    token.value = ''
    // 赋值更新：logged.value = false

    logged.value = false
    // 赋值更新：userId.value = 0

    userId.value = 0
  // 当前代码块结束

  }

  // 修复偏差4：演示登录改为调用后端真实接口签发 JWT，禁止本地生成 token @ticket 00-project
  // 定义函数 loginDemo：接收手机号与验证码，走真实 HTTP 登录闭环

  async function loginDemo(phone: string, code: string) {
    // 调用后端演示登录接口，取回真实 token 与 user_id

    const res = await apiDemoLogin(phone, code)
    // 持久化后端签发的真实 token

    setToken(res.token)
    // 赋值更新：token.value = res.token

    token.value = res.token
    // 赋值更新：userId.value = res.user_id

    userId.value = res.user_id
    // 赋值更新：logged.value = true

    logged.value = true
    // 返回 res 给调用方

    return res
  // 当前代码块结束

  }

  // 兼容语义化命名：loginByPhone 等价于真实演示登录 @ticket 00-project
  // 定义函数 loginByPhone：转调 loginDemo

  async function loginByPhone(phone: string, code: string) {
    // 返回 loginDemo(phone, code) 给调用方

    return loginDemo(phone, code)
  // 当前代码块结束

  }

  // 返回 { token, userId, logged, login, logout, loginDemo } 给调用方

  return { token, userId, logged, login, logout, loginDemo, loginByPhone }
// 执行 })，实现对应的前端逻辑

})
