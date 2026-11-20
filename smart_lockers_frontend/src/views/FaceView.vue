<script setup lang="ts">
// @ticket 04-物品寄存：人脸识别添加/开柜凭证
// @proto 原型对齐：06-用户界面规划-18 人脸识别（添加人脸识别/示例图）
import { ref } from 'vue'
import { useRouter } from 'vue-router'
// 声明变量 router，用于保存数据或函数引用

const router = useRouter()
// 声明变量 added，用于保存数据或函数引用

const added = ref(false)
// 定义函数/方法，封装可复用的处理逻辑：function addFace() {

function addFace() {
  // 赋值更新：added.value = true

  added.value = true
// 当前代码块结束

}
// 定义函数/方法，封装可复用的处理逻辑：function back() {

function back() {
  // 执行 router.back()，实现对应的前端逻辑

  router.back()
// 当前代码块结束

}
</script>

<template>
  <!-- div 元素，class 为 col -->

  <div class="col">
    <!-- div 元素，class 为 face-box -->

    <div class="face-box">
      <!-- div 元素，class 为 frame -->

      <div class="frame">人脸采集区域</div>
      <!-- p 元素，class 为 hint -->

      <p class="hint">请正对屏幕，保持光线充足，勿遮挡面部</p>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- div 元素，class 为 tips -->

    <div class="tips">
      <!-- div 元素 -->

      <div>· 人脸信息仅用于开柜身份核验</div>
      <!-- div 元素 -->

      <div>· 采集数据加密存储，可随时删除</div>
      <!-- div 元素 -->

      <div>· 添加后可在柜机前刷脸开柜</div>
    <!-- 闭合 div 元素 -->

    </div>
    <!-- button 元素，class 为 primary -->

    <button v-if="!added" class="primary" @click="addFace">添加人脸识别</button>
    <!-- div 元素，class 为 done -->

    <div v-else class="done">人脸识别已添加</div>
    <!-- button 元素，class 为 ghost -->

    <button class="ghost" @click="back">返回</button>
  <!-- 闭合 div 元素 -->

  </div>
</template>

<style scoped>
/* .col 样式定义 */

.col { display: flex; flex-direction: column; gap: 14px; }
/* .face-box 样式定义 */

.face-box { background: #fff; border-radius: 10px; padding: 16px; text-align: center; border: 1px solid #e5e7eb; }
/* .frame 样式定义 */

.frame { height: 200px; border: 2px dashed #93c5fd; border-radius: 50%; width: 200px; margin: 0 auto; display: flex; align-items: center; justify-content: center; color: #93c5fd; font-size: 13px; }
/* .hint 样式定义 */

.hint { font-size: 12px; color: #9ca3af; }
/* .tips 样式定义 */

.tips { font-size: 13px; color: #6b7280; line-height: 2; }
/* .primary 样式定义 */

.primary { height: 44px; border: 0; border-radius: 22px; background: #2563eb; color: #fff; font-size: 16px; }
/* .ghost 样式定义 */

.ghost { height: 40px; border: 1px solid #d1d5db; border-radius: 20px; background: #fff; color: #6b7280; font-size: 14px; }
/* .done 样式定义 */

.done { text-align: center; color: #16a34a; font-size: 15px; font-weight: 700; }
</style>

