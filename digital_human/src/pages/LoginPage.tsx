/** @ticket 07-数字人Coze：数字人登录页——按场景角色进入对应数字人页面（寄件引导 / 取件指引 / 智能客服），角色为数字人交互身份而非系统权限 */
import { DhRole } from './types'

// 角色卡片元数据：场景身份/描述/进入后的数字人形象
// 声明常量 ROLE_CARDS，用于保存运行过程中的数据

const ROLE_CARDS: { role: DhRole; icon: string; title: string; desc: string }[] = [
  // 寄件引导员：面向首次寄件用户，引导选点→选格→支付→投件全流程
  { role: 'sender', icon: '📦', title: '寄件引导员', desc: '引导寄件全流程：选择寄存点、选格口、支付与投件' },
  // 取件指引员：面向取件用户，讲解取件码、临时开柜、超时费用与异常开门
  { role: 'receiver', icon: '🔑', title: '取件指引员', desc: '取件码使用、临时开柜、超时补交与开门异常处理' },
  // 智能客服专员：常见问题解答（计费规则/优惠券/柜门故障/联系网点）
  { role: 'service', icon: '💬', title: '智能客服专员', desc: '计费规则、优惠券使用、柜门故障与联系网点' },
// 当前代码块结束

]

// 导出角色元数据（App 顶栏徽标与形象映射使用）
// 导出 const DH_ROLES，供其他模块引用

export const DH_ROLES = ROLE_CARDS

// 定义组件 Props：onLogin 回调携带所选场景角色
interface Props {
  // onLogin：角色选择完成回调
  onLogin: (role: DhRole) => void
}

// 导出默认组件 LoginPage：数字人场景角色登录页
export default function LoginPage({ onLogin }: Props) {
  // 渲染全屏登录容器
  return (
    // 登录页根容器
    <div className="dh-login">
      {/* 品牌标题区 */}
      <div className="dh-login-head">
        {/* 大标题 */}
        <h1>智能快寄柜 · 数字人</h1>
        {/* 副标题 */}
        <p>请选择数字人服务角色进入对应的交互页面</p>
      {/* 闭合标题区 */}
      </div>
      {/* 角色卡片列表 */}
      <div className="dh-login-cards">
        {/* 循环渲染三个场景角色卡片，点击即登录进入对应页面 */}
        {ROLE_CARDS.map((c) => (
          // 角色卡片：点击触发 onLogin
          <div key={c.role} className="dh-login-card" onClick={() => onLogin(c.role)}>
            {/* 图标 */}
            <span className="c-icon">{c.icon}</span>
            {/* 标题 */}
            <span className="c-title">{c.title}</span>
            {/* 描述 */}
            <span className="c-desc">{c.desc}</span>
            {/* 进入按钮文案 */}
            <span className="c-go">进入 ›</span>
          {/* 闭合卡片 */}
          </div>
        ))}
      {/* 闭合卡片列表 */}
      </div>
      {/* 底部说明 */}
      <p className="dh-login-foot">提示：不同角色对应不同的数字人话术库与服务场景，可在页面内随时退出切换</p>
    {/* 闭合根容器 */}
    </div>
  )
// 当前代码块结束
}
