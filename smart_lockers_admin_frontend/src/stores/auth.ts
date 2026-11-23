// @ticket B端-前端：登录态与角色权限状态（Pinia）
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api'
import { setToken } from '@/api/request'
import type { AdminUser, AdminRole } from '@/types/api'
import { PERMISSIONS } from '@/types/api'

// 演示登录会话在 localStorage 中的键，与真实 admin_token 隔离，避免旧令牌残留导致打开即进后台
const DEMO_SESSION_KEY = 'admin_demo_session'
// 演示会话结构：本地模拟的 token 与用户信息，全程不依赖后端凭证
interface DemoSession { token: string; user: AdminUser }

// 从 localStorage 读取演示会话，解析失败或不存在时返回 null，保证守卫判定安全
function readDemoSession(): DemoSession | null {
  try {
    // 读取原始字符串，空值直接视为未登录
    const raw = localStorage.getItem(DEMO_SESSION_KEY)
    // 无会话记录时返回 null
    if (!raw) return null
    // 解析为演示会话对象并返回
    return JSON.parse(raw) as DemoSession
  } catch {
    // 解析异常一律按未登录处理，避免脏数据阻塞登录页
    return null
  }
}

// 导出 const useAuthStore = defineStore('auth', () => {，供其他模块引用

export const useAuthStore = defineStore('auth', () => {
  // 声明变量 session = readDemoSession()，仅在本地演示会话存在时回填登录态
  const session = readDemoSession()
  // 声明变量 user = ref<AdminUser | null>(session?.user || null)，演示会话用户或空
  const user = ref<AdminUser | null>(session?.user || null)
  // 声明变量 token = ref<string | null>(session?.token || null)，演示会话令牌或空
  const token = ref<string | null>(session?.token || null)

  // 并发去重：缓存进行中的 /auth/me 请求，避免多次路由跳转重复回填

  let meRequest: Promise<AdminUser | null> | null = null

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

  // 定义函数/方法，封装可复用的处理逻辑：演示登录——本地生成 token 与用户，无需真实凭证与后端
  // 角色按账号名映射（与后端 demo_admins 三账号对齐），保证不同角色登录后菜单/页面权限隔离 @ticket B端-RBAC
  async function loginDemo(username: string) {
    // 声明变量 demoName，去掉首尾空格，空则回退到默认演示账号
    const demoName = username.trim() || 'site_admin'
    // 角色映射：operator/admin→平台运营，user→普通用户，site/site_admin→网点管理员（兜底 site_admin）
    const demoRole: AdminRole = /operator|^admin$/i.test(demoName)
      ? 'operator'
      : /^user$/i.test(demoName)
        ? 'user'
        : 'site_admin'
    // 普通用户与网点管理员的数据域网点编号（与后端 demo_admins 的 site_id 对齐）
    const demoSiteId = demoRole === 'site_admin' ? 1001 : 0
    // 声明变量 demoToken，本地拼接时间戳与随机串生成模拟令牌
    const demoToken = `demo-${Date.now()}-${Math.random().toString(36).slice(2)}`
    // 声明变量 demoUser，按映射后的角色构造演示用户信息
    const demoUser: AdminUser = {
      // 定义对象属性 id，演示用户编号
      id: 1,
      // 定义对象属性 username，取登录框输入的账号
      username: demoName,
      // 定义对象属性 role，按账号名映射得到角色
      role: demoRole,
      // 定义对象属性 site_id，演示账号所属网点编号
      site_id: demoSiteId,
      // 定义对象属性 status，演示账号为启用状态
      status: 1,
    }
    // 赋值更新：token.value 设为本地模拟令牌
    token.value = demoToken
    // 调用 setToken(token.value)，同步到请求层用于 Bearer 注入
    setToken(token.value)
    // 赋值更新：user.value 设为演示用户
    user.value = demoUser
    // 调用 localStorage.setItem，持久化演示会话以支持刷新后保持登录态
    localStorage.setItem(DEMO_SESSION_KEY, JSON.stringify({ token: demoToken, user: demoUser }))
    // 返回 user.value 给调用方
    return user.value
  // 当前代码块结束

  }

  // 定义函数/方法，封装可复用的处理逻辑：async function fetchMe() {

  async function fetchMe() {
    // 条件判断：无 token 时直接返回 null，避免无意义请求

    if (!token.value) return null
    // 条件判断：已有用户信息时直接返回，避免每次跳转重复回填

    if (user.value) return user.value
    // 条件判断：已有进行中的回填请求时直接复用，处理并发重复调用

    if (meRequest) return meRequest
    // 发起 /auth/me 请求并缓存 Promise，finally 确保无论成败都清空缓存，避免失败后永久锁定

    meRequest = api.me().then((res) => {
      // 赋值更新：将接口返回的用户信息回填到 user，供顶栏读取真实用户名与角色

      user.value = res.data || null
      // 返回 user.value 给调用方

      return user.value
    }).finally(() => {
      // 赋值更新：请求结束后清空缓存，允许后续需要时重新拉取

      meRequest = null
    })
    // 返回 meRequest 给调用方

    return meRequest
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
    // 调用 localStorage.removeItem，清除持久化的演示会话，退出后刷新将回到登录页

    localStorage.removeItem(DEMO_SESSION_KEY)
  // 当前代码块结束

  }

  // 返回 { user, token, role, permissions, loggedIn, can, login, loginDemo, fetchMe, logout } 给调用方

  return { user, token, role, permissions, loggedIn, can, login, loginDemo, fetchMe, logout }
// 执行 })，实现对应的前端逻辑

})
