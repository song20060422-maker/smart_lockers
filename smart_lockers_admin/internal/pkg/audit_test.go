// @ticket B端-审计单元测试
package pkg

import (
	"testing"

	"smart_lockers_admin/internal/config"
)

// 定义函数 TestRecordAuditAndSnapshot：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestRecordAuditAndSnapshot(t *testing.T) {
	// 为 auditMu.Lock(); auditLogs 赋值：nil; auditNextID = 1; auditMu.Unlock()

	auditMu.Lock(); auditLogs = nil; auditNextID = 1; auditMu.Unlock()
	// 执行 config.ResetForTest()，完成对应的业务处理

	config.ResetForTest()
	// 执行 RecordAudit("uid:1", "operator", "reconcile.run", "batch=1", "success")，完成对应的业务处理

	RecordAudit("uid:1", "operator", "reconcile.run", "batch=1", "success")
	// 执行 RecordAudit("uid:2", "site_admin", "alarm.handle", "4001", "success")，完成对应的业务处理

	RecordAudit("uid:2", "site_admin", "alarm.handle", "4001", "success")
	// 定义并初始化局部变量 recs，取值为 AuditRecords()

	recs := AuditRecords()
	// 条件判断：len(recs) != 2 || recs[0].Operator != "uid:1" || recs[1].Action != "alarm.handle" 时执行下面逻辑

	if len(recs) != 2 || recs[0].Operator != "uid:1" || recs[1].Action != "alarm.handle" {
		// 执行 t.Fatalf("unexpected audit records: %+v", recs)，完成对应的业务处理

		t.Fatalf("unexpected audit records: %+v", recs)
	// 当前代码块结束

	}
	// 快照应隔离：改返回值不影响内部
	// 为 recs[0].Operator 赋值："tampered"

	recs[0].Operator = "tampered"
	// 条件判断：AuditRecords()[0].Operator != "uid:1" { t.Fatal("snapshot not isolated") } 时执行下面逻辑

	if AuditRecords()[0].Operator != "uid:1" { t.Fatal("snapshot not isolated") }
	// 上限截断
	// 定义并初始化局部变量 cfg，取值为 config.Get(); cfg.Policy.AuditMaxRecords = 3

	cfg := config.Get(); cfg.Policy.AuditMaxRecords = 3
	// 循环处理：for i := 0; i < 5; i++ { RecordAudit("u", "r", "a", "t", "ok") }，逐轮迭代

	for i := 0; i < 5; i++ { RecordAudit("u", "r", "a", "t", "ok") }
	// 条件判断：got := len(AuditRecords()); got != 3 { t.Fatalf("audit trim want 3 got %d", got) } 时执行下面逻辑

	if got := len(AuditRecords()); got != 3 { t.Fatalf("audit trim want 3 got %d", got) }
	// 执行 config.ResetForTest()，完成对应的业务处理

	config.ResetForTest()
// 当前代码块结束

}
