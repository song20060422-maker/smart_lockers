/** @ticket 07-数字人Coze：精致数字人形象引擎——纯浏览器原生 Canvas 2D 绘制，免商业 SDK / 免授权凭证，支持呼吸眨眼待机与口型随语音播报开合 */

// 从统一类型模块导入 Live2DAction 类型，用于接收工单既有的表情/动作 payload

import type { Live2DAction } from '@/types'

// 定义并导出 CartoonAvatarEngine 类，封装数字人的绘制、动画与说话口型驱动

export class CartoonAvatarEngine {

  // 声明私有字段 canvas，保存待绘制的 canvas 元素引用

  private canvas: HTMLCanvasElement

  // 声明私有字段 ctx，保存 2D 绘图上下文

  private ctx: CanvasRenderingContext2D

  // 声明私有字段 host，保存渲染容器元素，用于监听尺寸变化

  private host: HTMLElement

  // 声明私有字段 raf，保存 requestAnimationFrame 句柄以便取消

  private raf = 0

  // 声明私有字段 running，标识动画循环是否正在运行

  private running = false

  // 声明私有字段 w，保存画布的 CSS 逻辑宽度

  private w = 0

  // 声明私有字段 h，保存画布的 CSS 逻辑高度

  private h = 0

  // 声明私有字段 dpr，保存设备像素比（上限 2，兼顾清晰与性能）

  private dpr = 1

  // 声明私有字段 observer，保存尺寸监听器实例

  private observer: ResizeObserver | null = null

  // 声明私有字段 action，保存当前最新的动作 payload

  private action: Live2DAction | null = null

  // 声明私有字段 emotion，保存当前情绪字符串

  private emotion = 'neutral'

  // 声明私有字段 duration，保存当前动作的持续时间（毫秒）

  private duration = 0

  // 声明私有字段 actionStart，保存当前动作开始的时间戳

  private actionStart = 0

  // 声明私有字段 lastBlinkStart，保存最近一次眨眼的起始时间戳

  private lastBlinkStart = -1000

  // 声明私有字段 nextBlink，保存下一次计划眨眼的时间戳

  private nextBlink = 0

  // 声明私有字段 speaking，标识当前是否处于语音播报状态（驱动说话口型）

  private speaking = false

  // 声明私有字段 speakStart，保存本轮播报起始时间戳，作为口型节奏相位基准

  private speakStart = 0

  // 声明私有字段 lastBoundary，保存最近一次语音边界事件的触发时间戳

  private lastBoundary = 0

  // 声明私有字段 gazeSeed，保存眼神微动的相位种子，使左右眼自然同步

  private gazeSeed = Math.random() * 1000

  // 定义构造函数，接收 canvas 与宿主容器并完成绘图上下文初始化

  constructor(canvas: HTMLCanvasElement, host: HTMLElement) {

    // 赋值更新：保存传入的 canvas 元素

    this.canvas = canvas

    // 赋值更新：保存传入的宿主容器

    this.host = host

    // 声明常量 ctx，尝试获取 2D 绘图上下文

    const ctx = canvas.getContext('2d')

    // 条件判断：若无法获取绘图上下文则抛出错误，避免后续空指针

    if (!ctx) {

      // 抛出错误：明确告知调用方当前环境不支持 2D 绘图

      throw new Error('CartoonAvatarEngine: 无法获取 2D 绘图上下文')

    }

    // 赋值更新：保存绘图上下文

    this.ctx = ctx

  // 构造函数结束

  }

  // 定义方法 resize，依据容器尺寸与设备像素比同步画布像素尺寸

  private resize = () => {

    // 声明变量 rect，读取画布当前在页面中的布局矩形

    const rect = this.canvas.getBoundingClientRect()

    // 声明变量 rawW，优先取画布布局宽度，回退到容器宽度与默认值

    const rawW = Math.round(rect.width || this.host.clientWidth || 480)

    // 声明变量 rawH，优先取画布布局高度，回退到容器高度与默认值

    const rawH = Math.round(rect.height || this.host.clientHeight || 480)

    // 赋值更新：修正宽度为至少 1 像素，避免零尺寸画布

    this.w = rawW < 1 ? 480 : rawW

    // 赋值更新：修正高度为至少 1 像素，避免零尺寸画布

    this.h = rawH < 1 ? 480 : rawH

    // 赋值更新：设备像素比取系统值与 2 的较小者，防止超高分辨率拖慢渲染

    this.dpr = Math.min(window.devicePixelRatio || 1, 2)

    // 赋值更新：按设备像素比设置画布真实像素宽度

    this.canvas.width = Math.round(this.w * this.dpr)

    // 赋值更新：按设备像素比设置画布真实像素高度

    this.canvas.height = Math.round(this.h * this.dpr)

    // 调用方法：重置绘图上下文变换矩阵，使后续绘制沿用 CSS 逻辑坐标

    this.ctx.setTransform(this.dpr, 0, 0, this.dpr, 0, 0)

  // resize 方法结束

  }

  // 定义方法 start，启动动画循环并挂载尺寸监听

  start = () => {

    // 条件判断：若已在运行则直接返回，避免重复启动多个循环

    if (this.running) {

      // 提前返回：保持单一动画循环

      return

    }

    // 赋值更新：切换运行标记为 true

    this.running = true

    // 样式设置：让画布宽度填满容器，CSS 自适应由 flex 布局提供

    this.canvas.style.width = '100%'

    // 样式设置：让画布高度填满容器

    this.canvas.style.height = '100%'

    // 样式设置：去除行内元素底部空隙，保证画布贴合容器

    this.canvas.style.display = 'block'

    // 调用方法：首次同步画布尺寸

    this.resize()

    // 赋值更新：创建尺寸监听器，容器尺寸变化时自动调用 resize

    this.observer = new ResizeObserver(this.resize)

    // 调用方法：开始观察宿主容器尺寸

    this.observer.observe(this.host)

    // 赋值更新：初始化首次自动眨眼时间（2 秒后）

    this.nextBlink = performance.now() + 2000

    // 定义函数 loop，作为动画帧回调持续绘制

    const loop = (now: number) => {

      // 条件判断：若运行标记已关闭则终止循环

      if (!this.running) {

        // 提前返回：不再请求下一帧

        return

      }

      // 调用方法：绘制当前帧

      this.draw(now)

      // 调用方法：请求下一帧，形成持续动画

      this.raf = requestAnimationFrame(loop)

    // loop 函数结束

    }

    // 调用方法：请求首帧，正式启动动画

    this.raf = requestAnimationFrame(loop)

  // start 方法结束

  }

  // 定义方法 stop，停止动画循环并释放监听器

  stop = () => {

    // 赋值更新：切换运行标记为 false，令动画帧回调自然退出

    this.running = false

    // 调用方法：取消尚未执行的动画帧

    cancelAnimationFrame(this.raf)

    // 条件判断：若存在尺寸监听器则断开观察

    if (this.observer) {

      // 调用方法：断开尺寸监听

      this.observer.disconnect()

      // 赋值更新：清空监听器引用

      this.observer = null

    }

  // stop 方法结束

  }

  // 定义方法 setAction，接收工单 payload 驱动表情与动作表现

  setAction = (action: Live2DAction | null) => {

    // 声明变量 prev，保存更新前的动作引用用于判断是否为新动作

    const prev = this.action

    // 赋值更新：保存最新动作

    this.action = action

    // 条件判断：若收到有效动作

    if (action) {

      // 赋值更新：同步情绪字段，缺省视为 neutral

      this.emotion = action.emotion ?? 'neutral'

      // 赋值更新：同步持续时间，缺省 1200 毫秒

      this.duration = action.duration_ms ?? 1200

    } else {

      // 赋值更新：无动作时情绪复位为 neutral

      this.emotion = 'neutral'

      // 赋值更新：无动作时持续时间为 0

      this.duration = 0

    }

    // 条件判断：仅当动作引用真正变化时才重置起始时间，避免同帧重复触发

    if (action !== prev) {

      // 赋值更新：记录新动作的起始时间戳

      this.actionStart = performance.now()

    }

  // setAction 方法结束

  }

  // 定义方法 setSpeaking，由语音播报状态驱动说话口型与声波光环

  setSpeaking = (speaking: boolean) => {

    // 条件判断：仅当状态发生翻转时才更新时间基准，避免同帧重复重置相位

    if (speaking !== this.speaking) {

      // 赋值更新：记录本轮播报起始时间，作为口型节奏的相位基准

      this.speakStart = performance.now()

    }

    // 赋值更新：保存最新播报状态

    this.speaking = speaking

  // setSpeaking 方法结束

  }

  // 定义方法 pulseMouth，在语音边界事件触发时注入一次张口脉冲，使口型贴合真实发音

  pulseMouth = () => {

    // 赋值更新：记录本次语音边界时间戳，绘制时据此生成短促张口包络

    this.lastBoundary = performance.now()

  // pulseMouth 方法结束

  }

  // 定义方法 draw，根据时间相位计算动态参数并绘制一帧画面

  private draw = (now: number) => {

    // 声明变量 ctx，取出绘图上下文局部引用

    const ctx = this.ctx

    // 调用方法：清除上一帧绘制内容，透明区域露出容器的渐变背景

    ctx.clearRect(0, 0, this.w, this.h)

    // 声明变量 t，将时间戳换算为秒，作为呼吸与摆动的相位基准

    const t = now / 1000

    // 声明变量 action，读取当前动作 payload

    const action = this.action

    // 声明变量 kind，读取动作类型，缺省 idle

    const kind = action?.action ?? 'idle'

    // 声明变量 active，判断当前动作是否仍处于持续时间内

    const active = !!action && (now - this.actionStart) < this.duration

    // 声明变量 progress，计算动作归一化进度（0 到 1）

    const progress = active ? Math.min(1, (now - this.actionStart) / this.duration) : 0

    // 声明变量 happy，判断当前情绪是否属于开心阵营

    const happy = this.emotion === 'happy' || this.emotion === 'positive' || kind === 'smile'

    // 声明变量 sad，判断当前情绪是否为难过

    const sad = this.emotion === 'sad'

    // 声明变量 breath，用正弦波生成呼吸起伏系数

    const breath = Math.sin(t * 1.6)

    // 声明变量 sway，用慢速正弦波生成身体轻微摆动角度

    const sway = Math.sin(t * 0.7) * 0.02

    // 条件判断：到达计划眨眼时间时启动一次新的眨眼

    if (now >= this.nextBlink) {

      // 赋值更新：记录本次眨眼起始时间

      this.lastBlinkStart = now

      // 赋值更新：安排下一次眨眼（2.4 到 5.4 秒之间随机）

      this.nextBlink = now + 2400 + Math.random() * 3000

    }

    // 声明变量 blink，按眨眼进度生成 0 到 1 的闭眼程度（正弦包络）

    const blink = Math.max(

      // 计算自动眨眼包络：眨眼窗口内呈半周期正弦

      this.autoBlink(now),

      // 计算动作眨眼：blink 动作时按进度生成闭眼包络

      active && kind === 'blink' ? Math.sin(progress * Math.PI) : 0,

    )

    // 声明变量 headBob，nod 动作时让头部上下点头

    const headBob = active && kind === 'nod' ? Math.sin(progress * Math.PI * 2) * 10 : 0

    // 声明变量 headTurn，shake 动作时让头部左右摇头

    const headTurn = active && kind === 'shake' ? Math.sin(progress * Math.PI * 4) * 0.13 : 0

    // 声明变量 transformScale，transform 动作时整体短暂放大

    const transformScale = active && kind === 'transform' ? 1 + Math.sin(progress * Math.PI) * 0.06 : 1

    // 调用方法：绘制浅蓝柔和背景光晕与漂浮光点，营造通透智能的视觉基调

    this.drawBackdrop(ctx, now)

    // 声明变量 cx，取得画布水平中心作为角色锚点

    const cx = this.w / 2

    // 声明变量 cy，取得画布垂直中心作为角色锚点

    const cy = this.h / 2

    // 声明变量 scale，按画布较小边计算整体缩放基准（基准高度 440）

    const scale = Math.min(this.w, this.h) / 440 * transformScale

    // 调用方法：保存当前绘图状态，准备绘制角色主体

    ctx.save()

    // 调用方法：将坐标原点移动到角色锚点

    ctx.translate(cx, cy)

    // 调用方法：按缩放基准放大坐标系，使角色尺寸随画布自适应

    ctx.scale(scale, scale)

    // 调用方法：按摆动角度轻微旋转，形成左右摇摆的空闲动态

    ctx.rotate(sway)

    // 调用方法：按呼吸系数对整个角色做纵向缩放，形成起伏感

    ctx.scale(1, 1 + breath * 0.012)

    // 声明变量 mouthOpen，计算当前口型张开程度（0 到 1）

    const mouthOpen = this.computeMouthOpen(now, happy)

    // 调用方法：绘制说话时围绕头部的声波光环

    this.drawAura(ctx, now, mouthOpen)

    // 调用方法：绘制角色身体部分（简约浅色上衣与品牌色点缀）

    this.drawBody(ctx, happy, sad, breath)

    // 调用方法：绘制脖颈

    this.drawNeck(ctx)

    // 调用方法：绘制后层头发（位于脸型之后）

    this.drawBackHair(ctx)

    // 调用方法：绘制脸型与耳朵

    this.drawFace(ctx)

    // 调用方法：绘制眉毛

    this.drawBrows(ctx, happy, sad)

    // 调用方法：绘制眼睛（含眨眼、眼神微动与高光）

    this.drawEyes(ctx, blink, happy, now)

    // 调用方法：绘制鼻部柔和阴影

    this.drawNose(ctx)

    // 调用方法：绘制嘴巴（含说话开合口型）

    this.drawMouth(ctx, happy, sad, mouthOpen)

    // 调用方法：绘制两侧柔和腮红

    this.drawBlush(ctx, happy)

    // 调用方法：绘制前层刘海与发丝高光

    this.drawFrontHair(ctx)

    // 调用方法：恢复绘图状态，结束角色主体绘制

    ctx.restore()

  // draw 方法结束

  }

  // 定义方法 computeMouthOpen，综合播报状态、音节节奏与语音边界脉冲计算口型张开程度

  private computeMouthOpen(now: number, happy: boolean): number {

    // 条件判断：未处于播报状态时返回闭合口型（保留开心表情的轻微上扬由绘制层处理）

    if (!this.speaking) {

      // 返回值：完全闭合

      return 0

    }

    // 声明变量 elapsed，计算自本轮播报开始经过的秒数

    const elapsed = (now - this.speakStart) / 1000

    // 声明变量 syllable，用复合正弦模拟连续音节的张合节奏，避免机械式匀速开合

    const syllable = 0.5 + 0.5 * Math.sin(elapsed * 13 + Math.sin(elapsed * 4.3) * 1.6)

    // 声明变量 sinceBoundary，计算距离最近一次语音边界的时间差

    const sinceBoundary = now - this.lastBoundary

    // 声明变量 boundaryPulse，在边界后 180 毫秒内生成一次短促张口脉冲

    const boundaryPulse = sinceBoundary >= 0 && sinceBoundary < 180

      ? Math.sin((1 - sinceBoundary / 180) * Math.PI / 2)

      : 0

    // 声明变量 base，开心情绪下说话幅度略柔和，其余保持自然幅度

    const base = happy ? 0.72 : 0.85

    // 返回值：叠加节奏、脉冲与基础张口并限制在 0 到 1 区间

    return Math.min(1, 0.18 + base * (0.45 * syllable + 0.55 * boundaryPulse))

  // computeMouthOpen 方法结束

  }

  // 定义方法 autoBlink，按眨眼窗口返回平滑的闭眼包络值

  private autoBlink(now: number): number {

    // 声明变量 since，计算距离最近一次眨眼起始的时间差

    const since = now - this.lastBlinkStart

    // 条件判断：仅在 150 毫秒的眨眼窗口内返回正闭眼值

    if (since >= 0 && since < 150) {

      // 返回值：半周期正弦包络，闭眼再睁眼

      return Math.sin((since / 150) * Math.PI)

    }

    // 返回值：窗口外保持睁眼

    return 0

  // autoBlink 方法结束

  }

  // 定义方法 drawBackdrop，绘制干净通透的浅蓝背景与缓慢漂浮的柔光粒子
  private drawBackdrop(ctx: CanvasRenderingContext2D, now: number) {

    // 声明变量 t，换算时间戳为秒作为漂浮相位
    const t = now / 1000

    // 声明变量 g，创建以画布中上部为中心的浅蓝径向渐变作为通透背景
    const g = ctx.createRadialGradient(this.w * 0.5, this.h * 0.4, 8, this.w * 0.5, this.h * 0.4, Math.max(this.w, this.h) * 0.75)

    // 渐变停靠点：中心使用柔和近白
    g.addColorStop(0, 'rgba(255,255,255,0.98)')

    // 渐变停靠点：中段使用干净浅蓝
    g.addColorStop(0.5, 'rgba(226,240,255,0.85)')

    // 渐变停靠点：外缘使用更浅的冰蓝
    g.addColorStop(1, 'rgba(214,233,252,0.35)')

    // 设置填充样式：应用背景渐变
    ctx.fillStyle = g

    // 调用方法：填充整个画布形成干净背景
    ctx.fillRect(0, 0, this.w, this.h)

    // 声明变量 halo，创建底部柔和光晕，托出角色
    const halo = ctx.createRadialGradient(this.w * 0.5, this.h * 0.78, 4, this.w * 0.5, this.h * 0.78, this.h * 0.34)

    // 渐变停靠点：光晕中心使用极淡品牌蓝
    halo.addColorStop(0, 'rgba(120,170,235,0.16)')

    // 渐变停靠点：光晕外缘淡出为透明
    halo.addColorStop(1, 'rgba(120,170,235,0)')

    // 设置填充样式：应用底部光晕
    ctx.fillStyle = halo

    // 调用方法：填充底部光晕
    ctx.fillRect(0, 0, this.w, this.h)

    // 调用方法：循环绘制 10 个缓慢漂浮的柔光粒子，保持画面干净不喧宾夺主
    for (let i = 0; i < 10; i++) {

      // 声明变量 px，用正弦让粒子水平缓慢漂移
      const px = (0.1 + 0.8 * ((i * 0.173 + 0.07) % 1)) * this.w + Math.sin(t * 0.35 + i) * 10

      // 声明变量 py，用余弦让粒子垂直缓慢漂浮
      const py = (0.12 + 0.76 * ((i * 0.233 + 0.04) % 1)) * this.h + Math.cos(t * 0.28 + i) * 14

      // 声明变量 r，按序号生成 1.5 到 4 像素的粒子半径
      const r = 1.5 + (i % 4)

      // 调用方法：开始绘制粒子路径
      ctx.beginPath()

      // 调用方法：绘制粒子圆形
      ctx.arc(px, py, r, 0, Math.PI * 2)

      // 设置填充样式：粒子使用低透明品牌蓝，柔和不刺眼
      ctx.fillStyle = 'rgba(90,150,225,0.16)'

      // 调用方法：填充粒子
      ctx.fill()

    }

  // drawBackdrop 方法结束

  }

  // 定义方法 drawAura，在说话时绘制围绕头部的柔和声波光环

  private drawAura(ctx: CanvasRenderingContext2D, now: number, mouthOpen: number) {

    // 条件判断：仅在播报状态绘制声波光环

    if (!this.speaking) {

      // 提前返回：待机时保持画面简洁

      return

    }

    // 声明变量 t，换算时间戳为秒作为声波扩散相位

    const t = now / 1000

    // 调用方法：循环绘制三层扩散声波圆环

    for (let i = 0; i < 3; i++) {

      // 声明变量 phase，计算每层圆环的扩散进度（0 到 1）

      const phase = (t * 0.7 + i / 3) % 1

      // 声明变量 radius，随扩散进度增大环半径

      const radius = 96 + phase * 70

      // 声明变量 alpha，随扩散进度淡出且受口型强度调制

      const alpha = (1 - phase) * 0.28 * (0.5 + mouthOpen * 0.5)

      // 调用方法：开始绘制声波圆环路径

      ctx.beginPath()

      // 调用方法：以头部中心绘制圆环

      ctx.arc(0, -58, radius, 0, Math.PI * 2)

      // 设置描边样式：使用带透明度的品牌蓝

      ctx.strokeStyle = `rgba(74,143,231,${alpha.toFixed(3)})`

      // 设置线宽：声波环保持纤细

      ctx.lineWidth = 2

      // 调用方法：描边声波圆环

      ctx.stroke()

    }

  // drawAura 方法结束

  }

  // 定义方法 drawBody，绘制简约奶白上衣、浅蓝圆领与品牌蓝丝带点缀
  private drawBody(ctx: CanvasRenderingContext2D, happy: boolean, sad: boolean, breath: number) {

    // 调用方法：保存绘图状态，隔离身体绘制变换

    ctx.save()

    // 声明变量 g，创建上衣的纵向渐变（肩部明亮向底部柔和浅蓝过渡）

    const g = ctx.createLinearGradient(0, 60, 0, 220)

    // 渐变停靠点：肩部使用明亮奶白

    g.addColorStop(0, happy ? '#ffffff' : sad ? '#eef1f6' : '#fdfbf7')

    // 渐变停靠点：中段使用柔和浅蓝白

    g.addColorStop(0.55, happy ? '#f4f9ff' : sad ? '#dfe5ee' : '#eef3fa')

    // 渐变停靠点：底部使用更实的浅蓝增加层次

    g.addColorStop(1, happy ? '#e4eefb' : sad ? '#cbd3df' : '#dbe8f7')

    // 设置填充样式：应用上衣渐变

    ctx.fillStyle = g

    // 调用方法：开始绘制圆润肩部轮廓（无硬朗领型）

    ctx.beginPath()

    // 调用方法：定位肩部左下起点

    ctx.moveTo(-158, 220)

    // 调用方法：用二次贝塞尔曲线描绘左肩到领口

    ctx.quadraticCurveTo(-150, 86, -60, 62)

    // 调用方法：用柔和弧线描绘领口凹陷

    ctx.quadraticCurveTo(0, 78, 60, 62)

    // 调用方法：用二次贝塞尔曲线描绘右肩

    ctx.quadraticCurveTo(150, 86, 158, 220)

    // 调用方法：闭合肩部路径

    ctx.closePath()

    // 调用方法：填充肩部轮廓

    ctx.fill()

    // 调用方法：保存状态，绘制肩部柔和受光高光，增加布料立体层次

    ctx.save()

    // 声明变量 shoulderGlow，创建肩部左侧的斜向高光渐变

    const shoulderGlow = ctx.createLinearGradient(-150, 80, 40, 200)

    // 渐变停靠点：高光起点使用半透明白

    shoulderGlow.addColorStop(0, 'rgba(255,255,255,0.55)')

    // 渐变停靠点：高光中段快速淡出

    shoulderGlow.addColorStop(0.5, 'rgba(255,255,255,0.06)')

    // 渐变停靠点：高光末端完全透明

    shoulderGlow.addColorStop(1, 'rgba(255,255,255,0)')

    // 设置填充样式：应用肩部高光渐变

    ctx.fillStyle = shoulderGlow

    // 调用方法：以肩部轮廓作为高光的裁剪形状

    ctx.beginPath()

    // 调用方法：重描肩部左起点

    ctx.moveTo(-158, 220)

    // 调用方法：重描左肩曲线

    ctx.quadraticCurveTo(-150, 86, -60, 62)

    // 调用方法：重描领口弧线

    ctx.quadraticCurveTo(0, 78, 60, 62)

    // 调用方法：重描右肩曲线

    ctx.quadraticCurveTo(150, 86, 158, 220)

    // 调用方法：闭合高光区域

    ctx.closePath()

    // 调用方法：填充肩部高光

    ctx.fill()

    // 调用方法：恢复高光绘制状态

    ctx.restore()

    // 调用方法：循环绘制简洁衣褶，让衣料不平板且随呼吸轻微起伏

    for (let i = 0; i < 2; i++) {

      // 声明变量 foldY，按序号计算褶皱起始高度

      const foldY = 118 + i * 34

      // 声明变量 foldX，左侧褶皱随呼吸轻微起伏

      const foldX = -120 + i * 8 + breath * 1.2

      // 调用方法：开始绘制左侧褶皱

      ctx.beginPath()

      // 调用方法：定位左侧褶皱起点

      ctx.moveTo(foldX, foldY)

      // 调用方法：用贝塞尔描绘褶皱的自然弯曲

      ctx.bezierCurveTo(foldX + 24, foldY + 10, foldX + 16, foldY + 30, foldX + 44, foldY + 48)

      // 设置描边样式：褶皱使用低透明蓝灰表现布纹阴影

      ctx.strokeStyle = 'rgba(120,150,190,0.16)'

      // 设置线宽：褶皱保持纤细

      ctx.lineWidth = 2

      // 设置线帽：圆头褶皱

      ctx.lineCap = 'round'

      // 调用方法：描边左侧褶皱

      ctx.stroke()

      // 调用方法：开始绘制右侧对称褶皱

      ctx.beginPath()

      // 调用方法：定位右侧褶皱起点

      ctx.moveTo(-foldX, foldY)

      // 调用方法：用贝塞尔描绘右侧褶皱

      ctx.bezierCurveTo(-foldX - 24, foldY + 10, -foldX - 16, foldY + 30, -foldX - 44, foldY + 48)

      // 调用方法：描边右侧褶皱

      ctx.stroke()

    }

    // 声明变量 collarG，创建浅蓝罗纹圆领的渐变

    const collarG = ctx.createLinearGradient(0, 62, 0, 96)

    // 渐变停靠点：领口顶部使用柔和浅蓝

    collarG.addColorStop(0, '#cfe0f6')

    // 渐变停靠点：领口底部使用略深品牌蓝

    collarG.addColorStop(1, '#a9c8ee')

    // 设置描边样式：应用圆领渐变

    ctx.strokeStyle = collarG

    // 设置线宽：圆领保持适中宽度体现罗纹质感

    ctx.lineWidth = 9

    // 调用方法：开始绘制圆领弧线

    ctx.beginPath()

    // 调用方法：定位圆领左端

    ctx.moveTo(-44, 66)

    // 调用方法：用二次曲线描绘圆领弧度

    ctx.quadraticCurveTo(0, 90, 44, 66)

    // 调用方法：描边圆领

    ctx.stroke()

    // 设置描边样式：领口内侧使用柔和暖灰表现深度

    ctx.strokeStyle = 'rgba(150,120,100,0.16)'

    // 设置线宽：内阴影保持纤细

    ctx.lineWidth = 2

    // 调用方法：开始绘制领口内阴影

    ctx.beginPath()

    // 调用方法：定位内阴影左端

    ctx.moveTo(-40, 70)

    // 调用方法：用二次曲线描绘内阴影弧度

    ctx.quadraticCurveTo(0, 92, 40, 70)

    // 调用方法：描边内阴影

    ctx.stroke()

    // 调用方法：保存状态，绘制品牌蓝丝带蝴蝶结并随呼吸轻微起伏

    ctx.save()

    // 调用方法：将原点移动到领口中央

    ctx.translate(0, 84 + breath * 0.8)

    // 声明变量 ribbonG，创建丝带渐变（增加垂坠质感）

    const ribbonG = ctx.createLinearGradient(-18, 0, 18, 0)

    // 渐变停靠点：丝带左侧使用明亮品牌蓝

    ribbonG.addColorStop(0, happy ? '#63a8f2' : sad ? '#93a0b4' : '#4a8fe7')

    // 渐变停靠点：丝带右侧使用沉稳品牌蓝

    ribbonG.addColorStop(1, happy ? '#3f83e0' : sad ? '#6b7789' : '#2f6fd0')

    // 设置填充样式：应用丝带渐变

    ctx.fillStyle = ribbonG

    // 调用方法：绘制左侧蝴蝶结

    ctx.beginPath()

    // 调用方法：定位左结中心

    ctx.moveTo(0, 0)

    // 调用方法：描绘左结上缘

    ctx.quadraticCurveTo(-18, -10, -18, 2)

    // 调用方法：描绘左结下缘收拢

    ctx.quadraticCurveTo(-18, 12, 0, 4)

    // 调用方法：闭合左结

    ctx.closePath()

    // 调用方法：填充左结

    ctx.fill()

    // 调用方法：绘制右侧蝴蝶结

    ctx.beginPath()

    // 调用方法：定位右结中心

    ctx.moveTo(0, 0)

    // 调用方法：描绘右结上缘

    ctx.quadraticCurveTo(18, -10, 18, 2)

    // 调用方法：描绘右结下缘收拢

    ctx.quadraticCurveTo(18, 12, 0, 4)

    // 调用方法：闭合右结

    ctx.closePath()

    // 调用方法：填充右结

    ctx.fill()

    // 调用方法：绘制蝴蝶结中心结

    ctx.beginPath()

    // 调用方法：以圆形描绘中心结

    ctx.arc(0, 2, 4.5, 0, Math.PI * 2)

    // 调用方法：填充中心结

    ctx.fill()

    // 调用方法：绘制左侧飘带

    ctx.beginPath()

    // 调用方法：定位左飘带起点

    ctx.moveTo(-3, 5)

    // 调用方法：描绘左飘带外缘

    ctx.quadraticCurveTo(-7, 20, -12, 26)

    // 调用方法：描绘左飘带内缘收拢

    ctx.quadraticCurveTo(-4, 22, -1, 8)

    // 调用方法：闭合左飘带

    ctx.closePath()

    // 调用方法：填充左飘带

    ctx.fill()

    // 调用方法：绘制右侧飘带

    ctx.beginPath()

    // 调用方法：定位右飘带起点

    ctx.moveTo(3, 5)

    // 调用方法：描绘右飘带外缘

    ctx.quadraticCurveTo(7, 18, 12, 22)

    // 调用方法：描绘右飘带内缘收拢

    ctx.quadraticCurveTo(4, 20, 1, 8)

    // 调用方法：闭合右飘带

    ctx.closePath()

    // 调用方法：填充右飘带

    ctx.fill()

    // 调用方法：恢复丝带绘制状态

    ctx.restore()

    // 调用方法：绘制胸前小小的品牌色圆点徽章，简约不刺眼

    ctx.beginPath()

    // 调用方法：以圆形描绘徽章外圈

    ctx.arc(66, 150, 5, 0, Math.PI * 2)

    // 设置填充样式：徽章外圈使用低饱和品牌蓝

    ctx.fillStyle = 'rgba(74,143,231,0.55)'

    // 调用方法：填充徽章外圈

    ctx.fill()

    // 调用方法：绘制徽章高光点

    ctx.beginPath()

    // 调用方法：以小圆描绘徽章高光

    ctx.arc(66, 150, 2, 0, Math.PI * 2)

    // 设置填充样式：高光使用近白

    ctx.fillStyle = 'rgba(255,255,255,0.85)'

    // 调用方法：填充徽章高光

    ctx.fill()

    // 调用方法：恢复绘图状态，结束身体绘制

    ctx.restore()

  // drawBody 方法结束

  }

  // 定义方法 drawNeck，绘制白皙脖颈及其柔和阴影
  private drawNeck(ctx: CanvasRenderingContext2D) {

    // 调用方法：保存绘图状态，隔离脖颈绘制

    ctx.save()

    // 声明变量 g，创建脖颈的横向圆柱渐变（两侧深、中间亮，制造立体感）

    const g = ctx.createLinearGradient(-26, 0, 26, 0)

    // 渐变停靠点：左侧使用深肤色作阴影

    g.addColorStop(0, '#efc3a3')

    // 渐变停靠点：中部使用明亮白皙肤色

    g.addColorStop(0.5, '#ffe0c6')

    // 渐变停靠点：右侧使用稍深肤色

    g.addColorStop(1, '#ecbe9c')

    // 设置填充样式：应用脖颈渐变

    ctx.fillStyle = g

    // 调用方法：开始绘制脖颈路径（上宽下略收的圆润形状）

    ctx.beginPath()

    // 调用方法：定位脖颈左上

    ctx.moveTo(-24, 18)

    // 调用方法：描绘左侧收拢曲线

    ctx.bezierCurveTo(-26, 40, -24, 62, -20, 76)

    // 调用方法：描绘脖颈底部圆润过渡

    ctx.quadraticCurveTo(0, 84, 20, 76)

    // 调用方法：描绘右侧收拢曲线

    ctx.bezierCurveTo(24, 62, 26, 40, 24, 18)

    // 调用方法：描绘脖颈顶部弧线

    ctx.quadraticCurveTo(0, 28, -24, 18)

    // 调用方法：闭合脖颈路径

    ctx.closePath()

    // 调用方法：填充脖颈

    ctx.fill()

    // 声明变量 shadow，创建下颌落在颈部的柔和投影渐变

    const shadow = ctx.createLinearGradient(0, 16, 0, 48)

    // 渐变停靠点：投影顶部使用暖棕色

    shadow.addColorStop(0, 'rgba(140,90,66,0.24)')

    // 渐变停靠点：投影底部淡出为透明

    shadow.addColorStop(1, 'rgba(140,90,66,0)')

    // 设置填充样式：应用下颌投影渐变

    ctx.fillStyle = shadow

    // 调用方法：以椭圆描绘下颌投影范围

    ctx.beginPath()

    // 调用方法：绘制投影椭圆

    ctx.ellipse(0, 30, 26, 20, 0, 0, Math.PI * 2)

    // 调用方法：填充下颌投影

    ctx.fill()

    // 声明变量 fade，创建脖颈底部的柔和淡出，让皮肤自然融入圆领

    const fade = ctx.createLinearGradient(0, 60, 0, 84)

    // 渐变停靠点：淡出顶部完全透明

    fade.addColorStop(0, 'rgba(255,224,198,0)')

    // 渐变停靠点：淡出底部使用柔和肤色过渡

    fade.addColorStop(1, 'rgba(255,224,198,0.55)')

    // 设置填充样式：应用底部过渡渐变

    ctx.fillStyle = fade

    // 调用方法：填充脖颈底部过渡区域

    ctx.fillRect(-26, 60, 52, 24)

    // 调用方法：恢复绘图状态，结束脖颈绘制

    ctx.restore()

  // drawNeck 方法结束

  }

  // 定义方法 drawBackHair，绘制垂至肩下的柔顺深棕长发后层
  private drawBackHair(ctx: CanvasRenderingContext2D) {

    // 调用方法：保存绘图状态，隔离后层头发绘制

    ctx.save()

    // 声明变量 g，创建深棕发色渐变（顶部暖棕高光向发尾深棕过渡）

    const g = ctx.createLinearGradient(0, -170, 0, 130)

    // 渐变停靠点：顶部使用暖棕高光

    g.addColorStop(0, '#5a3d2e')

    // 渐变停靠点：中部使用深棕主色

    g.addColorStop(0.45, '#3d2820')

    // 渐变停靠点：发尾使用更深的棕黑

    g.addColorStop(1, '#241813')

    // 设置填充样式：应用发色渐变

    ctx.fillStyle = g

    // 调用方法：开始绘制后层长发整体轮廓（头顶圆弧加两侧垂肩长发绺）

    ctx.beginPath()

    // 调用方法：定位左侧发际起点

    ctx.moveTo(-92, -10)

    // 调用方法：用贝塞尔曲线描绘左侧发际至头顶

    ctx.bezierCurveTo(-108, -120, -72, -168, 0, -168)

    // 调用方法：用贝塞尔曲线描绘右侧发际

    ctx.bezierCurveTo(72, -168, 108, -120, 92, -10)

    // 调用方法：描绘右侧垂肩长发绺外缘

    ctx.quadraticCurveTo(100, 70, 86, 132)

    // 调用方法：描绘右侧发梢圆润收拢

    ctx.quadraticCurveTo(74, 150, 62, 128)

    // 调用方法：描绘右侧发绺内缘回到发际

    ctx.quadraticCurveTo(66, 40, 52, -6)

    // 调用方法：描绘底部发丝过渡连接左侧

    ctx.quadraticCurveTo(0, 30, -52, -6)

    // 调用方法：描绘左侧发绺内缘

    ctx.quadraticCurveTo(-66, 40, -62, 128)

    // 调用方法：描绘左侧发梢圆润收拢

    ctx.quadraticCurveTo(-74, 150, -86, 132)

    // 调用方法：闭合左侧发绺外缘回到起点

    ctx.quadraticCurveTo(-100, 70, -92, -10)

    // 调用方法：闭合后层头发路径

    ctx.closePath()

    // 调用方法：填充后层头发

    ctx.fill()

    // 声明变量 gloss，创建后层头发的暖棕高光带，增加发丝层次

    const gloss = ctx.createLinearGradient(-78, -150, -30, -20)

    // 渐变停靠点：高光起点使用半透明暖棕

    gloss.addColorStop(0, 'rgba(180,130,95,0.30)')

    // 渐变停靠点：高光末端淡出为透明

    gloss.addColorStop(1, 'rgba(180,130,95,0)')

    // 设置填充样式：应用后发高光渐变

    ctx.fillStyle = gloss

    // 调用方法：描绘左侧发丝高光带

    ctx.beginPath()

    // 调用方法：定位高光带上缘起点

    ctx.moveTo(-64, -138)

    // 调用方法：用曲线描绘高光带外缘

    ctx.quadraticCurveTo(-86, -70, -72, 4)

    // 调用方法：用曲线描绘高光带内缘收拢

    ctx.quadraticCurveTo(-58, -60, -46, -126)

    // 调用方法：闭合高光带

    ctx.closePath()

    // 调用方法：填充高光带

    ctx.fill()

    // 设置描边样式：右侧发绺使用低透明暖棕反光

    ctx.strokeStyle = 'rgba(190,140,105,0.22)'

    // 设置线宽：反光保持纤细

    ctx.lineWidth = 3

    // 设置线帽：圆头反光

    ctx.lineCap = 'round'

    // 调用方法：绘制右侧发绺柔和反光

    ctx.beginPath()

    // 调用方法：定位反光起点

    ctx.moveTo(70, -110)

    // 调用方法：描绘反光弧线

    ctx.quadraticCurveTo(86, -50, 76, 30)

    // 调用方法：描边反光

    ctx.stroke()

    // 调用方法：恢复绘图状态，结束后层头发绘制

    ctx.restore()

  // drawBackHair 方法结束

  }

  // 定义方法 drawFace，绘制圆润可爱的精致脸型与耳朵
  private drawFace(ctx: CanvasRenderingContext2D) {

    // 调用方法：保存绘图状态，隔离脸型绘制

    ctx.save()

    // 声明变量 g，创建脸型肤色渐变（额头高光向自然白皙肤色柔和过渡）

    const g = ctx.createLinearGradient(0, -150, 0, 30)

    // 渐变停靠点：额头使用明亮暖白肤色

    g.addColorStop(0, '#fff2e4')

    // 渐变停靠点：脸颊使用自然白皙肤色

    g.addColorStop(0.55, '#ffe3c9')

    // 渐变停靠点：下巴使用柔和深肤色

    g.addColorStop(1, '#ffd2b2')

    // 设置填充样式：应用脸型渐变

    ctx.fillStyle = g

    // 调用方法：开始绘制脸型路径

    ctx.beginPath()

    // 调用方法：以椭圆描绘圆润脸型轮廓（水平半径 73，垂直半径 82，中心上移 58）

    ctx.ellipse(0, -58, 73, 82, 0, 0, Math.PI * 2)

    // 调用方法：填充脸型

    ctx.fill()

    // 声明变量 volume，创建面部柔和体积阴影，让扁平脸型获得圆润立体感

    const volume = ctx.createRadialGradient(-10, -76, 16, 0, -58, 96)

    // 渐变停靠点：额头偏左使用柔和高光

    volume.addColorStop(0, 'rgba(255,255,255,0.30)')

    // 渐变停靠点：中部完全透明

    volume.addColorStop(0.6, 'rgba(255,225,200,0)')

    // 渐变停靠点：脸缘使用暖棕阴影收拢轮廓

    volume.addColorStop(1, 'rgba(196,120,88,0.22)')

    // 设置填充样式：应用面部体积渐变

    ctx.fillStyle = volume

    // 调用方法：沿用脸型路径填充体积光影

    ctx.fill()

    // 设置描边样式：柔和暖色轮廓弱化脸型与发丝、脖颈的硬边界

    ctx.strokeStyle = 'rgba(214,158,124,0.28)'

    // 设置线宽：脸部轮廓保持纤细

    ctx.lineWidth = 1.4

    // 调用方法：描边脸型轮廓

    ctx.stroke()

    // 调用方法：绘制左侧耳朵

    ctx.beginPath()

    // 调用方法：以椭圆描绘左耳

    ctx.ellipse(-72, -52, 10, 17, -0.1, 0, Math.PI * 2)

    // 设置填充样式：耳朵使用稍深肤色

    ctx.fillStyle = '#f7cbaa'

    // 调用方法：填充左耳

    ctx.fill()

    // 调用方法：绘制右侧耳朵

    ctx.beginPath()

    // 调用方法：以椭圆描绘右耳

    ctx.ellipse(72, -52, 10, 17, 0.1, 0, Math.PI * 2)

    // 调用方法：填充右耳

    ctx.fill()

    // 调用方法：恢复绘图状态，结束脸型绘制

    ctx.restore()

  // drawFace 方法结束

  }

  // 定义方法 drawBrows，按情绪绘制不同倾斜度的柔和深棕弯眉
  private drawBrows(ctx: CanvasRenderingContext2D, happy: boolean, sad: boolean) {

    // 调用方法：保存绘图状态，隔离眉毛绘制

    ctx.save()

    // 设置描边样式：使用深棕色作为眉毛颜色

    ctx.strokeStyle = '#5a3b2c'

    // 设置线宽：眉毛描边宽度

    ctx.lineWidth = 4

    // 设置线帽：使用圆头让眉毛更柔和

    ctx.lineCap = 'round'

    // 声明变量 tilt，根据情绪决定眉毛倾斜方向（开心上扬、难过下垂）

    const tilt = happy ? -4 : sad ? 5 : 0

    // 调用方法：开始绘制左眉路径

    ctx.beginPath()

    // 调用方法：用二次曲线描绘左眉，形成自然弧度

    ctx.moveTo(-46, -94 + tilt)

    // 调用方法：控制左眉弧线

    ctx.quadraticCurveTo(-30, -102 + tilt, -16, -93)

    // 调用方法：描边左眉

    ctx.stroke()

    // 调用方法：开始绘制右眉路径

    ctx.beginPath()

    // 调用方法：用二次曲线描绘右眉

    ctx.moveTo(46, -94 + tilt)

    // 调用方法：控制右眉弧线

    ctx.quadraticCurveTo(30, -102 + tilt, 16, -93)

    // 调用方法：描边右眉

    ctx.stroke()

    // 调用方法：恢复绘图状态，结束眉毛绘制

    ctx.restore()

  // drawBrows 方法结束

  }

  // 定义方法 drawEyes，按眨眼与眼神微动绘制双眼
  private drawEyes(ctx: CanvasRenderingContext2D, blink: number, happy: boolean, now: number) {

    // 声明变量 openness，将闭眼程度换算为眼睛张开比例（0 闭合、1 张开）

    const openness = Math.max(0.04, 1 - blink)

    // 声明变量 gazeX，用慢速正弦生成左右眼神微动偏移量

    const gazeX = Math.sin(now / 1000 * 0.6 + this.gazeSeed) * 2.4

    // 声明变量 gazeY，用更慢的余弦生成上下眼神微动偏移量

    const gazeY = Math.cos(now / 1000 * 0.45 + this.gazeSeed) * 1.3

    // 调用方法：绘制左眼

    this.drawEye(ctx, -30, -64, openness, happy, gazeX, gazeY)

    // 调用方法：绘制右眼

    this.drawEye(ctx, 30, -64, openness, happy, gazeX, gazeY)

  // drawEyes 方法结束

  }

  // 定义方法 drawEye，绘制单只明亮大眼（眼白、暖棕虹膜、瞳孔、多重高光与睫毛）
  private drawEye(
    ctx: CanvasRenderingContext2D,
    x: number,
    y: number,
    openness: number,
    happy: boolean,
    gazeX: number,
    gazeY: number,
  ) {

    // 调用方法：保存绘图状态，隔离单眼绘制

    ctx.save()

    // 调用方法：将原点移动到眼睛中心

    ctx.translate(x, y)

    // 条件判断：开心情绪且眼睛张开时绘制弯弯笑眼

    if (happy && openness > 0.6) {

      // 设置描边样式：深棕色眼线

      ctx.strokeStyle = '#3a241a'

      // 设置线宽：笑眼描边宽度

      ctx.lineWidth = 4.5

      // 设置线帽：圆头眼线

      ctx.lineCap = 'round'

      // 调用方法：开始绘制笑眼弧线

      ctx.beginPath()

      // 调用方法：绘制向上弯的笑眼

      ctx.arc(0, 6, 18, Math.PI * 1.12, Math.PI * 1.88)

      // 调用方法：描边笑眼

      ctx.stroke()

    } else {

      // 设置填充样式：眼白使用微带冷调的近白色，更显清透

      ctx.fillStyle = '#fcfdff'

      // 调用方法：开始绘制眼白椭圆

      ctx.beginPath()

      // 调用方法：按张开比例绘制大眼白（纵向半径随闭眼压缩）

      ctx.ellipse(0, 0, 18, 20 * openness, 0, 0, Math.PI * 2)

      // 调用方法：填充眼白

      ctx.fill()

      // 条件判断：仅在眼睛张开超过阈值时绘制虹膜与高光

      if (openness > 0.35) {

        // 声明变量 ig，创建虹膜径向渐变（外深内亮的暖棕）

        const ig = ctx.createRadialGradient(gazeX, gazeY, 1, gazeX, gazeY, 12 * openness)

        // 渐变停靠点：虹膜中心使用明亮暖棕

        ig.addColorStop(0, '#d69a63')

        // 渐变停靠点：虹膜中段使用暖棕色

        ig.addColorStop(0.55, '#965e38')

        // 渐变停靠点：虹膜外缘使用深棕

        ig.addColorStop(1, '#4a2c1c')

        // 设置填充样式：应用虹膜渐变

        ctx.fillStyle = ig

        // 调用方法：开始绘制虹膜

        ctx.beginPath()

        // 调用方法：绘制大虹膜圆，随眼神偏移

        ctx.arc(gazeX, gazeY, 12 * openness, 0, Math.PI * 2)

        // 调用方法：填充虹膜

        ctx.fill()

        // 设置填充样式：瞳孔使用深棕黑

        ctx.fillStyle = '#241309'

        // 调用方法：开始绘制瞳孔

        ctx.beginPath()

        // 调用方法：绘制瞳孔圆

        ctx.arc(gazeX, gazeY, 5 * openness, 0, Math.PI * 2)

        // 调用方法：填充瞳孔

        ctx.fill()

        // 设置填充样式：主高光使用纯白

        ctx.fillStyle = '#ffffff'

        // 调用方法：开始绘制主高光点

        ctx.beginPath()

        // 调用方法：绘制虹膜左上方主高光

        ctx.arc(gazeX - 4, gazeY - 4 * openness, 3.2, 0, Math.PI * 2)

        // 调用方法：填充主高光

        ctx.fill()

        // 调用方法：开始绘制次高光点

        ctx.beginPath()

        // 调用方法：绘制虹膜右下方次高光

        ctx.arc(gazeX + 3.5, gazeY + 3.5 * openness, 1.7, 0, Math.PI * 2)

        // 调用方法：填充次高光

        ctx.fill()

        // 设置描边样式：底部弧形反光增强眼睛通透感

        ctx.strokeStyle = 'rgba(255,255,255,0.55)'

        // 设置线宽：反光保持纤细

        ctx.lineWidth = 1.6

        // 调用方法：开始绘制底部弧形反光

        ctx.beginPath()

        // 调用方法：沿虹膜下缘绘制反光弧

        ctx.arc(gazeX, gazeY + 1, 8 * openness, Math.PI * 0.15, Math.PI * 0.5)

        // 调用方法：描边底部反光

        ctx.stroke()

        // 设置描边样式：上眼线使用深棕色

        ctx.strokeStyle = '#3a241a'

        // 设置线宽：睫毛线宽

        ctx.lineWidth = 3.4

        // 设置线帽：圆头睫毛线

        ctx.lineCap = 'round'

        // 调用方法：开始绘制上睫毛弧线

        ctx.beginPath()

        // 调用方法：绘制覆盖眼白上缘的睫毛弧

        ctx.arc(0, 0, 19, Math.PI * 1.05, Math.PI * 1.95)

        // 调用方法：描边睫毛

        ctx.stroke()

        // 设置描边样式：下睫毛使用低透明深棕

        ctx.strokeStyle = 'rgba(58,36,26,0.5)'

        // 设置线宽：下睫毛更纤细

        ctx.lineWidth = 1.4

        // 调用方法：开始绘制下睫毛弧线

        ctx.beginPath()

        // 调用方法：绘制眼白下缘的下睫毛弧

        ctx.arc(0, 0, 19, Math.PI * 0.1, Math.PI * 0.75)

        // 调用方法：描边下睫毛

        ctx.stroke()

      }

    }

    // 调用方法：恢复绘图状态，结束单眼绘制

    ctx.restore()

  // drawEye 方法结束

  }

  // 定义方法 drawNose，绘制小巧鼻部柔和阴影与高光
  private drawNose(ctx: CanvasRenderingContext2D) {

    // 调用方法：保存绘图状态，隔离鼻部绘制

    ctx.save()

    // 设置填充样式：鼻影使用低透明度暖色

    ctx.fillStyle = 'rgba(214,150,120,0.42)'

    // 调用方法：开始绘制鼻影路径

    ctx.beginPath()

    // 调用方法：以小椭圆描绘鼻尖柔和阴影

    ctx.ellipse(0, -26, 4.5, 3, 0, 0, Math.PI * 2)

    // 调用方法：填充鼻影

    ctx.fill()

    // 设置填充样式：鼻尖高光使用半透明白

    ctx.fillStyle = 'rgba(255,255,255,0.5)'

    // 调用方法：开始绘制鼻尖高光

    ctx.beginPath()

    // 调用方法：以小椭圆描绘鼻尖高光

    ctx.ellipse(-1.2, -27.5, 1.6, 1.1, 0, 0, Math.PI * 2)

    // 调用方法：填充鼻尖高光

    ctx.fill()

    // 调用方法：恢复绘图状态，结束鼻部绘制

    ctx.restore()

  // drawNose 方法结束

  }

  // 定义方法 drawMouth，按情绪与说话口型绘制亲切微笑的嘴部
  private drawMouth(
    ctx: CanvasRenderingContext2D,
    happy: boolean,
    sad: boolean,
    mouthOpen: number,
  ) {

    // 调用方法：保存绘图状态，隔离嘴部绘制

    ctx.save()

    // 调用方法：将原点移动到嘴部中心

    ctx.translate(0, -6)

    // 条件判断：说话状态（存在张嘴幅度）时绘制张开的立体嘴型

    if (mouthOpen > 0.06) {

      // 声明变量 oh，将口型程度换算为嘴部张开高度

      const oh = 3 + mouthOpen * 16

      // 声明变量 ow，根据张口程度略微收窄嘴宽，形成自然椭圆

      const ow = 17 - mouthOpen * 3

      // 设置填充样式：口腔内部使用深玫红

      ctx.fillStyle = '#a8435a'

      // 调用方法：开始绘制口腔路径

      ctx.beginPath()

      // 调用方法：以椭圆描绘张开的口腔

      ctx.ellipse(0, 6, ow, oh, 0, 0, Math.PI * 2)

      // 调用方法：填充口腔

      ctx.fill()

      // 设置填充样式：上排牙齿使用近白色

      ctx.fillStyle = '#fff7f3'

      // 调用方法：开始绘制上排牙齿

      ctx.beginPath()

      // 调用方法：绘制贴合口腔上缘的牙齿弧条

      ctx.ellipse(0, 6 - oh * 0.55, ow * 0.78, oh * 0.32, 0, 0, Math.PI * 2)

      // 调用方法：填充牙齿

      ctx.fill()

      // 条件判断：张口较大时绘制舌尖细节

      if (mouthOpen > 0.45) {

        // 设置填充样式：舌尖使用柔和粉色

        ctx.fillStyle = '#e07f8a'

        // 调用方法：开始绘制舌尖

        ctx.beginPath()

        // 调用方法：绘制口腔下部的舌尖椭圆

        ctx.ellipse(0, 6 + oh * 0.5, ow * 0.5, oh * 0.28, 0, 0, Math.PI * 2)

        // 调用方法：填充舌尖

        ctx.fill()

      }

      // 设置描边样式：唇线使用柔和玫红

      ctx.strokeStyle = '#d16b7a'

      // 设置线宽：唇线描边宽度

      ctx.lineWidth = 2.5

      // 调用方法：开始绘制嘴唇外轮廓

      ctx.beginPath()

      // 调用方法：描绘嘴唇外缘椭圆

      ctx.ellipse(0, 6, ow + 1, oh + 1, 0, 0, Math.PI * 2)

      // 调用方法：描边嘴唇

      ctx.stroke()

    } else {

      // 设置描边样式：唇线使用柔和玫红

      ctx.strokeStyle = '#d16b7a'

      // 设置线宽：唇线描边宽度

      ctx.lineWidth = 3

      // 设置线帽：圆头唇线

      ctx.lineCap = 'round'

      // 调用方法：开始绘制闭嘴曲线

      ctx.beginPath()

      // 条件判断：开心情绪绘制上扬微笑

      if (happy) {

        // 调用方法：用二次曲线绘制微笑唇形

        ctx.moveTo(-15, 4)

        // 调用方法：控制微笑弧线中段

        ctx.quadraticCurveTo(0, 17, 15, 4)

        // 调用方法：用二次曲线闭合上唇

        ctx.quadraticCurveTo(0, 10, -15, 4)

      } else if (sad) {

        // 调用方法：用二次曲线绘制下垂唇形

        ctx.moveTo(-14, 8)

        // 调用方法：控制下垂弧线中段

        ctx.quadraticCurveTo(0, 0, 14, 8)

      } else {

        // 调用方法：绘制自然闭合的平静唇形

        ctx.moveTo(-13, 6)

        // 调用方法：用二次曲线连接右端

        ctx.quadraticCurveTo(0, 12, 13, 6)

      }

      // 调用方法：描边唇形

      ctx.stroke()

    }

    // 调用方法：恢复绘图状态，结束嘴部绘制

    ctx.restore()

  // drawMouth 方法结束

  }

  // 定义方法 drawBlush，在两颊绘制柔和粉嫩腮红与高光点
  private drawBlush(ctx: CanvasRenderingContext2D, happy: boolean) {

    // 调用方法：保存绘图状态，隔离腮红绘制

    ctx.save()

    // 声明变量 alpha，开心时腮红更明显、平时更淡雅

    const alpha = happy ? 0.34 : 0.22

    // 设置填充样式：使用柔和粉色半透明

    ctx.fillStyle = `rgba(242,140,150,${alpha.toFixed(3)})`

    // 调用方法：绘制左颊腮红

    ctx.beginPath()

    // 调用方法：以椭圆描绘左侧腮红

    ctx.ellipse(-50, -40, 15, 8.5, 0, 0, Math.PI * 2)

    // 调用方法：填充左侧腮红

    ctx.fill()

    // 调用方法：绘制右颊腮红

    ctx.beginPath()

    // 调用方法：以椭圆描绘右侧腮红

    ctx.ellipse(50, -40, 15, 8.5, 0, 0, Math.PI * 2)

    // 调用方法：填充右侧腮红

    ctx.fill()

    // 设置填充样式：腮红高光使用半透明白

    ctx.fillStyle = 'rgba(255,255,255,0.35)'

    // 调用方法：绘制左颊高光点

    ctx.beginPath()

    // 调用方法：以小圆描绘左侧高光

    ctx.arc(-52, -42, 2.5, 0, Math.PI * 2)

    // 调用方法：填充左侧高光

    ctx.fill()

    // 调用方法：绘制右颊高光点

    ctx.beginPath()

    // 调用方法：以小圆描绘右侧高光

    ctx.arc(48, -42, 2.5, 0, Math.PI * 2)

    // 调用方法：填充右侧高光

    ctx.fill()

    // 调用方法：恢复绘图状态，结束腮红绘制

    ctx.restore()

  // drawBlush 方法结束

  }

  // 定义方法 drawFrontHair，绘制空气刘海与柔暖发丝高光
  private drawFrontHair(ctx: CanvasRenderingContext2D) {

    // 调用方法：保存绘图状态，隔离前层头发绘制

    ctx.save()

    // 声明变量 g，创建空气刘海渐变（顶部暖棕高光向发尾深棕过渡）

    const g = ctx.createLinearGradient(0, -172, 0, -70)

    // 渐变停靠点：顶部使用暖棕高光

    g.addColorStop(0, '#5a3d2e')

    // 渐变停靠点：中部使用深棕主色

    g.addColorStop(0.5, '#3d2820')

    // 渐变停靠点：发尾使用深棕黑

    g.addColorStop(1, '#2a1b16')

    // 设置填充样式：应用刘海渐变

    ctx.fillStyle = g

    // 调用方法：开始绘制空气刘海整体路径（轻薄、眉上、几缕分开）

    ctx.beginPath()

    // 调用方法：定位左侧刘海起点

    ctx.moveTo(-78, -66)

    // 调用方法：用贝塞尔曲线描绘左侧发际至头顶

    ctx.bezierCurveTo(-88, -150, -52, -170, 0, -170)

    // 调用方法：用贝塞尔曲线描绘右侧发际

    ctx.bezierCurveTo(52, -170, 88, -150, 78, -66)

    // 调用方法：绘制右侧第一缕空气刘海底缘

    ctx.quadraticCurveTo(64, -92, 50, -84)

    // 调用方法：绘制右侧第二缕空气刘海底缘

    ctx.quadraticCurveTo(40, -100, 26, -88)

    // 调用方法：绘制中部空气刘海底缘

    ctx.quadraticCurveTo(14, -102, 2, -88)

    // 调用方法：绘制中左空气刘海底缘

    ctx.quadraticCurveTo(-10, -100, -24, -86)

    // 调用方法：绘制左侧第二缕空气刘海底缘

    ctx.quadraticCurveTo(-40, -98, -54, -84)

    // 调用方法：绘制左侧第一缕空气刘海底缘回到起点

    ctx.quadraticCurveTo(-68, -94, -78, -66)

    // 调用方法：闭合刘海路径

    ctx.closePath()

    // 调用方法：填充刘海

    ctx.fill()

    // 声明变量 sheen，创建刘海顶部的宽幅柔和高光，赋予丝绸般光泽

    const sheen = ctx.createLinearGradient(0, -170, 0, -116)

    // 渐变停靠点：高光顶部使用半透明暖棕

    sheen.addColorStop(0, 'rgba(210,160,120,0.28)')

    // 渐变停靠点：高光底部淡出为透明

    sheen.addColorStop(1, 'rgba(210,160,120,0)')

    // 设置填充样式：应用刘海高光渐变

    ctx.fillStyle = sheen

    // 调用方法：描绘刘海高光带

    ctx.beginPath()

    // 调用方法：定位高光带左侧起点

    ctx.moveTo(-62, -120)

    // 调用方法：用曲线描绘高光带上缘

    ctx.quadraticCurveTo(-38, -160, 0, -164)

    // 调用方法：用曲线描绘高光带上缘右侧

    ctx.quadraticCurveTo(38, -160, 62, -120)

    // 调用方法：用曲线描绘高光带下缘右侧

    ctx.quadraticCurveTo(30, -142, 0, -140)

    // 调用方法：用曲线描绘高光带下缘左侧

    ctx.quadraticCurveTo(-30, -142, -62, -120)

    // 调用方法：闭合高光带

    ctx.closePath()

    // 调用方法：填充高光带

    ctx.fill()

    // 设置描边样式：发丝高光使用半透明暖棕

    ctx.strokeStyle = 'rgba(205,155,115,0.4)'

    // 设置线宽：发丝高光保持纤细

    ctx.lineWidth = 2.4

    // 设置线帽：圆头发丝

    ctx.lineCap = 'round'

    // 调用方法：绘制左侧发丝高光

    ctx.beginPath()

    // 调用方法：描绘左侧弧形高光

    ctx.moveTo(-52, -136)

    // 调用方法：控制左侧高光弧线

    ctx.quadraticCurveTo(-32, -156, -6, -152)

    // 调用方法：描边左侧高光

    ctx.stroke()

    // 调用方法：绘制右侧发丝高光

    ctx.beginPath()

    // 调用方法：描绘右侧弧形高光

    ctx.moveTo(22, -154)

    // 调用方法：控制右侧高光弧线

    ctx.quadraticCurveTo(46, -150, 58, -128)

    // 调用方法：描边右侧高光

    ctx.stroke()

    // 调用方法：绘制中部细高光，进一步丰富发丝层次

    ctx.beginPath()

    // 调用方法：描绘中部短弧高光

    ctx.moveTo(-6, -158)

    // 调用方法：控制中部高光弧线

    ctx.quadraticCurveTo(6, -156, 18, -158)

    // 调用方法：描边中部高光

    ctx.stroke()

    // 调用方法：恢复绘图状态，结束前层头发绘制

    ctx.restore()

  // drawFrontHair 方法结束

  }


// CartoonAvatarEngine 类结束

}
