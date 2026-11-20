// @ticket B端-统一错误码：20xxx 分段，全项目唯一错误码表
package pkg

const (
	CodeOK      = 0
	ErrBad      = 10001 // 参数错误
	ErrAuth     = 10002 // 未授权/登录失效
	ErrForb     = 10003 // 禁止访问(角色权限不足)
	ErrNF       = 10004 // 资源不存在
	ErrInternal = 10005 // 系统内部错误
	ErrLimit    = 10006 // 限流/频控
	ErrDup      = 10007 // 重复请求(幂等命中)
	ErrLock     = 10008 // 并发冲突/锁失败

	// B端 网点管理(20)
	ErrSiteExists = 20001 // 网点编码已存在
	ErrSiteStatus = 20002 // 网点状态不允许该操作

	// B端 设备监控(21)
	ErrDeviceNF    = 21001 // 设备不存在
	ErrDeviceState = 21002 // 设备状态不允许该操作

	// B端 充值对账(22)
	ErrReconcileDiff  = 22001 // 对账不平
	ErrReconcileDup   = 22002 // 对账批次重复(幂等)
	ErrReconcileState = 22003 // 对账状态不允许

	// B端 告警处理(23)
	ErrAlarmNF    = 23001 // 告警不存在
	ErrAlarmState = 23002 // 告警状态不允许重复处理
)

type ErrInfo struct {
	Code    int
	Message string
}

var errMessages = map[int]string{
	CodeOK: "ok",
	ErrBad: "参数错误", ErrAuth: "登录已失效", ErrForb: "无权限执行该操作", ErrNF: "资源不存在",
	ErrInternal: "系统内部错误", ErrLimit: "操作过于频繁", ErrDup: "重复请求", ErrLock: "操作冲突，请重试",
	ErrSiteExists: "网点编码已存在", ErrSiteStatus: "网点状态不允许该操作",
	ErrDeviceNF: "设备不存在", ErrDeviceState: "设备状态不允许该操作",
	ErrReconcileDiff: "对账不平", ErrReconcileDup: "对账批次重复", ErrReconcileState: "对账状态不允许",
	ErrAlarmNF: "告警不存在", ErrAlarmState: "告警已处理，不可重复处理",
}

func GetErr(code int) ErrInfo {
	msg, ok := errMessages[code]
	if !ok {
		msg = "系统内部错误"
	}
	return ErrInfo{Code: code, Message: msg}
}
