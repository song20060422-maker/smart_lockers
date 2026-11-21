// @ticket B端-RBAC：三类角色权限矩阵与数据域约束
package middleware

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
)

// 权限点 @ticket B端-RBAC
// 开始定义一组常量的 const 声明块

const (
	// 为 PermSiteRead 赋值："site:read"

	PermSiteRead   = "site:read"
	// 为 PermSiteWrite 赋值："site:write"

	PermSiteWrite  = "site:write"
	// 为 PermDeviceRead 赋值："device:read"

	PermDeviceRead = "device:read"
	// 为 PermDeviceCtl 赋值："device:control"

	PermDeviceCtl  = "device:control"
	// 为 PermReconRead 赋值："reconcile:read"

	PermReconRead  = "reconcile:read"
	// 为 PermReconRun 赋值："reconcile:run"

	PermReconRun   = "reconcile:run"
	// 为 PermAlarmRead 赋值："alarm:read"

	PermAlarmRead  = "alarm:read"
	// 为 PermAlarmHandle 赋值："alarm:handle"

	PermAlarmHandle = "alarm:handle"
	// 为 PermAuditRead 赋值："audit:read"

	PermAuditRead  = "audit:read"
// 结束参数/表达式列表

)

// rolePermissions 角色->权限映射（可配置化：生产可迁移至配置/DB）@ticket B端-RBAC
// 声明变量 rolePermissions = map[model.AdminRole]map[string]bool{，用于保存运行过程中的数据

var rolePermissions = map[model.AdminRole]map[string]bool{
	// 执行 model.RoleNormalUser: {，完成对应的业务处理

	model.RoleNormalUser: {
		// 执行 PermSiteRead: true,，完成对应的业务处理

		PermSiteRead: true,
	// 定义结束

	},
	// 执行 model.RoleSiteAdmin: {，完成对应的业务处理

	model.RoleSiteAdmin: {
		// 执行 PermSiteRead: true, PermSiteWrite: true,，完成对应的业务处理

		PermSiteRead: true, PermSiteWrite: true,
		// 执行 PermDeviceRead: true, PermDeviceCtl: true,，完成对应的业务处理

		PermDeviceRead: true, PermDeviceCtl: true,
		// 执行 PermReconRead: true,，完成对应的业务处理

		PermReconRead: true,
		// 执行 PermAlarmRead: true, PermAlarmHandle: true,，完成对应的业务处理

		PermAlarmRead: true, PermAlarmHandle: true,
	// 定义结束

	},
	// 执行 model.RolePlatformOp: {，完成对应的业务处理

	model.RolePlatformOp: {
		// 执行 PermSiteRead: true, PermSiteWrite: true,，完成对应的业务处理

		PermSiteRead: true, PermSiteWrite: true,
		// 执行 PermDeviceRead: true, PermDeviceCtl: true,，完成对应的业务处理

		PermDeviceRead: true, PermDeviceCtl: true,
		// 执行 PermReconRead: true, PermReconRun: true,，完成对应的业务处理

		PermReconRead: true, PermReconRun: true,
		// 执行 PermAlarmRead: true, PermAlarmHandle: true,，完成对应的业务处理

		PermAlarmRead: true, PermAlarmHandle: true,
		// 执行 PermAuditRead: true,，完成对应的业务处理

		PermAuditRead: true,
	// 定义结束

	},
// 当前代码块结束

}

// HasPerm 角色是否拥有权限
// 定义函数 HasPerm：封装一段可复用的处理逻辑（参数：role model.AdminRole, perm string）

func HasPerm(role model.AdminRole, perm string) bool {
	// 返回 rolePermissions[role][perm] 给调用方，结束当前函数

	return rolePermissions[role][perm]
// 当前代码块结束

}

// RequirePerm 权限中间件 @ticket B端-RBAC
// 定义函数 RequirePerm：封装一段可复用的处理逻辑（参数：perm string）

func RequirePerm(perm string) gin.HandlerFunc {
	// 返回 func(c *gin.Context) { 给调用方，结束当前函数

	return func(c *gin.Context) {
		// 定义并初始化局部变量 u，取值为 User(c)

		u := User(c)
		// 条件判断：!HasPerm(u.Role, perm) 时执行下面逻辑

		if !HasPerm(u.Role, perm) {
			// 执行 pkg.AbortWithCode(c, pkg.ErrForb)，完成对应的业务处理

			pkg.AbortWithCode(c, pkg.ErrForb)
			// 返回  给调用方，结束当前函数

			return
		// 当前代码块结束

		}
		// 执行 c.Next()，完成对应的业务处理

		c.Next()
	// 当前代码块结束

	}
// 当前代码块结束

}

// SiteScope 数据域：网点管理员仅限本网点，平台运营可跨网点(0表示全部)@ticket B端-RBAC
// 定义函数 SiteScope：封装一段可复用的处理逻辑（参数：u CurrentUser）

func SiteScope(u CurrentUser) uint {
	// 条件判断：u.Role == model.RoleSiteAdmin 时执行下面逻辑

	if u.Role == model.RoleSiteAdmin {
		// 返回 u.SiteID 给调用方，结束当前函数

		return u.SiteID
	// 当前代码块结束

	}
	// 返回 0 给调用方，结束当前函数

	return 0
// 当前代码块结束

}
