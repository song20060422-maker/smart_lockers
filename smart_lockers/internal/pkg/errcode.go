// @ticket 00-project：统一错误码表，按业务域分段
package pkg

// 错误码分段：1xxxx 通用；11xxx 首页；12xxx 支付；13xxx 寄存；14xxx 订单；15xxx 个人中心；16xxx 通知；18xxx 运维
// 开始定义一组常量的 const 声明块

const (
	// 为 CodeOK 赋值：0

	CodeOK      = 0
	// 为 ErrBad 赋值：10001 // 参数错误

	ErrBad      = 10001 // 参数错误
	// 为 ErrAuth 赋值：10002 // 未授权/登录失效

	ErrAuth     = 10002 // 未授权/登录失效
	// 为 ErrForb 赋值：10003 // 禁止访问

	ErrForb     = 10003 // 禁止访问
	// 为 ErrNF 赋值：10004 // 资源不存在

	ErrNF       = 10004 // 资源不存在
	// 为 ErrInternal 赋值：10005 // 系统内部错误

	ErrInternal = 10005 // 系统内部错误
	// 为 ErrLimit 赋值：10006 // 限流/频控

	ErrLimit    = 10006 // 限流/频控
	// 为 ErrDup 赋值：10007 // 重复请求（幂等命中）

	ErrDup      = 10007 // 重复请求（幂等命中）
	// 为 ErrLock 赋值：10008 // 并发冲突/锁失败

	ErrLock     = 10008 // 并发冲突/锁失败

	// 首页（03）
	// 为 ErrLocFail 赋值：11001 // 定位失败

	ErrLocFail = 11001 // 定位失败
	// 为 ErrCityNF 赋值：11002 // 城市不存在

	ErrCityNF  = 11002 // 城市不存在

	// 支付（02）
	// 为 ErrBalanceLow 赋值：12001 // 余额不足

	ErrBalanceLow  = 12001 // 余额不足
	// 为 ErrRechargeNF 赋值：12002 // 充值单不存在

	ErrRechargeNF  = 12002 // 充值单不存在
	// 为 ErrPayCallback 赋值：12003 // 支付回调校验失败

	ErrPayCallback = 12003 // 支付回调校验失败

	// 寄存（04）
	// 为 ErrLockerBusy 赋值：13001 // 柜门不可用

	ErrLockerBusy  = 13001 // 柜门不可用
	// 为 ErrBillingRule 赋值：13002 // 收费规则未匹配

	ErrBillingRule = 13002 // 收费规则未匹配
	// 为 ErrDepositNF 赋值：13003 // 寄存记录不存在

	ErrDepositNF   = 13003 // 寄存记录不存在

	// 订单（05）
	// 为 ErrOrderNF 赋值：14001 // 订单不存在

	ErrOrderNF     = 14001 // 订单不存在
	// 为 ErrOrderState 赋值：14002 // 订单状态不允许该操作

	ErrOrderState  = 14002 // 订单状态不允许该操作
	// 为 ErrOpenAuthExp 赋值：14003 // 临时开柜授权失效

	ErrOpenAuthExp = 14003 // 临时开柜授权失效
	// 为 ErrMakeup 赋值：14004 // 补缴失败

	ErrMakeup      = 14004 // 补缴失败

	// 个人中心（06）
	// 为 ErrAvatarType 赋值：15001 // 头像格式错误

	ErrAvatarType = 15001 // 头像格式错误
	// 为 ErrAvatarSize 赋值：15002 // 头像过大

	ErrAvatarSize = 15002 // 头像过大
	// 为 ErrCouponExp 赋值：15003 // 卡券已过期/不可用

	ErrCouponExp  = 15003 // 卡券已过期/不可用

	// 通知（01）
	// 为 ErrNotifyTpl 赋值：16001 // 通知模板不存在

	ErrNotifyTpl = 16001 // 通知模板不存在

	// IoT 设备接入（19）
	// 为 ErrDeviceNF 赋值：19001 // 设备不存在

	ErrDeviceNF    = 19001 // 设备不存在
	// 为 ErrDeviceSecret 赋值：19002 // 设备密钥校验失败

	ErrDeviceSecret = 19002 // 设备密钥校验失败
	// 为 ErrDeviceToken 赋值：19003 // 设备令牌失效

	ErrDeviceToken = 19003 // 设备令牌失效
	// 为 ErrDeviceState 赋值：19004 // 设备状态不允许该操作

	ErrDeviceState = 19004 // 设备状态不允许该操作
	// 为 ErrCmdFail 赋值：19005 // 指令下发失败

	ErrCmdFail     = 19005 // 指令下发失败
	// 为 ErrOtaNF 赋值：19006 // OTA 任务不存在

	ErrOtaNF       = 19006 // OTA 任务不存在

	// 数字人账号登录（07）：账号或密码不匹配（含账号不存在/禁用），统一码避免账号枚举
	// 为 ErrDHLogin 赋值：17001 // 数字人账号或密码不正确

	ErrDHLogin = 17001 // 数字人账号或密码不正确
// 结束参数/表达式列表

)

// 定义结构体 ErrInfo：用于承载一组相关的字段数据

type ErrInfo struct {
	// 结构体字段 Code，类型为 int

	Code    int
	// 结构体字段 Message，类型为 string

	Message string
// 当前代码块结束

}

// 声明变量 errMessages = map[int]string{，用于保存运行过程中的数据

var errMessages = map[int]string{
	// 执行 CodeOK: "ok",，完成对应的业务处理

	CodeOK: "ok",
	// 执行 ErrBad: "参数错误", ErrAuth: "登录已失效", ErrForb: "禁止访问", ErrNF: "资源不存在",，完成对应的业务处理

	ErrBad: "参数错误", ErrAuth: "登录已失效", ErrForb: "禁止访问", ErrNF: "资源不存在",
	// 执行 ErrInternal: "系统内部错误", ErrLimit: "操作过于频繁", ErrDup: "重复请求", ErrLock: "操作冲突，请重试",，完成对应的业务处理

	ErrInternal: "系统内部错误", ErrLimit: "操作过于频繁", ErrDup: "重复请求", ErrLock: "操作冲突，请重试",
	// 执行 ErrLocFail: "定位失败", ErrCityNF: "城市不存在",，完成对应的业务处理

	ErrLocFail: "定位失败", ErrCityNF: "城市不存在",
	// 执行 ErrBalanceLow: "余额不足", ErrRechargeNF: "充值单不存在", ErrPayCallback: "支付回调校验失败",，完成对应的业务处理

	ErrBalanceLow: "余额不足", ErrRechargeNF: "充值单不存在", ErrPayCallback: "支付回调校验失败",
	// 执行 ErrLockerBusy: "柜门不可用", ErrBillingRule: "收费规则未匹配", ErrDepositNF: "寄存记录不存在",，完成对应的业务处理

	ErrLockerBusy: "柜门不可用", ErrBillingRule: "收费规则未匹配", ErrDepositNF: "寄存记录不存在",
	// 执行 ErrOrderNF: "订单不存在", ErrOrderState: "订单状态不允许该操作", ErrOpenAuthExp: "临时开柜授权失效", ErrMakeup: "补缴失败",，完成对应的业务处理

	ErrOrderNF: "订单不存在", ErrOrderState: "订单状态不允许该操作", ErrOpenAuthExp: "临时开柜授权失效", ErrMakeup: "补缴失败",
	// 执行 ErrAvatarType: "头像格式错误", ErrAvatarSize: "头像过大", ErrCouponExp: "卡券已过期或不可用",，完成对应的业务处理

	ErrAvatarType: "头像格式错误", ErrAvatarSize: "头像过大", ErrCouponExp: "卡券已过期或不可用",
	// 执行 ErrNotifyTpl: "通知模板不存在",，完成对应的业务处理

	ErrNotifyTpl: "通知模板不存在",
	// 执行 ErrDeviceNF: "设备不存在", ErrDeviceSecret: "设备密钥校验失败", ErrDeviceToken: "设备令牌失效",，完成对应的业务处理

	ErrDeviceNF: "设备不存在", ErrDeviceSecret: "设备密钥校验失败", ErrDeviceToken: "设备令牌失效",
	// 执行 ErrDeviceState: "设备状态不允许该操作", ErrCmdFail: "指令下发失败", ErrOtaNF: "OTA任务不存在",，完成对应的业务处理

	ErrDeviceState: "设备状态不允许该操作", ErrCmdFail: "指令下发失败", ErrOtaNF: "OTA任务不存在",

	ErrDHLogin: "账号或密码不正确",
// 当前代码块结束

}

// 定义函数 GetErr：封装一段可复用的处理逻辑（参数：code int）

func GetErr(code int) ErrInfo {
	// 定义并初始化局部变量 msg, ok，取值为 errMessages[code]

	msg, ok := errMessages[code]
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 为 msg 赋值："系统内部错误"

		msg = "系统内部错误"
	// 当前代码块结束

	}
	// 返回 ErrInfo{Code: code, Message: msg} 给调用方，结束当前函数

	return ErrInfo{Code: code, Message: msg}
// 当前代码块结束

}
