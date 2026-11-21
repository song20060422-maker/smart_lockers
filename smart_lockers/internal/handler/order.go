// @ticket 05-订单查看与取出：列表/详情/临时开柜/补缴/聚合 HTTP handler
// @proto 原型对齐：v-orders 存取订单 / v-opendoor 异常开门(临时开柜授权)；见 AGENT.MD 2.1
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// OrderList @ticket 05-订单查看与取出：GET /api/v1/orders?type=current|history&cursor=&limit=
// 定义函数 OrderList：封装一段可复用的处理逻辑（参数：c *gin.Context）

func OrderList(c *gin.Context) {
	// 定义并初始化局部变量 status，取值为 1

	status := 1
	// 条件判断：c.Query("type") == "history" 时执行下面逻辑

	if c.Query("type") == "history" {
		// 为 status 赋值：2

		status = 2
	// 当前代码块结束

	}
	// 定义并初始化局部变量 limit, _，取值为 strconv.Atoi(c.DefaultQuery("limit", "20"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	// 定义并初始化局部变量 data, code，取值为 service.CursorList(middleware.UID(c), status, c.Query("cursor"), limit)

	data, code := service.CursorList(middleware.UID(c), status, c.Query("cursor"), limit)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// OrderDetail @ticket 05-订单查看与取出：GET /api/v1/orders/:order_no
// 定义函数 OrderDetail：封装一段可复用的处理逻辑（参数：c *gin.Context）

func OrderDetail(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.OrderDetail(middleware.UID(c), c.Param("order_no"))

	data, code := service.OrderDetail(middleware.UID(c), c.Param("order_no"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// OngoingSummary @ticket 05-订单查看与取出：GET /api/v1/orders/summary/ongoing
// 定义函数 OngoingSummary：封装一段可复用的处理逻辑（参数：c *gin.Context）

func OngoingSummary(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.OngoingSummary(middleware.UID(c))

	data, code := service.OngoingSummary(middleware.UID(c))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// CreateOpenAuth @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/open-auth
// 定义函数 CreateOpenAuth：封装一段可复用的处理逻辑（参数：c *gin.Context）

func CreateOpenAuth(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 MaxTimes int `json:"max_times"`，完成对应的业务处理

		MaxTimes int `json:"max_times"`
		// 执行 TTLSec   int `json:"ttl_sec"`，完成对应的业务处理

		TTLSec   int `json:"ttl_sec"`
	// 当前代码块结束

	}
	// 忽略该调用的返回值，仅执行其副作用

	_ = c.ShouldBindJSON(&req)
	// 定义并初始化局部变量 data, code，取值为 service.CreateOpenAuth(middleware.UID(c), c.Param("order_no"), req.MaxTimes, req.TTLSec)

	data, code := service.CreateOpenAuth(middleware.UID(c), c.Param("order_no"), req.MaxTimes, req.TTLSec)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// UseOpenAuth @ticket 05-订单查看与取出：POST /api/v1/orders/open-auth/:token/use
// 定义函数 UseOpenAuth：封装一段可复用的处理逻辑（参数：c *gin.Context）

func UseOpenAuth(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.UseOpenAuth(c.Param("token"))

	data, code := service.UseOpenAuth(c.Param("token"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// MakeupPayment @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/makeup
// 定义函数 MakeupPayment：封装一段可复用的处理逻辑（参数：c *gin.Context）

func MakeupPayment(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 Amount int64 `json:"amount" binding:"required,gt=0"`，完成对应的业务处理

		Amount int64 `json:"amount" binding:"required,gt=0"`
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
	// 定义并初始化局部变量 data, code，取值为 service.MakeupPayment(middleware.UID(c), c.Param("order_no"), req.Amount, idemKey(c))

	data, code := service.MakeupPayment(middleware.UID(c), c.Param("order_no"), req.Amount, idemKey(c))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// CompleteOrder @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/complete
// 定义函数 CompleteOrder：封装一段可复用的处理逻辑（参数：c *gin.Context）

func CompleteOrder(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.CompleteOrder(middleware.UID(c), c.Param("order_no"))

	data, code := service.CompleteOrder(middleware.UID(c), c.Param("order_no"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}
