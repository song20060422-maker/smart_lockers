// @ticket B端-RBAC：登录接口，签发带角色的令牌
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
)

// 定义结构体 AuthHandler：用于承载一组相关的字段数据

type AuthHandler struct{}

// Login 演示账号登录（凭据来自配置占位/环境变量，不写死明文）@ticket B端-RBAC
// 定义方法 Login：接收者 h *AuthHandler，实现该类型对应的业务行为

func (h *AuthHandler) Login(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 结构体字段 Username，类型为 string，附带标签 `json:"username" binding:"required"`

		Username string `json:"username" binding:"required"`
		// 结构体字段 Password，类型为 string，附带标签 `json:"password" binding:"required"`

		Password string `json:"password" binding:"required"`
	// 当前代码块结束

	}
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 循环处理：for i, a := range config.Get().DemoAdmins，逐轮迭代

	for i, a := range config.Get().DemoAdmins {
		// 条件判断：a.Username == req.Username && a.Password == req.Password 时执行下面逻辑

		if a.Username == req.Username && a.Password == req.Password {
			// 定义并初始化局部变量 role，取值为 model.AdminRole(a.Role)

			role := model.AdminRole(a.Role)
			// 定义并初始化局部变量 token, err，取值为 middleware.GenerateToken(uint(i+1), role, a.SiteID)

			token, err := middleware.GenerateToken(uint(i+1), role, a.SiteID)
			// 条件判断：err != nil 时执行下面逻辑

			if err != nil {
				// 执行 pkg.Fail(c, pkg.ErrInternal)，完成对应的业务处理

				pkg.Fail(c, pkg.ErrInternal)
				// 返回  给调用方，结束当前函数

				return
			// 当前代码块结束

			}
			// 修复偏差2：登录响应增加 user 节点，字段与前端 AdminUser 一致（id/username/role/site_id/status）@ticket B端-RBAC

			pkg.OK(c, map[string]interface{}{
				// 保留顶层 token/role/site_id/username 兼容旧调用，新增 user 对象供前端直接使用

				"token": token, "role": role, "site_id": a.SiteID, "username": a.Username,
				// user 节点：id 为配置索引+1，与 Me 返回的 uid 对齐；status=1 表示启用 @ticket B端-RBAC

				"user": model.AdminUser{ID: uint(i + 1), Username: a.Username, Role: role, SiteID: a.SiteID, Status: 1},
			// 执行 })，完成对应的业务处理

			})
			// 返回  给调用方，结束当前函数

			return
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 执行 pkg.Fail(c, pkg.ErrAuth)，完成对应的业务处理

	pkg.Fail(c, pkg.ErrAuth)
// 当前代码块结束

}

// Me 当前身份
// 定义方法 Me：接收者 h *AuthHandler，实现该类型对应的业务行为

func (h *AuthHandler) Me(c *gin.Context) {
	// 定义并初始化局部变量 u，取值为 middleware.User(c)

	u := middleware.User(c)
	// 修复偏差3：按 uid 从演示账号配置反查 username，避免臆造字段；uid 即配置索引+1 @ticket B端-RBAC

	username := ""
	// 遍历演示账号，定位与当前 uid 对应的用户名

	for i, a := range config.Get().DemoAdmins {
		// uid 与配置索引+1 一致时命中

		if uint(i+1) == u.ID {
			// 取到用户名后结束查找

			username = a.Username
			// 跳出循环

			break
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 修复偏差3：返回与前端 AdminUser 对齐的完整字段——id（与 uid 同值）、username、status，保留 uid 兼容旧调用 @ticket B端-RBAC

	pkg.OK(c, map[string]interface{}{"uid": u.ID, "id": u.ID, "username": username, "role": u.Role, "site_id": u.SiteID, "status": 1})
// 当前代码块结束

}
