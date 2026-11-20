// @ticket B端-存储层单元测试（事务/幂等/CRUD/快照）
package store

import (
	"errors"
	"testing"

	"smart_lockers_admin/internal/model"
)

func TestNewSeedAndReads(t *testing.T) {
	s := New()
	if len(s.Sites()) != 2 || len(s.Devices(0)) != 2 || len(s.Recons(0)) != 1 || len(s.Alarms(0)) != 1 {
		t.Fatalf("seed missing: sites=%d dev=%d recon=%d alarm=%d", len(s.Sites()), len(s.Devices(0)), len(s.Recons(0)), len(s.Alarms(0)))
	}
	if _, ok := s.Site(1001); !ok { t.Fatal("site 1001 missing") }
	if _, ok := s.Site(9999); ok { t.Fatal("unexpected site") }
	if len(s.Devices(1001)) != 2 || len(s.Devices(9999)) != 0 { t.Fatal("device scope filter failed") }
	if _, ok := s.Device(2001); !ok { t.Fatal("device 2001 missing") }
	if _, ok := s.Device(9999); ok { t.Fatal("unexpected device") }
	if len(s.Recons(1001)) != 1 || len(s.Recons(9999)) != 0 { t.Fatal("recon scope filter failed") }
	if _, ok := s.Alarm(4001); !ok { t.Fatal("alarm 4001 missing") }
	if _, ok := s.Alarm(9999); ok { t.Fatal("unexpected alarm") }
}

func TestCRUD(t *testing.T) {
	s := New()
	s.AddSite(&model.Site{Code: "S3", Name: "n3"})
	sites := s.Sites()
	if len(sites) != 3 { t.Fatalf("add site failed: %d", len(sites)) }
	newID := sites[len(sites)-1].ID
	if !s.UpdateSite(newID, func(v *model.Site) bool { v.Status = 0; return true }) { t.Fatal("update site failed") }
	if v, _ := s.Site(newID); v.Status != 0 { t.Fatal("site status not updated") }
	if s.UpdateSite(9999, func(v *model.Site) bool { return true }) { t.Fatal("update missing site should fail") }
	if s.UpdateSite(newID, func(v *model.Site) bool { return false }) { t.Fatal("no-op mutate should report false") }

	if !s.UpdateDevice(2001, func(d *model.Device) bool { d.Status = 1; return true }) { t.Fatal("update device failed") }
	if d, _ := s.Device(2001); d.Status != 1 { t.Fatal("device not updated") }
	if s.UpdateDevice(9999, func(d *model.Device) bool { return true }) { t.Fatal("update missing device should fail") }
	if s.UpdateDevice(2001, func(d *model.Device) bool { return false }) { t.Fatal("no-op device should false") }

	s.AddRecon(&model.ReconcileRecord{BatchNo: "B2"})
	if len(s.Recons(0)) != 2 { t.Fatal("add recon failed") }
	if !s.UpdateAlarm(4001, func(a *model.Alarm) bool { a.Status = 1; return true }) { t.Fatal("update alarm failed") }
	if a, _ := s.Alarm(4001); a.Status != 1 { t.Fatal("alarm not updated") }
	if s.UpdateAlarm(9999, func(a *model.Alarm) bool { return true }) { t.Fatal("update missing alarm should fail") }
	if s.UpdateAlarm(4001, func(a *model.Alarm) bool { return false }) { t.Fatal("no-op alarm should false") }
}

func TestTxCommitAndRollback(t *testing.T) {
	s := New()
	before := len(s.Sites())
	err := s.Tx(func() error {
		s.AddSiteNoLock(&model.Site{Code: "TX", Name: "tx"})
		return nil
	})
	if err != nil || len(s.Sites()) != before+1 { t.Fatalf("commit failed err=%v n=%d", err, len(s.Sites())) }
	wantErr := errors.New("boom")
	err = s.Tx(func() error {
		s.AddSiteNoLock(&model.Site{Code: "ROLLBACK", Name: "rb"})
		return wantErr
	})
	if err != wantErr { t.Fatalf("rollback err=%v", err) }
	if len(s.Sites()) != before+1 { t.Fatalf("rollback left data n=%d", len(s.Sites())) }
}

func TestIdempotency(t *testing.T) {
	s := New()
	if s.CheckIdem("") { t.Fatal("empty key should be rejected") }
	if !s.CheckIdem("k1") { t.Fatal("first idem should pass") }
	if s.CheckIdem("k1") { t.Fatal("second idem should be rejected") }
	// 幂等键也纳入事务快照回滚
	err := s.Tx(func() error {
		if !s.CheckIdemNoLock("k2") { t.Fatal("k2 first") }
		return errors.New("rb")
	})
	if err == nil { t.Fatal("want rollback") }
	if !s.CheckIdem("k2") { t.Fatal("idem key should be rolled back") }
}
