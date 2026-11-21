// @ticket 04-物品寄存：柜子/试算/落单/换柜 HTTP handler
// @proto 原型对齐：v-fee 网点收费规则(QuoteFee) / v-orders 存取订单落单；见 AGENT.MD 2.1
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// QueryCabinets @ticket 04-物品寄存：GET /api/v1/deposit/cabinets
// 定义函数 QueryCabinets：封装一段可复用的处理逻辑（参数：c *gin.Context）

func QueryCabinets(c *gin.Context) {
	// 定义并初始化局部变量 siteID, _，取值为 strconv.ParseUint(c.Query("site_id"), 10, 64)

	siteID, _ := strconv.ParseUint(c.Query("site_id"), 10, 64)
	// 定义并初始化局部变量 data, code，取值为 service.QueryCabinets(uint(siteID), c.Query("size"))

	data, code := service.QueryCabinets(uint(siteID), c.Query("size"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// QuoteFee @ticket 04-物品寄存：GET /api/v1/deposit/quote?size=&hours=&coupon=
// 定义函数 QuoteFee：封装一段可复用的处理逻辑（参数：c *gin.Context）

func QuoteFee(c *gin.Context) {
	// 定义并初始化局部变量 hours, _，取值为 strconv.ParseFloat(c.Query("hours"), 64)

	hours, _ := strconv.ParseFloat(c.Query("hours"), 64)
	// 定义并初始化局部变量 coupon, _，取值为 strconv.ParseInt(c.Query("coupon"), 10, 64)

	coupon, _ := strconv.ParseInt(c.Query("coupon"), 10, 64)
	// 定义并初始化局部变量 data, code，取值为 service.QuoteFee(c.Query("size"), hours, coupon)

	data, code := service.QuoteFee(c.Query("size"), hours, coupon)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// CreateDeposit @ticket 04-物品寄存：POST /api/v1/deposit/orders（鉴权 + 幂等）
// 定义函数 CreateDeposit：封装一段可复用的处理逻辑（参数：c *gin.Context）

func CreateDeposit(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 DoorID uint    `json:"door_id" binding:"required"`，完成对应的业务处理

		DoorID uint    `json:"door_id" binding:"required"`
		// 执行 SiteID uint    `json:"site_id" binding:"required"`，完成对应的业务处理

		SiteID uint    `json:"site_id" binding:"required"`
		// 执行 Hours  float64 `json:"hours" binding:"required"`，完成对应的业务处理

		Hours  float64 `json:"hours" binding:"required"`
		// 执行 Coupon int64   `json:"coupon"`，完成对应的业务处理

		Coupon int64   `json:"coupon"`
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
	// 定义并初始化局部变量 key，取值为 idemKey(c)

	key := idemKey(c)
	// 定义并初始化局部变量 data, code，取值为 service.CreateDeposit(service.CreateDepositInput{

	data, code := service.CreateDeposit(service.CreateDepositInput{
		// 执行 UserID: middleware.UID(c), DoorID: req.DoorID, SiteID: req.SiteID,，完成对应的业务处理

		UserID: middleware.UID(c), DoorID: req.DoorID, SiteID: req.SiteID,
		// 执行 Hours: req.Hours, Coupon: req.Coupon, IdempotencyKey: key,，完成对应的业务处理

		Hours: req.Hours, Coupon: req.Coupon, IdempotencyKey: key,
	// 执行 })，完成对应的业务处理

	})
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// ChangeCabinet @ticket 04-物品寄存：POST /api/v1/deposit/orders/:order_no/change-cabinet
// 定义函数 ChangeCabinet：封装一段可复用的处理逻辑（参数：c *gin.Context）

func ChangeCabinet(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 NewDoorID uint `json:"new_door_id" binding:"required"`，完成对应的业务处理

		NewDoorID uint `json:"new_door_id" binding:"required"`
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
	// 定义并初始化局部变量 data, code，取值为 service.ChangeCabinet(c.Param("order_no"), req.NewDoorID)

	data, code := service.ChangeCabinet(c.Param("order_no"), req.NewDoorID)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}
