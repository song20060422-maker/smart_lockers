// @ticket B端-前端：轻提示
import { ref } from 'vue'
const message = ref('')
let timer: ReturnType<typeof setTimeout> | null = null
export function useNotify() {
  function notify(msg: string) {
    message.value = msg
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => (message.value = ''), 2500)
  }
  return { message, notify }
}
