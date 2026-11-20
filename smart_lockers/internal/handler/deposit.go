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
func QueryCabinets(c *gin.Context) {
	siteID, _ := strconv.ParseUint(c.Query("site_id"), 10, 64)
	data, code := service.QueryCabinets(uint(siteID), c.Query("size"))
	respond(c, data, code)
}

// QuoteFee @ticket 04-物品寄存：GET /api/v1/deposit/quote?size=&hours=&coupon=
func QuoteFee(c *gin.Context) {
	hours, _ := strconv.ParseFloat(c.Query("hours"), 64)
	coupon, _ := strconv.ParseInt(c.Query("coupon"), 10, 64)
	data, code := service.QuoteFee(c.Query("size"), hours, coupon)
	respond(c, data, code)
}

// CreateDeposit @ticket 04-物品寄存：POST /api/v1/deposit/orders（鉴权 + 幂等）
func CreateDeposit(c *gin.Context) {
	var req struct {
		DoorID uint    `json:"door_id" binding:"required"`
		SiteID uint    `json:"site_id" binding:"required"`
		Hours  float64 `json:"hours" binding:"required"`
		Coupon int64   `json:"coupon"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	key := idemKey(c)
	data, code := service.CreateDeposit(service.CreateDepositInput{
		UserID: middleware.UID(c), DoorID: req.DoorID, SiteID: req.SiteID,
		Hours: req.Hours, Coupon: req.Coupon, IdempotencyKey: key,
	})
	respond(c, data, code)
}

// ChangeCabinet @ticket 04-物品寄存：POST /api/v1/deposit/orders/:order_no/change-cabinet
func ChangeCabinet(c *gin.Context) {
	var req struct {
		NewDoorID uint `json:"new_door_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	data, code := service.ChangeCabinet(c.Param("order_no"), req.NewDoorID)
	respond(c, data, code)
}
