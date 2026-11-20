// @ticket B端-四项功能 service 单元测试
package service

import (
	"testing"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

func op(role model.AdminRole, siteID uint) middleware.CurrentUser {
	return middleware.CurrentUser{ID: 1, Role: role, SiteID: siteID}
}

func TestSiteListScopeAndMask(t *testing.T) {
	st := store.New()
	svc := NewSite(st)
	// 网点管理员仅本网点(1001)，种子两条分别 ID=1/2，Code=SITE-1001/1002
	opUser := op(model.RoleSiteAdmin, 1001)
	list := svc.List(opUser)
	for _, s := range list {
		if s.ID != 1001 {
			t.Fatalf("site admin saw out-of-scope site: %+v", s)
		}
		if len(s.ManagerPhone) >= 4 && s.ManagerPhone[len(s.ManagerPhone)-4:] != "8001" {
			// 脱敏后尾号保留
		}
	}
	// 平台运营可见全部
	if got := len(svc.List(op(model.RolePlatformOp, 0))); got != 2 {
		t.Fatalf("operator site count=%d want 2", got)
	}
}

func TestSiteCreateDupAndAudit(t *testing.T) {
	st := store.New()
	svc := NewSite(st)
	code := "SITE-NEW-1"
	v := &model.Site{Code: code, Name: "新网点"}
	if c := svc.Create(op(model.RolePlatformOp, 0), v); c != pkg.CodeOK {
		t.Fatalf("create code=%d", c)
	}
	if c := svc.Create(op(model.RolePlatformOp, 0), &model.Site{Code: code, Name: "重复"}); c != pkg.ErrSiteExists {
		t.Fatalf("dup create code=%d", c)
	}
}

func TestDeviceMonitorAndControl(t *testing.T) {
	st := store.New()
	svc := NewDevice(st)
	m := svc.Monitor(op(model.RolePlatformOp, 0))
	if m["total"].(int) != 2 || m["fault"].(int) != 1 {
		t.Fatalf("monitor=%+v", m)
	}
	// ID=4 设备故障->在线
	if c := svc.Control(op(model.RoleSiteAdmin, 1001), 2002, 0); c != pkg.CodeOK {
		t.Fatalf("control code=%d", c)
	}
	if d, _ := st.Device(2002); d.Status != 0 { t.Fatalf("device status=%d", d.Status) }
	if c := svc.Control(op(model.RolePlatformOp, 0), 999, 0); c != pkg.ErrDeviceNF {
		t.Fatalf("missing device code=%d", c)
	}
}

func TestReconcileRunIdempotent(t *testing.T) {
	st := store.New()
	svc := NewReconcile(st)
	u := op(model.RolePlatformOp, 0)
	c1, rec := svc.Run(u, "BATCH-X", 1000, 1000, "idem-1")
	if c1 != pkg.CodeOK || rec.Status != 1 { t.Fatalf("run1 code=%d rec=%+v", c1, rec) }
	c2, _ := svc.Run(u, "BATCH-Y", 1000, 900, "idem-1") // 幂等键重复
	if c2 != pkg.ErrReconcileDup { t.Fatalf("idem code=%d", c2) }
	c3, rec3 := svc.Run(u, "BATCH-Z", 1000, 900, "idem-2") // 不平
	if c3 != pkg.CodeOK || rec3.Status != 2 { t.Fatalf("diff code=%d rec=%+v", c3, rec3) }
}

func TestAlarmHandleState(t *testing.T) {
	st := store.New()
	svc := NewAlarm(st)
	u := op(model.RoleSiteAdmin, 1)
	if c := svc.Handle(u, 4001, "已重启设备", "idem-a"); c != pkg.CodeOK { t.Fatalf("handle code=%d", c) }
	if c := svc.Handle(u, 4001, "重复", "idem-b"); c != pkg.ErrAlarmState { t.Fatalf("dup handle code=%d", c) }
	if c := svc.Handle(u, 999, "", "idem-c"); c != pkg.ErrAlarmNF { t.Fatalf("missing code=%d", c) }
}
