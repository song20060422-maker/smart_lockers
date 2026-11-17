// @ticket 00-project：统一错误码表，按业务域分段
package pkg

// 错误码分段：1xxxx 通用；11xxx 首页；12xxx 支付；13xxx 寄存；14xxx 订单；15xxx 个人中心；16xxx 通知；18xxx 运维
const (
	CodeOK      = 0
	ErrBad      = 10001 // 参数错误
	ErrAuth     = 10002 // 未授权/登录失效
	ErrForb     = 10003 // 禁止访问
	ErrNF       = 10004 // 资源不存在
	ErrInternal = 10005 // 系统内部错误
	ErrLimit    = 10006 // 限流/频控
	ErrDup      = 10007 // 重复请求（幂等命中）
	ErrLock     = 10008 // 并发冲突/锁失败

	// 首页（03）
	ErrLocFail = 11001 // 定位失败
	ErrCityNF  = 11002 // 城市不存在

	// 支付（02）
	ErrBalanceLow  = 12001 // 余额不足
	ErrRechargeNF  = 12002 // 充值单不存在
	ErrPayCallback = 12003 // 支付回调校验失败

	// 寄存（04）
	ErrLockerBusy  = 13001 // 柜门不可用
	ErrBillingRule = 13002 // 收费规则未匹配
	ErrDepositNF   = 13003 // 寄存记录不存在

	// 订单（05）
	ErrOrderNF     = 14001 // 订单不存在
	ErrOrderState  = 14002 // 订单状态不允许该操作
	ErrOpenAuthExp = 14003 // 临时开柜授权失效
	ErrMakeup      = 14004 // 补缴失败

	// 个人中心（06）
	ErrAvatarType = 15001 // 头像格式错误
	ErrAvatarSize = 15002 // 头像过大
	ErrCouponExp  = 15003 // 卡券已过期/不可用

	// 通知（01）
	ErrNotifyTpl = 16001 // 通知模板不存在
)

type ErrInfo struct {
	Code    int
	Message string
}

var errMessages = map[int]string{
	CodeOK: "ok",
	ErrBad: "参数错误", ErrAuth: "登录已失效", ErrForb: "禁止访问", ErrNF: "资源不存在",
	ErrInternal: "系统内部错误", ErrLimit: "操作过于频繁", ErrDup: "重复请求", ErrLock: "操作冲突，请重试",
	ErrLocFail: "定位失败", ErrCityNF: "城市不存在",
	ErrBalanceLow: "余额不足", ErrRechargeNF: "充值单不存在", ErrPayCallback: "支付回调校验失败",
	ErrLockerBusy: "柜门不可用", ErrBillingRule: "收费规则未匹配", ErrDepositNF: "寄存记录不存在",
	ErrOrderNF: "订单不存在", ErrOrderState: "订单状态不允许该操作", ErrOpenAuthExp: "临时开柜授权失效", ErrMakeup: "补缴失败",
	ErrAvatarType: "头像格式错误", ErrAvatarSize: "头像过大", ErrCouponExp: "卡券已过期或不可用",
	ErrNotifyTpl: "通知模板不存在",
}

func GetErr(code int) ErrInfo {
	msg, ok := errMessages[code]
	if !ok {
		msg = "系统内部错误"
	}
	return ErrInfo{Code: code, Message: msg}
}
