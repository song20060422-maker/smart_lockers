// @ticket B端-前端：轻提示
import { ref } from 'vue'
// 声明变量 message = ref('')，用于保存数据或函数引用

const message = ref('')
// 声明变量 timer: ReturnType<typeof setTimeout> | null = null，用于保存数据或函数引用

let timer: ReturnType<typeof setTimeout> | null = null
// 导出 function useNotify() {，供其他模块引用

export function useNotify() {
  // 定义函数/方法，封装可复用的处理逻辑：function notify(msg: string) {

  function notify(msg: string) {
    // 赋值更新：message.value = msg

    message.value = msg
    // 条件判断：if (timer) clearTimeout(timer) 时执行下面逻辑

    if (timer) clearTimeout(timer)
    // 赋值更新：timer = setTimeout(() => (message.value = ''), 2500)

    timer = setTimeout(() => (message.value = ''), 2500)
  // 当前代码块结束

  }
  // 返回 { message, notify } 给调用方

  return { message, notify }
// 当前代码块结束

}
