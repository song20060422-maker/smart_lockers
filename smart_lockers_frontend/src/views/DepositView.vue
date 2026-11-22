<script setup lang="ts">
// @ticket 04-物品寄存：选柜/时长/优惠券/费用试算/落单/换柜
// @proto 原型对齐：v-fee 网点收费规则 / v-opendoor 异常换柜；见 AGENT.MD 2.1
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { queryCabinets, quoteFee, createDeposit, changeCabinet } from '@/api/deposit'
import { myCoupons } from '@/api/user'
import type { Cabinet, FeeResult, Coupon } from '@/types'
// 新增统一券数据工具：归一化后端 {available,expired} 对象并对空值做防御，杜绝 -NaN元
import { normalizeCoupons, couponOptionText, type NormalizedCoupon } from '@/utils/coupon'

// 声明变量 route = useRoute()，用于保存数据或函数引用

const route = useRoute()
// 声明变量 router = useRouter()，用于保存数据或函数引用

const router = useRouter()
// 声明变量 siteId = Number(route.query.site_id || 0)，用于保存数据或函数引用

const siteId = Number(route.query.site_id || 0)
// 声明变量 cabinets = ref<Cabinet[]>([])，用于保存数据或函数引用

const cabinets = ref<Cabinet[]>([])
// 声明变量 coupons = ref<NormalizedCoupon[]>([])，统一券口径；原 Coupon 类型保留以兼容其它逻辑

const coupons = ref<NormalizedCoupon[]>([])
// 声明变量 doorId = ref<number>(0)，用于保存数据或函数引用

const doorId = ref<number>(0)
// 声明变量 hours = ref<number>(1)，用于保存数据或函数引用

const hours = ref<number>(1)
// 声明变量 couponId = ref<number | undefined>(undefined)，用于保存数据或函数引用

const couponId = ref<number | undefined>(undefined)
// 声明变量 fee = ref<FeeResult | null>(null)，用于保存数据或函数引用

const fee = ref<FeeResult | null>(null)
// 声明变量 orderNo = ref('')，用于保存数据或函数引用

const orderNo = ref('')
// 声明变量 err = ref('')，用于保存数据或函数引用

const err = ref('')
// 声明变量 msg = ref('')，用于保存数据或函数引用

const msg = ref('')

// 调用 onMounted(async () => {，触发对应功能

onMounted(async () => {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 等待异步操作完成：cabinets.value = await queryCabinets({ site_id: siteId })

    cabinets.value = await queryCabinets({ site_id: siteId })
    // 等待异步操作完成：coupons.value = normalizeCoupons(await myCoupons())，统一券口径并空值防御

    coupons.value = normalizeCoupons(await myCoupons())
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 执行 })，实现对应的前端逻辑

})

// 定义函数/方法，封装可复用的处理逻辑：async function doQuote() {

async function doQuote() {
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 等待异步操作完成：fee.value = await quoteFee({ size: 'M', hours: hours.value, coupon: couponId.value })

    fee.value = await quoteFee({ size: 'M', hours: hours.value, coupon: couponId.value })
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}

// @ticket 04-物品寄存：敏感操作二次确认 + 幂等键防重复提交
// 定义函数/方法，封装可复用的处理逻辑：async function submitOrder() {

async function submitOrder() {
  // 条件判断：if (!doorId.value) { err.value = '请先选择柜门'; return } 时执行下面逻辑

  if (!doorId.value) { err.value = '请先选择柜门'; return }
  // 条件判断：if (!window.confirm('确认按当前费用下单并支付？')) return 时执行下面逻辑

  if (!window.confirm('确认按当前费用下单并支付？')) return
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 idem = `deposit-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`，用于保存数据或函数引用

    const idem = `deposit-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    // 声明变量 r = await createDeposit({，用于保存数据或函数引用

    const r = await createDeposit({
      // 定义对象属性 door_id，取值为 doorId.value, site_id: siteId, hours: hours.value, coupon: couponId.value

      door_id: doorId.value, site_id: siteId, hours: hours.value, coupon: couponId.value,
    // 执行 }, idem)，实现对应的前端逻辑

    }, idem)
    // 赋值更新：orderNo.value = r.order_no

    orderNo.value = r.order_no
    // 赋值更新：msg.value = '下单成功'

    msg.value = '下单成功'
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：async function doChangeCabinet() {

async function doChangeCabinet() {
  // 条件判断：if (!orderNo.value) { err.value = '暂无可换柜的订单'; return } 时执行下面逻辑

  if (!orderNo.value) { err.value = '暂无可换柜的订单'; return }
  // 声明变量 newDoor = window.prompt('输入新柜门 ID')，用于保存数据或函数引用

  const newDoor = window.prompt('输入新柜门 ID')
  // 条件判断：if (!newDoor) return 时执行下面逻辑

  if (!newDoor) return
  // 尝试执行可能出错的逻辑，配合 catch 进行错误处理

  try {
    // 声明变量 r = await changeCabinet(orderNo.value, Number(newDoor))，用于保存数据或函数引用

    const r = await changeCabinet(orderNo.value, Number(newDoor))
    // 赋值更新：doorId.value = r.door_id

    doorId.value = r.door_id
    // 赋值更新：msg.value = `已换柜至 ${r.door_id}`

    msg.value = `已换柜至 ${r.door_id}`
  // 捕获上述逻辑抛出的异常并进行错误处理

  } catch (e) { err.value = (e as Error).message }
// 当前代码块结束

}

// 定义函数/方法，封装可复用的处理逻辑：function fen(v: number): string { return (v / 100).toFixed(2) }

function fen(v: number): string { return (v / 100).toFixed(2) }
</script>

<template>
  <!-- 页面根容器，class 为 col，纵向排列寄存相关卡片 -->

  <div class="col">
    <!-- 卡片容器，class 为 card，展示可选柜子列表 -->

    <div class="card">
      <!-- 二级标题，显示“选择柜子” -->

      <h2>选择柜子</h2>
      <!-- 基于 cabinets 数组循环生成柜子条目，key 取 c.id，样式类为 list-item -->

      <div v-for="c in cabinets" :key="c.id" class="list-item">
        <!-- 行容器，class 为 row，横向排列单个柜子的编码、规格、可用量与选择按钮 -->

        <div class="row">
          <!-- 加粗元素，插值 c.code 显示柜子编码 -->

          <strong>{{ c.code }}</strong>
          <!-- 行内标签元素，class 为 tag，插值 c.size 显示柜子规格 -->

          <span class="tag">{{ c.size }}</span>
          <!-- 行内文本元素，class 为 muted，显示可用量并插值 c.available 与 c.total -->

          <span class="muted">可用 {{ c.available }}/{{ c.total }}</span>
          <!-- 选择按钮，动态绑定 class：doorId 不等于 c.id 时添加 btn-ghost 样式，点击时先将 doorId 设为 c.id 再调用 doQuote() 试算费用 -->

          <button :class="{ 'btn-ghost': doorId !== c.id }" @click="doorId = c.id; doQuote()">
            <!-- 按钮文本插值，当 doorId 等于 c.id 时显示“已选”，否则显示“选择” -->

            {{ doorId === c.id ? '已选' : '选择' }}
          <!-- 闭合选择按钮 button 元素 -->

          </button>
        <!-- 闭合当前柜子行容器 div 元素 -->

        </div>
      <!-- 闭合当前柜子条目 div 元素 -->

      </div>
      <!-- 空状态提示段落，当 cabinets 长度为 0 时渲染，class 为 muted，显示“该寄存点暂无可用柜子” -->

      <p v-if="!cabinets.length" class="muted">该寄存点暂无可用柜子</p>
    <!-- 闭合柜子列表卡片 div 元素 -->

    </div>

    <!-- 卡片容器，class 为 card，展示寄存费用设置与下单操作 -->

    <div class="card">
      <!-- 二级标题，显示“寄存设置” -->

      <h2>寄存设置</h2>
      <!-- 行容器，class 为 row，横向排列时长、优惠券与试算按钮 -->

      <div class="row">
        <!-- 表单标签，显示“时长(小时)”与数字输入框，v-model.number 双向绑定 hours，最小值 1，值变化时触发 doQuote -->

        <label>时长(小时) <input v-model.number="hours" type="number" min="1" @change="doQuote" /></label>
        <!-- 表单标签，显示“优惠券”文本（本行与后续行共同构成标签内容） -->

        <label>优惠券
          <!-- 下拉框，v-model.number 双向绑定 couponId，值变化时触发 doQuote -->

          <select v-model.number="couponId" @change="doQuote">
            <!-- 下拉选项，值为 undefined，显示“不使用” -->

            <option :value="undefined">不使用</option>
            <!-- 基于 coupons 数组循环生成优惠券选项，使用统一面值文案，杜绝 -NaN元 -->

            <option v-for="c in coupons" :key="c.id" :value="c.id">{{ couponOptionText(c) }}</option>
          <!-- 闭合优惠券下拉框 select 元素 -->

          </select>
        <!-- 闭合优惠券表单标签 label 元素 -->

        </label>
        <!-- 试算按钮，点击触发 doQuote 计算费用 -->

        <button @click="doQuote">试算费用</button>
      <!-- 闭合寄存设置行容器 div 元素 -->

      </div>
      <!-- 费用明细容器，class 为 list-item，v-if 依据 fee 是否有值决定是否渲染 -->

      <div v-if="fee" class="list-item">
        <!-- 段落元素，显示基础费与优惠金额，均经 fen 格式化 -->

        <p>基础费：{{ fen(fee.base_fee) }}元 / 优惠：{{ fen(fee.discount) }}元</p>
        <!-- 段落元素，显示押金与应付金额，应付金额以加粗 strong 元素突出显示 -->

        <p>押金：{{ fen(fee.deposit) }}元 / 应付：<strong>{{ fen(fee.payable) }}元</strong></p>
      <!-- 闭合费用明细容器 div 元素 -->

      </div>
      <!-- 行容器，class 为 row，横向排列下单与异常换柜按钮 -->

      <div class="row">
        <!-- 下单并支付按钮，当 fee 无值时禁用，点击触发 submitOrder -->

        <button :disabled="!fee" @click="submitOrder">下单并支付</button>
        <!-- 异常换柜按钮，class 为 btn-ghost，点击触发 doChangeCabinet -->

        <button class="btn-ghost" @click="doChangeCabinet">异常换柜</button>
      <!-- 闭合下单操作行容器 div 元素 -->

      </div>
      <!-- 提示段落，当 msg 有值时渲染，class 为 muted，插值 msg 显示提示内容 -->

      <p v-if="msg" class="muted">{{ msg }}</p>
    <!-- 闭合寄存设置卡片 div 元素 -->

    </div>
    <!-- ===== 补齐入口：预支付/支付完成原型界面跳转（只新增，不改动上方原有内容） ===== -->
    <!-- 支付流程卡片容器，class 为 card，放置费用确认与支付完成界面的跳转入口 -->
    <div class="card">
      <!-- 二级标题，显示“支付流程” -->
      <h2>支付流程</h2>
      <!-- 行容器，class 为 row，横向排列各支付界面入口链接 -->
      <div class="row">
        <!-- 链接：费用确认，跳转 /pay-confirm -->
        <router-link to="/pay-confirm" class="tag" style="text-decoration:none">费用确认</router-link>
        <!-- 链接：支付完成，跳转 /pay-success -->
        <router-link to="/pay-success" class="tag" style="text-decoration:none">支付完成</router-link>
        <!-- 链接：人脸识别，跳转 /face -->
        <router-link to="/face" class="tag" style="text-decoration:none">人脸识别</router-link>
      <!-- 闭合入口链接行容器 div 元素 -->
      </div>
    <!-- 闭合支付流程卡片 div 元素 -->
    </div>

    <!-- 错误提示段落，当 err 有值时渲染，class 为 err，插值 err 显示错误信息 -->

    <p v-if="err" class="err">{{ err }}</p>
  <!-- 闭合页面根容器 div 元素 -->

  </div>
</template>
