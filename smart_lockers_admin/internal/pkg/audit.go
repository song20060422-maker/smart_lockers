// @ticket B端-审计：资金/状态变更写操作留痕，可配置保留条数
package pkg

import (
	"sync"
	"time"

	"smart_lockers_admin/internal/config"
)

// AuditLog 审计记录（操作人、动作、对象、结果、时间；不含敏感明文）
// 定义结构体 AuditLog：用于承载一组相关的字段数据

type AuditLog struct {
	// 结构体字段 ID，类型为 uint，附带标签 `json:"id"`

	ID        uint      `json:"id"`
	// 结构体字段 Operator，类型为 string，附带标签 `json:"operator"`

	Operator  string    `json:"operator"`
	// 结构体字段 Role，类型为 string，附带标签 `json:"role"`

	Role      string    `json:"role"`
	// 结构体字段 Action，类型为 string，附带标签 `json:"action"`

	Action    string    `json:"action"`
	// 结构体字段 Target，类型为 string，附带标签 `json:"target"`

	Target    string    `json:"target"`
	// 结构体字段 Result，类型为 string，附带标签 `json:"result"`

	Result    string    `json:"result"`
	// 结构体字段 At，类型为 time.Time，附带标签 `json:"at"`

	At        time.Time `json:"at"`
// 当前代码块结束

}

// 开始定义一组包级/局部变量的 var 声明块

var (
	// 执行 auditMu     sync.Mutex，完成对应的业务处理

	auditMu     sync.Mutex
	// 执行 auditLogs   []AuditLog，完成对应的业务处理

	auditLogs   []AuditLog
	// 为 auditNextID uint 赋值：1

	auditNextID uint = 1
// 结束参数/表达式列表

)

// RecordAudit 追加审计记录，超过上限丢弃最旧记录 @ticket B端-审计
// 定义函数 RecordAudit：封装一段可复用的处理逻辑（参数：operator, role, action, target, result string）

func RecordAudit(operator, role, action, target, result string) {
	// 执行 auditMu.Lock()，完成对应的业务处理

	auditMu.Lock()
	// 延迟注册 auditMu.Unlock()，在当前函数返回前自动执行

	defer auditMu.Unlock()
	// 为 auditLogs 赋值：append(auditLogs, AuditLog{

	auditLogs = append(auditLogs, AuditLog{
		// 执行 ID: auditNextID, Operator: operator, Role: role,，完成对应的业务处理

		ID: auditNextID, Operator: operator, Role: role,
		// 执行 Action: action, Target: target, Result: result, At: time.Now(),，完成对应的业务处理

		Action: action, Target: target, Result: result, At: time.Now(),
	// 执行 })，完成对应的业务处理

	})
	// 执行 auditNextID++，完成对应的业务处理

	auditNextID++
	// 定义并初始化局部变量 max，取值为 config.Get().Policy.AuditMaxRecords

	max := config.Get().Policy.AuditMaxRecords
	// 条件判断：max <= 0 时执行下面逻辑

	if max <= 0 {
		// 为 max 赋值：500

		max = 500
	// 当前代码块结束

	}
	// 条件判断：len(auditLogs) > max 时执行下面逻辑

	if len(auditLogs) > max {
		// 为 auditLogs 赋值：auditLogs[len(auditLogs)-max:]

		auditLogs = auditLogs[len(auditLogs)-max:]
	// 当前代码块结束

	}
// 当前代码块结束

}

// AuditRecords 审计记录快照（只读）
// 定义函数 AuditRecords：封装一段可复用的处理逻辑（参数：）

func AuditRecords() []AuditLog {
	// 执行 auditMu.Lock()，完成对应的业务处理

	auditMu.Lock()
	// 延迟注册 auditMu.Unlock()，在当前函数返回前自动执行

	defer auditMu.Unlock()
	// 定义并初始化局部变量 out，取值为 make([]AuditLog, len(auditLogs))

	out := make([]AuditLog, len(auditLogs))
	// 执行 copy(out, auditLogs)，完成对应的业务处理

	copy(out, auditLogs)
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}
