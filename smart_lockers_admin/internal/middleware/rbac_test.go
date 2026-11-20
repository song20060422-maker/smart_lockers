// @ticket B端-RBAC 单元测试：角色权限矩阵、数据域、令牌角色携带
package middleware

import (
	"testing"

	"smart_lockers_admin/internal/model"
)

func TestRolePermissions(t *testing.T) {
	cases := []struct {
		role model.AdminRole
		perm string
		want bool
	}{
		{model.RoleNormalUser, PermSiteRead, true},
		{model.RoleNormalUser, PermSiteWrite, false},
		{model.RoleNormalUser, PermReconRun, false},
		{model.RoleSiteAdmin, PermSiteWrite, true},
		{model.RoleSiteAdmin, PermReconRun, false},
		{model.RoleSiteAdmin, PermAlarmHandle, true},
		{model.RolePlatformOp, PermReconRun, true},
		{model.RolePlatformOp, PermAuditRead, true},
	}
	for _, c := range cases {
		if got := HasPerm(c.role, c.perm); got != c.want {
			t.Fatalf("HasPerm(%s,%s)=%v want %v", c.role, c.perm, got, c.want)
		}
	}
}

func TestSiteScope(t *testing.T) {
	if got := SiteScope(CurrentUser{Role: model.RoleSiteAdmin, SiteID: 1001}); got != 1001 {
		t.Fatalf("site admin scope=%d", got)
	}
	if got := SiteScope(CurrentUser{Role: model.RolePlatformOp, SiteID: 1001}); got != 0 {
		t.Fatalf("operator scope should be 0, got %d", got)
	}
}

func TestTokenCarriesRole(t *testing.T) {
	tok, err := GenerateToken(9, model.RoleSiteAdmin, 1001)
	if err != nil { t.Fatalf("generate: %v", err) }
	claims, err := ParseToken(tok)
	if err != nil { t.Fatalf("parse: %v", err) }
	if claims.Role != model.RoleSiteAdmin || claims.SiteID != 1001 || claims.UserID != 9 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	tampered := tok[:len(tok)-2] + "xx"
	if _, err := ParseToken(tampered); err == nil { t.Fatal("tampered token accepted") }
}
