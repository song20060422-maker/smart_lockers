<script setup lang="ts">
// @ticket 00-project：登录页——手机号验证码登录 / 微信授权登录（演示，无需真实凭证）
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 声明变量 router = useRouter()，用于保存数据或函数引用

const router = useRouter()
// 声明变量 auth = useAuthStore()，用于保存数据或函数引用

const auth = useAuthStore()
// 声明变量 phone = ref('')，用于保存数据或函数引用

const phone = ref('')
// 声明变量 code = ref('')，用于保存数据或函数引用

const code = ref('')
// 声明变量 sent = ref(false)，用于保存数据或函数引用

const sent = ref(false)
// 声明变量 countdown = ref(0)，用于保存数据或函数引用

const countdown = ref(0)
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')
// 声明变量 loading = ref(false)，用于保存数据或函数引用

const loading = ref(false)

// @ticket 00-project：发送验证码（演示：任意 11 位手机号均可，倒计时模拟发送）
// 定义函数/方法，封装可复用的处理逻辑：function sendCode() {

function sendCode() {
  // 校验手机号是否为 11 位数字
  if (!/^1\d{10}$/.test(phone.value)) { err.value = '请输入正确的 11 位手机号'; return }
  // 清空错误信息

  err.value = ''
  // 标记验证码已发送

  sent.value = true
  // 启动 60 秒倒计时

  countdown.value = 60
  // 使用定时器每秒递减倒计时
  const timer = setInterval(() => {
    // 倒计时递减

    countdown.value--
    // 倒计时归零时清除定时器

    if (countdown.value <= 0) clearInterval(timer)
  // 定时器回调结束

  }, 1000)
// 当前代码块结束

}

// @ticket 00-project：手机号验证码登录（演示：验证码为 4~6 位数字即可）
// 定义函数/方法，封装可复用的处理逻辑：async function loginByPhone() {

async function loginByPhone() {
  // 校验手机号

  if (!/^1\d{10}$/.test(phone.value)) { err.value = '请输入正确的 11 位手机号'; return }
  // 校验验证码是否已发送

  if (!sent.value) { err.value = '请先获取验证码'; return }
  // 校验验证码格式

  if (!/^\d{4,6}$/.test(code.value)) { err.value = '请输入 4~6 位验证码'; return }
  // 清空错误信息

  err.value = ''
  // 开启加载状态

  loading.value = true
  // 调用演示登录（无需真实短信/后端）
  try {
    // 等待异步操作完成：await auth.loginByPhone(phone.value, code.value)

    await auth.loginByPhone(phone.value, code.value)
    // 登录成功跳转首页

    router.push('/home')
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
  // 无论成功失败都会执行的收尾逻辑

  finally { loading.value = false }
// 当前代码块结束

}

// @ticket 00-project：微信授权登录（演示：一键授权，无需真实微信凭证）
// 定义函数/方法，封装可复用的处理逻辑：async function loginByWechat() {

async function loginByWechat() {
  // 清空错误信息

  err.value = ''
  // 开启加载状态

  loading.value = true
  // 调用演示登录（模拟微信授权成功）
  try {
    // 等待异步操作完成：await auth.loginByPhone('13800138000', '123456')

    await auth.loginByPhone('13800138000', '123456')
    // 授权成功跳转首页

    router.push('/home')
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
  // 无论成功失败都会执行的收尾逻辑

  finally { loading.value = false }
// 当前代码块结束

}
</script>

<template>
  <!-- 登录页根容器，class 为 login-page，纵向居中排布品牌与登录表单 -->
  <div class="login-page">
    <!-- 品牌区域，class 为 brand-box，展示 Logo 与标题 -->
    <div class="brand-box">
      <!-- Logo 图标元素，class 为 logo，显示柜子图标 -->
      <div class="logo">▣</div>
      <!-- 应用标题元素，class 为 brand-name，显示“智能快寄柜” -->
      <div class="brand-name">智能快寄柜</div>
      <!-- 副标题元素，class 为 brand-sub，显示“寄存更轻松，取件更便捷” -->
      <div class="brand-sub">寄存更轻松，取件更便捷</div>
    <!-- 闭合品牌区域 div 元素 -->

    </div>
    <!-- 表单卡片，class 为 login-card，包含手机号验证码登录表单 -->
    <div class="login-card">
      <!-- 标题元素，class 为 card-title，显示“手机号登录” -->
      <div class="card-title">手机号登录</div>
      <!-- 手机号输入框，v-model 双向绑定 phone，类型 tel，最大长度 11，占位提示“请输入手机号” -->
      <input v-model="phone" type="tel" maxlength="11" placeholder="请输入手机号" />
      <!-- 验证码行容器，class 为 code-row，横向排列验证码输入框与获取按钮 -->
      <div class="code-row">
        <!-- 验证码输入框，v-model 双向绑定 code，最大长度 6，占位提示“请输入验证码” -->
        <input v-model="code" maxlength="6" placeholder="请输入验证码" />
        <!-- 获取验证码按钮，倒计时大于 0 时禁用，点击触发 sendCode，文本随倒计时变化 -->
        <button class="code-btn" :disabled="countdown > 0" @click="sendCode">
          <!-- 按钮文本：倒计时中显示剩余秒数，否则显示“获取验证码” -->
          {{ countdown > 0 ? countdown + 's' : '获取验证码' }}
        <!-- 闭合获取验证码按钮 button 元素 -->

        </button>
      <!-- 闭合验证码行容器 div 元素 -->

      </div>
      <!-- 演示提示段落，class 为 demo-tip，说明演示环境任意手机号与验证码均可 -->
      <p class="demo-tip">演示环境：任意 11 位手机号 + 任意 4~6 位验证码即可登录</p>
      <!-- 登录按钮，class 为 primary-btn，加载中禁用，点击触发 loginByPhone -->
      <button class="primary-btn" :disabled="loading" @click="loginByPhone">登录</button>
    <!-- 闭合表单卡片 div 元素 -->

    </div>
    <!-- 分隔区域，class 为 divider，展示“其他登录方式”文字 -->
    <div class="divider"><span>其他登录方式</span></div>
    <!-- 微信授权登录按钮，class 为 wechat-btn，加载中禁用，点击触发 loginByWechat -->
    <button class="wechat-btn" :disabled="loading" @click="loginByWechat">
      <!-- 微信图标元素，class 为 wechat-icon，显示微信标识 -->
      <span class="wechat-icon">●</span>
      <!-- 按钮文本，显示“微信授权登录” -->
      <span>微信授权登录</span>
    <!-- 闭合微信授权登录按钮 button 元素 -->

    </button>
    <!-- 错误提示段落，v-if 依据 err 是否有值决定是否渲染，class 为 err，显示错误信息 -->
    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合登录页根容器 div 元素 -->

  </div>
</template>

<style scoped>
/* 登录页根容器：铺满手机壳高度，纵向排布并居中 */
.login-page { min-height: 100vh; padding: 48px 24px 32px; display: flex; flex-direction: column; gap: 20px; background: linear-gradient(180deg, #eff6ff 0%, #f5f7fa 45%); }
/* 品牌区域：居中显示 Logo 与标题 */
.brand-box { text-align: center; }
/* Logo 图标：蓝色圆角方块 */
.logo { width: 64px; height: 64px; margin: 0 auto 12px; border-radius: 16px; background: #2563eb; color: #fff; font-size: 36px; line-height: 64px; }
/* 品牌名称样式 */
.brand-name { font-size: 22px; font-weight: 700; color: #1f2937; }
/* 品牌副标题样式 */
.brand-sub { margin-top: 6px; font-size: 13px; color: #6b7280; }
/* 表单卡片：白色圆角带阴影 */
.login-card { display: flex; flex-direction: column; gap: 12px; background: #fff; border-radius: 14px; padding: 20px; box-shadow: 0 4px 16px rgba(0, 0, 0, .06); }
/* 卡片标题样式 */
.card-title { font-size: 16px; font-weight: 600; color: #374151; }
/* 输入框通用样式：撑满宽度 */
.login-card input { width: 100%; }
/* 验证码行：输入框与按钮横向排列 */
.code-row { display: flex; gap: 8px; }
/* 验证码输入框占据剩余宽度 */
.code-row input { flex: 1; min-width: 0; }
/* 获取验证码按钮样式 */
.code-btn { white-space: nowrap; background: #eff6ff; color: #2563eb; }
/* 演示提示文字 */
.demo-tip { margin: 0; font-size: 12px; color: #9ca3af; }
/* 主登录按钮：撑满宽度 */
.primary-btn { width: 100%; padding: 12px; border-radius: 10px; font-size: 16px; }
/* 分隔区域：两侧横线 + 中间文字 */
.divider { display: flex; align-items: center; gap: 12px; color: #9ca3af; font-size: 12px; }
/* 分隔线左右横线 */
.divider::before, .divider::after { content: ''; flex: 1; height: 1px; background: #e5e7eb; }
/* 微信授权登录按钮：绿色，撑满宽度并居中内容 */
.wechat-btn { width: 100%; padding: 12px; border-radius: 10px; background: #22c55e; color: #fff; font-size: 16px; display: flex; align-items: center; justify-content: center; gap: 8px; }
/* 微信图标样式 */
.wechat-icon { font-size: 18px; }
</style>
