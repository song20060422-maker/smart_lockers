// @ticket B端-统一错误码：20xxx 分段，全项目唯一错误码表
package pkg

// 开始定义一组常量的 const 声明块

const (
	// 为 CodeOK 赋值：0

	CodeOK      = 0
	// 为 ErrBad 赋值：10001 // 参数错误

	ErrBad      = 10001 // 参数错误
	// 为 ErrAuth 赋值：10002 // 未授权/登录失效

	ErrAuth     = 10002 // 未授权/登录失效
	// 为 ErrForb 赋值：10003 // 禁止访问(角色权限不足)

	ErrForb     = 10003 // 禁止访问(角色权限不足)
	// 为 ErrNF 赋值：10004 // 资源不存在

	ErrNF       = 10004 // 资源不存在
	// 为 ErrInternal 赋值：10005 // 系统内部错误

	ErrInternal = 10005 // 系统内部错误
	// 为 ErrLimit 赋值：10006 // 限流/频控

	ErrLimit    = 10006 // 限流/频控
	// 为 ErrDup 赋值：10007 // 重复请求(幂等命中)

	ErrDup      = 10007 // 重复请求(幂等命中)
	// 为 ErrLock 赋值：10008 // 并发冲突/锁失败

	ErrLock     = 10008 // 并发冲突/锁失败

	// B端 网点管理(20)
	// 为 ErrSiteExists 赋值：20001 // 网点编码已存在

	ErrSiteExists = 20001 // 网点编码已存在
	// 为 ErrSiteStatus 赋值：20002 // 网点状态不允许该操作

	ErrSiteStatus = 20002 // 网点状态不允许该操作

	// B端 设备监控(21)
	// 为 ErrDeviceNF 赋值：21001 // 设备不存在

	ErrDeviceNF    = 21001 // 设备不存在
	// 为 ErrDeviceState 赋值：21002 // 设备状态不允许该操作

	ErrDeviceState = 21002 // 设备状态不允许该操作

	// B端 充值对账(22)
	// 为 ErrReconcileDiff 赋值：22001 // 对账不平

	ErrReconcileDiff  = 22001 // 对账不平
	// 为 ErrReconcileDup 赋值：22002 // 对账批次重复(幂等)

	ErrReconcileDup   = 22002 // 对账批次重复(幂等)
	// 为 ErrReconcileState 赋值：22003 // 对账状态不允许

	ErrReconcileState = 22003 // 对账状态不允许

	// B端 告警处理(23)
	// 为 ErrAlarmNF 赋值：23001 // 告警不存在

	ErrAlarmNF    = 23001 // 告警不存在
	// 为 ErrAlarmState 赋值：23002 // 告警状态不允许重复处理

	ErrAlarmState = 23002 // 告警状态不允许重复处理
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
	// 执行 ErrBad: "参数错误", ErrAuth: "登录已失效", ErrForb: "无权限执行该操作", ErrNF: "资源不存在",，完成对应的业务处理

	ErrBad: "参数错误", ErrAuth: "登录已失效", ErrForb: "无权限执行该操作", ErrNF: "资源不存在",
	// 执行 ErrInternal: "系统内部错误", ErrLimit: "操作过于频繁", ErrDup: "重复请求", ErrLock: "操作冲突，请重试",，完成对应的业务处理

	ErrInternal: "系统内部错误", ErrLimit: "操作过于频繁", ErrDup: "重复请求", ErrLock: "操作冲突，请重试",
	// 执行 ErrSiteExists: "网点编码已存在", ErrSiteStatus: "网点状态不允许该操作",，完成对应的业务处理

	ErrSiteExists: "网点编码已存在", ErrSiteStatus: "网点状态不允许该操作",
	// 执行 ErrDeviceNF: "设备不存在", ErrDeviceState: "设备状态不允许该操作",，完成对应的业务处理

	ErrDeviceNF: "设备不存在", ErrDeviceState: "设备状态不允许该操作",
	// 执行 ErrReconcileDiff: "对账不平", ErrReconcileDup: "对账批次重复", ErrReconcileState: "对账状态不允许",，完成对应的业务处理

	ErrReconcileDiff: "对账不平", ErrReconcileDup: "对账批次重复", ErrReconcileState: "对账状态不允许",
	// 执行 ErrAlarmNF: "告警不存在", ErrAlarmState: "告警已处理，不可重复处理",，完成对应的业务处理

	ErrAlarmNF: "告警不存在", ErrAlarmState: "告警已处理，不可重复处理",
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
