// @ticket B端-审计：资金/状态变更写操作留痕，可配置保留条数
package pkg

import (
	"sync"
	"time"

	"smart_lockers_admin/internal/config"
)

// AuditLog 审计记录（操作人、动作、对象、结果、时间；不含敏感明文）
type AuditLog struct {
	ID        uint      `json:"id"`
	Operator  string    `json:"operator"`
	Role      string    `json:"role"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Result    string    `json:"result"`
	At        time.Time `json:"at"`
}

var (
	auditMu     sync.Mutex
	auditLogs   []AuditLog
	auditNextID uint = 1
)

// RecordAudit 追加审计记录，超过上限丢弃最旧记录 @ticket B端-审计
func RecordAudit(operator, role, action, target, result string) {
	auditMu.Lock()
	defer auditMu.Unlock()
	auditLogs = append(auditLogs, AuditLog{
		ID: auditNextID, Operator: operator, Role: role,
		Action: action, Target: target, Result: result, At: time.Now(),
	})
	auditNextID++
	max := config.Get().Policy.AuditMaxRecords
	if max <= 0 {
		max = 500
	}
	if len(auditLogs) > max {
		auditLogs = auditLogs[len(auditLogs)-max:]
	}
}

// AuditRecords 审计记录快照（只读）
func AuditRecords() []AuditLog {
	auditMu.Lock()
	defer auditMu.Unlock()
	out := make([]AuditLog, len(auditLogs))
	copy(out, auditLogs)
	return out
}
