<script setup lang="ts">
// @ticket 04-物品寄存 + 05-订单查看与取出：开门异常页——重新开门/已解决/取消订单退费/更换格口/倒计时重试/联系管理员
// @proto 原型对齐：02-开门异常提示优化-1~7（异常图标+格口号大字+联系管理员电话+弹窗流程链）
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
// 声明变量 router，用于保存数据或函数引用
const router = useRouter()
// 声明变量 route，用于保存数据或函数引用
const route = useRoute()
// 场景：opened 正常已开 / abnormal 开门异常（对照原型主页面）
// 声明变量 scene = ref(...)，用于保存数据或函数引用
const scene = ref<string>((route.query.scene as string) || 'abnormal')
// 格口号（原型大字"15号"）
// 声明变量 doorNo = ref(...)，用于保存数据或函数引用
const doorNo = ref<string>((route.query.door_no as string) || '15')
// 网点管理员电话（原型"13456789099"）
// 声明变量 adminTel = ref(...)，用于保存数据或函数引用
const adminTel = ref<string>((route.query.admin_tel as string) || '13456789099')
// 弹窗类型：none / retryFail(重新开门失败→取消订单|更换柜门) / refunded(已取消订单并退回金额) / countdown(重新开门倒计时+联系管理员)
// 声明变量 dialog = ref(...)，用于保存数据或函数引用
const dialog = ref<'none' | 'retryFail' | 'refunded' | 'countdown'>('none')
// 更换格口底部弹层显隐（对照原型 02-5）
// 声明变量 showChange = ref(false)，用于保存数据或函数引用
const showChange = ref(false)
// 重新开门倒计时秒数（对照原型"重新开门（10）"）
// 声明变量 countdown = ref(0)，用于保存数据或函数引用
const countdown = ref(0)
// 声明变量 timer，定时器句柄
let timer: number | null = null

// 启动倒计时（对照原型 02-6：重新开门失败后 10 秒内禁点并显示剩余秒数）
// 定义函数/方法 startCountdown，封装可复用的处理逻辑
function startCountdown() {
  // 初始化 10 秒
  countdown.value = 10
  // 每秒递减，归零后自动清除
  timer = window.setInterval(() => {
    // 递减
    countdown.value--
    // 归零清理
    if (countdown.value <= 0 && timer) { window.clearInterval(timer); timer = null }
  // 当前代码块结束
  }, 1000)
// 当前代码块结束
}

// 挂载：异常场景默认弹出"重新开门失败"弹窗（对照原型 02-2 流程入口）
// 调用 onMounted(() => {，触发对应功能
onMounted(() => {
  // 异常场景展示重试失败弹窗
  if (scene.value === 'abnormal') dialog.value = 'retryFail'
// 当前代码块结束
})
// 卸载清理定时器
// 调用 onUnmounted(() => {，触发对应功能
onUnmounted(() => {
  // 清理倒计时定时器
  if (timer) window.clearInterval(timer)
// 当前代码块结束
})

// 重新开门：首次点击模拟失败→弹出失败弹窗；倒计时结束后点击→成功回到"柜门已开"
// 定义函数/方法 reopen，封装可复用的处理逻辑
function reopen() {
  // 倒计时中禁点
  if (countdown.value > 0) return
  // 演示逻辑：第一次失败弹倒计时窗，之后成功
  if (dialog.value === 'retryFail') {
    // 切换为倒计时弹窗（含联系管理员电话）
    dialog.value = 'countdown'
    // 启动 10 秒倒计时
    startCountdown()
    // 结束当前函数
    return
  // 当前代码块结束
  }
  // 重试成功：回到正常开门态
  scene.value = 'opened'
  // 关闭弹窗
  dialog.value = 'none'
// 当前代码块结束
}

// 已解决：返回订单列表
// 定义函数/方法 resolved，封装可复用的处理逻辑
function resolved() {
  // 执行 router.push('/orders')，实现对应的前端逻辑
  router.push('/orders')
// 当前代码块结束
}

// 取消订单并退费：弹出"已为您取消订单并退回已支付金额~"提示弹窗（对照原型 02-3/02-4）
// 定义函数/方法 cancelOrder，封装可复用的处理逻辑
function cancelOrder() {
  // 展示退款成功弹窗
  dialog.value = 'refunded'
// 当前代码块结束
}

// 更换柜门：打开底部"更换格口"弹层（对照原型 02-5）
// 定义函数/方法 changeDoor，封装可复用的处理逻辑
function changeDoor() {
  // 关闭失败弹窗，打开换柜弹层
  dialog.value = 'none'
  // 打开弹层
  showChange.value = true
// 当前代码块结束
}

// 确认换柜：关闭弹层并回到正常开门态（演示）
// 定义函数/方法 confirmChange，封装可复用的处理逻辑
function confirmChange() {
  // 关闭弹层
  showChange.value = false
  // 回到开门成功态
  scene.value = 'opened'
// 当前代码块结束
}

// 弹窗"好的"：关闭弹窗（退款弹窗确认后回订单列表）
// 定义函数/方法 closeDialog，封装可复用的处理逻辑
function closeDialog() {
  // 记录当前弹窗类型
  const cur = dialog.value
  // 关闭弹窗
  dialog.value = 'none'
  // 退款弹窗确认后返回订单列表
  if (cur === 'refunded') router.push('/orders')
// 当前代码块结束
}
</script>

<template>
  <!-- 页面根容器：柜门已开（异常态） -->
  <div class="door-page">
    <!-- 异常图标区（CSS 柜格+红色叉标，对照原型） -->
    <div class="icon-wrap">
      <!-- 柜格图形 -->
      <div class="cab"><i></i></div>
      <!-- 叉标 -->
      <span class="cross">✕</span>
    <!-- 闭合图标区 -->
    </div>
    <!-- 异常提示红字 -->
    <div class="alert-text">开门异常，请重新开门</div>
    <!-- 格口号大字卡 -->
    <div class="door-card">{{ doorNo }}号</div>
    <!-- 底部联系管理员与操作区 -->
    <div class="bottom-area">
      <!-- 联系网点管理员 -->
      <p class="contact">如遇问题请联系网点管理员：<a class="tel-link" :href="`tel:${adminTel}`">{{ adminTel }}</a></p>
      <!-- 操作按钮行：已解决 + 重新开门 -->
      <div class="btn-row">
        <!-- 已解决（浅蓝按钮） -->
        <button class="soft-btn" @click="resolved">已解决</button>
        <!-- 重新开门（主按钮） -->
        <button class="main-btn" @click="reopen">重新开门</button>
      <!-- 闭合按钮行 -->
      </div>
      <!-- 取消订单链接 -->
      <p class="cancel-link" @click="cancelOrder">不想用了，取消订单</p>
    <!-- 闭合底部区 -->
    </div>

    <!-- 全屏弹窗：重新开门失败（取消订单 / 更换柜门）对照原型 02-2 -->
    <div v-if="dialog === 'retryFail'" class="mask">
      <!-- 弹窗主体 -->
      <div class="dialog">
        <!-- 标题 -->
        <div class="dg-title">提示</div>
        <!-- 内容 -->
        <div class="dg-text">抱歉，重新开门失败。<br />您可以取消订单或更换其他柜门~</div>
        <!-- 按钮行 -->
        <div class="dg-btns">
          <!-- 取消订单（浅色） -->
          <button class="soft-btn" @click="cancelOrder">取消订单</button>
          <!-- 更换柜门（主色） -->
          <button class="main-btn" @click="changeDoor">更换柜门</button>
        <!-- 闭合按钮行 -->
        </div>
      <!-- 闭合弹窗 -->
      </div>
    <!-- 闭合遮罩 -->
    </div>

    <!-- 全屏弹窗：重新开门失败+倒计时+联系管理员 对照原型 02-6 -->
    <div v-if="dialog === 'countdown'" class="mask">
      <!-- 弹窗主体 -->
      <div class="dialog">
        <!-- 标题 -->
        <div class="dg-title">提示</div>
        <!-- 内容：失败原因 + 管理员电话 -->
        <div class="dg-text">抱歉，重新开门失败。<br />请联系网点管理员：<a class="tel-link">{{ adminTel }}</a></div>
        <!-- 倒计时禁点按钮（原型"重新开门（10）"浅色禁用态） -->
        <button class="cd-btn" :disabled="countdown > 0" @click="reopen">重新开门（{{ countdown }}）</button>
      <!-- 闭合弹窗 -->
      </div>
    <!-- 闭合遮罩 -->
    </div>

    <!-- 全屏弹窗：已取消订单并退回金额 对照原型 02-3/02-4 -->
    <div v-if="dialog === 'refunded'" class="mask">
      <!-- 弹窗主体 -->
      <div class="dialog">
        <!-- 标题 -->
        <div class="dg-title">提示</div>
        <!-- 内容 -->
        <div class="dg-text">已为您取消订单并退回<br />已支付金额~</div>
        <!-- 好的按钮 -->
        <button class="ok-btn" @click="closeDialog">好的</button>
      <!-- 闭合弹窗 -->
      </div>
    <!-- 闭合遮罩 -->
    </div>

    <!-- 底部弹层：更换格口 对照原型 02-5 -->
    <div v-if="showChange" class="mask bottom" @click.self="showChange = false">
      <!-- 弹层主体 -->
      <div class="sheet">
        <!-- 标题行 -->
        <div class="sheet-head">
          <!-- 标题 -->
          <span class="sheet-title">更换格口</span>
          <!-- 关闭 -->
          <span class="sheet-close" @click="showChange = false">✕</span>
        <!-- 闭合标题行 -->
        </div>
        <!-- 柜门编号卡 -->
        <div class="sheet-door">
          <!-- 标签 -->
          <span class="sd-label">柜门编号</span>
          <!-- 编号大字 -->
          <span class="sd-no">{{ doorNo }}号</span>
        <!-- 闭合编号卡 -->
        </div>
        <!-- 即将开门提示 -->
        <p class="sheet-tip">即将为您打开新的格口</p>
        <!-- 底部操作行 -->
        <div class="sheet-foot">
          <!-- 取消订单链接 -->
          <span class="cancel-link" @click="cancelOrder">不想用了，取消订单</span>
          <!-- 确认按钮 -->
          <button class="confirm-btn" @click="confirmChange">确认</button>
        <!-- 闭合操作行 -->
        </div>
      <!-- 闭合弹层 -->
      </div>
    <!-- 闭合遮罩 -->
    </div>
  <!-- 闭合根容器 -->
  </div>
</template>

<style scoped>
/* 页面根容器：纵向布局，浅色底 */
.door-page { min-height: calc(100vh - 110px); display: flex; flex-direction: column; align-items: center; padding-top: 6vh; }
/* 异常图标区 */
.icon-wrap { position: relative; margin-bottom: 12px; }
/* 柜格图形（蓝色圆角立方体+开门斜面） */
.cab { width: 96px; height: 74px; background: linear-gradient(160deg, #3b82f6, #1d4ed8); border-radius: 12px; position: relative; }
/* 柜门斜面 */
.cab i { position: absolute; left: 10px; top: 10px; width: 34px; height: 54px; background: rgba(255,255,255,.35); border-radius: 6px; transform: skewY(-6deg); }
/* 红色叉标 */
.cross { position: absolute; right: -8px; bottom: -8px; width: 30px; height: 30px; border-radius: 50%; background: #ff4d4f; color: #fff; font-size: 15px; text-align: center; line-height: 30px; border: 2px solid #fff; }
/* 异常提示红字 */
.alert-text { color: #ff4d4f; font-size: 15px; font-weight: 600; margin-bottom: 22px; }
/* 格口号大字卡 */
.door-card { width: 84%; background: #fff; border-radius: 14px; padding: 26px 0; text-align: center; font-size: 44px; font-weight: 800; color: #111827; margin-bottom: auto; }
/* 底部区 */
.bottom-area { width: 100%; text-align: center; padding-bottom: 12px; }
/* 联系管理员灰字 */
.contact { font-size: 13px; color: #9ca3af; margin: 0 0 16px; }
/* 电话链接 */
.tel-link { color: #2563eb; text-decoration: underline; }
/* 按钮行 */
.btn-row { display: flex; gap: 14px; justify-content: center; margin-bottom: 18px; }
/* 浅蓝按钮（已解决/取消订单） */
.soft-btn { background: #dbe7ff; color: #2563eb; border: none; border-radius: 22px; padding: 12px 34px; font-size: 15px; cursor: pointer; }
/* 主按钮（重新开门/更换柜门） */
.main-btn { background: #2563eb; color: #fff; border: none; border-radius: 22px; padding: 12px 34px; font-size: 15px; cursor: pointer; }
/* 取消订单链接 */
.cancel-link { color: #2563eb; font-size: 14px; text-decoration: underline; cursor: pointer; }
/* 全屏遮罩 */
.mask { position: fixed; inset: 0; background: rgba(0,0,0,.5); z-index: 60; display: flex; align-items: center; justify-content: center; }
/* 底部弹层遮罩变体 */
.mask.bottom { align-items: flex-end; }
/* 居中弹窗 */
.dialog { width: 78%; background: #fff; border-radius: 14px; padding: 22px 18px 18px; text-align: center; }
/* 弹窗标题 */
.dg-title { font-size: 16px; font-weight: 700; color: #111827; margin-bottom: 12px; }
/* 弹窗内容 */
.dg-text { font-size: 14px; color: #374151; line-height: 1.8; margin-bottom: 16px; }
/* 弹窗按钮行 */
.dg-btns { display: flex; gap: 12px; justify-content: center; }
/* 倒计时禁点按钮 */
.cd-btn { width: 100%; background: #dbe7ff; color: #93b8f5; border: none; border-radius: 10px; padding: 12px 0; font-size: 15px; }
/* 好的按钮 */
.ok-btn { width: 60%; background: #2563eb; color: #fff; border: none; border-radius: 22px; padding: 12px 0; font-size: 15px; }
/* 底部弹层 */
.sheet { width: 100%; max-width: 430px; background: #fff; border-radius: 16px 16px 0 0; padding: 16px 16px 22px; }
/* 弹层标题行 */
.sheet-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
/* 弹层标题 */
.sheet-title { font-size: 16px; font-weight: 700; }
/* 关闭按钮 */
.sheet-close { color: #6b7280; cursor: pointer; }
/* 柜门编号卡 */
.sheet-door { background: #f5f7fa; border-radius: 10px; padding: 14px; display: flex; align-items: center; justify-content: space-between; }
/* 编号标签 */
.sd-label { font-size: 12px; color: #9ca3af; }
/* 编号大字 */
.sd-no { font-size: 30px; font-weight: 800; color: #111827; }
/* 即将开门提示 */
.sheet-tip { font-size: 13px; color: #ff7a45; margin: 12px 0 4px; }
/* 弹层底部操作行 */
.sheet-foot { display: flex; align-items: center; justify-content: space-between; margin-top: 10px; }
/* 确认按钮 */
.confirm-btn { background: #2563eb; color: #fff; border: none; border-radius: 20px; padding: 10px 34px; font-size: 15px; cursor: pointer; }
</style>
