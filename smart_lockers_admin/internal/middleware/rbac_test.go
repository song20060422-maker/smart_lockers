// @ticket B端-RBAC 单元测试：角色权限矩阵、数据域、令牌角色携带
package middleware

import (
	"testing"

	"smart_lockers_admin/internal/model"
)

// 定义函数 TestRolePermissions：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestRolePermissions(t *testing.T) {
	// 定义并初始化局部变量 cases，取值为 []struct {

	cases := []struct {
		// 执行 role model.AdminRole，完成对应的业务处理

		role model.AdminRole
		// 执行 perm string，完成对应的业务处理

		perm string
		// 执行 want bool，完成对应的业务处理

		want bool
	// 执行 }{，完成对应的业务处理

	}{
		// 执行 {model.RoleNormalUser, PermSiteRead, true},，完成对应的业务处理

		{model.RoleNormalUser, PermSiteRead, true},
		// 执行 {model.RoleNormalUser, PermSiteWrite, false},，完成对应的业务处理

		{model.RoleNormalUser, PermSiteWrite, false},
		// 执行 {model.RoleNormalUser, PermReconRun, false},，完成对应的业务处理

		{model.RoleNormalUser, PermReconRun, false},
		// 执行 {model.RoleSiteAdmin, PermSiteWrite, true},，完成对应的业务处理

		{model.RoleSiteAdmin, PermSiteWrite, true},
		// 执行 {model.RoleSiteAdmin, PermReconRun, false},，完成对应的业务处理

		{model.RoleSiteAdmin, PermReconRun, false},
		// 执行 {model.RoleSiteAdmin, PermAlarmHandle, true},，完成对应的业务处理

		{model.RoleSiteAdmin, PermAlarmHandle, true},
		// 执行 {model.RolePlatformOp, PermReconRun, true},，完成对应的业务处理

		{model.RolePlatformOp, PermReconRun, true},
		// 执行 {model.RolePlatformOp, PermAuditRead, true},，完成对应的业务处理

		{model.RolePlatformOp, PermAuditRead, true},
	// 当前代码块结束

	}
	// 循环处理：for _, c := range cases，逐轮迭代

	for _, c := range cases {
		// 条件判断：got := HasPerm(c.role, c.perm); got != c.want 时执行下面逻辑

		if got := HasPerm(c.role, c.perm); got != c.want {
			// 执行 t.Fatalf("HasPerm(%s,%s)=%v want %v", c.role, c.perm, got, c.want)，完成对应的业务处理

			t.Fatalf("HasPerm(%s,%s)=%v want %v", c.role, c.perm, got, c.want)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestSiteScope：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestSiteScope(t *testing.T) {
	// 条件判断：got := SiteScope(CurrentUser{Role: model.RoleSiteAdmin, SiteID: 1001}); got != 1001 时执行下面逻辑

	if got := SiteScope(CurrentUser{Role: model.RoleSiteAdmin, SiteID: 1001}); got != 1001 {
		// 执行 t.Fatalf("site admin scope=%d", got)，完成对应的业务处理

		t.Fatalf("site admin scope=%d", got)
	// 当前代码块结束

	}
	// 条件判断：got := SiteScope(CurrentUser{Role: model.RolePlatformOp, SiteID: 1001}); got != 0 时执行下面逻辑

	if got := SiteScope(CurrentUser{Role: model.RolePlatformOp, SiteID: 1001}); got != 0 {
		// 执行 t.Fatalf("operator scope should be 0, got %d", got)，完成对应的业务处理

		t.Fatalf("operator scope should be 0, got %d", got)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestTokenCarriesRole：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestTokenCarriesRole(t *testing.T) {
	// 定义并初始化局部变量 tok, err，取值为 GenerateToken(9, model.RoleSiteAdmin, 1001)

	tok, err := GenerateToken(9, model.RoleSiteAdmin, 1001)
	// 条件判断：err != nil { t.Fatalf("generate: %v", err) } 时执行下面逻辑

	if err != nil { t.Fatalf("generate: %v", err) }
	// 定义并初始化局部变量 claims, err，取值为 ParseToken(tok)

	claims, err := ParseToken(tok)
	// 条件判断：err != nil { t.Fatalf("parse: %v", err) } 时执行下面逻辑

	if err != nil { t.Fatalf("parse: %v", err) }
	// 条件判断：claims.Role != model.RoleSiteAdmin || claims.SiteID != 1001 || claims.UserID != 9 时执行下面逻辑

	if claims.Role != model.RoleSiteAdmin || claims.SiteID != 1001 || claims.UserID != 9 {
		// 执行 t.Fatalf("unexpected claims: %+v", claims)，完成对应的业务处理

		t.Fatalf("unexpected claims: %+v", claims)
	// 当前代码块结束

	}
	// 定义并初始化局部变量 tampered，取值为 tok[:len(tok)-2] + "xx"

	tampered := tok[:len(tok)-2] + "xx"
	// 条件判断：_, err := ParseToken(tampered); err == nil { t.Fatal("tampered token accepted") } 时执行下面逻辑

	if _, err := ParseToken(tampered); err == nil { t.Fatal("tampered token accepted") }
// 当前代码块结束

}
