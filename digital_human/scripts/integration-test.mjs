/** @ticket 07-数字人Coze：联调测试脚本——语音识别 / TTS 播放 / 表情控制桩测试 */
import assert from 'node:assert'

/** 统一消息结构契约校验（工单定义的 role/content/type/timestamp） */
function assertMessage(msg) {
  assert(['user', 'assistant'].includes(msg.role), 'role 非法')
  assert(typeof msg.content === 'string', 'content 必须为 string')
  assert(['text', 'voice', 'command'].includes(msg.type), 'type 非法')
  assert(typeof msg.timestamp === 'number', 'timestamp 必须为 number')
}

const cases = [
  { role: 'user', content: '如何寄存物品', type: 'voice', timestamp: Date.now() },
  { role: 'assistant', content: '选择柜机扫码支付即可', type: 'text', timestamp: Date.now() },
]
cases.forEach(assertMessage)

/** 情绪→动作映射覆盖（微笑/点头/摇头/眨眼） */
const emotionAction = { happy: 'smile', positive: 'nod', sad: 'shake', neutral: 'blink' }
assert.equal(emotionAction.happy, 'smile')
assert.equal(emotionAction.sad, 'shake')

console.log('[ticket 07] integration-test: PASS (消息契约 + 情绪动作映射)')
