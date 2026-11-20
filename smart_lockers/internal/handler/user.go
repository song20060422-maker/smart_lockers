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
// 定义函数 OrderSummary：封装一段可复用的处理逻辑（参数：c *gin.Context）

func OrderSummary(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.OngoingSummary(middleware.UID(c))

	data, code := service.OngoingSummary(middleware.UID(c))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// UploadAvatar @ticket 06-个人中心：POST /api/v1/user/avatar
// 定义函数 UploadAvatar：封装一段可复用的处理逻辑（参数：c *gin.Context）

func UploadAvatar(c *gin.Context) {
	// 定义并初始化局部变量 ct，取值为 c.GetHeader("Content-Type")

	ct := c.GetHeader("Content-Type")
	// multipart 上传时从文件头取
	// 条件判断：file, err := c.FormFile("avatar"); err == nil 时执行下面逻辑

	if file, err := c.FormFile("avatar"); err == nil {
		// 为 ct 赋值：file.Header.Get("Content-Type")

		ct = file.Header.Get("Content-Type")
		// 条件判断：code := service.ValidateAvatar(ct, file.Size); code != pkg.CodeOK 时执行下面逻辑

		if code := service.ValidateAvatar(ct, file.Size); code != pkg.CodeOK {
			// 执行 pkg.Fail(c, code)，完成对应的业务处理

			pkg.Fail(c, code)
			// 返回  给调用方，结束当前函数

			return
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 URL         string `json:"url"`，完成对应的业务处理

		URL         string `json:"url"`
		// 执行 ContentType string `json:"content_type"`，完成对应的业务处理

		ContentType string `json:"content_type"`
		// 执行 Size        int64  `json:"size"`，完成对应的业务处理

		Size        int64  `json:"size"`
	// 当前代码块结束

	}
	// 忽略该调用的返回值，仅执行其副作用

	_ = c.ShouldBindJSON(&req)
	// 条件判断：req.ContentType != "" 时执行下面逻辑

	if req.ContentType != "" {
		// 为 ct 赋值：req.ContentType

		ct = req.ContentType
	// 当前代码块结束

	}
	// 定义并初始化局部变量 data, code，取值为 service.UpdateAvatar(middleware.UID(c), ct, req.Size, req.URL)

	data, code := service.UpdateAvatar(middleware.UID(c), ct, req.Size, req.URL)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// MyServices @ticket 06-个人中心：GET /api/v1/user/services
// 定义函数 MyServices：封装一段可复用的处理逻辑（参数：c *gin.Context）

func MyServices(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.MyServices(c.Query("role"))

	data, code := service.MyServices(c.Query("role"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// Coupons @ticket 06-个人中心：GET /api/v1/user/coupons
// 定义函数 Coupons：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Coupons(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.Coupons(middleware.UID(c))

	data, code := service.Coupons(middleware.UID(c))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// ValidateCoupon @ticket 06-个人中心：GET /api/v1/user/coupons/:code/validate
// 定义函数 ValidateCoupon：封装一段可复用的处理逻辑（参数：c *gin.Context）

func ValidateCoupon(c *gin.Context) {
	// 定义并初始化局部变量 discount, code，取值为 service.ValidateCoupon(c.Param("code"), middleware.UID(c))

	discount, code := service.ValidateCoupon(c.Param("code"), middleware.UID(c))
	// 条件判断：code != pkg.CodeOK 时执行下面逻辑

	if code != pkg.CodeOK {
		// 执行 pkg.Fail(c, code)，完成对应的业务处理

		pkg.Fail(c, code)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pkg.OK(c, map[string]interface{}{"discount": discount})，完成对应的业务处理

	pkg.OK(c, map[string]interface{}{"discount": discount})
// 当前代码块结束

}

// Login @ticket 06-个人中心：POST /api/v1/auth/login（签发令牌）
// 定义函数 Login：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Login(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 UserID uint `json:"user_id" binding:"required"`，完成对应的业务处理

		UserID uint `json:"user_id" binding:"required"`
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
	// 定义并初始化局部变量 tok, code，取值为 service.IssueToken(req.UserID)

	tok, code := service.IssueToken(req.UserID)
	// 条件判断：code != pkg.CodeOK 时执行下面逻辑

	if code != pkg.CodeOK {
		// 执行 pkg.Fail(c, code)，完成对应的业务处理

		pkg.Fail(c, code)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pkg.OK(c, map[string]interface{}{"token": tok})，完成对应的业务处理

	pkg.OK(c, map[string]interface{}{"token": tok})
// 当前代码块结束

}

// Logout @ticket 06-个人中心：POST /api/v1/auth/logout（令牌即时失效、幂等）
// 定义函数 Logout：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Logout(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.Logout(middleware.UID(c), middleware.TV(c))

	data, code := service.Logout(middleware.UID(c), middleware.TV(c))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// DemoLogin @ticket 00-project：POST /api/v1/auth/demo-login（手机号+验证码的可验证演示登录，签发真实 JWT）
// 定义函数 DemoLogin：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DemoLogin(c *gin.Context) {
	// 声明变量 req struct，接收手机号与验证码

	var req struct {
		// 手机号，11 位数字

		Phone string `json:"phone" binding:"required"`
		// 验证码，4~6 位数字

		Code string `json:"code" binding:"required"`
	// 当前代码块结束

	}
	// 绑定 JSON，失败返回参数错误

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 调用 service.DemoLogin 校验并签发真实 token

	uid, tok, code := service.DemoLogin(req.Phone, req.Code)
	// 签发失败返回对应错误码

	if code != pkg.CodeOK {
		// 执行 pkg.Fail(c, code)，完成对应的业务处理

		pkg.Fail(c, code)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 返回真实 token 与 user_id，供前端持久化登录态

	pkg.OK(c, map[string]interface{}{"token": tok, "user_id": uid, "phone": req.Phone})
// 当前代码块结束

}
