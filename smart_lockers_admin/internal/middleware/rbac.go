// @ticket B端-RBAC：三类角色权限矩阵与数据域约束
package middleware

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
)

// 权限点 @ticket B端-RBAC
const (
	PermSiteRead   = "site:read"
	PermSiteWrite  = "site:write"
	PermDeviceRead = "device:read"
	PermDeviceCtl  = "device:control"
	PermReconRead  = "reconcile:read"
	PermReconRun   = "reconcile:run"
	PermAlarmRead  = "alarm:read"
	PermAlarmHandle = "alarm:handle"
	PermAuditRead  = "audit:read"
)

// rolePermissions 角色->权限映射（可配置化：生产可迁移至配置/DB）@ticket B端-RBAC
var rolePermissions = map[model.AdminRole]map[string]bool{
	model.RoleNormalUser: {
		PermSiteRead: true,
	},
	model.RoleSiteAdmin: {
		PermSiteRead: true, PermSiteWrite: true,
		PermDeviceRead: true, PermDeviceCtl: true,
		PermReconRead: true,
		PermAlarmRead: true, PermAlarmHandle: true,
	},
	model.RolePlatformOp: {
		PermSiteRead: true, PermSiteWrite: true,
		PermDeviceRead: true, PermDeviceCtl: true,
		PermReconRead: true, PermReconRun: true,
		PermAlarmRead: true, PermAlarmHandle: true,
		PermAuditRead: true,
	},
}

// HasPerm 角色是否拥有权限
func HasPerm(role model.AdminRole, perm string) bool {
	return rolePermissions[role][perm]
}

// RequirePerm 权限中间件 @ticket B端-RBAC
func RequirePerm(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := User(c)
		if !HasPerm(u.Role, perm) {
			pkg.AbortWithCode(c, pkg.ErrForb)
			return
		}
		c.Next()
	}
}

// SiteScope 数据域：网点管理员仅限本网点，平台运营可跨网点(0表示全部)@ticket B端-RBAC
func SiteScope(u CurrentUser) uint {
	if u.Role == model.RoleSiteAdmin {
		return u.SiteID
	}
	return 0
}
