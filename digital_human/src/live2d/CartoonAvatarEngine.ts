/** @ticket 07-数字人Coze：零依赖/零授权 Canvas 卡通数字人自绘引擎（替换无 SDK 的空白 Live2D 渲染区） */

// 从统一类型模块导入 Live2DAction 类型，用于接收工单既有的表情/动作 payload

import type { Live2DAction } from '@/types'

// 定义并导出 CartoonAvatarEngine 类，封装纯 Canvas 2D 数字人的绘制与动画循环

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

    const breath = Math.sin(t * 1.8)

    // 声明变量 sway，用慢速正弦波生成身体轻微摆动角度

    const sway = Math.sin(t * 0.9) * 0.03

    // 条件判断：到达计划眨眼时间时启动一次新的眨眼

    if (now >= this.nextBlink) {

      // 赋值更新：记录本次眨眼起始时间

      this.lastBlinkStart = now

      // 赋值更新：安排下一次眨眼（2.2 到 5.2 秒之间随机）

      this.nextBlink = now + 2200 + Math.random() * 3000

    }

    // 声明变量 blink，按眨眼进度生成 0 到 1 的闭眼程度（正弦包络）

    const blink = Math.max(

      // 计算自动眨眼包络：眨眼窗口内呈半周期正弦

      this.autoBlink(now),

      // 计算动作眨眼：blink 动作时按进度生成闭眼包络

      active && kind === 'blink' ? Math.sin(progress * Math.PI) : 0,

    )

    // 声明变量 headBob，nod 动作时让头部上下点头

    const headBob = active && kind === 'nod' ? Math.sin(progress * Math.PI * 2) * 12 : 0

    // 声明变量 headTurn，shake 动作时让头部左右摇头

    const headTurn = active && kind === 'shake' ? Math.sin(progress * Math.PI * 4) * 0.16 : 0

    // 声明变量 transformScale，transform 动作时整体短暂放大

    const transformScale = active && kind === 'transform' ? 1 + Math.sin(progress * Math.PI) * 0.08 : 1

    // 声明变量 mouthOpen，说话回应期间让嘴部按高频正弦开合

    const mouthOpen = active ? (0.5 + 0.5 * Math.sin(t * 16)) * 9 * (kind === 'smile' ? 0.35 : 1) : 0

    // 声明变量 cx，取得画布水平中心作为角色锚点

    const cx = this.w / 2

    // 声明变量 cy，取得画布垂直中心作为角色锚点

    const cy = this.h / 2

    // 声明变量 scale，按画布较小边计算整体缩放基准（基准高度 420）

    const scale = Math.min(this.w, this.h) / 420 * transformScale

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

    // 调用方法：绘制角色身体部分

    this.drawBody(ctx, happy, sad)

    // 调用方法：绘制角色头部（含点头、摇头偏移）

    this.drawHead(ctx, headBob, headTurn, blink, happy, sad, mouthOpen, breath)

    // 调用方法：恢复绘图状态，结束角色主体绘制

    ctx.restore()

  // draw 方法结束

  }

  // 定义方法 autoBlink，按眨眼窗口返回平滑的闭眼包络值

  private autoBlink(now: number): number {

    // 声明变量 since，计算距离最近一次眨眼起始的时间差

    const since = now - this.lastBlinkStart

    // 条件判断：仅在 160 毫秒的眨眼窗口内返回正闭眼值

    if (since >= 0 && since < 160) {

      // 返回值：半周期正弦包络，闭眼再睁眼

      return Math.sin((since / 160) * Math.PI)

    }

    // 返回值：窗口外保持睁眼

    return 0

  // autoBlink 方法结束

  }

  // 定义方法 drawBody，绘制肩部、上衣与脖颈

  private drawBody(ctx: CanvasRenderingContext2D, happy: boolean, sad: boolean) {

    // 调用方法：保存绘图状态，隔离身体绘制变换

    ctx.save()

    // 声明变量 coat，根据情绪选择上衣主色（开心偏亮蓝，难过偏灰蓝）

    const coat = happy ? '#3b7bf5' : sad ? '#4a5670' : '#2f6df0'

    // 设置填充样式：应用上衣主色

    ctx.fillStyle = coat

    // 调用方法：开始绘制肩部轮廓路径

    ctx.beginPath()

    // 调用方法：定位到肩部左上起点

    ctx.moveTo(-150, 200)

    // 调用方法：用二次贝塞尔曲线描绘左肩到颈部

    ctx.quadraticCurveTo(-140, 70, -55, 55)

    // 调用方法：用水平线连接左右领口

    ctx.lineTo(55, 55)

    // 调用方法：用二次贝塞尔曲线描绘右肩

    ctx.quadraticCurveTo(140, 70, 150, 200)

    // 调用方法：闭合肩部路径回到起点

    ctx.closePath()

    // 调用方法：填充肩部轮廓

    ctx.fill()

    // 设置填充样式：使用浅肤色绘制脖颈

    ctx.fillStyle = '#f0c6a0'

    // 调用方法：开始绘制脖颈矩形

    ctx.fillRect(-26, 20, 52, 60)

    // 调用方法：恢复绘图状态，结束身体绘制

    ctx.restore()

  // drawBody 方法结束

  }

  // 定义方法 drawHead，绘制脸型、头发、眉眼、鼻子、嘴巴与腮红

  private drawHead(
    ctx: CanvasRenderingContext2D,
    headBob: number,
    headTurn: number,
    blink: number,
    happy: boolean,
    sad: boolean,
    mouthOpen: number,
    breath: number,
  ) {

    // 调用方法：保存绘图状态，隔离头部绘制变换

    ctx.save()

    // 调用方法：将原点移动到头部中心，并叠加点头偏移

    ctx.translate(0, -50 + headBob)

    // 调用方法：按摇头角度旋转头部

    ctx.rotate(headTurn)

    // 调用方法：按呼吸系数让头部做极轻微的纵向起伏

    ctx.scale(1, 1 + breath * 0.01)

    // 设置填充样式：使用肤色绘制脸型

    ctx.fillStyle = '#ffd9b3'

    // 调用方法：开始绘制脸型椭圆路径

    ctx.beginPath()

    // 调用方法：绘制脸部椭圆轮廓（水平半径 78，垂直半径 88）

    ctx.ellipse(0, 0, 78, 88, 0, 0, Math.PI * 2)

    // 调用方法：填充脸型

    ctx.fill()

    // 调用方法：绘制头发

    this.drawHair(ctx)

    // 调用方法：绘制眉毛

    this.drawBrows(ctx, happy, sad)

    // 调用方法：绘制眼睛（含眨眼）

    this.drawEyes(ctx, blink, happy)

    // 设置填充样式：用稍深肤色绘制鼻尖

    ctx.fillStyle = '#e8b48c'

    // 调用方法：开始绘制鼻尖圆点

    ctx.beginPath()

    // 调用方法：绘制鼻尖小圆

    ctx.arc(0, 6, 4, 0, Math.PI * 2)

    // 调用方法：填充鼻尖

    ctx.fill()

    // 调用方法：绘制嘴巴

    this.drawMouth(ctx, happy, sad, mouthOpen)

    // 条件判断：开心情绪时绘制腮红点缀

    if (happy) {

      // 调用方法：绘制两侧腮红

      this.drawBlush(ctx)

    }

    // 调用方法：恢复绘图状态，结束头部绘制

    ctx.restore()

  // drawHead 方法结束

  }

  // 定义方法 drawHair，绘制覆盖头顶的发型

  private drawHair(ctx: CanvasRenderingContext2D) {

    // 调用方法：保存绘图状态，隔离发型绘制

    ctx.save()

    // 设置填充样式：使用深棕色作为头发颜色

    ctx.fillStyle = '#3a2b26'

    // 调用方法：开始绘制头发路径

    ctx.beginPath()

    // 调用方法：定位到左鬓角起点

    ctx.moveTo(-80, 10)

    // 调用方法：用二次贝塞尔曲线描绘左侧发际到头顶

    ctx.quadraticCurveTo(-86, -104, 0, -104)

    // 调用方法：用二次贝塞尔曲线描绘右侧发际

    ctx.quadraticCurveTo(86, -104, 80, 10)

    // 调用方法：绘制刘海底缘中段，形成齐刘海造型

    ctx.quadraticCurveTo(40, -34, 0, -40)

    // 调用方法：绘制左侧刘海底缘

    ctx.quadraticCurveTo(-40, -34, -80, 10)

    // 调用方法：闭合头发路径

    ctx.closePath()

    // 调用方法：填充头发

    ctx.fill()

    // 调用方法：恢复绘图状态，结束发型绘制

    ctx.restore()

  // drawHair 方法结束

  }

  // 定义方法 drawBrows，按情绪绘制不同倾斜度的眉毛

  private drawBrows(ctx: CanvasRenderingContext2D, happy: boolean, sad: boolean) {

    // 调用方法：保存绘图状态，隔离眉毛绘制

    ctx.save()

    // 设置描边样式：使用深棕色作为眉毛颜色

    ctx.strokeStyle = '#4a3428'

    // 设置线宽：眉毛描边宽度

    ctx.lineWidth = 5

    // 设置线帽：使用圆头让眉毛更柔和

    ctx.lineCap = 'round'

    // 声明变量 tilt，根据情绪决定眉毛倾斜方向（开心上扬、难过下垂）

    const tilt = happy ? -6 : sad ? 7 : 0

    // 调用方法：开始绘制左眉路径

    ctx.beginPath()

    // 调用方法：定位左眉内端

    ctx.moveTo(-48, -34 + tilt)

    // 调用方法：连线到左眉外端

    ctx.lineTo(-18, -40)

    // 调用方法：描边左眉

    ctx.stroke()

    // 调用方法：开始绘制右眉路径

    ctx.beginPath()

    // 调用方法：定位右眉内端

    ctx.moveTo(48, -34 + tilt)

    // 调用方法：连线到右眉外端

    ctx.lineTo(18, -40)

    // 调用方法：描边右眉

    ctx.stroke()

    // 调用方法：恢复绘图状态，结束眉毛绘制

    ctx.restore()

  // drawBrows 方法结束

  }

  // 定义方法 drawEyes，按眨眼程度绘制睁眼或闭眼

  private drawEyes(ctx: CanvasRenderingContext2D, blink: number, happy: boolean) {

    // 声明变量 openness，将闭眼程度换算为眼睛张开比例（0 闭合、1 张开）

    const openness = Math.max(0.04, 1 - blink)

    // 调用方法：绘制左眼

    this.drawEye(ctx, -30, -14, openness, happy)

    // 调用方法：绘制右眼

    this.drawEye(ctx, 30, -14, openness, happy)

  // drawEyes 方法结束

  }

  // 定义方法 drawEye，绘制单只眼睛（含眼白、瞳孔与高光）

  private drawEye(
    ctx: CanvasRenderingContext2D,
    x: number,
    y: number,
    openness: number,
    happy: boolean,
  ) {

    // 调用方法：保存绘图状态，隔离单眼绘制

    ctx.save()

    // 调用方法：将原点移动到眼睛中心

    ctx.translate(x, y)

    // 条件判断：开心情绪时把眼睛画成弯弯的笑眼

    if (happy && openness > 0.6) {

      // 设置描边样式：深棕色眼线

      ctx.strokeStyle = '#2a211d'

      // 设置线宽：笑眼描边宽度

      ctx.lineWidth = 6

      // 设置线帽：圆头眼线

      ctx.lineCap = 'round'

      // 调用方法：开始绘制笑眼弧线

      ctx.beginPath()

      // 调用方法：绘制向上弯的笑眼

      ctx.arc(0, 6, 18, Math.PI * 1.15, Math.PI * 1.85)

      // 调用方法：描边笑眼

      ctx.stroke()

    } else {

      // 设置填充样式：眼白使用近白色

      ctx.fillStyle = '#ffffff'

      // 调用方法：开始绘制眼白椭圆

      ctx.beginPath()

      // 调用方法：按张开比例绘制眼白（纵向半径随闭眼压缩）

      ctx.ellipse(0, 0, 18, 22 * openness, 0, 0, Math.PI * 2)

      // 调用方法：填充眼白

      ctx.fill()

      // 条件判断：仅在眼睛张开超过阈值时绘制瞳孔与高光

      if (openness > 0.35) {

        // 设置填充样式：瞳孔使用深棕色

        ctx.fillStyle = '#2a211d'

        // 调用方法：开始绘制瞳孔

        ctx.beginPath()

        // 调用方法：绘制瞳孔圆

        ctx.arc(0, 0, 9 * openness, 0, Math.PI * 2)

        // 调用方法：填充瞳孔

        ctx.fill()

        // 设置填充样式：高光使用白色

        ctx.fillStyle = '#ffffff'

        // 调用方法：开始绘制高光点

        ctx.beginPath()

        // 调用方法：绘制瞳孔左上方的高光圆点

        ctx.arc(-3, -3 * openness, 3, 0, Math.PI * 2)

        // 调用方法：填充高光点

        ctx.fill()

      }

    }

    // 调用方法：恢复绘图状态，结束单眼绘制

    ctx.restore()

  // drawEye 方法结束

  }

  // 定义方法 drawMouth，按情绪与说话状态绘制嘴部

  private drawMouth(
    ctx: CanvasRenderingContext2D,
    happy: boolean,
    sad: boolean,
    mouthOpen: number,
  ) {

    // 调用方法：保存绘图状态，隔离嘴部绘制

    ctx.save()

    // 条件判断：说话状态（存在张嘴幅度）时绘制开合的椭圆嘴

    if (mouthOpen > 0.5) {

      // 设置填充样式：口腔使用深红色

      ctx.fillStyle = '#b23a48'

      // 调用方法：开始绘制张开的嘴部椭圆

      ctx.beginPath()

      // 调用方法：绘制嘴部椭圆，纵向半径由张嘴幅度决定

      ctx.ellipse(0, 46, 16, 4 + mouthOpen, 0, 0, Math.PI * 2)

      // 调用方法：填充嘴部

      ctx.fill()

    } else {

      // 设置描边样式：嘴线使用深红色

      ctx.strokeStyle = '#b23a48'

      // 设置线宽：嘴线描边宽度

      ctx.lineWidth = 4

      // 设置线帽：圆头嘴线

      ctx.lineCap = 'round'

      // 调用方法：开始绘制嘴部曲线

      ctx.beginPath()

      // 条件判断：开心情绪绘制上扬的微笑

      if (happy) {

        // 调用方法：绘制上扬弧线（微笑）

        ctx.arc(0, 40, 20, 0.15 * Math.PI, 0.85 * Math.PI)

      } else if (sad) {

        // 调用方法：绘制下垂弧线（难过）

        ctx.arc(0, 66, 20, 1.15 * Math.PI, 1.85 * Math.PI)

      } else {

        // 调用方法：绘制平直短线（平静）

        ctx.moveTo(-14, 50)

        // 调用方法：连线到右端形成平直嘴

        ctx.lineTo(14, 50)

      }

      // 调用方法：描边嘴部曲线

      ctx.stroke()

    }

    // 调用方法：恢复绘图状态，结束嘴部绘制

    ctx.restore()

  // drawMouth 方法结束

  }

  // 定义方法 drawBlush，在两颊绘制粉色腮红

  private drawBlush(ctx: CanvasRenderingContext2D) {

    // 调用方法：保存绘图状态，隔离腮红绘制

    ctx.save()

    // 设置填充样式：使用半透明粉色

    ctx.fillStyle = 'rgba(238, 120, 130, 0.35)'

    // 调用方法：绘制左颊腮红

    ctx.beginPath()

    // 调用方法：绘制左侧腮红椭圆

    ctx.ellipse(-48, 26, 16, 9, 0, 0, Math.PI * 2)

    // 调用方法：填充左侧腮红

    ctx.fill()

    // 调用方法：绘制右颊腮红

    ctx.beginPath()

    // 调用方法：绘制右侧腮红椭圆

    ctx.ellipse(48, 26, 16, 9, 0, 0, Math.PI * 2)

    // 调用方法：填充右侧腮红

    ctx.fill()

    // 调用方法：恢复绘图状态，结束腮红绘制

    ctx.restore()

  // drawBlush 方法结束

  }

// CartoonAvatarEngine 类结束

}