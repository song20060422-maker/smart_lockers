// @ticket 06-个人中心：头像/服务/卡券/退出登录 HTTP handler
// @proto 原型对齐：v-marketing 营销奖促中心(卡券 Coupons)；见 AGENT.MD 2.1
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// OrderSummary @ticket 06-个人中心：GET /api/v1/user/order-summary
func OrderSummary(c *gin.Context) {
	data, code := service.OngoingSummary(middleware.UID(c))
	respond(c, data, code)
}

// UploadAvatar @ticket 06-个人中心：POST /api/v1/user/avatar
func UploadAvatar(c *gin.Context) {
	ct := c.GetHeader("Content-Type")
	// multipart 上传时从文件头取
	if file, err := c.FormFile("avatar"); err == nil {
		ct = file.Header.Get("Content-Type")
		if code := service.ValidateAvatar(ct, file.Size); code != pkg.CodeOK {
			pkg.Fail(c, code)
			return
		}
	}
	var req struct {
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
		Size        int64  `json:"size"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.ContentType != "" {
		ct = req.ContentType
	}
	data, code := service.UpdateAvatar(middleware.UID(c), ct, req.Size, req.URL)
	respond(c, data, code)
}

// MyServices @ticket 06-个人中心：GET /api/v1/user/services
func MyServices(c *gin.Context) {
	data, code := service.MyServices(c.Query("role"))
	respond(c, data, code)
}

// Coupons @ticket 06-个人中心：GET /api/v1/user/coupons
func Coupons(c *gin.Context) {
	data, code := service.Coupons(middleware.UID(c))
	respond(c, data, code)
}

// ValidateCoupon @ticket 06-个人中心：GET /api/v1/user/coupons/:code/validate
func ValidateCoupon(c *gin.Context) {
	discount, code := service.ValidateCoupon(c.Param("code"), middleware.UID(c))
	if code != pkg.CodeOK {
		pkg.Fail(c, code)
		return
	}
	pkg.OK(c, map[string]interface{}{"discount": discount})
}

// Login @ticket 06-个人中心：POST /api/v1/auth/login（签发令牌）
func Login(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	tok, code := service.IssueToken(req.UserID)
	if code != pkg.CodeOK {
		pkg.Fail(c, code)
		return
	}
	pkg.OK(c, map[string]interface{}{"token": tok})
}

// Logout @ticket 06-个人中心：POST /api/v1/auth/logout（令牌即时失效、幂等）
func Logout(c *gin.Context) {
	data, code := service.Logout(middleware.UID(c), middleware.TV(c))
	respond(c, data, code)
}
