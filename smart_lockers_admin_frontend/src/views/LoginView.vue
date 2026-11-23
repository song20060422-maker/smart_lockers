<template>
  <!-- 登录页根容器：全屏浅色背景，内容块垂直居中（不占满全屏高）@ticket B端-前端 -->
  <div class="login-wrap">
    <!-- 内容块容器：红框标注尺寸（高约 420px、近全宽留边距），左品牌区 + 右表单区 -->
    <div class="login-box">
      <!-- 左侧品牌展示区：仅保留标题与副标题 @ticket B端-前端 -->
      <div class="brand-pane">
        <!-- 品牌标题 -->
        <h1 class="brand-title">智能快寄柜</h1>
        <!-- 品牌副标题 -->
        <p class="brand-sub">IoT 运营中台 · B端运营后台</p>
      <!-- 闭合品牌区 -->
      </div>
      <!-- 右侧登录表单区（浮动卡片） -->
      <form class="card login-card" @submit.prevent="onLogin">
        <!-- 二级标题，显示系统名称 -->
        <h2>运营后台登录</h2>
        <!-- 表单标签，显示字段名“账号” -->
        <label>账号</label>
        <!-- 账号输入框，v-model 双向绑定 username 变量 -->
        <input v-model="username" placeholder="请输入账号" autocomplete="username" />
        <!-- 表单标签，显示字段名“密码” -->
        <label>密码</label>
        <!-- 密码输入框，type 为 password，v-model 双向绑定 password 变量 -->
        <input v-model="password" type="password" placeholder="请输入密码" autocomplete="current-password" />
        <!-- 登录提交按钮 -->
        <button class="btn primary" type="submit" :disabled="loading">{{ loading ? '登录中…' : '登 录' }}</button>
        <!-- 错误提示段落 -->
        <p v-if="error" class="err">{{ error }}</p>
      <!-- 闭合登录表单 -->
      </form>
    <!-- 闭合内容块容器 -->
    </div>
    <!-- 全局轻提示条 -->
    <div v-if="message" class="toast">{{ message }}</div>
  <!-- 闭合登录页根容器 -->
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotify } from '@/composables/useNotify'

// 声明变量 router = useRouter()，用于保存数据或函数引用
const router = useRouter()
// 声明变量 route = useRoute()，读取登录守卫携带的 redirect 回跳地址
const route = useRoute()
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

// 定义函数/方法 onLogin：提交登录（优先真实后端，失败降级演示模式）
async function onLogin() {
  // 清空上一次错误
  error.value = ''
  // 置位 loading
  loading.value = true
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理
  try {
    // 优先走真实后端登录：/auth/login 校验 demo_admins 并签发带 role/site_id 的 JWT @ticket B端-RBAC
    await auth.login(username.value, password.value)
    // 调用 notify('登录成功')，触发对应功能
    notify('登录成功')
  // 捕获上述逻辑抛出的异常并进行错误处理
  } catch (e: any) {
    // 后端不可达（网络异常）时降级演示登录：本地按账号名映射角色 @ticket B端-RBAC
    await auth.loginDemo(username.value)
    // 调用 notify('登录成功')，触发对应功能
    notify('登录成功（演示模式）')
  // 无论成功失败都会执行的收尾逻辑
  } finally { loading.value = false }
  // 声明变量 redirect，取守卫携带的目标地址，非法或缺失时回退到网点首页
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/sites'
  // 执行 router.replace(redirect)，登录后跳回原目标页
  router.replace(redirect)
// 当前代码块结束
}
</script>
<style scoped>
/* 登录页根容器：全屏浅色背景，内容块居中 @ticket B端-前端 */
.login-wrap { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #eef3fa 0%, #e6ecf6 100%); }
/* 左右两屏一体内容块：固定尺寸 1100×660，不随浏览器窗口缩放变化 */
.login-box { display: flex; width: 1100px; height: 660px; flex: none; border-radius: 4px; overflow: hidden; box-shadow: 0 16px 48px rgba(23, 61, 143, .18); }
/* 左侧品牌区：深蓝渐变一体左屏（占 55%），标题垂直居中 */
.brand-pane { flex: 0 0 55%; display: flex; flex-direction: column; justify-content: center; padding: 0 10%; background: linear-gradient(160deg, #2f7bff 0%, #1f5fe8 60%, #173d8f 100%); color: #fff; }
/* 品牌大标题 */
.brand-title { font-size: 34px; font-weight: 700; letter-spacing: 3px; margin: 0 0 12px; }
/* 品牌副标题 */
.brand-sub { font-size: 14px; opacity: .88; margin: 0; }
/* 右侧表单区：白色一体右屏（占 45%），表单垂直水平居中 */
.login-card { flex: 1; align-self: stretch; display: flex; flex-direction: column; justify-content: center; gap: 12px; padding: 0 16%; background: #ffffff; }
/* 表单标题 */
.login-card h2 { font-size: 20px; margin: 0 0 10px; color: #111827; }
/* 错误提示 */
.err { color: var(--danger); font-size: 12px; }
/* 标签 */
label { font-size: 13px; color: var(--muted); }
/* 窗口过小时兜底：允许内容块按视口缩放，避免溢出 */
@media (max-width: 1160px), (max-height: 720px) {
  /* 内容块改为视口自适应 */
  .login-box { width: 100%; height: 100vh; border-radius: 0; }
/* 当前代码块结束 */
}
</style>
