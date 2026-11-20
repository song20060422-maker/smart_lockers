// @ticket 02-支付宝支付：余额/明细/充值/记录/回调 HTTP handler
// @proto 原型对齐：无直接对应(用户侧支付能力，原型未单列视图)；见 AGENT.MD 2.1
package handler

import (
	"math"
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

// AlipayNotify @ticket 02-支付宝支付：POST /api/v1/pay/alipay/notify（form-urlencoded + RSA2 验签）
// 验签通过后才进入幂等入账；按支付宝规范仅返回纯文本 success / failure
func AlipayNotify(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(200, "failure")
		return
	}
	if err := pkg.VerifyAlipayNotify(c.Request.Form); err != nil {
		pkg.Logf("[WARN] alipay notify verify failed: %v", err)
		c.String(200, "failure")
		return
	}
	chargeNo := c.PostForm("out_trade_no")
	tradeNo := c.PostForm("trade_no")
	status := c.PostForm("trade_status")
	total := c.PostForm("total_amount")
	if chargeNo == "" || tradeNo == "" || (status != "TRADE_SUCCESS" && status != "TRADE_FINISHED") {
		c.String(200, "failure")
		return
	}
	amountF, err := strconv.ParseFloat(total, 64)
	if err != nil || amountF <= 0 {
		c.String(200, "failure")
		return
	}
	amount := int64(math.Round(amountF * 100)) // 元 -> 分 @ticket 02-支付宝支付
	_, code := service.RechargeCallback(service.RechargeCallbackInput{
		ChargeNo: chargeNo, TradeNo: tradeNo, Amount: amount,
		Signature: c.PostForm("sign"), IdempotencyKey: chargeNo,
	})
	if code != pkg.CodeOK {
		c.String(200, "failure")
		return
	}
	c.String(200, "success")
}
