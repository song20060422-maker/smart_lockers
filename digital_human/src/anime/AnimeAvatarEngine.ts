/** @ticket 07-数字人Coze：动漫全身数字人动作引擎 v2（画面完全静止：说话时仅人物嘴部开合 + 动作姿势切换）
 * 原理：整幅立绘不做任何位移/缩放/旋转（画面稳定）；嘴部通过同图切片层裁剪显示，
 * 语音边界脉冲驱动嘴层在"闭口（原画）↔ 张口（口腔阴影）"间切换，实现画面中人物说话的效果；
 * 智能体动作指令触发整幅姿势立绘淡入淡出切换（人物做动作，画面本身不晃动）。
 */
// 姿势资源清单（相对 public 根路径）@ticket 07
// 导出 const POSES，供其他模块引用

export const POSES = {
  // 待机姿势：站立微笑
  idle: '/anime/idle.png',
  // 挥手：打招呼/问好
  wave: '/anime/wave.png',
  // 鞠躬：欢迎/致意（nod 点头映射到此姿势，动作幅度更明显）
  bow: '/anime/bow.png',
  // 指引：伸手指示方向/引导寄存
  point: '/anime/point.png',
  // 思考：手托下巴（shake 摇头映射到此姿势，配合表情）
  think: '/anime/think.png',
  // 开心：双手高举欢呼（作业完成/取件成功等正向反馈）
  happy: '/anime/happy.png',
  // 讲解：双手掌心向上展示（present 普通讲解态）
  present: '/anime/present.png',
// 当前代码块结束

} as const

// 智能体动作 → 姿势映射表（action/emotion 双通道，未命中回退 present）@ticket 07
// 声明变量 ACTION_POSES，用于保存运行过程中的数据

const ACTION_POSES: Record<string, string> = {
  // 点头 → 鞠躬姿势（更明显的身体动作）
  nod: 'bow',
  // 微笑 → 挥手（亲和表达）
  smile: 'wave',
  // 摇头 → 思考姿势（配合疑惑表情）
  shake: 'think',
  // 挥手
  wave: 'wave',
  // 指引
  point: 'point',
  // 思考
  think: 'think',
  // 开心
  happy: 'happy',
  // 讲解/展示
  present: 'present',
  // 通用讲解兜底
  explain: 'present',
// 当前代码块结束

}

// 情绪 → 姿势映射（动作未指定时的次级通道）@ticket 07
// 声明变量 EMOTION_POSES，用于保存运行过程中的数据

const EMOTION_POSES: Record<string, string> = {
  // 积极 → 开心
  positive: 'happy',
  // 中性 → 讲解
  neutral: 'present',
  // 疑惑 → 思考
  confused: 'think',
  // 抱歉 → 鞠躬
  apologetic: 'bow',
// 当前代码块结束

}

// 定义并导出类 AnimeAvatarEngine：画面静止 + 动作姿势演出引擎（无口型层）
export class AnimeAvatarEngine {
  // 声明属性 host：舞台宿主
  private host: HTMLElement
  // 声明属性 base：当前展示层（显示激活姿势图，完全静止）
  private base: HTMLElement
  // 声明属性 fade：过渡层（姿势切换时淡入用）
  private fade: HTMLElement
  // 声明属性 current：当前姿势键名
  private current = 'idle'
  // 声明属性 speaking：播报状态
  private speaking = false
  // 声明属性 backTimer：动作姿势回到待机的定时器
  private backTimer: number | null = null
  // 声明属性 poseDir：姿势集目录（角色切换预留）
  private poseDir: string

  // 构造函数：接收宿主与姿势集目录，构建图层并预加载姿势图
  constructor(host: HTMLElement, poseDir = '/anime') {
    // 记录宿主
    this.host = host
    // 记录姿势目录
    this.poseDir = poseDir
    // 清空宿主避免叠影
    host.innerHTML = ''
    // 创建展示层
    this.base = document.createElement('div')
    // 展示层样式类
    this.base.className = 'an-base'
    // 加入宿主
    host.appendChild(this.base)
    // 创建过渡层（初始透明）
    this.fade = document.createElement('div')
    // 过渡层样式类
    this.fade.className = 'an-fade'
    // 加入宿主
    host.appendChild(this.fade)
    // 预加载全部姿势图（首次切换不闪白）
    for (const p of Object.values(POSES)) {
      // 创建图片对象触发加载
      const img = new Image()
      // 赋值图片地址（触发预加载）
      img.src = p
    // 当前代码块结束

    }
    // 初始显示待机姿势
    this.setPose('idle')
  // 当前代码块结束
  }

  // poseUrl 拼接姿势图地址（支持姿势集目录）
  private poseUrl(key: string): string {
    // 目录 + 文件名
    return `${this.poseDir}/${key}.png`
  // 当前代码块结束
  }

  // setPose 直接设置姿势（无过渡）
  private setPose(key: string): void {
    // 更新当前姿势
    this.current = key
    // 展示层背景指向目标姿势
    this.base.style.backgroundImage = `url(${this.poseUrl(key)})`
  // 当前代码块结束
  }

  // transitionTo 带淡入过渡的姿势切换：过渡层淡入目标姿势后置换展示层
  private transitionTo(key: string, holdMs = 1600): void {
    // 目标与当前相同则跳过
    if (key === this.current) return
    // 过渡层载入目标姿势并淡入
    this.fade.style.backgroundImage = `url(${this.poseUrl(key)})`
    this.fade.style.opacity = '1'
    // 300ms 过渡完成后：置换展示层、隐藏过渡层
    window.setTimeout(() => {
      // 展示层切为目标姿势
      this.setPose(key)
      // 过渡层淡出
      this.fade.style.opacity = '0'
    // 当前代码块结束
    }, 300)
    // 动作保持 holdMs 后回到待机（取消上一次未执行的回退）
    if (this.backTimer) window.clearTimeout(this.backTimer)
    // 安排回退
    this.backTimer = window.setTimeout(() => this.setPose('idle'), holdMs)
  // 当前代码块结束
  }

  // start 启动（画面完全静止，无渲染循环，无额外初始化）
  start(): void {
    // 画面静止，无逐帧动画
  // 当前代码块结束
  }

  // setAction 接收智能体动作/情绪指令，映射到姿势并过渡切换 @ticket 07
  setAction(a: { action?: string; emotion?: string }): void {
    // 动作通道优先，其次情绪通道，兜底讲解姿势
    const key = (a.action && ACTION_POSES[a.action]) || (a.emotion && EMOTION_POSES[a.emotion]) || 'present'
    // 过渡切换（保持 1.8s 回待机）
    this.transitionTo(key, 1800)
  // 当前代码块结束
  }

  // setSpeaking 同步播报状态（仅记录；语音演出以动作姿势表达，无口型层）
  setSpeaking(s: boolean): void {
    // 记录状态
    this.speaking = s
  // 当前代码块结束
  }

  // pulseMouth 语音边界脉冲：不再驱动口型（已按需求移除），仅低概率切换亲和动作姿势 @ticket 07
  pulseMouth(): void {
    // 未播报时忽略
    if (!this.speaking) return
    // 播报中低概率切换亲和动作姿势（人物做动作，画面静止切换）
    if (Math.random() < 0.15) {
      // 随机挑选一个亲和动作姿势
      const pool = ['wave', 'present', 'happy']
      // 过渡到随机姿势（保持 1.2s）
      this.transitionTo(pool[Math.floor(Math.random() * pool.length)], 1200)
    // 当前代码块结束

    }
  // 当前代码块结束
  }

  // stop 清理定时器（无渲染循环）
  stop(): void {
    // 清理姿势回退定时器
    if (this.backTimer) window.clearTimeout(this.backTimer)
  // 当前代码块结束
  }
// 当前代码块结束
}
