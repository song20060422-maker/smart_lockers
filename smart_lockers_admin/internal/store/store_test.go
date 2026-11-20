// @ticket B端-存储层单元测试（事务/幂等/CRUD/快照）
package store

import (
	"errors"
	"testing"

	"smart_lockers_admin/internal/model"
)

// 定义函数 TestNewSeedAndReads：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestNewSeedAndReads(t *testing.T) {
	// 定义并初始化局部变量 s，取值为 New()

	s := New()
	// 条件判断：len(s.Sites()) != 2 || len(s.Devices(0)) != 2 || len(s.Recons(0)) != 1 || len(s.Alarms(0)) != 1 时执行下面逻辑

	if len(s.Sites()) != 2 || len(s.Devices(0)) != 2 || len(s.Recons(0)) != 1 || len(s.Alarms(0)) != 1 {
		// 执行 t.Fatalf("seed missing: sites=%d dev=%d recon=%d alarm=%d", len(s.Sites()), len(s.Devices(0)), len(s.Recons(0)), len(s.Alarms(0)))，完成对应的业务处理

		t.Fatalf("seed missing: sites=%d dev=%d recon=%d alarm=%d", len(s.Sites()), len(s.Devices(0)), len(s.Recons(0)), len(s.Alarms(0)))
	// 当前代码块结束

	}
	// 条件判断：_, ok := s.Site(1001); !ok { t.Fatal("site 1001 missing") } 时执行下面逻辑

	if _, ok := s.Site(1001); !ok { t.Fatal("site 1001 missing") }
	// 条件判断：_, ok := s.Site(9999); ok { t.Fatal("unexpected site") } 时执行下面逻辑

	if _, ok := s.Site(9999); ok { t.Fatal("unexpected site") }
	// 条件判断：len(s.Devices(1001)) != 2 || len(s.Devices(9999)) != 0 { t.Fatal("device scope filter failed") } 时执行下面逻辑

	if len(s.Devices(1001)) != 2 || len(s.Devices(9999)) != 0 { t.Fatal("device scope filter failed") }
	// 条件判断：_, ok := s.Device(2001); !ok { t.Fatal("device 2001 missing") } 时执行下面逻辑

	if _, ok := s.Device(2001); !ok { t.Fatal("device 2001 missing") }
	// 条件判断：_, ok := s.Device(9999); ok { t.Fatal("unexpected device") } 时执行下面逻辑

	if _, ok := s.Device(9999); ok { t.Fatal("unexpected device") }
	// 条件判断：len(s.Recons(1001)) != 1 || len(s.Recons(9999)) != 0 { t.Fatal("recon scope filter failed") } 时执行下面逻辑

	if len(s.Recons(1001)) != 1 || len(s.Recons(9999)) != 0 { t.Fatal("recon scope filter failed") }
	// 条件判断：_, ok := s.Alarm(4001); !ok { t.Fatal("alarm 4001 missing") } 时执行下面逻辑

	if _, ok := s.Alarm(4001); !ok { t.Fatal("alarm 4001 missing") }
	// 条件判断：_, ok := s.Alarm(9999); ok { t.Fatal("unexpected alarm") } 时执行下面逻辑

	if _, ok := s.Alarm(9999); ok { t.Fatal("unexpected alarm") }
// 当前代码块结束

}

// 定义函数 TestCRUD：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCRUD(t *testing.T) {
	// 定义并初始化局部变量 s，取值为 New()

	s := New()
	// 执行 s.AddSite(&model.Site{Code: "S3", Name: "n3"})，完成对应的业务处理

	s.AddSite(&model.Site{Code: "S3", Name: "n3"})
	// 定义并初始化局部变量 sites，取值为 s.Sites()

	sites := s.Sites()
	// 条件判断：len(sites) != 3 { t.Fatalf("add site failed: %d", len(sites)) } 时执行下面逻辑

	if len(sites) != 3 { t.Fatalf("add site failed: %d", len(sites)) }
	// 定义并初始化局部变量 newID，取值为 sites[len(sites)-1].ID

	newID := sites[len(sites)-1].ID
	// 条件判断：!s.UpdateSite(newID, func(v *model.Site) bool { v.Status = 0; return true }) { t.Fatal("update site failed") } 时执行下面逻辑

	if !s.UpdateSite(newID, func(v *model.Site) bool { v.Status = 0; return true }) { t.Fatal("update site failed") }
	// 条件判断：v, _ := s.Site(newID); v.Status != 0 { t.Fatal("site status not updated") } 时执行下面逻辑

	if v, _ := s.Site(newID); v.Status != 0 { t.Fatal("site status not updated") }
	// 条件判断：s.UpdateSite(9999, func(v *model.Site) bool { return true }) { t.Fatal("update missing site should fail") } 时执行下面逻辑

	if s.UpdateSite(9999, func(v *model.Site) bool { return true }) { t.Fatal("update missing site should fail") }
	// 条件判断：s.UpdateSite(newID, func(v *model.Site) bool { return false }) { t.Fatal("no-op mutate should report false") } 时执行下面逻辑

	if s.UpdateSite(newID, func(v *model.Site) bool { return false }) { t.Fatal("no-op mutate should report false") }

	// 条件判断：!s.UpdateDevice(2001, func(d *model.Device) bool { d.Status = 1; return true }) { t.Fatal("update device failed") } 时执行下面逻辑

	if !s.UpdateDevice(2001, func(d *model.Device) bool { d.Status = 1; return true }) { t.Fatal("update device failed") }
	// 条件判断：d, _ := s.Device(2001); d.Status != 1 { t.Fatal("device not updated") } 时执行下面逻辑

	if d, _ := s.Device(2001); d.Status != 1 { t.Fatal("device not updated") }
	// 条件判断：s.UpdateDevice(9999, func(d *model.Device) bool { return true }) { t.Fatal("update missing device should fail") } 时执行下面逻辑

	if s.UpdateDevice(9999, func(d *model.Device) bool { return true }) { t.Fatal("update missing device should fail") }
	// 条件判断：s.UpdateDevice(2001, func(d *model.Device) bool { return false }) { t.Fatal("no-op device should false") } 时执行下面逻辑

	if s.UpdateDevice(2001, func(d *model.Device) bool { return false }) { t.Fatal("no-op device should false") }

	// 执行 s.AddRecon(&model.ReconcileRecord{BatchNo: "B2"})，完成对应的业务处理

	s.AddRecon(&model.ReconcileRecord{BatchNo: "B2"})
	// 条件判断：len(s.Recons(0)) != 2 { t.Fatal("add recon failed") } 时执行下面逻辑

	if len(s.Recons(0)) != 2 { t.Fatal("add recon failed") }
	// 条件判断：!s.UpdateAlarm(4001, func(a *model.Alarm) bool { a.Status = 1; return true }) { t.Fatal("update alarm failed") } 时执行下面逻辑

	if !s.UpdateAlarm(4001, func(a *model.Alarm) bool { a.Status = 1; return true }) { t.Fatal("update alarm failed") }
	// 条件判断：a, _ := s.Alarm(4001); a.Status != 1 { t.Fatal("alarm not updated") } 时执行下面逻辑

	if a, _ := s.Alarm(4001); a.Status != 1 { t.Fatal("alarm not updated") }
	// 条件判断：s.UpdateAlarm(9999, func(a *model.Alarm) bool { return true }) { t.Fatal("update missing alarm should fail") } 时执行下面逻辑

	if s.UpdateAlarm(9999, func(a *model.Alarm) bool { return true }) { t.Fatal("update missing alarm should fail") }
	// 条件判断：s.UpdateAlarm(4001, func(a *model.Alarm) bool { return false }) { t.Fatal("no-op alarm should false") } 时执行下面逻辑

	if s.UpdateAlarm(4001, func(a *model.Alarm) bool { return false }) { t.Fatal("no-op alarm should false") }
// 当前代码块结束

}

// 定义函数 TestTxCommitAndRollback：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestTxCommitAndRollback(t *testing.T) {
	// 定义并初始化局部变量 s，取值为 New()

	s := New()
	// 定义并初始化局部变量 before，取值为 len(s.Sites())

	before := len(s.Sites())
	// 定义并初始化局部变量 err，取值为 s.Tx(func() error {

	err := s.Tx(func() error {
		// 执行 s.AddSiteNoLock(&model.Site{Code: "TX", Name: "tx"})，完成对应的业务处理

		s.AddSiteNoLock(&model.Site{Code: "TX", Name: "tx"})
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil || len(s.Sites()) != before+1 { t.Fatalf("commit failed err=%v n=%d", err, len(s.Sites())) } 时执行下面逻辑

	if err != nil || len(s.Sites()) != before+1 { t.Fatalf("commit failed err=%v n=%d", err, len(s.Sites())) }
	// 定义并初始化局部变量 wantErr，取值为 errors.New("boom")

	wantErr := errors.New("boom")
	// 为 err 赋值：s.Tx(func() error {

	err = s.Tx(func() error {
		// 执行 s.AddSiteNoLock(&model.Site{Code: "ROLLBACK", Name: "rb"})，完成对应的业务处理

		s.AddSiteNoLock(&model.Site{Code: "ROLLBACK", Name: "rb"})
		// 返回 wantErr 给调用方，结束当前函数

		return wantErr
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != wantErr { t.Fatalf("rollback err=%v", err) } 时执行下面逻辑

	if err != wantErr { t.Fatalf("rollback err=%v", err) }
	// 条件判断：len(s.Sites()) != before+1 { t.Fatalf("rollback left data n=%d", len(s.Sites())) } 时执行下面逻辑

	if len(s.Sites()) != before+1 { t.Fatalf("rollback left data n=%d", len(s.Sites())) }
// 当前代码块结束

}

// 定义函数 TestIdempotency：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestIdempotency(t *testing.T) {
	// 定义并初始化局部变量 s，取值为 New()

	s := New()
	// 条件判断：s.CheckIdem("") { t.Fatal("empty key should be rejected") } 时执行下面逻辑

	if s.CheckIdem("") { t.Fatal("empty key should be rejected") }
	// 条件判断：!s.CheckIdem("k1") { t.Fatal("first idem should pass") } 时执行下面逻辑

	if !s.CheckIdem("k1") { t.Fatal("first idem should pass") }
	// 条件判断：s.CheckIdem("k1") { t.Fatal("second idem should be rejected") } 时执行下面逻辑

	if s.CheckIdem("k1") { t.Fatal("second idem should be rejected") }
	// 幂等键也纳入事务快照回滚
	// 定义并初始化局部变量 err，取值为 s.Tx(func() error {

	err := s.Tx(func() error {
		// 条件判断：!s.CheckIdemNoLock("k2") { t.Fatal("k2 first") } 时执行下面逻辑

		if !s.CheckIdemNoLock("k2") { t.Fatal("k2 first") }
		// 返回 errors.New("rb") 给调用方，结束当前函数

		return errors.New("rb")
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err == nil { t.Fatal("want rollback") } 时执行下面逻辑

	if err == nil { t.Fatal("want rollback") }
	// 条件判断：!s.CheckIdem("k2") { t.Fatal("idem key should be rolled back") } 时执行下面逻辑

	if !s.CheckIdem("k2") { t.Fatal("idem key should be rolled back") }
// 当前代码块结束

}
