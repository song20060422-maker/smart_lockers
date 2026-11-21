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
// 定义函数 Balance：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Balance(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.GetBalance(middleware.UID(c))

	data, code := service.GetBalance(middleware.UID(c))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// BalanceLogs @ticket 02-支付宝支付：GET /api/v1/pay/balance/logs
// 定义函数 BalanceLogs：封装一段可复用的处理逻辑（参数：c *gin.Context）

func BalanceLogs(c *gin.Context) {
	// 定义并初始化局部变量 page, _，取值为 strconv.Atoi(c.DefaultQuery("page", "1"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// 定义并初始化局部变量 pageSize, _，取值为 strconv.Atoi(c.DefaultQuery("page_size", "20"))

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	// 定义并初始化局部变量 data, code，取值为 service.BalanceLogs(middleware.UID(c), page, pageSize)

	data, code := service.BalanceLogs(middleware.UID(c), page, pageSize)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// CreateRecharge @ticket 02-支付宝支付：POST /api/v1/pay/recharge
// 定义函数 CreateRecharge：封装一段可复用的处理逻辑（参数：c *gin.Context）

func CreateRecharge(c *gin.Context) {
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
	// 定义并初始化局部变量 data, code，取值为 service.CreateRecharge(middleware.UID(c), req.Amount)

	data, code := service.CreateRecharge(middleware.UID(c), req.Amount)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// RechargeRecords @ticket 02-支付宝支付：GET /api/v1/pay/recharge/records
// 定义函数 RechargeRecords：封装一段可复用的处理逻辑（参数：c *gin.Context）

func RechargeRecords(c *gin.Context) {
	// 定义并初始化局部变量 page, _，取值为 strconv.Atoi(c.DefaultQuery("page", "1"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// 定义并初始化局部变量 pageSize, _，取值为 strconv.Atoi(c.DefaultQuery("page_size", "20"))

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	// 定义并初始化局部变量 data, code，取值为 service.RechargeRecords(middleware.UID(c), page, pageSize)

	data, code := service.RechargeRecords(middleware.UID(c), page, pageSize)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// RechargeCallback @ticket 02-支付宝支付：POST /api/v1/pay/alipay/callback（回调幂等入账）
// 定义函数 RechargeCallback：封装一段可复用的处理逻辑（参数：c *gin.Context）

func RechargeCallback(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 ChargeNo  string `json:"charge_no" binding:"required"`，完成对应的业务处理

		ChargeNo  string `json:"charge_no" binding:"required"`
		// 执行 TradeNo   string `json:"trade_no" binding:"required"`，完成对应的业务处理

		TradeNo   string `json:"trade_no" binding:"required"`
		// 执行 Amount    int64  `json:"amount" binding:"required,gt=0"`，完成对应的业务处理

		Amount    int64  `json:"amount" binding:"required,gt=0"`
		// 执行 Signature string `json:"signature"`，完成对应的业务处理

		Signature string `json:"signature"`
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
	// 定义并初始化局部变量 data, code，取值为 service.RechargeCallback(service.RechargeCallbackInput{

	data, code := service.RechargeCallback(service.RechargeCallbackInput{
		// 执行 ChargeNo: req.ChargeNo, TradeNo: req.TradeNo, Amount: req.Amount,，完成对应的业务处理

		ChargeNo: req.ChargeNo, TradeNo: req.TradeNo, Amount: req.Amount,
		// 执行 Signature: req.Signature, IdempotencyKey: idemKey(c),，完成对应的业务处理

		Signature: req.Signature, IdempotencyKey: idemKey(c),
	// 执行 })，完成对应的业务处理

	})
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// AlipayNotify @ticket 02-支付宝支付：POST /api/v1/pay/alipay/notify（form-urlencoded + RSA2 验签）
// 验签通过后才进入幂等入账；按支付宝规范仅返回纯文本 success / failure
// 定义函数 AlipayNotify：封装一段可复用的处理逻辑（参数：c *gin.Context）

func AlipayNotify(c *gin.Context) {
	// 条件判断：err := c.Request.ParseForm(); err != nil 时执行下面逻辑

	if err := c.Request.ParseForm(); err != nil {
		// 执行 c.String(200, "failure")，完成对应的业务处理

		c.String(200, "failure")
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 条件判断：err := pkg.VerifyAlipayNotify(c.Request.Form); err != nil 时执行下面逻辑

	if err := pkg.VerifyAlipayNotify(c.Request.Form); err != nil {
		// 执行 pkg.Logf("[WARN] alipay notify verify failed: %v", err)，完成对应的业务处理

		pkg.Logf("[WARN] alipay notify verify failed: %v", err)
		// 执行 c.String(200, "failure")，完成对应的业务处理

		c.String(200, "failure")
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 chargeNo，取值为 c.PostForm("out_trade_no")

	chargeNo := c.PostForm("out_trade_no")
	// 定义并初始化局部变量 tradeNo，取值为 c.PostForm("trade_no")

	tradeNo := c.PostForm("trade_no")
	// 定义并初始化局部变量 status，取值为 c.PostForm("trade_status")

	status := c.PostForm("trade_status")
	// 定义并初始化局部变量 total，取值为 c.PostForm("total_amount")

	total := c.PostForm("total_amount")
	// 条件判断：chargeNo == "" || tradeNo == "" || (status != "TRADE_SUCCESS" && status != "TRADE_FINISHED") 时执行下面逻辑

	if chargeNo == "" || tradeNo == "" || (status != "TRADE_SUCCESS" && status != "TRADE_FINISHED") {
		// 执行 c.String(200, "failure")，完成对应的业务处理

		c.String(200, "failure")
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 amountF, err，取值为 strconv.ParseFloat(total, 64)

	amountF, err := strconv.ParseFloat(total, 64)
	// 条件判断：err != nil || amountF <= 0 时执行下面逻辑

	if err != nil || amountF <= 0 {
		// 执行 c.String(200, "failure")，完成对应的业务处理

		c.String(200, "failure")
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 amount，取值为 int64(math.Round(amountF * 100)) // 元 -> 分 @ticket 02-支付宝支付

	amount := int64(math.Round(amountF * 100)) // 元 -> 分 @ticket 02-支付宝支付
	// 定义并初始化局部变量 _, code，取值为 service.RechargeCallback(service.RechargeCallbackInput{

	_, code := service.RechargeCallback(service.RechargeCallbackInput{
		// 执行 ChargeNo: chargeNo, TradeNo: tradeNo, Amount: amount,，完成对应的业务处理

		ChargeNo: chargeNo, TradeNo: tradeNo, Amount: amount,
		// 执行 Signature: c.PostForm("sign"), IdempotencyKey: chargeNo,，完成对应的业务处理

		Signature: c.PostForm("sign"), IdempotencyKey: chargeNo,
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：code != pkg.CodeOK 时执行下面逻辑

	if code != pkg.CodeOK {
		// 执行 c.String(200, "failure")，完成对应的业务处理

		c.String(200, "failure")
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 c.String(200, "success")，完成对应的业务处理

	c.String(200, "success")
// 当前代码块结束

}
