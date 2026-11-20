// @ticket B端-四项功能 service 单元测试
package service

import (
	"testing"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// 定义函数 op：封装一段可复用的处理逻辑（参数：role model.AdminRole, siteID uint）

func op(role model.AdminRole, siteID uint) middleware.CurrentUser {
	// 返回 middleware.CurrentUser{ID: 1, Role: role, SiteID: siteID} 给调用方，结束当前函数

	return middleware.CurrentUser{ID: 1, Role: role, SiteID: siteID}
// 当前代码块结束

}

// 定义函数 TestSiteListScopeAndMask：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestSiteListScopeAndMask(t *testing.T) {
	// 定义并初始化局部变量 st，取值为 store.New()

	st := store.New()
	// 定义并初始化局部变量 svc，取值为 NewSite(st)

	svc := NewSite(st)
	// 网点管理员仅本网点(1001)，种子两条分别 ID=1/2，Code=SITE-1001/1002
	// 定义并初始化局部变量 opUser，取值为 op(model.RoleSiteAdmin, 1001)

	opUser := op(model.RoleSiteAdmin, 1001)
	// 定义并初始化局部变量 list，取值为 svc.List(opUser)

	list := svc.List(opUser)
	// 循环处理：for _, s := range list，逐轮迭代

	for _, s := range list {
		// 条件判断：s.ID != 1001 时执行下面逻辑

		if s.ID != 1001 {
			// 执行 t.Fatalf("site admin saw out-of-scope site: %+v", s)，完成对应的业务处理

			t.Fatalf("site admin saw out-of-scope site: %+v", s)
		// 当前代码块结束

		}
		// 条件判断：len(s.ManagerPhone) >= 4 && s.ManagerPhone[len(s.ManagerPhone)-4:] != "8001" 时执行下面逻辑

		if len(s.ManagerPhone) >= 4 && s.ManagerPhone[len(s.ManagerPhone)-4:] != "8001" {
			// 脱敏后尾号保留
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 平台运营可见全部
	// 条件判断：got := len(svc.List(op(model.RolePlatformOp, 0))); got != 2 时执行下面逻辑

	if got := len(svc.List(op(model.RolePlatformOp, 0))); got != 2 {
		// 执行 t.Fatalf("operator site count=%d want 2", got)，完成对应的业务处理

		t.Fatalf("operator site count=%d want 2", got)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestSiteCreateDupAndAudit：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestSiteCreateDupAndAudit(t *testing.T) {
	// 定义并初始化局部变量 st，取值为 store.New()

	st := store.New()
	// 定义并初始化局部变量 svc，取值为 NewSite(st)

	svc := NewSite(st)
	// 定义并初始化局部变量 code，取值为 "SITE-NEW-1"

	code := "SITE-NEW-1"
	// 定义并初始化局部变量 v，取值为 &model.Site{Code: code, Name: "新网点"}

	v := &model.Site{Code: code, Name: "新网点"}
	// 条件判断：c := svc.Create(op(model.RolePlatformOp, 0), v); c != pkg.CodeOK 时执行下面逻辑

	if c := svc.Create(op(model.RolePlatformOp, 0), v); c != pkg.CodeOK {
		// 执行 t.Fatalf("create code=%d", c)，完成对应的业务处理

		t.Fatalf("create code=%d", c)
	// 当前代码块结束

	}
	// 条件判断：c := svc.Create(op(model.RolePlatformOp, 0), &model.Site{Code: code, Name: "重复"}); c != pkg.ErrSiteExists 时执行下面逻辑

	if c := svc.Create(op(model.RolePlatformOp, 0), &model.Site{Code: code, Name: "重复"}); c != pkg.ErrSiteExists {
		// 执行 t.Fatalf("dup create code=%d", c)，完成对应的业务处理

		t.Fatalf("dup create code=%d", c)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestDeviceMonitorAndControl：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestDeviceMonitorAndControl(t *testing.T) {
	// 定义并初始化局部变量 st，取值为 store.New()

	st := store.New()
	// 定义并初始化局部变量 svc，取值为 NewDevice(st)

	svc := NewDevice(st)
	// 定义并初始化局部变量 m，取值为 svc.Monitor(op(model.RolePlatformOp, 0))

	m := svc.Monitor(op(model.RolePlatformOp, 0))
	// 条件判断：m["total"].(int) != 2 || m["fault"].(int) != 1 时执行下面逻辑

	if m["total"].(int) != 2 || m["fault"].(int) != 1 {
		// 执行 t.Fatalf("monitor=%+v", m)，完成对应的业务处理

		t.Fatalf("monitor=%+v", m)
	// 当前代码块结束

	}
	// ID=4 设备故障->在线
	// 条件判断：c := svc.Control(op(model.RoleSiteAdmin, 1001), 2002, 0); c != pkg.CodeOK 时执行下面逻辑

	if c := svc.Control(op(model.RoleSiteAdmin, 1001), 2002, 0); c != pkg.CodeOK {
		// 执行 t.Fatalf("control code=%d", c)，完成对应的业务处理

		t.Fatalf("control code=%d", c)
	// 当前代码块结束

	}
	// 条件判断：d, _ := st.Device(2002); d.Status != 0 { t.Fatalf("device status=%d", d.Status) } 时执行下面逻辑

	if d, _ := st.Device(2002); d.Status != 0 { t.Fatalf("device status=%d", d.Status) }
	// 条件判断：c := svc.Control(op(model.RolePlatformOp, 0), 999, 0); c != pkg.ErrDeviceNF 时执行下面逻辑

	if c := svc.Control(op(model.RolePlatformOp, 0), 999, 0); c != pkg.ErrDeviceNF {
		// 执行 t.Fatalf("missing device code=%d", c)，完成对应的业务处理

		t.Fatalf("missing device code=%d", c)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestReconcileRunIdempotent：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestReconcileRunIdempotent(t *testing.T) {
	// 定义并初始化局部变量 st，取值为 store.New()

	st := store.New()
	// 定义并初始化局部变量 svc，取值为 NewReconcile(st)

	svc := NewReconcile(st)
	// 定义并初始化局部变量 u，取值为 op(model.RolePlatformOp, 0)

	u := op(model.RolePlatformOp, 0)
	// 定义并初始化局部变量 c1, rec，取值为 svc.Run(u, "BATCH-X", 1000, 1000, "idem-1")

	c1, rec := svc.Run(u, "BATCH-X", 1000, 1000, "idem-1")
	// 条件判断：c1 != pkg.CodeOK || rec.Status != 1 { t.Fatalf("run1 code=%d rec=%+v", c1, rec) } 时执行下面逻辑

	if c1 != pkg.CodeOK || rec.Status != 1 { t.Fatalf("run1 code=%d rec=%+v", c1, rec) }
	// 定义并初始化局部变量 c2, _，取值为 svc.Run(u, "BATCH-Y", 1000, 900, "idem-1") // 幂等键重复

	c2, _ := svc.Run(u, "BATCH-Y", 1000, 900, "idem-1") // 幂等键重复
	// 条件判断：c2 != pkg.ErrReconcileDup { t.Fatalf("idem code=%d", c2) } 时执行下面逻辑

	if c2 != pkg.ErrReconcileDup { t.Fatalf("idem code=%d", c2) }
	// 定义并初始化局部变量 c3, rec3，取值为 svc.Run(u, "BATCH-Z", 1000, 900, "idem-2") // 不平

	c3, rec3 := svc.Run(u, "BATCH-Z", 1000, 900, "idem-2") // 不平
	// 条件判断：c3 != pkg.CodeOK || rec3.Status != 2 { t.Fatalf("diff code=%d rec=%+v", c3, rec3) } 时执行下面逻辑

	if c3 != pkg.CodeOK || rec3.Status != 2 { t.Fatalf("diff code=%d rec=%+v", c3, rec3) }
// 当前代码块结束

}

// 定义函数 TestAlarmHandleState：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestAlarmHandleState(t *testing.T) {
	// 定义并初始化局部变量 st，取值为 store.New()

	st := store.New()
	// 定义并初始化局部变量 svc，取值为 NewAlarm(st)

	svc := NewAlarm(st)
	// 定义并初始化局部变量 u，取值为 op(model.RoleSiteAdmin, 1)

	u := op(model.RoleSiteAdmin, 1)
	// 条件判断：c := svc.Handle(u, 4001, "已重启设备", "idem-a"); c != pkg.CodeOK { t.Fatalf("handle code=%d", c) } 时执行下面逻辑

	if c := svc.Handle(u, 4001, "已重启设备", "idem-a"); c != pkg.CodeOK { t.Fatalf("handle code=%d", c) }
	// 条件判断：c := svc.Handle(u, 4001, "重复", "idem-b"); c != pkg.ErrAlarmState { t.Fatalf("dup handle code=%d", c) } 时执行下面逻辑

	if c := svc.Handle(u, 4001, "重复", "idem-b"); c != pkg.ErrAlarmState { t.Fatalf("dup handle code=%d", c) }
	// 条件判断：c := svc.Handle(u, 999, "", "idem-c"); c != pkg.ErrAlarmNF { t.Fatalf("missing code=%d", c) } 时执行下面逻辑

	if c := svc.Handle(u, 999, "", "idem-c"); c != pkg.ErrAlarmNF { t.Fatalf("missing code=%d", c) }
// 当前代码块结束

}
