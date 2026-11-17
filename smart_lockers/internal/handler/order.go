// @ticket 05-订单查看与取出：列表/详情/临时开柜/补缴/聚合 HTTP handler
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// OrderList @ticket 05-订单查看与取出：GET /api/v1/orders?type=current|history&cursor=&limit=
func OrderList(c *gin.Context) {
	status := 1
	if c.Query("type") == "history" {
		status = 2
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	data, code := service.CursorList(middleware.UID(c), status, c.Query("cursor"), limit)
	respond(c, data, code)
}

// OrderDetail @ticket 05-订单查看与取出：GET /api/v1/orders/:order_no
func OrderDetail(c *gin.Context) {
	data, code := service.OrderDetail(middleware.UID(c), c.Param("order_no"))
	respond(c, data, code)
}

// OngoingSummary @ticket 05-订单查看与取出：GET /api/v1/orders/summary/ongoing
func OngoingSummary(c *gin.Context) {
	data, code := service.OngoingSummary(middleware.UID(c))
	respond(c, data, code)
}

// CreateOpenAuth @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/open-auth
func CreateOpenAuth(c *gin.Context) {
	var req struct {
		MaxTimes int `json:"max_times"`
		TTLSec   int `json:"ttl_sec"`
	}
	_ = c.ShouldBindJSON(&req)
	data, code := service.CreateOpenAuth(middleware.UID(c), c.Param("order_no"), req.MaxTimes, req.TTLSec)
	respond(c, data, code)
}

// UseOpenAuth @ticket 05-订单查看与取出：POST /api/v1/orders/open-auth/:token/use
func UseOpenAuth(c *gin.Context) {
	data, code := service.UseOpenAuth(c.Param("token"))
	respond(c, data, code)
}

// MakeupPayment @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/makeup
func MakeupPayment(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	data, code := service.MakeupPayment(middleware.UID(c), c.Param("order_no"), req.Amount, idemKey(c))
	respond(c, data, code)
}

// CompleteOrder @ticket 05-订单查看与取出：POST /api/v1/orders/:order_no/complete
func CompleteOrder(c *gin.Context) {
	data, code := service.CompleteOrder(middleware.UID(c), c.Param("order_no"))
	respond(c, data, code)
}
