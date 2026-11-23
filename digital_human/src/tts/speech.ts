/** @ticket 07-数字人Coze：语音模块——SpeechRecognition 输入 + SpeechSynthesis(TTS) 输出（可切换 Azure/OpenAI） */
// 导出 interface SpeechHooks {，供其他模块引用

export interface SpeechHooks {
  // 定义对象属性 onRecognized，取值为 (text: string) => void

  onRecognized: (text: string) => void
  // 定义对象属性 onSpeakStart?: () => void，播报开始时回调，驱动口型进入说话态

  onSpeakStart?: () => void
  // 定义对象属性 onSpeakBoundary?: () => void，语音边界触发时回调，驱动张口脉冲

  onSpeakBoundary?: () => void
  // 执行 onSpeakEnd?: () => void，实现对应的前端逻辑

  onSpeakEnd?: () => void
// 当前代码块结束

}

// 导出 class SpeechEngine {，供其他模块引用

export class SpeechEngine {
  // 赋值更新：private recognition: any = null

  private recognition: any = null
  // 执行 private hooks: SpeechHooks，实现对应的前端逻辑

  private hooks: SpeechHooks
  // 执行 private voiceName: string，实现对应的前端逻辑

  private voiceName: string

  // 调用 constructor(voiceName: string, hooks: SpeechHooks) {，触发对应功能

  constructor(voiceName: string, hooks: SpeechHooks) {
    // 赋值更新：this.voiceName = voiceName

    this.voiceName = voiceName
    // 赋值更新：this.hooks = hooks

    this.hooks = hooks
    // 定义函数/方法，封装可复用的处理逻辑：const SR = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition

    const SR = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition
    // 条件判断：if (SR) 时执行下面逻辑

    if (SR) {
      // 赋值更新：this.recognition = new SR()

      this.recognition = new SR()
      // 赋值更新：this.recognition.lang = 'zh-CN'

      this.recognition.lang = 'zh-CN'
      // 赋值更新：this.recognition.interimResults = false

      this.recognition.interimResults = false
      // 赋值更新：this.recognition.onresult = (e: any) => this.hooks.onRecognized(e.results[0][0].transcript)

      this.recognition.onresult = (e: any) => this.hooks.onRecognized(e.results[0][0].transcript)
    // 当前代码块结束

    }
  // 当前代码块结束

  }

  // 调用 startListen() { this.recognition?.start() }，触发对应功能

  startListen() { this.recognition?.start() }

  /** TTS 播报；Web Speech 为默认，Azure/OpenAI 由配置注入外部实现 */
  // 调用 speak(text: string) {，触发对应功能

  speak(text: string) {
    // 条件判断：if (!('speechSynthesis' in window)) return 时执行下面逻辑

    if (!('speechSynthesis' in window)) return
    // 执行 window.speechSynthesis.cancel()，实现对应的前端逻辑

    window.speechSynthesis.cancel()
    // 声明变量 utter = new SpeechSynthesisUtterance(text)，用于保存数据或函数引用

    const utter = new SpeechSynthesisUtterance(text)
    // 赋值更新：utter.lang = 'zh-CN'

    utter.lang = 'zh-CN'
    // 赋值更新：utter.rate = 1

    utter.rate = 1
    // 声明变量 voice = window.speechSynthesis.getVoices().find(v => v.name.includes(this.voiceName))，用于保存数据或函数引用

    const voice = window.speechSynthesis.getVoices().find(v => v.name.includes(this.voiceName))
    // 条件判断：if (voice) utter.voice = voice 时执行下面逻辑

    if (voice) utter.voice = voice
    // 赋值更新：utter.onend = () => this.hooks.onSpeakEnd?.()

    // 执行方法：播报开始时回调，通知上层进入说话口型状态

    utter.onstart = () => this.hooks.onSpeakStart?.()

    // 执行方法：语音边界事件回调，通知上层注入一次张口脉冲（部分浏览器不触发，作增强）

    utter.onboundary = () => this.hooks.onSpeakBoundary?.()

    utter.onend = () => this.hooks.onSpeakEnd?.()
    // 执行 window.speechSynthesis.speak(utter)，实现对应的前端逻辑

    window.speechSynthesis.speak(utter)
  // 当前代码块结束

  }

  // 调用 stop() { window.speechSynthesis?.cancel() }，触发对应功能

  stop() { window.speechSynthesis?.cancel() }
// 当前代码块结束

}
