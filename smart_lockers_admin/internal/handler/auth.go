// @ticket B端-RBAC：登录接口，签发带角色的令牌
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
)

type AuthHandler struct{}

// Login 演示账号登录（凭据来自配置占位/环境变量，不写死明文）@ticket B端-RBAC
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	for i, a := range config.Get().DemoAdmins {
		if a.Username == req.Username && a.Password == req.Password {
			role := model.AdminRole(a.Role)
			token, err := middleware.GenerateToken(uint(i+1), role, a.SiteID)
			if err != nil {
				pkg.Fail(c, pkg.ErrInternal)
				return
			}
			pkg.OK(c, map[string]interface{}{
				"token": token, "role": role, "site_id": a.SiteID, "username": a.Username,
			})
			return
		}
	}
	pkg.Fail(c, pkg.ErrAuth)
}

// Me 当前身份
func (h *AuthHandler) Me(c *gin.Context) {
	u := middleware.User(c)
	pkg.OK(c, map[string]interface{}{"uid": u.ID, "role": u.Role, "site_id": u.SiteID})
}
