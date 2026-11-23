/** @ticket 07-数字人Coze：写实真人数字人渲染引擎（照片级分层绑骨：眨眼/张口/点头/摇头，替代扁平卡通风格）
 * 原理：以写实人像为基底层，眼睛/嘴部区域使用同图切片层（DOM 裁剪），通过对切片层做纵向缩放
 * 实现"闭眼/张口"的视觉错觉，配合呼吸起伏与光影流动，呈现立体真实的柜机引导形象。
 */
// 定义并导出类 PhotorealAvatarEngine：封装写实数字人的分层动画驱动

export class PhotorealAvatarEngine {
  // 声明属性 host：宿主容器（舞台），所有图层挂载于其下
  private host: HTMLElement
  // 声明属性 base：基底层元素（完整人像）
  private base: HTMLElement
  // 声明属性 eyeLayer：眼部分片层（眨眼动画作用对象）
  private eyeLayer: HTMLElement
  // 声明属性 mouthLayer：嘴部分片层（口型动画作用对象）
  private mouthLayer: HTMLElement
  // 声明属性 glow：播报状态光环元素（说话时点亮）
  private glow: HTMLElement
  // 声明属性 portraitSrc：当前人像图地址（角色切换时更新）
  private portraitSrc: string
  // 声明属性 mouthOpen：当前张口幅度（0~1，随语音脉冲增减）
  private mouthOpen = 0
  // 声明属性 speaking：是否处于语音播报状态
  private speaking = false
  // 声明属性 blinkTimer：自动眨眼定时器句柄
  private blinkTimer: number | null = null
  // 声明属性 rafId：呼吸动画 requestAnimationFrame 句柄
  private rafId = 0
  // 声明属性 actionUntil：当前动作（点头/摇头）的结束时间戳
  private actionUntil = 0
  // 声明属性 actionKind：当前动作类型
  private actionKind = ''

  // 人像五官区域标定（相对图片宽高的比例，可按新形象微调）@ticket 07
  // 眼部区域：左右眼整体覆盖框
  private static EYE = { x: 0.33, y: 0.245, w: 0.34, h: 0.055 }
  // 嘴部区域：嘴唇覆盖框
  private static MOUTH = { x: 0.42, y: 0.325, w: 0.16, h: 0.045 }

  // 构造函数：接收宿主容器与人像图地址，构建图层结构
  constructor(host: HTMLElement, portraitSrc: string) {
    // 赋值更新：this.host = host
    this.host = host
    // 赋值更新：this.portraitSrc = portrait
    this.portraitSrc = portraitSrc
    // 清空宿主内容，避免重复挂载产生叠影
    host.innerHTML = ''

    // 创建基底层：完整人像
    this.base = document.createElement('div')
    // 基底层样式类
    this.base.className = 'ph-base'
    // 将基底层加入宿主
    host.appendChild(this.base)

    // 创建眼部分片层：同图裁剪，眨眼时纵向压缩
    this.eyeLayer = this.buildSlice(PhotorealAvatarEngine.EYE, 'ph-eye')
    // 创建嘴部分片层：同图裁剪，说话时纵向拉伸模拟张口
    this.mouthLayer = this.buildSlice(PhotorealAvatarEngine.MOUTH, 'ph-mouth')
    // 创建播报光环元素
    this.glow = document.createElement('div')
    // 光环样式类（默认熄灭）
    this.glow.className = 'ph-glow'
    // 光环加入宿主
    host.appendChild(this.glow)

    // 应用人像背景
    this.applyPortrait()
  // 当前代码块结束
  }

  // buildSlice 构建五官切片层：以同图 background 定位实现区域裁剪
  private buildSlice(r: { x: number; y: number; w: number; h: number }, cls: string): HTMLElement {
    // 创建切片元素
    const el = document.createElement('div')
    // 设置样式类
    el.className = `ph-slice ${cls}`
    // 保存原始区域标定，供背景定位换算使用
    el.dataset.rx = String(r.x)
    el.dataset.ry = String(r.y)
    el.dataset.rw = String(r.w)
    el.dataset.rh = String(r.h)
    // 将切片加入宿主
    this.host.appendChild(el)
    // 返回切片元素
    return el
  // 当前代码块结束
  }

  // applyPortrait 将人像与区域标定应用到所有图层（角色切换/尺寸变化时调用）
  private applyPortrait(): void {
    // 读取宿主宽高，换算区域像素
    const w = this.host.clientWidth || 480
    // 读取宿主高度
    const h = this.host.clientHeight || 720
    // 基底层背景：完整人像；与切片层使用相同的 (w×h) 拉伸，保证五官切片与基底像素级对齐 @ticket 07
    this.base.style.backgroundImage = `url(${this.portraitSrc})`
    this.base.style.backgroundSize = `${w}px ${h}px`
    this.base.style.backgroundPosition = '0 0'
    // 循环处理：对眼部/嘴部切片应用背景与定位，逐层迭代
    for (const el of [this.eyeLayer, this.mouthLayer]) {
      // 读取区域标定
      const r = { x: Number(el.dataset.rx), y: Number(el.dataset.ry), w: Number(el.dataset.rw), h: Number(el.dataset.rh) }
      // 切片定位：覆盖到对应区域
      el.style.left = `${r.x * w}px`
      el.style.top = `${r.y * h}px`
      // 切片尺寸
      el.style.width = `${r.w * w}px`
      el.style.height = `${r.h * h}px`
      // 切片背景：同一人像，通过 background-size/position 对齐基底层（实现"同图切片"）
      el.style.backgroundImage = `url(${this.portraitSrc})`
      el.style.backgroundSize = `${w}px ${h}px`
      el.style.backgroundPosition = `-${r.x * w}px -${r.y * h}px`
    // 当前代码块结束
    }
  // 当前代码块结束
  }

  // start 启动渲染循环：呼吸起伏 + 自动眨眼
  start(): void {
    // 记录起始时间
    const t0 = performance.now()
    // 定义呼吸循环帧回调
    const tick = (t: number) => {
      // 计算周期相位（4s 一个呼吸周期）
      const ph = ((t - t0) / 4000) % 1
      // 呼吸起伏：幅度 1.2% 的纵向缩放 + 轻微上移
      const s = 1 + 0.012 * Math.sin(ph * Math.PI * 2)
      // 应用到基底层（含动作叠加）
      this.applyTransform(s, t)
      // 请求下一帧
      this.rafId = requestAnimationFrame(tick)
    // 当前代码块结束
    }
    // 启动循环
    this.rafId = requestAnimationFrame(tick)
    // 自动眨眼：3~6s 随机间隔
    const scheduleBlink = () => {
      // 设置下次眨眼定时
      this.blinkTimer = window.setTimeout(() => {
        // 执行一次眨眼（眼层快速压缩再展开）
        this.blink()
        // 递归调度下一次
        scheduleBlink()
      // 当前代码块结束
      }, 3000 + Math.random() * 3000)
    // 当前代码块结束
    }
    // 启动眨眼调度
    scheduleBlink()
  // 当前代码块结束
  }

  // applyTransform 计算并应用基底层变换：呼吸缩放 + 动作（点头/摇头）叠加
  private applyTransform(breath: number, t: number): void {
    // 默认无动作偏移
    let rot = 0
    let dy = 0
    // 条件判断：存在未结束的动作时计算动作偏移
    if (t < this.actionUntil) {
      // 动作进度（0→1→0 的正半弦）
      const p = Math.sin(((this.actionUntil - t) / 1200) * Math.PI)
      // 点头：前倾位移；摇头：小幅左右摆动
      if (this.actionKind === 'nod') dy = 10 * p
      else if (this.actionKind === 'shake') rot = 2.2 * p
    // 当前代码块结束
    }
    // 应用合成变换到基底层
    this.base.style.transform = `translateY(${dy}px) rotate(${rot}deg) scale(${breath})`
  // 当前代码块结束
  }

  // blink 执行一次眨眼：眼层背景上移取额头皮肤区域（同图）覆盖眼部 + 压缩，形成"眼皮落下"错觉 @ticket 07
  private blink(): void {
    // 读取宿主高度用于背景偏移换算
    const h = this.host.clientHeight || 720
    // 眼层背景上移到额头皮肤区（眼区上方 0.06h），覆盖眼睛模拟闭眼
    this.eyeLayer.style.backgroundPosition = `-${Number(this.eyeLayer.dataset.rx) * (this.host.clientWidth || 480)}px -${(Number(this.eyeLayer.dataset.ry) - 0.06) * h}px`
    // 眼层压缩模拟眼皮落下
    this.eyeLayer.style.transform = 'scaleY(0.62)'
    // 130ms 后恢复：背景回到眼区原位、形变复位
    window.setTimeout(() => {
      // 重新应用眼部切片背景（回到睁眼态）
      this.applyPortrait()
      // 清除形变
      this.eyeLayer.style.transform = ''
    // 当前代码块结束
    }, 130)
  // 当前代码块结束
  }

  // setAction 接收智能体下发的动作/情绪指令（与 Live2D 引擎接口一致）
  setAction(a: { action?: string; emotion?: string }): void {
    // 记录动作类型（nod/shake 参与变换；smile 等仅靠光环与呼吸表达）
    this.actionKind = a.action || ''
    // 动作持续 1.2s
    this.actionUntil = performance.now() + 1200
  // 当前代码块结束
  }

  // setSpeaking 同步播报状态：点亮/熄灭光环，静止时复位口型
  setSpeaking(s: boolean): void {
    // 记录状态
    this.speaking = s
    // 光环点亮/熄灭
    this.glow.classList.toggle('on', s)
    // 停止播报时口型复位
    if (!s) {
      // 张口幅度清零
      this.mouthOpen = 0
      // 恢复嘴层切片背景（闭口态）并清除形变
      this.applyPortrait()
      this.mouthLayer.style.transform = ''
    // 当前代码块结束
    }
  // 当前代码块结束
  }

  // pulseMouth 语音边界脉冲：嘴层叠加"口腔阴影"并向下拉伸，模拟张口发音（误差<0.3s 由语音边界事件保证）@ticket 07
  pulseMouth(): void {
    // 未播报时忽略脉冲
    if (!this.speaking) return
    // 张口幅度阶梯上升（封顶 0.55）
    this.mouthOpen = Math.min(0.55, this.mouthOpen + 0.28)
    // 嘴层背景替换为口腔阴影渐变（深色椭圆），随张口幅度显现
    this.mouthLayer.style.backgroundImage = 'radial-gradient(ellipse at 50% 45%, rgba(70,32,30,0.85) 0%, rgba(90,40,36,0.55) 55%, transparent 78%)'
    this.mouthLayer.style.backgroundSize = '100% 100%'
    this.mouthLayer.style.backgroundPosition = '0 0'
    // 嘴层向下拉伸 + 轻微下移，形成开口动作
    this.mouthLayer.style.transform = `scaleY(${1 + this.mouthOpen}) translateY(${this.mouthOpen * 30}%)`
    // 180ms 后回落一半，形成开合节奏
    window.setTimeout(() => {
      // 幅度衰减
      this.mouthOpen = Math.max(0, this.mouthOpen - 0.2)
      // 幅度归零时恢复切片背景（闭口态），否则维持阴影并收缩
      if (this.mouthOpen <= 0.02) {
        // 恢复切片背景
        this.applyPortrait()
        // 清除形变
        this.mouthLayer.style.transform = ''
      // 上述条件不成立时，执行 else 分支逻辑

      } else {
        // 应用衰减后的口型
        this.mouthLayer.style.transform = `scaleY(${1 + this.mouthOpen}) translateY(${this.mouthOpen * 30}%)`
      // 当前代码块结束

      }
    // 当前代码块结束
    }, 180)
  // 当前代码块结束
  }

  // setPortrait 切换角色人像（品牌代言人/虚拟讲解员）
  setPortrait(src: string): void {
    // 更新图址
    this.portraitSrc = src
    // 重新应用图层背景与定位
    this.applyPortrait()
  // 当前代码块结束
  }

  // stop 停止动画循环并清理定时器
  stop(): void {
    // 取消渲染循环
    cancelAnimationFrame(this.rafId)
    // 清理眨眼定时器
    if (this.blinkTimer) window.clearTimeout(this.blinkTimer)
  // 当前代码块结束
  }
// 当前代码块结束
}
