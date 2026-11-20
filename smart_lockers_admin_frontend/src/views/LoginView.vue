<template>
  <!-- 登录页根容器，class 为 login-wrap，居中包裹登录卡片 -->

  <div class="login-wrap">
    <!-- 登录表单，class 为 card login-card，提交事件被 prevent 修饰阻止默认行为并调用 onLogin 方法 -->

    <form class="card login-card" @submit.prevent="onLogin">
      <!-- 二级标题，显示系统名称“智能快寄柜 · B端运营后台” -->

      <h2>智能快寄柜 · B端运营后台</h2>
      <!-- 段落元素，class 为 sub，显示登录提示说明文字 -->

      <p class="sub">请使用分配的运营账号登录</p>
      <!-- 表单标签，显示字段名“账号” -->

      <label>账号</label>
      <!-- 账号输入框，v-model 双向绑定 username 变量，占位提示为 operator，自动补全类型为 username -->

      <input v-model="username" placeholder="operator" autocomplete="username" />
      <!-- 表单标签，显示字段名“密码” -->

      <label>密码</label>
      <!-- 密码输入框，type 为 password，v-model 双向绑定 password 变量，占位提示为圆点，自动补全类型为 current-password -->

      <input v-model="password" type="password" placeholder="••••••" autocomplete="current-password" />
      <!-- 登录提交按钮，class 为 btn primary，type 为 submit，loading 为真时禁用，按钮文本在“登录中…”与“登录”间切换 -->

      <button class="btn primary" type="submit" :disabled="loading">{{ loading ? '登录中…' : '登录' }}</button>
      <!-- 错误提示段落，当 error 有值时渲染，class 为 err，插值 error 显示错误信息 -->

      <p v-if="error" class="err">{{ error }}</p>
    <!-- 闭合登录表单 form 元素 -->

    </form>
    <!-- 全局轻提示条，v-if 依据 message 是否有值决定是否渲染，class 为 toast，插值 message 显示提示内容 -->

    <div v-if="message" class="toast">{{ message }}</div>
  <!-- 闭合登录页根容器 div 元素 -->

  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'

// 声明变量 router = useRouter()，用于保存数据或函数引用

const router = useRouter()
// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 赋值更新：const { message, notify } = useNotify()

const { message, notify } = useNotify()
// 声明变量 username = ref('')，用于保存数据或函数引用

const username = ref('')
// 声明变量 password = ref('')，用于保存数据或函数引用

const password = ref('')
// 声明变量 loading = ref(false)，用于保存数据或函数引用

const loading = ref(false)
// 声明变量 error = ref('')，用于保存数据或函数引用

const error = ref('')
// 定义函数/方法，封装可复用的处理逻辑：async function onLogin() {

async function onLogin() {
  // 赋值更新：error.value = ''

  error.value = ''
  // 赋值更新：loading.value = true

  loading.value = true
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 等待异步操作完成：await auth.login(username.value.trim(), password.value)

    await auth.login(username.value.trim(), password.value)
    // 调用 notify('登录成功')，触发对应功能

    notify('登录成功')
    // 执行 router.replace('/sites')，实现对应的前端逻辑

    router.replace('/sites')
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e: any) {
    // 赋值更新：error.value = e?.message || '登录失败，请检查账号密码'

    error.value = e?.message || '登录失败，请检查账号密码'
  // 无论成功失败都会执行的收尾逻辑

  } finally { loading.value = false }
// 当前代码块结束

}
</script>
<style scoped>
.login-wrap { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #1f6feb22, #1d243022); }
.login-card { width: 360px; display: flex; flex-direction: column; gap: 12px; }
.login-card h2 { font-size: 18px; }
.sub, .err { color: var(--muted); font-size: 12px; }
.err { color: var(--danger); }
label { font-size: 13px; color: var(--muted); }
</style>
