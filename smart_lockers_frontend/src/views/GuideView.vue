<script setup lang="ts">
// @ticket 03-首页：寄存指南
// @proto 原型对齐：05-寄存指南长图、04-一期用户界面-6、06-用户界面规划-12 寄存指南图文步骤
import { ref, onMounted } from 'vue'
import { getGuide } from '@/api/home'
// 引入原型图 05 提取的 6 步图文与温馨提示，用于接口内容偏少时补齐
import { protoSteps, protoTips, type GuideStep } from '@/data/guideSteps'
// 执行 interface Step { title: string; desc: string }，实现对应的前端逻辑

interface Step { title: string; desc: string; visual?: string[] }
// 声明变量 defaultSteps:，用于保存数据或函数引用

const defaultSteps: Step[] = [
  // 执行 { title: '01 查询附近网点', desc: '打开首页，定位城市后查询周边寄存点，选择合适的柜机。' },，实现对应的前端逻辑

  { title: '01 查询附近网点', desc: '打开首页，定位城市后查询周边寄存点，选择合适的柜机。' },
  // 执行 { title: '02 扫码开始寄存', desc: '到达柜机后扫描柜身二维码，选择空闲柜门。' },，实现对应的前端逻辑

  { title: '02 扫码开始寄存', desc: '到达柜机后扫描柜身二维码，选择空闲柜门。' },
  // 执行 { title: '03 选择柜型与时长', desc: '按物品大小选择小/中/大格，设定寄存时长并支付。' },，实现对应的前端逻辑

  { title: '03 选择柜型与时长', desc: '按物品大小选择小/中/大格，设定寄存时长并支付。' },
  // 执行 { title: '04 放入物品关门', desc: '柜门开启后放入物品并关好柜门，保存取件凭证。' },，实现对应的前端逻辑

  { title: '04 放入物品关门', desc: '柜门开启后放入物品并关好柜门，保存取件凭证。' },
  // 执行 { title: '05 扫码取出', desc: '返回柜机输入取件码或在订单页点击开柜，取出物品。' },，实现对应的前端逻辑

  { title: '05 扫码取出', desc: '返回柜机输入取件码或在订单页点击开柜，取出物品。' },
// 执行 ]，实现对应的前端逻辑

]
// 声明变量 steps，用于保存数据或函数引用

const steps = ref<Step[]>(defaultSteps)
// 声明变量 extraContent，用于保存数据或函数引用

const extraContent = ref('')
// 声明变量 tip，用于保存数据或函数引用

const tip = ref('')
// 定义函数/方法，封装可复用的处理逻辑：function parseContent(content: string): Step[] {

function parseContent(content: string): Step[] {
  // 声明变量 lines，用于保存数据或函数引用

  const lines = content.split('\n').map((x) => x.trim()).filter(Boolean)
  // 声明变量 parsed:，用于保存数据或函数引用

  const parsed: Step[] = []
  // 执行 for (let i = 0; i < lines.length; i++) {，实现对应的前端逻辑

  for (let i = 0; i < lines.length; i++) {
    // 执行 parsed.push({ title: '步骤 ' + (i + 1), desc: lines[i] })，实现对应的前端逻辑

    parsed.push({ title: '步骤 ' + (i + 1), desc: lines[i] })
  // 当前代码块结束

  }
  // 返回 parsed 给调用方

  return parsed
// 当前代码块结束

}
// 补缺函数：以现有步骤为基础，不足原型 6 步时用原型图 05 的图文步骤补齐，不覆盖已有步骤

function fillSteps(base: Step[]): Step[] {
  // 复制现有步骤，避免修改原数组

  const merged: Step[] = base.slice()
  // 从现有长度补齐到原型步骤总数

  for (let i = merged.length; i < protoSteps.length; i++) {
    // 当前补齐项的原型来源

    const proto: GuideStep = protoSteps[i]
    // 追加原型步骤（标题/说明/示意图要点），保留 base 既有内容不动

    merged.push({ title: proto.title, desc: proto.desc, visual: proto.visual })
  // 当前代码块结束

  }
  // 返回补齐后的步骤数组

  return merged
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：async function load() {

async function load() {
  // 执行 try {，实现对应的前端逻辑

  try {
    // 声明变量 res，用于保存数据或函数引用

    const res = await getGuide()
    // 条件判断：if (res && res.content) { 时执行下面逻辑

    if (res && res.content) {
      // 声明变量 parsed，用于保存数据或函数引用

      const parsed = parseContent(res.content)
      // 条件判断：if (parsed.length) { 时执行下面逻辑

      if (parsed.length) {
        // 赋值更新：steps.value = parsed

        steps.value = parsed
        // 赋值更新：extraContent.value = ''

        extraContent.value = ''
      // 当前代码块结束

      }
    // 当前代码块结束

    }
  // 执行 } catch (e) {，实现对应的前端逻辑

  } catch (e) {
    // 赋值更新：tip.value = '当前展示演示数据'

    tip.value = '当前展示演示数据'
  // 当前代码块结束

  }
// 当前代码块结束

}
// 执行 onMounted(load)，实现对应的前端逻辑

onMounted(load)
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 banner -->

    <div class="banner">寄存无忧 · 安心随行</div>
    <!-- p 元素，class 为 tip -->

    <p v-if="tip" class="tip">{{ tip }}</p>
    <!-- div 元素，class 为 step -->

    <!-- 基于 steps 补齐后的 6 步循环生成步骤图文，key 取索引 -->

    <div v-for="(s, i) in fillSteps(steps)" :key="i" class="step">
      <!-- div 元素，class 为 step-title -->

      <div class="step-title">{{ s.title }}</div>
      <!-- div 元素，class 为 step-desc -->

      <div class="step-desc">{{ s.desc }}</div>
      <!-- 原型图 05 示意图文字要点：有 visual 时逐条展示，补齐图文 -->

      <div v-if="s.visual && s.visual.length" class="visual">
        <!-- 逐条展示示意图要点 -->

        <span v-for="(v, vi) in s.visual" :key="vi" class="visual-line">{{ v }}</span>
      <!-- 闭合 div 元素 -->

      </div>
      <!-- 无图文数据时的原占位框，保留既有文案 -->

      <div v-else class="ph">步骤示意图 {{ i + 1 }}</div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- 温馨提示（原有文案保留） -->

    <div class="warn">温馨提示：请勿存放贵重物品与违禁品；超时未取将按规则续费。</div>
    <!-- 原型图 05 补充的温馨提示（只新增，不删改上方既有提示） -->

    <div v-for="(t, ti) in protoTips" :key="ti" class="warn">{{ ti + 1 }}. {{ t }}</div>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 12px; }
/* .banner 样式定义 */

.banner { background: linear-gradient(135deg, #e8f4ff, #d6ebff); color: #2563eb; font-size: 18px; font-weight: 700; text-align: center; padding: 18px; border-radius: 12px; }
/* .tip 样式定义 */

.tip { font-size: 12px; color: #9ca3af; }
/* .step 样式定义 */

.step { background: #fff; border-radius: 10px; padding: 14px; border: 1px solid #e5e7eb; }
/* .step-title 样式定义 */

.step-title { font-size: 15px; font-weight: 700; color: #1d4ed8; }
/* .step-desc 样式定义 */

.step-desc { font-size: 13px; color: #4b5563; margin: 8px 0; line-height: 1.6; }
/* .ph 样式定义 */

.ph { height: 90px; border-radius: 8px; background: #eff6ff; color: #93c5fd; display: flex; align-items: center; justify-content: center; font-size: 13px; }
/* .visual 样式定义：原型图步骤示意图要点列表 */

.visual { background: #eff6ff; border-radius: 8px; padding: 10px 12px; display: flex; flex-direction: column; gap: 4px; }
/* .visual-line 样式定义 */

.visual-line { font-size: 12px; color: #2563eb; line-height: 1.5; }
/* .visual-line 前置圆点 */

.visual-line::before { content: '· '; color: #60a5fa; }
/* .warn 样式定义 */

.warn { background: #fef3c7; color: #92400e; font-size: 12px; padding: 12px; border-radius: 8px; line-height: 1.6; }
</style>

