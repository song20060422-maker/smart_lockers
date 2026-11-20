// @ticket B端-审计单元测试
package pkg

import (
	"testing"

	"smart_lockers_admin/internal/config"
)

func TestRecordAuditAndSnapshot(t *testing.T) {
	auditMu.Lock(); auditLogs = nil; auditNextID = 1; auditMu.Unlock()
	config.ResetForTest()
	RecordAudit("uid:1", "operator", "reconcile.run", "batch=1", "success")
	RecordAudit("uid:2", "site_admin", "alarm.handle", "4001", "success")
	recs := AuditRecords()
	if len(recs) != 2 || recs[0].Operator != "uid:1" || recs[1].Action != "alarm.handle" {
		t.Fatalf("unexpected audit records: %+v", recs)
	}
	// 快照应隔离：改返回值不影响内部
	recs[0].Operator = "tampered"
	if AuditRecords()[0].Operator != "uid:1" { t.Fatal("snapshot not isolated") }
	// 上限截断
	cfg := config.Get(); cfg.Policy.AuditMaxRecords = 3
	for i := 0; i < 5; i++ { RecordAudit("u", "r", "a", "t", "ok") }
	if got := len(AuditRecords()); got != 3 { t.Fatalf("audit trim want 3 got %d", got) }
	config.ResetForTest()
}
