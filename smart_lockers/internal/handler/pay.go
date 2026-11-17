// @ticket 02-支付宝支付：余额/明细/充值/记录/回调 HTTP handler
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// Balance @ticket 02-支付宝支付：GET /api/v1/pay/balance
func Balance(c *gin.Context) {
	data, code := service.GetBalance(middleware.UID(c))
	respond(c, data, code)
}

// BalanceLogs @ticket 02-支付宝支付：GET /api/v1/pay/balance/logs
func BalanceLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	data, code := service.BalanceLogs(middleware.UID(c), page, pageSize)
	respond(c, data, code)
}

// CreateRecharge @ticket 02-支付宝支付：POST /api/v1/pay/recharge
func CreateRecharge(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	data, code := service.CreateRecharge(middleware.UID(c), req.Amount)
	respond(c, data, code)
}

// RechargeRecords @ticket 02-支付宝支付：GET /api/v1/pay/recharge/records
func RechargeRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	data, code := service.RechargeRecords(middleware.UID(c), page, pageSize)
	respond(c, data, code)
}

// RechargeCallback @ticket 02-支付宝支付：POST /api/v1/pay/alipay/callback（回调幂等入账）
func RechargeCallback(c *gin.Context) {
	var req struct {
		ChargeNo  string `json:"charge_no" binding:"required"`
		TradeNo   string `json:"trade_no" binding:"required"`
		Amount    int64  `json:"amount" binding:"required,gt=0"`
		Signature string `json:"signature"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	data, code := service.RechargeCallback(service.RechargeCallbackInput{
		ChargeNo: req.ChargeNo, TradeNo: req.TradeNo, Amount: req.Amount,
		Signature: req.Signature, IdempotencyKey: idemKey(c),
	})
	respond(c, data, code)
}
