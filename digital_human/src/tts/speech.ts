/** @ticket 07-数字人Coze：语音模块——SpeechRecognition 输入 + SpeechSynthesis(TTS) 输出（可切换 Azure/OpenAI） */
export interface SpeechHooks {
  onRecognized: (text: string) => void
  onSpeakEnd?: () => void
}

export class SpeechEngine {
  private recognition: any = null
  private hooks: SpeechHooks
  private voiceName: string

  constructor(voiceName: string, hooks: SpeechHooks) {
    this.voiceName = voiceName
    this.hooks = hooks
    const SR = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition
    if (SR) {
      this.recognition = new SR()
      this.recognition.lang = 'zh-CN'
      this.recognition.interimResults = false
      this.recognition.onresult = (e: any) => this.hooks.onRecognized(e.results[0][0].transcript)
    }
  }

  startListen() { this.recognition?.start() }

  /** TTS 播报；Web Speech 为默认，Azure/OpenAI 由配置注入外部实现 */
  speak(text: string) {
    if (!('speechSynthesis' in window)) return
    window.speechSynthesis.cancel()
    const utter = new SpeechSynthesisUtterance(text)
    utter.lang = 'zh-CN'
    utter.rate = 1
    const voice = window.speechSynthesis.getVoices().find(v => v.name.includes(this.voiceName))
    if (voice) utter.voice = voice
    utter.onend = () => this.hooks.onSpeakEnd?.()
    window.speechSynthesis.speak(utter)
  }

  stop() { window.speechSynthesis?.cancel() }
}
